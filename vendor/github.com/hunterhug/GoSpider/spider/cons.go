/*
Copyright 2017 by GoSpider author. Email: gdccmcm14@live.com
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package spider

import "net/http"

const (
	// Default wait time
	WaitTime = 5

	// HTTP method
	GET      = "GET"
	POST     = "POST"
	POSTJSON = "POSTJSON"
	POSTXML  = "POSTXML"
	POSTFILE = "POSTFILE"
	PUT      = "PUT"
	PUTJSON  = "PUTJSON"
	PUTXML   = "PUTXML"
	PUTFILE  = "PUTFILE"
	DELETE   = "DELETE"
	OTHER    = "OTHER" // this stand for you can use other method this lib not own.

	// HTTP content type
	HTTPFORMContentType = "application/x-www-form-urlencoded"
	HTTPJSONContentType = "application/json"
	HTTPXMLContentType  = "text/xml"
	HTTPFILEContentType = "multipart/form-data"

	// Log mark
	CRITICAL = "CRITICAL"
	ERROR    = "ERROR"
	WARNING  = "WARNING"
	NOTICE   = "NOTICE"
	INFO     = "INFO"
	DEBUG    = "DEBUG"
)

var (
	// Browser User-Agent, Our default Http ua header!
	ourloveUa = "GolangSpider+(+http://cjhug.me+LoveYou/v2)"

	DefaultHeader = map[string][]string{
		"User-Agent": {
			ourloveUa,
		},
	}

	// DefaultTimeOut,http get and post No timeout
	DefaultTimeOut = 0
)

// Set global timeout, it can only by this way!
func SetGlobalTimeout(num int) {
	DefaultTimeOut = num
}

// Merge Cookie, not use
func MergeCookie(before []*http.Cookie, after []*http.Cookie) []*http.Cookie {
	cs := make(map[string]*http.Cookie)

	for _, b := range before {
		cs[b.Name] = b
	}

	for _, a := range after {
		if a.Value != "" {
			cs[a.Name] = a
		}
	}

	res := make([]*http.Cookie, 0, len(cs))

	for _, q := range cs {
		res = append(res, q)

	}

	return res

}

// Clone a header, If Not Ua, Set our Ua!
func CloneHeader(h map[string][]string) map[string][]string {
	if h == nil || len(h) == 0 {
		h = DefaultHeader
		return h
		//return map[string][]string{}
	}

	if len(h["User-Agent"]) == 0 {
		h["User-Agent"] = []string{ourloveUa}
	}
	return CopyM(h)
}
