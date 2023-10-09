package ecbank

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"learngo-pockets/money"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

const (
	// ErrClientSide is returned if the client side made a mistake in calling the server
	ErrClientSide = ecbankError("error client side")
	// ErrServerSide is returned if the server side failed to answer the request
	ErrServerSide = ecbankError("error server side")
	// ErrUnknownStatusCode is returned if the unknown status code is returned from the server side
	ErrUnknownStatusCode = ecbankError("error unknown status code")
	// ErrTimeout is returned if the request failed because of timeout
	ErrTimeout       = ecbankError("error timeout")
	clientErrorClass = 4
	serverErrorClass = 5
)

// checkStatusCode returns a different error depending on the returned status code.
func checkStatusCode(statusCode int) error {
	switch {
	case statusCode == http.StatusOK:
		return nil
	case httpStatusClass(statusCode) == clientErrorClass:
		return fmt.Errorf("%w: %d", ErrClientSide, statusCode)
	case httpStatusClass(statusCode) == serverErrorClass:
		return fmt.Errorf("%w: %d", ErrServerSide, statusCode)
	default:
		return fmt.Errorf("%w: %d", ErrUnknownStatusCode, statusCode)
	}
}

// httpStatusClass returns the class of a http status code.
func httpStatusClass(statusCode int) int {
	const httpErrorClassSize = 100
	return statusCode / httpErrorClassSize
}

// dumpResponseBody dumps a http response body to a file and returns the body
func dumpResponseBody(resp *http.Response, file *os.File) ([]byte, error) {
	body, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, fmt.Errorf("cannot dump http response: %w", err)
	}
	_, err = file.Write(body)
	if err != nil {
		return nil, fmt.Errorf("cannot write http response to file: %w", err)
	}
	err = file.Close()
	if err != nil {
		panic(err)
	}
	return body, nil
}

// EuroCentralBank can call the bank to retrieve exchange rates.
type EuroCentralBank struct {
	url    string
	client *http.Client
	cache  *os.File
}

// NewEuroCentralBank builds a EuroCentralBank that can fetch exchange rates within a given timeout.
func NewEuroCentralBank(timeout time.Duration, fileName string) EuroCentralBank {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	return EuroCentralBank{
		client: &http.Client{Timeout: timeout},
		cache:  file,
	}
}

// FetchExchangeRate fetches the ExchangeRate for the day and returns it.
func (ecb EuroCentralBank) FetchExchangeRate(source, target money.Currency) (money.ExchangeRate, error) {
	const euroxrefURL = "http://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

	rate, err := readRateFromResponse(source.Code(), target.Code(), ecb.cache)
	if err == nil {
		err = ecb.cache.Close()
		if err != nil {
			panic(err)
		}
		return rate, nil
	}

	if ecb.url == "" {
		ecb.url = euroxrefURL
	}

	resp, err := ecb.client.Get(ecb.url)
	if err != nil {
		var urlErr *url.Error
		if ok := errors.As(err, &urlErr); ok && urlErr.Timeout() {
			return 0, fmt.Errorf("%w: %s", ErrTimeout, err.Error())
		}
		return 0, fmt.Errorf("%w: %s", ErrClientSide, err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	if err := checkStatusCode(resp.StatusCode); err != nil {
		return 0, err
	}

	body, err := dumpResponseBody(resp, ecb.cache)
	if err != nil {
		return 0, err
	}

	rate, err = readRateFromResponse(source.Code(), target.Code(), bufio.NewReader(bytes.NewReader(body)))
	if err != nil {
		return 0, err
	}

	return rate, nil
}
