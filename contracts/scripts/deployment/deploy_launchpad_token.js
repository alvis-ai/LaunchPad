const hre = require("hardhat");
const { saveContractAddress } = require("../utils");

/**
 * 此脚本用于部署LaunchPad代币业务
 */
async function main() {
  const tokenName = "LaunchPad";
  const symbol = "LaunchPad";
  const totalSupply = "1000000000000000000000000000";
  const decimals = 18;

  const LaunchPad = await hre.ethers.getContractFactory("LaunchPadToken");
  //构造时初始化一定数量的LaunchPad代币
  const token = await LaunchPad.deploy(tokenName, symbol, totalSupply, decimals);
  await token.waitForDeployment();
  console.log("LaunchPad deployed to: ", await token.getAddress());

  saveContractAddress(hre.network.name, "LaunchPad-TOKEN", await token.getAddress());
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
