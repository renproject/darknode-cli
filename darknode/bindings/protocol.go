// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ProtocolMetaData contains all meta data concerning the Protocol contract.
var ProtocolMetaData = &bind.MetaData{
	ABI: "[{\"constant\":false,\"inputs\":[{\"internalType\":\"contractDarknodeRegistry\",\"name\":\"_newDarknodeRegistry\",\"type\":\"address\"}],\"name\":\"_updateDarknodeRegistry\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"contractShifterRegistry\",\"name\":\"_newShifterRegistry\",\"type\":\"address\"}],\"name\":\"_updateShifterRegistry\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"darknodePayment\",\"outputs\":[{\"internalType\":\"contractDarknodePayment\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"darknodePaymentStore\",\"outputs\":[{\"internalType\":\"contractDarknodePaymentStore\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"darknodeRegistry\",\"outputs\":[{\"internalType\":\"contractDarknodeRegistry\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"darknodeRegistryStore\",\"outputs\":[{\"internalType\":\"contractDarknodeRegistryStore\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"darknodeSlasher\",\"outputs\":[{\"internalType\":\"contractDarknodeSlasher\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_start\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_count\",\"type\":\"uint256\"}],\"name\":\"getShiftedTokens\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"string\",\"name\":\"_tokenSymbol\",\"type\":\"string\"}],\"name\":\"getShifterBySymbol\",\"outputs\":[{\"internalType\":\"contractIShifter\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_tokenAddress\",\"type\":\"address\"}],\"name\":\"getShifterByToken\",\"outputs\":[{\"internalType\":\"contractIShifter\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_start\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_count\",\"type\":\"uint256\"}],\"name\":\"getShifters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"string\",\"name\":\"_tokenSymbol\",\"type\":\"string\"}],\"name\":\"getTokenBySymbol\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"renToken\",\"outputs\":[{\"internalType\":\"contractRenToken\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"shifterRegistry\",\"outputs\":[{\"internalType\":\"contractShifterRegistry\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// ProtocolABI is the input ABI used to generate the binding from.
// Deprecated: Use ProtocolMetaData.ABI instead.
var ProtocolABI = ProtocolMetaData.ABI

// Protocol is an auto generated Go binding around an Ethereum contract.
type Protocol struct {
	ProtocolCaller     // Read-only binding to the contract
	ProtocolTransactor // Write-only binding to the contract
	ProtocolFilterer   // Log filterer for contract events
}

// ProtocolCaller is an auto generated read-only Go binding around an Ethereum contract.
type ProtocolCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProtocolTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ProtocolTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProtocolFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ProtocolFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProtocolSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ProtocolSession struct {
	Contract     *Protocol         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ProtocolCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ProtocolCallerSession struct {
	Contract *ProtocolCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ProtocolTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ProtocolTransactorSession struct {
	Contract     *ProtocolTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ProtocolRaw is an auto generated low-level Go binding around an Ethereum contract.
type ProtocolRaw struct {
	Contract *Protocol // Generic contract binding to access the raw methods on
}

// ProtocolCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ProtocolCallerRaw struct {
	Contract *ProtocolCaller // Generic read-only contract binding to access the raw methods on
}

// ProtocolTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ProtocolTransactorRaw struct {
	Contract *ProtocolTransactor // Generic write-only contract binding to access the raw methods on
}

// NewProtocol creates a new instance of Protocol, bound to a specific deployed contract.
func NewProtocol(address common.Address, backend bind.ContractBackend) (*Protocol, error) {
	contract, err := bindProtocol(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Protocol{ProtocolCaller: ProtocolCaller{contract: contract}, ProtocolTransactor: ProtocolTransactor{contract: contract}, ProtocolFilterer: ProtocolFilterer{contract: contract}}, nil
}

// NewProtocolCaller creates a new read-only instance of Protocol, bound to a specific deployed contract.
func NewProtocolCaller(address common.Address, caller bind.ContractCaller) (*ProtocolCaller, error) {
	contract, err := bindProtocol(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ProtocolCaller{contract: contract}, nil
}

// NewProtocolTransactor creates a new write-only instance of Protocol, bound to a specific deployed contract.
func NewProtocolTransactor(address common.Address, transactor bind.ContractTransactor) (*ProtocolTransactor, error) {
	contract, err := bindProtocol(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ProtocolTransactor{contract: contract}, nil
}

// NewProtocolFilterer creates a new log filterer instance of Protocol, bound to a specific deployed contract.
func NewProtocolFilterer(address common.Address, filterer bind.ContractFilterer) (*ProtocolFilterer, error) {
	contract, err := bindProtocol(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ProtocolFilterer{contract: contract}, nil
}

// bindProtocol binds a generic wrapper to an already deployed contract.
func bindProtocol(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ProtocolABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Protocol *ProtocolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Protocol.Contract.ProtocolCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Protocol *ProtocolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Protocol.Contract.ProtocolTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Protocol *ProtocolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Protocol.Contract.ProtocolTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Protocol *ProtocolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Protocol.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Protocol *ProtocolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Protocol.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Protocol *ProtocolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Protocol.Contract.contract.Transact(opts, method, params...)
}

// DarknodePayment is a free data retrieval call binding the contract method 0xb6b34c67.
//
// Solidity: function darknodePayment() view returns(address)
func (_Protocol *ProtocolCaller) DarknodePayment(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Protocol.contract.Call(opts, &out, "darknodePayment")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DarknodePayment is a free data retrieval call binding the contract method 0xb6b34c67.
//
// Solidity: function darknodePayment() view returns(address)
func (_Protocol *ProtocolSession) DarknodePayment() (common.Address, error) {
	return _Protocol.Contract.DarknodePayment(&_Protocol.CallOpts)
}

// DarknodePayment is a free data retrieval call binding the contract method 0xb6b34c67.
//
// Solidity: function darknodePayment() view returns(address)
func (_Protocol *ProtocolCallerSession) DarknodePayment() (common.Address, error) {
	return _Protocol.Contract.DarknodePayment(&_Protocol.CallOpts)
}

// DarknodePaymentStore is a free data retrieval call binding the contract method 0xbcfe9fe5.
//
// Solidity: function darknodePaymentStore() view returns(address)
func (_Protocol *ProtocolCaller) DarknodePaymentStore(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Protocol.contract.Call(opts, &out, "darknodePaymentStore")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DarknodePaymentStore is a free data retrieval call binding the contract method 0xbcfe9fe5.
//
// Solidity: function darknodePaymentStore() view returns(address)
func (_Protocol *ProtocolSession) DarknodePaymentStore() (common.Address, error) {
	return _Protocol.Contract.DarknodePaymentStore(&_Protocol.CallOpts)
}

// DarknodePaymentStore is a free data retrieval call binding the contract method 0xbcfe9fe5.
//
// Solidity: function darknodePaymentStore() view returns(address)
func (_Protocol *ProtocolCallerSession) DarknodePaymentStore() (common.Address, error) {
	return _Protocol.Contract.DarknodePaymentStore(&_Protocol.CallOpts)
}

// DarknodeRegistry is a free data retrieval call binding the contract method 0x9e45e0d0.
//
// Solidity: function darknodeRegistry() view returns(address)
func (_Protocol *ProtocolCaller) DarknodeRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Protocol.contract.Call(opts, &out, "darknodeRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DarknodeRegistry is a free data retrieval call binding the contract method 0x9e45e0d0.
//
// Solidity: function darknodeRegistry() view returns(address)
func (_Protocol *ProtocolSession) DarknodeRegistry() (common.Address, error) {
	return _Protocol.Contract.DarknodeRegistry(&_Protocol.CallOpts)
}

// DarknodeRegistry is a free data retrieval call binding the contract method 0x9e45e0d0.
//
// Solidity: function darknodeRegistry() view returns(address)
func (_Protocol *ProtocolCallerSession) DarknodeRegistry() (common.Address, error) {
	return _Protocol.Contract.DarknodeRegistry(&_Protocol.CallOpts)
}

// DarknodeRegistryStore is a free data retrieval call binding the contract method 0x981a58e9.
//
// Solidity: function darknodeRegistryStore() view returns(address)
func (_Protocol *ProtocolCaller) DarknodeRegistryStore(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Protocol.contract.Call(opts, &out, "darknodeRegistryStore")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DarknodeRegistryStore is a free data retrieval call binding the contract method 0x981a58e9.
//
// Solidity: function darknodeRegistryStore() view returns(address)
func (_Protocol *ProtocolSession) DarknodeRegistryStore() (common.Address, error) {
	return _Protocol.Contract.DarknodeRegistryStore(&_Protocol.CallOpts)
}

// DarknodeRegistryStore is a free data retrieval call binding the contract method 0x981a58e9.
//
// Solidity: function darknodeRegistryStore() view returns(address)
func (_Protocol *ProtocolCallerSession) DarknodeRegistryStore() (common.Address, error) {
	return _Protocol.Contract.DarknodeRegistryStore(&_Protocol.CallOpts)
}

// DarknodeSlasher is a free data retrieval call binding the contract method 0xc4259606.
//
// Solidity: function darknodeSlasher() view returns(address)
func (_Protocol *ProtocolCaller) DarknodeSlasher(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Protocol.contract.Call(opts, &out, "darknodeSlasher")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DarknodeSlasher is a free data retrieval call binding the contract method 0xc4259606.
//
// Solidity: function darknodeSlasher() view returns(address)
func (_Protocol *ProtocolSession) DarknodeSlasher() (common.Address, error) {
	return _Protocol.Contract.DarknodeSlasher(&_Protocol.CallOpts)
}

// DarknodeSlasher is a free data retrieval call binding the contract method 0xc4259606.
//
// Solidity: function darknodeSlasher() view returns(address)
func (_Protocol *ProtocolCallerSession) DarknodeSlasher() (common.Address, error) {
	return _Protocol.Contract.DarknodeSlasher(&_Protocol.CallOpts)
}

// GetShiftedTokens is a free data retrieval call binding the contract method 0xa47c84bd.
//
// Solidity: function getShiftedTokens(address _start, uint256 _count) view returns(address[])
func (_Protocol *ProtocolCaller) GetShiftedTokens(opts *bind.CallOpts, _start common.Address, _count *big.Int) ([]common.Address, error) {
	var out []interface{}
	err := _Protocol.contract.Call(opts, &out, "getShiftedTokens", _start, _count)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetShiftedTokens is a free data retrieval call binding the contract method 0xa47c84bd.
//
// Solidity: function getShiftedTokens(address _start, uint256 _count) view returns(address[])
func (_Protocol *ProtocolSession) GetShiftedTokens(_start common.Address, _count *big.Int) ([]common.Address, error) {
	return _Protocol.Contract.GetShiftedTokens(&_Protocol.CallOpts, _start, _count)
}

// GetShiftedTokens is a free data retrieval call binding the contract method 0xa47c84bd.
//
// Solidity: function getShiftedTokens(address _start, uint256 _count) view returns(address[])
func (_Protocol *ProtocolCallerSession) GetShiftedTokens(_start common.Address, _count *big.Int) ([]common.Address, error) {
	return _Protocol.Contract.GetShiftedTokens(&_Protocol.CallOpts, _start, _count)
}

// GetShifterBySymbol is a free data retrieval call binding the contract method 0x92a29e30.
//
// Solidity: function getShifterBySymbol(string _tokenSymbol) view returns(address)
func (_Protocol *ProtocolCaller) GetShifterBySymbol(opts *bind.CallOpts, _tokenSymbol string) (common.Address, error) {
	var out []interface{}
	err := _Protocol.contract.Call(opts, &out, "getShifterBySymbol", _tokenSymbol)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetShifterBySymbol is a free data retrieval call binding the contract method 0x92a29e30.
//
// Solidity: function getShifterBySymbol(string _tokenSymbol) view returns(address)
func (_Protocol *ProtocolSession) GetShifterBySymbol(_tokenSymbol string) (common.Address, error) {
	return _Protocol.Contract.GetShifterBySymbol(&_Protocol.CallOpts, _tokenSymbol)
}

// GetShifterBySymbol is a free data retrieval call binding the contract method 0x92a29e30.
//
// Solidity: function getShifterBySymbol(string _tokenSymbol) view returns(address)
func (_Protocol *ProtocolCallerSession) GetShifterBySymbol(_tokenSymbol string) (common.Address, error) {
	return _Protocol.Contract.GetShifterBySymbol(&_Protocol.CallOpts, _tokenSymbol)
}

// GetShifterByToken is a free data retrieval call binding the contract method 0x28c4f410.
//
// Solidity: function getShifterByToken(address _tokenAddress) view returns(address)
func (_Protocol *ProtocolCaller) GetShifterByToken(opts *bind.CallOpts, _tokenAddress common.Address) (common.Address, error) {
	var out []interface{}
	err := _Protocol.contract.Call(opts, &out, "getShifterByToken", _tokenAddress)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetShifterByToken is a free data retrieval call binding the contract method 0x28c4f410.
//
// Solidity: function getShifterByToken(address _tokenAddress) view returns(address)
func (_Protocol *ProtocolSession) GetShifterByToken(_tokenAddress common.Address) (common.Address, error) {
	return _Protocol.Contract.GetShifterByToken(&_Protocol.CallOpts, _tokenAddress)
}

// GetShifterByToken is a free data retrieval call binding the contract method 0x28c4f410.
//
// Solidity: function getShifterByToken(address _tokenAddress) view returns(address)
func (_Protocol *ProtocolCallerSession) GetShifterByToken(_tokenAddress common.Address) (common.Address, error) {
	return _Protocol.Contract.GetShifterByToken(&_Protocol.CallOpts, _tokenAddress)
}

// GetShifters is a free data retrieval call binding the contract method 0xdeeb2efe.
//
// Solidity: function getShifters(address _start, uint256 _count) view returns(address[])
func (_Protocol *ProtocolCaller) GetShifters(opts *bind.CallOpts, _start common.Address, _count *big.Int) ([]common.Address, error) {
	var out []interface{}
	err := _Protocol.contract.Call(opts, &out, "getShifters", _start, _count)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetShifters is a free data retrieval call binding the contract method 0xdeeb2efe.
//
// Solidity: function getShifters(address _start, uint256 _count) view returns(address[])
func (_Protocol *ProtocolSession) GetShifters(_start common.Address, _count *big.Int) ([]common.Address, error) {
	return _Protocol.Contract.GetShifters(&_Protocol.CallOpts, _start, _count)
}

// GetShifters is a free data retrieval call binding the contract method 0xdeeb2efe.
//
// Solidity: function getShifters(address _start, uint256 _count) view returns(address[])
func (_Protocol *ProtocolCallerSession) GetShifters(_start common.Address, _count *big.Int) ([]common.Address, error) {
	return _Protocol.Contract.GetShifters(&_Protocol.CallOpts, _start, _count)
}

// GetTokenBySymbol is a free data retrieval call binding the contract method 0xefa74f1f.
//
// Solidity: function getTokenBySymbol(string _tokenSymbol) view returns(address)
func (_Protocol *ProtocolCaller) GetTokenBySymbol(opts *bind.CallOpts, _tokenSymbol string) (common.Address, error) {
	var out []interface{}
	err := _Protocol.contract.Call(opts, &out, "getTokenBySymbol", _tokenSymbol)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetTokenBySymbol is a free data retrieval call binding the contract method 0xefa74f1f.
//
// Solidity: function getTokenBySymbol(string _tokenSymbol) view returns(address)
func (_Protocol *ProtocolSession) GetTokenBySymbol(_tokenSymbol string) (common.Address, error) {
	return _Protocol.Contract.GetTokenBySymbol(&_Protocol.CallOpts, _tokenSymbol)
}

// GetTokenBySymbol is a free data retrieval call binding the contract method 0xefa74f1f.
//
// Solidity: function getTokenBySymbol(string _tokenSymbol) view returns(address)
func (_Protocol *ProtocolCallerSession) GetTokenBySymbol(_tokenSymbol string) (common.Address, error) {
	return _Protocol.Contract.GetTokenBySymbol(&_Protocol.CallOpts, _tokenSymbol)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Protocol *ProtocolCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Protocol.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Protocol *ProtocolSession) Owner() (common.Address, error) {
	return _Protocol.Contract.Owner(&_Protocol.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Protocol *ProtocolCallerSession) Owner() (common.Address, error) {
	return _Protocol.Contract.Owner(&_Protocol.CallOpts)
}

// RenToken is a free data retrieval call binding the contract method 0x34246f9b.
//
// Solidity: function renToken() view returns(address)
func (_Protocol *ProtocolCaller) RenToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Protocol.contract.Call(opts, &out, "renToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RenToken is a free data retrieval call binding the contract method 0x34246f9b.
//
// Solidity: function renToken() view returns(address)
func (_Protocol *ProtocolSession) RenToken() (common.Address, error) {
	return _Protocol.Contract.RenToken(&_Protocol.CallOpts)
}

// RenToken is a free data retrieval call binding the contract method 0x34246f9b.
//
// Solidity: function renToken() view returns(address)
func (_Protocol *ProtocolCallerSession) RenToken() (common.Address, error) {
	return _Protocol.Contract.RenToken(&_Protocol.CallOpts)
}

// ShifterRegistry is a free data retrieval call binding the contract method 0xe6a12ca9.
//
// Solidity: function shifterRegistry() view returns(address)
func (_Protocol *ProtocolCaller) ShifterRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Protocol.contract.Call(opts, &out, "shifterRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ShifterRegistry is a free data retrieval call binding the contract method 0xe6a12ca9.
//
// Solidity: function shifterRegistry() view returns(address)
func (_Protocol *ProtocolSession) ShifterRegistry() (common.Address, error) {
	return _Protocol.Contract.ShifterRegistry(&_Protocol.CallOpts)
}

// ShifterRegistry is a free data retrieval call binding the contract method 0xe6a12ca9.
//
// Solidity: function shifterRegistry() view returns(address)
func (_Protocol *ProtocolCallerSession) ShifterRegistry() (common.Address, error) {
	return _Protocol.Contract.ShifterRegistry(&_Protocol.CallOpts)
}

// UpdateDarknodeRegistry is a paid mutator transaction binding the contract method 0x179091a4.
//
// Solidity: function _updateDarknodeRegistry(address _newDarknodeRegistry) returns()
func (_Protocol *ProtocolTransactor) UpdateDarknodeRegistry(opts *bind.TransactOpts, _newDarknodeRegistry common.Address) (*types.Transaction, error) {
	return _Protocol.contract.Transact(opts, "_updateDarknodeRegistry", _newDarknodeRegistry)
}

// UpdateDarknodeRegistry is a paid mutator transaction binding the contract method 0x179091a4.
//
// Solidity: function _updateDarknodeRegistry(address _newDarknodeRegistry) returns()
func (_Protocol *ProtocolSession) UpdateDarknodeRegistry(_newDarknodeRegistry common.Address) (*types.Transaction, error) {
	return _Protocol.Contract.UpdateDarknodeRegistry(&_Protocol.TransactOpts, _newDarknodeRegistry)
}

// UpdateDarknodeRegistry is a paid mutator transaction binding the contract method 0x179091a4.
//
// Solidity: function _updateDarknodeRegistry(address _newDarknodeRegistry) returns()
func (_Protocol *ProtocolTransactorSession) UpdateDarknodeRegistry(_newDarknodeRegistry common.Address) (*types.Transaction, error) {
	return _Protocol.Contract.UpdateDarknodeRegistry(&_Protocol.TransactOpts, _newDarknodeRegistry)
}

// UpdateShifterRegistry is a paid mutator transaction binding the contract method 0x715a6132.
//
// Solidity: function _updateShifterRegistry(address _newShifterRegistry) returns()
func (_Protocol *ProtocolTransactor) UpdateShifterRegistry(opts *bind.TransactOpts, _newShifterRegistry common.Address) (*types.Transaction, error) {
	return _Protocol.contract.Transact(opts, "_updateShifterRegistry", _newShifterRegistry)
}

// UpdateShifterRegistry is a paid mutator transaction binding the contract method 0x715a6132.
//
// Solidity: function _updateShifterRegistry(address _newShifterRegistry) returns()
func (_Protocol *ProtocolSession) UpdateShifterRegistry(_newShifterRegistry common.Address) (*types.Transaction, error) {
	return _Protocol.Contract.UpdateShifterRegistry(&_Protocol.TransactOpts, _newShifterRegistry)
}

// UpdateShifterRegistry is a paid mutator transaction binding the contract method 0x715a6132.
//
// Solidity: function _updateShifterRegistry(address _newShifterRegistry) returns()
func (_Protocol *ProtocolTransactorSession) UpdateShifterRegistry(_newShifterRegistry common.Address) (*types.Transaction, error) {
	return _Protocol.Contract.UpdateShifterRegistry(&_Protocol.TransactOpts, _newShifterRegistry)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _owner) returns()
func (_Protocol *ProtocolTransactor) Initialize(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _Protocol.contract.Transact(opts, "initialize", _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _owner) returns()
func (_Protocol *ProtocolSession) Initialize(_owner common.Address) (*types.Transaction, error) {
	return _Protocol.Contract.Initialize(&_Protocol.TransactOpts, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _owner) returns()
func (_Protocol *ProtocolTransactorSession) Initialize(_owner common.Address) (*types.Transaction, error) {
	return _Protocol.Contract.Initialize(&_Protocol.TransactOpts, _owner)
}
