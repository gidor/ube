package apicfg

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func (m ApiMethod) Handler(w http.ResponseWriter, r *http.Request) {
	// execute the api method : get info from request, and prepare a
	// RuntimeInfo for the execution then launch the Execute method

	info := new(RuntimeInfo)
	info.current = m.self
	info.context = make(ParamsType)
	info.headers = make(map[string]string)
	// get paramaters in path using gorilla mux
	for k, v := range mux.Vars(r) {
		info.context[k] = v
	}

	for k, v := range r.URL.Query() {
		if len(v) == 1 {
			info.context[k] = v[0]
		} else {
			info.context[k] = v
		}
	}
	if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
		contenttype := r.Header.Get("Content-Type")
		if contenttype == "application/x-www-form-urlencoded" {
			err := r.ParseForm()
			if err == nil {
				for k, v := range r.PostForm {
					if len(v) == 1 {
						info.context[k] = v[0]
					} else {
						info.context[k] = v
					}
				}
			}
		} else if contenttype == "application/json" || contenttype == "application/javascript" {
			// getparameters from a json marshaled object in body
			b, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err == nil {
				// Unmarshal
				var objmap ParamsType
				err = json.Unmarshal(b, &objmap)
				if err == nil {
					for k, v := range objmap {
						info.context[k] = v
					}
				}
			}

		} else if contenttype[0:19] == "multipart/form-data" {
			// TODO manage multipart form-data and file upload
			err := r.ParseMultipartForm(32 << 20)
			if err == nil {
				// r.FormFile()
				for k, v := range r.PostForm {
					if len(v) == 1 {
						info.context[k] = v[0]
					} else {
						info.context[k] = v
					}
				}
			}
		} else {
			err := r.ParseForm()
			if err == nil {
				for k, v := range r.PostForm {
					info.context[k] = v
				}
			}
		}

	}

	// execute method
	m.execute(info)
	// set headers
	for k, v := range info.headers {
		w.Header().Add(k, v)
	}
	if info.todo == jsonEncode {
		w.Header().Add("Content-Type", "application/json")
	}
	if info.status > 0 {
		w.WriteHeader(info.status)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	switch info.todo {
	case jsonEncode:
		b, err := json.Marshal(info.data)
		if err != nil {
			io.WriteString(w, fmt.Sprintf(`{"error": "while marshaling data: %s"}`, err.Error()))
		}
		// io.WriteString(w, fmt.Sprintf(`{"error": "while marshaling data: %s"}`, err.Error()))
		w.Write(b)
	}
}

func (m ApiMethod) execute(info *RuntimeInfo) {
	status := 200
	for i := 0; i < len(m.Actions); i++ {
		status = m.Actions[i].Execute(info)
		if status != 200 {
			break
		}
	}
	info.status = status
}
