package handler

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/shreyansh-ML/movieapp/file/internal/storage"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

type Files struct {
	log   hclog.Logger
	store storage.Storage
}

// NewFiles creates a new File handler
func NewFiles(s storage.Storage, l hclog.Logger) *Files {
	return &Files{store: s, log: l}
}

func (f *Files) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	//in_ur, _ := url.Parse(r.URL.Path)
	vars := mux.Vars(r)
	id := vars["id"]
	fn := vars["filename"]

	f.log.Info("Handle POST", "id", id, "filename", fn)

	// no need to check for invalid id or filename as the mux router will not send requests
	// here unless they have the correct parameters

	f.saveFile(id, fn, rw, r)
}

func (f *Files) saveFile(id, path string, rw http.ResponseWriter, r *http.Request) {
	f.log.Info("Save file for product", "id", id, "path", path)
	ctx := context.WithValue(context.Background(), "id", id)
	fp := filepath.Join(id, path)
	err := f.store.Save(ctx, fp, r.Body)
	if err != nil {
		f.log.Error("Unable to save file", "error", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
	}
}
