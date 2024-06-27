package temps

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

func TestFrameFiber() {

	test_tmpl, err := template.New("RenderData").Parse(testTemplateFiber)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	// #################################################
	err = os.MkdirAll("tests", os.ModePerm)
	if err != nil {
		panic(err)
	}

	for _, model := range RenderData.Models {

		folder_path := fmt.Sprintf("tests/%v_controller_test.go", model.Name)
		folder_path = strings.ToLower(folder_path)
		test_file, err := os.Create(folder_path)
		if err != nil {
			panic(err)
		}

		err = test_tmpl.Execute(test_file, model)
		if err != nil {
			panic(err)
		}
		test_file.Close()

	}
	// ###################################################################
	test_app_tmpl, err := template.New("RenderData").Parse(testAppTemplate)
	if err != nil {
		panic(err)
	}
	//
	test_app_file, err := os.Create("tests/test_app.go")
	if err != nil {
		panic(err)
	}

	err = test_app_tmpl.Execute(test_app_file, RenderData)
	if err != nil {
		panic(err)
	}
	defer test_app_file.Close()

	// ###################################################################
	// running go mod tidy finally
	if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
		fmt.Printf("error: %v \n", err)
	}
}

var testTemplateFiber = `
package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"{{.ProjectName}}.com/models"
	"{{.ProjectName}}.com/models/controllers"
)

// go test -coverprofile=coverage.out ./...
// go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html


// ##########################################################################
var tests{{.Name}}sPostID = []struct {
	name         string          //name of string
	description  string          // description of the test case
	route        string          // route path to test
	{{.LowerName}}_id      string          //path param
	post_data    models.{{.Name}}Post // patch_data
	expectedCode int             // expected HTTP status code
}{
	// First test case
	{
		name:        "post {{.Name}} - 1",
		description: "post Single {{.Name}}",
		route:       "/"+group_path+"/{{.LowerName}}",
		post_data: models.{{.Name}}Post{
			Name:        "New one Posted 3",
			Description: "Description of Name Posted neww333",
		},
		expectedCode: 200,
	},
	// Second test case
	{
		name:        "post {{.Name}} - 2",
		description: "post Single ",
		route:       "/"+group_path+"/{{.LowerName}}",
		post_data: models.{{.Name}}Post{
			Name:        "New one Posted 3",
			Description: "Description of Name Posted neww333",
		},
		expectedCode: 200,
	},
	// Second Third case
	{
		name:        "get {{.Name}} By ID check - 3",
		description: "get HTTP status 404, when {{.Name}} Does not exist",
		route:       "/"+group_path+"/{{.LowerName}}",
		post_data: models.{{.Name}}Post{
			Name:        "Name one",
			Description: "Description of Name one",
		},
		expectedCode: 500,
	},
}

func TestPost{{.Name}}sByID(t *testing.T) {

	ReturnTestApp()


	// Iterate through test single test cases
	for _, test := range tests{{.Name}}sPostID {
		t.Run(test.name, func(t *testing.T) {
			//  changing post data to json
			post_data, _ := json.Marshal(test.post_data)

			req := httptest.NewRequest(http.MethodPost, test.route, bytes.NewReader(post_data))

			// Add specfic headers if needed as below
			req.Header.Set("Content-Type", "application/json")
			// req.Header.Set("X-APP-TOKEN", "hi")

	

			resp, _ := TestApp.Test(req)


			var responseMap map[string]interface{}
			body, _ := io.ReadAll(resp.Body)
			uerr := json.Unmarshal(body, &responseMap)
			if uerr != nil {
				// fmt.Printf("Error marshaling response : %v", uerr)
				fmt.Println()
			}



			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
			//  running delete test if post is success
			if resp.StatusCode == 200 {
				t.Run("Checking the Delete Request Path for {{.Name}}s", func(t *testing.T) {

					test_route := fmt.Sprintf("%v/%v", test.route, responseMap["data"].(map[string]interface{})["id"])

					req_delete := httptest.NewRequest(http.MethodDelete, test_route, bytes.NewReader(post_data))

					// Add specfic headers if needed as below
					req_delete.Header.Set("Content-Type", "application/json")

					resp, _ := TestApp.Test(req_delete)

					assert.Equalf(t, 200, resp.StatusCode, test.description+"deleteing")
				})
			} else {
				t.Run("Checking the Delete Request Path for  that does not exit", func(t *testing.T) {

					test_route_1 := fmt.Sprintf("%v/:%v", test.route, 1000000)

					req_delete := httptest.NewRequest(http.MethodDelete, test_route_1, bytes.NewReader(post_data))

					// Add specfic headers if needed as below
					req_delete.Header.Set("Content-Type", "application/json")

					resp, _ := TestApp.Test(req_delete)
					assert.Equalf(t, 500, resp.StatusCode, test.description+"deleteing")
				})

				t.Run("Checking the Delete Request Path that is not valid", func(t *testing.T) {

					test_route_2 := fmt.Sprintf("%v/%v", test.route, "$$$")

					req_delete := httptest.NewRequest(http.MethodDelete, test_route_2, bytes.NewReader(post_data))

					// Add specfic headers if needed as below
					req_delete.Header.Set("Content-Type", "application/json")
					resp, _ := TestApp.Test(req_delete)

					
					assert.Equalf(t, 500, resp.StatusCode, test.description+"deleteing")
				})
			}
		})
	}

}

// ##########################################################################
var tests{{.Name}}sPatchID = []struct {
	name         string           //name of string
	description  string           // description of the test case
	route        string           // route path to test
	patch_data   models.{{.Name}}Patch // patch_data
	expectedCode int              // expected HTTP status code
}{
	// First test case
	{
		name:        "patch {{.Name}}s By ID check - 1",
		description: "patch Single {{.Name}} by ID",
		route:       "/"+group_path+"/{{.LowerName}}/1",
		patch_data: models.{{.Name}}Patch{
			Name:        "Name one eight",
			Description: "Description of Name one for test one",
		},
		expectedCode: 200,
	},

	// Second test case
	{
		name:        "get {{.Name}} By ID check - 2",
		description: "get HTTP status 404, when {{.Name}} Does not exist",
		route:       "/"+group_path+"/{{.LowerName}}/1000",
		patch_data: models.{{.Name}}Patch{
			Name:        "Name one eight",
			Description: "Description of Name one for test 3",
		},
		expectedCode: 404,
	},
	// Second test case
	{
		name:        "get {{.Name}} By ID check - 4",
		description: "get HTTP status 404, when {{.Name}} Does not exist",
		route:       "/"+group_path+"/{{.LowerName}}/@@",
		patch_data: models.{{.Name}}Patch{
			Name:        "Name one eight",
			Description: "Description of Name one for test 2",
		},
		expectedCode: 400,
	},
}

func TestPatch{{.Name}}sByID(t *testing.T) {

	ReturnTestApp()


	// Iterate through test single test cases
	for _, test := range tests{{.Name}}sPatchID {
		t.Run(test.name, func(t *testing.T) {

			//  changing post data to json
			patch_data, _ := json.Marshal(test.patch_data)
			
			req := httptest.NewRequest(http.MethodPatch,test.route, bytes.NewReader(patch_data))

			// Add specfic headers if needed as below
			req.Header.Set("Content-Type", "application/json")
			// req.Header.Set("X-APP-TOKEN", "hi")

			//  this is the response recorder
	
			resp, _ := TestApp.Test(req)

			// for debuging you can uncomment
			// fmt.Println("########")
			// fmt.Println(resp.StatusCode)
			// body, _ := io.ReadAll(resp.Result().Body)
			// fmt.Println(string(body))
			// fmt.Println("########")

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)

		})
	}

}

// ##########################################################################
// Define a structure for specifying input and output data
// of a single test case
var tests{{.Name}}sGet = []struct {
	name         string //name of string
	description  string // description of the test case
	route        string // route path to test
	expectedCode int    // expected HTTP status code
}{
	// First test case
	{
		name:         "get {{.Name}}s working - 1",
		description:  "get HTTP status 200",
		route:        "/"+group_path+"/{{.LowerName}}?page=1&size=10",
		expectedCode: 200,
	},
	// First test case
	{
		name:         "get {{.Name}}s working - 2",
		description:  "get HTTP status 200",
		route:        "/"+group_path+"/{{.LowerName}}?page=0&size=-5",
		expectedCode: 400,
	},
	// Second test case
	{
		name:         "get {{.Name}}s Working - 3",
		description:  "get HTTP status 404, when {{.Name}} Does not exist",
		route:        "/"+group_path+"/{{.LowerName}}?page=1&size=0",
		expectedCode: 400,
	},
}

func TestGet{{.Name}}s(t *testing.T) {
	ReturnTestApp()
	
	
	// Iterate through test single test cases
	for _, test := range tests{{.Name}}sGet {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.route, nil)

			// Add specfic headers if needed as below
			// req.Header.Set("X-APP-TOKEN", "hi")

			//  this is the response recorder
			

			resp, _ := TestApp.Test(req)

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)

		})
	}

}

// ##############################################################

var tests{{.Name}}sGetByID = []struct {
	name         string //name of string
	description  string // description of the test case
	route        string // route path to test
	expectedCode int    // expected HTTP status code
}{
	// First test case
	{
		name:         "get {{.Name}}s By ID check - 1",
		description:  "get Single {{.Name}} by ID",
		route:        "/"+group_path+"/{{.LowerName}}/1",
		expectedCode: 200,
	},

	// First test case
	{
		name:         "get {{.Name}}s By ID check - 2",
		description:  "get Single {{.Name}} by ID",
		route:        "/"+group_path+"/{{.LowerName}}/-1",
		expectedCode: 404,
	},
	// Second test case
	{
		name:         "get {{.Name}} By ID check - 3",
		description:  "get HTTP status 404, when {{.Name}} Does not exist",
		route:        "/"+group_path+"/{{.LowerName}}/1000",
		expectedCode: 404,
	},
}

func TestGet{{.Name}}sByID(t *testing.T) {

	ReturnTestApp()


	// Iterate through test single test cases
	for _, test := range tests{{.Name}}sGetByID {
		t.Run(test.name, func(t *testing.T) {
			
			req := httptest.NewRequest(http.MethodGet, test.route, nil)
		
			// Add specfic headers if needed as below
			// req.Header.Set("X-APP-TOKEN", "hi")

			resp, _ := TestApp.Test(req)
			

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)

		})
	}

}
`

