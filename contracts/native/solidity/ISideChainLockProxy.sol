pragma solidity ^0.5.0;

/**
 * @dev Interface of the LockProxy contract
 */
interface ISideChainLockProxy {
    function name() external view returns (string memory);
    function mint(bytes calldata argsBs, bytes calldata fromContractAddr, uint64 fromChainId) external returns (bool);
    function burn(uint64 toChainId, uint256 amount) external returns (bool);
    function approve(address spender, uint256 amount) external returns (bool);
    function allowance(address owner, address spender) external view returns (uint256);

    event BurnEvent(address fromAssetHash, address fromAddress, uint64 toChainId, bytes toAssetHash, bytes toAddress, uint256 amount);
    event MintEvent(address toAssetHash, address toAddress, uint256 amount);
    event Approval(address indexed owner, address indexed spender, uint256 value);
}
