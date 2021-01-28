package googleanalytics

import (
	"log"
	"net/http"
	"fmt"
	"io"
	"io/ioutil"
)

const (
	gaUrl        = "https://www.google-analytics.com/collect"
	gaVersion    = "1"
	defaultAgent = "peloton-tableau-connector"
)

var client = &http.Client{}

type Event struct {
	TrackingId     string
	CustomerId     string
	UserId         string
	EventType      string
	DocPath        string
	DocTitle       string
	DocHost        string
	IpOverride     string
	AppName        string
	AppVersion     string
	CampaignSource string
	CampaignMedium string
}

func TrackEvent(event Event) {
	if len(event.TrackingId) == 0 {
		log.Printf("google analytics tracking id is not set for the event")
		return
	} else {
		req, err := http.NewRequest(http.MethodPost, gaUrl, nil)
		if err != nil {
			log.Printf("error creating google analytics request, %s", err.Error())
			return
		}

		agent := defaultAgent

		// query parameters
		q := req.URL.Query()
		q.Add("v", gaVersion)
		q.Add("tid", event.TrackingId)
		if len(event.EventType) > 0 {
			q.Add("t", event.EventType)
		}
		if len(event.CustomerId) > 0 {
			q.Add("cid", event.CustomerId)
		}
		if len(event.UserId) > 0 {
			q.Add("uid", event.UserId)
		}
		if len(event.DocPath) > 0 {
			q.Add("dp", event.DocPath)
		}
		if len(event.DocTitle) > 0 {
			q.Add("dt", event.DocTitle)
		}
		if len(event.DocHost) > 0 {
			q.Add("dh", event.DocHost)
		}
		if len(event.IpOverride) > 0 {
			q.Add("uip", event.IpOverride)
		}
		if len(event.CampaignSource) > 0 {
			q.Add("cs", event.CampaignSource)
		}
		if len(event.CampaignMedium) > 0 {
			q.Add("cm", event.CampaignMedium)
		}
		if len(event.AppName) > 0 {
			q.Add("an", event.AppName)
			agent = event.AppName
		}
		if len(event.AppVersion) > 0 {
			q.Add("av", event.AppVersion)
			agent = fmt.Sprintf("%s/%s", agent, event.AppVersion)
		}

		req.URL.RawQuery = q.Encode()

		// headers
		req.Header.Add("Content-Length", "0")
		req.Header.Add("User-Agent", agent)
		log.Printf("posting %s %s event to google analytics with user-agent %s", event.DocPath, event.EventType, agent)

		res, err := client.Do(req)
		if err != nil {
			log.Printf("error posting event to google analytics, %s", err.Error())
			return
		}
		defer res.Body.Close()
		_, _ = io.Copy(ioutil.Discard, res.Body)

		log.Printf("response from google analytics: %d %s", res.StatusCode, res.Status)
	}
}
