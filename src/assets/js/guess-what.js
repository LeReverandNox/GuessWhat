function toDataURL(src, callback, outputFormat) {
    var img = new Image();
    img.crossOrigin = 'Anonymous';
    img.onload = function () {
        var canvas = document.createElement('CANVAS');
        var ctx = canvas.getContext('2d');
        var dataURL;
        canvas.height = this.naturalHeight;
        canvas.width = this.naturalWidth;
        ctx.drawImage(this, 0, 0);
        dataURL = canvas.toDataURL(outputFormat);
        callback(dataURL);
    };
    img.src = src;
    if (img.complete || img.complete === undefined) {
        img.src = "data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///ywAAAAAAQABAAACAUwAOw==";
        img.src = src;
    }
}

function init () {
    if (window["WebSocket"]) {
        var conn = new WebSocket("ws://" + document.location.host + "/ws");
        conn.onclose = function (evt) {
            console.log("Socket closed");
        };
        conn.onmessage = function (evt) {
            console.log(JSON.parse(evt.data));
        };

        var imgSrc = "https://www.gravatar.com/avatar/d50c83cc0c6523b4d3f6085295c953e0";
        var msg = JSON.stringify({
            action: "message",
            nickname: "boulou",
            message: "Bonjour, je suis ",
            x: "10",
            y: "10"

        });
        document.getElementById("hw").addEventListener("click", function () {
            conn.send(msg);
        })
        document.getElementById("change_nickname").addEventListener("click", function () {
            var msg = JSON.stringify({
                action: "set_nickname",
                nickname: document.getElementById("nickname").value
            })
            conn.send(msg);
        })
        document.getElementById("send_message").addEventListener("click", function () {
            var msg = JSON.stringify({
                action: "send_message",
                content: document.getElementById("message").value
            })
            conn.send(msg);
        })
    } else {
        console.log("WEBSOCKET NOT SUPPORTED");
    }
}

document.addEventListener("DOMContentLoaded", function () {
    init();
});
