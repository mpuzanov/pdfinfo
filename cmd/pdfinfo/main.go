/*

 */
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ledongthuc/pdf"
)

//PDFZap структура одной записи с результатом
type PDFZap struct {
	Path      string `json:"path"`
	FileCount int    `json:"file_count"`
	PageCount int    `json:"page_count"`
}

//PDFResult структура с результатом
type PDFResult []PDFZap

func (o *PDFResult) toJSON() string {
	b, _ := json.Marshal(o)
	return string(b)
}

func (o *PDFResult) toCSV() string {
	var s string
	for _, v := range *o {
		s += fmt.Sprintf("%s;%d;%d\n", v.Path, v.FileCount, v.PageCount)
	}
	return s
}

var (
	format      = flag.String("format", "json", "format output: json or csv")
	fileNameOut = flag.String("o", "", "FileName result")
	fileNameIn  = flag.String("i", "", "FileName input list pdf-files(or dirs)")
	formats     = map[string]string{"json": ".json", "csv": ".csv"}
)

func main() {
	flag.Parse()
	flag.Args()

	if _, exists := formats[*format]; !exists {
		fmt.Printf("format output <%s> unknown", *format)
		os.Exit(1)
	}

	var pdfinfo PDFResult
	pdfinfo = make(PDFResult, 0)
	files := flag.Args()

	if len(files) == 0 {

		if *fileNameIn != "" {
			file, err := os.Open(*fileNameIn)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				return
			}
			defer file.Close()
			if err := readFileFromInput(file, &files); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				return
			}
		} else {
			StdinInfo, _ := os.Stdin.Stat()
			if StdinInfo.Size() == 0 {
				prog := filepath.Base(os.Args[0])
				fmt.Println("Usage:", prog, "<pdf-file(dir)> ...")
				return
			}
			err := readFileFromInput(os.Stdin, &files)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				return
			}
		}
	}

	//fmt.Println(files)

	for _, file := range files {
		p, err := readPath(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		} else {
			pdfinfo = append(pdfinfo, p)
		}
	}
	var text string
	switch *format {
	case "json":
		text = fmt.Sprint(pdfinfo.toJSON())
	case "csv":
		text = fmt.Sprint(pdfinfo.toCSV())
	}

	if *fileNameOut == "" {
		fmt.Println(text)
	} else {
		//fmt.Println("file name:", *fileNameOut)
		//newFileName := strings.TrimSuffix(fileName, path.Ext(fileName)) + formats[*format]
		f, err := os.Create(*fileNameOut) // всегда создаём новый файл
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		fmt.Fprintf(f, "%s", text)
		fmt.Printf("Создан файл: %s\n", *fileNameOut)
	}
}

func readPath(path string) (out PDFZap, err error) {
	// проверим на существование
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return
	}
	out.Path = path
	// ищем файлы пдф в каталоге и в подкаталогах
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !f.IsDir() { // если файл

			if filepath.Ext(path) == ".pdf" { // проверяем, что pdf
				PageCount, err := getPageCountPdf(path)
				if err != nil {
					return err
				}
				out.PageCount += PageCount
				out.FileCount++
				if *fileNameOut == "" {
					fmt.Printf("%d\r", out.FileCount) // выводим на экран счётчик обработки файлов
				}
			}
		}
		return nil
	})
	return out, nil
}

//getFilePdf возвращает кол-во страниц в файле
func getPageCountPdf(fileName string) (PageCount int, err error) {
	f, r, err := pdf.Open(fileName)
	if err != nil {
		err = fmt.Errorf("error open file %s. %w", fileName, err)
		return
	}
	defer f.Close()
	PageCount = r.NumPage()
	return
}

//readFileFromInput
//Заполняем список файлов files для обработки из входного потока или файла
func readFileFromInput(f *os.File, files *[]string) (err error) {

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if len(text) == 0 {
			break
		}
		*files = append(*files, text)
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	return
}
