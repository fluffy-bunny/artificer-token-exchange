package auth

import (
	"echo-starter/internal/wellknown"

	middleware_oidc "github.com/fluffy-bunny/grpcdotnetgo/pkg/middleware/oidc"
	services_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/services/claimsprincipal"
)

// BuildGrpcEntrypointPermissionsClaimsMap ...
func BuildGrpcEntrypointPermissionsClaimsMap() map[string]*middleware_oidc.EntryPointConfig {
	entryPointClaimsBuilder := services_claimsprincipal.NewEntryPointClaimsBuilder()
	// HEALTH SERVICE START
	//---------------------------------------------------------------------------------------------------
	// health check is open to anyone
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(wellknown.HealthPath)
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(wellknown.ReadinessPath)
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(wellknown.LoginPath)
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(wellknown.LogoutPath)
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(wellknown.OIDCCallbackPath)
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(wellknown.HomePath)
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(wellknown.AboutPath)
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(wellknown.ErrorPath)
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(wellknown.UnauthorizedPath)
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen("/css*")
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen("/assets*")
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen("/js*")

	entryPointClaimsBuilder.GetClaimsConfig(wellknown.ArtistsPath).
		WithGrpcEntrypointPermissionsClaimFactsMapOR(
			services_claimsprincipal.NewClaimFactType(wellknown.ClaimTypeAuthenticated),
		)
	entryPointClaimsBuilder.AddMetaData(wellknown.ArtistsPath, map[string]interface{}{
		"onUnauthenticated": "login",
	})
	entryPointClaimsBuilder.GetClaimsConfig(wellknown.APIArtistsPath).
		WithGrpcEntrypointPermissionsClaimFactsMapOR(
			services_claimsprincipal.NewClaimFactType(wellknown.ClaimTypeAuthenticated),
		)
	entryPointClaimsBuilder.GetClaimsConfig(wellknown.APIArtistsIdPath).
		WithGrpcEntrypointPermissionsClaimFactsMapOR(
			services_claimsprincipal.NewClaimFactType(wellknown.ClaimTypeAuthenticated),
		)
	entryPointClaimsBuilder.GetClaimsConfig(wellknown.APIArtistsIdAlbumsPath).
		WithGrpcEntrypointPermissionsClaimFactsMapOR(
			services_claimsprincipal.NewClaimFactType(wellknown.ClaimTypeAuthenticated),
		)

	entryPointClaimsBuilder.GetClaimsConfig(wellknown.ProfilesPath).
		WithGrpcEntrypointPermissionsClaimFactsMapOR(
			services_claimsprincipal.NewClaimFactType(wellknown.ClaimTypeAuthenticated),
		)
	entryPointClaimsBuilder.AddMetaData(wellknown.ProfilesPath, map[string]interface{}{
		"onUnauthenticated": "login",
	})

	entryPointClaimsBuilder.GetClaimsConfig(wellknown.DeepPath).
		WithGrpcEntrypointPermissionsClaimFactsMapOR(
			services_claimsprincipal.NewClaimFactType(wellknown.ClaimTypeAuthenticated),
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
