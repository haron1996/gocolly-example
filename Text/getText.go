package text

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
)

// Define a struct to represent each entry in the JSON file
type Comment struct {
	VideoWebUrl       string `json:"videoWebUrl"`
	SubmittedVideoUrl string `json:"submittedVideoUrl"`
	Cid               string `json:"cid"`
	CreateTime        int64  `json:"createTime"`
	CreateTimeISO     string `json:"createTimeISO"`
	Text              string `json:"text"`
	DiggCount         int    `json:"diggCount"`
	RepliesToId       string `json:"repliesToId"`
	ReplyCommentTotal int    `json:"replyCommentTotal"`
	Uid               string `json:"uid"`
	UniqueId          string `json:"uniqueId"`
	AvatarThumbnail   string `json:"avatarThumbnail"`
}

func GetText() {
	// Read the JSON file
	data, err := os.ReadFile("comments.json")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// Define a slice to hold the comments
	var comments []Comment

	// Unmarshal the JSON data into the slice of comments
	err = json.Unmarshal(data, &comments)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	// Create a new CSV file
	file, err := os.Create("comments.csv")
	if err != nil {
		log.Fatalf("Error creating CSV file: %v", err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{"Text"})

	// Write each comment to the CSV file
	for _, comment := range comments {
		err := writer.Write([]string{comment.Text})
		if err != nil {
			log.Fatalf("Error writing to CSV: %v", err)
		}
	}
}
