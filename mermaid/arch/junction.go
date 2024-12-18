package arch

import "fmt"

// Junction are a special type of node which acts as a potential 4-way split between edges.
// Syntax: junction {junction id} (in {parent id})?
func (a *Architecture) Junction(junctionID string) *Architecture {
	a.body = append(a.body, fmt.Sprintf("    junction %s", junctionID))
	return a
}

// JunctionsInParent adds a junction in a parent group to the architecture diagram.
// Syntax: junction {junction id} in {parent id}
func (a *Architecture) JunctionsInParent(junctionID, parentGroupID string) *Architecture {
	a.body = append(a.body, fmt.Sprintf("    junction %s in %s", junctionID, parentGroupID))
	return a
}
