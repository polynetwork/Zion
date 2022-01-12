pragma solidity ^0.5.0;

import "./IMainChainLockProxy.sol";

contract MainChainWrapperTest {

    event WrapperLock(address indexed sender, address toAddress, uint64 toChainId, uint256 amount, uint256 fee);

    address public lockProxy;

    constructor() public {
        lockProxy = 0x7d79D936DA7833c7fe056eB450064f34A327DcA8;
    }

    // user lock, `tx.Value` == `amount` + `fee`, and this value will be transfered to `wrapper` contract first, 
    // in native contract the amount of `amount` will be transfer to `lockProxy` again, and the rest left in 
    // `wrapper` contract as handling fee.
    function lock(uint64 toChainId, address toAddress, uint256 amount, uint256 fee) external payable returns (bool) {
        uint256 total = amount + fee;
        require(msg.value == total, "msg.value should be equal to amount + fee!");

        require(IMainChainLockProxy(lockProxy).lock(toChainId, toAddress, amount), "lock failed");

        emit WrapperLock(msg.sender, toAddress, toChainId, amount, fee);
        return true;
    }

}