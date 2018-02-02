package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

func respondWithError(w http.ResponseWriter, err error) {
	log.Println("hook error: ", err.Error())
	http.Error(w, err.Error(), 503)
}

func handleHook(w http.ResponseWriter, r *http.Request) {
	var i interface{}
	switch et := r.Header.Get("X-Github-Event"); et {
	case "issues":
		i = github.IssueEvent{}
	case "issue_comment":
		i = github.IssueCommentEvent{}
	case "milestone":
		i = github.MilestoneEvent{}
	case "project":
		i = github.ProjectEvent{}
	case "project_card":
		i = github.ProjectCardEvent{}
	case "project_column":
		i = github.ProjectColumnEvent{}
	case "pull_request":
		i = github.PullRequestEvent{}
	case "pull_request_review":
		i = github.PullRequestReviewEvent{}
	case "pull_request_review_comment":
		i = github.PullRequestReviewCommentEvent{}
	default:
		log.Println("Unknown event type: ", et)
		return
	}

	d := json.NewDecoder(r.Body)
	if err := d.Decode(&i); err != nil {
		respondWithError(w, errors.Wrap(err, "decoding error"))
		return
	}

	b, err := json.Marshal(i)
	if err != nil {
		respondWithError(w, errors.Wrap(err, "encoding error"))
		return
	}

	deliveryChan := make(chan kafka.Event)
	t := kafka.TopicPartition{Topic: &kafkaTopic, Partition: kafka.PartitionAny}
	if err := kafkaProducer.Produce(&kafka.Message{TopicPartition: t, Value: b}, deliveryChan); err != nil {
		respondWithError(w, errors.Wrap(err, "writting to kafka"))
		return
	}
	dm := <-deliveryChan
	m, ok := dm.(*kafka.Message)
	if !ok {
		respondWithError(w, errors.New("unable to convert to *kafka.Message)"))
		return
	}
	if m.TopicPartition.Error != nil {
		respondWithError(w, errors.Wrap(err, "delivery error"))
		return
	}
	close(deliveryChan)
	log.Println("delivered event to kafka: ", string(b))
}
