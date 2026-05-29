import contractAddresses from '../../../contracts/deployments/contract-addresses.json';
/* eslint-disable import/no-anonymous-default-export */
export default [
  // {
  //   chainId: 56,
  //   depositTokenAddress: "0x2a50b9D094f3aceA0Ddd7E6077dD6dadA51ac36d",
  //   earnedTokenAddress: "0x9eBBEB7f6b842377714EAdD987CaA4510205107A",
  //   stakingAddress: "0xF67601B52cd39b6B2FFdDfBE02Aed8D92bC557ca",
  //   poolId: 2,
  //   available: true,
  //   depositSymbol: "LPT",
  //   earnedSymbol: "BRE",
  //   title: "BNB/BRE LPT",
  //   depositLogo:
  //     "https://bobabrewery.oss-ap-southeast-1.aliyuncs.com/bnb-logo.svg",
  //   earnedLogo:
  //     "https://bobabrewery.oss-ap-southeast-1.aliyuncs.com/Brewery32x32.png",
  //   getLptHref:
  //     "https://pancakeswap.finance/add/BNB/0x9eBBEB7f6b842377714EAdD987CaA4510205107A",
  //   aprRate: 3.15 / 20,
  //   aprUrl: "/boba/apr/bsc",
  // },
  // {
  //   chainId: 288,
  //   depositTokenAddress: "0x8Ebb37552fFaB6997baD219213019718abea2cE8",
  //   earnedTokenAddress: "0x3A93bD0fA10050d206370eeA53276542A105C885",
  //   stakingAddress: "0x4a050E34174F982Bcb9731871b3Af9140841b463",
  //   poolId: 2,
  //   available: true,
  //   depositSymbol: "LPT",
  //   earnedSymbol: "BRE",
  //   title: "ETH/BRE LPT",
  //   depositLogo:
  //     "https://bobabrewery.oss-ap-southeast-1.aliyuncs.com/eth-logo.png",
  //   earnedLogo:
  //     "https://bobabrewery.oss-ap-southeast-1.aliyuncs.com/Brewery32x32.png",
  //   getLptHref:
  //     "https://oolongswap.com/#/add/ETH/0x3A93bD0fA10050d206370eeA53276542A105C885",
  //   aprRate: 0.35 / 66.6,
  //   aprUrl: "/boba/apr/boba",
  // },
  { // 合约切换报错未解决
    chainId: 11155111,
    depositTokenAddress: contractAddresses['sepolia']['LaunchPad-TOKEN'],
    earnedTokenAddress: contractAddresses['sepolia']['LaunchPad-TOKEN'],
    stakingAddress: contractAddresses['sepolia']['FarmingLaunchPad'],
    poolId: 0,
    available: true,
    depositSymbol: "FLaunchPad",
    earnedSymbol: "LaunchPad",
    title: "Sepolia/LaunchPad FLaunchPad",
    depositLogo:
      "https://bobabrewery.oss-ap-southeast-1.aliyuncs.com/bnb-logo.svg",
    earnedLogo:
      "https://bobabrewery.oss-ap-southeast-1.aliyuncs.com/Brewery32x32.png",
    getLptHref:
      "https://pancakeswap.finance/add/BNB/0x9eBBEB7f6b842377714EAdD987CaA4510205107A",
    aprRate: 3.15 / 20,
    aprUrl: "/boba/apr/bsc",
  },
  {
    chainId: 31337,
    depositTokenAddress: contractAddresses['local']['LaunchPad-TOKEN'],
    earnedTokenAddress: contractAddresses['local']['LaunchPad-TOKEN'],
    stakingAddress: contractAddresses['local']['FarmingLaunchPad'],
    poolId: 0,
    available: true,
    depositSymbol: "FLaunchPad",
    earnedSymbol: "LaunchPad",
    title: "Local/LaunchPad FLaunchPad",
    depositLogo:
      "https://bobabrewery.oss-ap-southeast-1.aliyuncs.com/bnb-logo.svg",
    earnedLogo:
      "https://bobabrewery.oss-ap-southeast-1.aliyuncs.com/Brewery32x32.png",
    getLptHref:
      "https://pancakeswap.finance/add/BNB/0x9eBBEB7f6b842377714EAdD987CaA4510205107A",
    aprRate: 3.15 / 20,
    aprUrl: "/boba/apr/bsc",
  },
];
