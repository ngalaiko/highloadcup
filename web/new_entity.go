package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ngalayko/highloadcup/schema"
	"github.com/zenazn/goji/web"
)

// NewEntityHandler is a handler for /:entity/new
func (wb *Web) NewEntityHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	entity, err := parseEntity(c)
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

	if err := wb.db.Get(entity, val.IntID(), new(interface{})); err == nil {
		responseErr(w, fmt.Errorf("entity already exists"))
		return
	}

	if err := wb.db.CreateOrUpdate(val); err != nil {
		responseErr(w, err)
		return
	}

	responseJson(w, struct{}{})
}
