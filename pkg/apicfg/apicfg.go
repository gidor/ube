package apicfg

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/gidor/ube/pkg/infra"
	"gopkg.in/yaml.v2"
)

type Api struct {
	Name string      `yaml:"name"`
	Api  []ApiMethod `yaml:"api"`
}

type ApiMethod struct {
	Name    string       `yaml:"name"`
	Verbs   []string     `yaml:"nerbs"`
	Path    string       `yaml:"path"`
	Params  []string     `yaml:"params"`
	Actions []MethodStep `yaml:"actions"`
}

type MethodStep struct {
	Do     string     `yaml:"do"`
	Params ParamsType `yaml:"params"`
}

type ParamsType map[string]interface{}

var cfg Api

func init() {
	cfg := os.Getenv("ubecfg")
	if cfg == "" {

	}
}

// 		-
// 		-
// 			do: assert
// 			params:
// 				vars:
// 					- id
// 					- pippo

// 		-
// 			do: sql
// 			params:
// 				stmnt: select * from table where id = {{id}}
// 				out: data

// }

func GetApiCfg() (*Api, error) {
	cfg := infra.Getenv("ube_cfg", "---")
	if cfg == "---" {
		return (*Api)(nil), errors.New("ube_cfg not set")
	}
	var api Api
	content, err := ioutil.ReadAll(cfg)
	if err != nil {
		return (*Api)(nil), err
	}
	err = yaml.Unmarshal(content, api)
	if err != nil {
		return (*Api)(nil), err
	}
	return &api, nil

}
