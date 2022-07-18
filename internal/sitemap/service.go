package sitemap

import (
	"encoding/xml"
	"golang.org/x/net/context"
	"gositemap/internal/link"
	"gositemap/internal/urlcrawler"
	"gositemap/internal/workerpool"
	"io"
	"log"
	"os"
	"time"
)

type Loc struct {
	Value string `xml:"loc"`
}

type Params struct {
	UrlFlag        string
	OutputFilePath string
	MaxDepth       int
	Wp             workerpool.WorkerPool
	Uc             urlcrawler.UrlCrawler
	Sm             SiteMap
	Ctx            context.Context
}

type SiteMap interface {
	MarshalXML(file io.Writer) error
	CreateXML(outputFilePath string) error
	AddPages(pages link.HrefSlice)
}

type siteMap struct {
	Xmlns string `xml:"xmlns"`
	Urls  []Loc  `xml:"url"`
}

func NewSiteMap() SiteMap {
	sm := &siteMap{
		Xmlns: "https://www.sitemaps.org/schemas/sitemap/0.9",
	}
	return sm
}

func (sm *siteMap) AddPages(pages link.HrefSlice) {
	sm.Urls = make([]Loc, 0, len(pages))

	for _, page := range pages {
		sm.Urls = append(sm.Urls, Loc{Value: page})
	}
}

func (sm siteMap) MarshalXML(file io.Writer) error {
	encoder := xml.NewEncoder(file)
	encoder.Indent("", "    ")

	if err := encoder.Encode(sm); err != nil {
		return err
	}

	return nil
}

func (sm siteMap) CreateXML(outputFilePath string) error {
	file, err := os.Create(outputFilePath + "sitemap.xml")
	if err != nil {
		return err
	}
	_, err = file.WriteString(xml.Header)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(file)
	err = sm.MarshalXML(file)
	if err != nil {
		return err
	}
	return nil
}

func Run(params Params) error {
	start := time.Now()
	params.Wp.Run()
	pages, err := params.Uc.CrawlUrls(params.UrlFlag, params.MaxDepth)
	if err != nil {
		return err
	}
	duration := time.Since(start)

	params.Sm.AddPages(pages)
	err = params.Sm.CreateXML(params.OutputFilePath)
	if err != nil {
		return err
	}
	log.Printf("Time duration: %f seconds\n", duration.Seconds())
	log.Printf("Count of links visited: %d\n", params.Uc.PagesVisitedCount())

	// @todo ctx everywhere

	return nil
}
