# djangobot
[Curl](https://curl.haxx.se/) for [Django](https://www.djangoproject.com/). Make authenticated requests to a Django server.

# How does it work ?

Django authentication relies on two cookies: **csrftoken** and **sessionid**. Once you accquire these cookies, you can make authenticated requests just like the browser does.

Getting the **csrftoken** is easy. Just make a request to a page and the server sends back the cookie. 

**Sessionid**, however is tricky because most production servers configure it as a [secure HTTP only](https://docs.djangoproject.com/en/1.11/ref/settings/#std:setting-SESSION_COOKIE_SECURE) cookie. It is only sent if authentication is made securely ( via. HTTPS )

To make a secure connection we need SSL/TLS certificates. GO has a package called [autocert](https://godoc.org/golang.org/x/crypto/acme/autocert) which lets us accquire these certificates. Autocert provides automatic access to certificates from [Let's Encrypt](https://letsencrypt.org/) and any other ACME-based CA.

With certs in place, all we need is the authentication details and we're good to go.

# What are the potential use cases for this library ?

* API testing
* Web Scraping
* Automation
* Bot frameworks
* Mobile apps

See an example usage below.

# Installation

* Add ```github.com/hemantasapkota/djangobot``` as an import to your project.
* Execute ```go get github.com/hemantasapkota/djangobot```

# Usage
In this example, we'll authenticate with [Disqus](https://disqus.com/) which is built on top of Django. Let's inspect the parameters that get sent to the login endpoint.

![](disqus.png)

The query parameter is **next** and the form data items are **csrfmiddlewaretoken**, **username**, and **password**.

We'll do the same. But before being able to call the login endpoint we'll need to accquire the CSRF token. Let's go get it.

```go
bot := djangobot.With("https://disqus.com/profile/login/").
		 ForHost("disqus.com").
		 SetUsername("<<username>>").
		 SetPassword("<<password>>").
         	 LoadCookies()

if bot.Error != nil {
	panic(bot.Error)
}
```

Next, let's authenticate with the server. Django expects the csrf token to be sent as the **csrfmiddlewaretoken** form data. **Set()** sets the query parameters and **X()** sets the form data.

```go
client, err := bot.Set("next", "https://disqus.com/").
		   X("csrfmiddlewaretoken", bot.Cookie("csrftoken").Value).
		   X("username", bot.Username).
		   X("password", bot.Password).
		   Login()

if err != nil {
	panic(err)
}

sessionid := bot.Cookie("sessionid").Value
if sessionid == "" {
    panic("Authentication failed.")
}

```

Successful authentication creates the **sessionid** cookie and returns an http [client](https://github.com/parnurzeal/gorequest) object.

From this point on, the HTTP client can be used to make requests. It's important to note that all subsequent requests should have at least these headers: **User-Agent**, **Referrer**, **X-CSRFToken**, and **X-Requested-With**.

The ```bot.Requester()``` method is available to prepare requests with pre-filled headers. Example below.

### Changing your Discus account password

Let's put this library to use by changing our account's password.

```go
data := map[string]string{
	"email":        "<<your email address>>",
	"old_password": "<<your old password>>",
	"password":     "<<new password>>",
	"username":     "<<username>>",
}

_, body, _ := bot.Requester("PUT", "https://disqus.com/users/self/account/").
	      Client.
	      Send(data).
              End()

fmt.Println(body)

```

Please refer to the test file for more details.
