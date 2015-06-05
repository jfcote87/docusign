// Copyright 2015 James Cote and Liberty Fund, Inc.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// docusign implements a service to use the version 2 Docusign
// rest api. Api documentation may be found at:
// https://www.docusign.com/p/RESTAPIGuide/RESTAPIGuide.htm

// You must define an environment variable name DOCSUSIGN_CONFIG for the
// test to run properly.  The variable must be a json representation of a
// Config struct that provides the api key, user name, password and account id.
// An example is shown below:
//
// DOCUSIGN_CONFIG={"key":"XXXX-XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX","user":"xxxxx@xxxxx.xxx","pwd":"xxxxxxx","acctId":"XXXXXX"}
//
// If you wish to skip generating an oauth2 token, you may define an environment
// variable named DOCUSIGN_TOKEN which contains an existing token.
//
// A draft envelope will be created in the Docusign demo environment with the subject "Created by Go Test".
package docusign

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
	"time"

	"golang.org/x/net/context"
)

type InviteData struct {
	ConfName       string `json:"confName" xml:"ConfName" ds:"txtConfName,Fellow"`
	Location       string `json:"location" xml:"Location" ds:"txtLocation,Fellow"`
	ConfDate       string `json:"confDate" xml:"ConfDate" ds:"txtConfDate,Fellow"`
	Pid            int    `json:"pid" xml:"Pid" ds:"-"`
	Prefix         string `json:"prefix" xml:"Prefix" ds:"txtPrefix,Conferee"`
	LastName       string `json:"lName" xml:"LName" ds:"txtLName,Conferee"`
	FirstName      string `json:"fName" xml:"FName" ds:"txtFName,Conferee"`
	Middle         string `json:"middle" xml:"Middle" ds:"txtMiddle,Conferee"`
	Suffix         string `json:"suffix" xml:"Suffix" ds:"txtSuffix,Conferee"`
	Email          string `json:"email" xml:"Email" ds:"txtEmail,Conferee,email"`
	Aff            string `json:"aff" xml:"Aff" ds:"txtAff,Conferee"`
	AffTitle       string `json:"affTitle" xml:"AffTitle" ds:"txtAffTitle,Conferee"`
	HomeAddr       string `json:"homeAddr" xml:"HomeAddr" ds:"txtHomeAddr,Conferee"`
	HomePhone      string `json:"homePhone" xml:"HomePhone" ds:"txtHomePhone,Conferee"`
	HomeFax        string `json:"homeFax" xml:"HomeFax" ds:"txtHomeFax,Conferee"`
	BusAddr        string `json:"busAddr" xml:"BusAddr" ds:"txtBusAddr,Conferee"`
	BusPhone       string `json:"busPhone" xml:"BusPhone" ds:"txtBusPhone,Conferee"`
	BusFax         string `json:"busFax" xml:"BusFax" ds:"txtBusFax,Conferee"`
	Correspondence string `json:"correspondence" xml:"Correspondence" ds:"txtCorrespondence,Conferee,dl"`
	ConfereeList   string `json:"confereeList" xml:"ConfereeList" ds:"txtConfereeList,Conferee,dl"`
	Payment        string `json:"payment" xml:"Payment" ds:"txtPayment,Conferee,dl"`
	ListEmail      string `json:"listEmail" xml:"ListEmail" ds:"txtListEmail,Conferee,email"`
	TableTent      string `json:"TableTent" xml:"TableTent" ds:"txtTableTent,Conferee"`
	Diet           string `json:"diet" xml:"Diet" ds:"txtDiet,Conferee"`
	Hotel          string `json:"hotel" xml:"Hotel" ds:"txtHotel,Conferee"`
	AnythingElse   string `json:"anythingElse" xml:"AnythingElse" ds:"txtAnythingElse,Conferee"`
	Govt           string `json:"govt" xml:"Govt" ds:"rbUSCit,Conferee,rb"`
	Country        string `json:"country" xml:"Country" ds:"txtCountry,Conferee"`
	SSN            bool   `json:"ssn" xml:"SSN" ds:"cbSSN,Conferee,cb"`
}

var testEnvId = "04885769-590b-4412-8c4e-a527b4da1925"
var testTemplateId = "6b000a67-42cb-4d9a-820a-afa19c0e43cf"

