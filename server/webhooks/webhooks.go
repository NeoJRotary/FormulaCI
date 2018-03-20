package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"./webhooks"
)

// type webhookConf struct {
// 	From []string
// 	To   map[string]string
// }

// Manager webhooks manager
type webhooksManager struct{}

var webhooksM = webhooksManager{}

func (m *webhooksManager) init(mux *http.ServeMux) {
	mux.HandleFunc("/webhook/gitlab", m.getFromGitlab)
}

func (m *webhooksManager) getFromGitlab(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// if r.Method == http.MethodOptions {
	// 	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	w.Header().Set("Access-Control-Allow-Methods", "POST")
	// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Gitlab-Token")
	// 	w.WriteHeader(http.StatusOK)
	// 	return
	// } else
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Header.Get("X-Gitlab-Token") != git.webhookToken {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var hook webhooks.Gitlab
	err := decoder.Decode(&hook)
	if isErr(err) {
		return
	}

	if hook.ObjectKind != "push" && hook.ObjectKind != "merge_request" {
		return
	}

	// fmt.Println(repo.hub)

	hookBranch := strings.Replace(hook.Ref, "refs/heads/", "", 1)
	for _, r := range repo.hub["gitlab.com"] {
		if r.Name == hook.Repository.Name && r.Branch == hookBranch {
			go ci.trigger(r.Name, hookBranch)
			return
		}
	}
}

func (m *webhooksManager) send(hook formulaWebhook, hs tableHistory, triggers []*formulaTrigger) {
	switch hook.Type {
	case "slack":
		m.sendToSlack(hook, hs, triggers)
	}
}

func (m *webhooksManager) sendToSlack(hook formulaWebhook, hs tableHistory, triggers []*formulaTrigger) {
	prefix := "*[ " + hs.Repo + "/" + hs.Branch + " ]*"
	if hook.Prefix != "" {
		prefix = hook.Prefix
		prefix = strings.Replace(prefix, "${repo}", hs.Repo, -1)
		prefix = strings.Replace(prefix, "${branch}", hs.Branch, -1)
		prefix = ci.replaceTriggerVar(triggers, prefix)
	}

	temp := ""
	switch hs.Result {
	case -2:
		temp = prefix + " pipeline canceled."
	case -1:
		temp = prefix + " new pipeline created."
	case 0:
		temp = prefix + " pipeline failed"
	case 1:
		temp = fmt.Sprintf(prefix+" pipeline completed. Used %v seconds", hs.Dur)
	default:
		return
	}
	reader := strings.NewReader(jsonString(map[string]string{"text": temp}))
	_, err := http.Post(hook.URL, "application/json", reader)
	if err != nil {
		fmt.Println(err)
	}
}
