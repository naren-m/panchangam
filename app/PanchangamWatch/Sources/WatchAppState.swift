import Combine
import Foundation
import OSLog
import PanchangamShared
import WidgetKit

private let watchRefreshLogger = Logger(subsystem: "app.panchangam.watch", category: "refresh")

@MainActor
final class WatchAppState: ObservableObject {
    enum Status: Equatable {
        case idle
        case loading
        case loaded
        case synced
        case stale(String?)
        case waitingForSettings
        case failed(String)

        var text: String {
            switch self {
            case .idle:
                return "Idle"
            case .loading:
                return "Loading"
            case .loaded:
                return "Loaded"
            case .synced:
                return "Synced from iPhone"
            case .stale(let message):
                if let message, !message.isEmpty {
                    return "Showing cached result: \(message)"
                }
                return "Showing cached result"
            case .waitingForSettings:
                return "Open iPhone app to sync settings"
            case .failed(let message):
                return message
            }
        }
    }

    @Published var summary: TithiSummaryResponse?
    @Published var status: Status = .idle

    private let settingsStore: SettingsStore
    private let cache: TithiCache
    private let storageStatusText: String?

    init(
        settingsStore: SettingsStore = SettingsStore(defaults: UserDefaults(suiteName: PanchangamStorage.appGroupIdentifier) ?? .standard),
        cache: TithiCache = TithiCache(fileURL: PanchangamStorage.cacheFileURL())
    ) {
        self.settingsStore = settingsStore
        self.cache = cache
        storageStatusText = PanchangamStorage.appGroupStatusText()
    }

    func loadCachedSummary() {
        _ = showCachedFallback()
    }

    func loadCachedSummaryIfNeeded() {
        guard summary == nil else {
            return
        }

        _ = showCachedFallback()
    }

    func apply(syncedSummary summary: TithiSummaryResponse?) {
        guard let summary else {
            return
        }
        publish(summary: summary, status: .synced)
    }

    func refresh() async {
        guard status != .loading else {
            return
        }

        if let storageStatusText {
            watchRefreshLogger.error("Watch refresh failed: \(storageStatusText, privacy: .public)")
#if DEBUG && targetEnvironment(simulator)
            publish(summary: .watchSimulatorPreview, status: .stale("Simulator preview: \(storageStatusText)"))
#else
            status = .failed(storageStatusText)
#endif
            return
        }

        status = .loading
        do {
            guard settingsStore.hasSavedSettings() else {
#if DEBUG && targetEnvironment(simulator)
                publish(summary: .watchSimulatorPreview, status: .stale("Simulator preview: Waiting for iPhone settings"))
#else
                if !showCachedFallback() {
                    watchRefreshLogger.notice("Watch refresh waiting for iPhone settings")
                    status = .waitingForSettings
                }
#endif
                return
            }

            let settings = try settingsStore.load()
            try settings.validate()
            let response = try await PanchangamAPIClient(baseURL: settings.apiBaseURL).currentTithi(settings: settings)
            try cache.save(response)
            WidgetCenter.shared.reloadAllTimelines()
            publish(summary: response, status: .loaded)
            watchRefreshLogger.notice("Watch tithi refresh completed")
        } catch let error as APISettingsValidationError {
            watchRefreshLogger.error("Watch refresh failed: \(error.userMessage, privacy: .public)")
            if !showCachedFallback(reason: error.userMessage) {
                status = .failed(error.userMessage)
            }
        } catch let error as PanchangamAPIError {
            watchRefreshLogger.error("Watch refresh failed: \(error.userMessage, privacy: .public)")
            if !showCachedFallback(reason: error.userMessage) {
                status = .failed(error.userMessage)
            }
        } catch {
            watchRefreshLogger.error("Watch refresh failed with unexpected error")
            if !showCachedFallback(reason: "Refresh failed") {
                status = .failed("Refresh failed")
            }
        }
    }

    private func showCachedFallback(reason: String? = nil) -> Bool {
        guard let cached = try? cache.load() else {
            return false
        }

        publish(summary: cached.summary, status: .stale(reason))
        return true
    }

    private func publish(summary: TithiSummaryResponse, status: Status) {
        self.status = status
        self.summary = summary
    }
}
