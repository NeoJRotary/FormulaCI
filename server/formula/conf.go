package formula

import (
	"bytes"
	"io/ioutil"

	D "github.com/NeoJRotary/describe-go"
	yaml "gopkg.in/yaml.v2"
)

// Conf ...
type Conf struct {
	Repo     string
	Branch   string
	Mode     string
	Name     string
	Setup    []string
	Trigger  []Trigger
	Flow     []string
	Steps    map[string]Step
	Deploy   Deploy
	Webhooks []Webhook
}

// Trigger ...
type Trigger struct {
	Name    string
	Value   string
	Tag     string
	Changes []string
}

// Step ...
type Step struct {
	Env     map[string]string
	Trigger []Trigger
	Cmd     []string
}

// Deploy ...
type Deploy struct {
	Path       string
	Target     string
	Kubernetes struct {
		Type          string
		Namespace     string
		Name          string
		ContainerName string `yaml:"containerName"`
		ImageHub      string `yaml:"imageHub"`
		Image         string
	}
}

// Webhook ...
type Webhook struct {
	Type   string
	URL    string `yaml:"url"`
	Prefix string
}

// GetConfFromYAML ...
func GetConfFromYAML(path string) ([]*Conf, error) {
	b, err := ioutil.ReadFile(path)
	if D.IsErr(err) {
		return nil, err
	}

	bb := bytes.Split(b, []byte("---\n"))

	list := []*Conf{}
	for _, s := range bb {
		var f Conf
		err := yaml.Unmarshal(s, &f)
		if D.IsErr(err) {
			return nil, err
		}
		list = append(list, &f)
	}
	return list, nil
}
