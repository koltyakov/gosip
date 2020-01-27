package csom

import "testing"

func TestObject(t *testing.T) {

	o := NewObject(`<Property Id="{{.ID}}" ParentId="{{.ParentID}}" Name="Web" />`)

	t.Run("ID", func(t *testing.T) {
		o.SetID(2)
		if o.GetID() != 2 {
			t.Error("can't set/get object ID")
		}
	})

	t.Run("ParentID", func(t *testing.T) {
		o.SetParentID(1)
		if o.GetParentID() != 1 {
			t.Error("can't set/get object ParentID")
		}
	})

	t.Run("String", func(t *testing.T) {
		o.SetID(2)
		o.SetParentID(1)
		if o.String() != `<Property Id="2" ParentId="1" Name="Web" />` {
			t.Error("template compilation error")
		}
	})

	t.Run("CheckErr", func(t *testing.T) {
		o := NewObject(`<Property Id="{{.ID}}" ParentId="{{.IncorrectID}}" Name="Web" />`)
		_ = o.String()
		if o.CheckErr() == nil {
			t.Error("should throw an error")
		}
	})

}
