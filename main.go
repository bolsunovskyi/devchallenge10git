//59781114d585e4b6dbd3223b5ca393f7fc31e4b9
//devchallenge10t/test1

package main

import (
	"flag"
	"log"
	"gbot/git"
	"gbot/bot"
	"time"
)

var repository, token, intervalStr string
var interval time.Duration

func init() {
	flag.StringVar(&repository, "r", "", "Repository on github, format: author/repo, for example golang/go or angular/angular")
	flag.StringVar(&token, "t", "", "Access token from github, you can create it here: https://github.com/settings/tokens/new")
	flag.StringVar(&intervalStr, "i", "1m", "Tick interval, for example 1s, 5m, 1h, details here: https://golang.org/pkg/time/#ParseDuration")
	flag.Parse()

	git.CreateClient(token)

	var err error
	interval, err = time.ParseDuration(intervalStr)
	if err != nil {
		log.Fatalln(err)
	}

	if token == "" {
		log.Fatalln("No token supplied, use --help")
	}

	if repository == "" {
		log.Fatalln("No repository supplied, use --help")
	}
}

func runTick(owner *string, repo *string) {
	err := bot.Tick(*owner, *repo)
	if err != nil {
		log.Println(err)
	}
}

func main() {


	owner, repo, err := git.CheckAccess(repository)
	if err != nil {
		log.Fatalln(err)
	}

	//important for big intervals to run tick before main cycle
	runTick(owner, repo)

	ticker := time.NewTicker(interval)
	for range ticker.C {
		runTick(owner, repo)
	}
}