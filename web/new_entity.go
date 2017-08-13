package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"log"

	"github.com/zenazn/goji/web"

	"github.com/ngalayko/highloadcup/schema"
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

	if _, err := wb.db.Get(entity, val.IntID()); err == nil {
		responseErr(w, fmt.Errorf("entity already exists"))
		return
	}

	if err := wb.db.CreateOrUpdate(val); err != nil {
		responseErr(w, err)
		return
	}

	if newVisit, ok := val.(*schema.Visit); ok {
		go wb.onVisitInserted(newVisit)
	}

	responseJson(w, struct{}{})
}

func (wb *Web) onVisitInserted(visit *schema.Visit) {
	user, err := wb.db.GetUser(visit.UserID)
	if err != nil {
		log.Printf("error when updating visit relates (onInserted): %s", err)
	}

	location, err := wb.db.GetLocation(visit.UserID)
	if err != nil {
		log.Printf("error when updating visit relates (onInserted): %s", err)
	}

	user.VisitIDs = append(user.VisitIDs, visit.ID)
	if err := wb.db.CreateOrUpdate(user); err != nil {
		log.Printf("error when updating visit relates (onInserted): %s", err)
	}

	location.VisitIDs = append(location.VisitIDs, visit.ID)
	if err := wb.db.CreateOrUpdate(location); err != nil {
		log.Printf("error when updating visit relates (onInserted): %s", err)
	}
}
