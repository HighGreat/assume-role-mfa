## How to use

### Install
```shell
$ go get https://github.com/HighGreat/assume-role-mfa
```

### Prepare aws configuration
```
[default]
aws_access_key_id = ***     # Credentials.AccessKeyId
aws_secret_access_key = *** # Credentials.SecretAccessKey
aws_session_token = ***     # Credentials.SessionToken
mfa_serial = arn:aws:iam::***:mfa/***

[profile ***]
source_profile = default
role_arn = arn:aws:iam::***:role/***
mfa_serial = arn:aws:iam::***:mfa/***
```

### Run command
```
$ assume-role-mfa --profile *** --duration $((60*60)) --session_name ***
$ export AWS_PROFILE=***
```
