package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"github.com/xuri/excelize"
)

func main() {
	// create new excel file start
	f := excelize.NewFile()
	index := f.NewSheet("linksSheet")
	f.SetActiveSheet(index)


	f.SetCellValue("linksSheet", "A1", "links")
	// create new excel file end

	for _, url := range os.Args[1:] {
		links, err := findLinks(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse: %v\n", err)
		}
		for i, link := range links {
			fmt.Println(fmt.Sprintf("=======" ))
			f.SetCellValue("linksSheet", fmt.Sprintf("A%d", i+1), link)
			fmt.Println(fmt.Sprintf("%d  : %s", i, link ))
		}
	}

	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}

func findLinks(url string) ([]string, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, res.Status)
	}

	doc, err := html.Parse(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	return visit(nil, doc), nil
}

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}

	return links
}

// создаём эксель лист
// func createSheet(links []string)  {
// 	f := excelize.NewFile()

// 	index := f.NewSheet("Sheet1")

// 	f.SetActiveSheet(index)

// 	if err := f.SaveAs("Book1.xlsx"); err != nil {
// 		fmt.Println(err)
// 	}

// 	for i, link := range links {
// 		f.SetCellValue("Sheet1", fmt.Sprintf("A%v", i+1), link)
// 	}
// }