package temps

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

func TestFrameEcho() {

	// ############################################################

	test_tmpl, err := template.New("RenderData").Parse(testTemplateEcho)
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

}

var testTemplateEcho = `
package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
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
		route:       "/admin/{{.LowerName}}",
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
		route:       "/admin/{{.LowerName}}",
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
		route:       "/admin/{{.LowerName}}",
		post_data: models.{{.Name}}Post{
			Name:        "Name one",
			Description: "Description of Name one",
		},
		expectedCode: 500,
	},
}

func TestPost{{.Name}}sByID(t *testing.T) {

	// loading env file
	godotenv.Load(".test.env")

	// Setup Test APP
	TestApp := echo.New()

	// Iterate through test single test cases
	for _, test := range tests{{.Name}}sPostID {
		t.Run(test.name, func(t *testing.T) {
			//  changing post data to json
			post_data, _ := json.Marshal(test.post_data)

			req := httptest.NewRequest(http.MethodPost, test.route, bytes.NewReader(post_data))

			// Add specfic headers if needed as below
			req.Header.Set("Content-Type", "application/json")
			// req.Header.Set("X-APP-TOKEN", "hi")

			//  this is the response recorder
			resp := httptest.NewRecorder()

			//  create echo context to test the app function
			echo_contx := TestApp.NewContext(req, resp)
			echo_contx.SetPath(test.route)

			// Now testing the Get{{.Name}}s funciton
			controllers.Post{{.Name}}(echo_contx)

			var responseMap map[string]interface{}
			body, _ := io.ReadAll(resp.Body)
			uerr := json.Unmarshal(body, &responseMap)
			if uerr != nil {
				// fmt.Printf("Error marshaling response : %v", uerr)
				fmt.Println()
			}

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.Result().StatusCode, test.description)
			//  running delete test if post is success
			if resp.Result().StatusCode == 200 {
				t.Run("Checking the Delete Request Path for {{.Name}}s", func(t *testing.T) {

					test_route := fmt.Sprintf("%v/:%v", test.route, "{{.LowerName}}_id")

					req_delete := httptest.NewRequest(http.MethodDelete, test.route, bytes.NewReader(post_data))

					// Add specfic headers if needed as below
					req_delete.Header.Set("Content-Type", "application/json")

					//  this is the response recorder
					resp_delete := httptest.NewRecorder()

					//  create echo context to test the app function
					echo_contx_del := TestApp.NewContext(req_delete, resp_delete)
					echo_contx_del.SetPath(test_route)

					// seting path paramenters
					path_value := fmt.Sprintf("%v", responseMap["data"].(map[string]interface{})["id"])
					echo_contx_del.SetParamNames("{{.LowerName}}_id")
					echo_contx_del.SetParamValues(path_value)

					// Now testing the Get{{.Name}}s funciton
					controllers.Delete{{.Name}}(echo_contx_del)
					assert.Equalf(t, 200, resp.Result().StatusCode, test.description+"deleteing")
				})
			} else {
				t.Run("Checking the Delete Request Path for  that does not exit", func(t *testing.T) {

					test_route := fmt.Sprintf("%v/:%v", test.route, "{{.LowerName}}_id")

					req_delete := httptest.NewRequest(http.MethodDelete, test.route, bytes.NewReader(post_data))

					// Add specfic headers if needed as below
					req_delete.Header.Set("Content-Type", "application/json")

					//  this is the response recorder
					resp_delete := httptest.NewRecorder()

					//  create echo context to test the app function
					echo_contx_del := TestApp.NewContext(req_delete, resp_delete)
					echo_contx_del.SetPath(test_route)

					// seting path paramenters
					path_value := fmt.Sprintf("%v", 2000)
					echo_contx_del.SetParamNames("{{.LowerName}}_id")
					echo_contx_del.SetParamValues(path_value)

					// Now testing the Get{{.Name}}s funciton
					controllers.Delete{{.Name}}(echo_contx_del)
					assert.Equalf(t, 500, resp.Result().StatusCode, test.description+"deleteing")
				})

				t.Run("Checking the Delete Request Path that is not valid", func(t *testing.T) {

					test_route := fmt.Sprintf("%v/:%v", test.route, "{{.LowerName}}_id")

					req_delete := httptest.NewRequest(http.MethodDelete, test.route, bytes.NewReader(post_data))

					// Add specfic headers if needed as below
					req_delete.Header.Set("Content-Type", "application/json")

					//  this is the response recorder
					resp_delete := httptest.NewRecorder()

					//  create echo context to test the app function
					echo_contx_del := TestApp.NewContext(req_delete, resp_delete)
					echo_contx_del.SetPath(test_route)

					// seting path paramenters
					path_value := fmt.Sprintf("%v", "@@")
					echo_contx_del.SetParamNames("{{.LowerName}}_id")
					echo_contx_del.SetParamValues(path_value)

					// Now testing the Get{{.Name}}s funciton
					controllers.Delete{{.Name}}(echo_contx_del)
					assert.Equalf(t, 500, resp.Result().StatusCode, test.description+"deleteing")
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
	{{.LowerName}}_id      string           //path param
	patch_data   models.{{.Name}}Patch // patch_data
	expectedCode int              // expected HTTP status code
}{
	// First test case
	{
		name:        "patch {{.Name}}s By ID check - 1",
		description: "patch Single {{.Name}} by ID",
		route:       "/admin/{{.LowerName}}/:{{.LowerName}}_id",
		{{.LowerName}}_id:     "1",
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
		route:       "/admin/{{.LowerName}}/:{{.LowerName}}_id",
		{{.LowerName}}_id:     "100",
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
		route:       "/admin/{{.LowerName}}/:{{.LowerName}}_id",
		{{.LowerName}}_id:     "@@",
		patch_data: models.{{.Name}}Patch{
			Name:        "Name one eight",
			Description: "Description of Name one for test 2",
		},
		expectedCode: 400,
	},
}

func TestPatch{{.Name}}sByID(t *testing.T) {

	// loading env file
	godotenv.Load(".test.env")

	// Setup Test APP
	TestApp := echo.New()

	// Iterate through test single test cases
	for _, test := range tests{{.Name}}sPatchID {
		t.Run(test.name, func(t *testing.T) {

			//  changing post data to json
			patch_data, _ := json.Marshal(test.patch_data)

			req := httptest.NewRequest(http.MethodPatch, test.route, bytes.NewReader(patch_data))

			// Add specfic headers if needed as below
			req.Header.Set("Content-Type", "application/json")
			// req.Header.Set("X-APP-TOKEN", "hi")

			//  this is the response recorder
			resp := httptest.NewRecorder()

			//  create echo context to test the app function
			echo_contx := TestApp.NewContext(req, resp)
			echo_contx.SetPath(test.route)

			// seting path paramenters
			echo_contx.SetParamNames("{{.LowerName}}_id")
			echo_contx.SetParamValues(test.{{.LowerName}}_id)

			// Now testing the Get{{.Name}}s funciton
			controllers.Patch{{.Name}}(echo_contx)

			// fmt.Println("########")
			// fmt.Println(resp.Result().StatusCode)
			// body, _ := io.ReadAll(resp.Result().Body)
			// fmt.Println(string(body))
			// fmt.Println("########")

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.Result().StatusCode, test.description)

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
		route:        "/admin/{{.LowerName}}?page=1&size=10",
		expectedCode: 200,
	},
	// First test case
	{
		name:         "get {{.Name}}s working - 2",
		description:  "get HTTP status 200",
		route:        "/admin/{{.LowerName}}?page=0&size=-5",
		expectedCode: 400,
	},
	// Second test case
	{
		name:         "get {{.Name}}s Working - 3",
		description:  "get HTTP status 404, when {{.Name}} Does not exist",
		route:        "/admin/{{.LowerName}}?page=1&size=0",
		expectedCode: 400,
	},
}

func TestGet{{.Name}}s(t *testing.T) {
	// loading env file
	godotenv.Load(".test.env")

	// Setup Test APP
	TestApp := echo.New()

	// Iterate through test single test cases
	for _, test := range tests{{.Name}}sGet {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.route, nil)

			// Add specfic headers if needed as below
			// req.Header.Set("X-APP-TOKEN", "hi")

			//  this is the response recorder
			resp := httptest.NewRecorder()

			//  create echo context to test the app function
			echo_contx := TestApp.NewContext(req, resp)
			echo_contx.SetPath(test.route)
			// Now testing the Get{{.Name}}s funciton
			controllers.Get{{.Name}}s(echo_contx)

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.Result().StatusCode, test.description)

		})
	}

}

// ##############################################################

var tests{{.Name}}sGetByID = []struct {
	name         string //name of string
	description  string // description of the test case
	route        string // route path to test
	{{.LowerName}}_id      string // path parm
	expectedCode int    // expected HTTP status code
}{
	// First test case
	{
		name:         "get {{.Name}}s By ID check - 1",
		description:  "get Single {{.Name}} by ID",
		route:        "/admin/{{.LowerName}}/:{{.LowerName}}_id",
		{{.LowerName}}_id:      "1",
		expectedCode: 200,
	},

	// First test case
	{
		name:         "get {{.Name}}s By ID check - 2",
		description:  "get Single {{.Name}} by ID",
		route:        "/admin/{{.LowerName}}/:{{.LowerName}}_id",
		{{.LowerName}}_id:      "-1",
		expectedCode: 404,
	},
	// Second test case
	{
		name:         "get {{.Name}} By ID check - 3",
		description:  "get HTTP status 404, when {{.Name}} Does not exist",
		route:        "/admin/{{.LowerName}}/:{{.LowerName}}_id",
		{{.LowerName}}_id:      "100",
		expectedCode: 404,
	},
}

func TestGet{{.Name}}sByID(t *testing.T) {

	// loading env file
	godotenv.Load(".test.env")

	// Setup Test APP
	TestApp := echo.New()

	// Iterate through test single test cases
	for _, test := range tests{{.Name}}sGetByID {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.route, nil)

			// Add specfic headers if needed as below
			// req.Header.Set("X-APP-TOKEN", "hi")

			//  this is the response recorder
			resp := httptest.NewRecorder()

			//  create echo context to test the app function
			echo_contx := TestApp.NewContext(req, resp)
			echo_contx.SetPath(test.route)

			// seting path paramenters
			echo_contx.SetParamNames("{{.LowerName}}_id")
			echo_contx.SetParamValues(test.{{.LowerName}}_id)

			// Now testing the Get{{.Name}}s funciton
			controllers.Get{{.Name}}ByID(echo_contx)

			//  Finally asserting test cases
			assert.Equalf(t, test.expectedCode, resp.Result().StatusCode, test.description)

		})
	}

}
`
