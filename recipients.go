// Copyright 2015 James Cote and Liberty Fund, Inc.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package docusign

// RecipientList defines the recipients for an envelope
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Recipient%20Parameter.htm
type RecipientList struct {
	Agents              []Agent             `json:"agents,omitempty"`
	CarbonCopies        []CarbonCopy        `json:"carbonCopies,omitempty"`
	CertifiedDeliveries []CertifiedDelivery `json:"certifiedDeliveries,omitempty"`
	Editors             []Editor            `json:"editors,omitempty"`
	InPersonSigners     []InPersonSigner    `json:"inPersonSigners,omitempty"`
	Intermediaries      []Intermediary      `json:"intermediaries,omitempty"`
	Signers             []Signer            `json:"signers,omitempty"`
	RecipientCount      string              `json:"recipientCount,omitempty"`
}

// Values returns a NmVal slice contiaing the
// tabLabel and value for each tab in the RecipientList.
func (r RecipientList) Values() []NmVal {
	v := make([]NmVal, 0)
	for _, x := range r.InPersonSigners {
		v = append(v, x.Tabs.Values()...)
	}
	for _, x := range r.Signers {
		v = append(v, x.Tabs.Values()...)
	}
	return v
}

// EmailNotification contains the email message sent to a
// recipient.  If not set, the envelopes EmailBlurb and
// EmailSubject are used.
type EmailNotification struct {
	EmailBody         string `json:"emailBody,omitempty"`
	EmailSubject      string `json:"emailSubject,omitempty"`
	SupportedLanguage string `json:"supportedLanguage,omitempty"`
}

// IDCheckInformationInput specifies authentication check by name. See api
// documentation for specific values.
type IDCheckInformationInput struct {
	AddressInformationInput *AddressInformationInput `json:"addressInformationInput,omitempty"`
	DobInformationInput     *DobInformationInput     `json:"dobInformationInput,omitempty"`
	Ssn4InformationInput    *Ssn4InformationInput    `json:"ssn4InformationInput,omitempty"`
	Ssn9InformationInput    *Ssn9InformationInput    `json:"ssn9InformationInput,omitempty"`
}

type InformationInput struct {
	DisplayLevelCode  string `json:"displayLevelCode,omitempty"`
	ReceiveInResponse string `json:"receiveInResponse,omitempty"`
}

type AddressInformationInput struct {
	InformationInput
	AddressInformation *AddressInformation `json:"addressInformation,omitempty"`
}

type AddressInformation struct {
	Street1  string `json:"street1,omitempty"`
	Street2  string `json:"street2,omitempty"`
	City     string `json:"city,omitempty"`
	State    string `json:"state,omitempty"`
	Zip      string `json:"zip,omitempty"`
	ZipPlus4 string `json:"zipPlus4,omitempty"`
}

type DobInformationInput struct {
	InformationInput
	DateOfBirth string `json:"dateOfBirth,omitempty"`
}

type Ssn4InformationInput struct {
	InformationInput
	Ssn4 string `json:"ssn4,omitempty"`
}

type Ssn9InformationInput struct {
	InformationInput
	Ssn9 string `json:"ssn9,omitempty"`
}

type PhoneAuthentication struct {
	RecipMayProvideNumber       string   `json:"recipMayProvideNumber,omitempty"`
	ValidateRecipProvidedNumber string   `json:"validateRecipProvidedNumber,omitempty"`
	RecordVoicePrint            string   `json:"recordVoicePrint,omitempty"`
	SenderProvidedNumbers       []string `json:"senderProvidedNumbers,omitempty"`
}

type SamlAuthentication struct {
	SamlAssertionAttributes []NmVal `json:"samlAssertionAttributes,omitempty"`
}

// RecipientAttachement will be used to a specific file attachment
// for a recipient.
type RecipientAttachment struct {
	Label          string `json:"label,omitempty"`
	AttachmentType string `json:"attachmentType,omitempty"`
	Data           string `json:"label,omitempty"`
}

type SmsAuthentication struct {
	SenderProvidedNumbers []string `json:"smsAuthentication,omitempty"`
}

