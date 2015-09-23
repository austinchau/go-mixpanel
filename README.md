go-mixpanel
===========

Golang client for Mixpanel API


## Getting Started

```
package main

import (
	"fmt"

	mixpanel "github.com/austinchau/go-mixpanel"
)

func main() {
	mp := mixpanel.NewMixpanel()

	params := map[string]string{
		"from_date": "2015-01-01",
		"to_date":   "2015-01-01",
		"event":     "EventA,EventB", // omit "event" to retrieve all events
	}

	result, err := mp.ExportQuery(params)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", result)
}
```

## Install Commandline Tool


```
go install github.com/austinchau/go-mixpanel/mixpanelexport
```

```
mixpanelexport -start=2015-01-01 -end=2015-01-01 -event="EventA,EventB"
```


Expects mixpanel api credentials to be exported as environment variables.

```
export MIXPANEL_API_KEY="YOUR_KEY"
export MIXPANEL_SECRET="YOUR_SECRET"
```
