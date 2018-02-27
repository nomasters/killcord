pragma solidity 0.4.19;

contract killcord {
  string version;
  string publishedKey;
  string payloadEndpoint;
  uint lastCheckIn;
  address owner;
  address publisher;
  bool lockPublishedKey;
  bool lockPayloadEndpoint;

  // set the `owner` of the contract and log first `checkIn`
  function killcord(address p) public {
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
  function getLastCheckIn() public constant returns (uint) {
    return lastCheckIn;
  }

  // Outputs the `string` for the last `block.timestamp`
  // that registered to this contract on the blockchain.
  function getPayloadEndpoint() public constant returns (string) {
    return payloadEndpoint;
  }

  // This function is restricted to work with only the contract owner.
  // Sets the Payload Endpoint after checking max length of the string.
  // sets lockPayloadEndpoint to TRUE so that once set, this value can
  // not be changed.
  function setPayloadEndpoint(string s) public onlyOwner {
    uint max = 512;
    require(bytes(s).length <= max);
    require(lockPayloadEndpoint == false);
    payloadEndpoint = s;
    lockPayloadEndpoint = true;
  }

  // getKey() simply outputs the `publishedKey` saved to the blockChain
  function getKey() public constant returns (string) {
    return publishedKey;
  }

  function getOwner() public constant returns (address) {
    return owner;
  }

  function getPublisher() public constant returns (address) {
    return publisher;
  }

  function getVersion() public constant returns (string) {
    return version;
  }

  // This function is restricted to work with only the contract owner.
  function setKey(string k) public onlyOwnerOrPublisher {
    uint max = 128;
    require(bytes(k).length <= max);
    require(lockPublishedKey == false);
    publishedKey = k;
    lockPublishedKey = true;
  }
}