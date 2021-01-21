package main

import (
	"github.com/go-ini/ini"
)

type AwsConfig struct {
	RoleArn   string
	MfaSerial string
}

func loadAwsConfig(config_file_path string, profile string) (*AwsConfig, error) {
	config, err := ini.Load(config_file_path)
	if err != nil {
		return nil, err
	}

	section := config.Section("profile " + profile)

	return &AwsConfig{
		RoleArn:   section.Key("role_arn").Value(),
		MfaSerial: section.Key("mfa_serial").Value(),
	}, nil
}

func saveCredential(credentials_file_path string, profile string, credential *Credential) error {
	credentials, err := ini.Load(credentials_file_path)
	if err != nil {
		return err
	}

	section := credentials.Section(profile)

	section.Key("aws_access_key_id").SetValue(credential.AccessKeyId)
	section.Key("aws_secret_access_key").SetValue(credential.SecretAccessKey)
	section.Key("aws_session_token").SetValue(credential.SessionToken)

	if err := credentials.SaveTo(credentials_file_path); err != nil {
		return err
	}

	return nil
}
