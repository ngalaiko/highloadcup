package web

import (
	"net/http"

	"github.com/ngalayko/highloadcup/schema"
	"github.com/zenazn/goji/web"
)

// GetVisitsHandler is handler for /users/:id/visits
func (wb *Web) GetVisitsHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	id, err := parseId(c)
	if err != nil {
		responseErr(w, err)
		return
	}

	user := &schema.User{}
	if err := wb.db.Get(schema.EntityUsers, id, user); err != nil {
		responseErr(w, err)
		return
	}

	visits := []*schema.Visit{}
	if err := wb.db.GetIds(schema.EntityVisits, user.VisitIDs, &visits); err != nil {
		responseErr(w, err)
		return
	}

	responseJson(w, visits)
}
