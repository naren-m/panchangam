import { test, expect } from '@playwright/test';

test.describe('Comprehensive UI Stability Analysis', () => {
  test('Deep UI stability and flakiness detection', async ({ page }) => {
    const networkRequests = [];
    const failedRequests = [];
    const consoleErrors = [];
    const loadingStates = [];
    const dataRenderEvents = [];
    
    // Track all network activity
    page.on('request', request => {
      networkRequests.push({
        url: request.url(),
        method: request.method(),
        timestamp: Date.now()
      });
    });
    
    page.on('requestfailed', request => {
      failedRequests.push({
        url: request.url(),
        failure: request.failure()?.errorText,
        timestamp: Date.now()
      });
    });
    
    // Track console messages
    page.on('console', msg => {
      if (msg.type() === 'error' || msg.type() === 'warn') {
        consoleErrors.push({
          type: msg.type(),
          text: msg.text(),
          timestamp: Date.now()
        });
      }
    });
    
    console.log('🔍 Starting comprehensive UI stability test...');
    
    // Navigate and wait for initial load
    await page.goto('http://localhost:8086');
    await page.waitForSelector('.grid-cols-7', { timeout: 15000 });
    console.log('✅ Calendar grid loaded');
    
    // Test 1: Initial Load Stability
    console.log('📊 Test 1: Analyzing initial load stability...');
    
    // Wait for data to populate and track loading states
    for (let i = 0; i < 10; i++) {
      await page.waitForTimeout(1000);
      
      const loadingCells = await page.locator('.animate-pulse, [data-loading="true"]').count();
      const totalCells = await page.locator('.grid-cols-7 > div').count();
      const errorCells = await page.locator('text=Unavailable, text=Error, text=Failed').count();
      const dataCells = await page.locator('.grid-cols-7 > div').filter({ hasNotText: 'Unavailable' }).count();
      
      loadingStates.push({
        second: i + 1,
        loading: loadingCells,
        total: totalCells,
        errors: errorCells,
        withData: dataCells,
        percentage: Math.round((dataCells / totalCells) * 100)
      });
      
      if (dataCells > totalCells * 0.8) {
        console.log(`✅ Stable data load achieved at ${i + 1} seconds`);
        break;
      }
    }
    
    // Test 2: Navigation Stability
    console.log('📊 Test 2: Testing navigation stability...');
    const beforeNavRequests = networkRequests.length;
    
    // Test rapid navigation
    for (let i = 0; i < 3; i++) {
      await page.click('[data-testid="next-month"], .text-orange-600');
      await page.waitForTimeout(2000);
      
      const newErrorCells = await page.locator('text=Unavailable, text=Error').count();
      if (newErrorCells > 5) {
        console.log(`⚠️ Navigation ${i + 1} caused ${newErrorCells} error cells`);
      }
    }
    
    // Navigate back
    for (let i = 0; i < 3; i++) {
      await page.click('[data-testid="prev-month"], .text-orange-600');
      await page.waitForTimeout(2000);
    }
    
    const afterNavRequests = networkRequests.length;
    const navRequests = afterNavRequests - beforeNavRequests;
    
    // Test 3: Interaction Stability  
    console.log('📊 Test 3: Testing user interaction stability...');
    
    // Try clicking on calendar cells
    const clickableCells = await page.locator('.grid-cols-7 > div[role="button"]').count();
    if (clickableCells > 0) {
      await page.locator('.grid-cols-7 > div[role="button"]').first().click();
      await page.waitForTimeout(1000);
    }
    
    // Test 4: Stress Test - Rapid Interactions
    console.log('📊 Test 4: Stress testing rapid interactions...');
    const stressTestStart = Date.now();
    
    for (let i = 0; i < 5; i++) {
      await page.click('[data-testid="next-month"], .text-orange-600');
      await page.waitForTimeout(500);
      await page.click('[data-testid="prev-month"], .text-orange-600');
      await page.waitForTimeout(500);
    }
    
    await page.waitForTimeout(3000); // Allow settling
    const stressTestEnd = Date.now();
    
    // Final Analysis
    console.log('📊 Analyzing results...');
    
    const finalLoadingCells = await page.locator('.animate-pulse, [data-loading="true"]').count();
    const finalErrorCells = await page.locator('text=Unavailable, text=Error, text=Failed').count();
    const finalDataCells = await page.locator('.grid-cols-7 > div').filter({ hasNotText: 'Unavailable' }).count();
    const totalCells = await page.locator('.grid-cols-7 > div').count();
    
    // Filter API requests
    const apiRequests = networkRequests.filter(r => r.url.includes('panchangam'));
    const apiFailures = failedRequests.filter(r => r.url.includes('panchangam'));
    
    // Calculate request timing patterns
    const requestIntervals = [];
    for (let i = 1; i < apiRequests.length; i++) {
      requestIntervals.push(apiRequests[i].timestamp - apiRequests[i-1].timestamp);
    }
    
    const avgInterval = requestIntervals.length > 0 
      ? requestIntervals.reduce((a, b) => a + b, 0) / requestIntervals.length 
      : 0;
    
    // Print comprehensive analysis
    console.log('\n🔍 COMPREHENSIVE UI STABILITY ANALYSIS RESULTS:');
    console.log('=' .repeat(60));
    
    console.log('\n📊 LOADING PROGRESSION:');
    loadingStates.forEach(state => {
      console.log(`  ${state.second}s: ${state.withData}/${state.total} cells (${state.percentage}%) - ${state.loading} loading, ${state.errors} errors`);
    });
    
    console.log('\n🌐 NETWORK ANALYSIS:');
    console.log(`  Total requests: ${networkRequests.length}`);
    console.log(`  API requests: ${apiRequests.length}`);
    console.log(`  Failed requests: ${failedRequests.length}`);
    console.log(`  API failures: ${apiFailures.length}`);
    console.log(`  Navigation requests: ${navRequests}`);
    console.log(`  Average request interval: ${Math.round(avgInterval)}ms`);
    
    console.log('\n❌ ERROR ANALYSIS:');
    console.log(`  Console errors: ${consoleErrors.length}`);
    console.log(`  Final loading cells: ${finalLoadingCells}`);
    console.log(`  Final error cells: ${finalErrorCells}`);
    console.log(`  Final data cells: ${finalDataCells}/${totalCells} (${Math.round((finalDataCells/totalCells)*100)}%)`);
    
    // Show sample errors
    if (consoleErrors.length > 0) {
      console.log('\n  Sample console errors:');
      consoleErrors.slice(0, 3).forEach(err => {
        console.log(`    - ${err.type}: ${err.text.substring(0, 100)}...`);
      });
    }
    
    if (apiFailures.length > 0) {
      console.log('\n  Sample network failures:');
      apiFailures.slice(0, 3).forEach(fail => {
        console.log(`    - ${fail.url}: ${fail.failure}`);
      });
    }
    
    console.log('\n🎯 STABILITY METRICS:');
    const dataLoadSuccess = (finalDataCells / totalCells) * 100;
    const errorRate = (apiFailures.length / apiRequests.length) * 100;
    const isStable = finalLoadingCells < 3 && finalErrorCells < 5 && dataLoadSuccess > 85;
    
    console.log(`  Data load success: ${dataLoadSuccess.toFixed(1)}%`);
    console.log(`  API error rate: ${errorRate.toFixed(1)}%`);
    console.log(`  UI Stability: ${isStable ? '✅ STABLE' : '❌ UNSTABLE'}`);
    
    // Identify specific flakiness patterns
    console.log('\n🔍 FLAKINESS PATTERNS DETECTED:');
    
    if (finalLoadingCells > 5) {
      console.log(`  ⚠️ PERSISTENT LOADING: ${finalLoadingCells} cells still loading after test`);
    }
    
    if (finalErrorCells > 5) {
      console.log(`  ⚠️ PERSISTENT ERRORS: ${finalErrorCells} cells showing errors`);
    }
    
    if (apiFailures.length > apiRequests.length * 0.1) {
      console.log(`  ⚠️ HIGH API FAILURE RATE: ${apiFailures.length}/${apiRequests.length} requests failing`);
    }
    
    if (avgInterval < 50) {
      console.log(`  ⚠️ REQUEST FLOODING: Average interval ${Math.round(avgInterval)}ms too fast`);
    }
    
    if (consoleErrors.length > 10) {
      console.log(`  ⚠️ EXCESSIVE CONSOLE ERRORS: ${consoleErrors.length} errors detected`);
    }
    
    // Analyze loading progression for instability
    const unstableProgression = loadingStates.some((state, index) => {
      if (index > 0) {
        const prev = loadingStates[index - 1];
        return state.withData < prev.withData; // Data decreased = instability
      }
      return false;
    });
    
    if (unstableProgression) {
      console.log(`  ⚠️ UNSTABLE DATA LOADING: Data cells decreased during loading`);
    }
    
    console.log('\n' + '=' .repeat(60));
    
    await page.screenshot({ 
      path: 'test-results/comprehensive-ui-stability-final.png',
      fullPage: true 
    });
    
    // Final assertions for CI/CD validation
    expect(dataLoadSuccess).toBeGreaterThan(80); // At least 80% data load success
    expect(finalLoadingCells).toBeLessThan(5); // Max 5 cells still loading
    expect(finalErrorCells).toBeLessThan(10); // Max 10 error cells
    expect(apiFailures.length).toBeLessThan(apiRequests.length * 0.2); // Max 20% API failure rate
  });
  
  test('UI responsiveness and timing analysis', async ({ page }) => {
    console.log('⏱️ Testing UI timing and responsiveness...');
    
    const timingEvents = [];
    
    // Track timing events
    const startTime = Date.now();
    
    await page.goto('http://localhost:8086');
    timingEvents.push({ event: 'page_load', time: Date.now() - startTime });
    
    await page.waitForSelector('.grid-cols-7', { timeout: 10000 });
    timingEvents.push({ event: 'calendar_render', time: Date.now() - startTime });
    
    // Wait for first data to appear
    await page.waitForSelector('.grid-cols-7 > div:not(.animate-pulse)', { timeout: 10000 });
    timingEvents.push({ event: 'first_data', time: Date.now() - startTime });
    
    // Test navigation timing
    const navStart = Date.now();
    await page.click('[data-testid="next-month"], .text-orange-600');
    await page.waitForLoadState('networkidle', { timeout: 5000 });
    timingEvents.push({ event: 'navigation_complete', time: Date.now() - navStart });
    
    console.log('\n⏱️ TIMING ANALYSIS:');
    timingEvents.forEach(event => {
      console.log(`  ${event.event}: ${event.time}ms`);
    });
    
    // Timing assertions
    expect(timingEvents.find(e => e.event === 'calendar_render')?.time).toBeLessThan(5000);
    expect(timingEvents.find(e => e.event === 'first_data')?.time).toBeLessThan(8000);
    expect(timingEvents.find(e => e.event === 'navigation_complete')?.time).toBeLessThan(3000);
  });
});