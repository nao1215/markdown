package sequence

import "fmt"

// Activate add a participant to the sequence diagram.
func (d *Diagram) Activate(participant string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    activate %s", participant))
	return d
}

// Deactivate add a participant to the sequence diagram.
func (d *Diagram) Deactivate(participant string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    deactivate %s", participant))
	return d
}

// SyncRequestWithActivation add a request to the sequence diagram.
func (d *Diagram) SyncRequestWithActivation(from, to, message string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    %s->>+%s: %s", from, to, message))
	return d
}

// SyncRequestfWithActivation add a request to the sequence diagram.
func (d *Diagram) SyncRequestfWithActivation(from, to, format string, args ...any) *Diagram {
	return d.SyncRequestWithActivation(from, to, fmt.Sprintf(format, args...))
}

// SyncResponseWithActivation add a response to the sequence diagram.
func (d *Diagram) SyncResponseWithActivation(from, to, message string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    %s-->>-%s: %s", from, to, message))
	return d
}

// SyncResponsefWithActivation add a response to the sequence diagram.
func (d *Diagram) SyncResponsefWithActivation(from, to, format string, args ...any) *Diagram {
	return d.SyncResponseWithActivation(from, to, fmt.Sprintf(format, args...))
}

// AsyncRequestWithActivation add a async request to the sequence diagram.
func (d *Diagram) AsyncRequestWithActivation(from, to, message string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    %s->>+%s: %s", from, to, message))
	return d
}

// AsyncRequestfWithActivation add a async request to the sequence diagram.
func (d *Diagram) AsyncRequestfWithActivation(from, to, format string, args ...any) *Diagram {
	return d.AsyncRequestWithActivation(from, to, fmt.Sprintf(format, args...))
}

// AsyncResponseWithActivation add a async response to the sequence diagram.
func (d *Diagram) AsyncResponseWithActivation(from, to, message string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    %s-->>-%s: %s", from, to, message))
	return d
}

// AsyncResponsefWithActivation add a async response to the sequence diagram.
func (d *Diagram) AsyncResponsefWithActivation(from, to, format string, args ...any) *Diagram {
	return d.AsyncResponseWithActivation(from, to, fmt.Sprintf(format, args...))
}
