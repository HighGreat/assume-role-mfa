## How to use

### Install
```shell
$ go get https://github.com/HighGreat/assume-role-mfa
```

### Prepare aws configuration
```
...

[profile example-config]
source_profile = default
role_arn = arn:aws:iam::***:role/***
mfa_serial = arn:aws:iam::***:mfa/***
```

### Run command
```
$ AWS_PROFILE=default assume-role-mfa --profile example --duration $((60*60)) --session_name john
$ export AWS_PROFILE=example
```
