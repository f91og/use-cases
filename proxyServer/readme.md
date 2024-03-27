init project
- mkdir my-server && cd my-server
- npm init -y
- npm install express
- input server start cmd to package.json
  ```
  "scripts": {
  "start": "node server.js"
  },
  ```

将localhost变为一个代理服务器，如果客户端不能访问某个url但是自己的机器可以的话，可以将自己的机器变成中转服务器
- ngrok提供公网访问的url，对这个公网url的访问会转到本机
- 然后用代理服务器，将这个流量从本机发到目的url

实现方案
- server implemented by codes(express, python, etc)
- nginx
- 