package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func init() {
	actualizaparam()
}

func actualizaparam() {
	content, err := ioutil.ReadFile("param.txt")
	if err != nil {
		log.Println(err.Error())
	}
	//lines := strings.Split(string(content), "\n")//version linux
	lines := strings.Split(strings.Replace(string(content), "\r\n", "\n", -1), "\n") //versio windows
	param = lines

}

func conectadb() {
	bd, err := sql.Open("mysql", param[2]+param[3]+param[4]+param[5]+param[6]+param[7]+param[8])
	if err != nil {
		log.Println(err.Error())
	}
}

func cerrarSesiones() {
	for {
		log.Println("Eliminando Sesiones Antiguas")
		timenow := time.Now().Add(time.Hour * -2)
		for key, element := range dbSesiones {
			if element.TimeStamp.Before(timenow) {
				delete(dbSesiones, key)
			}
		}
		time.Sleep(time.Hour * 2)
	}
}

//____________________________________________________________________________________________
func main() {
	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/", handlerSesionCerrada).Methods("POST")                  //1.0
	r.HandleFunc("/CrearCuenta", handlerCrearCuenta).Methods("POST")         //1.0
	r.HandleFunc("/IniciarSesion", handlerIniciarSesion).Methods("POST")     //1.0
	r.HandleFunc("/CerrarSesion", handlerCerrarSesion).Methods("POST")       //1.0
	r.HandleFunc("/RecuperarClave1", handlerRecuperarClave1).Methods("POST") //1.0
	r.HandleFunc("/RecuperarClave2", handlerRecuperarClave2).Methods("POST") //1.0
	r.HandleFunc("/Quests", handlerQuest).Methods("POST")                    //1.0
	r.HandleFunc("/Quest", handlerQuest).Methods("POST")                     //1.0
	r.HandleFunc("/Inscribirse", handlerQuest).Methods("POST")               //1.0
	r.HandleFunc("/EnviarRespuesta", handlerQuest).Methods("POST")           //1.0
	server := http.Server{
		Addr:           param[0],
		Handler:        r,
		ReadTimeout:    3 * time.Second,
		WriteTimeout:   3 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("escuchando puerto" + param[0])

	server.ListenAndServe()
}

//____________________________________________________________________________________________

func handlerSesionCerrada(w http.ResponseWriter, r *http.Request) {

	respuesta, err := json.Marshal(Sesion{Sesion: "Cerrada"})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respuesta)

}

func handlerCrearCuenta(w http.ResponseWriter, r *http.Request)     {}
func handlerIniciarSesion(w http.ResponseWriter, r *http.Request)   {}
func handlerCerrarSesion(w http.ResponseWriter, r *http.Request)    {}
func handlerRecuperarClave1(w http.ResponseWriter, r *http.Request) {}
func handlerRecuperarClave2(w http.ResponseWriter, r *http.Request) {}
func handlerQuest(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(param[1])
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	Cuenta, ok := dbSesiones[c.Value]
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	respuesta, err := json.Marshal(questActivas)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respuesta)

}
