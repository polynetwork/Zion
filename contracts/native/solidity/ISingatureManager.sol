pragma solidity >=0.7.0 <0.9.0;

/**
 * @dev Interface of the signature manager
 */

interface ISignatureManager {
    function addSignature(address addr, uint256 sideChainID, bytes calldata subject, bytes calldata signature) external returns (bool);

    event AddSignatureQuorumEvent(bytes id, bytes subject, uint256 sideChainID);
}
