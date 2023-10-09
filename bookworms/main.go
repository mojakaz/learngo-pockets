package bookworms

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var filePath string
	flag.StringVar(&filePath, "path", "testdata/bookworms.json", "path to the JSON file to load")
	flag.Parse()
	bws, err := loadBookworms(filePath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to load bookworms: %s\n", err)
		os.Exit(1)
	}
	commonBooks := findCommonBooks(bws)
	displayBooks(commonBooks)
}
