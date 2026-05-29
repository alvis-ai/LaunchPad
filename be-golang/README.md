# LaunchPad-BE-GOLANG

用 Golang 语言主流的 Web 框架 Gin 重写后台服务，工程目录命名为 be-golang。

相关文档：
- [Gin 文档](https://gin-gonic.com/zh-cn/)
- [Gorm 文档](https://gorm.io/docs/)
- [Geth 官方文档](https://geth.ethereum.org/docs)

LaunchPad-Launchpad后台为IDO流程提供服务，farm流程无需启动后台。部署使用docker，方式同Java版launchPad。

IDO 启动流程回顾

1. 启动本地链
2. 执行`make ido`部署合约（如果系统内没有make命令，可直接看make file中的`npx ...`命令，在控制台按顺序执行）
3. 启动后台服务
4. IDO业务操作（sale end后，提取MCK代币大概有1分钟延迟，等待一下withdraw才可成功到账）

## 业务功能

后台提供产品合约管理、链上事件监听和签名服务。

产品合约管理支持查询产品信息、获取列表、更新信息等操作。通过产品ID可以查询到销售地址、代币地址、价格、时间等基础信息。

链上事件监听，监听销售工厂合约的部署事件（SaleDeployed），自动捕获新部署的销售合约。同时也监听销售合约的状态变化事件，包括SaleCreated、StartTimeSet、RegistrationTimeSet等。当监听到链上事件后，从链上查询销售合约信息并同步到数据库。

签名服务提供注册签名和参与签名功能，用于IDO注册流程和购买流程的权限验证。

数据存储使用MySQL，通过GORM进行数据库操作，存储产品合约的销售配置、代币信息、时间节点等数据。

## 部署项目

```
cd deployment
docker-compose up -d
```

## 目录

```
├── api
│   ├── handler
│   │   └── xxx.go        # 用户相关的业务逻辑处理函数
│   ├── request
│   │   └── xxx.go        # 请求参数的定义和验证
│   ├── response
│       └── xxx.go        # 响应数据的定义
├── config
│   └── config.go          # 配置文件的加载、管理
├── database
│   └── database.go        # 数据库的初始化和连接管理
├── docs                   # 项目文档或接口文档
├── middleware
│   └── auth.go            # 中间件，如认证、日志等
├── model
│   └── xxx.go            # GORM 模型定义
├── repository
│   └── xxx_repository.go # 数据库访问逻辑，封装数据库操作
├── route
│   └── router.go          # 路由定义
├── service
│   └── xxx_service.go    # 业务逻辑层，调用 repository 处理业务
├── utils
│   └──xxx_utils.go        # 各种工具类
├── main.go                # 主入口文件
├── go.mod                 # Go modules 管理文件
├── go.sum                 # Go modules 依赖文件
├── output.log             # 项目日志, 运行时会自动生成
├── README.md  
├── .env                   # 环境变量文件, 不用也可以去掉
└── .gitignore   
```

## 本地开发

安装依赖：
```bash
go get
```

启动项目：
```bash
go run main.go
```

#### 使用脚本生成 sql
    首先需要拿到链上 json 数据作为参数（make ido后的结果字符串），使用数据维护脚本更新数据。

    # 用法如下:
    #  sh generate_update_data.sh [json_file] [server_url]
    # 样例:
    #  sh generate_update_data.sh '{"saleAddress":"0x330981485Dbd4EAcD7f14AD4e6A1324B48B09995","saleToken":"0x998abeb3E57409262aE5b751f60747921B33613E","saleOwner":"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266","tokenPriceInEth":"100000000000","totalTokens":"10000000000000000000000000","saleEndTime":1767708887,"tokensUnlockTime":1767709117,"registrationStart":1767708697,"registrationEnd":1767708817,"saleStartTime":1767708827}' localhost:8080
