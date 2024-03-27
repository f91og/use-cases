const express = require('express');
const httpProxy = require('http-proxy');
const app = express();
const proxy = httpProxy.createProxyServer();

// 将请求转发到 https://chat.openai.com/
app.use('/', (req, res) => {
  proxy.web(req, res, { target: 'https://chat.openai.com/' });
});

// 启动服务器监听在 3000 端口
app.listen(3000, () => {
  console.log('服务器已启动，监听在端口 3000');
});
