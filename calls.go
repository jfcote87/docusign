// Copyright 2015 James Cote and Liberty Fund, Inc.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package docusign

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
)

// FolderList returns a list of the folders for the account, including the
// folder hierarchy.
// Optional query string: template={string} include or only
//
// RestApi documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get%20Folder%20List.htm
func (s *Service) FolderList(ctx context.Context, args ...FolderTemplateParam) (*FolderList, error) {
	var ret *FolderList
	q := make(url.Values)
	for _, nv := range args {
		q.Add(nv.Name, nv.Value)
	}
	return ret, (&Call{
		Method: "GET",
		URL:    &url.URL{Path: "folders", RawQuery: q.Encode()},
		Result: &ret,
	}).Do(ctx, s)
}

type FolderTemplateParam NmVal

var FolderTemplatesInclude = FolderTemplateParam{Name: "template", Value: "include"}
var FolderTemplatesOnly = FolderTemplateParam{Name: "template", Value: "only"}

// FolderEnvSearch returns a list of the envelopes in the specified folder. You can narrow
// the query by adding some optional items.
// additions: start_position={ integer}, from_date = {date/time}, to_date= {date/time},
// search_text={text}, status={envelope status}, owner_name={username}, owner_email={email}
//
// RestApi documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get%20Folder%20Envelope%20List.htm
func (s *Service) FolderEnvSearch(ctx context.Context, folderId string, args ...FolderEnvSearchParam) (*FolderEnvList, error) {
	var ret *FolderEnvList
	q := make(url.Values)
	for _, nv := range args {
		q.Add(nv.Name, nv.Value)
	}
	return ret, (&Call{
		Method: "GET",
		URL:    &url.URL{Path: fmt.Sprintf("folders/%s", folderId), RawQuery: q.Encode()},
		Result: &ret,
	}).Do(ctx, s)
}

type FolderEnvSearchParam NmVal

func FolderEnvSearchStartPosition(pos int) FolderEnvSearchParam {
	return FolderEnvSearchParam{Name: "startpostition", Value: strconv.Itoa(pos)}
}

func FolderEnvSearchFromDate(tm time.Time) FolderEnvSearchParam {
	return FolderEnvSearchParam{Name: "from_date", Value: DsQueryTimeFormat(tm)}
}

func FolderEnvSearchToDate(tm time.Time) FolderEnvSearchParam {
	return FolderEnvSearchParam{Name: "to_date", Value: DsQueryTimeFormat(tm)}
}

func FolderEnvSearchText(searchText string) FolderEnvSearchParam {
	return FolderEnvSearchParam{Name: "search_text", Value: searchText}
}

func FolderEnvSearchStatus(status string) FolderEnvSearchParam {
	return FolderEnvSearchParam{Name: "status", Value: status}
}

func FolderEnvSearchOwnerName(nm string) FolderEnvSearchParam {
	return FolderEnvSearchParam{Name: "owner_name", Value: nm}
}

func FolderEnvSearchOwnerEmail(email string) FolderEnvSearchParam {
	return FolderEnvSearchParam{Name: "owner_email", Value: email}
}

// EnvelopeSearch returns a list of envelopes that match the criteria specified in the query parameters.
// Query Parameters:
// start_position is the starting value for the list.
// count is the number of records to return in the cache. The number must be greater than 1 and less than or equal to 100.
// from_date is the start of the date range. If no value is provided, the default search is the previous 30 days.
// to_date is the end of the date range.
// order_by	Column used to sort the list. Valid values are listed EnvelopeSearchOrderBy* variables
// order sets the direction of the sort.  Valid values are EnvelopeSearchAsc or EnvelopeSearchDesc.
// include_recipients returns the recipient information when true.
//
// RestApi documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST API References/Get List of Envelopes in Folders.htm
func (s *Service) EnvelopeSearch(ctx context.Context, searchFld SearchFolder, args ...SearchFolderParam) (*FolderEnvList, error) {
	var ret *FolderEnvList
	q := make(url.Values)
	for _, nv := range args {
		q.Add(nv.Name, nv.Value)
	}
	return ret, (&Call{
		Method: "GET",
		URL:    &url.URL{Path: fmt.Sprintf("search_folders/%s", searchFld), RawQuery: q.Encode()},
		Result: &ret,
	}).Do(ctx, s)

}

// SearchFolder specifies strings describing the the type of available search folders
type SearchFolder string

var SearchFolderDrafts SearchFolder = "drafts"
var SearchFolderAwaitingSig SearchFolder = "awaiting_my_signature"
var SearchFolderOutForSig SearchFolder = "out_for_signature"
var SearchFolderCompleted SearchFolder = "completed"

