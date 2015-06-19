package docusign

import "strings"

// Tabs describes the data tabs for a recipient
type Tabs struct {
	ApproveTabs          []ApproveTab          `json:"ApproveTabs,omitempty"`
	CheckboxTabs         []CheckboxTab         `json:"checkboxTabs,omitempty"`
	CompanyTabs          []CompanyTab          `json:"companyTabs,omitempty"`
	DateSignedTabs       []DateSignedTab       `json:"dateSignedTabs,omitempty"`
	DateTabs             []DateTab             `json:"dateTabs,omitempty"`
	DeclineTabs          []DeclineTab          `json:"declineTabs,omitempty"`
	EmailTabs            []EmailTab            `json:"emailTabs,omitempty"`
	EnvelopeIdTabs       []EnvelopeIdTab       `json:"envelopeIdTabs,omitempty"`
	FullNameTabs         []FullNameTab         `json:"fullNameTabs,omitempty"`
	InitialHereTabs      []InitialHereTab      `json:"initialHereTabs,omitempty"`
	ListTabs             []ListTab             `json:"listTabs,omitempty"`
	NoteTabs             []NoteTab             `json:"noteTabs,omitempty"`
	NumberTabs           []NumberTab           `json:"numberTabs,omitempty"`
	RadioGroupTabs       []RadioGroupTab       `json:"radioGroupTabs,omitempty"`
	SignHereTabs         []SignHereTab         `json:"signHereTabs,omitempty"`
	SignerAttachmentTabs []SignerAttachmentTab `json:"signerAttachmentTabs,omitempty"`
	SsnTabs              []SsnTab              `json:"ssnTabs,omitempty"`
	TextTabs             []TextTab             `json:"textTabs,omitempty"`
	TitleTabs            []TitleTab            `json:"titleTabs,omitempty"`
	ZipTabs              []ZipTab              `json:"zipTabs,omitempty"`
}

func (t Tabs) Values() []NmVal {
	vals := make([]NmVal, 0, len(t.CheckboxTabs)+len(t.CompanyTabs)+len(t.DateTabs)+len(t.EmailTabs)+len(t.ListTabs)+len(t.NoteTabs)+len(t.NumberTabs)+len(t.RadioGroupTabs)+len(t.SsnTabs)+len(t.TextTabs)+len(t.TitleTabs)+len(t.ZipTabs))
	for _, v := range t.CheckboxTabs {
		vals = append(vals, v.NmVal())
	}
	for _, v := range t.CompanyTabs {
		vals = append(vals, v.NmVal())
	}
	for _, v := range t.DateTabs {
		vals = append(vals, v.NmVal())
	}
	for _, v := range t.EmailTabs {
		vals = append(vals, v.NmVal())
	}
	for _, v := range t.ListTabs {
		vals = append(vals, v.NmVal())
	}
	for _, v := range t.NoteTabs {
		vals = append(vals, v.NmVal())
	}
	for _, v := range t.NumberTabs {
		vals = append(vals, v.NmVal())
	}
	for _, v := range t.RadioGroupTabs {
		vals = append(vals, v.NmVal())
	}
	for _, v := range t.SsnTabs {
		vals = append(vals, v.NmVal())
	}
	for _, v := range t.TextTabs {
		vals = append(vals, v.NmVal())
	}
	for _, v := range t.TitleTabs {
		vals = append(vals, v.NmVal())
	}
	for _, v := range t.ZipTabs {
		vals = append(vals, v.NmVal())
	}
	return vals
}

type ValueTab interface {
	NmVal() NmVal
}

// BasePosTab contains positioning fields for tabs
type BasePosTab struct {
	AnchorIgnoreIfNotPresent string `json:"anchorIgnoreIfNotPresent,omitempty"`
	AnchorString             string `json:"anchorString,omitempty"`
	AnchorUnits              string `json:"anchorUnits,omitempty"`
	AnchorXOffset            string `json:"anchorXOffset,omitempty"`
	AnchorYOffset            string `json:"anchorYOffset,omitempty"`
	PageNumber               string `json:"pageNumber,omitempty"`
	XPosition                string `json:"xPosition,omitempty"`
	YPosition                string `json:"yPosition,omitempty"`
	TabId                    string `json:"tabId,omitempty"`
}

