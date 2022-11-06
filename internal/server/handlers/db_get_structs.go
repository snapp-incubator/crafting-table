package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/snapp-incubator/crafting-table/internal/structure"
)

type requestBody struct {
	Source string `json:"source"`
}

type responseBody struct {
	Status string                 `json:"status"`
	Data   map[string]interface{} `json:"data"`
}

func GetStructs() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		var rq requestBody
		err := json.NewDecoder(r.Body).Decode(&rq)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		structs, err := structure.BindAllStructs(rq.Source)
		if err != nil {
			fmt.Println("2")
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		structNames := make([]string, 0, len(structs))
		for _, s := range structs {
			structNames = append(structNames, s.Name)
		}

		var res responseBody
		res.Status = "success"
		res.Data = make(map[string]interface{})
		res.Data["structs"] = structNames

		err = json.NewEncoder(w).Encode(&res)
		if err != nil {
			fmt.Println("4")
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}
	}
}
