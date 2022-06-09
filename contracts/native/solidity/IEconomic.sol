pragma solidity ^0.5.0;

interface IEconomic {
    function name() external view returns (string memory);
    function totalSupply() external view returns (uint256);
    function reward() external view returns (bytes memory);
}