package src

import (
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
)

// Server ...
type Server struct {
	httpServer *mux.Router
	config     *config.Config
	kafka      *kafka.Writer
}

// InitServer ...
func InitServer(cfg *config.Config) *Server {
	r := mux.NewRouter()

	kafkaWriter := getKafkaWriter(cfg.KafkaUrl(), cfg.KafkaTopic())
	defer kafkaWriter.Close()

	return &Server{httpServer: r, config: cfg, kafka: kafkaWriter}
}

// Run ...
func (s *Server) Run() {
	s.httpServer.HandleFunc("/api/input-harga", s.createHarga).Methods("POST")

	srv := &http.Server{
		Handler:      s.httpServer,
		Addr:         fmt.Sprintf(":%d", s.config.Port()),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func (s *Server) createHarga(w http.ResponseWriter, r *http.Request) {
	params := dto.RequestInputHarga{}
	decoder := json.NewDecoder(r.Body)
	reffID, _ := shortid.Generate()
	if err := decoder.Decode(&params); err != nil {
		shared.RespondWithError(w, http.StatusBadRequest, true, reffID, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	payloadBytes, err := json.Marshal(&params)
	if err != nil {
		shared.RespondWithError(w, http.StatusBadRequest, true, reffID, err.Error())
		return
	}

	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("address-%s", r.RemoteAddr)),
		Value: payloadBytes,
	}
	err = s.kafka.WriteMessages(r.Context(), msg)
	if err != nil {
		log.Println(err.Error())
		shared.RespondWithError(w, http.StatusBadRequest, true, reffID, "Kafka not ready")
		return
	}

	shared.RespondSuccess(w, http.StatusBadRequest, false, reffID, nil)
}

func getKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}