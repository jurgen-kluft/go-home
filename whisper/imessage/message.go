package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// Thanks to linjunpop
// https://github.com/linjunpop/imessage/blob/master/lib/imessage/sender.rb
const sendFileScript = `
	on run argv
		set toAddress to first item of argv
		set theFilePath to second item of argv
		set theFile to POSIX file theFilePath
		tell application "System Events"
				if exists file theFilePath then
						tell application "Messages"
								send theFile to buddy toAddress of (service 1 whose service type is iMessage)
						end tell
				else
						error "File does not exist."
				end if
		end tell
	end run`

const sendMessageScript = `
	on run argv
		set toAddress to first item of argv
		set message to second item of argv
		tell application "Messages"
			send message to buddy toAddress of (service 1 whose service type is iMessage)
		end tell
	end run`

// Message corresponds to each row in 'message' table
type Message struct {
	RowID         string       `json:"RowID"`
	Text          *string      `json:"Text"`
	IsFromMe      bool         `json:"IsFromMe"`
	HasAttachment bool         `json:"HasAttachment"`
	Delivered     bool         `json:"Delivered"`
	Date          int          `json:"Date"`
	DateDelivered int          `json:"DateDelivered"`
	DateRead      int          `json:"DateRead"`
	Handle        Handle       `json:"Handle"`
	Attachments   []Attachment `json:"Attachment"`
}

// Handle might be expanded in the future, I'm not sure
type Handle struct {
	ID *string `json:"ID"`
}

func sendMessage(chatID string, message string, file *string) error {
	if chatID == "" {
		return fmt.Errorf("Empty chat ID")
	}
	chat, ok := server.ChatMap[chatID]
	if !ok {
		return fmt.Errorf("Chat does not exist in server. You may need to refresh, or check your ID")
	}

	// Parse recipients or file
	var (
		arg    []string
		script string
	)

	// Switch command if we're sending an attachment
	if file == nil {
		for ID := range chat.Recipients {
			arg = append(arg, ID)
		}
		script = sendMessageScript
	} else {
		path, err := filepath.Abs(*file)
		if err != nil {
			return fmt.Errorf("Failed opening file %s", *file)
		}
		arg = append(arg, path)
		script = sendFileScript
	}

	// Setup and Execute
	cmd := exec.Command("oascript", "-e", script, strings.Join(arg, ","))
	err := cmd.Run()
	if err != nil {
		log.Error().Msg(err.Error())
	}

	return nil
}

func handleSendMessage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	body, err := readBody(w, r)
	if err != nil {
		//format.WriteResponse(&w, r, format.JSONResponse{OK: false, Output: err.Error()})
		response := JSONResponse{OK: false, Output: err.Error()}
		response.Write(&w, r)
	}

	if body == nil {
		//format.WriteResponse(&w, r, format.JSONResponse{OK: false, Output: "Unknown request body"})
		response := JSONResponse{OK: false, Output: "Unknown request body"}
		response.Write(&w, r)
	}

	err = sendMessage(params["id"], string(body), nil)
	if err != nil {
		log.Error().Msg(err.Error())
		//format.WriteResponse(&w, r, format.JSONResponse{OK: false, Output: err.Error()})
		response := JSONResponse{OK: false, Output: err.Error()}
		response.Write(&w, r)
		return
	}

	//format.WriteResponse(&w, r, format.JSONResponse{OK: true, Output: "OK"})
	response := JSONResponse{OK: true, Output: "OK"}
	response.Write(&w, r)
}

func hasNewMessages(chatID string) bool {
	serverChat, ok := server.ChatMap[chatID]
	if !ok {
		return true
	}

	sql := fmt.Sprintf("SELECT message_date FROM chat_message_join LEFT OUTER JOIN chat ON chat.ROWID = chat_message_join.chat_id WHERE chat_id = %d AND chat.service_name = 'iMessage' ORDER BY message_date DESC LIMIT 1", serverChat.RowID)
	rows, err := query(sql)
	if err != nil {
		return true
	}
	defer rows.Close()
	var lastMessageDate int
	for rows.Next() {
		rows.Scan(&lastMessageDate)
	}

	if serverChat.LastMessageDate != lastMessageDate {
		log.Info().Msg(fmt.Sprintf("Row %d not found, last recorded date is %d while the server shows %d", serverChat.RowID, lastMessageDate, serverChat.LastMessageDate))
		return true
	}

	return false
}

func getAllMessages(chatID string) ([]Message, error) {
	// Default to handle ID, check if it's a group chat
	selector := "chat.room_name IS NULL AND handle.id"
	if strings.Contains(chatID, "chat") {
		selector = "chat.chat_identifier"
	}

	sql := fmt.Sprintf("SELECT DISTINCT message.ROWID, handle.id, message.text, message.is_from_me, message.cache_has_attachments, message.is_delivered, message.date, message.date_delivered, message.date_read FROM message LEFT OUTER JOIN chat ON chat.room_name = message.cache_roomnames LEFT OUTER JOIN handle ON handle.ROWID = message.handle_id WHERE message.service = 'iMessage' AND %s = '%s' ORDER BY message.date DESC LIMIT 50", selector, chatID)
	rows, err := query(sql)
	if err != nil {
		return nil, err
	}
	return parseMessageRows(rows), nil
}

func getLastMessage(chatID string) (Message, error) {
	// Default to handle ID, check if it's a group chat
	selector := "handle.id"
	if strings.Contains(chatID, "chat") {
		selector = "chat.chat_identifier"
	}

	sql := fmt.Sprintf("SELECT DISTINCT message.ROWID, handle.id, message.text, message.is_from_me, message.cache_has_attachments, message.is_delivered, message.date, message.date_delivered, message.date_read FROM message LEFT OUTER JOIN chat ON chat.room_name = message.cache_roomnames LEFT OUTER JOIN handle ON handle.ROWID = message.handle_id WHERE message.service = 'iMessage' AND %s = '%s' ORDER BY message.date DESC LIMIT 1", selector, chatID)
	rows, err := query(sql)
	if err != nil {
		return Message{}, err
	}

	messages := parseMessageRows(rows)
	if messages == nil {
		return Message{}, nil
	}
	return messages[0], nil
}

func parseMessageRows(rows *sql.Rows) []Message {
	var out []Message
	if rows == nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		m := Message{}
		err := rows.Scan(&m.RowID, &m.Handle.ID, &m.Text, &m.IsFromMe, &m.HasAttachment, &m.Delivered, &m.Date, &m.DateDelivered, &m.DateRead)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		if m.Handle.ID == nil {
			newID := "me"
			m.Handle.ID = &newID
		}
		if m.HasAttachment {
			attachments, err := getAttachment(m.RowID)
			if err == nil {
				m.Attachments = attachments
			} else {
				log.Error().Msg(err.Error())
			}
		}
		out = append(out, m)
	}
	return out
}
