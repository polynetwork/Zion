pragma solidity >=0.6.0 <0.9.0;

interface IMaasConfig {

    function name() external view returns (string memory);
    function changeOwner(address addr) external returns (bool);
    function getOwner() external view returns (address);

    function blockAccount(address addr, bool doBlock) external returns (bool);
    function isBlocked(address addr) external view returns (bool);
    function getBlacklist() external view returns (string memory);
    
    function enableNodeWhite(bool doEnable) external returns (bool);
    function isNodeWhiteEnabled() external view returns (bool);
    function setNodeWhite(address addr, bool isWhite) external returns (bool);
    function isInNodeWhite(address addr) external view returns (bool);
    function getNodeWhitelist() external view returns (string memory);
    
    function enableGasManage(bool doEnable) external returns (bool);
    function isGasManageEnabled() external view returns (bool);
    function setGasManager(address addr, bool isManager) external returns (bool);
    function isGasManager(address addr) external view returns (bool);
    function getGasManagerList() external view returns (string memory);
    
    event ChangeOwner(address indexed oldOwner, address indexed newOwner);
    event BlockAccount(address indexed addr, bool doBlock);
    event EnableNodeWhite(bool doEnable);
    event SetNodeWhite(address indexed addr, bool isWhite);
    event EnableGasManage(bool doEnable);
    event SetGasManager(address indexed addr, bool isManager);
}