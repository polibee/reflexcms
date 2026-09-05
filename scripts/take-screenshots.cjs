// Screenshot harness: uses the playwright-core module already installed at
// D:/codex/cms/tests-browser/node_modules and the ms-playwright browsers.
const path = require('path');
const pw = require('D:/codex/cms/tests-browser/node_modules/playwright-core');

(async () => {
  const browser = await pw.chromium.launch({
    executablePath: 'C:/Users/hongboli/AppData/Local/ms-playwright/chromium-1234/chrome-win64/chrome.exe'
  });
  const ctx = await browser.newContext({ viewport: { width: 1280, height: 800 } });
  const page = await ctx.newPage();

  const out = 'D:/codex/reflexcms/docs/screenshots/';

  // 1. public home (hybrid)
  await page.goto('http://localhost:3000/', { waitUntil: 'networkidle', timeout: 30000 }).catch(() => {});
  await page.waitForTimeout(1500);
  await page.screenshot({ path: out + 'home-hybrid.png' });
  console.log('shot home-hybrid');

  // 2. forums
  await page.goto('http://localhost:3000/forums', { waitUntil: 'networkidle', timeout: 30000 }).catch(() => {});
  await page.waitForTimeout(1200);
  await page.screenshot({ path: out + 'forums.png' });
  console.log('shot forums');

  // 3. topic detail
  await page.goto('http://localhost:3000/topics/1010', { waitUntil: 'networkidle', timeout: 30000 }).catch(() => {});
  await page.waitForTimeout(1500);
  await page.screenshot({ path: out + 'topic-detail.png' });
  console.log('shot topic-detail');

  // 4. user profile
  await page.goto('http://localhost:3000/users/1', { waitUntil: 'networkidle', timeout: 30000 }).catch(() => {});
  await page.waitForTimeout(1500);
  await page.screenshot({ path: out + 'user-profile.png' });
  console.log('shot user-profile');

  // 5. shop
  await page.goto('http://localhost:3000/shop', { waitUntil: 'networkidle', timeout: 30000 }).catch(() => {});
  await page.waitForTimeout(1200);
  await page.screenshot({ path: out + 'shop.png' });
  console.log('shot shop');

  // 6. admin dashboard (login if redirected)
  await page.goto('http://localhost:3000/login', { waitUntil: 'networkidle', timeout: 30000 }).catch(() => {});
  await page.waitForTimeout(800);
  await page.fill('input[type="email"]', 'admin@reflexcms.dev');
  await page.fill('input[type="password"]', 'ReflexCMS@2026');
  await page.click('button[type="submit"]');
  await page.waitForTimeout(2000);
  await page.goto('http://localhost:3000/admin', { waitUntil: 'networkidle', timeout: 30000 }).catch(() => {});
  await page.waitForTimeout(1500);
  await page.screenshot({ path: out + 'admin-dashboard.png' });
  console.log('shot admin-dashboard');

  // 7. admin products
  await page.goto('http://localhost:3000/admin/products', { waitUntil: 'networkidle', timeout: 30000 }).catch(() => {});
  await page.waitForTimeout(1200);
  await page.screenshot({ path: out + 'admin-products.png' });
  console.log('shot admin-products');

  // 8. admin settings
  await page.goto('http://localhost:3000/admin/settings', { waitUntil: 'networkidle', timeout: 30000 }).catch(() => {});
  await page.waitForTimeout(1200);
  await page.screenshot({ path: out + 'admin-settings.png' });
  console.log('shot admin-settings');

  await browser.close();
  console.log('ALL DONE');
})().catch(e => { console.error('ERR:', e.message); process.exit(1); });
