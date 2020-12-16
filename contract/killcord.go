// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// KillCordABI is the input ABI used to generate the binding from.
const KillCordABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"p\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"checkIn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getKey\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastCheckIn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPayloadEndpoint\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPublisher\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"kill\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"k\",\"type\":\"string\"}],\"name\":\"setKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"s\",\"type\":\"string\"}],\"name\":\"setPayloadEndpoint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// KillCordBin is the compiled bytecode used for deploying new contracts.
var KillCordBin = "0x608060405234801561001057600080fd5b50604051610c95380380610c958339818101604052602081101561003357600080fd5b810190808051906020019092919050505033600460006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600560006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506040518060400160405280600581526020017f302e302e3100000000000000000000000000000000000000000000000000000081525060009080519060200190610111929190610189565b5061012061012660201b60201c565b50610234565b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461018057600080fd5b42600381905550565b828054600181600116156101000203166002900490600052602060002090601f0160209004810192826101bf5760008555610206565b82601f106101d857805160ff1916838001178555610206565b82800160010185558215610206579182015b828111156102055782518255916020019190600101906101ea565b5b5090506102139190610217565b5090565b5b80821115610230576000816000905550600101610218565b5090565b610a52806102436000396000f3fe608060405234801561001057600080fd5b506004361061009e5760003560e01c80638f4f106e116100665780638f4f106e146101f15780639f417563146102ac578063af42d1061461032f578063b1bedda3146103ea578063dbf4ab4e146104085761009e565b80630d8e6e2c146100a3578063183ff0851461012657806341c0e1b51461013057806382678dd61461013a578063893d20e8146101bd575b600080fd5b6100ab61043c565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100eb5780820151818401526020810190506100d0565b50505050905090810190601f1680156101185780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b61012e6104de565b005b610138610541565b005b6101426105d6565b6040518080602001828103825283818151815260200191508051906020019080838360005b83811015610182578082015181840152602081019050610167565b50505050905090810190601f1680156101af5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6101c5610678565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6102aa6004803603602081101561020757600080fd5b810190808035906020019064010000000081111561022457600080fd5b82018360208201111561023657600080fd5b8035906020019184600183028401116401000000008311171561025857600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505091929192905050506106a2565b005b6102b4610767565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156102f45780820151818401526020810190506102d9565b50505050905090810190601f1680156103215780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6103e86004803603602081101561034557600080fd5b810190808035906020019064010000000081111561036257600080fd5b82018360208201111561037457600080fd5b8035906020019184600183028401116401000000008311171561039657600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610809565b005b6103f261093d565b6040518082815260200191505060405180910390f35b610410610947565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b606060008054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156104d45780601f106104a9576101008083540402835291602001916104d4565b820191906000526020600020905b8154815290600101906020018083116104b757829003601f168201915b5050505050905090565b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461053857600080fd5b42600381905550565b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461059b57600080fd5b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16ff5b606060018054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561066e5780601f106106435761010080835404028352916020019161066e565b820191906000526020600020905b81548152906001019060200180831161065157829003601f168201915b5050505050905090565b6000600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146106fc57600080fd5b60006102009050808251111561071157600080fd5b60001515600560159054906101000a900460ff1615151461073157600080fd5b8160029080519060200190610747929190610971565b506001600560156101000a81548160ff0219169083151502179055505050565b606060028054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156107ff5780601f106107d4576101008083540402835291602001916107ff565b820191906000526020600020905b8154815290600101906020018083116107e257829003601f168201915b5050505050905090565b6000600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141561086657600190505b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614156108c157600190505b60011515811515146108d257600080fd5b60006080905080835111156108e657600080fd5b60001515600560149054906101000a900460ff1615151461090657600080fd5b826001908051906020019061091c929190610971565b506001600560146101000a81548160ff021916908315150217905550505050565b6000600354905090565b6000600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b828054600181600116156101000203166002900490600052602060002090601f0160209004810192826109a757600085556109ee565b82601f106109c057805160ff19168380011785556109ee565b828001600101855582156109ee579182015b828111156109ed5782518255916020019190600101906109d2565b5b5090506109fb91906109ff565b5090565b5b80821115610a18576000816000905550600101610a00565b509056fea264697066735822122022cdd092195489e25280d8a632a28378c0025b74da2dcdd7e59a77bcb0b0565564736f6c63430007050033"

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
	return address, tx, &KillCord{KillCordCaller: KillCordCaller{contract: contract}, KillCordTransactor: KillCordTransactor{contract: contract}, KillCordFilterer: KillCordFilterer{contract: contract}}, nil
}

