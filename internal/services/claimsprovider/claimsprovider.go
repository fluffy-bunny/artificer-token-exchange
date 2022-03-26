package claimsprovider

import (
	contracts_claimsprovider "echo-starter/internal/contracts/claimsprovider"
	mocks_claimsprovider "echo-starter/internal/mocks/claimsprovider"
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
)

func assertImplementation() {
	var _ contracts_claimsprovider.IClaimsProvider = (*service)(nil)
}

var reflectType = reflect.TypeOf((*service)(nil))

// AddSingletonIClaimsProvider registers the *service as a singleton.
func AddSingletonIClaimsProvider(builder *di.Builder) {
	log.Info().Str("DI", "IClaimsProvider").Send()
	contracts_claimsprovider.AddSingletonIClaimsProvider(builder, reflectType)
}

// AddSingletonIClaimsProvider registers the *service as a singleton.
func AddSingletonIClaimsProviderMock(builder *di.Builder, ctrl *gomock.Controller) {
	log.Info().Str("DI", "IClaimsProvider - MOCK").Send()
	obj := mocks_claimsprovider.NewMockIClaimsProvider(ctrl)
	obj.EXPECT().GetClaims(gomock.Any()).DoAndReturn(func(userID string) ([]*contracts_claimsprincipal.Claim, error) {
		return []*contracts_claimsprincipal.Claim{
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
		}, nil
	})
	contracts_claimsprovider.AddSingletonIClaimsProviderByObj(builder, obj)
}

func (s *service) Ctor() {}
func (s *service) GetClaims(userID string) ([]*contracts_claimsprincipal.Claim, error) {
	return nil, errors.New("not implemented")
}
