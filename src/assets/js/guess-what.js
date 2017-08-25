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

        this.rooms = [];
        this.generalChatMessages = [];
        this.generalClients = [];
        this.roomChatMessages = [];
        this.roomClients = [];

        this.isInRoom = false;
    };

    GuessWhat.prototype.startEventListeners = function () {
        this.connexionButton.addEventListener("click", this.connect.bind(this));
        this.sendMessageButton.addEventListener("click", this.sendMessage.bind(this));
        this.createRoomButton.addEventListener("click", this.createRoom.bind(this));
        this.leaveRoomButton.addEventListener("click", this.leaveRoom.bind(this));
        this.startRoomButton.addEventListener("click", this.startRoom.bind(this));
        this.colorHolder.addEventListener("click", this.colorClick.bind(this));
        this.thicknessHolder.addEventListener("click", this.thicknessClick.bind(this));
        this.toolsHolder.addEventListener("click", this.toolsClick.bind(this));

        this.$roomsHolder.on("click", this.roomClick.bind(this));

        this.canvas.addEventListener("mousedown", this.onLocalCanvasMouseDown.bind(this));
        this.canvas.addEventListener("mousemove", this.onLocalCanvasMouseMove.bind(this));
        this.canvas.addEventListener("mouseup", this.onLocalCanvasMouseUp.bind(this));
        this.canvas.addEventListener("mouseout", this.onLocalCanvasMouseOut.bind(this));
    };

    GuessWhat.prototype.registerElements = function () {
        this.$connexionBlock = $("#connexion_block");
        this.$gameBlock = $("#game_block");

        this.$chatHolder = $("#chat_holder");
        this.$roomsHolder = $("#rooms_holder");
        this.$clientsHolder = $("#clients_holder");
        this.$drawerToolsHolder = $("#drawer_tools_holder");
        this.$roundInfoHolder = $("#round_info_holder");

        this.connexionButton = document.getElementById("connexion_button");
        this.nicknameInput = document.getElementById("nickname");

        this.sendMessageButton = document.getElementById("send_message");
        this.messageInput = document.getElementById("message");

        this.createRoomButton = document.getElementById("create_room");
        this.leaveRoomButton = document.getElementById("leave_room");
        this.startRoomButton = document.getElementById("start_room");
        this.roomInput = document.getElementById("room");

        this.colorHolder = document.getElementById("color_holder");
        this.thicknessHolder = document.getElementById("thickness_holder");
        this.toolsHolder = document.getElementById("tools_holder");

        this.secretWordP = document.getElementById("secret_word");
        this.countdownHolder = document.getElementById("countdown_holder");
        this.countdown = document.getElementById("countdown");

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

    GuessWhat.prototype.joinRoom = function (roomName) {
        if (!this.isConnected)
            return false

        var msg = JSON.stringify({
            action: "join_room",
            room: roomName
        });
        this.socket.send(msg);
    };


    GuessWhat.prototype.createRoom = function () {
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
                case "new_round_start":
                this.onNewRoundStart(data);
                break;
                case "round_end":
                this.onRoundEnd(data)
                break;
                case "clean_canvas":
                this.onCleanCanvas(data);
                break;
                case "reveal_letter":
                this.onRevealLetter(data);
                break;
                case "you_are_drawing":
                this.onYouAreDrawing(data);
                break;
                case "round_is_going":
                this.onRoundIsGoing(data);
                break;
                case "connexion_cb":
                this.onConnexionCb(data);
                break;
                // Chat events
                case "incoming_all_global_message":
                this.onIncomingAllGlobalMessage(data);
                break;
                case "incoming_global_message":
                this.onIncomingGlobalMessage(data);
                break;
                case "incoming_all_room_message":
                    this.onIncomingAllRoomMessage(data);
                    break;
                case "incoming_room_message":
                    this.onIncomingRoomMessage(data);
                    break;
                // Room events
                case "incoming_all_rooms":
                this.onIncomingAllRooms(data);
                break;
                case "incoming_room":
                this.onIncomingRoom(data);
                break;
                case "leaving_room":
                this.onLeavingRoom(data);
                break;
                case "leaving_room":
                this.onLeavingRoom(data);
                break;
                case "join_room_cb":
                    this.onJoinRoomCb(data);
                    break;
                case "leave_room_cb":
                    this.onLeaveRoomCb(data);
                    break;
                // Client events
            case "incoming_all_global_users":
                this.onIncomingAllGlobalUsers(data);
                break;
            case "incoming_client":
                this.onIncomingGlobalClient(data);
                break;
            case "leaving_client":
                this.onLeavingGlobalClient(data);
                break;
            case "incoming_all_room_clients":
                this.onIncomingAllRoomClients(data);
                break
            case "incoming_room_client":
                this.onIncomingRoomClient(data);
                break;
            case "leaving_room_client":
                this.onLeavingRoomClient(data);
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

        this.$connexionBlock.show();
        this.$gameBlock.hide();
        this.$drawerToolsHolder.hide();
        this.$roundInfoHolder.hide();
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
        this.$roomsHolder.hide();
        this.isInRoom = true;
        this.displayClients();
        this.displayChat();

        this.cleanCanvas();
    };

    GuessWhat.prototype.onLeaveRoomCb = function (e) {
        this.$roomsHolder.show();
        this.isInRoom = false;
        this.displayClients();
        this.displayChat();

        this.cleanCanvas();
        this.$roundInfoHolder.hide();
        this.stopTimer();
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

    GuessWhat.prototype.roomClick = function (e) {
        var tar = e.target;
        if (tar.tagName === "LI") {
            this.joinRoom(tar.dataset.room);
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
        this.cleanCanvas();
        this.$drawerToolsHolder.hide();

        this.initSecretWord(e.word_length);
        this.startTimer(e.room.RoundDuration);

        this.$roundInfoHolder.show();
    };

    GuessWhat.prototype.onRevealLetter = function (e) {
        var index = e.pos;
        var letter = e.letter;
        this.updateSecretWord(index, letter);
    };

    GuessWhat.prototype.onYouAreDrawing = function (e) {
        this.cleanCanvas();
        this.$drawerToolsHolder.show();

        this.secretWord = e.word.Value;
        this.secretWordP.innerHTML = this.secretWord;
        this.startTimer(e.room.RoundDuration);
        this.$roundInfoHolder.show();
    };

    GuessWhat.prototype.startTimer = function (duration) {
        this.countdownIntv = makeTimer(duration, this.countdown);
    }

    GuessWhat.prototype.stopTimer = function () {
        clearInterval(this.countdownIntv);
        this.countdown.innerHTML = "";
    };

    GuessWhat.prototype.onRoundEnd = function (e) {
        this.stopTimer();
        this.$roundInfoHolder.hide()
        console.log(e);
    }

    GuessWhat.prototype.onRoundIsGoing = function (e) {
        var self = this;
        this.startTimer(e.time_left);

        setTimeout(function () {
            self.initSecretWord(e.word_length);

            Object.keys(e.revealed_letters).map(function(index) {
                self.updateSecretWord(parseInt(index), e.revealed_letters[index]);
            });
            self.$roundInfoHolder.show();
        }, 1000)
    };

    GuessWhat.prototype.initSecretWord = function (length) {
        this.secretWord = "_".repeat(length);
        this.secretWordP.innerHTML = this.secretWord;
    };

    GuessWhat.prototype.updateSecretWord = function (index, letter) {
        this.secretWord = replaceAt(this.secretWord, index, letter)
        this.secretWordP.innerHTML = this.secretWord;
    };

    GuessWhat.prototype.onConnexionCb = function (data) {
        if (data.success) {
            this.$connexionBlock.hide();
            this.$gameBlock.show();
        } else {
            alert(data.reason);
        }
    };

    GuessWhat.prototype.displayRooms = function () {
        var rooms = this.rooms;
        var $ul = $("<ul id=\"rooms\"></ul>");
        var $li;
        rooms.map(function (room) {
            $li = $("<li class=\"room\" data-room=\"" + room.Name + "\">" + room.Name + "</li>");
            $ul.append($li);
        });

        this.$roomsHolder.children("ul").remove();
        this.$roomsHolder.append($ul);
    };

    GuessWhat.prototype.displayClients = function () {
        var clients = (this.isInRoom) ? this.roomClients : this.generalClients;
        var $ul = $("<ul id=\"clients\"></ul>");
        var $li;
        clients.map(function (client) {
            $li = $("<li class=\"client\">" + client.Nickname + "</li>");
            $ul.append($li);
        });

        this.$clientsHolder.children("ul").remove();
        this.$clientsHolder.append($ul);

    };

    GuessWhat.prototype.displayChat = function () {
        var messages = (this.isInRoom) ? this.roomChatMessages : this.generalChatMessages;
        var $ul = $("<ul id=\"chat_messages\"></ul>");
        var $li;
        messages.map(function (message) {
            $li = $("<li class=\"chat_message\">" + message.Sender.Nickname +" : " + message.Content + "</li>");
            $ul.append($li);
        });

        this.$chatHolder.children("ul").remove();
        this.$chatHolder.append($ul);
    };

    GuessWhat.prototype.displayInfos = function () {

    };

    GuessWhat.prototype.displayEndRound = function () {

    };

    GuessWhat.prototype.onIncomingAllGlobalMessage = function (e) {
        var messages = e.messages;
        this.generalChatMessages = messages;
        this.displayChat();
    };

    GuessWhat.prototype.onIncomingGlobalMessage = function (e) {
        this.generalChatMessages.push(e.message);
        this.displayChat();
    };

    GuessWhat.prototype.onIncomingAllRooms = function (e) {
        var rooms = e.rooms;
        this.rooms = rooms;
        this.displayRooms();
    }

    GuessWhat.prototype.onIncomingRoom = function (e) {
        this.rooms.push(e.room);
        this.displayRooms();
    }

    GuessWhat.prototype.onIncomingAllGlobalUsers = function (e) {
        var clients = e.clients;
        this.generalClients = clients;
        this.displayClients();
    };

    GuessWhat.prototype.onIncomingGlobalClient = function (e) {
        this.generalClients.push(e.client);
        this.displayClients();
    };

    GuessWhat.prototype.onLeavingGlobalClient = function (e) {
        var clientToRemove = e.client;
        this.generalClients = this.generalClients.filter(function (client, i) {
            if (client.Nickname !== clientToRemove.Nickname) {
                return true;
            }
        })
        this.displayClients();
    };

    GuessWhat.prototype.onLeavingRoom = function (e) {
        var roomToRemove = e.room;
        this.rooms = this.rooms.filter(function (room, i) {
            if (room.Name !== roomToRemove.Name) {
                return true;
            }
        })
        this.displayRooms();
    };

    GuessWhat.prototype.onIncomingAllRoomMessage = function (e) {
        var messages = e.messages;
        this.roomChatMessages = messages;
        this.displayChat();
    };

    GuessWhat.prototype.onIncomingRoomMessage = function (e) {
        this.roomChatMessages.push(e.message);
        this.displayChat();
    };

    GuessWhat.prototype.onIncomingAllRoomClients = function (e) {
        var clients = e.clients;
        this.roomClients = clients;
        this.displayClients();
    };


    GuessWhat.prototype.onIncomingRoomClient = function (e) {
        this.roomClients.push(e.client);
        this.displayClients();
    };

    GuessWhat.prototype.onLeavingRoomClient = function (e) {
        var clientToRemove = e.client;
        this.roomClients = this.roomClients.filter(function (client, i) {
            if (client.Nickname !== clientToRemove.Nickname) {
                return true;
            }
        })
        this.displayClients();
    };

    function replaceAt (string, index, replacement) {
        return string.substr(0, index) + replacement + string.substr(index + replacement.length);
    }
    function makeTimer(duration, display) {
        var intv = setInterval(function () {
            display.innerHTML = duration;
            if (duration  >= 1) {
                duration--
            }
        }, 1000);
        return intv
    }

    document.addEventListener("DOMContentLoaded", function () {
        var GW = new GuessWhat();
        GW.init();
    });

    global.GuessWhat = GuessWhat;
}(this));
