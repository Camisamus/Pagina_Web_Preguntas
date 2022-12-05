const urlParams = new URLSearchParams(window.location.search.substring(1));
const Inicios = {
    "*********": () => {},
    "Main.html": () => {},
    "icio.html": () => {
        cerrarSesion()
    },
    "enta.html": () => {},
    "enta.html": () => {},
    "lave.html": () => {},
    "ave1.html": () => {},
    "reso.html": () => {
        var quests = new sQuests();
        quests.consultar();
    },
    "uest.html": () => {
        var quests = new sQuests();
        quests.consultar();
        var quest = new sQuestsDetalle();
        quest.consultar();
    },

};
var sec = 0;
var min = 0;
var t;
$(document).ready(function() {
    return Inicios[window.location.href.replace(window.location.search.substring(0), "").slice(-9)]();
});

function CrearCuenta() {
    if ($("#password").val() == $("#password2").val()) {
        if ($("#cuenta").val() == "" || $("#email").val() == "" || $("#password").val() == "" || $("#password2").val() == "") { alert("Formulario incompeto!"); return false; }
        var crearCuenta = new sCrearCuenta();
        crearCuenta.consultar();
        return
    }
    alert("Claves no coinciden!");
}

function iniciarsesion() {
    if ($("#email").val() == "" || $("#password").val() == "") { alert("Formulario incompeto!"); return false; }
    var iniciarSesion = new sIngresar();
    iniciarSesion.consultar();

}

function cerrarSesion() {
    var salida = new sSalir()
    salida.consultar();
}

function recuperarClave1() {
    var RecuperarClave1 = new sRecuperarClave1();
    RecuperarClave1.consultar();
}

function recuperarclave2() {
    var RecuperarClave2 = new sRecuperarClave2();
    RecuperarClave2.consultar();
}


function sQuests() {

    this.source = 'http://localhost:8090/Quests';

    this.callback = null;
    this.extra = null;

    this.consultar = function() {
        var data = ""
        var that = this;

        $.ajax({
            url: this.source,
            xhrFields: {
                withCredentials: true,
                credentials: 'include'
            },
            data: "",
            method: 'POST',
            success: function(data) {
                if (data.Sesion == "Cerrada") { window.location.href = "../paginas/inicio.html"; return; }
                crearMenu(data);
            },
            error: function(data) {
                window.location.href = "../paginas/error.html"
            },
            async: true
        });
    };
}

function sQuestsDetalle() {

    this.source = 'http://localhost:8090/Quest';

    this.callback = null;
    this.extra = null;

    this.consultar = function() {
        var data = ""
        var that = this;

        $.ajax({
            url: this.source,
            xhrFields: {
                withCredentials: true,
                credentials: 'include'
            },
            data: JSON.stringify({
                IDQuest: urlParams.get('Quest'),
            }),
            method: 'POST',
            success: function(data) {
                console.log(data);
                armarPaginaQuest0(data)
            },
            error: function(data) {
                window.location.href = "../paginas/error.html"
            },
            async: true
        });
    };
}



function sCrearCuenta() {

    this.source = 'http://localhost:8090/CrearCuenta';

    this.callback = null;
    this.extra = null;

    this.consultar = function() {
        var data = ""
        var that = this;

        $.ajax({
            url: this.source,
            data: JSON.stringify({
                IDCuenta: "",
                NombreCuenta: $("#cuenta").val(),
                Email: $("#email").val(),
                Clave1: $("#password").val(),
                Clave2: $("#password2").val(),
                Estado: "",
            }),
            method: 'POST',
            success: function(data) {
                if (data.Sesion == "Cerrada" && location.href.slice(-11) != 'inicio.html') { window.location.href = "../paginas/inicio.html" }
            },
            error: function(data) {
                window.location.href = "../paginas/error.html"
            },
            async: true
        });
    };
}

