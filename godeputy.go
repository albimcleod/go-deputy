package godeputy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	baseURL  = "https://once.deputy.com/"
	tokenURL = "my/oauth/access_token"

//	ApiEndpointURL = "https://api.reckon.com"
)

var (
	defaultSendTimeout = time.Second * 30
)

type Deputy struct {
	StoreCode    string
	ClientID     string
	ClientSecret string
	RedirectURI  string
	EndPoint     string
	Timeout      time.Duration
}

// NewClient will create a Reckon client with default values
func NewClient(code string, clientID string, clientSecret string, redirectURI string, endpoint string) *Deputy {
	return &Deputy{
		StoreCode:    code,
		Timeout:      defaultSendTimeout,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		EndPoint:     endpoint,
	}
}

// AccessToken will get a new access token
func (v *Deputy) AccessToken() (string, string, time.Time, string, error) {

	u, _ := url.ParseRequestURI(baseURL)
	u.Path = tokenURL
	urlStr := fmt.Sprintf("%v", u)

	request := fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s&client_id=%s&client_secret=%s&scope=longlife_refresh_token", v.StoreCode, v.RedirectURI, v.ClientID, v.ClientSecret)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer([]byte(request)))

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, _ := client.Do(r)

	rawResBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", "", time.Now(), "", fmt.Errorf("%v", string(rawResBody))
	}

	if res.StatusCode == 200 {
		resp := &TokenResponse{}
		if err := json.Unmarshal(rawResBody, resp); err != nil {
			return "", "", time.Now(), "", err
		}
		return resp.AccessToken, resp.RefreshToken, time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second), resp.EndPoint, nil
	}

	return "", "", time.Now(), "", fmt.Errorf("Failed to get access token: %s", res.Status)
}

// RefreshToken will get a new refresg token
func (v *Deputy) RefreshToken(refreshtoken string, token string) (string, string, time.Time, error) {

	u, _ := url.ParseRequestURI("https://" + v.EndPoint + "/")
	u.Path = "oauth/access_token"
	urlStr := fmt.Sprintf("%v", u)

	request := fmt.Sprintf("grant_type=refresh_token&refresh_token=%s&redirect_uri=%s&client_id=%s&client_secret=%s&scope=longlife_refresh_token", refreshtoken, v.RedirectURI, v.ClientID, v.ClientSecret)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer([]byte(request)))

	r.Header = http.Header(make(map[string][]string))
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Authorization", "OAuth "+token)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, _ := client.Do(r)

	rawResBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", "", time.Now(), fmt.Errorf("%v", string(rawResBody))
	}

	if res.StatusCode == 200 {
		resp := &TokenResponse{}
		if err := json.Unmarshal(rawResBody, resp); err != nil {
			return "", "", time.Now(), err
		}
		return resp.AccessToken, resp.RefreshToken, time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second), nil
	}

	return "", "", time.Now(), fmt.Errorf("Failed to get refresh token: %s", res.Status)
}

// GetCompanies will return the sites of the authenticated company
func (v *Deputy) GetCompanies(token string) (Companies, error) {
	client := &http.Client{}
	client.CheckRedirect = checkRedirectFunc

	u, _ := url.ParseRequestURI("https://" + v.EndPoint + "/")
	u.Path = "api/v1/resource/Company/"
	urlStr := fmt.Sprintf("%v", u)

	r, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	r.Header = http.Header(make(map[string][]string))
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Authorization", "Bearer "+token)

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	rawResBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 200 {
		var resp Companies

		err = json.Unmarshal(rawResBody, &resp)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
	return nil, fmt.Errorf("Failed to get Deputy Company %s", res.Status)

}

// GetRosters will
func (v *Deputy) GetRosters(token string, companyId string, startDate string, endDate string) (Rosters, error) {

	client := &http.Client{}
	client.CheckRedirect = checkRedirectFunc

	u, _ := url.ParseRequestURI("https://" + v.EndPoint + "/")
	u.Path = "api/v1/resource/Roster/QUERY"
	urlStr := fmt.Sprintf("%v", u)

	request := fmt.Sprintf(`{"join":["OperationalUnitObject", "MatchedByTimesheetObject"],"search":{"dateStart":{"field":"Date","type":"ge","data":"%v"},"dateEnd":{"field":"Date","type":"le","data":"%v"}}}`, startDate, endDate)

	//
	fmt.Println("Search", request)

	r, err := http.NewRequest("POST", urlStr, bytes.NewBuffer([]byte(request)))
	if err != nil {
		return nil, err
	}

	r.Header = http.Header(make(map[string][]string))
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Authorization", "Bearer "+token)

	res, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	rawResBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	//fmt.Println("rawResBody", string(rawResBody))

	if res.StatusCode == 200 {
		var resp Rosters

		err = json.Unmarshal(rawResBody, &resp)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
	return nil, fmt.Errorf("Failed to get Deputy rosters %s", res.Status)

}

// GetEmployees will
func (v *Deputy) GetEmployees(token string, companyId string) (Employees, error) {

	client := &http.Client{}
	client.CheckRedirect = checkRedirectFunc

	u, _ := url.ParseRequestURI("https://" + v.EndPoint + "/")
	u.Path = "api/v1/resource/Employee/QUERY"
	urlStr := fmt.Sprintf("%v", u)

	request := fmt.Sprintf(`{"search":{"company":{"field":"Company","type":"eq","data":"%s"}},"join":["ContactObject"]}`, companyId)

	fmt.Println(request)

	r, err := http.NewRequest("POST", urlStr, bytes.NewBuffer([]byte(request)))
	if err != nil {
		return nil, err
	}

	r.Header = http.Header(make(map[string][]string))
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Authorization", "Bearer "+token)

	res, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	rawResBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	//fmt.Println("rawResBody", string(rawResBody))

	if res.StatusCode == 200 {
		var resp Employees

		err = json.Unmarshal(rawResBody, &resp)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
	return nil, fmt.Errorf("Failed to get Deputy rosters %s", res.Status)

}

func checkRedirectFunc(req *http.Request, via []*http.Request) error {
	if req.Header.Get("Authorization") == "" {
		req.Header.Add("Authorization", via[0].Header.Get("Authorization"))
	}
	return nil
}
