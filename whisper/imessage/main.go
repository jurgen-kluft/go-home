package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var server ServerType

func init() {
	var err error
	defaultZone, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Panic().Msg(err.Error())
	}

	// Configure logging
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().In(defaultZone)
	}

	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "Mon Jan 2 15:04:05"}
	log.Logger = zerolog.New(output).With().Caller().Timestamp().Logger()

	// Init Server
	server = ServerType{}
	server.Chats = make([]*Chat, 0)
	server.ChatMap = make(map[string]*Chat, 0)
	server.Attachments = make(map[string]Attachment, 0)
}

func main() {
	var chatdb string
	var intervalString string
	flag.StringVar(&chatdb, "db", "", "Filepath to your iMessage chat.db (typically ~/Library/Messages/chat.db)")
	flag.StringVar(&intervalString, "interval", "", "Interval in milliseconds to refresh chats")
	debug := flag.Bool("debug", false, "sets log level to debug")
	flag.Parse()

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
			fileparts := strings.Split(file, "/")
			filename := strings.Replace(fileparts[len(fileparts)-1], ".go", "", -1)
			return filename + ":" + strconv.Itoa(line)
		}
	}

	if intervalString == "" {
		intervalString = "2000"
	}
	i, err := strconv.Atoi(intervalString)
	if err != nil {
		log.Panic().Msg(err.Error())
	}
	interval := time.Millisecond * time.Duration(i)
	log.Info().Msg(fmt.Sprintf("Checking DB every %f seconds", interval.Seconds()))

	if chatdb == "" {
		log.Panic().Msg("--db Chat DB is required.")
	}

	server.SQLiteFile = chatdb
	log.Info().Msg(fmt.Sprintf("Doing startup parse of DB at location %s", server.SQLiteFile))

	// Start router and DB readers
	go startDatabaseReader(interval)
	startRouter()
}

func startDatabaseReader(interval time.Duration) {
	for {
		if hasBeenModified() {
			log.Info().Msg("DB parse requested by timeout")
			refreshChats()
		}
		time.Sleep(interval)
	}
}