function sIngresar() {

    this.source = 'http://localhost:8090/IniciarSesion';

    this.callback = null;
    this.extra = null;

    this.consultar = function() {
        var data = ""
        var that = this;

        $.ajax({
            url: this.source,
            xhrFields: {
                withCredentials: true,
            },
            data: JSON.stringify({
                Email: $("#email").val(),
                Clave1: $("#password").val(),
                Estado: "",
            }),
            method: 'POST',
            success: function(data) {
                if (data.Estado == "True" && location.href.slice(-12) != 'ingreso.html') {
                    document.cookie = data.Clave1 + "=" + data.Clave2;
                    window.location.href = "../paginas/ingreso.html";
                }
                if (data.Estado != "True") { alert(data.Estado) }
            },
            error: function(data) {
                window.location.href = "../paginas/error.html"
            },
            async: true
        });
    };
}

function sSalir() {

    this.source = 'http://localhost:8090/CerrarSesion';

    this.callback = null;
    this.extra = null;

    this.consultar = function() {
        var data = ""
        var that = this;

        $.ajax({
            url: this.source,
            xhrFields: {
                withCredentials: true,
            },
            data: {},
            method: 'POST',
            success: function(data) {
                if (data.Sesion == "Cerrada" && location.href.slice(-11) != 'inicio.html') { window.location.href = "../paginas/inicio.html" }
            },
            error: function(data) {
                window.location.href = "../paginas/error.html"
            },
            async: true
        });
    };
}


function sRecuperarClave1() {

    this.source = 'http://localhost:8090/RecuperarClave1';

    this.callback = null;
    this.extra = null;

    this.consultar = function() {
        var data = ""
        var that = this;

        $.ajax({
            url: this.source,
            data: JSON.stringify({
                Email: $("#email").val(),
            }),
            method: 'POST',
            success: function(data) {
                window.location.href = "../paginas/inicio.html"
            },
            error: function(data) {
                window.location.href = "../paginas/error.html"
            },
            async: true
        });
    };
}

function sRecuperarClave2() {

    this.source = 'http://localhost:8090/RecuperarClave2';

    this.callback = null;
    this.extra = null;

    this.consultar = function() {
        var data = ""
        var that = this;

        $.ajax({
            url: this.source,
            data: JSON.stringify({
                Clave1: $("#password").val(),
                Clave2: $("#password2").val(),
                Estado: urlParams.get('ClaveCambio'),
            }),
            method: 'POST',
            success: function(data) {
                alert(data.Estado);
                window.location.href = "../paginas/inicio.html"
            },
            error: function(data) {
                window.location.href = "../paginas/error.html"
            },
            async: true
        });
    };
}

function sInscribirEquipo() {

    this.source = 'http://localhost:8090/Inscribirse';

    this.callback = null;
    this.extra = null;

    this.consultar = function(p) {
        var data = ""
        var that = this;

        $.ajax({
            url: this.source,
            data: JSON.stringify({
                Respuestas: [{
                    IDPregunta: "" + p,
                    IDQuest: urlParams.get('Quest'),
                    Pregunta: {
                        IDPregunta: "" + p,
                        Respuesta: $("#Respuesta" + p).val(),
                    },
                }],
                IDEquipo: $("#inputEquipos").val(),
                IDQuest: urlParams.get('Quest'),
            }),
            method: 'POST',
            success: function(data) {
                if (data.IDEquipo == "Gano") {
                    alert("GANASTE!!!")
                    return
                }
                if (data.IDEquipo == "Aun no se permite volver a responder") {
                    alert(data.IDEquipo)
                    return
                }
                if (!data.Respuestas[0].Correcta) {
                    partir();
                    tiempo();
                    alert("Respuesta Incorrecta. siquiente oportunidad en 20 minutos.")
                    return
                }
            },
            error: function(data) {
                window.location.href = "../paginas/error.html"
            },
            async: true
        });
    };
}

