package web

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/valyala/fasthttp"

	"github.com/ngalayko/highloadcup/config"
	"github.com/ngalayko/highloadcup/database"
	"github.com/ngalayko/highloadcup/schema"
	"github.com/ngalayko/highloadcup/views"
	"strconv"
	"strings"
)

const (
	webCtxKey ctxKey = "ctx_key_for_web"
)

var (
	getEntityRegex    = regexp.MustCompile(`^/(?P<entity>\w+)/(?P<id>\d+)$`)
	getVisitsRegex    = regexp.MustCompile(`^/users/(?P<id>\d+)/visits$`)
	locationsAvgRegex = regexp.MustCompile(`^/locations/(?P<id>\d+)/avg$`)
	newEntityRegex    = regexp.MustCompile(`^/(?P<entity>\w+)/new$`)
)

type ctxKey string

// Web is a wb service
type Web struct {
	db     *database.DB
	views  *views.Views
	config *config.Config
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
		db:     database.FromContext(ctx),
		views:  views.FromContext(ctx),
		config: config.FromContext(ctx),
	}

	return w
}

func (wb *Web) HandleFastHttp(ctx *fasthttp.RequestCtx) {
	log.Printf("method: %s path: %s query: %v body: %v", ctx.Method(), ctx.Path(), ctx.QueryArgs(), ctx.Request.Body())

	parts := strings.Split(string(ctx.Path()), "/")

	fmt.Println(parts[1], parts[2])
	switch {
	case locationsAvgRegex.Match(ctx.Path()):

		if !ctx.IsGet() {
			ctx.NotFound()
		}

		wb.GetLocationsAvgHandler(ctx, parts[2])

	case getVisitsRegex.Match(ctx.Path()):

		if !ctx.IsGet() {
			ctx.NotFound()
		}

		wb.GetVisitsHandler(ctx, parts[2])

	case getEntityRegex.Match(ctx.Path()):

		switch {
		case ctx.IsGet():
			wb.GetEntityHandler(ctx, parts[2], parts[1])

		case ctx.IsPost():
			wb.NewEntityHandler(ctx, parts[1])

		default:
			ctx.NotFound()
		}

	case newEntityRegex.Match(ctx.Path()):
		if !ctx.IsPost() {
			ctx.NotFound()
		}

		wb.UpdateEntityHandler(ctx, parts[2], parts[1])

	default:
		ctx.NotFound()
	}
}

func (wb *Web) ServeHTTP() error {

	log.Println("listening on:", wb.config.ListenUrl)
	if err := fasthttp.ListenAndServe(wb.config.ListenUrl, wb.HandleFastHttp); err != nil {
		log.Printf("listenAndServe err: %s", err)
	}

	return nil
}

func responseErr(ctx *fasthttp.RequestCtx, err error) {
	ctx.SetStatusCode(http.StatusBadRequest)

	responseJson(ctx, struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	})
}

func responseJson(ctx *fasthttp.RequestCtx, val interface{}) {
	ctx.Response.Header.Set("Content-type", "application/json")

	data, err := json.Marshal(val)
	if err != nil {
		responseErr(ctx, err)
		return
	}

	ctx.Write(data)
}

func parseId(str string ) (uint32, error) {
	id, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("erorr parsing id: %s", err)
	}

	return uint32(id), nil
}

func parseEntity(str string) (entity schema.Entity, err error) {
	if err := entity.UnmarshalText([]byte(str)); err != nil {
		return schema.EntityUndefined, fmt.Errorf("erorr parsing entity: %s", err)
	}

	return entity, nil
}

func parseToDistance(ctx *fasthttp.RequestCtx) (uint32, error) {
	toDistance, err := ctx.QueryArgs().GetUint("toDistance")
	switch err {
	case nil:
		return uint32(toDistance), nil
	case fasthttp.ErrNoArgValue:
		return 0, nil
	default:
		return 0, err
	}
}

func parseCountry(ctx *fasthttp.RequestCtx) string {
	return string(ctx.QueryArgs().Peek("country"))
}

func parseToDate(ctx *fasthttp.RequestCtx) (int64, error) {
	toDate, err := ctx.QueryArgs().GetUint("toDate")
	switch err {
	case nil:
		return int64(toDate), nil
	case fasthttp.ErrNoArgValue:
		return 0, nil
	default:
		return 0, err
	}
}

func parseFromDate(ctx *fasthttp.RequestCtx) (int64, error) {
	fromDate, err := ctx.QueryArgs().GetUint("fromDate")
	switch err {
	case nil:
		return int64(fromDate), nil
	case fasthttp.ErrNoArgValue:
		return 0, nil
	default:
		return 0, err
	}
}

func parseFromAge(ctx *fasthttp.RequestCtx) (int, error) {
	fromAge, err := ctx.QueryArgs().GetUint("fromAge")
	switch err {
	case nil:
		return fromAge, nil
	case fasthttp.ErrNoArgValue:
		return 0, nil
	default:
		return 0, err
	}
}

func parseToAge(ctx *fasthttp.RequestCtx) (int, error) {
	toAge, err := ctx.QueryArgs().GetUint("toAge")
	switch err {
	case nil:
		return toAge, nil
	case fasthttp.ErrNoArgValue:
		return 0, nil
	default:
		return 0, err
	}
}

func parseGender(ctx *fasthttp.RequestCtx) (schema.Gender, error) {
	genderBytes := ctx.QueryArgs().Peek("gender")
	switch len(genderBytes) {
	case 0:
		return schema.GenderUndefined, nil
	default:
		var gender schema.Gender
		if err := gender.UnmarshalText(genderBytes); err != nil {
			return schema.GenderUndefined, fmt.Errorf("erorr parsing gender: %s", err)
		}

		return gender, nil
	}
}
