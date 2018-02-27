// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// KillCordABI is the input ABI used to generate the binding from.
const KillCordABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"getVersion\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"checkIn\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"kill\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getKey\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"s\",\"type\":\"string\"}],\"name\":\"setPayloadEndpoint\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getPayloadEndpoint\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"k\",\"type\":\"string\"}],\"name\":\"setKey\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getLastCheckIn\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getPublisher\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"p\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// KillCordBin is the compiled bytecode used for deploying new contracts.
const KillCordBin = `6060604052341561000f57600080fd5b604051602080610c7c8339810160405280805190602001909190505033600460006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600560006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506040805190810160405280600581526020017f302e302e31000000000000000000000000000000000000000000000000000000815250600090805190602001906100f892919061017f565b5061011461011a640100000000026104b2176401000000009004565b50610224565b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561017657600080fd5b42600381905550565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106101c057805160ff19168380011785556101ee565b828001600101855582156101ee579182015b828111156101ed5782518255916020019190600101906101d2565b5b5090506101fb91906101ff565b5090565b61022191905b8082111561021d576000816000905550600101610205565b5090565b90565b610a49806102336000396000f3006060604052600436106100a4576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630d8e6e2c146100a9578063183ff0851461013757806341c0e1b51461014c57806382678dd614610161578063893d20e8146101ef5780638f4f106e146102445780639f417563146102a1578063af42d1061461032f578063b1bedda31461038c578063dbf4ab4e146103b5575b600080fd5b34156100b457600080fd5b6100bc61040a565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100fc5780820151818401526020810190506100e1565b50505050905090810190601f1680156101295780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561014257600080fd5b61014a6104b2565b005b341561015757600080fd5b61015f610517565b005b341561016c57600080fd5b6101746105ae565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101b4578082015181840152602081019050610199565b50505050905090810190601f1680156101e15780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34156101fa57600080fd5b610202610656565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b341561024f57600080fd5b61029f600480803590602001908201803590602001908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505091905050610680565b005b34156102ac57600080fd5b6102b461074b565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156102f45780820151818401526020810190506102d9565b50505050905090810190601f1680156103215780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561033a57600080fd5b61038a600480803590602001908201803590602001908080601f016020809104026020016040519081016040528093929190818152602001838380828437820191505050505050919050506107f3565b005b341561039757600080fd5b61039f610930565b6040518082815260200191505060405180910390f35b34156103c057600080fd5b6103c861093a565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610412610964565b60008054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156104a85780601f1061047d576101008083540402835291602001916104a8565b820191906000526020600020905b81548152906001019060200180831161048b57829003601f168201915b5050505050905090565b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561050e57600080fd5b42600381905550565b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561057357600080fd5b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16ff5b6105b6610964565b60018054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561064c5780601f106106215761010080835404028352916020019161064c565b820191906000526020600020905b81548152906001019060200180831161062f57829003601f168201915b5050505050905090565b6000600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6000600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415156106de57600080fd5b6102009050808251111515156106f357600080fd5b60001515600560159054906101000a900460ff16151514151561071557600080fd5b816002908051906020019061072b929190610978565b506001600560156101000a81548160ff0219169083151502179055505050565b610753610964565b60028054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156107e95780601f106107be576101008083540402835291602001916107e9565b820191906000526020600020905b8154815290600101906020018083116107cc57829003601f168201915b5050505050905090565b60008060009050600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141561085557600190505b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614156108b057600190505b600115158115151415156108c357600080fd5b60809150818351111515156108d757600080fd5b60001515600560149054906101000a900460ff1615151415156108f957600080fd5b826001908051906020019061090f929190610978565b506001600560146101000a81548160ff021916908315150217905550505050565b6000600354905090565b6000600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b602060405190810160405280600081525090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106109b957805160ff19168380011785556109e7565b828001600101855582156109e7579182015b828111156109e65782518255916020019190600101906109cb565b5b5090506109f491906109f8565b5090565b610a1a91905b80821115610a165760008160009055506001016109fe565b5090565b905600a165627a7a723058208f11f21e43e60c2f2651683ad54927a9af0a6a6e2fd55ea0597814611475e73d0029`

