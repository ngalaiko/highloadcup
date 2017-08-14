package web

import (
	"time"

	"github.com/valyala/fasthttp"

	"github.com/ngalayko/highloadcup/helper"
	"github.com/ngalayko/highloadcup/schema"
)

// GetLocationsAvgHandler is a handler for /locations/:id/avg
func (wb *Web) GetLocationsAvgHandler(ctx *fasthttp.RequestCtx) {
	id, err := parseId(ctx)
	if err != nil {
		responseErr(ctx, err)
		return
	}

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

	fromAge, err := parseFromAge(ctx)
	if err != nil && fromAge == 0 {
		responseErr(ctx, err)
		return
	}

	toAge, err := parseToAge(ctx)
	if err != nil && toAge == 0 {
		responseErr(ctx, err)
		return
	}

	gender, err := parseGender(ctx)
	if err != nil && gender == schema.GenderUndefined {
		responseErr(ctx, err)
		return
	}

	location, err := wb.db.GetLocation(id)
	if err != nil {
		ctx.NotFound()
		return
	}

	visits, err := wb.db.GetVisits(location.VisitIDs)
	if err != nil {
		responseErr(ctx, err)
		return
	}
	var userIds []uint32
	for _, visit := range visits {
		userIds = append(userIds, visit.UserID)
	}

	var marks []uint8
	for _, visit := range visits {
		user, err := wb.db.GetUser(visit.UserID)
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

		if gender != schema.GenderUndefined && user.Gender != gender {
			continue
		}

		if fromAge != 0 && time.Now().Year()-time.Unix(user.BirthDate, 0).Year() <= fromAge {
			continue
		}

		if toAge != 0 && time.Now().Year()-time.Unix(user.BirthDate, 0).Year() >= toAge {
			continue
		}

		marks = append(marks, visit.Mark)
	}

	if len(marks) > 0 {
		responseJson(ctx, struct {
			Avg float64 `json:"avg"`
		}{
			Avg: helper.Avg(marks...),
		})
		return
	}

	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Write([]byte(`{"avg": 0.0}`))
}
