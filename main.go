package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
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

	config, err := loadAwsConfig(filepath.Join(homedir, AWS_CONFG_PATH), *profileName)
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("aws mfa code:")
	scanner.Scan()
	mfaCode := scanner.Text()

	credential, err := assumeRole(&AssumeRoleInput{
		DurationSeconnd: *duration,
		RoleArn:         config.RoleArn,
		RoleSessionName: *sessionName,
		SerialNumber:    config.MfaSerial,
		TokenCode:       mfaCode,
	})
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	if err := saveCredential(filepath.Join(homedir, AWS_CREDENTIALS_PATH), *profileName, credential); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}
