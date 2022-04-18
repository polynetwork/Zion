pragma solidity >=0.6.0 <0.9.0;

interface IMaasConfig {

    function name() external view returns (string memory);
    function changeOwner(address addr) external returns (bool);
    function getOwner() external view returns (address);
    function blockAccount(address addr, bool doBlock) external returns (bool);
    function isBlocked(address addr) external view returns (bool);
    
    event BlockAccount(address addr, bool doBlock);
}