package handlers

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/isaulin-svg/sprint6/internal/service"
)

func HandlerIndex(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles(filepath.Join(".", "index.html"))
	if err != nil {
		http.Error(w, "Не удалось загрузить страницу: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func HandlerUpload(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Ошибка парсинга формы: "+err.Error(), http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, "Ошибка получения файла: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Ошибка чтения файла: "+err.Error(), http.StatusInternalServerError)
		return
	}

	input := string(data)

	result, err := service.ConvertAutomatically(input)
	if err != nil {
		http.Error(w, "Ошибка конвертации: "+err.Error(), http.StatusInternalServerError)
		return
	}

	outputToClient := result

	if service.IsMorse(result) {
		outputToClient, _ = service.ConvertAutomatically(result)
	}

	ext := filepath.Ext(header.Filename)
	newFileName := fmt.Sprintf("result_%v%v", time.Now().UTC().Unix(), ext)
	newFile, err := os.Create(newFileName)
	if err != nil {
		http.Error(w, "Ошибка создания файла: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer newFile.Close()

	_, err = newFile.Write([]byte(result))
	if err != nil {
		http.Error(w, "Ошибка записи результата: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	fmt.Fprintf(w, "Файл успешно конвертирован!\n\nРезультат:\n%s", outputToClient)
}
