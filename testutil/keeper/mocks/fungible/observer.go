// Code generated by mockery v2.41.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	common "github.com/zeta-chain/zetacore/common"

	observertypes "github.com/zeta-chain/zetacore/x/observer/types"

	types "github.com/cosmos/cosmos-sdk/types"
)

// FungibleObserverKeeper is an autogenerated mock type for the FungibleObserverKeeper type
type FungibleObserverKeeper struct {
	mock.Mock
}

// GetParams provides a mock function with given fields: ctx
func (_m *FungibleObserverKeeper) GetParams(ctx types.Context) observertypes.Params {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetParams")
	}

	var r0 observertypes.Params
	if rf, ok := ret.Get(0).(func(types.Context) observertypes.Params); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(observertypes.Params)
	}

	return r0
}

// GetSupportedChains provides a mock function with given fields: ctx
func (_m *FungibleObserverKeeper) GetSupportedChains(ctx types.Context) []*common.Chain {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetSupportedChains")
	}

	var r0 []*common.Chain
	if rf, ok := ret.Get(0).(func(types.Context) []*common.Chain); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*common.Chain)
		}
	}

	return r0
}

// NewFungibleObserverKeeper creates a new instance of FungibleObserverKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFungibleObserverKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *FungibleObserverKeeper {
	mock := &FungibleObserverKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
