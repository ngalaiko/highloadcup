package web

import (
	"net/http"
	"time"

	"github.com/zenazn/goji/web"

	"github.com/ngalayko/highloadcup/helper"
	"github.com/ngalayko/highloadcup/schema"
)

// GetLocationsAvgHandler is a handler for /locations/:id/avg
func (wb *Web) GetLocationsAvgHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	id, err := parseId(c)
	if err != nil {
		responseErr(w, err)
		return
	}

	fromDate, err := parseFromDate(r)
	if err != nil && fromDate != 0 {
		responseErr(w, err)
		return
	}

	toDate, err := parseToDate(r)
	if err != nil && toDate != 0 {
		responseErr(w, err)
		return
	}

	fromAge, err := parseFromAge(r)
	if err != nil && fromDate != 0 {
		responseErr(w, err)
		return
	}

	toAge, err := parseToAge(r)
	if err != nil && toAge != 0 {
		responseErr(w, err)
		return
	}

	gender, err := parseGender(r)
	if err != nil && gender == schema.GenderUndefined {
		responseErr(w, err)
		return
	}

	location, err := wb.db.GetLocation(id)
	if err != nil {
		http.NotFound(w, r)
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
	w.Write([]byte(`{"avg": 0.0}`))
}
