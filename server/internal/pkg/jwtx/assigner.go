package jwtx

var (
	assign = func(claim TokenClaims) (string, error) {

		return "", nil
	}
)

func (p *JWTes256) Assign(claim TokenClaims) (string, error) {
	return assign(claim)
}

func (p *JWThs256) Assign(claim TokenClaims) (string, error) {
	return assign(claim)
}

func (p *JWTrs256) Assign(claim TokenClaims) (string, error) {
	return assign(claim)
}
