package srcds

type LogCategory uint8

const (
	SourceSeer LogCategory = 1 << iota
	SRCDSCvar
	SRCDSError
	SRCDSOther
	SRCDSLog
)

func (c LogCategory) Clear(flag LogCategory) LogCategory  { return c &^ flag }
func (c LogCategory) Has(flag LogCategory) bool           { return c&flag != 0 }
func (c LogCategory) Set(flag LogCategory) LogCategory    { return c | flag }
func (c LogCategory) Toggle(flag LogCategory) LogCategory { return c ^ flag }

func Log(c LogCategory, msg string) {
	if len(msg) == 0 {
		return
	}

	switch c {
	case SourceSeer:
		//fmt.Println("[SOURCESEER ]", msg)
	case SRCDSCvar:
		//fmt.Println("[SRCDS CVAR ]", msg)
	case SRCDSError:
		//fmt.Println("[SRCDS ERR  ]", msg)
	case SRCDSLog:
		//fmt.Println("[SRCDS LOG  ]", msg)
	case SRCDSOther:
		//fmt.Println("[SRCDS OTHER]", msg)
	}
}
