package srcds

//Clients is a collection of type Client
type Clients []Client

func (m *Clients) ClientDropped(client Client) {
	i := m.clientIndex(client)

	if i >= 0 {
		l := len(*m)

		if l > 1 {
			*m = append((*m)[:i], (*m)[i+1:]...)
		} else if l == 1 {
			*m = Clients{}
		}
	}
}

func (m Clients) clientIndex(client Client) int {

	for i := range m {
		if ClientsAreEquivalent(&m[i], &client) {
			return i
		}
	}

	return -1
}

func (m *Clients) ClientJoined(client Client) {
	if !m.HasClient(client) {
		*m = append(*m, client)
	}
}

func (m Clients) HasClient(client Client) bool {
	return m.clientIndex(client) > -1
}
