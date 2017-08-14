package web

import (
	"encoding/json"
	"log"

	"github.com/valyala/fasthttp"

	"github.com/ngalayko/highloadcup/schema"
)

// GetEntityHandler is a handler for /:entity/:id
func (wb *Web) UpdateEntityHandler(ctx *fasthttp.RequestCtx) {
	entity, err := parseEntity(ctx)
	if err != nil {
		responseErr(ctx, err)
		return
	}

	id, err := parseId(ctx)
	if err != nil {
		responseErr(ctx, err)
		return
	}

	val := schema.GetIEntity(entity)
	if err := json.Unmarshal(ctx.PostBody(), val); err != nil {
		responseErr(ctx, err)
		return
	}

	if err := val.Validate(); err != nil {
		responseErr(ctx, err)
		return
	}

	oldValue, err := wb.db.Get(entity, id)
	if err != nil {
		ctx.NotFound()
		return
	}

	oldVisit := &schema.Visit{}
	if oldVisitPtr, ok := oldValue.(*schema.Visit); ok {
		*oldVisit = *oldVisitPtr
	}

	if err := wb.db.CreateOrUpdate(val); err != nil {
		responseErr(ctx, err)
		return
	}

	if newVisit, ok := val.(*schema.Visit); ok {
		go wb.onVisitUpdated(oldVisit, newVisit)
	}

	responseJson(ctx, struct{}{})
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
