package server

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gitlab-cache/pkg/repository"
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
	key := r.Header.Get("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	content := r.Body
	defer content.Close()
	err := s.repository.WriteContent(key, content)
	if err != nil {
		msg := fmt.Sprintf("can't save content for key [%s] to repository. error [%s]", key, err)
		log.Error().Msg(msg)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, msg)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *server) pullHandler(w http.ResponseWriter, r *http.Request) {

}
