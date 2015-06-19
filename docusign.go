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
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/context"
)

const (
	Version   = "0.5"
	userAgent = "docusign-api-go-client/" + Version
	liveUrl   = "https://www.docusign.net/restapi/v2"
	testUrl   = "https://demo.docusign.net/restapi/v2"
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
	// Authorize attaches an authorization header to a request.
	Authorize(*http.Request)
	// AccountPath returns the url path fragment for the credenials account.
	AccountPath() string
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
	// The docusign account used by the login user.  This may be
	// found using the LoginInformation call.
	AccountId string `json:"account_id,omitempty"`
}

// Authorize update request with authorization parameters
func (o OauthCredential) Authorize(req *http.Request) {
	req.Header.Set("Authorization", "bearer "+o.AccessToken)
	if len(o.OnBehalfOfUser) > 0 {
		req.Header.Set("X-DocuSign-Act-As-User", o.OnBehalfOfUser)
	}
	return
}

// AccountPath returns the path fragment for the account
// associated with the credential.
func (o OauthCredential) AccountPath() string {
	return "/accounts/" + o.AccountId + "/"
}

// Revoke invalidates the token ensuring that an error will occur on an subsequent uses.
func (o OauthCredential) Revoke(ctx context.Context) error {
	settings := contextSettings(ctx)
	v := url.Values{
		"token": {o.AccessToken},
	}
	req, err := http.NewRequest("POST", settings.Endpoint+"/oauth2/revoke", bytes.NewBufferString(v.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := settings.Client.Do(req)
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
	//client := contextClient(ctx)
	settings := contextSettings(ctx)
	v := url.Values{
		"grant_type": []string{"password"},
		"client_id":  []string{c.IntegratorKey},
		"username":   []string{c.UserName},
		"password":   []string{c.Password},
		"scope":      []string{"api"},
	}
	req, err := http.NewRequest("POST", settings.Endpoint+"/oauth2/token", bytes.NewBufferString(v.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := settings.Client.Do(req)
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

// OauthCredentialOnBehalfOf returns an *OauthCredential for the user name specied by nm.  oauthCred
// must be a credential for a user with administrative rights on the account.
func (c *Config) OauthCredentialOnBehalfOf(ctx context.Context, oauthCred OauthCredential, nm string) (*OauthCredential, error) {
	settings := contextSettings(ctx)
	v := url.Values{
		"grant_type": []string{"password"},
		"client_id":  []string{c.IntegratorKey},
		"username":   []string{nm},
		"scope":      []string{"api"},
	}
	req, err := http.NewRequest("POST", settings.Endpoint+"/oauth2/token", bytes.NewBufferString(v.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	oauthCred.Authorize(req)

	res, err := settings.Client.Do(req)
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

// Account path returns the account fragment
func (c Config) AccountPath() string {
	return "/accounts/" + c.AccountId + "/"
}

// Service contains all rest methods and stores authorization
type Service struct {
	accountId  string // Docusign account id
	endpoint   string
	credential Credential
	client     *http.Client
	ctx        context.Context
	// logRawRequest and logRawResponse are used to capture the
	// json serialization for a call.  The string argument will
	// be a printf style string.
	logRawRequest  func(context.Context, string, ...interface{})
	logRawResponse func(context.Context, string, ...interface{})
}

// New intializes a new rest api service.  If client is nil then
// http.DefaultClient is assumed.
//func New(ctx context.Context, accountId string, credential Credential) *Service {
func New(ctx context.Context, credential Credential) *Service {
	settings := contextSettings(ctx)
	return &Service{
		ctx:            ctx,
		client:         settings.Client,
		credential:     credential,
		endpoint:       settings.Endpoint, //fmt.Sprintf("%s/accounts/%s/", settings.Endpoint, accountId),
		logRawRequest:  settings.LogRawRequest,
		logRawResponse: settings.LogRawResponse,
	}
}

func (s Service) newRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	var u string
	if strings.HasPrefix(urlStr, "/") {
		u = s.endpoint + urlStr
	} else {
		u = s.endpoint + s.credential.AccountPath() + urlStr
	}
	return http.NewRequest(method, u, body)
}

// UseDemoServer returns a Context that ensures that calls are made to Docusign's demo server.
func UseDemoServer(logRequest, logResponse func(context.Context, string, ...interface{})) context.Context {
	demoSettings := &ContextSetting{
		Endpoint:       testUrl,
		LogRawRequest:  logRequest,
		LogRawResponse: logResponse,
	}
	return context.WithValue(context.Background(), APISettings, demoSettings)
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
	Name  string `json:"name,omitempty" xml:"name,attr"`
	Value string `json:"value,omitempty" xml:",chardata"`
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

// doPdf writes the pdf stream to the provided outputWriter.
func (s *Service) doPdf(outputWriter io.Writer, method string, urlStr string, payload interface{}) error {
	var body io.Reader = nil
	var err error

	if payload != nil {
		var b []byte
		if s.logRawRequest != nil {
			if b, err = json.MarshalIndent(payload, "", "    "); err == nil {
				s.logRawRequest(s.ctx, "Request Body: %s", string(b))
			}
		} else {
			b, err = json.Marshal(payload)
		}
		if err != nil {
			return err
		}
		body = bytes.NewReader(b)
	}

	req, err := s.newRequest(method, urlStr, body) //http.NewRequest(method, s.endpoint+urlStr, body)
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
		res, err := s.client.Do(req)
		if err == nil {
			if err = checkResponseStatus(res); err == nil {
				defer res.Body.Close()
				_, err = io.Copy(outputWriter, res.Body)
			}
		}
		resCh <- err
	}()

	select {
	case <-s.ctx.Done():
		err = s.ctx.Err()
		if t, ok := s.client.Transport.(canceler); ok {
			t.CancelRequest(req)
		}
	case err = <-resCh:
	}
	return err
}

// do returns the json response from a rest api call
func (s *Service) do(method string, urlStr string, payload interface{}, returnValue interface{}, files ...*UploadFile) error {

	var body io.Reader = nil
	var contentType string
	var err error
	if len(files) > 0 {
		body, contentType = multiBody(payload, files)
	} else if payload != nil {
		// Prepare body
		var b []byte
		if s.logRawRequest != nil {
			if b, err = json.MarshalIndent(payload, "", "    "); err == nil {
				s.logRawRequest(s.ctx, "Request Body: %s", string(b))
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

	req, err := s.newRequest(method, urlStr, body) //http.NewRequest(method, s.endpoint+urlStr, body)
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
	if s.logRawRequest != nil {
		s.logRawRequest(s.ctx, "RequestURL: %s", req.URL.String())
		for k, v := range req.Header {
			s.logRawRequest(s.ctx, "%s: %v\n", k, v)
		}

	}
	resCh := make(chan error)
	go func() {
		res, err := s.client.Do(req)
		if err == nil {
			if err = checkResponseStatus(res); err == nil {
				defer res.Body.Close()
				if s.logRawResponse != nil {
					var b []byte
					if b, err = ioutil.ReadAll(res.Body); err == nil {
						s.logRawResponse(s.ctx, "%s", string(b))
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
	case <-s.ctx.Done():
		err = s.ctx.Err()
		if t, ok := s.client.Transport.(canceler); ok {
			t.CancelRequest(req)
		}
	case err = <-resCh:
	}
	return err
}

type canceler interface {
	CancelRequest(*http.Request)
}

// multiBody is used to format calls containing files.
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
