// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract side_chain_manager {
    struct BtcTxParamDetial {
	    uint64 PVersion; 
	    uint64 FeeRate;  
	    uint64 MinChange;
    }
    
    event EventRegisterSideChain(uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait);
    event EventApproveRegisterSideChain(uint64 ChainId);
    event EventUpdateSideChain(uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait);
    event EventApproveUpdateSideChain(uint64 ChainId);
    event EventQuitSideChain(uint64 ChainId);
    event EventApproveQuitSideChain(uint64 ChainId);
    event EventRegisterRedeem(string rk, string ContractAddress);
    event EventSetBtcTxParam(string rk, uint64 RedeemChainId, uint64 FeeRate, uint64 MinChange);
    
    function MethodContractName() public returns(string memory Name) {
        return Name;
    }

    function MethodRegisterSideChain(address Address, uint64 ChainId, uint64 Router, string memory Name, uint64 BlocksToWait, bytes memory CCMCAddress, bytes memory ExtraInfo)
       public returns (bool success){
	    return success;
    }
    
    function MethodApproveRegisterSideChain(uint64 Chainid, address Address) public returns (bool success) {
	    return success;
    }
    
    function MethodUpdateSideChain(address Address, uint64 ChainId, uint64 Router, string memory Name, uint64 BlocksToWait, bytes memory CCMCAddress, bytes memory ExtraInfo)
       public returns (bool success){
	    return success;
    }
    
    function MethodApproveUpdateSideChain(uint64 Chainid, address Address) public returns (bool success) {
	    return success;
    }
    
    function MethodQuitSideChain(uint64 Chainid, address Address) public returns (bool success) {
	    return success;
    }
    
    function MethodApproveQuitSideChain(uint64 Chainid, address Address) public returns (bool success) {
	    return success;
    }
    
    function MethodRegisterRedeem(uint64 RedeemChainID, uint64 ContractChainID, bytes memory Redeem, uint64 CVersion, bytes memory ContractAddress, bytes[] memory Signs) public returns (bool success) {
	    return success;
    }
    
    function MethodSetBtcTxParam(bytes memory Redeem, uint64 RedeemChainId, bytes[] memory Sigs, BtcTxParamDetial memory Detial) public returns (bool success) {
	    return success;
    }
}