// KillCord is an auto generated Go binding around an Ethereum contract.
type KillCord struct {
	KillCordCaller     // Read-only binding to the contract
	KillCordTransactor // Write-only binding to the contract
	KillCordFilterer   // Log filterer for contract events
}

// KillCordCaller is an auto generated read-only Go binding around an Ethereum contract.
type KillCordCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KillCordTransactor is an auto generated write-only Go binding around an Ethereum contract.
type KillCordTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KillCordFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KillCordFilterer struct {
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
	contract, err := bindKillCord(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KillCord{KillCordCaller: KillCordCaller{contract: contract}, KillCordTransactor: KillCordTransactor{contract: contract}, KillCordFilterer: KillCordFilterer{contract: contract}}, nil
}

// NewKillCordCaller creates a new read-only instance of KillCord, bound to a specific deployed contract.
func NewKillCordCaller(address common.Address, caller bind.ContractCaller) (*KillCordCaller, error) {
	contract, err := bindKillCord(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KillCordCaller{contract: contract}, nil
}

// NewKillCordTransactor creates a new write-only instance of KillCord, bound to a specific deployed contract.
func NewKillCordTransactor(address common.Address, transactor bind.ContractTransactor) (*KillCordTransactor, error) {
	contract, err := bindKillCord(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KillCordTransactor{contract: contract}, nil
}

// NewKillCordFilterer creates a new log filterer instance of KillCord, bound to a specific deployed contract.
func NewKillCordFilterer(address common.Address, filterer bind.ContractFilterer) (*KillCordFilterer, error) {
	contract, err := bindKillCord(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KillCordFilterer{contract: contract}, nil
}

// bindKillCord binds a generic wrapper to an already deployed contract.
func bindKillCord(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KillCordABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KillCord *KillCordRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
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
func (_KillCord *KillCordCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
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
// Solidity: function getKey() view returns(string)
func (_KillCord *KillCordCaller) GetKey(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _KillCord.contract.Call(opts, &out, "getKey")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetKey is a free data retrieval call binding the contract method 0x82678dd6.
//
// Solidity: function getKey() view returns(string)
func (_KillCord *KillCordSession) GetKey() (string, error) {
	return _KillCord.Contract.GetKey(&_KillCord.CallOpts)
}

// GetKey is a free data retrieval call binding the contract method 0x82678dd6.
//
// Solidity: function getKey() view returns(string)
func (_KillCord *KillCordCallerSession) GetKey() (string, error) {
	return _KillCord.Contract.GetKey(&_KillCord.CallOpts)
}

// GetLastCheckIn is a free data retrieval call binding the contract method 0xb1bedda3.
//
// Solidity: function getLastCheckIn() view returns(uint256)
func (_KillCord *KillCordCaller) GetLastCheckIn(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _KillCord.contract.Call(opts, &out, "getLastCheckIn")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLastCheckIn is a free data retrieval call binding the contract method 0xb1bedda3.
//
// Solidity: function getLastCheckIn() view returns(uint256)
func (_KillCord *KillCordSession) GetLastCheckIn() (*big.Int, error) {
	return _KillCord.Contract.GetLastCheckIn(&_KillCord.CallOpts)
}

// GetLastCheckIn is a free data retrieval call binding the contract method 0xb1bedda3.
//
// Solidity: function getLastCheckIn() view returns(uint256)
func (_KillCord *KillCordCallerSession) GetLastCheckIn() (*big.Int, error) {
	return _KillCord.Contract.GetLastCheckIn(&_KillCord.CallOpts)
}

// GetOwner is a free data retrieval call binding the contract method 0x893d20e8.
//
// Solidity: function getOwner() view returns(address)
func (_KillCord *KillCordCaller) GetOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KillCord.contract.Call(opts, &out, "getOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOwner is a free data retrieval call binding the contract method 0x893d20e8.
//
// Solidity: function getOwner() view returns(address)
func (_KillCord *KillCordSession) GetOwner() (common.Address, error) {
	return _KillCord.Contract.GetOwner(&_KillCord.CallOpts)
}

// GetOwner is a free data retrieval call binding the contract method 0x893d20e8.
//
// Solidity: function getOwner() view returns(address)
func (_KillCord *KillCordCallerSession) GetOwner() (common.Address, error) {
	return _KillCord.Contract.GetOwner(&_KillCord.CallOpts)
}

// GetPayloadEndpoint is a free data retrieval call binding the contract method 0x9f417563.
//
// Solidity: function getPayloadEndpoint() view returns(string)
func (_KillCord *KillCordCaller) GetPayloadEndpoint(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _KillCord.contract.Call(opts, &out, "getPayloadEndpoint")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetPayloadEndpoint is a free data retrieval call binding the contract method 0x9f417563.
//
// Solidity: function getPayloadEndpoint() view returns(string)
func (_KillCord *KillCordSession) GetPayloadEndpoint() (string, error) {
	return _KillCord.Contract.GetPayloadEndpoint(&_KillCord.CallOpts)
}

// GetPayloadEndpoint is a free data retrieval call binding the contract method 0x9f417563.
//
// Solidity: function getPayloadEndpoint() view returns(string)
func (_KillCord *KillCordCallerSession) GetPayloadEndpoint() (string, error) {
	return _KillCord.Contract.GetPayloadEndpoint(&_KillCord.CallOpts)
}

// GetPublisher is a free data retrieval call binding the contract method 0xdbf4ab4e.
//
// Solidity: function getPublisher() view returns(address)
func (_KillCord *KillCordCaller) GetPublisher(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _KillCord.contract.Call(opts, &out, "getPublisher")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetPublisher is a free data retrieval call binding the contract method 0xdbf4ab4e.
//
// Solidity: function getPublisher() view returns(address)
func (_KillCord *KillCordSession) GetPublisher() (common.Address, error) {
	return _KillCord.Contract.GetPublisher(&_KillCord.CallOpts)
}

// GetPublisher is a free data retrieval call binding the contract method 0xdbf4ab4e.
//
// Solidity: function getPublisher() view returns(address)
func (_KillCord *KillCordCallerSession) GetPublisher() (common.Address, error) {
	return _KillCord.Contract.GetPublisher(&_KillCord.CallOpts)
}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() view returns(string)
func (_KillCord *KillCordCaller) GetVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _KillCord.contract.Call(opts, &out, "getVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() view returns(string)
func (_KillCord *KillCordSession) GetVersion() (string, error) {
	return _KillCord.Contract.GetVersion(&_KillCord.CallOpts)
}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() view returns(string)
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
// Solidity: function setKey(string k) returns()
func (_KillCord *KillCordTransactor) SetKey(opts *bind.TransactOpts, k string) (*types.Transaction, error) {
	return _KillCord.contract.Transact(opts, "setKey", k)
}

// SetKey is a paid mutator transaction binding the contract method 0xaf42d106.
//
// Solidity: function setKey(string k) returns()
func (_KillCord *KillCordSession) SetKey(k string) (*types.Transaction, error) {
	return _KillCord.Contract.SetKey(&_KillCord.TransactOpts, k)
}

// SetKey is a paid mutator transaction binding the contract method 0xaf42d106.
//
// Solidity: function setKey(string k) returns()
func (_KillCord *KillCordTransactorSession) SetKey(k string) (*types.Transaction, error) {
	return _KillCord.Contract.SetKey(&_KillCord.TransactOpts, k)
}

// SetPayloadEndpoint is a paid mutator transaction binding the contract method 0x8f4f106e.
//
// Solidity: function setPayloadEndpoint(string s) returns()
func (_KillCord *KillCordTransactor) SetPayloadEndpoint(opts *bind.TransactOpts, s string) (*types.Transaction, error) {
	return _KillCord.contract.Transact(opts, "setPayloadEndpoint", s)
}

// SetPayloadEndpoint is a paid mutator transaction binding the contract method 0x8f4f106e.
//
// Solidity: function setPayloadEndpoint(string s) returns()
func (_KillCord *KillCordSession) SetPayloadEndpoint(s string) (*types.Transaction, error) {
	return _KillCord.Contract.SetPayloadEndpoint(&_KillCord.TransactOpts, s)
}

// SetPayloadEndpoint is a paid mutator transaction binding the contract method 0x8f4f106e.
//
// Solidity: function setPayloadEndpoint(string s) returns()
func (_KillCord *KillCordTransactorSession) SetPayloadEndpoint(s string) (*types.Transaction, error) {
	return _KillCord.Contract.SetPayloadEndpoint(&_KillCord.TransactOpts, s)
}
