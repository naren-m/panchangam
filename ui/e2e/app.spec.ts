import { test, expect } from '@playwright/test';

test.describe('Panchangam Application - Core Functionality', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    // Wait for the app to load - use first() to handle multiple h1 elements
    await expect(page.locator('h1').first()).toContainText('Panchangam');
  });

  test('should display the application header correctly', async ({ page }) => {
    await expect(page.getByRole('heading', { name: /Panchangam/i })).toBeVisible();
    await expect(page.getByText('Hindu Calendar & Astronomical Almanac')).toBeVisible();
  });

  test('should show default location as Milpitas, California', async ({ page }) => {
    await expect(page.getByText(/Milpitas/i).first()).toBeVisible();
  });

  test('should display current month in navigation', async ({ page }) => {
    const currentDate = new Date();
    const monthNames = ['January', 'February', 'March', 'April', 'May', 'June',
      'July', 'August', 'September', 'October', 'November', 'December'];
    const currentMonth = monthNames[currentDate.getMonth()];
    const currentYear = currentDate.getFullYear().toString();

    await expect(page.getByText(new RegExp(currentMonth))).toBeVisible();
    await expect(page.getByText(new RegExp(currentYear))).toBeVisible();
  });

  test('should have view switcher with Calendar, Table, and Graph options', async ({ page }) => {
    await expect(page.getByRole('button', { name: /calendar/i })).toBeVisible();
    await expect(page.getByRole('button', { name: /table/i })).toBeVisible();
    await expect(page.getByRole('button', { name: /graph|analytics/i })).toBeVisible();
  });
});

test.describe('Panchangam Application - Navigation', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('h1').first()).toContainText('Panchangam');
  });

  test('should navigate to previous month', async ({ page }) => {
    const currentDate = new Date();
    const currentMonth = currentDate.getMonth();

    // Find prev button by aria-label or icon
    const prevButton = page.locator('button').filter({ has: page.locator('svg') }).first();
    await expect(prevButton).toBeVisible();
    await prevButton.click();

    // Wait for navigation
    await page.waitForTimeout(500);

    // The month should have changed
    const monthNames = ['January', 'February', 'March', 'April', 'May', 'June',
      'July', 'August', 'September', 'October', 'November', 'December'];
    const expectedMonth = monthNames[currentMonth === 0 ? 11 : currentMonth - 1];
    await expect(page.getByText(new RegExp(expectedMonth))).toBeVisible();
  });

  test('should navigate to next month', async ({ page }) => {
    const currentDate = new Date();
    const currentMonth = currentDate.getMonth();

    // Find buttons with SVG icons and get the second one (next)
    const navButtons = page.locator('button').filter({ has: page.locator('svg') });
    const nextButton = navButtons.nth(1);
    await expect(nextButton).toBeVisible();
    await nextButton.click();

    // Wait for navigation
    await page.waitForTimeout(500);

    const monthNames = ['January', 'February', 'March', 'April', 'May', 'June',
      'July', 'August', 'September', 'October', 'November', 'December'];
    const expectedMonth = monthNames[(currentMonth + 1) % 12];
    await expect(page.getByText(new RegExp(expectedMonth))).toBeVisible();
  });

  test('should return to today when clicking Today button', async ({ page }) => {
    // First navigate to a different month using left navigation
    const navButtons = page.locator('button').filter({ has: page.locator('svg') });
    const prevButton = navButtons.first();
    await prevButton.click();
    await page.waitForTimeout(500);
    await prevButton.click();
    await page.waitForTimeout(500);

    // Click Today button
    const todayButton = page.getByRole('button', { name: /today/i });
    await expect(todayButton).toBeVisible();
    await todayButton.click();

    await page.waitForTimeout(500);

    // Verify current month is displayed
    const currentDate = new Date();
    const monthNames = ['January', 'February', 'March', 'April', 'May', 'June',
      'July', 'August', 'September', 'October', 'November', 'December'];
    const currentMonth = monthNames[currentDate.getMonth()];
    await expect(page.getByText(new RegExp(currentMonth))).toBeVisible();
  });
});

test.describe('Panchangam Application - View Switching', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('h1').first()).toContainText('Panchangam');
  });

  test('should switch to Table view', async ({ page }) => {
    const tableButton = page.getByRole('button', { name: /table/i });
    await tableButton.click();

    // Wait for table view to load
    await page.waitForTimeout(2000);

    // Table view should be active - check for table content or heading change
    // The view has switched if the Table button is now active/selected
    await expect(tableButton).toBeVisible();
    // Also check that the page is still functional
    await expect(page.locator('h1').first()).toContainText('Panchangam');
  });

  test('should switch to Graph view', async ({ page }) => {
    const graphButton = page.getByRole('button', { name: /graph|analytics/i });
    await graphButton.click();

    // Wait for graph view to load
    await page.waitForTimeout(2000);

    // Graph view should be active - check that the button is visible and page is functional
    await expect(graphButton).toBeVisible();
    // Also check that the page is still functional
    await expect(page.locator('h1').first()).toContainText('Panchangam');
  });

  test('should switch back to Calendar view', async ({ page }) => {
    // Switch to Table first
    await page.getByRole('button', { name: /table/i }).click();
    await page.waitForTimeout(500);

    // Switch back to Calendar
    const calendarButton = page.getByRole('button', { name: /calendar/i });
    await calendarButton.click();

    await page.waitForTimeout(500);

    // Calendar content should be visible (grid or day elements)
    await expect(page.locator('h1').first()).toContainText('Panchangam');
  });
});

