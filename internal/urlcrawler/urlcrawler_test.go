package urlcrawler

import (
	"bytes"
	"gositemap/internal/link"
	"gositemap/internal/workerpool"
	"gositemap/test/mocks"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"sync"
	"testing"
)

func TestNewUrlCrawler(t *testing.T) {
	wp := workerpool.NewWorkerPool(1)
	hc := http.Client{}

	type args struct {
		hc HTTPClient
		wp workerpool.WorkerPool
	}
	tests := []struct {
		name string
		args args
		want UrlCrawler
	}{
		{
			name: "Should be equal",
			args: args{
				hc: &hc,
				wp: wp,
			},
			want: NewUrlCrawler(&hc, wp),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUrlCrawler(tt.args.hc, tt.args.wp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUrlCrawler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getBaseUrl(t *testing.T) {
	type args struct {
		resp *http.Response
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should get the baseurl",
			args: args{
				resp: &http.Response{
					StatusCode: 200,
					Request: &http.Request{URL: &url.URL{
						Scheme: "https",
						Host:   "getaurox.com",
					}},
				},
			},
			want: "https://getaurox.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getBaseUrl(tt.args.resp); got != tt.want {
				t.Errorf("getBaseUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_urlCrawler_CrawlUrls(t *testing.T) {
	Client := &mocks.MockClient{}
	r := ioutil.NopCloser(bytes.NewReader([]byte(mocks.PageWithNoHref)))
	wp := workerpool.NewWorkerPool(1)
	wp.Run()
	mocks.GetDoFunc = func(string) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
			Request: &http.Request{URL: &url.URL{
				Scheme: "text",
				Host:   "https://isaprotege.com.br/",
			}},
		}, nil
	}

	type fields struct {
		hc HTTPClient
		wp workerpool.WorkerPool
	}
	type args struct {
		url      string
		maxDepth int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    link.HrefSlice
		wantErr bool
	}{
		{
			name: "Should get the urls",
			fields: fields{
				hc: Client,
				wp: wp,
			},
			args: args{
				url:      "https://isaprotege.com.br/",
				maxDepth: 1,
			},
			want:    link.HrefSlice{"https://isaprotege.com.br"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &urlCrawler{
				hc:                tt.fields.hc,
				wp:                tt.fields.wp,
				pagesVisitedCount: 0,
				pagesVisited:      make(link.Slice),
			}

			got, err := uc.CrawlUrls(tt.args.url, tt.args.maxDepth)
			if (err != nil) != tt.wantErr {
				t.Errorf("CrawlUrls() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CrawlUrls() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_urlCrawler_PagesVisitedCount(t *testing.T) {
	type fields struct {
		pagesVisitedCount int
		pagesVisited      link.Slice
		hc                HTTPClient
		wp                workerpool.WorkerPool
		wg                sync.WaitGroup
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Should be zero",
			fields: fields{
				pagesVisitedCount: 0,
				pagesVisited:      nil,
				hc:                nil,
				wp:                nil,
				wg:                sync.WaitGroup{},
			},
			want: 0,
		},
		{
			name: "Should be 2",
			fields: fields{
				pagesVisitedCount: 2,
				pagesVisited:      nil,
				hc:                nil,
				wp:                nil,
				wg:                sync.WaitGroup{},
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &urlCrawler{
				pagesVisitedCount: tt.fields.pagesVisitedCount,
				pagesVisited:      tt.fields.pagesVisited,
				hc:                tt.fields.hc,
				wp:                tt.fields.wp,
				wg:                tt.fields.wg,
			}
			if got := uc.PagesVisitedCount(); got != tt.want {
				t.Errorf("PagesVisitedCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_urlCrawler_addPageVisited(t *testing.T) {
	type fields struct {
		pagesVisitedCount int
		pagesVisited      link.Slice
		hc                HTTPClient
		wp                workerpool.WorkerPool
		wg                sync.WaitGroup
	}
	type args struct {
		page string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Should add",
			fields: fields{
				pagesVisitedCount: 0,
				pagesVisited:      make(link.Slice),
				hc:                nil,
				wp:                nil,
				wg:                sync.WaitGroup{},
			},
			args: args{
				page: "https://google.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &urlCrawler{
				pagesVisitedCount: tt.fields.pagesVisitedCount,
				pagesVisited:      tt.fields.pagesVisited,
				hc:                tt.fields.hc,
				wp:                tt.fields.wp,
				wg:                tt.fields.wg,
			}
			uc.addPageVisited(tt.args.page)
			if _, ok := uc.pagesVisited[tt.args.page]; !ok {
				t.Errorf("addPageVisited() not added")
			}
		})
	}
}

func Test_urlCrawler_addPageVisitedCount(t *testing.T) {
	type fields struct {
		pagesVisitedCount int
		pagesVisited      link.Slice
		hc                HTTPClient
		wp                workerpool.WorkerPool
		wg                sync.WaitGroup
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Should be 3",
			fields: fields{
				pagesVisitedCount: 2,
				pagesVisited:      nil,
				hc:                nil,
				wp:                nil,
				wg:                sync.WaitGroup{},
			},
			want: 3,
		},
		{
			name: "Should be 344",
			fields: fields{
				pagesVisitedCount: 343,
				pagesVisited:      nil,
				hc:                nil,
				wp:                nil,
				wg:                sync.WaitGroup{},
			},
			want: 344,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &urlCrawler{
				pagesVisitedCount: tt.fields.pagesVisitedCount,
				pagesVisited:      tt.fields.pagesVisited,
				hc:                tt.fields.hc,
				wp:                tt.fields.wp,
				wg:                tt.fields.wg,
			}
			uc.addPageVisitedCount()

			if uc.pagesVisitedCount != tt.want {
				t.Errorf("addPageVisitedCount() = %v, want %v", uc.pagesVisitedCount, tt.want)
			}
		})
	}
}

func Test_urlCrawler_createNextQueue(t *testing.T) {
	channel := make(chan link.HrefSlice, 1)
	channel <- link.HrefSlice{"https://google.com", "https://facebook.com"}
	type fields struct {
		pagesVisitedCount int
		pagesVisited      link.Slice
		hc                HTTPClient
		wp                workerpool.WorkerPool
		wg                sync.WaitGroup
	}
	type args struct {
		resultC  <-chan link.HrefSlice
		maxDepth int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   link.Slice
	}{
		{
			name: "",
			fields: fields{
				pagesVisitedCount: 0,
				pagesVisited:      nil,
				hc:                nil,
				wp:                nil,
				wg:                sync.WaitGroup{},
			},
			args: args{
				resultC:  channel,
				maxDepth: 1,
			},
			want: link.Slice{"https://google.com": struct{}{}, "https://facebook.com": struct{}{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &urlCrawler{
				pagesVisitedCount: tt.fields.pagesVisitedCount,
				pagesVisited:      tt.fields.pagesVisited,
				hc:                tt.fields.hc,
				wp:                tt.fields.wp,
				wg:                tt.fields.wg,
			}
			if got := uc.createNextQueue(tt.args.resultC, tt.args.maxDepth); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createNextQueue() = %v, want %v", got, tt.want)
			}
		})
	}
}
