package api

import (
	"fmt"
	"strconv"
	"time"
)

// Records represents SharePoint Item Records throught REST+CSOM API object struct
type Records struct {
	item *Item
}

// NewRecords ...
func NewRecords(item *Item) *Records {
	return &Records{item: item}
}

// IsRecord ...
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
	if fmt.Sprintf("%s", date) == "0001-01-01 00:00:00 +0000 UTC" {
		return false, nil
	}
	return true, nil
}

// RecordDate ...
func (records *Records) RecordDate() (time.Time, error) {
	data, err := records.item.Select("OData__vti_ItemDeclaredRecord").Get()
	if err != nil {
		return time.Time{}, err
	}
	res := &struct {
		RecordDate time.Time `json:"OData__vti_ItemDeclaredRecord"`
	}{}
	if err := data.Unmarshal(&res); err != nil {
		return time.Time{}, err
	}
	return res.RecordDate, nil
}

// Declare ...
func (records *Records) Declare() error {
	_, err := csomItemRecordMethod(records.item, "DeclareItemAsRecord", nil)
	return err
}

// DeclareWithDate ...
func (records *Records) DeclareWithDate(date time.Time) error {
	_, err := csomItemRecordMethod(records.item, "DeclareItemAsRecordWithDeclarationDate", &date)
	return err
}

// Undeclare ...
func (records *Records) Undeclare() error {
	_, err := csomItemRecordMethod(records.item, "UndeclareItemAsRecord", nil)
	return err
}

// func csomItemRecordMethod(item *Item, csomStaticMethod string, date *time.Time) ([]byte, error) {
// 	sp := NewHTTPClient(item.client)
// 	site := NewSP(item.client).Site().Conf(item.config)
// 	list := item.ParentList()
// 	web := item.ParentList().ParentWeb()

// 	var siteR SiteResp // Find a way to reduce requests number
// 	var webR WebResp
// 	var listR ListResp
// 	var itemR ItemResp
// 	errs := []error{}

// 	var wg sync.WaitGroup

// 	wg.Add(1)
// 	go func() {
// 		siteRR, err := site.Select("Id").Get()
// 		if err != nil {
// 			errs = append(errs, err)
// 		}
// 		siteR = siteRR
// 		wg.Done()
// 	}()

// 	wg.Add(1)
// 	go func() {
// 		webRR, err := web.Select("Id").Get()
// 		if err != nil {
// 			errs = append(errs, err)
// 		}
// 		webR = webRR
// 		wg.Done()
// 	}()

// 	wg.Add(1)
// 	go func() {
// 		listRR, err := list.Select("Id").Get()
// 		if err != nil {
// 			errs = append(errs, err)
// 		}
// 		listR = listRR
// 		wg.Done()
// 	}()

// 	wg.Add(1)
// 	go func() {
// 		itemRR, err := item.Select("Id").Get()
// 		if err != nil {
// 			errs = append(errs, err)
// 		}
// 		itemR = itemRR
// 		wg.Done()
// 	}()

// 	wg.Wait()

// 	if len(errs) > 0 {
// 		err := fmt.Errorf("")
// 		for _, e := range errs {
// 			if len(err.Error()) > 0 {
// 				err = fmt.Errorf("%s; ", err)
// 			}
// 			err = fmt.Errorf("%s %s", err, e)
// 		}
// 		return nil, err
// 	}

// 	timeParameter := ""
// 	if date != nil && csomStaticMethod == "DeclareItemAsRecordWithDeclarationDate" {
// 		timeParameter = fmt.Sprintf(`<Parameter Type="DateTime">%s</Parameter>`, date.Format(time.RFC3339))
// 	}
// 	body := []byte(trimMultiline(`
// 		<Request xmlns="http://schemas.microsoft.com/sharepoint/clientquery/2009" SchemaVersion="15.0.0.0" LibraryVersion="16.0.0.0" ApplicationName="Gosip">
// 			<Actions>
// 				<StaticMethod TypeId="{ea8e1356-5910-4e69-bc05-d0c30ed657fc}" Name="` + csomStaticMethod + `" Id="6">
// 					<Parameters>
// 						<Parameter ObjectPathId="5" />
// 						` + timeParameter + `
// 					</Parameters>
// 				</StaticMethod>
// 			</Actions>
// 			<ObjectPaths>
// 				<Identity Id="2" Name="740c6a0b-85e2-48a0-a494-e0f1759d4aa7:site:` + siteR.Data().ID + `:web:` + webR.Data().ID + `" />
// 				<Property Id="3" ParentId="2" Name="Lists" />
// 				<Method Id="4" ParentId="3" Name="GetById">
// 					<Parameters>
// 						<Parameter Type="String">` + listR.Data().ID + `</Parameter>
// 					</Parameters>
// 				</Method>
// 				<Method Id="5" ParentId="4" Name="GetItemById">
// 					<Parameters>
// 						<Parameter Type="Number">` + strconv.Itoa(itemR.Data().ID) + `</Parameter>
// 					</Parameters>
// 				</Method>
// 			</ObjectPaths>
// 		</Request>
// 	`))
// 	jsomResp, err := sp.ProcessQuery(item.client.AuthCnfg.GetSiteURL(), body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return jsomResp, nil
// }

func csomItemRecordMethod(item *Item, csomStaticMethod string, date *time.Time) ([]byte, error) {
	sp := NewHTTPClient(item.client)
	site := NewSP(item.client).Site().Conf(item.config)
	siteR, err := site.Select("Id").Get()
	if err != nil {
		return nil, err
	}
	itemR, err := item.Select("Id").Get()
	if err != nil {
		return nil, err
	}
	list := item.ParentList()
	listR, err := list.Select("Id").Get()
	if err != nil {
		return nil, err
	}
	web := item.ParentList().ParentWeb()
	webR, err := web.Select("Id").Get()
	if err != nil {
		return nil, err
	}
	timeParameter := ""
	if date != nil && csomStaticMethod == "DeclareItemAsRecordWithDeclarationDate" {
		timeParameter = fmt.Sprintf(`<Parameter Type="DateTime">%s</Parameter>`, date.Format(time.RFC3339))
	}
	body := []byte(trimMultiline(`
		<Request xmlns="http://schemas.microsoft.com/sharepoint/clientquery/2009" SchemaVersion="15.0.0.0" LibraryVersion="16.0.0.0" ApplicationName="Gosip">
			<Actions>
				<StaticMethod TypeId="{ea8e1356-5910-4e69-bc05-d0c30ed657fc}" Name="` + csomStaticMethod + `" Id="6">
					<Parameters>
						<Parameter ObjectPathId="5" />
						` + timeParameter + `
					</Parameters>
				</StaticMethod>
			</Actions>
			<ObjectPaths>
				<Identity Id="2" Name="740c6a0b-85e2-48a0-a494-e0f1759d4aa7:site:` + siteR.Data().ID + `:web:` + webR.Data().ID + `" />
				<Property Id="3" ParentId="2" Name="Lists" />
				<Method Id="4" ParentId="3" Name="GetById">
					<Parameters>
						<Parameter Type="String">` + listR.Data().ID + `</Parameter>
					</Parameters>
				</Method>
				<Method Id="5" ParentId="4" Name="GetItemById">
					<Parameters>
						<Parameter Type="Number">` + strconv.Itoa(itemR.Data().ID) + `</Parameter>
					</Parameters>
				</Method>
			</ObjectPaths>
		</Request>
	`))
	jsomResp, err := sp.ProcessQuery(item.client.AuthCnfg.GetSiteURL(), body)
	if err != nil {
		return nil, err
	}
	return jsomResp, nil
}
