package main

import (
    "fmt"
    // "io"
    // "log"
    "net/http"
    "os"
    "path"
   //  "strings"
)

func main() {
	//defer os.RemoveAll("dowloadedImages")
   os.Mkdir("dowloadedImages", 0750)

   imgUrl := "/sites/default/files/styles/768_width/public/articles/main-images/Makerspace-main-D73722_129.jpg?itok=k7L22pyv"

   fullUrl := "https://www.unlv.edu" + imgUrl
   r, e := http.Get(fullUrl)
   if e != nil {
      panic(e)
   }

   defer r.Body.Close()

   f, e := os.Create(path.Join("dowloadedImages", "something.jpg"))
   if e != nil {
      panic(e)
   }

   defer f.Close()
   f.ReadFrom(r.Body)
   fmt.Println(r.Body)
/*

   url := "/sites/default/files/styles/1200_width/public/hero_video/image/ComputerScience-HeroImage-D65457-10.PNG?itok=6zYYZOXB"
   splitStr := strings.Split(url, "/")

   lastSplit := splitStr[len(splitStr) - 1]

   locate := strings.Index(lastSplit, ".")

   var fileName string
   if strings.Contains(lastSplit, "jpeg") {
      fileName = lastSplit[0:(locate + 5)]
   } else {
      fileName = lastSplit[0:(locate + 4)]
   }

   fmt.Println(fileName)
   */

}
// package main

// import "fmt"

// func main() {
// 	defer fmt.Println("9")
// 	fmt.Println("0")
// 	defer fmt.Println("8")
// 	fmt.Println("1")
// 	if false {
// 		defer fmt.Println("not reachable")
// 	}
// 	defer func() {
// 		defer fmt.Println("7")
// 		fmt.Println("3")
// 		defer func() {
// 			fmt.Println("5")
// 			fmt.Println("6")
// 		}()
// 		fmt.Println("4")
// 	}()
// 	fmt.Println("2")
// 	return
// 	defer fmt.Println("not reachable")
// }

// package main

// import (
//     "fmt"
// )

// var arrtestStrings [3]string
// var slicetestStrings []string

// func main() {
//     arrtestStrings = [...]string{"apple", "banana", "kiwi"}
//     slicetestStrings = []string{"apple", "banana", "kiwi"}
//     fmt.Println(arrtestStrings)
//     fmt.Println(slicetestStrings)
// }