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
	"strings"

	"github.com/gocolly/colly"
)

type Product struct {
	Name      string `json:"name"`
	Link      string `json:"link"`
}
var pages int = 1
var products []Product

func main() {
	c := colly.NewCollector()

	
	c.OnHTML(".catalog-card__item ", func(e *colly.HTMLElement) {
    	link := e.ChildAttr("a", "href")
    	fullLink := e.Request.AbsoluteURL(link)
    	log.Println("Ссылка на товар:", fullLink)

		imgURL := e.ChildAttr("a img", "src")
		fullImgURL := e.Request.AbsoluteURL(imgURL)
		filename := fmt.Sprintf("LOFTIT%s", path.Base(imgURL))

		downloadImg(fullImgURL, filename)

		products = append(products, Product{
			Name: filename,
			Link: fullLink,
		})

		

		log.Println("Картинка:", fullImgURL, " ",  "Ссылка на товар:", fullLink, len(products), filename)
	})	

	for i := 1; i <= pages; i++{
		url := "https://loftit.ru/catalog/3/filter/category-is-%D1%83%D0%BB%D0%B8%D1%87%D0%BD%D1%8B%D0%B5%20%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8/in_category-is-%D1%83%D0%BB%D0%B8%D1%87%D0%BD%D1%8B%D0%B5%20%D0%BF%D0%BE%D0%B4%D0%B2%D0%B5%D1%81%D0%BD%D1%8B%D0%B5%20%D1%81%D0%B2%D0%B5%D1%82%D0%B8%D0%BB%D1%8C%D0%BD%D0%B8%D0%BA%D0%B8/apply/?PAGEN_1=%d"
		temp := strings.ReplaceAll(url, "%", "%%")
		temp = strings.Replace(temp, "%%d", "%d", 1)

		page := fmt.Sprintf(temp, i)
		log.Println(page)
		c.Visit(page)
	}

	
	saveToJson("alldata\\data\\podves.json", products)
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
