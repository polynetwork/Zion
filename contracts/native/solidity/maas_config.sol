pragma solidity >=0.6.0 <0.9.0;

interface IMaasConfig {
    function name() external view returns (string memory);
    function changeOwner(address addr) external returns (bool);
    function getOwner() external view returns (address);

    function blockAccount(address addr, bool doBlock) external returns (bool);
    function isBlocked(address addr) external view returns (bool);
    function getBlacklist() external view returns (string memory);
    
    function enableGasManage(bool doEnable) external returns (bool);
    function isGasManageEnabled() external view returns (bool);
    function setGasManager(address addr, bool isManager) external returns (bool);
    function isGasManager(address addr) external view returns (bool);
    function getGasManagerList() external view returns (string memory);

    function setGasUsers(address[] memory addrs, bool addOrRemove) external returns (bool);
    function isGasUser(address addr) external view returns (bool);
    function getGasUserList() external view returns (string memory);
    
    function setAdmins(address[] memory addrs, bool addOrRemove) external returns (bool);
    function isAdmin(address addr) external view returns (bool);
    function getAdminList() external view returns (string memory);
    
    event ChangeOwner(address indexed oldOwner, address indexed newOwner);
    event BlockAccount(address indexed addr, bool doBlock);
    event EnableGasManage(bool doEnable);
    event SetGasManager(address indexed addr, bool isManager);
    event SetGasUsers(address[] addrs, bool addOrRemove);
    event SetAdmins(address[] addrs, bool addOrRemove);
}
