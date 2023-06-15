terraform {
  required_providers {
    onecloud = {
      version = "0.1.0"
      source  = "xxx"
    }
  }
  backend "local" {
    path = "/home/atlantis/tfstate/atlantis-test.tfstate" # 这里需要match部署了atlantis服务器上的目录路径
  }
}


provider "onecloud" {
  #  client_id = "" # 如果需要从环境变量设置敏感数据，则不设置，否则设置为空就会取空值
  #  client_secret = ""
}