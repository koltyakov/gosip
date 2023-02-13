package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/koltyakov/gosip/csom"
)

// Records represents SharePoint Item Records via REST+CSOM API object struct
// Always use NewRecords constructor instead of &Records{}
type Records struct {
	item *Item
}

// NewRecords - Records struct constructor function
func NewRecords(item *Item) *Records {
	return &Records{item: item}
}

// IsRecord checks is current item is declared as a record
func (records *Records) IsRecord() (bool, error) {
	// // It is better using REST and OData__vti_ItemDeclaredRecord field value
	// jsomResp, err := csomItemRecordMethod(records.item, "IsRecord", nil)
	// if err != nil {
	// 	return false, err
	// }

	// arrRes := []interface{}{}
	// if err := json.Unmarshal(jsomResp, &arrRes); err != nil {
	// 	return false, err
	// }
	// if len(arrRes) < 3 {
	// 	return false, fmt.Errorf("can't parse CSOM response")
	// }

	// return arrRes[2].(bool), nil
	date, err := records.RecordDate()
	if err != nil {
		return false, err
	}
	if date.String() == "0001-01-01 00:00:00 +0000 UTC" {
		return false, nil
	}
	return true, nil
}

// RecordDate checks record declaration date of this item
func (records *Records) RecordDate() (time.Time, error) {
	data, err := records.item.Select("OData__vti_ItemDeclaredRecord").Get()
	if err != nil {
		if strings.Contains(err.Error(), "OData__vti_ItemDeclaredRecord") {
			return time.Time{}, nil // in place records is not configured in a list
		}
		return time.Time{}, err
	}
	res := &struct {
		RecordDate time.Time `json:"OData__vti_ItemDeclaredRecord"`
	}{}
	data = NormalizeODataItem(data)
	if err := json.Unmarshal(data, &res); err != nil {
		return time.Time{}, err
	}
	return res.RecordDate, nil
}

// Declare declares this item as a record (CSOM helper)
func (records *Records) Declare() error {
	_, err := csomItemRecordMethod(records.item, "DeclareItemAsRecord", nil, records.item.config)
	return err
}

// DeclareWithDate declares this item as a record with record declaration date (CSOM helper)
func (records *Records) DeclareWithDate(date time.Time) error {
	_, err := csomItemRecordMethod(records.item, "DeclareItemAsRecordWithDeclarationDate", &date, records.item.config)
	return err
}

// Undeclare undeclared this item as a record (the item is not a record after an action is done) (CSOM helper)
func (records *Records) Undeclare() error {
	_, err := csomItemRecordMethod(records.item, "UndeclareItemAsRecord", nil, records.item.config)
	return err
}

// csomItemRecordMethod constructs CSOM API process query to cover missed REST API functionality
func csomItemRecordMethod(item *Item, csomStaticMethod string, date *time.Time, config *RequestConfig) ([]byte, error) {
	client := NewHTTPClient(item.client)
	itemR, err := item.Select("Id").Get()
	if err != nil {
		return nil, err
	}
	list := item.ParentList()
	listR, err := list.Select("Id").Get()
	if err != nil {
		return nil, err
	}
	timeParameter := ""
	if date != nil && csomStaticMethod == "DeclareItemAsRecordWithDeclarationDate" {
		timeParameter = fmt.Sprintf(`<Parameter Type="DateTime">%s</Parameter>`, date.Format(time.RFC3339))
	}

	b := csom.NewBuilder()

	b.AddObject(csom.NewObjectProperty("Web"), nil)
	b.AddObject(csom.NewObjectProperty("Lists"), nil)
	b.AddObject(csom.NewObjectMethod("GetById", []string{`<Parameter Type="String">` + listR.Data().ID + `</Parameter>`}), nil)
	b.AddObject(csom.NewObjectMethod("GetItemById", []string{`<Parameter Type="Number">` + strconv.Itoa(itemR.Data().ID) + `</Parameter>`}), nil)
	b.AddAction(csom.NewAction(`
		<StaticMethod TypeId="{ea8e1356-5910-4e69-bc05-d0c30ed657fc}" Name="`+csomStaticMethod+`" Id="{{.ID}}">
			<Parameters>
				<Parameter ObjectPathId="{{.ObjectID}}" />
				`+timeParameter+`
			</Parameters>
		</StaticMethod>
	`), nil)

	csomPkg, err := b.Compile()
	if err != nil {
		return nil, err
	}

	jsomResp, err := client.ProcessQuery(item.client.AuthCnfg.GetSiteURL(), bytes.NewBuffer([]byte(csomPkg)), config)
	if err != nil {
		return nil, err
	}
	return jsomResp, nil
}
