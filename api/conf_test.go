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

	"github.com/koltyakov/gosip"
)

func TestConf(t *testing.T) {
	checkClient(t)

	hs := map[string]*RequestConfig{
		"nometadata":      HeadersPresets.Nometadata,
		"minimalmetadata": HeadersPresets.Minimalmetadata,
		"verbose":         HeadersPresets.Verbose,
	}

	apiConstructors := getAllConstructors(spClient)

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
	for _, constructor := range getAstConstructors() {
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

func TestModifiers(t *testing.T) {
	checkClient(t)

	modsMethods := []string{
		"Select",
		"Expand",
		"Filter",
		"Skip",
		"Top",
		"OrderBy",
	}

	apiConstructors := getAllConstructors(spClient)
	withNoModsMethod := []string{}

	for key, obj := range apiConstructors {
		mods := []string{}
		for _, modMethodName := range modsMethods {
			method := reflect.ValueOf(obj).MethodByName(modMethodName)
			if method.IsValid() {
				mods = append(mods, modMethodName)
			}
		}

		if len(mods) == 0 {
			withNoModsMethod = append(withNoModsMethod, key)
			continue
		}

		t.Run(strings.Replace(fmt.Sprintf("Mods/%T", obj), "*api.", "", 1), func(t *testing.T) {
			for _, modMethodName := range mods {
				switch modMethodName {
				case "Top":
					obj = reflect.ValueOf(obj).MethodByName(modMethodName).
						Call([]reflect.Value{reflect.ValueOf(1)})[0].
						Interface()
				case "OrderBy":
					obj = reflect.ValueOf(obj).MethodByName(modMethodName).
						Call([]reflect.Value{
							reflect.ValueOf("*"),
							reflect.ValueOf(true),
						})[0].
						Interface()
				default:
					obj = reflect.ValueOf(obj).MethodByName(modMethodName).
						Call([]reflect.Value{reflect.ValueOf("*")})[0].
						Interface()
				}
			}

			m := reflect.ValueOf(obj).Elem().FieldByName("modifiers")
			if !(m.IsValid() && !m.IsNil() && m.Elem().FieldByName("mods").Len() == len(mods)) {
				t.Error("wrong number of modifiers")
			}
		})

	}

	if len(withNoModsMethod) > 0 {
		t.Logf("the following constructors don't contain OData modifiers methods, but this is OK: %v\n", withNoModsMethod)
	}

}

func getAllConstructors(spClient *gosip.SPClient) map[string]interface{} {
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
	return apiConstructors
}

func getAstConstructors() []string {
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
