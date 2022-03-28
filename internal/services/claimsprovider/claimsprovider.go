package claimsprovider

import (
	contracts_claimsprovider "echo-starter/internal/contracts/claimsprovider"

	"echo-starter/internal/wellknown"

	"errors"
	"reflect"

	contracts_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/claimsprincipal"
	"github.com/golang/mock/gomock"

	"github.com/rs/zerolog/log"

	di "github.com/fluffy-bunny/sarulabsdi"
)

type (
	service struct {
	}
	serviceMock struct {
	}
)

func assertImplementation() {
	var _ contracts_claimsprovider.IClaimsProvider = (*service)(nil)
}

var mockProfileStore map[string][]*contracts_claimsprincipal.Claim

func init() {
	mockProfileStore = make(map[string][]*contracts_claimsprincipal.Claim)

	mockProfileStore[""] = []*contracts_claimsprincipal.Claim{}

	mockProfileStore["profile1"] = []*contracts_claimsprincipal.Claim{
		{
			Type:  wellknown.ClaimTypeDeep,
			Value: wellknown.ClaimValueRead,
		},
		{
			Type:  wellknown.ClaimTypeDeep,
			Value: wellknown.ClaimValueReadWrite,
		},
		{
			Type:  wellknown.ClaimTypeDeep,
			Value: wellknown.ClaimValueReadWriteAll,
		},
	}

	mockProfileStore["profile2"] = []*contracts_claimsprincipal.Claim{
		{
			Type:  wellknown.ClaimTypeDeep,
			Value: wellknown.ClaimValueRead,
		},
		{
			Type:  wellknown.ClaimTypeDeep,
			Value: wellknown.ClaimValueReadWrite,
		},
	}
	mockProfileStore["profile3"] = []*contracts_claimsprincipal.Claim{
		{
			Type:  wellknown.ClaimTypeDeep,
			Value: wellknown.ClaimValueRead,
		},
	}
}

var reflectType = reflect.TypeOf((*service)(nil))
var reflectTypeMock = reflect.TypeOf((*serviceMock)(nil))

// AddSingletonIClaimsProvider registers the *service as a singleton.
func AddSingletonIClaimsProvider(builder *di.Builder) {
	log.Info().Str("DI", "IClaimsProvider").Send()
	contracts_claimsprovider.AddSingletonIClaimsProvider(builder, reflectType)
}

func AddSingletonIClaimsProviderMock(builder *di.Builder, ctrl *gomock.Controller) {
	log.Info().Str("DI", "IClaimsProvider - MOCK").Send()
	contracts_claimsprovider.AddSingletonIClaimsProvider(builder, reflectTypeMock)
}
func (s *service) Ctor() {}
func (s *service) GetProfiles(userID string) ([]string, error) {
	return []string{"profile1", "profile2", "profile3"}, nil
}
func (s *service) GetClaims(userID string, profile string) ([]*contracts_claimsprincipal.Claim, error) {
	return nil, errors.New("not implemented")
}
func (s *serviceMock) GetProfiles(userID string) ([]string, error) {
	return []string{"profile1", "profile2", "profile3"}, nil
}
func (s *serviceMock) GetClaims(userID string, profile string) ([]*contracts_claimsprincipal.Claim, error) {
	claims, ok := mockProfileStore[profile]
	if !ok {
		return nil, errors.New("profile not found")
	}
	return claims, nil
}
