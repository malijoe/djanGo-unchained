package django

import (
	"encoding"
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v2"
)

type HandleFunc[Request any] func(ctx *gin.Context, serializer *Request) (code int, response any, err error)

func (f HandleFunc[Request]) Wrap() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request Request
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		code, response, err := f(ctx, &request)
		if err != nil {
			ctx.Error(err)
		}
		ctx.JSON(code, response)
	}
}

type Opt[R Serializer] func(*Handler[R])

func EnforceContentMatch[R Serializer]() Opt[R] {
	return func(h *Handler[R]) {
		h.contentMustMatch = true
		var r R
		if _, ok := any(r).(json.Unmarshaler); ok {
			h.acceptedContent = append(h.acceptedContent, MIMEJSON)
		}
		if _, ok := any(r).(yaml.Unmarshaler); ok {
			h.acceptedContent = append(h.acceptedContent, MIMEYAML)
		}
		if _, ok := any(r).(encoding.TextUnmarshaler); ok {
			h.acceptedContent = append(h.acceptedContent, MIMEPlain)
		}
		if _, ok := any(r).(xml.Unmarshaler); ok {
			h.acceptedContent = append(h.acceptedContent, MIMEXML)
		}
	}
}

type Handler[R Serializer] struct {
	get, post, put, patch, delete HandleFunc[R]
	contentMustMatch              bool
	acceptedContent               []string
}

func NewHandler[R Serializer](opts ...Opt[R]) *Handler[R] {
	h := new(Handler[R])
	for _, opt := range opts {
		opt(h)
	}
	return h
}

func (h *Handler[R]) contentMiddleware(ctx *gin.Context) {
	if len(h.acceptedContent) > 0 {
		content_type := ctx.GetHeader("Content-Type")
		if !slices.Contains(h.acceptedContent, content_type) {
			ctx.JSON(http.StatusNotAcceptable, "unaccepted content-type")
			return
		}
	}
	ctx.Next()
}

func (h *Handler[R]) Get(handle HandleFunc[R]) *Handler[R] {
	h.get = handle
	return h
}

func (h *Handler[R]) Post(handle HandleFunc[R]) *Handler[R] {
	h.post = handle
	return h
}

func (h *Handler[R]) Put(handle HandleFunc[R]) *Handler[R] {
	h.put = handle
	return h
}

func (h *Handler[R]) Patch(handle HandleFunc[R]) *Handler[R] {
	h.patch = handle
	return h
}

func (h *Handler[R]) Delete(handle HandleFunc[R]) *Handler[R] {
	h.delete = handle
	return h
}

func (h *Handler[R]) options(ctx *gin.Context) {
	metadata := make(map[string]map[string]Field)
	var r R
	sFields := r.Metadata()
	fieldMetadata := make(map[string]Field)
	for _, f := range sFields {
		fieldMetadata[f.Name] = f
	}
	var methods = make([]string, 0, 5)
	if h.get != nil {
		methods = append(methods, http.MethodGet)
	}
	if h.post != nil {
		methods = append(methods, http.MethodPost)
	}
	if h.put != nil {
		methods = append(methods, http.MethodPut)
	}
	if h.patch != nil {
		methods = append(methods, http.MethodPatch)
	}
	if h.delete != nil {
		methods = append(methods, http.MethodDelete)
	}

	for _, method := range methods {
		metadata[method] = fieldMetadata
	}

	ctx.JSON(http.StatusOK, metadata)
}

func (h *Handler[R]) asView(router *gin.RouterGroup) {
	if h.contentMustMatch {
		router.Use(h.contentMiddleware)
	}
	handleMap := map[string]HandleFunc[R]{
		http.MethodGet:    h.get,
		http.MethodPost:   h.post,
		http.MethodPut:    h.put,
		http.MethodPatch:  h.patch,
		http.MethodDelete: h.delete,
	}
	for verb, handler := range handleMap {
		if handler != nil {
			router.Handle(verb, "", handler.Wrap())
		}
	}
}
