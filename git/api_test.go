package git

import "testing"

func init() {
	CreateClient("59781114d585e4b6dbd3223b5ca393f7fc31e4b9")
}

func TestGetRepos(t *testing.T) {
	_, err := GetRepos()
	if err != nil {
		t.Error(err)
	}
}

func TestGetPullRequests(t *testing.T) {
	_, err := GetPullRequests("devchallenge10t", "test1")
	if err != nil {
		t.Error(err)
	}
}

func TestCheckAccess(t *testing.T) {
	_, _, err := CheckAccess("devchallenge10t/test1")
	if err != nil {
		t.Error(err)
	}
}

func TestGetTeams(t *testing.T) {
	_, err := GetTeams("devchallenge10t", "test1")
	if err != nil {
		t.Error(err)
	}
}

func TestGetOwners(t *testing.T) {
	_, err := GetOwners("devchallenge10t", "test1")
	if err != nil {
		t.Error(err)
	}
}