type SearchFolderParam NmVal

func EnvelopeSearchStartPosition(pos int) SearchFolderParam {
	return SearchFolderParam{Name: "startpostition", Value: strconv.Itoa(pos)}
}

func EnvelopeSearchCount(cnt int) SearchFolderParam {
	return SearchFolderParam{Name: "count", Value: strconv.Itoa(cnt)}
}

func EnvelopeSearchFromDate(tm time.Time) SearchFolderParam {
	return SearchFolderParam{Name: "from_date", Value: DsQueryTimeFormat(tm)}
}

func EnvelopeSearchToDate(tm time.Time) SearchFolderParam {
	return SearchFolderParam{Name: "to_date", Value: DsQueryTimeFormat(tm)}
}

var EnvelopeSearchOrderByActionRequired = SearchFolderParam{
	Name:  "order_by",
	Value: "action_required",
}
var EnvelopeSearchOrderByCreated = SearchFolderParam{
	Name:  "order_by",
	Value: "created",
}
var EnvelopeSearchOrderByCompleted = SearchFolderParam{
	Name:  "order_by",
	Value: "completed",
}
var EnvelopeSearchOrderBySent = SearchFolderParam{
	Name:  "order_by",
	Value: "sent",
}
var EnvelopeSearchOrderBySignerList = SearchFolderParam{
	Name:  "order_by",
	Value: "signer_list",
}
var EnvelopeSearchOrderByStatus = SearchFolderParam{
	Name:  "order_by",
	Value: "status",
}
var EnvelopeSearchOrderBySubject = SearchFolderParam{
	Name:  "order_by",
	Value: "subject",
}

var EnvelopeSearchOrderAsc = SearchFolderParam{
	Name:  "order",
	Value: "asc",
}
var EnvelopeSearchOrderDesc = SearchFolderParam{
	Name:  "order",
	Value: "desc",
}

var EnvelopeSearchIncludeRecipients = SearchFolderParam{
	Name:  "include_recipients",
	Value: "true",
}

// EnvelopeAuditEvents returns the events for this envelope.
//
// RestApi documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get%20Envelope%20Audit%20Events.htm
func (s *Service) EnvelopeAuditEvents(ctx context.Context, envId string) (*AuditEventList, error) {
	var ret *AuditEventList
	return ret, (&Call{
		Method: "GET",
		URL:    &url.URL{Path: fmt.Sprintf("envelopes/%s/audit_events", envId)},
		Result: &ret,
	}).Do(ctx, s)

}

// EnvelopeNotification returns the reminder and expiration information for the envelope.
//
// RestApi documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get%20Envelope%20Notification%20Information.htm
func (s *Service) EnvelopeNotification(ctx context.Context, envId string) (*Notification, error) {
	var ret *Notification
	return ret, (&Call{
		Method: "GET",
		URL:    &url.URL{Path: fmt.Sprintf("envelopes/%s/notification", envId)},
		Result: &ret,
	}).Do(ctx, s)
}

// EnvelopeCustomFields returns all custom field info in a Custom Field List
// RestApi documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get%20Envelope%20Custom%20Field%20Information.htm
func (s *Service) EnvelopeCustomFields(ctx context.Context, envId string) (*CustomFieldList, error) {
	var ret *CustomFieldList
	return ret, (&Call{
		Method: "GET",
		URL:    &url.URL{Path: fmt.Sprintf("envelopes/%s/custom_fields", envId)},
		Result: &ret,
	}).Do(ctx, s)

}

// EnvelopeAddCustomFields adds custom fields to an existing envelope.  Duplicates will return error in the
// ErrorDetails struct of the CustomField or ListCustomField item.
// RestApi documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Add%20Envelope%20Custom%20Fields%20to%20an%20Envelope.htm
func (s *Service) EnvelopeAddCustomFields(ctx context.Context, envId string, l *CustomFieldList) (*CustomFieldList, error) {
	var ret *CustomFieldList
	return ret, (&Call{
		Method:  "POST",
		URL:     &url.URL{Path: fmt.Sprintf("envelopes/%s/custom_fields", envId)},
		Payload: l,
		Result:  &ret,
	}).Do(ctx, s)

}

// EnvelopeModifyCustomFields modifies custom fields in CustomFieldList structure.  Id's are mandatory and errors
// are found in ErrorDetails struct of CustomField or ListCustomField items.  Nil ErrorDetails means success
// RestApi documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Modify%20Envelope%20Custom%20Fields%20for%20an%20Envelope.htm
func (s *Service) EnvelopeModifyCustomFields(ctx context.Context, envId string, l *CustomFieldList) (*CustomFieldList, error) {
	var ret *CustomFieldList
	return ret, (&Call{
		Method:  "PUT",
		URL:     &url.URL{Path: fmt.Sprintf("envelopes/%s/custom_fields", envId)},
		Payload: l,
		Result:  &ret,
	}).Do(ctx, s)
}

