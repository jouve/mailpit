package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/axllent/mailpit/data"
	"github.com/axllent/mailpit/server/websockets"
	"github.com/axllent/mailpit/storage"
	"github.com/gorilla/mux"
)

type messagesResult struct {
	Total  int            `json:"total"`
	Unread int            `json:"unread"`
	Count  int            `json:"count"`
	Start  int            `json:"start"`
	Items  []data.Summary `json:"items"`
}

// Return a list of available mailboxes
func apiMailboxStats(w http.ResponseWriter, _ *http.Request) {
	res := storage.StatsGet()

	bytes, _ := json.Marshal(res)
	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(bytes)
}

// List messages
func apiListMessages(w http.ResponseWriter, r *http.Request) {
	start, limit := getStartLimit(r)

	messages, err := storage.List(start, limit)
	if err != nil {
		httpError(w, err.Error())
		return
	}

	stats := storage.StatsGet()

	var res messagesResult

	res.Start = start
	res.Items = messages
	res.Count = len(res.Items)
	res.Total = stats.Total
	res.Unread = stats.Unread

	bytes, _ := json.Marshal(res)
	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(bytes)
}

// Search all messages
func apiSearchMessages(w http.ResponseWriter, r *http.Request) {
	search := strings.TrimSpace(r.URL.Query().Get("query"))
	if search == "" {
		fourOFour(w)
		return
	}

	messages, err := storage.Search(search)
	if err != nil {
		httpError(w, err.Error())
		return
	}

	stats := storage.StatsGet()

	var res messagesResult

	res.Start = 0
	res.Items = messages
	res.Count = len(messages)
	res.Total = stats.Total
	res.Unread = stats.Unread

	bytes, _ := json.Marshal(res)
	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(bytes)
}

// Open a message
func apiOpenMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	msg, err := storage.GetMessage(id)
	if err != nil {
		httpError(w, err.Error())
		return
	}

	bytes, _ := json.Marshal(msg)
	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(bytes)
}

// Download/view an attachment
func apiDownloadAttachment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	partID := vars["partID"]

	a, err := storage.GetAttachmentPart(id, partID)
	if err != nil {
		httpError(w, err.Error())
		return
	}
	fileName := a.FileName
	if fileName == "" {
		fileName = a.ContentID
	}

	w.Header().Add("Content-Type", a.ContentType)
	w.Header().Set("Content-Disposition", "filename=\""+fileName+"\"")
	_, _ = w.Write(a.Content)
}

// Download the full email source as plain text
func apiDownloadRaw(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	dl := r.FormValue("dl")

	data, err := storage.GetMessageRaw(id)
	if err != nil {
		httpError(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	if dl == "1" {
		w.Header().Set("Content-Disposition", "attachment; filename=\""+id+".eml\"")
	}
	_, _ = w.Write(data)
}

// Delete all messages
func apiDeleteAll(w http.ResponseWriter, r *http.Request) {
	err := storage.DeleteAllMessages()
	if err != nil {
		httpError(w, err.Error())
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	_, _ = w.Write([]byte("ok"))
}

// Delete all selected messages
func apiDeleteSelected(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var data struct {
		IDs []string
	}
	err := decoder.Decode(&data)
	if err != nil {
		panic(err)
	}

	ids := data.IDs

	for _, id := range ids {
		if err := storage.DeleteOneMessage(id); err != nil {
			httpError(w, err.Error())
			return
		}
	}

	w.Header().Add("Content-Type", "text/plain")
	_, _ = w.Write([]byte("ok"))
}

// Delete a single message
func apiDeleteOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	err := storage.DeleteOneMessage(id)
	if err != nil {
		httpError(w, err.Error())
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	_, _ = w.Write([]byte("ok"))
}

// Mark single message as unread
func apiUnreadOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	err := storage.MarkUnread(id)
	if err != nil {
		httpError(w, err.Error())
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	_, _ = w.Write([]byte("ok"))
}

// Mark all messages as read
func apiMarkAllRead(w http.ResponseWriter, r *http.Request) {
	err := storage.MarkAllRead()
	if err != nil {
		httpError(w, err.Error())
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	_, _ = w.Write([]byte("ok"))
}

// Mark selected message as read
func apiMarkSelectedRead(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var data struct {
		IDs []string
	}
	err := decoder.Decode(&data)
	if err != nil {
		panic(err)
	}

	ids := data.IDs

	for _, id := range ids {
		if err := storage.MarkRead(id); err != nil {
			httpError(w, err.Error())
			return
		}
	}

	w.Header().Add("Content-Type", "text/plain")
	_, _ = w.Write([]byte("ok"))
}

// Mark selected message as unread
func apiMarkSelectedUnread(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var data struct {
		IDs []string
	}
	err := decoder.Decode(&data)
	if err != nil {
		panic(err)
	}

	ids := data.IDs

	for _, id := range ids {
		if err := storage.MarkUnread(id); err != nil {
			httpError(w, err.Error())
			return
		}
	}

	w.Header().Add("Content-Type", "text/plain")
	_, _ = w.Write([]byte("ok"))
}

// Websocket to broadcast changes
func apiWebsocket(w http.ResponseWriter, r *http.Request) {
	websockets.ServeWs(websockets.MessageHub, w, r)
}
