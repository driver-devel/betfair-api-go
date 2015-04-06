package betfair

var testAccount Account

func GetTestAccount() *Account {
	a := testAccount
	return &a
}
