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
	itemHeight := 21
	svgHeight := len(repos) * itemHeight
	maxStars := 0
	for _, repo := range repos {
		if repo.Stars > maxStars {
			maxStars = repo.Stars
		}
	}
	maxDigits := len(strconv.Itoa(maxStars))
	if maxDigits > 3 {
		maxDigits = 3
	}
	starPadding := "  "
	if maxDigits == 2 {
		starPadding = " "
	}
	if maxDigits == 3 {
		starPadding = ""
	}
	svgString += fmt.Sprintf("<svg width='800' height='%d' xmlns='http://www.w3.org/2000/svg' xmlns:xlink='http://www.w3.org/1999/xlink'>", svgHeight)
	svgString += "<!--. this svg was generated with github-scoreboard https://github.com/donuts-are-good/github-scoreboard  .-->"
	svgString += "<rect width='800' height='" + strconv.Itoa(svgHeight) + "' fill='white' />"

	for i, repo := range repos {
		escapedDescription := html.EscapeString(repo.Description)
		indexPadding := "  "
		if i >= 9 {
			indexPadding = " "
		}
		starString := strconv.Itoa(repo.Stars)
		if repo.Stars >= 1000 {
			starString = fmt.Sprintf("%.1fk", float64(repo.Stars)/1000.0)
		}
		if repo.Stars >= 1000000 {
			starString = fmt.Sprintf("%.1fM", float64(repo.Stars)/1000000.0)
		}
		svgString += fmt.Sprintf("<a xlink:href='https://github.com/donuts-are-good/%s'><text x='20' y='%d' fill='black' font-family='monospace'>%d.%s⭐[%s]%s - %s - %s</text></a>", repo.Name, 20+20*i, i+1, indexPadding, starString, starPadding, repo.Name, escapedDescription)
	}
	svgString += "</svg>"
	return svgString
}
