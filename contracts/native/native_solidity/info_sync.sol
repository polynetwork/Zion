//SPDX-License-Identifier: Unlicense
pragma solidity >=0.7.0 <0.9.0;
pragma experimental ABIEncoderV2;

interface InfoSync {
  event SyncRootInfoEvent(uint64 chainID, uint32 height, uint256 BlockHeight);
  function name() external view returns(string memory);
  function syncRootInfo(uint64 chainID, bytes[] calldata rootInfos, bytes memory signature) external returns(bool);
  function getInfoHeight(uint64 chainID) external view returns(uint32);
  function getInfo(uint64 chainID, uint32 height) external view returns(bytes memory);
}