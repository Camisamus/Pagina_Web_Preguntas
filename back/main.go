package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var bd *sql.DB
var param []string
var questActivas []QuestMenu
var dbClavesInscrion = map[string]bool{}            //
var dbSesiones = map[string]Sesion{}                // session ID, user ID
var dbUsuarios = map[string]Cuenta{}                //
var dbIntentos = map[string]Intento{}               //
var dbClavesEnproceso = map[string]CuentaEnEspera{} //
var dbQuestCeradas = map[string]bool{}

//QuestMenu, objeto para array que lista las quest activas
type QuestMenu struct {
	ID_Quest     string `json:"IDQuest"`
	Nombre_Quest string `json:"NombreQuest"`
	Categoria    string `json:"Categoria"`
}

//QuestDetalle, objeto para array que lista las quest activas
type QuestDetalle struct {
	ID_Quest      string        `json:"IDQuest"`
	Representante Representante `json:"Representante"`
	Equipos       []Equipo      `json:"Equipo"`
	Quest         Quest         `json:"Quest"`
}

//Quest, estructura de la quest
type Quest struct {
	ID_Quest     string     `json:"IDQuest"`
	Nombre_Quest string     `json:"NombreQuest"`
	Premio       string     `json:"Premio"`
	Preguntas    []Pregunta `json:"Preguntas"`
	Categoria    string     `json:"Categoria"`
	Estado       string     `json:"Estado"`
	FechaInicio  string     `json:"FechaInicio"`
	Ganador      string     `json:"Ganador"`
}
type Pregunta struct {
	ID_Pregunta string `json:"IDPregunta"`
	ID_Quest    string `json:"IDQuest"`
	Pista       string `json:"Pista"`
	Pregunta    string `json:"Pregunta"`
	Respuesta   string `json:"Respuesta"`
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
	Nombre_Representante string `json:"NombreRepresentante"`
	Email_Representante  string `json:"EmailRepresentante"`
}

//Equipo, Equipoque participa en una quest
type Equipo struct {
	ID_Equipo          string    `json:"IDEquipo"`
	ID_Quest           string    `json:"IDQuest"`
	ID_Representante   string    `json:"ID_Representante"`
	Nombre_Equipo      string    `json:"NombreEquipo"`
	Rut_Respondable    string    `json:"RutRespondable"`
	Nombre_Respondable string    `json:"NombreRespondable"`
	Miembros_Equipo    []Miembro `json:"Miembros_Equipo"`
}

//Miembro,de un equipo
type Miembro struct {
	ID_Miembro     string `json:"IDMiembro"`
	ID_Equipo      string `json:"IDEquipo"`
	Nombre_Miembro string `json:"NombreMiembro"`
	Rut_Miembro    string `json:"RutMiembro"`
}

//Miembro,de un equipo
type CuentaEnEspera struct {
	Cuenta    Cuenta
	TimeStamp time.Time
}

//Intento
type Intento struct {
	Hora       time.Time    `json:"HoraInivcio"`
	Respuestas []Respuestas `json:"Respuestas"`
	ID_Equipo  string       `json:"IDEquipo"`
	ID_Quest   string       `json:"IDQuest"`
}

//Respuestas
type Respuestas struct {
	Hora        time.Time `json:"HoraInivcio"`
	ID_Pregunta string    `json:"IDPregunta"`
	ID_Quest    string    `json:"IDQuest"`
	Pregunta    Pregunta  `json:"Pregunta"`
	Correcta    bool      `json:"Correcta"`
}

func init() {
	actualizaparam()
	conectadb()
	go cargarClavesActivas()
	go cargarQestActivas()
	go cerrarSesiones()
	go clavesNoRecuperadas()
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
		if err != nil {
			log.Println("Errores al marcar Solicitudes: " + err.Error())
			return
		}
		defer tab.Close()
		aux1 := []QuestMenu{}
		for tab.Next() {
			aux2 := QuestMenu{}
			err = tab.Scan(&aux2.ID_Quest, &aux2.Nombre_Quest)
			if err != nil {
				log.Println("Error: " + err.Error())
				return
			}
			aux1 = append(aux1, aux2)
		}
		questActivas = aux1
		time.Sleep(time.Hour * 12)
	}
}

