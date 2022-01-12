pragma solidity >=0.7.0 <0.9.0;

contract header_sync {
    event syncHeader(uint64 chainID, uint64 height, string blockHash, uint256 BlockHeight);
    event OKEpochSwitchInfoEvent(uint64 chainID, string BlockHash, uint64 Height, string NextValidatorsHash, string InfoChainID, uint64 BlockHeight);
    
    function name() public returns(string memory Name) {
        return Name;
    }
    
    function syncGenesisHeader(uint64 ChainID, bytes memory GenesisHeader) public returns(bool success) {
        return success;
    }
    
    function syncBlockHeader(uint64 ChainID, address Address, bytes[] memory Headers) public returns(bool success) {
        return success;
    }

    function syncCrossChainMsg(uint64 ChainID, address Address, bytes[] memory CrossChainMsgs) public returns(bool success) {
        return success;
    }
}
