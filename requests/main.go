package requests

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"rickturner2001/houndmail/utils"
)

func NewFormDataRequest(reqUrl, boundary string, headers, fields map[string]string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if boundary != "" {
		err := writer.SetBoundary(boundary)
		if err != nil {
			log.Printf("Error setting boundary: %s", err)
			return nil, err

		}
	}

	for key, value := range fields {
		err := writer.WriteField(key, value)
		if err != nil {

			log.Printf("Error writing field %s: %s", key, err)
			return nil, err
		}
	}

	err := writer.Close()
	if err != nil {
		log.Printf("Error closing writer: %s", err)
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, reqUrl, body)
	if err != nil {
		err = utils.LogAndError(fmt.Sprintf("Request error: %s", err))
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

func ExtractCookie(c *http.Client, cookieKey string, endpoint string) (string, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		err = utils.LogAndError(fmt.Sprintf("Could not create cookie jar: %s", err))
		return "", err
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

	log.Printf("Retreiving cookie with key: %s at %s", cookieKey, endpoint)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		err = utils.LogAndError(fmt.Sprintf("Could not initilize request error: %s", err))
		return "", err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Mobile Safari/537.36")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("X-Csrftoken", token)

	if jarCookies[0] != nil {
		req.AddCookie(jarCookies[0])
	}

	log.Printf("Finishde Setting up headers: %s", req.Header.Get("X-Csrftoken"))

	res, err := c.Do(req)
	if err != nil {
		err = utils.LogAndError(fmt.Sprintf("Could not send request: %s", err))
		return "", err
	}

	cookies := res.Cookies()
	if len(cookies) == 0 {
		fmt.Println("No cookies were found")
	} else {
		for _, cookie := range cookies {
			fmt.Printf("%s: %s", cookie.Name, cookie.Value)
		}
	}

	return "", nil
}

// func (ac *APIClient) ContactProvider(p *providers.Provider) error {
// 	var r *http.Request
//
// 	if p.IsFormData {
// 		r, err := http.NewRequest(http.MethodGet, p.Endpoint, nil)
// 		if err != nil {
// 			log.Printf("Could not initialize a request to %s", p.Endpoint)
// 			return fmt.Errorf("Could not initialize a request to %s", p.Endpoint)
// 		}
// 	}else{
// 		r, err := http.NewRequest(http.MethodPost, p.Endpoint, )
// 	}
//
// 	for key, val := range p.Headers {
// 		r.Header.Set(key, val)
// 	}
//
// 	res, err := ac.Client.Do(r)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
