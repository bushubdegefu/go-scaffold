package temps

// Define a struct to match the structure of your JSON data
type Data struct {
	ProjectName string  `json:"project_name"`
	AppName     string  `json:"app_name"`
	BackTick    string  `json:"back_tick"`
	Models      []Model `json:"models"`
}

type Model struct {
	Name        string  `json:"name"`
	LowerName   string  `json:"lower_name"`
	BackTick    string  `json:"back_tick"`
	Fields      []Field `json:"fields"`
	ProjectName string  `json:"project_name"`
	AppName     string  `json:"app_name"`
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
	ProjectName string `json:"project_name"`
	AppName     string `json:"app_name"`
}
