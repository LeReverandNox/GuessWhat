/*jslint browser this */
/*global $ window URL WebSocket console */

(function (global) {
    "use strict";
    var GuessWhat = function () {
        this.isConnected = false;
        this.tool = {
            color: "#000000",
            thickness: 7
        };
        this.lastX;
        this.lastY;
        this.secretWord = "";
    };

    GuessWhat.prototype.startEventListeners = function () {
        this.connexionButton.addEventListener("click", this.connect.bind(this));
        this.sendMessageButton.addEventListener("click", this.sendMessage.bind(this));
        this.joinRoomButton.addEventListener("click", this.joinRoom.bind(this));
        this.leaveRoomButton.addEventListener("click", this.leaveRoom.bind(this));
        this.startRoomButton.addEventListener("click", this.startRoom.bind(this));
        this.colorHolder.addEventListener("click", this.colorClick.bind(this));
        this.thicknessHolder.addEventListener("click", this.thicknessClick.bind(this));
        this.toolsHolder.addEventListener("click", this.toolsClick.bind(this));

        this.canvas.addEventListener("mousedown", this.onLocalCanvasMouseDown.bind(this));
        this.canvas.addEventListener("mousemove", this.onLocalCanvasMouseMove.bind(this));
        this.canvas.addEventListener("mouseup", this.onLocalCanvasMouseUp.bind(this));
        this.canvas.addEventListener("mouseout", this.onLocalCanvasMouseOut.bind(this));
    };

    GuessWhat.prototype.registerElements = function () {
        this.connexionButton = document.getElementById("connexion_button");
        this.nicknameInput = document.getElementById("nickname");

        this.sendMessageButton = document.getElementById("send_message");
        this.messageInput = document.getElementById("message");

        this.joinRoomButton = document.getElementById("join_room");
        this.leaveRoomButton = document.getElementById("leave_room");
        this.startRoomButton = document.getElementById("start_room");
        this.roomInput = document.getElementById("room");

        this.colorHolder = document.getElementById("color_holder");
        this.thicknessHolder = document.getElementById("thickness_holder");
        this.toolsHolder = document.getElementById("tools_holder");

        this.secretWordP = document.getElementById("secret_word");

        this.canvas = document.getElementById("canvas");
        this.context = this.canvas.getContext("2d");
        this.context.lineCap = "round";
        this.context.lineJoin = "round";
    };

    GuessWhat.prototype.connect = function () {
        if (this.isConnected)
            return false

        var nickname = this.nicknameInput.value;
        if (window["WebSocket"]) {
            this.socket = new WebSocket("ws://" + document.location.host + "/ws?nickname=" + nickname);
            this.socket.onclose = this.onClose.bind(this);
            this.socket.onmessage = this.onMessage.bind(this);
            this.isConnected = true;
        } else {
            console.log("WEBSOCKET NOT SUPPORTED");
        }
    };

    GuessWhat.prototype.sendMessage = function () {
        if (!this.isConnected)
            return false

        var msg = JSON.stringify({
            action: "send_message",
            content: this.messageInput.value
        });
        this.socket.send(msg);
    };

    GuessWhat.prototype.joinRoom = function () {
        if (!this.isConnected)
            return false

        var msg = JSON.stringify({
            action: "join_room",
            room: this.roomInput.value
        });
        this.socket.send(msg);
    };

    GuessWhat.prototype.leaveRoom = function () {
        if (!this.isConnected)
            return false

        var msg = JSON.stringify({
            action: "leave_room",
            room: this.roomInput.value
        });
        this.socket.send(msg);
    };

    GuessWhat.prototype.startRoom = function () {
        if (!this.isConnected)
            return false

        var msg = JSON.stringify({
            action: "start_room",
            room: this.roomInput.value
        });
        this.socket.send(msg);
    };

    GuessWhat.prototype.onMessage = function (e) {
        var data = JSON.parse(e.data)
        var action = data.action;
        switch (action) {
            case "canvas_mouse_down":
                this.onCanvasMouseDown(data)
                break;
            case "canvas_mouse_move":
                this.onCanvasMouseMove(data)
                break;
            case "canvas_mouse_up":
                this.onCanvasMouseUp(data)
                break;
            case "incoming_room_image":
                this.onIncomingRoomImage(data)
                break;
            case "ask_for_image":
                this.onAskForImage(data)
                break;
            case "join_room_cb":
                this.onJoinRoomCb(data);
                break;
            case "leave_room_cb":
                this.onLeaveRoomCb(data);
                break;
            case "new_round_start":
                this.onNewRoundStart(data);
                break;
            case "clean_canvas":
                this.onCleanCanvas(data);
                break;
            case "new_round_start":
                this.onNewRoundStart(data);
                break;
            case "reveal_letter":
                this.onRevealLetter(data);
                break;
            case "you_are_drawing":
                this.onYouAreDrawing(data);
                break;
            default:
                console.log(data);
                break;
        }
    };

    GuessWhat.prototype.onClose = function (e) {
        this.isConnected = false;
        console.log("Socket closed");
    };

    GuessWhat.prototype.onLocalCanvasMouseDown = function (e) {
        if (!this.isConnected)
            return false

        this.localClick1 = true;
        this.lastX = e.offsetX;
        this.lastY = e.offsetY;

        var msg = JSON.stringify({
            action: "canvas_mouse_down",
            room: this.roomInput.value,
            toX: String(this.lastX),
            toY: String(this.lastY),
            thickness: String(this.tool.thickness),
            color: this.tool.color
        });
        this.socket.send(msg);
    };
    GuessWhat.prototype.onLocalCanvasMouseMove = function (e) {
        if (!this.isConnected)
            return false
        if (!this.localClick1)
            return false

        var toX = e.offsetX;
        var toY = e.offsetY
        var msg = JSON.stringify({
            action: "canvas_mouse_move",
            room: this.roomInput.value,
            fromX: String(this.lastX),
            fromY: String(this.lastY),
            toX: String(toX),
            toY: String(toY),
            thickness: String(this.tool.thickness),
            color: this.tool.color
        });

        this.socket.send(msg);

        this.lastX = toX;
        this.lastY = toY;
    };
    GuessWhat.prototype.onLocalCanvasMouseUp = function (e) {
        if (!this.isConnected)
            return false
        if (!this.localClick1)
            return false

        this.localClick1 = false;
    };
    GuessWhat.prototype.onLocalCanvasMouseOut = function (e) {
        if (!this.isConnected)
            return false
        if (!this.localClick1)
            return false

        var clean = function () {
            document.removeEventListener("mouseup", mouseUpHandler);
            this.canvas.removeEventListener("mouseenter", mouseEnterHandler);
        };
        var mouseUpHandler = function () {
            clean.call(this);
            this.onLocalCanvasMouseUp(e);

        };
        var mouseEnterHandler = function () {
            clean.call(this);
        };

        document.addEventListener("mouseup", mouseUpHandler.bind(this));
        this.canvas.addEventListener("mouseenter", mouseEnterHandler.bind(this));
    };

    GuessWhat.prototype.onCanvasMouseDown = function (e) {
        this.context.strokeStyle = e.color;
        this.context.fillStyle = e.color;
        this.context.lineWidth = e.thickness;

        this.context.beginPath();
        this.context.arc(e.toX, e.toY, e.thickness / 2, 0, 2 * Math.PI, true);
        this.context.fill();
    };

    GuessWhat.prototype.onCanvasMouseMove = function (e) {
        this.context.strokeStyle = e.color;
        this.context.fillStyle = e.color;
        this.context.lineWidth = e.thickness;

        this.context.beginPath();
        this.context.moveTo(e.fromX, e.fromY);
        this.context.lineTo(e.toX, e.toY);
        this.context.stroke();
    };

    GuessWhat.prototype.onCleanCanvas = function (e) {
        this.cleanCanvas();
    }

    GuessWhat.prototype.cleanCanvas = function () {
        this.context.clearRect(0, 0, this.canvas.width, this.canvas.height);
    };

    GuessWhat.prototype.init = function () {
        var self = this;
        this.registerElements();
        this.startEventListeners();
    };

    GuessWhat.prototype.getCanvasBase64 = function () {
        var b64 = this.canvas.toDataURL();
        return b64;
    };

    GuessWhat.prototype.drawBase64ToCanvas = function (base64) {
        var img = new Image();
        img.src = base64;
        var onLoadCb = function () {
            this.context.drawImage(img, 0, 0);
        };

        img.onload = onLoadCb.bind(this)
    };

    GuessWhat.prototype.onIncomingRoomImage = function (e) {
        this.drawBase64ToCanvas(e.room.Image);
    }

    GuessWhat.prototype.onAskForImage = function (e) {
        var msg = JSON.stringify({
            action: "send_image",
            room: this.roomInput.value,
            image: this.getCanvasBase64()
        });
        this.socket.send(msg);
    };

    GuessWhat.prototype.onJoinRoomCb = function (e) {
        this.cleanCanvas();
    };

    GuessWhat.prototype.onLeaveRoomCb = function (e) {
        this.cleanCanvas();
    };

    GuessWhat.prototype.onNewRoundStart = function (e) {
        this.cleanCanvas();
    };

    GuessWhat.prototype.colorClick = function (e) {
        var tar = e.target;
        if (tar.tagName === "BUTTON") {
            this.tool.color = "#" + tar.dataset.color;
        }
    };

    GuessWhat.prototype.thicknessClick = function (e) {
        var tar = e.target;
        if (tar.tagName === "BUTTON") {
            this.tool.thickness = tar.dataset.thickness;
        }
    };

    GuessWhat.prototype.toolsClick = function (e) {
        var tar = e.target;
        if (tar.tagName === "BUTTON") {
            switch (tar.dataset.tool) {
                case "clear":
                    this.askCleanCanvas();
                    break;
                default:
                    break;
            }
        }
    };

    GuessWhat.prototype.askCleanCanvas = function () {
        if (!this.isConnected)
            return false

        var msg = JSON.stringify({
            action: "clean_canvas",
            room: this.roomInput.value
        });
        this.socket.send(msg);
    };

    GuessWhat.prototype.onNewRoundStart = function (e) {
        this.secretWord = "_".repeat(e.word_length);
        this.secretWordP.innerHTML = this.secretWord;
    };

    GuessWhat.prototype.onRevealLetter = function (e) {
        var index = e.pos;
        var letter = e.letter;
        this.secretWord = replaceAt(this.secretWord, index, letter)

        this.secretWordP.innerHTML = this.secretWord;
    };

    GuessWhat.prototype.onYouAreDrawing = function (e) {
        this.secretWord = e.word.Value;
        this.secretWordP.innerHTML = this.secretWord;
    };

    function replaceAt (string, index, replacement) {
        return string.substr(0, index) + replacement + string.substr(index + replacement.length);
    }

    document.addEventListener("DOMContentLoaded", function () {
        var GW = new GuessWhat();
        GW.init();
    });

    global.GuessWhat = GuessWhat;
}(this));
