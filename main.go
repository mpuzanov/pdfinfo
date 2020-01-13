/*

 */
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	pi "pdfinfo/infopdf"
	"strings"
	"sync"
	"time"
)

var (
	format       = flag.String("format", "json", "format output: json or csv")
	resultToFile = flag.String("o", "", "FileName result")
	fileNameIn   = flag.String("i", "", "FileName input list pdf-files(or dirs)")
	formats      = map[string]string{"json": ".json", "csv": ".csv"}
)

func main() {
	flag.Parse()
	flag.Args()

	if _, exists := formats[*format]; !exists {
		fmt.Printf("format output <%s> unknown", *format)
		os.Exit(1)
	}

	var pdfinfo pi.PDFResult
	pdfinfo = make(pi.PDFResult, 0)
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
	countFiles := make(chan int)
	var n sync.WaitGroup

	for _, file := range files {
		n.Add(1)
		go func(file string) {
			defer n.Done()
			p, err := pi.ReadPath(file, countFiles)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			} else {
				pdfinfo = append(pdfinfo, p)
			}
		}(file)
	}
	go func() {
		n.Wait()
		close(countFiles)
	}()

	var showProgres bool
	if *resultToFile == "" {
		showProgres = true
	}
	var tick <-chan time.Time // Print the results periodically.
	if showProgres {
		tick = time.Tick(500 * time.Millisecond)
	}
	var nfiles int64

	func() {
		for {
			select {
			case _, ok := <-countFiles:
				if !ok {
					return // канал закрыт - выходим
				}
				nfiles++
			case <-tick:
				if showProgres {
					fmt.Printf("Обработано: %d\r", nfiles) // выводим на экран счётчик обработки файлов
				}
			}
		}
	}()

	err := printPdfInfo(pdfinfo, *format, *resultToFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

}

//readFileFromInput Заполняем список файлов files для обработки из входного потока или файла
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

//printPdfInfo Сохранение результатов в файл или на экран
func printPdfInfo(pdfinfo pi.PDFResult, format string, resultToFile string) error {
	var text string
	switch format {
	case "json":
		text = fmt.Sprint(pdfinfo.ToJSON())
	case "csv":
		text = fmt.Sprint(pdfinfo.ToCSV())
	}

	if resultToFile == "" {
		fmt.Println(text)
	} else {
		f, err := os.Create(resultToFile) // всегда создаём новый файл
		if err != nil {
			return err
		}
		defer f.Close()
		fmt.Fprintf(f, "%s", text)
		fmt.Printf("Результат записан в файл: %s\n", resultToFile)
	}
	return nil
}
