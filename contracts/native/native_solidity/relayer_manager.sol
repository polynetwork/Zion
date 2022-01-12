pragma solidity >=0.7.0 <0.9.0;

contract relayer_manager {
    event evtApproveRegisterRelayer(uint64 ID);
    event evtApproveRemoveRelayer(uint64 ID);
    event evtRegisterRelayer(uint64 applyID);
    event evtRemoveRelayer(uint64 removeID);

    function name() public returns(string memory Name) {
        return Name;
    }
    
    function registerRelayer(address[] memory AddressList, address Address) public returns(bool success) {
        return success;
    }
    
    function approveRegisterRelayer(uint64 ID, address Address) public returns(bool success) {
        return success;
    }

    function removeRelayer(address[] memory AddressList, address Address) public returns(bool success) {
        return success;
    }
    
    function approveRemoveRelayer(uint64 ID, address Address) public returns(bool success) {
        return success;
    }
}