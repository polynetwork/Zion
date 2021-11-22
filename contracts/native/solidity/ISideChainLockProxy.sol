pragma solidity ^0.5.0;

/**
 * @dev Interface of the LockProxy contract
 */
interface ISideChainLockProxy {
    function name() external view returns (string memory);
    // function initGenesisHeader(bytes calldata header, bytes calldata proof, bytes calldata extra, bytes calldata epoch) external returns (bool);
    // function changeEpoch(bytes calldata header, bytes calldata proof, bytes calldata extra, bytes calldata epoch) external returns (bool);
    function burn(uint64 toChainId, address toAddress, uint256 amount) external returns (bool);
    function verifyHeaderAndExecuteTx(bytes calldata header, bytes calldata rawCrossTx, bytes calldata proof, bytes calldata extra) external returns (bool);

    event InitGenesisBlockEvent(uint256 height, bytes header, bytes epoch);
    event ChangeEpochEvent(uint256 height, bytes header, bytes oldEpoch, bytes newEpoch);
    event BurnEvent(uint64 toChainId, address fromAddress, address toAddress, uint256 amount, bytes crossTxId);
    event MintEvent(uint64 toChainId, address fromAddress, address toAddress, uint256 amount);
}
