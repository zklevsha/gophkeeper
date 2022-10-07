package config

import "testing"

func TestServerConfig(t *testing.T) {
	testServerAddress := "1.1.1.1:443"
	testDSN := "postgres://username:password@localhost:5432/database_name"
	testKey := "verystrongkey"
	tt := []struct {
		name string
		args []string
		want ServerConfig
	}{
		{name: "only required flags", args: []string{"-d", testDSN, "-k", testKey},
			want: ServerConfig{ServerAddress: serverAddressDefault,
				DSN: testDSN, Key: testKey}},
		{name: "all flags", args: []string{"-a", testServerAddress,
			"-d", testDSN, "-k", testKey},
			want: ServerConfig{ServerAddress: testServerAddress,
				DSN: testDSN, Key: testKey}},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			res := GetServerConfig(tc.args)
			if res != tc.want {
				t.Errorf("ClientConfig mismatch: have: %v,  want: %v", res, tc.want)
			}
		})
	}
}

func TestGetServerConfigEnv(t *testing.T) {
	t.Run("Get server config with env variables", func(t *testing.T) {
		want := ServerConfig{ServerAddress: "testServ", DSN: "test_dsn", Key: "verystrongkey"}
		t.Setenv("GK_SERVER_ADDRESS", want.ServerAddress)
		t.Setenv("GK_DB_DSN", want.DSN)
		t.Setenv("GK_KEY", want.Key)
		have := GetServerConfig([]string{})

		if have != want {
			t.Errorf("ServerConfig mismatch: have: %v, want: %v", have, want)
		}
	})
}
