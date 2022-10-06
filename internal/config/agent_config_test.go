package config

import (
	"testing"
)

func TestAgentConfigEnv(t *testing.T) {
	testServerAddress := "1.2.3.4"
	want := AgentConfig{ServerAddress: testServerAddress}
	t.Run("Get agent config with env variables", func(t *testing.T) {
		t.Setenv("GK_SERVER_ADDRESS", testServerAddress)
		res := GetAgentConfig()
		if res != want {
			t.Errorf("AgentConfig mismatch: have: %v,  want: %v", res, want)
		}
	})

}

func TestAgentConfigNoEnv(t *testing.T) {
	want := AgentConfig{ServerAddress: serverAddressDefault}
	t.Run("Get agent config with no env variables", func(t *testing.T) {
		res := GetAgentConfig()
		if res != want {
			t.Errorf("AgentConfig mismatch: have: %v,  want: %v", res, want)
		}
	})
}
