package csom

import "testing"

func TestAction(t *testing.T) {

	a := NewAction(`<Query Id="{{.ID}}" ObjectPathId="{{.ObjectID}}"></Query>`)

	t.Run("ID", func(t *testing.T) {
		a.SetID(2)
		if a.GetID() != 2 {
			t.Error("can't set/get action ID")
		}
	})

	t.Run("ObjectID", func(t *testing.T) {
		a.SetObjectID(1)
		if a.GetObjectID() != 1 {
			t.Error("can't set/get object ObjectID")
		}
	})

	t.Run("String", func(t *testing.T) {
		a.SetID(2)
		a.SetObjectID(1)
		if a.String() != `<Query Id="2" ObjectPathId="1"></Query>` {
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

}
