package web

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/zenazn/goji/web"

	"github.com/ngalayko/highloadcup/config"
	"github.com/ngalayko/highloadcup/database"
	"github.com/ngalayko/highloadcup/schema"
	"github.com/ngalayko/highloadcup/views"
	"fmt"
)

const (
	webCtxKey ctxKey = "ctx_key_for_web"
)

type ctxKey string

// Web is a wb service
type Web struct {
	db *database.DB
	views *views.Views

	server *http.Server
}

// NewContext saves wb in given context
func NewContext(ctx context.Context, web interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if _, ok := web.(*Web); !ok {
		web = NewWeb(ctx)
	}

	return context.WithValue(ctx, webCtxKey, web)
}

// FromContext return wen from context
func FromContext(ctx context.Context) *Web {
	if service, ok := ctx.Value(webCtxKey).(*Web); ok {
		return service
	}

	return NewWeb(ctx)
}

// NewWeb creates new wb
func NewWeb(ctx context.Context) *Web {

	w := &Web{
		db: database.FromContext(ctx),
		views: views.FromContext(ctx),

		server: &http.Server{
			Addr: config.FromContext(ctx).ListenUrl,
		},
	}

	w.server.Handler = w.initMux()

	return w
}

func (wb *Web) initMux() http.Handler {
	mux := web.New()

	mux.Get("/locations/:id/avg", wb.GetLocationsAvgHandler)
	mux.Get("/users/:id/visits", wb.GetVisitsHandler)
	mux.Get("/:entity/:id", wb.GetEntityHandler)

	mux.Post("/:entity/new", wb.NewEntityHandler)
	mux.Post("/:entity/:id", wb.UpdateEntityHandler)

	return mux
}

func (wb *Web) ServeHTTP() error {

	log.Println("listening on:", wb.server.Addr)
	if err := wb.server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func responseErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)

	responseJson(w, struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	})
}

func responseJson(w http.ResponseWriter, val interface{}) {
	w.Header().Set("Content-type", "application/json")

	data, err := json.Marshal(val)
	if err != nil {
		responseErr(w, err)
		return
	}

	w.Write(data)
}

func parseId(c web.C) (uint32, error) {
	id, err := strconv.ParseInt(c.URLParams["id"], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("erorr parsing id: %s", err)
	}

	return uint32(id), nil
}

func parseEntity(c web.C) (entity schema.Entity, err error) {
	if err := entity.UnmarshalText([]byte(c.URLParams["entity"])); err != nil {
		return schema.EntityUndefined, fmt.Errorf("erorr parsing entity: %s", err)
	}

	return entity, nil
}

func parseToDistance(r *http.Request) (uint32, error) {
	distance, err := strconv.ParseInt(r.URL.Query().Get("toDistance"), 10, 32)
	if err != nil {
		return 0, fmt.Errorf("erorr parsing toDistance: %s", err)
	}
	return uint32(distance), nil
}

func parseCountry(r *http.Request) string {
	return r.URL.Query().Get("country")
}

func parseToDate(r *http.Request) (int64, error) {
	date, err := strconv.ParseInt(r.URL.Query().Get("toDate"), 10, 32)
	if err != nil {
		return 0, fmt.Errorf("erorr parsing toDate: %s", err)
	}
	return date, nil
}

func parseFromDate(r *http.Request) (int64, error) {
	date, err := strconv.ParseInt(r.URL.Query().Get("fromDate"), 10, 32)
	if err != nil {
		return 0, fmt.Errorf("erorr parsing fromDateDate: %s", err)
	}
	return date, nil
}

func parseFromAge(r *http.Request) (int, error) {
	age, err := strconv.ParseInt(r.URL.Query().Get("fromAge"), 10, 32)
	if err != nil {
		return 0, fmt.Errorf("erorr parsing fromAge: %s", err)
	}
	return int(age), nil
}

func parseToAge(r *http.Request) (int, error) {
	age, err := strconv.ParseInt(r.URL.Query().Get("toAge"), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("erorr parsing toAge: %s", err)
	}
	return int(age), err
}

func parseGender(r *http.Request) (gender schema.Gender, err error) {
	if err := gender.UnmarshalText([]byte(r.URL.Query().Get("gender"))); err != nil {
		return schema.GenderUndefined, fmt.Errorf("erorr parsing gender: %s", err)
	}
	return gender, err
}
