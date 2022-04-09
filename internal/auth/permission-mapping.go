package auth

import (
	"echo-starter/internal/wellknown"

	core_wellknown "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/wellknown"
	middleware_oidc "github.com/fluffy-bunny/grpcdotnetgo/pkg/middleware/oidc"
	services_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/services/claimsprincipal"
)

// BuildGrpcEntrypointPermissionsClaimsMap ...
func BuildGrpcEntrypointPermissionsClaimsMap() map[string]*middleware_oidc.EntryPointConfig {
	entryPointClaimsBuilder := services_claimsprincipal.NewEntryPointClaimsBuilder()
	// HEALTH SERVICE START
	//---------------------------------------------------------------------------------------------------
	// health check is open to anyone
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(wellknown.WebHookPath)
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(wellknown.HealthzPath)
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(wellknown.ReadyPath)
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(wellknown.LoginPath)
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(wellknown.LogoutPath)
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(wellknown.OIDCCallbackPath)
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(wellknown.OAuth2CallbackPath)
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(wellknown.HomePath)
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(wellknown.AboutPath)
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(wellknown.ErrorPath)
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(wellknown.UnauthorizedPath)
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(wellknown.APIDevPath)
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen("/static*")

	// GraphQL
	//---------------------------------------------------------------------------------------------------
	entryPointClaimsBuilder.GetClaimsConfig(wellknown.GraphQLEndpointPath).
		WithGrpcEntrypointPermissionsClaimFactsMapOR(
			services_claimsprincipal.NewClaimFactType(core_wellknown.ClaimTypeAuthenticated),
		)
	entryPointClaimsBuilder.GetClaimsConfig(wellknown.GraphiQLPath).
		WithGrpcEntrypointPermissionsClaimFactsMapOR(
			services_claimsprincipal.NewClaimFactType(core_wellknown.ClaimTypeAuthenticated),
		)
	entryPointClaimsBuilder.AddMetaData(wellknown.GraphiQLPath, map[string]interface{}{
		"onUnauthenticated": "login",
	})
	// ACCOUNTS
	//---------------------------------------------------------------------------------------------------
	entryPointClaimsBuilder.GetClaimsConfig(wellknown.AccountsPath).
		WithGrpcEntrypointPermissionsClaimFactsMapOR(
			services_claimsprincipal.NewClaimFactType(core_wellknown.ClaimTypeAuthenticated),
		)
	entryPointClaimsBuilder.AddMetaData(wellknown.AccountsPath, map[string]interface{}{
		"onUnauthenticated": "login",
	})
	entryPointClaimsBuilder.GetClaimsConfig(wellknown.APIAccountsPath).
		WithGrpcEntrypointPermissionsClaimFactsMapOR(
			services_claimsprincipal.NewClaimFactType(core_wellknown.ClaimTypeAuthenticated),
		)

	// ARTISTS
	//---------------------------------------------------------------------------------------------------
	entryPointClaimsBuilder.GetClaimsConfig(wellknown.ArtistsPath).
		WithGrpcEntrypointPermissionsClaimFactsMapOR(
			services_claimsprincipal.NewClaimFactType(core_wellknown.ClaimTypeAuthenticated),
		)
	entryPointClaimsBuilder.AddMetaData(wellknown.ArtistsPath, map[string]interface{}{
		"onUnauthenticated": "login",
	})
	entryPointClaimsBuilder.GetClaimsConfig(wellknown.APIArtistsPath).
		WithGrpcEntrypointPermissionsClaimFactsMapOR(
			services_claimsprincipal.NewClaimFactType(core_wellknown.ClaimTypeAuthenticated),
		)
	entryPointClaimsBuilder.GetClaimsConfig(wellknown.APIArtistsIdPath).
		WithGrpcEntrypointPermissionsClaimFactsMapOR(
			services_claimsprincipal.NewClaimFactType(core_wellknown.ClaimTypeAuthenticated),
		)
	entryPointClaimsBuilder.GetClaimsConfig(wellknown.APIArtistsIdAlbumsPath).
		WithGrpcEntrypointPermissionsClaimFactsMapOR(
			services_claimsprincipal.NewClaimFactType(core_wellknown.ClaimTypeAuthenticated),
		)

	entryPointClaimsBuilder.GetClaimsConfig(wellknown.ProfilesPath).
		WithGrpcEntrypointPermissionsClaimFactsMapOR(
			services_claimsprincipal.NewClaimFactType(core_wellknown.ClaimTypeAuthenticated),
		)
	entryPointClaimsBuilder.AddMetaData(wellknown.ProfilesPath, map[string]interface{}{
		"onUnauthenticated": "login",
	})

	entryPointClaimsBuilder.GetClaimsConfig(wellknown.DeepPath).
		WithGrpcEntrypointPermissionsClaimFactsMapOR(
			services_claimsprincipal.NewClaimFactType(core_wellknown.ClaimTypeAuthenticated),
		).GetChild().
		WithGrpcEntrypointPermissionsClaimFactsMapOR(
			services_claimsprincipal.NewClaimFactTypeAndValue(wellknown.ClaimTypeDeep, wellknown.ClaimValueRead),
			services_claimsprincipal.NewClaimFactTypeAndValue(wellknown.ClaimTypeDeep, wellknown.ClaimValueReadWrite),
			services_claimsprincipal.NewClaimFactTypeAndValue(wellknown.ClaimTypeDeep, wellknown.ClaimValueReadWriteAll),
		)
	entryPointClaimsBuilder.AddMetaData(wellknown.DeepPath, map[string]interface{}{
		"onUnauthenticated": "login",
	})
	cMap := entryPointClaimsBuilder.GrpcEntrypointClaimsMap
	return cMap
}
