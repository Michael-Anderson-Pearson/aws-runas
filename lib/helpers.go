package lib

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"os"
)

// Lookup the MFA devices configured for the calling user's IAM account
func LookupMfa(sess *session.Session) ([]*iam.MFADevice, error) {
	s := iam.New(sess)

	res, err := s.ListMFADevices(&iam.ListMFADevicesInput{})
	if err != nil {
		return nil, err
	}

	return res.MFADevices, nil
}

// Print prompt to enter MFA code and gather input
func PromptForMfa() string {
	var mfaCode string
	fmt.Print("Enter MFA Code: ")
	fmt.Scanln(&mfaCode)
	return mfaCode
}

// Return the location of the AWS SDK config file.  Use the value of
// the AWS_CONFIG_FILE environment variable, if available, otherwise
// use the SDK default location
func AwsConfigFile() string {
	c := defaults.SharedConfigFilename()
	e, ok := os.LookupEnv("AWS_CONFIG_FILE")
	if ok && len(e) > 0 {
		c = e
	}
	return c
}

func AwsSession(profile string) *session.Session {
	// Doing this kills the ability to use env vars, which may mess
	// with the -M option, requiring the ~/.aws/credentials file
	// Unset AWS credential env vars
	//env := []string{
	//	"AWS_ACCESS_KEY_ID", "AWS_ACCESS_KEY",
	//	"AWS_SECRET_ACCESS_KEY", "AWS_SECRET_KEY",
	//	"AWS_SESSION_TOKEN", "AWS_SECURITY_TOKEN",
	//}
	//for _, e := range env {
	//	os.Unsetenv(e)
	//}

	opts := session.Options{
		SharedConfigState:       session.SharedConfigEnable,
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
	}

	if len(profile) > 0 {
		opts.Profile = profile
	}

	return session.Must(session.NewSessionWithOptions(opts))
}