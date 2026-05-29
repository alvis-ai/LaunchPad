# 极简部署说明

本文档介绍了快速部署 IDO 的核心步骤，具体流程可以详见 launchpad-contract 和 launchpad-be 下面的 README

## 合约部署

合约部署使用 contracts/Makefile 对合约进行了快速部署，相关业务参数使用了可以快速启动的缺省值
前往 contracts 目录
将 contracts 目录下的.env.exmaple 改为.env
npx hardhat node
make ido

## 后端部署

需要确保在开发环境中安装好 jdk8, maven(3.6.3+) && docker(20.10.17+) 和 docker compose
前往 launchpad-be 目录
使用 deploy.sh 完成部署
sh deploy.sh
如果安装 docker 出现类似 The requested image's platform (linux/arm64) does not match the detected host platform (linux/amd64/v3) and no specific platform was requested 的错误，需要替换 portal-api 下 Dockerfile 里的 jdk 基础镜像为 amd 版本

## 前端部署

前端使用 yarn 进行构建和启动
yarn install
yarn dev
