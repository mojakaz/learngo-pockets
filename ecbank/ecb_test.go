package ecbank

import (
	"errors"
	"fmt"
	"learngo-pockets/money"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
)

const (
	eurofxref = `<gesmes:Envelope xmlns:gesmes="http://www.gesmes.org/xml/2002-08-01" xmlns="http://www.ecb.int/vocabulary/2002-08-01/eurofxref">
<gesmes:subject>Reference rates</gesmes:subject>
<gesmes:Sender>
<gesmes:name>European Central Bank</gesmes:name>
</gesmes:Sender>
<Cube>
<Cube time="2023-09-08">
<Cube currency="USD" rate="1.0704"/>
<Cube currency="JPY" rate="157.84"/>
<Cube currency="BGN" rate="1.9558"/>
<Cube currency="CZK" rate="24.452"/>
<Cube currency="DKK" rate="7.4591"/>
<Cube currency="GBP" rate="0.85735"/>
<Cube currency="HUF" rate="385.70"/>
<Cube currency="PLN" rate="4.6213"/>
<Cube currency="RON" rate="4.9630"/>
<Cube currency="SEK" rate="11.9040"/>
<Cube currency="CHF" rate="0.9543"/>
<Cube currency="ISK" rate="143.30"/>
<Cube currency="NOK" rate="11.4220"/>
<Cube currency="TRY" rate="28.7390"/>
<Cube currency="AUD" rate="1.6743"/>
<Cube currency="BRL" rate="5.3238"/>
<Cube currency="CAD" rate="1.4623"/>
<Cube currency="CNY" rate="7.8565"/>
<Cube currency="HKD" rate="8.3915"/>
<Cube currency="IDR" rate="16438.40"/>
<Cube currency="ILS" rate="4.1162"/>
<Cube currency="INR" rate="88.8610"/>
<Cube currency="KRW" rate="1428.51"/>
<Cube currency="MXN" rate="18.7019"/>
<Cube currency="MYR" rate="5.0057"/>
<Cube currency="NZD" rate="1.8127"/>
<Cube currency="PHP" rate="60.660"/>
<Cube currency="SGD" rate="1.4605"/>
<Cube currency="THB" rate="38.042"/>
<Cube currency="ZAR" rate="20.4370"/>
</Cube>
</Cube>
</gesmes:Envelope>`
	eurofxrefWrongFormat = `<?xml...>`
)

func mustParseCurrency(t *testing.T, code string) money.Currency {
	t.Helper()
	currency, err := money.ParseCurrency(code)
	if err != nil {
		t.Fatalf("cannot parse currency %s code", code)
	}
	return currency
}

func mustParseCache(t *testing.T, fileName string) *os.File {
	testFile, err := os.OpenFile("test", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Fatalf("could not open file: %s", err.Error())
	}
	return testFile
}

func TestEuroCentralBank_FetchExchangeRate_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, eurofxref)
		if err != nil {
			t.Fatalf("failed to write response: %s", err.Error())
		}
	}))
	defer ts.Close()

	proxyURL, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatalf("failed to parse proxy URL: %v", err)
	}

	tt := map[string]struct {
		sourceCode string
		targetCode string
		expected   money.ExchangeRate
		err        error
	}{
		"USD to RON": {
			sourceCode: "USD",
			targetCode: "RON",
			expected:   4.63658445441,
			err:        nil,
		},
		"USD to JPY": {
			sourceCode: "USD",
			targetCode: "JPY",
			expected:   147.458893871,
			err:        nil,
		},
		"SGD to CAD": {
			sourceCode: "SGD",
			targetCode: "CAD",
			expected:   1.0012324,
			err:        nil,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			ecb := EuroCentralBank{
				url:   ts.URL,
				cache: mustParseCache(t, "test"),
				client: &http.Client{
					Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
					Timeout:   time.Second,
				}}

			got, err := ecb.FetchExchangeRate(mustParseCurrency(t, tc.sourceCode), mustParseCurrency(t, tc.targetCode))
			if err != nil {
				t.Errorf("expected no error, got %v", err.Error())
			}
			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestEuroCentralBank_FetchExchangeRate_Failure(t *testing.T) {
	tt := map[string]struct {
		sourceCode string
		targetCode string
		f          func(w http.ResponseWriter, r *http.Request)
		err        error
	}{
		"500": {
			sourceCode: "USD",
			targetCode: "RON",
			f:          func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusInternalServerError) },
			err:        ErrServerSide,
		},
		"400": {
			sourceCode: "USD",
			targetCode: "RON",
			f:          func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusBadRequest) },
			err:        ErrClientSide,
		},
		"300": {
			sourceCode: "USD",
			targetCode: "RON",
			f:          func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusMultipleChoices) },
			err:        ErrUnknownStatusCode,
		},
		"No such currency 1": {
			sourceCode: "USD",
			targetCode: "XYZ",
			f: func(w http.ResponseWriter, r *http.Request) {
				_, err := fmt.Fprintln(w, eurofxref)
				if err != nil {
					t.Fatalf("failed to write response: %s", err.Error())
				}
			},
			err: ErrChangeRateNotFound,
		},
		"No such currency 2": {
			sourceCode: "XYZ",
			targetCode: "RON",
			f: func(w http.ResponseWriter, r *http.Request) {
				_, err := fmt.Fprintln(w, eurofxref)
				if err != nil {
					t.Fatalf("failed to write response: %s", err.Error())
				}
			},
			err: ErrChangeRateNotFound,
		},
		"wrong format": {
			sourceCode: "USD",
			targetCode: "RON",
			f: func(w http.ResponseWriter, r *http.Request) {
				_, err := fmt.Fprintln(w, eurofxrefWrongFormat)
				if err != nil {
					t.Fatalf("failed to write response: %s", err.Error())
				}
			},
			err: ErrUnexpectedFormat,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(tc.f))
			defer ts.Close()

			proxyURL, err := url.Parse(ts.URL)
			if err != nil {
				t.Fatalf("failed to parse proxy URL: %v", err)
			}

			ecb := EuroCentralBank{
				url:   ts.URL,
				cache: nil,
				client: &http.Client{
					Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
					Timeout:   time.Second,
				},
			}

			_, err = ecb.FetchExchangeRate(mustParseCurrency(t, tc.sourceCode), mustParseCurrency(t, tc.targetCode))
			if err == nil {
				t.Errorf("expected an error %v, got no error", tc.err)
			}
		})
	}
}

func TestEuroCentralBank_FetchExchangeRate_FailureCallingServer(t *testing.T) {
	t.Run("wrong url", func(t *testing.T) {
		ecb := EuroCentralBank{
			url:   "nosuchurl",
			cache: nil,
			client: &http.Client{
				Timeout: time.Second,
			},
		}
		_, err := ecb.FetchExchangeRate(mustParseCurrency(t, "USD"), mustParseCurrency(t, "RON"))
		if err == nil {
			t.Errorf("expected an error %v, got no error", ErrClientSide)
		}
	})
}

func TestEuroCentralBank_FetchExchangeRate_Timeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 2)
	}))
	defer ts.Close()

	proxyURL, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatalf("failed to parse proxy URL: %v", err)
	}

	ecb := EuroCentralBank{
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
			Timeout: time.Second,
		},
	}

	_, err = ecb.FetchExchangeRate(mustParseCurrency(t, "USD"), mustParseCurrency(t, "RON"))
	if !errors.Is(err, ErrTimeout) {
		t.Errorf("unexpected error: %v, expected %v", err, ErrTimeout)
	}
}
