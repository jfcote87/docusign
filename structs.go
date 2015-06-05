// Copyright 2015 James Cote and Liberty Fund, Inc.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package docusign

import (
	"time"
)

// ReturnUrlType is used to generate url for user signing
// see PostSenderView, PostRecipientView, PostEditView, PostCorrection
type ReturnUrlType string

type returnUrlStruct struct {
	ReturnUrl ReturnUrlType `json:"returnUrl"`
}

const ReturnUrlTypeSend = "send"
const ReturnUrlTypeSave = "save"
const ReturnUrlTypeCancel = "cancel"
const ReturnUrlTypeError = "error"
const ReturnUrlTypeSessionEnd = "sessionEnd"
const ReturnUrlTypeDecline = "decline"
const ReturnUrlTypeException = "exception"
const ReturnUrlTypeFaxPending = "fax_pending"
const ReturnUrlTypeIdCheckFailed = "id_check_failed"
const ReturnUrlTypeSessionTimeout = "session_timeout"
const ReturnUrlTypeSigningComplete = "signing_complete"
const ReturnUrlTypeTTLExpired = "ttl_expired"
const ReturnUrlTypeViewComplete = "view_complete"

// EnvUrl is a url for an envelope view.
type EnvUrl struct {
	Url string `json:"url,omitempty"`
}

// EnvRecipientView is used to create a url for a recipient.
// See PostRecipientView
type EnvRecipientView struct {
	ClientUserId          string        `json:"clientUserId,omitempty"`
	AuthenticationMethod  string        `json:"authenticationMethod,omitempty"`
	AssertionId           string        `json:"assertionId,omitempty"`
	AuthenticationInstant string        `json:"authenticationInstant,omitempty"`
	SecurityDomain        string        `json:"securityDomain,omitempty"`
	Email                 string        `json:"email,omitempty"`
	UserId                string        `json:"userId,omitempty"`
	UserName              string        `json:"userName,omitempty"`
	ReturnUrl             ReturnUrlType `json:"returnUrl,omitempty"`
}

// FolderEnvList is the response struct for Serivice.GetFolderEnvList()
// and contains a list of envelopes in the folder
type FolderEnvList struct {
	EndPosition   string       `json:"endPosition,omitempty"`
	ResultSetSize string       `json:"resultSetSize,omitempty"`
	StartPosition string       `json:"startPosition,omitempty"`
	TotalSetSize  string       `json:"totalSetSize,omitempty"`
	TotalRows     string       `json:"totalRows,omitempty"`
	NextUri       string       `json:"nextUri,omitempty"`
	PreviousUri   string       `json:"previousUri,omitempty"`
	FolderItems   []FolderItem `json:"folderItems,omitempty"`
}

// FolderTemplateList is the response struct for Service.GetTemplateList
type FolderTemplateList struct {
	FolderEnvList
	EnvelopeTemplates []TemplateItem `json:"EnvelopeTemplates,omitempty"`
}

// FolderList is the response struct for GetFolderList()
type FolderList struct {
	Folders []Folder `json:"folders,omitempty"`
	XX      string   `json:"-"`
}

// Folder Definition and list of Child Folders
type Folder struct {
	OwnerUserName   string            `json:"ownerUserName,omitempty"`
	OwnerEmail      string            `json:"ownerEmail,omitempty"`
	OwnerUserId     string            `json:"ownerUserId,omitempty"`
	Type            string            `json:"type,omitempty"`
	Name            string            `json:"name,omitempty"`
	Uri             string            `json:"uri,omitempty"`
	ParentFolderId  string            `json:"parentFolderId,omitempty"`
	ParentFolderUri string            `json:"parentFolderUri,omitempty"`
	FolderId        string            `json:"folderId,omitempty"`
	Folders         []Folder          `json:"folders,omitempty"`
	Filter          map[string]string `json:"filter,omitempty"`
}

