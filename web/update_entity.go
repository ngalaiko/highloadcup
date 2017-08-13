package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ngalayko/highloadcup/schema"
	"github.com/zenazn/goji/web"
)

// GetEntityHandler is a handler for /:entity/:id
func (wb *Web) UpdateEntityHandler(c web.C, w http.ResponseWriter, r *http.Request) {
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

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseErr(w, err)
		return
	}
	defer r.Body.Close()

	val := schema.GetIEntity(entity)
	if err := json.Unmarshal(data, val); err != nil {
		responseErr(w, err)
		return
	}

	if err := val.Validate(); err != nil {
		responseErr(w, err)
		return
	}

	if err := wb.db.Get(entity, id, new(interface{})); err != nil {
		http.NotFound(w, r)
		return
	}

	if err := wb.db.CreateOrUpdate(val); err != nil {
		responseErr(w, err)
		return
	}

	responseJson(w, struct{}{})
}