// EnvelopeRemoveCustomFields deletes Custom Fields using the Id field in CustomField items.
// RestApi documentatiion
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Remove%20Envelope%20Custom%20Fields%20from%20an%20Envelope.htm
func (s *Service) EnvelopeRemoveCustomFields(ctx context.Context, envId string, l *CustomFieldList) (*CustomFieldList, error) {
	var ret *CustomFieldList
	return ret, (&Call{
		Method:  "DELETE",
		URL:     &url.URL{Path: fmt.Sprintf("envelopes/%s/custom_fields", envId)},
		Payload: l,
		Result:  &ret,
	}).Do(ctx, s)

}

// EnvelopeDocuments returns a Document Asset List for an envelope
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get%20List%20of%20Envelope%20Documents.htm
func (s *Service) EnvelopeDocuments(ctx context.Context, envId string) (*DocumentAssetList, error) {
	var ret *DocumentAssetList
	return ret, (&Call{
		Method: "GET",
		URL:    &url.URL{Path: fmt.Sprintf("envelopes/%s/documents", envId)},
		Result: &ret,
	}).Do(ctx, s)
}

// EnvelopeDocument returns the pdf of a specific document from an envelope.  The returned
// response will either have a status code of 200.  Any other status code for a response
// will result in a nil response with a ResponseError detailing the http status code and
// docusign's error message.  Developer is expected to close the http.Response when finished
// processing.
// Only possible arg is show_changes={true/false}
//
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get%20Document%20from%20Envelope.htm
func (s *Service) EnvelopeDocument(ctx context.Context, envId string, docId string, args ...EnvelopeDocumentParam) (*http.Response, error) {
	q := make(url.Values)
	for _, nv := range args {
		q.Add(nv.Name, nv.Value)
	}
	var ret *http.Response
	return ret, (&Call{
		Method: "GET",
		URL:    &url.URL{Path: fmt.Sprintf("envelopes/%s/documents", envId), RawQuery: q.Encode()},
		Result: &ret,
	}).Do(ctx, s)
}

type EnvelopeDocumentParam NmVal

var EnvelopeDocumentShowChanges = EnvelopeDocumentParam{
	Name:  "show_changes",
	Value: "true",
}

// EnvelopeDocumentsCombined retrieves a PDF containing the combined content of all documents
// and the certificate via an http.Response. If the account has the Highlight Data Changes
// feature enabled,there is an option to request that any changes in the envelope be highlighted.
// The returned response will either have a status code of 200.  Any other status code
// for a response will result in a nil response with a ResponseError detailing the http
// status code and docusign's error message.  Developer is expected to close the http.Response
// when finished processing.
// Optional additions: certificate={true or false}, show_changes={true}, watermark={true or false}
//
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get%20Envelope%20Documents%20and%20Certificate.htm
func (s *Service) EnvelopeDocumentsCombined(ctx context.Context, envId string, args ...EnvelopeDocumentsCombinedParam) (*http.Response, error) {
	q := make(url.Values)
	for _, nv := range args {
		q.Add(nv.Name, nv.Value)
	}
	var ret *http.Response
	return ret, (&Call{
		Method: "GET",
		URL:    &url.URL{Path: fmt.Sprintf("envelopes/%s/documents/combined", envId), RawQuery: q.Encode()},
		Result: &ret,
	}).Do(ctx, s)
}

type EnvelopeDocumentsCombinedParam NmVal

var EnvelopeDocumentsCombinedCert = EnvelopeDocumentsCombinedParam{
	Name:  "certificate",
	Value: "true",
}
var EnvelopeDocumentCombinedShowChanges = EnvelopeDocumentsCombinedParam{
	Name:  "show_changes",
	Value: "true",
}
var EnvelopeDocumentsCombinedWatermark = EnvelopeDocumentsCombinedParam{
	Name:  "watermark",
	Value: "true",
}

// LoginInformation determine if a user is authenticated and to choose the account to be used
// for other operations. Each account associated with the login credentials is listed.
// optional paramenters:
// api_password	Boolean
// include_account_id_guid Boolean
// login_settings string [all,none]
func (s *Service) LoginInformation(ctx context.Context, args ...LoginInfoParam) (*LoginInfo, error) {
	q := make(url.Values)
	for _, nv := range args {
		q.Add(nv.Name, nv.Value)
	}
	var ret *LoginInfo
	return ret, (&Call{
		Method: "GET",
		URL:    &url.URL{Path: "/login_information", RawQuery: q.Encode()},
		Result: &ret,
	}).Do(ctx, s)
}

