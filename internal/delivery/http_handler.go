package delivery

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/cuongtranba/mynoti/internal/domain"
	"github.com/cuongtranba/mynoti/internal/usecase"
)

type Server struct {
	server *http.Server
}

type Comic struct {
	Url         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Html        string `json:"html"`
}

func NewServer(port string, comicUseCase usecase.ComicUseCase) *Server {
	mux := http.NewServeMux()
	mux.Handle("POST /subscribe", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var comic Comic
		if err := json.NewDecoder(r.Body).Decode(&comic); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err := comicUseCase.Subscribe(r.Context(), &domain.Comic{
			Url:         comic.Url,
			Name:        comic.Name,
			Description: comic.Description,
			Html:        comic.Html,
		})
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	return &Server{
		server: &http.Server{
			Addr:    port,
			Handler: mux,
		},
	}
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
