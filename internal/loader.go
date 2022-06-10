package internal

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:generate cp ../artefacts/quotes.json ./quotes.json
var (
	//go:embed private.json
	privateMetadata []byte
	//go:embed public.json
	publicMetadata []byte
	//go:embed quotes.json
	quotes []byte
)

type Loader struct {
	PrivateFields map[string]string
	PublicFields  map[string]string
	Quotes        []string
}

func loadFields(metadata []byte) (map[string]string, error) {
	type fieldsType struct {
		Fields map[string]string `json:"metadata"`
	}
	var fields fieldsType
	err := json.Unmarshal(metadata, &fields)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal embedded metadata: %v", err)
	}
	return fields.Fields, nil
}

func validateFieldsExistence(metadata []byte) error {
	switch {
	case cap(metadata) == 0:
		return fmt.Errorf("metadata is nil, populate fields")
	case !json.Valid(metadata):
		return fmt.Errorf("metadata is not valid JSON, verify formatting")
	}
	return nil
}

func loadQuotes(quotes []byte) ([]string, error) {
	var err error
	type quotesType struct {
		Quotes []string `json:"quotes"`
	}
	var q quotesType
	err = json.Unmarshal(quotes, &q)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal embedded quotes: %v", err)
	}
	return q.Quotes, nil
}

func (l *Loader) Load() error {
	var err error
	err = validateFieldsExistence(privateMetadata)
	if err != nil {
		return fmt.Errorf("failed to validate private metadata: %v", err)
	}
	err = validateFieldsExistence(publicMetadata)
	if err != nil {
		return fmt.Errorf("failed to validate public metadata: %v", err)
	}
	l.PrivateFields, err = loadFields(privateMetadata)
	if err != nil {
		return fmt.Errorf("failed to load private metadata: %v", err)
	}
	l.PublicFields, err = loadFields(publicMetadata)
	if err != nil {
		return fmt.Errorf("failed to load public metadata: %v", err)
	}
	l.Quotes, err = loadQuotes(quotes)
	if err != nil {
		return fmt.Errorf("failed to load quotes: %v", err)
	}
	return nil
}
