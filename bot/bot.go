package bot

import (
	"gbot/git"
	"strings"
	"regexp"
	"fmt"
	"github.com/google/go-github/github"
)

type approve_t struct {
	Path	string
	Name	string
}

type owner_t struct {
	Path	string
	Name	string
	Type	string
}

var teams 	map[string][]string
var owners	[]owner_t

const commentText = "approve mention"

func init() {
	clearData()
}

func clearData() {
	teams = make(map[string][]string)
	owners = make([]owner_t, 0)
}

func loadTeams(owner string, repo string) error {
	teamsContent, err := git.GetTeams(owner, repo)
	if err != nil {
		return err
	}

	lines := strings.Split(*teamsContent, "\n")
	teamRgx := regexp.MustCompile(`^\[(.+)\]$`)

	var currentTeam string
	for _, line := range lines {
		if teamRgx.MatchString(line) {
			currentTeam = teamRgx.FindStringSubmatch(line)[1]
			teams[currentTeam] = make([]string, 0)
			continue
		}
		if currentTeam != "" && line != "" {
			teams[currentTeam] = append(teams[currentTeam], line)
		}
	}

	return nil
}

func loadOwners(owner string, repo string) error {
	ownersContent, err := git.GetOwners(owner, repo)
	if err != nil {
		return err
	}

	lines := strings.Split(*ownersContent, "\n")
	ownerRgx := regexp.MustCompile(`^.+:.+@.+$`)

	for _, line := range lines {
		if ownerRgx.MatchString(line) {
			parts := strings.Split(line, ":")
			subParts := strings.Split(parts[1], "@")
			//owners[parts[0]] = parts[1]
			owners = append(owners, owner_t{
				Path:	parts[0],
				Name:	subParts[0],
				Type:	subParts[1],
			})
		}
	}

	return nil
}

func createMention(fileName string) []string {
	mentions := make([]string, 0)

	for _, owner := range owners {
		fileRgx := regexp.MustCompile(fmt.Sprintf("^%s.*",  regexp.QuoteMeta(owner.Path)))
		tFile := fmt.Sprintf("/%s", fileName)

		if fileRgx.MatchString(tFile) {
			if owner.Type == "developer" {
				mentions = append(mentions, fmt.Sprintf("@%s", owner.Name))
			} else if owner.Type == "team" {
				if val, ok := teams[owner.Name]; ok {
					for _, teamate := range val {
						mentions = append(mentions, fmt.Sprintf("@%s", teamate))
					}
				}
			}
		}
	}

	return mentions
}

func getCoolAndMentions(comments []*github.PullRequestComment) ([]approve_t, bool) {
	//Get all cool comments
	cool := make([]approve_t, 0)

	//Check if already mentioned, and get cool comments
	mentioned := false
	for _, cmnt := range comments {
		if strings.Index(*cmnt.Body, "cool") != -1 {
			cool = append(cool, approve_t{
				Name:	*cmnt.User.Login,
				Path:	*cmnt.Path,
			})
		}

		if strings.Index(*cmnt.Body, commentText) != -1 {
			mentioned = true
		}
	}

	return cool, mentioned
}

func Tick(owner string, repo string) error {
	clearData()
	err := loadOwners(owner, repo)
	if err != nil {
		return err
	}

	err = loadTeams(owner, repo)
	if err != nil {
		return err
	}

	requests, err := git.GetPullRequests(owner, repo)
	for _, rq := range requests {
		requestApproved := true

		comments, err := git.GetPullComments(owner, repo, *rq.Number)
		if err != nil {
			return err
		}
		cool, mentioned := getCoolAndMentions(comments)

		commits, err := git.GetPullCommits(owner, repo, *rq.Number)
		if err != nil {
			return err
		}

		for _, commit := range commits {

			for _, file := range commit.Files {

				mentions := createMention(*file.Filename)

				//if not already mentioned, do it
				if !mentioned {
					mentionText := fmt.Sprintf("%s %s", strings.Join(mentions, " "), commentText)


					err = git.CreateMention(owner, repo, mentionText, *rq.Number, *commit.SHA, *file.Filename, 1)
					if err != nil {
						return err
					}

					requestApproved = false
					continue
				}

				//check for request approving, cool comments must be the same as mentions
				approves := 0

				for _, approve := range cool {
					for _, m := range mentions {
						if approve.Path == *file.Filename && approve.Name == m[1:] {
							approves++
						}
					}
				}

				if approves != len(mentions) {
					requestApproved = false
				}

			}
		}

		if requestApproved {
			err := git.ApproveRequest(owner, repo, *rq.Number)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

