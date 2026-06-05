#!/bin/sh
set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
APP_DIR="$ROOT_DIR/app"
PROJECT_FILE="$APP_DIR/PanchangamMobile.xcodeproj/project.pbxproj"

fail() {
  printf 'verify-mobile: %s\n' "$1" >&2
  exit 1
}

require_file() {
  [ -f "$1" ] || fail "missing file: $1"
}

require_text() {
  file=$1
  text=$2
  grep -Fq -- "$text" "$file" || fail "missing '$text' in $file"
}

reject_text() {
  file=$1
  text=$2
  if grep -Fq -- "$text" "$file"; then
    fail "unexpected '$text' in $file"
  fi
}

require_plist_raw() {
  file=$1
  key=$2
  expected=$3
  actual=$(plutil -extract "$key" raw -o - "$file" 2>/dev/null || true)
  [ "$actual" = "$expected" ] || fail "$key in $file should be $expected, got ${actual:-missing}"
}

cd "$APP_DIR"
xcodegen generate --spec project.yml

require_file "$PROJECT_FILE"
require_file "$APP_DIR/PanchangamMobile.xcodeproj/xcshareddata/xcschemes/Panchangam.xcscheme"
require_file "$APP_DIR/PanchangamMobile.xcodeproj/xcshareddata/xcschemes/PanchangamWatch.xcscheme"
require_file "$APP_DIR/scripts/build-simulators.sh"
require_file "$APP_DIR/scripts/check-native-tooling.sh"
require_file "$APP_DIR/Panchangam/Assets.xcassets/Contents.json"
require_file "$APP_DIR/Panchangam/Assets.xcassets/AppIcon.appiconset/Contents.json"
require_file "$APP_DIR/PanchangamWatch/Assets.xcassets/Contents.json"
require_file "$APP_DIR/PanchangamWatch/Assets.xcassets/AppIcon.appiconset/Contents.json"
require_file "$APP_DIR/Panchangam/Assets.xcassets/AppIcon.appiconset/iphone-1024.png"
require_file "$APP_DIR/PanchangamWatch/Assets.xcassets/AppIcon.appiconset/watch-1024.png"
require_text "$ROOT_DIR/.gitignore" ".build/"

if ! git -C "$ROOT_DIR" check-ignore -q "app/PanchangamShared/.build/build.db"; then
  fail "SwiftPM build output should be ignored"
fi

plutil -lint \
  "$APP_DIR/Panchangam/Info.plist" \
  "$APP_DIR/PanchangamWatch/Info.plist" \
  "$APP_DIR/PanchangamWatchWidget/Info.plist" \
  "$APP_DIR/Panchangam/Panchangam.entitlements" \
  "$APP_DIR/PanchangamWatch/PanchangamWatch.entitlements" \
  "$APP_DIR/PanchangamWatchWidget/PanchangamWatchWidget.entitlements"

