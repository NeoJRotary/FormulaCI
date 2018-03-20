package main

// rows, _ := db.Query("SELECT * FROM config;")
// 	// rows.Columns()
// 	defer rows.Close()
// 	for rows.Next() {
// 		var key string
// 		var val string
// 		rows.Scan(&key, &val)
// 		switch key {
// 		case "git/email":
// 			git.email = val
// 			cmdEX.run("git", "config", "--global", "user.email", val)
// 		case "git/webhookToken":
// 			git.webhookToken = val
// 		case "gcloud/project":
// 			gcloud.project = val
// 			cmdEX.run("gcloud", "config", "set", "project", val)
// 		case "gcloud/gkeZone":
// 			gcloud.gkeZone = val
// 		case "gcloud/gkeName":
// 			gcloud.gkeName = val
// 			// case "gitlabEmail":
// 			// 	git.gitlabEmail = val
// 			// case "githubEmail":
// 			// 	git.githubEmail = val
// 		}
// 	}

// 	if gcloud.gkeName != "" {
// 		cmdEX.run("gcloud", "container", "clusters", "get-credentials", gcloud.gkeName, "--zone", gcloud.gkeZone, "--project", gcloud.project)
// 	}
