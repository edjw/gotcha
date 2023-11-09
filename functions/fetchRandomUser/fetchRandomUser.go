package fetchRandomUser

import (
	"encoding/json"
	"net/http"

	"github.com/edjw/gotcha/html/partials"
)

func FetchRandomUser() (*partials.PersonData, error) {
	url := "https://random-data-api.com/api/v2/users"
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var person partials.PersonData
	if err := json.NewDecoder(resp.Body).Decode(&person); err != nil {
		return nil, err
	}

	return &person, nil
}
