package temps

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
