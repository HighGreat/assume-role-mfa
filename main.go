package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/go-ini/ini"
)

const (
	AWS_CONFG_PATH       = ".aws/config"
	AWS_CREDENTIALS_PATH = ".aws/credentials"
)

var (
	duration    = flag.Int64("duration", 60*60, "session duration")
	profileName = flag.String("profile", "", "profile name in config file")
	sessionName = flag.String("session_name", "", "session name")
)

func main() {
	flag.Parse()

	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	config, err := ini.Load(filepath.Join(homedir, AWS_CONFG_PATH))
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	configSection := config.Section(fmt.Sprintf("profile %s", *profileName))

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("aws mfa code:")
	scanner.Scan()
	mfaCode := scanner.Text()

	output, err := assumeRole(&AssumeRoleInput{
		DurationSeconnd: *duration,
		RoleArn:         configSection.Key("role_arn").Value(),
		RoleSessionName: *sessionName,
		SerialNumber:    configSection.Key("mfa_serial").Value(),
		TokenCode:       mfaCode,
	})

	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	// save credentials
	credential, err := ini.Load(filepath.Join(homedir, AWS_CREDENTIALS_PATH))
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	credentialSection := credential.Section(*profileName)

	credentialSection.Key("aws_access_key_id").SetValue(output.AccessKeyId)
	credentialSection.Key("aws_secret_access_key").SetValue(output.SecretAccessKey)
	credentialSection.Key("aws_session_token").SetValue(output.SessionToken)

	if err := credential.SaveTo(filepath.Join(homedir, AWS_CREDENTIALS_PATH)); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

type AssumeRoleInput struct {
	DurationSeconnd int64
	RoleArn         string
	RoleSessionName string
	SerialNumber    string
	TokenCode       string
}

type Credential struct {
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
}

func assumeRole(input *AssumeRoleInput) (*Credential, error) {
	session, err := session.NewSessionWithOptions(session.Options{
		Profile: "default",
	})

	if err != nil {
		return nil, err
	}

	svc := sts.New(session)

	output, err := svc.AssumeRole(&sts.AssumeRoleInput{
		DurationSeconds: aws.Int64(input.DurationSeconnd),
		RoleArn:         aws.String(input.RoleArn),
		RoleSessionName: aws.String(input.RoleSessionName),
		SerialNumber:    aws.String(input.SerialNumber),
		TokenCode:       aws.String(input.TokenCode),
	})

	if err != nil {
		return nil, err
	}

	return &Credential{
		AccessKeyId:     *output.Credentials.AccessKeyId,
		SecretAccessKey: *output.Credentials.SecretAccessKey,
		SessionToken:    *output.Credentials.SessionToken,
	}, nil
}
