import playwright from 'playwright';

const userDataDir = `/Users/xue.a.yu/Library/Application Support/Google Chrome`;

const browser = await playwright.chromium.launchPersistentContext(userDataDir, {
    headless: true,
    slowMo: 500,
    executablePath: "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
    timeout: 600000
});

export default browser
