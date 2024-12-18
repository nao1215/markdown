package arch

import "fmt"

// Service adds a service to the architecture diagram.
// Syntax: service {service id}({icon name})[{title}]
func (a *Architecture) Service(serviceID string, icon Icon, title string) *Architecture {
	a.body = append(a.body, fmt.Sprintf("    service %s(%s)[%s]", serviceID, icon, title))
	return a
}

// ServiceInGroup adds a service in a group to the architecture diagram.
func (a *Architecture) ServiceInGroup(serviceID string, icon Icon, title, groupID string) *Architecture {
	a.body = append(a.body, fmt.Sprintf("    service %s(%s)[%s] in %s", serviceID, icon, title, groupID))
	return a
}
