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

func CurdFrameEcho() {

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
	"context"
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
	// Starting tracer context and tracer
	ctx := context.Background()
	tracer, span := observe.AppSpanner(ctx, "Get{{.Name}}s-root")
	defer span.End()

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
	result, err := common.PaginationPureModel(db, models.{{.Name}}{}, []models.{{.Name}}{}, uint(Page), uint(Limit), tracer)
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
	// Starting tracer context and tracer
	ctx := context.Background()
	tracer, span := observe.AppSpanner(ctx, "Get{{.Name}}ByID-root")
	defer span.End()


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
	if res := db.WithContext(tracer).Model(&models.{{.Name}}{}).Preload(clause.Associations).Where("id = ?", id).First(&{{.LowerName}}s); res.Error != nil {
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
// @Tags {{.Name}}s
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param {{.LowerName}} body models.{{.Name}}Post true "Add {{.Name}}"
// @Success 200 {object} common.ResponseHTTP{data=models.{{.Name}}Post}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /{{.LowerName}}s [post]
func Post{{.Name}}(contx echo.Context) error {
	// Starting tracer context and tracer
	ctx := context.Background()
	tracer, span := observe.AppSpanner(ctx, "Post{{.Name}}-root")
	defer span.End()

	// Database connection
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
	tx := db.WithContext(tracer).Begin()

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
// @Tags {{.Name}}s
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param {{.LowerName}} body models.{{.Name}}Post true "Patch {{.Name}}"
// @Param id path int true "{{.Name}} ID"
// @Success 200 {object} common.ResponseHTTP{data=models.{{.Name}}Post}
// @Failure 400 {object} common.ResponseHTTP{}
// @Failure 500 {object} common.ResponseHTTP{}
// @Router /{{.LowerName}}/{{ "{" }}{{.LowerName}}_id{{ "}" }} [patch]
func Patch{{.Name}}(contx echo.Context) error {
	// // Starting tracer context and tracer
	ctx := context.Background()
	tracer, span := observe.AppSpanner(ctx, "Patch{{.Name}}-root")
	defer span.End()

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
	tx := db.WithContext(tracer).Begin()

	// Check if the record exists
	if err := db.WithContext(tracer).First(&{{.LowerName}}, {{.LowerName}}.ID).Error; err != nil {
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
	if err := db.WithContext(tracer).Model(&{{.LowerName}}).UpdateColumns(*patch_{{.LowerName}}).Error; err != nil {
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
// @Tags {{.Name}}s
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "{{.Name}} ID"
// @Success 200 {object} common.ResponseHTTP{}
// @Failure 404 {object} common.ResponseHTTP{}
// @Failure 503 {object} common.ResponseHTTP{}
// @Router /{{.LowerName}}/{{ "{" }}{{.LowerName}}_id{{ "}" }} [delete]
func Delete{{.Name}}(contx echo.Context) error {
	// Starting tracer context and tracer
	ctx := context.Background()
	tracer, span := observe.AppSpanner(ctx, "Delete{{.Name}}-root")
	defer span.End()

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
	tx := db.WithContext(tracer).Begin()

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
func Add{{.FieldName}}{{.ParentName}}s(contx echo.Context) error {
	// Starting tracer context and tracer
	ctx := context.Background()
	tracer, span := observe.AppSpanner(ctx, "Add{{.FieldName}}{{.ParentName}}s-root")
	defer span.End()	

	//  database connection
	db := database.ReturnSession()

	// validate path params
	{{.LowerFieldName}}_id, err := strconv.Atoi(contx.Param("{{.LowerFieldName}}_id"))
	if err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// validate path params
	{{.LowerParentName}}_id, err := strconv.Atoi(contx.Param("{{.LowerParentName}}_id"))
	if err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// fetching {{.LowerParentName}} to be added
	var {{.LowerParentName}} models.{{.ParentName}}
	{{.LowerParentName}}.ID = uint({{.LowerParentName}}_id)
	if res := db.Find(&{{.LowerParentName}}); res.Error != nil {
		return contx.JSON(http.StatusServiceUnavailable, common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	//  {{.LowerParentName}}ending assocation
	var {{.LowerFieldName}} models.{{.FieldName}}
	{{.LowerFieldName}}.ID = uint({{.LowerFieldName}}_id)
	if err := db.Find(&{{.LowerFieldName}}); err.Error != nil {
		return contx.JSON(http.StatusNotFound, common.ResponseHTTP{
			Success: false,
			Message: "Record not Found",
			Data:    err.Error.Error(),
		})
	}

	tx := db.WithContext(tracer).Begin()
	if err := db.WithContext(tracer).Model(&{{.LowerFieldName}}).Association("{{.ParentName}}s").Append(&{{.LowerParentName}}); err != nil {
		tx.Rollback()
		return contx.JSON(http.StatusNotFound, common.ResponseHTTP{
			Success: false,
			Message: "{{.ParentName}}ending {{.ParentName}} Failed",
			Data:    err.Error(),
		})
	}
	tx.Commit()

	// return value if transaction is sucessfull
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
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
func Delete{{.FieldName}}{{.ParentName}}s(contx echo.Context) error {
	// Starting tracer context and tracer
	ctx := context.Background()
	tracer, span := observe.AppSpanner(ctx, "Delete{{.FieldName}}{{.ParentName}}s-root")
	defer span.End()

	//Connect to Database   
	db := database.ReturnSession()
	
	// validate path params
	{{.LowerFieldName}}_id, err := strconv.Atoi(contx.Param("{{.LowerFieldName}}_id"))
	if err != nil || {{.LowerFieldName}}_id == 0 {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	{{.LowerParentName}}_id, err := strconv.Atoi(contx.Param("{{.LowerParentName}}_id"))
	if err != nil || {{.LowerParentName}}_id == 0 {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}
	// fetching {{.LowerParentName}} to be deleted
	var {{.LowerParentName}} models.{{.ParentName}}
	{{.LowerParentName}}.ID = uint({{.LowerParentName}}_id)
	if res := db.Find(&{{.LowerParentName}}); res.Error != nil {
		return contx.JSON(http.StatusServiceUnavailable, common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// fettchng {{.LowerFieldName}}
	var {{.LowerFieldName}} models.{{.FieldName}}
	{{.LowerFieldName}}.ID = uint({{.LowerFieldName}}_id)
	if err := db.Find(&{{.LowerFieldName}}); err.Error != nil {
		return contx.JSON(http.StatusNotFound, common.ResponseHTTP{
			Success: false,
			Message: "Record not Found",
			Data:    err.Error.Error(),
		})
	}

	// removing {{.LowerParentName}}
	tx := db.WithContext(tracer).Begin()
	if err := db.WithContext(tracer).Model(&{{.LowerFieldName}}).Association("{{.ParentName}}s").Delete(&{{.LowerParentName}}); err != nil {
		tx.Rollback()
		return contx.JSON(http.StatusNonAuthoritativeInfo, common.ResponseHTTP{
			Success: false,
			Message: "Please Try Again Something Unexpected H{{.LowerParentName}}ened",
			Data:    err.Error(),
		})
	}

	tx.Commit()

	// return value if transaction is sucessfull
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
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
// @Param {{.LowerFieldName}}_id path int true "{{.ParentName}} ID"
// @Param {{.LowerParentName}}_id query int true "{{.FieldName}} ID"
// @Failure 400 {object} common.ResponseHTTP{}
// @Router /{{.LowerParentName}}{{.LowerFieldName}}/{{ "{" }}{{.LowerFieldName}}_id{{ "}" }} [patch]
func Add{{.FieldName}}{{.ParentName}}s(contx echo.Context) error {
	// Starting tracer context and tracer
	ctx := context.Background()
	tracer, span := observe.AppSpanner(ctx, "Add{{.FieldName}}{{.ParentName}}s-root")
	defer span.End()	

	//  database connection
	db := database.ReturnSession()

	// validate path params
	{{.LowerFieldName}}_id, err := strconv.Atoi(contx.Param("{{.LowerFieldName}}_id"))
	if err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// fetching Endpionts
	var {{.LowerFieldName}} models.{{.FieldName}}
	if res := db.WithContext(tracer).Model(&models.{{.FieldName}}{}).Where("id = ?", {{.LowerFieldName}}_id).First(&{{.LowerFieldName}}); res.Error != nil {
		return contx.JSON(http.StatusServiceUnavailable, common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// fetching {{.LowerFieldName}} to be added
	{{.LowerParentName}}_id, _ := strconv.Atoi(contx.QueryParam("{{.LowerParentName}}_id"))
	var {{.LowerParentName}} models.{{.ParentName}}
	if res := db.WithContext(tracer).Model(&models.{{.ParentName}}{}).Where("id = ?", {{.LowerParentName}}_id).First(&{{.LowerParentName}}); res.Error != nil {
		return contx.JSON(http.StatusServiceUnavailable, common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// startng update transaction

	tx := db.WithContext(tracer).Begin()
	//  Adding one to many Relation
	if err := db.WithContext(tracer).Model(&{{.LowerParentName}}).Association("{{.FieldName}}s").Append(&{{.LowerFieldName}}); err != nil {
		tx.Rollback()
		return contx.JSON(http.StatusNotFound, common.ResponseHTTP{
			Success: false,
			Message: "Error Adding Record",
			Data:    err.Error(),
		})
	}
	tx.Commit()

	// return value if transaction is sucessfull
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
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
// @Router /{{.LowerParentName}}{{.LowerFieldName}}/{{ "{" }}{{.LowerFieldName}}_id{{ "}" }} [delete]
func Delete{{.FieldName}}{{.ParentName}}s(contx echo.Context) error {
	// Starting tracer context and tracer
	ctx := context.Background()
	tracer, span := observe.AppSpanner(ctx, "Delete{{.FieldName}}{{.ParentName}}s-root")
	defer span.End()

	// Database Connection
	db := database.ReturnSession()

	// validate path params
	{{.LowerFieldName}}_id, err := strconv.Atoi(contx.Param("{{.LowerFieldName}}_id"))
	if err != nil {
		return contx.JSON(http.StatusBadRequest, common.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// Getting {{.FieldName}}
	var {{.LowerFieldName}} models.{{.FieldName}}
	if res := db.WithContext(tracer).Model(&models.{{.FieldName}}{}).Where("id = ?", {{.LowerFieldName}}_id).First(&{{.LowerFieldName}}); res.Error != nil {
		return contx.JSON(http.StatusServiceUnavailable, common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// fetching {{.LowerParentName}} to be added
	var {{.LowerParentName}} models.{{.ParentName}}
	{{.LowerParentName}}_id, _ := strconv.Atoi(contx.QueryParam("{{.LowerParentName}}_id"))
	if res := db.WithContext(tracer).Model(&models.{{.ParentName}}{}).Where("id = ?", {{.LowerParentName}}_id).First(&{{.LowerParentName}}); res.Error != nil {
		return contx.JSON(http.StatusServiceUnavailable, common.ResponseHTTP{
			Success: false,
			Message: res.Error.Error(),
			Data:    nil,
		})
	}

	// Removing {{.FieldName}} From {{.ParentName}}
	tx := db.WithContext(tracer).Begin()
	if err := db.WithContext(tracer).Model(&{{.LowerParentName}}).Association("{{.FieldName}}s").Delete(&{{.LowerFieldName}}); err != nil {
		tx.Rollback()
		return contx.JSON(http.StatusNotFound, common.ResponseHTTP{
			Success: false,
			Message: "Record not Found",
			Data:    err.Error(),
		})
	}
	tx.Commit()

	// return value if transaction is sucessfull
	return contx.JSON(http.StatusOK, common.ResponseHTTP{
		Success: true,
		Message: "Success Deleteing a {{.FieldName}} From {{.ParentName}}.",
		Data:    {{.LowerParentName}},
	})
}


{{ end}}
{{ end}}


`
