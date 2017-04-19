package djangobot

import (
	"crypto/tls"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"strings"
)

type Bot struct {
	RequestUrl  string
	QueryParams map[string]string
	FormData    map[string]string

	Username string
	Password string

	Client *gorequest.SuperAgent
	Cookies map[string]*http.Cookie

	hosts  []string
	Error error
}

func With(requestUrl string) *Bot {
	return &Bot{
		QueryParams: make(map[string]string),
		FormData:    make(map[string]string),
		RequestUrl:  requestUrl,
	}
}

func (c *Bot) Set(key string, val string) *Bot {
	c.QueryParams[key] = val
	return c
}

func (c *Bot) X(key string, val string) *Bot {
	c.FormData[key] = val
	return c
}

func (c *Bot) AddHost(host string) *Bot {
	c.hosts = append(c.hosts, host)
	return c
}

func (c *Bot) SetUsername(username string) *Bot {
	c.Username = username
	return c
}

func (c *Bot) SetPassword(password string) *Bot {
	c.Password = password
	return c
}

func (c *Bot) referrer() (string, string) {
	return "Referer", c.RequestUrl
}

func (c *Bot) userAgent() (string, string) {
	return "User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36`
}

func (c *Bot) cookieKey(cookie *http.Cookie) string {
	if cookie == nil {
		return "nookie"
	}
	return strings.Split(cookie.Raw, "=")[0]
}

func (c *Bot) Cookie(key string) *http.Cookie {
	cookie, ok := c.Cookies[key]
	if !ok {
		return &http.Cookie{}
	}
	return cookie
}

func (c *Bot) appendCookies(resp *http.Response) {
	if c.Cookies == nil {
		c.Cookies = make(map[string]*http.Cookie)
	}

	for _, cookie := range resp.Cookies() {
		// Get key
		c.Cookies[c.cookieKey(cookie)] = cookie
		// Append cookies to the client ( Should've been appended after the requests automatically, but does not )
		c.Client = c.Client.AddCookie(cookie)
	}
}

func (c *Bot) LoadCookies() *Bot {
	var body []byte
	var errs []error
	var resp *http.Response

	var tlsConfig *tls.Config

	c.Client = gorequest.New().SetCurlCommand(false).SetDebug(false)

	// Get certs for SSL connection
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(c.hosts...),
	}

	tlsConfig = &tls.Config{GetCertificate: m.GetCertificate}

	c.Client = c.Client.TLSClientConfig(tlsConfig).Get(c.RequestUrl).Set(c.userAgent())

	resp, body, errs = c.Client.EndBytes()

	if errs != nil {
		c.Error = errs[0]
	}

	// Add cookies from the response to the client
	c.appendCookies(resp)

	_ = resp
	_ = body
	_ = errs

	return c
}

func (c *Bot) Login() (*gorequest.SuperAgent, error) {
	var body []byte
	var errs []error
	var resp *http.Response

	var rPolicy = func(req gorequest.Request, via []gorequest.Request) error {
		return http.ErrUseLastResponse
	}

	// Set query parameters
	url := c.RequestUrl
	if len(c.QueryParams) > 0 {
		url += "?" + formEncode(c.QueryParams)
	}

	if len(c.FormData) > 0 {
		c.Client = c.Client.
			Post(url).
			Set(c.userAgent()).
			Set(c.referrer()).
			Send(formEncode(c.FormData)).
			RedirectPolicy(rPolicy)
	}

	resp, body, errs = c.Client.EndBytes()
	if errs != nil {
		return nil, errs[0]
	}

	c.appendCookies(resp)

	_ = resp
	_ = body
	_ = errs

	return c.Client, nil
}

func formEncode(data map[string]string) string {
	values := make([]string, 0)
	for key, value := range data {
		values = append(values, fmt.Sprintf(`%s=%s`, key, value))
	}
	return strings.Join(values, "&")
}