type LoginInfoParam NmVal

var LoginInformationIncludeApiPassword = LoginInfoParam{
	Name:  "api_password",
	Value: "true",
}
var LoginInformationAcctGUID = LoginInfoParam{
	Name:  "include_account_id_guid",
	Value: "true",
}
var LoginInformationSettingsAll = LoginInfoParam{
	Name:  "login_settings",
	Value: "all",
}
var LoginInformationSettingsNone = LoginInfoParam{
	Name:  "login_settings",
	Value: "none",
}

// EnvelopeCreate adds an envelope.  The Status field determines whether the envelope is saved as a Draft
// or sent.
// RestApi Documentation
//
func (s *Service) EnvelopeCreate(ctx context.Context, env *Envelope, files ...*UploadFile) (*EnvelopeResponse, error) {
	var ret *EnvelopeResponse
	return ret, (&Call{
		Method:  "POST",
		URL:     &url.URL{Path: "envelopes"},
		Payload: env,
		Result:  &ret,
		Files:   files,
	}).Do(ctx, s)

}

// EnvelopeStatusChanges returns envelope status changes for all envelopes. The information returned can be
// modified by adding query strings to limit the request to check between certain dates and times, or for certain envelopes,
// or for certain status codes. It is recommended that you use one or more of the query strings in order to limit the size of the response.
//
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST API References/Get Envelope Status Changes.htm
//
// Optional query strings: from_date={dateTime}, to_date={dateTime}, status={status code}, from_to_status={changed or any or list of statuses},
// envelopeId={envelopeId}, custom_field={envelope custom field name}={envelope custom field value}, transaction_ids={transactionIds (comma separated)}
//
// Use DsQueryTimeFormat to formate dateTime query arguments
func (s *Service) EnvelopeStatusChanges(ctx context.Context, args ...EnvelopeStatusChangesParam) (*EnvelopeList, error) {
	q := make(url.Values)
	for _, nv := range args {
		q.Add(nv.Name, nv.Value)
	}
	var ret *EnvelopeList
	return ret, (&Call{
		Method: "GET",
		URL:    &url.URL{Path: "envelopes", RawQuery: q.Encode()},
		Result: &ret,
	}).Do(ctx, s)
}

type EnvelopeStatusChangesParam NmVal

func StatusChangeFromDate(t time.Time) EnvelopeStatusChangesParam {
	return EnvelopeStatusChangesParam{Name: "from_date", Value: DsQueryTimeFormat(t)}
}

func StatusChangeToDate(t time.Time) EnvelopeStatusChangesParam {
	return EnvelopeStatusChangesParam{Name: "to_date", Value: DsQueryTimeFormat(t)}
}

func StatusChangeStatusCode(status string) EnvelopeStatusChangesParam {
	return EnvelopeStatusChangesParam{Name: "status", Value: status}
}

func StatusChangeFromToStatus(status string) EnvelopeStatusChangesParam {
	return EnvelopeStatusChangesParam{Name: "from_to_status", Value: status}
}

func StatusChangeEnvelope(envId string) EnvelopeStatusChangesParam {
	return EnvelopeStatusChangesParam{Name: "envelopeId", Value: envId}
}

func StatusChangeCustomField(fldName, fldValue string) EnvelopeStatusChangesParam {
	return EnvelopeStatusChangesParam{Name: "custom_field", Value: fldName + "=" + fldValue}
}

func StatusChangeTransactionId(transactionIds ...string) EnvelopeStatusChangesParam {
	return EnvelopeStatusChangesParam{Name: "transaction_ids", Value: strings.Join(transactionIds, ",")}
}

// EnvelopeStatus returns returns the overall status for a single envelope.
//
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST API References/Get Envelope Status for One Envelope.htm
func (s *Service) EnvelopeStatus(ctx context.Context, envId string) (res *EnvelopeUris, err error) {
	var ret *EnvelopeUris
	c := &Call{
		Method: "GET",
		URL:    &url.URL{Path: "envelopes/" + envId},
		Result: &ret,
	}
	return ret, c.Do(ctx, s)
}

// EnvelopeStatusMulti returns the status for the requested envelopes.
//
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get%20Envelope%20Status%20for%20more%20than%20one%20envelope.htm
func (s *Service) EnvelopeStatusMulti(ctx context.Context, envIds ...string) ([]EnvelopeUris, error) {
	var retVal struct {
		Envelopes     []EnvelopeUris `json:"envelopes"`
		ResultSetSize string         `json:"resultSetSize"`
	}
	envList := map[string][]string{"envelopeIds": envIds}
	return retVal.Envelopes, (&Call{
		Method:  "PUT",
		URL:     &url.URL{Path: "envelopes/status", RawQuery: "envelope_ids=request_body"},
		Payload: envList,
		Result:  &retVal,
	}).Do(ctx, s)
}

