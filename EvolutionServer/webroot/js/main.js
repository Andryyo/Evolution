var game = new Phaser.Game(1000, 800, Phaser.AUTO, 'game_holder', { preload: preload, create: create, update: update, render: render});
var gameOverlay;
var cardHeight = 254;
var cardWidth = 182;
var cardEdgeWidth = 35;
var controlAreaWidth = 170;
var handArea;
var mainArea;
var controlArea;
var foodBank;
var hand = null;
var players = null;
var availableActions = null;
var currentPlayerId;
var playerId;
var selectionRect;
var socket;
var voteStart = false;
var selectionArrow;

var MESSAGE_EXECUTED_ACTION = 0
var MESSAGE_CHOICES_LIST = 1
var	MESSAGE_NAME = 2
var	MESSAGE_CHOICE_NUM = 3
var MESSAGE_LOBBIES_LIST = 4
var	MESSAGE_NEW_LOBBY = 5
var	MESSAGE_JOIN_LOBBY = 6
var MESSAGE_VOTE_START = 7

function preload() {
	game.load.spritesheet('cards','assets/spritesheet.png',cardWidth,cardHeight,20);
	game.load.image('back','assets/back.png');
	game.load.image('table','assets/bg_texture___wood_by_nortago.jpg');
	game.load.image('bronze','assets/bronze.png');
	game.load.image('copper','assets/copper.png');
	game.load.image('pass','assets/pass.png');
	game.load.image('end turn', 'assets/End turn.png')
	game.load.image('vote', 'assets/vote.png')
}

function create() {
	game.input.addMoveCallback(mouseMoveCallback,this);
	game.add.tileSprite(0, 0, game.width, game.height, 'table');
	mainArea = new Phaser.Rectangle(10, 10, game.width-20, game.height-cardHeight-10);
	handArea = new Phaser.Rectangle(10, game.height-cardHeight+10, game.width-controlAreaWidth-30, cardHeight-20);
	controlArea = new Phaser.Rectangle(game.width-controlAreaWidth-10, game.height-cardHeight+10, controlAreaWidth, cardHeight-20);
	game.add.button(controlArea.x + 10, controlArea.y + 50, 'pass', pass, this);
	game.add.button(controlArea.x + 10, controlArea.y + 110, 'end turn', endTurn, this);
	game.add.button(controlArea.x + 10, controlArea.y + 170, 'vote', vote, this);
	selectionRect = game.add.graphics();
	game.physics.startSystem(Phaser.Physics.ARCADE);
	gameOverlay = game.add.graphics(0, 0);
	gameOverlay.lineStyle(2, 0xFFFFFF, 1);
	gameOverlay.drawRoundedRect(mainArea.x, mainArea.y, mainArea.width, mainArea.height, 3);
	gameOverlay.drawRoundedRect(handArea.x, handArea.y, handArea.width, handArea.height, 3);
	gameOverlay.drawRoundedRect(controlArea.x, controlArea.y, controlArea.width, controlArea.height, 3);
	foodBank = game.add.graphics();
	foodBank.x = mainArea.halfWidth;
	foodBank.y = mainArea.halfHeight;
	foodBank.lineStyle(0);
	hand = game.add.group();
	hand.x = handArea.x;
	hand.y = handArea.y;
	players = game.add.group();
	socket = new WebSocket("ws://127.0.0.1/socket");
	//socket = new WebSocket("ws://93.188.39.118:8081/socket");
	//socket = new WebSocket("ws://82.193.120.243:80/socket");
	socket.onopen = onSocketOpen;
	socket.onmessage = onSocketMessage;
}

function vote() {
	voteStart = !voteStart;
	var message = {
    		Type: MESSAGE_VOTE_START,
    		Value: voteStart
    	};
    socket.send(JSON.stringify(message));
}

function pass() {
	if (availableActions == null) {
		return false;
	};
	var action = {
		Type: "Pass",
		Arguments: {}
	};
	return executeAction(action)
}

function endTurn() {
	if (availableActions == null) {
		return false;
	};
	var action = {
		Type: "End turn",
		Arguments: {}
	};
	return executeAction(action)
}

function update() {
	if (hand!= null) {
		hand.forEach(function(card) {
			if (card.input.overDuration() > 500 && !card.input.isDragged) {
				if (!card.flipped) {
                	card.scale.y = 1;
                	card.scale.x = 1;
                } else {
                	card.scale.x = -1;
                	card.scale.y = -1;
                }

			}
		}, this);
	}
}

