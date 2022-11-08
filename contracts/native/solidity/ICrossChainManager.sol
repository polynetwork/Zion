pragma solidity >=0.7.0 <0.9.0;

interface ICrossChainManager {

    event makeProof(string merkleValueHex, uint64 BlockHeight, string key);
    event ReplenishEvent(string[] txHashes, uint64 chainID);
    event MultiSign(uint64 fromChainId, uint64 toChainId, string txHash, string payment, uint32 sequence);
    event RippleTx(uint64 fromChainId, uint64 toChainId, string txHash, string txJson, uint32 sequence);

    function name() external view returns(string memory Name);
    
    function importOuterTransfer(uint64 SourceChainID, uint32 Height, bytes memory Proof, bytes memory Extra, bytes memory Signature) external returns(bool success);

    function multiSignRipple(uint64 ToChainId, bytes calldata AssetAddress, uint64 FromChainId, bytes calldata TxHash, string calldata TxJson) external returns(bool success);

    function reconstructRippleTx(uint64 FromChainId, bytes calldata TxHash, uint64 ToChainId) external returns(bool success);
  
    function checkDone(uint64 chainID, bytes memory crossChainID) external view returns(bool success);

    function BlackChain(uint64 ChainID) external returns(bool success);

    function WhiteChain(uint64 ChainID) external returns(bool success);

    function replenish(uint64 chainID, string[] calldata txHashes) external returns(bool success);
}
