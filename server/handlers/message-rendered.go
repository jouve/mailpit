package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/axllent/mailpit/config"
	"github.com/axllent/mailpit/internal/storage"
	"github.com/gorilla/mux"
)

// GetMessageHTML (method: GET) returns a rendered version of a message's HTML part
func GetMessageHTML(w http.ResponseWriter, r *http.Request) {
	// swagger:route GET /view/{ID}.html testing GetMessageHTML
	//
	// # Render message HTML part
	//
	// Renders just the message's HTML part which can be used for UI integration testing.
	// Attached inline images are modified to link to the API provided they exist.
	// Note that is the message does not contain a HTML part then an 404 error is returned.
	//
	// The ID can be set to `latest` to return the latest message.
	//
	//	Produces:
	//	- text/html
	//
	//	Schemes: http, https
	//
	//	Parameters:
	//	  + name: ID
	//	    in: path
	//	    description: Database ID or latest
	//	    required: true
	//	    type: string
	//
	//	Responses:
	//		200: HTMLResponse
	//		default: ErrorResponse

	vars := mux.Vars(r)

	id := vars["id"]

	if id == "latest" {
		messages, err := storage.List(0, 1)
		if err != nil {
			httpError(w, err.Error())
			return
		}

		if len(messages) == 0 {
			w.WriteHeader(404)
			fmt.Fprint(w, "Message not found")
			return
		}

		id = messages[0].ID
	}

	msg, err := storage.GetMessage(id)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprint(w, "Message not found")
		return
	}
	if msg.HTML == "" {
		w.WriteHeader(404)
		fmt.Fprint(w, "This message does not contain a HTML part")
		return
	}

	html := linkInlinedImages(msg)
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(html))
}

// GetMessageText (method: GET) returns a message's text part
func GetMessageText(w http.ResponseWriter, r *http.Request) {
	// swagger:route GET /view/{ID}.txt testing GetMessageText
	//
	// # Render message text part
	//
	// Renders just the message's text part which can be used for UI integration testing.
	//
	// The ID can be set to `latest` to return the latest message.
	//
	//	Produces:
	//	- text/plain
	//
	//	Schemes: http, https
	//
	//	Parameters:
	//	  + name: ID
	//	    in: path
	//	    description: Database ID or latest
	//	    required: true
	//	    type: string
	//
	//	Responses:
	//		200: TextResponse
	//		default: ErrorResponse

	vars := mux.Vars(r)

	id := vars["id"]

	if id == "latest" {
		messages, err := storage.List(0, 1)
		if err != nil {
			httpError(w, err.Error())
			return
		}

		if len(messages) == 0 {
			w.WriteHeader(404)
			fmt.Fprint(w, "Message not found")
			return
		}

		id = messages[0].ID
	}

	msg, err := storage.GetMessage(id)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprint(w, "Message not found")
		return
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(msg.Text))
}

// This will remap all attachment images with relative paths
func linkInlinedImages(msg *storage.Message) string {
	html := msg.HTML

	for _, a := range msg.Inline {
		if a.ContentID != "" {
			re := regexp.MustCompile(`(?i)(=["\']?)(cid:` + regexp.QuoteMeta(a.ContentID) + `)(["|\'|\\s|\\/|>|;])`)
			u := config.Webroot + "api/v1/message/" + msg.ID + "/part/" + a.PartID
			matches := re.FindAllStringSubmatch(html, -1)
			for _, m := range matches {
				html = strings.ReplaceAll(html, m[0], m[1]+u+m[3])
			}
		}
	}

	for _, a := range msg.Attachments {
		if a.ContentID != "" {
			re := regexp.MustCompile(`(?i)(=["\']?)(cid:` + regexp.QuoteMeta(a.ContentID) + `)(["|\'|\\s|\\/|>|;])`)
			u := config.Webroot + "api/v1/message/" + msg.ID + "/part/" + a.PartID
			matches := re.FindAllStringSubmatch(html, -1)
			for _, m := range matches {
				html = strings.ReplaceAll(html, m[0], m[1]+u+m[3])
			}
		}
	}

	return html
}
