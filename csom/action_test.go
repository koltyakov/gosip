package csom

import "testing"

func TestAction(t *testing.T) {

	act := NewAction(`<Query Id="{{.ID}}" ObjectPathId="{{.ObjectID}}"></Query>`)

	t.Run("ID", func(t *testing.T) {
		act.SetID(2)
		if act.GetID() != 2 {
			t.Error("can't set/get action ID")
		}
	})

	t.Run("ObjectID", func(t *testing.T) {
		act.SetObjectID(1)
		if act.GetObjectID() != 1 {
			t.Error("can't set/get object ObjectID")
		}
	})

	t.Run("String", func(t *testing.T) {
		act.SetID(2)
		act.SetObjectID(1)
		if act.String() != `<Query Id="2" ObjectPathId="1"></Query>` {
			t.Error("template compilation error")
		}
	})

	t.Run("CheckErr", func(t *testing.T) {
		a := NewAction(`<Query Id="{{.ID}}" ObjectPathId="{{.IncorrectID}}"></Query>`)
		_ = a.String()
		if a.CheckErr() == nil {
			t.Error("should throw an error")
		}
	})

	t.Run("NewActionIdentityQuery", func(t *testing.T) {
		shouldBe := `<ObjectIdentityQuery Id="2" ObjectPathId="1" />`
		o := NewActionIdentityQuery()
		o.SetID(2)
		o.SetObjectID(1)
		if o.String() != shouldBe {
			t.Error("wrong object property")
		}
		if err := o.CheckErr(); err != nil {
			t.Error(err)
		}
	})

	t.Run("NewActionMethod", func(t *testing.T) {
		shouldBe := `<Method Id="2" ObjectPathId="1" Name="Update"><Parameters><Parameter /></Parameters></Method>`
		o := NewActionMethod("Update", []string{"<Parameter />"})
		o.SetID(2)
		o.SetObjectID(1)
		if o.String() != shouldBe {
			t.Error("wrong object property")
		}
		if err := o.CheckErr(); err != nil {
			t.Error(err)
		}
	})

	t.Run("NewQueryWithProps", func(t *testing.T) {
		shouldBe := `<Query Id="2" ObjectPathId="1"><Query SelectAllProperties="true"><Properties><Property /></Properties></Query></Query>`
		o := NewQueryWithProps([]string{"<Property />"})
		o.SetID(2)
		o.SetObjectID(1)
		if o.String() != shouldBe {
			t.Error("wrong object property")
		}
		if err := o.CheckErr(); err != nil {
			t.Error(err)
		}
	})

}
