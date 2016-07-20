package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

// Configuration for this application
type Configuration struct {
	User   string `json:"user"`
	Key    string `json:"key"`
	ApiUrl string `json:"apiUrl"`
}

type GistContent struct {
	Content string `json:"content"`
}

type Gist struct {
	Description string                 `json:"description"`
	Public      bool                   `json:"public"`
	Files       map[string]GistContent `json:"files"`
}

// Can add more data here if necessary
type GistCreationResponse struct {
	Id      string `json:"id"`
	HtmlUrl string `json:"html_url"`
}

// Read stdin until EOF
func readInput() string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(os.Stdin)
	return buf.String()
}

func createGist(configuration Configuration, filename string, content string) (string, error) {
	gist := Gist{
		// Description: "Gist",
		Public: false,
		Files: map[string]GistContent{
			filename: GistContent{
				Content: content,
			},
		},
	}
	gistBytes, jsonMarhsalError := json.Marshal(gist)
	if jsonMarhsalError != nil {
		return "", jsonMarhsalError
	}

	request, requestError := http.NewRequest(
		"POST",
		configuration.ApiUrl + "/gists",
		bytes.NewBuffer(gistBytes))
	if requestError != nil {
		return "", requestError
	}
	request.SetBasicAuth(configuration.User, configuration.Key)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	} else {
		defer response.Body.Close()
		var structuredResponse GistCreationResponse
		decoderErr := json.NewDecoder(response.Body).Decode(&structuredResponse)
		if decoderErr != nil {
			return "", decoderErr
		}
		return structuredResponse.HtmlUrl, nil
	}
}

func readConfigurationFile() (Configuration, error) {
	var configuration Configuration

	defaultApiUrl := "https://api.github.com"
	configurationPath := os.Getenv("HOME") + "/.gistrc"
	if _, err := os.Stat(configurationPath); os.IsNotExist(err) {
		return configuration, errors.New(
			"Unable to find ~/.gistrc, please create one following the form:\n" +
				"{\n\t\"key\": \"yourApiKey\",\n" +
				"\t\"user\": \"yourUsername\", # Default: current user\n" +
				"\t\"apiUrl\": \"github instance apiUrl\" # Default: " + defaultApiUrl + "\n}")
	}

	// Read the configuration file
	file, err := ioutil.ReadFile(configurationPath)
	if err != nil {
		return configuration, err
	}
	err = json.Unmarshal(file, &configuration)
	if err != nil {
		return configuration, err
	}

	// Set default values
	if configuration.User == "" {
		configuration.User = os.Getenv("USER")
	}
	if configuration.ApiUrl == "" {
		configuration.ApiUrl = defaultApiUrl
	}
	if configuration.Key == "" {
		return configuration, errors.New("Github API key must be configured in ~/.gistrc")
	}

	return configuration, nil
}

func main() {
	file := flag.String("f", "", "File to upload to gist")
	flag.Parse();
	configuration, err := readConfigurationFile()
	if err != nil {
		log.Fatal(err)
		return
	}
	var filename string
	var content string
	if *file == "" {
		filename = "file"
		content = readInput()
	} else {
		_, filename = path.Split(*file)
		file, err := ioutil.ReadFile(*file)
		if err != nil {
			log.Fatal(err)
		}
		content = string(file)
	}

	createdGistUrl, creationError := createGist(configuration, filename, content)
	if creationError != nil {
		log.Fatal(creationError)
	}
	fmt.Printf("%s\n", createdGistUrl)
}