func (s *Service) EnvelopeSetDocuments(ctx context.Context, envId string, dl *DocumentList, files ...*UploadFile) (*DocumentAssetList, error) {
	var ret *DocumentAssetList
	return ret, (&Call{
		Method:  "PUT",
		URL:     &url.URL{Path: fmt.Sprintf("envelopes/%s/documents", envId)},
		Files:   files,
		Payload: dl,
		Result:  &ret,
	}).Do(ctx, s)

}

func (s *Service) EnvelopeRemoveDocuments(ctx context.Context, envId string, dl *DocumentList) (*DocumentAssetList, error) {
	var ret *DocumentAssetList
	return ret, (&Call{
		Method:  "DELETE",
		URL:     &url.URL{Path: fmt.Sprintf("envelopes/%s/documents", envId)},
		Payload: dl,
		Result:  &ret,
	}).Do(ctx, s)

}

// DocumentAddCustomFields creates new custom fields on a specific document.  Errors are returned in the
// DocumentFieldList ErrorDetails field.
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Add%20Custom%20Document%20Fields%20to%20an%20Envelope%20Document.htm
func (s *Service) DocumentAddCustomFields(ctx context.Context, envId, docId string, dfl *DocumentFieldList) (*DocumentFieldList, error) {
	var ret *DocumentFieldList
	return ret, (&Call{
		Method:  "POST",
		URL:     &url.URL{Path: fmt.Sprintf("envelopes/%s/documents/%s/fields", envId, docId)},
		Payload: dfl,
		Result:  &ret,
	}).Do(ctx, s)
}

// DocumentModifyCustomFields modifies existing custom document fields for an existing envelope. Errors are returned in the
// DocumentFieldList ErrorDetails field.
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Modify%20Custom%20Document%20Fields%20for%20an%20Envelope%20Document.htm
func (s *Service) DocumentModifyCustomFields(ctx context.Context, envId, docId string, dfl *DocumentFieldList) (*DocumentFieldList, error) {
	var ret *DocumentFieldList
	return ret, (&Call{
		Method:  "PUT",
		URL:     &url.URL{Path: fmt.Sprintf("envelopes/%s/documents/%s/fields", envId, docId)},
		Payload: dfl,
		Result:  &ret,
	}).Do(ctx, s)
}

// DocumentRemoveCustomFields delete existing custom document fields for an existing envelope. Errors are returned in the
// DocumentFieldList ErrorDetails field.
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Delete%20Custom%20Document%20Fields%20from%20an%20Envelope%20Document.htm
func (s *Service) DocumentRemoveCustomFields(ctx context.Context, envId, docId string, dfl *DocumentFieldList) (*DocumentFieldList, error) {
	var ret *DocumentFieldList
	return ret, (&Call{
		Method:  "DELETE",
		URL:     &url.URL{Path: fmt.Sprintf("envelopes/%s/documents/%s/fields", envId, docId)},
		Payload: dfl,
		Result:  &ret,
	}).Do(ctx, s)
}

// Recipients
// Optional query strings: include_tabs={true or false}, include_extended={true or false}
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get%20Envelope%20Recipient%20Status.htm
func (s *Service) Recipients(ctx context.Context, envId string, args ...RecipientsParam) (*RecipientList, error) {
	q := make(url.Values)
	for _, nv := range args {
		q.Add(nv.Name, nv.Value)
	}
	var ret *RecipientList
	return ret, (&Call{
		Method: "GET",
		URL:    &url.URL{Path: fmt.Sprintf("envelopes/%s/recipients", envId), RawQuery: q.Encode()},
		Result: &ret,
	}).Do(ctx, s)
}

type RecipientsParam NmVal

var RecipientsIncludeTabs = RecipientsParam{
	Name:  "include_tabs",
	Value: "true",
}
var RecipientsIncludeExtended = RecipientsParam{
	Name:  "include_extended",
	Value: "true",
}
var RecipientsIncludeAnchorTabs = RecipientsParam{
	Name:  "include_anchor_tab_locations",
	Value: "true",
}
var RecipientsResend = RecipientsParam{
	Name:  "resend_envelope",
	Value: "true",
}