// FolderItem describes an envelope in a FolderEnvList
type FolderItem struct {
	Name            string        `json:"name,omitempty"`
	CreatedDateTime time.Time     `json:"createdDateTime,omitempty"`
	EnvelopeId      string        `json:"envelopeId,omitempty"`
	EnvelopeUri     string        `json:"envelopeUri,omitempty"`
	OwnerName       string        `json:"ownerName,omitempty"`
	SenderEmail     string        `json:"senderEmail,omitempty"`
	SenderName      string        `json:"senderName,omitempty"`
	SentDateTime    time.Time     `json:"sentDateTime,omitempty"`
	Status          string        `json:"status,omitempty"`
	Subject         string        `json:"subject,omitempty"`
	Recipients      RecipientList `json:"recipients,omitempty"`
}

// Notificaton is the response struct for GetEnvelopeNotification
type Notification struct {
	UseAccountDefaults string      `json:"useAccountDefaults,omitempty"`
	Reminders          *Reminder   `json:"reminders,omitempty"`
	Expirations        *Expiration `json:"expirations,omitempty"`
}

type Reminder struct {
	ReminderEnabled   string `json:"reminderEnabled,omitempty"`
	ReminderDelay     string `json:"reminderDelay,omitempty"`     // Number of days
	ReminderFrequency string `json:"reminderFrequency,omitempty"` // Number of intervals
}

type Expiration struct {
	ExpireEnabled string `json:"expireEnabled,omitempty"`
	ExpireAfter   string `json:"expireAfter,omitempty"` // Number of days until expiration
	ExpireWarn    string `json:"expireWarn,omitempty"`  // Number of days until warning
}

type BccEmail struct {
	BccEmailAddressId string `json:"bccEmailAddressId,omitempty"`
	Email             string `json:"email,omitempty"`
}

type ServerTemplate struct {
	Sequence   string `json:"sequence,omitempty"`
	TemplateId string `json:"templateId,omitempty"`
}

type InlineTemplate struct {
	Sequence  string     `json:"sequence,omitempty"`
	Documents []Document `json:"documents,omitempty"`
}
type CompositeTemplate struct {
	CompositeTemplateId         string           `json:"compositeTemplateId,omitempty"`
	ServerTemplates             []ServerTemplate `json:"serverTemplates,omitempty"`
	InlineTemplates             []InlineTemplate `json:"inlineTemplates,omitempty"`
	PdfMetaDataTemplateSequence string           `json:"pdfMetaDataTemplateSequence,omitempty"`
	Document                    []Document       `json:"document,omitempty"`
}

// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Send%20an%20Envelope.htm
type Envelope struct {
	Accessibility           string              `json:"accessibility,omitempty"`
	AllowMarkup             string              `json:"allowMarkup,omitempty"`
	AllowReassign           string              `json:"allowReassign,omitempty"`
	AllowRecipientRecursion string              `json:"allowRecipientRecursion,omitempty"`
	Asynchronous            string              `json:"asynchronous,omitempty"`
	AuthoritativeCopy       string              `json:"authoritativeCopy,omitempty"`
	AutoNavigation          string              `json:"autoNavigation,omitempty"`
	BrandId                 string              `json:"brandId,omitempty"`
	EmailBlurb              string              `json:"emailBlurb,omitempty"`
	EmailSubject            string              `json:"emailSubject,omitempty"`
	EnableWetSign           string              `json:"enableWetSign,omitempty"`
	EnforceSignerVisibility string              `json:"enforceSignerVisibility,omitempty"`
	EnvelopeIdStamping      string              `json:"envelopeIdStamping,omitempty"`
	MessageLock             string              `json:"messageLock,omitempty"`
	Notification            *Notification       `json:"notification,omitempty"`
	RecipientsLock          string              `json:"recipientsLock,omitempty"`
	SigningLocation         string              `json:"signingLocation,omitempty"`
	Status                  string              `json:"status,omitempty"`
	TransactionId           string              `json:"transactionId,omitempty"`
	UseDisclosure           bool                `json:"useDisclosure,omitempty"`
	CustomFields            *CustomFieldList    `json:"customFields,omitempty"`
	Documents               []Document          `json:"documents,omitempty"`
	Recipients              *RecipientList      `json:"recipients,omitempty"`
	EventNotification       *EventNotification  `json:"eventNotification,omitempty"`
	EmailSettings           *EmailSetting       `json:"emailSettings,omitempty"`
	TemplateId              string              `json:"templateId,omitempty"`
	TemplateRoles           []TemplateRole      `json:"templateRoles,omitempty"`
	CompositeTemplates      []CompositeTemplate `json:"compositeTemplates,omitempty"`
}

