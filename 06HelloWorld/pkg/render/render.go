package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func RenderTemplate(w http.ResponseWriter, tmpl string) {

	//create
	tc,err:=createTemplateCache()
	if err!=nil{
		log.Fatal(err)
	}
	//get from cache
	t,ok:=tc[tmpl]
	if !ok{
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	err=t.Execute(buf,nil)

	if err!=nil{
		log.Println(err)
	}

	//render
	_,err = buf.WriteTo(w)
	if err!=nil{
		log.Println(err)
	}

	// parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl,"./templates/base.layout.tmpl")
	// err = parsedTemplate.Execute(w, nil)
	// if err != nil {
	// 	fmt.Println("Error parsing template: ", err)
	// 	return
	// }
}

// var tc = make(map[string]*template.Template)

// func RenderTemplate(w http.ResponseWriter, t string){
// 	var tmpl *template.Template
// 	var err error

// 	//check to see if we already have the template in our cache

// 	_,inMap:= tc[t]
// 	if !inMap{
// 		log.Println("Creating template and adding to cache")
// 		createTemplateCache(t)
// 		if err!=nil{
// 			log.Println(err)
// 		}
// 	} else{
// 		log.Println("Using cached template")
// 	}

// 	tmpl = tc[t]
// 	err = tmpl.Execute(w,nil)
// }

// func createTemplateCache(t string) error{
// 	templates:= []string{
// 		fmt.Sprintf("./templates/%s",t),
// 		"./templates/base.layout.tmpl",
// 	}
// 	tmpl,err:= template.ParseFiles(templates...)
// 	if err!=nil{
// 		return err
// 	}
// 	tc[t] = tmpl
// 	return nil
// }

func createTemplateCache()(map[string]*template.Template,error){
	myCache:= map[string]*template.Template{}

	pages,err:= filepath.Glob("./templates/*page.tmpl")
	if err!=nil{
		return myCache,err
	}
	//range through all files ending with *.page.tmpl
	for _,page:=range pages{
		name:=filepath.Base(page)
		ts, err:=template.New(name).ParseFiles(page)

		if err!=nil{
			return myCache,err
		}
		matches, err:=filepath.Glob("./templates/*.layout.tmpl")
		if err!=nil{
			return myCache,err
		}
		if len(matches)>0{
			ts, err= ts.ParseGlob("./templates/*.layout.tmpl")
			if err!=nil{
				return myCache,err
			}
		}
		myCache[name] = ts

	}
	return myCache,nil
}