package web

import (
	"net/http"

	"github.com/ngalayko/highloadcup/schema"
	"github.com/zenazn/goji/web"
)

// GetLocationsAvgHandler is a handler for /locations/:id/avg
func (wb *Web) GetLocationsAvgHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	id, err := parseId(c)
	if err != nil {
		responseErr(w, err)
		return
	}

	location := &schema.Location{}
	if err := wb.db.Get(schema.EntityLocations, id, location); err != nil {
		responseErr(w, err)
		return
	}

	responseJson(w, struct {
		Avg float32 `json:"avg"`
	}{
		Avg: location.AvgMark,
	})
}
