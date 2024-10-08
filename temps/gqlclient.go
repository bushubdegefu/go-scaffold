package temps

import (
	"fmt"
	"os"
	"text/template"
)

func GQLClientFrame() {
	//  this is creating manger file inside the manager folder
	// ############################################################
	devf_tmpl, err := template.New("RenderData").Parse(gqlClientTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("manager", os.ModePerm)
	if err != nil {
		panic(err)
	}

	devf_file, err := os.Create("store.js")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer devf_file.Close()

	err = devf_tmpl.Execute(devf_file, RenderData)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

}

var gqlClientTemplate = `
import axios from "axios";
import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";

// https://goblue-back.onrender.com/api/v1
export const btmClient = axios.create({
  baseURL: "/api/v1",
  timeout: 10000,
});
const postURL = "/admin";

{{ range .Models}}
//#######################################################
//  graph {{.LowerName}} store and requests
//#######################################################
export const use{{.Name}}SchemaStore = create((set, get) => ({
  {{.LowerName}}s: [],
  {{.LowerName}}: null,{{ range .Relations }}
 {{.LowerParentName}}{{.LowerFieldName}}s: [],{{end}}
  page: 1,
  size: 15,
  get{{.LowerName}}s: async () => {
  	const pdata ={
   		query: {{.BackTick}} query { {{.LowerName}}s(page: Int!, size: Int!) {
            {{range .Fields}} {{.LowerName}}
            {{end}}
        }{{.BackTick}},
        variables: {
        	page: get().page,
        	size: get().size
        }
    }

    await btmClient
      .request({
        method: "POST",
        url: postURL,
        headers: {
          "Content-Type": "application/json",
          "X-APP-TOKEN": token,
        },
        data: pdata,
      })
      .then(function (response) {
        console.log(response.data);
        set((state) => ({
          ...state,
          {{.LowerName}}s: response?.data?.{{.LowerName}}s,
        }));
      })
      .catch((response, error) => {
        const responseError = response?.data?.details;
        console.log(responseError);
      });
  },
  get{{.LowerName}}: async (id) => {
    const pdata={
    	query: {{.BackTick}} query {  {{.LowerName}}(id: Int!) {
          {{range .Fields}} {{.LowerName}}
          {{end}}}
        }{{.BackTick}},
        variables: {
        id: id
      }
    }
    await btmClient
      .request({
        method: "POST",
        url: postURL,
        headers: {
          "Content-Type": "application/json",
          "X-APP-TOKEN": token,
        },
        data: pdata,
      })
      .then(function (response) {
        console.log(response.data);
        set((state) => ({
          ...state,
          {{.LowerName}}: response?.data?.{{.LowerName}},
        }));
      })
      .catch((response, error) => {
        const responseError = response?.data?.details;
        console.log(responseError);
      });
  },
  create{{.LowerName}}: async (data) => {
    const pdata = {
	    mutation: {{.BackTick}} mutation {
	        create{{.LowerName}}(input: Create{{.Name}}Input) {
	        {{range .Fields}} {{.LowerName}}
	        {{end}}
	        }}
	    {{.BackTick}},
	    variables: { input: {
		 {{range .Fields}} {{if .Post}} {{.LowerName}}: data.{{.LowerName}},
			{{end}}{{end}}}
			}
    }

    await btmClient
      .request({
        method: "POST",
        url: postURL,
        headers: {
          "Content-Type": "application/json",
          "X-APP-TOKEN": token,
        },
        data: pdata,
      })
      .then(function (response) {
        console.log(response.data);
        set((state) => ({
          ...state,
          {{.LowerName}}: response?.data?.{{.LowerName}},
        }));
      })
      .catch((response, error) => {
        const responseError = response?.data?.details;
        console.log(responseError);
      });
  },
  update{{.LowerName}}: async (data) => {
    const pdata = {
    mutation: {{.BackTick}} mutation {
          update{{.LowerName}}(input: Update{{.Name}}Input) {
          {{range .Fields}} {{.LowerName}}
          {{end}}
          }}
      {{.BackTick}},
      variables: { input: {
		 {{range .Fields}} {{if .Put}} {{.LowerName}}: data.{{.LowerName}},
			{{end}}{{end}}}
			}
    }

    await btmClient
      .request({
        method: "POST",
        url: postURL,
        headers: {
          "Content-Type": "application/json",
          "X-APP-TOKEN": token,
        },
        data: pdata,
      })
      .then(function (response) {
        console.log(response.data);
      })
      .catch((response, error) => {
        const responseError = response?.data?.details;
        console.log(responseError);
      });
  },
  delete{{.LowerName}}: async (id) => {
    const pdata = {
    	mutation: {{.BackTick}} mutation {
          delete{{.LowerName}}(id: Int!)
        }{{.BackTick}},
      	variables: {
	       	id: id
	       }
    }

    await btmClient
      .request({
        method: "POST",
        url: postURL,
        headers: {
          "Content-Type": "application/json",
          "X-APP-TOKEN": token,
        },
        data: pdata,
      })
      .then(function (response) {
        console.log(response.data);
      })
      .catch((response, error) => {
        const responseError = response?.data?.details;
        console.log(responseError);
      });
  },
  // ######################################
  // relation OTM/MTM
  {{ range .Relations }}
  get{{.LowerParentName}}{{.LowerFieldName}}s: async ({{.LowerParentName}}Id, {{.LowerFieldName}}Id, page, size) => {
      const pdata = {
      	query: {{.BackTick}} query {
            {{.LowerParentName}}{{.LowerFieldName}}s( {{.LowerParentName}}_id: Int!,  {{.LowerFieldName}}_id: Int!, page: Int!, size: Int!) {
           	{{range .ParentFields}} {{.LowerName}}
            {{end}}}
            }
          }{{.BackTick}},
        variables: {
        	{{.LowerParentName}}_id: {{.LowerParentName}}Id,
        	{{.LowerFieldName}}_id: {{.LowerFieldName}}Id,
         	page: page,
          	size: size,
           }
        }

      await btmClient
        .request({
          method: "POST",
          url: postURL,
          headers: {
            "Content-Type": "application/json",
            "X-APP-TOKEN": token,
          },
          data: pdata,
        })
        .then(function (response) {
          console.log(response.data);
          set((state) => ({
            ...state,
            {{.LowerParentName}}{{.LowerFieldName}}s: response?.data?.{{.LowerFieldName}},
          }));
        })
        .catch((response, error) => {
          const responseError = response?.data?.details;
          console.log(responseError);
        });
    },
  create{{.LowerParentName}}{{.LowerFieldName}}s: async ({{.LowerParentName}}Id, {{.LowerFieldName}}Id) => {
	const pdata = {
		mutation: {{.BackTick}} mutation {
	        create{{.LowerParentName}}{{.LowerFieldName}}({{.LowerParentName}}_id: Int!, {{.LowerFieldName}}_id: Int! ) {
	        }}
	    {{.BackTick}},
		variables: {
			{{.LowerParentName}}_id: {{.LowerParentName}}Id,
		 	{{.LowerFieldName}}_id: {{.LowerFieldName}}Id,
			}
		}

	    await btmClient
	    .request({
	        method: "POST",
	        url: postURL,
	        headers: {
	        "Content-Type": "application/json",
	        "X-APP-TOKEN": token,
	        },
	        data: pdata,
	    })
	    .then(function (response) {
	        console.log(response.data);
	    })
	    .catch((response, error) => {
	        const responseError = response?.data;
	        console.log(responseError);
    });
    },
    delete{{.LowerParentName}}{{.LowerFieldName}}s: async ({{.LowerParentName}}Id, {{.LowerFieldName}}Id) => {
    const pdata = {
		mutation: {{.BackTick}} mutation {
	        delete{{.LowerParentName}}{{.LowerFieldName}}({{.LowerParentName}}_id: Int!, {{.LowerFieldName}}_id: Int! ) {
	        }}
	    {{.BackTick}},
		variables: {
			{{.LowerParentName}}_id: {{.LowerParentName}}Id,
		 	{{.LowerFieldName}}_id: {{.LowerFieldName}}Id,
			}
		}

      await btmClient
        .request({
          method: "POST",
          url: postURL,
          headers: {
            "Content-Type": "application/json",
            "X-APP-TOKEN": token,
          },
          data: pdata,
        })
        .then(function (response) {
          console.log(response.data);
        })
        .catch((response, error) => {
          const responseError = response?.data?.details;
          console.log(responseError);
        });
    },
  {{end}}

}));{{end}}
`
