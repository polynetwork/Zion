pragma solidity >=0.7.0 <0.9.0;

interface INodeManager {
    function createValidator(address consensusAddress, address signerAddress, address proposalAddress, int commission, int initStake, string calldata desc) external returns(bool success);
    function updateValidator(address consensusAddress, address signerAddress, address proposalAddress, string calldata desc) external returns(bool success);
    function updateCommission(address consensusAddress, int commission) external returns(bool success);
    function stake(address consensusAddress, int amount) external returns(bool success);
    function unStake(address consensusAddress, int amount) external returns(bool success);
    function withdraw() external returns(bool success);
    function cancelValidator(address consensusAddress) external returns(bool success);
    function withdrawValidator(address consensusAddress) external returns(bool success);
    function changeEpoch() external returns(bool success);
    function withdrawStakeRewards(address consensusAddress) external returns(bool success);
    function withdrawCommission(address consensusAddress) external returns(bool success);
    function endBlock() external returns(bool success);
    function getGlobalConfig() external view returns (bytes memory);
    function getCommunityInfo() external view returns (bytes memory);
    function getCurrentEpochInfo() external view returns (bytes memory);
    function getEpochInfo(int id) external view returns (bytes memory);
    function getAllValidators() external view returns (bytes memory);
    function getValidator(address consensusAddress) external view returns (bytes memory);
    function getStakeInfo(address consensusAddress, address stakeAddress) external view returns (bytes memory);
    function getUnlockingInfo(address stakeAddress) external view returns (bytes memory);
    function getStakeStartingInfo(address consensusAddress, address stakeAddress) external view returns (bytes memory);
    function getAccumulatedCommission(address consensusAddress) external view returns (bytes memory);
    function getValidatorSnapshotRewards(address consensusAddress, uint64 period) external view returns (bytes memory);
    function getValidatorAccumulatedRewards(address consensusAddress) external view returns (bytes memory);
    function getValidatorOutstandingRewards(address consensusAddress) external view returns (bytes memory);
    function getTotalPool() external view returns (bytes memory);
    function getOutstandingRewards() external view returns (bytes memory);

    event CreateValidator(string consensusAddress);
    event UpdateValidator(string consensusAddress);
    event UpdateCommission(string consensusAddress);
    event Stake(string consensusAddress, string amount);
    event UnStake(string consensusAddress, string amount);
    event Withdraw(string caller, string amount);
    event CancelValidator(string consensusAddress);
    event WithdrawValidator(string consensusAddress, string selfStake);
    event ChangeEpoch(string epochID);
    event WithdrawStakeRewards(string consensusAddress, string caller, string rewards);
    event WithdrawCommission(string consensusAddress, string commission);
}