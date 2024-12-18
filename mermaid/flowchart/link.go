package flowchart

import "fmt"

// LinkWithArrowHead adds a link with an arrow head to the flowchart.
func (f *Flowchart) LinkWithArrowHead(from, to string) *Flowchart {
	f.body = append(f.body, fmt.Sprintf("    %s-->%s", from, to))
	return f
}

// LinkWithArrowHeadAndText adds a link with an arrow head and text to the flowchart.
func (f *Flowchart) LinkWithArrowHeadAndText(from, to, text string) *Flowchart {
	f.body = append(f.body, fmt.Sprintf("    %s-->|\"%s\"|%s", from, text, to))
	return f
}

// OpenLink adds an open link to the flowchart.
func (f *Flowchart) OpenLink(from, to string) *Flowchart {
	f.body = append(f.body, fmt.Sprintf("    %s --- %s", from, to))
	return f
}

// OpenLinkWithText adds an open link with text to the flowchart.
func (f *Flowchart) OpenLinkWithText(from, to, text string) *Flowchart {
	f.body = append(f.body, fmt.Sprintf("    %s---|\"%s\"|%s", from, text, to))
	return f
}

// DottedLink adds a dotted link to the flowchart.
func (f *Flowchart) DottedLink(from, to string) *Flowchart {
	f.body = append(f.body, fmt.Sprintf("    %s-.->%s", from, to))
	return f
}

// DottedLinkWithText adds a dotted link with text to the flowchart.
func (f *Flowchart) DottedLinkWithText(from, to, text string) *Flowchart {
	f.body = append(f.body, fmt.Sprintf("    %s-. \"%s\" .-> %s", from, text, to))
	return f
}

// ThickLink adds a thick link to the flowchart.
func (f *Flowchart) ThickLink(from, to string) *Flowchart {
	f.body = append(f.body, fmt.Sprintf("    %s ==> %s", from, to))
	return f
}

// ThickLinkWithText adds a thick link with text to the flowchart.
func (f *Flowchart) ThickLinkWithText(from, to, text string) *Flowchart {
	f.body = append(f.body, fmt.Sprintf("    %s == \"%s\" ==> %s", from, text, to))
	return f
}

// InvisibleLink adds an invisible link to the flowchart.
func (f *Flowchart) InvisibleLink(from, to string) *Flowchart {
	f.body = append(f.body, fmt.Sprintf("    %s ~~~ %s", from, to))
	return f
}
