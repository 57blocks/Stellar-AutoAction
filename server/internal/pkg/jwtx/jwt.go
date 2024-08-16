package jwtx

func init() {
	ES256 = &JWTes256{
		Header: Header{
			Algorithm: AlgES256,
			Type:      "JWT",
		},
	}

	HS256 = &JWThs256{
		Header: Header{
			Algorithm: AlgHS256,
			Type:      "JWT",
		},
	}

	RS256 = &JWTrs256{
		Header: Header{
			Algorithm: AlgRS256,
			Type:      "JWT",
		},
	}
}

var (
	ES256 *JWTes256
	HS256 *JWThs256
	RS256 *JWTrs256
)

type (
	Parser interface {
		Parse(s string) (string, bool)
	}

	Assigner interface {
		Assign() (string, error)
	}

	JWTes256 struct {
		Header Header
	}
	JWThs256 struct {
		Header Header
	}
	JWTrs256 struct {
		Header Header
	}

	Header struct {
		Algorithm Algorithm `json:"alg"`
		Type      string    `json:"typ"`
	}
)

// Algorithm is the type of the algorithm used to sign the token
type Algorithm string

const (
	AlgES256 Algorithm = "ES256"
	AlgHS256 Algorithm = "HS256"
	AlgRS256 Algorithm = "RS256"
)
