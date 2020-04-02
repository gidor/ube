package apicfg

import (
	"errors"
	"strings"

	"github.com/cbroglie/mustache"
	"github.com/gidor/ube/pkg/infra"
)

func (d MethodStep) sql(data *RuntimeInfo) int {
	var stmnt string
	var err error
	ret := 200
	log, err := infra.GetLogger()

	switch d.Params["stmnt"].(type) {
	case string:
		stmnt = d.Params["stmnt"].(string)
		stmnt, err = mustache.Render(stmnt, data.context)
	case []string:
		stmnt = strings.Join(d.Params["stmnt"].([]string), " ")
		stmnt, err = mustache.Render(stmnt, data.context)
	default:
		err = errors.New("stmnt not defined in params for action sql")
	}
	if err != nil {
		if log != nil {
			log.Log.Error(err.Error())
			ret = 500
		}
	} else {
		result, err := data.current.Connection.QueryRecord(stmnt)

		if err != nil {
			if log != nil {
				log.Log.Error(err.Error())
				ret = 500
			}
		} else {
			data.data = result
		}
	}
	return ret
}

func (d MethodStep) assert(data *RuntimeInfo) int {
	return 200
}

func (d MethodStep) unknown(data *RuntimeInfo) int {
	return 200
}

func (step *MethodStep) Execute(data *RuntimeInfo) int {
	ret := 200
	switch step.Do {
	case "sql":
		ret = step.sql(data)
	case "assert":
		ret = step.assert(data)
	default:
		ret = step.unknown(data)
	}
	return ret
}
