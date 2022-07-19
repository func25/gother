package gother

type account struct {
	*GotherClient
	pri string
}

func NewAccount(privateKey string) *account {
	return &account{
		pri:          privateKey,
		GotherClient: Client,
	}
}

func (a *account) InjectPrivate(privateKey string) *account {
	a.pri = privateKey
	return a
}

func (c account) Keccak256Sign(data ...[]byte) (str string, err error) {
	return Keccak256Sign(c.pri, data...)
}

func (c account) Sign(data []byte) (str string, err error) {
	return Sign(c.pri, data)
}
