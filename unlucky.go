package main

import (
	"fmt"
	//"io"
	"strings"
	"os"
	"time"
	//"path"
	"net/http"
	"golang.org/x/net/html"
)

func main() {
	urlList := []string {
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

	start:= time.Now()

	urlTxt, err:= os.Create("foundUrls.txt")
	if err != nil {
		os.Exit(1)
	}

	defer urlTxt.Close()

	imgTxt, err:= os.Create("foundImages.txt")
	if err != nil {
		os.Exit(1)
	}

	defer imgTxt.Close()

	for i:= 0; i < len(urlList); i++ {
		//open the given URl
		reader, err := http.Get(urlList[i])
		if err != nil {
			break
		}

		defer reader.Body.Close()

		tokenizer := html.NewTokenizer(reader.Body)
		//tokenizer := html.NewTokenizerFragment(reader.Body, "img")
		//loop through html body
		for {
			nextToken := tokenizer.Next()
			//if eof reached, exit
			if nextToken == html.ErrorToken {
				// Returning io.EOF indicates success.
				break
			}

			token := tokenizer.Token()

			//check if the text associated with the tag is "img"
			if strings.Contains(token.String(), "<img") {
				parser, err := html.Parse(strings.NewReader(token.String()))
				if err != nil {
					break
				}

				var stringParse func(*html.Node)
				stringParse = func(parser *html.Node) {
					if parser.Type == html.ElementNode && parser.Data == "img" {
						for _, a := range parser.Attr {
							if a.Key == "src" {
								img :=[]byte(a.Val + "\n")
								_, err := imgTxt.Write(img)
								if err != nil {
									os.Exit(1)
								}

								break
							}
						}
					}
					for i := parser.FirstChild; i != nil; i = i.NextSibling {
						stringParse(i)
					}
				}

				stringParse(parser)
			}
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("Downloads completed in %s \n", elapsed)
}