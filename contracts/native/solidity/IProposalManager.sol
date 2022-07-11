pragma solidity >=0.7.0 <0.9.0;

interface IProposalManager {
    function updateNodeManagerGlobalConfig(string calldata maxCommissionChange, string calldata minInitialStake, uint maxDescLength, uint blockPerEpoch, uint consensusValidatorNum,
        uint voterValidatorNum, uint expireHeight) external returns(bool success);

    event UpdateNodeManagerGlobalConfig();
}
