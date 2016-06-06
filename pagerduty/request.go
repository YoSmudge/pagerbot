package pagerduty

import(
  "fmt"
  "sync"
  "time"
  "net/url"
  "encoding/json"
  log "github.com/Sirupsen/logrus"
  "gopkg.in/jmcvetta/napping.v3"
)

type response interface{
  Add(interface{})
}

type pagination struct{
  Limit       int
  Offset      int
  Total       int
}

// Request a list of items (users, schedules etc.) and paginate
// Pagerduty has a weird response format so we make the assumption that the
// data we want is the same key as the URL (I.E. schedules, users etc.)
// There is a super weird JSON hack in here to decode the responses properly
func (a *Api) requestThing(thing string, r interface{}) error{
  var err error
  var items []interface{}

  var page int = 0
  var perPage int = 25

  for{
    q := url.Values{}
    q.Set("limit", fmt.Sprintf("%d", perPage))
    q.Set("offset", fmt.Sprintf("%d", page*perPage))

    path := fmt.Sprintf("%s?%s", thing, q.Encode())

    var resp *napping.Response
    resp, err = a.request(path)
    if err != nil {
      return err
    }

    var pages pagination
    resp.Unmarshal(&pages)

    var rawResponse map[string]interface{}
    resp.Unmarshal(&rawResponse)

    for _,i := range rawResponse[thing].([]interface{}){
      items = append(items, i)
    }

    if len(items) == pages.Total{
      break
    }

    page += 1
  }

  itemsEncoded, _ := json.Marshal(&items)
  json.Unmarshal(itemsEncoded, &r)
  // Told you

  return nil
}

var requestLock sync.Mutex
const requestPause time.Duration = 2*time.Second
var limiter <-chan time.Time = time.Tick(requestPause)

// Send request to Pagerduty
// Ensures requests are no more frequent than `requestPause` to permit Goroutines calling without hitting rate limit
func (a *Api) request(path string) (*napping.Response, error){
  requestLock.Lock()
  defer func(){
    go func(){
      <- limiter
      requestLock.Unlock()
    }()
  }()

  var resp *napping.Response
  var err error
  var u url.URL

  u.Host = fmt.Sprintf("%s.pagerduty.com", a.org)
  u.Scheme = "https"

  fullPath := fmt.Sprintf("%s/api/v1/%s", u.String(), path)

  resp, err = a.http.Get(fullPath, &url.Values{}, nil, nil)
  if err != nil || resp.Status() != 200 {
    fields := log.Fields{
      "url": fullPath,
      "error": err,
    }
    var status int
    if err == nil{
      fields["status"] = resp.Status()
      status = resp.Status()
    }
    log.WithFields(fields).Warning("Error from Pagerduty")

    return resp, fmt.Errorf("Got error from Pagerduty: %d (%s)", status, err)
  }

  log.WithFields(log.Fields{
    "url": fullPath,
    "status": resp.Status(),
  }).Debug("Pagerduty request")

  return resp, nil
}
