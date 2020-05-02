package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// Simply check if an error occurred
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// The c++ ternary operator
func ternary(cond bool, t, f interface{}) interface{} {
	if cond {
		return t
	} else {
		return f
	}
}

// A json convertible data map
type object map[string]interface{}

func main() {
	fmt.Println("Sending discord message...")

	// Needed for the random profile image
	rand.Seed(time.Now().Unix())

	// Getting the flags
	success := flag.Bool("success", false, "If set the message will be green")
	hookUrl := flag.String("hook", os.Getenv("WEBHOOK_URL"), "The discord webhook url")
	message := flag.String("message", "¯\\_(ツ)_/¯", "The text of the message")
	flag.Parse()

	// The request payload based on https://discordapp.com/developers/docs/resources/webhook
	payload := object{
		"avatar_url": fmt.Sprintf("https://travis-ci.org/images/logos/TravisCI-Mascot-%d.png", rand.Intn(3)+1),
		"embeds": []object{
			{
				"title": *message,
				"color": ternary(*success, 0x81f207, 0xf2073a),
				"author": object{
					"name":     fmt.Sprintf("Job %s on branch %s at %s", os.Getenv("TRAVIS_JOB_NUMBER"), os.Getenv("TRAVIS_BRANCH"), os.Getenv("TRAVIS_REPO_SLUG")),
					"url":      os.Getenv("TRAVIS_BUILD_WEB_URL"),
					"icon_url": fmt.Sprintf("https://travis-ci.org/images/logos/TravisCI-Mascot-%s.png", ternary(*success, "blue", "red")),
				},
				"fields": []object{
					{
						"name":   "Commit",
						"value":  os.Getenv("TRAVIS_COMMIT"),
						"inline": true,
					},
					{
						"name":   "Commit message",
						"value":  os.Getenv("TRAVIS_COMMIT_MESSAGE"),
						"inline": true,
					},
					{
						"name":   "Branch",
						"value":  os.Getenv("TRAVIS_BRANCH"),
						"inline": true,
					},
					{
						"name":   "Stage",
						"value":  os.Getenv("TRAVIS_BUILD_STAGE_NAME"),
						"inline": true,
					},
				},
			},
		},
	}
	pl, err := json.Marshal(payload)
	check(err)
	resp, err := http.Post(*hookUrl, "application/json", bytes.NewBuffer(pl))
	check(err)
	if resp.StatusCode != 204 {
		panic("Sending unsuccessful")
	}
	fmt.Println("Message sent successfully")
}
