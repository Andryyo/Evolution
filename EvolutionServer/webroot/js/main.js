var game = new Phaser.Game(1000, 800, Phaser.AUTO, 'game_holder', { preload: preload, create: create, update: update, render: render});
var gameOverlay;
var cardHeight = 254;
var cardWidth = 182;
var cardEdgeWidth = 38;
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
var selectionRect = null;

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
	game.load.image('end turn', 'assets/End turn.png');
	game.load.image('vote', 'assets/vote.png');
	game.load.image('chain', 'assets/copper-chain-btf-0292-sm.png');
}

function create() {
	game.input.addMoveCallback(mouseMoveCallback,this);
	game.input.onUp.add(mouseUp, this);
	game.add.tileSprite(0, 0, game.width, game.height, 'table');
	mainArea = new Phaser.Rectangle(10, 10, game.width-20, game.height-cardHeight-10);
	handArea = new Phaser.Rectangle(10, game.height-cardHeight+10, game.width-controlAreaWidth-30, cardHeight-20);
	controlArea = new Phaser.Rectangle(game.width-controlAreaWidth-10, game.height-cardHeight+10, controlAreaWidth, cardHeight-20);
	game.add.button(controlArea.x + 10, controlArea.y + 50, 'pass', pass, this);
	game.add.button(controlArea.x + 10, controlArea.y + 110, 'end turn', endTurn, this);
	game.add.button(controlArea.x + 10, controlArea.y + 170, 'vote', vote, this);
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
	socket = new WebSocket("ws://127.0.0.1:8081/socket");
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

function executeAddCreatureAction(cardId) {
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
	var action = {
		Type: "Add single property",
		Arguments: {
			Creature: creatureId,
			Property: propertyId
		}
	};
	return executeAction(action)
}

function executeAddPairPropertyAction(firstCreatureId, secondCreatureId, propertyId) {
	var action = {
		Type: "Add pair property",
		Arguments: {
			Pair: [
				firstCreatureId,
				secondCreatureId
			],
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
	hand.removeAll(true);
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
	if (selectionArrow != null) {
		selectionArrow.arrow.destroy();
		selectionArrow = null;
	}
	players.removeAll(true);
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
	this.Id = creatureDTO.Id;
	for (var i in creatureDTO.Cards) {
		var card = new Card(creatureDTO.Cards[i], 0, cardEdgeWidth/2 * i);
		game.add.existing(card);
		card.inputEnabled = true;
		this.add(card);
		card.selection = null;
		card.events.onInputOver.add(propertyOver, card);
		card.events.onInputOut.add(propertyOut, card);
	}
	var back = new Phaser.Sprite(game, 0, creatureDTO.Cards.length*cardEdgeWidth/2, 'back');
	back.anchor.setTo(0.5, 0.5);
    back.scale.setTo(0.5, 0.5);
	game.add.existing(back);
	this.add(back);
};

Creature.prototype = Object.create(Phaser.Group.prototype);
Creature.prototype.constructor = Creature;

function propertyOver(card, pointer) {
	if (card.selection == null) {
		card.selection = game.add.graphics();
		card.parent.parent.add(card.selection);
		var creature = card.parent;
		card.selection.lineStyle(1, 0x000000, 1);
		card.selection.drawRoundedRect(creature.position.x + 4 - cardWidth/4, creature.position.y + card.position.y - cardHeight/4  + 4 , cardWidth/2-8, cardEdgeWidth/2-8, 3);
		var property = card.getActiveProperty();
		if (property.pair) {
			var pairCard = getPairProperty(card);
			if (pairCard != null) {
				var creature = pairCard.parent;
				card.selection.drawRoundedRect(creature.position.x + 4 - cardWidth/4, creature.position.y + pairCard.position.y - cardHeight/4 + 4, cardWidth/2-8, cardEdgeWidth/2-8, 3);
			}
		}
	}
}

function propertyOut(card, pointer) {
	card.selection.destroy();
	card.selection = null;
}

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
	if (selectionRect != null) {
		selectionRect.destroy();
		selectionRect = null;
	}
	var creature = getIntersectedCreature(card.getBounds());
	if (Phaser.Rectangle.intersects(card.getBounds(),mainArea)) {
		if (creature != null) {
			var property = card.getActiveProperty();
			if (!property.pair) {
				if (executeAddPropertyAction(creature.Id, property.Id)) {
					return;
				} else {
					card.position = card.input.dragStartPoint.clone();
					return;
				}
			} else {
				var arguments = {
					firstCreature: creature,
					property: property
				};
				startSelection(creature, arguments, onSelectSecondPairCreature);
				return;
			}
		}
		if (executeAddCreatureAction(card.Id)) {
			return;
		}
	}
	card.position = card.input.dragStartPoint.clone();
}

function cardDragUpdate(card) {
	var intersectedCreature = getIntersectedCreature(card.getBounds());
	if (intersectedCreature != null) {
		//var bounds = intersectedCreature.getLocalBounds();
		var bounds = new Phaser.Rectangle(0, 0, cardWidth/2, cardHeight/2);
		if (selectionRect != null) {
			selectionRect.destroy();
		}
		selectionRect = game.add.graphics();
		intersectedCreature.parent.add(selectionRect);
		selectionRect.lineStyle(2, 0xFFFFFF, 1);
		selectionRect.drawRoundedRect(-cardWidth/4-10, -cardHeight/4-10, bounds.width+20, bounds.height+20);
		selectionRect.x = intersectedCreature.x;
		selectionRect.y = intersectedCreature.y;
		game.world.bringToTop(selectionRect);
	} else {
		if (selectionRect != null) {
			selectionRect.destroy();
			selectionRect = null;
		}
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

function getCreatureAtPoint(point) {
	if (Phaser.Rectangle.containsPoint(mainArea, point)) {
		for (var i in players.children) {
			for (var j in players.getChildAt(i).children) {
				var creature = players.getChildAt(i).getChildAt(j);
				var bounds = creature.getBounds();
				if (Phaser.Rectangle.containsPoint(bounds, point)) {
					return creature;
				}
			}
		}
	}
	return null;
}

function getPairProperty(firstProperty) {
	var player = firstProperty.parent.parent;
	for (var j in player.children) {
		var creature = player.getChildAt(j);
		for (var k = 0; k<creature.children.length-1; k++) {
			if (firstProperty.parent.Id == creature.getChildAt(k).parent.Id) {
				continue;
			}
			if (creature.getChildAt(k).Id == firstProperty.Id) {
				return creature.getChildAt(k);
			}
		}
	}
	return null;
}

function getCardAtPoint(point) {
	var creature = getCreatureAtPoint(point);
	if (creature == null) {
		return null;
	}
	for (var i = creature.children.length-2; i=>0; i++) {
		if (Phaser.Rectangle.containsPoint(creature.getChildAt(i).getBounds(), point)) {
			return creature.getChildAt(i);
		}
	}
	return null;
}

Card = function(cardDTO, x, y) {
	Phaser.Sprite.call(this, game, x, y, 'cards');
	this.anchor.setTo(0.5, 0.5);
	this.scale.setTo(0.5, 0.5);
	game.physics.arcade.enable(this);
    this.inputEnabled = true;
    this.Id = cardDTO.Id;
    this.properties = cardDTO.Properties;
    this.flipped = false;
	if (cardDTO.ActiveProperty.Id != cardDTO.Properties[0].Id) {
		this.flipped = true;
		this.scale.y *= -1;
		this.scale.x *= -1;
	}
	this.getActiveProperty = function() {
		if (this.properties.length == 1 || !this.flipped) {
				return this.properties[0];
			} else {
				return this.properties[1];
		}
	};
	this.properties[0].pair = false;
    if ($.inArray("Communication", this.properties[0].Traits) != -1) {
		this.properties[0].pair = true;
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
		this.properties[0].pair = true;
      	this.frame = 7;
    } else if ($.inArray("Cooperation", this.properties[0].Traits) != -1 && $.inArray("Fat tissue", this.properties[1].Traits) != -1) {
		this.properties[0].pair = true;
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
		this.properties[0].pair = true;
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
		updateSelectionArrow(x, y);
	}
}

function startSelection(startObject, arguments, onSelect) {
	var arrow = game.add.group();
	arrow.x = startObject.getBounds().x + startObject.getBounds().width/2;
	arrow.y = startObject.getBounds().y + startObject.getBounds().height/2;
	var line = game.add.tileSprite(-6, 0, 12, 1, 'chain');
	arrow.add(line);
	selectionArrow = {
		arrow: arrow,
		arguments: arguments,
		onSelect: onSelect
	};
	updateSelectionArrow(game.input.mousePointer.x, game.input.mousePointer.y);
}

function updateSelectionArrow(x, y) {
	var group = selectionArrow.arrow;
	var length = Math.sqrt((x-group.x)*(x-group.x) + (y-group.y)*(y-group.y));
	group.getChildAt(0).height = length;
	var angle = Math.atan((x-group.x)/(y-group.y));
	if (group.y > y) {
		angle += Math.PI;
	}
	if (y != group.y) {
		group.rotation = - angle;
	}
	var sprite = group.getChildAt(0)
}

function mouseUp(pointer) {
	if (selectionArrow != null) {
		selectionArrow.onSelect(selectionArrow.arguments, pointer);
		selectionArrow.arrow.destroy();
		selectionArrow = null;
	}
}

function onSelectSecondPairCreature(arguments, pointer) {
	var firstCreature = arguments.firstCreature;
	var secondCreature = getCreatureAtPoint(pointer.position);
	var property = arguments.property;
	if (firstCreature == null || secondCreature == null) {
		return
	}
	executeAddPairPropertyAction(firstCreature.Id, secondCreature.Id, property.Id);
}