package djangobot

import (
	"testing"
)

func TestLogin(t *testing.T) {

	// Before running this test, make sure to set your username and password below.

	username := ""
	password := ""

	if username == "" && password == "" {
		t.Log("Cannot execute test. Please supply username and password.")
		return
	}

	// Disqus
	bot := With("https://disqus.com/profile/login/").
		ForHost("disqus.com").
		SetUsername(username).
		SetPassword(password)

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
		return
	}

	t.Log("Authenticated with session id: ", sessionid)

	// Add details to the form data
	data := map[string]string{}

	//data := map[string]string{
	//	"email":        "",
	//	"old_password": "",
	//	"password":     "",
	//	"username":     "",
	//}

	if len(data) == 0 {
		t.Log("Cannot change password. Form data empty.")
		return
	}

	_, body, _ := bot.Requester("PUT", "https://disqus.com/users/self/account/").
		Client.
		Send(data).
		End()

	t.Log(body)

}
