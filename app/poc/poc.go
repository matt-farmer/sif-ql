package main

import (
	"encoding/json"
	"fmt"

	"github.com/nsip/nias2/xml"
	"github.com/playlyfe/go-graphql"
)

func main() {

	schema := `
	type NAPEvent {
		EventID: String
		SPRefID: String
		PSI: String
		#SchoolRefID: String
		#SchoolID: String
	}
	type School {
		events: [NAPEvent]
	}
	`
	event := &xml.NAPEvent{
		EventID:     "event1",
		SPRefID:     "student1",
		PSI:         "psi1",
		SchoolRefID: "school_refid1",
		SchoolID:    "12345",
	}

	resolvers := map[string]interface{}{}
	resolvers["NAPEvent/EventID"] = func(params *graphql.ResolveParams) (interface{}, error) {

		if napEvent, ok := params.Source.(*xml.NAPEvent); ok {
			return napEvent.EventID + ":source", nil
		} else {
			return event.EventID + ":var", nil
		}
	}
	resolvers["NAPEvent/SPRefID"] = func(params *graphql.ResolveParams) (interface{}, error) {

		if napEvent, ok := params.Source.(*xml.NAPEvent); ok {
			return napEvent.SPRefID + ":source", nil
		} else {
			return event.SPRefID + ":var", nil
		}
	}
	resolvers["NAPEvent/PSI"] = func(params *graphql.ResolveParams) (interface{}, error) {

		if napEvent, ok := params.Source.(*xml.NAPEvent); ok {
			return napEvent.PSI + ":source", nil
		} else {
			return event.PSI + ":var", nil
		}
	}

	resolvers["School/events"] = func(params *graphql.ResolveParams) (interface{}, error) {

		return []interface{}{event, event}, nil

	}

	context := map[string]interface{}{}
	variables := map[string]interface{}{}
	executor, err := graphql.NewExecutor(schema, "School", "", resolvers)
	executor.ResolveType = func(value interface{}) string {
		switch value.(type) {
		case *xml.NAPEvent:
			return "NAPEvent"
		}
		return ""
	}
	query := `
	{
		events {
			EventID
			SPRefID
			PSI
		}
	}`
	result, err := executor.Execute(context, query, variables, "")
	if err != nil {
		panic(err)
	}

	jsonResult, err := json.Marshal(result)

	fmt.Println(string(jsonResult), "\n", err)

}
