package scrape

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/dgraph-io/dgo/v2"
)

// FilesIds for handling the currently scraped file
type FilesIds struct {
	id       int
	filename string
}

//InitialConfig is the configuration passed into the scraper
type InitialConfig struct {
	StartRange   int
	EndRange     int
	SlowInterval time.Duration
}

// The point of this section is to concurrently download ical files from a specified ID, and cache them on the system.
// This should also provide a function which can delete a cache if it exists

// GetIds adds all the urls to be scraped to a channel
func GetIds(config *InitialConfig, chIds chan int) error {
	if config == nil {
		return ErrConfig
	}
	for i := config.StartRange; i < config.EndRange; i++ {
		chIds <- i
	}
	close(chIds)
	log.Println("Closed ids to download channel")
	return nil
}

// FormatURL formats the id to the actual ical file url
func FormatURL(id int) string {
	return fmt.Sprintf("https://www.kent.ac.uk/timetabling/ical/%d.ics", id)
}

// FormatFilename formats the integer id into the filename path
func FormatFilename(id int) string {
	return filepath.Join(getIcalDir(), fmt.Sprintf("%d.ics", id))
}

func getIcalDir() string {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exPath := filepath.Join(filepath.Dir(ex), "ical_cache")
	return exPath
}

// DownloadFile makes the request and saves the result to a file
func DownloadFile(fid FilesIds) error {

	url := FormatURL(fid.id)

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	} else if resp.StatusCode != 200 {
		if resp.StatusCode >= 500 {
			return ErrUniversityPanicking
		}
		return ErrInvalidID
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(fid.filename)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	return err
}

// downloadAll downloads all of the urls in the channel
func downloadAll(chIds chan int, chFiles chan FilesIds) {
	//Set up cache directory
	os.Mkdir(getIcalDir(), os.FileMode(0755))

	var downloadWG sync.WaitGroup

	numberOfWorkers := 1
	for i := 0; i < numberOfWorkers; i++ {
		downloadWG.Add(1)
		go worker(chIds, i, chFiles, &downloadWG)
	}
	downloadWG.Wait()
	close(chFiles)
	log.Println("------- downloadAll completed ------")
}

func worker(queue chan int, worknumber int, chFiles chan FilesIds, wg *sync.WaitGroup) {
	for id := range queue {
		fid := FilesIds{
			id:       id,
			filename: FormatFilename(id),
		}
		err := DownloadFile(fid)

		if err == nil {
			chFiles <- fid
		} else if err != ErrInvalidID {
			log.Println(err)
		}
	}
	// log.Printf("Exiting download worker %d", worknumber)
	wg.Done()
}

// processAll scrapes all files in the channel
func processAll(c *dgo.Dgraph, chFiles chan FilesIds) {
	var processWG sync.WaitGroup
	var eventMX = &sync.Mutex{}
	for filename := range chFiles {
		processWG.Add(1)
		go ProcessFile(c, filename, eventMX, &processWG)
	}
	processWG.Wait()
	log.Println("------- processAll completed ------")
}

// ProcessFile sends the file to be scraped, and once that is complete, it deletes the cached file
func ProcessFile(c *dgo.Dgraph, fid FilesIds, mx *sync.Mutex, wg *sync.WaitGroup) {
	err := ParseCal(c, fid, mx)
	if err != nil {
		log.Fatal(err)
	}

	// Remove the cache
	os.Remove(fid.filename)
	wg.Done()
}

// FuckIt Runs the main scraping program
func FuckIt(config *InitialConfig, c *dgo.Dgraph) {
	//Channels
	chIds := make(chan int, 1000) //Channel of ids to be downloaded

	chFiles := make(chan FilesIds, 100) //Channel of filenames to be scraped

	go GetIds(config, chIds) //Populate the URLs to be downloaded

	// Download all of the urls in the channel
	go downloadAll(chIds, chFiles)

	// Whilst this is happening, scrape all the files in the other channel
	processAll(c, chFiles)
}
