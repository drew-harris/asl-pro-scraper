package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
)

func main() {
	color.HiCyan("Getting words from input.txt")
	words, err := getWords("./input.txt")
	if err != nil {
		panic(err)
	}
	missed := downloadWords(words)
	fmt.Println("\n\nManually Required Words: ")
	for _, word := range missed {
		fmt.Println("  " + word + "\n")
	}
}

func getWords(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	var words []string
	for scanner.Scan() {
		words = append(words, strings.ReplaceAll(strings.ToLower(scanner.Text()), " ", "_"))
	}

	if err := scanner.Err(); err != nil {
		log.Panic(err)
		return nil, err
	}

	return words, nil
}

func downloadWords(words []string) []string {
	client := &http.Client{}
	var missed []string
	for _, word := range words {
		url := fmt.Sprintf("http://www.aslpro.cc/main/%s/%s.mp4", word[0:1], word)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Referer", url)
		res, error := client.Do(req)
		if error != nil {
			fmt.Println(error.Error())
			missed = append(missed, word)
			color.Red("Could not download word: " + word)
			continue
		}
		defer res.Body.Close()
		if res.ContentLength < 5000 {
			missed = append(missed, word)
			color.Red("Could not download word: " + word)
			continue
		}

		out, err := os.Create("./output/" + word + ".mp4")
		if err != nil {
			missed = append(missed, word)
			color.Red("Could not download word: " + word)
			fmt.Println(color.RedString(err.Error()))
			continue
		}

		defer out.Close()

		_, err = io.Copy(out, res.Body)
		if err != nil {
			missed = append(missed, word)
			color.Red("Could not download word: " + word)
		}
		color.Green("Downloaded word: " + word)

	}
	return missed

}
