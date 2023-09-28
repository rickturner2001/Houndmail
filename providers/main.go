package providers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"rickturner2001/houndmail/utils"
	"strings"
)

type Provider struct {
	Name string `json:"name"`
}

func NewProvider(name string) *Provider {
	return &Provider{
		Name: name,
	}
}

func (p *Provider) HandleResponse(c *http.Client) {
	log.Printf("Handling provider response for %s", p.Name)
	switch p.Name {
	case "instagram":
		instagramInternalHandler(c)
	}
}

// TODO: clean this up
func instagramInternalHandler(c *http.Client) (bool, error) {
	reqUrl := "https://www.instagram.com/api/v1/web/accounts/web_create_ajax/attempt/"

	data := url.Values{}

	data.Set("email", "someemail@gmail.com")
	data.Set("first_name", "")
	data.Set("username", "")
	data.Set("opt_into_one_tap", "true")

	req, err := http.NewRequest(http.MethodPost, reqUrl, strings.NewReader(data.Encode()))

	jar, err := cookiejar.New(nil)
	if err != nil {
		err = utils.LogAndError(fmt.Sprintf("Could not create cookie jar: %s", err))
		return false, err
	}

	token := utils.GenToken(32, false)

	jarCookies := []*http.Cookie{{
		Name:  "csrftoken",
		Value: token,
	}}

	jar.SetCookies(&url.URL{
		Scheme: "https",
		Host:   "instagram.com",
	}, jarCookies)

	c.Jar = jar

	if err != nil {
		err = utils.LogAndError(fmt.Sprintf("Could not construct form data request: %s", err))
		return false, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Mobile Safari/537.36")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("X-Csrftoken", token)
	req.Header.Add("Referer", "https://www.instagram.com/accounts/emailsignup/")

	res, err := c.Do(req)
	if err != nil {
		err = utils.LogAndError(fmt.Sprintf("Could not send request to instagram: %s", err))
		return false, err
	}
}

func GetProviderByName(name string) (*Provider, error) {
	switch name {
	case "instagram":
		return NewProvider(name), nil
	default:
		err := utils.LogAndError(fmt.Sprintf("%s is not a valid provider", name))
		return nil, err
	}
}
