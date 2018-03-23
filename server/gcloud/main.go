package gcloud

import (
	"io/ioutil"

	"../executer"
	"../litedb"
	D "github.com/NeoJRotary/describe-go"
)

// Agent ...
type Agent struct {
	db                *litedb.DB
	execEX            *executer.Exec
	DataPath          string
	ServiceAccountKey string
	Project           string
	GKEZone           string
	GKEName           string
}

// NewAgent ...
func NewAgent(dataPath string, db *litedb.DB, execEX *executer.Exec) (agn *Agent, err error) {
	D.RecoverErr(nil)

	config, err := db.GetAllConfig()
	D.CheckErr(err)

	agn.db = db
	agn.execEX = execEX
	agn.Project = db.GetConfig("gcloud/Project")
	agn.GKEZone = db.GetConfig("gcloud/GKEZone")
	agn.GKEName = db.GetConfig("gcloud/GKEName")

	return agn, nil
}

// UpdateServiceAccountKey ...
func (agn *Agent) UpdateServiceAccountKey(b []byte) (err error) {
	D.RecoverErr(nil)

	path := agn.DataPath + "GCPServiceAccountKey.json"
	err = ioutil.WriteFile(path, b, 0777)
	D.CheckErr(err)

	_, err = agn.execEX.Run("gcloud", "auth", "activate-service-account", "--key-file", path)
	D.CheckErr(err)

	return nil
}

// GetServiceAccountKey ...
func (agn *Agent) GetServiceAccountKey() ([]byte, error) {
	return ioutil.ReadFile(agn.DataPath + "GCPServiceAccountKey.json")
}

// SetProject ...
func (agn *Agent) SetProject(name string) (err error) {
	D.RecoverErr(nil)

	_, err = agn.execEX.Run("gcloud", "config", "set", "project", name)
	D.CheckErr(err)

	err = agn.db.UpdateConfig("gcloud/Project", name)
	D.CheckErr(err)

	return nil
}

func (agn *Agent) SetGKE(name, zone string) error {

}

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
