// Copyright 2015 James Cote and Liberty Fund, Inc.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// docusign implements a service to use the version 2 Docusign
// rest api. Api documentation may be found at:
// https://www.docusign.com/p/RESTAPIGuide/RESTAPIGuide.htm
package docusign

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"time"

	"golang.org/x/net/context"
)

const (
	Version   = "0.5"
	userAgent = "docusign-api-go-client/" + Version
	liveUrl   = "https://www.docusign.net/restapi/v2/"
	testUrl   = "https://demo.docusign.net/restapi/v2/"
)

var (
	baseUrl string = testUrl

	// LogRawResponse will copy all api responses to log, if true
	LogRawResponse bool = false
	LogRawRequest  bool = false
)

// DSBool is used to fix problem of capitalized DSBooleans in json. Unmarshals
// "True" and "true" as true, any other value returns false
type DSBool bool

func (d *DSBool) UnmarshalJSON(b []byte) error {
	*d = DSBool(b[0] == 0x22 && (b[1] == 0x54 || b[1] == 0x74))
	return nil
}

// Credential add an authorization header(s) for a rest http request
type Credential interface {
	Authorize(*http.Request)
}

// OauthCredential provides authorization for rest request via
// docusign's oauth protocol
//
// Documentation: https://www.docusign.com/p/RESTAPIGuide/RESTAPIGuide.htm#OAuth2/OAuth2 Authentication Support in DocuSign REST API.htm
type OauthCredential struct {
	AccessToken string `json:"access_token,omitempty"`
	Scope       string `json:"scope,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
	// Used to set X-DocuSign-Act-As-User header if non-empty
	OnBehalfOfUser string `json:"onBehalf,omitempty"`
	AccountId      string `json:"account_id,omitempty"`
}

// Authorize update request with authorization parameters
func (o OauthCredential) Authorize(req *http.Request) {
	req.Header.Set("Authorization", "bearer "+o.AccessToken)
	if len(o.OnBehalfOfUser) > 0 {
		req.Header.Set("X-DocuSign-Act-As-User", o.OnBehalfOfUser)
	}
	return
}

// Revoke invalidates the token ensuring that an error will occur on an subsequent uses.
func (o OauthCredential) Revoke(ctx context.Context) error {
	client, err := contextClient(ctx)
	if err != nil {
		return err
	}
	v := url.Values{
		"token": {o.AccessToken},
	}
	req, err := http.NewRequest("POST", baseUrl+"oauth2/revoke", bytes.NewBufferString(v.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	return checkResponseStatus(res)
}

// Config provides methods to authenticate via a user/password combination.  It may also
// be used to generate an OauthCredential.
//
// Documentation:  https://www.docusign.com/p/RESTAPIGuide/RESTAPIGuide.htm#SOBO/Send On Behalf Of Functionality in the DocuSign REST API.htm
type Config struct {
	IntegratorKey  string `json:"key"`
	UserName       string `json:"user"`
	Password       string `json:"pwd"`
	OnBehalfOfUser string `json:"behalfOf,omitempty"`
	AccountId      string `json:"acctId,omitempty"`
}

// OauthCredential retrieves an OauthCredential  from docusign
// using the username and password from Config. The returned
// token does not have a expiration although it may be revoked
// via
func (c *Config) OauthCredential(ctx context.Context) (*OauthCredential, error) {
	client, err := contextClient(ctx)
	if err != nil {
		return nil, err
	}

	v := url.Values{
		"grant_type": []string{"password"},
		"client_id":  []string{c.IntegratorKey},
		"username":   []string{c.UserName},
		"password":   []string{c.Password},
		"scope":      []string{"api"},
	}
	req, err := http.NewRequest("POST", baseUrl+"oauth2/token", bytes.NewBufferString(v.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if client == nil {
		client = http.DefaultClient
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if err = checkResponseStatus(res); err != nil {
		return nil, err
	}
	var tk *OauthCredential
	err = json.NewDecoder(res.Body).Decode(&tk)
	return tk, err
}

// Authorize adds authorization headers to a rest request using user/password functionality.
func (c Config) Authorize(req *http.Request) {
	var onBehalfOf string
	if len(c.OnBehalfOfUser) > 0 {
		onBehalfOf = "<SendOnBehalfOf>" + c.OnBehalfOfUser + "</SendOnBehalfOf>"
	}
	authString := "<DocuSignCredentials>" + onBehalfOf +
		"<Username>" + c.UserName + "</Username><Password>" +
		c.Password + "</Password><IntegratorKey>" +
		c.IntegratorKey + "</IntegratorKey></DocuSignCredentials>"
	req.Header.Set("X-DocuSign-Authentication", authString)
	return
}

// Service contains all rest methods and stores authorization
type Service struct {
	accountId  string // Docusign account id
	baseUrl    string
	credential Credential
	endpoint   string
	client     *http.Client
	ctx        context.Context
}

// New intializes a new rest api service.  If client is nil then
// http.DefaultClient is assumed.
func New(accountId string, credential Credential) *Service {
	var ctx context.Context
	endpoint, ok := ctx.Value(APIEndpoint).(string)
	if !ok {
		endpoint = liveUrl
	}
	cl, _ := contextClient(ctx)
	return &Service{
		//	ctx:        ctx,
		client:     cl,
		accountId:  fmt.Sprintf("accounts/%s/", accountId),
		credential: credential,
		endpoint:   endpoint,
		baseUrl:    fmt.Sprintf("%s/accounts/%s/", baseUrl, accountId),
	}
}

// UseDemoServer changes tht DefaultCtx so that calls are made to Docusign's demo server.
func UseDemoServer() context.Context {
	DefaultCtx = context.WithValue(context.Background(), APIEndpoint, testUrl)
	return DefaultCtx
}

func UseProductionServer() {
	DefaultCtx = context.Background()
}

// Upload file describes an a document attachment for uploading
type UploadFile struct {
	// mime type of content
	ContentType string
	// file name to display in envelope
	FileName string
	// envelope documentId
	Id string
	// document order for envelope
	Order string
	// reader for creating file
	Data io.Reader
}

// NmVal is a generic name value pair struct used throughout the api.
type NmVal struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// ResponseError is generated when docusign returns an http error.
//
// Documentation: https://www.docusign.com/p/RESTAPIGuide/RESTAPIGuide.htm#Error Code/Error Code Information.htm
type ResponseError struct {
	Err         string `json:"errorCode,omitempty"`
	Description string `json:"message,omitempty"`
	Status      int    `json:"-"`
}

// UnmarshalJSON allows different versions of response error to be unmarshalled.
func (r *ResponseError) UnmarshalJSON(b []byte) error {
	t := make(map[string]string)
	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}
	for k, v := range t {
		if k == "errorCode" || k == "error" {
			r.Err = v
		} else if k == "message" || k == "error_description" {
			r.Description = v
		}
	}
	return nil
}

func (r ResponseError) Error() string {
	return fmt.Sprintf("Status: %d  %s: %s", r.Status, r.Err, r.Description)
}

// CheckErr looks at the response for a 200 or 201.  If not it will
// decode the json into a Response Error.  Returns nil on  success.
// https://www.docusign.com/p/RESTAPIGuide/RESTAPIGuide.htm#Error Code/Error Code Information.htm
/*func CheckResponse(res *http.Response) (err error) {
	if res.StatusCode != 200 && res.StatusCode != 201 {
		defer res.Body.Close()
		re := &ResponseError{Status: res.StatusCode}
		bx := bytes.NewBuffer(make([]byte, 0))
		io.Copy(bx, res.Body)
		b2 := bytes.NewBuffer(bx.Bytes())
		log.Printf("Error Msg: %d %s", res.ContentLength, string(bx.Bytes()))

		err = json.NewDecoder(b2).Decode(&re)
		if err == nil {
			err = re
		}
	}
	return err
}*/

// DsQueryTimeFormat returns a string in the correct format for a querystring format
func DsQueryTimeFormat(t time.Time) string {
	return t.Format("01/02/2006 15:04")
}

// createQueryString converts a slice of NmVal into a url encoded querystring.
func createQueryString(args []NmVal) string {
	if len(args) > 0 {

		q := url.Values{}
		for _, nv := range args {
			if _, ok := q[nv.Name]; !ok {
				q[nv.Name] = []string{nv.Value}
			} else {
				q[nv.Name] = append(q[nv.Name], nv.Value)
			}
		}
		return "?" + q.Encode()
	}
	return ""
}

// checkResponseStatus looks at the response for a 200 or 201.  If not it will
// decode the json into a Response Error.  Returns nil on  success.
// https://www.docusign.com/p/RESTAPIGuide/RESTAPIGuide.htm#Error Code/Error Code Information.htm
func checkResponseStatus(res *http.Response) (err error) {
	if res.StatusCode != 200 && res.StatusCode != 201 {
		re := &ResponseError{Status: res.StatusCode}
		if res.ContentLength > 0 {
			defer res.Body.Close()
			err = json.NewDecoder(res.Body).Decode(re)
			if err != nil {
				re.Description = err.Error()
			}
		}
		err = re
	}
	return
}

// doPdf writes the pdf stream to outputWriter
func (s *Service) doPdf(ctx context.Context, outputWriter io.Writer, method string, urlStr string, payload interface{}) error {
	var body io.Reader = nil

	client, err := contextClient(ctx)
	if err != nil {
		return err
	}

	if payload != nil {
		var b []byte
		if LogRawRequest {
			if b, err = json.MarshalIndent(payload, "", "    "); err == nil {
				log.Printf("Request Body: %s", string(b))
			}
		} else {
			b, err = json.Marshal(payload)
		}
		if err != nil {
			return err
		}
		body = bytes.NewReader(b)
	}

	bUrl, ok := ctx.Value(APIEndpoint).(string)
	if !ok {
		bUrl = liveUrl
	}
	req, err := http.NewRequest(method, bUrl+s.accountId+urlStr, body)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Accept", "application/pdf")
	if payload != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	s.credential.Authorize(req)
	resCh := make(chan error)
	go func() {
		res, err := client.Do(req)
		if err == nil {
			if err = checkResponseStatus(res); err == nil {
				defer res.Body.Close()
				_, err = io.Copy(outputWriter, res.Body)
			}
		}
		resCh <- err
	}()

	select {
	case <-ctx.Done():
		err = ctx.Err()
		if t, ok := client.Transport.(canceler); ok {
			t.CancelRequest(req)
		}
	case err = <-resCh:
	}
	return err
}

// do returns the json response from a rest api call
func (s *Service) do(ctx context.Context, method string, urlStr string, payload interface{}, returnValue interface{}, files ...*UploadFile) error {

	var body io.Reader = nil
	var contentType string
	client, err := contextClient(ctx)
	if err != nil {
		return err
	}
	if len(files) > 0 {
		body, contentType = multiBody(payload, files)
	} else if payload != nil {
		// Prepare body
		var b []byte
		if LogRawRequest {
			if b, err = json.MarshalIndent(payload, "", "    "); err == nil {
				log.Printf("Request Body: %s", string(b))
			}
		} else {
			b, err = json.Marshal(payload)
		}
		if err == nil {
			body = bytes.NewReader(b)
			contentType = "application/json"
		}
	}
	if err != nil {
		return err
	}

	bUrl, ok := ctx.Value(APIEndpoint).(string)
	if !ok {
		bUrl = liveUrl
	}
	req, err := http.NewRequest(method, bUrl+s.accountId+urlStr, body)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", userAgent)
	if len(files) == 0 {
		req.Header.Add("Accept", "application/json")
	}
	if len(contentType) > 0 {
		req.Header.Set("Content-Type", contentType)
	}
	s.credential.Authorize(req)
	if LogRawRequest {
		log.Printf("RequestURL: %s", req.URL.String())
		for k, v := range req.Header {
			log.Printf("%s: %v\n", k, v)
		}

	}
	resCh := make(chan error)
	go func() {
		res, err := client.Do(req)
		if err == nil {
			if err = checkResponseStatus(res); err == nil {
				defer res.Body.Close()
				if LogRawResponse {
					var b []byte
					if b, err = ioutil.ReadAll(res.Body); err == nil {
						log.Printf("%s", string(b))
						err = json.Unmarshal(b, returnValue)
					}
				} else {
					err = json.NewDecoder(res.Body).Decode(returnValue)
				}
			}
		}
		resCh <- err
		close(resCh)
	}()

	select {
	case <-ctx.Done():
		err = ctx.Err()
		if t, ok := client.Transport.(canceler); ok {
			t.CancelRequest(req)
		}
	case err = <-resCh:
	}
	return err
}

type canceler interface {
	CancelRequest(*http.Request)
}

func multiBody(payload interface{}, files []*UploadFile) (io.Reader, string) {
	pr, pw := io.Pipe()
	mpw := multipart.NewWriter(pw)
	// write json payload first
	go func() {
		var err error
		var ptw io.Writer
		var fcnt int
		defer func() {
			if err != nil {
				pr.CloseWithError(fmt.Errorf("batch: multiPart Error: %v", err))
			}
			for fcnt < len(files) {
				if closer, ok := files[fcnt].Data.(io.Closer); ok {
					closer.Close()
				}
				fcnt++
			}
			mpw.Close()
			pw.Close()
		}()
		if payload != nil {
			mh := textproto.MIMEHeader{
				"Content-Type":        []string{"application/json"},
				"Content-Disposition": []string{"form-data"},
			}
			if ptw, err = mpw.CreatePart(mh); err != nil {
				return
			}
			if err = json.NewEncoder(ptw).Encode(payload); err != nil {
				return
			}
		}
		for _, f := range files {
			mh := textproto.MIMEHeader{
				"Content-Type":        []string{f.ContentType},
				"Content-Disposition": []string{fmt.Sprintf("file; filename=\"%s\";documentid=%s", f.FileName, f.Id)},
			}
			if ptw, err = mpw.CreatePart(mh); err == nil {
				_, err = io.Copy(ptw, f.Data)
				if closer, ok := f.Data.(io.Closer); ok {
					closer.Close()
				}
				fcnt++
			}
			if err != nil {
				return
			}
		}
	}()
	return pr, "multipart/form-data; boundary=" + mpw.Boundary()
}
