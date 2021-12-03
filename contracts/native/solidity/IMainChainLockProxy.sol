pragma solidity ^0.5.0;

/**
 * @dev Interface of the LockProxy contract
 */
interface IMainChainLockProxy {
    function name() external view returns (string memory);
    function getSideChainLockAmount(uint64 chainId) external view returns (uint256);
    function lock(uint64 toChainId, address toAddress, uint256 amount) external payable returns (bool);
    function approve(address spender, uint256 amount) external returns (bool);
    function allowance(address owner, address spender) external view returns (uint256);

    event LockEvent(address fromAssetHash, address fromAddress, uint64 toChainId, bytes toAssetHash, bytes toAddress, uint256 amount);
    event UnlockEvent(address toAssetHash, address toAddress, uint256 amount);
    event CrossChainEvent(address indexed sender, bytes txId, address proxyOrAssetContract, uint64 toChainId, bytes toContract, bytes rawdata);
    event VerifyHeaderAndExecuteTxEvent(uint64 fromChainID, bytes toContract, bytes crossChainTxHash, bytes fromChainTxHash);
    event Approval(address indexed owner, address indexed spender, uint256 value);
}
