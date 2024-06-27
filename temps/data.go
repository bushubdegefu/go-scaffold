package temps

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Define a struct to match the structure of your JSON data
type Data struct {
	ProjectName string  `json:"project_name"`
	AppName     string  `json:"app_name"`
	BackTick    string  `json:"back_tick"`
	Models      []Model `json:"models"`
}

type Model struct {
	Name        string         `json:"name"`
	LowerName   string         `json:"lower_name"`
	RlnModel    []string       `json:"rln_model"` // value to one of the models defined in the config json file
	BackTick    string         `json:"back_tick"`
	Fields      []Field        `json:"fields"`
	ProjectName string         `json:"project_name"`
	AppName     string         `json:"app_name"`
	Relations   []Relationship `json:"relations"`
}

type Relationship struct {
	ParentName      string `json:"parent_name"`
	LowerParentName string `json:"lower_parent_name"`
	FieldName       string `json:"field_name"`
	LowerFieldName  string `json:"lower_field_name"`
	MtM             bool   `json:"mtm"`
	OtM             bool   `json:"otm"`
	MtO             bool   `json:"mto"`
}

type Field struct {
	BackTick    string `json:"back_tick"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Annotation  string `json:"annotation"`
	CurdFlag    string `json:"curd_flag"`
	Get         bool   `json:"get"`
	Post        bool   `json:"post"`
	Patch       bool   `json:"patch"`
	Put         bool   `json:"put"`
	OtM         bool   `json:"otm"`
	MtM         bool   `json:"mtm"`
	ProjectName string `json:"project_name"`
	AppName     string `json:"app_name"`
}

var RenderData Data

func LoadData() {
	// Open the JSON file
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return
	}
	defer file.Close() // Defer closing the file until the function returns

	// Decode the JSON content into the data structure
	var data Data
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&RenderData)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	// setting default value for config data file
	//  GetPostPatchPut
	// "Get$Post$Patch$Put$OtM$MtM"

	for i := 0; i < len(data.Models); i++ {
		data.Models[i].LowerName = strings.ToLower(data.Models[i].Name)
		data.Models[i].AppName = data.AppName
		data.Models[i].ProjectName = data.ProjectName
		rl_list := make([]Relationship, 0)
		for k := 0; k < len(data.Models[i].RlnModel); k++ {
			rmf := strings.Split(data.Models[i].RlnModel[k], "$")
			cur_relation := Relationship{
				ParentName:      data.Models[i].Name,
				LowerParentName: data.Models[i].LowerName,
				FieldName:       rmf[0],
				LowerFieldName:  strings.ToLower(rmf[0]),
				MtM:             rmf[1] == "mtm",
				OtM:             rmf[1] == "otm",
				MtO:             rmf[1] == "mto",
			}
			rl_list = append(rl_list, cur_relation)
			data.Models[i].Relations = rl_list
		}

		for j := 0; j < len(data.Models[i].Fields); j++ {
			data.Models[i].Fields[j].BackTick = "`"
			cf := strings.Split(data.Models[i].Fields[j].CurdFlag, "$")

			data.Models[i].Fields[j].Get, _ = strconv.ParseBool(cf[0])
			data.Models[i].Fields[j].Post, _ = strconv.ParseBool(cf[1])
			data.Models[i].Fields[j].Patch, _ = strconv.ParseBool(cf[2])
			data.Models[i].Fields[j].Put, _ = strconv.ParseBool(cf[3])
			data.Models[i].Fields[j].AppName = data.AppName
			data.Models[i].Fields[j].ProjectName = data.ProjectName

		}
	}
}
