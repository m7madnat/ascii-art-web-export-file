package main

import (
	"bufio"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

var tmpl *template.Template

func main() {
	tmpl = template.Must(template.ParseGlob("*.html"))
	fmt.Println("Server is running at http://localhost:8080/")
	fs := http.FileServer(http.Dir("templates"))
	http.Handle("/templates/", http.StripPrefix("/templates", fs))
	http.HandleFunc("/", asciiArtHandler)
	http.HandleFunc("/export", handleExport)
	http.ListenAndServe(":8080", nil)
}

func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	bannerStyle := r.FormValue("banner_style")

	file, _ := os.Open(fmt.Sprintf("%s.txt", bannerStyle))

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	asciiChrs := make(map[int][]string)
	dec := 31
	for _, line := range lines {
		if line == "" {
			dec++
		} else {
			asciiChrs[dec] = append(asciiChrs[dec], line)
		}
	}
	result := generateAsciiArt(text, asciiChrs)

	tmpl.ExecuteTemplate(w, "index.html", result)
}

func handleExport(w http.ResponseWriter, r *http.Request) {
	asciiArt := r.URL.Query().Get("result")
	w.Write([]byte(asciiArt))
}

func generateAsciiArt(text string, asciiChrs map[int][]string) string {
	var result strings.Builder
	lines := strings.Split(text, "\r\n")
	for _, line := range lines {
		words := strings.Split(line, " ")
		for j := 0; j < len(asciiChrs[32]); j++ {
			for _, word := range words {
				for _, letter := range word {
					result.WriteString(asciiChrs[int(letter)][j])
				}
				result.WriteString("  ")
			}
			result.WriteString("\n")
		}
	}
	return result.String()
}
