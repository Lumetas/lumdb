package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

const app_path string = "C:\\lumdb\\"
const path_sepirator string = "\\"
const port string = ":8234"

func path_exist(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else {
		return false
	}
}

func create_dir(path string) {
	if !path_exist(path) {
		os.MkdirAll(path, 0777)
	}
}

func create_file(path string) {
	if !path_exist(path) {
		os.Create(path)
	}
}

func writeStringToFile(filePath string, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func readFileContent(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func init_app(password string) {
	create_dir(app_path)
	if !path_exist(app_path + "lumdb") {
		create_file(app_path + "lumdb")
		writeStringToFile(app_path+"lumdb", password)
	}
}

func deleteFile(filePath string, w http.ResponseWriter) error {
	if !path_exist(filePath) {
		return fmt.Errorf("Файл не существует: %s", filePath)
	}

	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
	return nil
}

func deleteDirectory(dirPath string, w http.ResponseWriter) error {
	if !path_exist(dirPath) {
		return fmt.Errorf("Папка не существует: %s", dirPath)
	}

	err := os.RemoveAll(dirPath)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
	return nil
}

func test_string_substring(str string, substr string) bool {

	matched, err := regexp.MatchString(substr, str)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return false
	}

	return matched
}

func select_write(db string, table string, query string, lines string, index bool, w http.ResponseWriter) bool {
	if !path_exist(app_path + db + path_sepirator + table) {
		return false
	}
	count, err := strconv.Atoi(lines)
	if err == nil {
		file, err := os.Open(app_path + db + path_sepirator + table)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		current := 0
		w.WriteHeader(http.StatusOK)

		for scanner.Scan() {
			if count == 0 {
				return false
			}
			if test_string_substring(scanner.Text(), query) {
				if index {
					fmt.Fprintln(w, current)
				} else {
					fmt.Fprintln(w, scanner.Text())
				}
				count -= 1
			}

			current++
		}

		return true
	} else {
		return false
	}
}

func create_db(name string, w http.ResponseWriter) {
	if path_exist(app_path) {

		if !path_exist(app_path + name) {
			create_dir(app_path + name)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "ok, created")
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "already")
		}

	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "application not initialysed")
	}
}

func add_write(db string, table string, data string, w http.ResponseWriter) error {
	filePath := app_path + db + path_sepirator + table
	if !path_exist(app_path + db + path_sepirator + table) {
		return fmt.Errorf("Ошибка при открытии файла: %v")
	}
	content := data
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Ошибка при открытии файла: %v", err)
	}
	defer file.Close()

	_, err = file.Seek(0, os.SEEK_END)
	if err != nil {
		return fmt.Errorf("Ошибка при перемещении в конец файла: %v", err)
	}

	_, err = file.Write([]byte(content + "\n"))
	if err != nil {
		return fmt.Errorf("Ошибка при записи в файл: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "save data")
	return nil
}

func read_line(path string, target_line int) string {

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	current := 0
	for scanner.Scan() {
		if current == target_line {
			return scanner.Text()
		}
		current++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return "false"
}

func get_write_by_id(db string, table string, id string, w http.ResponseWriter) {
	intid, err := strconv.Atoi(id)
	if err == nil {
		if !path_exist(app_path + db + path_sepirator + table) {
			return
		}
		result := read_line(app_path+db+path_sepirator+table, intid)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, result)
	}
}

func create_table(db string, name string, w http.ResponseWriter) {
	if path_exist(app_path) {
		if path_exist(app_path + db) {
			if !path_exist(app_path + db + path_sepirator + name) {
				create_file(app_path + db + path_sepirator + name)
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, "ok, created")
			} else {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, "already")
			}
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "db not")
		}

	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "application not initialysed")
	}
}

func main() {
	http.HandleFunc("/", handlePost)
	fmt.Println("Starting server on " + port)
	http.ListenAndServe(port, nil)
}

func rewriteLine(filePath, lineNumberStr, newLineContent string, w http.ResponseWriter) error {
	// Проверяем, существует ли файл
	if !path_exist(filePath) {
		return fmt.Errorf("Файл не существует: %s", filePath)
	}

	// Преобразуем номер строки в целое число
	lineNumber, err := strconv.Atoi(lineNumberStr)
	if err != nil {
		return err
	}

	// Открываем файл для чтения
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Создаем временный файл для записи
	tempFile, err := os.CreateTemp("", "temp_file.txt")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Читаем и записываем строки в новый файл, обновляя нужную строку
	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(tempFile)
	currentLine := 0
	for scanner.Scan() {
		if currentLine == lineNumber {
			_, err := writer.WriteString(newLineContent + "\n")
			if err != nil {
				return err
			}
		} else {
			_, err := writer.WriteString(scanner.Text() + "\n")
			if err != nil {
				return err
			}
		}
		currentLine++
	}

	// Проверяем, была ли ошибка при чтении файла
	if err := scanner.Err(); err != nil {
		return err
	}

	// Записываем изменения в исходный файл
	err = writer.Flush()
	if err != nil {
		return err
	}

	// Перезаписываем исходный файл
	err = os.Truncate(filePath, 0)
	if err != nil {
		return err
	}

	_, err = tempFile.Seek(0, 0)
	if err != nil {
		return err
	}

	file, err = os.OpenFile(filePath, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, tempFile)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
	return nil
}

