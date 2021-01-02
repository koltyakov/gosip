package csom

import (
	"testing"
)

func TestCSOMBuilder(t *testing.T) {
	b := NewBuilder()

	b.AddObject(NewObject(`<Property Id="{{.ID}}" ParentId="{{.ParentID}}" Name="Web" />`), nil)
	b.AddObject(NewObject(`
		<Method Id="{{.ID}}" ParentId="{{.ParentID}}" Name="GetFolderByServerRelativeUrl">
			<Parameters>
				<Parameter Type="String">/sites/site/Lists/List</Parameter>
			</Parameters>
		</Method>
	`), nil)
	b.AddAction(NewAction(`
		<Query Id="{{.ID}}" ObjectPathId="{{.ObjectID}}">
			<Query SelectAllProperties="true">
				<Properties />
			</Query>
		</Query>
	`), nil)

	pkg, err := b.Compile()
	if err != nil {
		t.Error(err)
	}

	csomPkg := `<Request xmlns="http://schemas.microsoft.com/sharepoint/clientquery/2009" SchemaVersion="15.0.0.0" LibraryVersion="16.0.0.0" ApplicationName="Gosip"><Actions><Query Id="3" ObjectPathId="2"><Query SelectAllProperties="true"><Properties /></Query></Query></Actions><ObjectPaths><StaticProperty Id="0" TypeId="{3747adcd-a3c3-41b9-bfab-4a64dd2f1e0a}" Name="Current" /><Property Id="1" ParentId="0" Name="Web" /><Method Id="2" ParentId="1" Name="GetFolderByServerRelativeUrl"><Parameters><Parameter Type="String">/sites/site/Lists/List</Parameter></Parameters></Method></ObjectPaths></Request>`
	if pkg != csomPkg {
		t.Error("incorrect package")
	}
}

func TestCSOMGetObjectID(t *testing.T) {
	b := NewBuilder()

	b.AddObject(NewObject(`<Property Id="{{.ID}}" ParentId="{{.ParentID}}" Name="Web" />`), nil)
	obj := NewObject(`
		<Method Id="{{.ID}}" ParentId="{{.ParentID}}" Name="GetFolderByServerRelativeUrl">
			<Parameters>
				<Parameter Type="String">/sites/site/Lists/List</Parameter>
			</Parameters>
		</Method>
	`)
	b.AddObject(obj, nil)

	objID, err := b.GetObjectID(obj)
	if err != nil {
		t.Error(err)
	}

	if objID != 2 {
		t.Error("wrong object ID")
	}

	incorrectObj := NewObject(`<Property Id="{{.ID}}" ParentId="{{.Incorrect}}" Name="Web" />`)
	b.AddObject(incorrectObj, nil)

	if _, err := b.GetObjectID(incorrectObj); err == nil {
		t.Error("should throw an error")
	}
}

func TestCSOMGetObjects(t *testing.T) {
	b := NewBuilder()

	b.AddObject(NewObject(`<Property Id="{{.ID}}" ParentId="{{.ParentID}}" Name="Web" />`), nil)
	b.AddObject(NewObject(`<Property Id="{{.ID}}" ParentId="{{.ParentID}}" Name="Lists" />`), nil)

	if len(b.GetObjects()) != 3 {
		t.Error("wrong objects count")
	}
}

func TestCSOMCompileError(t *testing.T) {
	b := NewBuilder()

	b.AddObject(NewObject(`<Property Id="{{.ID}}" ParentId="{{.ParentID}}" Name="Web" />`), nil)
	b.AddObject(NewObject(`
		<Method Id="{{.ID}}" ParentId="{{.IncorrectID}}" Name="GetFolderByServerRelativeUrl">
			<Parameters>
				<Parameter Type="String">/sites/site/Lists/List</Parameter>
			</Parameters>
		</Method>
	`), nil)

	if _, err := b.Compile(); err == nil {
		t.Error("should throw an error")
	}
}

func TestCSOMClone(t *testing.T) {
	b := NewBuilder().(*builder)

	b.AddObject(NewObject(`<Property Id="{{.ID}}" ParentId="{{.ParentID}}" Name="Web" />`), nil)
	b.AddAction(NewAction(`<Query Id="{{.ID}}" ObjectPathId="{{.ObjectID}}" />`), nil)

	nb := b.Clone().(*builder)
	b.AddObject(NewObject(`<Property Id="{{.ID}}" ParentId="{{.ParentID}}" Name="Web" />`), nil)
	b.AddAction(NewAction(`<Query Id="{{.ID}}" ObjectPathId="{{.ObjectID}}" />`), nil)

	if len(b.objects) == len(nb.objects) || len(b.actions) == len(nb.actions) {
		t.Error("error cloning CSOM builder")
	}
}
