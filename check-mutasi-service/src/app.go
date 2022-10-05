package src

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/check-mutasi-service/config"
	"github.com/check-mutasi-service/shared"
	"github.com/check-mutasi-service/src/dto"
	"github.com/check-mutasi-service/src/model"
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
	s.httpServer.HandleFunc("/api/mutasi", s.detailTransaction).Methods("POST")

	srv := &http.Server{
		Handler:      s.httpServer,
		Addr:         fmt.Sprintf(":%d", s.config.Port()),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func (s *Server) detailTransaction(w http.ResponseWriter, r *http.Request) {
	params := dto.RequestTransaction{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		shared.RespondWithError(w, http.StatusBadRequest, true, "", "Invalid reques payload")
		return
	}
	defer r.Body.Close()

	resp := []model.Transaksi{}
	db := s.config.DB().Master().SelectRaw(`SELECT * from transaksi`)

	saldo := 0.0
	if params.Norek != "" {
		rekeningData := model.Rekening{}
		dbRek := s.config.DB().Master().SelectRaw(`SELECT * from rekening`).Where("norek = ?", params.Norek)
		_, err := dbRek.LoadContext(context.Background(), &rekeningData)
		if err != nil {
			shared.RespondWithError(w, http.StatusBadRequest, true, "", "Harga not found")
			return
		}
		saldo = rekeningData.Saldo
		db = db.Where("rekening_id = ?", rekeningData.ID)
	}

	if params.StartDate != 0 {
		db = db.Where("created at >= ?", params.StartDate)
	}

	if params.EndDate != 0 {
		db = db.Where("created at <= ?", params.EndDate)
	}

	_, err := db.LoadContext(context.Background(), &resp)
	if err != nil {
		shared.RespondWithError(w, http.StatusBadRequest, true, "", "Harga not found")
		return
	}

	result := dto.ResponseTransaction{}
	for _, val := range resp {

		result.Data = append(result.Data, dto.DetailTransaction{
			Date:         val.CreatedAt,
			Type:         val.Type,
			Gram:         val.Gram,
			HargaTopup:   val.BuybackPrice,
			HargaBuyback: val.BuybackPrice,
			Saldo:        saldo,
		})
	}

	shared.RespondSuccess(w, http.StatusOK, true, "", result)

}
