package main

// import (
// 	"fmt"
// 	"io/ioutil"
// )

// type gcloudAPI struct {
// 	project string
// 	gkeZone string
// 	gkeName string
// }

// var gcloud = gcloudAPI{}

// func (g *gcloudAPI) getInfo(data interface{}, resFunc wsResFunc) {
// 	var authKey string

// 	b, err := ioutil.ReadFile(dataPath + "gcpServiceAccountKey.json")
// 	if !isErr(err) {
// 		authKey = string(b)
// 	}

// 	obj := map[string]interface{}{
// 		"authKey": authKey,
// 		"project": gcloud.project,
// 		"gkeZone": gcloud.gkeZone,
// 		"gkeName": gcloud.gkeName,
// 	}
// 	resFunc(200, obj, "")
// }

// func (g *gcloudAPI) setAuthKey(data interface{}, resFunc wsResFunc) {
// 	// fmt.Println(data)
// 	b := []byte(data.(string))
// 	err := ioutil.WriteFile(dataPath+"gcpServiceAccountKey.json", b, 0777)
// 	if isErr(err) {
// 		resFunc(500, nil, err.Error())
// 	}

// 	fmt.Println("\n--- SetServiceAccountKey --- ")
// 	_, err = cmdEX.run("gcloud", "auth", "activate-service-account", "--key-file", dataPath+"gcpServiceAccountKey.json")
// 	if isErr(err) {
// 		resFunc(400, nil, err.Error())
// 		return
// 	}
// 	resFunc(200, "", "")
// 	fmt.Println()
// }

// func (g *gcloudAPI) setProject(data interface{}, resFunc wsResFunc) {
// 	p := data.(string)
// 	_, err := cmdEX.run("gcloud", "config", "set", "project", p)
// 	if isErr(err) {
// 		resFunc(400, nil, err.Error())
// 		return
// 	}
// 	updateConfig("gcloud/project", p)
// 	g.project = p
// 	resFunc(200, "", "")
// }

// func (g *gcloudAPI) setGKE(data interface{}, resFunc wsResFunc) {
// 	d := data.(map[string]interface{})
// 	zone := d["zone"].(string)
// 	name := d["name"].(string)
// 	_, err := cmdEX.run("gcloud", "container", "clusters", "get-credentials", name, "--zone", zone, "--project", g.project)
// 	if isErr(err) {
// 		resFunc(400, nil, err.Error())
// 		return
// 	}
// 	updateConfig("gcloud/gkeZone", zone)
// 	updateConfig("gcloud/gkeName", name)
// 	g.gkeName = name
// 	g.gkeZone = zone
// 	resFunc(200, "connected", "")
// }
