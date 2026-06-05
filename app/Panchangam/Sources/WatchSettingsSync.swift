import Combine
import Foundation
import OSLog
import PanchangamShared
@preconcurrency import WatchConnectivity

private let watchSyncLogger = Logger(subsystem: "app.panchangam.iphone", category: "watch-sync")

@MainActor
final class WatchSettingsSync: NSObject, ObservableObject, WCSessionDelegate {
    @Published var lastSyncText: String?

    private static let syncFormatter: DateFormatter = {
        let formatter = DateFormatter()
        formatter.dateStyle = .none
        formatter.timeStyle = .short
        return formatter
    }()

    private let session: WCSession?
    private var pendingSettings: APISettings?
    private var pendingSummary: TithiSummaryResponse?

    override init() {
        if WCSession.isSupported() {
            session = WCSession.default
        } else {
            session = nil
            lastSyncText = "Watch sync unavailable"
            watchSyncLogger.error("WatchConnectivity is unavailable")
        }
        super.init()
        session?.delegate = self
        session?.activate()
    }

    func send(settings: APISettings, latestSummary: TithiSummaryResponse? = nil) {
        pendingSettings = settings
        pendingSummary = latestSummary

        guard let session, session.activationState == .activated else {
            lastSyncText = "Watch sync inactive"
            watchSyncLogger.notice("Watch settings sync queued while session is inactive")
            return
        }

        sendPendingSettings(using: session)
    }

    nonisolated func session(_ session: WCSession, activationDidCompleteWith activationState: WCSessionActivationState, error: Error?) {
        let didFail = error != nil
        Task { @MainActor in
            if didFail {
                lastSyncText = "Watch sync failed"
                watchSyncLogger.error("WatchConnectivity activation failed")
            } else if activationState == .activated, let session = self.session {
                watchSyncLogger.notice("WatchConnectivity activated")
                sendPendingSettings(using: session)
            } else {
                lastSyncText = "Watch sync inactive"
                watchSyncLogger.notice("WatchConnectivity activation ended inactive")
            }
        }
    }

    nonisolated func sessionDidBecomeInactive(_ session: WCSession) {}

    nonisolated func sessionDidDeactivate(_ session: WCSession) {
        session.activate()
    }

    private func sendPendingSettings(using session: WCSession) {
        guard let settings = pendingSettings else {
            lastSyncText = "Watch sync ready"
            watchSyncLogger.notice("Watch settings sync ready")
            return
        }

        guard session.isPaired else {
            lastSyncText = "No paired Apple Watch"
            watchSyncLogger.error("Watch settings sync failed: no paired Apple Watch")
            return
        }

        guard session.isWatchAppInstalled else {
            lastSyncText = "Install watch app"
            watchSyncLogger.error("Watch settings sync failed: watch app is not installed")
            return
        }

        do {
            let payload = try SettingsSyncPayload(settings: settings, latestSummary: pendingSummary).dictionary()
            try session.updateApplicationContext(payload)
            pendingSettings = nil
            pendingSummary = nil
            lastSyncText = "Watch settings synced at \(Self.syncFormatter.string(from: Date()))"
            watchSyncLogger.notice("Watch settings sync sent")
        } catch {
            lastSyncText = "Watch sync failed"
            watchSyncLogger.error("Watch settings sync failed while sending context")
        }
    }
}
