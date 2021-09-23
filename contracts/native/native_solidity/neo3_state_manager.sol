pragma solidity >=0.7.0 <0.9.0;

contract neo3_state_manager {
    event EventApproveRegisterStateValidator(uint64 ID);
    event EventApproveRemoveStateValidator(uint64 ID);

    function MethodContractName() public returns(string memory Name) {
        return Name;
    }
    
    function MethodGetCurrentStateValidator() public returns(bytes memory Validator) {
        return Validator;
    }
    
    function MethodRegisterStateValidator(string[] memory StateValidators, address Address) public returns(bool success) {
        return success;
    }

    function MethodApproveRegisterStateValidator(uint64 ID, address Address) public returns(bool success) {
        return success;
    }
    
    function MethodRemoveStateValidator(string[] memory StateValidators, address Address) public returns(bool success) {
        return success;
    }
    
    function MethodApproveRemoveStateValidator(uint64 ID, address Address) public returns(bool success) {
        return success;
    }
}