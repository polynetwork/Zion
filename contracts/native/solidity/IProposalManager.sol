pragma solidity >=0.7.0 <0.9.0;

interface IProposalManager {
    function propose(bytes calldata content) external returns(bool success);
    function proposeConfig(bytes calldata content) external returns(bool success);
    function voteProposal(int ID) external returns(bool success);
    function getProposal(int ID) external view returns(bytes memory);
    function getProposalList() external view returns(bytes memory);
    function getConfigProposalList() external view returns(bytes memory);

    event Propose(string ID, string caller, string stake, string content);
    event ProposeConfig(string ID, string caller, string stake, string content);
    event VoteProposal(string ID);
}
