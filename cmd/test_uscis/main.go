package main

import (
	"fmt"
	"github.com/gonyyi/areq"
	"time"
)

func main() {
	err, a, b := CheckN400Status()
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Printf("From: %s to %s\n", a, b)
}

const (
	uscisMemphis = "https://egov.uscis.gov/processing-times/api/processingtime/N-400/MEM"
)


func CheckN400Status() (err error, low string, high string) {
	// curl 'https://egov.uscis.gov/processing-times/api/processingtime/N-400/MEM' \
	//  -H 'Connection: keep-alive' \
	//  -H 'Accept: application/json, text/plain, */*' \
	//  -H 'Referer: https://egov.uscis.gov/processing-times/' \
	//  -H 'Accept-Language: en-US,en;q=0.9' \

	type resp struct {
		Data struct {
			ProcessingTime struct {
				FormName   string      `json:"form_name"` // "N-400"
				Range      []struct {
					// {"unit":"Months","unit_en":"Months","unit_es":"Meses","value":14}
					// {"unit":"Months","unit_en":"Months","unit_es":"Meses","value":8}
					Unit   string `json:"unit"`
					Value  int    `json:"value"`
				} `json:"range"`
			} `json:"processing_time"`
		} `json:"data"`
	}

	var out resp

	req := areq.NewRequest().SetTimeout(15*time.Second, -1, -1, -1)

	err = req.Req("GET", uscisMemphis,
		areq.Do.ReqHeaderAdd("Accept", "application/json"),
		areq.Do.ReqHeaderAdd("Referer", "https://egov.uscis.gov/processing-times/"),
		areq.Do.ResJSONTo(&out),
	)

	if err != nil {
		return err, "", ""
	}

	var intLow, intHigh int
	var unitLow, unitHigh string
	for _, v := range out.Data.ProcessingTime.Range {
		if intHigh == 0 {
			intHigh = v.Value
			unitHigh = v.Unit
		}
		if intLow == 0 {
			intLow = v.Value
			unitLow = v.Unit
		}

		if v.Value > intHigh {
			intHigh = v.Value
			unitHigh = v.Unit
		}
		if v.Value < intLow {
			intLow = v.Value
			unitLow = v.Unit
		}
	}

	return nil, fmt.Sprintf("%d %s", intLow, unitLow), fmt.Sprintf("%d %s", intHigh, unitHigh)
}
