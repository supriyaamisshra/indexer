package fetcher

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"time"

	"go.uber.org/zap"
)

type RequestArgs struct {
	url    string
	method string
	params map[string]string
	header map[string]string
	body   []byte
}

func sendRequest(client *http.Client, args RequestArgs) ([]byte, error) {
	var req *http.Request
	var err error

	switch args.method {
	case "GET":
		req, err = http.NewRequest(args.method, args.url, nil)
		if err != nil {
			return nil, err
		}
		for k, v := range args.header {
			req.Header.Add(k, v)
		}
		query := req.URL.Query()
		for k, v := range args.params {
			query.Add(k, v)
		}
		req.URL.RawQuery = query.Encode()

	case "POST":
		req, err = http.NewRequest(args.method, args.url, bytes.NewBuffer(args.body))
		if err != nil {
			return nil, err
		}
		for k, v := range args.header {
			req.Header.Add(k, v)
		}

	default:
		return nil, fmt.Errorf("Unsupported method %s", args.method)
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Response code: %d", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func httpClient() *http.Client {
	client := new(http.Client)
	var transport http.RoundTripper = &http.Transport{
		Proxy:              http.ProxyFromEnvironment,
		DisableKeepAlives:  false,
		DisableCompression: false,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 300 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client.Transport = transport
	return client
}

func isAddress(address string) bool {
	return regexp.MustCompile("^(0x)?[0-9a-fA-F]{40}$").MatchString(address)
}

func convertTwitterHandle(inputHandle string) string {
	retHandle := inputHandle

	// Solution 1.1 - some inputs begin with "https://twitter.com"
	re1_1, _ := regexp.Compile(`\bhttps://twitter.com/`)
	if re1_1.MatchString(retHandle) && len(retHandle) > 20 {
		retHandle = retHandle[20:]
	}

	// Solution 1.2 - some inputs begin with "https://twitter/"
	re1_2, _ := regexp.Compile(`\bhttps://twitter/`)
	if re1_2.MatchString(retHandle) && len(retHandle) > 16 {
		retHandle = retHandle[16:]
	}

	// Solution 1.3 - some inputs begin with "www.twitter.com/"
	re1_3, _ := regexp.Compile(`\bwww.twitter.com/`)
	if re1_3.MatchString(retHandle) && len(retHandle) > 16 {
		retHandle = retHandle[16:]
	}

	// Solution 2 - some inputs begin with "@"
	re2, _ := regexp.Compile(`[@]{1}`)
	if re2.MatchString(retHandle) && len(retHandle) > 1 {
		retHandle = retHandle[1:]
	}

	// Solution 3 - some inputs begin with "/"
	if retHandle[0] == '/' {
		retHandle = retHandle[1:]
	}

	// Solution 4 - some inputs ends with "/"
	if retHandle[len(retHandle)-1] == '/' {
		retHandle = retHandle[:len(retHandle)-1]
	}

	// Final Check - if retHandle still contains special characters, report it
	// will fix them in future
	re0, _ := regexp.Compile(`\W+`)
	if re0.MatchString(retHandle) {
		zap.L().With(zap.Any("handle", retHandle)).Error("Unqualified Twitter Handle")
	}

	return retHandle
}
