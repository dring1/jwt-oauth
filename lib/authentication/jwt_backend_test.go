package authentication

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestJWT(t *testing.T) {
	Convey("Given some arbitrary string", t, func() {
		invalidPrivatePem := `
-----BEGIN RSA PRIVATE KEY-----
Proc-Type: 4,ENCRYPTED
DEK-Info: DES-EDE3-CBC,C11FAFDC8551B267

-----END RSA PRIVATE KEY-----
    `
		key, err := getPrivateKey([]byte(invalidPrivatePem))
		So(err, ShouldNotBeNil)
		So(key, ShouldBeNil)
	})

	Convey("Given some arbitrary string", t, func() {
		invalidPrivatePem := `
-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDh/rwNKtzYz24402FAeyWvSp6Z3L3q8VAn5AAXuG/q4jhM8jFw
MQ3TG161sqJIpvCAJ7B/LwR71lcAGzSgDCLmOA6S9MUlbqgUF16bAC3w8SG/1j3L
o6cPRhcVgWU71sN3M90M35YEUB7oz56ttCoMFjgBe6sZW7aGWpM8aMzs2wIDAQAB
AoGBANx/lW1nf7kOkmVhYIbwYHFqZdqLdMWxktqI68o5CwFqnTH+MFxdkDaEguDX
HN2z+/2eO+ersT0+gP46jRsHHFgbjqQTswcZOmIxICaKlpcizGWqaDe1bQ6XXVpB
6sI3Jr26DtntSxgB/j6rd1CibnsQCidyIMUl1cRF+9hltjCxAkEA9+wvod10V+Sm
pJGr3rbi4kxfbb3JH59BzN3sy4aGY2+yA5VxbH9dgyYFNebYp7oiftMu+tz2lU9F
6KZdWNyWowJBAOlbqROcvODelQBq0vVr2qJ0cL5YT1ySgIkKpakfXq577TBjtlpa
uA9IQteqUuPBp8Extv5XFKj3J8zRWGVmjGkCQQDkZCrFPNOvHK7/sEra0zRUMPNA
j7O2c+oUJuW74OPwurcNYiCpSPQGm4H1VAKHEYwxta7z35cxmWPXnVslP6FtAkBc
SyRT3WnWhjHoOFe3OTD/j44HumWo90he6xcaDI4l9F2bBdTZZ4fkg2/sXDDsY2s1
vbPiZA6HxTi4iROtByIBAkEA3v+sZgUyN+ioUX+T+ZjLorgwt+SSo8bUr/yMK7aQ
7H34avELRlm6u3PD9kHaQJ3rk+RZm5LiWXzbZeOI+bT2Ug==
-----END RSA PRIVATE KEY-----
      `
		key, err := getPrivateKey([]byte(invalidPrivatePem))
		So(err, ShouldBeNil)
		So(key, ShouldNotBeNil)
	})
}
