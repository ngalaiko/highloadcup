package web

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/valyala/fasthttp"

	"github.com/ngalayko/highloadcup/schema"
)

// NewEntityHandler is a handler for /:entity/new
func (wb *Web) NewEntityHandler(ctx *fasthttp.RequestCtx) {
	entity, err := parseEntity(ctx)
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

	if _, err := wb.db.Get(entity, val.IntID()); err == nil {
		responseErr(ctx, fmt.Errorf("entity already exists"))
		return
	}

	if err := wb.db.CreateOrUpdate(val); err != nil {
		responseErr(ctx, err)
		return
	}

	if newVisit, ok := val.(*schema.Visit); ok {
		go wb.onVisitInserted(newVisit)
	}

	responseJson(ctx, struct{}{})
}

func (wb *Web) onVisitInserted(visit *schema.Visit) {
	wb.onVisitInsertedUpdateLocation(visit)
	wb.onVisitInsertedUpdateUser(visit)
}

func (wb *Web) onVisitInsertedUpdateLocation(visit *schema.Visit) {
	location, err := wb.db.GetLocation(visit.UserID)
	if err != nil {
		log.Printf("error when updating visit relates (onInserted): %s", err)
		return
	}

	location.VisitIDs = append(location.VisitIDs, visit.ID)
	if err := wb.db.CreateOrUpdate(location); err != nil {
		log.Printf("error when updating visit relates (onInserted): %s", err)
	}
}

func (wb *Web) onVisitInsertedUpdateUser(visit *schema.Visit) {
	user, err := wb.db.GetUser(visit.UserID)
	if err != nil {
		log.Printf("error when updating visit relates (onInserted): %s", err)
		return
	}

	user.VisitIDs = append(user.VisitIDs, visit.ID)
	if err := wb.db.CreateOrUpdate(user); err != nil {
		log.Printf("error when updating visit relates (onInserted): %s", err)
	}
}
