//SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../interfaces/IAdmin.sol";
import "./LaunchPadSale.sol";


contract SalesFactory {

    IAdmin public admin;
    address public allocationStaking;

    mapping (address => bool) public isSaleCreatedThroughFactory;

    mapping(address => address) public saleOwnerToSale;
    mapping(address => address) public tokenToSale;

    // Expose so query can be possible only by position as well
    address [] public allSales;

    event SaleDeployed(address saleContract);
    event SaleOwnerAndTokenSetInFactory(address sale, address saleOwner, address saleToken);

    modifier onlyAdmin {
        require(admin.isAdmin(msg.sender), "Only Admin can deploy sales");
        _;
    }

    constructor (address _adminContract, address _allocationStaking)  {
        admin = IAdmin(_adminContract);
        allocationStaking = _allocationStaking;
    }

    // Set allocation staking contract address.
    function setAllocationStaking(address _allocationStaking) public onlyAdmin {
        require(_allocationStaking != address(0));
        allocationStaking = _allocationStaking;
    }


    function deploySale()
    external
    onlyAdmin
    {
        LaunchPadSale sale = new LaunchPadSale(address(admin), allocationStaking);

        isSaleCreatedThroughFactory[address(sale)] = true;
        allSales.push(address(sale));

        emit SaleDeployed(address(sale));
    }

    function setSaleOwnerAndToken(address saleOwner, address saleToken) external {
        require(isSaleCreatedThroughFactory[msg.sender], "setSaleOwnerAndToken: Contract not eligible.");
        require(saleOwner != address(0), "Sale owner can not be zero address.");
        require(saleToken != address(0), "Sale token can not be zero address.");
        require(saleOwnerToSale[saleOwner] == address(0), "Sale owner already set.");
        require(tokenToSale[saleToken] == address(0), "Sale token already set.");

        saleOwnerToSale[saleOwner] = msg.sender;
        tokenToSale[saleToken] = msg.sender;

        emit SaleOwnerAndTokenSetInFactory(msg.sender, saleOwner, saleToken);
    }

    // Function to return number of pools deployed
    function getNumberOfSalesDeployed() external view returns (uint) {
        return allSales.length;
    }

    // Function
    function getLastDeployedSale() external view returns (address) {
        //
        if(allSales.length > 0) {
            return allSales[allSales.length - 1];
        }
        return address(0);
    }


    // Function to get all sales
    function getAllSales(uint startIndex, uint endIndex) external view returns (address[] memory) {
        require(endIndex > startIndex, "Bad input");

        address[] memory sales = new address[](endIndex - startIndex);
        uint index = 0;

        for(uint i = startIndex; i < endIndex; i++) {
            sales[index] = allSales[i];
            index++;
        }

        return sales;
    }

}