function sRespuesta() {

    this.source = 'http://localhost:8090/EnviarRespuesta';

    this.callback = null;
    this.extra = null;

    this.consultar = function(p) {
        var data = ""
        var that = this;

        $.ajax({
            url: this.source,
            data: JSON.stringify({
                Respuestas: [{
                    IDPregunta: "" + p,
                    IDQuest: urlParams.get('Quest'),
                    Pregunta: {
                        IDPregunta: "" + p,
                        Respuesta: $("#Respuesta" + p).val(),
                    },
                }],
                IDEquipo: $("#inputEquipos").val(),
                IDQuest: urlParams.get('Quest'),
            }),
            method: 'POST',
            success: function(data) {
                if (data.IDEquipo == "Gano") {
                    alert("GANASTE!!!")
                    return
                }
                if (data.IDEquipo == "Aun no se permite volver a responder") {
                    alert(data.IDEquipo)
                    return
                }
                if (!data.Respuestas[0].Correcta) {
                    partir();
                    tiempo();
                    alert("Respuesta Incorrecta. siquiente oportunidad en 20 minutos.")
                    return
                }
            },
            error: function(data) {
                window.location.href = "../paginas/error.html"
            },
            async: true
        });
    };
}

function crearMenu(data) {
    let lista = document.getElementById("navbardiv");
    let cerrarSesion = document.createElement("input");
    cerrarSesion.setAttribute("type", "button");
    cerrarSesion.setAttribute("value", "Cerrar Sesion");
    cerrarSesion.setAttribute("onclick", "cerrarSesion()");
    lista.appendChild(cerrarSesion);
    let ul = document.createElement("ul");
    let liIngreso = document.createElement("li");
    let aIngreso = document.createElement("a");
    let textoIngreso = document.createTextNode("Ingreso");
    aIngreso.title = "Volver";
    aIngreso.href = "../paginas/ingreso.html";
    aIngreso.appendChild(textoIngreso);
    liIngreso.appendChild(aIngreso);
    ul.appendChild(liIngreso);
    for (var i = 0; i < data.length; i++) {
        let li = document.createElement("li");
        let a = document.createElement("a");
        let texto = document.createTextNode(data[i].NombreQuest);
        a.title = data[i].NombreQuest;
        a.href = "../paginas/quest.html?Quest=" + data[i].IDQuest;
        a.appendChild(texto);
        li.appendChild(a);
        ul.appendChild(li);
    }
    lista.appendChild(ul);
}

function armarPaginaQuest0(data) {
    let bodydiv = document.getElementById("contentdiv");
    let divDatos = llebarDivDatos(data);
    let divEquipos = llebarDivEquipos(data);
    let divpreguntas = llebarDivpreguntas(data);
    reloj = document.createElement("h1")
    reloj.setAttribute("id", "temporisador")
    timer = document.createElement("time")
    timer.appendChild(document.createTextNode("00:00"))
    reloj.appendChild(timer)
    divEquipos.appendChild(reloj)
    bodydiv.appendChild(divDatos);
    bodydiv.appendChild(divEquipos);
    cambioEquipo()
    bodydiv.appendChild(divpreguntas);
}

function llebarDivDatos(data) {
    let divDatos = document.createElement("div");
    let titulo = document.createElement("h3");
    text = document.createTextNode(data.Quest.NombreQuest);
    titulo.appendChild(text);

    let btninscribir = document.createElement("input");
    btninscribir.setAttribute("type", "button")
    btninscribir.setAttribute("value", "Inscribir")
    btninscribir.setAttribute("id", "btnInscribir")
    btninscribir.setAttribute("onclick", "inscribir()")
    divDatos.appendChild(titulo);
    divDatos.appendChild(btninscribir);
    return divDatos
}