func delete_write(db string, table string, id string, w http.ResponseWriter) error {
	lineToDelete, err := strconv.Atoi(id)
	if err != nil || !path_exist(app_path+db+path_sepirator+table) {
		return fmt.Errorf("Какая-то ошибка")
	}
	filePath := app_path + db + path_sepirator + table
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Создаем временный файл для записи
	tempFile, err := os.CreateTemp("", "temp_file.txt")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Читаем и записываем строки в новый файл, пропуская строку, которую нужно удалить
	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(tempFile)
	currentLine := 0
	for scanner.Scan() {
		if currentLine != lineToDelete {
			_, err := writer.WriteString(scanner.Text() + "\n")
			if err != nil {
				return err
			}
		}
		currentLine++
	}

	// Проверяем, была ли ошибка при чтении файла
	if err := scanner.Err(); err != nil {
		return err
	}

	// Записываем изменения в исходный файл
	err = writer.Flush()
	if err != nil {
		return err
	}

	// Перезаписываем исходный файл
	err = os.Truncate(filePath, 0)
	if err != nil {
		return err
	}

	_, err = tempFile.Seek(0, 0)
	if err != nil {
		return err
	}

	file, err = os.OpenFile(filePath, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, tempFile)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
	return nil
}

func listDirectoryContents(dirPath string, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusOK)

	// Проверяем, существует ли директория
	if !path_exist(dirPath) {
		fmt.Fprintln(w, "not")
		return fmt.Errorf("Директория не существует: %s", dirPath)
	}

	// Получаем список файлов и папок в директории
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	// Выводим содержимое директории
	for _, file := range files {
		if file.IsDir() {
			fmt.Fprintf(w, "%s\n", file.Name())
		} else {
			fmt.Fprintf(w, "%s\n", file.Name())
		}
	}

	return nil
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Считываем данные из формы
		delete_db := r.FormValue("delete_db")
		delete_table := r.FormValue("delete_table")
		show_tables := r.FormValue("show_tables")
		show_databases := r.FormValue("show_databases")
		password := r.FormValue("password")
		create_data_base := r.FormValue("create_data_base")
		create_tab := r.FormValue("create_table")
		init := r.FormValue("init")
		data := r.FormValue("data")
		db := r.FormValue("db")
		table := r.FormValue("table")
		read_by_id := r.FormValue("read_by_id")
		search := r.FormValue("select")
		count := r.FormValue("count")
		index := r.FormValue("index")
		delete_string := r.FormValue("remove")
		edit := r.FormValue("edit")

		if init != "" {
			init_app(init)
		}
		if file, err := readFileContent(app_path + "lumdb"); file != password || err != nil {
			if init == "" {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, "invalid password")
				return
			}
		}

		if db == "lumdb" || delete_db == "lumdb" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "not")
			return
		}

		if show_databases != "" {
			listDirectoryContents(app_path, w)
			return
		}

		if show_tables != "" && db != "" {
			listDirectoryContents(app_path+db, w)
			return
		}

		if delete_db != "" {
			deleteDirectory(app_path+delete_db, w)
			return
		}

		if delete_table != "" && db != "" {
			deleteFile(app_path+db+path_sepirator+delete_table, w)
			return
		}

		if db != "" && table != "" && edit != "" && data != "" {
			rewriteLine(app_path+db+path_sepirator+table, edit, data, w)
			return
		}

		if create_data_base != "" {
			create_db(create_data_base, w)
			return
		}

		if db != "" && create_tab != "" {
			create_table(db, create_tab, w)
			return
		}

		if data != "" && db != "" && table != "" {
			add_write(db, table, data, w)
			return
		}

		if db != "" && table != "" && read_by_id != "" {
			get_write_by_id(db, table, read_by_id, w)
			return
		}

		if db != "" && table != "" && search != "" {
			if index == "true" {
				select_write(db, table, search, count, true, w)
			} else {
				select_write(db, table, search, count, false, w)
			}
			return
		}

		if db != "" && table != "" && delete_string != "" {
			delete_write(db, table, delete_string, w)
			return
		}

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
