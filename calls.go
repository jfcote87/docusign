// Copyright 2015 James Cote and Liberty Fund, Inc.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package docusign

import (
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// FolderList returns a list of the folders for the account, including the
// folder hierarchy.
// Optional query string: template={string} include or only
//
// RestApi documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get%20Folder%20List.htm
func (s *Service) FolderList(args ...FolderTemplateParam) (fl *FolderList, err error) {
	q := make(url.Values)
	for _, nv := range args {
		q.Add(nv.Name, nv.Value)
	}
	var queryStr string
	if len(q) > 0 {
		queryStr = "?" + q.Encode()
	}
	urlStr := "folders" + queryStr
	err = s.do("GET", urlStr, nil, &fl)
	return
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
func (s *Service) FolderEnvSearch(folderId string, args ...FolderEnvSearchParam) (f *FolderEnvList, err error) {
	q := make(url.Values)
	for _, nv := range args {
		q.Add(nv.Name, nv.Value)
	}
	var queryStr string
	if len(q) > 0 {
		queryStr = "?" + q.Encode()
	}
	urlStr := fmt.Sprintf("folders/%s", folderId) + queryStr
	err = s.do("GET", urlStr, nil, &f)
	return
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
func (s *Service) EnvelopeSearch(searchFld SearchFolder, args ...SearchFolderParam) (f *FolderEnvList, err error) {
	q := make(url.Values)
	for _, nv := range args {
		q.Add(nv.Name, nv.Value)
	}
	var queryStr string
	if len(q) > 0 {
		queryStr = "?" + q.Encode()
	}
	urlStr := fmt.Sprintf("search_folders/%s", searchFld) + queryStr
	err = s.do("GET", urlStr, nil, &f)
	return
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
func (s *Service) EnvelopeAuditEvents(envId string) (a *AuditEventList, err error) {
	urlStr := fmt.Sprintf("envelopes/%s/audit_events", envId)
	err = s.do("GET", urlStr, nil, &a)
	return
}

// EnvelopeNotification returns the reminder and expiration information for the envelope.
//
// RestApi documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get%20Envelope%20Notification%20Information.htm
func (s *Service) EnvelopeNotification(envId string) (n *Notification, err error) {
	urlStr := fmt.Sprintf("envelopes/%s/notification", envId)
	err = s.do("GET", urlStr, nil, &n)
	return
}

// EnvelopeCustomFields returns all custom field info in a Custom Field List
// RestApi documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get%20Envelope%20Custom%20Field%20Information.htm
func (s *Service) EnvelopeCustomFields(envId string) (cl *CustomFieldList, err error) {
	urlStr := fmt.Sprintf("envelopes/%s/custom_fields", envId)
	err = s.do("GET", urlStr, nil, &cl)
	return
}

// EnvelopeAddCustomFields adds custom fields to an existing envelope.  Duplicates will return error in the
// ErrorDetails struct of the CustomField or ListCustomField item.
// RestApi documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Add%20Envelope%20Custom%20Fields%20to%20an%20Envelope.htm
func (s *Service) EnvelopeAddCustomFields(envId string, l *CustomFieldList) (cl *CustomFieldList, err error) {
	urlStr := fmt.Sprintf("envelopes/%s/custom_fields", envId)
	err = s.do("POST", urlStr, l, &cl)
	return
}

// EnvelopeModifyCustomFields modifies custom fields in CustomFieldList structure.  Id's are mandatory and errors
// are found in ErrorDetails struct of CustomField or ListCustomField items.  Nil ErrorDetails means success
// RestApi documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Modify%20Envelope%20Custom%20Fields%20for%20an%20Envelope.htm
func (s *Service) EnvelopeModifyCustomFields(envId string, l *CustomFieldList) (cl *CustomFieldList, err error) {
	urlStr := fmt.Sprintf("envelopes/%s/custom_fields", envId)
	err = s.do("PUT", urlStr, l, &cl)
	return
}

// EnvelopeRemoveCustomFields deletes Custom Fields using the Id field in CustomField items.
// RestApi documentatiion
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Remove%20Envelope%20Custom%20Fields%20from%20an%20Envelope.htm
func (s *Service) EnvelopeRemoveCustomFields(envId string, l *CustomFieldList) (cl *CustomFieldList, err error) {
	urlStr := fmt.Sprintf("envelopes/%s/custom_fields", envId)
	err = s.do("DELETE", urlStr, l, &cl)
	return
}

// EnvelopeDocuments returns a Document Asset List for an envelope
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get%20List%20of%20Envelope%20Documents.htm
func (s *Service) EnvelopeDocuments(envId string) (docList *DocumentAssetList, err error) {
	urlStr := fmt.Sprintf("envelopes/%s/documents", envId)
	err = s.do("GET", urlStr, nil, &docList)
	return
}

// EnvelopeDocument returns a specific document from an envelope.  The outputBuffer is an io.Writer used for
// saving the pdf file.  The pdf will be written to the outputBuffer io.Writer.
// Only possible arg is show_changes={true/false}
//
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get%20Document%20from%20Envelope.htm
func (s *Service) EnvelopeDocument(outputBuffer io.Writer, envId string, docId string, args ...EnvelopeDocumentParam) (err error) {
	q := make(url.Values)
	for _, nv := range args {
		q.Add(nv.Name, nv.Value)
	}
	var queryStr string
	if len(q) > 0 {
		queryStr = "?" + q.Encode()
	}
	urlStr := fmt.Sprintf("envelopes/%s/documents/%s", envId, docId) + queryStr
	err = s.doPdf(outputBuffer, "GET", urlStr, nil)
	return
}

type EnvelopeDocumentParam NmVal

var EnvelopeDocumentShowChanges = EnvelopeDocumentParam{
	Name:  "show_changes",
	Value: "true",
}

// EnvelopeDocumentsCombined retrieves a PDF containing the combined content of all documents
// and the certificate. If the account has the Highlight Data Changes feature enabled,
// there is an option to request that any changes in the envelope be highlighted. The pdf will
// be written to the outputBuffer io.Writer.
// Optional additions: certificate={true or false}, show_changes={true}, watermark={true or false}
//
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get%20Envelope%20Documents%20and%20Certificate.htm
func (s *Service) EnvelopeDocumentsCombined(outputBuffer io.Writer, envId string, args ...EnvelopeDocumentsCombinedParam) (err error) {
	q := make(url.Values)
	for _, nv := range args {
		q.Add(nv.Name, nv.Value)
	}
	var queryStr string
	if len(q) > 0 {
		queryStr = "?" + q.Encode()
	}
	urlStr := fmt.Sprintf("envelopes/%s/documents/combined", envId) + queryStr
	err = s.doPdf(outputBuffer, "GET", urlStr, nil)
	return
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
func (s *Service) LoginInformation(args ...LoginInfoParam) (li *LoginInfo, err error) {
	q := make(url.Values)
	for _, nv := range args {
		q.Add(nv.Name, nv.Value)
	}
	var queryStr string
	if len(q) > 0 {
		queryStr = "?" + q.Encode()
	}
	urlStr := "/login_information" + queryStr
	fmt.Printf("Q: %s\n", queryStr)
	err = s.do("GET", urlStr, nil, &li)
	return
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
func (s *Service) EnvelopeCreate(env *Envelope, files ...*UploadFile) (envResp *EnvelopeResponse, err error) {
	urlStr := ("envelopes")
	err = s.do("POST", urlStr, env, &envResp, files...)
	return
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
func (s *Service) EnvelopeStatusChanges(args ...EnvelopeStatusChangesParam) (el *EnvelopeList, err error) {
	q := make(url.Values)
	for _, nv := range args {
		q.Add(nv.Name, nv.Value)
	}
	var queryStr string
	if len(q) > 0 {
		queryStr = "?" + q.Encode()
	}
	urlStr := "envelopes" + queryStr
	err = s.do("GET", urlStr, nil, &el)
	return
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
func (s *Service) EnvelopeStatus(envId string) (res *EnvelopeUris, err error) {
	urlStr := "envelopes/" + envId
	err = s.do("GET", urlStr, nil, &res)
	return
}

// EnvelopeStatusMulti returns the status for the requested envelopes.
//
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get%20Envelope%20Status%20for%20more%20than%20one%20envelope.htm
func (s *Service) EnvelopeStatusMulti(envIds ...string) (res []EnvelopeUris, err error) {
	var retVal struct {
		Envelopes     []EnvelopeUris `json:"envelopes"`
		ResultSetSize string         `json:"resultSetSize"`
	}
	envList := map[string][]string{"envelopeIds": envIds}
	urlStr := "envelopes/status?envelope_ids=request_body"
	err = s.do("PUT", urlStr, envList, &retVal)
	if err != nil {
		return nil, err
	}
	return retVal.Envelopes, nil
}

func (s *Service) EnvelopeSetDocuments(envId string, dl *DocumentList, files ...*UploadFile) (res *DocumentAssetList, err error) {
	urlStr := fmt.Sprintf("envelopes/%s/documents", envId)
	err = s.do("PUT", urlStr, dl, &res, files...)
	return
}

func (s *Service) EnvelopeRemoveDocuments(envId string, dl *DocumentList) (res *DocumentAssetList, err error) {
	urlStr := fmt.Sprintf("envelopes/%s/documents", envId)
	err = s.do("DELETE", urlStr, dl, &res)
	return
}

// DocumentAddCustomFields creates new custom fields on a specific document.  Errors are returned in the
// DocumentFieldList ErrorDetails field.
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Add%20Custom%20Document%20Fields%20to%20an%20Envelope%20Document.htm
func (s *Service) DocumentAddCustomFields(envId, docId string, dfl *DocumentFieldList) (res *DocumentFieldList, err error) {
	urlStr := fmt.Sprintf("envelopes/%s/documents/%s/fields", envId, docId)
	err = s.do("POST", urlStr, dfl, &res)
	return
}

// DocumentModifyCustomFields modifies existing custom document fields for an existing envelope. Errors are returned in the
// DocumentFieldList ErrorDetails field.
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Modify%20Custom%20Document%20Fields%20for%20an%20Envelope%20Document.htm
func (s *Service) DocumentModifyCustomFields(envId, docId string, dfl *DocumentFieldList) (res *DocumentFieldList, err error) {
	urlStr := fmt.Sprintf("envelopes/%s/documents/%s/fields", envId, docId)
	err = s.do("PUT", urlStr, dfl, &res)
	return
}

// DocumentRemoveCustomFields delete existing custom document fields for an existing envelope. Errors are returned in the
// DocumentFieldList ErrorDetails field.
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Delete%20Custom%20Document%20Fields%20from%20an%20Envelope%20Document.htm
func (s *Service) DocumentRemoveCustomFields(envId, docId string, dfl *DocumentFieldList) (res *DocumentFieldList, err error) {
	urlStr := fmt.Sprintf("envelopes/%s/documents/%s/fields", envId, docId)
	err = s.do("DELETE", urlStr, dfl, &res)
	return
}

// Recipients
// Optional query strings: include_tabs={true or false}, include_extended={true or false}
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get%20Envelope%20Recipient%20Status.htm
func (s *Service) Recipients(envId string, args ...RecipientsParam) (res *RecipientList, err error) {
	q := make(url.Values)
	for _, nv := range args {
		q.Add(nv.Name, nv.Value)
	}
	var queryStr string
	if len(q) > 0 {
		queryStr = "?" + q.Encode()
	}
	urlStr := fmt.Sprintf("envelopes/%s/recipients", envId) + queryStr
	err = s.do("GET", urlStr, nil, &res)
	return
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
var RecipientsResend = RecipientsParam{
	Name:  "resend_envelope",
	Value: "true",
}

// RecipientsAdd
// If an error occurred during the operation, recipient struct will contain an ErrorDetail
// Optional addition: resend_envelope {true or false}
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Modify%20or%20Correct%20and%20Resend%20Recipient%20Information.htm
func (s *Service) RecipientsAdd(envId string, rl *RecipientList, args ...RecipientsParam) (res *RecipientList, err error) {
	q := make(url.Values)
	for _, nv := range args {
		q.Add(nv.Name, nv.Value)
	}
	var queryStr string
	if len(q) > 0 {
		queryStr = "?" + q.Encode()
	}
	urlStr := fmt.Sprintf("envelopes/%s/recipients", envId) + queryStr
	err = s.do("POST", urlStr, rl, &res)
	return
}

// ModifyRecipients
// Optional addition: resend_envelope {true or false}
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get%20Envelope%20Recipient%20Status.htm
func (s *Service) RecipientsModify(envId string, rl *RecipientList, args ...RecipientsParam) (res *RecipientUpdateResult, err error) {
	q := make(url.Values)
	for _, nv := range args {
		q.Add(nv.Name, nv.Value)
	}
	var queryStr string
	if len(q) > 0 {
		queryStr = "?" + q.Encode()
	}
	urlStr := fmt.Sprintf("envelopes/%s/recipients", envId) + queryStr
	err = s.do("PUT", urlStr, rl, &res)
	return
}

// RecipientsRemove
// Optional addition: resend_envelope {true or false}
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Delete%20Recipients%20from%20an%20Envelope.htm
func (s *Service) RecipientsRemove(envId string, rl *RecipientList) (res *RecipientList, err error) {
	urlStr := fmt.Sprintf("envelopes/%s/recipients", envId)
	err = s.do("DELETE", urlStr, rl, &res)
	return
}

// RecipientTabs
// RestApiDocumentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get%20Tab%20Information%20for%20a%20Recipient.htm
func (s *Service) RecipientTabs(envId string, recipId string) (res *Tabs, err error) {
	urlStr := fmt.Sprintf("envelopes/%s/recipients/%s/tabs", envId, recipId)
	err = s.do("GET", urlStr, nil, &res)
	return
}

// RecipientTabsAdd adds tabs to a recipient The response returns the success or failure of each document being added
// to the envelope and the envelope ID. Failed operations will add the ErrorDetails structure containing
// an error code and message. If ErrorDetails is nil, then the operation was successful for that item.
// RestApiDocumentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Add%20Tabs%20for%20a%20Recipient.htm
func (s *Service) RecipientTabsAdd(envId string, recipId string, tb *Tabs) (res *Tabs, err error) {
	urlStr := fmt.Sprintf("envelopes/%s/recipients/%s/tabs", envId, recipId)
	err = s.do("POST", urlStr, tb, &res)
	return
}

// RecipientTabsModify
// The parameters used to modify tabs are the same as those used in an envelope, but you can only modify existing tabs
// and the tabId must be included.
// RestApiDocumentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Modify%20Tabs%20for%20a%20Recipient.htm
func (s *Service) RecipientTabsModify(envId string, recipId string, tb *Tabs) (res *Tabs, err error) {
	urlStr := fmt.Sprintf("envelopes/%s/recipients/%s/tabs", envId, recipId)
	err = s.do("PUT", urlStr, tb, &res)
	return
}

// RecipientTabsRemove
// If an error occurred during the DELETE operation for any of the recipients, that recipient will contain
// an error node with an errorCode and message.
// RestApiDocumentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Delete%20Recipients%20from%20an%20Envelope.htm
func (s *Service) RecipientTabsRemove(envId string, recipId string, tb *Tabs) (res *Tabs, err error) {
	urlStr := fmt.Sprintf("envelopes/%s/recipients/%s/tabs", envId, recipId)
	err = s.do("DELETE", urlStr, tb, &res)
	return
}

// TemplateSearch retrieves the list of templates for the specified account
// Optional query strings: folder={string}, folder_ids={GUID, GUID}, include={string}, count={integer},
// start_position={integer}, from_date={date/time}, to_date={date/time}, used_from_date={date/time},
// used_to_date={date/time}, search_text={string}, order={string}, order_by={string}, user_filter={string},
// shared_by_me={true/false}
// RestApiDocumentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get%20List%20of%20Templates.htm
func (s *Service) TemplateSearch(args ...TemplateSearchParam) (res *FolderTemplateList, err error) {
	q := make(url.Values)
	for _, nv := range args {
		q.Add(nv.Name, nv.Value)
	}
	var queryStr string
	if len(q) > 0 {
		queryStr = "?" + q.Encode()
	}
	urlStr := "templates" + queryStr
	err = s.do("GET", urlStr, nil, &res)
	return
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
func (s *Service) EnvelopeCorrection(envId string, retUrlType ReturnUrlType, suppressNavigation bool) (res *EnvUrl, err error) {
	var correctType struct {
		Type ReturnUrlType `json:"returnUrl"`
		Nav  string        `json:"suppressNavigation,omitempty"`
	}
	correctType.Type = retUrlType
	if suppressNavigation {
		correctType.Nav = "true"
	}
	urlStr := fmt.Sprintf("envelopes/%s/views/correct", envId)
	err = s.do("POST", urlStr, &correctType, &res)
	return
}

// RecipientView returns a URL to start a Recipient view of the DocuSign UI.
//
// RestApiDocumentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Post%20Recipient%20View.htm
func (s *Service) RecipientView(envId string, er *EnvRecipientView) (res *EnvUrl, err error) {
	urlStr := fmt.Sprintf("envelopes/%s/views/recipient", envId)
	err = s.do("POST", urlStr, er, &res)
	return
}

// SenderView returns a URL to start the sender view of the DocuSign UI.
//
// RestApiDocumentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Post%20Sender%20View.htm
func (s *Service) SenderView(envId string, retType ReturnUrlType) (res *EnvUrl, err error) {
	payload := &returnUrlStruct{ReturnUrl: retType}
	urlStr := fmt.Sprintf("envelopes/%s/views/sender", envId)
	err = s.do("POST", urlStr, payload, &res)
	return
}

// EditView returns a URL to start the edit view of the DocuSign UI.
//
// RestApiDocumentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Post%20Edit%20View.htm
func (s *Service) EditView(envId string, retType ReturnUrlType) (res *EnvUrl, err error) {
	payload := &returnUrlStruct{ReturnUrl: retType}
	urlStr := fmt.Sprintf("envelopes/%s/views/edit", envId)
	err = s.do("POST", urlStr, payload, &res)
	return
}

// EnvelopeTemplates returns a list of templates used by an envelope
//
// RestApiDocumentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get%20List%20of%20Templates%20used%20in%20an%20Envelope.htm
func (s *Service) EnvelopeTemplates(envId string) (res *TemplateList, err error) {
	urlStr := fmt.Sprintf("envelopes/%s/templates", envId)
	err = s.do("GET", urlStr, nil, &res)
	return
}

// EnvelopeMove move the specified envelope to the folder specified int toFolderId
//
// RestApiDocumentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Move%20Envelope.htm
func (s *Service) EnvelopeMove(toFolderId string, envIds ...string) (err error) {
	data := map[string]interface{}{"envelopeIds": envIds}
	var retData interface{}
	urlStr := fmt.Sprintf("folders/%s", toFolderId)
	err = s.do("PUT", urlStr, data, &retData)
	return
}

// AccountCustomFields retrieves a list of envelope custom fields associated with the account.
//
// RestApiDocumentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get List of Account Custom Fields.htm
func (s *Service) AccountCustomFields() (res *CustomFieldList, err error) {
	err = s.do("GET", "custom_fields", nil, &res)
	return
}

// GetTemplate returns field data for the specified template
//
// RestApiDocumentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Get Template.htm
func (s *Service) GetTemplate(id string) (res *Template, err error) {
	urlStr := fmt.Sprintf("templates/%s", id)
	err = s.do("GET", urlStr, nil, &res)
	return
}
