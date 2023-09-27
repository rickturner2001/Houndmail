package providers

import (
	"fmt"
	"log"
	"net/http"
	"rickturner2001/houndmail/requests"
	"rickturner2001/houndmail/utils"
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

func instagramInternalHandler(c *http.Client) {
	log.Printf("Internal instagram handler")
	requests.ExtractCookie(c, "csrftoken", "https://instagram.com")
	// reqUrl = "https://www.instagram.com/api/v1/web/accounts/web_create_ajax/attempt/"

	// req, err  != http.NewformdataNewFormDataRequest(reqUrl, "")
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
