package poeclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ByChanderZap/exile-tracker/utils"
	"github.com/rs/zerolog"
)

const (
	BaseURL               = "https://www.pathofexile.com"
	CharactersEndpoint    = "/character-window/get-characters"
	PassiveSkillsEndpoint = "/character-window/get-passive-skills"
	ItemsEndpoint         = "/character-window/get-items"

	UserAgent = "Oath exile-tracker/0.0.1 (contact: neryt.alexander@gmail.com)"
)

type POEClient struct {
	httpClient *http.Client
	baseURL    string
	userAgent  string
	log        zerolog.Logger
}

type RateLimitInfo struct {
	Policy     string
	Rules      []string
	Limits     map[string]RateLimit
	States     map[string]RateLimitState
	RetryAfter int
}

type RateLimit struct {
	Max         int
	Period      int
	Restriction int
}

type RateLimitState struct {
	CurrentHits int
	Period      int
	Restricted  int
}

type POEError struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"Message"`
	} `json:"error"`
}

func NewPoeClient(timeout time.Duration) *POEClient {
	return &POEClient{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		baseURL:   BaseURL,
		userAgent: UserAgent,
		log:       utils.ChildLogger("poe-client"),
	}
}

func (pc *POEClient) makeRequest(endpoint string, params map[string]string) (*http.Response, error) {
	//parse url
	u, err := url.Parse(pc.baseURL + endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("couldnt create request %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", pc.userAgent)

	res, err := pc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	rateLimitInfo := pc.parseRateLimitHeaders(res)
	if rateLimitInfo.RetryAfter > 0 {
		pc.log.Warn().
			Int("retry_after", rateLimitInfo.RetryAfter).
			Msg("Rate limited, should retry after seconds")
	}

	if res.StatusCode >= 400 {
		defer res.Body.Close()

		if res.StatusCode == 429 {
			return nil, fmt.Errorf("rate limit reached, retry after %d", rateLimitInfo.RetryAfter)
		}

		body, _ := io.ReadAll(res.Body)
		var poeError POEError
		err = json.Unmarshal(body, &poeError)
		if err != nil {
			return nil, fmt.Errorf("couldnt unmarshall error response %w", err)
		}

	}

	return res, nil
}

func (pc *POEClient) parseRateLimitHeaders(res *http.Response) RateLimitInfo {
	info := RateLimitInfo{
		Limits: make(map[string]RateLimit),
		States: make(map[string]RateLimitState),
	}

	info.Policy = res.Header.Get("X-Rate-Limit-Policy")

	if rules := res.Header.Get("X-Rate-Limit-Rules"); rules != "" {
		sRules := strings.Split(rules, ",")
		info.Rules = sRules
	}

	if retryAfter := res.Header.Get("Retry-After"); retryAfter != "" {
		if retryInt, err := strconv.Atoi(retryAfter); err == nil {
			info.RetryAfter = retryInt
		}
	}

	for _, rule := range info.Rules {
		if limitHeader := res.Header.Get("X-Rate-Limit-" + rule); limitHeader != "" {
			info.Limits[rule] = RateLimit{}
		}

		if stateHeader := res.Header.Get("X-Rate-Limit-" + rule + "-State"); stateHeader != "" {
			info.States[rule] = RateLimitState{}
		}
	}

	return info
}

func (pc *POEClient) GetCharacters(acName string, realm string) (*http.Response, error) {
	params := map[string]string{
		"accountName": acName,
		"realm":       realm,
	}

	res, err := pc.makeRequest(CharactersEndpoint, params)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (pc *POEClient) GetPassiveSkills(acName string, character string, realm string) (*http.Response, error) {
	params := map[string]string{
		"accountName": acName,
		"character":   character,
		"realm":       realm,
	}
	res, err := pc.makeRequest(PassiveSkillsEndpoint, params)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (pc *POEClient) GetItems(acName string, character string, realm string) (*http.Response, error) {
	params := map[string]string{
		"accountName": acName,
		"character":   character,
		"realm":       realm,
	}
	res, err := pc.makeRequest(ItemsEndpoint, params)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (pc *POEClient) GetItemsJson(acName string, character string, realm string) ([]byte, error) {
	res, err := pc.GetItems(acName, character, realm)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error while decoding get items response %w", err)
	}
	return data, nil
}

func (pc *POEClient) GetPassiveSkillsJson(acName string, character string, realm string) ([]byte, error) {
	res, err := pc.GetPassiveSkills(acName, character, realm)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error while decoding get passive skills response %w", err)
	}
	return data, nil
}

func (pc *POEClient) GetCharactersJson(acName string, realm string) ([]byte, error) {
	res, err := pc.GetCharacters(acName, realm)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error while decoding get characters response %w", err)
	}
	return data, nil
}
