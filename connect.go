// Copyright 2015 James Cote and Liberty Fund, Inc.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// connect.go contains definitions for Docusign's Connect
// Service.
// see https://www.docusign.com/developer-center/explore/connect
package docusign

import (
	"strings"
	"time"
)

// DSTime handles the multiple datetime formats that DS returns
// in the Connect service.
type DSTime string

func (d DSTime) Time() (tm time.Time) {
	if strings.HasSuffix(string(d), "Z") {
		tm, _ = time.Parse(time.RFC3339Nano, string(d))

	} else {
		tm, _ = time.Parse("2006-01-02T15:04:05.999999999", string(d))
	}
	return tm
}

// Connect Data is the top level struct for a Connect status message.
type ConnectData struct {
	EnvelopeStatus EnvelopeStatusXml `xml:"EnvelopeStatus" json:"envelopeStatus,omitempty"`
	DocumentPdfs   []DocumentPdfXml  `xml:"DocumentPDFs>DocumentPDF" json:"documentPdfs,omitempty"`
}

// EnvelopeStatusXML contains envelope information.
type EnvelopeStatusXml struct {
	TimeGenerated      DSTime               `xml:"TimeGenerated" json:"timeGenerated,omitempty"`
	EnvelopeID         string               `xml:"EnvelopeID" json:"envelopeID,omitempty"`
	Subject            string               `xml:"Subject" json:"subject,omitempty"`
	UserName           string               `xml:"UserName" json:"userName,omitempty"`
	Email              string               `xml:"Email" json:"email,omitempty"`
	Status             string               `xml:"Status" json:"status,omitempty"`
	Created            DSTime               `xml:"Created" json:"created,omitempty"`
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

// DocumentPdfXML contains the base64 enocded pdf files.
type DocumentPdfXml struct {
	Name     string `xml:"Name" json:"name,omitempty"`
	PDFBytes string `xml:"PDFBytes" json:"pdfBytes,omitempty"`
}

// RecipientStatusXml contains data describing each recipient for an envelope.
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

// DocumentStatusXml describes each document in the envelope.
type DocumentStatusXml struct {
	ID           string `xml:"ID" json:"id,omitempty"`
	Name         string `xml:"Name"RecipientAttachment json:"name,omitempty"`
	TemplateName string `xml:"TemplateName" json:"templateName,omitempty"`
	Sequence     string `xml:"Sequence" json:"sequence,omitempty"`
}

// NmValXml contains the key/value pair for the form values.
type NmValXml struct {
	Name  string `xml:"name,attr" json:"name,omitempty"`
	Value string `xml:"value" json:"value,omitempty"`
}

type CustomFieldXml struct {
	Name     string `xml:"Name" json:"name,omitempty"`
	Value    string `xml:"Value" json:"value,omitempty"`
	Show     bool   `xml:"Show" json:"show,omitempty"`
	Required bool   `xml:"Required" json:"required,omitempty"`
}

type RecipientAttachmentXml struct {
	Data  string `xml:"Data" json:"data,omitempty"`
	Label string `xml:"Label" json:"label,omitempty"`
}

// TabStatusXML describes the properties of each recipient tab
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
