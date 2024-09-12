package constant

type OAuthHeader string

const (
	AuthHeader OAuthHeader = "Authorization"
	APIKey     OAuthHeader = "API-Key"
)

func (o OAuthHeader) Str() string {
	return string(o)
}

type OAuthCtxKey string

const (
	ClaimRaw OAuthCtxKey = "claim_raw"
	ClaimSub OAuthCtxKey = "claim_sub"
	ClaimIss OAuthCtxKey = "claim_iss"
)

func (o OAuthCtxKey) Str() string {
	return string(o)
}
