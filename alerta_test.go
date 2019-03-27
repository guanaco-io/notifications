package main

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

var mockResponse = `
{
  "alerts": [
    {
      "attributes": {
        "ip": "172.18.0.1"
      },
      "correlate": [
        "tilroy-to-crmSuccess",
        "tilroy-to-crmFailure",
        "tilroy-to-crmWarning"
      ],
      "createTime": "2019-03-27T06:38:44.385Z",
      "customer": null,
      "duplicateCount": 41,
      "environment": "Production",
      "event": "tilroy-to-crmSuccess",
      "group": "Misc",
      "history": [
        {
          "event": "tilroy-to-crmSuccess",
          "href": "http://integratie.tjcc.be/alerta/api/alert/3b00007d-a9c0-4af6-b20e-4d0332586d18",
          "id": "3b00007d-a9c0-4af6-b20e-4d0332586d18",
          "severity": "normal",
          "status": null,
          "text": "",
          "type": "severity",
          "updateTime": "2019-03-27T06:38:44.385Z",
          "value": null
        }
      ],
      "href": "http://integratie.tjcc.be/alerta/api/alert/3b00007d-a9c0-4af6-b20e-4d0332586d18",
      "id": "3b00007d-a9c0-4af6-b20e-4d0332586d18",
      "lastReceiveId": "65ee2a24-2d57-4e8d-b976-d8fe00b0607d",
      "lastReceiveTime": "2019-03-27T12:47:27.570Z",
      "origin": "uwsgi/40a6a9ba9422",
      "previousSeverity": "indeterminate",
      "rawData": null,
      "receiveTime": "2019-03-27T06:38:44.386Z",
      "repeat": true,
      "resource": "UnmappedType:Customer",
      "service": [
        "servicemix",
        "tilroy"
      ],
      "severity": "normal",
      "status": "closed",
      "tags": [],
      "text": "",
      "timeout": 604800,
      "trendIndication": "noChange",
      "type": "exceptionAlert",
      "value": null
    },
    {
      "attributes": {
        "ip": "172.18.0.1"
      },
      "correlate": [
        "tilroy-to-crmSuccess",
        "tilroy-to-crmFailure",
        "tilroy-to-crmWarning"
      ],
      "createTime": "2019-03-27T12:47:27.520Z",
      "customer": null,
      "duplicateCount": 0,
      "environment": "Production",
      "event": "tilroy-to-crmSuccess",
      "group": "Misc",
      "history": [
        {
          "event": "tilroy-to-crmSuccess",
          "href": "http://integratie.tjcc.be/alerta/api/alert/21a09544-5abe-46db-841e-4d120ef730ab",
          "id": "21a09544-5abe-46db-841e-4d120ef730ab",
          "severity": "normal",
          "status": null,
          "text": "",
          "type": "severity",
          "updateTime": "2019-03-27T12:47:27.520Z",
          "value": null
        }
      ],
      "href": "http://integratie.tjcc.be/alerta/api/alert/21a09544-5abe-46db-841e-4d120ef730ab",
      "id": "21a09544-5abe-46db-841e-4d120ef730ab",
      "lastReceiveId": "21a09544-5abe-46db-841e-4d120ef730ab",
      "lastReceiveTime": "2019-03-27T12:47:27.521Z",
      "origin": "uwsgi/40a6a9ba9422",
      "previousSeverity": "indeterminate",
      "rawData": null,
      "receiveTime": "2019-03-27T12:47:27.521Z",
      "repeat": false,
      "resource": "TILROY/CUSTOMER/4566669",
      "service": [
        "servicemix",
        "tilroy"
      ],
      "severity": "normal",
      "status": "closed",
      "tags": [],
      "text": "",
      "timeout": 604800,
      "trendIndication": "noChange",
      "type": "exceptionAlert",
      "value": null
    },
    {
      "attributes": {
        "ip": "172.18.0.1"
      },
      "correlate": [
        "tilroy-to-crmSuccess",
        "tilroy-to-crmFailure",
        "tilroy-to-crmWarning"
      ],
      "createTime": "2019-03-27T06:38:42.202Z",
      "customer": null,
      "duplicateCount": 27,
      "environment": "Production",
      "event": "tilroy-to-crmSuccess",
      "group": "Misc",
      "history": [
        {
          "event": "tilroy-to-crmSuccess",
          "href": "http://integratie.tjcc.be/alerta/api/alert/8bebbfd1-b35b-4eaa-a83e-2fbfcb5accb1",
          "id": "8bebbfd1-b35b-4eaa-a83e-2fbfcb5accb1",
          "severity": "normal",
          "status": null,
          "text": "",
          "type": "severity",
          "updateTime": "2019-03-27T06:38:42.202Z",
          "value": null
        }
      ],
      "href": "http://integratie.tjcc.be/alerta/api/alert/8bebbfd1-b35b-4eaa-a83e-2fbfcb5accb1",
      "id": "8bebbfd1-b35b-4eaa-a83e-2fbfcb5accb1",
      "lastReceiveId": "21815ed5-c170-4cf1-ac65-8cd7dd743dcb",
      "lastReceiveTime": "2019-03-27T12:47:25.462Z",
      "origin": "uwsgi/40a6a9ba9422",
      "previousSeverity": "indeterminate",
      "rawData": null,
      "receiveTime": "2019-03-27T06:38:42.203Z",
      "repeat": true,
      "resource": "UnmappedType:byte[]",
      "service": [
        "servicemix",
        "tilroy"
      ],
      "severity": "normal",
      "status": "closed",
      "tags": [],
      "text": "",
      "timeout": 604800,
      "trendIndication": "noChange",
      "type": "exceptionAlert",
      "value": null
    } ],
  "autoRefresh": true,
  "lastTime": "2019-03-27T01:33:57.277Z",
  "more": false,
  "page": 1,
  "pageSize": 1000,
  "pages": 1,
  "severityCounts": {
    "minor": 17,
    "warning": 1
  },
  "status": "ok",
  "statusCounts": {
    "open": 18
  },
  "total": 18
}
`

func TestUnMarshallAlertsResponse(t *testing.T) {

	var alertsResponse AlertsResponse
	//alertsResponse := AlertsResponse{}

	err := json.Unmarshal([]byte(mockResponse), &alertsResponse)

	if err != nil {
		t.Fatalf("cannot unmarshal data: %v", err)
	}

	fmt.Println(alertsResponse)
	log.Printf("Parsed response: %v", alertsResponse)
}

func TestGeneric(t *testing.T) {

	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(mockResponse), &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat)
}