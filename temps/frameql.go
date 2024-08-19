package temps

import (
	"fmt"
	"os"
	"text/template"
)

func GraphFrame() {
	// ####################################################
	//  graph template
	schema_tmpl, err := template.New("RenderData").Parse(graphSchemaTemplate)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("gschema", os.ModePerm)
	if err != nil {
		panic(err)
	}

	schema_file, err := os.Create("gschema/schema.graphqls")
	if err != nil {
		panic(err)
	}
	defer schema_file.Close()

	err = schema_tmpl.Execute(schema_file, RenderData)
	if err != nil {
		panic(err)
	}
}

var graphSchemaTemplate = `
# Define the input type for pagination
input PaginationInput {
  page: Int!   # Page number
  limit: Int!  # Number of items per page
}

# Define the type for pagination information
type PageInfo {
  hasNextPage: Boolean!
  hasPreviousPage: Boolean!
  startCursor: String
  endCursor: String
}

{{range .Models}}
type {{.Name}}Connection {
  edges: [{{.Name}}!]!
  pageInfo: PageInfo!
}

type {{.Name}} {
	{{range .Fields}} {{.LowerName}} {{.Type}}!
	{{end}}}

input Create{{.Name}}Input {
	{{range .Fields}} {{if .Post}} {{.LowerName}}! {{.Type}}{{end}} {{end}}}

input Update{{.Name}}Input {
	{{range .Fields}} {{if .Put}} {{.LowerName}}! {{.Type}}{{end}}
	{{end}}}

{{end}}


# Define the queries
type Query {
{{range .Models}} # Retrieve a paginated list of apps
  #create paginated items
  {{.LowerName}}s(pagination: PaginationInput!): {{.Name}}Connection!

  # Retrieve a specific app by its ID
  {{.Name}}(id: ID!): {{.Name}}{{end}}}

# Define the mutations
type Mutation {
  {{range .Models}}# Create a new app
  #create object
  create(input: Create{{.Name}}Input!): App!

  # Update an existing app
  update{{.LowerName}}(input: Update{{.Name}}Input!): App!

  # Delete an app by its ID
  delete{{.LowerName}}(id: ID!): Boolean!{{end}}
}
`
