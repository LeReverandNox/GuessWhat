/*jslint browser this */
/*global $ window URL WebSocket console */

(function (global) {
    "use strict";
    var GuessWhat = function () {

    };

    GuessWhat.prototype.startEventListeners = function () {
        this.changeNicknameButton.addEventListener("click", this.changeNickname.bind(this));
        this.sendMessageButton.addEventListener("click", this.sendMessage.bind(this));
        this.joinRoomButton.addEventListener("click", this.joinRoom.bind(this));
        this.leaveRoomButton.addEventListener("click", this.leaveRoom.bind(this));
    };

    GuessWhat.prototype.registerElements = function () {
        this.changeNicknameButton = document.getElementById("change_nickname");
        this.nicknameInput = document.getElementById("nickname");

        this.sendMessageButton = document.getElementById("send_message");
        this.messageInput = document.getElementById("message");

        this.joinRoomButton = document.getElementById("join_room");
        this.leaveRoomButton = document.getElementById("leave_room");
        this.roomInput = document.getElementById("room");
    };

    GuessWhat.prototype.changeNickname = function () {
        var msg = JSON.stringify({
            action: "set_nickname",
            nickname: this.nicknameInput.value
        });
        this.socket.send(msg);
    };

    GuessWhat.prototype.sendMessage = function () {
        var msg = JSON.stringify({
            action: "send_message",
            content: this.messageInput.value
        });
        this.socket.send(msg);
    };

    GuessWhat.prototype.joinRoom = function () {
        var msg = JSON.stringify({
            action: "join_room",
            room: this.roomInput.value
        });
        this.socket.send(msg);
    };

    GuessWhat.prototype.leaveRoom = function () {
        var msg = JSON.stringify({
            action: "leave_room",
            room: this.roomInput.value
        });
        this.socket.send(msg);
    };

    GuessWhat.prototype.onMessage = function (e) {
        console.log(JSON.parse(e.data));
    };

    GuessWhat.prototype.onClose = function (e) {
        console.log("Socket closed");
    };

    GuessWhat.prototype.startSocket = function (cb) {
        if (window["WebSocket"]) {
            this.socket = new WebSocket("ws://" + document.location.host + "/ws");
            this.socket.onclose = this.onClose;
            this.socket.onmessage = this.onMessage;

            return cb(false);
        } else {
            console.log("WEBSOCKET NOT SUPPORTED");
            return cb(true);
        }
    };

    GuessWhat.prototype.init = function () {
        var self = this;
        this.startSocket(function (err) {
            if (err) {
                return alert("Ca va pas marcher...");
            }
            self.registerElements();
            self.startEventListeners();
        });
    };

    document.addEventListener("DOMContentLoaded", function () {
        var GW = new GuessWhat();
        GW.init();
    });

    global.GuessWhat = GuessWhat;
}(this));