type CustomFieldList struct {
	ListCustomFields []ListCustomField `json:"listCustomFields,omitempty"`
	TextCustomFields []CustomField     `json:"textCustomFields,omitempty"`
}

type CustomField struct {
	Id           string         `json:"fieldId,omitempty"`
	Name         string         `json:"name,omitempty"`
	Required     string         `json:"required,omitempty"`
	Show         string         `json:"show,omitempty"`
	Value        string         `json:"value,omitempty"`
	ErrorDetails *ResponseError `json:errorDetails,omitempty"`
}

type ListCustomField struct {
	CustomField
	ListItems []string `json:"listItems,omitempty"`
}

// RestApi documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Document%20Parameters.htm
type Document struct {
	Name                    string     `json:"name,omitempty"`
	DocumentId              string     `json:"documentId,omitempty"`
	RemoteUrl               string     `json:"remoteUrl,omitempty"`
	Order                   string     `json:"order,omitempty"`
	TransformPdfFields      string     `json:"transformPdfFields,omitempty"`
	DocumentFields          []NmVal    `json:"documentFields,omitempty"`
	EncryptedWithKeyManager string     `json:"encryptedWithKeyManager,omitempty"`
	Pages                   string     `json:"pages,omitempty"`
	FileExtension           string     `json:"fileExtension,omitempty"`
	DocumentBase64          string     `json:"documentBase64,omitempty"`
	Matchboxes              []Matchbox `json:"matchboxes,omitempty"`
}

// Matchbox describes the area used for template matching
type Matchbox struct {
	PageNumber string `json:"pageNumber,omitempty"`
	XPosition  string `json:"xPosition,omitempty"`
	YPosition  string `json:"yPosition,omitempty"`
	Width      string `json:"width,omitempty"`
	Height     string `json:"height,omitempty"`
}

type DocumentList struct {
	//EnvelopeId string     `json:"envelopeId,omitempty"`
	Documents []Document `json:"documents,omitempty"`
}

type DocumentFieldList struct {
	DocumentFields []CustomDocumentField `json:documentFields,omitempty"`
}

type CustomDocumentField struct {
	NmVal
	ErrorDetails *ResponseError `json:errorDetails,omitempty"`
}

type EventNotification struct {
	Url                               string           `json:"url,omitempty"`
	LoggingEnabled                    string           `json:"loggingEnabled,omitempty"`
	RequireAcknowledgment             string           `json:"requireAcknowledgment,omitempty"`
	UseSoapInterface                  string           `json:"useSoapInterface,omitempty"`
	SoapNameSpace                     string           `json:"soapNameSpace,omitempty"`
	IncludeCertificateWithSoap        string           `json:"includeCertificateWithSoap,omitempty"`
	SignMessageWithX509Cert           string           `json:"signMessageWithX509Cert,omitempty"`
	IncludeDocuments                  string           `json:"includeDocuments,omitempty"`
	IncludeEnvelopeVoidReason         string           `json:"includeEnvelopeVoidReason,omitempty"`
	IncludeTimeZone                   string           `json:"includeTimeZone,omitempty"`
	IncludeSenderAccountAsCustomField string           `json:"includeSenderAccountAsCustomField,omitempty"`
	IncludeDocumentFields             string           `json:"includeDocumentFields,omitempty"`
	IncludeCertificateOfCompletion    string           `json:"includeCertificateOfCompletion,omitempty"`
	EnvelopeEvents                    []EnvelopeEvent  `json:"envelopeEvents,omitempty"`
	RecipientEvents                   []RecipientEvent `json:"recipientEvents,omitempty"`
}

