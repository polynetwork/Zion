pragma solidity ^0.5.0;

/**
 * @dev Interface of the LockProxy contract
 */
interface IEthCrossChainManager {
    function name() external view returns (string memory);
    function crossChain(uint64 toChainId, bytes calldata toContract, bytes calldata method, bytes calldata txData) external returns (bool);
    function verifyHeaderAndExecuteTx(bytes calldata proof, bytes calldata rawHeader, bytes calldata headerProof, bytes calldata curRawHeader,bytes calldata headerSig) external returns (bool);

    event CrossChainEvent(address indexed sender, bytes txId, address proxyOrAssetContract, uint64 toChainId, bytes toContract, bytes rawdata);
    event VerifyHeaderAndExecuteTxEvent(uint64 fromChainID, bytes toContract, bytes crossChainTxHash, bytes fromChainTxHash);
}
