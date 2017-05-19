package main

import (
	"encoding/json"
	"fmt"

	"github.com/nsip/nias2/naprr"
	"github.com/nsip/nias2/xml"
	"github.com/playlyfe/go-graphql"
)

func main() {

	sr := naprr.NewStreamReader()
	sd := sr.GetSchoolData("21212")

	schema := `
	type TestDisruption {
		Event: String
	}

	type NAPEvent {
		EventID: String
		SPRefID: String
		PSI: String
		SchoolRefID: String
		SchoolID: String
		TestID: String
		NAPTestLocalID: String
		TestDisruptionList: [TestDisruption]
	}

	type School {
		events: [NAPEvent]
	}
	`
	// event := &xml.NAPEvent{
	// 	EventID:     "event1",
	// 	SPRefID:     "student1",
	// 	PSI:         "psi1",
	// 	SchoolRefID: "school_refid1",
	// 	SchoolID:    "123456",
	// }

	// fmt.Printf("\n\nTestDisruptionList: %#v\n\n", event.TestDisruptionList)
	// fmt.Printf("\n\nTestDisruption: %#v\n\n", event.TestDisruptionList.TestDisruption)
	// fmt.Printf("\n\nEvent: %#v\n\n", event)

	// tdl := struct {
	// 	TestDisruption []struct {
	// 		Event string
	// 	}
	// }{
	// 	// {
	// 	[]struct{}{{Event: "e1"}},
	// 	// },
	// }

	resolvers := map[string]interface{}{}

	resolvers["NAPEvent/TestDisruptionList"] = func(params *graphql.ResolveParams) (interface{}, error) {
		disruptionEvents := []interface{}{}
		if napEvent, ok := params.Source.(xml.NAPEvent); ok {
			for _, event := range napEvent.TestDisruptionList.TestDisruption {
				disruptionEvents = append(disruptionEvents, event)
			}
		}
		return disruptionEvents, nil
	}

	// resolvers["TestDisruption"] = func(params *graphql.ResolveParams) (interface{}, error) {

	// 	fmt.Println("\n\n-- TD got Called --\n\n")

	// 	if napEvent, ok := params.Source.(TestDisruption); ok {
	// 		return napEvent, nil
	// 	} else {
	// 		return "empty:var", nil
	// 	}
	// }

	// resolvers["NAPEvent/PSI"] = func(params *graphql.ResolveParams) (interface{}, error) {

	// 	if napEvent, ok := params.Source.(*xml.NAPEvent); ok {
	// 		return napEvent.PSI + ":source", nil
	// 	} else {
	// 		return event.PSI + ":var", nil
	// 	}
	// }
	// resolvers["NAPEvent/SchoolRefID"] = func(params *graphql.ResolveParams) (interface{}, error) {

	// 	if napEvent, ok := params.Source.(*xml.NAPEvent); ok {
	// 		return napEvent.SchoolRefID + ":source", nil
	// 	} else {
	// 		return event.SchoolRefID + ":var", nil
	// 	}
	// }
	// resolvers["NAPEvent/SchoolID"] = func(params *graphql.ResolveParams) (interface{}, error) {

	// 	if napEvent, ok := params.Source.(*xml.NAPEvent); ok {
	// 		return napEvent.SchoolID + ":source", nil
	// 	} else {
	// 		return event.SchoolID + ":var", nil
	// 	}
	// }

	resolvers["School/events"] = func(params *graphql.ResolveParams) (interface{}, error) {

		events := []interface{}{}
		for _, event := range sd.Events {
			// if event.TestDisruptionList.TestDisruption != nil {
			events = append(events, event)
			// }
		}
		return events, nil

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
			#SPRefID
			PSI
			SchoolRefID
			SchoolID
			TestID
			NAPTestLocalID
			TestDisruptionList {
				Event
			}
		}
	}`
	result, err := executor.Execute(context, query, variables, "")
	if err != nil {
		panic(err)
	}

	jsonResult, err := json.Marshal(result)

	fmt.Println(string(jsonResult), "\n", err)

}
