import { test } from '@playwright/test';

test('Test with cache busting', async ({ page }) => {
  console.log('🔄 Testing with cache busting...');
  
  // Clear cache and reload
  await page.goto('about:blank');
  
  // Go to the app with cache busting
  const timestamp = Date.now();
  await page.goto(`http://192.168.68.138:8086?_cb=${timestamp}`, { 
    waitUntil: 'networkidle',
    timeout: 30000 
  });
  
  // Force reload with no cache
  await page.reload({ waitUntil: 'networkidle' });
  
  await page.waitForSelector('.grid-cols-7', { timeout: 15000 });
  console.log('✅ Application loaded with cache busting');
  
  // Find and click location button
  const locationButton = page.locator('button').filter({ hasText: /Chennai|Mumbai|Delhi|New York|Los Angeles|Location/ }).first();
  
  if (await locationButton.isVisible({ timeout: 5000 })) {
    console.log('📍 Location button found, clicking...');
    await locationButton.click();
    await page.waitForTimeout(2000);
    
    // Check the placeholder text first
    const searchInput = page.locator('input[placeholder*="Search"]').first();
    if (await searchInput.isVisible()) {
      const placeholder = await searchInput.getAttribute('placeholder');
      console.log(`🔍 Search placeholder: "${placeholder}"`);
      
      if (placeholder?.includes('worldwide')) {
        console.log('✅ Updated placeholder detected - new code is loaded!');
      } else {
        console.log(`❌ Old placeholder detected - cache issue persists`);
      }
    }
    
    const bodyText = await page.textContent('body');
    console.log('\n🇺🇸 Looking for US cities with cache busting:');
    console.log(`  Contains "US Cities": ${bodyText.includes('US Cities')}`);
    console.log(`  Contains "🇺🇸": ${bodyText.includes('🇺🇸')}`);
    console.log(`  Contains "New York": ${bodyText.includes('New York')}`);
    console.log(`  Contains "Los Angeles": ${bodyText.includes('Los Angeles')}`);
    console.log(`  Contains "America/": ${bodyText.includes('America/')}`);
    
    // Check for the Indian Cities header too
    console.log(`  Contains "🇮🇳": ${bodyText.includes('🇮🇳')}`);
    console.log(`  Contains "Indian Cities": ${bodyText.includes('Indian Cities')}`);
    
    await page.screenshot({ 
      path: 'cache-bust-test.png', 
      fullPage: true 
    });
    console.log('\n📸 Screenshot saved as cache-bust-test.png');
  }
});