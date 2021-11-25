pragma solidity ^0.5.0;

/**
 * @dev Interface of the LockProxy contract
 */
interface IMainChainLockProxy {
    function name() external view returns (string memory);
    // function bindProxyHash(uint64 toChainId, bytes calldata targetProxyHash) external returns (bool);
    // function getProxyHash(uint64 toChainId) external view returns (bytes memory);
    // function bindAssetHash(address fromAssetHash, uint64 toChainId, bytes calldata toAssetHash) external returns (bool);
    // function getAssetHash(address fromAssetHash, uint64 toChainId) external view returns (bytes memory);
    // function bindCaller(uint64 toChainId, bytes calldata caller) external returns (bool);
    // function getCaller(uint64 toChainId) external view returns (bytes memory);
    function getSideChainLockAmount(uint64 chainId) external view returns (uint);
    function lock(address fromAssetHash, uint64 toChainId, bytes calldata toAddress, uint256 amount) external payable returns (bool);

    // event BindProxyEvent(uint64 toChainId, bytes targetProxyHash);
    // event BindAssetEvent(address fromAssetHash, uint64 toChainId, bytes targetProxyHash, uint initialAmount);
    // event BindCaller(uint64 toChainId, bytes caller);
    event LockEvent(address fromAssetHash, address fromAddress, uint64 toChainId, bytes toAssetHash, bytes toAddress, uint256 amount);
    event UnlockEvent(address toAssetHash, address toAddress, uint256 amount);
    event CrossChainEvent(address indexed sender, bytes txId, address proxyOrAssetContract, uint64 toChainId, bytes toContract, bytes rawdata);
    event VerifyHeaderAndExecuteTxEvent(uint64 fromChainID, bytes toContract, bytes crossChainTxHash, bytes fromChainTxHash);
}