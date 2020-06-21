package scrape

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/dgraph-io/dgo/v200"
)

// FilesIds for handling the currently scraped file
type FilesIds struct {
	id       int
	filename string
}

//InitialConfig is the configuration passed into the scraper
type InitialConfig struct {
	Url          string
	StartRange   int
	EndRange     int
	SlowInterval time.Duration
	MaxAge       time.Duration
	//DownloadPool is the max concurrent ical downloads
	DownloadPool int
	//ProcessPool is the max concurrent files being processed
	ProcessPool int
	//EventProcessPool is the number of workers spawned to process the events within a file being parsed
	//If you have 3 event process workers, and 4 process workers, then you'll have 12 concurrent event workers
	EventProcessPool int
	DBClient         *dgo.Dgraph
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

// // FormatFilename formats the integer id into the filename path
// func FormatFilename(id int) string {
// 	return filepath.Join(getIcalDir(), fmt.Sprintf("%d.ics", id))
// }

// func getIcalDir() string {
// 	ex, err := os.Executable()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	exPath := filepath.Join(filepath.Dir(ex), "ical_cache")
// 	return exPath
// }

// DownloadFile makes the request and saves the result to a file
func DownloadFile(id int) (*FilesIds, error) {

	url := FormatURL(id)

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		if resp.StatusCode >= 500 {
			return nil, ErrUniversityPanicking
		}
		return nil, ErrInvalidID
	}
	defer resp.Body.Close()

	// Create the file
	tmpfile, err := ioutil.TempFile("", "*.ics")
	if err != nil {
		return nil, err
	}

	// Write the body to file
	_, err = io.Copy(tmpfile, resp.Body)

	defer tmpfile.Close()

	// Create mapping to temp file
	fid := &FilesIds{
		id:       id,
		filename: tmpfile.Name(),
	}

	return fid, nil
}

// //CreateDownloadDir does what it says on the tin
// func CreateDownloadDir() error {
// 	//Set up cache directory
// 	return os.Mkdir(getIcalDir(), os.FileMode(0755))
// }

// downloadAll downloads all of the urls in the channel
func downloadAll(chIds chan int, chFiles chan FilesIds, config *InitialConfig) {
	var downloadWG sync.WaitGroup

	numberOfWorkers := config.DownloadPool
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
		fid, err := DownloadFile(id)

		if err == nil {
			chFiles <- *fid
		} else if err != ErrInvalidID {
			log.Println(err)
		}
	}
	wg.Done()
}

// processAll scrapes all files in the channel
func processAll(c *dgo.Dgraph, chFiles chan FilesIds, config *InitialConfig) {
	var processWG sync.WaitGroup
	var eventMX = &sync.Mutex{}

	numberOfWorkers := config.ProcessPool
	for i := 0; i < numberOfWorkers; i++ {
		processWG.Add(1)
		go processWorker(c, chFiles, eventMX, &processWG, config)
	}
	processWG.Wait()
	log.Println("------- processAll completed ------")
}

func processWorker(c *dgo.Dgraph, chFiles chan FilesIds, mx *sync.Mutex, wg *sync.WaitGroup, config *InitialConfig) {
	for filename := range chFiles {
		err := ProcessFile(c, filename, mx, config)

		if err != nil {
			log.Fatal(err)
		}
	}
	wg.Done()
}

// ProcessFile sends the file to be scraped, and once that is complete, it deletes the cached file
func ProcessFile(c *dgo.Dgraph, fid FilesIds, mx *sync.Mutex, config *InitialConfig) error {
	err := ParseCal(c, fid, mx, config)

	// Remove the cache
	os.Remove(fid.filename)
	return err
}

// FuckIt Runs the main scraping program
func FuckIt(config *InitialConfig, c *dgo.Dgraph) {
	//Channels
	chIds := make(chan int, 100) //Channel of ids to be downloaded

	chFiles := make(chan FilesIds, 50) //Channel of filenames to be scraped

	go GetIds(config, chIds) //Populate the URLs to be downloaded

	// Download all of the urls in the channel
	go downloadAll(chIds, chFiles, config)

	// Whilst this is happening, scrape all the files in the other channel
	processAll(c, chFiles, config)
}
