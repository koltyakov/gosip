package api

// BasePermissions - Low/High pair of base permissions
type BasePermissions struct {
	High int64 `json:"High,string"`
	Low  int64 `json:"Low,string"`
}

// RoleAssigment role asigments model
type RoleAssigment struct {
	Member *struct {
		LoginName     string
		PrincipalType int
	}
	RoleDefinitionBindings []*RoleDefInfo
}

// PermissionKind type
type PermissionKind int64

// PermissionKinds enumerator
var PermissionKinds = struct {
	// Has no permissions on the Site. Not available through the user interface
	EmptyMask PermissionKind

	// View items in lists, documents in document libraries, and Web discussion comments
	ViewListItems PermissionKind

	// Add items to lists, documents to document libraries, and Web discussion comments
	AddListItems PermissionKind

	// Edit items in lists, edit documents in document libraries, edit Web discussion comments in documents, and customize Web Part Pages in document libraries
	EditListItems PermissionKind

	// Delete items from a list, documents from a document library, and Web discussion comments in documents
	DeleteListItems PermissionKind

	// Approve a minor version of a list item or document
	ApproveItems PermissionKind

	// View the source of documents with server-side file handlers
	OpenItems PermissionKind

	// View past versions of a list item or document
	ViewVersions PermissionKind

	// Delete past versions of a list item or document
	DeleteVersions PermissionKind

	// Discard or check in a document which is checked out to another user
	CancelCheckout PermissionKind

	// Create, change, and delete personal views of lists
	ManagePersonalViews PermissionKind

	// Create and delete lists, add or remove columns in a list, and add or remove public views of a list
	ManageLists PermissionKind

	// View forms, views, and application pages, and enumerate lists
	ViewFormPages PermissionKind

	/**
	* Make content of a list or document library retrieveable for anonymous users through SharePoint search.
	* The list permissions in the site do not change.
	 */
	AnonymousSearchAccessList PermissionKind

	// Allow users to open a Site, list, or folder to access items inside that container
	Open PermissionKind

	// View pages in a Site
	ViewPages PermissionKind

	// Add, change, or delete HTML pages or Web Part Pages, and edit the Site using a Windows SharePoint Services compatible editor
	AddAndCustomizePages PermissionKind

	// Apply a theme or borders to the entire Site
	ApplyThemeAndBorder PermissionKind

	// Apply a style sheet (.css file) to the Site
	ApplyStyleSheets PermissionKind

	// View reports on Site usage
	ViewUsageData PermissionKind

	// Create a Site using Self-Service Site Creation
	CreateSSCSite PermissionKind

	// Create subsites such as team sites, Meeting Workspace sites, and Document Workspace sites
	ManageSubwebs PermissionKind

	// Create a group of users that can be used anywhere within the site collection
	CreateGroups PermissionKind

	// Create and change permission levels on the Site and assign permissions to users and groups
	ManagePermissions PermissionKind

	// Enumerate files and folders in a Site using Microsoft Office SharePoint Designer and WebDAV interfaces
	BrowseDirectories PermissionKind

	// View information about users of the Site
	BrowseUserInfo PermissionKind

	// Add or remove personal Web Parts on a Web Part Page
	AddDelPrivateWebParts PermissionKind

	// Update Web Parts to display personalized information
	UpdatePersonalWebParts PermissionKind

	/**
	* Grant the ability to perform all administration tasks for the Site as well as
	* manage content, activate, deactivate, or edit properties of Site scoped Features
	* through the object model or through the user interface (UI). When granted on the
	* root Site of a Site Collection, activate, deactivate, or edit properties of
	* site collection scoped Features through the object model. To browse to the Site
	* Collection Features page and activate or deactivate Site Collection scoped Features
	* through the UI, you must be a Site Collection administrator.
	 */
	ManageWeb PermissionKind

	/**
	* Content of lists and document libraries in the Web site will be retrieveable for anonymous users through
	* SharePoint search if the list or document library has AnonymousSearchAccessList set
	 */
	AnonymousSearchAccessWebLists PermissionKind

	// Use features that launch client applications. Otherwise, users must work on documents locally and upload changes
	UseClientIntegration PermissionKind

	// Use SOAP, WebDAV, or Microsoft Office SharePoint Designer interfaces to access the Site
	UseRemoteAPIs PermissionKind

	// Manage alerts for all users of the Site
	ManageAlerts PermissionKind

	// Create e-mail alerts
	CreateAlerts PermissionKind

	// Allows a user to change his or her user information, such as adding a picture
	EditMyUserInfo PermissionKind

	// Enumerate permissions on Site, list, folder, document, or list item
	EnumeratePermissions PermissionKind

	// Has all permissions on the Site. Not available through the user interface
	FullMask PermissionKind
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

// HasPermissions checks if base permissions include permissions kind bits
func HasPermissions(basePermissions BasePermissions, permissionsKind PermissionKind) bool {
	if permissionsKind == 0 {
		return true
	}

	if permissionsKind == PermissionKinds.FullMask {
		return (basePermissions.High&32767) == 32767 && basePermissions.Low == 65535
	}

	perm := permissionsKind - 1
	num := int64(1)

	if perm >= 0 && perm < 32 {
		num = num << perm
		return 0 != (basePermissions.Low & num)
	} else if perm >= 32 && perm < 64 {
		num = num<<perm - 32
		return 0 != (basePermissions.High & num)
	}

	return false
}
