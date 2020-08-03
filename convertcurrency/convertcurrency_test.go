package convertcurrency

import (
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_convertFromRateData(t *testing.T) {
	currencyEps := 1e-3

	type args struct {
		value   float64
		fromCur string
		toCur   string
		date    string
		rates   rateData
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "basic case",
			args: args{
				value:   14.50,
				fromCur: "USD",
				toCur:   "CHF",
				date:    "2020-06-08",
				rates:   mockRateData,
			},
			want:    13.96,
			wantErr: false,
		},
		{
			name: "convert to EUR",
			args: args{
				value:   14.50,
				fromCur: "USD",
				toCur:   "EUR",
				date:    "2020-06-08",
				rates:   mockRateData,
			},
			want:    12.85,
			wantErr: false,
		},
		{
			name: "identical currency",
			args: args{
				value:   14.50,
				fromCur: "USD",
				toCur:   "USD",
				date:    "2020-06-08",
				rates:   mockRateData,
			},
			want:    14.50,
			wantErr: false,
		},
		{
			name: "value negative",
			args: args{
				value:   -14.50,
				fromCur: "USD",
				toCur:   "CHF",
				date:    "2020-06-08",
				rates:   mockRateData,
			},
			wantErr: true,
		},
		{
			name: "invalid date",
			args: args{
				value:   14.50,
				fromCur: "USD",
				toCur:   "CHF",
				date:    "today",
				rates:   mockRateData,
			},
			wantErr: true,
		},
		{
			name: "invalid currency to convert from",
			args: args{
				value:   14.50,
				fromCur: "dollars",
				toCur:   "CHF",
				date:    "2020-06-08",
				rates:   mockRateData,
			},
			wantErr: true,
		},
		{
			name: "invalid currency to convert to",
			args: args{
				value:   14.50,
				fromCur: "USD",
				toCur:   "francs",
				date:    "2020-06-08",
				rates:   mockRateData,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertFromRateData(tt.args.value, tt.args.fromCur, tt.args.toCur, tt.args.date, tt.args.rates)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertFromRateData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if math.Abs(got-tt.want) >= currencyEps {
				t.Errorf("convertFromRateData() = %v, want %v", got, tt.want)
			}
		})
	}
}

// not testing standard error handling in fetch()
func Test_fetch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, mockXML)
	}))
	defer ts.Close()

	rates := newRateData()
	err := rates.fetch(ts.URL)
	fmt.Println(err)
	if !reflect.DeepEqual(rates, mockRateData) {
		t.Errorf("rate data after fetch() = %v, want %v", rates.rate, mockRateData)
	}
}

var mockRateData = rateData{
	rate: map[rateKey]float64{
		{"2020-06-08", "USD"}: 1.1285,
		{"2020-06-08", "CHF"}: 1.0861,
		{"2020-06-05", "CHF"}: 1.0866,
	},
}

var mockXML = `<gesmes:Envelope xmlns:gesmes="http://www.gesmes.org/xml/2002-08-01" xmlns="http://www.ecb.int/vocabulary/2002-08-01/eurofxref">
<gesmes:subject>Reference rates</gesmes:subject>
<gesmes:Sender>
<gesmes:name>European Central Bank</gesmes:name>
</gesmes:Sender>
<Cube>
<Cube time="2020-06-08">
<Cube currency="USD" rate="1.1285"/>
<Cube currency="CHF" rate="1.0861"/>
</Cube>
<Cube time="2020-06-05">
<Cube currency="CHF" rate="1.0866"/>
</Cube>
</Cube>
</gesmes:Envelope>`
