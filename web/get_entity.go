package web

import (
	"github.com/valyala/fasthttp"
)

// GetEntityHandler is a handler for /:entity/:id
func (wb *Web) GetEntityHandler(ctx *fasthttp.RequestCtx, idStr string, entityStr string) {

	entity, err := parseEntity(entityStr)
	if err != nil {
		responseErr(ctx, err)
		return
	}

	id, err := parseId(idStr)
	if err != nil {
		ctx.NotFound()
		return
	}

	result, err := wb.db.Get(entity, id)
	if err != nil {
		ctx.NotFound()
		return
	}

	responseJson(ctx, result)
}
