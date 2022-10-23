package main

import (
	"fmt"
    "net/http"
    "os"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "strings"
    "path/filepath"
    "strconv"
	"io"
	"log"
	"html/template"
)

type BudgetItem struct {
		Name	string
		Amount	int64
}

var budgetItemsCache [25]BudgetItem


func main() {
    total := 0
    current_value := 0
	//budgetItemsCache := LoadData()
	LoadData()

	for _, val := range budgetItemsCache {
		if val.Name != "" {
			fmt.Print(val.Name, "\t", val.Amount, "\n")
			total += int(val.Amount)
		}
	}
	fmt.Print("total:\t",total)

    r := chi.NewRouter()
    r.Use(middleware.Logger)

    r.Put("/add_row", func(w http.ResponseWriter, r *http.Request) {
		u := template.New("row_total").Parse
	    row := `
			<tr>
		    	<td class="px-5" hx-post='/update_cell' hx-trigger='input' contenteditable='true'>Item</td>
		    	<td>$</td>
		    	<td hx-post='/update_cell' hx-trigger='input' contenteditable='true'>0</td>
				<td>
				<button class="px-5 text-red-600" hx-delete="/delete-row">
				x
				</button>
				</td>
		    </tr>
			
			<tr id="total">
				<td></td>
				<td class="align-right"><b>$</b></td>
				<td class="py-5 hx-trigger="input from:td" hx-gets="/get_total"><b>0</b></td>
			</tr>
			`

	    w.Write([]byte(row))
    })

    r.Get("/send_data", func(w http.ResponseWriter, r *http.Request) {
	    s := strconv.Itoa(current_value)
		w.Write([]byte(s))
		fmt.Printf("\n\nsend data: %s\n\n", s)

    })

    r.Get("/update_cell", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\n\nupdate cell: %s\n\n", r)

	})

    r.Get("/get_total", func(w http.ResponseWriter, r *http.Request) {
		total = 0
		for _, t := range budgetItemsCache {
			total += int(t.Amount)
		}
		fmt.Print("total:\t",total)
	    s := strconv.Itoa(total)
	    w.Write([]byte("<b>" + s + "</b>"))
    })
    
    workDir, _ := os.Getwd()
    filesDir := http.Dir(filepath.Join(workDir, "."))
    FileServer(r, "/", filesDir)

    http.ListenAndServe(":3000", r)
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
	}
/*
func UpdateData(amount int, item string) {
	defer f.Close()
    f, err := os.OpenFile("data", os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println(err)
        return
    }
	build_string := item + "  " + amount
    //newLine := _, err = fmt.Fprintln(f, newLine)
    if err != nil {
      fmt.Println(err)
    }
}
*/

func LoadData() {
    f, err := os.Open("data")
    if err != nil {
        fmt.Println(err)
        return
    }

    buf := make([]byte, 2500)

    for {
        n, err := f.Read(buf)
        if err != nil && err != io.EOF {
            log.Fatal(err)
        }

        if err == io.EOF {
            break
        }
		
		var r_buf [25]byte
	    var r BudgetItem	
		var l1_count = 0 // index of r_buf
		var l2_count = 0 // index of budgetItemsCache 

		for _, val := range buf[:n] {
			// ascii ',' = 44
			// val '\n' = 10
			
	 		if val == byte(44) {
				s := string(r_buf[:l1_count])
				r.Name = s
				//fmt.Println(r.Name)
				l1_count = 0
				r_buf = [25]byte{}

			} else if val == byte(10) {
				s := string(r_buf[:l1_count])	
				int2, _ := strconv.ParseInt(s, 10, 32)
				r.Amount = int2
				budgetItemsCache[l2_count] = r
				l2_count++

				r_buf = [25]byte{}
				l1_count = 0

			} else {
				r_buf[l1_count] = val
				l1_count++
			}

		}
		//for _, x:= range budgetItemsCache[:l2_count] {
		//fmt.Print(x.Name," ",x.Amount)
		//fmt.Print("\n")
	}
}
/*
func PreventFloats() {*/
