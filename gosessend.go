package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Sender struct {
	verbose bool
}

const DefaultRegion = "us-west-1"

func auth(sender *Sender) (*ses.SES, error) {
	log.Debug().Msg("auth")

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = os.Getenv("AMAZON_REGION")
	}
	if region == "" {
		region = DefaultRegion
	}

	log.Debug().Str("Region", region).Send()

	config := &aws.Config{
		Region:                        aws.String(region),
		CredentialsChainVerboseErrors: &sender.verbose,
		Credentials:                   credentials.NewSharedCredentials("", ""),
	}

	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	log.Debug().Msg("Session created")

	_, err = sess.Config.Credentials.Get()
	if err != nil {
		return nil, err
	}

	svc := ses.New(sess)

	return svc, nil
}

func send(svc *ses.SES, rawEmail []byte) error {
	log.Debug().Msg("Sending ... ")

	input := &ses.SendRawEmailInput{
		FromArn: aws.String(""),
		RawMessage: &ses.RawMessage{
			Data: rawEmail,
		},
		ReturnPathArn: aws.String(""),
		Source:        aws.String(""),
		SourceArn:     aws.String(""),
	}

	result, err := svc.SendRawEmail(input)
	if err != nil {
		return err
	}

	log.Debug().Msg("send is complete")

	log.Info().Msg(result.String())

	return nil
}

func readFile(fileName string) (b []byte, err error) {
	log.Debug().Msg("Reading file ...")

	b, err = ioutil.ReadFile(fileName)
	if err == nil {
		log.Debug().Int("rawSize", len(b)).Send()
	}

	return
}

func printUsage() {
	exe := filepath.Base(os.Args[0])

	fmt.Printf(`  Use: ./%s mail-file
`, exe)
}

func checkArgs() (string, bool) {
	if len(os.Args) < 2 {
		kingpin.Usage()
	}
	verboseA := kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()
	fileNameA := kingpin.Arg("mail-file", "Raw mail file.").Required().String()
	kingpin.Parse()

	return *fileNameA, *verboseA
}

func main() {
	fileName, verbose := checkArgs()

	sender := &Sender{verbose: verbose}

	if sender.verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: zerolog.TimeFormatUnix})

	log.Debug().Str("file", fileName).Msg("Using")

	rawEmail, err := readFile(fileName)
	if err != nil {
		log.Err(err).Send()
		return
	}

	svc, err := auth(sender)
	if err != nil {
		log.Err(err).Send()
		return
	}

	err = send(svc, rawEmail)
	if err != nil {
		log.Err(err).Send()
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == "InvalidClientTokenId" {
				log.Error().Msg("Probably wrong key in $HOME/.aws/credentials")
			}
		}
		return
	}
}
