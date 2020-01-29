package csom

import (
	"bytes"
	"fmt"
	"text/template"
)

// Object CSOM XML object path node builder interface
type Object interface {
	String() string
	SetID(id int)
	GetID() int
	SetParentID(parentID int)
	GetParentID() int
	CheckErr() error
}

type object struct {
	template string
	id       int
	parentID int
	err      error
}

// NewObject creates CSOM XML object path node builder instance
func NewObject(template string) Object {
	o := &object{}
	o.template = template
	return o
}

// NewObjectProperty creates CSOM XML object path node builder instance
func NewObjectProperty(propertyName string) Object {
	return NewObject(fmt.Sprintf(`<Property Id="{{.ID}}" ParentId="{{.ParentID}}" Name="%s" />`, propertyName))
}

// NewObjectMethod creates CSOM XML object path node builder instance
func NewObjectMethod(methodName string, parameters []string) Object {
	params := ""
	for _, param := range parameters {
		params += param
	}
	return NewObject(fmt.Sprintf(`
		<Method Id="{{.ID}}" ParentId="{{.ParentID}}" Name="%s">
			<Parameters>%s</Parameters>
		</Method>
	`, methodName, trimMultiline(params)))
}

// NewObjectIdentity creates CSOM XML object path node builder instance
func NewObjectIdentity(identityPath string) Object {
	return NewObject(`<Identity Id="{{.ID}}" Name="` + identityPath + `" />`)
}

func (o *object) String() string {
	o.err = nil

	t, _ := template.New("objectPath").Parse(o.template)

	data := &struct {
		ID       int
		ParentID int
	}{
		ID:       o.GetID(),
		ParentID: o.GetParentID(),
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		o.err = err
		return o.template
	}

	return trimMultiline(tpl.String())
}

func (o *object) SetID(id int) {
	o.id = id
}

func (o *object) GetID() int {
	return o.id
}

func (o *object) SetParentID(parentID int) {
	o.parentID = parentID
}

func (o *object) GetParentID() int {
	return o.parentID
}

func (o *object) CheckErr() error {
	return o.err
}
