package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var acceptOrigin = []string{
	"localhost:8099",
	"127.0.0.1:8099",
	"192.168.99.100:8099",
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == r.Host {
			return true
		}

		for _, o := range acceptOrigin {
			if o == r.Host {
				return true
			}
		}
		return false
	},
}

// APIServer ...
type APIServer struct {
	conn *websocket.Conn
	mux  *http.ServeMux
}

type wsEvt struct {
	Name string      `json:"evt"`
	I    int         `json:"i"`
	Data interface{} `json:"data"`
}

type wsResFunc func(int, interface{}, string) error

// Start ...
func Start(mux *http.ServeMux) {
	svr := formulaServer{mux: mux}
	mux.HandleFunc("/sys-ws", svr.wsUpgrade)
	webhooksM.init(mux)
}

func (apis *APIServer) wsUpgrade(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if isErr(err) {
		log.Println(err)
		return
	}
	defer c.Close()
	svr.conn = c
	for {
		mt, msg, err := c.ReadMessage()
		if isErr(err) {
			log.Println(err)
			break
		}
		if mt != websocket.TextMessage {
			continue
		}
		evt := wsEvt{}
		err = json.Unmarshal(msg, &evt)
		if isErr(err) {
			log.Println(err)
			continue
		}
		go svr.handler(evt)
	}
}

func (apis *APIServer) getResFunc(evt wsEvt) wsResFunc {
	return func(status int, data interface{}, errMsg string) error {
		if data == nil {
			data = ""
		}
		obj := map[string]interface{}{
			"evt":    evt.Name,
			"i":      evt.I,
			"status": status,
			"data":   data,
			"error":  errMsg,
		}
		b, err := json.Marshal(obj)
		if isErr(err) {
			return err
		}
		err = svr.conn.WriteMessage(websocket.TextMessage, b)
		if isErr(err) {
			return err
		}
		return nil
	}
}

func (apis *APIServer) handler(evt wsEvt) {
	resFunc := svr.getResFunc(evt)
	switch evt.Name {
	case "gcloud/getInfo":
		gcloud.getInfo(evt.Data, resFunc)
	case "gcloud/setAuthKey":
		gcloud.setAuthKey(evt.Data, resFunc)
	case "gcloud/setProject":
		gcloud.setProject(evt.Data, resFunc)
	case "gcloud/setGKE":
		gcloud.setGKE(evt.Data, resFunc)
	case "git/getInfo":
		git.getInfo(evt.Data, resFunc)
	case "git/setEmail":
		git.setEmail(evt.Data, resFunc)
	case "git/setWebhookToken":
		git.setWebhookToken(evt.Data, resFunc)
	case "git/generateSSH":
		git.generateSSH(evt.Data, resFunc)
	case "repo/getList":
		repo.getList(evt.Data, resFunc)
	case "repo/add":
		repo.add(evt.Data, resFunc)
	case "repo/remove":
		repo.remove(evt.Data, resFunc)
	case "repo/trigger":
		repo.trigger(evt.Data, resFunc)
	case "formula/getHistory":
		ci.getHistory(evt.Data, resFunc)
	case "ping":
		resFunc(200, "pong", "")
	default:
		resFunc(404, nil, "Not Found")
	}
}
