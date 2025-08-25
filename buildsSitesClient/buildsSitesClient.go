package buildsSitesClient

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// SiteInfo holds info for each build site
type SiteInfo struct {
	Label      string
	ID         string
	CodeOut    string
	PostURL    string
	PostFields string
	LinkURL    string
}

// SitesUrl provides static references to each supported site
var SitesUrl = struct {
	POBBin   SiteInfo
	PoeNinja SiteInfo
	Poedb    SiteInfo
}{
	POBBin: SiteInfo{
		Label:      "pobb.in",
		ID:         "POBBin",
		CodeOut:    "https://pobb.in/",
		PostURL:    "https://pobb.in/pob/",
		PostFields: "",
		LinkURL:    "pobb.in/",
	},
	PoeNinja: SiteInfo{
		Label:      "PoeNinja",
		ID:         "PoeNinja",
		CodeOut:    "",
		PostURL:    "https://poe.ninja/pob/api/api_post.php",
		PostFields: "api_paste_code=",
		LinkURL:    "poe.ninja/pob/",
	},
	Poedb: SiteInfo{
		Label:      "poedb.tw",
		ID:         "PoEDB",
		CodeOut:    "",
		PostURL:    "https://poedb.tw/pob/api/gen",
		PostFields: "",
		LinkURL:    "poedb.tw/pob/",
	},
}

func UploadBuild(buildCode string, site SiteInfo) (string, error) {
	if site.PostURL == "" {
		return "", fmt.Errorf("no post URL for site %s", site.Label)
	}
	postBody := site.PostFields + buildCode

	req, err := http.NewRequest("POST", site.PostURL, bytes.NewBufferString(postBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "exile-tracker/0.0.1 (contact: neryt.alexander@gmail.com)")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		return fmt.Sprintf("%s%s", site.CodeOut, string(body)), nil
	}

	return "", fmt.Errorf("upload failed: %s", string(body))
}
