package temps

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

func MongoCurdFrameFiber() {

	// ############################################################
	curd_tmpl, err := template.New("RenderData").Parse(mongocurdTemplateFiber)
	if err != nil {
		panic(err)
	}
	// Create the models directory if it does not exist
	// #################################################
	err = os.MkdirAll("controllers", os.ModePerm)
	if err != nil {
		panic(err)
	}
	for _, model := range RenderData.Models {
		folder_path := fmt.Sprintf("controllers/%v_controller.go", model.Name)
		folder_path = strings.ToLower(folder_path)
		curd_file, err := os.Create(folder_path)
		if err != nil {
			panic(err)
		}

		err = curd_tmpl.Execute(curd_file, model)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		curd_file.Close()
	}
}

var mongocurdTemplateFiber = `
package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"{{.ProjectName}}.com/common"
	"{{.ProjectName}}.com/nosqlmodels"
	"{{.ProjectName}}.com/observe"
)

// Get{{.Name}}is a function to get a {{.Name}}s by ID
// @Summary Get {{.Name}}s
// @Description Get {{.Name}}s
// @Tags {{.Name}}s
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security Refresh
// @Param page query int true "page"
// @Param size query int true "page size"
// @Success 200 {object} common.ResponsePagination{data=[]models.{{.Name}}Get}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /{{.LowerName}} [get]
func Get{{.Name}}s(contx *fiber.Ctx) error {

	//  Getting tracer context
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)


	//  Getting Database connection
	db, _ := contx.Locals("db").(*mongo.Database)


	//  parsing Query Prameters
	Page, _ := strconv.Atoi(contx.Query("page"))
	Limit, _ := strconv.Atoi(contx.Query("size"))
	//  checking if query parameters  are correct
	if Page == 0 || Limit == 0 {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Not Allowed, Bad request",
			Data:    nil,
		})
	}


	// ###################################################################
	// here is were mongo db curd calls from mongo driver are to be called
	// ###################################################################

	// returning result if all the above completed successfully
	return contx.Status(http.StatusOK).JSON(result)
}

// Get{{.Name}}ByID is a function to get a {{.Name}}s by ID
// @Summary Get {{.Name}} by ID
// @Description Get {{.LowerName}} by ID
// @Tags {{.Name}}s
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param {{.LowerName}}_id path int true "{{.Name}} ID"
// @Success 200 {object} common.ResponseHTTP{data=models.{{.Name}}Get}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /{{.LowerName}}/{{ "{" }}{{.LowerName}}_id{{ "}" }} [get]
func Get{{.Name}}ByID(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)


	//  Getting Database connection
	db, _ := contx.Locals("db").(*mongo.Database)


	//  parsing Query Prameters
	id, err := strconv.Atoi(contx.Params("{{.LowerName}}_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}


	// Preparing and querying database using Gorm
	var {{.LowerName}}s_get models.{{.Name}}Get
	var {{.LowerName}}s models.{{.Name}}
	// ###################################################################
	// here is were mongo db curd calls from mongo driver are to be called
	// ###################################################################
	//
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Error retrieving {{.Name}}",
			Data:    nil,
		})
	}

	// filtering response data according to filtered defined struct
	mapstructure.Decode({{.LowerName}}s, &{{.LowerName}}s_get)

	//  Finally returing response if All the above compeleted successfully
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success got one {{.LowerName}}.",
		Data:    &{{.LowerName}}s_get,
	})
}

// Add {{.Name}} to data
// @Summary Add a new {{.Name}}
// @Description Add {{.Name}}
// @Tags {{.Name}}s
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param {{.LowerName}} body models.{{.Name}}Post true "Add {{.Name}}"
// @Success 200 {object} common.ResponseHTTP{data=models.{{.Name}}Post}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /{{.LowerName}} [post]
func Post{{.Name}}(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)


	// Getting Database Connection
	db, _ := contx.Locals("db").(*mongo.Database)



	// validator initialization
	validate := validator.New()

	//validating post data
	posted_{{.LowerName}} := new(models.{{.Name}}Post)

	//first parse request data
	if err := contx.BodyParser(&posted_{{.LowerName}}); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validate structure
	if err := validate.Struct(posted_{{.LowerName}}); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}


	// ###################################################################
	// here is were mongo db curd calls from mongo driver are to be called
	// ###################################################################

	// return data if transaction is sucessfull
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "{{.Name}} created successfully.",
		Data:    {{.LowerName}},
	})
}

// Patch {{.Name}} to data
// @Summary Patch {{.Name}}
// @Description Patch {{.Name}}
// @Tags {{.Name}}s
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param {{.LowerName}} body models.{{.Name}}Post true "Patch {{.Name}}"
// @Param {{.LowerName}}_id path int true "{{.Name}} ID"
// @Success 200 {object} common.ResponseHTTP{data=models.{{.Name}}Post}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /{{.LowerName}}/{{ "{" }}{{.LowerName}}_id{{ "}" }} [patch]
func Patch{{.Name}}(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)


	// Get database connection
	db, _ := contx.Locals("db").(*mongo.Database)

	// ###################################################################
	// here is were mongo db curd calls from mongo driver are to be called
	// ###################################################################

	//  initialize data validator
	validate := validator.New()

	// validate path params
	id, err := strconv.Atoi(contx.Params("{{.LowerName}}_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// validate data struct
	patch_{{.LowerName}} := new(models.{{.Name}}Patch)
	if err := contx.BodyParser(&patch_{{.LowerName}}); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validating
	if err := validate.Struct(patch_{{.LowerName}}); err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// ###################################################################
	// here is were mongo db curd calls from mongo driver are to be called
	// ###################################################################


	// Return  success response
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "{{.Name}} updated successfully.",
		Data:    {{.LowerName}},
	})
}

// Delete{{.Name}}s function removes a {{.LowerName}} by ID
// @Summary Remove {{.Name}} by ID
// @Description Remove {{.LowerName}} by ID
// @Tags {{.Name}}s
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param {{.LowerName}}_id path int true "{{.Name}} ID"
// @Success 200 {object} common.ResponseHTTP{}
// @Failure 404 {object} common.ResponseHTTP{}
// @Failure 503 {object} common.ResponseHTTP{}
// @Router /{{.LowerName}}/{{ "{" }}{{.LowerName}}_id{{ "}" }} [delete]
func Delete{{.Name}}(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)


	// Getting Database connection
	db, err := contx.Locals("db").(*mongo.Database)
	if err != nil {
		return ctx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// ###################################################################
	// here is were mongo db curd calls from mongo driver are to be called
	// ###################################################################

	// validate path params
	id, err := strconv.Atoi(contx.Params("{{.LowerName}}_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}


	// perform delete operation if the object exists
	tx := db.WithContext(tracer.Tracer).Begin()



	// ###################################################################
	// here is were mongo db curd calls from mongo driver are to be called
	// ###################################################################

	// Return success respons
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "{{.Name}} deleted successfully.",
		Data:    {{.LowerName}},
	})
}

`
