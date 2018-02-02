package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/heroku/x/hkafka"
)

var (
	kafkaProducer *kafka.Producer
	kafkaTopic    string
)

func main() {
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "8080"
	}

	kafkaTopic = os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		log.Fatal("KAFAK_TOPIC is empty or missing")
	}

	hkc, err := hkafka.NewConfigFromEnv()
	if err != nil {
		log.Fatal("error reading kafka config: ", err)
	}

	bl, err := hkc.BrokerAddresses()
	if err != nil {
		log.Fatal("error constructing broker list: ", err)
	}

	if err := hkc.WriteDefaultSSLFiles(); err != nil {
		log.Fatal(err)
	}

	h := os.Getenv("HOME")

	kafkaProducer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(bl, ","),
		"security.protocol":        "ssl",
		"ssl.key.location":         filepath.Join(h, hkafka.DefaultClientCertKeyFileName),
		"ssl.certificate.location": filepath.Join(h, hkafka.DefaultClientCertFileName),
		"ssl.ca.location":          filepath.Join(h, hkafka.DefaultRootCAFileName),
		"socket.keepalive.enable":  "true",
	})
	if err != nil {
		log.Fatal("error creating kafka producer: ", err)
	}

	http.HandleFunc("/webhooks", handleHook)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Work In Progress\n")
		io.WriteString(w, "\n")
		io.WriteString(w, "See github.com/freeformz/github-projector\n")
	})

	http.ListenAndServe(":"+port, nil)
}
