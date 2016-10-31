package routes

// given a valid jwt
// generate a new token
// blacklist the token with a TTL until it expires
type RefreshTokenRoute struct {
	Route
}

func (r *RefreshTokenRoute) CompileRoute() (*Route, error) {
	return &r.Route, nil
}
