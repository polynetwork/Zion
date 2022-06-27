pragma solidity =0.5.16;

contract INodeManager {
    function createValidator(string calldata consensusPubkey, address proposalAddress, int commission, int initStake, string calldata desc) external returns(bool success);
    function updateValidator(string calldata consensusPubkey, address proposalAddress, string calldata desc) external returns(bool success);
    function updateCommission(string calldata consensusPubkey, int commission) external returns(bool success);
    function stake(string calldata consensusPubkey, int amount) external returns(bool success);
    function unStake(string calldata consensusPubkey, int amount) external returns(bool success);
    function withdraw() external returns(bool success);
    function cancelValidator(string calldata consensusPubkey) external returns(bool success);
    function withdrawValidator(string calldata consensusPubkey) external returns(bool success);
    function changeEpoch() external returns(bool success);
    function withdrawStakeRewards(string calldata consensusPubkey) external returns(bool success);
    function withdrawCommission(string calldata consensusPubkey) external returns(bool success);
    function endBlock() external returns(bool success);
    function getGlobalConfig() external view returns (bytes memory);
    function getCommunityInfo() external view returns (bytes memory);
    function getCurrentEpochInfo() external view returns (bytes memory);

    event CreateValidator(string consensusPubkey);
    event UpdateValidator(string consensusPubkey);
    event UpdateCommission(string consensusPubkey);
    event Stake(string consensusPubkey, string amount);
    event UnStake(string consensusPubkey, string amount);
    event Withdraw(string caller, string amount);
    event CancelValidator(string consensusPubkey);
    event WithdrawValidator(string consensusPubkey, string selfStake);
    event ChangeEpoch(string epochID);
    event WithdrawStakeRewards(string rewards);
    event WithdrawCommission(string consensusPubkey, string commission);
}