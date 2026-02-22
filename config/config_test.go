package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test.toml")
	require.NoError(t, err)

	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	testData := `
		[chain]
		chain_id = "100"
		url = "https://l1.example.com"
		starting_height = 100
		confirmation_depth = 10
		polling_interval = 5000

		[db]
		host = "127.0.0.1"
		port = 5432
		user = "postgres"
		password = "postgres"
	    name = "indexer"

		[http]
		host = "127.0.0.1"
		port = 8080

		[metrics]
		host = "127.0.0.1"
		port = 7300
	`

	data := []byte(testData)
	err = os.WriteFile(tmpfile.Name(), data, 0o600)
	require.NoError(t, err)

	defer os.Remove(tmpfile.Name())

	err = tmpfile.Close()
	require.NoError(t, err)

	conf, err := LoadConfig(tmpfile.Name())
	require.NoError(t, err)

	require.Equal(t, conf.Chain.ChainID, "100")
	require.Equal(t, conf.Chain.PollingInterval, uint(5000))
	require.Equal(t, conf.Chain.URL, "https://l1.example.com")
	require.Equal(t, conf.Chain.StartingHeight, uint64(100))
	require.Equal(t, conf.Chain.ConfirmationDepth, uint(10))
	require.Equal(t, conf.DB.Host, "127.0.0.1")
	require.Equal(t, conf.DB.Port, 5432)
	require.Equal(t, conf.DB.User, "postgres")
	require.Equal(t, conf.DB.Password, "postgres")
	require.Equal(t, conf.DB.Name, "indexer")
	require.Equal(t, conf.HTTPServer.Host, "127.0.0.1")
	require.Equal(t, conf.HTTPServer.Port, 8080)
	require.Equal(t, conf.MetricsServer.Host, "127.0.0.1")
	require.Equal(t, conf.MetricsServer.Port, 7300)
}
