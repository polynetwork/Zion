pragma solidity ^0.5.0;

/**
 * @dev Interface of the LockProxy contract
 */
interface IEthCrossChainManager {
    function crossChain(uint64 toChainId, bytes calldata toContract, bytes calldata method, bytes calldata txData) external returns (bool);
    function verifyHeaderAndExecuteTx(bytes calldata proof, bytes calldata rawHeader, bytes calldata headerProof, bytes calldata curRawHeader,bytes calldata headerSig) external returns (bool);
}
