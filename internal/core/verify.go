package core

import (
	"fmt"
)

func (f *fuzzy) VerifyUser(jwtToken string) (map[string]any, error) {
	claims, err := f.jwtIssuer.Validate(jwtToken)
	if err != nil {
		return nil, fmt.Errorf("jwt validate: %w", err)
	}

	return map[string]any(claims), nil
}
