package flowchart

import (
	"errors"
	"fmt"
)

// Subgraph opens a subgraph. Everything added until the matching SubgraphEnd is
// drawn inside it.
//
// Subgraphs nest, and the pair is explicit rather than a callback taking a
// nested builder: the chain stays flat, which is how every other builder in
// this library reads, and how mermaid/c4 opens a boundary.
//
// The title is quoted and escaped the same way a node label is, so it can hold
// whatever a label can. mermaid takes an unquoted title too, but only for the
// titles that need no quoting, which is not something a caller should have to
// think about.
func (f *Flowchart) Subgraph(id, title string) *Flowchart {
	if f.err != nil {
		return f
	}

	f.body = append(f.body, fmt.Sprintf(`%ssubgraph %s["%s"]`, f.indent(), id, escapePlainText(title)))
	f.depth++
	return f
}

// SubgraphEnd closes the subgraph opened last.
//
// Calling it outside a subgraph is an error rather than a silent no-op: a chain
// that has closed more subgraphs than it opened is not the diagram its author
// meant, and the "end" it would write loses the whole drawing.
func (f *Flowchart) SubgraphEnd() *Flowchart {
	if f.err != nil {
		return f
	}
	if f.depth == 0 {
		f.setError(errors.New("SubgraphEnd was called outside a subgraph; there is nothing to close"))
		return f
	}

	f.depth--
	f.body = append(f.body, f.indent()+"end")
	return f
}

// Direction is which way a subgraph is laid out.
//
// The flowchart's own direction is set with an option instead, WithOrientalTopDown
// and its siblings, which is how this package has always spelled it. A subgraph
// needs a value rather than an option, because it is set partway through a
// chain, and the name matches what mermaid/state, mermaid/class and
// mermaid/requirement call the same thing.
type Direction string

const (
	// DirectionTB lays the subgraph out from the top to the bottom.
	DirectionTB Direction = "TB"
	// DirectionBT lays the subgraph out from the bottom to the top.
	DirectionBT Direction = "BT"
	// DirectionLR lays the subgraph out from the left to the right.
	DirectionLR Direction = "LR"
	// DirectionRL lays the subgraph out from the right to the left.
	DirectionRL Direction = "RL"
)

// SubgraphDirection sets which way the subgraph opened last is laid out,
// independently of the flowchart around it.
//
// It is the one thing a subgraph is for beyond grouping: a chart that runs top
// to bottom can hold a row of steps that runs left to right.
func (f *Flowchart) SubgraphDirection(direction Direction) *Flowchart {
	if f.err != nil {
		return f
	}
	if f.depth == 0 {
		f.setError(errors.New("SubgraphDirection was called outside a subgraph"))
		return f
	}

	f.body = append(f.body, fmt.Sprintf("%sdirection %s", f.indent(), direction))
	return f
}
