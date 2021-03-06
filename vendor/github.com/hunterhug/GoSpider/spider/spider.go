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

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/hunterhug/GoTool/util"
)

// New a spider, if ipstring is a proxy address, New a proxy client.
// Proxy address such as:
// 		http://[user]:[password@]ip:port, [] stand it can choose or not
// 		socks5://127.0.0.1:1080
func NewSpider(ipstring interface{}) (*Spider, error) {
	sp := new(Spider)
	sp.Header = http.Header{}
	sp.Data = url.Values{}
	sp.BData = []byte{}
	if ipstring != nil {
		client, err := NewProxyClient(strings.ToLower(ipstring.(string)))
		sp.Client = client
		sp.Ipstring = ipstring.(string)
		return sp, err
	} else {
		client, err := NewClient()
		sp.Client = client
		sp.Ipstring = "localhost"
		return sp, err
	}

}

// Alias Name for NewSpider
func New(ipstring interface{}) (*Spider, error) {
	return NewSpider(ipstring)
}

// New Spider by Your Client
func NewSpiderByClient(client *http.Client) *Spider {
	sp := new(Spider)
	sp.Header = http.Header{}
	sp.Data = url.Values{}
	sp.BData = []byte{}
	sp.Client = client
	return sp
}

// New API Spider, No Cookie Keep.
func NewAPI() *Spider {
	return NewSpiderByClient(NoCookieClient)
}

// Auto decide which method, Default Get.
func (sp *Spider) Go() (body []byte, e error) {
	switch strings.ToUpper(sp.Method) {
	case POST:
		return sp.Post()
	case POSTJSON:
		return sp.PostJSON()
	case POSTXML:
		return sp.PostXML()
	case POSTFILE:
		return sp.PostFILE()
	case PUT:
		return sp.Put()
	case PUTJSON:
		return sp.PutJSON()
	case PUTXML:
		return sp.PutXML()
	case PUTFILE:
		return sp.PutFILE()
	case DELETE:
		return sp.Delete()
	case OTHER:
		return []byte(""), errors.New("please use method OtherGo(method, content type)")
	default:
		return sp.Get()
	}
}

func (sp *Spider) GoByMethod(method string) (body []byte, e error) {
	return sp.SetMethod(method).Go()
}

// This make effect only your spider exec serial! Attention!
// Change Your Raw data To string
func (sp *Spider) ToString() string {
	if sp.Raw == nil {
		return ""
	}
	return string(sp.Raw)
}

// This make effect only your spider exec serial! Attention!
// Change Your JSON'like Raw data to string
func (sp *Spider) JsonToString() (string, error) {
	if sp.Raw == nil {
		return "", nil
	}
	temp, err := util.JsonBack(sp.Raw)
	if err != nil {
		return "", err
	}
	return string(temp), nil
}

// Main method I make!
func (sp *Spider) sent(method, contenttype string, binary bool) (body []byte, e error) {
	// Lock it for save
	sp.mux.Lock()
	defer sp.mux.Unlock()

	// Before FAction we can change or add something before Go()
	if sp.BeforeAction != nil {
		sp.BeforeAction(sp.Ctx, sp)
	}

	// Wait if must
	if sp.Wait > 0 {
		Wait(sp.Wait)
	}

	// For debug
	Logger.Debugf("[GoSpider] %s %s", method, sp.Url)

	// New a Request
	var request = &http.Request{}

	// If binary parm value is true and BData is not empty
	// suit for POSTJSON(), POSTFILE()
	if len(sp.BData) != 0 && binary {
		pr := ioutil.NopCloser(bytes.NewReader(sp.BData))
		request, _ = http.NewRequest(method, sp.Url, pr)
	} else if len(sp.Data) != 0 { // such POST() from table form
		pr := ioutil.NopCloser(strings.NewReader(sp.Data.Encode()))
		request, _ = http.NewRequest(method, sp.Url, pr)
	} else {
		request, _ = http.NewRequest(method, sp.Url, nil)
	}

	// Clone Header, I add some HTTP header!
	request.Header = CloneHeader(sp.Header)

	// In fact contenttype must not empty
	if contenttype != "" {
		request.Header.Set("Content-Type", contenttype)
	}
	sp.Request = request

	// Debug for RequestHeader
	OutputMaps("Request header", request.Header)

	// Tolerate abnormal way to create a Spider
	if sp.Client == nil {
		sp.Client = Client
	}

	// Do it
	response, err := sp.Client.Do(request)
	if err != nil {
		// I count Error time
		sp.Errortimes++
		return nil, err
	}

	// Close it attention response may be nil
	if response != nil {
		defer response.Body.Close()
	}

	// Debug
	OutputMaps("Response header", response.Header)
	Logger.Debugf("[GoSpider] %v %s", response.Proto, response.Status)

	// Read output
	body, e = ioutil.ReadAll(response.Body)
	sp.Raw = body

	sp.UrlStatuscode = response.StatusCode
	sp.Preurl = sp.Url
	sp.Response = response
	sp.Fetchtimes++

	// After action
	if sp.AfterAction != nil {
		sp.AfterAction(sp.Ctx, sp)
	}
	return
}

// Get method
func (sp *Spider) Get() (body []byte, e error) {
	sp.Clear()
	return sp.sent(GET, "", false)
}

func (sp *Spider) Delete() (body []byte, e error) {
	sp.Clear()
	return sp.sent(DELETE, "", false)
}

// Post Almost include bellow:
/*
	"application/x-www-form-urlencoded"
	"application/json"
	"text/xml"
	"multipart/form-data"
*/
func (sp *Spider) Post() (body []byte, e error) {
	return sp.sent(POST, HTTPFORMContentType, false)
}

func (sp *Spider) PostJSON() (body []byte, e error) {
	return sp.sent(POST, HTTPJSONContentType, true)
}

func (sp *Spider) PostXML() (body []byte, e error) {
	return sp.sent(POST, HTTPXMLContentType, true)
}

func (sp *Spider) PostFILE() (body []byte, e error) {
	return sp.sent(POST, HTTPFILEContentType, true)

}

// Put
func (sp *Spider) Put() (body []byte, e error) {
	return sp.sent(PUT, HTTPFORMContentType, false)
}

func (sp *Spider) PutJSON() (body []byte, e error) {
	return sp.sent(PUT, HTTPJSONContentType, true)
}

func (sp *Spider) PutXML() (body []byte, e error) {
	return sp.sent(PUT, HTTPXMLContentType, true)
}

func (sp *Spider) PutFILE() (body []byte, e error) {
	return sp.sent(PUT, HTTPFILEContentType, true)

}

// Other Method
/*
     Method         = "OPTIONS"                ; Section 9.2
                    | "GET"                    ; Section 9.3
                    | "HEAD"                   ; Section 9.4
                    | "POST"                   ; Section 9.5
                    | "PUT"                    ; Section 9.6
                    | "DELETE"                 ; Section 9.7
                    | "TRACE"                  ; Section 9.8
                    | "CONNECT"                ; Section 9.9
                    | extension-method
   extension-method = token
     token          = 1*<any CHAR except CTLs or separators>


// content type
	"application/x-www-form-urlencoded"
	"application/json"
	"text/xml"
	"multipart/form-data"
*/
func (sp *Spider) OtherGo(method, contenttype string) (body []byte, e error) {
	return sp.sent(method, contenttype, true)
}
