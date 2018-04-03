package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/fabioberger/airtable-go"
	"github.com/kamaln7/airtable-highrise/highrise"
)

var (
	highriseTeam  = flag.String("highrise.team", "", "highrise team name")
	highriseKey   = flag.String("highrise.key", "", "highrise API key")
	airtableBase  = flag.String("airtable.base", "", "airtable base ID")
	airtableKey   = flag.String("airtable.key", "", "airtable API key")
	airtableTable = flag.String("airtable.table", "Contacts", "airtable table name")
)

type airtableContact struct {
	ID     string
	Fields struct {
		First, Last, HighriseID string
	}
}

func main() {
	flag.Parse()

	air, err := airtable.New(*airtableKey, *airtableBase)
	if err != nil {
		log.Fatalf("could not connect to airtable: %v\n", err)
	}

	hr := highrise.New(*highriseKey, *highriseTeam)

	hrContacts, err := hr.GetContacts()
	if err != nil {
		log.Fatalf("could not get highrise contacts: %v\n", err)
	}

	var airContacts []airtableContact
	err = air.ListRecords(*airtableTable, &airContacts, airtable.ListParameters{
		FilterByFormula: "{Highrise ID} = ''",
		Fields:          []string{"First", "Last"},
	})
	if err != nil {
		log.Fatalf("could not get airtable contacts: %v\n", err)
	}

	log.Printf("found %d airtable contacts without highrise info\n", len(airContacts))
	for _, contact := range airContacts {
		id := ""
		for _, candidate := range hrContacts {
			if candidate.First == contact.Fields.First && candidate.Last == contact.Fields.Last {
				// found an existing highrise contact
				id = candidate.ID
				break
			}
		}

		if id == "" {
			// create a new Highrise contact
			hrContact, err := hr.CreateContact(contact.Fields.First, contact.Fields.Last)
			if err != nil {
				log.Printf("could not create new highrise contact (%s %s): %v\n", contact.Fields.First, contact.Fields.Last, err)
				continue
			}

			id = hrContact.ID
		}

		err = air.UpdateRecord(*airtableTable, contact.ID, map[string]interface{}{"Highrise ID": id}, &contact)
		if err != nil {
			log.Printf("could not update airtable record (%s %s): %v\n", contact.Fields.First, contact.Fields.Last, err)
			continue
		}

		fmt.Printf("updated %s %s (highrise %s)\n", contact.Fields.First, contact.Fields.Last, id)
	}
}
