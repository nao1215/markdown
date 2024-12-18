package flowchart

import "fmt"

// Node adds a node to the flowchart.
func (f *Flowchart) Node(name string) *Flowchart {
	f.body = append(f.body, fmt.Sprintf("    %s", name))
	return f
}

// NodeWithText adds a node with text to the flowchart.
// Unicode characters are supported.
func (f *Flowchart) NodeWithText(name, text string) *Flowchart {
	f.body = append(f.body, fmt.Sprintf("    %s[\"%s\"]", name, text))
	return f
}

// NodeWithMarkdown adds a node with markdown text to the flowchart.
func (f *Flowchart) NodeWithMarkdown(name, markdownText string) *Flowchart {
	f.body = append(f.body, fmt.Sprintf("    %s[\"`%s`\"]", name, markdownText))
	return f
}

// NodeWithNewLines adds a node with new lines to the flowchart.
func (f *Flowchart) NodeWithNewLines(name, textWithNewLines string) *Flowchart {
	f.body = append(f.body, fmt.Sprintf("    %s[\"`%s`\"]", name, textWithNewLines))
	return f
}

// RoundEdgesNode adds a node with round edges to the flowchart.
func (f *Flowchart) RoundEdgesNode(name, text string) *Flowchart {
	f.body = append(f.body, fmt.Sprintf("    %s(\"%s\")", name, text))
	return f
}

// StadiumNode adds a node with stadium shape to the flowchart.
func (f *Flowchart) StadiumNode(name, text string) *Flowchart {
	f.body = append(f.body, fmt.Sprintf("    %s([\"%s\"])", name, text))
	return f
}

// SubroutineNode adds a node with subroutine shape to the flowchart.
func (f *Flowchart) SubroutineNode(name, text string) *Flowchart {
	f.body = append(f.body, fmt.Sprintf("    %s[[\"%s\"]]", name, text))
	return f
}

// CylindricalNode adds a node with cylindrical shape to the flowchart.
func (f *Flowchart) CylindricalNode(name, text string) *Flowchart {
	f.body = append(f.body, fmt.Sprintf("    %s[(\"%s\")]", name, text))
	return f
}

// DatabaseNode adds a node with database shape to the flowchart.
// This method is same as CylindricalShapeNode()
func (f *Flowchart) DatabaseNode(name, text string) *Flowchart {
	return f.CylindricalNode(name, text)
}

// CircleNode adds a node with circle shape to the flowchart.
func (f *Flowchart) CircleNode(name, text string) *Flowchart {
	f.body = append(f.body, fmt.Sprintf("    %s((\"%s\"))", name, text))
	return f
}

// AsymmetricNode adds a node with asymmetric shape to the flowchart.
func (f *Flowchart) AsymmetricNode(name, text string) *Flowchart {
	f.body = append(f.body, fmt.Sprintf("    %s>\"%s\"]", name, text))
	return f
}

// RhombusNode adds a node with rhombus shape to the flowchart.
func (f *Flowchart) RhombusNode(name, text string) *Flowchart {
	f.body = append(f.body, fmt.Sprintf("    %s{\"%s\"}", name, text))
	return f
}

// HexagonNode adds a node with hexagon shape to the flowchart.
func (f *Flowchart) HexagonNode(name, text string) *Flowchart {
	f.body = append(f.body, fmt.Sprintf("    %s{{\"%s\"}}", name, text))
	return f
}

// ParallelogramNode adds a node with parallelogram shape to the flowchart.
func (f *Flowchart) ParallelogramNode(name, text string) *Flowchart {
	f.body = append(f.body, fmt.Sprintf("    %s[/\"%s\"/]", name, text))
	return f
}

// ParallelogramAltNode adds a node with parallelogram shape to the flowchart.
func (f *Flowchart) ParallelogramAltNode(name, text string) *Flowchart {
	f.body = append(f.body, fmt.Sprintf("    %s[\\\"%s\"\\]", name, text))
	return f
}

// TrapezoidNode adds a node with trapezoid shape to the flowchart.
func (f *Flowchart) TrapezoidNode(name, text string) *Flowchart {
	f.body = append(f.body, fmt.Sprintf("    %s[/\"%s\"\\]", name, text))
	return f
}

// TrapezoidAltNode adds a node with trapezoid shape to the flowchart.
func (f *Flowchart) TrapezoidAltNode(name, text string) *Flowchart {
	f.body = append(f.body, fmt.Sprintf("    %s[\\\"%s\"/]", name, text))
	return f
}

// DoubleCircleNode adds a node with double circle shape to the flowchart.
func (f *Flowchart) DoubleCircleNode(name, text string) *Flowchart {
	f.body = append(f.body, fmt.Sprintf("    %s(((\"%s\")))", name, text))
	return f
}
