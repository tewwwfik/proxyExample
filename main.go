package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

func main() {
	reverseProxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Printf("[reverse proxy server] received request at: %s\n", time.Now())

		//get request data
		reqData, err := getRequestData(req)
		if err != nil {
			log.Printf(err.Error())
			http.Error(rw, err.Error(), http.StatusNoContent)
			return
		}

		//parse target server
		targetServer, err := url.Parse(reqData.Url)
		if err != nil {
			log.Printf(err.Error())
			http.Error(rw, err.Error(), http.StatusNoContent)
			return
		}

		// set req Host, URL and Request URI to forward a request to the target server
		req.Host = targetServer.Host
		req.URL.Host = targetServer.Host
		req.URL.Scheme = targetServer.Scheme
		req.RequestURI = ""

		//clear header
		for k, _ := range req.Header {
			req.Header.Del(k)
		}

		//set header
		for k, v := range reqData.Headers {
			req.Header.Set(k,v)
		}

		// save the response from the target server
		targetServerResponse, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf(err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer targetServerResponse.Body.Close()
		defer req.Body.Close()

		//combine response data
		resData := newResponseData()
		resData.Id = uuid.NewString()
		resData.Status = targetServerResponse.Status

		numOfBytes, err := io.Copy(ioutil.Discard, targetServerResponse.Body)
		if err != nil {
			log.Printf(err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		resData.Lenght = numOfBytes

		//copy header info from target server
		for k, v := range targetServerResponse.Header {
			resData.Headers[k] = strings.Join(v, ", ")
		}

		resDataJson, err := json.Marshal(resData)
		if err != nil {
			log.Printf(err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		//return response
		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(resDataJson)
	})

	fmt.Printf("starting proxy :8080")
	log.Fatal(http.ListenAndServe(":8080", reverseProxy))
}

//for the post method 'body' will be needed
type requestData struct {
	Method  string
	Url     string
	Headers map[string]string
}

type responseData struct {
	Id      string
	Status  string
	Headers map[string]string
	Lenght  int64
}

func newResponseData() *responseData {
	rv := new(responseData)
	rv.Headers = make(map[string]string, 10)

	return rv
}

func getRequestData(req *http.Request) (*requestData, error) {
	buff, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	req.Body = ioutil.NopCloser(bytes.NewBuffer(buff))

	reqData := new(requestData)

	err = json.Unmarshal(buff, reqData)
	if err != nil {
		return nil, err
	}

	return reqData, nil
}
