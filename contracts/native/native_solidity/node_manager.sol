pragma solidity >=0.7.0 <0.9.0;

contract node_manager {
    struct VBFTPeerInfo {
        uint32 Index;
        string PeerPubkey;
        string Address;
    }
    struct UpdateConfigParam  {
	    Configuration Config;
    }
    struct Configuration {
        uint32 BlockMsgDelay;
        uint32 HashMsgDelay;
        uint32 PeerHandshakeTimeout;
        uint32 MaxBlockChangeView;
    }
    
    event evtRegisterCandidate(string Pubkey);
    event evtUnRegisterCandidate(string Pubkey);
    event evtApproveCandidate(string Pubkey);
    event evtBlackNode(string[] PubkeyList);
    event evtWhiteNode(string Pubkey);
    event evtQuitNode(string Pubkey);
    event evtCommitDpos();
    event evtUpdateConfig(Configuration Config);
    event CheckConsensusSignsEvent(uint64 signs);
    
    function name() public returns(string memory Name) {
        return Name;
    }
    
    function initConfig(uint32 BlockMsgDelay, uint32 HashMsgDelay, uint32 PeerHandshakeTimeout, uint32 MaxBlockChangeView, string memory VrfValue, string memory VrfProof, VBFTPeerInfo memory Peers) public returns(bool success) {
        return success;
    }
    
    function registerCandidate(string memory PeerPubkey, address Address) public returns(bool success) {
        return success;
    }

    function unRegisterCandidate(string memory PeerPubkey, address Address) public returns(bool success) {
        return success;
    }
    
    function quitNode(string memory PeerPubkey, address Address) public returns(bool success) {
        return success;
    }

    function approveCandidate(string memory PeerPubkey, address Address) public returns(bool success) {
        return success;
    }
    
    function blackNode(string[] memory PeerPubkeyList, address Address) public returns(bool success) {
        return success;
    }

    function whiteNode(string memory PeerPubkey, address Address) public returns(bool success) {
        return success;
    }
    
    function updateConfig(UpdateConfigParam memory ConfigParam) public returns(bool success) {
        return success;
    }

    function commitDpos() public returns(bool success) {
        return success;
    }
}