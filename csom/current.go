package csom

type current struct{}

func (cw *current) String() string {
	return `<StaticProperty Id="0" TypeId="{3747adcd-a3c3-41b9-bfab-4a64dd2f1e0a}" Name="Current" />`
}

func (cw *current) SetID(id int) {}

func (cw *current) GetID() int {
	return 0
}

func (cw *current) SetParentID(parentID int) {}

func (cw *current) GetParentID() int {
	return -1
}

func (cw *current) CheckErr() error {
	return nil
}
