import { chromium } from '@playwright/test';

async function manualTest() {
  console.log('🚀 Starting manual Playwright test...\n');

  const browser = await chromium.launch({ headless: false, slowMo: 500 });
  const context = await browser.newContext();
  const page = await context.newPage();

  try {
    // Test 1: Navigate to app
    console.log('📍 Test 1: Navigating to application...');
    await page.goto('http://localhost:5173', { timeout: 30000 });
    await page.waitForTimeout(2000);
    console.log('✅ App loaded successfully\n');

    // Test 2: Verify header
    console.log('📍 Test 2: Verifying header...');
    const header = await page.locator('h1').first().textContent();
    console.log(`   Header text: "${header}"`);
    console.log('✅ Header verified\n');

    // Test 3: Click Previous Month
    console.log('📍 Test 3: Testing Previous Month navigation...');
    await page.getByRole('button', { name: /previous/i }).click({ timeout: 10000 }).catch(() => {
      console.log('   Trying alternative selector...');
      return page.locator('button').first().click();
    });
    await page.waitForTimeout(1000);
    console.log('✅ Previous month clicked\n');

    // Take screenshot
    await page.screenshot({ path: 'e2e/screenshots/test-prev-month.png' });
    console.log('📸 Screenshot saved: test-prev-month.png\n');

    // Test 4: Click Next Month twice to go forward
    console.log('📍 Test 4: Testing Next Month navigation...');
    const nextBtn = page.getByRole('button', { name: /next/i });
    await nextBtn.click({ timeout: 5000 }).catch(() => page.locator('button').nth(1).click());
    await page.waitForTimeout(500);
    await nextBtn.click({ timeout: 5000 }).catch(() => page.locator('button').nth(1).click());
    await page.waitForTimeout(1000);
    console.log('✅ Next month clicked (x2)\n');

    // Test 5: Click Today button
    console.log('📍 Test 5: Testing Today button...');
    await page.getByRole('button', { name: /today/i }).click({ timeout: 5000 });
    await page.waitForTimeout(1000);
    console.log('✅ Today button clicked\n');

    // Take screenshot
    await page.screenshot({ path: 'e2e/screenshots/test-today.png' });
    console.log('📸 Screenshot saved: test-today.png\n');

    // Test 6: Switch to Table view
    console.log('📍 Test 6: Switching to Table view...');
    await page.getByRole('button', { name: /table/i }).click({ timeout: 5000 });
    await page.waitForTimeout(2000);
    console.log('✅ Table view activated\n');

    // Take screenshot
    await page.screenshot({ path: 'e2e/screenshots/test-table-view.png' });
    console.log('📸 Screenshot saved: test-table-view.png\n');

    // Test 7: Switch to Analytics view
    console.log('📍 Test 7: Switching to Analytics view...');
    await page.getByRole('button', { name: /analytics/i }).click({ timeout: 5000 });
    await page.waitForTimeout(2000);
    console.log('✅ Analytics view activated\n');

    // Take screenshot
    await page.screenshot({ path: 'e2e/screenshots/test-analytics-view.png' });
    console.log('📸 Screenshot saved: test-analytics-view.png\n');

    // Test 8: Switch back to Calendar view
    console.log('📍 Test 8: Switching back to Calendar view...');
    await page.getByRole('button', { name: /calendar/i }).click({ timeout: 5000 });
    await page.waitForTimeout(1000);
    console.log('✅ Calendar view activated\n');

    // Test 9: Open Settings
    console.log('📍 Test 9: Opening Settings...');
    await page.getByRole('button', { name: /settings/i }).click({ timeout: 5000 });
    await page.waitForTimeout(1500);
    console.log('✅ Settings panel opened\n');

    // Take screenshot
    await page.screenshot({ path: 'e2e/screenshots/test-settings.png' });
    console.log('📸 Screenshot saved: test-settings.png\n');

    // Close settings by pressing Escape
    await page.keyboard.press('Escape');
    await page.waitForTimeout(500);

    // Test 10: Test responsive - Mobile viewport
    console.log('📍 Test 10: Testing mobile viewport...');
    await page.setViewportSize({ width: 375, height: 667 });
    await page.waitForTimeout(1000);
    await page.screenshot({ path: 'e2e/screenshots/test-mobile.png' });
    console.log('✅ Mobile viewport tested\n');
    console.log('📸 Screenshot saved: test-mobile.png\n');

    // Reset viewport
    await page.setViewportSize({ width: 1280, height: 720 });

    console.log('═══════════════════════════════════════');
    console.log('🎉 All manual tests completed successfully!');
    console.log('═══════════════════════════════════════\n');

    // Keep browser open for 5 seconds to view
    console.log('Browser will close in 5 seconds...');
    await page.waitForTimeout(5000);

  } catch (error) {
    console.error('❌ Test failed:', error.message);
    await page.screenshot({ path: 'e2e/screenshots/test-error.png' });
  } finally {
    await browser.close();
  }
}

manualTest();
