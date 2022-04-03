package shared

const (
	LoginParamsSessionKey = "_auth_3be0_LoginParamsSessionKey"
	AuthStateSessionKey   = "_auth_3be0_AuthStateSessionKey"
	ProfileSessionKey     = "_auth_3be0_ProfileSessionKey"
)

type (
	LoginParms struct {
		RedirectURL string `json:"redirect_url" xml:"redirect_url" form:"redirect_url" query:"redirect_url"`
	}
)