// BaseTab fields
type BaseTab struct {
	DocumentID   string         `json:"documentID,omitempty"`
	TabLabel     string         `json:"tabLabel,omitempty"`
	ErrorDetails *ResponseError `json:"errorDetails,omitempty"`
}

// template related fields
type BaseTemplateTab struct {
	RecipientID      string `json:"recipientID,omitempty"`
	TemplateLocked   DSBool `json:"templateLocked,omitempty"`
	TemplateRequired DSBool `json:"templaterequired"`
}

// Conditional value fields
type BaseConditionalTab struct {
	ConditionalParentLabel string `json:"conditionalParentLabel,omitempty"`
	ConditionalParentValue string `json:"conditionalParentValue,omitempty"`
}

// Style fields
type BaseStyleTab struct {
	Bold      DSBool `json:"bold,omitempty"`
	Font      string `json:"font,omitempty"`
	FontColor string `json:"fontColor,omitempty"`
	FontSize  string `json:"fontSize,omitempty"`
	Italic    DSBool `json:"italic,omitempty"`
	Name      string `json:"name,omitempty"`
	Underline DSBool `json:"underline,omitempty"`
}

// Approve button tab
type ApproveTab struct {
	BaseTab
	BasePosTab
	BaseStyleTab
	BaseTemplateTab
	BaseConditionalTab
	ButtonText string `json:"buttonText,omitempty"`
	Height     int    `json:"height,omitempty"`
	Width      int    `json:"width,omitempty"`
}

// Checkbox
type CheckboxTab struct {
	BaseTab
	BasePosTab
	BaseStyleTab
	BaseTemplateTab
	BaseConditionalTab
	Locked                          DSBool `json:"locked"`
	MergeFieldXml                   string `json:"mergeFieldXml,omitempty"`
	RequireInitialOnSharedTabChange DSBool `json:"requireInitialOnSharedTabChange,omitempty"`
	Selected                        DSBool `json:"selected,omitempty"`
	Shared                          DSBool `json:"shared,omitempty"`
}

func (c CheckboxTab) NmVal() NmVal {
	retVal := NmVal{Name: c.TabLabel}
	if c.Selected {
		retVal.Value = "X"
	}
	return retVal
}

// Company information tab - just a stylized text tab
type CompanyTab struct {
	BaseTab
	BasePosTab
	BaseStyleTab
	BaseTemplateTab
	BaseConditionalTab
	ConcealValueOnDocument DSBool `json:"concealValueOnDocument,omitempty"`
	DisableAutoSize        DSBool `json:"disableAutoSize,omitempty"`
	Locked                 DSBool `json:"locked"`
	Required               DSBool `json:"required"`
	Value                  string `json:"value,omitempty"`
	Width                  int    `json:"width,omitempty"`
}

func (c CompanyTab) NmVal() NmVal {
	return NmVal{Name: c.TabLabel, Value: c.Value}
}

// Auto filled date tab
type DateSignedTab struct {
	BaseTab
	BasePosTab
	BaseStyleTab
	BaseTemplateTab
	BaseConditionalTab
	Value string `json:"value,omitempty"`
}

// User updateable date value tab
type DateTab struct {
	BaseTab
	BasePosTab
	BaseStyleTab
	BaseTemplateTab
	BaseConditionalTab
	ConcealValueOnDocument          DSBool `json:"concealValueOnDocument,omitempty"`
	DisableAutoSize                 DSBool `json:"disableAutoSize,omitempty"`
	Height                          string `json:"height,omitempty"`
	Locked                          DSBool `json:"locked"`
	MergeFieldXml                   string `json:"mergeFieldXml,omitempty"`
	Required                        DSBool `json:"required"`
	RequireInitialOnSharedTabChange DSBool `json:"requireInitialOnSharedTabChange,omitempty"`
	Shared                          DSBool `json:"shared,omitempty"`
	Value                           string `json:"value,omitempty"`
	Width                           int    `json:"width,omitempty"`
}

