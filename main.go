package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var conn_db_stands string = "root:111@tcp(localhost:3306)/stands"
var showThemesGlob string

type Demo struct {
	Id            int
	Namesys       string
	Datecreate    string
	Dateupdate    string
	Datecheck     string
	Path          string
	Nameprop      string
	Actual        int
	Places        string
	Numberpp      int
	Check_success string
}
type Namesyses struct {
	Id      int
	Namesys string
}
type Themes struct {
	Id    int
	Theme string
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		return
	}
}
func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/tabs.html", "templates/index.html")
	if err != nil {
		log.Print(err)
	}
	vars0 := mux.Vars(r)
	vars1 := vars0["id"]
	fmt.Println(vars1)

	t.ExecuteTemplate(w, "index", vars1)
}
func listcheck(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/listcheck.html", "templates/tabs.html")
	if err != nil {
		log.Print(err)
	}
	demo_db, err := sql.Open("mysql", conn_db_stands)
	if err != nil {
		log.Println(err)
	}
	defer demo_db.Close()

	rows, err := demo_db.QueryContext(context.Background(), fmt.Sprintf("select id, namesys, datecreate, dateupdate, datecheck, path, nameprop, actual, places, numberpp from demosys1 where datecheck <= dateupdate or actual = 0 order by id desc"), []any{}...)
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	var demos []Demo
	for rows.Next() {
		var demo Demo
		if err := rows.Scan(&demo.Id, &demo.Namesys, &demo.Datecreate,
			&demo.Dateupdate, &demo.Datecheck, &demo.Path, &demo.Nameprop,
			&demo.Actual, &demo.Places, &demo.Numberpp); err != nil {
			log.Println(err)
		}
		demos = append(demos, demo)
	}
	namesyses_db, err := sql.Open("mysql", conn_db_stands)
	if err != nil {
		log.Println(err)
	}
	defer namesyses_db.Close()

	namesyses_rows, err := namesyses_db.QueryContext(context.Background(), fmt.Sprintf("select * from namesyses"), []any{}...)
	if err != nil {
		log.Println(err)
	}
	defer namesyses_rows.Close()

	var namesyses []Namesyses
	for namesyses_rows.Next() {
		var namesyse Namesyses
		if err := namesyses_rows.Scan(&namesyse.Id, &namesyse.Namesys); err != nil {
			log.Println(err)
		}
		namesyses = append(namesyses, namesyse)
	}

	// Чтение из файла
	data_out, err := ioutil.ReadFile("assets/txt/example.txt")
	if err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return
	}

	// Печать содержимого файла
	fmt.Println("Содержимое файла:")
	fmt.Println(string(data_out))

	essenses := struct {
		Demos     []Demo
		Namesyses []Namesyses
		Data_out  []byte
	}{
		Demos:     demos,
		Namesyses: namesyses,
		Data_out:  data_out,
	}

	t.ExecuteTemplate(w, "listcheck", essenses)
}
func create(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/create.html", "templates/tabs.html"))
	demo_db, err := sql.Open("mysql", conn_db_stands)
	if err != nil {
		log.Println(err)
	}
	defer demo_db.Close()

	rows, err := demo_db.QueryContext(context.Background(), fmt.Sprintf("select id, namesys, datecreate, dateupdate, datecheck, path, nameprop, actual, places, numberpp from demosys1 order by id desc"), []any{}...)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var demos []Demo
	for rows.Next() {
		var demo Demo
		if err := rows.Scan(&demo.Id, &demo.Namesys, &demo.Datecreate,
			&demo.Dateupdate, &demo.Datecheck, &demo.Path, &demo.Nameprop,
			&demo.Actual, &demo.Places, &demo.Numberpp); err != nil {
			log.Println(err)
		}
		demos = append(demos, demo)
	}
	namesyses_db, err := sql.Open("mysql", conn_db_stands)
	if err != nil {
		log.Println(err)
	}
	defer namesyses_db.Close()

	namesyses_rows, err := namesyses_db.QueryContext(context.Background(), fmt.Sprintf("select * from namesyses"), []any{}...)
	if err != nil {
		log.Println(err)
	}
	defer namesyses_rows.Close()

	var namesyses []Namesyses
	for namesyses_rows.Next() {
		var namesyse Namesyses
		if err := namesyses_rows.Scan(&namesyse.Id, &namesyse.Namesys); err != nil {
			log.Println(err)
		}
		namesyses = append(namesyses, namesyse)
	}

	essenses := struct {
		Demos     []Demo
		Namesyses []Namesyses
	}{
		Demos:     demos,
		Namesyses: namesyses,
	}

	t.Execute(w, essenses)
}
func createAction(w http.ResponseWriter, r *http.Request) {

	now := time.Now()
	datecreate := now.Format("2006-01-02 15:04:05")
	dateupdate := datecreate
	datecheck := "2025-03-15"

	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(1001)

	namesys := r.FormValue("namesys")
	nameprop := r.FormValue("nameprop")
	path := randomNum + 1 // Используйте одно случайное число
	actual1 := r.FormValue("actual1")
	places := r.FormValue("places")

	count0 := r.FormValue("count")
	count, err := strconv.Atoi(count0)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid count value", http.StatusBadRequest)
		return
	}
	numberpp := count + 1

	db, err := sql.Open("mysql", conn_db_stands)
	checkErr(err)
	defer db.Close()

	// Подготовка параметров для пакетного добавления
	for i := 0; i < count; i++ {
		numberpp--
		_, err = db.Exec("INSERT INTO demosys1 (namesys, datecreate, dateupdate, datecheck, path, nameprop, actual, places, numberpp) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
			namesys, datecreate, dateupdate, datecheck, path, nameprop, actual1, places, numberpp)
		checkErr(err) // Проверяем ошибку после каждого выполнения
	}

	fmt.Println("Сработало добавление строк")
	fmt.Println("Даты:", datecreate, dateupdate)

	// Открытие файла для добавления (если файла нет, он будет создан)
	file, err := os.OpenFile("assets/txt/example.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close() // Закрытие файла в конце работы

	// Запись новой строки в файл
	dateTimeStr := now.Format("2006-01-02 15:04:05")
	str := "::: Создание темы. Название системы: " + namesys
	перевод_каретки := "\n"

	_, err = file.WriteString(dateTimeStr + str + перевод_каретки)
	if err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
		return
	}
	fmt.Println("Данные успешно добавлены в файл.")

	http.Redirect(w, r, "/create", http.StatusSeeOther)
}
func createNamesyses(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/createnamesyses.html", "templates/tabs.html")
	if err != nil {
		log.Print(err)
	}

	namesyses_db, err := sql.Open("mysql", conn_db_stands)
	if err != nil {
		log.Println(err)
	}
	defer namesyses_db.Close()

	namesyses_rows, err := namesyses_db.QueryContext(context.Background(), fmt.Sprintf("select * from namesyses"), []any{}...)
	if err != nil {
		log.Println(err)
	}
	defer namesyses_rows.Close()

	var namesyses []Namesyses
	for namesyses_rows.Next() {
		var namesyse Namesyses
		if err := namesyses_rows.Scan(&namesyse.Id, &namesyse.Namesys); err != nil {
			log.Println(err)
		}
		namesyses = append(namesyses, namesyse)
	}

	essenses := struct {
		Namesyses []Namesyses
	}{

		Namesyses: namesyses,
	}

	t.ExecuteTemplate(w, "createnamesyses", essenses)
}
func createnamesysesAction(w http.ResponseWriter, r *http.Request) {
	namesys := r.FormValue("namesys")

	db, err := sql.Open("mysql", conn_db_stands)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Установка данных

	// Пример параметризованного запроса INSERT, который защищает от SQL инъекций

	_, err = db.Exec("insert into namesyses (namesys) values (?)", namesys)

	if err != nil {
		panic(err)
	}

	fmt.Println("Сработало добавление одной строки названия демосистемы")

	http.Redirect(w, r, "/createnamesyses", http.StatusSeeOther)
}
func createThemes(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/createthemes.html", "templates/tabs.html")
	if err != nil {
		log.Print(err)
	}

	themes_db, err := sql.Open("mysql", conn_db_stands)
	if err != nil {
		log.Println(err)
	}
	defer themes_db.Close()

	themes_rows, err := themes_db.QueryContext(context.Background(), fmt.Sprintf("select * from themes"), []any{}...)
	if err != nil {
		log.Println(err)
	}
	defer themes_rows.Close()

	var themes []Themes
	for themes_rows.Next() {
		var theme Themes
		if err := themes_rows.Scan(&theme.Id, &theme.Theme); err != nil {
			log.Println(err)
		}
		themes = append(themes, theme)
	}

	essenses := struct {
		Themes []Themes
	}{

		Themes: themes,
	}

	t.ExecuteTemplate(w, "createthemes", essenses)
}
func createthemesAction(w http.ResponseWriter, r *http.Request) {
	themes := r.FormValue("themes")

	db, err := sql.Open("mysql", conn_db_stands)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Установка данных

	// Пример параметризованного запроса INSERT, который защищает от SQL инъекций

	_, err = db.Exec("insert into themes (theme) values (?)", themes)

	if err != nil {
		panic(err)
	}

	fmt.Println("Сработало добавление одной строки названия темы")

	http.Redirect(w, r, "/createthemes", http.StatusSeeOther)
}
func showNamesys(w http.ResponseWriter, r *http.Request) {
	log.Print("showNamesys:")
	vars := mux.Vars(r)
	t, err := template.ParseFiles("templates/shownamesys.html", "templates/tabs.html")
	if err != nil {
		log.Println(err)
	}
	db, err := sql.Open("mysql", conn_db_stands)
	if err != nil {
		panic(err)
	}

	res, err := func() (*sql.Rows, error) {
		var query string = fmt.Sprintf("select * from `namesyses` where `id` = '%s'", vars["id"])
		return db.QueryContext(context.Background(), query, []any{}...)
	}()
	if err != nil {
		panic(err)
	}

	showNamesys := Namesyses{

		Id:      0,
		Namesys: "",
	}
	for res.Next() {
		var post22 Namesyses
		err = res.Scan(&post22.Id, &post22.Namesys)
		if err != nil {
			panic(err)
		}
		showNamesys = post22
		defer res.Close()
	}

	t.ExecuteTemplate(w, "shownamesys", showNamesys)
}
func showThemes(w http.ResponseWriter, r *http.Request) {
	log.Print("showThemes:")
	vars := mux.Vars(r)
	t, err := template.ParseFiles("templates/showthemes.html", "templates/tabs.html")
	if err != nil {
		log.Println(err)
	}
	db, err := sql.Open("mysql", conn_db_stands)
	if err != nil {
		panic(err)
	}

	res, err := func() (*sql.Rows, error) {
		var query string = fmt.Sprintf("select * from `themes` where `id` = '%s'", vars["id"])
		return db.QueryContext(context.Background(), query, []any{}...)
	}()
	if err != nil {
		panic(err)
	}

	showThemes := Themes{

		Id:    0,
		Theme: "",
	}
	for res.Next() {
		var post22 Themes
		err = res.Scan(&post22.Id, &post22.Theme)
		if err != nil {
			panic(err)
		}
		showThemes = post22
		fmt.Println("showThemes:", showThemes)
		showThemesGlob = showThemes.Theme
		fmt.Println("showThemesGlob:", showThemesGlob)
		defer res.Close()
	}

	t.ExecuteTemplate(w, "showthemes", showThemes)
}
func editnamesysAction(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	id := r.FormValue("id")
	namesys := r.FormValue("namesys")
	db3, err := sql.Open("mysql", conn_db_stands)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer db3.Close()

	result, err := db3.Exec("update namesyses set `namesys`=? where `id`=?",
		namesys, id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Println("Сработало изменение данных в функции EditHandler и изменено записей: ", rows)
	http.Redirect(w, r, "/createnamesyses", http.StatusSeeOther)
}
func editthemesAction(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	id := r.FormValue("id")
	themes := r.FormValue("themes")
	db3, err := sql.Open("mysql", conn_db_stands)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer db3.Close()

	result, err := db3.Exec("update themes set `theme`=? where `id`=?",
		themes, id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Println("Сработало изменение данных в функции editthemesAction и изменено записей: ", rows)
	http.Redirect(w, r, "/createthemes", http.StatusSeeOther)
}
func deletenamesysAction(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	id := r.FormValue("id")

	db3, err := sql.Open("mysql", conn_db_stands)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer db3.Close()

	_, err = db3.Exec("delete from namesyses where id=?", id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	namesys := ""
	logfiles(id, namesys)

	// rows, err := result.RowsAffected()
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, err)
	// 	return
	// }

	// fmt.Println("Сработало изменение данных в функции EditHandler и изменено записей: ", rows)
	http.Redirect(w, r, "/createnamesyses", http.StatusSeeOther)
}
func deletethemesAction(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	id := r.FormValue("id")

	db3, err := sql.Open("mysql", conn_db_stands)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer db3.Close()

	_, err = db3.Exec("delete from themes where id=?", id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// rows, err := result.RowsAffected()
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, err)
	// 	return
	// }

	// fmt.Println("Сработало изменение данных в функции EditHandler и изменено записей: ", rows)
	http.Redirect(w, r, "/createthemes", http.StatusSeeOther)
}
func changethemeincardAction(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	id := r.FormValue("id")
	nameprop := r.FormValue("nameprop")
	db3, err := sql.Open("mysql", conn_db_stands)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer db3.Close()

	result, err := db3.Exec("update demosys1 set `nameprop`=? where `id`=?",
		nameprop, id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Println("Сработало изменение данных в функции changethemeincardAction и изменено записей: ", rows)
	http.Redirect(w, r, "/sys/"+id, http.StatusSeeOther)
}
func changenumberpocketAction(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	id := r.FormValue("id")
	numberpp := r.FormValue("numberpp")
	db3, err := sql.Open("mysql", conn_db_stands)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer db3.Close()

	result, err := db3.Exec("update demosys1 set `numberpp`=? where `id`=?",
		numberpp, id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Println("Сработало изменение данных в функции changenumberpocketAction и изменено записей: ", rows)
	http.Redirect(w, r, "/sys/"+id, http.StatusSeeOther)
}
func changedateupdateAction(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	id := r.FormValue("id")
	dateupdate := r.FormValue("dateupdate")
	db3, err := sql.Open("mysql", conn_db_stands)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer db3.Close()

	result, err := db3.Exec("update demosys1 set `dateupdate`=? where `id`=?",
		dateupdate, id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Println("Сработало изменение данных в функции changedateupdateAction и изменено записей: ", rows)
	http.Redirect(w, r, "/sys/"+id, http.StatusSeeOther)
}
func changedatecheckAction(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	id := r.FormValue("id")
	datecheck := r.FormValue("datecheck")
	db3, err := sql.Open("mysql", conn_db_stands)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer db3.Close()

	result, err := db3.Exec("update demosys1 set `datecheck`=? where `id`=?",
		datecheck, id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Println("Сработало изменение данных в функции changedatecheckAction и изменено записей: ", rows)
	http.Redirect(w, r, "/sys/"+id, http.StatusSeeOther)
}
func changedatecheckAction01(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	id := r.FormValue("id")

	now := time.Now()
	datecheck := now.Format("2006-01-02 15:04:05")

	actual := r.FormValue("actual")
	db3, err := sql.Open("mysql", conn_db_stands)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer db3.Close()

	result, err := db3.Exec("update demosys1 set `datecheck`=?, `actual`=? where `id`=?",
		datecheck, actual, id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// Открытие файла для добавления (если файла нет, он будет создан)
	file, err := os.OpenFile("assets/txt/example.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close() // Закрытие файла в конце работы

	// Запись новой строки в файл
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	dateTimeStr := now.Format("2006-01-02 15:04:05")
	str := "::: Адрес:" + ip + " Проверка КС Антропово: " + id + "." + " Состояние листовки: " + actual + "."
	перевод_каретки := "\n"

	_, err = file.WriteString(dateTimeStr + str + перевод_каретки)
	if err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
		return
	}
	fmt.Println("Данные успешно добавлены в файл.")

	fmt.Println("Сработало изменение данных в функции changedatecheckAction и изменено записей: ", rows)
	http.Redirect(w, r, "/sys01/"+id, http.StatusSeeOther)
}
func district01(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/district01.html", "templates/tabs.html")
	checkErr(err)

	demo_db, err := sql.Open("mysql", conn_db_stands)
	checkErr(err)
	defer demo_db.Close()

	rows, err := demo_db.QueryContext(context.Background(), "SELECT id, namesys, datecreate, dateupdate, datecheck, path, nameprop, actual, places, numberpp, CASE WHEN dateupdate > datecheck THEN 'link-danger' WHEN dateupdate <= datecheck AND actual = '1' THEN 'link-success' ELSE 'link-secondary' END AS check_success FROM demosys1 WHERE places = ? ORDER BY id DESC", "КС Антропово")
	checkErr(err)
	defer rows.Close()

	var demos []Demo
	for rows.Next() {
		var demo Demo
		if err := rows.Scan(&demo.Id, &demo.Namesys, &demo.Datecreate, &demo.Dateupdate,
			&demo.Datecheck, &demo.Path, &demo.Nameprop, &demo.Actual, &demo.Places, &demo.Numberpp, &demo.Check_success); err != nil {
			log.Println(err)
			continue // продолжаем обрабатывать остальные строки
		}
		demos = append(demos, demo)
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
	}

	namesyses_rows, err := demo_db.QueryContext(context.Background(), "SELECT * FROM namesyses")
	checkErr(err)
	defer namesyses_rows.Close()

	var namesyses []Namesyses
	for namesyses_rows.Next() {
		var namesyse Namesyses
		if err := namesyses_rows.Scan(&namesyse.Id, &namesyse.Namesys); err != nil {
			log.Println(err)
			continue
		}
		namesyses = append(namesyses, namesyse)
	}

	essenses := struct {
		Demos     []Demo
		Namesyses []Namesyses
	}{
		Demos:     demos,
		Namesyses: namesyses,
	}

	t.ExecuteTemplate(w, "district01", essenses)
}
func showSys(w http.ResponseWriter, r *http.Request) {
	log.Print("showSys:")
	vars := mux.Vars(r)
	t, err := template.ParseFiles("templates/showsys.html", "templates/tabs.html")
	if err != nil {
		log.Println(err)
	}
	db, err := sql.Open("mysql", conn_db_stands)
	if err != nil {
		panic(err)
	}

	res, err := func() (*sql.Rows, error) {
		var query string = fmt.Sprintf("select * from `demosys1` where `id` = '%s'", vars["id"])
		return db.QueryContext(context.Background(), query, []any{}...)
	}()
	if err != nil {
		panic(err)
	}

	showSys := Demo{

		Id:         0,
		Namesys:    "",
		Datecreate: "",
		Dateupdate: "",
		Datecheck:  "",
		Path:       "",
		Nameprop:   "",
		Actual:     0,
		Places:     "",
		Numberpp:   0,
	}
	for res.Next() {
		var post22 Demo
		err = res.Scan(&post22.Id, &post22.Namesys, &post22.Datecreate,
			&post22.Dateupdate, &post22.Datecheck, &post22.Path, &post22.Nameprop,
			&post22.Actual, &post22.Places, &post22.Numberpp)
		if err != nil {
			panic(err)
		}
		defer res.Close()
		showSys = post22
	}

	demo_db, err := sql.Open("mysql", conn_db_stands)
	if err != nil {
		log.Println(err)
	}
	defer demo_db.Close()

	rows, err := demo_db.QueryContext(context.Background(), fmt.Sprintf("select id, namesys, datecreate, dateupdate, datecheck, path, nameprop, actual, places, numberpp from demosys1"), []any{}...)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var demos []Demo
	for rows.Next() {
		var demo Demo
		if err := rows.Scan(&demo.Id, &demo.Namesys, &demo.Datecreate,
			&demo.Dateupdate, &demo.Datecheck, &demo.Path, &demo.Nameprop,
			&demo.Actual, &demo.Places, &demo.Numberpp); err != nil {
			log.Println(err)
		}
		demos = append(demos, demo)
	}
	themes_db, err := sql.Open("mysql", conn_db_stands)
	if err != nil {
		log.Println(err)
	}
	defer themes_db.Close()

	themes_rows, err := themes_db.QueryContext(context.Background(), fmt.Sprintf("select id, theme from themes"), []any{}...)
	if err != nil {
		log.Println(err)
	}
	defer themes_rows.Close()

	var themes []Themes
	for themes_rows.Next() {
		var theme Themes
		if err := themes_rows.Scan(&theme.Id, &theme.Theme); err != nil {
			log.Println(err)
		}
		themes = append(themes, theme)
	}

	ss := "test"
	essenses := struct {
		ShowSys Demo
		Ss      string
		Demos   []Demo
		Themes  []Themes
	}{
		ShowSys: showSys,
		Ss:      ss,
		Demos:   demos,
		Themes:  themes,
	}

	showThemesGlob = showSys.Nameprop
	fmt.Println("showThemesGlob from Nameprop:", showThemesGlob)

	t.ExecuteTemplate(w, "showsys", essenses)
}
func showSys01(w http.ResponseWriter, r *http.Request) {
	log.Print("showSys:")
	vars := mux.Vars(r)
	t, err := template.ParseFiles("templates/showsys01.html", "templates/tabs.html")
	if err != nil {
		log.Println(err)
	}
	db, err := sql.Open("mysql", conn_db_stands)
	if err != nil {
		panic(err)
	}

	res, err := func() (*sql.Rows, error) {
		var query string = fmt.Sprintf("select * from `demosys1` where `id` = '%s'", vars["id"])
		return db.QueryContext(context.Background(), query, []any{}...)
	}()
	if err != nil {
		panic(err)
	}

	showSys := Demo{

		Id:         0,
		Namesys:    "",
		Datecreate: "",
		Dateupdate: "",
		Datecheck:  "",
		Path:       "",
		Nameprop:   "",
		Actual:     0,
		Places:     "",
		Numberpp:   0,
	}
	for res.Next() {
		var post22 Demo
		err = res.Scan(&post22.Id, &post22.Namesys, &post22.Datecreate,
			&post22.Dateupdate, &post22.Datecheck, &post22.Path, &post22.Nameprop,
			&post22.Actual, &post22.Places, &post22.Numberpp)
		if err != nil {
			panic(err)
		}
		defer res.Close()
		showSys = post22
	}

	demo_db, err := sql.Open("mysql", conn_db_stands)
	if err != nil {
		log.Println(err)
	}
	defer demo_db.Close()

	rows, err := demo_db.QueryContext(context.Background(), fmt.Sprintf("select id, namesys, datecreate, dateupdate, datecheck, path, nameprop, actual, places, numberpp from demosys1"), []any{}...)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var demos []Demo
	for rows.Next() {
		var demo Demo
		if err := rows.Scan(&demo.Id, &demo.Namesys, &demo.Datecreate,
			&demo.Dateupdate, &demo.Datecheck, &demo.Path, &demo.Nameprop,
			&demo.Actual, &demo.Places, &demo.Numberpp); err != nil {
			log.Println(err)
		}
		demos = append(demos, demo)
	}
	themes_db, err := sql.Open("mysql", conn_db_stands)
	if err != nil {
		log.Println(err)
	}
	defer themes_db.Close()

	themes_rows, err := themes_db.QueryContext(context.Background(), fmt.Sprintf("select id, theme from themes"), []any{}...)
	if err != nil {
		log.Println(err)
	}
	defer themes_rows.Close()

	var themes []Themes
	for themes_rows.Next() {
		var theme Themes
		if err := themes_rows.Scan(&theme.Id, &theme.Theme); err != nil {
			log.Println(err)
		}
		themes = append(themes, theme)
	}

	ss := "test"
	essenses := struct {
		ShowSys Demo
		Ss      string
		Demos   []Demo
		Themes  []Themes
	}{
		ShowSys: showSys,
		Ss:      ss,
		Demos:   demos,
		Themes:  themes,
	}
	showThemesGlob = showSys.Nameprop
	fmt.Println("showThemesGlob from Nameprop:", showThemesGlob)

	t.ExecuteTemplate(w, "showsys01", essenses)
}
func district02(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/district02.html")
	if err != nil {
		log.Print(err)
	}

	t.ExecuteTemplate(w, "district02", nil)
}
func pdfHandlerBackup(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", conn_db_stands)
	if err != nil {
		panic(err)
	}
	pdf5 := r.FormValue("id")

	res, err := func() (*sql.Rows, error) {
		var query string = fmt.Sprintf("select * from `themes` where `id` = '%s'", pdf5)
		return db.QueryContext(context.Background(), query, []any{}...)
	}()
	if err != nil {
		panic(err)
	}

	var themes []Themes // Объявляем themes как срез структур Themes

	for res.Next() {
		var post22 Themes
		err := res.Scan(&post22.Id, &post22.Theme)
		if err != nil {
			panic(err)
		}
		themes = append(themes, post22) // Добавляем каждый post к срезу themes
	}
	defer res.Close()

	// Закрываем результат после обработки

	// Выводим результат
	fmt.Println(themes) // Это выведет срез themes
	for _, theme := range themes {
		fmt.Println(theme.Theme) // Выводит значение поля Theme для каждого элемента

		xxx := showThemesGlob
		fmt.Println("xxx:", xxx)
		// Открываем файл PDF
		pdfFile, err := os.Open("assets/" + xxx + "/example.pdf")
		if err != nil {
			http.Error(w, "Could not open PDF file", http.StatusInternalServerError)
			return
		}
		defer pdfFile.Close()

		// Устанавливаем правильный Content-Type для PDF
		w.Header().Set("Content-Type", "application/pdf")
		// w.Header().Set("Content-Type", "application/pdf")
		// Content-Type:[application/vnd.openxmlformats-officedocument.spreadsheetml.sheet]]
		//Content-Type:[application/vnd.openxmlformats-officedocument.presentationml.presentation]]

		// Копируем содержимое файла в HTTP-ответ
		_, err = io.Copy(w, pdfFile)
		if err != nil {
			http.Error(w, "Could not copy PDF content to response", http.StatusInternalServerError)
			return
		}
	}
}
func pdfHandler(w http.ResponseWriter, r *http.Request) {
	xxx := showThemesGlob
	fmt.Println("xxx:", xxx)
	// Открываем файл PDF
	pdfFile, err := os.Open("assets/pdf/" + xxx + "/example.pdf")
	if err != nil {
		http.Error(w, "Could not open PDF file", http.StatusInternalServerError)
		return
	}
	defer pdfFile.Close()

	// Устанавливаем правильный Content-Type для PDF
	w.Header().Set("Content-Type", "application/pdf")
	// w.Header().Set("Content-Type", "application/pdf")
	// Content-Type:[application/vnd.openxmlformats-officedocument.spreadsheetml.sheet]]
	//Content-Type:[application/vnd.openxmlformats-officedocument.presentationml.presentation]]

	// Копируем содержимое файла в HTTP-ответ
	_, err = io.Copy(w, pdfFile)
	if err != nil {
		http.Error(w, "Could not copy PDF content to response", http.StatusInternalServerError)
		return
	}
}

func twosqlrequestAction(w http.ResponseWriter, r *http.Request) {

	// Установите соединение с базой данных
	db, err := sql.Open("mysql", conn_db_stands)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	rows, err := db.Query("update demosys1 set `nameprop`='МСК2' where `nameprop`='МСК'")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	http.Redirect(w, r, "/createthemes", http.StatusSeeOther)
}
func twoAction(w http.ResponseWriter, r *http.Request) {
	editthemesAction(w, r)
	twosqlrequestAction(w, r)
}
func addfiledownloadedAction(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	vars := mux.Vars(r)
	id := vars["id"]
	addPic(w, r)

	//Обновление колонки datecheck при обновлении листовки
	now := time.Now()
	dateupdate := now.Format("2006-01-02 15:04:05")
	nameprop := r.FormValue("theme")
	actual := 0
	db3, err := sql.Open("mysql", conn_db_stands)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer db3.Close()

	result, err := db3.Exec("update demosys1 set `dateupdate`=?, `actual`=? where `nameprop`=?",
		dateupdate, actual, nameprop)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Println("Сработало изменение данных в функции changedatecheckAction и изменено записей: ", rows)

	// Открытие файла для добавления (если файла нет, он будет создан)
	file, err := os.OpenFile("assets/txt/example.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close() // Закрытие файла в конце работы

	// Запись новой строки в файл
	dateTimeStr := now.Format("2006-01-02 15:04:05")
	str := "::: Обновление листовки в режиме настроек листовок: " + id
	перевод_каретки := "\n"

	_, err = file.WriteString(dateTimeStr + str + перевод_каретки)
	if err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
		return
	}
	fmt.Println("Данные успешно добавлены в файл.")

	fmt.Println("id in addfiledownloadedAction:", id)

	http.Redirect(w, r, "/themes/"+id, http.StatusSeeOther)
}
func addPic(w http.ResponseWriter, r *http.Request) {
	log.Println("addPic")
	Who_ip(w, r)

	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
	}
	defer file.Close()
	log.Println("Handler.FileName", handler.Filename)
	fmt.Printf("MIME Header: %+v\n", handler.Header)
	theme := r.FormValue("theme")
	// themes := r.FormValue("themes")
	// branch := r.FormValue("branch")
	// nameF := r.FormValue("namef")

	// Обрабатываем текущую дату, которая используется при создании записи
	now := time.Now()
	alkaa := now.Format("2006-01-02")

	fmt.Println(alkaa)
	contentType := handler.Header["Content-Type"][0]
	if contentType == "application/pdf" {

		// Создаем директорию с помощью os.Mkdir
		dirPath := "assets/pdf/" + theme
		err := os.Mkdir(dirPath, 0755) // 0755 - это права доступа к директории
		if err != nil {
			fmt.Printf("Ошибка при создании директории с помощью os.Mkdir: %v\n", err)
		} else {
			fmt.Println("Директория успешно создана с помощью os.Mkdir")
		}

		tempFile, err := os.Create("assets/pdf/" + theme + "/example.pdf")
		if err != nil {
			fmt.Println(err)
		}
		log.Println("tempFile:", *tempFile)

		defer tempFile.Close()

		// Копируем данные из загруженного файла в локальный файл
		if _, err := io.Copy(tempFile, file); err != nil {
			http.Error(w, "Ошибка при сохранении файла", http.StatusInternalServerError)
			return
		}

	} else if contentType == "application/x-zip-compressed" {

		tempFile, err := os.CreateTemp("assets"+"/"+"pdf"+"/"+"1"+"/"+handler.Filename, "*.zip")
		if err != nil {
			fmt.Println(err)
		}
		log.Println("tempFile:", tempFile)

		defer tempFile.Close()

		fileBytes, err := io.ReadAll(io.Reader(file))
		if err != nil {
			fmt.Println(err)
		}

		tempFile.Write(fileBytes)
	} else {
		fmt.Fprintf(w, "<html><br> <br><h1>Неподдерживаемый тип загружаемого файла</h1><br><h3> Файлы могут быть добавлены пока только в формате zip и pdf.</h3><br><br><h3>Файл не загружен на сервер. Добавьте файлы в zip-папку и загрузите ее.</h3></html>")
	}
}
func Who_ip(w http.ResponseWriter, r *http.Request) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	fmt.Println("ip-адрес:", ip)
}
func logactionusers(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/demo1.html", "templates/tabs.html"))

	// Чтение из файла
	data_out, err := ioutil.ReadFile("assets/txt/example.txt")
	if err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return
	}

	// Печать содержимого файла
	fmt.Println("Содержимое файла:")
	fmt.Println(string(data_out))

	essenses := struct {
		Data_out []byte
	}{

		Data_out: data_out,
	}

	t.Execute(w, essenses)
}
func logfiles(id string, namesys string) {
	// Открытие файла для добавления (если файла нет, он будет создан)
	file, err := os.OpenFile("assets/txt/example.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close() // Закрытие файла в конце работы

	// Запись новой строки в файл
	now := time.Now()
	dateTimeStr := now.Format("2006-01-02 15:04:05")
	str := "::: Удаление темы. Идентификатор:" + id + namesys
	перевод_каретки := "\n"

	_, err = file.WriteString(dateTimeStr + str + перевод_каретки)
	if err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
		return
	}
	fmt.Println("Данные успешно добавлены в файл.")
}
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/pdf/{id:[0-9]+}", pdfHandler)
	r.HandleFunc("/logactionusers", logactionusers)

	r.HandleFunc("/create", create)
	r.HandleFunc("/listcheck", listcheck)
	r.HandleFunc("/createAction", createAction)
	r.HandleFunc("/createnamesyses", createNamesyses)

	r.HandleFunc("/createnamesysesAction", createnamesysesAction)
	r.HandleFunc("/createthemes", createThemes)
	r.HandleFunc("/createthemesAction", createthemesAction)

	r.HandleFunc("/addfiledownloadedAction/{id:[0-9]+}", addfiledownloadedAction)

	r.HandleFunc("/sys/{id:[0-9]+}", showSys)
	r.HandleFunc("/sys01/{id:[0-9]+}", showSys01)

	r.HandleFunc("/namesys/{id:[0-9]+}", showNamesys)
	r.HandleFunc("/themes/{id:[0-9]+}", showThemes)

	r.HandleFunc("/editnamesysAction", editnamesysAction)
	r.HandleFunc("/editthemesAction", editthemesAction)
	r.HandleFunc("/deletenamesysAction", deletenamesysAction)
	r.HandleFunc("/deletethemesAction", deletethemesAction)

	r.HandleFunc("/changethemeincardAction", changethemeincardAction)
	r.HandleFunc("/changenumberpocketAction", changenumberpocketAction)
	r.HandleFunc("/changedateupdateAction", changedateupdateAction)
	r.HandleFunc("/changedatecheckAction", changedatecheckAction)
	r.HandleFunc("/changedatecheckAction01", changedatecheckAction01)

	r.HandleFunc("/twosqlrequestAction", twosqlrequestAction)
	r.HandleFunc("/twoAction", twoAction)

	r.HandleFunc("/district01", district01)
	r.HandleFunc("/district02/{id:[0-9]+}", BasicAuth02(district02))
	r.HandleFunc("/", index)
	http.Handle("/", r)

	httpDirForFileServer := http.Dir("./assets/")
	forStripPrefix := http.StripPrefix("/assets/", http.FileServer(httpDirForFileServer))
	http.Handle("/assets/", forStripPrefix)

	fmt.Println("Loaded on port :8089 stands4")
	http.ListenAndServeTLS(":8089", "/etc/letsencrypt/live/this.myftp.org/fullchain.pem", "/etc/letsencrypt/live/this.myftp.org/privkey.pem", nil)
}

// BasicAuth - Middleware для проверки базовой авторизации
func BasicAuth01(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || !checkCredentials01(username, password) {
			// Если учетные данные неверны, отправляем ответ 401
			w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your credentials"`)
			fmt.Println("Запрос авторизации 01")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// checkCredentials01 - проверка учетных данных
func checkCredentials01(username, password string) bool {
	// Здесь простая проверка, замените на вашу логику
	return username == "ks01" && password == "ks01"
}

// BasicAuth - Middleware для проверки базовой авторизации
func BasicAuth02(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || !checkCredentials02(username, password) {
			// Если учетные данные неверны, отправляем ответ 401
			w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your credentials"`)
			fmt.Println("Запрос авторизации 02")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// checkCredentials01 - проверка учетных данных
func checkCredentials02(username, password string) bool {
	// Здесь простая проверка, замените на вашу логику
	return username == "ks02" && password == "ks02"
}
