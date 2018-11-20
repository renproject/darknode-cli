// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

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

// BindingsABI is the input ABI used to generate the binding from.
const BindingsABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"}],\"name\":\"isPendingRegistration\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numDarknodesNextEpoch\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"},{\"name\":\"_publicKey\",\"type\":\"bytes\"},{\"name\":\"_bond\",\"type\":\"uint256\"}],\"name\":\"register\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_nextMinimumBond\",\"type\":\"uint256\"}],\"name\":\"updateMinimumBond\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numDarknodes\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"}],\"name\":\"getDarknodeOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"nextSlasher\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"}],\"name\":\"isPendingDeregistration\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_start\",\"type\":\"address\"},{\"name\":\"_count\",\"type\":\"uint256\"}],\"name\":\"getPreviousDarknodes\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"nextMinimumEpochInterval\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minimumEpochInterval\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_prover\",\"type\":\"address\"},{\"name\":\"_challenger1\",\"type\":\"address\"},{\"name\":\"_challenger2\",\"type\":\"address\"}],\"name\":\"slash\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"}],\"name\":\"isRefundable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"previousEpoch\",\"outputs\":[{\"name\":\"epochhash\",\"type\":\"uint256\"},{\"name\":\"blocknumber\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"nextMinimumBond\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_nextMinimumEpochInterval\",\"type\":\"uint256\"}],\"name\":\"updateMinimumEpochInterval\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"nextMinimumPodSize\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numDarknodesPreviousEpoch\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentEpoch\",\"outputs\":[{\"name\":\"epochhash\",\"type\":\"uint256\"},{\"name\":\"blocknumber\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"}],\"name\":\"isRegisteredInPreviousEpoch\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"}],\"name\":\"isDeregistered\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_nextMinimumPodSize\",\"type\":\"uint256\"}],\"name\":\"updateMinimumPodSize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"}],\"name\":\"deregister\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"}],\"name\":\"getDarknodePublicKey\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ren\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"epoch\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"store\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minimumBond\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"slasher\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_slasher\",\"type\":\"address\"}],\"name\":\"updateSlasher\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"}],\"name\":\"getDarknodeBond\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferStoreOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"}],\"name\":\"isRegistered\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minimumPodSize\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"}],\"name\":\"isDeregisterable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_start\",\"type\":\"address\"},{\"name\":\"_count\",\"type\":\"uint256\"}],\"name\":\"getDarknodes\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"}],\"name\":\"refund\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_darknodeID\",\"type\":\"address\"}],\"name\":\"isRefunded\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_VERSION\",\"type\":\"string\"},{\"name\":\"_renAddress\",\"type\":\"address\"},{\"name\":\"_storeAddress\",\"type\":\"address\"},{\"name\":\"_minimumBond\",\"type\":\"uint256\"},{\"name\":\"_minimumPodSize\",\"type\":\"uint256\"},{\"name\":\"_minimumEpochInterval\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_darknodeID\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_bond\",\"type\":\"uint256\"}],\"name\":\"LogDarknodeRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_darknodeID\",\"type\":\"address\"}],\"name\":\"LogDarknodeDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"LogDarknodeOwnerRefunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"LogNewEpoch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"previousMinimumBond\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"nextMinimumBond\",\"type\":\"uint256\"}],\"name\":\"LogMinimumBondUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"previousMinimumPodSize\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"nextMinimumPodSize\",\"type\":\"uint256\"}],\"name\":\"LogMinimumPodSizeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"previousMinimumEpochInterval\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"nextMinimumEpochInterval\",\"type\":\"uint256\"}],\"name\":\"LogMinimumEpochIntervalUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"previousSlasher\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"nextSlasher\",\"type\":\"address\"}],\"name\":\"LogSlasherUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"}],\"name\":\"OwnershipRenounced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

// Bindings is an auto generated Go binding around an Ethereum contract.
type Bindings struct {
	BindingsCaller     // Read-only binding to the contract
	BindingsTransactor // Write-only binding to the contract
	BindingsFilterer   // Log filterer for contract events
}

