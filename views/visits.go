package views

import (
	"github.com/ngalayko/highloadcup/schema"
)


type VisitsView struct {
	Visits []*VisitView `json:"visits"`
}

type VisitView struct {
	Mark uint8 `json:"mark"`
	VisitedAt int64 `json:"visited_at"`
	Place string `json:"place"`
}

func (v *Views) FillVisitsViews(visits []*schema.Visit) (*VisitsView, error) {

	result := &VisitsView{}

	for _, visit := range visits {

		location, err := v.db.GetLocation(visit.LocationID)
		if err != nil {
			return nil, err
		}

		result.Visits = append(result.Visits, &VisitView{
			Mark: visit.Mark,
			VisitedAt: visit.VisitedAt,
			Place: location.Place,
		})
	}

	return result, nil
}

