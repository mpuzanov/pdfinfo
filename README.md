# Программа для получения информации о PDF-файле  

*Написана для подсчёта количества страниц в большом количестве переданных для печати файлов.*  

**Аргументы:**   
имя файла или каталог(с подкаталогами), по которому будет поиск всех pdf-файлов

**На выходе:**  
***Формат CSV***  
	/media/sf_H_DRIVE/_Печать;316;13488
***Формат JSON***  
	[{"path":"/media/sf_H_DRIVE/_Печать","file_count":316,"page_count":13488}]

**Примеры:**  
***for Linux*** 

	./pdfinfo -format=csv -i=/media/sf_H_DRIVE/in_file.txt -o=/media/sf_H_DRIVE/out_file.txt
	./pdfinfo < files/in_file_linux.txt > files/out.json 2>files/err.txt
	./pdfinfo files/test.pdf
	./pdfinfo "/media/sf_H_DRIVE/_Печать"
	./pdfinfo "/media/sf_H_DRIVE/_Печать/ивц/2 этап 23дома/10 лет Октября 2/2.pdf"

***for Windows***  

	./pdfinfo test.pdf "H:\_Печать\ивц\2 этап 23дома\10 лет Октября 2\2.pdf"
	./pdfinfo test.pdf "H:\_Печать"
	./pdfinfo "H:\_Печать"
	./pdfinfo "H:\_Печать" > out.json
	./pdfinfo < in_file_win.txt
