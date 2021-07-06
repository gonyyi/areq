package main

import (
	"fmt"
	"github.com/gonyyi/areq"
	"time"
)

// curl 'https://egov.uscis.gov/processing-times/api/processingtime/N-400/MEM' \
//  -H 'Connection: keep-alive' \
//  -H 'Accept: application/json, text/plain, */*' \
//  -H 'Referer: https://egov.uscis.gov/processing-times/' \
//  -H 'Accept-Language: en-US,en;q=0.9' \
func main() {
	err, a, b := check()
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Printf("From: %s to %s\n", a, b)
}

const (
	uscisMemphis = "https://egov.uscis.gov/processing-times/api/processingtime/N-400/MEM"
)

type resp struct {
	Data struct {
		ProcessingTime struct {
			FormInfoEn string      `json:"form_info_en"`
			FormInfoEs string      `json:"form_info_es"`
			FormName   string      `json:"form_name"` // "N-400"
			FormNoteEn interface{} `json:"form_note_en"`
			FormNoteEs interface{} `json:"form_note_es"`
			OfficeCode string      `json:"office_code"`
			Range      []struct {
				// {"unit":"Months","unit_en":"Months","unit_es":"Meses","value":14}
				// {"unit":"Months","unit_en":"Months","unit_es":"Meses","value":8}
				Unit   string `json:"unit"`
				UnitEn string `json:"unit_en"`
				Value  int    `json:"value"`
			} `json:"range"`
			Subtypes []struct {
				FormType        string `json:"form_type"`
				PublicationDate string `json:"publication_date"` // "June 23, 2021"
				Range           []struct {
					Unit   string `json:"unit"`
					UnitEn string `json:"unit_en"`
					UnitEs string `json:"unit_es"`
					Value  int    `json:"value"`
				} `json:"range"`
				ServiceRequestDate   string      `json:"service_request_date"`
				ServiceRequestDateEn string      `json:"service_request_date_en"`
				ServiceRequestDateEs string      `json:"service_request_date_es"`
				SubtypeInfoEn        string      `json:"subtype_info_en"`
				SubtypeInfoEs        string      `json:"subtype_info_es"`
				SubtypeNoteEn        interface{} `json:"subtype_note_en"`
				SubtypeNoteEs        interface{} `json:"subtype_note_es"`
			} `json:"subtypes"`
		} `json:"processing_time"`
	} `json:"data"`
	Message string `json:"message"`
}

func check() (err error, low string, high string) {
	var out resp

	req := areq.NewRequest()
	req.SetDialTimeout(15*time.Second, 60*time.Second)

	err = req.Req("GET", uscisMemphis,
		areq.Do.SetClientTimeout(15*time.Second),
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
