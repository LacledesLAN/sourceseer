package srcds

const (
	clientFlagAlpha ClientFlag = 1 << iota
	clientFlagBravo
	clientFlagCharlie
	clientFlagDelta
	clientFlagEcho
	clientFlagGolf
	clientFlagHotel
	clientFlagIndia
	clientFlagJuliett
	clientFlagKilo
	clientFlagLima
	clientFlagMike
	clientFlagNovember
	clientFlagOscar
	clientFlagPapa
	clientFlagQuedec
)

var allFlags = [16]ClientFlag{
	clientFlagAlpha, clientFlagBravo, clientFlagCharlie, clientFlagDelta, clientFlagEcho, clientFlagGolf, clientFlagHotel, clientFlagIndia,
	clientFlagJuliett, clientFlagKilo, clientFlagLima, clientFlagMike, clientFlagNovember, clientFlagOscar, clientFlagPapa, clientFlagQuedec,
}
