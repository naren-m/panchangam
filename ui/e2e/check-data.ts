import { chromium } from '@playwright/test';

async function checkCalendar() {
  const browser = await chromium.launch({ headless: false });
  const page = await browser.newPage();

  await page.goto('http://localhost:5173');
  console.log('Page loaded, waiting for data to load...');

  // Wait for loading to complete - look for progress bar to disappear or data to appear
  try {
    await page.waitForSelector('[role="progressbar"]', { state: 'hidden', timeout: 30000 });
    console.log('Loading indicator hidden');
  } catch {
    console.log('Progress bar check timed out, continuing...');
  }

  // Additional wait for data rendering
  await page.waitForTimeout(5000);

  // Take screenshot
  await page.screenshot({ path: 'e2e/screenshots/calendar-with-data.png', fullPage: true });
  console.log('Screenshot saved to e2e/screenshots/calendar-with-data.png');

  // Check if calendar has actual content (not just skeleton)
  const bodyText = await page.textContent('body');

  // Look for tithi names that indicate data loaded
  const tithiNames = ['Pratipada', 'Dwitiya', 'Tritiya', 'Chaturthi', 'Panchami',
    'Shashthi', 'Saptami', 'Ashtami', 'Navami', 'Dashami', 'Ekadashi',
    'Dwadashi', 'Thrayodashi', 'Chaturdashi', 'Purnima', 'Amavasya'];

  const foundTithi = tithiNames.find(t => bodyText?.includes(t));

  if (foundTithi) {
    console.log(`✅ SUCCESS: Calendar data is showing! Found tithi: ${foundTithi}`);
  } else if (bodyText?.includes('Loading')) {
    console.log('⚠️ Still showing loading state');
  } else if (bodyText?.includes('Backend server')) {
    console.log('❌ Error: Backend server unavailable message showing');
  } else {
    console.log('❓ Unknown state - checking for date numbers...');
    // Look for date numbers 1-31
    const hasNumbers = /\b([1-9]|[12][0-9]|3[01])\b/.test(bodyText || '');
    if (hasNumbers) {
      console.log('✅ Calendar appears to have date numbers');
    }
  }

  await page.waitForTimeout(3000);
  await browser.close();
}

checkCalendar();
