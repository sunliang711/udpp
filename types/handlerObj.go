package types

import "net/http"

type HandlerObj struct {
	Path    string
	H       http.Handler
	Methods []string
	Usage   string
}