func (d DateTab) NmVal() NmVal {
	return NmVal{Name: d.TabLabel, Value: d.Value}
}

// Decline button tab
type DeclineTab struct {
	BaseTab
	BasePosTab
	BaseStyleTab
	BaseTemplateTab
	BaseConditionalTab
	ButtonText string `json:"buttonText,omitempty"`
	Height     int    `json:"height,omitempty"`
	Width      int    `json:"width,omitempty"`
}

// User updateable email address
type EmailAddressTab struct {
	BaseTab
	BasePosTab
	BaseStyleTab
	BaseTemplateTab
	BaseConditionalTab
}

// User email display tab
type EmailTab struct {
	BaseTab
	BasePosTab
	BaseStyleTab
	BaseTemplateTab
	BaseConditionalTab
	ConcealValueOnDocument          DSBool `json:"concealValueOnDocument,omitempty"`
	DisableAutoSize                 DSBool `json:"disableAutoSize,omitempty"`
	Height                          int    `json:"height,omitempty"`
	Locked                          DSBool `json:"locked"`
	MergeFieldXml                   string `json:"mergeFieldXml,omitempty"`
	Required                        DSBool `json:"required"`
	RequireInitialOnSharedTabChange DSBool `json:"requireInitialOnSharedTabChange,omitempty"`
	Shared                          DSBool `json:"shared,omitempty"`
	Value                           string `json:"value,omitempty"`
	Width                           int    `json:"width,omitempty"`
}

func (e EmailTab) NmVal() NmVal {
	return NmVal{Name: e.TabLabel, Value: e.Value}
}

// Id display tab
type EnvelopeIdTab struct {
	BaseTab
	BasePosTab
}

type FirstNameTab struct {
	BaseTab
	BasePosTab
	BaseStyleTab
	BaseTemplateTab
	BaseConditionalTab
}

type FormulaTab struct {
	BaseTab
	BasePosTab
	BaseStyleTab
	BaseTemplateTab
	ConcealValueOnDocument DSBool `json:"concealValueOnDocument,omitempty"`
	DisableAutoSize        DSBool `json:"disableAutoSize,omitempty"`
	Formula                string `json:"formula,omitempty"`
	Height                 int    `json:"height,omitempty"`
	IsPaymentAmount        DSBool `json:"isPaymentAmount,omitempty"`
	Locked                 DSBool `json:"locked"`
	MergeFieldXml          string `json:"mergeFieldXml,omitempty"`
	Required               DSBool `json:"required"`
	RoundDecimalPlaces     string `json:"roundDecimalPlaces,omitempty"`
	Value                  string `json:"value,omitempty"`
	Width                  int    `json:"width,omitempty"`
}

func (f FormulaTab) NmVal() NmVal {
	return NmVal{Name: f.TabLabel, Value: f.Value}
}

type FullNameTab struct {
	BaseTab
	BasePosTab
	BaseStyleTab
	BaseTemplateTab
	BaseConditionalTab
}

type InitialHereTab struct {
	BaseTab
	BasePosTab
	BaseStyleTab
	BaseTemplateTab
	BaseConditionalTab
	Optional   DSBool  `json:"optional,omitempty"`
	ScaleValue float64 `json:"scaleValue,omitempty"`
}

type LastNameTab struct {
	BaseTab
	BasePosTab
	BaseStyleTab
	BaseTemplateTab
	BaseConditionalTab
}

type ListItem struct {
	Selected DSBool `json:"selected,omitempty"`
	Text     string `json:"text,omitempty"`
	Value    string `json:"value,omitempty"`
}

type ListTab struct {
	BaseTab
	BasePosTab
	BaseStyleTab
	BaseTemplateTab
	BaseConditionalTab
	ListItems                       []ListItem `json:"listItems,omitempty"`
	Locked                          DSBool     `json:"locked"`
	MergeFieldXml                   string     `json:"mergeFieldXml,omitempty"`
	Required                        DSBool     `json:"required"`
	RequireInitialOnSharedTabChange DSBool     `json:"requireInitialOnSharedTabChange,omitempty"`
	senderRequired                  DSBool     `json:"senderRequired,omitempty,omitempty"`
	Shared                          DSBool     `json:"shared,omitempty"`
	Value                           string     `json:"value,omitempty"`
	Width                           int        `json:"width,omitempty"`
}

