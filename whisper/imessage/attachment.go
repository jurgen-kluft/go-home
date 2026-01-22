package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// Attachment corresponds to a row in the 'attachment' table
type Attachment struct {
	RowID         *string `json:"RowID"`
	MessageID     *string `json:"MessageID"`
	Filename      *string `json:"Filename"`
	MIMEType      *string `json:"MIMEType"`
	TransferState int     `json:"TransferState"`
	TotalBytes    int     `json:"TotalBytes"`
}

type JSONResponse struct {
	Output interface{} `json:"output,omitempty"`
	Status string      `json:"status,omitempty"`
	OK     bool        `json:"ok"`
	Method string      `json:"method,omitempty"`
	ID     int         `json:"id,omitempty"`
}

// Write to an http writer, adding extra info and HTTP status as needed
func (response *JSONResponse) Write(w *http.ResponseWriter, r *http.Request) {
	// Deref writer
	writer := *w

	writer.Header().Set("Content-Type", "application/json")

	// Add string Status if it doesn't exist, add appropriate headers
	if response.OK {
		if response.Status == "" {
			response.Status = "success"
		}
		writer.WriteHeader(http.StatusOK)
	} else {
		if response.Status == "" {
			response.Status = "fail"
			writer.WriteHeader(http.StatusBadRequest)
		} else if response.Status == "error" {
			writer.WriteHeader(http.StatusNoContent)
		} else {
			writer.WriteHeader(http.StatusBadRequest)
		}
	}

	// Write out this response
	json.NewEncoder(writer).Encode(response)
}

func handleAttachmentsGetAll(w http.ResponseWriter, r *http.Request) {
	server.lock.RLock()
	defer server.lock.RUnlock()

	response := JSONResponse{
		OK:     true,
		Output: server.Attachments,
	}

	response.Write(&w, r)
}

func handleAttachmentsGet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	server.lock.RLock()
	defer server.lock.RUnlock()

	attachment, ok := server.Attachments[params["id"]]
	if !ok {
		em := fmt.Sprintf("Could not fetch attachment id %s. It likely doesn't exist", params["id"])
		log.Error().Msg(em)
		response := JSONResponse{OK: false, Output: ""}
		response.Write(&w, r)
		return
	}

	response := JSONResponse{OK: true, Output: attachment}
	response.Write(&w, r)
}

func parseAttachmentRows(rows *sql.Rows) []Attachment {
	var out []Attachment
	if rows == nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		m := Attachment{}
		err := rows.Scan(&m.MessageID, &m.RowID, &m.Filename, &m.MIMEType, &m.TransferState, &m.TotalBytes)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		out = append(out, m)

		// Ensure attachment is registered with the server
		server.Attachments[*m.MessageID] = m
	}
	return out
}

func getAttachment(messageID string) ([]Attachment, error) {
	sql := fmt.Sprintf("SELECT message_attachment_join.message_id, attachment.ROWID, attachment.filename, attachment.mime_type, attachment.transfer_state, attachment.total_bytes FROM attachment LEFT OUTER JOIN message_attachment_join ON message_attachment_join.attachment_id = attachment.ROWID WHERE message_attachment_join.message_id = %s ORDER BY attachment.created_date DESC", messageID)
	rows, err := query(sql)
	if err != nil {
		return nil, err
	}
	return parseAttachmentRows(rows), nil
}
