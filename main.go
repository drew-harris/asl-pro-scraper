package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/drew-harris/asl-pro/database"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter tag to use: ")
	scanner.Scan()
	tag := scanner.Text()
	err := database.SaveTag(tag)
	if err != nil {
		panic(err)
	}
	fmt.Println("Saved tag: " + tag)

	for {
		fmt.Print("\n1. Download From Input.txt\n2. Download individual words\n3. Exit\nWhich would you like to do: ")
		scanner.Scan()
		choice := scanner.Text()
		fmt.Print("\n")
		if choice == "1" {
			downloadFromFile(tag)
		} else if choice == "2" {
			typeToDownload(scanner, tag)
		} else {
			break
		}
	}
}

func typeToDownload(scanner *bufio.Scanner, tag string) {
	for {
		fmt.Print("Enter word: ")
		scanner.Scan()
		word := scanner.Text()
		if word == "exit" {
			break
		}
		err := downloadWord(word)
		if err != nil {
			fmt.Println(word + ": " + err.Error())
		} else {
			fmt.Println("Downloaded word: " + word)
			database.SaveWord(word, tag)
			fmt.Println("Saved word to DB: " + word)
		}
	}
}

func downloadFromFile(tag string) {
	fmt.Println("Getting words from input.txt")
	words, err := getWords("./input.txt")
	if err != nil {
		panic(err)
	}

	var missed []string
	for _, word := range words {
		err = downloadWord(word)
		if err != nil {
			fmt.Println(word + ": " + err.Error())
			missed = append(missed, word)
		} else {
			fmt.Println("Downloaded word: " + word)
			database.SaveWord(word, tag)
			fmt.Println("Saved word to DB: " + word)
		}
	}

	fmt.Println("\n\nManually Required Words: ")
	for _, word := range missed {
		fmt.Println(word)
	}
	fmt.Print("\n\n")

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

func downloadWord(word string) error {
	client := &http.Client{}
	url := fmt.Sprintf("http://www.aslpro.cc/main/%s/%s.mp4", word[0:1], word)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Referer", url)
	res, error := client.Do(req)
	if error != nil {
		return errors.New("could not generate http request")
	}
	defer res.Body.Close()
	if res.ContentLength < 5000 {
		return errors.New("word not found")
	}

	out, err := os.Create("./output/" + word + ".mp4")
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return errors.New("could not copy data to path")
	}

	return nil

}
