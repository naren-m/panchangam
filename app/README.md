# Panchangam Mobile

This folder contains the iPhone app, Apple Watch app, and Apple Watch WidgetKit complication.

## Requirements

- Full Xcode with iOS and watchOS simulator SDKs.
- XcodeGen.
- A development team configured in `project.yml` before signing a device build.

Install Xcode from the App Store or Apple Developer if `/Applications/Xcode.app` is missing.
Select full Xcode before simulator builds:

```sh
sudo xcode-select -s /Applications/Xcode.app/Contents/Developer
```

Install XcodeGen:

```sh
brew install xcodegen
```

Generate the Xcode project:

```sh
cd app
xcodegen generate --spec project.yml
open PanchangamMobile.xcodeproj
```

## Targets

- `Panchangam`: iPhone app for settings, location, and refresh.
- `PanchangamWatch`: Apple Watch app for current tithi details.
- `PanchangamWatchWidget`: WidgetKit complication for watch faces.
- `PanchangamShared`: shared Swift package for API calls, models, cache, settings, and formatting.

## Bundle Metadata

- iPhone app: `app.panchangam.Panchangam`.
- Watch app: `app.panchangam.Panchangam.watchkitapp`.
- Watch complication extension: `app.panchangam.Panchangam.watchkitapp.widgets`.
- Shared app group: `group.app.panchangam`.

The watch app `Info.plist` sets `WKCompanionAppBundleIdentifier` to the iPhone app bundle ID.
The watch app also sets `WKRunsIndependentlyOfCompanionApp` so the watch app and complication can refresh after settings are synced.

## Preview

The browser preview is only for visual review:

```sh
open preview/index.html
```

Use Xcode previews or a simulator build for the real app render.

## Use

1. Start the Panchangam backend and HTTP gateway, then confirm `/api/v1/tithi/current` responds.
2. Open the `Panchangam` iPhone app.
3. Set API Base URL to the gateway URL. For a physical iPhone or Apple Watch, use the Mac LAN address instead of `localhost` or `127.0.0.1`.
4. Confirm latitude, longitude, timezone, region, method, locale, and calendar system.
5. Tap Refresh. The iPhone app saves settings, fetches the current tithi, caches it, and sends settings plus the latest summary to the watch.
6. Check Watch Sync in the iPhone app. It should show that settings were synced, or a clear paired-watch/install error.
7. Open the watch app once if the complication says `Open iPhone app`. The watch app loads synced settings, cached summary, and then refreshes independently.
8. Add the Panchangam Tithi complication to a watch face. Inline, circular, rectangular, and corner families are supported.
9. Tap the complication to open the watch app detail view for tithi number, paksha, nakshatra, yoga, karana, sunrise, sunset, abhijit, remaining time, and calculation metadata.

## Local Checks

From the repository root, run the app verifier:

```sh
sh app/scripts/verify-mobile.sh
```

Run the shared package checks:

```sh
cd app/PanchangamShared
swift run PanchangamSharedChecks
swift build
```

On a Mac with full Xcode, build the shared schemes:

```sh
sh app/scripts/build-simulators.sh
```

Set `IOS_DESTINATION` or `WATCH_DESTINATION` to override the default generic simulator destinations.
