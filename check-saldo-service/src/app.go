package src

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/check-saldo-service/config"
	"github.com/check-saldo-service/shared"
	"github.com/check-saldo-service/src/dto"
	"github.com/check-saldo-service/src/model"
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
	s.httpServer.HandleFunc("/api/saldo", s.detailSaldo).Methods("POST")

	srv := &http.Server{
		Handler:      s.httpServer,
		Addr:         fmt.Sprintf(":%d", s.config.Port()),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func (s *Server) detailSaldo(w http.ResponseWriter, r *http.Request) {
	params := dto.RequestSaldo{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		shared.RespondWithError(w, http.StatusBadRequest, true, "", "Invalid reques payload")
		return
	}
	defer r.Body.Close()

	resp := model.Rekening{}
	db := s.config.DB().Master().SelectRaw(`SELECT * from rekening`).Where("norek = ?", params.Norek)

	_, err := db.LoadContext(context.Background(), &resp)
	if err != nil {
		shared.RespondWithError(w, http.StatusBadRequest, true, "", "Saldo not found")
		return
	}

	result := dto.ResponseSaldo{
		Norek: resp.Norek,
		Saldo: resp.Saldo,
	}

	shared.RespondSuccess(w, http.StatusOK, false, "", result)
}