// RecipientsAdd
// If an error occurred during the operation, recipient struct will contain an ErrorDetail
// Optional addition: resend_envelope {true or false}
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Modify%20or%20Correct%20and%20Resend%20Recipient%20Information.htm
func (s *Service) RecipientsAdd(ctx context.Context, envId string, rl *RecipientList, args ...RecipientsParam) (*RecipientList, error) {
	q := make(url.Values)
	for _, nv := range args {
		q.Add(nv.Name, nv.Value)
	}
	var ret *RecipientList
	return ret, (&Call{
		Method:  "POST",
		URL:     &url.URL{Path: fmt.Sprintf("envelopes/%s/recipients", envId), RawQuery: q.Encode()},
		Payload: rl,
		Result:  &ret,
	}).Do(ctx, s)

}

// ModifyRecipients
// Optional addition: resend_envelope {true or false}
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get%20Envelope%20Recipient%20Status.htm
func (s *Service) RecipientsModify(ctx context.Context, envId string, rl *RecipientList, args ...RecipientsParam) (*RecipientUpdateResult, error) {
	q := make(url.Values)
	for _, nv := range args {
		q.Add(nv.Name, nv.Value)
	}
	var ret *RecipientUpdateResult
	return ret, (&Call{
		Method:  "PUT",
		URL:     &url.URL{Path: fmt.Sprintf("envelopes/%s/recipients", envId), RawQuery: q.Encode()},
		Payload: rl,
		Result:  &ret,
	}).Do(ctx, s)

}

// RecipientsRemove
// Optional addition: resend_envelope {true or false}
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Delete%20Recipients%20from%20an%20Envelope.htm
func (s *Service) RecipientsRemove(ctx context.Context, envId string, rl *RecipientList) (res *RecipientList, err error) {
	var ret *RecipientList
	return ret, (&Call{
		Method:  "DELETE",
		URL:     &url.URL{Path: fmt.Sprintf("envelopes/%s/recipients", envId)},
		Payload: rl,
		Result:  &ret,
	}).Do(ctx, s)
}

// RecipientTabs
// RestApiDocumentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get%20Tab%20Information%20for%20a%20Recipient.htm
func (s *Service) RecipientTabs(ctx context.Context, envId string, recipId string) (*Tabs, error) {
	var ret *Tabs
	return ret, (&Call{
		Method: "GET",
		URL:    &url.URL{Path: fmt.Sprintf("envelopes/%s/recipients/%s/tabs", envId, recipId)},
		Result: &ret,
	}).Do(ctx, s)
}

// RecipientTabsAdd adds tabs to a recipient The response returns the success or failure of each document being added
// to the envelope and the envelope ID. Failed operations will add the ErrorDetails structure containing
// an error code and message. If ErrorDetails is nil, then the operation was successful for that item.
// RestApiDocumentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Add%20Tabs%20for%20a%20Recipient.htm
func (s *Service) RecipientTabsAdd(ctx context.Context, envId string, recipId string, tb *Tabs) (*Tabs, error) {
	var ret *Tabs
	return ret, (&Call{
		Method:  "POST",
		URL:     &url.URL{Path: fmt.Sprintf("envelopes/%s/recipients/%s/tabs", envId, recipId)},
		Payload: tb,
		Result:  &ret,
	}).Do(ctx, s)
}

// RecipientTabsModify
// The parameters used to modify tabs are the same as those used in an envelope, but you can only modify existing tabs
// and the tabId must be included.
// RestApiDocumentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Modify%20Tabs%20for%20a%20Recipient.htm
func (s *Service) RecipientTabsModify(ctx context.Context, envId string, recipId string, tb *Tabs) (*Tabs, error) {
	var ret *Tabs
	return ret, (&Call{
		Method:  "PUT",
		URL:     &url.URL{Path: fmt.Sprintf("envelopes/%s/recipients/%s/tabs", envId, recipId)},
		Payload: tb,
		Result:  &ret,
	}).Do(ctx, s)
}

// RecipientTabsRemove
// If an error occurred during the DELETE operation for any of the recipients, that recipient will contain
// an error node with an errorCode and message.
// RestApiDocumentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Delete%20Recipients%20from%20an%20Envelope.htm
func (s *Service) RecipientTabsRemove(ctx context.Context, envId string, recipId string, tb *Tabs) (*Tabs, error) {
	var ret *Tabs
	return ret, (&Call{
		Method:  "DELETE",
		URL:     &url.URL{Path: fmt.Sprintf("envelopes/%s/recipients/%s/tabs", envId, recipId)},
		Payload: tb,
		Result:  &ret,
	}).Do(ctx, s) //urlStr := fmt.Sprintf("envelopes/%s/recipients/%s/tabs", envId, recipId)
}

