package main

import (
	"GoPassCAPTCHA/request"
	"net/http/cookiejar"
)

func work() {
	client := request.New(UserAgent)

	cookieJar, _ := cookiejar.New(nil)
	client.SetCookieJar(cookieJar).
		SetHeader("User-Agent", UserAgent)

	// Step 1: Get initial cookies and CSRF token
	resp, err := client.Request().
		EnableTrace().
		Get(LoginURL)
	if err != nil {
		panic(err)
	}

	html := resp.String()
	csrftoken := extractCSRFToken(html)
	cookies := resp.Cookies()
	var session, route string
	for _, cookie := range cookies {
		switch cookie.Name {
		case "JSESSIONID":
			session = cookie.Value
		case "route":
			route = cookie.Value
		}
	}

	if csrftoken == "" || session == "" || route == "" {
		panic("Failed to fetch required tokens")
	}

	crackerCaptcha(&client)
}

func main() {
	for i := 0; i < 100; i++ {
		work()
	}
}
