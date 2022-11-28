const Inicios = {
    "Main.html": () => {},
    "icio.html": () => {},
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
    },

};
$(document).ready(function() {
    return Inicios[location.href.slice(-9)]();
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
    var iniciarSesion = new Singresar();
    iniciarSesion.consultar();

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
            data: {},
            method: 'POST',
            success: function(data) {
                if (data.Estado == "True") { window.location.href = "../paginas/ingreso.html" }

                //alert(data.Sesion);
                //_data = jQuery.parseJSON(data);
                //if (that.callback) that.callback(_data, that.extra);
            },
            error: function(data) {
                //alert(data);
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

                //alert(data.Sesion);
                //_data = jQuery.parseJSON(data);
                //if (that.callback) that.callback(_data, that.extra);
            },
            error: function(data) {
                //alert(data);
                window.location.href = "../paginas/error.html"
            },
            async: true
        });
    };
}

function Singresar() {

    this.source = 'http://localhost:8090/IniciarSesion';

    this.callback = null;
    this.extra = null;

    this.consultar = function() {
        var data = ""
        var that = this;

        $.ajax({
            url: this.source,
            data: JSON.stringify({
                Email: $("#email").val(),
                Clave1: $("#password").val(),
                Estado: "",
            }),
            method: 'POST',
            success: function(data) {
                if (data.Sesion == "Cerrada" && location.href.slice(-11) != 'inicio.html') { window.location.href = "../paginas/inicio.html" }

                //alert(data.Sesion);
                //_data = jQuery.parseJSON(data);
                //if (that.callback) that.callback(_data, that.extra);
            },
            error: function(data) {
                //alert(data);
                window.location.href = "../paginas/error.html"
            },
            async: true
        });
    };
}