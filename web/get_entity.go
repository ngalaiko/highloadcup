package web

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

// GetEntityHandler is a handler for /:entity/:id
func (wb *Web) GetEntityHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	entity, err := parseEntity(c)
	if err != nil {
		responseErr(r, w, err)
		return
	}

	id, err := parseId(c)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	result, err := wb.db.Get(entity, id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	responseJson(w, result)
}