// TemplateSearch retrieves the list of templates for the specified account
// Optional query strings: folder={string}, folder_ids={GUID, GUID}, include={string}, count={integer},
// start_position={integer}, from_date={date/time}, to_date={date/time}, used_from_date={date/time},
// used_to_date={date/time}, search_text={string}, order={string}, order_by={string}, user_filter={string},
// shared_by_me={true/false}
// RestApiDocumentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get%20List%20of%20Templates.htm
func (s *Service) TemplateSearch(ctx context.Context, args ...TemplateSearchParam) (*FolderTemplateList, error) {
	q := make(url.Values)
	for _, nv := range args {
		q.Add(nv.Name, nv.Value)
	}
	var ret *FolderTemplateList
	return ret, (&Call{
		Method: "GET",
		URL:    &url.URL{Path: "templates", RawQuery: q.Encode()},
		Result: &ret,
	}).Do(ctx, s)
}

type TemplateSearchParam NmVal

func TemplateSearchFolder(folder string) TemplateSearchParam {
	return TemplateSearchParam{Name: "folder", Value: folder}
}

func TemplateSearchFolderIds(folder ...string) TemplateSearchParam {
	return TemplateSearchParam{Name: "folder", Value: strings.Join(folder, ",")}
}

func TemplateSearchInclude(recipients, folders, documents, customFields, notifications bool) TemplateSearchParam {
	vals := make([]string, 0, 5)
	if recipients {
		vals = append(vals, "recipients")
	}
	if folders {
		vals = append(vals, "folders")
	}
	if documents {
		vals = append(vals, "documents")
	}
	if customFields {
		vals = append(vals, "custom_fields")
	}
	if notifications {
		vals = append(vals, "notifications")
	}
	return TemplateSearchParam{Name: "include", Value: strings.Join(vals, ",")}
}

func TemplateSearchCount(count int) TemplateSearchParam {
	return TemplateSearchParam{Name: "count", Value: strconv.Itoa(count)}
}

func TemplateSearchStartPosition(pos int) TemplateSearchParam {
	return TemplateSearchParam{Name: "start_position", Value: strconv.Itoa(pos)}
}

func TemplateSearchFromDate(t time.Time) TemplateSearchParam {
	return TemplateSearchParam{Name: "from_date", Value: DsQueryTimeFormat(t)}
}

func TemplateSearchToDate(t time.Time) TemplateSearchParam {
	return TemplateSearchParam{Name: "to_date", Value: DsQueryTimeFormat(t)}
}

func TemplateSearchUsedFromDate(t time.Time) TemplateSearchParam {
	return TemplateSearchParam{Name: "used_from_date", Value: DsQueryTimeFormat(t)}
}

func TemplateSearchUsedToDate(t time.Time) TemplateSearchParam {
	return TemplateSearchParam{Name: "used_to_date", Value: DsQueryTimeFormat(t)}
}

func TemplateSearchSearch(searchText string) TemplateSearchParam {
	return TemplateSearchParam{Name: "search_text", Value: searchText}
}

var TemplateSearchOrderAsc = TemplateSearchParam{
	Name:  "order",
	Value: "asc",
}
var TemplateSearchOrderDesc = TemplateSearchParam{
	Name:  "order",
	Value: "desc",
}

var TemplateSearchOrderByName = TemplateSearchParam{
	Name:  "orderby",
	Value: "name",
}
var TemplateSearchOrderByModified = TemplateSearchParam{
	Name:  "orderby",
	Value: "modified",
}
var TemplateSearchOrderByUsed = TemplateSearchParam{
	Name:  "orderby",
	Value: "used",
}

var TemplateSearchFilterOwned = TemplateSearchParam{
	Name:  "user_filter",
	Value: "owned_by_me",
}
var TemplateSearchFilterSharedWithMe = TemplateSearchParam{
	Name:  "user_filter",
	Value: "shared_with_me",
}
var TemplateSearchFilterAll = TemplateSearchParam{
	Name:  "user_filter",
	Value: "all",
}
var TemplateSearchSharedByMe = TemplateSearchParam{
	Name:  "shared_by_me",
	Value: "true",
}
var TemplateSearchNotSharedByMe = TemplateSearchParam{
	Name:  "shared_by_me",
	Value: "false",
}

// EnvelopeCorrection returns a URL to start the correction view of the DocuSign UI.
//
// RestApiDocumentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Post%20Envelope%20Correction.htm
func (s *Service) EnvelopeCorrection(ctx context.Context, envId string, retUrlType ReturnUrlType, suppressNavigation bool) (*EnvUrl, error) {
	var correctType struct {
		Type ReturnUrlType `json:"returnUrl"`
		Nav  string        `json:"suppressNavigation,omitempty"`
	}
	correctType.Type = retUrlType
	if suppressNavigation {
		correctType.Nav = "true"
	}
	var ret *EnvUrl
	return ret, (&Call{
		Method:  "POST",
		URL:     &url.URL{Path: fmt.Sprintf("envelopes/%s/views/correct", envId)},
		Payload: &correctType,
		Result:  &ret,
	}).Do(ctx, s)

}

