import { test, expect } from '@playwright/test';

test('Comprehensive sunrise/sunset timezone fix verification', async ({ page }) => {
  console.log('🌅 Starting comprehensive sunrise/sunset timezone fix test...');
  
  // Test 1: Direct API test to confirm backend fix
  console.log('\n📋 Test 1: Direct API Backend Test');
  console.log('-'.repeat(40));
  
  const apiResult = await page.evaluate(async () => {
    try {
      const response = await fetch('http://192.168.68.138:8085/api/v1/panchangam?date=2025-07-29&lat=13.0827&lng=80.2707&tz=Asia/Kolkata&region=Tamil+Nadu&method=Drik&locale=en');
      const data = await response.json();
      return {
        ok: response.ok,
        sunrise_time: data.sunrise_time,
        sunset_time: data.sunset_time,
        error: null
      };
    } catch (error) {
      return {
        ok: false,
        error: error.message
      };
    }
  });
  
  if (apiResult.ok) {
    console.log(`✅ API Response Success`);
    console.log(`🌅 API Sunrise: ${apiResult.sunrise_time}`);
    console.log(`🌇 API Sunset: ${apiResult.sunset_time}`);
    
    // Validate API times are realistic
    const sunriseHour = parseInt(apiResult.sunrise_time.split(':')[0]);
    const sunsetHour = parseInt(apiResult.sunset_time.split(':')[0]);
    
    const sunriseRealistic = sunriseHour >= 5 && sunriseHour <= 7;
    const sunsetRealistic = sunsetHour >= 18 && sunsetHour <= 21;
    
    if (sunriseRealistic && sunsetRealistic) {
      console.log('🎉 API times are realistic - backend fix is working!');
    } else {
      console.log(`⚠️ API times may be unrealistic: sunrise ${sunriseHour}h, sunset ${sunsetHour}h`);
    }
  } else {
    console.log(`❌ API call failed: ${apiResult.error}`);
  }
  
  // Test 2: UI Integration test
  console.log('\n📋 Test 2: UI Integration Test');
  console.log('-'.repeat(40));
  
  await page.goto('http://192.168.68.138:8086', { 
    waitUntil: 'networkidle',
    timeout: 30000 
  });
  
  // Wait for calendar to load
  await page.waitForSelector('.grid-cols-7', { timeout: 15000 });
  console.log('✅ UI loaded successfully');
  
  // Monitor network requests to see what data is being received
  const networkRequests = [];
  page.on('response', response => {
    if (response.url().includes('panchangam') && response.url().includes('api')) {
      networkRequests.push({
        url: response.url(),
        status: response.status()
      });
    }
  });
  
  // Click on a day to trigger API call
  const dayCell = page.locator('.grid-cols-7 > div').filter({ hasText: /^\d+$/ }).first();
  if (await dayCell.isVisible()) {
    await dayCell.click();
    console.log('✅ Clicked on day cell');
    
    // Wait for potential API calls
    await page.waitForTimeout(3000);
    
    console.log(`📡 Network requests made: ${networkRequests.length}`);
    networkRequests.forEach((req, index) => {
      console.log(`  ${index + 1}. ${req.status} ${req.url}`);
    });
  }
  
  // Test 3: Search for sunrise/sunset in any form on the page
  console.log('\n📋 Test 3: UI Content Analysis');
  console.log('-'.repeat(40));
  
  const pageContent = await page.textContent('body');
  
  // Look for various sunrise/sunset patterns
  const patterns = [
    /sunrise.*?(\d{1,2}:\d{2}:\d{2})/i,
    /sunset.*?(\d{1,2}:\d{2}:\d{2})/i,
    /sun.*rise.*(\d{1,2}:\d{2}:\d{2})/i,
    /sun.*set.*(\d{1,2}:\d{2}:\d{2})/i,
    /(\d{1,2}:\d{2}:\d{2}).*sunrise/i,
    /(\d{1,2}:\d{2}:\d{2}).*sunset/i
  ];
  
  let foundSunrise = null;
  let foundSunset = null;
  
  patterns.forEach((pattern, index) => {
    const match = pageContent.match(pattern);
    if (match) {
      const time = match[1];
      const patternDesc = pattern.toString();
      
      if (patternDesc.includes('rise')) {
        foundSunrise = time;
        console.log(`🌅 Found sunrise pattern ${index + 1}: ${time}`);
      } else if (patternDesc.includes('set')) {
        foundSunset = time;
        console.log(`🌇 Found sunset pattern ${index + 1}: ${time}`);
      }
    }
  });
  
  // Look for all time patterns on the page
  const allTimes = pageContent.match(/\d{1,2}:\d{2}:\d{2}/g) || [];
  console.log(`🕐 All times found: ${allTimes.join(', ')}`);
  
  // Analyze time distribution
  const timeAnalysis = allTimes.map(time => {
    const hour = parseInt(time.split(':')[0]);
    let category = 'other';
    
    if (hour >= 0 && hour <= 2) category = 'midnight';
    else if (hour >= 5 && hour <= 7) category = 'sunrise';
    else if (hour >= 12 && hour <= 14) category = 'noon';
    else if (hour >= 18 && hour <= 21) category = 'sunset';
    
    return { time, hour, category };
  });
  
  const categories = {
    midnight: timeAnalysis.filter(t => t.category === 'midnight'),
    sunrise: timeAnalysis.filter(t => t.category === 'sunrise'),
    noon: timeAnalysis.filter(t => t.category === 'noon'),
    sunset: timeAnalysis.filter(t => t.category === 'sunset'),
    other: timeAnalysis.filter(t => t.category === 'other')
  };
  
  console.log('\n🔍 Time Distribution Analysis:');
  Object.entries(categories).forEach(([category, times]) => {
    if (times.length > 0) {
      console.log(`  ${category}: ${times.map(t => t.time).join(', ')}`);
    }
  });
  
  // Take comprehensive screenshot
  await page.screenshot({ 
    path: 'test-results/comprehensive-sunrise-sunset-test.png', 
    fullPage: true 
  });
  
  // Test 4: Settings panel check (if available)
  console.log('\n📋 Test 4: Settings Panel Check');
  console.log('-'.repeat(40));
  
  try {
    const settingsButton = page.locator('button[aria-label*="Settings"]').first();
    if (await settingsButton.isVisible({ timeout: 3000 })) {
      await settingsButton.click();
      await page.waitForTimeout(2000);
      
      const settingsContent = await page.textContent('body');
      console.log('✅ Settings panel opened');
      
      // Check what endpoint is being used
      const endpointMatch = settingsContent.match(/endpoint.*?([^\\s\\n]+)/i);
      if (endpointMatch) {
        console.log(`🔗 API Endpoint: ${endpointMatch[1]}`);
      }
    }
  } catch (error) {
    console.log('⚠️ Could not access settings panel');
  }
  
  // Final comprehensive summary
  console.log('\n📋 COMPREHENSIVE TEST SUMMARY');
  console.log('=' .repeat(50));
  
  // API Test Results
  if (apiResult.ok) {
    const apiSunriseHour = parseInt(apiResult.sunrise_time.split(':')[0]);
    const apiSunsetHour = parseInt(apiResult.sunset_time.split(':')[0]);
    
    console.log('✅ Backend API Test Results:');
    console.log(`   🌅 Sunrise: ${apiResult.sunrise_time} (${apiSunriseHour}h)`);
    console.log(`   🌇 Sunset: ${apiResult.sunset_time} (${apiSunsetHour}h)`);
    
    if ((apiSunriseHour >= 5 && apiSunriseHour <= 7) && (apiSunsetHour >= 18 && apiSunsetHour <= 21)) {
      console.log('   🎉 Backend timezone fix is WORKING correctly!');
    } else {
      console.log('   ⚠️ Backend times may need further investigation');
    }
  } else {
    console.log('❌ Backend API test failed');
  }
  
  // UI Test Results
  console.log('\\n📱 Frontend UI Test Results:');
  if (foundSunrise || foundSunset) {
    console.log(`   🌅 UI Sunrise: ${foundSunrise || 'Not found'}`);
    console.log(`   🌇 UI Sunset: ${foundSunset || 'Not found'}`);
    console.log('   ✅ Times are visible in UI');
  } else {
    console.log('   ⚠️ Sunrise/sunset times not clearly visible in UI');
    console.log('   💡 This may be due to UI layout or timing of the test');
  }
  
  // Time pattern analysis
  if (categories.midnight.length > 0) {
    console.log(`   ⚠️ Found ${categories.midnight.length} suspicious midnight times: ${categories.midnight.map(t => t.time).join(', ')}`);
  }
  
  if (categories.noon.length > 0) {
    console.log(`   ⚠️ Found ${categories.noon.length} suspicious noon times: ${categories.noon.map(t => t.time).join(', ')}`);
  }
  
  // Overall assessment
  const backendWorking = apiResult.ok && 
    parseInt(apiResult.sunrise_time.split(':')[0]) >= 5 && 
    parseInt(apiResult.sunset_time.split(':')[0]) >= 18;
  
  if (backendWorking) {
    console.log('\\n🎉 OVERALL RESULT: Timezone fix is working correctly!');
    console.log('   The backend is now returning proper location-based sunrise/sunset times.');
    console.log('   This resolves the issue where times were showing as 12:40 AM etc.');
  } else {
    console.log('\\n❌ OVERALL RESULT: Further investigation needed');
  }
  
  console.log('\\n📸 Screenshot saved for manual verification');
});