// Recipient contains the common fields for all recipient types
type Recipient struct {
	Name                                  string               `json:"name,omitempty"`
	AccessCode                            string               `json:"accessCode,omitempty"`
	AddAccessCodeToEmail                  DSBool               `json:"addAccessCodeToEmail,omitempty"`
	ClientUserId                          string               `json:"clientUserId,omitempty"`
	EmbeddedRecipientStartURL             string               `json:"embeddedRecipientStartURL,omitempty"`
	CustomFields                          string               `json:"customFields,omitempty"`
	EmailNotification                     *EmailNotification   `json:"emailNotification,omitempty"`
	ExcludedDocuments                     string               `json:"excludedDocuments,omitempty"`
	IdCheckConfigurationName              string               `json:"idCheckConfigurationName,omitempty"`
	IDCheckInformationInput               string               `json:"iDCheckInformationInput,omitempty"`
	InheritEmailNotificationConfiguration DSBool               `json:"inheritEmailNotificationConfiguration,omitempty"`
	Note                                  string               `json:"note,omitempty"`
	PhoneAuthentication                   *PhoneAuthentication `json:"phoneAuthentication,omitempty"`
	RecipientAttachments                  *RecipientAttachment `json:"recipientAttachment,omitempty"`
	RecipientCaptiveInfo                  string               `json:"recipientCaptiveInfo,omitempty"`
	RecipientId                           string               `json:"recipientId,omitempty"`
	RequireIdLookup                       DSBool               `json:"requireIdLookup,omitempty"`
	RoleName                              string               `json:"roleName,omitempty"`
	RoutingOrder                          string               `json:"routingOrder,omitempty"`
	SamlAuthentication                    *SamlAuthentication  `json:"samlAuthentication,omitempty"`
	SmsAuthentication                     *SmsAuthentication   `json:"smsAuthentication,omitempty"`
	SocialAuthentications                 DSBool               `json:"socialAuthentications,omitempty"`
	TemplateAccessCodeRequired            DSBool               `json:"templateAccessCodeRequired,omitempty"`
	TemplateLocked                        DSBool               `json:"templateLocked,omitempty"`
	TemplateRequired                      DSBool               `json:"templateRequired,omitempty"`
	ErrorDetails                          *ResponseError       `json:"errorDetails,omitempty"`
}

// EmailRecipient adds email field to base recipient structure
type EmailRecipient struct {
	Recipient
	Email string `json:"email,omitempty"`
}

// Agent can add name and email information for recipients that appear after the recipient in routing order.
// RestApi Documetation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Recipients/Agent%20Recipient.htm
type Agent struct {
	EmailRecipient
	CanEditRecipientEmails DSBool `json:"canEditRecipientEmails,omitempty"`
	CanEditRecipientNames  DSBool `json:"canEditRecipientNames,omitempty"`
}

// CarbonCopy receives a copy of the envelope when the envelope reaches the recipient’s order in the process flow and when the envelope is completed.
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Recipients/Carbon%20Copies%20Recipient.htm
type CarbonCopy struct {
	EmailRecipient
}

// CertifiedDeliveryr receives the completed documents for the envelope to be completed
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Recipients/Certified%20Deliveries%20Recipient.htm
type CertifiedDelivery struct {
	EmailRecipient
	CanEditRecipientEmails DSBool `json:"canEditRecipientEmails,omitempty"`
	CanEditRecipientNames  DSBool `json:"canEditRecipientNames,omitempty"`
}

// Editor can add name and email information, add or change the routing order and set authentication options for the remaining recipients.
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Recipients/Editors%20Recipient.htm
type Editor struct {
	EmailRecipient
	CanEditRecipientEmails DSBool `json:"canEditRecipientEmails,omitempty"`
	CanEditRecipientNames  DSBool `json:"canEditRecipientNames,omitempty"`
}

// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Recipients/In%20Person%20Signers%20Recipient.htm
type InPersonSigner struct {
	Recipient
	BaseSigner
	HostEmail string `json:"hostEmail,omitempty"`
	HostName  string `json:"hostName,omitempty"`
}

// This recipient can, but is not required to, add name and email information for recipients at the same or subsequent level in the routing order
//
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Recipients/Intermediaries%20Recipient.htm
type Intermediary struct {
	EmailRecipient
	CanEditRecipientEmails DSBool `json:"canEditRecipientEmails,omitempty"`
	CanEditRecipientNames  DSBool `json:"canEditRecipientNames,omitempty"`
}

// BaseSigner contains common fields of all signer types
type BaseSigner struct {
	AutoNavigation     string `json:"autoNavigation,omitempty"`
	DefaultRecipient   string `json:"defaultRecipient,omitempty"`
	SignInEachLocation string `json:"signInEachLocation,omitempty"`
	SignerEmail        string `json:"signerEmail,omitempty"`
	SignerName         string `json:"signerName,omitempty"`
	Tabs               *Tabs  `json:"tabs,omitempty"`
}

// Use this action if your recipient must sign, initial, date or add data to form fields on the documents in the envelope.
//
// RestApi Documentation
// https://www.docusign.com/p/RESTAPIGuide/Content/REST%20API%20References/Recipients/Signers%20Recipient.htm
type Signer struct {
	EmailRecipient
	BaseSigner
	IsBulkRecipient   string            `json:"isBulkRecipient,omitempty"`
	BulkRecipientsUri string            `json:"bulkRecipientsUri,omitempty"`
	DeliveryMethod    string            `json:"deliveryMethod,omitempty"`
	DeliveredDateTime string            `json:"deliveredDateTime,omitempty"`
	SignedDateTime    string            `json:"signedDateTime,omitempty"`
	OfflineAttributes map[string]string `json:"offlineAttributes,omitempty"`
}

// RecipeintUpdateResult is returned via the RecipientsModify call and returns
// a list of recipient ids and a corresponding error detail for each modification.
type RecipientUpdateResult struct {
	recipientUpdateResults []Recipient `json:"recipientUpdateResults"`
}
