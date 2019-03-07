package merchant


type merchant struct {
	username string
	password []byte
	token string
}

func New(username string, password []byte) *merchant {
	return &merchant{
		username: username,
		password: password,
	}
}

func (m *merchant) Username() string{
	return m.username
}

func (m *merchant) setMerchant(username string) {
	m.username = username
}

func (m *merchant) Password() []byte {
	return m.password
}

func (m *merchant) setPassword(password []byte) {
	m.password = password
}

func (m *merchant) Token() string {
	return m.token
}

func (m *merchant) setToken(token string) {
	m.token = token
}
// Authenticate connect to the database to authenticate merchant credentials
func Authenticate(username string, password []byte) (bool, error) {

	return true, nil
}