import SwiftUI

@main
struct PanchangamWatchApp: App {
    @StateObject private var state = WatchAppState()
    @StateObject private var receiver = WatchSettingsReceiver()
    @State private var didRunStartupTask = false

    var body: some Scene {
        WindowGroup {
            WatchContentView(state: state, receiver: receiver)
                .task {
                    guard !didRunStartupTask else {
                        return
                    }

                    didRunStartupTask = true
                    receiver.onSettingsChanged = { summary in
                        state.apply(syncedSummary: summary)
                        if summary == nil {
                            Task { await state.refresh() }
                        }
                    }
                    receiver.loadReceivedApplicationContext()
                    state.loadCachedSummaryIfNeeded()
                    await state.refresh()
                }
        }
    }
}