var testAppTemplate = `
package tests

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"{{.ProjectName}}.com/models/controllers"
)

var (

	TestApp *fiber.App
	group_path string
)

func nextFunc(contx *fiber.Ctx) error {
	return contx.Next()
}

// initalaizing the app
func ReturnTestApp() {

	// loading env file
	godotenv.Load(".test.env")
	
	TestApp = fiber.New()

	group_path = "/api/v1"
	
	app := TestApp.Group(group_path)
	
	{{range .Models}}
		app.Get("/{{.LowerName}}",nextFunc).Name("get_all_{{.LowerName}}s").Get("/{{.LowerName}}", controllers.Get{{.Name}}s)
		app.Get("/{{.LowerName}}/:{{.LowerName}}_id",nextFunc).Name("get_one_{{.LowerName}}s").Get("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Get{{.Name}}ByID)
		app.Post("/{{.LowerName}}",nextFunc).Name("post_{{.LowerName}}").Post("/{{.LowerName}}", controllers.Post{{.Name}})
		app.Patch("/{{.LowerName}}/:{{.LowerName}}_id",nextFunc).Name("patch_{{.LowerName}}").Patch("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Patch{{.Name}})
		app.Delete("/{{.LowerName}}/:{{.LowerName}}_id",nextFunc).Name("delete_{{.LowerName}}").Delete("/{{.LowerName}}/:{{.LowerName}}_id", controllers.Delete{{.Name}}).Name("delete_{{.LowerName}}")

		{{range .Relations}}
		app.Post("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id/:{{.LowerParentName}}_id",nextFunc).Name("add_{{.LowerFieldName}}{{.LowerParentName}}").Post("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id/:{{.LowerParentName}}_id",controllers.Add{{.FieldName}}{{.ParentName}}s)
		app.Delete("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id/:{{.LowerParentName}}_id",nextFunc).Name("delete_{{.LowerFieldName}}{{.LowerParentName}}").Delete("/{{.LowerFieldName}}{{.LowerParentName}}/:{{.LowerFieldName}}_id/:{{.LowerParentName}}_id",controllers.Delete{{.FieldName}}{{.ParentName}}s)
		{{end}}
	{{end}}

	

}


`
