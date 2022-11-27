package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var bd *sql.DB
var param []string
var questActivas []QuestMenu
var dbSesiones = map[string]Sesion{} // session ID, user ID

//QuestMenu, objeto para array que lista las quest activas
type QuestMenu struct {
	ID_Quest     string `json:"IDQuest"`
	Nombre_Quest string `json:"NombreQuest"`
}

//Sesion, objeto para manejar sesiones
type Sesion struct {
	Sesion    string
	TimeStamp time.Time
}

//Sesion, objeto para manejar sesiones
type Cuenta struct {
	IDCuenta     string `json:"IDCuenta"`
	NombreCuenta string `json:"NombreCuenta"`
	Email        string `json:"Email"`
	Clave1       string `json:"Clave1"`
	Clave2       string `json:"Clave2"`
	Estado       string `json:"Estado"`
}

//Representeante, Usuario dueño de la cuenta y Representante de los equipos
type Representante struct {
	ID_Representante     string `json:"ID_Representante"`
	Nombre_Representante string `json:"NombreQuest"`
}

//Equipo, Equipoque participa en una quest
type Equipo struct {
	ID_Equipo     string `json:"IDEquipo"`
	ID_Quest      string `json:"IDQuest"`
	Nombre_Equipo string `json:"NombreQuest"`
}

//Miembro,de un equipo
type Miembro struct {
	ID_Miembro     string `json:"IDMiembro"`
	ID_Equipo      string `json:"IDEquipo"`
	Nombre_Miembro string `json:"NombreMiembro"`
	Rut_Miembro    string `json:"RutMiembro"`
}

func init() {
	actualizaparam()
	conectadb()
	go cargarQestActivas()
	go cerrarSesiones()
}

func actualizaparam() {
	content, err := ioutil.ReadFile("parametros.txt")
	if err != nil {
		log.Println(err.Error())
	}
	//lines := strings.Split(string(content), "\n")//version linux
	lines := strings.Split(strings.Replace(string(content), "\r\n", "\n", -1), "\n") //versio windows
	param = lines

}

func conectadb() {

	bda, err := sql.Open("mysql", param[2]+param[3]+param[4]+param[5]+param[6]+param[7]+param[8])
	if err != nil {
		log.Println(err.Error())
	}
	bd = bda
}

func cargarQestActivas() {
	for {
		log.Println("Cargando Pendientes")
		tab, err := bd.Query("select Q.ID_QUEST, Q.NOMBRE_QUEST  from QUEST Q where Q.ESTADO_QUEST = ?", "A")
		defer tab.Close()
		if err != nil {
			log.Println("Errores al marcar Solicitudes: " + err.Error())
			return
		}
		aux1 := []QuestMenu{}
		for tab.Next() {
			aux2 := QuestMenu{}
			err = tab.Scan(&aux2.ID_Quest, &aux2.Nombre_Quest)
			aux1 = append(aux1, aux2)
		}
		questActivas = aux1
		time.Sleep(time.Hour * 12)
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
	enableCORS(r)
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

	log.Println("escuchando puerto: " + param[0])

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func enableCORS(router *mux.Router) {
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	}).Methods(http.MethodOptions)
	router.Use(middlewareCors)
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			// Just put some headers to allow CORS...
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			// and call next handler!
			next.ServeHTTP(w, req)
		})
}

//____________________________________________________________________________________________

func handlerSesionCerrada(w http.ResponseWriter, r *http.Request) {

	respuesta, err := json.Marshal(Sesion{
		Sesion:    "Cerrada",
		TimeStamp: time.Time{},
	})
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

func handlerCrearCuenta(w http.ResponseWriter, r *http.Request) {

	nuevaCuenta := Cuenta{}
	err := json.NewDecoder(r.Body).Decode(&nuevaCuenta)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}
	resultado, err := CrearCuenta(nuevaCuenta)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}
	respuesta, err := json.Marshal(resultado)
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
func handlerIniciarSesion(w http.ResponseWriter, r *http.Request)   {}
func handlerCerrarSesion(w http.ResponseWriter, r *http.Request)    {}
func handlerRecuperarClave1(w http.ResponseWriter, r *http.Request) {}
func handlerRecuperarClave2(w http.ResponseWriter, r *http.Request) {}
func handlerQuest(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(param[1])
	if err != nil {
		handlerSesionCerrada(w, r)
		return
	}
	_, ok := dbSesiones[c.Value]
	if !ok {
		handlerSesionCerrada(w, r)
		return
	}

	respuesta, err := json.Marshal(questActivas)
	if err != nil {
		log.Println("Error: " + err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respuesta)

}
func CrearCuenta(nuevaCuenta Cuenta) (Cuenta, error) {
	var cuentaCreada = Cuenta{Estado: "False"}
	if nuevaCuenta.Clave1 == nuevaCuenta.Clave2 {
		db1, err := bd.Begin()
		if err != nil {
			log.Println("Error: " + err.Error())
			return cuentaCreada, err
		}
		tab1, err := db1.Query("insert into cuenta (EMAIL, CLAVE, NOMBRE, ESTADO) values (?,?,?,'1')", nuevaCuenta.Email, nuevaCuenta.Clave1, nuevaCuenta.NombreCuenta)
		defer tab1.Close()
		if err != nil {
			db1.Rollback()
			log.Println("Error: " + err.Error())
			return cuentaCreada, err
		}
		tab, err := bd.Query("select LAST_INSERT_ID()")
		defer tab.Close()
		if err != nil {
			db1.Rollback()
			log.Println("Error: " + err.Error())
			return cuentaCreada, err
		}
		for tab.Next() {
			err = tab.Scan(&cuentaCreada.IDCuenta)
			cuentaCreada = nuevaCuenta
			cuentaCreada.Estado = "True"
		}

		err = db1.Commit()
		if err != nil {
			log.Println("Error: " + err.Error())
			return cuentaCreada, err
		}
	}
	return cuentaCreada, nil
}
