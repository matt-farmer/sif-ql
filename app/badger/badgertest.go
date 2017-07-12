// badgertest.go

package main

import (
	"bytes"
	goxml "encoding/xml"
	"fmt"
	"log"

	"github.com/dgraph-io/badger"
	"github.com/nsip/nias2/naprr"
	"github.com/nsip/nias2/xml"
)

var ge = naprr.GobEncoder{}
var kv *badger.KV

func main() {

	kv = openDB()
	// defer kv.Close()

	// ingestResultsFile("master_nap.xml.zip")

	runIterator()

}

func runIterator() {
	fmt.Println("Calling iterator")

	itrOpt := badger.IteratorOptions{
		PrefetchSize: 1000,
		FetchValues:  true,
		Reverse:      false,
	}
	itr := kv.NewIterator(itrOpt)

	searchKey := []byte("StudentPersonal:")

	for itr.Seek(searchKey); itr.Valid(); itr.Next() {
		item := itr.Item()
		key := item.Key()
		val := item.Value()

		// search exact keys only, assumes key marker :
		if !bytes.Contains(key, searchKey) {
			break
		}

		fmt.Printf("%s - %+v\n", key, val)
	}

}

func openDB() *badger.KV {
	log.Println("opening database...")
	opt := badger.DefaultOptions
	opt.Dir = "/tmp/badger"
	opt.ValueDir = "/tmp/badger"
	kv, dbErr := badger.NewKV(&opt)
	if dbErr != nil {
		log.Fatalln("DB Create error: ", dbErr)
	}
	return kv
}

