**usage**
```sh
# download manga from chapter 1 to chapter 19
node mangaSaver.js url 19 1

# download manga from start to chapter 19, default start chapter is 1
node mangaSaver.js url 19
```

**cmd**
- directly run js: `node asakai.js`
- record code: `npx playwright codegen --channel chrome url --save-storage`
- `pnpm install commander`

**issues**
- https://stackoverflow.com/questions/63588714/node9374-warning-to-load-an-es-module-set-type-module

**note**
- js定义函数并导出
  ```js
  function writeBase64toImg(base64Str, filePath, fileName, format) {
    //......
  }

  export default writeBase64toImg // 这种情况下导出体的名字要和上面的一致
  ``` 
- js保存网页图片的几种方法: https://intoli.com/blog/saving-images/
- mocha单元测试工具
- npx
  - npx: allows you to run a command-line tool without having to install it globally on your system, or without having to specify the path to the tool's executable
  - When you run a command using npx, it looks for the command locally in the ./node_modules/.bin directory of your project, and if it doesn't find it, it downloads and installs it **temporarily**, and then runs it
- 为什么不能使用 `const assert = require("assert");`? <- commonJs module plan vs es module plan
- eslint提示不对，提示 xxx is not defined, 但是事实上是能找到的
- save canvas image
  - https://piccoma.com/web/viewer/6081/912240
  - https://intoli.com/blog/saving-images/
- [Does page.click automatically page.waitForNavigation?](https://github.com/microsoft/playwright/issues/2078)
- 等待元素
  ```js
  await page.click(searchResultSelector);
  // Locate the full title with a unique string
  const textSelector = await page.waitForSelector(
    'text/Customize and automate'
  );
  ```

import foo from 'module' imports the default export as foo => 导入module文件里的default export，并将其赋值给foo
import {foo} from 'module' imports the named export foo => 导入module文件里的named export foo