package web

import (
	"net/http"

	"github.com/ngalayko/highloadcup/helper"
	"github.com/ngalayko/highloadcup/schema"
	"github.com/zenazn/goji/web"
	"time"
)

// GetLocationsAvgHandler is a handler for /locations/:id/avg
func (wb *Web) GetLocationsAvgHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	id, err := parseId(c)
	if err != nil {
		responseErr(w, err)
		return
	}

	fromDate := parseFromDate(r)
	toDate := parseToDate(r)
	fromAge := parseFromAge(r)
	toAge := parseToAge(r)
	gender := parseGender(r)

	location, err := wb.db.GetLocation(id)
	if err != nil {
		responseErr(w, err)
		return
	}

	visits, err := wb.db.GetVisits(location.VisitIDs)
	if err != nil {
		responseErr(w, err)
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
			responseErr(w, err)
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
		responseJson(w, struct {
			Avg float64 `json:"avg"`
		}{
			Avg: helper.Avg(marks...),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"avg": 0.0}"`))
}
