pragma solidity >=0.7.0 <0.9.0;

contract neo3_state_manager {
    event evtApproveRegisterStateValidator(uint64 ID);
    event evtApproveRemoveStateValidator(uint64 ID);

    function name() public returns(string memory Name) {
        return Name;
    }
    
    function getCurrentStateValidator() public returns(bytes memory Validator) {
        return Validator;
    }
    
    function registerStateValidator(string[] memory StateValidators, address Address) public returns(bool success) {
        return success;
    }

    function approveRegisterStateValidator(uint64 ID, address Address) public returns(bool success) {
        return success;
    }
    
    function removeStateValidator(string[] memory StateValidators, address Address) public returns(bool success) {
        return success;
    }
    
    function approveRemoveStateValidator(uint64 ID, address Address) public returns(bool success) {
        return success;
    }
}