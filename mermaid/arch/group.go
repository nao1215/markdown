package arch

import "fmt"

// Group adds a group to the architecture diagram.
// Syntax: group {group id}({icon name})[{title}]
func (a *Architecture) Group(groupID string, icon Icon, title string) *Architecture {
	a.body = append(a.body, fmt.Sprintf("    group %s(%s)[%s]", groupID, icon, title))
	return a
}

// GroupInParentGroup adds a group in a parent group to the architecture diagram.
// Syntax: group {group id}({icon name})[{title}] in {parent group id}
func (a *Architecture) GroupInParentGroup(groupID string, icon Icon, title, parentGroupID string) *Architecture {
	a.body = append(a.body, fmt.Sprintf("    group %s(%s)[%s] in %s", groupID, icon, title, parentGroupID))
	return a
}
