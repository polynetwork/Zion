pragma solidity >=0.7.0 <0.9.0;

/**
 * @dev Interface of the SideChainManager contract
 */

interface ISideChainManager {
    event RegisterSideChain(uint64 ChainId, uint64 Router, string Name);
    event ApproveRegisterSideChain(uint64 ChainId);
    event UpdateSideChain(uint64 ChainId, uint64 Router, string Name);
    event ApproveUpdateSideChain(uint64 ChainId);
    event QuitSideChain(uint64 ChainId);
    event ApproveQuitSideChain(uint64 ChainId);

    struct SideChain {
        address owner;
        uint64 chainID;
        uint64 router;
        string name;
        bytes CCMCAddress;
        bytes extraInfo;
    }

    function getSideChain(uint64 chainID) external view returns(SideChain memory sidechain);
    
    function registerSideChain(uint64 chainID, uint64 router, string calldata name, bytes calldata CCMCAddress, bytes calldata extraInfo) external;
    
    function approveRegisterSideChain(uint64 chainID) external returns (bool success);
    
    function updateSideChain(uint64 chainID, uint64 router, string calldata name, bytes calldata CCMCAddress, bytes calldata extraInfo) external;
    
    function approveUpdateSideChain(uint64 chainID) external returns (bool success);
    
    function quitSideChain(uint64 chainID) external;
    
    function approveQuitSideChain(uint64 chainID) external returns (bool success);

    function updateFee(uint64 chainID, uint64 viewNum, int fee) external returns (bool success);

    function registerAsset(uint64 chainID, uint64[] calldata AssetMapKey, bytes[] calldata AssetMapValue, uint64[] calldata LockProxyMapKey, bytes[] calldata LockProxyMapValue) external returns (bool success);

    function getFee(uint64 chainID) external view returns (bytes memory);
}
