# QuipExporter
 
## Running

```
$ cp config-sample.yml config.yml # and make changes
$ go run main.go run -vv
```

## Output

The exporter tool scrapes:
- documents as both HTML and also their exportable types: docx, xlsx, pdf
- conversation threads
- users

It uses the filesystem as a repository and stores everything in two forms:
- `output/archive`: a tree structure that mirrors the structure in Quip
- `output/data`: flat structure of all entities indexed by ID

The `archive` output is intended for "offline" browsing, while the `data` output will be used to build a minimal quip-like app that can also render the comment threads.

## TODO

- [x] Base scraper logic that downloads all entities
- [x] Throttling logic to stay below the rate-limit line
- [x] Export blobs as well
- [ ] Archive browser app
