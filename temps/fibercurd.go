package temps

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

func CurdFrameFiber() {

	// ############################################################

	curd_tmpl, err := template.New("RenderData").Parse(curdTemplateFiber)
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

var curdTemplateFiber = `
package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"{{.ProjectName}}.com/common"
	"{{.ProjectName}}.com/models"
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
	db, _ := contx.Locals("db").(*gorm.DB)

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


	//  querying result with pagination using gorm function
	result, err := common.PaginationPureModel(db, models.{{.Name}}{}, []models.{{.Name}}{}, uint(Page), uint(Limit), tracer.Tracer)
	if err != nil {
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Failed to get all {{.Name}}.",
			Data:    "something",
		})
	}

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
	db, _ := contx.Locals("db").(*gorm.DB)

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
	if res := db.WithContext(tracer.Tracer).Model(&models.{{.Name}}{}).Preload(clause.Associations).Where("id = ?", id).First(&{{.LowerName}}s); res.Error != nil {
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
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
	db, _ := contx.Locals("db").(*gorm.DB)


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

	//  initiate -> {{.LowerName}}
	{{.LowerName}} := new(models.{{.Name}})
	{{.LowerName}}.Name = posted_{{.LowerName}}.Name
	{{.LowerName}}.Description = posted_{{.LowerName}}.Description

	//  start transaction to database
	tx := db.WithContext(tracer.Tracer).Begin()

	// add  data using transaction if values are valid
	if err := tx.Create(&{{.LowerName}}).Error; err != nil {
		tx.Rollback()
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "{{.Name}} Creation Failed",
			Data:    err,
		})
	}

	// close transaction
	tx.Commit()

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
	db, _ := contx.Locals("db").(*gorm.DB)

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

	// startng update transaction
	var {{.LowerName}} models.{{.Name}}
	{{.LowerName}}.ID = uint(id)
	tx := db.WithContext(tracer.Tracer).Begin()

	// Check if the record exists
	if err := db.WithContext(tracer.Tracer).Where("id = ? ", id).First(&{{.LowerName}}).Error; err != nil {
		tx.Rollback()
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Update the record
	if err := db.WithContext(tracer.Tracer).Model(&{{.LowerName}}).UpdateColumns(*patch_{{.LowerName}}).Error; err != nil {
		tx.Rollback()
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

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
	db, _ := contx.Locals("db").(*gorm.DB)

	// get deleted {{.LowerName}} attributes to return
	var {{.LowerName}} models.{{.Name}}

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

	// first getting {{.LowerName}} and checking if it exists
	if err := db.WithContext(tracer.Tracer).Where("id = ?", id).First(&{{.LowerName}}).Error; err != nil {
		tx.Rollback()
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Delete the {{.LowerName}}
	if err := db.WithContext(tracer.Tracer).Delete(&{{.LowerName}}).Error; err != nil {
		tx.Rollback()
		return contx.Status(http.StatusInternalServerError).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Error deleting {{.LowerName}}",
			Data:    nil,
		})
	}

	// Commit the transaction
	tx.Commit()

	// Return success respons
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "{{.Name}} deleted successfully.",
		Data:    {{.LowerName}},
	})
}

// ################################################################
// Relationship Based Endpoints
// ################################################################

{{ range .Relations }}
{{if .MtM}}

// Add {{.FieldName}} to {{.ParentName}}
// @Summary Add {{.ParentName}} to {{.FieldName}}
// @Description Add {{.FieldName}} {{.ParentName}}
// @Tags {{.FieldName}}{{.ParentName}}s
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param {{.LowerFieldName}}_id path int true "{{.FieldName}} ID"
// @Param {{.LowerParentName}}_id path int true "{{.ParentName}} ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /{{.LowerFieldName}}{{.LowerParentName}}/{{ "{" }}{{.LowerFieldName}}_id{{ "}" }}/{{ "{" }}{{.LowerParentName}}_id{{ "}" }} [post]
func Add{{.FieldName}}{{.ParentName}}s(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)


	// database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	// validate path params
	{{.LowerFieldName}}_id, err := strconv.Atoi(contx.Params("{{.LowerFieldName}}_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// validate path params
	{{.LowerParentName}}_id, err := strconv.Atoi(contx.Params("{{.LowerParentName}}_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// fetching {{.LowerParentName}} to be added
	var {{.LowerParentName}} models.{{.ParentName}}
	if res := db.WithContext(tracer.Tracer).Where("id = ?", {{.LowerParentName}}_id ).First(&{{.LowerParentName}}); res.Error != nil {
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	//  {{.LowerParentName}}ending assocation
	var {{.LowerFieldName}} models.{{.FieldName}}
	if err := db.WithContext(tracer.Tracer).Where("id = ?",{{.LowerFieldName}}_id).First(&{{.LowerFieldName}}); err.Error != nil {
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Record not Found",
			Data:    err.Error.Error(),
		})
	}

	tx := db.WithContext(tracer.Tracer).Begin()
	if err := db.WithContext(tracer.Tracer).Model(&{{.LowerFieldName}}).Association("{{.ParentName}}s").Append(&{{.LowerParentName}}); err != nil {
		tx.Rollback()
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: "{{.ParentName}}ending {{.ParentName}} Failed",
			Data:    err.Error(),
		})
	}
	tx.Commit()

	// return value if transaction is sucessfull
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success Creating a {{.LowerParentName}} {{.FieldName}}.",
		Data:    {{.LowerParentName}},
	})
}

// Delete {{.ParentName}} to {{.FieldName}}
// @Summary Add {{.ParentName}}
// @Description Delete {{.FieldName}} {{.ParentName}}
// @Tags {{.FieldName}}s
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param {{.LowerFieldName}}_id path int true "{{.FieldName}} ID"
// @Param {{.LowerParentName}}_id path int true "{{.ParentName}} ID"
// @Success 200 {object} common.ResponseHTTP{data=models.{{.ParentName}}Post}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /{{.LowerFieldName}}{{.LowerParentName}}/{{ "{" }}{{.LowerFieldName}}_id{{ "}" }}/{{ "{" }}{{.LowerParentName}}_id{{ "}" }} [delete]
func Delete{{.FieldName}}{{.ParentName}}s(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)


	//Connect to Database
	db, _ := contx.Locals("db").(*gorm.DB)

	// validate path params
	{{.LowerFieldName}}_id, err := strconv.Atoi(contx.Params("{{.LowerFieldName}}_id"))
	if err != nil || {{.LowerFieldName}}_id == 0 {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	{{.LowerParentName}}_id, err := strconv.Atoi(contx.Params("{{.LowerParentName}}_id"))
	if err != nil || {{.LowerParentName}}_id == 0 {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}
	// fetching {{.LowerParentName}} to be deleted
	var {{.LowerParentName}} models.{{.ParentName}}
	if res := db.WithContext(tracer.Tracer).Where("id = ?", {{.LowerParentName}}_id).First(&{{.LowerParentName}}); res.Error != nil {
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// fettchng {{.LowerFieldName}}
	var {{.LowerFieldName}} models.{{.FieldName}}
	if err := db.WithContext(tracer.Tracer).Where("id = ?",{{.LowerFieldName}}_id).First(&{{.LowerFieldName}}); err.Error != nil {
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Record not Found",
			Data:    err.Error.Error(),
		})
	}

	// removing {{.LowerParentName}}
	tx := db.WithContext(tracer.Tracer).Begin()
	if err := db.WithContext(tracer.Tracer).Model(&{{.LowerFieldName}}).Association("{{.ParentName}}s").Delete(&{{.LowerParentName}}); err != nil {
		tx.Rollback()
		return contx.Status(http.StatusNonAuthoritativeInfo).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Please Try Again Something Unexpected H{{.LowerParentName}}ened",
			Data:    err.Error(),
		})
	}

	tx.Commit()

	// return value if transaction is sucessfull
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success Removing a {{.LowerParentName}} from {{.LowerFieldName}}.",
		Data:    {{.LowerParentName}},
	})
}


{{ end}}
{{ end}}

{{ range .Relations }}
{{if .OtM}}

// Add {{.ParentName}} {{.FieldName}}
// @Summary Add {{.ParentName}} to {{.FieldName}}
// @Description Add {{.ParentName}} to {{.FieldName}}
// @Tags {{.ParentName}}s
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param {{.LowerFieldName}}_id path int true "{{.FieldName}} ID"
// @Param {{.LowerParentName}}_id query int true " {{.ParentName}} ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /{{.LowerFieldName}}{{.LowerParentName}}/{{ "{" }}{{.LowerFieldName}}_id{{ "}" }} [patch]
func Add{{.FieldName}}{{.ParentName}}s(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)


	// connect
	db, _ := contx.Locals("db").(*gorm.DB)

	// validate path params
	{{.LowerFieldName}}_id, err := strconv.Atoi(contx.Params("{{.LowerFieldName}}_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// fetching Endpionts
	var {{.LowerFieldName}} models.{{.FieldName}}
	if res := db.WithContext(tracer.Tracer).Model(&models.{{.FieldName}}{}).Where("id = ?", {{.LowerFieldName}}_id).First(&{{.LowerFieldName}}); res.Error != nil {
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// fetching {{.LowerFieldName}} to be added
	{{.LowerParentName}}_id, _ := strconv.Atoi(contx.Query("{{.LowerParentName}}_id"))
	var {{.LowerParentName}} models.{{.ParentName}}
	if res := db.WithContext(tracer.Tracer).Model(&models.{{.ParentName}}{}).Where("id = ?", {{.LowerParentName}}_id).First(&{{.LowerParentName}}); res.Error != nil {
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// startng update transaction

	tx := db.WithContext(tracer.Tracer).Begin()
	//  Adding one to many Relation
	if err := db.WithContext(tracer.Tracer).Model(&{{.LowerParentName}}).Association("{{.FieldName}}s").Append(&{{.LowerFieldName}}); err != nil {
		tx.Rollback()
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Error Adding Record",
			Data:    err.Error(),
		})
	}
	tx.Commit()

	// return value if transaction is sucessfull
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success Adding a {{.FieldName}} to {{.ParentName}}.",
		Data:    {{.LowerParentName}},
	})
}

// Delete {{.ParentName}} {{.FieldName}}
// @Summary Delete {{.ParentName}} {{.FieldName}}
// @Description Delete {{.ParentName}} {{.FieldName}}
// @Tags {{.ParentName}}s
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param {{.LowerFieldName}}_id path int true "{{.ParentName}} ID"
// @Param {{.LowerParentName}}_id query int true "{{.FieldName}} ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /{{.LowerFieldName}}{{.LowerParentName}}/{{ "{" }}{{.LowerFieldName}}_id{{ "}" }} [delete]
func Delete{{.FieldName}}{{.ParentName}}s(contx *fiber.Ctx) error {

	// Starting tracer context and tracer
	ctx := contx.Locals("tracer")
	tracer, _ := ctx.(*observe.RouteTracer)


	//  database connection
	db, _ := contx.Locals("db").(*gorm.DB)

	// validate path params
	{{.LowerFieldName}}_id, err := strconv.Atoi(contx.Params("{{.LowerFieldName}}_id"))
	if err != nil {
		return contx.Status(http.StatusBadRequest).JSON(common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Getting {{.FieldName}}
	var {{.LowerFieldName}} models.{{.FieldName}}
	if res := db.WithContext(tracer.Tracer).Model(&models.{{.FieldName}}{}).Where("id = ?", {{.LowerFieldName}}_id).First(&{{.LowerFieldName}}); res.Error != nil {
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// fetching {{.LowerParentName}} to be added
	var {{.LowerParentName}} models.{{.ParentName}}
	{{.LowerParentName}}_id, _ := strconv.Atoi(contx.Query("{{.LowerParentName}}_id"))
	if res := db.WithContext(tracer.Tracer).Model(&models.{{.ParentName}}{}).Where("id = ?", {{.LowerParentName}}_id).First(&{{.LowerParentName}}); res.Error != nil {
		return contx.Status(http.StatusServiceUnavailable).JSON(common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// Removing {{.FieldName}} From {{.ParentName}}
	tx := db.WithContext(tracer.Tracer).Begin()
	if err := db.WithContext(tracer.Tracer).Model(&{{.LowerParentName}}).Association("{{.FieldName}}s").Delete(&{{.LowerFieldName}}); err != nil {
		tx.Rollback()
		return contx.Status(http.StatusNotFound).JSON(common.ResponseHTTP{
			Success: false,
			Message: "Record not Found",
			Data:    err.Error(),
		})
	}
	tx.Commit()

	// return value if transaction is sucessfull
	return contx.Status(http.StatusOK).JSON(common.ResponseHTTP{
		Success: true,
		Message: "Success Deleteing a {{.FieldName}} From {{.ParentName}}.",
		Data:    {{.LowerParentName}},
	})
}


{{ end}}
{{ end}}


`
