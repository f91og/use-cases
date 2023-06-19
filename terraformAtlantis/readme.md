其实跟着这个https://www.vultr.com/docs/run-terraform-in-automation-with-atlantis/ 来做就可以了
deploy steps by systemd:
1. set user atlantis
2. setup args for command `atlantis server`,  prepare `repos.yaml`, tfstate and tfvars folder under /home/atlantis
   1. remember to set backend config in tf codes to point to above tfstate folder
3. set environment variable(should match the name defined in go codes) in atlantis.service in systemd
4. set ExecStart for execute `atlantis server`, then `systemctl daemon-reload`, `systemctl start atlantis`
5. setup custom provider if you use

坑
- bitbucket server的情形，配置时要记得在atlantis里也要配置它的url，这里官方文档的说明有缺失
- 使用官方提供的那个sts里，要加那个ATLANTIS_PORT环境变量，文档里没有给这个
- 为什么 allow apply after approve 不工作？🤮，因为没有在atlantis server的启动参数里设定，需要在repos.yaml里设置，不是默认启用的
- systemd来运行的模式下，如何为将secret信息设置为atlantis的环境变量？
  - 在systemd里设置Environment，注意前面不要加ATLANTIS_，这部分的实际使用也和官网有出入
- 文档上关于本地state file的设置貌似有错误，实际上是可以支持本地state的，在terraform里指定backend配置就可以了
- 如何把custom provider给设置到atlantis里去？ => 把build出来的provider放到机器上的atlantis的工作目录里就可以了，和本地执行是一样的，不过应该cicd自动部署上去的
- local_exec的docker权限问题 => 需要重启atlantis服务，在assign了权限之后


