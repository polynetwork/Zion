pragma solidity >=0.7.0 <0.9.0;

contract relayer_manager {
    event EventApproveRegisterRelayer(uint64 ID);
    event EventApproveRemoveRelayer(uint64 ID);
    event EventRegisterRelayer(uint64 applyID);
    event EventRemoveRelayer(uint64 removeID);

    function MethodContractName() public returns(string memory Name) {
        return Name;
    }
    
    function MethodRegisterRelayer(address[] memory AddressList, address Address) public returns(bool success) {
        return success;
    }
    
    function MethodApproveRegisterRelayer(uint64 ID, address Address) public returns(bool success) {
        return success;
    }

    function MethodRemoveRelayer(address[] memory AddressList, address Address) public returns(bool success) {
        return success;
    }
    
    function MethodApproveRemoveRelayer(uint64 ID, address Address) public returns(bool success) {
        return success;
    }
}