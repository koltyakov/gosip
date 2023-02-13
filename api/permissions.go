package api

// BasePermissions - Low/High pair of base permissions
type BasePermissions struct {
	High int64 `json:"High,string"`
	Low  int64 `json:"Low,string"`
}

// RoleAssigment role assignments model
type RoleAssigment struct {
	Member *struct {
		LoginName     string
		PrincipalType int
	}
	RoleDefinitionBindings []*RoleDefInfo
}

// PermissionKind enumerator
var PermissionKind = struct {
	// Has no permissions on the Site. Not available through the user interface
	EmptyMask int64

	// View items in lists, documents in document libraries, and Web discussion comments
	ViewListItems int64

	// Add items to lists, documents to document libraries, and Web discussion comments
	AddListItems int64

	// Edit items in lists, edit documents in document libraries, edit Web discussion comments in documents, and customize Web Part Pages in document libraries
	EditListItems int64

	// Delete items from a list, documents from a document library, and Web discussion comments in documents
	DeleteListItems int64

	// Approve a minor version of a list item or document
	ApproveItems int64

	// View the source of documents with server-side file handlers
	OpenItems int64

	// View past versions of a list item or document
	ViewVersions int64

	// Delete past versions of a list item or document
	DeleteVersions int64

	// Discard or check in a document which is checked out to another user
	CancelCheckout int64

	// Create, change, and delete personal views of lists
	ManagePersonalViews int64

	// Create and delete lists, add or remove columns in a list, and add or remove public views of a list
	ManageLists int64

	// View forms, views, and application pages, and enumerate lists
	ViewFormPages int64

	/**
	* Make content of a list or document library retrievable for anonymous users through SharePoint search.
	* The list permissions in the site do not change.
	 */
	AnonymousSearchAccessList int64

	// Allow users to open a Site, list, or folder to access items inside that container
	Open int64

	// View pages in a Site
	ViewPages int64

	// Add, change, or delete HTML pages or Web Part Pages, and edit the Site using a Windows SharePoint Services compatible editor
	AddAndCustomizePages int64

	// Apply a theme or borders to the entire Site
	ApplyThemeAndBorder int64

	// Apply a style sheet (.css file) to the Site
	ApplyStyleSheets int64

	// View reports on Site usage
	ViewUsageData int64

	// Create a Site using Self-Service Site Creation
	CreateSSCSite int64

	// Create subsites such as team sites, Meeting Workspace sites, and Document Workspace sites
	ManageSubwebs int64

	// Create a group of users that can be used anywhere within the site collection
	CreateGroups int64

	// Create and change permission levels on the Site and assign permissions to users and groups
	ManagePermissions int64

	// Enumerate files and folders in a Site using Microsoft Office SharePoint Designer and WebDAV interfaces
	BrowseDirectories int64

	// View information about users of the Site
	BrowseUserInfo int64

	// Add or remove personal Web Parts on a Web Part Page
	AddDelPrivateWebParts int64

	// Update Web Parts to display personalized information
	UpdatePersonalWebParts int64

	/**
	* Grant the ability to perform all administration tasks for the Site as well as
	* manage content, activate, deactivate, or edit properties of Site scoped Features
	* through the object model or through the user interface (UI). When granted on the
	* root Site of a Site Collection, activate, deactivate, or edit properties of
	* site collection scoped Features through the object model. To browse to the Site
	* Collection Features page and activate or deactivate Site Collection scoped Features
	* through the UI, you must be a Site Collection administrator.
	 */
	ManageWeb int64

	/**
	* Content of lists and document libraries in the Web site will be retrievable for anonymous users through
	* SharePoint search if the list or document library has AnonymousSearchAccessList set
	 */
	AnonymousSearchAccessWebLists int64

	// Use features that launch client applications. Otherwise, users must work on documents locally and upload changes
	UseClientIntegration int64

	// Use SOAP, WebDAV, or Microsoft Office SharePoint Designer interfaces to access the Site
	UseRemoteAPIs int64

	// Manage alerts for all users of the Site
	ManageAlerts int64

	// Create e-mail alerts
	CreateAlerts int64

	// Allows a user to change his or her user information, such as adding a picture
	EditMyUserInfo int64

	// Enumerate permissions on Site, list, folder, document, or list item
	EnumeratePermissions int64

	// Has all permissions on the Site. Not available through the user interface
	FullMask int64
}{
	EmptyMask:                     0,
	ViewListItems:                 1,
	AddListItems:                  2,
	EditListItems:                 3,
	DeleteListItems:               4,
	ApproveItems:                  5,
	OpenItems:                     6,
	ViewVersions:                  7,
	DeleteVersions:                8,
	CancelCheckout:                9,
	ManagePersonalViews:           10,
	ManageLists:                   12,
	ViewFormPages:                 13,
	AnonymousSearchAccessList:     14,
	Open:                          17,
	ViewPages:                     18,
	AddAndCustomizePages:          19,
	ApplyThemeAndBorder:           20,
	ApplyStyleSheets:              21,
	ViewUsageData:                 22,
	CreateSSCSite:                 23,
	ManageSubwebs:                 24,
	CreateGroups:                  25,
	ManagePermissions:             26,
	BrowseDirectories:             27,
	BrowseUserInfo:                28,
	AddDelPrivateWebParts:         29,
	UpdatePersonalWebParts:        30,
	ManageWeb:                     31,
	AnonymousSearchAccessWebLists: 32,
	UseClientIntegration:          37,
	UseRemoteAPIs:                 38,
	ManageAlerts:                  39,
	CreateAlerts:                  40,
	EditMyUserInfo:                41,
	EnumeratePermissions:          63,
	FullMask:                      65,
}

// HasPermissions checks if base permissions include permissions kind mask
// permissionsKind is represented with in64 value (use PermissionKind struct helper as enumerator)
func HasPermissions(basePermissions BasePermissions, permissionsKind int64) bool {
	if permissionsKind == 0 {
		return true
	}

	perm := uint64(permissionsKind - 1)
	num := uint64(1)
	low := uint64(basePermissions.Low)
	high := uint64(basePermissions.High)

	if permissionsKind == PermissionKind.FullMask {
		return (high&32767) == 32767 && low == 65535
	}

	// if perm >= 0 && perm < 32 {
	if perm < 32 {
		num = num << perm
		return 0 != (low & num)
	} else if perm >= 32 && perm < 64 {
		num = num<<perm - 32
		return 0 != (high & num)
	}

	return false
}
