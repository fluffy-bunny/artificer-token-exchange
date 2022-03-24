// Code generated by MockGen. DO NOT EDIT.
// Source: echo-starter/internal/contracts/auth (interfaces: IOIDCAuthenticator)

// Package auth is a generated GoMock package.
package auth

import (
	context "context"
	reflect "reflect"

	oidc "github.com/coreos/go-oidc/v3/oidc"
	gomock "github.com/golang/mock/gomock"
	oauth2 "golang.org/x/oauth2"
)

// MockIOIDCAuthenticator is a mock of IOIDCAuthenticator interface.
type MockIOIDCAuthenticator struct {
	ctrl     *gomock.Controller
	recorder *MockIOIDCAuthenticatorMockRecorder
}

// MockIOIDCAuthenticatorMockRecorder is the mock recorder for MockIOIDCAuthenticator.
type MockIOIDCAuthenticatorMockRecorder struct {
	mock *MockIOIDCAuthenticator
}

// NewMockIOIDCAuthenticator creates a new mock instance.
func NewMockIOIDCAuthenticator(ctrl *gomock.Controller) *MockIOIDCAuthenticator {
	mock := &MockIOIDCAuthenticator{ctrl: ctrl}
	mock.recorder = &MockIOIDCAuthenticatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIOIDCAuthenticator) EXPECT() *MockIOIDCAuthenticatorMockRecorder {
	return m.recorder
}

// AuthCodeURL mocks base method.
func (m *MockIOIDCAuthenticator) AuthCodeURL(arg0 string, arg1 ...oauth2.AuthCodeOption) string {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AuthCodeURL", varargs...)
	ret0, _ := ret[0].(string)
	return ret0
}

// AuthCodeURL indicates an expected call of AuthCodeURL.
func (mr *MockIOIDCAuthenticatorMockRecorder) AuthCodeURL(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthCodeURL", reflect.TypeOf((*MockIOIDCAuthenticator)(nil).AuthCodeURL), varargs...)
}

// VerifyIDToken mocks base method.
func (m *MockIOIDCAuthenticator) VerifyIDToken(arg0 context.Context, arg1 *oauth2.Token) (*oidc.IDToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyIDToken", arg0, arg1)
	ret0, _ := ret[0].(*oidc.IDToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyIDToken indicates an expected call of VerifyIDToken.
func (mr *MockIOIDCAuthenticatorMockRecorder) VerifyIDToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyIDToken", reflect.TypeOf((*MockIOIDCAuthenticator)(nil).VerifyIDToken), arg0, arg1)
}