function llebarDivEquipos(data) {
    let divEquipos = document.createElement("div");
    var inputEquipos = document.createElement("SELECT");
    inputEquipos.setAttribute("id", "inputEquipos");
    inputEquipos.setAttribute("onchange", "cambioEquipo()");
    divEquipos.appendChild(inputEquipos);
    for (i = 0; i < data.Equipo.length; i++) {
        opcion = document.createElement("option");
        opcion.setAttribute("value", data.Equipo[i].IDEquipo);
        opcion.appendChild(document.createTextNode(data.Equipo[i].NombreRespondable + " - " + data.Equipo[i].RutRespondable));
        inputEquipos.appendChild(opcion)
        divEquipo0 = document.createElement("div");
        divEquipo0.setAttribute("id", "divEquipo" + data.Equipo[i].IDEquipo);
        divEquipo0.setAttribute("Class", "divEquipo");
        divEquipo0.setAttribute("Style", "display: none")
        divEquipo1 = document.createElement("div");
        divEquipo2 = document.createElement("div");
        divEquipo3 = document.createElement("div");
        text1 = document.createTextNode(data.Equipo[i].NombreEquipo);
        divEquipo1.appendChild(text1);
        text2 = document.createTextNode(data.Equipo[i].NombreRespondable + " ");
        divEquipo2.appendChild(text2);
        text3 = document.createTextNode(data.Equipo[i].RutRespondable);
        divEquipo2.appendChild(text3);
        for (j = 0; j < data.Equipo[i].Miembros_Equipo.length; j++) {
            text4 = document.createTextNode(data.Equipo[i].Miembros_Equipo[j].NombreMiembro + " ");
            divEquipo3.appendChild(text4);
            text5 = document.createTextNode(data.Equipo[i].Miembros_Equipo[j].RutMiembro);
            divEquipo3.appendChild(text5);
            divEquipo3.appendChild(document.createElement("br"));
        }
        divEquipo0.appendChild(divEquipo1);
        divEquipo0.appendChild(divEquipo2);
        divEquipo0.appendChild(divEquipo3);
        divEquipos.appendChild(divEquipo0);
    }
    return divEquipos
}

function llebarDivpreguntas(data) {
    let divpreguntas = document.createElement("div");
    for (i = 0; i < data.Quest.Preguntas.length; i++) {
        divPregunta0 = document.createElement("div");
        divPregunta0.setAttribute("class", "inputter");
        divPregunta0.appendChild(document.createTextNode("Pregunta " + (i + 1) + ": " + data.Quest.Preguntas[i].Pregunta))
        divPregunta1 = document.createElement("div");
        divPregunta1.setAttribute("class", "inputter");
        divPregunta1.appendChild(document.createTextNode("Pistas : " + data.Quest.Preguntas[i].Pista))
        divPregunta2 = document.createElement("div");
        divPregunta2.setAttribute("class", "inputter");
        camporespuesta = document.createElement("input")
        camporespuesta.setAttribute("type", "text")
        camporespuesta.setAttribute("id", "Respuesta" + data.Quest.Preguntas[i].IDPregunta)
        divPregunta2.appendChild(camporespuesta)
        divPregunta3 = document.createElement("div");
        divPregunta3.setAttribute("class", "inputter");
        botonRespuesta = document.createElement("input")
        botonRespuesta.setAttribute("type", "button")
        botonRespuesta.setAttribute("value", "Enviar Respuesta")
        botonRespuesta.setAttribute("onclick", "enviarRespuesta(" + data.Quest.Preguntas[i].IDPregunta + ")")
        divPregunta3.appendChild(botonRespuesta)
        divpreguntas.appendChild(divPregunta0);
        divpreguntas.appendChild(divPregunta1);
        divpreguntas.appendChild(divPregunta2);
        divpreguntas.appendChild(divPregunta3);
    }
    return divpreguntas
}
numeroDeMiembros = 0

