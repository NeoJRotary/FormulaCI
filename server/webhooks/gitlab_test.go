package webhooks

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestGitlabWebhook(t *testing.T) {
	temp := `{
		"object_kind": "push",
		"ref": "refs/heads/master",
		"repository":{
			"name": "admin",
			"url": "git@example.com:mike/diaspora.git"
		},
		"commits": [
			{
				"id": "b6568db1bc1dcd7f8b4d5a946b0b91f9dacd7327",
				"message": "Update Catalan translation to e38cb41.",
				"timestamp": "2011-12-12T14:27:31+02:00",
				"url": "http://example.com/mike/diaspora/commit/b6568db1bc1dcd7f8b4d5a946b0b91f9dacd7327",
				"author": {
					"name": "Jordi Mallach",
					"email": "jordi@softcatala.org"
				},
				"added": ["CHANGELOG"],
				"modified": [".formulaci.yaaml"],
				"removed": []
			},
			{
				"id": "da1560886d4f094c3e6c9ef40349f7d38b5d27d7",
				"message": "fixed readme",
				"timestamp": "2012-01-03T23:36:29+02:00",
				"url": "http://example.com/mike/diaspora/commit/da1560886d4f094c3e6c9ef40349f7d38b5d27d7",
				"author": {
					"name": "GitLab dev user",
					"email": "gitlabdev@dv6700.(none)"
				},
				"added": ["CHANGELOG"],
				"modified": ["app/controller/application.rb"],
				"removed": []
			}
		],
		"total_commits_count": 4
	}`

	req, err := http.NewRequest("POST", "http://localhost:8099/webhook/gitlab", strings.NewReader(temp))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("X-Gitlab-Token", "12345")

	client := &http.Client{}
	client.Do(req)

	// res, err := http.Post(
	// 	"http://localhost:8099/webhook/gitlab",
	// 	"application/json",
	// 	strings.NewReader(temp))

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// defer res.Body.Close()
}
