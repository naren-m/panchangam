// Runtime configuration for the Panchangam application
// API_ENDPOINT is intentionally not set so the app uses window.location.origin
// This allows the same build to work at any domain (nginx proxies /api/ to gateway)
window.__RUNTIME_CONFIG__ = {
  // API_ENDPOINT: Leave empty to use window.location.origin at runtime
  BUILD_TIME: new Date().toISOString(),
  VERSION: "1.0.0"
};