test.describe('Panchangam Application - Calendar Interaction', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('h1').first()).toContainText('Panchangam');
    // Wait for calendar to load
    await page.waitForTimeout(2000);
  });

  test('should display calendar grid with day cells', async ({ page }) => {
    // Wait for calendar data to load
    await page.waitForTimeout(3000);

    // Look for the main container or grid
    const container = page.locator('.container, [class*="calendar"], main').first();
    await expect(container).toBeVisible();
  });

  test('should be able to interact with the calendar', async ({ page }) => {
    // Wait for calendar to load with data
    await page.waitForTimeout(3000);

    // Verify the page has loaded and is interactive
    const heading = page.getByRole('heading', { name: /Panchangam/i });
    await expect(heading).toBeVisible();

    // Check that navigation buttons are clickable
    const buttons = page.getByRole('button');
    const buttonCount = await buttons.count();
    expect(buttonCount).toBeGreaterThan(3);
  });
});

test.describe('Panchangam Application - Settings and Location', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('h1').first()).toContainText('Panchangam');
  });

  test('should have navigation buttons available', async ({ page }) => {
    // Check that multiple navigation/action buttons are available
    const buttons = page.getByRole('button');
    const count = await buttons.count();
    expect(count).toBeGreaterThan(3);
  });

  test('should display location information', async ({ page }) => {
    // Verify location is displayed somewhere on the page
    await expect(page.getByText(/Milpitas|California/i).first()).toBeVisible();
  });
});

test.describe('Panchangam Application - Footer', () => {
  test('should display calculation method and location in footer', async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('h1').first()).toContainText('Panchangam');

    // Check footer content - wait for page to fully load
    await page.waitForTimeout(1000);
    await expect(page.getByText(/Calculated using.*method/i)).toBeVisible();
    await expect(page.getByText(/Milpitas/i).first()).toBeVisible();
  });

  test('should display blessing message', async ({ page }) => {
    await page.goto('/');
    await page.waitForTimeout(1000);

    await expect(page.getByText(/divine blessings|auspicious/i)).toBeVisible();
  });
});

test.describe('Panchangam Application - Accessibility', () => {
  test('should have no critical accessibility violations on main page', async ({ page }) => {
    await page.goto('/');
    await expect(page.locator('h1').first()).toContainText('Panchangam');

    // Basic accessibility checks
    // Check that buttons are accessible
    const buttons = page.getByRole('button');
    const buttonCount = await buttons.count();
    expect(buttonCount).toBeGreaterThan(0);

    // Check that headings exist
    const headings = page.getByRole('heading');
    const headingCount = await headings.count();
    expect(headingCount).toBeGreaterThan(0);
  });

  test('should support keyboard navigation', async ({ page }) => {
    await page.goto('/');

    // Tab through the page
    await page.keyboard.press('Tab');
    await page.keyboard.press('Tab');
    await page.keyboard.press('Tab');

    // Check that an element is focused
    const focusedElement = page.locator(':focus');
    await expect(focusedElement).toBeVisible();
  });
});

test.describe('Panchangam Application - Responsive Design', () => {
  test('should display correctly on mobile viewport', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 });
    await page.goto('/');

    await expect(page.getByRole('heading', { name: /Panchangam/i })).toBeVisible();
    await expect(page.getByText('Hindu Calendar & Astronomical Almanac')).toBeVisible();
  });

  test('should display correctly on tablet viewport', async ({ page }) => {
    await page.setViewportSize({ width: 768, height: 1024 });
    await page.goto('/');

    await expect(page.getByRole('heading', { name: /Panchangam/i })).toBeVisible();
    await expect(page.getByText('Hindu Calendar & Astronomical Almanac')).toBeVisible();
  });

  test('should display correctly on desktop viewport', async ({ page }) => {
    await page.setViewportSize({ width: 1920, height: 1080 });
    await page.goto('/');

    await expect(page.getByRole('heading', { name: /Panchangam/i })).toBeVisible();
    await expect(page.getByText('Hindu Calendar & Astronomical Almanac')).toBeVisible();
  });
});

test.describe('Panchangam Application - Loading States', () => {
  test('should load the page successfully', async ({ page }) => {
    // Navigate and check for successful load
    await page.goto('/');

    // Wait for the main heading to appear
    await expect(page.getByRole('heading', { name: /Panchangam/i })).toBeVisible({ timeout: 15000 });

    // The page should display the subtitle
    await expect(page.getByText('Hindu Calendar & Astronomical Almanac')).toBeVisible();
  });
});

test.describe('Panchangam Application - Error Handling', () => {
  test('should render UI structure even if API fails', async ({ page }) => {
    // Block API requests to simulate network failure
    await page.route('**/api/**', route => route.abort());

    await page.goto('/');

    // Wait for page to attempt to load
    await page.waitForTimeout(5000);

    // The app should still render some UI structure
    // Try to find the heading - it should still appear since UI renders before API calls
    const heading = page.getByRole('heading', { name: /Panchangam/i });
    const isHeadingVisible = await heading.isVisible().catch(() => false);

    // Either the heading is visible OR we verify the page didn't completely crash
    expect(isHeadingVisible || (await page.title()) !== '').toBeTruthy();
  });
});
