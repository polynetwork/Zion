pragma solidity ^0.5.0;

import "./ISideChainLockProxy.sol";

contract SideChainWrapperTest {

    address public lockProxy;
    uint64 public toChainId;

    event WrapperBurn(address indexed sender, uint256 amount, uint256 fee);

    constructor() public {
        lockProxy = 0x7d79D936DA7833c7fe056eB450064f34A327DcA8;
        toChainId = 1;
    }

    // user burn, the `tx.value` will be transferred to `wrapper` contract first,
    // and in native contract call, the amount of `amount` will be burned, and the 
    // rest of `tx.value` will be left as handling fee.
    function burn(uint256 amount, uint256 fee) external payable returns (bool) {
        uint256 total = amount + fee;
        require(msg.value == total, "tx.value should be equal to sum of amount and fee");
        
        require(ISideChainLockProxy(lockProxy).burn(toChainId, amount), "lockProxy.burn failed");

        emit WrapperBurn(msg.sender, amount, fee);
        return true;
    }
}