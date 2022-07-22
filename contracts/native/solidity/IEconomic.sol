pragma solidity >=0.7.0 <0.9.0;

interface IEconomic {
    function name() external view returns (string memory);
    function totalSupply() external view returns (uint256);
    function reward() external view returns (bytes memory);
}