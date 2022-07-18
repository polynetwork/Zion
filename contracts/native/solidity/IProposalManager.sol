pragma solidity >=0.7.0 <0.9.0;

interface IProposalManager {
    function propose(uint8 pType, bytes calldata content, int stake) external returns(bool success);
    function setActiveProposal() external returns(bool success);
    function voteActiveProposal(int ID) external returns(bool success);
    function getActiveProposal() external view returns(bytes memory);
    function getProposalList() external view returns(bytes memory);

    event Propose(string caller, uint8 pType, string stake, string content);
    event VoteActiveProposal(string ID, uint8 pType);
}
