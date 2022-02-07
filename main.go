package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func init() {
	for _, x := range []*env{
		&host,
		&token,
		&channel,
		&maxRetry,

		&repoFullname,
		&remoteURL,
		&commitSha,
		&commitBranch,
		&commitMessage,
		&commitLink,
		&commitAuthorName,
		&buildEvent,
		&buildStatus,
		&buildNumber,
		&stageName,
	} {
		x.Value = os.Getenv(x.EnvVar)
	}
}

// report things: repo(name),time,commit(author/hash/msg),branch,build result(event/)
func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	c := &http.Client{
		Transport: tr,
	}

	maxRetryCount, err := strconv.Atoi(maxRetry.Value)
	if err != nil {
		maxRetryCount = 0
	}
	if maxRetryCount == 0 {
		maxRetryCount = 3
	}
	emoji := " ✅"
	if strings.Contains(buildStatus.Value, "failure") {
		emoji = " ❌"
	}
	msg := fmt.Sprintf(`CI build: %s (branch: %s) (stage: %s)
commit: [click to see diff](%s) 
author: %s
commit message: %s
event: %s:%s
%s%s`,
		repoFullname, commitBranch, stageName,
		commitLink, commitAuthorName, commitMessage,
		buildEvent, buildNumber,
		buildStatus, emoji)
	osExitCode := 1

	payload, err := json.Marshal(map[string]string{
		"text":    msg,
		"channel": channel.Value,
	})
	if err != nil {
		panic(err)
	}

	url := fmt.Sprintf("%s/hooks/%s", host, token)

	for i := 0; i < maxRetryCount; i++ {
		fmt.Printf("Posting to %s: %s\n", url, string(payload))
		r, err := c.Post(url, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			log.Printf("create post failed (i:%d), %v", i, err.Error())
			time.Sleep(time.Duration(i) * time.Second)
			continue
		}
		if r.StatusCode == http.StatusOK {
			osExitCode = 0
			log.Println("post success")
			break
		} else {
			log.Printf("status code: %d", r.StatusCode)
			continue
		}
	}
	os.Exit(osExitCode)
}

var (
	//mattermost env
	host = env{
		EnvVar: "host",
	}
	token = env{
		EnvVar: "token",
	}
	channel = env{
		EnvVar: "channel",
	}
	maxRetry = env{
		Usage:  "max retry when post message failed",
		EnvVar: "maxRetry",
	}

	//drone env
	repoFullname = env{
		Name:   "repo.fullname",
		Usage:  "repository full name",
		EnvVar: "DRONE_REPO",
	}
	remoteURL = env{
		Name:   "remote.url",
		Usage:  "git remote url",
		EnvVar: "DRONE_REMOTE_URL",
	}
	commitSha = env{
		Name:   "commit.sha",
		Usage:  "git commit sha",
		EnvVar: "DRONE_COMMIT_SHA",
	}
	commitBranch = env{
		Name:   "commit.branch",
		Value:  "master",
		Usage:  "git commit branch",
		EnvVar: "DRONE_COMMIT_BRANCH",
	}
	commitMessage = env{
		Name:   "commit.message",
		Usage:  "git commit message",
		EnvVar: "DRONE_COMMIT_MESSAGE",
	}
	commitLink = env{
		Name:   "commit.link",
		Usage:  "git commit link",
		EnvVar: "DRONE_COMMIT_LINK",
	}
	commitAuthorName = env{
		Name:   "commit.author.name",
		Usage:  "git author name",
		EnvVar: "DRONE_COMMIT_AUTHOR",
	}
	buildEvent = env{
		Name:   "build.event",
		Value:  "push",
		Usage:  "build event",
		EnvVar: "DRONE_BUILD_EVENT",
	}
	buildStatus = env{
		Name:   "build.status",
		Usage:  "build status",
		Value:  "success",
		EnvVar: "DRONE_BUILD_STATUS",
	}
	buildNumber = env{
		Name:   "build.number",
		Usage:  "build number",
		EnvVar: "DRONE_BUILD_NUMBER",
	}
	stageName = env{
		Name:   "stage.name",
		Usage:  "drone stage name",
		EnvVar: "DRONE_STAGE_NAME",
	}
)
