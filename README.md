# gosessend
Small golang SES sending email tool

## Usage
```
usage: gosessend.exe [<flags>] <raw-mail-file>

Flags:
      --help     Show context-sensitive help (also try --help-long and
                 --help-man).
  -v, --verbose  Verbose mode.

Args:
  <raw-mail-file>  Raw mail file.
```

## Build

1) Setup Go by instruction from https://golang.org/dl/
2) 
```bash
git clone https://github.com/inv2004/gosessend
```
```bash
cd gosessend
```
```bash
go get -d .
```
```bash
go build
```
3)
Create ``$HOME/.aws/credentials`` for ``default`` profile: https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html
