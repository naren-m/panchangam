import { test, expect } from '@playwright/test';

test('Diagnose data loading issue', async ({ page }) => {
  const networkRequests = [];
  const failedRequests = [];
  const consoleMessages = [];
  
  // Capture all network activity
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
  
  // Capture console messages
  page.on('console', msg => {
    consoleMessages.push({
      type: msg.type(),
      text: msg.text(),
      timestamp: Date.now()
    });
  });
  
  console.log('🔍 Starting data loading diagnosis...');
  
  // Navigate to the application
  try {
    await page.goto('http://localhost:8086', { timeout: 10000 });
    console.log('✅ Page loaded successfully');
  } catch (error) {
    console.log('❌ Page load failed:', error.message);
    return;
  }
  
  // Wait for calendar grid to appear
  try {
    await page.waitForSelector('.grid-cols-7', { timeout: 10000 });
    console.log('✅ Calendar grid found');
  } catch (error) {
    console.log('❌ Calendar grid not found:', error.message);
    await page.screenshot({ path: 'no-calendar-grid.png' });
  }
  
  // Check for loading states
  const loadingElements = await page.locator('[data-loading="true"], .animate-pulse').count();
  console.log(`📊 Loading elements found: ${loadingElements}`);
  
  // Check for error messages
  const errorElements = await page.locator('text=Unavailable, text=Error, text=Failed').count();
  console.log(`❌ Error elements found: ${errorElements}`);
  
  // Check for data cells
  const dataCells = await page.locator('.grid-cols-7 > div').count();
  console.log(`📱 Total calendar cells: ${dataCells}`);
  
  // Wait and monitor data loading for 10 seconds
  console.log('⏳ Monitoring data loading for 10 seconds...');
  
  for (let i = 0; i < 10; i++) {
    await page.waitForTimeout(1000);
    
    const currentLoadingElements = await page.locator('[data-loading="true"], .animate-pulse').count();
    const currentErrorElements = await page.locator('text=Unavailable, text=Error, text=Failed').count();
    const currentDataCells = await page.locator('.grid-cols-7 > div').filter({ hasNotText: 'Unavailable' }).count();
    
    console.log(`  ${i + 1}s: Loading: ${currentLoadingElements}, Errors: ${currentErrorElements}, Data: ${currentDataCells}/${dataCells}`);
    
    if (currentLoadingElements === 0 && currentDataCells > dataCells * 0.8) {
      console.log(`✅ Data loading completed at ${i + 1} seconds`);
      break;
    }
  }
  
  // Get current state
  const finalLoadingElements = await page.locator('[data-loading="true"], .animate-pulse').count();
  const finalErrorElements = await page.locator('text=Unavailable, text=Error, text=Failed').count();
  const finalDataCells = await page.locator('.grid-cols-7 > div').filter({ hasNotText: 'Unavailable' }).count();
  
  // Take a screenshot
  await page.screenshot({ path: 'data-loading-diagnosis.png', fullPage: true });
  
  // Analyze results
  console.log('\n🔍 DIAGNOSIS RESULTS:');
  console.log('=' .repeat(50));
  
  console.log('\n📊 FINAL STATE:');
  console.log(`  Loading elements: ${finalLoadingElements}`);
  console.log(`  Error elements: ${finalErrorElements}`);
  console.log(`  Data cells: ${finalDataCells}/${dataCells} (${Math.round((finalDataCells/dataCells)*100)}%)`);
  
  console.log('\n🌐 NETWORK ACTIVITY:');
  console.log(`  Total requests: ${networkRequests.length}`);
  console.log(`  Failed requests: ${failedRequests.length}`);
  
  // Show API requests specifically
  const apiRequests = networkRequests.filter(r => r.url.includes('panchangam') || r.url.includes('api'));
  console.log(`  API requests: ${apiRequests.length}`);
  
  if (apiRequests.length > 0) {
    console.log('\n  Recent API requests:');
    apiRequests.slice(-5).forEach(req => {
      console.log(`    ${req.method} ${req.url}`);
    });
  }
  
  if (failedRequests.length > 0) {
    console.log('\n  Failed requests:');
    failedRequests.forEach(req => {
      console.log(`    ${req.url}: ${req.failure}`);
    });
  }
  
  console.log('\n💬 CONSOLE MESSAGES:');
  if (consoleMessages.length > 0) {
    consoleMessages.slice(-10).forEach(msg => {
      console.log(`  [${msg.type}] ${msg.text}`);
    });
  } else {
    console.log('  No console messages');
  }
  
  // Check runtime configuration
  const runtimeConfig = await page.evaluate(() => window.__RUNTIME_CONFIG__);
  console.log('\n⚙️ RUNTIME CONFIG:', runtimeConfig);
  
  // Check if data is actually in the DOM
  const cellsWithContent = await page.locator('.grid-cols-7 > div').evaluateAll(elements => {
    return elements.map(el => ({
      textContent: el.textContent?.trim(),
      hasData: el.textContent && el.textContent.trim() !== '' && !el.textContent.includes('Unavailable')
    }));
  });
  
  const cellsWithData = cellsWithContent.filter(cell => cell.hasData);
  console.log(`\n📝 CELL CONTENT ANALYSIS:`);
  console.log(`  Cells with content: ${cellsWithData.length}/${cellsWithContent.length}`);
  
  if (cellsWithData.length > 0) {
    console.log('  Sample cell content:');
    cellsWithData.slice(0, 3).forEach((cell, index) => {
      console.log(`    Cell ${index + 1}: "${cell.textContent?.substring(0, 50)}..."`);
    });
  }
  
  // Determine the issue
  console.log('\n🎯 DIAGNOSIS:');
  if (finalLoadingElements > 5) {
    console.log('  ⚠️ ISSUE: Data is still loading - slow API responses or network issues');
  } else if (finalErrorElements > 5) {
    console.log('  ❌ ISSUE: Multiple error states - API connectivity or configuration problem');
  } else if (finalDataCells < dataCells * 0.5) {
    console.log('  ⚠️ ISSUE: Low data population - partial API failures or missing data');
  } else if (apiRequests.length === 0) {
    console.log('  ❌ ISSUE: No API requests detected - frontend not making API calls');
  } else if (failedRequests.length > 0) {
    console.log('  ❌ ISSUE: Network request failures detected');
  } else {
    console.log('  ✅ STATUS: Data appears to be loading normally');
  }
  
  console.log('\n' + '=' .repeat(50));
});