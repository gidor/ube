package apicfg

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/gidor/ube/pkg/infra"
	"gopkg.in/yaml.v2"
)

type Api struct {
	Name       string         `yaml:"name"`
	Api        []ApiMethod    `yaml:"api,flow"`
	Connection DataConnection `yaml:"dataconnection,flow"`
}

type DataConnection struct {
	Type             string `yaml:"type"`
	ConnectionString string `yaml:"string"`
	pool             *sql.DB
	lasterror        error
}

type ApiMethod struct {
	Name    string       `yaml:"name"`
	Verbs   []string     `yaml:"verbs,flow"`
	Path    string       `yaml:"path"`
	Params  []string     `yaml:"params"`
	Actions []MethodStep `yaml:"actions,flow"`
	self    *Api
}

type MethodStep struct {
	Do     string     `yaml:"do"`
	Params ParamsType `yaml:"params,flow"`
}

type RuntimeInfo struct {
	current *Api
	context ParamsType
	data    interface{}
	headers map[string]string
	status  int
	todo    outWriteOperation
}

type outWriteOperation int8

const (
	jsonEncode outWriteOperation = 0
	stream     outWriteOperation = 1
	copy       outWriteOperation = 2
)

type ParamsType map[string]interface{}

var cfg Api

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
	// api := Api{}
	// var mapi yaml.MapSlice

	content, err := ioutil.ReadFile(cfg)
	if err != nil {
		return (*Api)(nil), err
	}
	err = yaml.Unmarshal(content, &api)
	if err != nil {
		return (*Api)(nil), err
	}
	// initialize the refrence to the api for each method
	for i := 0; i < len(api.Api); i++ {
		api.Api[i].self = &api
		// method := api.Api[i]
		// method.self = &api
	}
	api.Connection.open()
	return &api, nil

}

func (api Api) Method(name string) (*ApiMethod, error) {
	// get method by name
	for i := 0; i < len(api.Api); i++ {
		method := api.Api[i]
		if method.Name == name {
			return &method, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Not found %s", name))
}

// func (api Api) Register() {
// 	// get method by name
// 	for i := 0; i < len(api.Api); i++ {
// 		method := api.Api[i]
// 		if method.Name == name {
// 			return &method, nil
// 		}
// 	}
// 	return nil, errors.New(fmt.Sprintf("Not found %s", name))
// }
