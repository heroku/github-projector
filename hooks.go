package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"github.com/heroku/x/hredis/redigo"
	"github.com/pkg/errors"

	"github.com/garyburd/redigo/redis"
	"github.com/heroku/x/hredis"
)

var rp *redis.Pool

func init() {
	rurl, err := hredis.RedissURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal("Error getting redis url:", err)
	}
	rp, err = redigo.NewRedisPoolFromURL(rurl)
	if err != nil {
		log.Fatalf("error making redis pool from url (%s): %s", rurl, err.Error())
	}
}

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
	if err := d.Decode(i); err != nil {
		respondWithError(w, errors.Wrap(err, "decoding error"))
		return
	}

	b, err := json.Marshal(i)
	if err != nil {
		respondWithError(w, errors.Wrap(err, "encoding error"))
		return
	}
	log.Printf("%s", b)
}
