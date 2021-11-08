pragma solidity ^0.5.0;

/**
 * @dev Interface of the LockProxy contract
 */
interface ILockProxy {
    function name() external view returns (string memory);
    function bindProxyHash(uint64 toChainId, bytes calldata targetProxyHash) external returns (bool);
    function bindAssetHash(address fromAssetHash, uint64 toChainId, bytes calldata toAssetHash) external returns (bool);
    function crossChain(uint64 toChainId, bytes calldata toContract, bytes calldata method, bytes calldata txData) external returns (bool);
    function verifyHeaderAndExecuteTx(bytes calldata header, bytes calldata proof) external returns (bool);
    function lock(address fromAssetHash, uint64 toChainId, bytes calldata toAddress, uint256 amount) external payable returns (bool);
    function unlock(bytes calldata argsBs, bytes calldata fromContractAddr, uint64 fromChainId) external returns (bool);

    event BindProxyEvent(uint64 toChainId, bytes targetProxyHash);
    event BindAssetEvent(address fromAssetHash, uint64 toChainId, bytes targetProxyHash, uint initialAmount);
    event UnlockEvent(address toAssetHash, address toAddress, uint256 amount);
    event LockEvent(address fromAssetHash, address fromAddress, uint64 toChainId, bytes toAssetHash, bytes toAddress, uint256 amount);
    event CrossChainEvent(address indexed sender, bytes txId, address proxyOrAssetContract, uint64 toChainId, bytes toContract, string method, bytes rawdata);
    event VerifyHeaderAndExecuteTxEvent(uint64 fromChainID, bytes toContract, bytes crossChainTxHash, bytes fromChainTxHash);
}