package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type config struct {
	akamaiHost          string
	akamaiAccessToken   string
	akamaiClientToken   string
	akamaiClientSecret  string
	akamaiPurgeMethod   string
	akamaiPurgeNetwork  string
	akamaiPurgeHostname string

	githubCommitSHA  string
	githubToken      string
	githubOrg        string
	githubRepository string
	githubBranch     string
	siteURL          string
}

func main() {
	start := time.Now()

	configEnv := config{
		akamaiHost:          getEnv("AKAMAI_HOST", ""),
		akamaiAccessToken:   getEnv("AKAMAI_ACCESS_TOKEN", ""),
		akamaiClientToken:   getEnv("AKAMAI_CLIENT_TOKEN", ""),
		akamaiClientSecret:  getEnv("AKAMAI_CLIENT_SECRET", ""),
		akamaiPurgeMethod:   getEnv("AKAMAI_PURGE_METHOD", "invalidate"),
		akamaiPurgeNetwork:  getEnv("AKAMAI_PURGE_NETWORK", "production"),
		akamaiPurgeHostname: getEnv("AKAMAI_PURGE_HOSTNAME", "www.foo.com"),

		githubCommitSHA:  getEnv("GITHUB_COMMIT_SHA", ""),
		githubToken:      getEnv("GITHUB_TOKEN", ""),
		githubOrg:        getEnv("GITHUB_ORGANIZATION", ""),
		githubRepository: getEnv("GITHUB_REPOSITORY", ""),
		githubBranch:     getEnv("GITHUB_BRANCH", "master"),
	}
	configEdge := edgegrid.Config{
		Host:         configEnv.akamaiHost,
		ClientToken:  configEnv.akamaiClientToken,
		ClientSecret: configEnv.akamaiClientSecret,
		AccessToken:  configEnv.akamaiAccessToken,
		MaxBody:      1024,
		Debug:        true,
	}

	akamaiPurge(configEnv, configEdge, getFilesFromCommitGithub(configEnv))

	elapsed := time.Since(start)
	log.Println("Time elapsed: ", elapsed)

}

func getFilesFromCommitGithub(config config) []string {
	var listFiles []string
	var commitSHA string

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	opt := &github.CommitsListOptions{SHA: config.githubBranch}
	listCommits, _, err := client.Repositories.ListCommits(ctx, config.githubOrg, config.githubRepository, opt)
	if err != nil {
		log.Fatal(ctx, "Error in get list of commits: %v", err)
	}
	commitSHA = listCommits[0].GetSHA()
	if !isEmpty(config.githubCommitSHA) {
		commitSHA = config.githubCommitSHA
	}

	lastCommit, _, err := client.Repositories.GetCommit(ctx, config.githubOrg, config.githubRepository, commitSHA)
	if err != nil {
		log.Fatal(ctx, "Error in get commit: %v", err)
	}

	for _, files := range lastCommit.Files {
		if isPurgeableAsset(files.GetFilename()) {
			listFiles = append(listFiles, akamaiMakeUrl(config.akamaiPurgeHostname, files.GetFilename()))
		}
	}

	return listFiles
}

func akamaiMakeUrl(akamaiHostname string, fileName string) string {
	var outputUrl string
	outputUrl = "https://" + akamaiHostname + "/" + fileName
	return outputUrl
}

func isPurgeableAsset(Filename string) bool {
	var fileExtensionRegex = regexp.MustCompile(`(?i)\.(jpg|png|gif|js|webp|css|jpeg|html|scss)$`)
	var status bool
	status = false
	if fileExtensionRegex.MatchString(Filename) {
		status = true
	}
	return status
}

func isEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func akamaiPurge(configEnv config, configEdge edgegrid.Config, assets []string) {

	type akamaiRequest struct {
		Hostname string   `json:"hostname"`
		Objects  []string `json:"objects"`
	}

	akamaiBodyRequest := akamaiRequest{
		Hostname: configEnv.akamaiPurgeHostname,
		Objects:  assets,
	}

	bodyJson, err := json.Marshal(akamaiBodyRequest)
	if err != nil {
		fmt.Println("error:", err)
	}

	req, _ := client.NewRequest(configEdge, "POST", fmt.Sprintf("https://%s/ccu/v3/%s/url/%s", configEdge.Host, configEnv.akamaiPurgeMethod, configEnv.akamaiPurgeNetwork), bytes.NewBuffer(bodyJson))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := client.Do(configEdge, req)

	defer resp.Body.Close()
	byt, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byt))

}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
