package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fabioberger/airtable-go"
	"github.com/satori/go.uuid"
)

var (
	airtableBase  = flag.String("airtable.base", "", "airtable base ID")
	airtableKey   = flag.String("airtable.key", "", "airtable API key")
	airtableTable = flag.String("airtable.table", "Contacts", "airtable table name")
	outputPath    = flag.String("output.path", "", "path to output csv files. empty for stdout")
	outputURL    = flag.String("output.url", "", "url to download output files")
)

type airtableContact struct {
	ID     string
	Fields struct {
		First, Last, Email, Company, HighriseID string
	}
}

func main() {
	flag.Parse()
	*outputPath = strings.TrimRight(*outputPath, "/")
	*outputURL = strings.TrimRight(*outputURL, "/")

	air, err := airtable.New(*airtableKey, *airtableBase)
	if err != nil {
		log.Fatalf("could not connect to airtable: %v\n", err)
	}

	var airContacts []airtableContact
	err = air.ListRecords(*airtableTable, &airContacts, airtable.ListParameters{
		FilterByFormula: "{Highrise ID} = ''",
		Fields:          []string{"First", "Last", "Email", "Company"},
	})
	if err != nil {
		log.Fatalf("could not get airtable contacts: %v\n", err)
	}

	fmt.Fprintf(os.Stderr, "found %d airtable contacts without highrise info\n", len(airContacts))

	output := os.Stdout
	var filePath, fileName string
	if *outputPath != "" {
		fileName = fmt.Sprintf("contacts-%s.csv", uuid.Must(uuid.NewV4()))
		filePath = fmt.Sprintf("%s/%s", *outputPath, fileName)
		output, err = os.Create(filePath)
		if err != nil {
			log.Fatalf("couldn't open file %s for writing: %v\n", filePath, err)
		}
	}
	writer := csv.NewWriter(output)
	defer writer.Flush()
	writer.Write([]string{"First name", "Last name", "Company", "Email address - Work", "Tags"})
	for _, contact := range airContacts {
		str := []string{contact.Fields.First, contact.Fields.Last, contact.Fields.Company, contact.Fields.Email, "open-source"}

		writer.Write(str)
	}
	if *outputPath != "" {
		if *outputURL != "" {
			fmt.Printf("Wrote CSV to %s/%s\n", *outputURL, fileName)
		} else {
			fmt.Printf("Wrote CSV to %s\n", filePath)
		}
	}
}
