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

	country := parseCountry(r)
	fromDate, fromDateErr := parseFromDate(r)
	toDate, toDateErr := parseToDate(r)
	toDistance, toDistanceErr := parseToDistance(r)

	user, err := wb.db.GetUser(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Println(3)

	visits, err := wb.db.GetVisits(user.VisitIDs)
	if err != nil {
		responseErr(w, err)
		return
	}

	fmt.Println(4)

	var result []*schema.Visit
	for _, visit := range visits {
		location, err := wb.db.GetLocation(visit.LocationID)
		if err != nil {
			responseErr(w, err)
			return
		}

		if fromDateErr != nil && visit.VisitedAt <= fromDate {
			continue
		}

		if toDateErr != nil && visit.VisitedAt >= toDate {
			continue
		}

		if country != "" && location.Country != country {
			continue
		}

		if toDistanceErr != nil && location.Distance >= toDistance {
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
