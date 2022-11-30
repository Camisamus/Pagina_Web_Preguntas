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