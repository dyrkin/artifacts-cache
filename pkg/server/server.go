package server

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"gitlab-cache/pkg/repository"
	"io"
	"net/http"
)

type server struct {
	port       int
	repository repository.Repository
}

func NewServer(port int, repository repository.Repository) *server {
	return &server{port: port, repository: repository}
}

func (s *server) Serve() error {
	http.HandleFunc("/pull", s.pullHandler)
	http.HandleFunc("/push", s.pushHandler)

	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}

func (s *server) pushHandler(w http.ResponseWriter, r *http.Request) {
	subset := r.Header.Get("subset")
	name := r.Header.Get("name")
	if subset == "" || name == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "name or subset header is empty")
		return
	}
	content := r.Body
	defer content.Close()
	err := s.repository.WriteContent(subset, name, content)
	if err != nil {
		msg := fmt.Sprintf("can't save content for name [%s] to repository. error [%s]", name, err)
		log.Error().Msg(msg)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, msg)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *server) pullHandler(w http.ResponseWriter, r *http.Request) {
	subset := r.URL.Query().Get("subset")
	filter := r.URL.Query().Get("filter")
	if subset == "" || filter == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "filter or subset query parameter is empty")
		return
	}
	content, err := s.repository.FindContent(subset, filter)
	if err != nil && errors.Is(err, repository.NoContentOnServerError) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		msg := fmt.Sprintf("can't find content content in subset [%s] by filter [%s] to repository. error [%s]", subset, filter, err)
		log.Error().Msg(msg)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, msg)
		return
	} else {
		defer content.Close()
		_, _ = io.Copy(w, content)
	}
}
