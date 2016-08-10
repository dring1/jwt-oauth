package config

var environments = map[string]string{
	"production": "production.json",
	"staging":    "staging.json",
	"test":       "test.json",
}

type Cfg struct {
	JWTExpirationDelta int
	PrivateKey         []byte
	PublicKey          []byte
	Port               int
	GitHubClientID     string
	GitHubClientSecret string
	OauthRedirectURL   string
}

func NewConfig(opts ...func(*Cfg) error) (*Cfg, error) {
	c := &Cfg{
		JWTExpirationDelta: 60,
		PrivateKey:         make([]byte, 10),
		PublicKey:          make([]byte, 10),
		Port:               8080,
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}
