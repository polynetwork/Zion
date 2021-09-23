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
    
    event EventRegisterCandidate(string Pubkey);
    event EventUnRegisterCandidate(string Pubkey);
    event EventApproveCandidate(string Pubkey);
    event EventBlackNode(string[] PubkeyList);
    event EventWhiteNode(string Pubkey);
    event EventQuitNode(string Pubkey);
    event EventCommitDpos();
    event EventUpdateConfig(Configuration Config);
    event CheckConsensusSignsEvent(uint64 signs);
    
    function MethodContractName() public returns(string memory Name) {
        return Name;
    }
    
    function MethodInitConfig(uint32 BlockMsgDelay, uint32 HashMsgDelay, uint32 PeerHandshakeTimeout, uint32 MaxBlockChangeView, string memory VrfValue, string memory VrfProof, VBFTPeerInfo memory Peers) public returns(bool success) {
        return success;
    }
    
    function MethodRegisterCandidate(string memory PeerPubkey, address Address) public returns(bool success) {
        return success;
    }

    function MethodUnRegisterCandidate(string memory PeerPubkey, address Address) public returns(bool success) {
        return success;
    }
    
    function MethodQuitNode(string memory PeerPubkey, address Address) public returns(bool success) {
        return success;
    }

    function MethodApproveCandidate(string memory PeerPubkey, address Address) public returns(bool success) {
        return success;
    }
    
    function MethodBlackNode(string[] memory PeerPubkeyList, address Address) public returns(bool success) {
        return success;
    }

    function MethodWhiteNode(string memory PeerPubkey, address Address) public returns(bool success) {
        return success;
    }
    
    function MethodUpdateConfig(UpdateConfigParam memory ConfigParam) public returns(bool success) {
        return success;
    }

    function MethodCommitDpos() public returns(bool success) {
        return success;
    }
}