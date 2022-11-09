package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/snapp-incubator/crafting-table/internal/repository"

	"github.com/snapp-incubator/crafting-table/internal/structure"
)

type requestBodyGenerateManifestAndRepo struct {
	Source              string                      `json:"source"`
	Destination         string                      `json:"destination"`
	DestinationFileName string                      `json:"destination_file_name"`
	PackageName         string                      `json:"package_name"`
	StructName          string                      `json:"struct_name"`
	TableName           string                      `json:"table_name"`
	DBLibrary           string                      `json:"db_library"`
	Get                 []structure.GetVariable     `json:"get"`
	Update              []structure.UpdateVariables `json:"update"`
	Create              structure.CreateVariables   `json:"create"`
	Test                bool                        `json:"create_test"`
	//Join        []structure.JoinVariables   `json:"join"`
}

func GenerateManifestAndRepo() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		var rq requestBodyGenerateManifestAndRepo
		err := json.NewDecoder(r.Body).Decode(&rq)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		rq.Destination = rq.Destination + "/" + rq.DestinationFileName + ".go"

		b, err := json.Marshal(&rq)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		var manifest repository.Repo
		err = json.Unmarshal(b, &manifest)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		if err := repository.Generate(manifest.Source, manifest.Destination, manifest.PackageName,
			manifest.StructName, &manifest.Get, &manifest.Update, &manifest.Join,
			manifest.Create.Enable, manifest.Test); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		var res responseBody
		res.Status = "success"
		res.Data = make(map[string]interface{})
		res.Data["message"] = "db repository created successfully"

		err = json.NewEncoder(w).Encode(&res)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}
	}
}
