import SwiftUI

@main
struct PanchangamApp: App {
    @StateObject private var appState = PanchangamAppState()
    @StateObject private var locationService = LocationService()
    @StateObject private var watchSync = WatchSettingsSync()

    var body: some Scene {
        WindowGroup {
            ContentView(state: appState, locationService: locationService, watchSync: watchSync)
        }
    }
}
