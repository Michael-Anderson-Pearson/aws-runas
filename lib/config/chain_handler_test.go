package config

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func ExampleChainConfigHandler_Config() {
	os.Unsetenv("AWS_REGION")
	os.Unsetenv("AWS_DEFAULT_PROFILE")
	os.Unsetenv("SESSION_TOKEN_DURATION")
	os.Setenv("AWS_CONFIG_FILE", "test/aws.cfg")
	os.Setenv("AWS_PROFILE", "has_bad_role")
	os.Setenv("CREDENTIALS_DURATION", "4h")
	defer os.Unsetenv("AWS_PROFILE")
	defer os.Unsetenv("CREDENTIALS_DURATION")

	c := new(AwsConfig)
	h := DefaultConfigHandler(DefaultConfigHandlerOpts)
	if err := h.Config(c); err != nil {
		fmt.Printf("Unexpected error during Config(): %v\n", err)
	}

	fmt.Println(c.Name)
	fmt.Println(c.RoleArn)
	fmt.Println(c.GetMfaSerial())
	fmt.Println(c.GetCredDuration())
	// Output:
	// has_bad_role
	// aws:iam::012345678901:mfa/my_iam_user
	// 12345678
	// 4h0m0s
}

func TestChainConfigHandler_Add(t *testing.T) {
	h := DefaultConfigHandler(DefaultConfigHandlerOpts).(*ChainConfigHandler)
	h.Add(NewCmdlineConfigHandler(new(ConfigHandlerOpts), new(CmdlineOptions)))

	if len(h.handlers) != 3 {
		t.Errorf("Expected 3 handlers after Add(), got: %d", len(h.handlers))
	}
}

func TestChainConfigHandler_ConfigNil(t *testing.T) {
	defer func() {
		if x := recover(); x != nil {
			t.Errorf("Unexpected panic calling Config() with nil argument")
		}
	}()
	h := DefaultConfigHandler(DefaultConfigHandlerOpts)
	if err := h.Config(nil); err != nil {
		t.Errorf("Unexpected error calling Config() will nil argument")
	}
}

func TestChainConfigHandler_Config(t *testing.T) {
	os.Unsetenv("AWS_REGION")
	os.Unsetenv("AWS_DEFAULT_PROFILE")
	os.Unsetenv("SESSION_TOKEN_DURATION")
	os.Setenv("AWS_CONFIG_FILE", "test/aws.cfg")
	os.Setenv("AWS_PROFILE", "has_role")
	os.Setenv("CREDENTIALS_DURATION", "4h")
	defer os.Unsetenv("AWS_PROFILE")
	defer os.Unsetenv("CREDENTIALS_DURATION")

	c := new(AwsConfig)
	h := DefaultConfigHandler(DefaultConfigHandlerOpts)
	if err := h.Config(c); err != nil {
		t.Errorf("Unexpected error during Config(): %v", err)
	}

	if c.GetCredDuration() != 4*time.Hour || c.GetSessionDuration() != 18*time.Hour {
		t.Errorf("Unexpected duration values from config & env: %+v", c)
	}
}
