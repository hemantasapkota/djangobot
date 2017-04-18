# djangobot
Convert your django app into a headless web client

# Installation

```
go get github.com/hemantasapkota/djangobot
```

# Usage
In this example, we'll authenticate with Disqus ( https://disqus.com/ ) which is built on top of Django.

## Step 1: Load Cookies

This step loads the CSRF cookie.

```
	bot := djangobot.With("https://disqus.com/profile/login/").
		               AddHost("disqus.com").
		               SetUsername("<<username>>").
		               SetPassword("<<password>>").
                   LoadCookies()
                   
```

## Step 2: Authenticate with the server

Django expects the csrf token to be sent as the **csrfmiddlewaretoken** HTTP header.

```
client, err := bot.Set("next", "https://disqus.com/").
		                 X("csrfmiddlewaretoken", bot.Cookies["csrftoken"].Value).
		                 X("username", bot.Username).
		                 X("password", bot.Password).
		                 Login()

if err != nil {
	panic(err)
}
```

### Step 3: Check authentication result

Successful authentication creates the **sessionid** cookie.

```
cookie, ok := bot.Cookies["sessionid"]
if !ok {
	panic("Authentication failed.")
}
  
 fmt.Println(cookie.Value)
  
```

# Use cases

* Developer testing
* Web Scraping
* API
