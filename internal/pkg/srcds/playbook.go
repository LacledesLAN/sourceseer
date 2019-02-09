package srcds

// Playbook represents the rules to apply to a srcds instance
type Playbook func(SRCDS) error
