/* eslint-disable no-empty */
/* eslint-disable no-unused-vars */
import process from 'node:process';
import {Command} from 'commander';
import cliProgress from 'cli-progress';
import fs from 'fs';
import downloadWebpImg from './utils.js';
import launchBrowser from './browser.js';

const program = new Command();
program
  .version("0.0.1")
  .option("-h, --headless", "run headless browser") // 命令参数中出现了--headless才是true
  .option("-s, --startUrl [string]", "start chapter url")
  .option("-c, --startChapter <number>", "start chapter num", 1)
  .option("-t, --toChapter <number>", "to chapter num", 4)
  .option("-w, --timeout <number>", "wait time", 60000)
  .parse(process.argv);

const options = program.opts();
console.log(options);

if (!(options.startUrl || options.toChapter))  {
    console.log("please input start url and to chapter num");
    process.exit(1);
}

const startChapterUrl = options.startUrl;
const startChapter = parseInt(options.startChapter);
const toChapter = parseInt(options.toChapter);

// const matchChapter = startChapterUrl.match(/\/(\d+)\.html/);
// const chapterNumber = matchChapter ? matchChapter[1] : null;
// const matchPrefix = startChapterUrl.match(/(.*\/)\d+\.html/);
// const prefix = matchPrefix ? matchPrefix[1] : null;

// if (chapterNumber !== null && prefix !== null) {
//     console.log("chapter start number", chapterNumber);
//     console.log("manga url prefix", prefix);
// } else {
//     console.log('无法提取chapter start number or manga url prefix', chapterNumber, prefix);
//     process.exit(1)
// }

// const chapterLinks = [];

// const chapterNum = parseInt(chapterNumber)
// for (let i = 0; i < numCounts; i++) {
//     const link = `${prefix}${chapterNum + i}.html`;
//     chapterLinks.push(link);
// }

const headers = {
    'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7',
    'accept-language': 'zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6'
};

const browser = await launchBrowser(options.headless, options.timeout)
const page = await browser.newPage();
await page.setExtraHTTPHeaders(headers);

// first page of the chapter
await page.goto(startChapterUrl, { timeout: 600000 });
const manga = await page.$eval('h1 a', element => element.textContent);
console.log(`start to download manga ${manga}`);

const mangaFolder = `./manga/${manga}`;

const restrictButton = await page.$('#checkAdult');
if (restrictButton) {
    await restrictButton.click();
}

for (let i = startChapter; i <= toChapter; i++) {
    const currentChapter = await page.$eval('h2', element => element.textContent);
    const lastPage = parseInt(await page.$eval('#pageSelect', select => select.options[select.options.length - 1].value));

    if (!verifyChapter(currentChapter, `第${i.toString().padStart(2, '0')}卷`, lastPage)) break;

    const url = page.url().split('#')[0];
    console.info(`saving ${currentChapter}, total ${lastPage} pages, url: ${url}`);
    const bar = new cliProgress.SingleBar({}, cliProgress.Presets.shades_classic);
    bar.start(lastPage, 0);

    const chapterFolder = `${mangaFolder}/${currentChapter}`
    if (!fs.existsSync(chapterFolder)) {
        fs.mkdirSync(chapterFolder, { recursive: true });
    }

    // save images for a chapter
    for (let j = 1; j <= lastPage; j++) {
        // Get the image url for the current page
        const imageUrl = await page.$eval('#mangaFile', img => img.src);

        // console.log(`saving ${j}/${lastPage} image for chapter${i}: ${imageUrl}`);

        const image = `${chapterFolder}/${j.toString().padStart(3, '0')}_img.png`;
        try {
            if (!fs.existsSync(image)) {
                await downloadWebpImg(imageUrl, {
                    'accept': 'image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8',
                    'referer': url
                }, image);
            }
            bar.increment();
        } catch (err) {
            console.error(`saving ${image}: ${imageUrl} failed`, err);
            process.exit(1);
        }

        if (j !== lastPage) {
            const nextButton = await page.$('#next'); // selector: elementId
            await nextButton.click();
            await page.waitForLoadState();
        }
    }

    // if (images.length != pages) {
    //     fs.writeFileSync(`./${manga}/chapter_${i}_failed.txt`, images.join('\n'), 'utf8');
    // } else {
    //     fs.writeFileSync(`./${manga}/chapter_${i}.txt`, images.join('\n'), 'utf8');
    // }
    // console.log(chapterImageMap);

    // zip all images
    // const zipResult = shell.exec(`zip -0 -rq ./${mangaFolder}/chapter${i}.zip ./${mangaFolder}/*.png`);
    // if (zipResult.code !== 0) {
    //     console.error(`Error in zip images for chapter${i}`, zipResult.stderr);
    //     process.exit(1);
    // }

    // shell.exec(`rm ./${mangaFolder}/*.png`);

    // console.info(`saved chapter${i}`);
    bar.stop();

    if (i >= toChapter) break;
    await page.click('a.btn-red.nextC');
    await page.waitForNavigation({timeout: 60000});
}

await browser.close()

function verifyChapter(currentChapter, chapterSaving, lastPage) {
    if (lastPage < 70) {
        console.error(`no enough chapter pages: ${lastPage}`); return false;
    }
    if (currentChapter.startsWith("第") && currentChapter !== chapterSaving) {
        console.error(`chapter error, current: ${currentChapter}, saving: ${chapterSaving}`); return false;
    }

    return true;
}