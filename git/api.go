package git

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"fmt"
	"errors"
	"encoding/base64"
)

var client *github.Client

func CreateClient(token string) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client = github.NewClient(tc)
}

func GetRepos() ([]*github.Repository, error) {
	repos, _, err := client.Repositories.List("", nil)
	if err != nil {
		return nil, err
	}

	return repos, nil
}

func GetPullComments(owner string, repo string, number int) ([]*github.PullRequestComment, error) {
	comments, _, err := client.PullRequests.ListComments(owner, repo, number, nil)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func GetPullRequests(owner string, repo string) ([]*github.PullRequest, error) {
	requests, _, err := client.PullRequests.List(owner, repo, nil)
	if err != nil {
		return nil, err
	}

	return requests, nil
}

func CheckAccess(repository string) (repoOwner *string, repoName *string, err error) {
	repos, err := GetRepos()
	if err != nil {
		return nil, nil, err
	}

	for _, r := range repos {
		if repository == fmt.Sprintf("%s/%s", *r.Owner.Login, *r.Name) {
			repoOwner = r.Owner.Login
			repoName = r.Name
			return
		}
	}

	return nil, nil, errors.New("You have no access to given repository")
}

func getContent(owner string, repo string, path string) (*string, error) {
	file, _, _, err := client.Repositories.GetContents(owner, repo, path, nil)
	if err != nil {
		return nil, err
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(*file.Content)
	if err != nil {
		return nil, err
	}
	data := string(decodedBytes)

	return &data, nil
}

func GetOwners(owner string, repo string) (*string, error) {
	data, err := getContent(owner, repo, ".owners")
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetTeams(owner string, repo string) (*string, error) {
	data, err := getContent(owner, repo, ".teams")
	if err != nil {
		return nil, err
	}

	return data, nil
}

func CreateMention(owner string, repo string, text string, rq int, commitID string, path string, position int) error {
	_, _, err := client.PullRequests.CreateComment(owner, repo, rq, &github.PullRequestComment{
		Body: &text,
		CommitID: &commitID,
		Path: &path,
		Position: &position})
	if err != nil {
		return err
	}

	return nil
}

func GetPullCommits(owner string, repo string, num int) ([]*github.RepositoryCommit, error) {
	commits, _, err := client.PullRequests.ListCommits(owner, repo, num, nil)
	if err != nil {
		return nil, err
	}

	filledCommits := make([]*github.RepositoryCommit, 0)
	for _, cmt := range commits {
		rc, _, err := client.Repositories.GetCommit(owner, repo, *cmt.SHA)
		if err != nil {
			return nil, err
		}
		filledCommits = append(filledCommits, rc)
	}

	return filledCommits, nil
}

func ApproveRequest(owner string, repo string, num int) error {
	_, _, err := client.PullRequests.Merge(owner, repo, num, "bot merge", nil)
	if err != nil {
		return err
	}

	return nil
}