type EnvelopeEvent struct {
	EnvelopeEventStatusCode string `json:"envelopeEventStatusCode,omitempty"`
	IncludeDocuments        string `json:"includeDocuments,omitempty"`
}

type RecipientEvent struct {
	RecipientEventStatusCode string `json:"RecipientEventStatusCode,omitempty"`
	IncludeDocuments         string `json:"includeDocuments,omitempty"`
}

// Documentation: https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Add%20Email%20Setting%20Overrides%20to%20an%20Envelope.htm
type EmailSetting struct {
	ReplyEmailAddressOverride string `json:"replyEmailAddressOverride,omitempty"`
	ReplyEmailNameOverride    string `json:"replyEmailNameOverride,omitempty"`
	BccEmailAddresses         string `json:"bccEmailAddresses,omitempty"`
}

type TemplateRole struct {
	Email              string             `json:"email,omitempty"`
	Name               string             `json:"name,omitempty"`
	RoleName           string             `json:"roleName,omitempty"`
	ClientUserId       string             `json:"clientUserId,omitempty"`
	DefaultRecipient   string             `json:"defaultRecipient,omitempty"`
	RoutingOrder       string             `json:"routingOrder,omitempty"`
	AccessCode         string             `json:"accessCode,omitempty"`
	InPersonSignerName string             `json:"inPersonSignerName,omitempty"`
	EmailNotification  *EmailNotification `json:"emailNotification,omitempty"`
	Tabs               *Tabs              `json:"tabs,omitempty"`
}

type EnvelopeList struct {
	Envelopes []EnvelopeUris `json:"envelopes"`
}
type EnvelopeUris struct {
	AllowReassign         string    `json:"allowReassign,omitempty"`
	CertificateUri        string    `json:"certificateUri,omitempty"`
	CreatedDateTime       time.Time `json:"createdDateTime,omitempty"`
	CustomFieldsUri       string    `json:"customFieldsUri,omitempty"`
	DocumentsCombinedUri  string    `json:"documentsCombinedUri,omitempty"`
	DocumentsUri          string    `json:"documentsUri,omitempty"`
	EmailBlurb            string    `json:"emailBlurb,omitempty"`
	EmailSubject          string    `json:"emailSubject,omitempty"`
	EnableWetSign         string    `json:"enableWetSign,omitempty"`
	EnvelopeId            string    `json:"envelopeId,omitempty"`
	EnvelopeUri           string    `json:"envelopeUri,omitempty"`
	LastModifiedDateTime  time.Time `json:"lastModifiedDateTime,omitempty"`
	NotificationUri       string    `json:"notificationUri,omitempty"`
	PurgeState            string    `json:"purgeState,omitempty"`
	RecipientsUri         string    `json:"recipientsUri,omitempty"`
	Status                string    `json:"status,omitempty"`
	StatusChangedDateTime time.Time `json:"statusChangedDateTime,omitempty"`
	TemplatesUri          string    `json:"templatesUri,omitempty"`
}

type AuditEventList struct {
	AuditEvents []AuditEvent `json:"auditEvents,omitempty"`
}

type AuditEvent struct {
	EventFields []NmVal `json:"eventFields,omitempty"`
}

type DocumentAsset struct {
	Name         string         `json:"name,omitempty"`
	Type         string         `json:"type,omitempty"`
	DocumentId   string         `json:"documentId,omitempty"`
	Order        string         `json:"order,omitempty"`
	Pages        string         `json:"pages,omitempty"`
	Uri          string         `json:"uri:omitempty"`
	ErrorDetails *ResponseError `json:errorDetails,omitempty"`
}

type DocumentAssetList struct {
	EnvelopeId        string          `json:"EnvelopeId,omitempty"`
	EnvelopeDocuments []DocumentAsset `json:"envelopeDocuments,omitempty"`
}

type LoginInfo struct {
	ApiPassword string `json:"apiPassword,omitempty"`
}

