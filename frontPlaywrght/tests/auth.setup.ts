import {test as setup, expect} from '@playwright/test';
import * as dotenv from 'dotenv';

const authFile = '.auth/user.json';

setup('authenticate', async ({page}) => {
    dotenv.config();

    // Perform authentication steps. Replace these actions with your own.
    await page.goto('http://localhost:50003/player/login');


    await page.locator("#email").fill(process.env.LOGIN_USER);
    await page.locator("#password").fill(process.env.LOGIN_PASS);
    await page.locator("button[name='action']").click();
    // Wait until the page receives the cookies.
    //
    // Sometimes login flow sets cookies in the process of several redirects.
    // Wait for the final URL to ensure that the cookies are actually set.
    await page.waitForURL('http://localhost:50003/profile');

    await page.context().storageState({path: authFile});
});