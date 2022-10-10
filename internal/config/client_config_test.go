package config

import (
	"testing"
)

func TestClientConfig(t *testing.T) {
	testServerAddress := "1.1.1.1:443"

	tt := []struct {
		name string
		args []string
		want ClientConfig
	}{
		{name: "no flags", args: []string{},
			want: ClientConfig{ServerAddress: serverAddressDefault}},
		{name: "all flags", args: []string{"-a", testServerAddress},
			want: ClientConfig{ServerAddress: testServerAddress}},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			res := GetClientConfig(tc.args)
			if res != tc.want {
				t.Errorf("ClientConfig mismatch: have: %v,  want: %v", res, tc.want)
			}
		})
	}
}

func TestClientConfigEnv(t *testing.T) {
	testServerAddress := "1.2.3.4"
	want := ClientConfig{ServerAddress: testServerAddress}
	t.Run("Get client config with env variables", func(t *testing.T) {
		t.Setenv("GK_SERVER_ADDRESS", testServerAddress)
		res := GetClientConfig([]string{})
		if res != want {
			t.Errorf("ClientConfig mismatch: have: %v,  want: %v", res, want)
		}
	})

}

func TestClientConfigNoEnv(t *testing.T) {
	want := ClientConfig{ServerAddress: serverAddressDefault}
	t.Run("Get client config with no env variables", func(t *testing.T) {
		res := GetClientConfig([]string{})
		if res != want {
			t.Errorf("ClientConfig mismatch: have: %v,  want: %v", res, want)
		}
	})
}