// RecipientView returns a URL to start a Recipient view of the DocuSign UI.
//
// RestApiDocumentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Post%20Recipient%20View.htm
func (s *Service) RecipientView(ctx context.Context, envId string, er *EnvRecipientView) (*EnvUrl, error) {
	var ret *EnvUrl
	return ret, (&Call{
		Method:  "POST",
		URL:     &url.URL{Path: fmt.Sprintf("envelopes/%s/views/recipient", envId)},
		Payload: er,
		Result:  &ret,
	}).Do(ctx, s)
}

// SenderView returns a URL to start the sender view of the DocuSign UI.
//
// RestApiDocumentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Post%20Sender%20View.htm
func (s *Service) SenderView(ctx context.Context, envId string, retType ReturnUrlType) (*EnvUrl, error) {
	var ret *EnvUrl
	return ret, (&Call{
		Method:  "POST",
		URL:     &url.URL{Path: fmt.Sprintf("envelopes/%s/views/sender", envId)},
		Payload: &returnUrlStruct{ReturnUrl: retType},
		Result:  &ret,
	}).Do(ctx, s)

}

// EditView returns a URL to start the edit view of the DocuSign UI.
//
// RestApiDocumentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Post%20Edit%20View.htm
func (s *Service) EditView(ctx context.Context, envId string, retType ReturnUrlType) (*EnvUrl, error) {
	var ret *EnvUrl
	return ret, (&Call{
		Method:  "POST",
		URL:     &url.URL{Path: fmt.Sprintf("envelopes/%s/views/edit", envId)},
		Payload: &returnUrlStruct{ReturnUrl: retType},
		Result:  &ret,
	}).Do(ctx, s)

}

// EnvelopeTemplates returns a list of templates used by an envelope
//
// RestApiDocumentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get%20List%20of%20Templates%20used%20in%20an%20Envelope.htm
func (s *Service) EnvelopeTemplates(ctx context.Context, envId string) (*TemplateList, error) {
	var ret *TemplateList
	return ret, (&Call{
		Method: "GET",
		URL:    &url.URL{Path: fmt.Sprintf("envelopes/%s/templates", envId)},
		Result: &ret,
	}).Do(ctx, s)

}

// EnvelopeMove move the specified envelope to the folder specified int toFolderId
//
// RestApiDocumentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Move%20Envelope.htm
func (s *Service) EnvelopeMove(ctx context.Context, toFolderId string, envIds ...string) error {
	data := map[string]interface{}{"envelopeIds": envIds}
	return (&Call{
		Method:  "PUT",
		URL:     &url.URL{Path: fmt.Sprintf("folders/%s", toFolderId)},
		Payload: data,
	}).Do(ctx, s)
}

// AccountCustomFields retrieves a list of envelope custom fields associated with the account.
//
// RestApiDocumentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get List of Account Custom Fields.htm
func (s *Service) AccountCustomFields(ctx context.Context) (*CustomFieldList, error) {
	var ret *CustomFieldList
	return ret, (&Call{
		Method: "GET",
		URL:    &url.URL{Path: "custom_fields"},
		Result: &ret,
	}).Do(ctx, s)

}

// GetTemplate returns field data for the specified template
//
// RestApiDocumentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get Template.htm
func (s *Service) GetTemplate(ctx context.Context, id string) (*Template, error) {
	var ret *Template
	return ret, (&Call{
		Method: "GET",
		URL:    &url.URL{Path: fmt.Sprintf("templates/%s", id)},
		Result: &ret,
	}).Do(ctx, s)
}

// VoidEnvelope voids and existing envelope.
//
// RestApiDocumentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Void Envelope.htm
func (s *Service) Void(ctx context.Context, envId string, reason string) error {
	c := &Call{
		Method:  "PUT",
		URL:     &url.URL{Path: fmt.Sprintf("envelopes/%s", envId)},
		Payload: map[string]string{"status": "voided", "voidedReason": reason},
	}
	return c.Do(ctx, s)
}

// Remind sends a reminder to an envelope recipient.
//
// RestApiDocumentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References//Void Envelope.htm
func (s *Service) Remind(ctx context.Context, envId string, rl *RecipientList) error {
	ret := make(map[string]interface{})
	c := &Call{
		Method:  "PUT",
		URL:     &url.URL{Path: fmt.Sprintf("envelopes/%s/recipients", envId), RawQuery: "resend_envelope=true"},
		Payload: rl,
		Result:  &ret,
	}
	return c.Do(ctx, s)

}
