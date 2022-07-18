package gositemap

import (
	"flag"
	"gositemap/internal/sitemap"
	"gositemap/internal/urlcrawler"
	"gositemap/internal/workerpool"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
)

type AppEnv struct {
	Hc             urlcrawler.HTTPClient
	UrlFlag        string
	NumWorkers     int
	MaxDepth       int
	OutputFilePath string
}

func (app *AppEnv) Run() error {
	wp := workerpool.NewWorkerPool(app.NumWorkers)
	uc := urlcrawler.NewUrlCrawler(app.Hc, wp)
	sm := sitemap.NewSiteMap()

	params := sitemap.Params{
		UrlFlag:        app.UrlFlag,
		OutputFilePath: app.OutputFilePath,
		MaxDepth:       app.MaxDepth,
		Wp:             wp,
		Uc:             uc,
		Sm:             sm,
		Ctx:            nil,
	}
	return sitemap.Run(params)
}

func (app *AppEnv) validate(fl *flag.FlagSet) error {
	if app.MaxDepth < 1 {
		log.Printf("got bad output type, max-depth should be bigger or equal 0 : %d\n", app.MaxDepth)
		fl.Usage()
		return flag.ErrHelp
	}

	if app.NumWorkers < 1 {
		log.Printf("got bad output type, flag parallel has to be greater or equal 1, you requested: %d\n", app.NumWorkers)
		fl.Usage()
		return flag.ErrHelp
	}

	numCpu := runtime.NumCPU()

	if app.NumWorkers > numCpu {
		log.Printf("got bad output type, you requested more workers than you have :(, you requested: %d, you have : %d\n", app.NumWorkers, numCpu)
		fl.Usage()
		return flag.ErrHelp
	}
	_, err := url.ParseRequestURI(app.UrlFlag)

	if err != nil {
		log.Printf("got bad output type, failed to open directory, flag url has to a valid url (including https://), you sent: %q\n", app.UrlFlag)
		fl.Usage()
		return flag.ErrHelp
	}

	dir, err := os.Stat(app.OutputFilePath)
	if err != nil {
		log.Printf("got bad output type, flag output-file has to a valid output-file path (ex: /home/), you sent: %q\n", app.OutputFilePath)
		fl.Usage()
		return flag.ErrHelp
	}
	if !dir.IsDir() {
		log.Printf("got bad output type, flag output-file is not a directory , you sent: %q\n", app.OutputFilePath)
		fl.Usage()
		return flag.ErrHelp
	}

	return nil
}

func (app *AppEnv) fromArgs(args []string) error {
	app.Hc = &http.Client{}
	fl := flag.NewFlagSet("go-site-map", flag.ContinueOnError)
	fl.IntVar(
		&app.NumWorkers, "parallel", 1, "Number of parallel workers to navigate through site",
	)
	fl.StringVar(
		&app.UrlFlag, "url", "https://getaurox.com/", "the url you want to build a sitemap for",
	)
	fl.IntVar(
		&app.MaxDepth, "max-depth", 1, "max depth of url navigation recursion",
	)
	fl.StringVar(
		&app.OutputFilePath, "output-file", "./", "output file path",
	)

	if err := fl.Parse(args); err != nil {
		return err
	}

	return app.validate(fl)
}

func CLI(args []string) int {
	var app AppEnv
	err := app.fromArgs(args)
	if err != nil {
		return 2
	}
	if err = app.Run(); err != nil {
		log.Printf("Runtime error: %v\n", err)
		return 1
	}
	return 0
}
