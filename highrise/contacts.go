package highrise

import (
	"bytes"
	"encoding/xml"
)

type Contact struct {
	XMLName xml.Name `xml:"person"`
	ID      string   `xml:"id,omitempty"`
	First   string   `xml:"first-name"`
	Last    string   `xml:"last-name"`
}

type Contacts []Contact

type xmlPeople struct {
	XMLName xml.Name  `xml:"people"`
	People  []Contact `xml:"person"`
}

func (c *Client) GetContacts() (Contacts, error) {
	response, err := c.request("GET", "/people.xml", nil)
	if err != nil {
		return nil, err
	}

	var people xmlPeople
	err = xml.NewDecoder(response.Body).Decode(&people)
	if err != nil {
		return nil, err
	}

	return people.People, err
}

func (c *Client) CreateContact(first, last string) (*Contact, error) {
	contact := &Contact{
		First: first,
		Last:  last,
	}

	body, err := xml.Marshal(contact)
	if err != nil {
		return nil, err
	}

	response, err := c.request("POST", "/people.xml", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	err = xml.NewDecoder(response.Body).Decode(&contact)
	if err != nil {
		return nil, err
	}

	return contact, nil
}
