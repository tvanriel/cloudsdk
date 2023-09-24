package amqp_test

import (
	"testing"

	"github.com/tvanriel/cloudsdk/amqp"
	"gotest.tools/v3/assert"
)

func TestDsn(t *testing.T) {
	Address := "localhost"
	Username := "guest"
	Password := "password"

	onlyHostname := amqp.Configuration{
		Address: Address,
	}
	assert.Equal(t, onlyHostname.Dsn(), "amqp://localhost/")
	hostnameWithUsername := amqp.Configuration{
		Address:  Address,
		Username: Username,
	}
	assert.Equal(t, hostnameWithUsername.Dsn(), "amqp://guest@localhost/")
	hostnameWithUsernameAndPassword := amqp.Configuration{
		Address:  Address,
		Username: Username,
		Password: Password,
	}
	assert.Equal(t, hostnameWithUsernameAndPassword.Dsn(), "amqp://guest:password@localhost/")

	tlsHostname := amqp.Configuration{
		TLS:     true,
		Address: Address,
	}

	assert.Equal(t, tlsHostname.Dsn(), "amqps://localhost/")
	tlsHostnameWithUsername := amqp.Configuration{
		Address:  Address,
		Username: Username,
		TLS:      true,
	}
	assert.Equal(t, tlsHostnameWithUsername.Dsn(), "amqps://guest@localhost/")
	tlsHostnameWithUsernameAndPassword := amqp.Configuration{
		Address:  Address,
		Username: Username,
		Password: Password,
		TLS:      true,
	}
	assert.Equal(t, tlsHostnameWithUsernameAndPassword.Dsn(), "amqps://guest:password@localhost/")

}
