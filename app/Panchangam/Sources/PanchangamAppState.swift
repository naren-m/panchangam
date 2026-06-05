import Combine
import CoreLocation
import Foundation
import OSLog
import PanchangamShared
import WidgetKit

private let iphoneRefreshLogger = Logger(subsystem: "app.panchangam.iphone", category: "refresh")

@MainActor
final class PanchangamAppState: ObservableObject {
    enum Status: Equatable {
        case idle
        case loading
        case loaded(Date)
        case stale(String?)
        case failed(String)

        var text: String {
            switch self {
            case .idle:
                return "Idle"
            case .loading:
                return "Loading"
            case .loaded(let date):
                return "Updated \(Self.statusFormatter.string(from: date))"
            case .stale(let message):
                if let message, !message.isEmpty {
                    return "Showing cached result: \(message)"
                }
                return "Showing cached result"
            case .failed(let message):
                return message
            }
        }

        private static let statusFormatter: DateFormatter = {
            let formatter = DateFormatter()
            formatter.dateStyle = .none
            formatter.timeStyle = .short
            return formatter
        }()
    }

    @Published var summary: TithiSummaryResponse?
    @Published var status: Status = .idle
    @Published var apiBaseURLText: String
    @Published var latitudeText: String
    @Published var longitudeText: String
    @Published var timezoneText: String
    @Published var regionText: String
    @Published var methodText: String
    @Published var localeText: String
    @Published var calendarSystemText: String
    @Published var usesManualLocation = true

    private let settingsStore: SettingsStore
    private let cache: TithiCache

    init(
        settingsStore: SettingsStore = SettingsStore(defaults: UserDefaults(suiteName: PanchangamStorage.appGroupIdentifier) ?? .standard),
        cache: TithiCache = TithiCache(fileURL: PanchangamStorage.cacheFileURL())
    ) {
        self.settingsStore = settingsStore
        self.cache = cache

        let settings = (try? settingsStore.load()) ?? APISettings()
        apiBaseURLText = settings.apiBaseURL.absoluteString
        latitudeText = Self.decimalString(settings.latitude)
        longitudeText = Self.decimalString(settings.longitude)
        timezoneText = settings.timezone
        regionText = settings.region
        methodText = settings.method
        localeText = settings.locale
        calendarSystemText = settings.calendarSystem ?? ""
        _ = showCachedFallback()
    }

    func refresh(sync: WatchSettingsSync? = nil) async {
        guard status != .loading else {
            return
        }

        var attemptedSettings: APISettings?

        do {
            let settings = try makeSettings()
            attemptedSettings = settings
            status = .loading
            try settingsStore.save(settings)
            sync?.send(settings: settings)

            let client = PanchangamAPIClient(baseURL: settings.apiBaseURL)
            let response = try await client.currentTithi(settings: settings)
            try cache.save(response)
            WidgetCenter.shared.reloadAllTimelines()
            sync?.send(settings: settings, latestSummary: response)
            publish(summary: response, status: .loaded(Date()))
            iphoneRefreshLogger.notice("Current tithi refresh completed")
        } catch let error as APISettingsValidationError {
            iphoneRefreshLogger.error("Current tithi refresh failed: \(error.userMessage, privacy: .public)")
            if !showAndSyncCachedFallback(reason: error.userMessage, settings: attemptedSettings, sync: sync) {
                status = .failed(error.userMessage)
            }
        } catch let error as PanchangamAPIError {
            iphoneRefreshLogger.error("Current tithi refresh failed: \(error.userMessage, privacy: .public)")
            if !showAndSyncCachedFallback(reason: error.userMessage, settings: attemptedSettings, sync: sync) {
                status = .failed(error.userMessage)
            }
        } catch {
            iphoneRefreshLogger.error("Current tithi refresh failed with unexpected error")
            if !showAndSyncCachedFallback(reason: "Refresh failed", settings: attemptedSettings, sync: sync) {
                status = .failed("Refresh failed")
            }
        }
    }

    func loadCachedSummary() {
        _ = showCachedFallback()
    }

    func syncCachedToWatch(sync: WatchSettingsSync) {
        guard settingsStore.hasSavedSettings(),
              let settings = try? settingsStore.load() else {
            return
        }

        let cached = try? cache.load()
        sync.send(settings: settings, latestSummary: cached?.summary)
    }

    func apply(location: CLLocation) {
        usesManualLocation = false
        latitudeText = Self.decimalString(location.coordinate.latitude)
        longitudeText = Self.decimalString(location.coordinate.longitude)
        timezoneText = TimeZone.current.identifier
    }

    func makeSettings() throws -> APISettings {
        try APISettings.parse(
            apiBaseURLText: apiBaseURLText,
            latitudeText: latitudeText,
            longitudeText: longitudeText,
            timezoneText: timezoneText,
            regionText: regionText,
            methodText: methodText,
            localeText: localeText,
            calendarSystemText: calendarSystemText
        )
    }

    var developmentHostWarning: String? {
        (try? makeSettings())?.developmentHostWarning
    }

    var storageStatusText: String? {
        PanchangamStorage.appGroupStatusText()
    }

    var hasSavedSettings: Bool {
        settingsStore.hasSavedSettings()
    }

    private static func decimalString(_ value: Double) -> String {
        let formatter = NumberFormatter()
        formatter.locale = Locale(identifier: "en_US_POSIX")
        formatter.numberStyle = .decimal
        formatter.maximumFractionDigits = 8
        formatter.minimumFractionDigits = 0
        formatter.usesGroupingSeparator = false
        return formatter.string(from: NSNumber(value: value)) ?? String(value)
    }

    private func showCachedFallback(reason: String? = nil) -> CachedTithiSummary? {
        guard let cached = try? cache.load() else {
            return nil
        }

        publish(summary: cached.summary, status: .stale(reason))
        return cached
    }

    private func showAndSyncCachedFallback(reason: String, settings: APISettings?, sync: WatchSettingsSync?) -> Bool {
        guard let cached = showCachedFallback(reason: reason) else {
            return false
        }

        if let settings {
            sync?.send(settings: settings, latestSummary: cached.summary)
        }

        return true
    }

    private func publish(summary: TithiSummaryResponse, status: Status) {
        self.status = status
        self.summary = summary
    }
}