func cargarClavesActivas() {
	for {
		log.Println("Cargando Codigos")
		tab, err := bd.Query("select i.CODIGO  from inscripciones i where i.EQUIPO = 0")
		if err != nil {
			log.Println("Errores al marcar Solicitudes: " + err.Error())
			return
		}
		defer tab.Close()
		aux1 := 0
		for tab.Next() {
			aux2 := ""
			err = tab.Scan(&aux2)
			if err != nil {
				log.Println("Error: " + err.Error())
				return
			}
			aux1++
			dbClavesInscrion[aux2] = true
		}
		db1, err := bd.Begin()
		if err != nil {
			log.Println("Error: " + err.Error())
			return
		}
		for i := aux1; i < 50; i++ {
			sID := uuid.NewV4()
			_, err := db1.Exec("INSERT INTO inscripciones (CODIGO, EQUIPO) VALUES( ? , 0);", sID.String())
			if err != nil {
				log.Println("Errores al marcar Solicitudes: " + err.Error())
				return
			}
			dbClavesInscrion[sID.String()] = true
		}
		db1.Commit()
		time.Sleep(time.Hour * 12)
	}
}

func cerrarSesiones() {
	for {
		log.Println("Eliminando Sesiones Antiguas")
		timenow := time.Now().Add(time.Hour * -2)
		for key, element := range dbSesiones {
			if element.TimeStamp.Before(timenow) {
				delete(dbUsuarios, element.Sesion)
				delete(dbSesiones, key)
			}
		}
		time.Sleep(time.Hour * 2)
	}
}

func clavesNoRecuperadas() {
	for {
		log.Println("Eliminando Claves No Recuperadas")
		timenow := time.Now().Add(time.Hour * -24)
		for key, element := range dbClavesEnproceso {
			if element.TimeStamp.Before(timenow) {
				delete(dbClavesEnproceso, key)
			}
		}
		time.Sleep(time.Hour * 24)
	}
}

func bajarIntento(intento string) {
	time.Sleep(time.Minute * 20)
	delete(dbIntentos, intento)
}

func EliminatIntentos() {
	for {
		log.Println("Eliminando Claves No Recuperadas")
		timenow := time.Now().Add(time.Minute * -20)
		for key, element := range dbIntentos {
			if element.Hora.Before(timenow) {
				delete(dbIntentos, key)
			}
		}
		time.Sleep(time.Minute * 20)
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
	r.HandleFunc("/Quest", handlerQuestDetalle).Methods("POST")              //1.0
	r.HandleFunc("/Inscribirse", handlerInscribirse).Methods("POST")         //1.0
	r.HandleFunc("/EnviarRespuesta", handlerRespuesta).Methods("POST")       //1.0

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
			w.Header().Set("Access-Control-Allow-Origin", "null")
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
func handlerIniciarSesion(w http.ResponseWriter, r *http.Request) {

	ingreso := Cuenta{}
	err := json.NewDecoder(r.Body).Decode(&ingreso)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}
	resultado, err := ingresar(ingreso)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}
	if resultado.Estado == "True" {
		sID := uuid.NewV4()
		c := &http.Cookie{
			Name:     param[1],
			Value:    sID.String(),
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
			Expires:  time.Now().Add(time.Hour + 2),
			Path:     "/Quests",
		}
		http.SetCookie(w, c)
		c2 := &http.Cookie{
			Name:     param[1],
			Value:    sID.String(),
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
			Expires:  time.Now().Add(time.Hour + 2),
			Path:     "/CerrarSesion",
		}
		http.SetCookie(w, c2)
		c3 := &http.Cookie{
			Name:     param[1],
			Value:    sID.String(),
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
			Expires:  time.Now().Add(time.Hour + 2),
			Path:     "/Quest",
		}
		http.SetCookie(w, c3)
		c4 := &http.Cookie{
			Name:     param[1],
			Value:    sID.String(),
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
			Expires:  time.Now().Add(time.Hour + 2),
			Path:     "/Inscribirse",
		}
		http.SetCookie(w, c4)
		ses := Sesion{}
		ses.Sesion = resultado.Email
		ses.TimeStamp = time.Now()
		dbSesiones[sID.String()] = ses
		dbUsuarios[resultado.Email] = resultado
		resultado.Clave1 = param[1]
		resultado.Clave2 = sID.String()
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
func handlerCerrarSesion(w http.ResponseWriter, r *http.Request) {

	c, err := r.Cookie(param[1])
	if err != nil {
		handlerSesionCerrada(w, r)
		return
	}
	sesion, ok := dbSesiones[c.Value]
	if !ok {
		handlerSesionCerrada(w, r)
		return
	}
	_, ok = dbUsuarios[sesion.Sesion]
	if !ok {
		handlerSesionCerrada(w, r)
		return
	}

	delete(dbUsuarios, sesion.Sesion)
	delete(dbSesiones, c.Value)
	respuesta, err := json.Marshal(Sesion{
		Sesion:    "Cerrada",
		TimeStamp: time.Time{},
	})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respuesta)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respuesta)
}

