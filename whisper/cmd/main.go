package main

import (
	"log"
	"os"
	"strings"

	imessage "github.com/jurgen-kluft/go-home/whisper"
)

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	userName := "obnosis5"
	iChatDBLocation := "/Users/" + userName + "/Library/Messages/chat.db"
	c := &imessage.Config{
		SQLPath:   iChatDBLocation,                               // Set this correctly
		QueueSize: 10,                                            // 10-20 is fine. If your server is super busy, tune this.
		Retries:   3,                                             // run the applescript up to this many times to send a message. 3 works well.
		DebugLog:  log.New(os.Stdout, "[DEBUG] ", log.LstdFlags), // Log debug messages.
		ErrorLog:  log.New(os.Stderr, "[ERROR] ", log.LstdFlags), // Log errors.
	}
	s, err := imessage.Init(c)
	checkErr(err)

	done := make(chan imessage.Incoming) // Make a channel to receive incoming messages.
	s.IncomingChan(".*", done)           // Bind to all incoming messages.
	err = s.Start()                      // Start outgoing and incoming message go routines.
	checkErr(err)
	log.Print("waiting for msgs")

	for msg := range done { // wait here for messages to come in.
		if len(msg.Text) < 60 {
			log.Println("id:", msg.RowID, "from:", msg.From, "attachment?", msg.File, "msg:", msg.Text)
		} else {
			log.Println("id:", msg.RowID, "from:", msg.From, "length:", len(msg.Text))
		}
		if strings.HasPrefix(msg.Text, "Help") {
			// Reply to any incoming message that has the word "Help" as the first word.
			s.Send(imessage.Outgoing{Text: "no help for you", To: msg.From})
		}
	}
}
