package config

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewConfig(t *testing.T) {
	Convey("Given some config", t, func() {
		env := "test"
		cfg, err := NewConfig(env)
		So(err, ShouldBeNil)
		So(cfg, ShouldNotBeNil)
		So(cfg.PrivateKeyPath, ShouldEqual, "keys/test/private_key.pem")
	})

	Convey("Global config", t, func() {
		So(Cfg, ShouldNotBeNil)
		So(Cfg.PrivateKeyPath, ShouldEqual, "keys/test/private_key.pem")
	})
}
