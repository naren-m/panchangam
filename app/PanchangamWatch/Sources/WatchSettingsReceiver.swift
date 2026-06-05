import Combine
import Foundation
import OSLog
import PanchangamShared
@preconcurrency import WatchConnectivity
import WidgetKit

private let settingsSyncLogger = Logger(subsystem: "app.panchangam.watch", category: "settings-sync")

@MainActor
final class WatchSettingsReceiver: NSObject, ObservableObject, WCSessionDelegate {
    private enum SettingsSyncResult: Sendable {
        case received(SettingsSyncPayload)
        case failed
    }

    @Published var statusText = "Waiting for iPhone settings"

    var onSettingsChanged: ((TithiSummaryResponse?) -> Void)?

    private let settingsStore = SettingsStore(defaults: UserDefaults(suiteName: PanchangamStorage.appGroupIdentifier) ?? .standard)
    private let cache = TithiCache(fileURL: PanchangamStorage.cacheFileURL())
    private let storageStatusText = PanchangamStorage.appGroupStatusText()
    private let session: WCSession?

    override init() {
        if WCSession.isSupported() {
            session = WCSession.default
        } else {
            session = nil
            statusText = "Watch sync unavailable"
            settingsSyncLogger.error("WatchConnectivity is unavailable")
        }
        super.init()
        if let storageStatusText {
            statusText = storageStatusText
            settingsSyncLogger.error("Settings sync storage unavailable: \(storageStatusText, privacy: .public)")
        }
        session?.delegate = self
        session?.activate()
    }

    func loadReceivedApplicationContext() {
        guard let applicationContext = session?.receivedApplicationContext, !applicationContext.isEmpty else {
            return
        }

        settingsSyncLogger.notice("Loading received iPhone settings context")
        apply(syncResult: decode(applicationContext: applicationContext))
    }

    nonisolated func session(_ session: WCSession, activationDidCompleteWith activationState: WCSessionActivationState, error: Error?) {
        let didFail = error != nil
        Task { @MainActor in
            if didFail {
                statusText = "Watch sync failed"
                settingsSyncLogger.error("WatchConnectivity activation failed")
            } else if let storageStatusText {
                statusText = storageStatusText
                settingsSyncLogger.error("Settings sync storage unavailable: \(storageStatusText, privacy: .public)")
            } else if activationState == .activated {
                statusText = "Watch sync ready"
                settingsSyncLogger.notice("WatchConnectivity activated")
                loadReceivedApplicationContext()
            } else {
                statusText = "Watch sync inactive"
                settingsSyncLogger.notice("WatchConnectivity activation ended inactive")
            }
        }
    }

    nonisolated func session(_ session: WCSession, didReceiveApplicationContext applicationContext: [String: Any]) {
        let result = decode(applicationContext: applicationContext)
        Task { @MainActor in
            settingsSyncLogger.notice("Received iPhone settings context")
            apply(syncResult: result)
        }
    }

    private nonisolated func decode(applicationContext: [String: Any]) -> SettingsSyncResult {
        do {
            let payload = try SettingsSyncPayload(dictionary: applicationContext)
            return .received(payload)
        } catch {
            return .failed
        }
    }

    private func apply(syncResult: SettingsSyncResult) {
        if let storageStatusText {
            statusText = storageStatusText
            return
        }

        switch syncResult {
        case .received(let payload):
            do {
                try settingsStore.save(payload.settings)
                if let latestSummary = payload.latestSummary {
                    try cache.save(latestSummary)
                }
                WidgetCenter.shared.reloadAllTimelines()
                statusText = "Settings received"
                settingsSyncLogger.notice("Settings sync applied")
                onSettingsChanged?(payload.latestSummary)
            } catch {
                statusText = "Settings sync failed"
                settingsSyncLogger.error("Settings sync failed while applying context")
            }
        case .failed:
            statusText = "Settings sync failed"
            settingsSyncLogger.error("Settings sync context decode failed")
        }
    }
}
