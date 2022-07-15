pragma solidity >=0.7.0 <0.9.0;

interface IProposalManager {
    function updateNodeManagerGlobalConfig(string calldata maxCommissionChange, string calldata minInitialStake, uint64 maxDescLength, uint64 blockPerEpoch, uint64 consensusValidatorNum,
        uint64 voterValidatorNum, uint64 expireHeight) external returns(bool success);

    event UpdateNodeManagerGlobalConfig();
}
