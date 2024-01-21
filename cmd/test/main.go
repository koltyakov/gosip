package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/koltyakov/gosip/api"
	m "github.com/koltyakov/gosip/test/manual"
)

func main() {

	strategy := flag.String("strategy", "fba", "Auth strategy code")
	flag.Parse()

	client, err := m.GetTestClient(*strategy)
	if err != nil {
		log.Fatal(err)
	}

	// Define requests hook handlers
	// client.Hooks = &gosip.HookHandlers{
	// 	OnError: func(e *gosip.HookEvent) {
	// 		fmt.Println("\n======= On Error ========")
	// 		fmt.Printf("URL: %s\n", e.Request.URL)
	// 		fmt.Printf("StatusCode: %d\n", e.StatusCode)
	// 		fmt.Printf("Error: %s\n", e.Error)
	// 		fmt.Printf("took %f seconds\n", time.Since(e.StartedAt).Seconds())
	// 		fmt.Printf("=========================\n\n")
	// 	},
	// 	OnRetry: func(e *gosip.HookEvent) {
	// 		fmt.Println("\n======= On Retry ========")
	// 		fmt.Printf("URL: %s\n", e.Request.URL)
	// 		fmt.Printf("StatusCode: %d\n", e.StatusCode)
	// 		fmt.Printf("Error: %s\n", e.Error)
	// 		fmt.Printf("took %f seconds\n", time.Since(e.StartedAt).Seconds())
	// 		fmt.Printf("=========================\n\n")
	// 	},
	// 	OnRequest: func(e *gosip.HookEvent) {
	// 		if e.Error == nil {
	// 			fmt.Println("\n====== On Request =======")
	// 			fmt.Printf("URL: %s\n", e.Request.URL)
	// 			fmt.Printf("auth injection took %f seconds\n", time.Since(e.StartedAt).Seconds())
	// 			fmt.Printf("=========================\n\n")
	// 		}
	// 	},
	// 	OnResponse: func(e *gosip.HookEvent) {
	// 		if e.Error == nil {
	// 			fmt.Println("\n====== On Response =======")
	// 			fmt.Printf("URL: %s\n", e.Request.URL)
	// 			fmt.Printf("StatusCode: %d\n", e.StatusCode)
	// 			fmt.Printf("took %f seconds\n", time.Since(e.StartedAt).Seconds())
	// 			fmt.Printf("==========================\n\n")
	// 		}
	// 	},
	// }

	// Manual test code is below

	sp := api.NewSP(client)
	res, err := sp.Web().Select("Title").Get(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", res.Data().Title)

	l := sp.Web().Lists().GetByTitle("Calendar01")
	ii, err := l.Items().Select("Id,Title,Created,Editor/Title").Expand("Editor").Get(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range ii.ToMap() {
		fmt.Printf("%v#\n", i)
	}

	// l := sp.Web().Lists().GetByTitle("Calendar01")

	// vd, err := l.Views().DefaultView().Get()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// _ = l.Views().GetByTitle("Test_View_03").Delete()

	// viewXML := vd.Data().ListViewXML
	// cleanOut := []string{
	// 	`Name="{1232726A-7909-4791-82CF-375D7FC6EF1B}"`,
	// 	`MobileUrl="_layouts/15/mobile/viewdaily.aspx"`,
	// 	`Url="/sites/ci/Lists/Calendar01/calendar.aspx"`,
	// 	`DefaultView="TRUE"`,
	// 	`DisplayName="Calendar"`,
	// }
	// for _, cc := range cleanOut {
	// 	viewXML = strings.Replace(viewXML, cc, ``, -1)
	// }

	// viewXML = strings.Replace(viewXML, "<View ", `<View ViewType="CALENDAR" `, -1)

	// fmt.Println(viewXML)

	// viewXML := `<View Name="{EC2C6A0E-1D06-4BF3-A9D5-D9003EF981CC}" MobileView="TRUE" Type="CALENDAR" TabularView="FALSE" Scope="Recursive" RecurrenceRowset="TRUE" DisplayName="Test_View_03" Url="/sites/ci/Lists/Calendar01/Test_View_03.aspx" Level="1" BaseViewID="2" ContentTypeID="0x" MobileUrl="_layouts/15/mobile/viewdaily.aspx" ImageUrl="/_layouts/15/images/events.png?rev=47" ><Query><Where><DateRangesOverlap><FieldRef Name="EventDate" /><FieldRef Name="EndDate" /><FieldRef Name="RecurrenceID" /><Value Type="DateTime"><Month /></Value></DateRangesOverlap></Where></Query><ViewFields><FieldRef Name="EventDate" /><FieldRef Name="EndDate" /><FieldRef Name="Title" /><FieldRef Name="fRecurrence" Explicit="TRUE" /></ViewFields><RowLimit>0</RowLimit><Aggregations Value="Off" /><CalendarViewStyles>&lt;CalendarViewStyle  Title='Day' Type='day' Template='CalendarViewdayChrome' Sequence='1' Default='FALSE' /&gt;&lt;CalendarViewStyle  Title='Week' Type='week' Template='CalendarViewweekChrome' Sequence='2' Default='FALSE' /&gt;&lt;CalendarViewStyle  Title='Month' Type='month' Template='CalendarViewmonthChrome' Sequence='3' Default='TRUE' /&gt;</CalendarViewStyles><ViewData><FieldRef Name="Title" Type="CalendarMonthTitle" /><FieldRef Name="Title" Type="CalendarWeekTitle" /><FieldRef Name="Location" Type="CalendarWeekLocation" /><FieldRef Name="Title" Type="CalendarDayTitle" /><FieldRef Name="Location" Type="CalendarDayLocation" /></ViewData><Toolbar Type="Standard"/></View>`

	// meta := map[string]interface{}{
	// 	"Title":        "Test_View_03",
	// 	"ListViewXml":  viewXML,
	// 	"ViewData":     `<FieldRef Name="Title" Type="CalendarMonthTitle" /><FieldRef Name="Title" Type="CalendarWeekTitle" /><FieldRef Name="Location" Type="CalendarWeekLocation" /><FieldRef Name="Title" Type="CalendarDayTitle" /><FieldRef Name="Location" Type="CalendarDayLocation" />`,
	// 	"ViewQuery":    `<Where><DateRangesOverlap><FieldRef Name="EventDate" /><FieldRef Name="EndDate" /><FieldRef Name="RecurrenceID" /><Value Type="DateTime"><Month /></Value></DateRangesOverlap></Where>`,
	// 	"BaseViewId":   "2",
	// 	"ViewType":     "CALENDAR",
	// 	"ViewTypeKind": 524288 | 8193,
	// 	"baseViewId":   "2",
	// 	"TabularView":  false,
	// 	"Paged":        false,
	// }
	// body, _ := json.Marshal(meta)

	// if _, err := l.Views().Add(body); err != nil {
	// 	fmt.Printf("error while adding a view: %s\n", err)
	// }

	// if _, err := sp.Web().Lists().GetByTitle("NotExisting").Get(); err != nil {
	// 	log.Fatal(err)
	// }

}
