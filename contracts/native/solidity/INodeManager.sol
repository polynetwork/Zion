pragma solidity =0.5.16;

contract INodeManager {

    function createValidator(string calldata consensusPubkey, address proposalAddress, int commission, int initStake, string calldata desc) external returns(bool success);

    function updateValidator(string calldata consensusPubkey, address proposalAddress, int commission, string calldata desc) external returns(bool success);

    function stake(string calldata consensusPubkey, int amount) external returns(bool success);

    function unStake(string calldata consensusPubkey, int amount) external returns(bool success);

    function withdraw() external returns(bool success);

    function cancelValidator(string calldata consensusPubkey) external returns(bool success);

    function withdrawValidator(string calldata consensusPubkey) external returns(bool success);

    function changeEpoch() external returns(bool success);

    function withdrawStakeRewards(string calldata consensusPubkey) external returns(bool success);

    function withdrawCommission(string calldata consensusPubkey) external returns(bool success);

    function beginBlock() external returns(bool success);
}


//pragma solidity ^0.5.0;
//
//interface INodeManager {
//    function name() external view returns (string memory);
//    function propose(uint64 startHeight, bytes calldata peers) external returns (bool);
//    function vote(uint64 epochID, bytes calldata epochHash) external returns (bool);
//    function epoch() external view returns (bytes memory);
//    function getChangingEpoch() external view returns (bytes memory);
//    function getEpochByID(uint64 epochID) external view returns (bytes memory);
//    function proof(uint64 epochID) external view returns (bytes memory);
//
//    event Proposed(bytes epoch);
//    event Voted(uint64 epochID, bytes epochHash, uint64 votedNumber, uint64 groupSize);
//    event EpochChanged(bytes epoch, bytes nextEpoch);
//    event ConsensusSigned(string method, bytes input, address signer, uint64 size);
//}