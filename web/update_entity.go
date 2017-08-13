package web

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/zenazn/goji/web"

	"github.com/ngalayko/highloadcup/schema"
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

	oldValue, err := wb.db.Get(entity, id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	oldVisit := &schema.Visit{}
	if oldVisitPtr, ok := oldValue.(*schema.Visit); ok {
		*oldVisit = *oldVisitPtr
	}

	if err := wb.db.CreateOrUpdate(val); err != nil {
		responseErr(w, err)
		return
	}

	if newVisit, ok := val.(*schema.Visit); ok {
		go wb.onVisitUpdated(oldVisit, newVisit)
	}

	responseJson(w, struct{}{})
}

func (wb *Web) onVisitUpdated(oldVisit *schema.Visit, newVisit *schema.Visit) {
	wb.onVisitUpdatedUpdateLocation(oldVisit, newVisit)
	wb.onVisitUpdatedUpdateUser(oldVisit, newVisit)

	wb.onVisitInserted(newVisit)
}

func (wb *Web) onVisitUpdatedUpdateLocation(oldVisit *schema.Visit, newVisit *schema.Visit) {
	location, err := wb.db.GetLocation(oldVisit.UserID)
	if err != nil {
		log.Printf("error when updating visit relates (onUpdated): %s", err)
		return
	}

	var newLocationVisitIds []uint32
	for _, visitID := range location.VisitIDs {
		if visitID == oldVisit.ID {
			continue
		}

		newLocationVisitIds = append(newLocationVisitIds, visitID)
	}
	location.VisitIDs = newLocationVisitIds

	if err := wb.db.CreateOrUpdate(location); err != nil {
		log.Printf("error when updating visit relates (onUpdated): %s", err)
	}
}

func (wb *Web) onVisitUpdatedUpdateUser(oldVisit *schema.Visit, newVisit *schema.Visit) {
	user, err := wb.db.GetUser(oldVisit.UserID)
	if err != nil {

		log.Printf("error when updating visit relates (onUpdated): %s", err)
		return
	}

	var newUserVisitIds []uint32
	for _, visitID := range user.VisitIDs {
		if visitID == oldVisit.ID {
			continue
		}

		newUserVisitIds = append(newUserVisitIds, visitID)
	}
	user.VisitIDs = newUserVisitIds

	if err := wb.db.CreateOrUpdate(user); err != nil {
		log.Printf("error when updating visit relates (onUpdated): %s", err)
	}
}