// BindingsCaller is an auto generated read-only Go binding around an Ethereum contract.
type BindingsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BindingsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BindingsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BindingsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BindingsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BindingsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BindingsSession struct {
	Contract     *Bindings         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BindingsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BindingsCallerSession struct {
	Contract *BindingsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// BindingsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BindingsTransactorSession struct {
	Contract     *BindingsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// BindingsRaw is an auto generated low-level Go binding around an Ethereum contract.
type BindingsRaw struct {
	Contract *Bindings // Generic contract binding to access the raw methods on
}

// BindingsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BindingsCallerRaw struct {
	Contract *BindingsCaller // Generic read-only contract binding to access the raw methods on
}

// BindingsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BindingsTransactorRaw struct {
	Contract *BindingsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBindings creates a new instance of Bindings, bound to a specific deployed contract.
func NewBindings(address common.Address, backend bind.ContractBackend) (*Bindings, error) {
	contract, err := bindBindings(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Bindings{BindingsCaller: BindingsCaller{contract: contract}, BindingsTransactor: BindingsTransactor{contract: contract}, BindingsFilterer: BindingsFilterer{contract: contract}}, nil
}

// NewBindingsCaller creates a new read-only instance of Bindings, bound to a specific deployed contract.
func NewBindingsCaller(address common.Address, caller bind.ContractCaller) (*BindingsCaller, error) {
	contract, err := bindBindings(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BindingsCaller{contract: contract}, nil
}

// NewBindingsTransactor creates a new write-only instance of Bindings, bound to a specific deployed contract.
func NewBindingsTransactor(address common.Address, transactor bind.ContractTransactor) (*BindingsTransactor, error) {
	contract, err := bindBindings(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BindingsTransactor{contract: contract}, nil
}

// NewBindingsFilterer creates a new log filterer instance of Bindings, bound to a specific deployed contract.
func NewBindingsFilterer(address common.Address, filterer bind.ContractFilterer) (*BindingsFilterer, error) {
	contract, err := bindBindings(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BindingsFilterer{contract: contract}, nil
}

// bindBindings binds a generic wrapper to an already deployed contract.
func bindBindings(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BindingsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bindings *BindingsRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Bindings.Contract.BindingsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bindings *BindingsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bindings.Contract.BindingsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bindings *BindingsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bindings.Contract.BindingsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bindings *BindingsCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Bindings.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bindings *BindingsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bindings.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bindings *BindingsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bindings.Contract.contract.Transact(opts, method, params...)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_Bindings *BindingsCaller) VERSION(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Bindings.contract.Call(opts, out, "VERSION")
	return *ret0, err
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_Bindings *BindingsSession) VERSION() (string, error) {
	return _Bindings.Contract.VERSION(&_Bindings.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_Bindings *BindingsCallerSession) VERSION() (string, error) {
	return _Bindings.Contract.VERSION(&_Bindings.CallOpts)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() constant returns(epochhash uint256, blocknumber uint256)
func (_Bindings *BindingsCaller) CurrentEpoch(opts *bind.CallOpts) (struct {
	Epochhash   *big.Int
	Blocknumber *big.Int
}, error) {
	ret := new(struct {
		Epochhash   *big.Int
		Blocknumber *big.Int
	})
	out := ret
	err := _Bindings.contract.Call(opts, out, "currentEpoch")
	return *ret, err
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() constant returns(epochhash uint256, blocknumber uint256)
func (_Bindings *BindingsSession) CurrentEpoch() (struct {
	Epochhash   *big.Int
	Blocknumber *big.Int
}, error) {
	return _Bindings.Contract.CurrentEpoch(&_Bindings.CallOpts)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() constant returns(epochhash uint256, blocknumber uint256)
func (_Bindings *BindingsCallerSession) CurrentEpoch() (struct {
	Epochhash   *big.Int
	Blocknumber *big.Int
}, error) {
	return _Bindings.Contract.CurrentEpoch(&_Bindings.CallOpts)
}

// GetDarknodeBond is a free data retrieval call binding the contract method 0xba0f5b20.
//
// Solidity: function getDarknodeBond(_darknodeID address) constant returns(uint256)
func (_Bindings *BindingsCaller) GetDarknodeBond(opts *bind.CallOpts, _darknodeID common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Bindings.contract.Call(opts, out, "getDarknodeBond", _darknodeID)
	return *ret0, err
}

// GetDarknodeBond is a free data retrieval call binding the contract method 0xba0f5b20.
//
// Solidity: function getDarknodeBond(_darknodeID address) constant returns(uint256)
func (_Bindings *BindingsSession) GetDarknodeBond(_darknodeID common.Address) (*big.Int, error) {
	return _Bindings.Contract.GetDarknodeBond(&_Bindings.CallOpts, _darknodeID)
}

// GetDarknodeBond is a free data retrieval call binding the contract method 0xba0f5b20.
//
// Solidity: function getDarknodeBond(_darknodeID address) constant returns(uint256)
func (_Bindings *BindingsCallerSession) GetDarknodeBond(_darknodeID common.Address) (*big.Int, error) {
	return _Bindings.Contract.GetDarknodeBond(&_Bindings.CallOpts, _darknodeID)
}

// GetDarknodeOwner is a free data retrieval call binding the contract method 0x1cedf8a3.
//
// Solidity: function getDarknodeOwner(_darknodeID address) constant returns(address)
func (_Bindings *BindingsCaller) GetDarknodeOwner(opts *bind.CallOpts, _darknodeID common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Bindings.contract.Call(opts, out, "getDarknodeOwner", _darknodeID)
	return *ret0, err
}

// GetDarknodeOwner is a free data retrieval call binding the contract method 0x1cedf8a3.
//
// Solidity: function getDarknodeOwner(_darknodeID address) constant returns(address)
func (_Bindings *BindingsSession) GetDarknodeOwner(_darknodeID common.Address) (common.Address, error) {
	return _Bindings.Contract.GetDarknodeOwner(&_Bindings.CallOpts, _darknodeID)
}

// GetDarknodeOwner is a free data retrieval call binding the contract method 0x1cedf8a3.
//
// Solidity: function getDarknodeOwner(_darknodeID address) constant returns(address)
func (_Bindings *BindingsCallerSession) GetDarknodeOwner(_darknodeID common.Address) (common.Address, error) {
	return _Bindings.Contract.GetDarknodeOwner(&_Bindings.CallOpts, _darknodeID)
}

// GetDarknodePublicKey is a free data retrieval call binding the contract method 0x84d2688c.
//
// Solidity: function getDarknodePublicKey(_darknodeID address) constant returns(bytes)
func (_Bindings *BindingsCaller) GetDarknodePublicKey(opts *bind.CallOpts, _darknodeID common.Address) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _Bindings.contract.Call(opts, out, "getDarknodePublicKey", _darknodeID)
	return *ret0, err
}

// GetDarknodePublicKey is a free data retrieval call binding the contract method 0x84d2688c.
//
// Solidity: function getDarknodePublicKey(_darknodeID address) constant returns(bytes)
func (_Bindings *BindingsSession) GetDarknodePublicKey(_darknodeID common.Address) ([]byte, error) {
	return _Bindings.Contract.GetDarknodePublicKey(&_Bindings.CallOpts, _darknodeID)
}

// GetDarknodePublicKey is a free data retrieval call binding the contract method 0x84d2688c.
//
// Solidity: function getDarknodePublicKey(_darknodeID address) constant returns(bytes)
func (_Bindings *BindingsCallerSession) GetDarknodePublicKey(_darknodeID common.Address) ([]byte, error) {
	return _Bindings.Contract.GetDarknodePublicKey(&_Bindings.CallOpts, _darknodeID)
}

// GetDarknodes is a free data retrieval call binding the contract method 0xec5325c1.
//
// Solidity: function getDarknodes(_start address, _count uint256) constant returns(address[])
func (_Bindings *BindingsCaller) GetDarknodes(opts *bind.CallOpts, _start common.Address, _count *big.Int) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _Bindings.contract.Call(opts, out, "getDarknodes", _start, _count)
	return *ret0, err
}

// GetDarknodes is a free data retrieval call binding the contract method 0xec5325c1.
//
// Solidity: function getDarknodes(_start address, _count uint256) constant returns(address[])
func (_Bindings *BindingsSession) GetDarknodes(_start common.Address, _count *big.Int) ([]common.Address, error) {
	return _Bindings.Contract.GetDarknodes(&_Bindings.CallOpts, _start, _count)
}

// GetDarknodes is a free data retrieval call binding the contract method 0xec5325c1.
//
// Solidity: function getDarknodes(_start address, _count uint256) constant returns(address[])
func (_Bindings *BindingsCallerSession) GetDarknodes(_start common.Address, _count *big.Int) ([]common.Address, error) {
	return _Bindings.Contract.GetDarknodes(&_Bindings.CallOpts, _start, _count)
}

// GetPreviousDarknodes is a free data retrieval call binding the contract method 0x43846074.
//
// Solidity: function getPreviousDarknodes(_start address, _count uint256) constant returns(address[])
func (_Bindings *BindingsCaller) GetPreviousDarknodes(opts *bind.CallOpts, _start common.Address, _count *big.Int) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _Bindings.contract.Call(opts, out, "getPreviousDarknodes", _start, _count)
	return *ret0, err
}

// GetPreviousDarknodes is a free data retrieval call binding the contract method 0x43846074.
//
// Solidity: function getPreviousDarknodes(_start address, _count uint256) constant returns(address[])
func (_Bindings *BindingsSession) GetPreviousDarknodes(_start common.Address, _count *big.Int) ([]common.Address, error) {
	return _Bindings.Contract.GetPreviousDarknodes(&_Bindings.CallOpts, _start, _count)
}

// GetPreviousDarknodes is a free data retrieval call binding the contract method 0x43846074.
//
// Solidity: function getPreviousDarknodes(_start address, _count uint256) constant returns(address[])
func (_Bindings *BindingsCallerSession) GetPreviousDarknodes(_start common.Address, _count *big.Int) ([]common.Address, error) {
	return _Bindings.Contract.GetPreviousDarknodes(&_Bindings.CallOpts, _start, _count)
}

// IsDeregisterable is a free data retrieval call binding the contract method 0xe1878925.
//
// Solidity: function isDeregisterable(_darknodeID address) constant returns(bool)
func (_Bindings *BindingsCaller) IsDeregisterable(opts *bind.CallOpts, _darknodeID common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Bindings.contract.Call(opts, out, "isDeregisterable", _darknodeID)
	return *ret0, err
}

// IsDeregisterable is a free data retrieval call binding the contract method 0xe1878925.
//
// Solidity: function isDeregisterable(_darknodeID address) constant returns(bool)
func (_Bindings *BindingsSession) IsDeregisterable(_darknodeID common.Address) (bool, error) {
	return _Bindings.Contract.IsDeregisterable(&_Bindings.CallOpts, _darknodeID)
}

// IsDeregisterable is a free data retrieval call binding the contract method 0xe1878925.
//
// Solidity: function isDeregisterable(_darknodeID address) constant returns(bool)
func (_Bindings *BindingsCallerSession) IsDeregisterable(_darknodeID common.Address) (bool, error) {
	return _Bindings.Contract.IsDeregisterable(&_Bindings.CallOpts, _darknodeID)
}

// IsDeregistered is a free data retrieval call binding the contract method 0x8020fc1f.
//
// Solidity: function isDeregistered(_darknodeID address) constant returns(bool)
func (_Bindings *BindingsCaller) IsDeregistered(opts *bind.CallOpts, _darknodeID common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Bindings.contract.Call(opts, out, "isDeregistered", _darknodeID)
	return *ret0, err
}

// IsDeregistered is a free data retrieval call binding the contract method 0x8020fc1f.
//
// Solidity: function isDeregistered(_darknodeID address) constant returns(bool)
func (_Bindings *BindingsSession) IsDeregistered(_darknodeID common.Address) (bool, error) {
	return _Bindings.Contract.IsDeregistered(&_Bindings.CallOpts, _darknodeID)
}

// IsDeregistered is a free data retrieval call binding the contract method 0x8020fc1f.
//
// Solidity: function isDeregistered(_darknodeID address) constant returns(bool)
func (_Bindings *BindingsCallerSession) IsDeregistered(_darknodeID common.Address) (bool, error) {
	return _Bindings.Contract.IsDeregistered(&_Bindings.CallOpts, _darknodeID)
}

// IsPendingDeregistration is a free data retrieval call binding the contract method 0x303ee989.
//
// Solidity: function isPendingDeregistration(_darknodeID address) constant returns(bool)
func (_Bindings *BindingsCaller) IsPendingDeregistration(opts *bind.CallOpts, _darknodeID common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Bindings.contract.Call(opts, out, "isPendingDeregistration", _darknodeID)
	return *ret0, err
}

// IsPendingDeregistration is a free data retrieval call binding the contract method 0x303ee989.
//
// Solidity: function isPendingDeregistration(_darknodeID address) constant returns(bool)
func (_Bindings *BindingsSession) IsPendingDeregistration(_darknodeID common.Address) (bool, error) {
	return _Bindings.Contract.IsPendingDeregistration(&_Bindings.CallOpts, _darknodeID)
}

// IsPendingDeregistration is a free data retrieval call binding the contract method 0x303ee989.
//
// Solidity: function isPendingDeregistration(_darknodeID address) constant returns(bool)
func (_Bindings *BindingsCallerSession) IsPendingDeregistration(_darknodeID common.Address) (bool, error) {
	return _Bindings.Contract.IsPendingDeregistration(&_Bindings.CallOpts, _darknodeID)
}

// IsPendingRegistration is a free data retrieval call binding the contract method 0x040fa051.
//
// Solidity: function isPendingRegistration(_darknodeID address) constant returns(bool)
func (_Bindings *BindingsCaller) IsPendingRegistration(opts *bind.CallOpts, _darknodeID common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Bindings.contract.Call(opts, out, "isPendingRegistration", _darknodeID)
	return *ret0, err
}

// IsPendingRegistration is a free data retrieval call binding the contract method 0x040fa051.
//
// Solidity: function isPendingRegistration(_darknodeID address) constant returns(bool)
func (_Bindings *BindingsSession) IsPendingRegistration(_darknodeID common.Address) (bool, error) {
	return _Bindings.Contract.IsPendingRegistration(&_Bindings.CallOpts, _darknodeID)
}

// IsPendingRegistration is a free data retrieval call binding the contract method 0x040fa051.
//
// Solidity: function isPendingRegistration(_darknodeID address) constant returns(bool)
func (_Bindings *BindingsCallerSession) IsPendingRegistration(_darknodeID common.Address) (bool, error) {
	return _Bindings.Contract.IsPendingRegistration(&_Bindings.CallOpts, _darknodeID)
}

// IsRefundable is a free data retrieval call binding the contract method 0x5aebd1cb.
//
// Solidity: function isRefundable(_darknodeID address) constant returns(bool)
func (_Bindings *BindingsCaller) IsRefundable(opts *bind.CallOpts, _darknodeID common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Bindings.contract.Call(opts, out, "isRefundable", _darknodeID)
	return *ret0, err
}

// IsRefundable is a free data retrieval call binding the contract method 0x5aebd1cb.
//
// Solidity: function isRefundable(_darknodeID address) constant returns(bool)
func (_Bindings *BindingsSession) IsRefundable(_darknodeID common.Address) (bool, error) {
	return _Bindings.Contract.IsRefundable(&_Bindings.CallOpts, _darknodeID)
}

// IsRefundable is a free data retrieval call binding the contract method 0x5aebd1cb.
//
// Solidity: function isRefundable(_darknodeID address) constant returns(bool)
func (_Bindings *BindingsCallerSession) IsRefundable(_darknodeID common.Address) (bool, error) {
	return _Bindings.Contract.IsRefundable(&_Bindings.CallOpts, _darknodeID)
}

// IsRefunded is a free data retrieval call binding the contract method 0xffc9152e.
//
// Solidity: function isRefunded(_darknodeID address) constant returns(bool)
func (_Bindings *BindingsCaller) IsRefunded(opts *bind.CallOpts, _darknodeID common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Bindings.contract.Call(opts, out, "isRefunded", _darknodeID)
	return *ret0, err
}

// IsRefunded is a free data retrieval call binding the contract method 0xffc9152e.
//
// Solidity: function isRefunded(_darknodeID address) constant returns(bool)
func (_Bindings *BindingsSession) IsRefunded(_darknodeID common.Address) (bool, error) {
	return _Bindings.Contract.IsRefunded(&_Bindings.CallOpts, _darknodeID)
}

// IsRefunded is a free data retrieval call binding the contract method 0xffc9152e.
//
// Solidity: function isRefunded(_darknodeID address) constant returns(bool)
func (_Bindings *BindingsCallerSession) IsRefunded(_darknodeID common.Address) (bool, error) {
	return _Bindings.Contract.IsRefunded(&_Bindings.CallOpts, _darknodeID)
}

// IsRegistered is a free data retrieval call binding the contract method 0xc3c5a547.
//
// Solidity: function isRegistered(_darknodeID address) constant returns(bool)
func (_Bindings *BindingsCaller) IsRegistered(opts *bind.CallOpts, _darknodeID common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Bindings.contract.Call(opts, out, "isRegistered", _darknodeID)
	return *ret0, err
}

// IsRegistered is a free data retrieval call binding the contract method 0xc3c5a547.
//
// Solidity: function isRegistered(_darknodeID address) constant returns(bool)
func (_Bindings *BindingsSession) IsRegistered(_darknodeID common.Address) (bool, error) {
	return _Bindings.Contract.IsRegistered(&_Bindings.CallOpts, _darknodeID)
}

// IsRegistered is a free data retrieval call binding the contract method 0xc3c5a547.
//
// Solidity: function isRegistered(_darknodeID address) constant returns(bool)
func (_Bindings *BindingsCallerSession) IsRegistered(_darknodeID common.Address) (bool, error) {
	return _Bindings.Contract.IsRegistered(&_Bindings.CallOpts, _darknodeID)
}

// IsRegisteredInPreviousEpoch is a free data retrieval call binding the contract method 0x7be266da.
//
// Solidity: function isRegisteredInPreviousEpoch(_darknodeID address) constant returns(bool)
func (_Bindings *BindingsCaller) IsRegisteredInPreviousEpoch(opts *bind.CallOpts, _darknodeID common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Bindings.contract.Call(opts, out, "isRegisteredInPreviousEpoch", _darknodeID)
	return *ret0, err
}

// IsRegisteredInPreviousEpoch is a free data retrieval call binding the contract method 0x7be266da.
//
// Solidity: function isRegisteredInPreviousEpoch(_darknodeID address) constant returns(bool)
func (_Bindings *BindingsSession) IsRegisteredInPreviousEpoch(_darknodeID common.Address) (bool, error) {
	return _Bindings.Contract.IsRegisteredInPreviousEpoch(&_Bindings.CallOpts, _darknodeID)
}

// IsRegisteredInPreviousEpoch is a free data retrieval call binding the contract method 0x7be266da.
//
// Solidity: function isRegisteredInPreviousEpoch(_darknodeID address) constant returns(bool)
func (_Bindings *BindingsCallerSession) IsRegisteredInPreviousEpoch(_darknodeID common.Address) (bool, error) {
	return _Bindings.Contract.IsRegisteredInPreviousEpoch(&_Bindings.CallOpts, _darknodeID)
}

// MinimumBond is a free data retrieval call binding the contract method 0xaa7517e1.
//
// Solidity: function minimumBond() constant returns(uint256)
func (_Bindings *BindingsCaller) MinimumBond(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Bindings.contract.Call(opts, out, "minimumBond")
	return *ret0, err
}

// MinimumBond is a free data retrieval call binding the contract method 0xaa7517e1.
//
// Solidity: function minimumBond() constant returns(uint256)
func (_Bindings *BindingsSession) MinimumBond() (*big.Int, error) {
	return _Bindings.Contract.MinimumBond(&_Bindings.CallOpts)
}

// MinimumBond is a free data retrieval call binding the contract method 0xaa7517e1.
//
// Solidity: function minimumBond() constant returns(uint256)
func (_Bindings *BindingsCallerSession) MinimumBond() (*big.Int, error) {
	return _Bindings.Contract.MinimumBond(&_Bindings.CallOpts)
}

// MinimumEpochInterval is a free data retrieval call binding the contract method 0x55cacda5.
//
// Solidity: function minimumEpochInterval() constant returns(uint256)
func (_Bindings *BindingsCaller) MinimumEpochInterval(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Bindings.contract.Call(opts, out, "minimumEpochInterval")
	return *ret0, err
}

// MinimumEpochInterval is a free data retrieval call binding the contract method 0x55cacda5.
//
// Solidity: function minimumEpochInterval() constant returns(uint256)
func (_Bindings *BindingsSession) MinimumEpochInterval() (*big.Int, error) {
	return _Bindings.Contract.MinimumEpochInterval(&_Bindings.CallOpts)
}

// MinimumEpochInterval is a free data retrieval call binding the contract method 0x55cacda5.
//
// Solidity: function minimumEpochInterval() constant returns(uint256)
func (_Bindings *BindingsCallerSession) MinimumEpochInterval() (*big.Int, error) {
	return _Bindings.Contract.MinimumEpochInterval(&_Bindings.CallOpts)
}

// MinimumPodSize is a free data retrieval call binding the contract method 0xc7dbc2be.
//
// Solidity: function minimumPodSize() constant returns(uint256)
func (_Bindings *BindingsCaller) MinimumPodSize(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Bindings.contract.Call(opts, out, "minimumPodSize")
	return *ret0, err
}

// MinimumPodSize is a free data retrieval call binding the contract method 0xc7dbc2be.
//
// Solidity: function minimumPodSize() constant returns(uint256)
func (_Bindings *BindingsSession) MinimumPodSize() (*big.Int, error) {
	return _Bindings.Contract.MinimumPodSize(&_Bindings.CallOpts)
}

// MinimumPodSize is a free data retrieval call binding the contract method 0xc7dbc2be.
//
// Solidity: function minimumPodSize() constant returns(uint256)
func (_Bindings *BindingsCallerSession) MinimumPodSize() (*big.Int, error) {
	return _Bindings.Contract.MinimumPodSize(&_Bindings.CallOpts)
}

// NextMinimumBond is a free data retrieval call binding the contract method 0x60a22fe4.
//
// Solidity: function nextMinimumBond() constant returns(uint256)
func (_Bindings *BindingsCaller) NextMinimumBond(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Bindings.contract.Call(opts, out, "nextMinimumBond")
	return *ret0, err
}

// NextMinimumBond is a free data retrieval call binding the contract method 0x60a22fe4.
//
// Solidity: function nextMinimumBond() constant returns(uint256)
func (_Bindings *BindingsSession) NextMinimumBond() (*big.Int, error) {
	return _Bindings.Contract.NextMinimumBond(&_Bindings.CallOpts)
}

// NextMinimumBond is a free data retrieval call binding the contract method 0x60a22fe4.
//
// Solidity: function nextMinimumBond() constant returns(uint256)
func (_Bindings *BindingsCallerSession) NextMinimumBond() (*big.Int, error) {
	return _Bindings.Contract.NextMinimumBond(&_Bindings.CallOpts)
}

// NextMinimumEpochInterval is a free data retrieval call binding the contract method 0x455dc46d.
//
// Solidity: function nextMinimumEpochInterval() constant returns(uint256)
func (_Bindings *BindingsCaller) NextMinimumEpochInterval(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Bindings.contract.Call(opts, out, "nextMinimumEpochInterval")
	return *ret0, err
}

// NextMinimumEpochInterval is a free data retrieval call binding the contract method 0x455dc46d.
//
// Solidity: function nextMinimumEpochInterval() constant returns(uint256)
func (_Bindings *BindingsSession) NextMinimumEpochInterval() (*big.Int, error) {
	return _Bindings.Contract.NextMinimumEpochInterval(&_Bindings.CallOpts)
}

// NextMinimumEpochInterval is a free data retrieval call binding the contract method 0x455dc46d.
//
// Solidity: function nextMinimumEpochInterval() constant returns(uint256)
func (_Bindings *BindingsCallerSession) NextMinimumEpochInterval() (*big.Int, error) {
	return _Bindings.Contract.NextMinimumEpochInterval(&_Bindings.CallOpts)
}

// NextMinimumPodSize is a free data retrieval call binding the contract method 0x702c25ee.
//
// Solidity: function nextMinimumPodSize() constant returns(uint256)
func (_Bindings *BindingsCaller) NextMinimumPodSize(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Bindings.contract.Call(opts, out, "nextMinimumPodSize")
	return *ret0, err
}

// NextMinimumPodSize is a free data retrieval call binding the contract method 0x702c25ee.
//
// Solidity: function nextMinimumPodSize() constant returns(uint256)
func (_Bindings *BindingsSession) NextMinimumPodSize() (*big.Int, error) {
	return _Bindings.Contract.NextMinimumPodSize(&_Bindings.CallOpts)
}

// NextMinimumPodSize is a free data retrieval call binding the contract method 0x702c25ee.
//
// Solidity: function nextMinimumPodSize() constant returns(uint256)
func (_Bindings *BindingsCallerSession) NextMinimumPodSize() (*big.Int, error) {
	return _Bindings.Contract.NextMinimumPodSize(&_Bindings.CallOpts)
}

// NextSlasher is a free data retrieval call binding the contract method 0x21a2ad3a.
//
// Solidity: function nextSlasher() constant returns(address)
func (_Bindings *BindingsCaller) NextSlasher(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Bindings.contract.Call(opts, out, "nextSlasher")
	return *ret0, err
}

// NextSlasher is a free data retrieval call binding the contract method 0x21a2ad3a.
//
// Solidity: function nextSlasher() constant returns(address)
func (_Bindings *BindingsSession) NextSlasher() (common.Address, error) {
	return _Bindings.Contract.NextSlasher(&_Bindings.CallOpts)
}

// NextSlasher is a free data retrieval call binding the contract method 0x21a2ad3a.
//
// Solidity: function nextSlasher() constant returns(address)
func (_Bindings *BindingsCallerSession) NextSlasher() (common.Address, error) {
	return _Bindings.Contract.NextSlasher(&_Bindings.CallOpts)
}

// NumDarknodes is a free data retrieval call binding the contract method 0x1460e603.
//
// Solidity: function numDarknodes() constant returns(uint256)
func (_Bindings *BindingsCaller) NumDarknodes(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Bindings.contract.Call(opts, out, "numDarknodes")
	return *ret0, err
}

// NumDarknodes is a free data retrieval call binding the contract method 0x1460e603.
//
// Solidity: function numDarknodes() constant returns(uint256)
func (_Bindings *BindingsSession) NumDarknodes() (*big.Int, error) {
	return _Bindings.Contract.NumDarknodes(&_Bindings.CallOpts)
}

// NumDarknodes is a free data retrieval call binding the contract method 0x1460e603.
//
// Solidity: function numDarknodes() constant returns(uint256)
func (_Bindings *BindingsCallerSession) NumDarknodes() (*big.Int, error) {
	return _Bindings.Contract.NumDarknodes(&_Bindings.CallOpts)
}

// NumDarknodesNextEpoch is a free data retrieval call binding the contract method 0x0847e9fa.
//
// Solidity: function numDarknodesNextEpoch() constant returns(uint256)
func (_Bindings *BindingsCaller) NumDarknodesNextEpoch(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Bindings.contract.Call(opts, out, "numDarknodesNextEpoch")
	return *ret0, err
}

// NumDarknodesNextEpoch is a free data retrieval call binding the contract method 0x0847e9fa.
//
// Solidity: function numDarknodesNextEpoch() constant returns(uint256)
func (_Bindings *BindingsSession) NumDarknodesNextEpoch() (*big.Int, error) {
	return _Bindings.Contract.NumDarknodesNextEpoch(&_Bindings.CallOpts)
}

// NumDarknodesNextEpoch is a free data retrieval call binding the contract method 0x0847e9fa.
//
// Solidity: function numDarknodesNextEpoch() constant returns(uint256)
func (_Bindings *BindingsCallerSession) NumDarknodesNextEpoch() (*big.Int, error) {
	return _Bindings.Contract.NumDarknodesNextEpoch(&_Bindings.CallOpts)
}

// NumDarknodesPreviousEpoch is a free data retrieval call binding the contract method 0x71740d16.
//
// Solidity: function numDarknodesPreviousEpoch() constant returns(uint256)
func (_Bindings *BindingsCaller) NumDarknodesPreviousEpoch(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Bindings.contract.Call(opts, out, "numDarknodesPreviousEpoch")
	return *ret0, err
}

// NumDarknodesPreviousEpoch is a free data retrieval call binding the contract method 0x71740d16.
//
// Solidity: function numDarknodesPreviousEpoch() constant returns(uint256)
func (_Bindings *BindingsSession) NumDarknodesPreviousEpoch() (*big.Int, error) {
	return _Bindings.Contract.NumDarknodesPreviousEpoch(&_Bindings.CallOpts)
}

// NumDarknodesPreviousEpoch is a free data retrieval call binding the contract method 0x71740d16.
//
// Solidity: function numDarknodesPreviousEpoch() constant returns(uint256)
func (_Bindings *BindingsCallerSession) NumDarknodesPreviousEpoch() (*big.Int, error) {
	return _Bindings.Contract.NumDarknodesPreviousEpoch(&_Bindings.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Bindings *BindingsCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Bindings.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Bindings *BindingsSession) Owner() (common.Address, error) {
	return _Bindings.Contract.Owner(&_Bindings.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Bindings *BindingsCallerSession) Owner() (common.Address, error) {
	return _Bindings.Contract.Owner(&_Bindings.CallOpts)
}

// PreviousEpoch is a free data retrieval call binding the contract method 0x5cdaab48.
//
// Solidity: function previousEpoch() constant returns(epochhash uint256, blocknumber uint256)
func (_Bindings *BindingsCaller) PreviousEpoch(opts *bind.CallOpts) (struct {
	Epochhash   *big.Int
	Blocknumber *big.Int
}, error) {
	ret := new(struct {
		Epochhash   *big.Int
		Blocknumber *big.Int
	})
	out := ret
	err := _Bindings.contract.Call(opts, out, "previousEpoch")
	return *ret, err
}

// PreviousEpoch is a free data retrieval call binding the contract method 0x5cdaab48.
//
// Solidity: function previousEpoch() constant returns(epochhash uint256, blocknumber uint256)
func (_Bindings *BindingsSession) PreviousEpoch() (struct {
	Epochhash   *big.Int
	Blocknumber *big.Int
}, error) {
	return _Bindings.Contract.PreviousEpoch(&_Bindings.CallOpts)
}

// PreviousEpoch is a free data retrieval call binding the contract method 0x5cdaab48.
//
// Solidity: function previousEpoch() constant returns(epochhash uint256, blocknumber uint256)
func (_Bindings *BindingsCallerSession) PreviousEpoch() (struct {
	Epochhash   *big.Int
	Blocknumber *big.Int
}, error) {
	return _Bindings.Contract.PreviousEpoch(&_Bindings.CallOpts)
}

// Ren is a free data retrieval call binding the contract method 0x8a9b4067.
//
// Solidity: function ren() constant returns(address)
func (_Bindings *BindingsCaller) Ren(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Bindings.contract.Call(opts, out, "ren")
	return *ret0, err
}

// Ren is a free data retrieval call binding the contract method 0x8a9b4067.
//
// Solidity: function ren() constant returns(address)
func (_Bindings *BindingsSession) Ren() (common.Address, error) {
	return _Bindings.Contract.Ren(&_Bindings.CallOpts)
}

// Ren is a free data retrieval call binding the contract method 0x8a9b4067.
//
// Solidity: function ren() constant returns(address)
func (_Bindings *BindingsCallerSession) Ren() (common.Address, error) {
	return _Bindings.Contract.Ren(&_Bindings.CallOpts)
}

// Slasher is a free data retrieval call binding the contract method 0xb1344271.
//
// Solidity: function slasher() constant returns(address)
func (_Bindings *BindingsCaller) Slasher(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Bindings.contract.Call(opts, out, "slasher")
	return *ret0, err
}

// Slasher is a free data retrieval call binding the contract method 0xb1344271.
//
// Solidity: function slasher() constant returns(address)
func (_Bindings *BindingsSession) Slasher() (common.Address, error) {
	return _Bindings.Contract.Slasher(&_Bindings.CallOpts)
}

// Slasher is a free data retrieval call binding the contract method 0xb1344271.
//
// Solidity: function slasher() constant returns(address)
func (_Bindings *BindingsCallerSession) Slasher() (common.Address, error) {
	return _Bindings.Contract.Slasher(&_Bindings.CallOpts)
}

// Store is a free data retrieval call binding the contract method 0x975057e7.
//
// Solidity: function store() constant returns(address)
func (_Bindings *BindingsCaller) Store(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Bindings.contract.Call(opts, out, "store")
	return *ret0, err
}

// Store is a free data retrieval call binding the contract method 0x975057e7.
//
// Solidity: function store() constant returns(address)
func (_Bindings *BindingsSession) Store() (common.Address, error) {
	return _Bindings.Contract.Store(&_Bindings.CallOpts)
}

// Store is a free data retrieval call binding the contract method 0x975057e7.
//
// Solidity: function store() constant returns(address)
func (_Bindings *BindingsCallerSession) Store() (common.Address, error) {
	return _Bindings.Contract.Store(&_Bindings.CallOpts)
}

// Deregister is a paid mutator transaction binding the contract method 0x84ac33ec.
//
// Solidity: function deregister(_darknodeID address) returns()
func (_Bindings *BindingsTransactor) Deregister(opts *bind.TransactOpts, _darknodeID common.Address) (*types.Transaction, error) {
	return _Bindings.contract.Transact(opts, "deregister", _darknodeID)
}

// Deregister is a paid mutator transaction binding the contract method 0x84ac33ec.
//
// Solidity: function deregister(_darknodeID address) returns()
func (_Bindings *BindingsSession) Deregister(_darknodeID common.Address) (*types.Transaction, error) {
	return _Bindings.Contract.Deregister(&_Bindings.TransactOpts, _darknodeID)
}

// Deregister is a paid mutator transaction binding the contract method 0x84ac33ec.
//
// Solidity: function deregister(_darknodeID address) returns()
func (_Bindings *BindingsTransactorSession) Deregister(_darknodeID common.Address) (*types.Transaction, error) {
	return _Bindings.Contract.Deregister(&_Bindings.TransactOpts, _darknodeID)
}

// Epoch is a paid mutator transaction binding the contract method 0x900cf0cf.
//
// Solidity: function epoch() returns()
func (_Bindings *BindingsTransactor) Epoch(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bindings.contract.Transact(opts, "epoch")
}

// Epoch is a paid mutator transaction binding the contract method 0x900cf0cf.
//
// Solidity: function epoch() returns()
func (_Bindings *BindingsSession) Epoch() (*types.Transaction, error) {
	return _Bindings.Contract.Epoch(&_Bindings.TransactOpts)
}

// Epoch is a paid mutator transaction binding the contract method 0x900cf0cf.
//
// Solidity: function epoch() returns()
func (_Bindings *BindingsTransactorSession) Epoch() (*types.Transaction, error) {
	return _Bindings.Contract.Epoch(&_Bindings.TransactOpts)
}

// Refund is a paid mutator transaction binding the contract method 0xfa89401a.
//
// Solidity: function refund(_darknodeID address) returns()
func (_Bindings *BindingsTransactor) Refund(opts *bind.TransactOpts, _darknodeID common.Address) (*types.Transaction, error) {
	return _Bindings.contract.Transact(opts, "refund", _darknodeID)
}

// Refund is a paid mutator transaction binding the contract method 0xfa89401a.
//
// Solidity: function refund(_darknodeID address) returns()
func (_Bindings *BindingsSession) Refund(_darknodeID common.Address) (*types.Transaction, error) {
	return _Bindings.Contract.Refund(&_Bindings.TransactOpts, _darknodeID)
}

// Refund is a paid mutator transaction binding the contract method 0xfa89401a.
//
// Solidity: function refund(_darknodeID address) returns()
func (_Bindings *BindingsTransactorSession) Refund(_darknodeID common.Address) (*types.Transaction, error) {
	return _Bindings.Contract.Refund(&_Bindings.TransactOpts, _darknodeID)
}

// Register is a paid mutator transaction binding the contract method 0x0aeb6b40.
//
// Solidity: function register(_darknodeID address, _publicKey bytes, _bond uint256) returns()
func (_Bindings *BindingsTransactor) Register(opts *bind.TransactOpts, _darknodeID common.Address, _publicKey []byte, _bond *big.Int) (*types.Transaction, error) {
	return _Bindings.contract.Transact(opts, "register", _darknodeID, _publicKey, _bond)
}

// Register is a paid mutator transaction binding the contract method 0x0aeb6b40.
//
// Solidity: function register(_darknodeID address, _publicKey bytes, _bond uint256) returns()
func (_Bindings *BindingsSession) Register(_darknodeID common.Address, _publicKey []byte, _bond *big.Int) (*types.Transaction, error) {
	return _Bindings.Contract.Register(&_Bindings.TransactOpts, _darknodeID, _publicKey, _bond)
}

// Register is a paid mutator transaction binding the contract method 0x0aeb6b40.
//
// Solidity: function register(_darknodeID address, _publicKey bytes, _bond uint256) returns()
func (_Bindings *BindingsTransactorSession) Register(_darknodeID common.Address, _publicKey []byte, _bond *big.Int) (*types.Transaction, error) {
	return _Bindings.Contract.Register(&_Bindings.TransactOpts, _darknodeID, _publicKey, _bond)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Bindings *BindingsTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bindings.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Bindings *BindingsSession) RenounceOwnership() (*types.Transaction, error) {
	return _Bindings.Contract.RenounceOwnership(&_Bindings.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Bindings *BindingsTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Bindings.Contract.RenounceOwnership(&_Bindings.TransactOpts)
}

// Slash is a paid mutator transaction binding the contract method 0x563bf264.
//
// Solidity: function slash(_prover address, _challenger1 address, _challenger2 address) returns()
func (_Bindings *BindingsTransactor) Slash(opts *bind.TransactOpts, _prover common.Address, _challenger1 common.Address, _challenger2 common.Address) (*types.Transaction, error) {
	return _Bindings.contract.Transact(opts, "slash", _prover, _challenger1, _challenger2)
}

// Slash is a paid mutator transaction binding the contract method 0x563bf264.
//
// Solidity: function slash(_prover address, _challenger1 address, _challenger2 address) returns()
func (_Bindings *BindingsSession) Slash(_prover common.Address, _challenger1 common.Address, _challenger2 common.Address) (*types.Transaction, error) {
	return _Bindings.Contract.Slash(&_Bindings.TransactOpts, _prover, _challenger1, _challenger2)
}

// Slash is a paid mutator transaction binding the contract method 0x563bf264.
//
// Solidity: function slash(_prover address, _challenger1 address, _challenger2 address) returns()
func (_Bindings *BindingsTransactorSession) Slash(_prover common.Address, _challenger1 common.Address, _challenger2 common.Address) (*types.Transaction, error) {
	return _Bindings.Contract.Slash(&_Bindings.TransactOpts, _prover, _challenger1, _challenger2)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_Bindings *BindingsTransactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _Bindings.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_Bindings *BindingsSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _Bindings.Contract.TransferOwnership(&_Bindings.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_Bindings *BindingsTransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _Bindings.Contract.TransferOwnership(&_Bindings.TransactOpts, _newOwner)
}

// TransferStoreOwnership is a paid mutator transaction binding the contract method 0xc2250a99.
//
// Solidity: function transferStoreOwnership(_newOwner address) returns()
func (_Bindings *BindingsTransactor) TransferStoreOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _Bindings.contract.Transact(opts, "transferStoreOwnership", _newOwner)
}

// TransferStoreOwnership is a paid mutator transaction binding the contract method 0xc2250a99.
//
// Solidity: function transferStoreOwnership(_newOwner address) returns()
func (_Bindings *BindingsSession) TransferStoreOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _Bindings.Contract.TransferStoreOwnership(&_Bindings.TransactOpts, _newOwner)
}

// TransferStoreOwnership is a paid mutator transaction binding the contract method 0xc2250a99.
//
// Solidity: function transferStoreOwnership(_newOwner address) returns()
func (_Bindings *BindingsTransactorSession) TransferStoreOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _Bindings.Contract.TransferStoreOwnership(&_Bindings.TransactOpts, _newOwner)
}

// UpdateMinimumBond is a paid mutator transaction binding the contract method 0x0ff9aafe.
//
// Solidity: function updateMinimumBond(_nextMinimumBond uint256) returns()
func (_Bindings *BindingsTransactor) UpdateMinimumBond(opts *bind.TransactOpts, _nextMinimumBond *big.Int) (*types.Transaction, error) {
	return _Bindings.contract.Transact(opts, "updateMinimumBond", _nextMinimumBond)
}

// UpdateMinimumBond is a paid mutator transaction binding the contract method 0x0ff9aafe.
//
// Solidity: function updateMinimumBond(_nextMinimumBond uint256) returns()
func (_Bindings *BindingsSession) UpdateMinimumBond(_nextMinimumBond *big.Int) (*types.Transaction, error) {
	return _Bindings.Contract.UpdateMinimumBond(&_Bindings.TransactOpts, _nextMinimumBond)
}

// UpdateMinimumBond is a paid mutator transaction binding the contract method 0x0ff9aafe.
//
// Solidity: function updateMinimumBond(_nextMinimumBond uint256) returns()
func (_Bindings *BindingsTransactorSession) UpdateMinimumBond(_nextMinimumBond *big.Int) (*types.Transaction, error) {
	return _Bindings.Contract.UpdateMinimumBond(&_Bindings.TransactOpts, _nextMinimumBond)
}

// UpdateMinimumEpochInterval is a paid mutator transaction binding the contract method 0x63b851b9.
//
// Solidity: function updateMinimumEpochInterval(_nextMinimumEpochInterval uint256) returns()
func (_Bindings *BindingsTransactor) UpdateMinimumEpochInterval(opts *bind.TransactOpts, _nextMinimumEpochInterval *big.Int) (*types.Transaction, error) {
	return _Bindings.contract.Transact(opts, "updateMinimumEpochInterval", _nextMinimumEpochInterval)
}

// UpdateMinimumEpochInterval is a paid mutator transaction binding the contract method 0x63b851b9.
//
// Solidity: function updateMinimumEpochInterval(_nextMinimumEpochInterval uint256) returns()
func (_Bindings *BindingsSession) UpdateMinimumEpochInterval(_nextMinimumEpochInterval *big.Int) (*types.Transaction, error) {
	return _Bindings.Contract.UpdateMinimumEpochInterval(&_Bindings.TransactOpts, _nextMinimumEpochInterval)
}

// UpdateMinimumEpochInterval is a paid mutator transaction binding the contract method 0x63b851b9.
//
// Solidity: function updateMinimumEpochInterval(_nextMinimumEpochInterval uint256) returns()
func (_Bindings *BindingsTransactorSession) UpdateMinimumEpochInterval(_nextMinimumEpochInterval *big.Int) (*types.Transaction, error) {
	return _Bindings.Contract.UpdateMinimumEpochInterval(&_Bindings.TransactOpts, _nextMinimumEpochInterval)
}

// UpdateMinimumPodSize is a paid mutator transaction binding the contract method 0x80a0c461.
//
// Solidity: function updateMinimumPodSize(_nextMinimumPodSize uint256) returns()
func (_Bindings *BindingsTransactor) UpdateMinimumPodSize(opts *bind.TransactOpts, _nextMinimumPodSize *big.Int) (*types.Transaction, error) {
	return _Bindings.contract.Transact(opts, "updateMinimumPodSize", _nextMinimumPodSize)
}

// UpdateMinimumPodSize is a paid mutator transaction binding the contract method 0x80a0c461.
//
// Solidity: function updateMinimumPodSize(_nextMinimumPodSize uint256) returns()
func (_Bindings *BindingsSession) UpdateMinimumPodSize(_nextMinimumPodSize *big.Int) (*types.Transaction, error) {
	return _Bindings.Contract.UpdateMinimumPodSize(&_Bindings.TransactOpts, _nextMinimumPodSize)
}

// UpdateMinimumPodSize is a paid mutator transaction binding the contract method 0x80a0c461.
//
// Solidity: function updateMinimumPodSize(_nextMinimumPodSize uint256) returns()
func (_Bindings *BindingsTransactorSession) UpdateMinimumPodSize(_nextMinimumPodSize *big.Int) (*types.Transaction, error) {
	return _Bindings.Contract.UpdateMinimumPodSize(&_Bindings.TransactOpts, _nextMinimumPodSize)
}

// UpdateSlasher is a paid mutator transaction binding the contract method 0xb3139d38.
//
// Solidity: function updateSlasher(_slasher address) returns()
func (_Bindings *BindingsTransactor) UpdateSlasher(opts *bind.TransactOpts, _slasher common.Address) (*types.Transaction, error) {
	return _Bindings.contract.Transact(opts, "updateSlasher", _slasher)
}

// UpdateSlasher is a paid mutator transaction binding the contract method 0xb3139d38.
//
// Solidity: function updateSlasher(_slasher address) returns()
func (_Bindings *BindingsSession) UpdateSlasher(_slasher common.Address) (*types.Transaction, error) {
	return _Bindings.Contract.UpdateSlasher(&_Bindings.TransactOpts, _slasher)
}

// UpdateSlasher is a paid mutator transaction binding the contract method 0xb3139d38.
//
// Solidity: function updateSlasher(_slasher address) returns()
func (_Bindings *BindingsTransactorSession) UpdateSlasher(_slasher common.Address) (*types.Transaction, error) {
	return _Bindings.Contract.UpdateSlasher(&_Bindings.TransactOpts, _slasher)
}

// BindingsLogDarknodeDeregisteredIterator is returned from FilterLogDarknodeDeregistered and is used to iterate over the raw logs and unpacked data for LogDarknodeDeregistered events raised by the Bindings contract.
type BindingsLogDarknodeDeregisteredIterator struct {
	Event *BindingsLogDarknodeDeregistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BindingsLogDarknodeDeregisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BindingsLogDarknodeDeregistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BindingsLogDarknodeDeregistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BindingsLogDarknodeDeregisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BindingsLogDarknodeDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BindingsLogDarknodeDeregistered represents a LogDarknodeDeregistered event raised by the Bindings contract.
type BindingsLogDarknodeDeregistered struct {
	DarknodeID common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogDarknodeDeregistered is a free log retrieval operation binding the contract event 0x2dc89de5703d2c341a22ebfc7c4d3f197e5e1f0c19bc2e1135f387163cb927e4.
//
// Solidity: event LogDarknodeDeregistered(_darknodeID address)
func (_Bindings *BindingsFilterer) FilterLogDarknodeDeregistered(opts *bind.FilterOpts) (*BindingsLogDarknodeDeregisteredIterator, error) {

	logs, sub, err := _Bindings.contract.FilterLogs(opts, "LogDarknodeDeregistered")
	if err != nil {
		return nil, err
	}
	return &BindingsLogDarknodeDeregisteredIterator{contract: _Bindings.contract, event: "LogDarknodeDeregistered", logs: logs, sub: sub}, nil
}

// WatchLogDarknodeDeregistered is a free log subscription operation binding the contract event 0x2dc89de5703d2c341a22ebfc7c4d3f197e5e1f0c19bc2e1135f387163cb927e4.
//
// Solidity: event LogDarknodeDeregistered(_darknodeID address)
func (_Bindings *BindingsFilterer) WatchLogDarknodeDeregistered(opts *bind.WatchOpts, sink chan<- *BindingsLogDarknodeDeregistered) (event.Subscription, error) {

	logs, sub, err := _Bindings.contract.WatchLogs(opts, "LogDarknodeDeregistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BindingsLogDarknodeDeregistered)
				if err := _Bindings.contract.UnpackLog(event, "LogDarknodeDeregistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// BindingsLogDarknodeOwnerRefundedIterator is returned from FilterLogDarknodeOwnerRefunded and is used to iterate over the raw logs and unpacked data for LogDarknodeOwnerRefunded events raised by the Bindings contract.
type BindingsLogDarknodeOwnerRefundedIterator struct {
	Event *BindingsLogDarknodeOwnerRefunded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BindingsLogDarknodeOwnerRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BindingsLogDarknodeOwnerRefunded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BindingsLogDarknodeOwnerRefunded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BindingsLogDarknodeOwnerRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BindingsLogDarknodeOwnerRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BindingsLogDarknodeOwnerRefunded represents a LogDarknodeOwnerRefunded event raised by the Bindings contract.
type BindingsLogDarknodeOwnerRefunded struct {
	Owner  common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterLogDarknodeOwnerRefunded is a free log retrieval operation binding the contract event 0x96ab9e56c79eee4a72db6e2879cbfbecdba5c65b411f4861824e66b89df19764.
//
// Solidity: event LogDarknodeOwnerRefunded(_owner address, _amount uint256)
func (_Bindings *BindingsFilterer) FilterLogDarknodeOwnerRefunded(opts *bind.FilterOpts) (*BindingsLogDarknodeOwnerRefundedIterator, error) {

	logs, sub, err := _Bindings.contract.FilterLogs(opts, "LogDarknodeOwnerRefunded")
	if err != nil {
		return nil, err
	}
	return &BindingsLogDarknodeOwnerRefundedIterator{contract: _Bindings.contract, event: "LogDarknodeOwnerRefunded", logs: logs, sub: sub}, nil
}

// WatchLogDarknodeOwnerRefunded is a free log subscription operation binding the contract event 0x96ab9e56c79eee4a72db6e2879cbfbecdba5c65b411f4861824e66b89df19764.
//
// Solidity: event LogDarknodeOwnerRefunded(_owner address, _amount uint256)
func (_Bindings *BindingsFilterer) WatchLogDarknodeOwnerRefunded(opts *bind.WatchOpts, sink chan<- *BindingsLogDarknodeOwnerRefunded) (event.Subscription, error) {

	logs, sub, err := _Bindings.contract.WatchLogs(opts, "LogDarknodeOwnerRefunded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BindingsLogDarknodeOwnerRefunded)
				if err := _Bindings.contract.UnpackLog(event, "LogDarknodeOwnerRefunded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// BindingsLogDarknodeRegisteredIterator is returned from FilterLogDarknodeRegistered and is used to iterate over the raw logs and unpacked data for LogDarknodeRegistered events raised by the Bindings contract.
type BindingsLogDarknodeRegisteredIterator struct {
	Event *BindingsLogDarknodeRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BindingsLogDarknodeRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BindingsLogDarknodeRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BindingsLogDarknodeRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BindingsLogDarknodeRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BindingsLogDarknodeRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BindingsLogDarknodeRegistered represents a LogDarknodeRegistered event raised by the Bindings contract.
type BindingsLogDarknodeRegistered struct {
	DarknodeID common.Address
	Bond       *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogDarknodeRegistered is a free log retrieval operation binding the contract event 0xd2819ba4c736158371edf0be38fd8d1fc435609832e392f118c4c79160e5bd7b.
//
// Solidity: event LogDarknodeRegistered(_darknodeID address, _bond uint256)
func (_Bindings *BindingsFilterer) FilterLogDarknodeRegistered(opts *bind.FilterOpts) (*BindingsLogDarknodeRegisteredIterator, error) {

	logs, sub, err := _Bindings.contract.FilterLogs(opts, "LogDarknodeRegistered")
	if err != nil {
		return nil, err
	}
	return &BindingsLogDarknodeRegisteredIterator{contract: _Bindings.contract, event: "LogDarknodeRegistered", logs: logs, sub: sub}, nil
}

// WatchLogDarknodeRegistered is a free log subscription operation binding the contract event 0xd2819ba4c736158371edf0be38fd8d1fc435609832e392f118c4c79160e5bd7b.
//
// Solidity: event LogDarknodeRegistered(_darknodeID address, _bond uint256)
func (_Bindings *BindingsFilterer) WatchLogDarknodeRegistered(opts *bind.WatchOpts, sink chan<- *BindingsLogDarknodeRegistered) (event.Subscription, error) {

	logs, sub, err := _Bindings.contract.WatchLogs(opts, "LogDarknodeRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BindingsLogDarknodeRegistered)
				if err := _Bindings.contract.UnpackLog(event, "LogDarknodeRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// BindingsLogMinimumBondUpdatedIterator is returned from FilterLogMinimumBondUpdated and is used to iterate over the raw logs and unpacked data for LogMinimumBondUpdated events raised by the Bindings contract.
type BindingsLogMinimumBondUpdatedIterator struct {
	Event *BindingsLogMinimumBondUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BindingsLogMinimumBondUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BindingsLogMinimumBondUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BindingsLogMinimumBondUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BindingsLogMinimumBondUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BindingsLogMinimumBondUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BindingsLogMinimumBondUpdated represents a LogMinimumBondUpdated event raised by the Bindings contract.
type BindingsLogMinimumBondUpdated struct {
	PreviousMinimumBond *big.Int
	NextMinimumBond     *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterLogMinimumBondUpdated is a free log retrieval operation binding the contract event 0x7c77c94944e9e4e5b0d46f1297127d060020792687cd743401d782346c68f655.
//
// Solidity: event LogMinimumBondUpdated(previousMinimumBond uint256, nextMinimumBond uint256)
func (_Bindings *BindingsFilterer) FilterLogMinimumBondUpdated(opts *bind.FilterOpts) (*BindingsLogMinimumBondUpdatedIterator, error) {

	logs, sub, err := _Bindings.contract.FilterLogs(opts, "LogMinimumBondUpdated")
	if err != nil {
		return nil, err
	}
	return &BindingsLogMinimumBondUpdatedIterator{contract: _Bindings.contract, event: "LogMinimumBondUpdated", logs: logs, sub: sub}, nil
}

// WatchLogMinimumBondUpdated is a free log subscription operation binding the contract event 0x7c77c94944e9e4e5b0d46f1297127d060020792687cd743401d782346c68f655.
//
// Solidity: event LogMinimumBondUpdated(previousMinimumBond uint256, nextMinimumBond uint256)
func (_Bindings *BindingsFilterer) WatchLogMinimumBondUpdated(opts *bind.WatchOpts, sink chan<- *BindingsLogMinimumBondUpdated) (event.Subscription, error) {

	logs, sub, err := _Bindings.contract.WatchLogs(opts, "LogMinimumBondUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BindingsLogMinimumBondUpdated)
				if err := _Bindings.contract.UnpackLog(event, "LogMinimumBondUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// BindingsLogMinimumEpochIntervalUpdatedIterator is returned from FilterLogMinimumEpochIntervalUpdated and is used to iterate over the raw logs and unpacked data for LogMinimumEpochIntervalUpdated events raised by the Bindings contract.
type BindingsLogMinimumEpochIntervalUpdatedIterator struct {
	Event *BindingsLogMinimumEpochIntervalUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BindingsLogMinimumEpochIntervalUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BindingsLogMinimumEpochIntervalUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BindingsLogMinimumEpochIntervalUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BindingsLogMinimumEpochIntervalUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BindingsLogMinimumEpochIntervalUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BindingsLogMinimumEpochIntervalUpdated represents a LogMinimumEpochIntervalUpdated event raised by the Bindings contract.
type BindingsLogMinimumEpochIntervalUpdated struct {
	PreviousMinimumEpochInterval *big.Int
	NextMinimumEpochInterval     *big.Int
	Raw                          types.Log // Blockchain specific contextual infos
}

// FilterLogMinimumEpochIntervalUpdated is a free log retrieval operation binding the contract event 0xb218cde2730b79a0667ddf869466ee66a12ef56fe65fa4986a590f8a7108c9de.
//
// Solidity: event LogMinimumEpochIntervalUpdated(previousMinimumEpochInterval uint256, nextMinimumEpochInterval uint256)
func (_Bindings *BindingsFilterer) FilterLogMinimumEpochIntervalUpdated(opts *bind.FilterOpts) (*BindingsLogMinimumEpochIntervalUpdatedIterator, error) {

	logs, sub, err := _Bindings.contract.FilterLogs(opts, "LogMinimumEpochIntervalUpdated")
	if err != nil {
		return nil, err
	}
	return &BindingsLogMinimumEpochIntervalUpdatedIterator{contract: _Bindings.contract, event: "LogMinimumEpochIntervalUpdated", logs: logs, sub: sub}, nil
}

// WatchLogMinimumEpochIntervalUpdated is a free log subscription operation binding the contract event 0xb218cde2730b79a0667ddf869466ee66a12ef56fe65fa4986a590f8a7108c9de.
//
// Solidity: event LogMinimumEpochIntervalUpdated(previousMinimumEpochInterval uint256, nextMinimumEpochInterval uint256)
func (_Bindings *BindingsFilterer) WatchLogMinimumEpochIntervalUpdated(opts *bind.WatchOpts, sink chan<- *BindingsLogMinimumEpochIntervalUpdated) (event.Subscription, error) {

	logs, sub, err := _Bindings.contract.WatchLogs(opts, "LogMinimumEpochIntervalUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BindingsLogMinimumEpochIntervalUpdated)
				if err := _Bindings.contract.UnpackLog(event, "LogMinimumEpochIntervalUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// BindingsLogMinimumPodSizeUpdatedIterator is returned from FilterLogMinimumPodSizeUpdated and is used to iterate over the raw logs and unpacked data for LogMinimumPodSizeUpdated events raised by the Bindings contract.
type BindingsLogMinimumPodSizeUpdatedIterator struct {
	Event *BindingsLogMinimumPodSizeUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BindingsLogMinimumPodSizeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BindingsLogMinimumPodSizeUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BindingsLogMinimumPodSizeUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BindingsLogMinimumPodSizeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BindingsLogMinimumPodSizeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BindingsLogMinimumPodSizeUpdated represents a LogMinimumPodSizeUpdated event raised by the Bindings contract.
type BindingsLogMinimumPodSizeUpdated struct {
	PreviousMinimumPodSize *big.Int
	NextMinimumPodSize     *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterLogMinimumPodSizeUpdated is a free log retrieval operation binding the contract event 0x6d520e46e5714982ddf8cb6216bcb3e1c1d5b79d337afc305335f819394f5d6a.
//
// Solidity: event LogMinimumPodSizeUpdated(previousMinimumPodSize uint256, nextMinimumPodSize uint256)
func (_Bindings *BindingsFilterer) FilterLogMinimumPodSizeUpdated(opts *bind.FilterOpts) (*BindingsLogMinimumPodSizeUpdatedIterator, error) {

	logs, sub, err := _Bindings.contract.FilterLogs(opts, "LogMinimumPodSizeUpdated")
	if err != nil {
		return nil, err
	}
	return &BindingsLogMinimumPodSizeUpdatedIterator{contract: _Bindings.contract, event: "LogMinimumPodSizeUpdated", logs: logs, sub: sub}, nil
}

// WatchLogMinimumPodSizeUpdated is a free log subscription operation binding the contract event 0x6d520e46e5714982ddf8cb6216bcb3e1c1d5b79d337afc305335f819394f5d6a.
//
// Solidity: event LogMinimumPodSizeUpdated(previousMinimumPodSize uint256, nextMinimumPodSize uint256)
func (_Bindings *BindingsFilterer) WatchLogMinimumPodSizeUpdated(opts *bind.WatchOpts, sink chan<- *BindingsLogMinimumPodSizeUpdated) (event.Subscription, error) {

	logs, sub, err := _Bindings.contract.WatchLogs(opts, "LogMinimumPodSizeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BindingsLogMinimumPodSizeUpdated)
				if err := _Bindings.contract.UnpackLog(event, "LogMinimumPodSizeUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// BindingsLogNewEpochIterator is returned from FilterLogNewEpoch and is used to iterate over the raw logs and unpacked data for LogNewEpoch events raised by the Bindings contract.
type BindingsLogNewEpochIterator struct {
	Event *BindingsLogNewEpoch // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BindingsLogNewEpochIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BindingsLogNewEpoch)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BindingsLogNewEpoch)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BindingsLogNewEpochIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BindingsLogNewEpochIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BindingsLogNewEpoch represents a LogNewEpoch event raised by the Bindings contract.
type BindingsLogNewEpoch struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterLogNewEpoch is a free log retrieval operation binding the contract event 0xeff7e281fe3b4211ed1f0a5e6419bcc40f4552974f771357e66926421f0a58e8.
//
// Solidity: event LogNewEpoch()
func (_Bindings *BindingsFilterer) FilterLogNewEpoch(opts *bind.FilterOpts) (*BindingsLogNewEpochIterator, error) {

	logs, sub, err := _Bindings.contract.FilterLogs(opts, "LogNewEpoch")
	if err != nil {
		return nil, err
	}
	return &BindingsLogNewEpochIterator{contract: _Bindings.contract, event: "LogNewEpoch", logs: logs, sub: sub}, nil
}

// WatchLogNewEpoch is a free log subscription operation binding the contract event 0xeff7e281fe3b4211ed1f0a5e6419bcc40f4552974f771357e66926421f0a58e8.
//
// Solidity: event LogNewEpoch()
func (_Bindings *BindingsFilterer) WatchLogNewEpoch(opts *bind.WatchOpts, sink chan<- *BindingsLogNewEpoch) (event.Subscription, error) {

	logs, sub, err := _Bindings.contract.WatchLogs(opts, "LogNewEpoch")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BindingsLogNewEpoch)
				if err := _Bindings.contract.UnpackLog(event, "LogNewEpoch", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// BindingsLogSlasherUpdatedIterator is returned from FilterLogSlasherUpdated and is used to iterate over the raw logs and unpacked data for LogSlasherUpdated events raised by the Bindings contract.
type BindingsLogSlasherUpdatedIterator struct {
	Event *BindingsLogSlasherUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BindingsLogSlasherUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BindingsLogSlasherUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BindingsLogSlasherUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BindingsLogSlasherUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BindingsLogSlasherUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BindingsLogSlasherUpdated represents a LogSlasherUpdated event raised by the Bindings contract.
type BindingsLogSlasherUpdated struct {
	PreviousSlasher common.Address
	NextSlasher     common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterLogSlasherUpdated is a free log retrieval operation binding the contract event 0x933228a1c3ba8fadd3ce47a9db5b898be647f89af99ba7c1b9a655f59ea306c8.
//
// Solidity: event LogSlasherUpdated(previousSlasher address, nextSlasher address)
func (_Bindings *BindingsFilterer) FilterLogSlasherUpdated(opts *bind.FilterOpts) (*BindingsLogSlasherUpdatedIterator, error) {

	logs, sub, err := _Bindings.contract.FilterLogs(opts, "LogSlasherUpdated")
	if err != nil {
		return nil, err
	}
	return &BindingsLogSlasherUpdatedIterator{contract: _Bindings.contract, event: "LogSlasherUpdated", logs: logs, sub: sub}, nil
}

// WatchLogSlasherUpdated is a free log subscription operation binding the contract event 0x933228a1c3ba8fadd3ce47a9db5b898be647f89af99ba7c1b9a655f59ea306c8.
//
// Solidity: event LogSlasherUpdated(previousSlasher address, nextSlasher address)
func (_Bindings *BindingsFilterer) WatchLogSlasherUpdated(opts *bind.WatchOpts, sink chan<- *BindingsLogSlasherUpdated) (event.Subscription, error) {

	logs, sub, err := _Bindings.contract.WatchLogs(opts, "LogSlasherUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BindingsLogSlasherUpdated)
				if err := _Bindings.contract.UnpackLog(event, "LogSlasherUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// BindingsOwnershipRenouncedIterator is returned from FilterOwnershipRenounced and is used to iterate over the raw logs and unpacked data for OwnershipRenounced events raised by the Bindings contract.
type BindingsOwnershipRenouncedIterator struct {
	Event *BindingsOwnershipRenounced // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BindingsOwnershipRenouncedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BindingsOwnershipRenounced)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BindingsOwnershipRenounced)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BindingsOwnershipRenouncedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BindingsOwnershipRenouncedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BindingsOwnershipRenounced represents a OwnershipRenounced event raised by the Bindings contract.
type BindingsOwnershipRenounced struct {
	PreviousOwner common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipRenounced is a free log retrieval operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: event OwnershipRenounced(previousOwner indexed address)
func (_Bindings *BindingsFilterer) FilterOwnershipRenounced(opts *bind.FilterOpts, previousOwner []common.Address) (*BindingsOwnershipRenouncedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _Bindings.contract.FilterLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BindingsOwnershipRenouncedIterator{contract: _Bindings.contract, event: "OwnershipRenounced", logs: logs, sub: sub}, nil
}

// WatchOwnershipRenounced is a free log subscription operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: event OwnershipRenounced(previousOwner indexed address)
func (_Bindings *BindingsFilterer) WatchOwnershipRenounced(opts *bind.WatchOpts, sink chan<- *BindingsOwnershipRenounced, previousOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _Bindings.contract.WatchLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BindingsOwnershipRenounced)
				if err := _Bindings.contract.UnpackLog(event, "OwnershipRenounced", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// BindingsOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Bindings contract.
type BindingsOwnershipTransferredIterator struct {
	Event *BindingsOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BindingsOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BindingsOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BindingsOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BindingsOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BindingsOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BindingsOwnershipTransferred represents a OwnershipTransferred event raised by the Bindings contract.
type BindingsOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_Bindings *BindingsFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*BindingsOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Bindings.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BindingsOwnershipTransferredIterator{contract: _Bindings.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_Bindings *BindingsFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BindingsOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Bindings.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BindingsOwnershipTransferred)
				if err := _Bindings.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
