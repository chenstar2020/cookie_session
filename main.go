package main

import (
	"html/template"
	"log"
	"net/http"
)


var (
	globalSessions *Manager  //session全局管理
)

func init(){
	var err error
	globalSessions, err = NewManager("memory", "gosessionid", 3600)
	if err != nil {
		panic("not register memory")
	}
}

func main(){
	http.HandleFunc("/login", login)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func login(w http.ResponseWriter, r *http.Request){
	sess := globalSessions.SessionStart(w, r)
	r.ParseForm()

	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.tpl")
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		t.Execute(w, sess.Get("username"))
	}else{
		sess.Set("username", r.Form["username"])
		http.Redirect(w, r, "/", 302)
	}
}