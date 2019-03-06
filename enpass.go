package enpass2gopass

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
)

type (
	// Export is the JSON export from enpass.
	Export struct {
		Items []EnpassItem `json:"items"`
	}

	// EnpassItem an item in enpass.
	EnpassItem struct {
		Category     string        `json:"category"`
		Fields       []EnpassField `json:"fields"`
		Note         string        `json:"note"`
		SubTitle     string        `json:"subtitle"`
		TemplateType string        `json:"template_type"`
		Title        string        `json:"title"`
	}

	// EnpassField a field in an item.
	EnpassField struct {
		Label string `json:"label"`
		Type  string `json:"type"`
		Value string `json:"value"`
	}
)

// NewExport read an Enpass JSON export.
func NewExport(location string) (*Export, error) {
	f, err := os.Open(location)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	var e Export
	if err := dec.Decode(&e); err != nil {
		return nil, err
	}

	return &e, nil
}

// Transfer transfers Enpass items into gopass.
func (e Export) Transfer() {
	for _, i := range e.Items {
		if err := i.Insert(); err != nil {
			log.Println(err)
		}
	}
}

// Insert an item into Gopass.
func (i EnpassItem) Insert() error {
	path := fmt.Sprintf("%s/%s/%s", i.Category, i.Title, i.SubTitle)
	log.Println("Inserting", path)
	return Insert(path, i.Values())
}

// Values gets the items values. A password will take priority at
// the start of the slice for clipboard compatibility. Following the password
// we will format the output for gobridge compat.
func (i EnpassItem) Values() []string {
	if i.Category == "note" {
		return []string{i.Note}
	}

	values := []string{}

	for _, f := range i.Fields {
		if f.Type == "password" {
			continue
		}
		if f.Value == "" {
			continue
		}
		if f.Type == "email" {
			values = append(values, fmt.Sprintf("user: %s", f.Value))
		} else {
			values = append(values, fmt.Sprintf("%s: %s", f.Type, f.Value))
		}
	}
	sort.Strings(values)
	values = append([]string{i.Password(), "---"}, values...)

	return values
}

// Password gets the password for the item.
func (i EnpassItem) Password() string {
	for _, f := range i.Fields {
		if f.Type == "password" {
			return f.Value
		}
	}
	return ""
}
