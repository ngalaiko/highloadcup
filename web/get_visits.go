package web

import (
	"net/http"
	"sort"

	"github.com/zenazn/goji/web"

	"github.com/ngalayko/highloadcup/schema"
	"fmt"
)

// GetVisitsHandler is handler for /users/:id/visits
func (wb *Web) GetVisitsHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	id, err := parseId(c)
	if err != nil {
		responseErr(w, err)
		return
	}

	fromDate, err := parseFromDate(r)
	if err != nil {
		responseErr(w, err)
		return
	}

	toDate, err := parseToDate(r)
	if err != nil {
		responseErr(w, err)
		return
	}

	country := parseCountry(r)

	toDistance, err := parseToDistance(r)
	fmt.Println(toDistance, err)
	if err != nil {
		responseErr(w, err)
		return
	}

	user, err := wb.db.GetUser(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	visits, err := wb.db.GetVisits(user.VisitIDs)
	if err != nil {
		responseErr(w, err)
		return
	}

	var result []*schema.Visit
	for _, visit := range visits {
		location, err := wb.db.GetLocation(visit.LocationID)
		if err != nil {
			responseErr(w, err)
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
		responseErr(w, err)
		return
	}

	responseJson(w, views)
}