func (l ListTab) NmVal() NmVal {
	retVal := make([]string, 0)
	for _, li := range l.ListItems {
		if li.Selected {
			retVal = append(retVal, li.Text)
		}
	}
	return NmVal{Name: l.TabLabel, Value: strings.Join(retVal, "\r\n")}
}

type NoteTab struct {
	BaseTab
	BasePosTab
	BaseStyleTab
	BaseTemplateTab
	BaseConditionalTab
	Height int    `json:"height,omitempty"`
	Shared DSBool `json:"shared,omitempty"`
	Value  string `json:"value,omitempty"`
	Width  int    `json:"width,omitempty"`
}

func (n NoteTab) NmVal() NmVal {
	return NmVal{Name: n.TabLabel, Value: n.Value}
}

type NumberTab struct {
	BaseTab
	BasePosTab
	BaseStyleTab
	BaseTemplateTab
	BaseConditionalTab
	ConcealValueOnDocument          DSBool `json:"concealValueOnDocument,omitempty"`
	DisableAutoSize                 DSBool `json:"disableAutoSize,omitempty"`
	Height                          int    `json:"height,omitempty"`
	Locked                          DSBool `json:"locked"`
	MergeFieldXml                   string `json:"mergeFieldXml,omitempty"`
	Required                        DSBool `json:"required"`
	RequireInitialOnSharedTabChange DSBool `json:"requireInitialOnSharedTabChange,omitempty"`
	Shared                          DSBool `json:"shared,omitempty"`
	Value                           string `json:"value,omitempty"`
	Width                           int    `json:"width,omitempty"`
}

func (n NumberTab) NmVal() NmVal {
	return NmVal{Name: n.TabLabel, Value: n.Value}
}

type Radio struct {
	BasePosTab
	Locked   DSBool `json:"locked"`
	Required DSBool `json:"required"`
	Selected DSBool `json:"selected,omitempty"`
	Value    string `json:"value,omitempty"`
}

type RadioGroupTab struct {
	ConditionalParentLabel          string  `json:"conditionalParentLabel,omitempty"`
	ConditionalParentValue          string  `json:"conditionalParentValue,omitempty"`
	DocumentID                      string  `json:"documentID,omitempty"`
	GroupName                       string  `json:"groupName,omitempty"`
	Radios                          []Radio `json:"radios,omitempty"`
	RecipientID                     string  `json:"recipientID,omitempty"`
	RequireInitialOnSharedTabChange DSBool  `json:"requireInitialOnSharedTabChange,omitempty"`
	Shared                          DSBool  `json:"shared,omitempty"`
	TemplateLocked                  DSBool  `json:"templateLocked,omitempty"`
	TemplateRequired                DSBool  `json:"templaterequired"`
}

func (rg RadioGroupTab) NmVal() NmVal {
	retVal := NmVal{Name: rg.GroupName}
	for _, r := range rg.Radios {
		if r.Selected {
			retVal.Value = r.Value
			break
		}
	}
	return retVal
}

type SignerAttachmentTab struct {
	BaseTab
	BasePosTab
	BaseTemplateTab
	BaseConditionalTab
	Optional DSBool `json:"optional,omitempty"`
	//Required DSBool `json:"required"`
}

type SignHereTab struct {
	BaseTab
	BasePosTab
	BaseStyleTab
	BaseTemplateTab
	BaseConditionalTab
	Optional   DSBool  `json:"optional,omitempty"`
	ScaleValue float64 `json:"scaleValue,omitempty"`
}

