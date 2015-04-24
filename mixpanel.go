package mixpanel

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	DEFAULT_EXPIRE_IN_DAYS int64 = 5
)

type Mixpanel struct {
	ApiKey  string
	Secret  string
	Format  string
	BaseUrl string
}

type EventQueryResult struct {
	LegendSize int `json:"legend_size"`
	Data       struct {
		Series []string                  `json:"series"`
		Values map[string]map[string]int `json:"values"`
	} `json:data`
}

type ExportQueryResult struct {
	Event      string                 `json:event`
	Properties map[string]interface{} `json:properties`
}

type SegmentationQueryResult struct {
	LegendSize int `json:"legend_size"`
	Data       struct {
		Series []string                  `json:"series"`
		Values map[string]map[string]int `json:"values"`
	} `json:data`
}

type TopEventsResult struct {
	Type   string `json:"type"`
	Events []struct {
		Amount           int     `json:"amount"`
		Event            string  `json:"event"`
		PercentageChange float64 `json:"percent_change"`
	} `json:events`
}

type CommonEventsResult []string

func NewMixpanel(key string, secret string) *Mixpanel {
	m := new(Mixpanel)
	m.Secret = secret
	m.ApiKey = key
	m.Format = "json"
	m.BaseUrl = "http://mixpanel.com/api/2.0"
	return m
}

func (m *Mixpanel) AddExpire(params *map[string]string) {
	if (*params)["expire"] == "" {
		(*params)["expire"] = fmt.Sprintf("%d", ExpireInDays(DEFAULT_EXPIRE_IN_DAYS))
	}
}

func (m *Mixpanel) AddSig(params *map[string]string) {
	keys := make([]string, 0)

	(*params)["api_key"] = m.ApiKey
	(*params)["format"] = m.Format

	for k, _ := range *params {
		keys = append(keys, k)
	}
	sort.StringSlice(keys).Sort()
	// fmt.Println(s)

	var buffer bytes.Buffer
	for _, key := range keys {
		value := (*params)[key]
		buffer.WriteString(fmt.Sprintf("%s=%s", key, value))
	}
	buffer.WriteString(m.Secret)
	// fmt.Println(buffer.String())

	hash := md5.New()
	hash.Write(buffer.Bytes())
	sigHex := fmt.Sprintf("%x", hash.Sum([]byte{}))
	(*params)["sig"] = sigHex
}

func (m *Mixpanel) MakeRequest(action string, params map[string]string) ([]byte, error) {
  event, ok := params["event"]
  delete(params, "event")
  if ok && event != "" {
    events := strings.Split(event, ",")
    bytes, err := json.Marshal(events)
    if err != nil {
      return []byte{}, err
    }
    params["event"] = string(bytes)
  }
  
  m.AddExpire(&params)
  m.AddSig(&params)
  
  var buffer bytes.Buffer
  for key,value := range params {
    value = url.QueryEscape(value)
    buffer.WriteString(fmt.Sprintf("%s=%s&", key, value))
  }
  
  uri := fmt.Sprintf("%s/%s?%s", m.BaseUrl, action, buffer.String())
  uri = uri[:len(uri)-1]
  // fmt.Println(uri)
  
  var bytes []byte  
  client := new(http.Client)
  req, err := http.NewRequest("GET", uri, nil)
  if err != nil {
    return bytes, err
  }
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
  resp, err := client.Do(req)
  if err != nil {
    return bytes, err
  }
  // fmt.Printf("%+v",resp)
  defer resp.Body.Close()
  bytes, err = ioutil.ReadAll(resp.Body)
  // fmt.Println(string(bytes))
  return bytes, err
}
  
func (m *Mixpanel) EventQuery(params map[string]string) (EventQueryResult,error) {
  m.BaseUrl = "http://mixpanel.com/api/2.0"
  var result EventQueryResult
  bytes, err := m.MakeRequest("events/properties", params)
  if err != nil {
    return result, err
  }
  err = json.Unmarshal(bytes, &result)
  return result, err
}

func (m *Mixpanel) ExportQuery(params map[string]string) ([]ExportQueryResult, error) {
  m.BaseUrl = "http://data.mixpanel.com/api/2.0"
  var results []ExportQueryResult
  bytes, err := m.MakeRequest("export", params)
  if err != nil {
    return results, err
  }
  str := string(bytes)
  for _, s := range strings.Split(str, "\n") {
    var result ExportQueryResult
    err := json.Unmarshal([]byte(s),&result)
    if err != nil {
      // ignore this error, extra line or one event being bad
    }
    results = append(results, result)
  }
  return results, nil
}

func (m *Mixpanel) PeopleQuery(params map[string]string) (map[string]interface{}, error) {
  var result map[string]interface{}
  m.BaseUrl = "http://mixpanel.com/api/2.0"
  bytes, err := m.MakeRequest("engage", params)
  if err != nil {
    return result, err
  }
  json.Unmarshal(bytes,&result)
  return result, nil
}

func (m *Mixpanel) UserInfo(id string) (map[string]interface{}, error) {
  params := map[string]string{
    "distinct_id": id,
  }
  var result map[string]interface{}
  result, err := m.PeopleQuery(params)
  if err != nil {
    return result, err
  }
  if len(result["results"].([]interface{})) == 0 {
    return make(map[string]interface{}), nil
  }
  return result["results"].([]interface{})[0].(map[string]interface{})["$properties"].(map[string]interface{}), nil
}

func (m *Mixpanel) SegmentationQuery(params map[string]string) (SegmentationQueryResult, error) {
	m.BaseUrl = "http://mixpanel.com/api/2.0"
	bytes, err := m.MakeRequest("segmentation", params)

	var result SegmentationQueryResult
	if err != nil {
    return result, err
	}

	err = json.Unmarshal(bytes, &result)
	return result, err
}

func (m *Mixpanel) TopEvents(params map[string]string) (TopEventsResult, error) {
	m.BaseUrl = "http://mixpanel.com/api/2.0"
  
	var result TopEventsResult
	bytes, err := m.MakeRequest("events/top", params)
	if err != nil {
    return result, err
	}

	err = json.Unmarshal(bytes, &result)
	return result, err

}

func (m *Mixpanel) MostCommonEventsLast31Days(params map[string]string) (CommonEventsResult, error) {
	m.BaseUrl = "http://mixpanel.com/api/2.0"
	bytes, err := m.MakeRequest("events/names", params)

	var result CommonEventsResult
	if err != nil {
    return result, err
	}

	err = json.Unmarshal(bytes, &result)
	return result, err
}

func ExpireInDays(days int64) int64 {
	return time.Now().Add(time.Duration(int64(time.Hour) * days * 24)).Unix()
}

func ExpireInHours(hours int64) int64 {
	return time.Now().Add(time.Duration(int64(time.Hour) * hours)).Unix()
}