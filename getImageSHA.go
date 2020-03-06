package main

import (
    "errors"
    "fmt"
    "net/http"
    "github.com/BurntSushi/toml"
    "os"
    "flag"
)

const ACCEPT_HEADER = "application/vnd.docker.distribution.manifest.v2+json"
const CREDENTIALS_FILE = ".credentials"

type Registry struct {
    Host       string `toml:"nexus_host"`
    Username   string `toml:"nexus_username"`
    Password   string `toml:"nexus_password"`
    Repository string `toml:"nexus_repository"`
}

func NewRegistry() (Registry, error) {
	r := Registry{}
	if _, err := os.Stat(CREDENTIALS_FILE); os.IsNotExist(err) {
		return r, errors.New(fmt.Sprintf("%s file not found\n", CREDENTIALS_FILE))
	} else if err != nil {
		return r, err
	}

	if _, err := toml.DecodeFile(CREDENTIALS_FILE, &r); err != nil {
		return r, err
	}
	return r, nil
}

func getImageSHA(image string, tag string, host string, repository string, username string, password string) (string, error) {
	client := &http.Client{}

	url := fmt.Sprintf("%s/repository/%s/v2/%s/manifests/%s", host, repository, image, tag)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(username, password)
	req.Header.Add("Accept", ACCEPT_HEADER)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("HTTP Code: %d", resp.StatusCode))
	}

	return resp.Header.Get("docker-content-digest"), nil
}

func main() {

image := flag.String("image","wrong","image name")
tag := flag.String("tag","nil","tag name")
flag.Parse()

output,err := NewRegistry()
if err == nil {
    sha,sha_err := getImageSHA(*image,*tag,output.Host,output.Repository,output.Username,output.Password)
        if sha_err == nil {
            fmt.Print(sha)
    } else {fmt.Println(sha_err)}
}

}
