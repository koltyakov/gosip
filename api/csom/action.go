package csom

import (
	"bytes"
	"text/template"
)

// Action CSOM XML action node builder interface
type Action interface {
	String() string
	SetID(id int)
	GetID() int
	SetObjectID(objectID int)
	GetObjectID() int
	CheckErr() error
}

type action struct {
	template string
	id       int
	objectID int
	err      error
}

// NewAction creates CSOM XML action node builder instance
func NewAction(template string) Action {
	a := &action{}
	a.template = template
	return a
}

func (a *action) String() string {
	a.err = nil

	template, _ := template.New("action").Parse(a.template)
	// if err != nil {
	// 	a.err = err
	// 	return a.template
	// }

	data := &struct {
		ID       int
		ObjectID int
	}{
		ID:       a.GetID(),
		ObjectID: a.GetObjectID(),
	}

	var tpl bytes.Buffer
	if err := template.Execute(&tpl, data); err != nil {
		a.err = err
		return a.template
	}

	return trimMultiline(tpl.String())
}

func (a *action) SetID(id int) {
	a.id = id
}

func (a *action) GetID() int {
	return a.id
}

func (a *action) SetObjectID(objectID int) {
	a.objectID = objectID
}

func (a *action) GetObjectID() int {
	return a.objectID
}

func (a *action) CheckErr() error {
	return a.err
}
