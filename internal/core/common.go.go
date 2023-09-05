package core

type fuzzy struct {
	repo Repository
}

// NewFuzzy is a constructor function for the fuzzy type
func NewFuzzy(db Repository) *fuzzy {
	return &fuzzy{
		repo: db,
	}
}
