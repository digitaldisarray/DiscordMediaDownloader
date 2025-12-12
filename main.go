package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
		return
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	// Unmarshal byte array into Messages object
	var messages Messages
	json.Unmarshal(byteValue, &messages)

	// Get absolute downloads directory path
	downloadsDir := "downloads"
	err = os.MkdirAll(downloadsDir, 0755)
	if err != nil {
		fmt.Println("Error creating downloads directory:", err)
		return
	}

	// Extract data from parsed JSON
	var messageIds []string
	for _, message := range messages.Messages {
		if message.Content == "" && len(message.Attachments) > 0 {
			for _, attachment := range message.Attachments {
				url := attachment.Url
				name := attachment.FileName
				if !IsMediaFile(name) {
					fmt.Println("Rejecting URL:", url)
					continue
				}

				fmt.Println("Downloading (attach):", url)
				safePath := GetUniqueFilePath(filepath.Join(downloadsDir, name))
				DownloadFile(safePath, url)
				messageIds = append(messageIds, message.Id)
			}
		} else if strings.HasPrefix(message.Content, "http") {
			url := ExtractUrl(message.Content)
			name := ExtractFileName(url)
			if !IsMediaFile(name) {
				fmt.Println("Rejecting URL:", url)
				continue
			}

			fmt.Println("Downloading (url msg):", url)
			safePath := GetUniqueFilePath(filepath.Join(downloadsDir, name))
			DownloadFile(safePath, url)
			messageIds = append(messageIds, message.Id)
		}
	}

	fmt.Println("Deleted Messages:", strings.Join(messageIds, ","))
}

func IsMediaFile(name string) bool {
	extensions := []string{".png", ".jpg", ".jpeg", ".gif", ".mp4", ".webm"}
	for _, ext := range extensions {
		if strings.HasSuffix(name, ext) {
			return true
		}
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

func GetUniqueFilePath(filepath string) string {
	ext := filepath[strings.LastIndex(filepath, "."):]
	base := filepath[:strings.LastIndex(filepath, ".")]

	counter := 1
	newPath := filepath

	for {
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			break
		}
		newPath = fmt.Sprintf("%s_%d%s", base, counter, ext)
		counter++
	}

	return newPath
}

func DownloadFile(filepath string, url string) error {
	// Get data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

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
