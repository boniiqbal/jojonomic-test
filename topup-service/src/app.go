package src

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	kafka "github.com/segmentio/kafka-go"
	"github.com/teris-io/shortid"

	"github.com/input-harga-storage-service/config"
	"github.com/input-harga-storage-service/shared"
	"github.com/input-harga-storage-service/src/dto"
	"github.com/input-harga-storage-service/src/model"
)

// Server ...
type Server struct {
	httpServer *mux.Router
	config     *config.Config
	kafka      *kafka.Conn
}

// InitServer ...
func InitServer(cfg *config.Config) *Server {
	r := mux.NewRouter()

	kafkaWriter := createKafkaConn(cfg.KafkaUrl(), cfg.KafkaTopic())

	return &Server{httpServer: r, config: cfg, kafka: kafkaWriter}
}

// Run ...
func (s *Server) Run() {
	s.httpServer.HandleFunc("/api/topup", s.createTopup).Methods("POST")

	srv := &http.Server{
		Handler:      s.httpServer,
		Addr:         fmt.Sprintf(":%d", s.config.Port()),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func (s *Server) createTopup(w http.ResponseWriter, r *http.Request) {
	params := dto.RequestTopup{}
	decoder := json.NewDecoder(r.Body)
	reffID, _ := shortid.Generate()
	if err := decoder.Decode(&params); err != nil {
		shared.RespondWithError(w, http.StatusBadRequest, true, reffID, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	resp := new(model.Harga)
	db := s.config.DB().Master().SelectRaw(`SELECT * from harga`).Where("topup_price = ?", params.Harga)

	_, err := db.LoadContext(context.Background(), &resp)
	if err != nil {
		shared.RespondWithError(w, http.StatusBadRequest, true, "", err.Error())
		return
	}

	if resp == nil {
		shared.RespondWithError(w, http.StatusBadRequest, true, "", "Harga tidak ditemukan")
		return
	}

	payloadBytes, err := json.Marshal(&params)
	if err != nil {
		shared.RespondWithError(w, http.StatusBadRequest, true, reffID, err.Error())
		return
	}

	s.kafka.SetWriteDeadline(time.Now().Add(10 * time.Second))
	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("address-%s", r.RemoteAddr)),
		Value: payloadBytes,
	}
	_, err = s.kafka.WriteMessages(msg)
	if err != nil {
		log.Println(err.Error())
		shared.RespondWithError(w, http.StatusBadRequest, true, reffID, "Kafka not ready")
		return
	}

	shared.RespondSuccess(w, http.StatusOK, false, reffID, nil)
}

func createKafkaConn(kafkaURL, topic string) *kafka.Conn {
	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaURL, topic, 0)
	if err != nil {
		log.Fatal(err.Error())
	}

	return conn
}
