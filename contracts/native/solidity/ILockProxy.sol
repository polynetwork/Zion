pragma solidity ^0.5.0;

/**
 * @dev Interface of the LockProxy contract
 */
interface ILockProxy {
    function name() external view returns (string memory);
    function bindProxyHash(uint64 toChainId, bytes calldata targetProxyHash) external returns (bool);
    function bindAssetHash(address fromAssetHash, uint64 toChainId, bytes calldata toAssetHash) external returns (bool);
    function lock(address fromAssetHash, uint64 toChainId, bytes calldata toAddress, uint256 amount) external returns (bool);
    function unlock(bytes calldata argsBs, bytes calldata fromContractAddr, uint64 fromChainId) external returns (bool);
}
