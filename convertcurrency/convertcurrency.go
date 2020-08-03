package convertcurrency

import (
	"encoding/xml"
	"fmt"
	"math"
	"net/http"
)

const rateURL = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml"

type rateKey struct {
	date     string
	currency string
}

// map used because of potential future batch currency convertions etc.
type rateData struct {
	rate map[rateKey]float64
}

// fields must be exported for both structs because of unmarshaling / decoding
type envelopes struct {
	Data []struct {
		Date  string `xml:"time,attr"`
		Rates []struct {
			Currency string  `xml:"currency,attr"`
			Rate     float64 `xml:"rate,attr"`
		} `xml:"Cube"`
	} `xml:"Cube>Cube"`
}

// ConvertCurrency converts from one currency to another via ECB exchange rates updated daily.
func ConvertCurrency(value float64, fromCur string, toCur string, date string) (float64, error) {
	rates := newRateData()
	err := rates.fetch(rateURL)
	if err != nil {
		return 0, err
	}

	retval, err := convertFromRateData(value, fromCur, toCur, date, rates)
	if err != nil {
		return 0, err
	}
	return retval, nil
}

func newRateData() (rd rateData) {
	return rateData{
		rate: make(map[rateKey]float64),
	}
}

func convertFromRateData(value float64, fromCur string, toCur string, date string, rates rateData) (float64, error) {
	// sanity checks
	if value < 0 {
		return 0, fmt.Errorf("trying to convert negative value")
	}
	if fromCur == toCur {
		return value, nil
	}

	if fromCur != "EUR" {
		// convert from expected currency to EUR
		if rate, ok := rates.rate[rateKey{date, fromCur}]; ok {
			value /= rate
		} else {
			return 0, fmt.Errorf("invalid date or currency to convert from")
		}
	}
	if toCur != "EUR" {
		// convert from EUR to expected currency
		if rate, ok := rates.rate[rateKey{date, toCur}]; ok {
			value *= rate
		} else {
			return 0, fmt.Errorf("invalid date or currency to convert to")
		}
	}

	return math.Round(value*100) / 100, nil
}

func (rd *rateData) fetch(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to GET currency data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code not 200: %v", resp.StatusCode)
	}

	rawXML := resp.Body
	var decodedXML envelopes

	err = xml.NewDecoder(rawXML).Decode(&decodedXML)
	if err != nil {
		return fmt.Errorf("failed to decode XML: %v", err)
	}

	for _, d := range decodedXML.Data {
		for _, r := range d.Rates {
			rd.rate[rateKey{d.Date, r.Currency}] = r.Rate
		}
	}

	return nil
}
