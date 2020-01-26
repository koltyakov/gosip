package csom

import (
	"strings"
)

// Builder ...
type Builder interface {
	AddObject(object Object, parent Object) Builder
	AddAction(action Action, parent Object) Builder
	Compile() (string, error)
}

type builder struct {
	objects []*objectsEdge
	actions []*actionEdge
}

type objectsEdge struct {
	Current Object
	Parent  Object
}

type actionEdge struct {
	Action Action
	Object Object
}

// NewBuilder ...
func NewBuilder() Builder {
	b := &builder{}
	b.AddObject(&current{}, nil)
	return b
}

func (b *builder) AddObject(object Object, parent Object) Builder {
	if parent == nil && len(b.objects) > 0 {
		parent = b.objects[len(b.objects)-1].Current
	}
	b.objects = append(b.objects, &objectsEdge{
		Current: object,
		Parent:  parent,
	})
	return b
}

func (b *builder) AddAction(action Action, object Object) Builder {
	if object == nil && len(b.objects) > 0 {
		object = b.objects[len(b.objects)-1].Current
	}
	b.actions = append(b.actions, &actionEdge{
		Action: action,
		Object: object,
	})
	return b
}

func (b *builder) Compile() (string, error) {
	objects := ""
	actions := ""
	for i, edge := range b.objects {
		if i > 1 {
			if edge.Parent.GetID() == 0 {
				edge.Parent.SetID(b.nextObjectID())
			}
		}
		if i > 0 {
			if edge.Current.GetID() == 0 {
				edge.Current.SetID(b.nextObjectID())
				edge.Current.SetParentID(edge.Parent.GetID())
			}
		}
		objects += edge.Current.String()
	}
	for _, edge := range b.actions {
		if edge.Action.GetID() == 0 {
			edge.Action.SetID(b.nextActionID())
			edge.Action.SetObjectID(edge.Object.GetID())
		}
		actions += edge.Action.String()
	}
	csomPkg := trimMultiline(`
		<Request xmlns="http://schemas.microsoft.com/sharepoint/clientquery/2009" SchemaVersion="15.0.0.0" LibraryVersion="16.0.0.0" ApplicationName="Javascript Library">
			<Actions>` + actions + `</Actions>
			<ObjectPaths>` + objects + `</ObjectPaths>
		</Request>
	`)
	return csomPkg, nil
}

func (b *builder) nextObjectID() int {
	nextID := 0
	for _, edge := range b.objects {
		if edge.Parent != nil && nextID <= edge.Parent.GetID() {
			nextID = edge.Parent.GetID() + 1
		}
		if nextID <= edge.Current.GetID() {
			nextID = edge.Current.GetID() + 1
		}
	}
	return nextID
}

func (b *builder) nextActionID() int {
	nextID := b.nextObjectID()
	for _, edge := range b.actions {
		if nextID <= edge.Action.GetID() {
			nextID = edge.Action.GetID() + 1
		}
	}
	return nextID
}

func trimMultiline(multi string) string {
	res := ""
	for _, line := range strings.Split(multi, "\n") {
		res += strings.Trim(line, "\t")
	}
	return res
}
