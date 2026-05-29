/* eslint-disable react-hooks/rules-of-hooks */
import { useEffect, useState } from 'react'
import { useAppSelector, useAppDispatch } from "@src/redux/hooks";
import { useWeb3React } from '@web3-react/core'
import {
  providers,
  Contract,
} from 'ethers'
import {
  VALID_CHAINS,
  tokenAbi,
  stakingPoolAddresses,
  VALID_CHAIN_IDS,
} from '@src/config'
import abiJSON from '@src/util/abis.json'
import chains from '@src/util/chain_id'

import { message } from 'antd'

import * as contractActions from '@src/redux/modules/contract'
import * as walletActions from '@src/redux/modules/wallet'

import { useLocalStorage } from './useLocalStorage'
import { connectorLocalStorageKey, walletLocalStorageKey } from '@src/types/Connector'

const validChains = VALID_CHAINS;

function checkMetaMask() {
  return window && window.ethereum && window.ethereum.isMetaMask;
}

function normalizeChainId(chainId: number | string | undefined): number {
  if (!chainId) {
    return 0;
  }
  if (typeof chainId === 'string') {
    return chainId.startsWith('0x') ? parseInt(chainId, 16) : parseInt(chainId, 10);
  }
  return chainId;
}

export function listenToWallet() {
  const dispatch = useAppDispatch();
  let walletAddress = useAppSelector(state => state && state.contract.walletAddress);
  let signer = useAppSelector(state => state && state.contract.signer);
  let chain = useAppSelector(state => state && state.wallet.chain);

  const { chainId, account, connector } = useWeb3React()

  const { getLocal } = useLocalStorage();

  useEffect(function mount() {
    // auto connect to wallet iff user has connected
    if (~~getLocal(connectorLocalStorageKey)) {
      // getAccount({hideError: true});
      // window.ethereum && window.ethereum.on('connect', (connectInfo: any) => {
      //   getAccount({hideError: true});
      // });

    }
    // getNetwork();
    window.ethereum && window.ethereum.on('accountsChanged', handleAccountsChanged);
    window.ethereum && window.ethereum.on('chainChanged', handleChainChanged);
    // window.ethereum && window.ethereum.on('disconnect', disconnect);
    window.ethereum && window.ethereum.on('error', (res) => message.error('Transaction Error'));
  }, []);

  useEffect(() => {
    handleAccountsChanged(account)
  }, [account]);

  useEffect(() => {
    handleChainChanged(chainId);
  }, [chainId]);

  useEffect(() => {
    // if (!walletAddress || !chain) {
    if (!walletAddress || !chain) {
      dispatch(contractActions.setSigner(null));
      return;
    }
    if (typeof window !== "undefined" && connector) {
      (async () => {
        const _provider = new providers.Web3Provider(await connector.getProvider());
        // Client-side-only code
        const _signer = _provider.getSigner();
        if (_signer) {
          dispatch(contractActions.setSigner(_signer));
        } else {
          // do nothing
          dispatch(contractActions.setSigner(null));
        }
      })();
    }
  }, [connector, walletAddress, chain])

  useEffect(() => {
    if (!chain) {
      clearContracts();
      return;
    }

    if (!signer) {
      clearContracts();
      return;
    } else {
      const depositTokenAddress = stakingPoolAddresses.find(v => v.chainId == chain.chainId)?.depositTokenAddress;
      if (!depositTokenAddress) {
        clearContracts();
        return;
      }
      const depositTokenContract = new Contract(depositTokenAddress, tokenAbi, signer);
      const earnedTokenAddress = stakingPoolAddresses.find(v => v.chainId == chain.chainId)?.earnedTokenAddress;
      if (!earnedTokenAddress) {
        clearContracts();
        return;
      }

      const breContract = new Contract(earnedTokenAddress, tokenAbi, signer);

      dispatch(contractActions.setContracts({
        depositTokenContract,
        breContract,
      }))
    }
  }, [signer, chain]);


  function clearContracts() {
    dispatch(contractActions.setContracts({
      depositTokenContract: null,
      breContract: null,
      stakingContract: null,
      saleContract: null,
    }))
  }

  async function handleAccountsChanged(account: string) {
    if (!account) {
      // MetaMask is locked or the user has not connected any accounts
      // message.warning('Please connect to your wallet.');
      dispatch(contractActions.setWalletAddress(null));
    } else if (account !== walletAddress) {
      // Do any other work!
      dispatch(contractActions.setWalletAddress(account));
      message.success({
        content: 'You have connected to account ' + account,
        duration: 1
      });
    }
  }

  async function handleChainChanged(chainId: number | string) {
    const normalizedChainId = normalizeChainId(chainId);
    const chain = validChains.find(v => v.chainId == normalizedChainId);
    if (!VALID_CHAIN_IDS.includes(normalizedChainId) || !chain) {
      if (normalizedChainId) {
        message.error('Please connect to Sepolia or local Hardhat network.');
      }
      clearContracts();
      dispatch(walletActions.setChain(null));
      return;
    }
    dispatch(walletActions.setChain(chain));
  }

  return;
}

