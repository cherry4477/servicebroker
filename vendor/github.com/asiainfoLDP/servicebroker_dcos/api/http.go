package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"errors"
	"io"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
)

func httpPost(url string, bodyType string, body interface{}, credential ...string) ([]byte, error) {
	return httpAction("POST", url, bodyType, body, credential...)
}

func httpGet(url string, credential ...string) ([]byte, error) {
	var resp *http.Response
	var err error
	if len(credential) == 2 {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("[http] err %s, %s\n", url, err)
		}
		req.Header.Set(credential[0], credential[1])
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("http get err:%s", err.Error())
			return nil, err
		}
		switch resp.StatusCode {
		case 502:
			fmt.Printf("unknown err %v", err)
			fmt.Println(" %s %s", url, credential, resp.StatusCode)
		case 504:
			fmt.Printf("unknown err %v", err)
			fmt.Println(" %s %s", url, credential, resp.StatusCode)
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("read body  err %v", err)
			}
			fmt.Printf("%s", string(b))

		case 404:
			return nil, ErrNotFound
		case 401:
			return nil, ErrUnauthorized
		case 200:
			return ioutil.ReadAll(resp.Body)
		}
		if resp.StatusCode < 200 || resp.StatusCode > 300 {
			return nil, fmt.Errorf("[http get] status err %s, %d\n", url, resp.StatusCode)
		}
	} else {
		resp, err = http.Get(url)
		if err != nil {
			fmt.Printf("http get err:%s", err.Error())
			return nil, err
		}
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("[http get] status err %s, %d\n", url, resp.StatusCode)
		}
	}

	return ioutil.ReadAll(resp.Body)
}

func httpAction(method, url string, bodyType string, body interface{}, credential ...string) ([]byte, error) {

	var req *http.Request
	var err error
	switch t := body.(type) {
	case []byte:
		req, err = http.NewRequest(method, url, bytes.NewBuffer(t))
	case io.Reader:
		req, err = http.NewRequest(method, url, t)
	}

	if err != nil {
		return nil, fmt.Errorf("[http] err %s, %s\n", url, err)
	}

	var resp *http.Response
	req.Header.Set("Content-Type", bodyType)
	if len(credential) == 2 {
		req.Header.Set(credential[0], credential[1])
	}
	resp, err = http.DefaultClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("[http] err %s, %s\n", url, err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[http] read err %s, %s\n", url, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		return nil, fmt.Errorf("[http] status err %s, %d\n", url, resp.StatusCode)
	}

	if resp.StatusCode >= 400 {
		return b, fmt.Errorf("[http] status err %s, %d\n", url, resp.StatusCode)
	}

	return b, nil
}

func httpDelete(url string, credential ...string) ([]byte, error) {
	var resp *http.Response
	var err error
	if len(credential) == 2 {
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			return nil, fmt.Errorf("[http] err %s, %s\n", url, err)
		}
		req.Header.Set(credential[0], credential[1])
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("http get err:%s", err.Error())
			return nil, err
		}
		switch resp.StatusCode {
		case 404:
			return nil, ErrNotFound
		case 200:
			return ioutil.ReadAll(resp.Body)
		}
		if resp.StatusCode < 200 || resp.StatusCode > 300 {
			return nil, fmt.Errorf("[http get] status err %s, %d\n", url, resp.StatusCode)
		}
	} else {
		resp, err = http.Get(url)
		if err != nil {
			fmt.Printf("http get err:%s", err.Error())
			return nil, err
		}
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("[http get] status err %s, %d\n", url, resp.StatusCode)
		}
	}

	return ioutil.ReadAll(resp.Body)
}

const (
	ContentType_Form = "application/x-www-form-urlencoded"
	ContentType_Json = "application/json"
)
