package e2e

import (
	"bytes"
	"gositemap/internal/urlcrawler"
	"gositemap/internal/workerpool"
	"gositemap/pkg/gositemap"
	"gositemap/test/mocks"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
)

func Test_appEnv_run(t *testing.T) {
	Client := &mocks.MockClient{}
	r := ioutil.NopCloser(bytes.NewReader([]byte(mocks.PageWithNoHref)))

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
		hc             urlcrawler.HTTPClient
		urlFlag        string
		numWorkers     int
		maxDepth       int
		outputFilePath string
		wp             workerpool.WorkerPool
	}

	var tests = []struct {
		name          string
		fields        fields
		fileToCompare string
		wantErr       bool
	}{
		{
			name: "Simple End to End",
			fields: fields{
				hc:             Client,
				urlFlag:        "https://isaprotege.com.br/",
				numWorkers:     1,
				maxDepth:       1,
				outputFilePath: "./",
			},
			fileToCompare: "../mocks/sitemaps/sitemap.xml",
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &gositemap.AppEnv{
				Hc:             tt.fields.hc,
				UrlFlag:        tt.fields.urlFlag,
				NumWorkers:     tt.fields.numWorkers,
				MaxDepth:       tt.fields.maxDepth,
				OutputFilePath: tt.fields.outputFilePath,
			}
			err := app.Run()

			f1, err2 := ioutil.ReadFile("sitemap.xml")
			if err2 != nil {
				t.Errorf("run() the client have not generated the file")
			}

			f2, _ := ioutil.ReadFile(tt.fileToCompare)

			isEqual := bytes.Equal(f1, f2)

			_ = os.Remove("sitemap.xml")
			if ((err != nil) != tt.wantErr) || !isEqual {
				t.Errorf("run() error = %v, wantErr %v, file is equal: %v", err, tt.wantErr, isEqual)
			}
		})
	}
}
