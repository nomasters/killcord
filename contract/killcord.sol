// SPDX-License-Identifier: Unlicense
pragma solidity ^0.7.0;

contract killcord {
  string version;
  string publishedKey;
  string payloadEndpoint;
  uint lastCheckIn;
  address payable owner;
  address publisher;
  bool lockPublishedKey;
  bool lockPayloadEndpoint;

  // set the `owner` of the contract and log first `checkIn`
  constructor(address p) {
    owner = msg.sender;
    publisher = p;
    version = "0.0.1";
    checkIn();
  }

  // a function modifier used to restrict most write functions to only
  // the contract owner
  modifier onlyOwner {
    require(msg.sender == owner);
    _;
  }

  // a function modifier used to restrict publishing the key to only
  // the owner or publisher addresses
  modifier onlyOwnerOrPublisher {
    bool ok = false;
    if (msg.sender == publisher) {
      ok = true;
    }
    if (msg.sender == owner) {
      ok = true;
    }
    require(ok == true);
    _;
  }

  // This function is restricted to work with only the contract owner.
  // friends don't let friends deploy contracts that can't be killed
  function kill() public onlyOwner {
    selfdestruct(owner);
  }

  // This function is restricted to work with only the contract owner.
  // `block.timestamp` is known to tolerate datestamp drift of up to
  // 900 seconds at the time of this writing, consider then when
  // setting TTL thresholds for the publisher.
  function checkIn() public onlyOwner {
    lastCheckIn = block.timestamp;
  }

  // Outputs the `uint` for the last `block.timestamp`
  // that registered to this contract on the blockchain.
  function getLastCheckIn() public view returns (uint) {
    return lastCheckIn;
  }

  // Outputs the `string` for the last `block.timestamp`
  // that registered to this contract on the blockchain.
  function getPayloadEndpoint() public view returns (string memory) {
    return payloadEndpoint;
  }

  // This function is restricted to work with only the contract owner.
  // Sets the Payload Endpoint after checking max length of the string.
  // sets lockPayloadEndpoint to TRUE so that once set, this value can
  // not be changed.
  function setPayloadEndpoint(string memory s) public onlyOwner {
    uint max = 512;
    require(bytes(s).length <= max);
    require(lockPayloadEndpoint == false);
    payloadEndpoint = s;
    lockPayloadEndpoint = true;
  }

  // getKey() simply outputs the `publishedKey` saved to the blockChain
  function getKey() public view returns (string memory) {
    return publishedKey;
  }

  function getOwner() public view returns (address) {
    return owner;
  }

  function getPublisher() public view returns (address) {
    return publisher;
  }

  function getVersion() public view returns (string memory) {
    return version;
  }

  // This function is restricted to work with only the contract owner.
  function setKey(string memory k) public onlyOwnerOrPublisher {
    uint max = 128;
    require(bytes(k).length <= max);
    require(lockPublishedKey == false);
    publishedKey = k;
    lockPublishedKey = true;
  }
}