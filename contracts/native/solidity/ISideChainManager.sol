pragma solidity >=0.7.0 <0.9.0;

/**
 * @dev Interface of the SideChainManager contract
 */

interface ISideChainManager {
    event RegisterSideChain(uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait);
    event ApproveRegisterSideChain(uint64 ChainId);
    event UpdateSideChain(uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait);
    event ApproveUpdateSideChain(uint64 ChainId);
    event QuitSideChain(uint64 ChainId);
    event ApproveQuitSideChain(uint64 ChainId);
    event RegisterRedeem(string rk, string ContractAddress);

    struct SideChain {
        address owner;
        uint64 chainID;
        uint64 router;
        string name;
        uint64 blocksToWait;
        bytes CCMCAddress;
        bytes extraInfo;
    }

    struct BtcTxParamDetail {
	    uint64 PVersion; 
	    uint64 feeRate;  
	    uint64 minChange;
    }

    function getSideChain(uint64 chainID) external view returns(SideChain memory sidechain);
    
    function registerSideChain(uint64 chainID, uint64 router, string calldata name, uint64 blocksToWait, bytes calldata CCMCAddress, bytes calldata extraInfo) external;
    
    function approveRegisterSideChain(uint64 chainID) external returns (bool success);
    
    function updateSideChain(uint64 chainID, uint64 router, string calldata name, uint64 blocksToWait, bytes calldata CCMCAddress, bytes calldata extraInfo) external;
    
    function approveUpdateSideChain(uint64 chainID) external returns (bool success);
    
    function quitSideChain(uint64 chainID) external;
    
    function approveQuitSideChain(uint64 chainID) external returns (bool success);
    
    function registerRedeem(uint64 redeemChainID, uint64 contractChainID, bytes calldata redeem, uint64 CVersion, bytes calldata contractAddress, bytes[] calldata signs) external returns (bool success);

    function setBtcTxParam(bytes calldata redeem, uint64 redeemChainID, bytes[] calldata sigs, BtcTxParamDetail calldata detail) external returns (bool success);
}
