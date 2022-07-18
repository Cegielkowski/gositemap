package urlcrawler

import (
	"gositemap/internal/link"
	"gositemap/internal/workerpool"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

type urlCrawler struct {
	pagesVisitedCount int
	pagesVisited      link.Slice
	hc                HTTPClient
	wp                workerpool.WorkerPool
	wg                sync.WaitGroup
}

type UrlCrawler interface {
	CrawlUrls(url string, maxDepth int) (link.HrefSlice, error)
	PagesVisitedCount() int
}

func NewUrlCrawler(hc HTTPClient, wp workerpool.WorkerPool) UrlCrawler {
	uc := &urlCrawler{
		hc:                hc,
		wp:                wp,
		pagesVisitedCount: 0,
		pagesVisited:      make(link.Slice),
	}

	return uc
}

func (uc *urlCrawler) addPageVisited(page string) {
	uc.pagesVisited[page] = struct{}{}
}

func (uc *urlCrawler) PagesVisitedCount() int {
	return uc.pagesVisitedCount
}

func (uc *urlCrawler) addPageVisitedCount() {
	uc.pagesVisitedCount += 1
}

func (uc *urlCrawler) get(urlStr string, resultC chan<- link.HrefSlice) {
	resp, err := uc.hc.Get(urlStr)

	if err == nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Fatalln(err)
			}
		}(resp.Body)

		uc.addPageVisitedCount()

		hrefSlice, err := link.GetHrefs(resp.Body, getBaseUrl(resp))
		if err != nil {
			log.Fatalln(err)
		}

		resultC <- hrefSlice
	}
}

func getBaseUrl(resp *http.Response) string {
	reqUrl := resp.Request.URL

	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	return baseUrl.String()
}

func (uc *urlCrawler) CrawlUrls(url string, maxDepth int) (link.HrefSlice, error) {
	queue := link.Slice{
		url: struct{}{},
	}

	uc.iterateQueue(queue, maxDepth)

	return uc.pagesVisited.ToSliceSafe(url), nil
}

func (uc *urlCrawler) iterateQueue(queue link.Slice, maxDepth int) {
	resultC := make(chan link.HrefSlice, len(queue))

	for page := range queue {
		uc.wg.Add(1)

		page = strings.TrimSuffix(page, "/")

		uc.wp.AddTask(func() {
			defer uc.wg.Done()
			uc.get(page, resultC)
		})

		uc.addPageVisited(page)
	}

	uc.wg.Wait()

	if nextQueue := uc.createNextQueue(resultC, maxDepth); len(nextQueue) > 0 {
		maxDepth--
		uc.iterateQueue(nextQueue, maxDepth)
	}
}

func (uc *urlCrawler) createNextQueue(resultC <-chan link.HrefSlice, maxDepth int) link.Slice {
	nextQueue := make(link.Slice)
	if maxDepth > 0 {
		for i := 0; i < len(resultC); i++ {
			res := <-resultC
			for _, v := range res {
				if _, ok := uc.pagesVisited[v]; ok || v == "" {
					continue
				}
				nextQueue[v] = struct{}{}
			}
		}
	}

	return nextQueue
}
