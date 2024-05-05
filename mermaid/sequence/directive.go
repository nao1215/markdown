package sequence

import (
	"fmt"
	"strings"
)

// AutoNumber add auto number to the sequence diagram.
func (d *Diagram) AutoNumber() *Diagram {
	d.body = append(d.body, "    autonumber")
	return d
}

// BoxStart add a box to the sequence diagram.
func (d *Diagram) BoxStart(participant []string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    box %s", strings.Join(participant, " & ")))
	return d
}

// BoxEnd add a box to the sequence diagram.
func (d *Diagram) BoxEnd() *Diagram {
	d.body = append(d.body, "    end")
	return d
}

// Participant add a participant to the sequence diagram.
func (d *Diagram) Participant(participant string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    participant %s", participant))
	return d
}

// Actor add a participant to the sequence diagram.
func (d *Diagram) Actor(actor string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    actor %s", actor))
	return d
}

// CreateParticipant add a participant to the sequence diagram.
func (d *Diagram) CreateParticipant(participant string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    create participant %s", participant))
	return d
}

// DestroyParticipant add a participant to the sequence diagram.
func (d *Diagram) DestroyParticipant(participant string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    destroy %s", participant))
	return d
}

// CreateActor add a participant to the sequence diagram.
func (d *Diagram) CreateActor(actor string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    create actor %s", actor))
	return d
}

// DestroyActor add a participant to the sequence diagram.
func (d *Diagram) DestroyActor(actor string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    destroy %s", actor))
	return d
}
