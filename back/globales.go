package main

import (
	"database/sql"
	"time"
)

var param []string
var questActivas []QuestMenu
var bd *sql.DB
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

//Representeante, Usuario dueño de la cuenta y Representante de los equipos
type Representante struct {
	ID_Representante     string `json:"IDQuest"`
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
