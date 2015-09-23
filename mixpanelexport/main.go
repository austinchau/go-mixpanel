package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	mixpanel "github.com/austinchau/go-mixpanel"
)

var (
	m      *mixpanel.Mixpanel
	key    string
	secret string
	event  string
	start  string
	end    string
	output string
)

func init() {
	flag.StringVar(&event, "event", "", "Event Name")
	flag.StringVar(&start, "start", "2015-01-01", "Start Date")
	flag.StringVar(&end, "end", "2015-01-01", "End Date")
	flag.StringVar(&output, "output", "", "Output File")
	flag.Parse()
}

func main() {
	if start == "" || end == "" {
		flag.Usage()
		os.Exit(1)
	}

	m = mixpanel.NewMixpanel(key, secret)
	params := map[string]string{
		"from_date": start,
		"to_date":   end,
		"event":     event,
	}

	m.BaseUrl = "http://data.mixpanel.com/api/2.0"
	bytes, err := m.MakeRequest("export", params)
	if err != nil {
		panic(err)
	}
	if output == "" {
		fmt.Println(string(bytes))
	} else {
		err := ioutil.WriteFile(output, bytes, 0644)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
}
