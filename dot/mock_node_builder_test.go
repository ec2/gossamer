// Code generated by MockGen. DO NOT EDIT.
// Source: node.go
//
// Generated by this command:
//
//	mockgen -source=node.go -destination=mock_node_builder_test.go -package=dot
//

// Package dot is a generated GoMock package.
package dot

import (
	reflect "reflect"

	config "github.com/ChainSafe/gossamer/config"
	core "github.com/ChainSafe/gossamer/dot/core"
	digest "github.com/ChainSafe/gossamer/dot/digest"
	network "github.com/ChainSafe/gossamer/dot/network"
	rpc "github.com/ChainSafe/gossamer/dot/rpc"
	state "github.com/ChainSafe/gossamer/dot/state"
	sync "github.com/ChainSafe/gossamer/dot/sync"
	system "github.com/ChainSafe/gossamer/dot/system"
	types "github.com/ChainSafe/gossamer/dot/types"
	babe "github.com/ChainSafe/gossamer/lib/babe"
	grandpa "github.com/ChainSafe/gossamer/lib/grandpa"
	keystore "github.com/ChainSafe/gossamer/lib/keystore"
	runtime "github.com/ChainSafe/gossamer/lib/runtime"
	gomock "go.uber.org/mock/gomock"
)

// MocknodeBuilderIface is a mock of nodeBuilderIface interface.
type MocknodeBuilderIface struct {
	ctrl     *gomock.Controller
	recorder *MocknodeBuilderIfaceMockRecorder
}

// MocknodeBuilderIfaceMockRecorder is the mock recorder for MocknodeBuilderIface.
type MocknodeBuilderIfaceMockRecorder struct {
	mock *MocknodeBuilderIface
}

// NewMocknodeBuilderIface creates a new mock instance.
func NewMocknodeBuilderIface(ctrl *gomock.Controller) *MocknodeBuilderIface {
	mock := &MocknodeBuilderIface{ctrl: ctrl}
	mock.recorder = &MocknodeBuilderIfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocknodeBuilderIface) EXPECT() *MocknodeBuilderIfaceMockRecorder {
	return m.recorder
}

// createBABEService mocks base method.
func (m *MocknodeBuilderIface) createBABEService(config *config.Config, st *state.Service, ks KeyStore, cs *core.Service, telemetryMailer Telemetry) (*babe.Service, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "createBABEService", config, st, ks, cs, telemetryMailer)
	ret0, _ := ret[0].(*babe.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// createBABEService indicates an expected call of createBABEService.
func (mr *MocknodeBuilderIfaceMockRecorder) createBABEService(config, st, ks, cs, telemetryMailer any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "createBABEService", reflect.TypeOf((*MocknodeBuilderIface)(nil).createBABEService), config, st, ks, cs, telemetryMailer)
}

// createBlockVerifier mocks base method.
func (m *MocknodeBuilderIface) createBlockVerifier(st *state.Service) *babe.VerificationManager {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "createBlockVerifier", st)
	ret0, _ := ret[0].(*babe.VerificationManager)
	return ret0
}

// createBlockVerifier indicates an expected call of createBlockVerifier.
func (mr *MocknodeBuilderIfaceMockRecorder) createBlockVerifier(st any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "createBlockVerifier", reflect.TypeOf((*MocknodeBuilderIface)(nil).createBlockVerifier), st)
}

// createCoreService mocks base method.
func (m *MocknodeBuilderIface) createCoreService(config *config.Config, ks *keystore.GlobalKeystore, st *state.Service, net *network.Service) (*core.Service, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "createCoreService", config, ks, st, net)
	ret0, _ := ret[0].(*core.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// createCoreService indicates an expected call of createCoreService.
func (mr *MocknodeBuilderIfaceMockRecorder) createCoreService(config, ks, st, net any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "createCoreService", reflect.TypeOf((*MocknodeBuilderIface)(nil).createCoreService), config, ks, st, net)
}

// createDigestHandler mocks base method.
func (m *MocknodeBuilderIface) createDigestHandler(st *state.Service) (*digest.Handler, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "createDigestHandler", st)
	ret0, _ := ret[0].(*digest.Handler)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// createDigestHandler indicates an expected call of createDigestHandler.
func (mr *MocknodeBuilderIfaceMockRecorder) createDigestHandler(st any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "createDigestHandler", reflect.TypeOf((*MocknodeBuilderIface)(nil).createDigestHandler), st)
}

// createGRANDPAService mocks base method.
func (m *MocknodeBuilderIface) createGRANDPAService(config *config.Config, st *state.Service, ks KeyStore, net *network.Service, telemetryMailer Telemetry) (*grandpa.Service, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "createGRANDPAService", config, st, ks, net, telemetryMailer)
	ret0, _ := ret[0].(*grandpa.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// createGRANDPAService indicates an expected call of createGRANDPAService.
func (mr *MocknodeBuilderIfaceMockRecorder) createGRANDPAService(config, st, ks, net, telemetryMailer any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "createGRANDPAService", reflect.TypeOf((*MocknodeBuilderIface)(nil).createGRANDPAService), config, st, ks, net, telemetryMailer)
}

