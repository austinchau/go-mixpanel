go-mixpanel
===========

Golang client for Mixpanel API


## Getting Started

``` 
package main

import (
  "fmt"
  "github.com/austinchau/go-mixpanel"
)

func main() {
  key = "YOUR_API_KEY"
  secret = "YOUR_API_SECRET"
  mp := mixpanel.NewMixpanel(key, secret)
  
  params := map[string]string{
    "from_date": "2015-01-01",
    "to_date": "2015-01-01",
    "event": "EventA,EventB" // omit "event" to retrieve all events
  }

  result, err := mp.ExportQuery(params)
  if err != nil {
    panic(err)
  }  
  fmt.Printf("%+v", result)
}
```

** To install commandline tool to dump data using MixPanel Data Export API **


```
go install github.com/austinchau/go-mixpanel/mixpanelexport
```

```
mixpanelexport -start=2015-01-01 -end=2015-01-01 -key="YOUR_API_KEY" -secret="YOUR_API_SECRET" -event="EventA,EventB"
```