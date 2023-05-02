import fs from 'fs';
import { Buffer } from 'node:buffer';
import fetch from 'node-fetch';

// function writeBase64toImg(base64Str, filePath, fileName, format) {
//     const path = filePath + '/' + fileName + '.' + format;
//     const dataBuffer = Buffer.from(base64Str, 'base64')
//     fs.writeFile(path, dataBuffer, function (err) {//用fs写入文件
//         if (err) {
//             console.log(err);
//         } else {
//             console.log('写入成功！');
//         }
//     })
// }

async function downloadWebpImg(url, headers, filePath) {
    const response = await fetch(url, {headers: headers});

    if (!response.ok) {
        throw new Error(`Could not fetch image, (status ${response.status}`);
    }

    const data = await response.arrayBuffer()
    const dataBuffer = Buffer.from(data, 'base64')
    fs.writeFileSync(filePath, dataBuffer, function (err) {//用fs写入文件
        if (err) {
            throw new Error(`Could not write image, ${err}`);
        } else {
            console.log(`image ${url} wrote to file ${filePath}`);
        }
    })
}

export default downloadWebpImg