func TestMain(m *testing.M) {
	//LogRawResponse = true

	UseDemoServer()
	DefaultCtx = context.WithValue(DefaultCtx, HTTPClient, http.DefaultClient)
	os.Exit(m.Run())
}

func TestCalls(t *testing.T) {
	t.Skip()
	var cfg *Config
	ctx := DefaultCtx

	testConfigString := os.Getenv("DOCUSIGN_CONFIG")

	err := json.Unmarshal([]byte(testConfigString), &cfg)
	if err != nil {
		t.Errorf("Unable to unmarshal DOCUSIGN_CONFIG: %v", err)
		return
	}

	if cfg.UserName == "" || cfg.Password == "" || cfg.IntegratorKey == "" || cfg.AccountId == "" {
		t.Errorf("Invalid Config")
		return

	}

	testToken := os.Getenv("DOCUSIGN_TOKEN")
	var c *OauthCredential

	if testToken > "" {
		c = &OauthCredential{AccessToken: testToken}
	} else {
		c, err = cfg.OauthCredential(DefaultCtx)
		if err != nil {
			t.Errorf("Ouauth2 token fail: %v", err)
			return
		}
		t.Logf("Token: %s\n", c.AccessToken)
		defer func() {
			if err := c.Revoke(DefaultCtx); err != nil {
				t.Errorf("Revoke token failed: %v", err)
			}
		}()
	}
	sv := New(cfg.AccountId, c)

	_, err = sv.GetTemplate(ctx, testTemplateId)
	if err != nil {
		t.Errorf("GetTemplate: %v", err)
		return
	}

	r, err := sv.TemplateSearch(DefaultCtx)
	if err != nil {
		t.Errorf("TemplateSearch: %v", err)
		return
	}

	for _, et := range r.EnvelopeTemplates {
		t.Logf("%s %s", et.TemplateId, et.Name)
	}

	// Get Draft Folder
	//LogRawRequest = true
	var draftFolder string
	fl, err := sv.FolderList(DefaultCtx, FolderTemplatesIncluede)
	if err != nil {
		t.Errorf("GetFolderList: %v", err)
		return
	}
	for _, fd := range fl.Folders {
		fmt.Printf("Folder: %s  %s\n", fd.FolderId, fd.Name)
		if fd.Name == "Draft" {
			draftFolder = fd.FolderId
		}
	}
	if draftFolder == "" {
		t.Errorf("Unable to find Draft folder")
		return
	}

	_, err = sv.AccountCustomFields(DefaultCtx)
	if err != nil {
		t.Errorf("AccountCustomFields error: %v", err)
		return
	}

	euris, err := sv.EnvelopeStatusChanges(DefaultCtx, StatusChangeToDate(time.Now()), StatusChangeFromDate(time.Now().AddDate(0, 0, -1)),
		StatusChangeStatusCode("created"), StatusChangeFromToStatus("created"), StatusChangeCustomField("PID", "123456"))
	//(time.Now().Add(time.Hour*24*-30)), StatusChangeToDate(time.Now()))
	if err != nil {
		t.Errorf("EnvelopeStatusChanges error: %v", err)
		return
	}
	fmt.Printf("Total Status changes returned - %d\n", len(euris.Envelopes))

	exs, err := sv.EnvelopeSearch(DefaultCtx, SearchFolderDrafts, EnvelopeSearchCount(3), EnvelopeSearchFromDate(time.Now().AddDate(0, -1, 0)),
		EnvelopeSearchToDate(time.Now()), EnvelopeSearchIncludeRecipients)
	if err != nil {
		t.Errorf("EnvelopeSearch error: %v", err)
		return
	}
	fmt.Printf("Total Envelopes Returned: %d %s\n", len(exs.FolderItems), exs.TotalRows)

	//STOP HERE
	return

	testEnv := testEnvelopePayload(cfg.UserName)

	file, err := os.Open("testdata/TestDocument.pdf")
	if err != nil {
		t.Errorf("Unable to open TestDocument.pdf: %v", err)
	}
	defer file.Close()
	u := &UploadFile{
		ContentType: "application/pdf",
		FileName:    "TestData.pdf",
		Id:          "1",
		Data:        file,
	}

	ex, err := sv.EnvelopeCreate(DefaultCtx, testEnv, u)
	if err != nil {
		t.Errorf("CreateEnvelope: %v", err)
		return
	}
	testEnvId = ex.EnvelopeId

	aTab := &Tabs{
		SignerAttachmentTabs: []SignerAttachmentTab{
			SignerAttachmentTab{
				BaseTab: BaseTab{
					DocumentID: "1",
					TabLabel:   "attTab",
				},
				BasePosTab: BasePosTab{
					AnchorString:  "SignatureA:",
					AnchorXOffset: "240",
					AnchorYOffset: "10",
					AnchorUnits:   "pixels",
					PageNumber:    "1",
					TabId:         "9985fd9a-a660-4ff3-983d-eb43706d496d",
				},
				BaseTemplateTab: BaseTemplateTab{
					RecipientID: "1",
				},
				Optional: true,
			},
		},
		TextTabs: []TextTab{
			TextTab{
				BaseTab: BaseTab{
					DocumentID: "1",
					TabLabel:   "deleteThisTab",
				},
				BasePosTab: BasePosTab{
					PageNumber: "1",
					XPosition:  "300",
					YPosition:  "350",
				},
				BaseTemplateTab: BaseTemplateTab{
					RecipientID: "1",
				},
			},
		},
	}
	aTab, err = sv.RecipientTabsAdd(DefaultCtx, testEnvId, "1", aTab)
	if err != nil {
		t.Errorf("Add Tabs error: %v", err)
		return
	}
	var deleteTabId string
	if len(aTab.TextTabs) == 1 {
		deleteTabId = aTab.TextTabs[0].TabId
	}

	recList, err := sv.Recipients(DefaultCtx, testEnvId, RecipientsIncludeTabs)
	if err != nil {
		t.Errorf("GetRecipients error: %v\n", err)
		return
	}
	if recList == nil || len(recList.Signers) != 2 {
		t.Errorf("Invalid recipients returned.")
		return
	}

	mTabs := &Tabs{
		RadioGroupTabs: recList.Signers[1].Tabs.RadioGroupTabs,
		ListTabs:       recList.Signers[1].Tabs.ListTabs,
		TextTabs: []TextTab{
			TextTab{Value: "ASFDAFD", BasePosTab: BasePosTab{TabId: "e611bf5f-339c-4ed0-8c71-87ec7f77fdc5"}},
		},
	}
	for i, rd := range mTabs.RadioGroupTabs[0].Radios {
		if rd.Value == "val2" {
			mTabs.RadioGroupTabs[0].Radios[i].Selected = true
		} else {
			mTabs.RadioGroupTabs[0].Radios[i].Selected = false

		}
	}

	for i, li := range mTabs.ListTabs[0].ListItems {
		xval := DSBool(false)
		if li.Value == "Y" {
			xval = true
		}
		mTabs.ListTabs[0].ListItems[i].Selected = xval
	}
	mTabs.ListTabs[0].Value = "Y Val"
	mTabs, err = sv.RecipientTabsModify(DefaultCtx, testEnvId, "2", mTabs)
	if err != nil {
		t.Errorf("Modify Tabs Error: %v", err)
		return
	}
	if len(mTabs.TextTabs) != 1 || mTabs.TextTabs[0].ErrorDetails == nil {
		t.Errorf("Wanted INVALID_TAB_OPERATION on TextTab[0]; got nil")
		return
	}

	rTabs := &Tabs{
		TextTabs: []TextTab{
			TextTab{
				BasePosTab: BasePosTab{
					TabId: deleteTabId,
				},
			},
		},
	}
	rTabs, err = sv.RecipientTabsRemove(DefaultCtx, testEnvId, "1", rTabs)
	if err != nil {
		t.Errorf("Error Deleting Tab: %v", err)
		return
	}

	newRecipients := &RecipientList{
		Signers: []Signer{
			Signer{
				EmailRecipient: EmailRecipient{
					Email: "extraRep@example.com",
					Recipient: Recipient{
						Name:              "Extra Name",
						Note:              "This is the ,Note for Extra Name",
						EmailNotification: &EmailNotification{EmailBody: "This is the recipient 3 email blurb", EmailSubject: "This is the Subject for recipient 3"},
						RecipientId:       "3",
						RoleName:          "Role3",
						RoutingOrder:      "6",
					},
				},
			},
		},
		CarbonCopies: []CarbonCopy{
			CarbonCopy{
				EmailRecipient: EmailRecipient{
					Email: "cc@example.com",
					Recipient: Recipient{
						Name:              "CC Name",
						Note:              "This is the ,Note for CCName",
						EmailNotification: &EmailNotification{EmailBody: "This is the recipient 4 email blurb", EmailSubject: "This is the Subject for recipient 4"},
						RecipientId:       "4",
						RoleName:          "Role4",
						RoutingOrder:      "5",
					},
				},
			},
		},
	}

	newRecipients, err = sv.RecipientsAdd(DefaultCtx, testEnvId, newRecipients)
	if err != nil {
		t.Errorf("Recipients Add Error: %v", err)
		return
	}

	for i := range newRecipients.Signers {
		if newRecipients.Signers[i].RecipientId == "3" {
			fmt.Printf("Here\n")
			newRecipients.Signers[i].Name = "Modified Name"
		}
	}
	modRec, err := sv.RecipientsModify(DefaultCtx, testEnvId, newRecipients)
	if err != nil {
		t.Errorf("Recipients Modify Error: %v", err)
		return
	}
	for _, rur := range modRec.recipientUpdateResults {
		if rur.ErrorDetails != nil && rur.ErrorDetails.Err == "SUCCESS" {
			continue
		}
		t.Errorf("RecipientsModify error: %v", rur.ErrorDetails)
		return
	}

	//sv.RecipientTabsRemove(ctx, envId, recipId, rl)

	//
	/*
				var draftFolder string

			ret, err := sv.GetEnvelopeStatus(DefaultCtx, testEnv)
			if err != nil {
				t.Errorf("GetEnvelop error: %v", err)
				return
			}


		recList, err := sv.GetRecipients(DefaultCtx, testEnv, GetRecipientsIncludeTabs)
		if err != nil {
			t.Errorf("GetRecipients error: %v\n", err)
			//return
		}
		fmt.Printf("Total tabs: %d\n", len(recList.Signers[1].Tabs.TextTabs))

		newRecipients := RecipientList{
			Signers: []Signer{
				Signer{
					EmailRecipient: EmailRecipient{
						Email: "extraRep@example.com",
						Recipient: Recipient{
							Name:              "Extra Name",
							Note:              "This is the ,Note for Extra Name",
							EmailNotification: &EmailNotification{EmailBody: "This is the recipient 3 email blurb", EmailSubject: "This is the Subject for recipient 3"},
							RecipientId:       "3",
							RoleName:          "Role3",
							RoutingOrder:      "6",
						},
					},
				},
			},
			CarbonCopies: []CarbonCopy{
				CarbonCopy{
					EmailRecipient: EmailRecipient{
						Email: "cc@example.com",
						Recipient: Recipient{
							Name:              "CC Name",
							Note:              "This is the ,Note for CCName",
							EmailNotification: &EmailNotification{EmailBody: "This is the recipient 4 email blurb", EmailSubject: "This is the Subject for recipient 4"},
							RecipientId:       "4",
							RoleName:          "Role4",
							RoutingOrder:      "5",
						},
					},
				},
			},
		}

		rx, err := sv.ModifyRecipients(DefaultCtx, testEnv, &newRecipients)
		if err != nil {
			t.Errorf("Add Recipeints Error: %v", err)
			return
		}
		fmt.Printf("Total Recipients: %s\n", rx.RecipientCount)

		delRec := RecipientList{
			Signers: []Signer{
				Signer{
					EmailRecipient: EmailRecipient{
						Recipient: Recipient{
							RecipientId: "3",
						},
					},
				},
			},
		}
		rx, err = sv.RemoveRecipients(DefaultCtx, testEnv, &delRec)
		if err != nil {
			t.Errorf("Del recipients Error: %v", err)
			return
		}
	*/
	/* nTabs := &Tabs{
		TextTabs: []TextTab{
			TextTab{
				BaseTab: BaseTab{
					DocumentID: "1",
					TabLabel:   "txtTextFieldY",
				},
				BasePosTab: BasePosTab{
					TabId:      "0c5261cc-ffb5-4159-bd14-82f57cc04ce0",
					XPosition:  "289",
					YPosition:  "290",
					PageNumber: "1",
				},
				BaseTemplateTab: BaseTemplateTab{
					RecipientID: "1",
				},
				Value: "YValue",
			},
		},
	}
	txa, err := sv.AddTabs(DefaultCtx, testEnv, "1", nTabs)
	if err != nil {
		t.Errorf("AddTabs Error: %v", err)
		return
	}
	fmt.Printf("Len Text Tabs: %d\n", len(txa.TextTabs)) */
	/* tx, err := sv.GetTabs(DefaultCtx, testEnv, "2")
	if err != nil {
		t.Errorf("GetTabs Error: %v", err)
		return
	}
	fmt.Printf("%#v\n", tx)
	*/
	/*

	*/
	/*
		fel, err := sv.GetFolderEnvList(DefaultCtx, baseTestFolder)
		if err != nil {
			t.Errorf("Err: %v", err)
			return
		}
		if len(fel.FolderItems) <= 0 {
			t.Errorf("GetFolderEnvList failed to return folders")
		}
		for _, xE := range fel.FolderItems {
			fmt.Printf("Nm: %s %s %s %s\n", xE.Name, xE.EnvelopeId, xE.Status, xE.Subject)
		}

		aevt, err := sv.GetEnvelopeAuditEvents(DefaultCtx, baseTestEnv)
		if err != nil {
			t.Errorf("GetEnvelopeAuditEvents failed: %v", err)

		}
		if len(aevt.AuditEvents) <= 0 {
			t.Errorf("GetEnvelopeAuditEvents failed: No records returned")
		}

		nt, err := sv.GetEnvelopeNotification(DefaultCtx, baseTestEnv)
		if err != nil {
			t.Errorf("GetEnvelopeNotification failed: %v", err)
			return
		}
		fmt.Printf("%v\n", nt)
		if len(nt.Reminders.ReminderEnabled) <= 0 {
			t.Errorf("GetEnvelopeNotification zero: ")
		}

		cf, err := sv.GetEnvelopeCustomFields(DefaultCtx, baseTestEnv)
		if err != nil {
			t.Errorf("GetEnvelopeCustomFields failed: %v", err)
			return
		}

		if len(cf.TextCustomFields) <= 0 {
			t.Errorf("GetEnvelopeCustomFields failed: ")
		}
		docList, err := sv.GetEnvelopeDocumentList(DefaultCtx, baseTestEnv)
		if err != nil {
			t.Errorf("GetEnvelopeDocumentList failed: %v", err)
			return
		}

		if len(docList.EnvelopeDocuments) <= 0 {
			t.Errorf("GetEnvelopeDocumentList failed: ")
		}

	*/
	/*
		fpdf, err := os.OpenFile("test.pdf", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
		if err != nil {
			t.Errorf("File Open: %v", err)
			return
		}

		_, err = fpdf.Write([]byte("ABCAFDAFD\n"))
		if err != nil {
			t.Errorf("%v", err)
		}
		defer fpdf.Close()

		err = sv.GetEnvelopeDocument(fpdf, baseTestEnv, "1", NmVal{Name: "show_changes", Value: "true"})
		if err != nil {
			t.Errorf("GetEnvelopeDocument: %v", err)
			return
		}
		fmt.Printf("%v\n", docList)
	*/
	/*
		err = sv.GetEnvelopeDocumentCombined(baseTestEnv, fpdf, NmVal{Name: "show_changes", Value: "true"}, NmVal{Name: "certificate", Value: "true"})
		if err != nil {
			t.Errorf("GetEnvelopeDocument: %v", err)
			return
		}
	*/
	/*
			fpdf, err := os.OpenFile("HotelExpenseRpt.doc", os.O_RDONLY, os.ModePerm) //os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
			if err != nil {
				t.Errorf("File Open: %v", err)
				return
			}
			dl := DocumentList{}
			dl.Documents = make([]Document, 1, 1)
			dl.Documents[0].DocumentId = "4"
			dl.Documents[0].Order = "4"
			dl.Documents[0].Name = "Uploadedx.doc"
			fu := UploadFile{}
			fu.ContentType = "application/msword"
			fu.Data = fpdf
			fu.FileName = "Uploadedx.doc"
			fu.Id = "4"
			_, err = sv.SetEnvelopeDocuments(baseTestEnv, dl, fu)


		dl := DocumentList{}
		dl.Documents = make([]Document, 1, 1)
		dl.Documents[0].DocumentId = "3"
		_, err = sv.RemoveEnvelopeDocuments(baseTestEnv, &dl)
		if err != nil {
			t.Errorf("ERR %v", err)

		}
	*/
	/*
		cl := &CustomFieldList{}
		cl.TextCustomFields = []CustomField{CustomField{Name: "CustFld1", Value: "XWhatever"}, CustomField{Name: "CustFld2", Value: "XWhateverYea"}}
		lx := ListCustomField{} //Name: "XName", Value: "X"}
		lx.Name = "ListFld1"
		lx.Value = "X"
		lx.ListItems = []string{"A", "B", "C", "D"}
		cl.ListCustomFields = []ListCustomField{lx}

		cl, err = sv.AddEnvelopeCustomFields(baseTestEnv, cl)
		if err != nil {
			t.Errorf("Err: %v", err)
		}

		cl, err = sv.GetEnvelopeCustomFields(baseTestEnv)
		if err != nil {
			t.Errorf("Err: %v", err)
		}
		for i := range cl.TextCustomFields {
			cl.TextCustomFields[i].Value = ""
			cl.TextCustomFields[i].Name = ""
			cl.TextCustomFields[i].Required = ""
			cl.TextCustomFields[i].Show = ""
		}
		for i := range cl.ListCustomFields {
			cl.ListCustomFields[i].Value = ""
			cl.ListCustomFields[i].Name = ""
			cl.ListCustomFields[i].Required = ""
			cl.ListCustomFields[i].Show = ""
			cl.ListCustomFields[i].ListItems = []string{}
		}

		_, err = sv.RemoveEnvelopeCustomFields(baseTestEnv, cl)
		if err != nil {
			t.Errorf("Err: %v", err)
		}

		return
		//f, err := Get
	*/
}

