package csom

import "testing"

func TestObject(t *testing.T) {

	template := `<Property Id="{{.ID}}" ParentId="{{.ParentID}}" Name="Web" />`
	o := NewObject(template)

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

	t.Run("Template", func(t *testing.T) {
		o.SetID(2)
		o.SetParentID(1)
		if o.Template() != template {
			t.Error("error getting object's template")
		}
	})

	t.Run("CheckErr", func(t *testing.T) {
		o := NewObject(`<Property Id="{{.ID}}" ParentId="{{.IncorrectID}}" Name="Web" />`)
		_ = o.String()
		if o.CheckErr() == nil {
			t.Error("should throw an error")
		}
	})

	t.Run("NewObjectProperty", func(t *testing.T) {
		shouldBe := `<Property Id="2" ParentId="1" Name="Web" />`
		o := NewObjectProperty("Web")
		o.SetID(2)
		o.SetParentID(1)
		if o.String() != shouldBe {
			t.Error("wrong object property")
		}
		if err := o.CheckErr(); err != nil {
			t.Error(err)
		}
	})

	t.Run("NewObjectMethod", func(t *testing.T) {
		shouldBe := `<Method Id="2" ParentId="1" Name="Add"><Parameters><Parameter /></Parameters></Method>`
		o := NewObjectMethod("Add", []string{"<Parameter />"})
		o.SetID(2)
		o.SetParentID(1)
		if o.String() != shouldBe {
			t.Error("wrong object property")
		}
		if err := o.CheckErr(); err != nil {
			t.Error(err)
		}
	})

	t.Run("NewObjectIdentity", func(t *testing.T) {
		shouldBe := `<Identity Id="2" Name=":identity:" />`
		o := NewObjectIdentity(":identity:")
		o.SetID(2)
		o.SetParentID(1)
		if o.String() != shouldBe {
			t.Error("wrong object property")
		}
		if err := o.CheckErr(); err != nil {
			t.Error(err)
		}
	})

}
