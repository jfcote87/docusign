package docusign

import (
	"strings"
	"time"
)

type DSTime string

func (d DSTime) Time() (time.Time, error) {
	if strings.HasSuffix(string(d), "Z") {
		return time.Parse(time.RFC3339Nano, string(d))

	}
	return time.Parse("2006-01-02T15:04:05.999999999", string(d))
}

type DocuSignEnvelopeInformation struct {
	EnvelopeStatus EnvelopeStatusXml `xml:"EnvelopeStatus" json:"envelopeStatus,omitempty"`
	DocumentPdfs   []DocumentPdfXml  `xml:"DocumentPDFs>DocumentPDF" json:"documentPdfs,omitempty"`
}

type EnvelopeStatusXml struct {
	TimeGenerated      DSTime               `xml:"TimeGenerated" json:"timeGenerated,omitempty"`
	EnvelopeID         string               `xml:"EnvelopeID" json:"envelopeID,omitempty"`
	Subject            string               `xml:"Subject" json:"subject,omitempty"`
	UserName           string               `xml:"UserName" json:"userName,omitempty"`
	Email              string               `xml:"Email" json:"email,omitempty"`
	Status             string               `xml:"Status" json:"status,omitempty"`
	Created            DSTime               `xml:"Created" json:"created,omitempty`
	Sent               DSTime               `xml:"Sent" json:"sent,omitempty"`
	Delivered          DSTime               `xml:"Delivered" json:"delivered,omitempty"`
	Signed             DSTime               `xml:"Signed" json:"signed,omitempty"`
	Completed          DSTime               `xml:"Completed" json:"completed,omitempty"`
	ACStatus           string               `xml:"ACStatus" json:"acStatus,omitempty"`
	ACStatusDate       string               `xml:"ACStatusDate" json:"acStatusDate,omitempty"`
	ACHolder           string               `xml:"ACHolder" json:"acHolder,omitempty"`
	ACHolderEmail      string               `xml:"ACHolderEmail" json:"acHolderEmail,omitempty"`
	ACHolderLocation   string               `xml:"ACHolderLocation" json:"acHolderLocation,omitempty"`
	SigningLocation    string               `xml:"SigningLocation" json:"signingLocation,omitempty"`
	SenderIPAddress    string               `xml:"SenderIPAddress" json:"senderIPAddress,omitempty"`
	EnvelopePDFHash    string               `xml:"EnvelopePDFHash" json:"envelopePDFHash,omitempty"`
	AutoNavigation     bool                 `xml:"AutoNavigation" json:"autoNavigation,omitempty"`
	EnvelopeIdStamping bool                 `xml:"EnvelopeIdStamping" json:"envelopeIdStamping,omitempty"`
	AuthoritativeCopy  bool                 `xml:"AuthoritativeCopy" json:"authoritativeCopy,omitempty"`
	RecipientStatuses  []RecipientStatusXml `xml:"RecipientStatuses>RecipientStatus" json:"recipientStatuses,omitempty"`
	CustomFields       []CustomFieldXml     `xml:"CustomFields>CustomField" json:"customFields,omitempty"`
	DocumentStatuses   []DocumentStatusXml  `xml:"DocumentStatuses>DocumentStatus" json:"documentStatuses,omitempty"`
}

type DocumentPdfXml struct {
	Name     string `xml:"Name" json:"name,omitempty"`
	PDFBytes string `xml:"PDFBytes" json:"pdfBytes,omitempty"`
}

type RecipientStatusXml struct {
	Type                string                 `xml:"Type" json:"type,omitempty"`
	Email               string                 `xml:"Email" json:"email,omitempty"`
	UserName            string                 `xml:"UserName" json:"userName,omitempty"`
	RoutingOrder        string                 `xml:"RoutingOrder" json:"routingOrder,omitempty"`
	Sent                DSTime                 `xml:"Sent" json:"sent,omitempty"`
	Delivered           DSTime                 `xml:"Delivered" json:"delivered,omitempty"`
	Signed              DSTime                 `xml:"Signed" json:"signed,omitempty"`
	DeclineReason       string                 `xml:"DeclineReason" json:"declineReason,omitempty"`
	Status              string                 `xml:"Status" json:"status,omitempty"`
	RecipientIPAddress  string                 `xml:"RecipientIPAddress" json:"recipientIPAdress,omitempty"`
	CustomFields        []CustomFieldXml       `xml:"CustomFields>CustomField" json:"customFields,omitempty"`
	AccountStatus       string                 `xml:"AccountStatus" json:"accountStatus,omitempty"`
	RecipientId         string                 `xml:"RecipientId" json:"recipientId,omitempty"`
	TabStatuses         []TabStatusXml         `xml:"TabStatuses>TabStatus" json:"tabStatuses,omitempty"`
	FormData            []NmValXml             `xml:"FormData>xfdf>fields>field" json:"formData,omitempty"`
	RecipientAttachment RecipientAttachmentXml `xml:"RecipientAttachment>Attachment" json:"recipientAttachment,omitempty"`
}

type FormDataXml struct {
	Fields []string `xml:"Recipient json:"fields,omitempty"`
}

type DocumentStatusXml struct {
	ID           string `xml:"ID" json:"id,omitempty"`
	Name         string `xml:"Name"RecipientAttachment json:"name,omitempty"`
	TemplateName string `xml:"TemplateName" json:"templateName,omitempty"`
	Sequence     string `xml:"Sequence" json:"sequence,omitempty"`
}
type NmValXml struct {
	Name  string `xml:"name,attr" json:"name,omitempty"`
	Value string `xml:"value" json:"value,omitempty"`
}

type CustomFieldXml struct {
	Name     string `xml:"Name" json:"name,omitempty"`
	Value    string `xml:"value" json:"value,omitempty"`
	Show     bool   `xml:"Show" json:"show,omitempty"`
	Required bool   `xml:"Required" json:"required,omitempty"`
}

type RecipientAttachmentXml struct {
	Data  string `xml:"Data" json:"data,omitempty"`
	Label string `xml:"Label" json:"label,omitempty"`
}

type TabStatusXml struct {
	TabType           string `xml:"TabType" json:"tabType,omitempty"`
	Status            string `xml:"Status" json:"status,omitempty"`
	XPosition         string `xml:"XPosition" json:"xPosition,omitempty"`
	YPosition         string `xml:"YPosition" json:"yPosition,omitempty"`
	TabLabel          string `xml:"TabLabel" json:"tabLabel,omitempty"`
	TabName           string `xml:"TabName" json:"tabName,omitempty"`
	TabValue          string `xml:"TabValue" json:"tabValue,omitempty"`
	DocumentID        string `xml:"DocumentID" json:"documentID,omitempty"`
	PageNumber        string `xml:"PageNumber" json:"pageNumber,omitempty"`
	OriginalValue     string `xml:"OriginalValue" json:"originalValue,omitempty"`
	ValidationPattern string `xml:"ValidationPattern" json:"validationPattern,omitempty"`
	ListValues        string `xml:"ListValues" json:"listValues,omitempty"`
	ListSelectedValue string `xml:"ListSelectedValue" json:"listSelectedValue,omitempty"`
	CustomTabType     string `xml:"CustomTabType" json:"customTabType,omitempty"`
}
