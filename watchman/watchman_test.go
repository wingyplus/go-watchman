package watchman

import (
	"encoding/json"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

func TestTrigger(t *testing.T) {
	wd, _ := os.Getwd()
	t.Log(wd)

	client := Connect()
	defer client.Close()

	client.Trigger(wd, Trigger{
		Name:       "gotrigger",
		Expression: []string{"pcre", ".go"},
		Command:    []string{"go", "test"},
	})

	trigger := triggerList(wd)[0]
	if trigger.Name != "gotrigger" {
		t.Error("Expect name is gotrigger")
	}
	if !reflect.DeepEqual(trigger.Expression, []string{"pcre", ".go"}) {
		t.Error("Expect expression is ['pcre', '.go']")
	}
	if !reflect.DeepEqual(trigger.Command, []string{"go", "test"}) {
		t.Error("Expect command is `go test`")
	}

	triggerDel(wd, "gotrigger")
}

type triggerListResult struct {
	Version  string    `json:"version"`
	Triggers []Trigger `json:"triggers"`
}

func triggerList(rootpath string) []Trigger {
	var result triggerListResult
	out, _ := exec.Command("watchman", "trigger-list", rootpath).Output()
	json.Unmarshal(out, &result)
	return result.Triggers
}

func triggerDel(rootpath, name string) {
	exec.Command("watchman", "trigger-del", rootpath, name).Run()
}
