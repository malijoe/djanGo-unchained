package django

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
)

const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
	MIMEPROTOBUF          = "application/x-protobuf"
	MIMEMSGPACK           = "application/x-msgpack"
	MIMEMSGPACK2          = "application/msgpack"
	MIMEYAML              = "application/x-yaml"
	MIMETOML              = "application/toml"
)

var (
	ErrorNoBody = errors.New("invalid request: no body")
)

type Decoder func(any) error

func (d Decoder) Decode(i any) error {
	return d(i)
}

type Parser interface {
	MediaTypes() []string
	Parse(*http.Request) (Decoder, error)
}

type BodyParser interface {
	Parser
	ParseBody([]byte) Decoder
}

type UriParser interface {
	MediaTypes() []string
	ParseUri(map[string][]string) Decoder
}

type JSONParser struct{}

func (JSONParser) Parse(request *http.Request) (Decoder, error) {
	if request == nil || request.Body == nil {
		return nil, ErrorNoBody
	}
	return getJSONDecoder(request.Body), nil
}

func (JSONParser) ParseBody(stream []byte) Decoder {
	return getJSONDecoder(bytes.NewReader(stream))
}

func (JSONParser) MediaTypes() []string {
	return []string{MIMEJSON}
}

func getJSONDecoder(r io.Reader) Decoder {
	decoder := json.NewDecoder(r)
	return func(i any) error {
		return decoder.Decode(i)
	}
}

type XMLParser struct{}

func (XMLParser) MediaTypes() []string {
	return []string{MIMEXML, MIMEXML2}
}

func (XMLParser) Parse(request *http.Request) (Decoder, error) {
	if request == nil || request.Body == nil {
		return nil, ErrorNoBody
	}
	return getXMLDecoder(request.Body), nil
}

func (XMLParser) ParseBody(stream []byte) Decoder {
	return getXMLDecoder(bytes.NewReader(stream))
}

func getXMLDecoder(r io.Reader) Decoder {
	decoder := xml.NewDecoder(r)
	return func(i any) error {
		return decoder.Decode(i)
	}
}
