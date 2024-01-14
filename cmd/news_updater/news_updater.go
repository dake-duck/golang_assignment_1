package main

import (
	db "DAKExDUCK/assignment_1/internal/db"
	pgsql "DAKExDUCK/assignment_1/pkg/models/pgsql"
	"encoding/json"
	"flag"
	"fmt"
)

func UpdateNewsFromMoodle() {
	dsn := flag.String("dsn", "user:pass@host:port", "MySQL data source name")
	flag.Parse()

	db, err := db.OpenDB(*dsn)
	if err != nil {
		fmt.Println(err)
	}
	model := pgsql.NewsModel{DB: db}

	moodleAPI := NewMoodleAPI()
	token := "741fb115dcfa9074af6b8634cbf6febe"
	params := map[string]string{
		"forumid": "17",
		"perpage": "10",
		"page":    "0",
	}

	result, err := moodleAPI.MakeRequest("mod_forum_get_forum_discussions", token, params)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var posts []any
	posts, ok := result["discussions"].([]any)
	fmt.Println(posts)
	if ok == false {
		jsonStr, err := json.Marshal(result)
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
		} else {
			fmt.Printf("Error: no key or type conversion failed! %s", jsonStr)
		}
		return
	}

	for i := 0; i < len(posts); i++ {
		post, ok := posts[i].(map[string]any)
		if ok == false {
			jsonStr, err := json.Marshal(post)
			if err != nil {
				fmt.Printf("Error: %s", err.Error())
			} else {
				fmt.Printf("Error: no key or type conversion failed! %s", jsonStr)
			}
			return
		}

		moodle_id := post["id"].(float64)
		header_html := post["name"]
		body_html := post["message"]
		if_attachement := post["attachment"].(bool)
		attachements := []map[string]any{}
		if !if_attachement {
			attachements = post["attachments"].([]map[string]any)
		}

		fmt.Println(moodle_id)
		fmt.Println(body_html)
		fmt.Println(header_html)
		fmt.Println(attachements)
		fmt.Println("= = = = = = = = = =")
		model.Insert(
			int(moodle_id),
			header_html.(string),
			body_html.(string),
			attachements,
		)
	}
}

func main() {
	UpdateNewsFromMoodle()
}
