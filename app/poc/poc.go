package main

import (
	"fmt"

	"github.com/playlyfe/go-graphql"
)

func main() {
	schema := `
	interface Pet {
	    name: String
	}
	type Dog implements Pet {
	    name: String
	    woofs: Boolean
	}
	type Cat implements Pet {
	    name: String
	    meows: Boolean
	}
	type QueryRoot {
	    pets: [Pet]
	}
	`
	resolvers := map[string]interface{}{}
	resolvers["QueryRoot/pets"] = func(params *graphql.ResolveParams) (interface{}, error) {
		return []map[string]interface{}{
			{
				"__typename": "Dog",
				"name":       "Odie",
				"woofs":      true,
			},
			{
				"__typename": "Cat",
				"name":       "Garfield",
				"meows":      false,
			},
		}, nil
	}
	context := map[string]interface{}{}
	variables := map[string]interface{}{}
	executor, err := graphql.NewExecutor(schema, "QueryRoot", "", resolvers)
	executor.ResolveType = func(value interface{}) string {
		if object, ok := value.(map[string]interface{}); ok {
			return object["__typename"].(string)
		}
		return ""
	}
	query := `{
		pets {
			name
			... on Dog {
				woofs
			}
			... on Cat {
				meows
			}
		}
	}`
	result, err := executor.Execute(context, query, variables, "")
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n\n %v \n\n", result)
}
