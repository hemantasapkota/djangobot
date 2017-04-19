package djangobot

import "testing"

func TestLogin(t *testing.T) {

	// Before running this test, make sure to set your username and password below.

	// Disqus
	bot := With("https://disqus.com/profile/login/").
		AddHost("disqus.com").
		SetUsername("").
		SetPassword("")

	if bot.Error != nil {
		panic(bot.Error)
	}

	_, err := bot.LoadCookies().
		Set("next", "https://disqus.com/").
		X("csrfmiddlewaretoken", bot.Cookies["csrftoken"].Value).
		X("username", bot.Username).
		X("password", bot.Password).
		Login()

	if err != nil {
		panic(err)
	}

	 cookie, ok := bot.Cookies["sessionid"]
	 if !ok {
	 	t.Error("Authentication failed.")
	 }

	 t.Log(cookie.Value)

}
