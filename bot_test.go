package djangobot

import "testing"

func TestLogin(t *testing.T) {

	// Before running this test, make sure to set your username and password below.

	// Disqus
	bot := With("https://disqus.com/profile/login/").
		ForHost("disqus.com").
		SetUsername("").
		SetPassword("")

	if bot.Error != nil {
		panic(bot.Error)
	}

	_, err := bot.LoadCookies().
		Set("next", "https://disqus.com/"). // Set the next parameter
		X("csrfmiddlewaretoken", bot.Cookie("csrftoken").Value).
		X("username", bot.Username).
		X("password", bot.Password).
		Login()

	if err != nil {
		panic(err)
	}

	 sessionid := bot.Cookie("sessionid").Value
	 if sessionid == "" {
	 	t.Error("Authentication failed.")
	 }

	 t.Log(sessionid)

}