func ingestResultsFile(resultsFilePath string) {

	// open the data file for streaming read
	xmlFile, err := naprr.OpenResultsFile(resultsFilePath)
	if err != nil {
		log.Fatalln("unable to open results file")
	}

	log.Println("Reading data file...")

	entries := []*badger.Entry{}

	decoder := goxml.NewDecoder(xmlFile)
	totalTests := 0
	totalTestlets := 0
	totalTestItems := 0
	totalTestScoreSummarys := 0
	totalEvents := 0
	totalResponses := 0
	totalCodeFrames := 0
	totalSchools := 0
	totalStudents := 0
	var inElement string
	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		// Inspect the type of the token just read.
		switch se := t.(type) {
		case goxml.StartElement:
			// If we just read a StartElement token
			inElement = se.Name.Local
			// ...handle by type
			switch inElement {
			case "NAPTest":
				var t xml.NAPTest
				decoder.DecodeElement(&t, &se)
				gt, err := ge.Encode(t)
				if err != nil {
					log.Println("Unable to gob-encode nap test: ", err)
				}

				// {NAPTest} = object
				entries = badger.EntriesSet(entries, []byte(t.TestID), gt)

				// NAPTest-type:{id} = id
				key := []byte("NAPTest:" + t.TestID)
				entries = badger.EntriesSet(entries, key, []byte(t.TestID))

				totalTests++

			case "NAPTestlet":
				var tl xml.NAPTestlet
				decoder.DecodeElement(&tl, &se)
				gtl, err := ge.Encode(tl)
				if err != nil {
					log.Println("Unable to gob-encode nap testlet: ", err)
				}

				// {NAPTestlet} = object
				entries = badger.EntriesSet(entries, []byte(tl.TestletID), gtl)

				// NAPTestlet-type:{id} = {id}
				key := []byte("NAPTestlet:" + tl.TestletID)
				entries = badger.EntriesSet(entries, key, []byte(tl.TestletID))

				totalTestlets++

			case "NAPTestItem":
				var ti xml.NAPTestItem
				decoder.DecodeElement(&ti, &se)
				gti, err := ge.Encode(ti)
				if err != nil {
					log.Println("Unable to gob-encode nap test item: ", err)
				}

				// {NAPTestItem} = object
				entries = badger.EntriesSet(entries, []byte(ti.ItemID), gti)

				// NapTestItem-type:{id} = {id}
				key := []byte("NAPTestItem:" + ti.ItemID)
				entries = badger.EntriesSet(entries, key, []byte(ti.ItemID))

				totalTestItems++

			case "NAPTestScoreSummary":
				var tss xml.NAPTestScoreSummary
				decoder.DecodeElement(&tss, &se)
				gtss, err := ge.Encode(tss)
				if err != nil {
					log.Println("Unable to gob-encode nap test-score-summary: ", err)
				}

				// {NAPTestScoreSummary} = object
				entries = badger.EntriesSet(entries, []byte(tss.SummaryID), gtss)
				// NAPTestScoreSummary-type:{id} = {id}
				key := []byte("NAPTestScoreSummary:" + tss.SummaryID)
				entries = badger.EntriesSet(entries, key, []byte(tss.SummaryID))

				// {school}:NAPTestScoreSummary-type:{id} = {id}
				key = []byte(tss.SchoolInfoRefId + ":NAPTestScoreSummary:" + tss.SummaryID)
				entries = badger.EntriesSet(entries, key, []byte(tss.SummaryID))

				totalTestScoreSummarys++

			case "NAPEventStudentLink":
				var e xml.NAPEvent
				decoder.DecodeElement(&e, &se)
				ge, err := ge.Encode(e)
				if err != nil {
					log.Println("Unable to gob-encode nap event link: ", err)
				}

				// {NAPEventStudentLink} = object
				entries = badger.EntriesSet(entries, []byte(e.EventID), ge)
				// NAPEventStudentLink-type:{id} = {id}
				key := []byte("NAPEventStudentLink:" + e.EventID)
				entries = badger.EntriesSet(entries, key, []byte(e.EventID))

				// {school}:NAPEventStudentLink-type:{id} = {id}
				key = []byte(e.SchoolRefID + ":NAPEventStudentLink:" + e.EventID)
				entries = badger.EntriesSet(entries, key, []byte(e.EventID))

				totalEvents++

			case "NAPStudentResponseSet":
				var r xml.NAPResponseSet
				decoder.DecodeElement(&r, &se)
				gr, err := ge.Encode(r)
				if err != nil {
					log.Println("Unable to gob-encode student response set: ", err)
				}

				// {response-id} = object
				entries = badger.EntriesSet(entries, []byte(r.ResponseID), gr)
				// response-type:{id} = {id}
				key := []byte("NAPStudentResponseSet:" + r.ResponseID)
				entries = badger.EntriesSet(entries, key, []byte(r.ResponseID))

				// {test}:NAPStudentResponseSet-type:{student} = {id}
				key = []byte(r.TestID + ":NAPStudentResponseSet:" + r.StudentID)
				entries = badger.EntriesSet(entries, key, []byte(r.ResponseID))

				totalResponses++

			case "NAPCodeFrame":
				var cf xml.NAPCodeFrame
				decoder.DecodeElement(&cf, &se)
				gcf, err := ge.Encode(cf)
				if err != nil {
					log.Println("Unable to gob-encode nap codeframe: ", err)
				}

				// {NAPCodeFrame-id} = object
				entries = badger.EntriesSet(entries, []byte(cf.RefId), gcf)

				// NAPCodeFrame-type:{id} = {id}
				key := []byte("NAPCodeFrame:" + cf.RefId)
				entries = badger.EntriesSet(entries, key, []byte(cf.RefId))

				totalCodeFrames++

			case "SchoolInfo":
				var si xml.SchoolInfo
				decoder.DecodeElement(&si, &se)
				gsi, err := ge.Encode(si)
				if err != nil {
					log.Println("Unable to gob-encode schoolinfo: ", err)
				}

				// {SchoolInfo-id} = object
				entries = badger.EntriesSet(entries, []byte(si.RefId), gsi)

				// SchoolInfo-type:{id} = {id}
				key := []byte("SchoolInfo:" + si.RefId)
				entries = badger.EntriesSet(entries, key, []byte(si.RefId))

				// ASL lookup
				key = []byte(si.ACARAId)
				entries = badger.EntriesSet(entries, key, []byte(si.RefId))

				totalSchools++

			case "StudentPersonal":
				var sp xml.RegistrationRecord
				decoder.DecodeElement(&sp, &se)
				gsp, err := ge.Encode(sp)
				if err != nil {
					log.Println("Unable to gob-encode studentpersonal: ", err)
				}

				// {StudentPersonal-id} = object
				entries = badger.EntriesSet(entries, []byte(sp.RefId), gsp)

				// StudentPersonal-type:{id} = {id}
				key := []byte("StudentPersonal:" + sp.RefId)
				entries = badger.EntriesSet(entries, key, []byte(sp.RefId))

				// {ASL-school-id}:StudentPersonal-type:{id} = {id}
				key = []byte(sp.ASLSchoolId + ":StudentPersonal:" + sp.RefId)
				entries = badger.EntriesSet(entries, key, []byte(sp.RefId))

				totalStudents++

			}
		default:
		}

	}

	// fmt.Println("entries:")
	// for _, entry := range entries {
	// 	fmt.Printf("key: %s\n val: %s\n\n", entry.Key, entry.Value)
	// }

	batcherr := kv.BatchSet(entries)
	if batcherr != nil {
		log.Fatalln("batch error: ", batcherr)
	}

	log.Println("Data file read complete...")
	log.Printf("Total tests: %d \n", totalTests)
	log.Printf("Total codeframes: %d \n", totalCodeFrames)
	log.Printf("Total testlets: %d \n", totalTestlets)
	log.Printf("Total test items: %d \n", totalTestItems)
	log.Printf("Total test score summaries: %d \n", totalTestScoreSummarys)
	log.Printf("Total events: %d \n", totalEvents)
	log.Printf("Total responses: %d \n", totalResponses)
	log.Printf("Total schools: %d \n", totalSchools)
	log.Printf("Total students: %d \n", totalStudents)

	log.Printf("ingestion complete for %s", resultsFilePath)

}
