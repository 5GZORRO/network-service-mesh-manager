// Package Nsmm provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version (devel) DO NOT EDIT.
package Nsmm

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"AQAA////LUnvJlcAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