function render() {
}

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

function updateLobbiesList(lobbies) {
	var select = document.getElementById("lobbies");
	while (select.hasChildNodes()) {
		select.removeChild(select.lastChild);
	}
	for (var i in lobbies) {
		var option = document.createElement("button");
		option.type = "button"
		option.onclick=function (event) {
			connectToLobby(event.target.	lobbyId);
			$("#overlay").hide();
		}
		option.className="list-group-item";
		option.innerHTML = "Lobby " + lobbies[i].Id + ": " + lobbies[i].PlayersCount + " players";
		option.lobbyId = lobbies[i].Id;
		select.appendChild(option);
	}
};

function executeAction(action) {
	for (var i in availableActions) {
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

function executeAddCreatureAction(cardId) {
	if (availableActions == null) {
		return false;
	};
	var action = {
		Type: "Add creature",
		Arguments: {
			Card: cardId,
			Player: currentPlayerId
		}
	};
	return executeAction(action)
}

function executeAddPropertyAction(creatureId, propertyId) {
	if (availableActions == null) {
		return false;
	};
	var action = {
		Type: "Add single property",
		Arguments: {
			Creature: creatureId,
			Property: propertyId
		}
	};
	return executeAction(action)
}

function showAction(action) {
	updateGameState(action.State)
}

function updateGameState(state) {
	currentPlayerId=state.CurrentPlayerId;
	playerId = state.PlayerId;
	localStorage.setItem("PlayerId", playerId);
	updateFoodBank(state.FoodBank);
	updatePlayers(state.Players);
	updateHand(state.PlayerCards);
}

function updateFoodBank(count) {
	foodBank.clear();
	foodBank.beginFill(0xFF0000, 1);
	var rectangle = new Phaser.Rectangle(-50, -50, 100, 100);
	for (var i = 0; i<count; i++) {
		foodBank.drawCircle(rectangle.randomX, rectangle.randomY, 10);
	}
	foodBank.endFill();
}

function updateHand(handDTO) {
	hand.removeAll();
	var y = handArea.halfHeight;
	var startX = (handArea.width-(cardWidth*handDTO.length/2*3/2))/2;
	if (startX < 0) {
		startX = cardWidth/4;
	}
	var offset = (handArea.width-startX*2)/(handDTO.length);
	
	for (var i in handDTO) {
		var card = new Card(handDTO[i], startX + (+i + +0.5)*offset, y);
		card.events.onInputOver.add(cardOver, card);
    	card.events.onInputOut.add(cardOut, card);
	    card.events.onInputUp.add(cardUp, card);
	    card.events.onDragStart.add(cardDragStart, card);
	    card.events.onDragStop.add(cardDragStop, card);
	    card.events.onDragUpdate.add(cardDragUpdate, card);
	    card.input.enableDrag();
		hand.add(card);
	}
}

function updatePlayers(playersDTO) {
	players.removeAll();
	var startAngle = 180;
	var deltaAngle = 360/playersDTO.length;
	var radiusX = mainArea.halfWidth - cardHeight/4;
	var radiusY = mainArea.halfHeight - cardHeight/4;
	var playerIndex = 0;
	for (var i in playersDTO) {
		if (playersDTO[i].Id == playerId) {
			playerIndex = i;
			break;
		}
	}
	var angle = 0;
	for (var i = playerIndex; i<playersDTO.length; i++) {
		var playersCreatures = new PlayerCreatures(playersDTO[i], mainArea.halfWidth-Math.sin(angle*Math.PI/180)*radiusX, mainArea.halfHeight+Math.cos(angle*Math.PI/180)*radiusY, angle)
        game.add.existing(playersCreatures);
        players.add(playersCreatures);
        angle += deltaAngle;
	}
	for (var i = 0; i<playerIndex; i++) {
		var playersCreatures = new PlayerCreatures(playersDTO[i], mainArea.halfWidth-Math.sin(angle*Math.PI/180)*radiusX, mainArea.halfHeight+Math.cos(angle*Math.PI/180)*radiusY, angle)
        game.add.existing(playersCreatures);
        players.add(playersCreatures);
        angle += deltaAngle;
	}
}

PlayerCreatures = function(playerDTO, x, y, angle) {
	Phaser.Group.call(this, game);
	this.x = x;
	this.y = y;
	this.angle = angle;
	var totalCreatureWidthHalf = cardWidth/2 * playerDTO.Creatures.length/2;
	for (var i in playerDTO.Creatures) {
		var creature = new Creature(playerDTO.Creatures[i], (+i + +1)*cardWidth/2-totalCreatureWidthHalf, 0);
		game.add.existing(creature);
		this.add(creature);
	}
};

PlayerCreatures.prototype = Object.create(Phaser.Group.prototype);
PlayerCreatures.prototype.constructor = PlayerCreatures;

Creature = function(creatureDTO, x, y) {
	Phaser.Group.call(this, game);
	this.x = x-cardWidth/4;
	this.y = y-cardHeight/4;
	this.id = creatureDTO.Id;
	for (var i in creatureDTO.Cards) {
		var card = new Card(creatureDTO.Cards[i], 0, cardEdgeWidth/2 * i);
		game.add.existing(card);
		this.add(card);
		card.events.onInputDown.add(function (card) {
			var arrow = game.add.group();
			arrow.x = card.getBounds().x + card.getBounds().width/2;
			arrow.y = card.getBounds().y + card.getBounds().height/2;
			var line = game.add.tileSprite(-2, 0, 4, 200, 'copper');
			arrow.add(line);
			selectionArrow = {
				arrow: arrow,
				startObject: card
			};
		}, card);
	}
	var back = new Phaser.Sprite(game, 0, creatureDTO.Cards.length*cardEdgeWidth/2, 'back');
	back.anchor.setTo(0.5, 0.5);
    back.scale.setTo(0.5, 0.5);
	game.add.existing(back);
	this.add(back);
};

Creature.prototype = Object.create(Phaser.Group.prototype);
Creature.prototype.constructor = Creature;

function cardOver(card, pointer) {
	card.bringToTop();
}

function cardUp(card, pointer) {
	if (card.input.pointerTimeUp()-card.input.pointerTimeDown() < 70) {
		card.flipped = !card.flipped;
		card.scale.y *= -1;
		card.scale.x *= -1;
	}
}

function cardOut(card, pointer) {
	card.anchor.y = 0.5;
	if (!card.flipped) {
    	card.scale.setTo(0.5, 0.5);
    } else {
    	card.scale.setTo(-0.5, -0.5);
    }
}

function cardDragStart(card) {
	card.anchor.y = 0.5;
	if (!card.flipped) {
    	card.scale.setTo(0.5, 0.5);
    } else {
    	card.scale.setTo(-0.5, -0.5);
    }
}

function cardDragStop(card) {
	selectionRect.clear();
	var creature = getIntersectedCreature(card.getBounds());
	if (Phaser.Rectangle.intersects(card.getBounds(),mainArea)) {
		if (creature != null) {
			if (card.properties.length == 1 || !card.flipped) {
				var property = card.properties[0];
			} else {
				var property = card.properties[1];
			}
			if (executeAddPropertyAction(creature.id, property.Id)) {
				return;
			} else {
				card.position = card.input.dragStartPoint.clone();
				return;
			}
		}
		if (executeAddCreatureAction(card.id)) {
			return;
		}
	}
	card.position = card.input.dragStartPoint.clone();
}

function cardDragUpdate(card) {
	var intersectedCreature = getIntersectedCreature(card.getBounds());
	selectionRect.clear();
	if (intersectedCreature != null) {
		intersectedCreature.parent.add(selectionRect);
		selectionRect.lineStyle(2, 0xFFFFFF, 1);
 		selectionRect.moveTo(-cardWidth/4-10, -cardHeight/4-10);
		selectionRect.lineTo(-cardWidth/4-10, +cardHeight/4+10);
		selectionRect.moveTo(cardWidth/4+10, -cardHeight/4-10);
		selectionRect.lineTo(cardWidth/4+10, +cardHeight/4+10);
		selectionRect.x = intersectedCreature.x;
		selectionRect.y = intersectedCreature.y;
		game.world.bringToTop(selectionRect);
	}
}

function getIntersectedCreature(rectangle) {
	var maxIntersectObject;
	var maxIntersectArea = 0;
	if (Phaser.Rectangle.intersects(rectangle,mainArea)) {
		for (var i in players.children) {
			for (var j in players.getChildAt(i).children) {
				var creature = players.getChildAt(i).getChildAt(j);
				var bounds = creature.getBounds();
				var intersectionRect = Phaser.Rectangle.intersection(rectangle, bounds);
				if (! intersectionRect.empty) {
					var intersectionRectArea = intersectionRect.width * intersectionRect.height;
					if (intersectionRectArea > maxIntersectArea) {
						maxIntersectArea = intersectionRectArea;
						maxIntersectObject = creature;
					}
				}
			}
		}
	}
	return maxIntersectObject;
}

Card = function(cardDTO, x, y) {
	Phaser.Sprite.call(this, game, x, y, 'cards');
	this.anchor.setTo(0.5, 0.5);
	this.scale.setTo(0.5, 0.5);
	game.physics.arcade.enable(this);
    this.inputEnabled = true;
    this.id = cardDTO.Id;
    this.properties = cardDTO.Properties;
    this.flipped = false;
	if (cardDTO.ActiveProperty.Id != cardDTO.Properties[0].Id) {
		this.flipped = true;
		this.scale.y *= -1;
		this.scale.x *= -1;
	}
    if ($.inArray("Communication", this.properties[0].Traits) != -1) {
		this.frame = 0;
	} else if ($.inArray("High body weight", this.properties[0].Traits) != -1 && $.inArray("Fat tissue", this.properties[1].Traits) != -1) {
	 	this.frame = 1;
	} else if ($.inArray("High body weight", this.properties[0].Traits) != -1 && $.inArray("Carnivorous", this.properties[1].Traits) != -1) {
        this.frame = 2;
    } else if ($.inArray("Sharp vision", this.properties[0].Traits) != -1) {
        this.frame = 3;
    } else if ($.inArray("Grazing", this.properties[0].Traits) != -1) {
        this.frame = 4;
    } else if ($.inArray("Parasite", this.properties[0].Traits) != -1 && $.inArray("Carnivorous", this.properties[1].Traits) != -1) {
      	this.frame = 5;
    } else if ($.inArray("Burrowing", this.properties[0].Traits) != -1) {
        this.frame = 6;
    } else if ($.inArray("Cooperation", this.properties[0].Traits) != -1 && $.inArray("Carnivorous", this.properties[1].Traits) != -1) {
      	this.frame = 7;
    } else if ($.inArray("Cooperation", this.properties[0].Traits) != -1 && $.inArray("Fat tissue", this.properties[1].Traits) != -1) {
    	this.frame = 8;
    } else if ($.inArray("Poisonous", this.properties[0].Traits) != -1) {
    	this.frame = 9;
    } else if ($.inArray("Camouflage", this.properties[0].Traits) != -1) {
        this.frame = 10;
    } else if ($.inArray("Hibernation", this.properties[0].Traits) != -1) {
        this.frame = 11;
    } else if ($.inArray("Mimicry", this.properties[0].Traits) != -1) {
        this.frame = 12;
    } else if ($.inArray("Symbiosys", this.properties[0].Traits) != -1) {
        this.frame = 13;
    } else if ($.inArray("Scavenger", this.properties[0].Traits) != -1) {
       this.frame = 14;
    } else if ($.inArray("Piracy", this.properties[0].Traits) != -1) {
       this.frame = 15;
    } else if ($.inArray("Tail loss", this.properties[0].Traits) != -1) {
       this.frame = 16;
    } else if ($.inArray("Running", this.properties[0].Traits) != -1) {
       this.frame = 17;
    } else if ($.inArray("Swimming", this.properties[0].Traits) != -1) {
       this.frame = 18;
    } else if ($.inArray("Parasite", this.properties[0].Traits) != -1 && $.inArray("Fat tissue", this.properties[1].Traits) != -1) {
       this.frame = 19;
    } else {
    	alert("Unknown card");
    }
    game.add.existing(this)
};

Card.prototype = Object.create(Phaser.Sprite.prototype);
Card.prototype.constructor = Card;

function myOnKeyPress(e) {
	if (e.keyCode == 13) {
		var command = document.getElementById("command").value
		socket.send(command)
		var textArea = document.getElementById("log")
		textArea.value = textArea.value + '\n' + command
		document.getElementById("command").value = ""
		return false
	}
};

function mouseMoveCallback(pointer, x, y) {
	if (selectionArrow != null) {
		var group = selectionArrow.arrow;
		var angle = Math.atan((x-group.x)/(y-group.y));
		if (group.y > y) {
			angle += Math.PI;
		}
		if (y != group.y) {
			group.rotation = - angle;
		}
		var sprite = group.getChildAt(0)
	}
}