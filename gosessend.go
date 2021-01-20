package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	log.Debug().Int("rawSize", len(b)).Send()
	return
}

func printUsage() {
	fmt.Print(`  Use: ./gosessend mail-file
`)
}

func checkArgs() (string, bool) {
	if len(os.Args) <= 1 {
		printUsage()
		os.Exit(1)
	} else if len(os.Args) == 2 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		printUsage()
		os.Exit(0)
	} else if len(os.Args) == 2 && (os.Args[1] == "--verbose" || os.Args[1] == "-v") {
		printUsage()
		os.Exit(1)
	}
	if len(os.Args) > 3 {
		fmt.Println("too many arguments")
		printUsage()
		os.Exit(1)
	}
	if len(os.Args) == 3 {
		if os.Args[1] == "-v" || os.Args[1] == "--verbose" {
			return os.Args[2], true
		} else if os.Args[2] == "-v" || os.Args[2] == "--verbose" {
			return os.Args[1], true
		}
	}
	return os.Args[1], false
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
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == "InvalidClientTokenId" {
				log.Error().Msg("Probably wrong key in $HOME/.aws/credentials")
			}
		}
		log.Err(err).Send()
		return
	}
}
