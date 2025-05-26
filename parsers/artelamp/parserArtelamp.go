package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gocolly/colly"
)

type Product struct {
	Name      string `json:"name"`
	Link      string `json:"link"`
}

var products []Product

var pages int = 8

func main() {
	c := colly.NewCollector()
	c.OnHTML(".listing_content_catalog_units", func(e *colly.HTMLElement) {
		e.ForEach(".unit", func(i int, h *colly.HTMLElement) {
			link := h.ChildAttr(".element .imgwr a", "href")
			fullLink := h.Request.AbsoluteURL(link)

			imgURL := h.ChildAttr(".element .imgwr a img", "data-src")
			fullImgURL := h.Request.AbsoluteURL(imgURL)
			
			filename := fmt.Sprintf("ARTELAMP%s", extractName(fullImgURL))

			downloadImg(fullImgURL, filename)

			products = append(products, Product{
				Name: filename,
				Link: fullLink,
			})

			log.Println("Картинка:", fullImgURL, " ",  "Ссылка на товар:", fullLink, len(products), " ", filename)
		})
				
	})	
	for i := 1; i <= pages; i++ {
		
		page := fmt.Sprintf("https://artelamp.ru/catalog/magnitnyie-trekovyie-sistemyi/magnitnyie-trekovyie-svetilniki/page_%d", i)
		err := c.Visit(page)

		if err != nil {
			log.Fatal(err)
		}
	}
	
	
	saveToJson("alldata\\data\\magn.json", products)
}

func downloadImg(src string, name string) {
	res, err := http.Get(src)

	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	filepath := filepath.Join("allimgs", name)

	out, err := os.Create(filepath)
	if err != nil {
		log.Println(err)
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		log.Println(err)
	}
	
}

func saveToJson(filename string, data interface{}) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println("Ошибка открытия файла:", err)
		return
	}
	defer file.Close()

	var exData []interface{}
	
	stat, _ := file.Stat()
	if stat.Size() != 0 {
		json.NewDecoder(file).Decode(&exData)
	}

	exData = append(exData, data)


	if _, err := file.Seek(0, 0); err != nil {
		log.Println("Ошибка seek:", err)
		return
	}
	if err := file.Truncate(0); err != nil {
		log.Println("Ошибка truncate:", err)
		return
	}


	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")

	err = encoder.Encode(exData)

	if err != nil {
		log.Println(err)
	}

}

func extractName(url string) string {
	count := 0
	name := []rune{}
	for _, let := range url {
		if let == '/' {
			count += 1
		}

		if count == 6 && let != '/'{
			name = append(name, let)
		}
	}
	filename := fmt.Sprintf("%s.jpeg", string(name))
	return filename
}
