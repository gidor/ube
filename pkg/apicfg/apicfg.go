package apicfg

import (
	"os"
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

func GetApiCfg() {

}
