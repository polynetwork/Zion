// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract side_chain_manager {
    struct BtcTxParamDetial {
	    uint64 PVersion; 
	    uint64 FeeRate;  
	    uint64 MinChange;
    }
    
    event evtRegisterSideChain(uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait);
    event evtApproveRegisterSideChain(uint64 ChainId);
    event evtUpdateSideChain(uint64 ChainId, uint64 Router, string Name, uint64 BlocksToWait);
    event evtApproveUpdateSideChain(uint64 ChainId);
    event evtQuitSideChain(uint64 ChainId);
    event evtApproveQuitSideChain(uint64 ChainId);
    event evtRegisterRedeem(string rk, string ContractAddress);
    event evtSetBtcTxParam(string rk, uint64 RedeemChainId, uint64 FeeRate, uint64 MinChange);
    
    function name() public returns(string memory Name) {
        return Name;
    }

    function registerSideChain(address Address, uint64 ChainId, uint64 Router, string memory Name, uint64 BlocksToWait, bytes memory CCMCAddress, bytes memory ExtraInfo)
       public returns (bool success){
	    return success;
    }
    
    function approveRegisterSideChain(uint64 Chainid, address Address) public returns (bool success) {
	    return success;
    }
    
    function updateSideChain(address Address, uint64 ChainId, uint64 Router, string memory Name, uint64 BlocksToWait, bytes memory CCMCAddress, bytes memory ExtraInfo)
       public returns (bool success){
	    return success;
    }
    
    function approveUpdateSideChain(uint64 Chainid, address Address) public returns (bool success) {
	    return success;
    }
    
    function quitSideChain(uint64 Chainid, address Address) public returns (bool success) {
	    return success;
    }
    
    function approveQuitSideChain(uint64 Chainid, address Address) public returns (bool success) {
	    return success;
    }
    
    function registerRedeem(uint64 RedeemChainID, uint64 ContractChainID, bytes memory Redeem, uint64 CVersion, bytes memory ContractAddress, bytes[] memory Signs) public returns (bool success) {
	    return success;
    }
    
    function setBtcTxParam(bytes memory Redeem, uint64 RedeemChainId, bytes[] memory Sigs, BtcTxParamDetial memory Detial) public returns (bool success) {
	    return success;
    }
}
