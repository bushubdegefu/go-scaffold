package temps

import (
	"fmt"
	"os"
	"text/template"
)

func GQLServerClientFrame() {
	//  this is creating manger file inside the manager folder
	// ############################################################
	devf_tmpl, err := template.New("RenderData").Parse(gqlServerClientTemplate)
	if err != nil {
		panic(err)
	}

	// Create the models directory if it does not exist
	err = os.MkdirAll("manager", os.ModePerm)
	if err != nil {
		panic(err)
	}

	devf_file, err := os.Create("nextstore.js")
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

var gqlServerClientTemplate = `
import axios from "axios";
import { jwtDecode } from "jwt-decode";
import { cookies } from "next/headers";


const postURL = "project";
const BaseURL = "http://localhost:8500/api/v1";


{{ range .Models}}
//#######################################################
//  graph {{.LowerName}}  requests
//#######################################################

  export async function get_{{.LowerName}}s(){
	  const cookieStore = cookies();
	  const token = cookieStore.get("access_token")?.value;
	  const current_project = JSON.parse(cookieStore.get("current_project")?.value);
	  const page_size = cookieStore.get("page_size")?.value;
	  const current_page = cookieStore.get("current_page")?.value;

	  	const pdata ={
	   		query: {{.BackTick}} query { {{.LowerName}}s($page: Int!, $size: Int!){
	     		{{.LowerName}}s($page: $page, $size: $size) {
	            {{range .Fields}} {{.LowerName}}
	            {{end}}
	            }
	        }{{.BackTick}},
	        variables: {
	        	page:current_page,
	        	size: page_size
	        }
	    }

	    try {
	      const response = await axios.post(
	        {{.BackTick}}${BaseURL}/${postURL}/${current_project?.id}{{.BackTick}},
	        pdata, // Directly passing the data
	        {
	          headers: {
	            "Content-Type": "application/json",
	            "X-APP-TOKEN": token,
	          },
	        },
	      );

	      // If the response is successful, return the {{.LowerName}}s data
	      return response?.data?.data?.{{.LowerName}}s || [];
	    } catch (error) {
	      // Log error details for debugging
	      console.error("Error fetching {{.LowerName}}s:", error);

	      // Check if error response contains specific details and return them
	      if (error?.response?.data) {
	        return [error];
	      }

	      // If no specific error details, return a generic message
	      return ["Unkown error happened",];
	  }
}

export  async function get_{{.LowerName}}({{.LowerName}}_id){
	const cookieStore = cookies();
  	const token = cookieStore.get("access_token")?.value;
   	const current_project = JSON.parse(cookieStore.get("current_project")?.value);

    const pdata={
    	query: {{.BackTick}} query { {{.LowerName}}($id: Int!)  {
     		{{.LowerName}}($id: $id) {
          		{{range .Fields}} {{.LowerName}}
            	{{end}}}
          }
        }{{.BackTick}},
        variables: {
        id: {{.LowerName}}
      }
    }
    try {
      const response = await axios.post(
        {{.BackTick}}${BaseURL}/${postURL}/${current_project?.id}{{.BackTick}},
        pdata, // Directly passing the data
        {
          headers: {
            "Content-Type": "application/json",
            "X-APP-TOKEN": token,
          },
        },
      );

      // If the response is successful, return the {{.LowerName}}s data
      return response?.data?.data?.{{.LowerName}}s || [];
    } catch (error) {
      // Log error details for debugging
      console.error("Error fetching {{.LowerName}}:", error);

      // Check if error response contains specific details and return them
      if (error?.response?.data) {
        return [error];
      }

      // If no specific error details, return a generic message
      return ["Unkown error happened",];
  }
}

export async function create_{{.LowerName}}(data){
	 const cookieStore = cookies();
	 const token = cookieStore.get("access_token")?.value;
	 const current_project = JSON.parse(cookieStore.get("current_project")?.value);

    const pdata = {
	    query: {{.BackTick}} mutation create{{.LowerName}}($input: Create{{.Name}}Input) {
	        create{{.LowerName}}($input: $input) {
	        {{range .Fields}} {{.LowerName}}
	        {{end}}
	        }
		}{{.BackTick}},
	    variables: { input: {
		 {{range .Fields}} {{if .Post}} {{.LowerName}}: data.{{.LowerName}},
			{{end}}{{end}}}
			}
    }

    try {
      const response = await axios.post(
        {{.BackTick}}${BaseURL}/${postURL}/${current_project?.id}{{.BackTick}},
        pdata, // Directly passing the data
        {
          headers: {
            "Content-Type": "application/json",
            "X-APP-TOKEN": token,
          },
        },
      );

      // If the response is successful, return the {{.LowerName}} data
      return response?.data?.data?.{{.LowerName}} || [];
    } catch (error) {
      // Log error details for debugging
      console.error("Error Creating {{.LowerName}}:", error);

      // Check if error response contains specific details and return them
      if (error?.response?.data) {
        return [error,];
      }

      // If no specific error details, return a generic message
      return ["Unkown error happened",];
  }
}

  export async function update_{{.LowerName}}(data){
	  const cookieStore = cookies();
	  const token = cookieStore.get("access_token")?.value;
	  const current_project = JSON.parse(cookieStore.get("current_project")?.value);

    const pdata = {
    	query: {{.BackTick}} mutation { update{{.LowerName}}($input: Update{{.Name}}Input) {
     		update{{.LowerName}}($input: $input)
        		{{range .Fields}} {{.LowerName}}
          	{{end}}
          }
        }{{.BackTick}},
      variables: { input: {
		 {{range .Fields}} {{if .Put}} {{.LowerName}}: data.{{.LowerName}},
			{{end}}{{end}}}
			}
    }

    try {
      const response = await axios.post(
        {{.BackTick}}${BaseURL}/${postURL}/${current_project?.id}{{.BackTick}},
        pdata, // Directly passing the data
        {
          headers: {
            "Content-Type": "application/json",
            "X-APP-TOKEN": token,
          },
        },
      );

      // If the response is successful, return the {{.LowerName}} data
      return response?.data?.data?.{{.LowerName}}s || [];
    } catch (error) {
      // Log error details for debugging
      console.error("Error updating {{.LowerName}}:", error);

      // Check if error response contains specific details and return them
      if (error?.response?.data) {
        return [error];
      }

      // If no specific error details, return a generic message
      return ["Unkown error happened",];
  }
}

export async function delete_{{.LowerName}}(id){
	 const cookieStore = cookies();
	 const token = cookieStore.get("access_token")?.value;
	 const current_project = JSON.parse(cookieStore.get("current_project")?.value);

    const pdata = {
    	query: {{.BackTick}} mutation delete{{.LowerName}}($id: Int!) {
          delete{{.LowerName}}($id: $id)
        }{{.BackTick}},
      	variables: {
	       	id: id
	       }
    }

    try {
      const response = await axios.post(
        {{.BackTick}}${BaseURL}/${postURL}/${current_project?.id}{{.BackTick}},
        pdata, // Directly passing the data
        {
          headers: {
            "Content-Type": "application/json",
            "X-APP-TOKEN": token,
          },
        },
      );

      // If the response is successful, return the {{.LowerName}} data
      return response?.data?.data?.{{.LowerName}} || [];
    } catch (error) {
      // Log error details for debugging
      console.error("Error Deleteing {{.LowerName}}:", error);

      // Check if error response contains specific details and return them
      if (error?.response?.data) {
        return [error];
      }

      // If no specific error details, return a generic message
      return ["Unkown error happened",];
  }
}

  // ######################################
  // relation OTM/MTM
  // ######################################
  {{ range .Relations }}
export async function get_{{.LowerParentName}}{{.LowerFieldName}}s({{.LowerParentName}}Id, {{.LowerFieldName}}Id, page, size){
	const cookieStore = cookies();
  	const token = cookieStore.get("access_token")?.value;
   	const current_project = JSON.parse(cookieStore.get("current_project")?.value);
    const page_size = cookieStore.get("page_size")?.value;
    const current_page = cookieStore.get("current_page")?.value;

      const pdata = {
      	query: {{.BackTick}} query {{.LowerParentName}}{{.LowerFieldName}}s(${{.LowerParentName}}_id: Int!, $page: Int!, $size: Int!) {
            {{.LowerParentName}}{{.LowerFieldName}}s( ${{.LowerParentName}}_id: ${{.LowerParentName}}_id, $page: $page, $size: $size) {
           	{{range .ParentFields}} {{.LowerName}}
            {{end}}}
            }
          }{{.BackTick}},
        variables: {
        	{{.LowerParentName}}_id: {{.LowerParentName}}Id,
         	page: current_page,
          	size: page_size,
           }
        }

        try {
          const response = await axios.post(
            {{.BackTick}}${BaseURL}/${postURL}/${current_project?.id}{{.BackTick}},
            pdata, // Directly passing the data
            {
              headers: {
                "Content-Type": "application/json",
                "X-APP-TOKEN": token,
              },
            },
          );

          // If the response is successful, return the {{.LowerFieldName}}s data
          return response?.data?.data?.{{.LowerFieldName}}s || [];
        } catch (error) {
          // Log error details for debugging
          console.error("Error fetching {{.LowerFieldName}}s:", error);

          // Check if error response contains specific details and return them
          if (error?.response?.data) {
            return [error];
          }

          // If no specific error details, return a generic message
          return ["Unkown error happened",];
    }
}

export async function create_{{.LowerParentName}}{{.LowerFieldName}}s({{.LowerParentName}}Id, {{.LowerFieldName}}Id){
	const cookieStore = cookies();
  	const token = cookieStore.get("access_token")?.value;
   	const current_project = JSON.parse(cookieStore.get("current_project")?.value);

	const pdata = {
		query: {{.BackTick}} mutation create{{.LowerParentName}}{{.LowerFieldName}}(${{.LowerParentName}}_id: Int!, ${{.LowerFieldName}}_id: Int! ) {
	        create{{.LowerParentName}}{{.LowerFieldName}}(${{.LowerParentName}}_id: ${{.LowerParentName}}_id, ${{.LowerFieldName}}_id: ${{.LowerFieldName}}_id)
		}{{.BackTick}},
		variables: {
			{{.LowerParentName}}_id: {{.LowerParentName}}Id,
		 	{{.LowerFieldName}}_id: {{.LowerFieldName}}Id,
			}
		}

	  try {
	      const response = await axios.post(
	        {{.BackTick}}${BaseURL}/${postURL}/${current_project?.id}{{.BackTick}},
	        pdata, // Directly passing the data
	        {
	          headers: {
	            "Content-Type": "application/json",
	            "X-APP-TOKEN": token,
	          },
	        },
	      );

	      // If the response is successful, return the added  {{.LowerFieldName}} on to {{.LowerParentName}} data
	      return response?.data?.data?.{{.LowerFieldName}} || [];
	    } catch (error) {
	      // Log error details for debugging
	      console.error("Error adding {{.LowerParentName}} to  {{.LowerFieldName}}:", error);

	      // Check if error response contains specific details and return them
	      if (error?.response?.data) {
	        return [error];
	      }

	      // If no specific error details, return a generic message
	      return ["Unkown error happened",];
	    }
}

export   async function delete_{{.LowerParentName}}{{.LowerFieldName}}s({{.LowerParentName}}Id, {{.LowerFieldName}}Id){
	const cookieStore = cookies();
 	const token = cookieStore.get("access_token")?.value;
  	const current_project = JSON.parse(cookieStore.get("current_project")?.value);

    const pdata = {
		query: {{.BackTick}} mutation {
	        delete{{.LowerParentName}}{{.LowerFieldName}}(${{.LowerParentName}}_id: Int!, ${{.LowerFieldName}}_id: Int! ) {
				delete{{.LowerParentName}}{{.LowerFieldName}}(${{.LowerParentName}}_id: ${{.LowerParentName}}_id, ${{.LowerFieldName}}_id: ${{.LowerParentName}}_id)
		} {{.BackTick}},
		variables: {
			{{.LowerParentName}}_id: {{.LowerParentName}}Id,
		 	{{.LowerFieldName}}_id: {{.LowerFieldName}}Id,
			}
		}

		try {
		      const response = await axios.post(
		        {{.BackTick}}${BaseURL}/${postURL}/${current_project?.id}{{.BackTick}},
		        pdata, // Directly passing the data
		        {
		          headers: {
		            "Content-Type": "application/json",
		            "X-APP-TOKEN": token,
		          },
		        },
		      );

		      // If the response is successful, return the deleted {{.LowerFieldName}} from {{.LowerParentName}} data
		      return response?.data?.data?.{{.LowerFieldName}} || [];
	    } catch (error) {
	      // Log error details for debugging
	      console.error("Error Deleteing {{.LowerFieldName}}  from {{.LowerParentName}}:", error);

	      // Check if error response contains specific details and return them
	      if (error?.response?.data) {
	        return [error];
	      }

	      // If no specific error details, return a generic message
	      return ["Unkown error happened",];
	    }
    }
  {{end}}
{{end}}
`
