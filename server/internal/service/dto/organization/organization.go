package organization

type (
	RespRelatedKeyPairs struct {
		_               struct{}
		JWTPairs        `json:"jwt_pairs"`
		CubeSignerPairs []CubeSignerPairs `json:"cubesigner_pairs"`
	}

	CubeSignerPairs struct {
		_       struct{}
		Private string `json:"private"`
		Public  string `json:"public"`
	}

	JWTPairs struct {
		_       struct{}
		Private string `json:"private"`
		Public  string `json:"public"`
	}
)
