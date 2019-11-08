/*
go run . test.pdf
go run . "H:\!Печать"
./pdfinfo "H:\!Печать" > out.txt
*/
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ledongthuc/pdf"
)

func main() {
	prog := filepath.Base(os.Args[0])
	if len(os.Args) <= 1 {
		fmt.Println("Usage:", prog, "\"pdf-file or DIR\"")
		return
	}
	var sumPageCount, fileCount int
	path := os.Args[1]
	// ищем файлы пдф каталоге и в подкаталогах
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {

			if filepath.Ext(path) == ".pdf" {
				PageCount, err := getFilePdf(path) // Read pdf file
				if err != nil {
					return err
				}
				sumPageCount += PageCount
				fileCount++
			}
		}
		return nil
	})
	outstr := fmt.Sprintf("Path: \"%s\", FileCount: %d, PageCount: %d", path, fileCount, sumPageCount)
	fmt.Println(outstr)
}

//getFilePdf возвращает кол-во страниц в файле
func getFilePdf(path string) (int, error) {
	f, r, err := pdf.Open(path)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return 0, err
	}
	return r.NumPage(), nil
}
