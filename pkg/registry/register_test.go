// +build integration

package registry

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	CONSUL_ADDRESS = "CONSUL_ADDR"
)

func consulAddr() string {
	return os.Getenv(CONSUL_ADDRESS)
}

func TestClient_Register(t *testing.T) {

	client, err := New(consulAddr())
	require.NoError(t, err)

	_, err = client.Register("dwarf", 3000)
	require.NoError(t, err)
}

func TestClient_DeRegister(t *testing.T) {

	client, err := New(consulAddr())
	require.NoError(t, err)

	id, err = client.Register("dwarf", 3000)
	require.NoError(t, err)

	require.NoError(t, client.DeRegister(id))
}
