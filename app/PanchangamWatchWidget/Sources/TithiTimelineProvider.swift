import Foundation
import OSLog
import PanchangamShared
import WidgetKit

private let complicationLogger = Logger(subsystem: "app.panchangam.watch", category: "complication")

private final class TimelineCompletionBox: @unchecked Sendable {
    private let completion: (Timeline<TithiEntry>) -> Void

    init(_ completion: @escaping (Timeline<TithiEntry>) -> Void) {
        self.completion = completion
    }

    func call(_ timeline: Timeline<TithiEntry>) {
        completion(timeline)
    }
}

struct TithiEntry: TimelineEntry {
    var date: Date
    var state: TithiWidgetState
}

enum TithiWidgetState {
    case loading
    case valid(TithiSummaryResponse)
    case stale(TithiSummaryResponse)
    case error(String)

    var summary: TithiSummaryResponse? {
        switch self {
        case .loading, .error:
            return nil
        case .valid(let summary), .stale(let summary):
            return summary
        }
    }

    var isStale: Bool {
        if case .stale = self {
            return true
        }
        return false
    }
}

struct TithiTimelineProvider: TimelineProvider {
    private let policy = TithiTimelinePolicy()

    func placeholder(in context: Context) -> TithiEntry {
        TithiEntry(date: Date(), state: .loading)
    }

    func getSnapshot(in context: Context, completion: @escaping (TithiEntry) -> Void) {
        completion(TithiEntry(date: Date(), state: .valid(.preview)))
    }

    func getTimeline(in context: Context, completion: @escaping (Timeline<TithiEntry>) -> Void) {
        let completionBox = TimelineCompletionBox(completion)
        Task {
            let now = Date()
            if let storageStatusText = PanchangamStorage.appGroupStatusText() {
                complicationLogger.error("Complication timeline failed: \(storageStatusText, privacy: .public)")
#if DEBUG && targetEnvironment(simulator)
                let preview = TithiSummaryResponse.watchSimulatorPreview
                completionBox.call(Timeline(
                    entries: entries(for: .stale(preview), now: now),
                    policy: .after(policy.refreshDate(for: preview, now: now))
                ))
                complicationLogger.notice("Complication timeline using simulator preview")
#else
                completionBox.call(Timeline(
                    entries: [TithiEntry(date: now, state: .error(storageStatusText))],
                    policy: .after(policy.refreshDate(for: nil, now: now))
                ))
#endif
                return
            }

            let settingsStore = SettingsStore(defaults: UserDefaults(suiteName: PanchangamStorage.appGroupIdentifier) ?? .standard)
            let cache = TithiCache(fileURL: PanchangamStorage.cacheFileURL())

            do {
                guard settingsStore.hasSavedSettings() else {
                    if let cached = try? cache.load() {
                        let refreshDate = policy.refreshDate(for: cached.summary, now: now)
                        completionBox.call(Timeline(
                            entries: entries(for: .stale(cached.summary), now: now),
                            policy: .after(refreshDate)
                        ))
                        return
                    }

#if DEBUG && targetEnvironment(simulator)
                    let preview = TithiSummaryResponse.watchSimulatorPreview
                    completionBox.call(Timeline(
                        entries: entries(for: .stale(preview), now: now),
                        policy: .after(policy.refreshDate(for: preview, now: now))
                    ))
                    complicationLogger.notice("Complication timeline using simulator settings preview")
#else
                    completionBox.call(Timeline(
                        entries: [TithiEntry(date: now, state: .error("Open iPhone app"))],
                        policy: .after(policy.refreshDate(for: nil, now: now))
                    ))
                    complicationLogger.notice("Complication timeline waiting for iPhone settings")
#endif
                    return
                }

                let settings = try settingsStore.load()
                let response = try await PanchangamAPIClient(baseURL: settings.apiBaseURL).currentTithi(settings: settings, at: now)
                try cache.save(response, fetchedAt: now)
                let refreshDate = policy.refreshDate(for: response, now: now)
                completionBox.call(Timeline(entries: entries(for: .valid(response), now: now), policy: .after(refreshDate)))
                complicationLogger.notice("Complication timeline refreshed")
            } catch {
                if let cached = try? cache.load() {
                    let refreshDate = policy.refreshDate(for: cached.summary, now: now)
                    completionBox.call(Timeline(
                        entries: entries(for: .stale(cached.summary), now: now),
                        policy: .after(refreshDate)
                    ))
                    complicationLogger.error("Complication timeline using cached tithi after refresh failure")
                    return
                }

                completionBox.call(Timeline(
                    entries: [TithiEntry(date: now, state: .error("Unavailable"))],
                    policy: .after(policy.refreshDate(for: nil, now: now))
                ))
                complicationLogger.error("Complication timeline unavailable")
            }
        }
    }

    private func entries(for state: TithiWidgetState, now: Date) -> [TithiEntry] {
        policy.entryDates(for: state.summary, now: now).map { date in
            TithiEntry(date: date, state: state)
        }
    }
}
