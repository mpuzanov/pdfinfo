package infopdf

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ledongthuc/pdf"
)

//PDFZap структура одной записи с результатом
type PDFZap struct {
	Path      string `json:"path,omitempty"`       // путь к файлу или к папке
	FileCount int    `json:"file_count,omitempty"` // кол-во файлов
	PageCount int    `json:"page_count,omitempty"` // кол-во страниц
	FileSize  int64  `json:"file_size,omitempty"`  // размер в байтах
}

//PDFResult структура с результатом
type PDFResult []PDFZap

//ToJSON ...
func (o *PDFResult) ToJSON() string {
	b, _ := json.Marshal(o)
	return string(b)
}

//ToCSV ...
func (o *PDFResult) ToCSV() string {
	var s string
	for _, v := range *o {
		s += fmt.Sprintf("%s;%d;%d;%d\n", v.Path, v.FileCount, v.PageCount, v.FileSize)
	}
	return s
}

//ReadPath ...
func ReadPath(path string, coutFiles chan<- int) (out PDFZap, err error) {
	// проверим на существование
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return
	}

	out.Path = path
	// ищем файлы пдф в каталоге и в подкаталогах
	err = filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !f.IsDir() { // если файл

			if filepath.Ext(path) == ".pdf" { // проверяем, что pdf
				PageCount, err := GetPageCountPdf(path)
				if err != nil {
					return err
				}
				out.FileSize += f.Size()
				out.PageCount += PageCount
				out.FileCount++

				coutFiles <- 1 // отправляем в канал подсчёта кол-ва файлов
			}
		}
		return nil
	})
	if err != nil {
		return out, err
	}
	return out, nil
}

//GetPageCountPdf возвращает кол-во страниц в файле
func GetPageCountPdf(fileName string) (PageCount int, err error) {
	f, r, err := pdf.Open(fileName)
	if err != nil {
		err = fmt.Errorf("error open file %s. %w", fileName, err)
		return
	}
	defer f.Close()
	PageCount = r.NumPage()
	return
}
