package bot

import (
	"testing"
	"gbot/git"
)

func init() {
	git.CreateClient("59781114d585e4b6dbd3223b5ca393f7fc31e4b9")
}

func TestLoadTeams(t *testing.T) {
	err := loadTeams("devchallenge10t", "test1")
	if err != nil {
		t.Error(err)
	}
}

func TestLoadOwners(t *testing.T) {
	err := loadOwners("devchallenge10t", "test1")
	if err != nil {
		t.Error(err)
	}
}
