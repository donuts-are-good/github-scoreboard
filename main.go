package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
)

type Repository struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Stars       int    `json:"stargazers_count,omitempty"`
	Forks       int    `json:"forks_count,omitempty"`
	Watchers    int    `json:"watchers_count,omitempty"`
}

func main() {
	username := "donuts-are-good"
	url := fmt.Sprintf("https://api.github.com/users/%s/repos", username)

	res, err := http.Get(url)
	handle(err)
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	handle(err)

	var repos []Repository
	err = json.Unmarshal(body, &repos)
	handle(err)

	sort.Slice(repos, func(i, j int) bool {
		return repos[i].Stars > repos[j].Stars
	})

	svgString := generateSVG(repos)

	file, err := os.Create("repo_list.svg")
	handle(err)
	defer file.Close()

	_, err = file.WriteString(svgString)
	handle(err)
}

func handle(err error) {
	if err != nil {
		log.Println(err)
		return
	}
}

func generateSVG(repos []Repository) string {
	var svgString string

	itemHeight := 20
	svgHeight := len(repos) * itemHeight
	svgString += fmt.Sprintf("<svg width='800' height='%d' xmlns='http://www.w3.org/2000/svg' xmlns:xlink='http://www.w3.org/1999/xlink'>", svgHeight)
	svgString += "<rect width='800' height='" + strconv.Itoa(svgHeight) + "' fill='white' />"

	for i, repo := range repos {
		escapedDescription := html.EscapeString(repo.Description)
		if i < 9 {
			svgString += fmt.Sprintf("<a xlink:href='https://github.com/donuts-are-good/%s'><text x='20' y='%d' fill='black' font-family='monospace'>%d.  ⭐[%d] - %s - %s</text></a>", repo.Name, 20+20*i, i+1, repo.Stars, repo.Name, escapedDescription)
		} else {
			svgString += fmt.Sprintf("<a xlink:href='https://github.com/donuts-are-good/%s'><text x='20' y='%d' fill='black' font-family='monospace'>%d. ⭐[%d] - %s - %s</text></a>", repo.Name, 20+20*i, i+1, repo.Stars, repo.Name, escapedDescription)
		}
	}

	svgString += "</svg>"

	return svgString
}
