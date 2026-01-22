package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

// ServerType holds global state of the iMessage server
type ServerType struct {
	Chats        []*Chat
	ChatMap      map[string]*Chat
	DB           *sql.DB
	SQLiteFile   string
	LastModified time.Time
	Attachments  map[string]Attachment
	lock         sync.RWMutex
}

func (s *ServerType) openDB() {
	var err error
	server.DB, err = sql.Open("sqlite3", server.SQLiteFile)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
}

func query(SQL string) (*sql.Rows, error) {
	log.Debug().Msg(SQL)

	// Open new connection to the DB if it's been closed.
	// Some connections are persistent, to speed up rapid fire queries
	if err := server.DB.Ping(); err != nil {
		log.Info().Msg("Creating new DB connection")
		server.openDB()
		defer server.DB.Close()
	}

	rows, err := server.DB.Query(SQL)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func hasBeenModified() bool {
	info, err := os.Stat(server.SQLiteFile)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	if info.ModTime().After(server.LastModified) {
		server.LastModified = info.ModTime()
		return true
	}
	return false
}

func readBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return nil, fmt.Errorf("Error reading body: %v", err)
	}

	// Put body back
	r.Body.Close() //  must close
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return body, nil
}

// Get group chats?
//sql := "SELECT DISTINCT chat.ROWID, chat.chat_identifier, chat.guid, chat.display_name FROM message LEFT OUTER JOIN chat ON chat.room_name = message.cache_roomnames LEFT OUTER JOIN handle ON handle.ROWID = message.handle_id WHERE message.is_from_me = 0 AND chat.service_name = 'iMessage' AND message.handle_id > 0 ORDER BY message.date DESC"
