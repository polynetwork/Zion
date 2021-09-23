pragma solidity >=0.7.0 <0.9.0;

contract cross_chain_manager {

    event btcTxMultiSignEvent(bytes TxHash, bytes MultiSign);
    event btcTxToRelayEvent(uint64 FromChainID, uint64 ChainID, string buf, string FromTxHash, string RedeemKey);
    event makeBtcTxEvent(string rk, string buf, uint64[] amts);
    event NOTIFY_MAKE_PROOF_EVENT(string merkleValueHex, uint64 BlockHeight, string key);

    function MethodContractName() public returns(string memory Name) {
        return Name;
    }

    function MethodImportOuterTransfer(uint64 SourceChainID, uint32 Height, bytes memory Proof, bytes memory RelayerAddress, bytes memory Extra, bytes memory HeaderOrCrossChainMsg) public returns(bool success) {
        return success;
    }

    function MethodMultiSign(uint64 ChainID, string memory RedeemKey, bytes memory TxHash, string memory Address, bytes[] memory Signs) public returns(bool success) {
        return success;
    }

    function MethodBlackChain(uint64 ChainID) public returns(bool success) {
        return success;
    }

    function MethodWhiteChain(uint64 ChainID) public returns(bool success) {
        return success;
    }
}
