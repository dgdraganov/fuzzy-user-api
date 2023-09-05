package verify

type Registry interface {
	VerifyUser(jwtToken string) (map[string]any, error)
}
