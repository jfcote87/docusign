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
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

const (
	Version   = "0.5"
	userAgent = "docusign-api-go-client/" + Version
)

/*  Documentation: https://docs.docusign.com/esign/

All REST API endpoints have the following base:
	https://{server}.docusign.net/restapi/v2

DocuSign hosts multiple geo-dispersed ISO 27001-certified and SSAE 16-audited data centers. For example, account holders in North America might have the following baseUrl:
	https://na2.docusign.net/restapi/v2/accounts/{accountId}
Whereas European users might access the following baseUrl:
	https://eu.docusign.net/restapi/v2/accounts/{accountId}

EXAMPLES
	"https://www.docusign.net/restapi/v2"  (deprecated?)
	"https://n2.docusign.net/restapi/v2"   (north america)
	"https://eu.docusign.net/restapi/v2"   (europe)
	"https://demo.docusign.net/restapi/v2" (sandbox)

*/
var (
	baseURL = &url.URL{Scheme: "https", User: (*url.Userinfo)(nil), Host: "", Path: "/restapi/v2"}
)

// DSBool is used to fix problem of capitalized DSBooleans in json. Unmarshals
// "True" and "true" as true, any other value returns false
type DSBool bool

func (d *DSBool) UnmarshalJSON(b []byte) error {
	*d = DSBool(b[0] == 0x22 && (b[1] == 0x54 || b[1] == 0x74))
	return nil
}

// dsResolveURL resolves a relative url.
// the host parameter determines which docusign server(s) to hit
//   EX: prod north america, prod europe, demo
// the accountID is used to finish the url's path.
func dsResolveURL(ref *url.URL, host string, accountID string) {
	baseURL.Host = host
	ref.Scheme = baseURL.Scheme
	ref.Host = baseURL.Host

	if strings.HasPrefix(ref.Path, "/") {
		ref.Path = baseURL.Path + ref.Path
	} else {
		ref.Path = baseURL.Path + "/accounts/" + accountID + "/" + ref.Path
	}
}

// Credential add an authorization header(s) for a rest http request
type Credential interface {
	// Authorize attaches an authorization header to a request and
	// and fixes the URL to the appropriate host.
	Authorize(*http.Request, string)
}

// OauthCredential provides authorization for rest request via
// docusign's oauth protocol
//
// Documentation: https://www.docusign.com/p/RESTAPIGuide/RESTAPIGuide.htm#OAuth2/OAuth2 Authentication Support in DocuSign REST API.htm
type OauthCredential struct {
	// The docusign account used by the login user.  This may be
	// found using the LoginInformation call.
	AccountId   string `json:"account_id,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
	Scope       string `json:"scope,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
	Host        string `json:"host,omitempty"`
}

// Authorize update request with authorization parameters
func (o OauthCredential) Authorize(req *http.Request, onBehalfOf string) {
	dsResolveURL(req.URL, o.Host, o.AccountId)

	var auth string
	if o.TokenType == "" {
		auth = "bearer " + o.AccessToken
	} else {
		auth = o.TokenType + " " + o.AccessToken
	}
	req.Header.Set("Authorization", auth)
	if onBehalfOf != "" {
		req.Header.Set("X-DocuSign-Act-As-User", onBehalfOf)
	}
	return
}

