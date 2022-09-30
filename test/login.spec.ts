// example.spec.ts
import { test, expect } from '@playwright/test'
import { env } from '../package.json'
const VITE_SERVER_ADDRESS = `http://127.0.0.1:${env.PORT || 3344}`

test("login sync as david.johansson", async ({page}) => {
  // Go to http://localhost:3344/
  await page.goto('http://localhost:3344/');

  // Go to http://localhost:3344/#/
  await page.goto('http://localhost:3344/#/');

  // Go to http://localhost:3344/#/login
  await page.goto('http://localhost:3344/#/login');

  // Click input[type="email"]
  await page.locator('input[type="email"]').click();

  // Fill input[type="email"]
  await page.locator('input[type="email"]').fill('base@user.com');

  // Press Enter
  await page.locator('input[type="email"]').press('Enter');

  // Fill input[type="password"]
  await page.locator('input[type="password"]').fill('s0m3_p4ssw0rd');

  // Press Enter
  await Promise.all([
    page.waitForNavigation(/*{ url: 'http://localhost:3344/#/' }*/),
    page.locator('input[type="password"]').press('Enter')
  ]);

  await page.locator('.ql-editor').type('Ahet', { delay: 25 })
  await page.locator('.ql-editor').press("Space");
  await page.locator('.ql-editor').type('mhet', { delay: 250})
  await page.locator('.ql-editor').press("Period");
  await page.locator('.ql-editor:has-text("Allmänhet möjlighet.")').click();
  await page.locator('.ql-editor').press('F4', { timeout: 125 });
  await page.locator('.ql-editor').type('Samhälls', { delay: 25 })
  // Press F2 with modifiers
  await page.locator('.ql-editor').press('Shift+F2');
  // Fill [placeholder="Förkortning"]
  //
  await page.locator('[placeholder="Förkortning"]').fill('shs');
  await page.locator('[placeholder="Förkortning"]').press('Enter');
  await page.locator('.ql-editor:has-text("samhälls")').press('F4');
});

test('test local as info@sttylus.se then as base user', async ({ page }) => {

  // Go to http://localhost:3344/
  await page.goto('http://localhost:3344/');

  // Go to http://localhost:3344/#/
  await page.goto('http://localhost:3344/#/');

  // Go to http://localhost:3344/#/login
  await page.goto('http://localhost:3344/#/login');

  // Double click text=Info
  await Promise.all([
    page.waitForNavigation(/*{ url: 'http://localhost:3344/#/' }*/),
    page.locator('text=Info').dblclick()
  ]);


  await page.locator('.ql-editor').type('Ahet', { delay: 25 })
  await page.locator('.ql-editor').press("Space");
  await page.locator('.ql-editor').type('mhet', { delay: 250})
  await page.locator('.ql-editor').press("Period");
  await page.locator('.ql-editor:has-text("Allmänhet möjlighet.")').click();

  await page.locator('.ql-editor').press('F4', { timeout: 125 });
  //
  // Press F5
  await page.locator('.ql-editor').press('F5');
  // Click text=Avsluta STTylus
  await page.locator('text=Avsluta STTylus').click();
  // Click button:has-text("Logga ut")
  await page.locator('button:has-text("Logga ut")').click();
  await expect(page).toHaveURL('http://localhost:3344/#/login');

  // Double click text=base user
  await Promise.all([
    page.waitForNavigation(/*{ url: 'http://localhost:3344/#/' }*/),
    page.locator('text=base user').dblclick()
  ]);

  await page.locator('.ql-editor').type('Väf', { delay: 25 })
  await page.locator('.ql-editor').press('Control+ ');
  await page.locator('.ql-editor').type('s', { delay: 25 })
  await page.locator('.ql-editor').press('Control+ ');
  await page.locator('.ql-editor').type('sht', { delay: 25 })
  await page.locator('.ql-editor').press('Space');

});
