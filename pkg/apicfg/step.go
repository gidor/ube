package apicfg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cbroglie/mustache"
	"github.com/gidor/ube/pkg/infra"
)

func valueAsString(container ParamsType, key string, defval string, context ParamsType) (string, bool) {
	value, ok := container[key]
	if !ok {
		return defval, false
	}
	switch value.(type) {
	case string:
		r, e := mustache.Render(value.(string), context)
		if e != nil {
			return defval, false
		}
		return r, true
	default:
		return fmt.Sprintf("%v", value), true
	}
}

func valueAsInt(container ParamsType, key string, defval int, context ParamsType) (int, bool) {
	value, ok := container[key]
	if !ok {
		return defval, false
	}
	switch value.(type) {
	case int:
		return value.(int), true
	case string:
		strval, err := mustache.Render(value.(string), context)
		if err != nil {
			return defval, false
		}
		intval, err := strconv.Atoi(strval)
		if err != nil {
			return defval, false
		}
		return intval, true
	default:
		return defval, true
	}
}

func valueAsString2(value interface{}, context ParamsType) (string, error) {
	switch value.(type) {
	case string:
		return mustache.Render(value.(string), context)
	default:
		return fmt.Sprintf("%v", value), nil
	}
}

func (d MethodStep) sql(data *RuntimeInfo) int {
	var stmnt string
	var err error
	ret := 200
	log, err := infra.GetLogger()
	out, _ := valueAsString(d.Params, "out", "data", data.context)
	outdata := out == "data"
	switch d.Params["stmnt"].(type) {
	case string:
		stmnt = d.Params["stmnt"].(string)
		stmnt, err = mustache.Render(stmnt, data.context)
	case []string:
		stmnt = strings.Join(d.Params["stmnt"].([]string), " ")
		stmnt, err = mustache.Render(stmnt, data.context)
	default:
		log.Log.Error("stmnt not defined in params for action sql")
		err = errors.New("stmnt not defined in params for action sql")
	}
	if err != nil {
		if log != nil {
			log.Log.Error(err.Error())
			ret = 500
		}
	} else {
		// fmt.Println(stmnt)
		log.Log.Debug("sql:" + stmnt)
		result, err := data.current.Connection.QueryRecord(stmnt)

		if err != nil {
			if log != nil {
				log.Log.Error(err.Error())
				ret = 500
			}
		} else {
			if outdata {
				data.data = result
				if result == nil {
					ret = 404
				}
			} else {
				if len(result) > 0 {
					m := result[0]
					for k, v := range m {
						data.context[k] = v
					}
				}
			}

		}
	}
	return ret
}

func (d MethodStep) assertdefined(data *RuntimeInfo) int {
	reterror, _ := valueAsInt(d.Params, "onerror", 400, data.context)
	vars, ok := d.Params["vars"]
	if !ok {
		return 200
	} else {
		switch vars.(type) {
		case string:

			if _, ok = data.context[vars.(string)]; ok {
				return 200
			} else {
				return reterror
			}
		case []string:
			for i := 0; i < len(vars.([]string)); i++ {
				key := vars.([]string)[i]
				if _, ok = data.context[key]; !ok {
					return reterror
				}
			}
			return 200
		default:
			return 200
		}
	}
}

func (d MethodStep) assertequal(data *RuntimeInfo) int {
	reterror, _ := valueAsInt(d.Params, "onerror", 400, data.context)
	vars, ok := d.Params["vars"]
	if !ok {
		return 200
	} else {
		switch vars.(type) {
		case map[string]interface{}:
			// m := vars.(map[string]interface{})
			for k, _ := range vars.(map[string]interface{}) {
				val, _ := valueAsString(vars.(map[string]interface{}), k, "", data.context)
				if cv, ok := data.context[k]; ok {
					if val != fmt.Sprintf("%v", cv) {
						return reterror
					}
				} else {
					return reterror
				}
			}
			return 200
		default:
			return 200
		}
	}
}

func (d MethodStep) assertnotequal(data *RuntimeInfo) int {
	reterror, _ := valueAsInt(d.Params, "onerror", 400, data.context)
	vars, ok := d.Params["vars"]
	if !ok {
		return 200
	} else {
		switch vars.(type) {
		case map[string]interface{}:
			for k, v := range vars.(map[string]interface{}) {
				if cv, ok := data.context[k]; ok {
					if fmt.Sprintf("%v", v) == fmt.Sprintf("%v", cv) {
						return reterror
					}
				} else {
					return reterror
				}
			}
			return 200
		default:
			return 200
		}
	}
}

func (d MethodStep) defaultvalue(data *RuntimeInfo) int {
	vars, ok := d.Params["vars"]
	if !ok {
		return 200
	} else {
		switch vars.(type) {
		case map[string]interface{}:
			for k, v := range vars.(map[string]interface{}) {
				if _, ok := data.context[k]; !ok {
					switch v.(type) {
					case string:
						data.context[k], _ = mustache.Render(v.(string), data.context)
					default:
						data.context[k] = fmt.Sprintf("%v", v)
					}
				}
			}
		}
	}
	return 200
}
func (d MethodStep) setheader(data *RuntimeInfo) int {
	for k, v := range d.Params {
		switch v.(type) {
		case string:
			data.headers[k], _ = mustache.Render(v.(string), data.context)
		default:
			data.headers[k] = fmt.Sprintf("%v", v)
		}
	}
	return 200
}

func (d MethodStep) unknown(data *RuntimeInfo) int {
	return 200
}

/*
 execute a method step the stp actioon are
  - sql:
  		execute an sql and put the rsult in dat or context
  - settheader:
  		set heders
  - assertdefined:
  		assert vairiable are defined in context
  - assertequal:
  		assert vairiable are defined in context with value
  - assertnotequal:
  		assert vairiable are defined in context  with value non equal to
  - default:
  		set value for variable in context oly if the variable is not defined
*/
func (step *MethodStep) Execute(data *RuntimeInfo) int {
	ret := 200
	switch step.Do {
	case "sql":
		ret = step.sql(data)
	case "setheader":
		ret = step.setheader(data)
	case "assertdefined":
		ret = step.assertdefined(data)
	case "assertequal":
		ret = step.assertequal(data)
	case "assertnotequal":
		ret = step.assertnotequal(data)
	case "default":
		ret = step.defaultvalue(data)
	default:
		ret = step.unknown(data)
	}
	return ret
}
