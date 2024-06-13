package temps

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"
)

func CurdFrame() {

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
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	// setting default value for config data file
	//  GetPostPatchPut
	// "Get$Post$Patch$Put"

	for i := 0; i < len(data.Models); i++ {
		data.Models[i].LowerName = strings.ToLower(data.Models[i].Name)
		data.Models[i].AppName = data.AppName
		data.Models[i].ProjectName = data.ProjectName

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

	// ############################################################

	curd_tmpl, err := template.New("data").Parse(curdTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	// #################################################
	err = os.MkdirAll("models/controllers", os.ModePerm)
	if err != nil {
		panic(err)
	}

	for _, model := range data.Models {

		folder_path := fmt.Sprintf("models/controllers/%v_controller.go", model.Name)
		folder_path = strings.ToLower(folder_path)
		curd_file, err := os.Create(folder_path)
		if err != nil {
			panic(err)
		}

		err = curd_tmpl.Execute(curd_file, model)
		if err != nil {
			panic(err)
		}
		curd_file.Close()

	}

	// running go mod tidy finally
	if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
		fmt.Printf("error: %v \n", err)
	}
}

var curdTemplate = `
package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"{{.ProjectName}}.com/common"
	"{{.ProjectName}}.com/database"
	"{{.ProjectName}}.com/models"
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
func Get{{.Name}}s(contx echo.Context) error {

	//  parsing Query Prameters
	Page, _ := strconv.Atoi(contx.QueryParam("page"))
	Limit, _ := strconv.Atoi(contx.QueryParam("size"))
	//  checking if query parameters  are correct
	if Page == 0 || Limit == 0 {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: "Not Allowed, Bad request",
			Data:    nil,
		})
	}

	//  Getting Database connection
	db := database.ReturnSession()

	//  querying result with pagination using gorm function
	result, err := common.PaginationPureModel(db, models.{{.Name}}{}, []models.{{.Name}}{}, uint(Page), uint(Limit))
	if err != nil {
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: "Success get all {{.Name}}.",
			Data:    "something",
		})
	}

	// returning result if all the above completed successfully
	return contx.JSON(http.StatusOK, result)
}

// Get{{.Name}}ByID is a function to get a {{.Name}}s by ID
// @Summary Get {{.Name}} by ID
// @Description Get {{.LowerName}} by ID
// @Tags {{.Name}}s
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "{{.Name}} ID"
// @Success 200 {object} common.ResponseHTTP{data=models.{{.Name}}Get}
// @Failure 404 {object} common.ResponseHTTP{}
// @Router /{{.LowerName}}/{{ "{" }}{{.LowerName}}_id{{ "}" }} [get]
func Get{{.Name}}ByID(contx echo.Context) error {

	//  parsing Query Prameters
	id, err := strconv.Atoi(contx.Param("{{.LowerName}}_id"))
	if err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	//  Getting Database connection
	db := database.ReturnSession()

	// Preparing and querying database using Gorm
	var {{.LowerName}}s_get models.{{.Name}}Get
	var {{.LowerName}}s models.{{.Name}}
	if res := db.Model(&models.{{.Name}}{}).Preload(clause.Associations).Where("id = ?", id).First(&{{.LowerName}}s); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return contx.JSON(http.StatusNotFound, common.ResponseHTTP{
				Success: false,
				Message: "{{.Name}} not found",
				Data:    nil,
			})
		}
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: "Error retrieving {{.Name}}",
			Data:    nil,
		})
	}

	// filtering response data according to filtered defined struct
	mapstructure.Decode({{.LowerName}}s, &{{.LowerName}}s_get)

	//  Finally returing response if All the above compeleted successfully
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "Success got one {{.LowerName}}.",
		Data:    &{{.LowerName}}s_get,
	})
}

// Add {{.Name}} to data
// @Summary Add a new {{.Name}}
// @Description Add {{.Name}}
// @Tags {{.Name}}
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param {{.LowerName}} body {{.Name}}Post true "Add {{.Name}}"
// @Success 200 {object} common.ResponseHTTP{data={{.Name}}Post}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /{{.LowerName}}s [post]
func Post{{.Name}}(contx echo.Context) error {
	//  parsing Query Prameters
	db := database.ReturnSession()

	// validator initialization
	validate := validator.New()

	//validating post data
	posted_{{.LowerName}} := new(models.{{.Name}}Post)

	//first parse request data
	if err := contx.Bind(&posted_{{.LowerName}}); err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validate structure
	if err := validate.Struct(posted_{{.LowerName}}); err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
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
	tx := db.Begin()

	// add  data using transaction if values are valid
	if err := tx.Create(&{{.LowerName}}).Error; err != nil {
		tx.Rollback()
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: "{{.Name}} Creation Failed",
			Data:    err,
		})
	}

	// close transaction
	tx.Commit()

	// return data if transaction is sucessfull
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "{{.Name}} created successfully.",
		Data:    {{.LowerName}},
	})
}

// Patch {{.Name}} to data
// @Summary Patch {{.Name}}
// @Description Patch {{.Name}}
// @Tags {{.Name}}
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param {{.LowerName}} body {{.Name}}Post true "Patch {{.Name}}"
// @Param id path int true "{{.Name}} ID"
// @Success 200 {object} common.ResponseHTTP{data=models.{{.Name}}Post}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /{{.LowerName}}/{{ "{" }}{{.LowerName}}_id{{ "}" }} [patch]
func Patch{{.Name}}(contx echo.Context) error {

	// Get database connection
	db := database.ReturnSession()

	//  initialize data validator
	validate := validator.New()

	// validate path params
	id, err := strconv.Atoi(contx.Param("{{.LowerName}}_id"))
	if err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// validate data struct
	patch_{{.LowerName}} := new(models.{{.Name}}Patch)
	if err := contx.Bind(&patch_{{.LowerName}}); err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// then validating
	if err := validate.Struct(patch_{{.LowerName}}); err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// startng update transaction
	var {{.LowerName}} models.{{.Name}}
	{{.LowerName}}.ID = uint(id)
	tx := db.Begin()

	// Check if the record exists
	if err := db.First(&{{.LowerName}}, {{.LowerName}}.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If the record doesn't exist, return an error response
			tx.Rollback()
			return contx.JSON(http.StatusNotFound, common.ResponseHTTP{
				Success: false,
				Message: "{{.Name}} not found",
				Data:    nil,
			})
		}
		// If there's an unexpected error, return an internal server error response
		tx.Rollback()
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Update the record
	if err := db.Model(&{{.LowerName}}).UpdateColumns(*patch_{{.LowerName}}).Error; err != nil {
		tx.Rollback()
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Return  success response
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "{{.Name}} updated successfully.",
		Data:    {{.LowerName}},
	})
}

// Delete{{.Name}}s function removes a {{.LowerName}} by ID
// @Summary Remove {{.Name}} by ID
// @Description Remove {{.LowerName}} by ID
// @Tags {{.Name}}
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "{{.Name}} ID"
// @Success 200 {object} common.ResponseHTTP{}
// @Failure 404 {object} common.ResponseHTTP{}
// @Failure 503 {object} common.ResponseHTTP{}
// @Router /{{.LowerName}}/{{ "{" }}{{.LowerName}}_id{{ "}" }} [delete]
func Delete{{.Name}}(contx echo.Context) error {

	// get deleted {{.LowerName}} attributes to return
	var {{.LowerName}} models.{{.Name}}

	// validate path params
	id, err := strconv.Atoi(contx.Param("{{.LowerName}}_id"))
	if err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Getting Database connection
	db := database.ReturnSession()

	// perform delete operation if the object exists
	tx := db.Begin()

	// first getting {{.LowerName}} and checking if it exists
	if err := db.Where("id = ?", id).First(&{{.LowerName}}).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return contx.JSON(http.StatusNotFound, common.ResponseHTTP{
				Success: false,
				Message: "{{.Name}} not found",
				Data:    nil,
			})
		}
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: "Error retrieving {{.LowerName}}",
			Data:    nil,
		})
	}

	// Delete the {{.LowerName}}
	if err := db.Delete(&{{.LowerName}}).Error; err != nil {
		tx.Rollback()
		return contx.JSON(http.StatusInternalServerError, common.ResponseHTTP{
			Success: false,
			Message: "Error deleting {{.LowerName}}",
			Data:    nil,
		})
	}

	// Commit the transaction
	tx.Commit()

	// Return success respons
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "{{.Name}} deleted successfully.",
		Data:    {{.LowerName}},
	})
}

`
