package srcds

// Args are the command line options for the Source Dedicated Server executable
type Args struct {
	Insecure  bool `long:"insecure" description:"Will start the server without Valve Anti-Cheat." hidden:"true"`
	Bots      bool `long:"bots" description:"Used to enable bots" hidden:"true"`
	NoRestart bool `long:"norestart" description:"Won't attempt to restart failed servers."`
}

// AsSlice returns the command line options stored in a slice with individual values properly formatted for SRCDS
func (o Args) AsSlice() []string {
	var r []string

	if o.Insecure {
		r = append(r, "-insecure")
	}

	if !o.Bots {
		r = append(r, "-nobots")
	}

	if o.NoRestart {
		r = append(r, "-norestart")
	}

	return r
}
