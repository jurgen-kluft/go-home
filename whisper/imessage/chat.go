package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// Chat corresponds to each row in the 'chat' table, along with a lock for the Messages
type Chat struct {
	RowID           int             `json:"RowID"`
	ID              string          `json:"ID"`
	Name            string          `json:"Name"`
	DisplayName     string          `json:"DisplayName"`
	ServiceName     string          `json:"ServiceName"`
	Messages        []Message       `json:"Messages"`
	Recipients      map[string]bool `json:"Recipients"`
	LastMessageDate int             `json:"LastMessageDate"`
	lock            sync.RWMutex
}

func handleChatGetAll(w http.ResponseWriter, r *http.Request) {
	refreshChats()

	// Sort chats
	chats := make(map[int]*Chat, 0)
	for _, chat := range server.ChatMap {
		chats[chat.LastMessageDate] = chat
	}

	//format.WriteResponse(&w, r, format.JSONResponse{OK: true, Output: chats})
	response := JSONResponse{OK: true, Output: chats}
	response.Write(&w, r)
}

func handleChatGet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	server.lock.Lock()
	defer server.lock.Unlock()

	server.openDB()
	defer server.DB.Close()

	chat, err := getAllMessages(params["id"])
	if err != nil {
		log.Error().Msg(err.Error())
		//format.WriteResponse(&w, r, format.JSONResponse{OK: false, Output: err.Error()})
		response := JSONResponse{OK: false, Output: err.Error()}
		response.Write(&w, r)
		return
	}

	//format.WriteResponse(&w, r, format.JSONResponse{OK: true, Output: chat})
	response := JSONResponse{OK: true, Output: chat}
	response.Write(&w, r)
}

func handleChatGetLast(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	server.lock.Lock()
	defer server.lock.Unlock()

	server.openDB()
	defer server.DB.Close()

	message, err := getLastMessage(params["id"])
	if err != nil {
		log.Error().Msg(err.Error())
		//format.WriteResponse(&w, r, format.JSONResponse{OK: false, Output: err.Error()})
		response := JSONResponse{OK: false, Output: err.Error()}
		response.Write(&w, r)
		return
	}

	response := JSONResponse{OK: true, Output: message}
	response.Write(&w, r)
}

func parseChatRows(rows *sql.Rows) []*Chat {
	var out []*Chat
	if rows == nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		c := Chat{Recipients: make(map[string]bool, 0)}
		rows.Scan(&c.RowID, &c.ID, &c.Name, &c.DisplayName)
		out = append(out, &c)
	}
	return out
}

func refreshChats() {
	server.lock.Lock()
	defer server.lock.Unlock()

	server.openDB()
	defer server.DB.Close()

	sql := "SELECT DISTINCT chat.ROWID, chat.chat_identifier, chat.guid, chat.display_name FROM chat LEFT OUTER JOIN chat_message_join ON chat.ROWID = chat_message_join.chat_id WHERE chat.service_name = 'iMessage' ORDER BY chat_message_join.message_date DESC"
	rows, err := query(sql)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	newChats := parseChatRows(rows)

	start := time.Now()
	updatedChats := 0
	for _, chat := range newChats {
		// Only fetch entire chat history if there are new messages
		if !hasNewMessages(chat.ID) {
			continue
		}

		updatedChats++
		chat.Messages, err = getAllMessages(chat.ID)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		if len(chat.Messages) > 0 {
			// Ensure chat has these messages registered
			for _, message := range chat.Messages {
				chat.Recipients[*message.Handle.ID] = true
			}
			chat.LastMessageDate = chat.Messages[0].Date
			server.ChatMap[chat.ID] = chat
		}
	}

	log.Debug().Msg(fmt.Sprintf("Took %dms to parse %d rows", time.Since(start).Milliseconds(), updatedChats))
}
