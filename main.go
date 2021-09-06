package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Messages struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	Id          string       `json:"id"`
	Content     string       `json:"content"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Url      string `json:"url"`
	FileName string `json:"fileName"`
}

func main() {

	// Read messages.json into byte array
	jsonFile, err := os.Open("messages.json")
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// Unmarshal byte array into Messages object
	var messages Messages
	json.Unmarshal(byteValue, &messages)

	// Extract data from parsed json
	for _, message := range messages.Messages {
		if message.Content == "" && message.Attachments != nil {
			for _, attachment := range message.Attachments {
				url := attachment.Url
				name := attachment.FileName
				if !IsMediaFile(name) {
					fmt.Println("Rejecting URL: " + url)
					continue
				}

				fmt.Println("Downloading (attach) : " + url)
				DownloadFile("downloads\\"+name, url)
			}
		} else if strings.HasPrefix(message.Content, "http") {
			url := ExtractUrl(message.Content)
			name := ExtractFileName(url)
			if !IsMediaFile(name) {
				fmt.Println("Rejecting URL: " + url)
				continue
			}

			fmt.Println("Downloading (url msg): " + url)
			DownloadFile("downloads\\"+name, url)
		}
	}

}

func IsMediaFile(name string) bool {
	if strings.HasSuffix(name, ".png") || strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".jpeg") || strings.HasSuffix(name, ".gif") || strings.HasSuffix(name, ".mp4") || strings.HasSuffix(name, ".webm") {
		return true
	}

	return false
}

func ExtractFileName(url string) string {
	return url[strings.LastIndex(url, "/")+1:]
}

func ExtractUrl(content string) string {
	extracted := strings.Split(content, " ")[0]
	questionRemoved := strings.Split(extracted, "?")[0]
	return questionRemoved
}

func DownloadFile(filepath string, url string) error {
	// Get data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// If file already exists
	_, err = os.Stat(filepath)
	coutner := 1
	for !os.IsNotExist(err) {
		// File exists

		// Add counter
		index := strings.LastIndex(filepath, ".")
		filepath = filepath[:index] + fmt.Sprintf("_%d", coutner) + filepath[index:]

		// Check with new name
		_, err = os.Stat(filepath)
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
