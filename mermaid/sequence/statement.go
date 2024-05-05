package sequence

import (
	"fmt"
)

// LoopStart add a loop to the sequence diagram.
func (d *Diagram) LoopStart(description string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    loop %s", description))
	return d
}

// LoopEnd add a loop to the sequence diagram.
func (d *Diagram) LoopEnd() *Diagram {
	d.body = append(d.body, "    end")
	return d
}

// AltStart add a alt to the sequence diagram.
func (d *Diagram) AltStart(description string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    alt %s", description))
	return d
}

// AltElse add a alt to the sequence diagram.
func (d *Diagram) AltElse(description string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    else %s", description))
	return d
}

// AltEnd add a alt to the sequence diagram.
func (d *Diagram) AltEnd() *Diagram {
	d.body = append(d.body, "    end")
	return d
}

// OptStart add a opt to the sequence diagram.
func (d *Diagram) OptStart(description string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    opt %s", description))
	return d
}

// OptEnd add a opt to the sequence diagram.
func (d *Diagram) OptEnd() *Diagram {
	d.body = append(d.body, "    end")
	return d
}

// ParallelStart add a parallel to the sequence diagram.
func (d *Diagram) ParallelStart(description string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    par %s", description))
	return d
}

// ParallelAnd add a parallel to the sequence diagram.
func (d *Diagram) ParallelAnd(description string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    and %s", description))
	return d
}

// ParallelEnd add a parallel to the sequence diagram.
func (d *Diagram) ParallelEnd() *Diagram {
	d.body = append(d.body, "    end")
	return d
}

// CriticalStart add a critical to the sequence diagram.
func (d *Diagram) CriticalStart(description string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    critical %s", description))
	return d
}

// CriticalOption add a critical opiton to the sequence diagram.
func (d *Diagram) CriticalOption(description string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    option %s", description))
	return d
}

// CriticalEnd add a critical to the sequence diagram.
func (d *Diagram) CriticalEnd() *Diagram {
	d.body = append(d.body, "    end")
	return d
}

// BreakStart add a break to the sequence diagram.
func (d *Diagram) BreakStart(description string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    break %s", description))
	return d
}

// BreakEnd add a break to the sequence diagram.
func (d *Diagram) BreakEnd() *Diagram {
	d.body = append(d.body, "    end")
	return d
}