func testEnvelopePayload(userName string) *Envelope {
	return &Envelope{
		Status: "created",
		CustomFields: &CustomFieldList{
			TextCustomFields: []CustomField{
				CustomField{Name: "PID", Value: "123456"},
				CustomField{Name: "Project", Value: "P1"},
			},
		},
		Documents: []Document{
			Document{
				DocumentFields: []NmVal{
					NmVal{Name: "Pid", Value: "122312"},
					NmVal{Name: "DocType", Value: "TestDoc"},
				},
				DocumentId: "1",
				Name:       "TestDoc.pdf",
				Order:      "1",
			},
		},
		EmailSubject: "Created by Go Test",
		EmailBlurb:   "Dear Person: Please read <strong>this</strong>.",
		Recipients: &RecipientList{
			Signers: []Signer{
				Signer{
					EmailRecipient: EmailRecipient{
						Email: userName,
						Recipient: Recipient{
							Name:              "My Name",
							Note:              "This is the ,Note for My Name",
							EmailNotification: &EmailNotification{EmailBody: "This is the recipient 1 email blurb", EmailSubject: "This is the Subject for recipient 1"},
							RecipientId:       "1",
							RoleName:          "Role1",
							RoutingOrder:      "1",
						},
					},
					BaseSigner: BaseSigner{
						Tabs: &Tabs{
							TextTabs: []TextTab{
								TextTab{
									BaseTab: BaseTab{
										DocumentID: "1",
										TabLabel:   "txtTextFieldA",
									},
									BasePosTab: BasePosTab{
										AnchorString:  "TextFieldA:",
										AnchorXOffset: "40",
										AnchorYOffset: "-7",
										AnchorUnits:   "pixels",
										PageNumber:    "1",
									},
									BaseTemplateTab: BaseTemplateTab{
										RecipientID: "1",
									},
									Value: "Value 1",
								},
							},
							SignHereTabs: []SignHereTab{
								SignHereTab{
									BaseTab: BaseTab{
										DocumentID: "1",
										TabLabel:   "signHereA",
									},
									BasePosTab: BasePosTab{
										AnchorString:  "SignatureA:",
										AnchorXOffset: "40",
										AnchorYOffset: "-7",
										AnchorUnits:   "pixels",
										PageNumber:    "1",
									},
									BaseTemplateTab: BaseTemplateTab{
										RecipientID: "1",
									},
								},
							},
							DateSignedTabs: []DateSignedTab{
								DateSignedTab{
									BaseTab: BaseTab{
										DocumentID: "1",
										TabLabel:   "dtSignedA",
									},
									BasePosTab: BasePosTab{
										AnchorString: "DateSignedA:",
										PageNumber:   "1",
									},
									BaseTemplateTab: BaseTemplateTab{
										RecipientID: "1",
									},
								},
							},
						},
					},
				},
				Signer{
					EmailRecipient: EmailRecipient{
						Email: "abc@example.com",
						Recipient: Recipient{
							Name:              "XXX YYYY",
							Note:              "Note for Recipient 2",
							EmailNotification: &EmailNotification{EmailBody: "This is the recipient 2 email blurb", EmailSubject: "This is the Subject for recipient 2"},
							RecipientId:       "2",
							RoleName:          "Role2",
							RoutingOrder:      "2",
						},
					},
					BaseSigner: BaseSigner{
						Tabs: &Tabs{
							TextTabs: []TextTab{
								TextTab{
									BaseTab: BaseTab{
										DocumentID: "1",
										TabLabel:   "txtTextFieldB",
									},
									BasePosTab: BasePosTab{
										AnchorString:  "TextFieldB:",
										AnchorXOffset: "40",
										AnchorYOffset: "-7",
										AnchorUnits:   "pixels",
										PageNumber:    "1",
									},
									BaseTemplateTab: BaseTemplateTab{
										RecipientID: "2",
									},
									Value: "Value 2",
								},
							},
							SignHereTabs: []SignHereTab{
								SignHereTab{
									BaseTab: BaseTab{
										DocumentID: "1",
										TabLabel:   "signHereA",
									},
									BasePosTab: BasePosTab{
										AnchorString:  "SignatureB:",
										AnchorXOffset: "40",
										AnchorYOffset: "-7",
										AnchorUnits:   "pixels",
										PageNumber:    "1",
									},
									BaseTemplateTab: BaseTemplateTab{
										RecipientID: "2",
									},
								},
							},
							DateSignedTabs: []DateSignedTab{
								DateSignedTab{
									BaseTab: BaseTab{
										DocumentID: "1",
										TabLabel:   "dtSignedB",
									},
									BasePosTab: BasePosTab{
										AnchorString:  "DateSignedB:",
										AnchorXOffset: "40",
										AnchorYOffset: "-7",
										AnchorUnits:   "pixels",
										PageNumber:    "1",
									},
									BaseTemplateTab: BaseTemplateTab{
										RecipientID: "2",
									},
								},
							},
							CheckboxTabs: []CheckboxTab{
								CheckboxTab{
									BaseTab: BaseTab{
										DocumentID: "1",
										TabLabel:   "cbTest",
									},
									BasePosTab: BasePosTab{
										AnchorString:  "Checkbox:",
										AnchorXOffset: "40",
										AnchorYOffset: "-7",
										AnchorUnits:   "pixels",
										PageNumber:    "1",
									},
									BaseTemplateTab: BaseTemplateTab{
										RecipientID: "2",
									},
									Selected: true,
								},
							},
							RadioGroupTabs: []RadioGroupTab{
								RadioGroupTab{
									GroupName:   "rbGrp",
									RecipientID: "2",
									DocumentID:  "1",
									Radios: []Radio{
										Radio{
											BasePosTab: BasePosTab{
												AnchorString:  "rbA",
												AnchorXOffset: "28",
												AnchorYOffset: "-7",
												AnchorUnits:   "pixels",
												PageNumber:    "1",
											},
											Selected: false,
											Value:    "val1",
										},
										Radio{
											BasePosTab: BasePosTab{
												AnchorString:  "rbB",
												AnchorXOffset: "28",
												AnchorYOffset: "-7",
												AnchorUnits:   "pixels",
												PageNumber:    "1",
											},
											Selected: true,
											Value:    "val2",
										},
										Radio{
											BasePosTab: BasePosTab{
												AnchorString:  "rbC",
												AnchorXOffset: "28",
												AnchorYOffset: "-7",
												AnchorUnits:   "pixels",
												PageNumber:    "1",
											},
											Selected: false,
											Value:    "val3",
										},
									},
								},
							},
							ListTabs: []ListTab{
								ListTab{
									BaseTab: BaseTab{
										DocumentID: "1",
										TabLabel:   "dlDrop",
									},
									BasePosTab: BasePosTab{
										AnchorString:  "DropdownList:",
										AnchorXOffset: "40",
										AnchorYOffset: "-7",
										AnchorUnits:   "pixels",
										PageNumber:    "1",
									},
									BaseTemplateTab: BaseTemplateTab{
										RecipientID: "2",
									},
									//Value: "X",
									ListItems: []ListItem{
										ListItem{
											Selected: true,
											Text:     "X Val",
											Value:    "X",
										},
										ListItem{
											Selected: false,
											Text:     "Y Val",
											Value:    "Y",
										},
										ListItem{
											Selected: false,
											Text:     "Z Val",
											Value:    "Z",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

}

func TestXML(t *testing.T) {
	_ = bytes.NewBufferString("")
	f, err := os.Open("testdata/connect.xml")
	if err != nil {
		t.Fatalf("Open Connect.xml: %v", err)
		return
	}
	var v *DocuSignEnvelopeInformation = &DocuSignEnvelopeInformation{}
	decoder := xml.NewDecoder(f)
	err = decoder.Decode(v)
	if err != nil {
		t.Fatalf("XML Decode: %v", err)
		return
	}
	if v.EnvelopeStatus.DocumentStatuses[0].Name != "Docusign1.pdf" {
		t.Errorf("Invalid Document Name in Connect XML: %s", v.EnvelopeStatus.DocumentStatuses[0].Name)
	}
	return
}

func TestMultiBody(t *testing.T) {
	var payload struct {
		A string `json:"a,omitempty"`
		B int    `json:"b,omitempty"`
	}
	payload.A = "A"
	payload.B = 999
	files := []*UploadFile{
		&UploadFile{Data: newReadCloser("XXXX"), ContentType: "text/plain", FileName: "fn1", Id: "1"},
		&UploadFile{Data: newReadCloser("XXXX"), ContentType: "text/plain", FileName: "fn2", Id: "2"},
		&UploadFile{Data: newReadCloser("XXXX"), ContentType: "text/plain", FileName: "fn3", Id: "3"},
	}
	r, ct := multiBody(payload, files)

	defer r.(io.ReadCloser).Close()

	mpr := multipart.NewReader(r, ct[30:])

	pt, err := mpr.NextPart()
	if err != nil {
		t.Errorf("Unable to parse part from multireader: %v", err)
		return
	}

	payload.A = ""
	payload.B = 0
	if err := json.NewDecoder(pt).Decode(&payload); err != nil {
		t.Errorf("JSON Unmarshal: %v", err)
		return
	} else {
		if payload.A != "A" || payload.B != 999 {
			t.Errorf("Expect A=A and B=999; got %s %d", payload.A, payload.B)
			return
		}
	}

	for cnt := 0; cnt < len(files); cnt++ {
		if pt, err = mpr.NextPart(); err != nil {
			t.Errorf("Unable to parse multipart reader: %v", err)
			return
		}
		if pt.Header.Get("content-disposition") != fmt.Sprintf("file; filename=\"%s\";documentid=%s", files[cnt].FileName, files[cnt].Id) {
			t.Errorf("Invalid content-dispostion: %s", pt.Header.Get("content-dispostion"))
		}
		bx := make([]byte, 4)
		if _, err = pt.Read(bx); err != nil {
			t.Errorf("Expected EOF: got %v", err)
		} else if string(bx) != "XXXX" {
			t.Errorf("expectd XXXX; got %s", string(bx))
		}
	}
}

func newReadCloser(s string) io.ReadCloser {
	return byteReadCloser{Buffer: bytes.NewBufferString(s)}
}

type byteReadCloser struct {
	*bytes.Buffer
}

func (b byteReadCloser) Close() error {
	return nil
}
