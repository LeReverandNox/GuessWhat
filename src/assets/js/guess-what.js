/*jslint browser this */
/*global $ window URL WebSocket console */

(function (global) {
    "use strict";
    var GuessWhat = function () {
        this.isConnected = false;
        this.tool = {
            color: "#000000",
            thickness: 10
        };
    };

    GuessWhat.prototype.startEventListeners = function () {
        this.connexionButton.addEventListener("click", this.connect.bind(this));
        this.sendMessageButton.addEventListener("click", this.sendMessage.bind(this));
        this.joinRoomButton.addEventListener("click", this.joinRoom.bind(this));
        this.leaveRoomButton.addEventListener("click", this.leaveRoom.bind(this));

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
        this.roomInput = document.getElementById("room");

        this.canvas = document.getElementById("canvas");
        this.context = this.canvas.getContext("2d");
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
        var msg = JSON.stringify({
            action: "canvas_mouse_down",
            room: this.roomInput.value,
            x: String(e.layerX),
            y: String(e.layerY),
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

        var msg = JSON.stringify({
            action: "canvas_mouse_move",
            room: this.roomInput.value,
            x: String(e.layerX),
            y: String(e.layerY),
            thickness: String(this.tool.thickness),
            color: this.tool.color
        });
        this.socket.send(msg);
    };
    GuessWhat.prototype.onLocalCanvasMouseUp = function (e) {
        if (!this.isConnected)
            return false
        if (!this.localClick1)
            return false

        this.localClick1 = false;
        var msg = JSON.stringify({
            action: "canvas_mouse_up",
            room: this.roomInput.value,
            x: String(e.layerX),
            y: String(e.layerY),
            thickness: String(this.tool.thickness),
            color: this.tool.color
        });
        this.socket.send(msg);
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
        this.click1 = true;
        this.context.beginPath();
        this.context.moveTo(e.x, e.y);
    };

    GuessWhat.prototype.onCanvasMouseMove = function (e) {
        if (this.click1 === true) {

            this.context.lineCap = "round";
            this.context.lineJoin = "round";
            this.context.strokeStyle = e.color;
            this.context.fillStyle = e.color;
            this.context.lineWidth = e.thickness;

            this.context.lineTo(e.x, e.y);
            this.context.stroke();
        }
    };

    GuessWhat.prototype.onCanvasMouseUp = function (e) {
        this.click1 = false;
    };

    GuessWhat.prototype.cleanCanvas = function () {
        this.context.clearRect(0, 0, this.canvas.width, this.canvas.height);
    };

    GuessWhat.prototype.init = function () {
        var self = this;
        this.registerElements();
        this.startEventListeners();
    };

    document.addEventListener("DOMContentLoaded", function () {
        var GW = new GuessWhat();
        GW.init();
    });

    global.GuessWhat = GuessWhat;
}(this));
