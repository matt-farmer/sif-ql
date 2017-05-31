// qlserver.go
package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/nsip/nias2/naprr"
	"github.com/nsip/nias2/xml"
	"github.com/playlyfe/go-graphql"
)

var sr = naprr.NewStreamReader()
var executor *graphql.Executor

// var sd = sr.GetSchoolData("21212")

var schema = `
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
		SchoolSector: String
		System: String
		SchoolGeoLocation: String
		ReportingSchoolName: String
		JurisdictionID: String
		ParticipationCode: String
		ParticipationText: String
		Device: String
		Date: String
		StartTime: String
		LapsedTimeTest: String
		ExemptionReason: String
		PersonalDetailsChanged: String
		PossibleDuplicate: String
		DOBRange: String
		TestDisruptionList: [TestDisruption]
	}

	type School {
		events(onlyDisruptions: Boolean): [NAPEvent]
	}

	type SchoolQuery {
		getSchoolData(acaraID: String): School
	}
	`

//
// wrapper type to capture graphql input
//
type GQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

func buildResolvers() map[string]interface{} {

	resolvers := map[string]interface{}{}
	resolvers["NAPEvent/TestDisruptionList"] = func(params *graphql.ResolveParams) (interface{}, error) {
		disruptionEvents := []interface{}{}

		// log.Printf("params: %#v\n\n", params)

		if napEvent, ok := params.Source.(xml.NAPEvent); ok {
			for _, event := range napEvent.TestDisruptionList.TestDisruption {
				disruptionEvents = append(disruptionEvents, event)
			}
		}
		return disruptionEvents, nil
	}
	resolvers["School/events"] = func(params *graphql.ResolveParams) (interface{}, error) {

		events := []interface{}{}

		// log.Printf("params: %#v\n\n", params)

		if sd, ok := params.Source.(*naprr.SchoolData); ok {

			if params.Args["onlyDisruptions"].(bool) {
				for _, event := range sd.Events {
					if event.TestDisruptionList.TestDisruption != nil {
						events = append(events, event)
					}
				}
			} else {
				for _, event := range sd.Events {
					events = append(events, event)
				}
			}
		}
		return events, nil
	}
	resolvers["SchoolQuery/getSchoolData"] = func(params *graphql.ResolveParams) (interface{}, error) {
		// get the school data
		schoolID := params.Args["acaraID"].(string)
		sd := sr.GetSchoolData(schoolID)
		return sd, nil
	}

	return resolvers
}

func buildExecutor() *graphql.Executor {

	executor, err := graphql.NewExecutor(schema, "SchoolQuery", "", buildResolvers())
	if err != nil {
		log.Fatalln("Cannot create Executor: ", err)
	}

	executor.ResolveType = func(value interface{}) string {
		log.Printf("resolve: %#v\n\n", value)
		switch value.(type) {
		case *xml.NAPEvent:
			return "NAPEvent"
		}
		return ""
	}

	return executor
}

//
// simple web server to support gql queries & web ui (graphiql)
//
func graphQLHandler(c echo.Context) error {

	grq := new(GQLRequest)
	if err := c.Bind(grq); err != nil {
		return err
	}

	query := grq.Query
	variables := grq.Variables
	gqlContext := map[string]interface{}{}
	// log.Printf("variables: %v\n\n", variables)
	result, err := executor.Execute(gqlContext, query, variables, "")
	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, result)

}

func main() {

	executor = buildExecutor()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/graphql", graphQLHandler)

	e.Logger.Fatal(e.Start(":1329"))
}
