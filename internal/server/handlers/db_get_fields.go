package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/snapp-incubator/crafting-table/internal/structure"
)

type requestBodyGetFieldsOFStruct struct {
	Source     string `json:"source"`
	StructName string `json:"struct_name"`
}

func GetFieldsOFStruct() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		var rq requestBodyGetFieldsOFStruct
		err := json.NewDecoder(r.Body).Decode(&rq)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		s, err := structure.BindStruct(rq.Source, rq.StructName)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		fieldsNames := make([]string, 0, len(s.Fields))
		for _, field := range s.Fields {
			fieldsNames = append(fieldsNames, field.Name)
		}

		var res responseBody
		res.Status = "success"
		res.Data = make(map[string]interface{})
		res.Data["fields"] = fieldsNames

		err = json.NewEncoder(w).Encode(&res)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}
	}
}
