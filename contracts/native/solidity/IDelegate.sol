pragma solidity ^0.5.0;

interface IDelegate {
    function approve(address spender, uint256 amount) external returns (bool);
    function allowance(address owner, address spender) external view returns (uint256);

    event Approval(address indexed owner, address indexed spender, uint256 value);
}