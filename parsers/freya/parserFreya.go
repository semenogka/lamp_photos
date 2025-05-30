package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gocolly/colly"
)

type Product struct {
	Name      string `json:"name"`
	Link      string `json:"link"`
}

var products []Product

func main() {
	c := colly.NewCollector()

	
	c.OnHTML(".catalog-grid__item ", func(e *colly.HTMLElement) {
    	link := e.ChildAttr("a", "href")
    	fullLink := e.Request.AbsoluteURL(link)
    	log.Println("Ссылка на товар:", fullLink)

		imgURL := e.ChildAttr("a .catalog-card__picture picture source", "srcset")
		fullImgURL := e.Request.AbsoluteURL(imgURL)
		filename := fmt.Sprintf("FREYA%s", path.Base(imgURL))

		downloadImg(fullImgURL, filename)

		products = append(products, Product{
			Name: filename,
			Link: fullLink,
		})

		

		log.Println("Картинка:", fullImgURL, " ",  "Ссылка на товар:", fullLink, len(products), filename)
	})	

	err := c.Visit("https://freya-light.com/products/dekorativnyy_svet/aksessuar/?SHOWALL_1=1")

	if err != nil {
		log.Fatal(err)
	}
	
	saveToJson("alldata\\data\\aks.json", products)
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
