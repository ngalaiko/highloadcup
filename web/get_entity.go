package web

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

// GetEntityHandler is a handler for /:entity/:id
func (wb *Web) GetEntityHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	entity, err := parseEntity(c)
	if err != nil {
		responseErr(w, err)
		return
	}

	id, err := parseId(c)
	if err != nil {
		responseErr(w, err)
		return
	}

	result, err := wb.db.GetBytes(entity, id)
	if err != nil {
		responseErr(w, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(result)
}
