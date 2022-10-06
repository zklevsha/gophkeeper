package config

import (
	"testing"
)

func TestClientConfigEnv(t *testing.T) {
	testServerAddress := "1.2.3.4"
	want := ClientConfig{ServerAddress: testServerAddress}
	t.Run("Get agent config with env variables", func(t *testing.T) {
		t.Setenv("GK_SERVER_ADDRESS", testServerAddress)
		res := GetClientConfig()
		if res != want {
			t.Errorf("ClientConfig mismatch: have: %v,  want: %v", res, want)
		}
	})

}

func TestClientConfigNoEnv(t *testing.T) {
	want := ClientConfig{ServerAddress: serverAddressDefault}
	t.Run("Get agent config with no env variables", func(t *testing.T) {
		res := GetClientConfig()
		if res != want {
			t.Errorf("ClientConfig mismatch: have: %v,  want: %v", res, want)
		}
	})
}