export const useWallet = () => {

  const dispatch = useAppDispatch();
  const loading = useAppSelector(state => state && state.contract.loading);
  const signer = useAppSelector(state => state && state.contract.signer);
  const walletShowed = useAppSelector(state => state.wallet.show);
  const chain = useAppSelector(state => state.wallet.chain);
  const isWalletInstalled = useAppSelector(state => state.wallet.isWalletInstalled);
  const walletAddress = useAppSelector(state => state.contract.walletAddress);
  const depositTokenContract = useAppSelector(state => state.contract.depositTokenContract);
  const breContract = useAppSelector(state => state.contract.breContract);

  const saleContract = useAppSelector(state => state.contract.saleContract)
  const [saleAddress, setSaleAddress] = useState('');

  const { setLocal } = useLocalStorage();
  /**
   * Init saleContract when saleAddress/userWalletAddress changes
   */
  useEffect(() => {
    if (!signer || !saleAddress) {
      return;
    }
    const saleContract = new Contract(saleAddress, abiJSON['hardhat']['LaunchPadSale'], signer);
    dispatch(contractActions.setContracts({
      saleContract,
    }));
  }, [saleAddress, signer])


  async function getAccount(options?) {
    const hideError = options && options.hideError;
    const showError = !hideError;
    if (!checkMetaMask()) {
      showError && message.error({
        content: 'Please install metamask!',
      })
      dispatch(walletActions.setisWalletInstalled(false));
      return Promise.reject();
    } else {
      dispatch(walletActions.setisWalletInstalled(true));
    }
    if (!window.ethereum.isConnected()) {
      showError && message.error({
        content: 'Please connect to metamask!',
      })
      return Promise.reject();
    }
    try {

      const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' })
      dispatch(contractActions.setWalletAddress(accounts[0]));

      return Promise.resolve(accounts[0]);
    } catch (e) {
      console.error('connect error', e)
    }
  }

  async function getNetwork() {
    if (!checkMetaMask()) {
      message.error({
        content: 'Please install metamask!',
      })
      dispatch(walletActions.setisWalletInstalled(false));
      return Promise.reject();
    } else {
      dispatch(walletActions.setisWalletInstalled(true));
    }
    try {
      const chainId = await window.ethereum.request({ method: 'eth_chainId' });
      const normalizedChainId = normalizeChainId(chainId);
      const chain = validChains.find(v => v.chainId == normalizedChainId);
      if (!chain || !VALID_CHAIN_IDS.includes(normalizedChainId)) {
        dispatch(walletActions.setChain(null));
        return Promise.reject();
      }
      dispatch(walletActions.setChain(chain));
    } catch (e) {
      console.error(e);
    }
    return Promise.resolve(chain);
  }

  function setLoading(data) {
    dispatch(contractActions.setLoading(data));
  }

  async function addToken(tokenAddress, symbolName) {
    await window.ethereum.request({
      method: "wallet_watchAsset",
      params: {
        type: "ERC20", // Initially only supports ERC20, but eventually more!
        options: {
          address: tokenAddress, // The address that the token is at.
          symbol: symbolName, // A ticker symbol or shorthand, up to 5 chars.
          decimals: 18 // The number of decimals in the token
        }
      }
    });
  }

  function showWallet(value?) {
    if (value === undefined) {
      dispatch(walletActions.showWallet(!walletShowed));
      return;
    }
    dispatch(walletActions.showWallet(!!value));
    return;
  }

  async function switchNetwork(chainId: number) {
    const chain = chains.find(chain => chain.chainId == chainId) || validChains[0];
    const params = [
      {
        chainId: `0x${chain.chainId.toString(16)}`,
        chainName: chain.name,
        nativeCurrency: {
          name: chain.nativeCurrency.name,
          symbol: chain.nativeCurrency.symbol,
          decimals: chain.nativeCurrency.decimals,
        },
        rpcUrls: chain.rpc,
        blockExplorerUrls: [`${chain.infoURL}/`],
      },
    ];
    try {
      await window.ethereum.request({
        method: 'wallet_switchEthereumChain',
        params: [{ chainId: '0x' + (chain.chainId).toString(16) }],
      });
      await getNetwork();
    } catch (switchError) {
      console.error('switch network error', switchError)
      // This error code indicates that the chain has not been added to MetaMask.
      if (switchError.code === 4902) {
        try {
          await window.ethereum.request({
            method: 'wallet_addEthereumChain',
            params: params,
          });
        } catch (addError) {
          console.error('add chain error', addError)
          // handle "add" error
        }
      }
      // handle other "switch" errors
    }
  }


  async function connect() {
    await getAccount()
      .then(getNetwork)
      .then((chain) => {
        message.success({
          content: 'Welcome, you\'re currently connected to metamask',
          duration: 1
        });
        // auto connect next time
        setLocal('auto_connect_wallet', 1);
      });
    return Promise.resolve();
  }

  function disconnect() {
    dispatch(contractActions.setWalletAddress(null));
    dispatch(walletActions.setChain(null));

    setLocal('auto_connect_wallet', 0);
    return Promise.resolve();
  }


  return {
    isWalletInstalled,
    walletAddress,
    depositTokenContract,
    breContract,
    saleContract,
    loading,
    saleAddress,
    validChains,
    signer,

    setSaleAddress,
    getAccount,
    setLoading,
    addToken,
    showWallet,
    walletShowed,
    chain,
    switchNetwork,
    connect,
    disconnect,

    isConnected: walletAddress && chain,
  }
}
