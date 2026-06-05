# Panchangam Watch

This is a small watch-only SwiftUI app for testing the Panchangam watch face on Apple Watch.

## Build

```bash
xcodegen generate
xcodebuild build -project PanchangamWatch.xcodeproj -scheme PanchangamWatch -destination 'platform=watchOS Simulator,name=Apple Watch Series 11 (46mm)' -derivedDataPath build/DerivedData
```

## Test On A Watch

1. Open `PanchangamWatch.xcodeproj` in Xcode.
2. Add your Apple account in Xcode Settings.
3. Select the `PanchangamWatch` target and set your team under Signing.
4. Pair the iPhone with Xcode, and make sure the Apple Watch appears under the paired iPhone.
5. Enable Developer Mode on the Apple Watch.
6. Select the real Apple Watch as the run destination.
7. Press `Cmd+R`.

The app tries `http://127.0.0.1:8080/api/v1/tithi/current` in the simulator. If the API is unavailable, it keeps showing built-in demo data.
