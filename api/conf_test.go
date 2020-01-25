package api

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestConf(t *testing.T) {
	checkClient(t)

	hs := map[string]*RequestConfig{
		"nometadata":      HeadersPresets.Nometadata,
		"minimalmetadata": HeadersPresets.Minimalmetadata,
		"verbose":         HeadersPresets.Verbose,
	}

	apiConstructors := map[string]interface{}{
		"SP":               NewSP(spClient),
		"Web":              NewWeb(spClient, "", nil),
		"Webs":             NewWebs(spClient, "", nil),
		"List":             NewList(spClient, "", nil),
		"Lists":            NewLists(spClient, "", nil),
		"Attachment":       NewAttachment(spClient, "", nil),
		"Attachments":      NewAttachments(spClient, "", nil),
		"AssociatedGroups": NewAssociatedGroups(spClient, "", nil),
		"Changes":          NewChanges(spClient, "", nil),
		"ContentType":      NewContentType(spClient, "", nil),
		"ContentTypes":     NewContentTypes(spClient, "", nil),
		"Context":          NewContext(spClient, "", nil),
		"CustomAction":     NewCustomAction(spClient, "", nil),
		"CustomActions":    NewCustomActions(spClient, "", nil),
		"EventReceivers":   NewEventReceivers(spClient, "", nil),
		"Features":         NewFeatures(spClient, "", nil),
		"Field":            NewField(spClient, "", nil),
		"FieldLink":        NewFieldLink(spClient, "", nil),
		"FieldLinks":       NewFieldLinks(spClient, "", nil),
		"Fields":           NewFields(spClient, "", nil, ""),
		"File":             NewFile(spClient, "", nil),
		"Files":            NewFiles(spClient, "", nil),
		"Folder":           NewFolder(spClient, "", nil),
		"Folders":          NewFolders(spClient, "", nil),
		"Group":            NewGroup(spClient, "", nil),
		"Groups":           NewGroups(spClient, "", nil),
		"Item":             NewItem(spClient, "", nil),
		"Items":            NewItems(spClient, "", nil),
		"Profiles":         NewProfiles(spClient, "", nil),
		"Properties":       NewProperties(spClient, "", nil, ""),
		"Records":          NewRecords(NewItem(spClient, "", nil)),
		"RecycleBin":       NewRecycleBin(spClient, "", nil),
		"RecycleBinItem":   NewRecycleBinItem(spClient, "", nil),
		"RoleDefinitions":  NewRoleDefinitions(spClient, "", nil),
		"Roles":            NewRoles(spClient, "", nil),
		"Search":           NewSearch(spClient, "", nil),
		"Site":             NewSite(spClient, "", nil),
		"User":             NewUser(spClient, "", nil),
		"Users":            NewUsers(spClient, "", nil),
		"Utility":          NewUtility(spClient, "", nil),
		"View":             NewView(spClient, "", nil),
		"Views":            NewViews(spClient, "", nil),
		"HTTPClient":       NewHTTPClient(spClient),
		"ODataMods":        NewODataMods(),
	}

	withNoContMethod := []string{}

	for key, obj := range apiConstructors {
		method := reflect.ValueOf(obj).MethodByName("Conf")
		if method.IsValid() {
			t.Run(strings.Replace(fmt.Sprintf("Conf/%T", obj), "*api.", "", 1), func(t *testing.T) {
				for key, preset := range hs {
					config := method.
						Call([]reflect.Value{reflect.ValueOf(preset)})[0].
						Elem().
						FieldByName("config")
					if fmt.Sprintf("%+v", config) != fmt.Sprintf("%+v", preset) {
						t.Errorf("can't config %v", key)
					}
				}
			})
		} else {
			withNoContMethod = append(withNoContMethod, key)
		}
	}

	if len(withNoContMethod) > 0 {
		t.Logf("the following constructors don't contain Conf method, but this is OK: %v\n", withNoContMethod)
	}

	missedConstructors := []string{}
	for _, constructor := range getAllConstructors() {
		found := false
		for key := range apiConstructors {
			if key == strings.Replace(constructor, "New", "", 1) {
				found = true
			}
		}
		if !found {
			missedConstructors = append(missedConstructors, constructor)
		}
	}
	if len(missedConstructors) > 0 {
		t.Logf("the following API constructors are not covered: %v\n", missedConstructors)
	}

}

func getAllConstructors() []string {
	_, filename, _, _ := runtime.Caller(1)
	pkgPath := path.Dir(filename)

	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, pkgPath, nil, 0)
	if err != nil {
		fmt.Println("Failed to parse package:", err)
		os.Exit(1)
	}

	constructors := []string{}

	for _, pack := range packs {
		for _, f := range pack.Files {
			for _, s := range f.Scope.Objects {
				if s.Kind.String() == "func" {
					if strings.Index(s.Name, "New") == 0 {
						constructors = append(constructors, s.Name)
					}
				}
			}
		}
	}

	return constructors
}
