package csom

import (
	"bytes"
	"text/template"
)

// Object ...
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

// NewObject ...
func NewObject(template string) Object {
	o := &object{}
	o.template = template
	return o
}

func (o *object) String() string {
	o.err = nil

	template, err := template.New("objectPath").Parse(o.template)
	if err != nil {
		o.err = err
		return o.template
	}

	data := &struct {
		ID       int
		ParentID int
	}{
		ID:       o.GetID(),
		ParentID: o.GetParentID(),
	}

	var tpl bytes.Buffer
	if err := template.Execute(&tpl, data); err != nil {
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