async function inscribir() {
    numeroDeMiembros = 0



    fondomodal = document.createElement("div");
    fondomodal.setAttribute("class", "fondoModal");
    fondomodal.setAttribute("id", "fondoModal");
    fondomodal.setAttribute("onclick", "cerrarModal()");
    divInscribirse0 = document.createElement("div");
    divInscribirse0.setAttribute("id", "divinscribirse0");
    divInscribirse0.setAttribute("style", "    background-color: rgba(255, 255, 255, 1);    margin: 15% auto;    border-radius: 10px;    padding: 15px;    border: 1px solid #888;    width: 80%;    min-width: 400px;    max-width: 600px;    color: black;");
    divInscribirse0.setAttribute("onclick", "event.stopPropagation();");
    divInscribirse1 = document.createElement("div");
    divInscribirse1.setAttribute("class", "inputter");
    divInscribirse2 = document.createElement("div");
    divInscribirse2.setAttribute("class", "inputter");
    inputRut = document.createElement("input");
    inputRut.setAttribute("id", "inputRutResponsable")
    inputRut.setAttribute("style", "margin: 10px");
    inputNombre = document.createElement("input");
    inputNombre.setAttribute("id", "inputNombreResponsable")
    inputNombre.setAttribute("style", "margin: 10px");
    botonAgregarMiembro = document.createElement("input");
    botonAgregarMiembro.setAttribute("type", "button");
    botonAgregarMiembro.setAttribute("style", "margin: 10px");
    botonAgregarMiembro.setAttribute("value", "Agregar Miembro");
    botonAgregarMiembro.setAttribute("onclick", "agregarCampoMiembro()");
    botonEnviarEquipo = document.createElement("input");
    botonEnviarEquipo.setAttribute("type", "button");
    botonEnviarEquipo.setAttribute("style", "margin: 10px");
    botonEnviarEquipo.setAttribute("value", "Enviar Equipo");
    botonEnviarEquipo.setAttribute("onclick", "EnviarEquipo()");
    divInscribirse0.appendChild(divInscribirse1);
    divInscribirse1.appendChild(document.createTextNode("Rut Responsable"));
    divInscribirse1.appendChild(inputRut);
    divInscribirse1.appendChild(document.createElement("br"));
    divInscribirse1.appendChild(document.createTextNode("Nombre Responsable"));
    divInscribirse1.appendChild(inputNombre);
    divInscribirse0.appendChild(divInscribirse2);
    divInscribirse2.appendChild(botonAgregarMiembro);
    divInscribirse2.appendChild(botonEnviarEquipo);
    fondomodal.appendChild(divInscribirse0);

    document.body.appendChild(fondomodal);

}

function agregarCampoMiembro() {
    numeroDeMiembros++
    divInscribirse0 = document.getElementById("divinscribirse0");
    divNuevoMiembro = document.createElement("div");
    divNuevoMiembro.setAttribute("class", "inputter");
    divNuevoMiembro.appendChild(document.createTextNode("Rut Miembro"));
    inputRut = document.createElement("input");
    inputRut.setAttribute("type", "text");
    inputRut.setAttribute("id", "inputRutMiembro_" + numeroDeMiembros)
    inputRut.setAttribute("style", "margin: 10px");
    divNuevoMiembro.appendChild(inputRut);
    divNuevoMiembro.appendChild(document.createElement("br"));
    divNuevoMiembro.appendChild(document.createTextNode("Nombre Miembro"));
    inputNombre = document.createElement("input");
    inputNombre.setAttribute("type", "text")
    inputNombre.setAttribute("id", "inputNombreMiembro_" + numeroDeMiembros)
    inputNombre.setAttribute("style", "margin: 10px");
    divNuevoMiembro.appendChild(inputNombre);

    divInscribirse0.appendChild(divNuevoMiembro)
}

function cerrarModal() {
    modal = document.getElementById("fondoModal")
    modal.innerHTML = ""
    modal.remove();
}

function enviarRespuesta(IDPregunta) {
    var respuesta = new sRespuesta;
    respuesta.consultar(IDPregunta);
}

function cambioEquipo() {
    divsEquipos = document.getElementsByClassName("divEquipo")
    for (var i = 0; i < divsEquipos.length; i++) {
        divsEquipos[i].setAttribute("Style", "display: none")
    }
    seleleccionado = document.getElementById("divEquipo" + document.getElementById("inputEquipos").value)
    seleleccionado.setAttribute("Style", "")

}

function tick() {
    sec--;
    if (sec < 0) {
        sec = 59;
        min--
        if (min < 0) {
            clearTimeout(t)
        }
    }
}


function partir() {
    document.getElementById("temporisador").textContent = "20:00";
    min = 20;
    sec = 0;
}

function tiempo() {
    if (min <= 0 && sec <= 0) {
        partir()
        return;
    }
    tick();
    document.getElementById("temporisador").textContent = min + ":" + (sec > 9 ? sec : "0" + sec)
    t = setTimeout(tiempo, 1000);
}