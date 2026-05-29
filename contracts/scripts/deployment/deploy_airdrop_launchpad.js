const hre = require("hardhat");
const { saveContractAddress, getSavedContractAddresses } = require("../utils");

/**
 * 此脚本用于部署空投业务
 */
async function main() {
  // get launchpad token address from contract address file
  const launchpadTokenAddress =
    getSavedContractAddresses()[hre.network.name]["LaunchPad-TOKEN"];
  console.log("launchpadTokenAddress: ", launchpadTokenAddress);

  const air = await hre.ethers.getContractFactory("Airdrop");
  const Air = await air.deploy(launchpadTokenAddress);
  await Air.waitForDeployment();
  console.log("Air deployed to: ", await Air.getAddress());

  saveContractAddress(hre.network.name, "Airdrop-LaunchPad", await Air.getAddress());
  // send launchpad token to airdrop contract
  const launchpadToken = await hre.ethers.getContractAt("LaunchPadToken", launchpadTokenAddress);
  let tx = await launchpadToken.transfer(
    await Air.getAddress(),
    ethers.parseEther("10000")
  );
  // wait for transfer
  await tx.wait();
  // get airdrop balance of launchpad token
  const balance = await launchpadToken.balanceOf(await Air.getAddress());
  console.log("Airdrop balance of LaunchPad token: ", ethers.formatEther(balance));
  // test airdrop
  tx = await Air.withdrawTokens();
  await tx.wait();
  // get airdrop balance of launchpad token
  const balanceAfter = await launchpadToken.balanceOf(await Air.getAddress());
  console.log(
    "Airdrop balance of LaunchPad token after withdrawTokens: ",
    ethers.formatEther(balanceAfter)
  );
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
