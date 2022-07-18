# Golang Sitemap Builder

### Usage
- Run the program using `go run main.go <optional flags>`, you will find a new file, name: sitemap.xml
- Optional flags are these:
  - `-url=<your domain of choice>` The default is `"https://getaurox.com/"`.
  - `-parallel=<your number of choice (>= 1)>` The default is `1`.
  - `-max-depth=<your number of choice (>= 1)>` The default is `1`.
  - `-output-file=<your directory of choice>` The default is `"./"`.

**Note:** Giving a depth of 2+ can cause a long run time
