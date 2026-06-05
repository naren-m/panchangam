# Panchangam Mandala Preview

This folder contains a browser-only preview for quick visual review.

Open `index.html` in a browser to see:

- The Mandala-only Apple Watch face preview.
- The Mandala-style iPhone app screen.
- Current tithi details opened by tapping the watch face.
- Valid, stale, loading, and error timeline states.
- The backend-host warning used for physical iPhone and Apple Watch testing.

This is not a replacement for Xcode. The real app still needs Xcode preview or simulator verification:

```sh
cd app
xcodegen generate
open PanchangamMobile.xcodeproj
```
