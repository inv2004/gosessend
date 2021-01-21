# gosessend
Small golang SES sending email tool

## Download
Rolling release is here https://github.com/inv2004/gosessend/releases/tag/rolling

## Usage
```
usage: gosessend.exe [<flags>] <raw-mail-file>

Flags:
      --help     Show context-sensitive help (also try --help-long and
                 --help-man).
  -v, --verbose  Verbose mode.
  -j, --json     print json for send-raw-email tool.

Args:
  <raw-mail-file>  Raw mail file.
```

## Build

1) Setup Go by instruction from https://golang.org/dl/
2) 
```bash
git clone https://github.com/inv2004/gosessend
cd gosessend
go get -d .
go build
```
3)
Create ``$HOME/.aws/credentials`` for ``default`` profile: https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html
