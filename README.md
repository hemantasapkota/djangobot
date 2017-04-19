# djangobot
CURL for Django. Make authenticated requests to a Django server.

# How does it work ?

Django authentication relies on two cookies: **csrfmiddlewaretoken** and **sessionid**. Once you accquire these cookies, you can make authenticated requests just like the browser does.

Getting the **csrfmiddlewaretoken** is easy. Just make a request to a page and the server sends back the cookie. **sessionid**, however is tricky because it's a secure HTTP only cookie. It is only sent if authentication is made securely ( via. HTTPS )

To make a secure connection we need SSL/TLS certificates. GO has a package called [autocert](https://godoc.org/golang.org/x/crypto/acme/autocert) which lets us accquire these certificates.

With certs in place, all we need is the authentication details and we're good to go.

# 

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