// DeployKillCord deploys a new Ethereum contract, binding an instance of KillCord to it.
func DeployKillCord(auth *bind.TransactOpts, backend bind.ContractBackend, p common.Address) (common.Address, *types.Transaction, *KillCord, error) {
	parsed, err := abi.JSON(strings.NewReader(KillCordABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(KillCordBin), backend, p)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &KillCord{KillCordCaller: KillCordCaller{contract: contract}, KillCordTransactor: KillCordTransactor{contract: contract}}, nil
}

// KillCord is an auto generated Go binding around an Ethereum contract.
type KillCord struct {
	KillCordCaller     // Read-only binding to the contract
	KillCordTransactor // Write-only binding to the contract
}

// KillCordCaller is an auto generated read-only Go binding around an Ethereum contract.
type KillCordCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KillCordTransactor is an auto generated write-only Go binding around an Ethereum contract.
type KillCordTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KillCordSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KillCordSession struct {
	Contract     *KillCord         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KillCordCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KillCordCallerSession struct {
	Contract *KillCordCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// KillCordTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KillCordTransactorSession struct {
	Contract     *KillCordTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// KillCordRaw is an auto generated low-level Go binding around an Ethereum contract.
type KillCordRaw struct {
	Contract *KillCord // Generic contract binding to access the raw methods on
}

// KillCordCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KillCordCallerRaw struct {
	Contract *KillCordCaller // Generic read-only contract binding to access the raw methods on
}

// KillCordTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KillCordTransactorRaw struct {
	Contract *KillCordTransactor // Generic write-only contract binding to access the raw methods on
}

// NewKillCord creates a new instance of KillCord, bound to a specific deployed contract.
func NewKillCord(address common.Address, backend bind.ContractBackend) (*KillCord, error) {
	contract, err := bindKillCord(address, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KillCord{KillCordCaller: KillCordCaller{contract: contract}, KillCordTransactor: KillCordTransactor{contract: contract}}, nil
}

// NewKillCordCaller creates a new read-only instance of KillCord, bound to a specific deployed contract.
func NewKillCordCaller(address common.Address, caller bind.ContractCaller) (*KillCordCaller, error) {
	contract, err := bindKillCord(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &KillCordCaller{contract: contract}, nil
}

// NewKillCordTransactor creates a new write-only instance of KillCord, bound to a specific deployed contract.
func NewKillCordTransactor(address common.Address, transactor bind.ContractTransactor) (*KillCordTransactor, error) {
	contract, err := bindKillCord(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &KillCordTransactor{contract: contract}, nil
}

// bindKillCord binds a generic wrapper to an already deployed contract.
func bindKillCord(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KillCordABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KillCord *KillCordRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _KillCord.Contract.KillCordCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KillCord *KillCordRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KillCord.Contract.KillCordTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KillCord *KillCordRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KillCord.Contract.KillCordTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KillCord *KillCordCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _KillCord.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KillCord *KillCordTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KillCord.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KillCord *KillCordTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KillCord.Contract.contract.Transact(opts, method, params...)
}

// GetKey is a free data retrieval call binding the contract method 0x82678dd6.
//
// Solidity: function getKey() constant returns(string)
func (_KillCord *KillCordCaller) GetKey(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _KillCord.contract.Call(opts, out, "getKey")
	return *ret0, err
}

// GetKey is a free data retrieval call binding the contract method 0x82678dd6.
//
// Solidity: function getKey() constant returns(string)
func (_KillCord *KillCordSession) GetKey() (string, error) {
	return _KillCord.Contract.GetKey(&_KillCord.CallOpts)
}

// GetKey is a free data retrieval call binding the contract method 0x82678dd6.
//
// Solidity: function getKey() constant returns(string)
func (_KillCord *KillCordCallerSession) GetKey() (string, error) {
	return _KillCord.Contract.GetKey(&_KillCord.CallOpts)
}

// GetLastCheckIn is a free data retrieval call binding the contract method 0xb1bedda3.
//
// Solidity: function getLastCheckIn() constant returns(uint256)
func (_KillCord *KillCordCaller) GetLastCheckIn(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _KillCord.contract.Call(opts, out, "getLastCheckIn")
	return *ret0, err
}

// GetLastCheckIn is a free data retrieval call binding the contract method 0xb1bedda3.
//
// Solidity: function getLastCheckIn() constant returns(uint256)
func (_KillCord *KillCordSession) GetLastCheckIn() (*big.Int, error) {
	return _KillCord.Contract.GetLastCheckIn(&_KillCord.CallOpts)
}

// GetLastCheckIn is a free data retrieval call binding the contract method 0xb1bedda3.
//
// Solidity: function getLastCheckIn() constant returns(uint256)
func (_KillCord *KillCordCallerSession) GetLastCheckIn() (*big.Int, error) {
	return _KillCord.Contract.GetLastCheckIn(&_KillCord.CallOpts)
}

// GetOwner is a free data retrieval call binding the contract method 0x893d20e8.
//
// Solidity: function getOwner() constant returns(address)
func (_KillCord *KillCordCaller) GetOwner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KillCord.contract.Call(opts, out, "getOwner")
	return *ret0, err
}

// GetOwner is a free data retrieval call binding the contract method 0x893d20e8.
//
// Solidity: function getOwner() constant returns(address)
func (_KillCord *KillCordSession) GetOwner() (common.Address, error) {
	return _KillCord.Contract.GetOwner(&_KillCord.CallOpts)
}

// GetOwner is a free data retrieval call binding the contract method 0x893d20e8.
//
// Solidity: function getOwner() constant returns(address)
func (_KillCord *KillCordCallerSession) GetOwner() (common.Address, error) {
	return _KillCord.Contract.GetOwner(&_KillCord.CallOpts)
}

// GetPayloadEndpoint is a free data retrieval call binding the contract method 0x9f417563.
//
// Solidity: function getPayloadEndpoint() constant returns(string)
func (_KillCord *KillCordCaller) GetPayloadEndpoint(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _KillCord.contract.Call(opts, out, "getPayloadEndpoint")
	return *ret0, err
}

// GetPayloadEndpoint is a free data retrieval call binding the contract method 0x9f417563.
//
// Solidity: function getPayloadEndpoint() constant returns(string)
func (_KillCord *KillCordSession) GetPayloadEndpoint() (string, error) {
	return _KillCord.Contract.GetPayloadEndpoint(&_KillCord.CallOpts)
}

// GetPayloadEndpoint is a free data retrieval call binding the contract method 0x9f417563.
//
// Solidity: function getPayloadEndpoint() constant returns(string)
func (_KillCord *KillCordCallerSession) GetPayloadEndpoint() (string, error) {
	return _KillCord.Contract.GetPayloadEndpoint(&_KillCord.CallOpts)
}

// GetPublisher is a free data retrieval call binding the contract method 0xdbf4ab4e.
//
// Solidity: function getPublisher() constant returns(address)
func (_KillCord *KillCordCaller) GetPublisher(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _KillCord.contract.Call(opts, out, "getPublisher")
	return *ret0, err
}

// GetPublisher is a free data retrieval call binding the contract method 0xdbf4ab4e.
//
// Solidity: function getPublisher() constant returns(address)
func (_KillCord *KillCordSession) GetPublisher() (common.Address, error) {
	return _KillCord.Contract.GetPublisher(&_KillCord.CallOpts)
}

// GetPublisher is a free data retrieval call binding the contract method 0xdbf4ab4e.
//
// Solidity: function getPublisher() constant returns(address)
func (_KillCord *KillCordCallerSession) GetPublisher() (common.Address, error) {
	return _KillCord.Contract.GetPublisher(&_KillCord.CallOpts)
}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() constant returns(string)
func (_KillCord *KillCordCaller) GetVersion(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _KillCord.contract.Call(opts, out, "getVersion")
	return *ret0, err
}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() constant returns(string)
func (_KillCord *KillCordSession) GetVersion() (string, error) {
	return _KillCord.Contract.GetVersion(&_KillCord.CallOpts)
}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() constant returns(string)
func (_KillCord *KillCordCallerSession) GetVersion() (string, error) {
	return _KillCord.Contract.GetVersion(&_KillCord.CallOpts)
}

// CheckIn is a paid mutator transaction binding the contract method 0x183ff085.
//
// Solidity: function checkIn() returns()
func (_KillCord *KillCordTransactor) CheckIn(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KillCord.contract.Transact(opts, "checkIn")
}

// CheckIn is a paid mutator transaction binding the contract method 0x183ff085.
//
// Solidity: function checkIn() returns()
func (_KillCord *KillCordSession) CheckIn() (*types.Transaction, error) {
	return _KillCord.Contract.CheckIn(&_KillCord.TransactOpts)
}

// CheckIn is a paid mutator transaction binding the contract method 0x183ff085.
//
// Solidity: function checkIn() returns()
func (_KillCord *KillCordTransactorSession) CheckIn() (*types.Transaction, error) {
	return _KillCord.Contract.CheckIn(&_KillCord.TransactOpts)
}

// Kill is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: function kill() returns()
func (_KillCord *KillCordTransactor) Kill(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KillCord.contract.Transact(opts, "kill")
}

// Kill is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: function kill() returns()
func (_KillCord *KillCordSession) Kill() (*types.Transaction, error) {
	return _KillCord.Contract.Kill(&_KillCord.TransactOpts)
}

// Kill is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: function kill() returns()
func (_KillCord *KillCordTransactorSession) Kill() (*types.Transaction, error) {
	return _KillCord.Contract.Kill(&_KillCord.TransactOpts)
}

// SetKey is a paid mutator transaction binding the contract method 0xaf42d106.
//
// Solidity: function setKey(k string) returns()
func (_KillCord *KillCordTransactor) SetKey(opts *bind.TransactOpts, k string) (*types.Transaction, error) {
	return _KillCord.contract.Transact(opts, "setKey", k)
}

// SetKey is a paid mutator transaction binding the contract method 0xaf42d106.
//
// Solidity: function setKey(k string) returns()
func (_KillCord *KillCordSession) SetKey(k string) (*types.Transaction, error) {
	return _KillCord.Contract.SetKey(&_KillCord.TransactOpts, k)
}

// SetKey is a paid mutator transaction binding the contract method 0xaf42d106.
//
// Solidity: function setKey(k string) returns()
func (_KillCord *KillCordTransactorSession) SetKey(k string) (*types.Transaction, error) {
	return _KillCord.Contract.SetKey(&_KillCord.TransactOpts, k)
}

// SetPayloadEndpoint is a paid mutator transaction binding the contract method 0x8f4f106e.
//
// Solidity: function setPayloadEndpoint(s string) returns()
func (_KillCord *KillCordTransactor) SetPayloadEndpoint(opts *bind.TransactOpts, s string) (*types.Transaction, error) {
	return _KillCord.contract.Transact(opts, "setPayloadEndpoint", s)
}

// SetPayloadEndpoint is a paid mutator transaction binding the contract method 0x8f4f106e.
//
// Solidity: function setPayloadEndpoint(s string) returns()
func (_KillCord *KillCordSession) SetPayloadEndpoint(s string) (*types.Transaction, error) {
	return _KillCord.Contract.SetPayloadEndpoint(&_KillCord.TransactOpts, s)
}

// SetPayloadEndpoint is a paid mutator transaction binding the contract method 0x8f4f106e.
//
// Solidity: function setPayloadEndpoint(s string) returns()
func (_KillCord *KillCordTransactorSession) SetPayloadEndpoint(s string) (*types.Transaction, error) {
	return _KillCord.Contract.SetPayloadEndpoint(&_KillCord.TransactOpts, s)
}
