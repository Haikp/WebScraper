/*
Description:

	Given a list of URLs to webscrape
	Tokenize and parse the url, see its html files
	obtain full, redirecting links and source images
	create textfiles to contain found images and urls
	create directory and download all images
	use go routines to parallelize the process
*/
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
	"golang.org/x/net/html"
)

// global variable, tracks if a routine has finished
var routineTracker []bool

func main() {
	urlList := []string{
		"https://www.unlv.edu/cs",
		"https://www.unlv.edu/engineering",
		"https://www.unlv.edu/engineering/advising-center",
		"https://www.unlv.edu/engineering/about",
		"https://www.unlv.edu/engineering/academic-programs",
		"https://www.unlv.edu/ceec",
		"https://ece.unlv.edu/",
		"https://www.unlv.edu/me",
		"https://www.unlv.edu/rotc",
		"https://www.unlv.edu/afrotc",
		"https://www.unlv.edu/eed",
		"https://www.unlv.edu/engineering/mendenhall",
		"https://www.unlv.edu/engineering/uas",
		"https://www.unlv.edu/engineering/solar",
		"https://www.unlv.edu/engineering/techcommercialization",
		"https://www.unlv.edu/engineering/railroad",
		"https://www.unlv.edu/engineering/future-students",
		"https://www.physics.unlv.edu/",
	}

	start := time.Now()

	//create text file to hold urls
	urlTxt, err := os.Create("foundUrls.txt")
	if err != nil {
		panic(err)
	}

	defer urlTxt.Close()

	//create text file to hold images
	imgTxt, err := os.Create("foundImages.txt")
	if err != nil {
		panic(err)
	}

	defer imgTxt.Close()

	//create channels to find
	urlCh := make(chan []byte)
	imgCh := make(chan []byte)

	for i := 0; i < len(urlList); i++ {
		//call individual go routines for each website
		go webScrape(urlList[i], urlCh, imgCh, i)
		//create an entry for each routine called
		routineTracker = append(routineTracker, false)
	}

	moreUrls := true
	moreImgs := true
	var urlInput, imgInput []byte

	//creates a new directory for the downloaded images
	dirName := "downloadedImages"
	os.Mkdir(dirName, 0750)

	//add a url to the hashmap, make sure that duplicates do not get re-entered
	dupCheck := make(map[string]bool)

	//loop while the channels are both still open
	for moreUrls == true || moreImgs == true {
		select {
		//upon an entry for url is found
		case urlInput, moreUrls = <-urlCh:
			//if there are not duplicates
			if moreUrls == true && dupCheck[string(urlInput)] != true {
				dupCheck[string(urlInput)] = true
				//write to the url text file
				_, err := urlTxt.Write(urlInput)
				if err != nil {
					panic(err)
				}
			}
			//upon an entry found for images
		case imgInput, moreImgs = <-imgCh:
			//check if there are no duplicates
			if moreImgs == true && dupCheck[string(imgInput)] != true {
				dupCheck[string(imgInput)] = true
				//write to image file
				_, err := imgTxt.Write(imgInput)
				if err != nil {
					panic(err)
				}

				//obtain valid url of the image
				splitUrl := strings.Split("https://www.unlv.edu"+string(imgInput), "\n")
				fullUrl := splitUrl[0]

				//request access to the url
				validUrl, err := http.Get(fullUrl)
				if err != nil {
					//edge case, if the website has a different domain
					splitUrl = strings.Split("https://www.physics.unlv.edu/"+string(imgInput), "\n")
					fullUrl = splitUrl[0]
					validUrl, err = http.Get(fullUrl)
					if err != nil {
						panic(err)
					}
				}

				defer validUrl.Body.Close()

				//name of the image will be the associating .png/.jpg/.jpeg within the image name
				splitStr := strings.Split(string(imgInput), "/")
				lastSplit := splitStr[len(splitStr)-1]

				locate := strings.Index(lastSplit, ".")
				var fileName string
				//slice the string name of the image to get the final fileName to use
				if strings.Contains(lastSplit, "jpeg") {
					fileName = lastSplit[0:(locate + 5)]
				} else {
					fileName = lastSplit[0:(locate + 4)]
				}

				//create the file within the directory created earlier
				dlImage, err := os.Create(path.Join(dirName, fileName))
				if err != nil {
					panic(err)
				}

				defer dlImage.Close()

				//copy/download the image into the file
				_, err = io.Copy(dlImage, validUrl.Body)
				if err != nil {
					panic(err)
				}
			}
		default:
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("Downloads completed in %s \n", elapsed)
}

func webScrape(url string, urlCh chan []byte, imgCh chan []byte, pid int) {
	//request access to the given URL
	reader, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer reader.Body.Close()

	//begin searching through the html file with a tokenizer
	tokenizer := html.NewTokenizer(reader.Body)

	//search till end of file
	for {
		nextToken := tokenizer.Next()
		//if eof reached, exit
		if nextToken == html.ErrorToken {
			// Returning io.EOF indicates success.
			break
		}

		token := tokenizer.Token()
		var stringParse func(*html.Node)
		//search the html until "<img" or "<a" is found
		if strings.Contains(token.String(), "<img") {
			//parse the appropiate line of code
			parser, err := html.Parse(strings.NewReader(token.String()))
			if err != nil {
				break
			}

			//recursively look through the parsed string
			stringParse = func(parser *html.Node) {
				for _, a := range parser.Attr {
					//upon finding "src", store the text within it, send to appropriate channel
					if a.Key == "src" {
						img := []byte(a.Val + "\n")
						//send to image channel
						imgCh <- img
						break
					}
				}

				//recursion loop
				for i := parser.FirstChild; i != nil; i = i.NextSibling {
					stringParse(i)
				}
			}

			stringParse(parser)
			//else statement to check if "<a" tag is found instead
		} else if strings.Contains(token.String(), "<a") {
			//parse the appropiate line of code
			parser, err := html.Parse(strings.NewReader(token.String()))
			if err != nil {
				break
			}

			//recursively look through the parsed string
			stringParse = func(parser *html.Node) {
				for _, a := range parser.Attr {
					//search for "href" and obtain its contents
					if a.Key == "href" {
						//make sure the link within "href" contains a full link
						if strings.Contains(a.Val, "http") {
							fullUrl := []byte(a.Val + "\n")
							//send to url channel
							urlCh <- fullUrl
						}

						break
					}
				}

				//recursion loop
				for i := parser.FirstChild; i != nil; i = i.NextSibling {
					stringParse(i)
				}
			}

			stringParse(parser)
		}
	}

	//once a routine has finished, set its corresponding entry in the global slice to true
	routineTracker[pid] = true

	//once all routines have finished, the LAST routine will close the channels
	for i := 0; i < len(routineTracker); i++ {
		//if an entry in the array of routines have not finished, exit the function
		if routineTracker[i] == false {
			return
		}
	}

	//close the channels ONLY IF all routines have finished, the last routine will close channels
	close(urlCh)
	close(imgCh)
}
