pragma solidity >=0.7.0 <0.9.0;

contract cross_chain_manager {

    event btcTxMultiSignEvent(bytes TxHash, bytes MultiSign);
    event btcTxToRelayEvent(uint64 FromChainID, uint64 ChainID, string buf, string FromTxHash, string RedeemKey);
    event makeBtcTxEvent(string rk, string buf, uint64[] amts);
    event makeProof(string merkleValueHex, uint64 BlockHeight, string key);

    function name() public returns(string memory Name) {
        return Name;
    }
    
    function importOuterTransfer(uint64 SourceChainID, uint32 Height, bytes memory Proof, bytes memory Extra, bytes memory Signature) public returns(bool success) {
        return success;
    }

    function checkDone(uint64 chainID, bytes memory crossChainID) public view returns(bool success) {
        return success;
    }
    
    function MultiSign(uint64 ChainID, string memory RedeemKey, bytes memory TxHash, string memory Address, bytes[] memory Signs) public returns(bool success) {
        //why does the type of this address is string?
        return success;
    }
    
    function BlackChain(uint64 ChainID) public returns(bool success) {
        return success;
    }
    
    function WhiteChain(uint64 ChainID) public returns(bool success) {
        return success;
    }
}
