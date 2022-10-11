package src

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	kafka "github.com/segmentio/kafka-go"
	"github.com/teris-io/shortid"

	"github.com/input-harga-storage-service/config"
	"github.com/input-harga-storage-service/src/model"
)

// Server ...
type Server struct {
	config *config.Config
	kafka  *kafka.Reader
}

// InitServer ...
func InitServer(cfg *config.Config) *Server {

	kafkaReader := getKafkaReader(os.Getenv("KAFKA_URL"), os.Getenv("KAFKA_TOPIC"), os.Getenv("KAFKA_GROUP_ID"))

	return &Server{config: cfg, kafka: kafkaReader}
}

// Run ...
func (s *Server) Run() {
	ctx := context.Background()
	for {
		m, err := s.kafka.FetchMessage(ctx)
		if err != nil {
			break
		}
		log.Printf("message at topic/partition/offset %v/%v/%v: %s\n", m.Topic, m.Partition, m.Offset, string(m.Key))
		var harga model.Harga
		var request model.RequestInputHarga
		if err := json.Unmarshal(m.Value, &request); err != nil {
			fmt.Printf("unmarshall data error : %s", err.Error())
		}

		reffID, _ := shortid.Generate()
		harga.ID = reffID
		harga.TopupPrice = int64(request.HargaTopup)
		harga.BuybackPrice = int64(request.HargaBuyback)
		harga.UserID = request.AdminID
		harga.CreatedAt = time.Now()
		_, err = s.config.DB().Master().InsertInto(`harga`).Columns(
			`topup_price`,
			`buyback_price`,
			`id`,
			`user_id`,
			`created_at`,
		).Record(harga).ExecContext(ctx)

		if err != nil {
			fmt.Printf("Failed error : %s", err.Error())
		}

		if err := s.kafka.CommitMessages(ctx, m); err != nil {
			log.Fatal("failed to commit messages:", err)
		}
	}
}

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 0,    // 1KB
		MaxBytes: 10e6, // 10MB
	})
}
