package web

import (
	"sort"

	"github.com/valyala/fasthttp"

	"github.com/ngalayko/highloadcup/schema"
)

// GetVisitsHandler is handler for /users/:id/visits
func (wb *Web) GetVisitsHandler(ctx *fasthttp.RequestCtx) {
	id, err := parseId(ctx)
	if err != nil {
		responseErr(ctx, err)
		return
	}

	country := parseCountry(ctx)
	fromDate, err := parseFromDate(ctx)
	if err != nil && fromDate == 0 {
		responseErr(ctx, err)
		return
	}
	toDate, err := parseToDate(ctx)
	if err != nil && toDate == 0 {
		responseErr(ctx, err)
		return
	}
	toDistance, err := parseToDistance(ctx)
	if err != nil && toDistance == 0 {
		responseErr(ctx, err)
		return
	}

	user, err := wb.db.GetUser(id)
	if err != nil {
		ctx.NotFound()
		return
	}

	visits, err := wb.db.GetVisits(user.VisitIDs)
	if err != nil {
		responseErr(ctx, err)
		return
	}

	var result []*schema.Visit
	for _, visit := range visits {
		location, err := wb.db.GetLocation(visit.LocationID)
		if err != nil {
			responseErr(ctx, err)
			return
		}

		if fromDate != 0 && visit.VisitedAt <= fromDate {
			continue
		}

		if toDate != 0 && visit.VisitedAt >= toDate {
			continue
		}

		if country != "" && location.Country != country {
			continue
		}

		if toDistance != 0 && location.Distance >= toDistance {
			continue
		}

		result = append(result, visit)
	}

	if len(result) == 0 {
		result = []*schema.Visit{}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].VisitedAt < result[j].VisitedAt
	})

	views, err := wb.views.FillVisitsViews(result)
	if err != nil {
		responseErr(ctx, err)
		return
	}

	responseJson(ctx, views)
}
