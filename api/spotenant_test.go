package api

import (
	"strings"
	"testing"
	"time"
)

// csomSitePropertiesMock is a sample CSOM response from
// GetSitePropertiesFromSharePointByFilters with two child items.
const csomSitePropertiesMock = `[
	{"SchemaVersion":"15.0.0.0","LibraryVersion":"16.0.27104.12006","ErrorInfo":null,"TraceCorrelationId":"a1b2c3d4-e5f6-7890-abcd-ef1234567890"},
	3,
	{"IsNull":false},
	4,
	{"IsNull":false},
	5,
	{
		"_ObjectType_":"Microsoft.Online.SharePoint.TenantAdministration.SPOSitePropertiesEnumerable",
		"NextStartIndex":-1,
		"NextStartIndexFromSharePoint":"SPSiteQuery,d4e5f6a7-b8c9-0123-defa-234567890123,e5f6a7b8-c9d0-1234-efab-345678901234,f6a7b8c9-d0e1-2345-fabc-456789012345",
		"_Child_Items_":[
			{
				"_ObjectType_":"Microsoft.Online.SharePoint.TenantAdministration.SiteProperties",
				"_ObjectIdentity_":"a1b2c3d4-e5f6-7890-abcd-ef1234567890|b2c3d4e5-f6a7-8901-bcde-f12345678901:c3d4e5f6-a7b8-9012-cdef-123456789012\nSiteProperties\nhttps%3a%2f%2fcontoso-my.sharepoint.com%2fpersonal%2fjdoe_contoso_com",
				"AllowDownloadingNonWebViewableFiles":false,
				"AllowEditing":false,
				"AllowSelfServiceUpgrade":true,
				"CompatibilityLevel":15,
				"ConditionalAccessPolicy":0,
				"CreatedTime":"\/Date(2025,6,18,3,15,55,40)\/",
				"DefaultLinkPermission":0,
				"DefaultSharingLinkType":0,
				"DenyAddAndCustomizePages":2,
				"Description":null,
				"GroupId":"\/Guid(00000000-0000-0000-0000-000000000000)\/",
				"HasHolds":false,
				"HubSiteId":"\/Guid(00000000-0000-0000-0000-000000000000)\/",
				"IsHubSite":false,
				"IsTeamsChannelConnected":false,
				"IsTeamsConnected":false,
				"LastContentModifiedDate":"\/Date(2026,2,20,5,1,7,280)\/",
				"Lcid":1033,
				"LockIssue":null,
				"LockReason":0,
				"LockState":"Unlock",
				"Owner":"jdoe@contoso.com",
				"OwnerEmail":null,
				"OwnerLoginName":null,
				"OwnerName":null,
				"PWAEnabled":1,
				"RelatedGroupId":"\/Guid(00000000-0000-0000-0000-000000000000)\/",
				"SensitivityLabel":"\/Guid(00000000-0000-0000-0000-000000000000)\/",
				"SharingCapability":2,
				"SharingDomainRestrictionMode":0,
				"SiteDefinedSharingCapability":2,
				"SiteId":"\/Guid(a7b8c9d0-e1f2-3456-abcd-567890123456)\/",
				"Status":"Active",
				"StorageMaximumLevel":1048576,
				"StorageQuotaType":null,
				"StorageUsage":42,
				"StorageWarningLevel":943718,
				"Template":"SPSPERS#12",
				"TeamsChannelType":0,
				"Title":"John Doe",
				"Url":"https:\u002f\u002fcontoso-my.sharepoint.com\u002fpersonal\u002fjdoe_contoso_com",
				"WebsCount":0
			},
			{
				"_ObjectType_":"Microsoft.Online.SharePoint.TenantAdministration.SiteProperties",
				"_ObjectIdentity_":"a1b2c3d4-e5f6-7890-abcd-ef1234567890|b2c3d4e5-f6a7-8901-bcde-f12345678901:c3d4e5f6-a7b8-9012-cdef-123456789012\nSiteProperties\nhttps%3a%2f%2fcontoso-my.sharepoint.com%2fpersonal%2fasmith_contoso_com",
				"AllowDownloadingNonWebViewableFiles":false,
				"AllowEditing":false,
				"AllowSelfServiceUpgrade":true,
				"CompatibilityLevel":15,
				"ConditionalAccessPolicy":0,
				"CreatedTime":"\/Date(2024,3,10,14,22,30,0)\/",
				"DefaultLinkPermission":0,
				"DefaultSharingLinkType":0,
				"DenyAddAndCustomizePages":2,
				"Description":null,
				"GroupId":"\/Guid(00000000-0000-0000-0000-000000000000)\/",
				"HasHolds":false,
				"HubSiteId":"\/Guid(00000000-0000-0000-0000-000000000000)\/",
				"IsHubSite":false,
				"IsTeamsChannelConnected":false,
				"IsTeamsConnected":false,
				"LastContentModifiedDate":"\/Date(2026,1,15,12,0,0,0)\/",
				"Lcid":1033,
				"LockIssue":null,
				"LockReason":0,
				"LockState":"Unlock",
				"Owner":"asmith@contoso.com",
				"OwnerEmail":null,
				"OwnerLoginName":null,
				"OwnerName":null,
				"PWAEnabled":1,
				"RelatedGroupId":"\/Guid(00000000-0000-0000-0000-000000000000)\/",
				"SensitivityLabel":"\/Guid(00000000-0000-0000-0000-000000000000)\/",
				"SharingCapability":2,
				"SharingDomainRestrictionMode":0,
				"SiteDefinedSharingCapability":2,
				"SiteId":"\/Guid(b8c9d0e1-f2a3-4567-bcde-678901234567)\/",
				"Status":"Active",
				"StorageMaximumLevel":1048576,
				"StorageQuotaType":null,
				"StorageUsage":8192,
				"StorageWarningLevel":943718,
				"Template":"SPSPERS#12",
				"TeamsChannelType":0,
				"Title":"Alice Smith",
				"Url":"https:\u002f\u002fcontoso-my.sharepoint.com\u002fpersonal\u002fasmith_contoso_com",
				"WebsCount":0
			}
		]
	}
]`

