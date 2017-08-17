
function init () {
    var nums = "0123456789";
    function generateNick() {
        var str = "";
        for (var i = 0; i < 5; i++) {
            str += nums[Math.floor(Math.random() * nums.length)]
        }
        return str;
    }
    var nick = generateNick();

    if (window["WebSocket"]) {
        var conn = new WebSocket("ws://" + document.location.host + "/ws");
        conn.onclose = function (evt) {
            console.log("CLOSE");
            console.log(evt);
        };
        conn.onmessage = function (evt) {
            console.log(JSON.parse(evt.data));
        };
        var msg = JSON.stringify({
            type: "message",
            content: "boulou",
            message: "Bonjour, je suis " + nick
        });
        document.getElementById("hw").addEventListener("click", function () {
            conn.send(msg);
        })
    } else {
        console.log("WEBSOCKET NOT SUPPORTED");
    }
}

document.addEventListener("DOMContentLoaded", init);
