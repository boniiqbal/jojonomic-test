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

	"github.com/topup-storage-service/config"
	"github.com/topup-storage-service/src/dto"
	"github.com/topup-storage-service/src/model"
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
		var params dto.RequestTopup
		if err := json.Unmarshal(m.Value, &params); err != nil {
			fmt.Printf("unmarshall data error : %s", err.Error())
		}

		resp := model.Rekening{}
		db := s.config.DB().Master().SelectRaw(`SELECT * from rekening`).Where("norek = ?", params.Norek)
		_, err = db.LoadContext(context.Background(), &resp)
		if err != nil {
			fmt.Printf("Failed error : %s", err.Error())
		}

		harga := model.Harga{}
		dbHarga := s.config.DB().Master().SelectRaw(`SELECT * from harga`).Where("topup_price = ?", params.Harga)
		_, err = dbHarga.LoadContext(context.Background(), &resp)
		if err != nil {
			fmt.Printf("Failed error : %s", err.Error())
		}

		if resp.ID != "" {
			saldo := resp.Saldo + params.Gram
			_, err = s.config.DB().Master().Update(`rekening`).
				Set(`saldo`, saldo).
				Where(`norek = ?`, params.Norek).Exec()
			if err != nil {
				fmt.Printf("Failed error : %s", err.Error())
			}

			_, err = s.config.DB().Master().InsertInto(`transaksi`).Columns(
				`rekening_id`,
				`gram`,
				`type`,
				`topup_price`,
				`buyback_price`,
				`created_at`,
			).Values(resp.ID, params.Gram, "topup", harga.TopupPrice, harga.BuybackPrice, time.Now().Unix()).ExecContext(ctx)

			if err != nil {
				fmt.Printf("Failed error : %s", err.Error())
			}
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
		MinBytes: 1e3,  // 1KB
		MaxBytes: 10e6, // 10MB
	})
}
