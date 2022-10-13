package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type EnpassJson struct {
	CustomIcons []struct {
		Data string `json:"data"`
		UUID string `json:"uuid"`
	} `json:"custom_icons"`
	Folders []struct {
		Icon       string `json:"icon"`
		ParentUUID string `json:"parent_uuid"`
		Title      string `json:"title"`
		UpdatedAt  int    `json:"updated_at"`
		UUID       string `json:"uuid"`
	} `json:"folders"`
	Items []struct {
		Archived   int    `json:"archived"`
		AutoSubmit int    `json:"auto_submit"`
		Category   string `json:"category"`
		CreatedAt  int    `json:"createdAt"`
		Favorite   int    `json:"favorite"`
		Fields     []struct {
			Deleted        int    `json:"deleted"`
			Label          string `json:"label"`
			Order          int    `json:"order"`
			Sensitive      int    `json:"sensitive"`
			Type           string `json:"type"`
			UID            int    `json:"uid"`
			UpdatedAt      int    `json:"updated_at"`
			Value          string `json:"value"`
			ValueUpdatedAt int    `json:"value_updated_at"`
		} `json:"fields"`
		Icon struct {
			Fav   string `json:"fav"`
			Image struct {
				File string `json:"file"`
			} `json:"image"`
			Type int    `json:"type"`
			UUID string `json:"uuid"`
		} `json:"icon"`
		Note         string `json:"note"`
		Subtitle     string `json:"subtitle"`
		TemplateType string `json:"template_type"`
		Title        string `json:"title"`
		Trashed      int    `json:"trashed"`
		UpdatedAt    int    `json:"updated_at"`
		UUID         string `json:"uuid"`
	} `json:"items"`
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalln("\n./enpass-to-keepassxc-macM1 input_file.json output_file.csv\nOR\ngo run main.go input_file.json output_file.csv")
	}
	input := os.Args[1]
	output := os.Args[2]
	b, err := os.ReadFile(input)
	if err != nil {
		log.Fatalln(err)
	}
	enpassData := &EnpassJson{}
	err = json.Unmarshal(b, enpassData)
	if err != nil {
		log.Fatalln(err)
	}

	keepassxc_records := [][]string{
		{"Group", "Title", "Username", "Password", "URL", "Notes", "TOTP", "Icon", "Last Modified", "Created"},
	}

	for _, item := range enpassData.Items {
		var username, password, url, totp string

		for _, fl := range item.Fields {
			if fl.Type == "username" || fl.Type == "email" {
				if fl.Value != "" {
					username = fl.Value
				}
			}
			if fl.Type == "password" {
				password = fl.Value
			}
			if fl.Type == "url" {
				url = fl.Value
			}
			if fl.Type == "totp" {
				totp = fl.Value
			}
		}
		record := []string{
			item.Category,
			item.Title,
			username,
			password,
			url,
			item.Note,
			totp,
			"",
			fmt.Sprintf("%d", item.UpdatedAt),
			fmt.Sprintf("%d", item.CreatedAt),
		}
		keepassxc_records = append(keepassxc_records, record)
	}
	f, err := os.Create(output)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	err = w.WriteAll(keepassxc_records)
	if err := w.Error(); err != nil {
		log.Fatalln(err)
	}
	log.Println("Success")
}