func handlerRecuperarClave1(w http.ResponseWriter, r *http.Request) {

	Cuenta := Cuenta{}
	err := json.NewDecoder(r.Body).Decode(&Cuenta)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}
	err = recuperarClave1(Cuenta)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}
func handlerRecuperarClave2(w http.ResponseWriter, r *http.Request) {
	Cuenta := Cuenta{}
	err := json.NewDecoder(r.Body).Decode(&Cuenta)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}
	clavevieja, ok := dbClavesEnproceso[Cuenta.Estado]
	if (Cuenta.Clave1 != Cuenta.Clave2) || (Cuenta.Clave2 == "") || (!ok) {
		respuesta, err := json.Marshal(Sesion{
			Sesion: "No se genero una clave Nueva",
		})
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(respuesta)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respuesta)

	}
	err = recuperarClave2(Cuenta, clavevieja.Cuenta)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}

	delete(dbClavesEnproceso, Cuenta.Estado)
	respuesta, err := json.Marshal(Sesion{
		Sesion: "Ya puede ingresar con su nueva clave",
	})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respuesta)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respuesta)

}

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

func handlerQuestDetalle(w http.ResponseWriter, r *http.Request) {
	questDetalle := QuestDetalle{}
	err := json.NewDecoder(r.Body).Decode(&questDetalle)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}
	c, err := r.Cookie(param[1])
	if err != nil {
		handlerSesionCerrada(w, r)
		return
	}
	sesion, ok := dbSesiones[c.Value]
	if !ok {
		handlerSesionCerrada(w, r)
		return
	}
	usuario, ok := dbUsuarios[sesion.Sesion]
	if !ok {
		handlerSesionCerrada(w, r)
		return
	}
	questDetalle.Representante = Representante{
		ID_Representante:     usuario.IDCuenta,
		Nombre_Representante: usuario.NombreCuenta,
		Email_Representante:  usuario.Email,
	}
	equiposEncontrados, err := buscarEquiposActivos(questDetalle.ID_Quest, questDetalle.Representante.ID_Representante)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}
	if len(equiposEncontrados) < 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
		return
	}
	questSeleccionada, err := buscarQuest(questDetalle.ID_Quest)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}
	questSeleccionada.Preguntas = limpiaRespueestas(questSeleccionada.Preguntas)
	questDetalle.Equipos = equiposEncontrados
	questDetalle.Quest = questSeleccionada
	respuesta, err := json.Marshal(questDetalle)
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

