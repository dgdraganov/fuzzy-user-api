package core

type fuzzy struct {
	repo      Repository
	jwtIssuer JwtIssuer
}

// NewFuzzy is a constructor function for the fuzzy type
func NewFuzzy(db Repository, issuer JwtIssuer) *fuzzy {
	return &fuzzy{
		repo:      db,
		jwtIssuer: issuer,
	}
}
