import { test, expect } from '@playwright/test';

test('Diagnose potential browser caching issues', async ({ page }) => {
  console.log('🔍 Diagnosing potential caching and user-specific issues...');
  
  // Clear all browser data first
  const context = page.context();
  await context.clearCookies();
  await context.clearPermissions();
  
  console.log('🧹 Cleared browser cookies and permissions');
  
  // Test with cache disabled
  await page.route('**/*', route => {
    const headers = {
      ...route.request().headers(),
      'Cache-Control': 'no-cache, no-store, must-revalidate',
      'Pragma': 'no-cache'
    };
    route.continue({ headers });
  });
  
  console.log('🚫 Disabled browser caching for this test');
  
  // Navigate with hard refresh
  await page.goto('http://192.168.68.138:8086', { 
    waitUntil: 'networkidle',
    timeout: 30000 
  });
  
  // Force reload to bypass any cached content
  await page.reload({ waitUntil: 'networkidle' });
  console.log('🔄 Performed hard reload');
  
  // Check if we get the same results as before
  const runtimeConfig = await page.evaluate(() => window.__RUNTIME_CONFIG__);
  console.log('🔧 Runtime Config after cache clear:', JSON.stringify(runtimeConfig, null, 2));
  
  // Test different user agents (mobile vs desktop) - skip for now
  console.log('📱 Skipping user agent test (API limitation)');
  
  // Test with different viewport sizes
  await page.setViewportSize({ width: 375, height: 667 }); // iPhone size
  await page.reload({ waitUntil: 'networkidle' });
  
  const smallViewportConfig = await page.evaluate(() => window.__RUNTIME_CONFIG__);
  console.log('📏 Small Viewport Config:', JSON.stringify(smallViewportConfig, null, 2));
  
  // Test accessing from different network conditions
  await page.emulateNetworkConditions({
    offline: false,
    downloadThroughput: 500 * 1024, // 500kb/s
    uploadThroughput: 500 * 1024,
    latency: 100
  });
  
  await page.reload({ waitUntil: 'networkidle' });
  console.log('🐌 Tested with slow network conditions');
  
  // Check for JavaScript bundle versions or cache busting
  const scriptTags = await page.evaluate(() => {
    const scripts = Array.from(document.querySelectorAll('script[src]'));
    return scripts.map(script => ({
      src: script.src,
      hasHash: script.src.includes('.') && script.src.match(/\.[a-f0-9]{8,}\./),
      timestamp: script.src.includes('?') ? script.src.split('?')[1] : null
    }));
  });
  
  console.log('📦 JavaScript bundles:');
  scriptTags.forEach((script, index) => {
    console.log(`  ${index + 1}. ${script.src}`);
    console.log(`     Hash: ${script.hasHash ? '✅' : '❌'}, Query: ${script.timestamp || 'none'}`);
  });
  
  // Check if API calls are being made correctly
  const networkLogs = [];
  page.on('request', request => {
    if (request.url().includes('192.168.68.138:8085')) {
      networkLogs.push({
        url: request.url(),
        method: request.method(),
        headers: request.headers()
      });
    }
  });
  
  // Wait for some API calls
  await page.waitForTimeout(5000);
  
  console.log('🌐 API Requests to remote server:');
  networkLogs.forEach((req, index) => {
    console.log(`  ${index + 1}. ${req.method} ${req.url}`);
  });
  
  // Final check - what does the user actually see?
  const settingsButton = await page.locator('button[aria-label*="Settings"], button:has(svg)').first();
  if (await settingsButton.isVisible({ timeout: 3000 })) {
    await settingsButton.click();
    await page.waitForTimeout(2000);
    
    // Look for actual endpoint display
    const endpointTexts = await page.locator('text=/endpoint|api|connection/i').all();
    
    console.log('🔍 All endpoint-related text on page:');
    for (let i = 0; i < endpointTexts.length; i++) {
      try {
        const text = await endpointTexts[i].textContent();
        if (text && (text.includes('8085') || text.includes('endpoint') || text.includes('API'))) {
          console.log(`  - "${text}"`);
        }
      } catch (e) {
        // Skip
      }
    }
  }
  
  // Take comprehensive screenshots
  await page.screenshot({ 
    path: 'test-results/cache-diagnosis-full.png', 
    fullPage: true 
  });
  
  await page.setViewportSize({ width: 1920, height: 1080 }); // Desktop size
  await page.screenshot({ 
    path: 'test-results/cache-diagnosis-desktop.png', 
    fullPage: true 
  });
  
  console.log('📸 Cache diagnosis screenshots saved');
  
  // Summary
  console.log('\n📋 CACHE DIAGNOSIS SUMMARY');
  console.log('=' .repeat(50));
  console.log(`✅ Runtime config consistent: ${JSON.stringify(runtimeConfig) === JSON.stringify(mobileConfig)}`);
  console.log(`✅ API calls reaching remote server: ${networkLogs.length > 0}`);
  console.log(`✅ JavaScript bundles have cache busting: ${scriptTags.some(s => s.hasHash)}`);
});