// createNetworkService mocks base method.
func (m *MocknodeBuilderIface) createNetworkService(config *config.Config, stateSrvc *state.Service, telemetryMailer Telemetry) (*network.Service, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "createNetworkService", config, stateSrvc, telemetryMailer)
	ret0, _ := ret[0].(*network.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// createNetworkService indicates an expected call of createNetworkService.
func (mr *MocknodeBuilderIfaceMockRecorder) createNetworkService(config, stateSrvc, telemetryMailer any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "createNetworkService", reflect.TypeOf((*MocknodeBuilderIface)(nil).createNetworkService), config, stateSrvc, telemetryMailer)
}

// createRPCService mocks base method.
func (m *MocknodeBuilderIface) createRPCService(params rpcServiceSettings) (*rpc.HTTPServer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "createRPCService", params)
	ret0, _ := ret[0].(*rpc.HTTPServer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// createRPCService indicates an expected call of createRPCService.
func (mr *MocknodeBuilderIfaceMockRecorder) createRPCService(params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "createRPCService", reflect.TypeOf((*MocknodeBuilderIface)(nil).createRPCService), params)
}

// createRuntimeStorage mocks base method.
func (m *MocknodeBuilderIface) createRuntimeStorage(st *state.Service) (*runtime.NodeStorage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "createRuntimeStorage", st)
	ret0, _ := ret[0].(*runtime.NodeStorage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// createRuntimeStorage indicates an expected call of createRuntimeStorage.
func (mr *MocknodeBuilderIfaceMockRecorder) createRuntimeStorage(st any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "createRuntimeStorage", reflect.TypeOf((*MocknodeBuilderIface)(nil).createRuntimeStorage), st)
}

// createStateService mocks base method.
func (m *MocknodeBuilderIface) createStateService(config *config.Config) (*state.Service, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "createStateService", config)
	ret0, _ := ret[0].(*state.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// createStateService indicates an expected call of createStateService.
func (mr *MocknodeBuilderIfaceMockRecorder) createStateService(config any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "createStateService", reflect.TypeOf((*MocknodeBuilderIface)(nil).createStateService), config)
}

// createSystemService mocks base method.
func (m *MocknodeBuilderIface) createSystemService(cfg *types.SystemInfo, stateSrvc *state.Service) (*system.Service, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "createSystemService", cfg, stateSrvc)
	ret0, _ := ret[0].(*system.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// createSystemService indicates an expected call of createSystemService.
func (mr *MocknodeBuilderIfaceMockRecorder) createSystemService(cfg, stateSrvc any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "createSystemService", reflect.TypeOf((*MocknodeBuilderIface)(nil).createSystemService), cfg, stateSrvc)
}

// initNode mocks base method.
func (m *MocknodeBuilderIface) initNode(config *config.Config) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "initNode", config)
	ret0, _ := ret[0].(error)
	return ret0
}

// initNode indicates an expected call of initNode.
func (mr *MocknodeBuilderIfaceMockRecorder) initNode(config any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "initNode", reflect.TypeOf((*MocknodeBuilderIface)(nil).initNode), config)
}

// loadRuntime mocks base method.
func (m *MocknodeBuilderIface) loadRuntime(config *config.Config, ns *runtime.NodeStorage, stateSrvc *state.Service, ks *keystore.GlobalKeystore, net *network.Service) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "loadRuntime", config, ns, stateSrvc, ks, net)
	ret0, _ := ret[0].(error)
	return ret0
}

// loadRuntime indicates an expected call of loadRuntime.
func (mr *MocknodeBuilderIfaceMockRecorder) loadRuntime(config, ns, stateSrvc, ks, net any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "loadRuntime", reflect.TypeOf((*MocknodeBuilderIface)(nil).loadRuntime), config, ns, stateSrvc, ks, net)
}

// newSyncService mocks base method.
func (m *MocknodeBuilderIface) newSyncService(config *config.Config, st *state.Service, finalityGadget BlockJustificationVerifier, verifier *babe.VerificationManager, cs *core.Service, net *network.Service, telemetryMailer Telemetry) (*sync.Service, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "newSyncService", config, st, finalityGadget, verifier, cs, net, telemetryMailer)
	ret0, _ := ret[0].(*sync.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// newSyncService indicates an expected call of newSyncService.
func (mr *MocknodeBuilderIfaceMockRecorder) newSyncService(config, st, finalityGadget, verifier, cs, net, telemetryMailer any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "newSyncService", reflect.TypeOf((*MocknodeBuilderIface)(nil).newSyncService), config, st, finalityGadget, verifier, cs, net, telemetryMailer)
}