func TestBuildSitePropertiesQuery(t *testing.T) {
	filter := &SitePropertiesFilter{
		Filter:              "Url -like 'contoso-my.sharepoint.com/personal/'",
		IncludePersonalSite: IncludePersonalSiteInclude,
		Template:            "SPSPERS",
	}

	xml, err := buildSitePropertiesQuery(filter, "0")
	if err != nil {
		t.Fatalf("buildSitePropertiesQuery: %v", err)
	}
	if xml == "" {
		t.Fatal("expected non-empty XML")
	}

	for _, want := range []string{
		"contoso-my.sharepoint.com/personal/",
		"SPSPERS",
		"GetSitePropertiesFromSharePointByFilters",
		spoTenantTypeID,
		spoSitePropertiesFilterTypeID,
	} {
		if !strings.Contains(xml, want) {
			t.Errorf("XML should contain %q", want)
		}
	}
}


func TestParseSitePropertiesResponse(t *testing.T) {
	t.Run("mocked response", func(t *testing.T) {
		result, err := parseSitePropertiesResponse([]byte(csomSitePropertiesMock))
		if err != nil {
			t.Fatalf("parseSitePropertiesResponse: %v", err)
		}

		// Pagination token
		if strings.TrimSpace(result.nextStartIndex) == "" {
			t.Error("nextStartIndex should be non-empty for continuation token")
		}
		if !strings.HasPrefix(result.nextStartIndex, "SPSiteQuery,") {
			t.Errorf("nextStartIndex should start with 'SPSiteQuery,', got %q", result.nextStartIndex[:20])
		}

		// Child items count
		if len(result.sites) != 2 {
			t.Fatalf("len(sites) = %d, want 2", len(result.sites))
		}

		// First site  - verify Data() unmarshals correctly
		site := result.sites[0].Data()
		if site.Title != "John Doe" {
			t.Errorf("Title = %q, want %q", site.Title, "John Doe")
		}
		if site.URL != "https://contoso-my.sharepoint.com/personal/jdoe_contoso_com" {
			t.Errorf("URL = %q, want decoded unicode path", site.URL)
		}
		if site.Owner != "jdoe@contoso.com" {
			t.Errorf("Owner = %q, want %q", site.Owner, "jdoe@contoso.com")
		}
		if site.Template != "SPSPERS#12" {
			t.Errorf("Template = %q, want %q", site.Template, "SPSPERS#12")
		}
		if site.LockState != "Unlock" {
			t.Errorf("LockState = %q, want %q", site.LockState, "Unlock")
		}
		if site.StorageUsage != 42 {
			t.Errorf("StorageUsage = %d, want 42", site.StorageUsage)
		}
		if site.SiteID != "/Guid(a7b8c9d0-e1f2-3456-abcd-567890123456)/" {
			t.Errorf("SiteID = %q, want wrapped GUID", site.SiteID)
		}
		if site.CleanSiteID() != "a7b8c9d0-e1f2-3456-abcd-567890123456" {
			t.Errorf("CleanSiteID() = %q, want unwrapped GUID", site.CleanSiteID())
		}

		// CSOM date parsing  - CreatedTime: /Date(2025,6,18,3,15,55,40)/ = July 18, 2025
		wantCreated := time.Date(2025, time.July, 18, 3, 15, 55, 40_000_000, time.UTC)
		if got := site.CreatedTimeDate(); !got.Equal(wantCreated) {
			t.Errorf("CreatedTimeDate() = %v, want %v", got, wantCreated)
		}

		// LastContentModifiedDate: /Date(2026,2,20,5,1,7,280)/ = March 20, 2026
		wantModified := time.Date(2026, time.March, 20, 5, 1, 7, 280_000_000, time.UTC)
		if got := site.LastContentModifiedDateTime(); !got.Equal(wantModified) {
			t.Errorf("LastContentModifiedDateTime() = %v, want %v", got, wantModified)
		}

		// Normalized() returns raw bytes for custom unmarshalling
		raw := result.sites[0].Normalized()
		if len(raw) == 0 {
			t.Error("Normalized() should return non-empty bytes")
		}

		// Second site - verify Data() unmarshals correctly
		site2 := result.sites[1].Data()
		if site2.Title != "Alice Smith" {
			t.Errorf("sites[1].Title = %q, want %q", site2.Title, "Alice Smith")
		}
		if site2.URL != "https://contoso-my.sharepoint.com/personal/asmith_contoso_com" {
			t.Errorf("sites[1].URL = %q, want decoded unicode path", site2.URL)
		}
		if site2.Owner != "asmith@contoso.com" {
			t.Errorf("sites[1].Owner = %q, want %q", site2.Owner, "asmith@contoso.com")
		}
		if site2.Template != "SPSPERS#12" {
			t.Errorf("sites[1].Template = %q, want %q", site2.Template, "SPSPERS#12")
		}
		if site2.LockState != "Unlock" {
			t.Errorf("sites[1].LockState = %q, want %q", site2.LockState, "Unlock")
		}
		if site2.StorageUsage != 8192 {
			t.Errorf("sites[1].StorageUsage = %d, want 8192", site2.StorageUsage)
		}
		if site2.SiteID != "/Guid(b8c9d0e1-f2a3-4567-bcde-678901234567)/" {
			t.Errorf("sites[1].SiteID = %q, want wrapped GUID", site2.SiteID)
		}
		if site2.CleanSiteID() != "b8c9d0e1-f2a3-4567-bcde-678901234567" {
			t.Errorf("sites[1].CleanSiteID() = %q, want unwrapped GUID", site2.CleanSiteID())
		}

		// CSOM date parsing - CreatedTime: /Date(2024,3,10,14,22,30,0)/ = April 10, 2024
		wantCreated2 := time.Date(2024, time.April, 10, 14, 22, 30, 0, time.UTC)
		if got := site2.CreatedTimeDate(); !got.Equal(wantCreated2) {
			t.Errorf("sites[1].CreatedTimeDate() = %v, want %v", got, wantCreated2)
		}

		// LastContentModifiedDate: /Date(2026,1,15,12,0,0,0)/ = February 15, 2026
		wantModified2 := time.Date(2026, time.February, 15, 12, 0, 0, 0, time.UTC)
		if got := site2.LastContentModifiedDateTime(); !got.Equal(wantModified2) {
			t.Errorf("sites[1].LastContentModifiedDateTime() = %v, want %v", got, wantModified2)
		}

		raw2 := result.sites[1].Normalized()
		if len(raw2) == 0 {
			t.Error("sites[1].Normalized() should return non-empty bytes")
		}
	})

	t.Run("no more pages", func(t *testing.T) {
		input := `[
			{"SchemaVersion":"15.0.0.0","LibraryVersion":"16.0.27104.12006","ErrorInfo":null,"TraceCorrelationId":"abc-123"},
			3,
			{"IsNull":false},
			4,
			{"IsNull":false},
			5,
			{"_ObjectType_":"Microsoft.Online.SharePoint.TenantAdministration.SPOSitePropertiesEnumerable","NextStartIndex":-1,"NextStartIndexFromSharePoint":"","_Child_Items_":[]}
		]`
		result, err := parseSitePropertiesResponse([]byte(input))
		if err != nil {
			t.Fatalf("parseSitePropertiesResponse: %v", err)
		}
		if strings.TrimSpace(result.nextStartIndex) != "" {
			t.Error("nextStartIndex should be empty for last page")
		}
		if len(result.sites) != 0 {
			t.Errorf("len(sites) = %d, want 0", len(result.sites))
		}
	})

	t.Run("empty array", func(t *testing.T) {
		_, err := parseSitePropertiesResponse([]byte(`[]`))
		if err == nil {
			t.Error("expected error for empty array")
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		_, err := parseSitePropertiesResponse([]byte(`not json`))
		if err == nil {
			t.Error("expected error for invalid json")
		}
	})
}

func TestParseCSOMDate(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  time.Time
	}{
		{
			"mocked created time from CSOM response",
			"/Date(2025,6,18,3,15,55,40)/",
			time.Date(2025, time.July, 18, 3, 15, 55, 40_000_000, time.UTC),
		},
		{
			"mocked modified time from CSOM response",
			"/Date(2026,2,20,5,1,7,280)/",
			time.Date(2026, time.March, 20, 5, 1, 7, 280_000_000, time.UTC),
		},
		{
			"december boundary",
			"/Date(2024,11,31,23,59,59,999)/",
			time.Date(2024, time.December, 31, 23, 59, 59, 999_000_000, time.UTC),
		},
		{
			"january zero-based month",
			"/Date(2024,0,15,10,30,0,0)/",
			time.Date(2024, time.January, 15, 10, 30, 0, 0, time.UTC),
		},
		{"empty string", "", time.Time{}},
		{"missing fields", "/Date(2024,0,15)/", time.Time{}},
		{"non-numeric", "/Date(abc,0,15,10,30,0,0)/", time.Time{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCSOMDate(tt.input)
			if !got.Equal(tt.want) {
				t.Errorf("parseCSOMDate(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestSitePropertiesCleanSiteID(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"CSOM format", "/Guid(a7b8c9d0-e1f2-3456-abcd-567890123456)/", "a7b8c9d0-e1f2-3456-abcd-567890123456"},
		{"already clean", "a7b8c9d0-e1f2-3456-abcd-567890123456", "a7b8c9d0-e1f2-3456-abcd-567890123456"},
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp := &SiteProperties{SiteID: tt.input}
			if got := sp.CleanSiteID(); got != tt.want {
				t.Errorf("CleanSiteID(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
