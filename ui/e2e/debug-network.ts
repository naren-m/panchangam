import { chromium } from '@playwright/test';

async function debugNetwork() {
  const browser = await chromium.launch({ headless: false });
  const page = await browser.newPage();

  // Log all network requests
  page.on('request', request => {
    if (request.url().includes('api') || request.url().includes('panchangam')) {
      console.log('>> REQUEST:', request.method(), request.url());
    }
  });

  // Log all network responses
  page.on('response', response => {
    if (response.url().includes('api') || response.url().includes('panchangam')) {
      console.log('<< RESPONSE:', response.status(), response.url());
    }
  });

  // Log console messages
  page.on('console', msg => {
    if (msg.type() === 'error' || msg.text().includes('API') || msg.text().includes('error')) {
      console.log('CONSOLE:', msg.type(), msg.text());
    }
  });

  // Log page errors
  page.on('pageerror', error => {
    console.log('PAGE ERROR:', error.message);
  });

  await page.goto('http://localhost:5173');
  console.log('Page loaded, monitoring network...');

  // Wait and monitor
  await page.waitForTimeout(15000);

  console.log('\n--- Final check ---');
  const loadingText = await page.textContent('body');
  if (loadingText?.includes('Loading')) {
    console.log('Still in loading state');
  }
  if (loadingText?.includes('Backend')) {
    console.log('Backend error shown');
  }

  await page.screenshot({ path: 'e2e/screenshots/debug-network.png' });
  await browser.close();
}

debugNetwork();
