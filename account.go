package gother

type account struct {
	Client *GotherClient
	pri    string
}

func NewAccount(privateKey string) *account {
	return &account{
		pri:    privateKey,
		Client: Client,
	}
}

func (a *account) InjectPrivate(privateKey string) *account {
	a.pri = privateKey
	return a
}

func (c account) SignRaw(data ...[]byte) (str string, err error) {
	return SignRaw(c.pri, data...)
}
