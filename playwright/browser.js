import playwright from 'playwright';

const userDataDir = `/Users/xue.a.yu/Library/Application Support/Google Chrome`;
const executablePath = `/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`;

// const browser = await playwright.chromium.launchPersistentContext(userDataDir, {
//     headless: true,
//     slowMo: 500,
//     executablePath: executablePath,
//     timeout: 600000
// });

// export default browser

export default async function launchBrowser(headless = false, timeout = 600000) {
    const context = await playwright.chromium.launchPersistentContext(userDataDir, {
        headless,
        slowMo: 500,
        executablePath: executablePath,
        timeout: timeout
    });

    return context;
}