require_text "$PROJECT_FILE" "PRODUCT_BUNDLE_IDENTIFIER = app.panchangam.Panchangam;"
require_text "$PROJECT_FILE" "PRODUCT_BUNDLE_IDENTIFIER = app.panchangam.Panchangam.watchkitapp;"
require_text "$PROJECT_FILE" "PRODUCT_BUNDLE_IDENTIFIER = app.panchangam.Panchangam.watchkitapp.widgets;"
require_text "$APP_DIR/project.yml" "CFBundleDisplayName: Panchangam Tithi"
require_text "$PROJECT_FILE" "APPLICATION_EXTENSION_API_ONLY = YES;"
require_text "$PROJECT_FILE" "ASSETCATALOG_COMPILER_APPICON_NAME = AppIcon;"
require_text "$PROJECT_FILE" "Assets.xcassets"
require_text "$PROJECT_FILE" "Embed Watch Content"
require_text "$PROJECT_FILE" 'dstPath = "$(CONTENTS_FOLDER_PATH)/Watch";'
require_text "$PROJECT_FILE" "PanchangamWatch.app in Embed Watch Content"
require_text "$PROJECT_FILE" "Embed Foundation Extensions"
require_text "$PROJECT_FILE" "PanchangamWatchWidget.appex in Embed Foundation Extensions"
require_text "$PROJECT_FILE" "PBXTargetDependency"
require_text "$APP_DIR/PanchangamWatch/Info.plist" "WKCompanionAppBundleIdentifier"
require_text "$APP_DIR/PanchangamWatch/Info.plist" "app.panchangam.Panchangam"
require_plist_raw "$APP_DIR/PanchangamWatch/Info.plist" "WKRunsIndependentlyOfCompanionApp" "true"
require_plist_raw "$APP_DIR/PanchangamWatchWidget/Info.plist" "CFBundleDisplayName" "Panchangam Tithi"
require_text "$APP_DIR/PanchangamWatch/Info.plist" "CFBundleURLTypes"
require_text "$APP_DIR/PanchangamWatch/Info.plist" "panchangam"
require_text "$APP_DIR/Panchangam/Panchangam.entitlements" "group.app.panchangam"
require_text "$APP_DIR/PanchangamWatch/PanchangamWatch.entitlements" "group.app.panchangam"
require_text "$APP_DIR/PanchangamWatchWidget/PanchangamWatchWidget.entitlements" "group.app.panchangam"
require_text "$APP_DIR/Panchangam/Info.plist" "NSAllowsLocalNetworking"
require_text "$APP_DIR/PanchangamWatch/Info.plist" "NSAllowsLocalNetworking"
require_text "$APP_DIR/PanchangamWatchWidget/Info.plist" "NSAllowsLocalNetworking"
require_text "$APP_DIR/Panchangam/Sources/PanchangamAppState.swift" "try APISettings.parse("
require_text "$APP_DIR/Panchangam/Sources/PanchangamAppState.swift" "guard status != .loading else {"
require_text "$APP_DIR/Panchangam/Sources/PanchangamAppState.swift" "case stale(String?)"
require_text "$APP_DIR/Panchangam/Sources/PanchangamAppState.swift" '"Showing cached result: \(message)"'
require_text "$APP_DIR/Panchangam/Sources/PanchangamAppState.swift" "showCachedFallback()"
require_text "$APP_DIR/Panchangam/Sources/PanchangamAppState.swift" "showCachedFallback(reason:"
require_text "$APP_DIR/Panchangam/Sources/PanchangamAppState.swift" "showAndSyncCachedFallback(reason:"
require_text "$APP_DIR/Panchangam/Sources/PanchangamAppState.swift" "latestSummary: response"
require_text "$APP_DIR/Panchangam/Sources/PanchangamAppState.swift" "var hasSavedSettings: Bool"
require_text "$APP_DIR/Panchangam/Sources/PanchangamAppState.swift" "settingsStore.hasSavedSettings()"
require_text "$APP_DIR/Panchangam/Sources/PanchangamAppState.swift" "syncCachedToWatch"
require_text "$APP_DIR/Panchangam/Sources/PanchangamAppState.swift" "guard settingsStore.hasSavedSettings(),"
require_text "$APP_DIR/Panchangam/Sources/PanchangamAppState.swift" "latestSummary: cached?.summary"
require_text "$APP_DIR/Panchangam/Sources/PanchangamAppState.swift" "latestSummary: cached.summary"
require_text "$APP_DIR/Panchangam/Sources/PanchangamAppState.swift" "publish(summary: response, status: .loaded(Date()))"
require_text "$APP_DIR/Panchangam/Sources/PanchangamAppState.swift" "publish(summary: cached.summary, status: .stale(reason))"
require_text "$APP_DIR/Panchangam/Sources/PanchangamAppState.swift" "private func publish(summary:"
require_text "$APP_DIR/Panchangam/Sources/PanchangamAppState.swift" "import OSLog"
require_text "$APP_DIR/Panchangam/Sources/PanchangamAppState.swift" "Logger(subsystem: \"app.panchangam.iphone\", category: \"refresh\")"
require_text "$APP_DIR/Panchangam/Sources/PanchangamAppState.swift" "iphoneRefreshLogger.notice"
require_text "$APP_DIR/Panchangam/Sources/PanchangamAppState.swift" "iphoneRefreshLogger.error"
require_text "$APP_DIR/Panchangam/Sources/PanchangamAppState.swift" "import WidgetKit"
require_text "$APP_DIR/Panchangam/Sources/PanchangamAppState.swift" "WidgetCenter.shared.reloadAllTimelines()"
reject_text "$APP_DIR/Panchangam/Sources/PanchangamAppState.swift" "summary = response"
reject_text "$APP_DIR/Panchangam/Sources/PanchangamAppState.swift" "summary = cached.summary"
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" "watchSync.lastSyncText"
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" "calculationSection"
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" "applyAndRefresh(location:"
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" "state.syncCachedToWatch(sync: watchSync)"
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" "await state.refresh(sync: watchSync)"
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" "if state.hasSavedSettings {"
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" "@State private var didRunStartupTask = false"
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" "guard !didRunStartupTask else {"
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" "didRunStartupTask = true"
reject_text "$APP_DIR/Panchangam/Sources/ContentView.swift" "if state.summary == nil"
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" "developmentHostWarning"
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" "storageStatusText"
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" 'TaraRow(label: "Storage"'
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" ".keyboardType(.numbersAndPunctuation)"
reject_text "$APP_DIR/Panchangam/Sources/ContentView.swift" ".keyboardType(.decimalPad)"
require_text "$APP_DIR/Panchangam/Sources/PanchangamAppState.swift" "storageStatusText"
require_text "$APP_DIR/Panchangam/Sources/PanchangamAppState.swift" "PanchangamStorage.appGroupStatusText()"
require_text "$APP_DIR/Panchangam/Sources/LocationService.swift" "@preconcurrency import CoreLocation"
require_text "$APP_DIR/Panchangam/Sources/LocationService.swift" "private func requestLocationIfAllowed(status:"
require_text "$APP_DIR/Panchangam/Sources/LocationService.swift" "manager.desiredAccuracy = kCLLocationAccuracyKilometer"
require_text "$APP_DIR/Panchangam/Sources/LocationService.swift" "let latitude = location.coordinate.latitude"
require_text "$APP_DIR/Panchangam/Sources/LocationService.swift" "CLLocationCoordinate2D(latitude: latitude, longitude: longitude)"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/APISettings.swift" "developmentHostWarning"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/APISettings.swift" "Mac LAN address"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/APISettings.swift" "invalidCalendarSystem"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/APISettings.swift" "Calendar system must be Purnimanta or Amanta"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/APISettings.swift" "trimmedCalendarSystem.lowercased()"
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" "taraSpace0"
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" "private let phoneStars"
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" "phoneStarField"
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" ".position(x: proxy.size.width * star.x"
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" "Current Tithi"
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" 'sectionLabel("Mandala")'
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" 'sectionLabel("Calculation")'
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" 'Picker("Calendar System", selection: $state.calendarSystemText)'
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" 'Text("Auto").tag("")'
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" 'Text("Purnimanta").tag("Purnimanta")'
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" 'Text("Amanta").tag("Amanta")'
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" ".pickerStyle(.segmented)"
reject_text "$APP_DIR/Panchangam/Sources/ContentView.swift" 'styledTextField("Calendar System"'
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" "TimelineView(.periodic"
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" "formatter.detailRows(for: summary, now: context.date)"
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" "formatter.nakshatraIndex(for:"
require_text "$APP_DIR/Panchangam/Sources/ContentView.swift" "phoneMandalaTick(index: index, activeIndex:"
reject_text "$APP_DIR/Panchangam/Sources/ContentView.swift" 'Text("10:50")'
require_text "$APP_DIR/Panchangam/Sources/WatchSettingsSync.swift" "pendingSettings"
require_text "$APP_DIR/Panchangam/Sources/WatchSettingsSync.swift" "pendingSummary"
require_text "$APP_DIR/Panchangam/Sources/WatchSettingsSync.swift" "@preconcurrency import WatchConnectivity"
require_text "$APP_DIR/Panchangam/Sources/WatchSettingsSync.swift" "import OSLog"
require_text "$APP_DIR/Panchangam/Sources/WatchSettingsSync.swift" "Logger(subsystem: \"app.panchangam.iphone\", category: \"watch-sync\")"
require_text "$APP_DIR/Panchangam/Sources/WatchSettingsSync.swift" "watchSyncLogger.notice"
require_text "$APP_DIR/Panchangam/Sources/WatchSettingsSync.swift" "watchSyncLogger.error"
require_text "$APP_DIR/Panchangam/Sources/WatchSettingsSync.swift" "private static let syncFormatter"
require_text "$APP_DIR/Panchangam/Sources/WatchSettingsSync.swift" "Watch settings synced at"
require_text "$APP_DIR/Panchangam/Sources/WatchSettingsSync.swift" "Self.syncFormatter.string(from: Date())"
require_text "$APP_DIR/Panchangam/Sources/WatchSettingsSync.swift" "session.isPaired"
require_text "$APP_DIR/Panchangam/Sources/WatchSettingsSync.swift" "session.isWatchAppInstalled"
require_text "$APP_DIR/Panchangam/Sources/WatchSettingsSync.swift" "No paired Apple Watch"
require_text "$APP_DIR/Panchangam/Sources/WatchSettingsSync.swift" "Install watch app"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/PanchangamAPIClient.swift" "currentTithi(settings: APISettings, at: Date? = nil)"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/PanchangamAPIClient.swift" "makeTithiSummaryURL(settings: APISettings, at: Date? = nil)"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/PanchangamAPIClient.swift" "if let at"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/PanchangamAPIClient.swift" "PanchangamAPIError.invalidSummary"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/SettingsSyncPayload.swift" "latestSummary"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/SettingsSyncPayload.swift" "try settings.validate()"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/SettingsSyncPayload.swift" "try latestSummary?.validate()"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/SettingsStore.swift" "hasSavedSettings()"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/SettingsStore.swift" "try settings.validate()"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiCache.swift" "appGroupContainerAvailable"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiCache.swift" "appGroupStatusText"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiCache.swift" "try summary.validate()"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiCache.swift" "try cached.summary.validate()"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiFormatter.swift" "Cached"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiFormatter.swift" '("Tithi Number", String(summary.tithi.number))'
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiFormatter.swift" '("Traditional Name", summary.tithi.traditionalName)'
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiFormatter.swift" '("Paksha", summary.tithi.paksha)'
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiFormatter.swift" '("Paksha Day", String(summary.tithi.pakshaDay))'
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiFormatter.swift" '("Tithi Type", summary.tithi.type)'
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiFormatter.swift" '("Abhijit", abhijitText(for: summary))'
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiFormatter.swift" '("Remaining", remainingText(for: summary, now: now))'
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiFormatter.swift" '("Generated", timeFormatter.string(from: summary.generatedAt))'
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiFormatter.swift" '("Next Refresh", timeFormatter.string(from: summary.nextRefreshAt))'
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiFormatter.swift" '("Region", summary.calculation.region)'
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiFormatter.swift" '("Timezone", summary.calculation.timezone)'
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiFormatter.swift" '("Calendar System", summary.calculation.calendarSystem)'
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiFormatter.swift" '("Method", summary.calculation.method)'
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiFormatter.swift" '("Locale", summary.calculation.locale)'
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiFormatter.swift" "complicationAccessibilityLabel("
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiFormatter.swift" "nakshatraIndex(for summary:"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiFormatter.swift" "Anuradha"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiFormatter.swift" "ended "
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiFormatter.swift" "durationText(minutes:"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiFormatter.swift" "timeStatusText(for:"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiFormatter.swift" "complicationInlineText("
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiFormatter.swift" "complicationSecondaryText("
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiFormatter.swift" "tithiProgress(for summary:"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiFormatter.swift" "progressText(for summary:"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiFormatter.swift" '("Progress", progressText(for: summary, now: now))'
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiTimelinePolicy.swift" "entryCadence"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiTimelinePolicy.swift" "entryDates(for summary:"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiSummary.swift" "TithiSummaryValidationError"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiSummary.swift" "public func validate() throws"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiSummary.swift" "case invalidGeneratedWindow"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiSummary.swift" "Generated time must be within the tithi start and end time"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiSummary.swift" "generatedAt >= tithi.startTime"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiSummary.swift" "watchSimulatorPreview"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiSummary.swift" "#if DEBUG"
require_text "$APP_DIR/PanchangamShared/Checks/PanchangamSharedChecks/main.swift" "generated before tithi start"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" "receiver.statusText"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" "taraWash"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" "private let watchStars"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" "watchStarField"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" ".frame(width: 154, height: 154)"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" ".offset(y: -76)"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" ".font(.system(size: 40, weight: .regular, design: .rounded))"
reject_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" 'Text("Mandala")'
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" "Paksha"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" "Text(formatter.inlineText(for: summary))"
reject_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" "summary.tithi.traditionalName.isEmpty ? summary.tithi.name : summary.tithi.traditionalName"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" "formatter.remainingText(for: summary, now: context.date)"
reject_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" "formatter.remainingText(for: summary, now: summary.generatedAt)"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" "Mandala"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" "mandalaTick"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" "formatter.nakshatraIndex(for:"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" "mandalaTick(index: index, activeIndex:"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" '.accessibilityLabel("Open Mandala tithi details")'
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" ".accessibilityAddTraits(.isButton)"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" ".accessibilityAction {"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" "sheet(item:"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" "WatchTithiDetailView"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" "statusText: selected.statusText"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" "statusText: state.status.text"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" "formatter.detailRows"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" "formatter.detailRows(for: summary, now: context.date)"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" "Current Tithi"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" "let statusText: String"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" "Text(statusText)"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" "opensDetailsWhenSummaryArrives"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" ".onChange(of: state.summary)"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" ".onOpenURL"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" "openTithiDetails(from:"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" 'url.scheme == "panchangam"'
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" 'url.host == "tithi"'
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" 'url.path == "/current"'
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" "state.loadCachedSummary()"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" "Task { await state.refresh() }"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" '.accessibilityLabel("\(summary.day.abhijitMuhurta.name) \(abhijitText(summary))")'
reject_text "$APP_DIR/PanchangamWatch/Sources/WatchContentView.swift" ".accessibilityLabel(receiver.statusText)"
require_text "$APP_DIR/PanchangamWatch/Sources/PanchangamWatchApp.swift" "state.apply(syncedSummary: summary)"
require_text "$APP_DIR/PanchangamWatch/Sources/PanchangamWatchApp.swift" "if summary == nil"
require_text "$APP_DIR/PanchangamWatch/Sources/PanchangamWatchApp.swift" "receiver.loadReceivedApplicationContext()"
require_text "$APP_DIR/PanchangamWatch/Sources/PanchangamWatchApp.swift" "state.loadCachedSummaryIfNeeded()"
require_text "$APP_DIR/PanchangamWatch/Sources/PanchangamWatchApp.swift" "@State private var didRunStartupTask = false"
require_text "$APP_DIR/PanchangamWatch/Sources/PanchangamWatchApp.swift" "guard !didRunStartupTask else {"
require_text "$APP_DIR/PanchangamWatch/Sources/PanchangamWatchApp.swift" "didRunStartupTask = true"
reject_text "$APP_DIR/PanchangamWatch/Sources/PanchangamWatchApp.swift" "state.loadCachedSummary()"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchAppState.swift" "apply(syncedSummary"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchAppState.swift" "func loadCachedSummaryIfNeeded()"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchAppState.swift" "guard summary == nil else"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchAppState.swift" "guard status != .loading else {"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchAppState.swift" "publish(summary: summary, status: .synced)"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchAppState.swift" "publish(summary: response, status: .loaded)"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchAppState.swift" "publish(summary: cached.summary, status: .stale(reason))"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchAppState.swift" "private func publish(summary:"
reject_text "$APP_DIR/PanchangamWatch/Sources/WatchAppState.swift" "summary = response"
reject_text "$APP_DIR/PanchangamWatch/Sources/WatchAppState.swift" "summary = cached.summary"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchAppState.swift" "showCachedFallback()"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchAppState.swift" "showCachedFallback(reason:"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchAppState.swift" "case stale(String?)"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchAppState.swift" '"Showing cached result: \(message)"'
require_text "$APP_DIR/PanchangamWatch/Sources/WatchAppState.swift" "case waitingForSettings"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchAppState.swift" "Open iPhone app to sync settings"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchAppState.swift" "PanchangamStorage.appGroupStatusText()"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchAppState.swift" "WidgetCenter.shared.reloadAllTimelines()"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchAppState.swift" "#if DEBUG && targetEnvironment(simulator)"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchAppState.swift" "Simulator preview:"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchAppState.swift" "Simulator preview: Waiting for iPhone settings"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchAppState.swift" ".watchSimulatorPreview"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchAppState.swift" "import OSLog"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchAppState.swift" "Logger(subsystem: \"app.panchangam.watch\", category: \"refresh\")"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchAppState.swift" "watchRefreshLogger.notice"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchAppState.swift" "watchRefreshLogger.error"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchSettingsReceiver.swift" "((TithiSummaryResponse?) -> Void)?"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchSettingsReceiver.swift" "@preconcurrency import WatchConnectivity"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchSettingsReceiver.swift" "import OSLog"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchSettingsReceiver.swift" "Logger(subsystem: \"app.panchangam.watch\", category: \"settings-sync\")"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchSettingsReceiver.swift" "settingsSyncLogger.notice"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchSettingsReceiver.swift" "settingsSyncLogger.error"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchSettingsReceiver.swift" "loadReceivedApplicationContext"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchSettingsReceiver.swift" "receivedApplicationContext"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchSettingsReceiver.swift" "private enum SettingsSyncResult: Sendable"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchSettingsReceiver.swift" "PanchangamStorage.appGroupStatusText()"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchSettingsReceiver.swift" "private nonisolated func decode(applicationContext:"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchSettingsReceiver.swift" "let didFail = error != nil"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchSettingsReceiver.swift" "payload.latestSummary"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchSettingsReceiver.swift" "onSettingsChanged?(payload.latestSummary)"
require_text "$APP_DIR/PanchangamWatch/Sources/WatchSettingsReceiver.swift" "apply(syncResult:"
reject_text "$APP_DIR/PanchangamWatch/Sources/WatchSettingsReceiver.swift" "apply(applicationContext:"
require_text "$APP_DIR/PanchangamShared/Sources/PanchangamShared/TithiFormatter.swift" "remainingText(for:"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiComplicationViews.swift" "now: entry.date"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiComplicationViews.swift" "isStale: entry.state.isStale"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiComplicationViews.swift" "taraMoon"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiComplicationViews.swift" "formatter.complicationInlineText("
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiComplicationViews.swift" "formatter.complicationSecondaryText("
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiComplicationViews.swift" "mandalaComplicationRing(activeIndex:"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiComplicationViews.swift" "progress: tithiProgress"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiComplicationViews.swift" "formatter.tithiProgress(for:"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiComplicationViews.swift" "activeNakshatraIndex"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiComplicationViews.swift" "formatter.nakshatraIndex(for:"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiComplicationViews.swift" "ForEach(0..<27"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiComplicationViews.swift" ".accessibilityLabel(accessibilityLabel)"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiComplicationViews.swift" ".containerBackground(taraSpace0, for: .widget)"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiComplicationViews.swift" '.widgetURL(URL(string: "panchangam://tithi/current"))'
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiComplicationViews.swift" "formatter.complicationAccessibilityLabel("
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiTimelineProvider.swift" "var isStale"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiTimelineProvider.swift" "PanchangamStorage.appGroupStatusText()"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiTimelineProvider.swift" "settingsStore.hasSavedSettings()"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiTimelineProvider.swift" "Open iPhone app"
reject_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiTimelineProvider.swift" '"Open iPhone"'
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiTimelineProvider.swift" "private func entries(for state:"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiTimelineProvider.swift" "policy.entryDates(for: state.summary, now: now).map"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiTimelineProvider.swift" "import OSLog"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiTimelineProvider.swift" "Logger(subsystem: \"app.panchangam.watch\", category: \"complication\")"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiTimelineProvider.swift" "complicationLogger.notice"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiTimelineProvider.swift" "complicationLogger.error"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiTimelineProvider.swift" "TimelineCompletionBox"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiTimelineProvider.swift" "@unchecked Sendable"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiTimelineProvider.swift" "completionBox.call"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiTimelineProvider.swift" "#if DEBUG && targetEnvironment(simulator)"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiTimelineProvider.swift" "TithiSummaryResponse.watchSimulatorPreview"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiTimelineProvider.swift" "Complication timeline using simulator preview"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/TithiTimelineProvider.swift" "Complication timeline using simulator settings preview"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/PanchangamComplicationWidget.swift" 'configurationDisplayName("Panchangam Tithi")'
require_text "$APP_DIR/PanchangamWatchWidget/Sources/PanchangamComplicationWidget.swift" "watch face complication"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/PanchangamComplicationWidget.swift" ".supportedFamilies(["
require_text "$APP_DIR/PanchangamWatchWidget/Sources/PanchangamComplicationWidget.swift" ".accessoryInline"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/PanchangamComplicationWidget.swift" ".accessoryCircular"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/PanchangamComplicationWidget.swift" ".accessoryRectangular"
require_text "$APP_DIR/PanchangamWatchWidget/Sources/PanchangamComplicationWidget.swift" ".accessoryCorner"
reject_text "$APP_DIR/PanchangamWatchWidget/Sources/PanchangamComplicationWidget.swift" ".containerBackground(.background, for: .widget)"
require_text "$APP_DIR/preview/index.html" "Watch Sync"
require_text "$APP_DIR/preview/index.html" "Backend Host"
require_text "$APP_DIR/preview/index.html" "<h2>Calculation</h2>"
require_text "$APP_DIR/preview/index.html" "Calendar System"
require_text "$APP_DIR/preview/index.html" "calendar-segments"
require_text "$APP_DIR/preview/index.html" ">Auto<"
require_text "$APP_DIR/preview/index.html" ">Purnimanta<"
require_text "$APP_DIR/preview/index.html" ">Amanta<"
require_text "$APP_DIR/preview/index.html" "Locale"
require_text "$APP_DIR/preview/index.html" "receiverStatusText"
require_text "$APP_DIR/preview/index.html" "mandala-face"
require_text "$APP_DIR/preview/index.html" "mandala-ring"
require_text "$APP_DIR/preview/index.html" "watch-detail-sheet"
require_text "$APP_DIR/preview/index.html" 'aria-modal="true"'
reject_text "$APP_DIR/preview/index.html" "watch face preview"
reject_text "$APP_DIR/preview/index.html" "companion app preview"
reject_text "$APP_DIR/preview/index.html" "Authoritative render path"
reject_text "$APP_DIR/preview/index.html" "This preview is a browser aid"
require_text "$APP_DIR/preview/README.md" "Mandala-only"
reject_text "$APP_DIR/preview/README.md" "Inline, circular, rectangular, and corner"
reject_text "$APP_DIR/preview/index.html" "face-controls"
reject_text "$APP_DIR/preview/index.html" "Previous face"
reject_text "$APP_DIR/preview/index.html" "Next face"
reject_text "$APP_DIR/preview/index.html" "Swipe or use dots"
require_text "$APP_DIR/preview/styles.css" "--space-0: #05070f"
require_text "$APP_DIR/preview/styles.css" "tara-case"
require_text "$APP_DIR/preview/styles.css" "starfield"
require_text "$APP_DIR/preview/styles.css" "width: min(100%, 462px);"
require_text "$APP_DIR/preview/styles.css" "aspect-ratio: 462 / 565;"
require_text "$APP_DIR/preview/styles.css" "width: 300px;"
require_text "$APP_DIR/preview/styles.css" "height: 300px;"
require_text "$APP_DIR/preview/styles.css" "translateY(-149px)"
require_text "$APP_DIR/preview/styles.css" "font-size: 70px;"
reject_text "$APP_DIR/preview/styles.css" "face-dots"
require_text "$APP_DIR/preview/preview.js" 'sample.remainingText'
require_text "$APP_DIR/preview/preview.js" 'sample.abhijit'
require_text "$APP_DIR/preview/preview.js" 'sample.traditionalName'
require_text "$APP_DIR/preview/preview.js" 'sample.tithiNumber'
require_text "$APP_DIR/preview/preview.js" 'sample.tithiType'
require_text "$APP_DIR/preview/preview.js" 'sample.progressText'
require_text "$APP_DIR/preview/preview.js" 'sample.nextRefreshAt'
require_text "$APP_DIR/preview/preview.js" 'sample.timezone'
require_text "$APP_DIR/preview/preview.js" "Watch settings synced at 12:00"
require_text "$APP_DIR/preview/preview.js" '["Tithi Number", String(sample.tithiNumber)]'
require_text "$APP_DIR/preview/preview.js" '["Traditional Name", sample.traditionalName]'
require_text "$APP_DIR/preview/preview.js" '["Paksha", sample.paksha]'
require_text "$APP_DIR/preview/preview.js" '["Paksha Day", String(sample.pakshaDay)]'
require_text "$APP_DIR/preview/preview.js" '["Tithi Type", sample.tithiType]'
require_text "$APP_DIR/preview/preview.js" '["Sunrise", sample.sunrise]'
require_text "$APP_DIR/preview/preview.js" '["Sunset", sample.sunset]'
require_text "$APP_DIR/preview/preview.js" '["Abhijit", sample.abhijit]'
require_text "$APP_DIR/preview/preview.js" '["Progress", sample.progressText]'
require_text "$APP_DIR/preview/preview.js" '["Remaining", sample.remainingText]'
require_text "$APP_DIR/preview/preview.js" '["Generated", sample.generatedAt]'
require_text "$APP_DIR/preview/preview.js" '["Next Refresh", sample.nextRefreshAt]'
require_text "$APP_DIR/preview/preview.js" '["Region", sample.region]'
require_text "$APP_DIR/preview/preview.js" '["Timezone", sample.timezone]'
require_text "$APP_DIR/preview/preview.js" '["Calendar System", sample.calendarSystem]'
require_text "$APP_DIR/preview/preview.js" '["Locale", sample.locale]'
require_text "$APP_DIR/preview/preview.js" "Mac LAN address"
require_text "$APP_DIR/preview/preview.js" "renderMandalaTicks"
require_text "$APP_DIR/preview/preview.js" "nakshatraIndex()"
reject_text "$APP_DIR/preview/preview.js" "index === 7"
require_text "$APP_DIR/preview/preview.js" "renderWatchDetails"
require_text "$APP_DIR/preview/preview.js" "openWatchDetails"
reject_text "$APP_DIR/preview/preview.js" "circularText"
reject_text "$APP_DIR/preview/preview.js" "renderRectangularLines"
require_text "$APP_DIR/scripts/build-simulators.sh" "-scheme Panchangam"
require_text "$APP_DIR/scripts/build-simulators.sh" "-scheme PanchangamWatch"
require_text "$APP_DIR/scripts/build-simulators.sh" "check-native-tooling.sh"
require_text "$APP_DIR/scripts/build-simulators.sh" "CODE_SIGNING_ALLOWED=NO"
require_text "$APP_DIR/scripts/build-simulators.sh" "CODE_SIGNING_REQUIRED=NO"
require_text "$APP_DIR/scripts/check-native-tooling.sh" "xcode-select -p"
require_text "$APP_DIR/scripts/check-native-tooling.sh" "/Applications/Xcode.app/Contents/Developer"
require_text "$APP_DIR/scripts/check-native-tooling.sh" "install Xcode"
require_text "$APP_DIR/scripts/check-native-tooling.sh" 'sudo xcode-select -s $xcode_developer_dir'
require_text "$APP_DIR/scripts/check-native-tooling.sh" "xcrun --find simctl"
require_text "$ROOT_DIR/scripts/verify-current-tithi-api.sh" "current-response.json"
require_text "$ROOT_DIR/scripts/verify-current-tithi-api.sh" "validate_current_response"
require_text "$ROOT_DIR/scripts/verify-current-tithi-api.sh" "generated_at should be within"
require_text "$APP_DIR/README.md" "sh app/scripts/build-simulators.sh"
require_text "$APP_DIR/README.md" "Install Xcode"
require_text "$APP_DIR/README.md" "sudo xcode-select -s /Applications/Xcode.app/Contents/Developer"
require_text "$APP_DIR/README.md" "## Use"
require_text "$APP_DIR/README.md" "Set API Base URL"
require_text "$APP_DIR/README.md" "Tap Refresh"
require_text "$APP_DIR/README.md" "Watch Sync"
require_text "$APP_DIR/README.md" "Add the Panchangam Tithi complication"
require_text "$APP_DIR/README.md" "Tap the complication"
reject_text "$APP_DIR/README.md" "sh scripts/build-simulators.sh"

if grep -Fq '"Invalid settings"' "$APP_DIR/Panchangam/Sources/PanchangamAppState.swift"; then
  fail "PanchangamAppState should report parser error messages instead of 'Invalid settings'"
fi

find \
  "$APP_DIR/Panchangam/Sources" \
  "$APP_DIR/PanchangamWatch/Sources" \
  "$APP_DIR/PanchangamWatchWidget/Sources" \
  "$APP_DIR/PanchangamShared/Sources" \
  "$APP_DIR/PanchangamShared/Checks" \
  -name '*.swift' \
  -exec swiftc -parse {} \;

cd "$ROOT_DIR"
node --check app/preview/preview.js

cd "$APP_DIR/PanchangamShared"
swift run --scratch-path /tmp/panchangam-shared-checks-build PanchangamSharedChecks
swift build --scratch-path /tmp/panchangam-shared-build

if native_tooling_status=$("$APP_DIR/scripts/check-native-tooling.sh" 2>&1); then
  printf '%s\n' "$native_tooling_status"
  "$APP_DIR/scripts/build-simulators.sh"
else
  printf '%s\n' "$native_tooling_status"
  printf 'verify-mobile: skipped iOS/watchOS simulator build\n'
fi
