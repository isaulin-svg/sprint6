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

// HandlerIndex — возвращает HTML-страницу index.html.
func HandlerIndex(w http.ResponseWriter, r *http.Request) {
	// Загружаем index.html из корня проекта
	tmpl, err := template.ParseFiles(filepath.Join(".", "index.html"))
	if err != nil {
		http.Error(w, "Не удалось загрузить страницу: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

// HandlerUpload — обрабатывает загрузку файла и конвертирует содержимое.
func HandlerUpload(w http.ResponseWriter, r *http.Request) {
	// Парсим форму (до 10 МБ)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Ошибка парсинга формы: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем файл из формы
	file, header, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, "Ошибка получения файла: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Читаем данные файла
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Ошибка чтения файла: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Конвертируем содержимое (текст ↔ морзе)
	result, err := service.ConvertAutomatically(string(data))
	if err != nil {
		http.Error(w, "Ошибка конвертации: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Создаём новый файл с результатом
	ext := filepath.Ext(header.Filename)
	newFileName := fmt.Sprintf("result_%v%v", time.Now().UTC().Unix(), ext)
	newFile, err := os.Create(newFileName)
	if err != nil {
		http.Error(w, "Ошибка создания файла: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer newFile.Close()

	// Записываем результат в файл
	_, err = newFile.Write([]byte(result))
	if err != nil {
		http.Error(w, "Ошибка записи результата: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ пользователю
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Файл успешно конвертирован!\n\nРезультат:\n%s", result)
}
