package handlers

import (
	"net/http"
	"path"
)

func getBaseURL(req *http.Request, paths ...string) string {
	return path.Join("http://"+req.Host, path.Join(paths...))
}
