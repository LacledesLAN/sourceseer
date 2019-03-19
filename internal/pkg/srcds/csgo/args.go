package csgo

import (
	"strconv"

	"github.com/lacledeslan/sourceseer/internal/pkg/srcds"
)

//Args contains the command line options for the Counter-Strike: Global Offensive Server executable
type Args struct {
	srcds.Args
	GameMode         int    `long:"game_mode" description:"Used with game_type to change the csgo game you are playing (e.g. arms race, competitive, etc)" hidden:"true" default:"1"`
	GameType         int    `long:"game_type" description:"Used with game_mode to change the csgo game you are playing (e.g. arms race, competitive, etc)" hidden:"true" default:"0"`
	Hostname         string `long:"hostname" description:"Set the server's host name" hidden:"true"`
	TeamName1        string `long:"mp_teamname_1" description:"The the name for team 1 (starts as CT)"`
	TeamName2        string `long:"mp_teamname_2" description:"The the name for team 2 (starts as Terrorist)"`
	LANMode          int    `long:"sv_lan" description:"If set to 1, server is only available in Local Area Network (LAN)." choice:"0" choice:"1" default:"1" hidden:"true"`
	RConPassword     string `long:"rcon_password" description:"This command will authenticate you for rcon with the specified password." default:"0"`
	Map              string `long:"map" description:"Set the starting map" default:"de_cbble"`
	Password         string `long:"sv_password" description:"Set a password required to connect to the server. Set to 0 to disable" default:"0"`
	TickRate         int    `long:"tickrate" description:"" choice:"64" choice:"128" default:"128"`
	TVName           string `long:"tv_name" description:"GOTV host name"`
	TVPassword       string `long:"tv_password" description:"GOTV password for all clients."`
	TVRelayPassword  string `long:"tv_relaypassword" description:"GOTV password for relay proxies" hidden:"true"`
	UseRemoteConsole bool   `long:"userrcon" description:"Enables Remove Console" default:"true"`
}

//AsSlice returns the command line options stored in a slice with individual values propertly formatted for SRCDS
func (o Args) AsSlice() []string {
	r := append(o.Args.AsSlice(), "-game csgo", "+game_mode "+strconv.Itoa(o.GameMode), "+game_type "+strconv.Itoa(o.GameType), "+map "+o.Map)

	if o.Hostname != "" {
		r = append(r, `+hostname "`+o.Hostname+`"`)
	}

	if o.TeamName1 != "" {
		r = append(r, `+mp_teamname_1 "`+o.TeamName1+`"`)
	}

	if o.TeamName2 != "" {
		r = append(r, `+mp_teamname_2 "`+o.TeamName2+`"`)
	}

	if o.LANMode != 0 {
		r = append(r, "+sv_lan "+strconv.Itoa(o.LANMode))
	}

	if len(o.RConPassword) > 0 && o.RConPassword == "0" {
		r = append(r, `+rcon_password "`+o.RConPassword+`"`)
	}

	if o.Password != "" {
		r = append(r, `+sv_password "`+o.Password+`"`)
	}

	r = append(r, "-tickrate "+strconv.Itoa(o.TickRate))

	if o.TVName != "" {
		r = append(r, `+tv_name "`+o.TVName+`"`)
	}

	if o.TVPassword != "" {
		r = append(r, `+tv_password "`+o.TVPassword+`"`)
	}

	if o.TVRelayPassword != "" {
		r = append(r, `+tv_relaypassword "`+o.TVRelayPassword+`"`)
	}

	if o.TVRelayPassword != "" {
		r = append(r, ` +tv_relaypassword "`+o.TVRelayPassword+`"`)
	}

	if o.UseRemoteConsole {
		r = append(r, "-usercon")
	}

	return r
}