func handlerInscribirse(w http.ResponseWriter, r *http.Request) {
	nuevoEquipo := Equipo{}
	err := json.NewDecoder(r.Body).Decode(&nuevoEquipo)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}
	c, err := r.Cookie(param[1])
	if err != nil {
		handlerSesionCerrada(w, r)
		return
	}
	sesion, ok := dbSesiones[c.Value]
	if !ok {
		handlerSesionCerrada(w, r)
		return
	}
	usuario, ok := dbUsuarios[sesion.Sesion]
	if !ok {
		handlerSesionCerrada(w, r)
		return
	}
	nuevoEquipo.ID_Representante = usuario.IDCuenta
	ok, err = ValidarClave(nuevoEquipo.ID_Equipo)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}
	if !ok {

		respuesta, err := json.Marshal("Clave Invalida")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respuesta)
		return
	}
	id, err := agregarEquipo(nuevoEquipo)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}
	err = usarClave(nuevoEquipo.ID_Equipo, id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}
	respuesta, err := json.Marshal("Equipo Inscrito")
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

func handlerRespuesta(w http.ResponseWriter, r *http.Request) {
	nuevaRespuesta := Intento{}
	err := json.NewDecoder(r.Body).Decode(&nuevaRespuesta)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}
	_, ok := dbQuestCeradas[nuevaRespuesta.ID_Quest]
	if ok {
		respuesta, err := json.Marshal("{ID_Equipo: Tarde}")
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
	resultado, err := revisarRespuesta(nuevaRespuesta)
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
func ValidarClave(clave string) (bool, error) {
	ok := dbClavesInscrion[clave]
	if !ok {
		return false, nil
	}
	dbClavesInscrion[clave] = false
	return ok, nil
}
func usarClave(clave string, equipo string) error {
	db1, err := bd.Begin()
	if err != nil {
		log.Println("Error: " + err.Error())
		return err
	}
	_, err = db1.Query("UPDATE inscripciones set EQUIPO= ? where  CODIGO = ?;", equipo, clave)
	if err != nil {
		db1.Rollback()
		log.Println("Error: " + err.Error())
		return err
	}

	sID := uuid.NewV4()
	_, err = db1.Exec("INSERT INTO inscripciones (CODIGO, EQUIPO) VALUES( ? , 0);", sID.String())
	if err != nil {
		log.Println("Errores al marcar Solicitudes: " + err.Error())
		return err
	}

	db1.Commit()
	dbClavesInscrion[sID.String()] = true
	delete(dbClavesInscrion, clave)
	return nil
}
func agregarEquipo(equipo Equipo) (string, error) {
	db1, err := bd.Begin()
	if err != nil {
		log.Println("Error: " + err.Error())
		return "", err
	}
	tab1, err := db1.Query("INSERT INTO equipo (ID_QUEST, ID_REPRESENTANTE, NOMBRE_EQUIPO, RUT_RESPONSABLE, NOMBRE_RESPONSABLE) VALUES(?, ?, ?, ?, ?);", equipo.ID_Quest, equipo.ID_Representante, equipo.Nombre_Equipo, equipo.Rut_Respondable, equipo.Nombre_Respondable)
	if err != nil {
		db1.Rollback()
		log.Println("Error: " + err.Error())
		return "", err
	}
	tab1.Close()
	tab, err := db1.Query("select LAST_INSERT_ID()")
	if err != nil {
		db1.Rollback()
		log.Println("Error: " + err.Error())
		return "", err
	}
	defer tab.Close()
	for tab.Next() {
		err = tab.Scan(&equipo.ID_Equipo)
		if err != nil {
			db1.Rollback()
			log.Println("Error: " + err.Error())
			return "", err
		}
	}
	for _, element := range equipo.Miembros_Equipo {

		tab2, err := db1.Query("INSERT INTO miembro (ID_EQUIPO, RUT_MIEMBRO, NOMBRE_MIEMBRO) VALUES(?, ?, ?);", equipo.ID_Equipo, element.Rut_Miembro, element.Nombre_Miembro)
		if err != nil {
			db1.Rollback()
			log.Println("Error: " + err.Error())
			return "", err
		}
		tab2.Close()
	}

	err = db1.Commit()
	if err != nil {
		db1.Rollback()
		log.Println("Error: " + err.Error())
		return "", err
	}
	return equipo.ID_Equipo, nil
}

func revisarRespuesta(respuesta Intento) (Intento, error) {
	llave := respuesta.ID_Quest + "_" + respuesta.ID_Equipo
	actual, ok := dbIntentos[llave]
	if !ok {
		actual = Intento{
			Hora:       time.Now(),
			Respuestas: []Respuestas{},
			ID_Equipo:  respuesta.ID_Equipo,
			ID_Quest:   respuesta.ID_Quest,
		}
	}
	for _, element := range actual.Respuestas {
		if respuesta.Respuestas[0].ID_Pregunta == element.ID_Pregunta {
			actual.ID_Equipo = "Aun no se permite volver a responder"
			return actual, nil
		}
	}
	preguntas, err := buscarPreguntas(respuesta.ID_Quest)
	if err != nil {
		return actual, err
	}
	var pregunta = Pregunta{}
	for _, element := range preguntas {
		if respuesta.Respuestas[0].ID_Pregunta == element.ID_Pregunta {
			pregunta = element
		}
	}
	var nuevaRespuesta = Respuestas{
		Hora:        time.Now(),
		ID_Pregunta: respuesta.Respuestas[0].ID_Pregunta,
		ID_Quest:    respuesta.ID_Quest,
		Correcta:    (pregunta.Respuesta == respuesta.Respuestas[0].Pregunta.Respuesta),
	}
	actual.Respuestas = append(actual.Respuestas, nuevaRespuesta)
	respuesta.Respuestas[0].Correcta = nuevaRespuesta.Correcta
	if len(actual.Respuestas) == len(preguntas) {
		ganador := true
		for _, element := range actual.Respuestas {
			if !element.Correcta {
				ganador = false
			}
		}
		if ganador {
			dbQuestCeradas[respuesta.ID_Quest] = true
			var newQuestActivas []QuestMenu
			for _, element := range questActivas {
				if element.ID_Quest != respuesta.ID_Quest {
					newQuestActivas = append(newQuestActivas, element)
				}
			}
			questActivas = newQuestActivas
			err := declararGanador(respuesta.ID_Equipo, respuesta.ID_Quest)
			respuesta.ID_Equipo = "Gano"
			if err != nil {
				log.Println("Equipo Ganador No Registrado equipo : " + respuesta.ID_Equipo + "  hora:  " + time.Now().String())
				return respuesta, err
			}
		}
	}
	dbIntentos[llave] = actual
	go bajarIntento(llave)
	return respuesta, nil
}
func declararGanador(equipo string, ID_Quest string) error {

	db1, err := bd.Begin()
	if err != nil {
		log.Println("Error: " + err.Error())
		return err
	}
	_, err = db1.Query("UPDATE quest set ESTADO_QUEST='F', GANADOR = ? where  ID_QUEST = ?;", equipo, ID_Quest)
	if err != nil {
		db1.Rollback()
		log.Println("Error: " + err.Error())
		return err
	}
	db1.Commit()
	return nil
}
func CrearCuenta(nuevaCuenta Cuenta) (Cuenta, error) {
	var cuentaCreada = Cuenta{Estado: "False"}
	if nuevaCuenta.Clave1 == nuevaCuenta.Clave2 {

		contraseñaPlanaComoByte1 := []byte(nuevaCuenta.Clave1)
		hash1, err := bcrypt.GenerateFromPassword(contraseñaPlanaComoByte1, 11) //DefaultCost es 10
		if err != nil {
			log.Println("Error: " + err.Error())
			return cuentaCreada, err
		}
		nuevaCuenta.Clave1 = string(hash1)
		db1, err := bd.Begin()
		if err != nil {
			log.Println("Error: " + err.Error())
			return cuentaCreada, err
		}
		tab1, err := db1.Query("insert into cuenta (EMAIL, CLAVE, NOMBRE, ESTADO) values (?,?,?,'1')", nuevaCuenta.Email, nuevaCuenta.Clave1, nuevaCuenta.NombreCuenta)
		if err != nil {
			db1.Rollback()
			log.Println("Error: " + err.Error())
			return cuentaCreada, err
		}
		defer tab1.Close()
		tab, err := bd.Query("select LAST_INSERT_ID()")
		if err != nil {
			db1.Rollback()
			log.Println("Error: " + err.Error())
			return cuentaCreada, err
		}
		defer tab.Close()
		for tab.Next() {
			err = tab.Scan(&cuentaCreada.IDCuenta)
			if err != nil {
				db1.Rollback()
				log.Println("Error: " + err.Error())
				return cuentaCreada, err
			}
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
func ingresar(cuenta Cuenta) (Cuenta, error) {
	var cuentaingresada = Cuenta{Estado: "False"}

	db1, err := bd.Begin()
	if err != nil {
		log.Println("Error: " + err.Error())
		return cuentaingresada, err
	}
	tab1, err := db1.Query("SELECT c.ID_CUENTA,	c.NOMBRE , c.EMAIL, c.CLAVE, c.ESTADO FROM cuenta c where c.EMAIL = ?", cuenta.Email)
	if err != nil {
		db1.Rollback()
		log.Println("Error: " + err.Error())
		return cuentaingresada, err
	}
	defer tab1.Close()
	encontrado := false
	for tab1.Next() {
		err = tab1.Scan(&cuentaingresada.IDCuenta, &cuentaingresada.NombreCuenta, &cuentaingresada.Email, &cuentaingresada.Clave1, &cuentaingresada.Estado)
		encontrado = true
		if err != nil {
			db1.Rollback()
			log.Println("Error: " + err.Error())
			return cuentaingresada, err
		}
	}
	err = db1.Commit()

	if err != nil {
		log.Println("Error: " + err.Error())
		return cuentaingresada, err
	}

	error := bcrypt.CompareHashAndPassword([]byte(cuentaingresada.Clave1), []byte(cuenta.Clave1))
	if error != nil {
		cuentaingresada.Estado = "Email o contraseña Incorrectos"
		return cuentaingresada, err
	}
	if !encontrado {
		cuentaingresada.Estado = "Email o contraseña Incorrectos"
		return cuentaingresada, err
	}
	cuentaingresada.Estado = "True"

	return cuentaingresada, nil
}

func recuperarClave1(cuenta Cuenta) error {
	var cuentaingresada = Cuenta{Estado: "False"}

	db1, err := bd.Begin()
	if err != nil {
		log.Println("Error: " + err.Error())
		return err
	}
	tab1, err := db1.Query("SELECT c.ID_CUENTA,	c.NOMBRE , c.EMAIL, c.CLAVE, c.ESTADO FROM cuenta c where c.EMAIL = ?", cuenta.Email)
	if err != nil {
		db1.Rollback()
		log.Println("Error: " + err.Error())
		return err
	}
	defer tab1.Close()
	encontrado := false
	for tab1.Next() {
		err = tab1.Scan(&cuentaingresada.IDCuenta, &cuentaingresada.NombreCuenta, &cuentaingresada.Email, &cuentaingresada.Clave1, &cuentaingresada.Estado)
		encontrado = true
		if err != nil {
			db1.Rollback()
			log.Println("Error: " + err.Error())
			return err
		}
	}
	err = db1.Commit()
	if err != nil {
		db1.Rollback()
		log.Println("Error: " + err.Error())
		return err
	}
	if encontrado {
		sID := uuid.NewV4()
		dbClavesEnproceso[sID.String()] = CuentaEnEspera{
			Cuenta:    cuentaingresada,
			TimeStamp: time.Now(),
		}

		link := "https://" + param[12] + "/recuperarclave2.html?ClaveCambio=" + sID.String()
		msg := "Su Link Es:  " + link
		log.Println("Contraseña Habilitada para crearse :" + cuentaingresada.NombreCuenta)
		return enviarmail(msg, cuentaingresada.Email, "Link Para Cambiar Clave :")

	}
	return nil
}

func recuperarClave2(nuevasClaves Cuenta, cuentaAcutal Cuenta) error {
	contraseñaPlanaComoByte1 := []byte(nuevasClaves.Clave1)
	hash1, err := bcrypt.GenerateFromPassword(contraseñaPlanaComoByte1, 11) //DefaultCost es 10
	if err != nil {
		log.Println("Error: " + err.Error())
		return err
	}
	nuevasClaves.Clave1 = string(hash1)
	db1, err := bd.Begin()
	if err != nil {
		log.Println("Error: " + err.Error())
		return err
	}
	_, err = db1.Query("UPDATE cuenta SET CLAVE = ? where EMAIL = ?", nuevasClaves.Clave1, cuentaAcutal.Email)
	if err != nil {
		db1.Rollback()
		log.Println("Error: " + err.Error())
		return err
	}
	err = db1.Commit()
	if err != nil {
		db1.Rollback()
		log.Println("Error: " + err.Error())
		return err
	}

	return nil
}

func buscarEquiposActivos(ID_Quest string, Representante string) ([]Equipo, error) {
	equiposEncontrados := []Equipo{}
	db1, err := bd.Begin()
	if err != nil {
		log.Println("Error: " + err.Error())
		return equiposEncontrados, err
	}
	tab1, err := db1.Query("SELECT e.ID_EQUIPO, e.NOMBRE_EQUIPO, e.RUT_RESPONSABLE, e.NOMBRE_RESPONSABLE FROM equipo e WHERE e.ID_QUEST = ? AND e.ID_REPRESENTANTE = ?", ID_Quest, Representante)
	if err != nil {
		db1.Rollback()
		log.Println("Error: " + err.Error())
		return equiposEncontrados, err
	}
	defer tab1.Close()
	for tab1.Next() {
		equipo := Equipo{
			ID_Quest: ID_Quest,
		}
		err = tab1.Scan(&equipo.ID_Equipo, &equipo.Nombre_Equipo, &equipo.Rut_Respondable, &equipo.Nombre_Respondable)
		if err != nil {
			db1.Rollback()
			log.Println("Error: " + err.Error())
			return equiposEncontrados, err
		}
		miembros, err := buscarMiembro(equipo.ID_Equipo)
		if err != nil {
			db1.Rollback()
			log.Println("Error: " + err.Error())
			return equiposEncontrados, err
		}
		equipo.Miembros_Equipo = miembros
		equiposEncontrados = append(equiposEncontrados, equipo)
	}
	err = db1.Commit()

	if err != nil {
		log.Println("Error: " + err.Error())
		return equiposEncontrados, err
	}

	return equiposEncontrados, nil
}

func buscarMiembro(ID_Equipo string) ([]Miembro, error) {
	miembros := []Miembro{}
	db1, err := bd.Begin()
	if err != nil {
		log.Println("Error: " + err.Error())
		return miembros, err
	}
	tab1, err := db1.Query("SELECT m.ID_MIEMBRO, m.NOMBRE_MIEMBRO, m.Rut_MIEMBRO FROM miembro m WHERE m.ID_Equipo = ?", ID_Equipo)
	if err != nil {
		db1.Rollback()
		log.Println("Error: " + err.Error())
		return miembros, err
	}
	defer tab1.Close()
	for tab1.Next() {
		miembro := Miembro{
			ID_Equipo: ID_Equipo,
		}
		err = tab1.Scan(&miembro.ID_Miembro, &miembro.Nombre_Miembro, &miembro.Rut_Miembro)
		if err != nil {
			db1.Rollback()
			log.Println("Error: " + err.Error())
			return miembros, err
		}
		miembros = append(miembros, miembro)
	}
	err = db1.Commit()
	if err != nil {
		log.Println("Error: " + err.Error())
		return miembros, err
	}
	return miembros, nil
}

func buscarQuest(ID_Quest string) (Quest, error) {
	quest := Quest{
		ID_Quest: ID_Quest,
	}
	db1, err := bd.Begin()
	if err != nil {
		log.Println("Error: " + err.Error())
		return quest, err
	}
	tab1, err := db1.Query("SELECT e.NOMBRE_QUEST, e.ESTADO_QUEST, e.FECHA_INICIO, e.PREMIO, e.GANADOR, e.CATEGORIA FROM quest e WHERE e.ID_QUEST = ? ", ID_Quest)
	if err != nil {
		db1.Rollback()
		log.Println("Error: " + err.Error())
		return quest, err
	}
	defer tab1.Close()
	for tab1.Next() {
		err = tab1.Scan(&quest.Nombre_Quest, &quest.Estado, &quest.FechaInicio, &quest.Premio, &quest.Ganador, &quest.Categoria)
		if err != nil {
			db1.Rollback()
			log.Println("Error: " + err.Error())
			return quest, err
		}
		preguntas, err := buscarPreguntas(quest.ID_Quest)
		if err != nil {
			db1.Rollback()
			log.Println("Error: " + err.Error())
			return quest, err
		}
		quest.Preguntas = preguntas
	}
	err = db1.Commit()

	if err != nil {
		log.Println("Error: " + err.Error())
		return quest, err
	}

	return quest, nil
}

func buscarPreguntas(ID_Quest string) ([]Pregunta, error) {
	preguntas := []Pregunta{}
	db1, err := bd.Begin()
	if err != nil {
		log.Println("Error: " + err.Error())
		return preguntas, err
	}
	tab1, err := db1.Query("SELECT ID_PREGUNTA, PISTA, PREGUNTA, RESPUESTA FROM pregunta WHERE ID_Quest = ?", ID_Quest)
	if err != nil {
		db1.Rollback()
		log.Println("Error: " + err.Error())
		return preguntas, err
	}
	defer tab1.Close()
	for tab1.Next() {
		pregunta := Pregunta{
			ID_Quest: ID_Quest,
		}
		err = tab1.Scan(&pregunta.ID_Pregunta, &pregunta.Pista, &pregunta.Pregunta, &pregunta.Respuesta)
		if err != nil {
			db1.Rollback()
			log.Println("Error: " + err.Error())
			return preguntas, err
		}
		preguntas = append(preguntas, pregunta)
	}
	err = db1.Commit()

	if err != nil {
		log.Println("Error: " + err.Error())
		return preguntas, err
	}
	return preguntas, nil
}

func limpiaRespueestas(preguntas []Pregunta) []Pregunta {
	preguntasn := []Pregunta{}
	for _, element := range preguntas {
		element.Respuesta = ""
		preguntasn = append(preguntasn, element)
	}
	return preguntasn
}

func enviarmail(contenido string, destin string, motivo string) error {
	// Set up authentication information.
	log.Println(contenido)
	auth := smtp.PlainAuth(
		"",
		param[14],
		param[15],
		param[13],
	)
	// Connect to the server, authenticate, set the sender and recipient,
	msg := "From: " + param[14] + "\n" +
		"To: " + destin + "\n" + "Subject: " + motivo + " :  \n\n" + contenido
	// and send the email all in one step.
	err := smtp.SendMail(
		param[13]+param[16],
		auth,
		param[15],
		[]string{destin},
		[]byte(msg),
	)
	if err != nil {
		log.Println("-->ERROR al Enviar Email :" + err.Error())
		//return err //////////////////////////////////////////////////////////////////esta linea debe ser reactivada
	}
	log.Println("Email Enviado a :" + destin)
	return nil
}
