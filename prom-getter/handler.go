package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
)

type Prometheus struct {
	Status string `json:"status"`
	Data   struct {
		Result []struct {
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

// Handle a serverless request
func Handle(req []byte) string {
	var endpoint = os.Getenv("PROMETHEUS_URL")
	var metric = os.Getenv("METRIC")

	u, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}

	u.Path = path.Join(u.Path, "api/v1/query")
	q := u.Query()
	q.Set("query", metric)
	u.RawQuery = q.Encode()
	res, err := http.Get(u.String())

	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Sprintf("Error: StatusCode=%d, url=%s", res.StatusCode, u.String())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}
	bytes := []byte(body)
	var p Prometheus
	json.Unmarshal(bytes, &p)
	return fmt.Sprintf("%s", p.Data.Result[0].Value[1])
}
