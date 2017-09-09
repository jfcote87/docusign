package docusign_test

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Wealthforge-Technologies/docusign"
	"golang.org/x/net/context"
)

func ExampleConfig() {
	ctx := context.Background()
	config := &docusign.Config{
		UserName:      "YOUR_USER_NAME",
		Password:      "YOUR_PASSWORD",
		IntegratorKey: "YOUR_INTEGRATOR_KEY",
		AccountId:     "YOUR_ACCOUNT_ID",
		Host:          "YOUR_HOST", // NOTE: set to 'demo.docusign.net' for non-prod tests
	}

	// create service using config as credential
	sv := docusign.New(config, "")

	folderList, err := sv.FolderList(ctx, docusign.FolderTemplatesInclude)
	if err != nil {
		log.Fatal(err)
	}
	for _, fld := range folderList.Folders {
		fmt.Printf("%s: %s\n", fld.FolderId, fld.Name)
	}

	// obtain a new oauth credential
	token, err := config.OauthCredential(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Token: %s\n", token.AccessToken)
}

func ExampleOauthCredential() {
	ctx := context.Background()
	cred := &docusign.OauthCredential{
		AccessToken: "SAVED_ACCESS_TOKEN",
		AccountId:   "YOUR_ACCOUNT_ID",
		TokenType:   "bearer",
	}

	sv := docusign.New(cred, "")

	results, err := sv.EnvelopeSearch(ctx, docusign.SearchFolderCompleted,
		docusign.EnvelopeSearchFromDate(time.Now().Add(-time.Hour*72)),
		docusign.EnvelopeSearchIncludeRecipients)
	if err != nil {
		log.Fatal(err)
	}
	for _, env := range results.FolderItems {
		fmt.Printf("%s: %s\n", env.Name, env.EnvelopeId)
	}
}

func ExampleOnBehalfOf(ctx context.Context, sv *docusign.Service, userEmail string) (string, error) {
	info, err := sv.OnBehalfOf(userEmail).LoginInformation(ctx,
		docusign.LoginInformationSettingsAll,
		docusign.LoginInformationIncludeApiPassword)
	if err != nil {
		return "", err
	}
	return info.LoginAccounts[0].Name, nil

}

func ExampleEnvelopeCreate(ctx context.Context, sv *docusign.Service, userID string) {
	f, err := os.Open("FILE_NAME")
	if err != nil {
		log.Fatal(err)
	}
	uploadDoc := docusign.UploadFile{
		ContentType: "application/pdf",
		FileName:    "contract.pdf",
		Id:          "1",
		Data:        f,
	}

	env := &docusign.Envelope{
		Status: "created",
		CustomFields: &docusign.CustomFieldList{
			TextCustomFields: []docusign.CustomField{
				docusign.CustomField{Name: "Project", Value: "P1"},
			},
		},
		Documents: []docusign.Document{
			docusign.Document{
				DocumentId: "1",
				Name:       "TestDoc.pdf",
				Order:      "1",
			},
		},
		EmailSubject: "Created by Go Test",
		EmailBlurb:   "Dear Person: Please read this.",
		Recipients: &docusign.RecipientList{
			Signers: []docusign.Signer{
				docusign.Signer{
					EmailRecipient: docusign.EmailRecipient{
						Email: "USER_EMAIL",
						Recipient: docusign.Recipient{
							Name:         "My Name",
							Note:         "This is the ,Note for My Name",
							RecipientId:  "1",
							RoleName:     "Role1",
							RoutingOrder: "1",
						},
					},
					BaseSigner: docusign.BaseSigner{
						Tabs: &docusign.Tabs{
							TextTabs: []docusign.TextTab{
								docusign.TextTab{
									BaseTab: docusign.BaseTab{
										DocumentID: "1",
										TabLabel:   "txtTextFieldA",
									},
									BasePosTab: docusign.BasePosTab{
										AnchorString:  "TextFieldA:",
										AnchorXOffset: "40",
										AnchorYOffset: "-7",
										AnchorUnits:   "pixels",
										PageNumber:    "1",
									},
									BaseTemplateTab: docusign.BaseTemplateTab{
										RecipientID: "1",
									},
									Value: "Value 1",
								},
							},
							SignHereTabs: []docusign.SignHereTab{
								docusign.SignHereTab{
									BaseTab: docusign.BaseTab{
										DocumentID: "1",
										TabLabel:   "signHereA",
									},
									BasePosTab: docusign.BasePosTab{
										AnchorString:  "SignatureA:",
										AnchorXOffset: "40",
										AnchorYOffset: "-7",
										AnchorUnits:   "pixels",
										PageNumber:    "1",
									},
									BaseTemplateTab: docusign.BaseTemplateTab{
										RecipientID: "1",
									},
								},
							},
							DateSignedTabs: []docusign.DateSignedTab{
								docusign.DateSignedTab{
									BaseTab: docusign.BaseTab{
										DocumentID: "1",
										TabLabel:   "dtSignedA",
									},
									BasePosTab: docusign.BasePosTab{
										AnchorString: "DateSignedA:",
										PageNumber:   "1",
									},
									BaseTemplateTab: docusign.BaseTemplateTab{
										RecipientID: "1",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	newEnvelope, err := sv.OnBehalfOf(userID).EnvelopeCreate(ctx, env, &uploadDoc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("EnvelopeID: %s\n", newEnvelope.EnvelopeId)

}