type SsnTab struct {
	BaseTab
	BasePosTab
	BaseStyleTab
	BaseTemplateTab
	BaseConditionalTab
	ConcealValueOnDocument          DSBool `json:"concealValueOnDocument,omitempty"`
	DisableAutoSize                 DSBool `json:"disableAutoSize,omitempty"`
	Height                          int    `json:"height,omitempty"`
	Locked                          DSBool `json:"locked"`
	MergeFieldXml                   string `json:"mergeFieldXml,omitempty"`
	Required                        DSBool `json:"required"`
	RequireInitialOnSharedTabChange DSBool `json:"requireInitialOnSharedTabChange,omitempty"`
	Shared                          DSBool `json:"shared,omitempty"`
	Value                           string `json:"value,omitempty"`
	Width                           int    `json:"width,omitempty"`
}

func (s SsnTab) NmVal() NmVal {
	return NmVal{Name: s.TabLabel, Value: s.Value}
}

type TextTab struct {
	BaseTab
	BasePosTab
	BaseStyleTab
	BaseTemplateTab
	BaseConditionalTab
	ConcealValueOnDocument          DSBool `json:"concealValueOnDocument,omitempty"`
	DisableAutoSize                 DSBool `json:"disableAutoSize,omitempty"`
	Height                          int    `json:"height,omitempty"`
	IsPaymentAmount                 DSBool `json:"isPaymentAmount,omitempty"`
	Locked                          DSBool `json:"locked"`
	MergeFieldXml                   string `json:"mergeFieldXml,omitempty"`
	Required                        DSBool `json:"required"`
	RequireInitialOnSharedTabChange DSBool `json:"requireInitialOnSharedTabChange,omitempty"`
	senderRequired                  DSBool `json:"senderRequired,omitempty"`
	Shared                          DSBool `json:"shared,omitempty"`
	ValidationMessage               string `json:"validationMessage,omitempty"`
	ValidationPattern               string `json:"validationPattern,omitempty"`
	Value                           string `json:"value,omitempty"`
	Width                           int    `json:"width,omitempty"`
}

func (t TextTab) NmVal() NmVal {
	return NmVal{Name: t.TabLabel, Value: t.Value}
}

type TitleTab struct {
	BaseTab
	BasePosTab
	BaseStyleTab
	BaseTemplateTab
	BaseConditionalTab
	ConcealValueOnDocument          DSBool `json:"concealValueOnDocument,omitempty"`
	DisableAutoSize                 DSBool `json:"disableAutoSize,omitempty"`
	Height                          int    `json:"height,omitempty"`
	Locked                          DSBool `json:"locked"`
	MergeFieldXml                   string `json:"mergeFieldXml,omitempty"`
	Required                        DSBool `json:"required"`
	RequireInitialOnSharedTabChange DSBool `json:"requireInitialOnSharedTabChange,omitempty"`
	senderRequired                  DSBool `json:"senderRequired,omitempty"`
	Shared                          DSBool `json:"shared,omitempty"`
	ValidationMessage               string `json:"validationMessage,omitempty"`
	ValidationPattern               string `json:"validationPattern,omitempty"`
	Value                           string `json:"value,omitempty"`
	Width                           int    `json:"width,omitempty"`
}

func (t TitleTab) NmVal() NmVal {
	return NmVal{Name: t.TabLabel, Value: t.Value}
}

type ZipTab struct {
	BaseTab
	BasePosTab
	BaseStyleTab
	BaseTemplateTab
	BaseConditionalTab
	ConcealValueOnDocument          DSBool `json:"concealValueOnDocument,omitempty"`
	DisableAutoSize                 DSBool `json:"disableAutoSize,omitempty"`
	Height                          int    `json:"height,omitempty"`
	Locked                          DSBool `json:"locked"`
	MergeFieldXml                   string `json:"mergeFieldXml,omitempty"`
	Required                        DSBool `json:"required"`
	RequireInitialOnSharedTabChange DSBool `json:"requireInitialOnSharedTabChange,omitempty"`
	Shared                          DSBool `json:"shared,omitempty"`
	Value                           string `json:"value,omitempty"`
	Width                           int    `json:"width,omitempty"`
}

func (z ZipTab) NmVal() NmVal {
	return NmVal{Name: z.TabLabel, Value: z.Value}
}
