package goutils

import (
	"io/fs"
	"net/http"
)

func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteNotFound(w, nil)
	})
}

type SPAHandler struct {
	fsys            fs.FS
	indexHTML       []byte
	fileServer      http.Handler
	prepareResponse func(http.ResponseWriter, *http.Request, bool)
}

func NewSPAHandler(fsys fs.FS, dirPath string, indexPath string, prepareResponse func(http.ResponseWriter, *http.Request, bool)) (*SPAHandler, error) {
	if dirPath != "" {
		var err error
		if fsys, err = fs.Sub(fsys, dirPath); err != nil {
			return nil, err
		}
	}

	indexHTML, err := fs.ReadFile(fsys, indexPath)
	if err != nil {
		return nil, err
	}

	if prepareResponse == nil {
		prepareResponse = func(http.ResponseWriter, *http.Request, bool) {}
	}

	return &SPAHandler{
		fsys:            fsys,
		indexHTML:       indexHTML,
		fileServer:      http.FileServer(http.FS(fsys)),
		prepareResponse: prepareResponse,
	}, nil
}

func (h *SPAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if file, err := fs.Stat(h.fsys, r.URL.Path); err != nil || file.IsDir() {
		h.prepareResponse(w, r, true)
		w.WriteHeader(http.StatusOK)
		w.Write(h.indexHTML)
		return
	}
	h.prepareResponse(w, r, false)
	h.fileServer.ServeHTTP(w, r)
}
