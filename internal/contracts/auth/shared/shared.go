package shared

const (
	LoginParamsSessionKey = "e50416d9-3be0-4380-877e-48a6a04b548a"
	AuthStateSessionKey   = "c8356a32-f727-4c5a-8071-a54cd5bfda6a"
	ProfileSessionKey     = "b0752457-c75e-472e-a2e7-1249d56c027a"
)

type (
	LoginParms struct {
		RedirectURL string `json:"redirect_url" xml:"redirect_url" form:"redirect_url" query:"redirect_url"`
	}
)
