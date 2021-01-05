package csom

type current struct{}

// Template returns objetc's template
func (cw *current) Template() string {
	return `<StaticProperty Id="0" TypeId="{3747adcd-a3c3-41b9-bfab-4a64dd2f1e0a}" Name="Current" />`
}

// String returns object as string
func (cw *current) String() string { return cw.Template() }

// SetID sets ID
// noinspection GoUnusedParameter
func (cw *current) SetID(id int) {}

// GetID sets ID
func (cw *current) GetID() int { return 0 }

// SetParentID sets parent ID
// noinspection GoUnusedParameter
func (cw *current) SetParentID(parentID int) {}

// GetParentID gets parent ID
func (cw *current) GetParentID() int { return -1 }

// CheckErr checks errors
func (cw *current) CheckErr() error { return nil }
