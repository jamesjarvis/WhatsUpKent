package scrape

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// The point of this section is to concurrently download ical files from a specified ID, and cache them on the system.
// This should also provide a function which can delete a cache if it exists

// Adds all the urls to be scraped to a channel
func GetIds(chIds chan int) {
	START := 130000
	END := 150000

	for i := START; i < END; i++ {
		chIds <- i
	}
}

// Formats the id to the actual ical file url
func FormatURL(id int) string {
	// https://www.kent.ac.uk/timetabling/ical/132056.ics
	return fmt.Sprintf("https://www.kent.ac.uk/timetabling/ical/%d.ics", id)
}

func FormatFilename(id int) string {
	return filepath.Join("ical_cache", fmt.Sprintf("%d.ics", id))
}

// Makes the request and saves the result to a file
func DownloadFile(id int, filepath string) error {

	url := FormatURL(id)

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	} else if resp.StatusCode != 200 {
		if resp.StatusCode == 500 {
			fmt.Printf("%d returned %d", id, resp.StatusCode)
		}
		return errors.New("Invalid ID")
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		fmt.Print(err)
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	return err
}

// Downloads all of the urls in the channel
func DownloadAll(chIds chan int, chFiles chan string) {
	//Set up cache directory
	os.Mkdir("ical_cache", os.FileMode(0755))

	numberOfWorkers := 1
	for i := 0; i < numberOfWorkers; i++ {
		go worker(chIds, i, chFiles)
	}
}

func worker(queue chan int, worknumber int, chFiles chan string) {
	for true {
		select {
		case id := <-queue:
			filepath := FormatFilename(id)
			err := DownloadFile(id, filepath)
			fmt.Println("doing work!", id, "worknumber", worknumber)
			if err == nil {
				chFiles <- filepath
			}
		}
	}
}

// Sends the file to be scraped, and once that is complete, it deletes the cached file
func ProcessFile(filename string) {
	ParseCal(filename)

	// Remove the cache
	os.Remove(filename)
}

// Scrapes all files in the channel
func ProcessAll(chFiles chan string) {
	for filename := range chFiles {
		go ProcessFile(filename)
	}
}

//Runs the main program
func FuckIt() {
	//Channels
	chIds := make(chan int)      //Channel of ids to be downloaded
	chFiles := make(chan string) //Channel of filenames to be scraped

	go GetIds(chIds) //Populate the URLs to be downloaded

	// Download all of the urls in the channel
	DownloadAll(chIds, chFiles)

	// Whilst this is happening, scrape all the files in the other channel
	ProcessAll(chFiles)
}