type LoginAccount struct {
	AccountId            string  `json:"accountId,omitempty"`
	BaseUrl              string  `json:"baseUrl,omitempty"`
	Email                string  `json:"email,omitempty"`
	IsDefault            string  `json:"isDefault,omitempty"`
	LoginAccountSettings []NmVal `json:"loginAccountSettings,omitempty"`
	LoginUserSettings    []NmVal `json:"loginUserSettings,omitempty"`
	Name                 string  `json:"Name,omitempty"`
	SiteDescription      string  `json:"siteDescription,omitempty"`
	UserId               string  `json:"userId,omitempty"`
	UserName             string  `json:"userName,omitempty"`
}

type EnvelopeResponse struct {
	EnvelopeId     string    `json:"envelopeId,omitempty"`
	Status         string    `json:"status,omitempty"`
	StatusDateTime time.Time `json:"statusDateTime,omitempty"`
	Uri            string    `json:"uri,omitempty"`
}

// Return structure for GetEnvelopeTemplate call
type TemplateList struct {
	Templates []TemplateItem `json:"templates,omitempty"`
}

type TemplateItem struct {
	Name       string `json:"name,omitempty"`
	TemplateId string `json:"templateId,omitempty"`
	Uri        string `json:"uri,omitempty"`
}

type Template struct {
	EnvelopeTemplateDefinition TemplateDefinition `json:"envelopeTemplateDefinition,omitempty"`
	Accessibility              string             `json:"accessibility,omitempty"`
	AllowMarkup                string             `json:"allowMarkup,omitempty"`
	AllowReassign              string             `json:"allowReassign,omitempty"`
	AllowRecipientRecursion    string             `json:"allowRecipientRecursion,omitempty"`
	Asynchronous               string             `json:"asynchronous,omitempty"`
	AuthoritativeCopy          string             `json:"authoritativeCopy,omitempty"`
	AutoNavigation             string             `json:"autoNavigation,omitempty"`
	BrandId                    string             `json:"brandId,omitempty"`
	EmailBlurb                 string             `json:"emailBlurb,omitempty"`
	EmailSubject               string             `json:"emailSubject,omitempty"`
	EnableWetSign              string             `json:"enableWetSign,omitempty"`
	EnforceSignerVisibility    string             `json:"enforceSignerVisibility,omitempty"`
	EnvelopeIdStamping         string             `json:"envelopeIdStamping,omitempty"`
	MessageLock                string             `json:"messageLock,omitempty"`
	Notification               *Notification      `json:"notification,omitempty"`
	RecipientsLock             string             `json:"recipientsLock,omitempty"`
	SigningLocation            string             `json:"signingLocation,omitempty"`
	CustomFields               *CustomFieldList   `json:"customFields,omitempty"`
	Documents                  []Document         `json:"documents,omitempty"`
	Recipients                 *RecipientList     `json:"recipients,omitempty"`
	EventNotification          *EventNotification `json:"eventNotification,omitempty"`
}

type TemplateDefinition struct {
	TemplateId     string             `json:"templateId,omitempty"`
	Name           string             `json:"name,omitempty"`
	Shared         string             `json:"shared,omitempty"`
	Password       string             `json:"password,omitempty"`
	Description    string             `json:"description,omitempty"`
	LastModified   time.Time          `json:"lastModified,omitempty"`
	LastModifiedBy TemplateModifiedBy `json:"lastModifiedBy,omitempty"`
	PageCount      int                `json:"pageCount,omitempty"`
	FolderName     string             `json:"folderName,omitempty"`
	FolderId       string             `json:"folderId,omitempty"`
	Owner          TemplateOwner      `json:"owner,omitempty"`
}

type TemplateModifiedBy struct {
	UserName string `json:"userName,omitempty"`
	UserId   string `json:"userId,omitempty"`
	Email    string `json:"email,omitempty"`
	URI      string `json:"uri,omitempty"`
}

type TemplateOwner struct {
	UserName   string `json:"userName,omitempty"`
	UserId     string `json:"userId,omitempty"`
	Email      string `json:"email,omitempty"`
	URI        string `json:"uri,omitempty"`
	UserType   string `json:"userType,omitempty"`
	UserStatus string `json:"userStatus,omitempty"`
}
