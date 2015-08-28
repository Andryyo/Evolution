var socket;

socket = new WebSocket("ws://127.0.0.1:8081/socket");
//socket = new WebSocket("ws://93.188.39.118:8081/socket");
//socket = new WebSocket("ws://82.193.120.243:80/socket");
socket.onopen = onSocketOpen;
socket.onmessage = onSocketMessage;

function onSocketOpen(event) {
	var textArea = document.getElementById("log");
    textArea.value = "";
};

function onSocketMessage(event) {
	var textArea = document.getElementById("log");
	textArea.value = textArea.value + '\n' + event.data;
	textArea.scrollTop = textArea.scrollHeight;
	var obj = JSON.parse(event.data);
	if (obj.Type == MESSAGE_EXECUTED_ACTION) {
		showAction(obj.Value);
	}
	if (obj.Type == MESSAGE_CHOICES_LIST) {
		updateGameState(obj.Value.State)
		availableActions = obj.Value.Actions;
	}
	if (obj.Type == MESSAGE_LOBBIES_LIST) {
		updateLobbiesList(obj.Value);
	}
};

function connectToLobby(lobbyId) {
	var playerId = localStorage.getItem("PlayerId")
	message = {
		Type: MESSAGE_JOIN_LOBBY,
		Value: {
			LobbyId: lobbyId,
			PlayerId: playerId
		}}
	socket.send(JSON.stringify(message))
};

function createLobby() {
	message = {
		Type: MESSAGE_NEW_LOBBY,
		Value: null}
	socket.send(JSON.stringify(message))
};


function executeAction(action) {
	if (availableActions == null) {
		return false;
	}
	for (var i in availableActions) {
		var tmp1 = JSON.stringify(action);
		var tmp2 = JSON.stringify(availableActions[i]);
		if (JSON.stringify(availableActions[i]) === JSON.stringify(action)) {
			availableActions = null;
			response = {
				Type: MESSAGE_CHOICE_NUM,
				Value:i
			}
			socket.send(JSON.stringify(response));
			return true;
		}
	}
	return false;
}