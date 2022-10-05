package src

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/check-harga-service/config"
	"github.com/check-harga-service/shared"
	"github.com/check-harga-service/src/dto"
	"github.com/check-harga-service/src/model"
)

// Server ...
type Server struct {
	httpServer *mux.Router
	config     *config.Config
}

// InitServer ...
func InitServer(cfg *config.Config) *Server {
	r := mux.NewRouter()

	return &Server{httpServer: r, config: cfg}
}

// Run ...
func (s *Server) Run() {
	s.httpServer.HandleFunc("/api/check-harga", s.detailGold).Methods("GET")

	srv := &http.Server{
		Handler:      s.httpServer,
		Addr:         fmt.Sprintf(":%d", s.config.Port()),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func (s *Server) detailGold(w http.ResponseWriter, r *http.Request) {
	resp := model.Harga{}
	db := s.config.DB().Master().SelectRaw(`SELECT * from harga`)
	_, err := db.LoadContext(context.Background(), &resp)
	if err != nil {
		shared.RespondWithError(w, http.StatusBadRequest, true, "", "Harga not found")
		return
	}

	result := dto.ResponseDetailGold{
		HargaBuyback: resp.BuybackPrice,
		HargaTopup:   resp.TopupPrice,
	}

	shared.RespondSuccess(w, http.StatusOK, true, "", result)

}
