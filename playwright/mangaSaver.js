/* eslint-disable no-empty */
/* eslint-disable no-unused-vars */
import browser from './browser.js'
// import fetch from 'node-fetch'
import process from 'node:process';
import fs from 'fs';
import downloadWebpImg from './utils.js';
import shell from 'shelljs'

const args = process.argv.slice(2);

if (args.length < 2) {
    console.log("please input start url and start chapter num");
    process.exit(1);
}

const startURL = args[0];
const toChapter = parseInt(args[1]);
const startChapter = args.length == 3 ? args[2] : 1;


// const matchChapter = startURL.match(/\/(\d+)\.html/);
// const chapterNumber = matchChapter ? matchChapter[1] : null;
// const matchPrefix = startURL.match(/(.*\/)\d+\.html/);
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

const page = await browser.newPage();
await page.setExtraHTTPHeaders(headers);

// first page of the chapter
await page.goto(startURL, { timeout: 600000 });
const manga = await page.$eval('h1 a', element => element.textContent);
console.log(`start to parse manga ${manga}`);

if (!fs.existsSync(manga)) {
    fs.mkdirSync(manga, { recursive: true });
}

const restrictButton = await page.$('#checkAdult');
if (restrictButton) {
    await restrictButton.click();
    // await page.waitForLoadState();
}

for (let i = startChapter; i <= toChapter; i++) {
    console.info(`saving chapter${i}`)
    // Create a map to store chapter number and image url
    // const images = [];

    const lastPagesValue = await page.$eval('#pageSelect', select => select.options[select.options.length - 1].value);
    if (!lastPagesValue) {
        console.error(`get pages for chapter${i} failed`)
        break;
    } else {
        console.info(`pages for chapter${i}:`, lastPagesValue)
    }
    const lastPage = parseInt(lastPagesValue);
    const url = page.url().split('#')[0]

    // save image for a chapter
    for (let j = 1; j <= lastPage; j++) {
        // Get the image url for the current page
        const imageUrl = await page.$eval('#mangaFile', img => img.src);

        // Add the chapter number and image url to the map
        // images.push(imageUrl);
        console.log(`saving ${j}/${lastPage} image for chapter${i}: ${imageUrl}`);

        const image = `./${manga}/${(j).toString().padStart(3, '0')}_img.png`
        try {
            if (!fs.existsSync(image)) {
                await downloadWebpImg(imageUrl, {
                    'accept': 'image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8',
                    'referer': url
                }, image);
            } else {
                console.log(`skip image ${image}`)
            }
        } catch (err) {
            console.error(`saving ${image}: ${imageUrl} failed`, err)
            process.exit(1)
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

    console.info(`merging chapter${i}`);
    // zip all images 
    const zipResult = shell.exec(`zip -0 -rq ./${manga}/chapter${i}.zip ./${manga}/*.png`);
    if (zipResult.code !== 0) {
        console.error(`Error in zip images for chapter${i}`, zipResult.stderr);
        process.exit(1);
    }

    // merge to 1 pdf: magick ./第01卷/*.webp output.pdf
    // const chapter = `chapter${i.toString().padStart(2, '0')}.pdf`;
    // const magicResult = shell.exec(`magick ./${manga}/*.png ./${manga}/${chapter}`);
    // if (magicResult.code !== 0) {
    //     console.error(`Error in magick ${chapter}`, magicResult.stderr);
    //     process.exit(1);
    // } else {
    //     shell.exec(`rm ./${manga}/*.png`);
    // }
    shell.exec(`rm ./${manga}/*.png`);
    console.info(`saved chapter${i}`);

    if (i >= toChapter) break;
    await page.click('a.btn-red.nextC');
    await page.waitForNavigation({ timeout: 600000 });

    // const maxRetries = 4;
    // let retries = 0;
    // while (retries < maxRetries) {
    //     try {
    //         await page.waitForNavigation({ timeout: 600000 });
    //         break
    //     } catch (error) {
    //         console.log(`navigation timeout, retrying... (${i}/${maxretries})`);
    //         await page.reload();
    //     }
    // }

    // if (retries >= maxretries) {
    //     console.log('navigation failed after max retries');
    //     break;
    // }
}

await browser.close()