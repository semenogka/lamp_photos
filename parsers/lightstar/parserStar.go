package main

import (
	"encoding/json"
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

	
	c.OnHTML(".grid__item", func(e *colly.HTMLElement) {
    	link := e.ChildAttr("header .card__title a ", "href")
    	fullLink := e.Request.AbsoluteURL(link)
    	log.Println("Ссылка на товар:", fullLink)

		imgURL := e.ChildAttr("header .card__image-wrapper img", "src")
		fullImgURL := e.Request.AbsoluteURL(imgURL)
		filename := path.Base(imgURL)

		downloadImg(fullImgURL, filename)

		products = append(products, Product{
			Name: filename,
			Link: fullLink,
		})

		

		log.Println("Картинка:", fullImgURL, " ",  "Ссылка на товар:", fullLink, filename, len(products))
	})	

	err := c.Visit("https://lightstar.ru/kupit-lustry-svetilniki-optom/potolochnye-svetilniki/?limit=all")

	if err != nil {
		log.Fatal(err)
	}
	
	saveToJson("alldata\\data\\potol.json", products)
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
