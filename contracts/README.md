### Developement instructions

- `$ yarn install` - _Install all dependencies_
- `$ echo PK="PRIVATE_KEY" > .env` - _Add testing private key_
- `$ npx hardhat compile` - _Compile all contracts_
- `$ npx hardhat test` - _Run all tests_

- Migrations are inside `scripts/` folder.
- Tests are inside `test/` folder.

node 版本：v18.19.1

部署流程

1. `npx hardhat run --network local scripts/deployment/deploy_launchpad_token.js`
2. `npx hardhat run --network local scripts/deployment/deploy_airdrop_launchpad.js`
3. `npx hardhat run --network local scripts/deployment/deploy_farm.js`
4. `npx hardhat run --network local scripts/deployment/deploy_ido.js`
5. `npx hardhat run --network local scripts/deployment/deploy_sales_token.js`
6. `npx hardhat run --network local scripts/deployment/deploy_sales.js`
7. `npx hardhat run --network local scripts/deployment/deploy_tge.js`

## makefile 命令解释

将脚本命令简化成了 make 直接可执行的命令，写到了 makefile 中

需要安装`make`编译工具，如果没有安装，直接执行 makefile 文件中的 js 命令也可以

命令列表：

- farm
  LaunchPad 代币&空投&Farm 及指令
- ido
  LaunchPad 代币&空投&Farm&IDO 流程及指令

#### 使用脚本生成 sql
    首先需要拿到链上 json 数据作为参数（make ido后的结果字符串），使用数据维护脚本更新数据。

    # 用法如下:
    #  sh generate_update_data.sh [json_file] [server_url]
    # 样例:
    #  sh generate_update_data.sh '{"saleAddress":"0x330981485Dbd4EAcD7f14AD4e6A1324B48B09995","saleToken":"0x998abeb3E57409262aE5b751f60747921B33613E","saleOwner":"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266","tokenPriceInEth":"100000000000","totalTokens":"10000000000000000000000000","saleEndTime":1767708887,"tokensUnlockTime":1767709117,"registrationStart":1767708697,"registrationEnd":1767708817,"saleStartTime":1767708827}' localhost:8080
