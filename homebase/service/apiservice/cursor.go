package apiservice

import (
	"encoding/json"
)

type Cursor string

func (Cursor) ImplementsGraphQLType(name string) bool {
	return name == "Cursor"
}

func (c *Cursor) UnmarshalGraphQL(input interface{}) error {
	*c = Cursor(input.(string))
	return nil
}

func (c Cursor) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(c))
}
