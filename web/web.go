package web

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/ngalayko/highloadcup/config"
	"github.com/ngalayko/highloadcup/database"
	"github.com/ngalayko/highloadcup/schema"
	"github.com/zenazn/goji/web"
)

const (
	webCtxKey ctxKey = "ctx_key_for_web"
)

type ctxKey string

// Web is a wb service
type Web struct {
	db *database.DB

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
		return 0, err
	}

	return uint32(id), nil
}

func parseEntity(c web.C) (schema.Entity, error) {
	var entity schema.Entity

	if err := entity.UnmarshalText([]byte(c.URLParams["entity"])); err != nil {
		return entity, err
	}

	return entity, nil
}
