const hre = require("hardhat");
const { ethers } = hre;
const ethUtil = require("ethereumjs-util");

/**
 * 脚本用于模拟 Java 后端的签名逻辑
 * Java 后端的签名方式与 Solidity 合约中的方式不同：
 * 
 * Java: 
 * 1. 去掉 0x 前缀
 * 2. 连接两个地址（小写）
 * 3. keccak256(连接后的字符串)
 * 4. 加上 Ethereum 消息前缀后签名
 * 
 * Solidity:
 * 1. abi.encodePacked(user, address(this))
 * 2. keccak256(编码后的 bytes)
 * 3. 加上 Ethereum 消息前缀后 recover
 */

// ==================== 配置区域 ====================
// 请在这里填写需要签名的参数

// 用户地址（必需）
const USER_ADDRESS = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266";

// Sale 合约地址（必需）
const SALE_ADDRESS = "0x8acd85898458400f7db866d53fcff6f0d49741ff";

// 签名者私钥（必需，用于签名的私钥，recover 出来的地址必须是 admin）
const SIGNER_PRIVATE_KEY = process.env.SIGNER_PRIVATE_KEY || "";

// ==================== 配置区域结束 ====================

function generateSignature(digest, privateKey) {
  // prefix with "\x19Ethereum Signed Message:\n32"
  const prefixedHash = ethUtil.hashPersonalMessage(ethUtil.toBuffer(digest));

  // sign message
  const { v, r, s } = ethUtil.ecsign(
    prefixedHash,
    Buffer.from(privateKey.replace("0x", ""), "hex")
  );

  // generate signature by concatenating r(32), s(32), v(1) in this order
  const vb = Buffer.from([v]);
  const signature = Buffer.concat([r, s, vb]);

  return "0x" + signature.toString("hex");
}

/**
 * 模拟 Java 后端的签名逻辑
 * Java 代码逻辑：
 * String contractAddr = Numeric.cleanHexPrefix(contractAddress);
 * String userAddr = Numeric.cleanHexPrefix(userAddress);
 * String concat = userAddr.concat(contractAddr).toLowerCase();
 * String hex = Numeric.prependHexPrefix(concat);
 * sign = encodeService.sign(hex); // 内部会 keccak256(hex) 然后签名
 */
function signRegistrationJavaStyle(userAddress, contractAddress, privateKey) {
  // 1. 去掉 0x 前缀
  const contractAddr = contractAddress.replace(/^0x/i, "");
  const userAddr = userAddress.replace(/^0x/i, "");
  
  // 2. 连接两个地址（小写）
  const concat = (userAddr + contractAddr).toLowerCase();
  
  // 3. 加上 0x 前缀
  const hex = "0x" + concat;
  
  console.log("\n=== Java Style Signing Process ===");
  console.log(`Step 1 - User Address (no prefix): ${userAddr}`);
  console.log(`Step 2 - Contract Address (no prefix): ${contractAddr}`);
  console.log(`Step 3 - Concatenated (lowercase): ${concat}`);
  console.log(`Step 4 - Final hex string: ${hex}`);
  
  // 4. keccak256(hexString) - Java 的 Hash.sha3 就是 keccak256
  const messageHash = ethers.keccak256(hex);
  console.log(`Step 5 - keccak256(hex): ${messageHash}`);
  
  // 5. 生成签名（会自动加上 Ethereum 消息前缀）
  return generateSignature(messageHash, privateKey);
}

/**
 * Solidity 合约的签名逻辑（用于对比）
 */
function signRegistrationSolidityStyle(userAddress, contractAddress, privateKey) {
  // compute keccak256(abi.encodePacked(user, address(this)))
  const digest = ethers.keccak256(
    ethers.solidityPacked(
      ["address", "address"],
      [userAddress, contractAddress]
    )
  );
  
  console.log("\n=== Solidity Style Signing Process ===");
  console.log(`User Address: ${userAddress}`);
  console.log(`Contract Address: ${contractAddress}`);
  console.log(`keccak256(abi.encodePacked(user, contract)): ${digest}`);
  
  return generateSignature(digest, privateKey);
}

async function main() {
  // 获取网络名称
  const network = hre.network.name;
  console.log(`Network: ${network}`);

  // 验证参数
  if (!USER_ADDRESS || !ethers.isAddress(USER_ADDRESS)) {
    console.error("Error: Invalid USER_ADDRESS");
    process.exit(1);
  }

  if (!SALE_ADDRESS || !ethers.isAddress(SALE_ADDRESS)) {
    console.error("Error: Invalid SALE_ADDRESS");
    process.exit(1);
  }

  if (!SIGNER_PRIVATE_KEY || !SIGNER_PRIVATE_KEY.startsWith("0x")) {
    console.error("Error: Invalid SIGNER_PRIVATE_KEY");
    process.exit(1);
  }

  // 验证私钥对应的地址
  const wallet = new ethers.Wallet(SIGNER_PRIVATE_KEY);
  const signerAddress = wallet.address;
  console.log("\n=== Simulate Java Backend Signature ===");
  console.log(`User Address: ${USER_ADDRESS}`);
  console.log(`Sale Address: ${SALE_ADDRESS}`);
  console.log(`Signer Address (will be recovered): ${signerAddress}`);

  // 生成 Java 风格的签名
  const javaSignature = signRegistrationJavaStyle(USER_ADDRESS, SALE_ADDRESS, SIGNER_PRIVATE_KEY);
  
  // 生成 Solidity 风格的签名（用于对比）
  const soliditySignature = signRegistrationSolidityStyle(USER_ADDRESS, SALE_ADDRESS, SIGNER_PRIVATE_KEY);

  console.log("\n=== Generated Signatures ===");
  console.log(`Java Style Signature: ${javaSignature}`);
  console.log(`Solidity Style Signature: ${soliditySignature}`);
  console.log(`Signatures are ${javaSignature === soliditySignature ? "SAME" : "DIFFERENT"}`);

  // 验证 Java 风格的签名在合约中是否有效
  try {
    const sale = await hre.ethers.getContractAt("LaunchPadSale", SALE_ADDRESS);
    
    console.log("\n=== Verification Results ===");
    
    const javaIsValid = await sale.checkRegistrationSignature(javaSignature, USER_ADDRESS);
    console.log(`Java Style Signature is valid: ${javaIsValid}`);
    
    const solidityIsValid = await sale.checkRegistrationSignature(soliditySignature, USER_ADDRESS);
    console.log(`Solidity Style Signature is valid: ${solidityIsValid}`);
    
    if (!javaIsValid) {
      console.log("\n⚠️  WARNING: Java backend signature will NOT work with the contract!");
      console.log("The Java backend uses a different signing method than the contract expects.");
      console.log("You need to fix the Java backend to use abi.encodePacked instead of string concatenation.");
    }
  } catch (error) {
    console.log("\n=== Verification Skipped ===");
    console.log(`Could not verify signature: ${error.message}`);
  }

  console.log("\n=== Java Backend Return Value ===");
  console.log(`If you call the Java API with:`);
  console.log(`  userAddress: ${USER_ADDRESS}`);
  console.log(`  contractAddress: ${SALE_ADDRESS}`);
  console.log(`\nThe Java backend will return:`);
  console.log(`  ${javaSignature}`);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
