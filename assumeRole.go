package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

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