// Revoke invalidates the token ensuring that an error will occur on an subsequent uses.
func (o OauthCredential) Revoke(ctx context.Context) error {
	v := url.Values{
		"token": {o.AccessToken},
	}
	req, err := http.NewRequest("POST", "/oauth2/revoke", bytes.NewBufferString(v.Encode()))
	if err != nil {
		return err
	}

	dsResolveURL(req.URL, o.Host, o.AccountId)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := ctxhttp.Do(ctx, contextClient(ctx), req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return checkResponseStatus(res)
}

// Config provides methods to authenticate via a user/password combination.
// It may also be used to generate an OauthCredential.
// Documentation:  https://www.docusign.com/p/RESTAPIGuide/RESTAPIGuide.htm#SOBO/Send On Behalf Of Functionality in the DocuSign REST API.htm
type Config struct {
	// The docusign account used by the login user.  This may be
	// found using the LoginInformation call.
	AccountId     string `json:"acctId,omitempty"`
	IntegratorKey string `json:"key"`
	UserName      string `json:"user"`
	Password      string `json:"pwd"`
	Host          string `json:"host,omitempty"`
}

// OauthCredential retrieves an OauthCredential  from docusign
// using the username and password from Config. The returned
// token does not have a expiration although it may be revoked
// via
func (c *Config) OauthCredential(ctx context.Context) (*OauthCredential, error) {
	v := url.Values{
		"grant_type": []string{"password"},
		"client_id":  []string{c.IntegratorKey},
		"username":   []string{c.UserName},
		"password":   []string{c.Password},
		"scope":      []string{"api"},
	}
	req, err := http.NewRequest("POST", "/oauth2/token", bytes.NewBufferString(v.Encode()))
	if err != nil {
		return nil, err
	}
	dsResolveURL(req.URL, c.Host, c.AccountId)

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := ctxhttp.Do(ctx, contextClient(ctx), req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if err = checkResponseStatus(res); err != nil {
		return nil, err
	}
	var tk *OauthCredential
	if err = json.NewDecoder(res.Body).Decode(&tk); err == nil {
		tk.Host = c.Host
		tk.AccountId = c.AccountId
	}
	return tk, err
}

// OauthCredentialOnBehalfOf returns an *OauthCredential for the user name specied by nm.  oauthCred
// must be a credential for a user with administrative rights on the account.
func (c *Config) OauthCredentialOnBehalfOf(ctx context.Context, oauthCred OauthCredential, nm string) (*OauthCredential, error) {
	v := url.Values{
		"grant_type": []string{"password"},
		"client_id":  []string{c.IntegratorKey},
		"username":   []string{nm},
		"scope":      []string{"api"},
	}
	req, err := http.NewRequest("POST", "/oauth2/token", bytes.NewBufferString(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	oauthCred.Authorize(req, "")

	res, err := ctxhttp.Do(ctx, contextClient(ctx), req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if err = checkResponseStatus(res); err != nil {
		return nil, err
	}
	var tk *OauthCredential
	if err = json.NewDecoder(res.Body).Decode(&tk); err == nil {
		tk.Host = c.Host
		tk.AccountId = c.AccountId
	}
	return tk, err
}

// Authorize adds authorization headers to a rest request using user/password functionality.
func (c Config) Authorize(req *http.Request, onBehalfOf string) {
	dsResolveURL(req.URL, c.Host, c.AccountId)
	if onBehalfOf != "" {
		onBehalfOf = "<SendOnBehalfOf>" + onBehalfOf + "</SendOnBehalfOf>"
	}
	authString := "<DocuSignCredentials>" + onBehalfOf +
		"<Username>" + c.UserName + "</Username><Password>" +
		c.Password + "</Password><IntegratorKey>" +
		c.IntegratorKey + "</IntegratorKey></DocuSignCredentials>"
	req.Header.Set("X-DocuSign-Authentication", authString)
	return
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

// checkResponseStatus looks at the response for a 200 or 201.  If not it will
// decode the json into a Response Error.  Returns nil on  success.
// https://www.docusign.com/p/RESTAPIGuide/RESTAPIGuide.htm#Error Code/Error Code Information.htm
func checkResponseStatus(res *http.Response) (err error) {
	if res.StatusCode != 200 && res.StatusCode != 201 {
		re := &ResponseError{Status: res.StatusCode}
		if res.ContentLength > 0 {
			err = json.NewDecoder(res.Body).Decode(re)
			if err != nil {
				re.Description = err.Error()
			}
		}
		err = re
	}
	return
}

// Service contains all rest methods and stores authorization
type Service struct {
	credential Credential
	onBehalfOf string
}

// New intializes a new rest api service.  If client is nil then
// http.DefaultClient is assumed.
//func New(ctx context.Context, accountId string, credential Credential) *Service {
func New(credential Credential, onBehalfOf string) *Service {
	return &Service{credential: credential, onBehalfOf: onBehalfOf}
}

// OnBehalfOf returns a new Service set to authenticate then
// onBehalfOf userId (email address).  The original Service
// credential must be an administrator.
func (s Service) OnBehalfOf(onBehalfOf string) *Service {
	s.onBehalfOf = onBehalfOf
	return &s
}

// Call provides all needed fields to make a call.  To debug
// a call simply set the Result to an **http.Response.
type Call struct {
	Method string

	Payload interface{}
	// Result may be either
	Result interface{}
	// uploaded files for a call
	Files []*UploadFile
	// relative url for the call
	URL *url.URL
}

// Do executes the call.  Response data is encoded into
// the call's Result.  If Result is a **http.Response, the
// response is returned without processing.
func (c Call) Do(ctx context.Context, s *Service) error {
	var body io.Reader
	var ct string
	var raw **http.Response
	logger := contextLogger(ctx)

	if len(c.Files) > 0 {
		// formatted body for file upload
		body, ct = multiBody(c.Payload, c.Files)
	} else if c.Payload != nil {
		// Prepare body
		b, err := json.Marshal(c.Payload)
		if err != nil {
			return err
		}
		body, ct = bytes.NewReader(b), "application/json"
	}

	req, err := http.NewRequest(c.Method, "", body)
	if err != nil {
		return err
	}
	req.URL = c.URL
	s.credential.Authorize(req, s.onBehalfOf)
	req.Header.Add("User-Agent", userAgent)

	if len(ct) > 0 {
		req.Header.Set("Content-Type", ct)
	}
	if c.Result != nil {
		if raw, _ = c.Result.(**http.Response); raw == nil {
			req.Header.Set("accept", "application/json")
		}
	}

	if logger != nil {
		logger.LogRequest(ctx, c.Payload, req)
	}

	res, err := ctxhttp.Do(ctx, contextClient(ctx), req)
	if err != nil {
		return err
	}
	if err = checkResponseStatus(res); err != nil {
		res.Body.Close()
		return err
	}
	if raw != nil {
		*raw = res
		return nil
	}

	defer res.Body.Close()
	body = res.Body
	if logger != nil {
		body = logger.LogResponse(ctx, res)
	}
	if c.Result != nil {
		err = json.NewDecoder(body).Decode(c.Result)
	}
	return err
}

// multiBody is used to format calls containing files as a multipart/form-data body.
//
func multiBody(payload interface{}, files []*UploadFile) (io.Reader, string) {
	pr, pw := io.Pipe()
	mpw := multipart.NewWriter(pw)

	go func() {
		var err error
		var ptw io.Writer
		defer func() {
			if err != nil {
				pr.CloseWithError(fmt.Errorf("batch: multiPart Error: %v", err))
			}
			// Close input files
			for _, f := range files {
				if closer, ok := f.Data.(io.Closer); ok {
					closer.Close()
				}
			}
			mpw.Close()
			pw.Close()
		}()
		// write json payload first
		if payload != nil {
			mh := textproto.MIMEHeader{
				"Content-Type":        []string{"application/json"},
				"Content-Disposition": []string{"form-data"},
			}
			if ptw, err = mpw.CreatePart(mh); err == nil {
				err = json.NewEncoder(ptw).Encode(payload)
			}
			if err != nil {
				return
			}
		}

		for _, f := range files {
			mh := textproto.MIMEHeader{
				"Content-Type":        []string{f.ContentType},
				"Content-Disposition": []string{fmt.Sprintf("file; filename=\"%s\";documentid=%s", f.FileName, f.Id)},
			}
			if ptw, err = mpw.CreatePart(mh); err == nil {
				if _, err = io.Copy(ptw, f.Data); err != nil {
					break
				}
			}
		}
		return
	}()
	return pr, "multipart/form-data; boundary=" + mpw.Boundary()
}
