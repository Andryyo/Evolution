var game = new Phaser.Game(1000, 800, Phaser.AUTO, 'game_holder', { preload: preload, create: create, update: update, render: render});
var gameOverlay;
var cardHeight = 254;
var cardWidth = 182;
var cardEdgeWidth = 40;
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
var voteStart = false;
var selectionArrow;
var selectionRect = null;
var currentGameState = null;
var messagesQueue = [];
var messageInProcessing = null;

MESSAGE_EXECUTED_ACTION = 0
MESSAGE_CHOICES_LIST = 1
MESSAGE_NAME = 2
MESSAGE_CHOICE_NUM = 3
MESSAGE_LOBBIES_LIST = 4
MESSAGE_NEW_LOBBY = 5
MESSAGE_JOIN_LOBBY = 6
MESSAGE_VOTE_START = 7

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
	foodBank = game.add.group();
	foodBank.x = mainArea.halfWidth;
	foodBank.y = mainArea.halfHeight;
	hand = game.add.group();
	hand.x = handArea.x;
	hand.y = handArea.y;
	players = game.add.group();
}

function addMessage(message) {
	messagesQueue.push(message);
}

function processMessage(message) {
	if (message.Type == MESSAGE_EXECUTED_ACTION) {
		//showAction(message.Value);
		updateGameState(message.Value.State)
		return
	}
	if (message.Type == MESSAGE_CHOICES_LIST) {
		updateGameState(message.Value.State)
		availableActions = message.Value.Actions;
		messageInProcessing = null;
		return
	}
	if (message.Type == MESSAGE_LOBBIES_LIST) {
		initLobbiesList(message.Value);
		messageInProcessing = null;
		return
	}
	updateGameState(message.Value.State)
}

function showAction(msg) {
	switch(msg.Action.Type) {
		case "Add creature":
			var player = findPlayer(msg.Action.Value.Player)
			break;	
		default:
			messageInProcessing = null;
	}
}

function updateGameState(state) {
	//if (currentGameState == null) {
		currentGameState = state;
		initGameState(state);
		return
	//}
}

function update() {
	if (messageInProcessing == null && messagesQueue.length != 0) {
		messageInProcessing = messagesQueue.shift();
		processMessage(messageInProcessing);
	}
}

function render() {
}

function initLobbiesList(lobbies) {
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

function initGameState(state) {
	currentPlayerId=state.CurrentPlayerId;
	playerId = state.PlayerId;
	localStorage.setItem("PlayerId", playerId);
	updateFoodBank(state.FoodBank);
	initPlayers(state.Players);
	initHand(state.PlayerCards);
}

function updateFoodBank(count) {
	while (foodBank.children.length > count) {
		foodBank.removeChildAt(0);
	}
	//foodBank.clear();
	while (foodBank.children.length < count) {
		var item = game.add.graphics();
		item.lineStyle(0);
		foodBank.add(item)
		var rectangle = new Phaser.Rectangle(-50, -50, 100, 100);
		var x = rectangle.randomX;
		var y = rectangle.randomY;
		item.beginFill(0xFFFFFF, 1);
		item.drawCircle(x, y, 21)
		item.endFill();
		item.beginFill(0xFF0000, 1);
		item.drawCircle(x, y, 20);
		item.endFill();
	}
}

function initHand(handDTO) {
	hand.removeAll(true);
	var y = handArea.halfHeight;
	var count = 0;
	for (var i in handDTO) {
		count ++;
	}
	var startX = (handArea.width-(cardWidth*count/2*3/2))/2;
	if (startX < 0) {
		startX = cardWidth/4;
	}
	var offset = (handArea.width-startX*2)/(count);
	var num = 0.0;
	
	for (var i in handDTO) {
		var card = new Card(handDTO[i], startX + +num + 0.5)*offset, y);
		card.events.onInputOver.add(cardOver, card);
    	card.events.onInputOut.add(cardOut, card);
	    card.events.onInputUp.add(cardUp, card);
	    card.events.onDragStart.add(cardDragStart, card);
	    card.events.onDragStop.add(cardDragStop, card);
	    card.events.onDragUpdate.add(cardDragUpdate, card);
	    card.input.enableDrag();
		hand.add(card);
		num++;
	}
}

function initPlayers(playersDTO) {
	if (selectionArrow != null) {
		selectionArrow.arrow.destroy();
		selectionArrow = null;
	}
	players.removeAll(true);
	var startAngle = 180;
	var count = 0;
	for (var i in playersDTO) {
		count++;
	}
	var deltaAngle = 360/count;
	var radiusX = mainArea.halfWidth - cardHeight/4;
	var radiusY = mainArea.halfHeight - cardHeight/4;
	var angle = 0;
	var playersCreatures = new PlayerCreatures(playersDTO[playerId], mainArea.halfWidth-Math.sin(angle*Math.PI/180)*radiusX, mainArea.halfHeight+Math.cos(angle*Math.PI/180)*radiusY, angle)
    game.add.existing(playersCreatures);
   	players.add(playersCreatures);
    angle += deltaAngle;
	for (var i in playersDTO) {
		if (i == playerId)
			continue;
		var playersCreatures = new PlayerCreatures(playersDTO[i], mainArea.halfWidth-Math.sin(angle*Math.PI/180)*radiusX, mainArea.halfHeight+Math.cos(angle*Math.PI/180)*radiusY, angle)
        game.add.existing(playersCreatures);
        players.add(playersCreatures);
        angle += deltaAngle;
	}
}

function findPlayer(id) {
	for (var i in players.children) {
		if (players.getChildAt(i).Id == id) {
			return players.getChildAt(i);
		}
	}
	return null;
}

PlayerCreatures = function(playerDTO, x, y, angle) {
	Phaser.Group.call(this, game);
	this.Id = playerDTO.Id;
	this.x = x;
	this.y = y;
	this.angle = angle;
	var num = 0;
	for (var i in playerDTO.Creatures) {
		var creature = new Creature(playerDTO.Creatures[i]);
		creature.x = (num + +1 - playerDTO.Creatures.length)*(cardWidth/2 + cardEdgeWidth/2)
		creature.y = 0
		game.add.existing(creature);
		this.add(creature);
		num++;
	}
};

PlayerCreatures.prototype = Object.create(Phaser.Group.prototype);
PlayerCreatures.prototype.constructor = PlayerCreatures;

Creature = function(creatureDTO) {
	Phaser.Group.call(this, game);
	this.Id = creatureDTO.Id;
	this.Traits = creatureDTO.Traits;
	for (var i in creatureDTO.Cards) {
		var card = new Card(creatureDTO.Cards[i], 0, cardEdgeWidth/2 * i);
		game.add.existing(card);
		card.inputEnabled = true;
		this.add(card);
		card.selection = null;
		if ($.inArray("Used", card.getActiveProperty().Traits) != -1) {
			card.rotation = Math.PI/2;
		} else {
			card.events.onInputOver.add(propertyOver, card);
			card.events.onInputOut.add(propertyOut, card);
			addPropertyEvents(card);
		}
	}
	var back = new Phaser.Sprite(game, 0, creatureDTO.Cards.length*cardEdgeWidth/2, 'back');
	back.inputEnabled = true;
	back.anchor.setTo(0.5, 0.5);
    back.scale.setTo(0.5, 0.5);
    back.events.onInputUp.add(function (card) {
    	executeActionGrabFood(card.parent.Id);
    }, card);
    back.events.onInputOut.add(propertyOut, back);
    back.events.onInputOver.add(backOver, back);
	game.add.existing(back);
	this.add(back);
	var backBounds = new Phaser.Rectangle(-cardWidth/8, -cardHeight/8, cardWidth/4, cardHeight/4);
	this.Food = game.add.graphics();
	this.Food.x = back.x;
	this.Food.y = back.y;
	this.add(this.Food);
    this.Food.beginFill(0xFF0000, 1);
    for (var i in creatureDTO.Traits) {
        if (creatureDTO.Traits[i] == "Food") {
        	this.Food.drawCircle(backBounds.randomX, backBounds.randomY, 20);
       	}
    }
    this.Food.endFill();
    this.Food.beginFill(0x0000FF, 1);
    for (var i in creatureDTO.Traits) {
    	if (creatureDTO.Traits[i] == "Additional food") {
        	this.Food.drawCircle(backBounds.randomX, backBounds.randomY, 20);
        }
    }
    this.Food.endFill();
    this.Food.beginFill(0xFFFF00, 1);
    for (var i in creatureDTO.Traits) {
    	if (creatureDTO.Traits[i] == "Fat") {
            this.Food.drawCircle(backBounds.randomX, backBounds.randomY, 20);
        }
    }
    this.Food.endFill();
		/*this.Food = game.add.group();
	this.AdditionalFood = game.add.group();
	this.Fat = game.add.group();
	this.Food.x = back.x;
	this.Food.y = back.y;
	this.AdditionalFood.x = back.x;
	this.AdditionalFood.y = back.y;
	this.Fat.x = back.x;
	this.Fat.y = back.y;
	//this.add(this.Food);
	//this.add(this.AdditionalFood);
	//this.add(this.Fat);
	
    for (var i in creatureDTO.Traits) {
        if (creatureDTO.Traits[i] == "Food") {
			var item = game.add.graphics();
			this.Food.add(item);
			var x = backBounds.randomX;
			var y = backBounds.randomY;
			item.beginFill(0xFFFFFF, 1);
        	item.drawCircle(x, y, 20);
			item.endFill();
			item.beginFill(0xFF0000, 1);
        	item.drawCircle(x, y, 20);
			item.endFill();
       	}
    }
    for (var i in creatureDTO.Traits) {
    	if (creatureDTO.Traits[i] == "Additional food") {
			var item = game.add.graphics();
			this.AdditionalFood.add(item);
        	var x = backBounds.randomX;
			var y = backBounds.randomY;
			item.beginFill(0xFFFFFF, 1);
        	item.drawCircle(x, y, 21);
			item.endFill();
			item.beginFill(0x0000FF, 1);
        	item.drawCircle(x, y, 20);
			item.endFill();
       	}
    }
    for (var i in creatureDTO.Traits) {
    	if (creatureDTO.Traits[i] == "Fat") {
			var item = game.add.graphics();
			this.Fat.add(item);
            var x = backBounds.randomX;
			var y = backBounds.randomY;
			item.beginFill(0xFFFFFF, 1);
        	item.drawCircle(x, y, 21);
			item.endFill();
			item.beginFill(0xFFFF00, 1);
        	item.drawCircle(x, y, 20);
			item.endFill();
        }
    }*/
};

Creature.prototype = Object.create(Phaser.Group.prototype);
Creature.prototype.constructor = Creature;

function addPropertyEvents(card) {
	var traits = card.getActiveProperty().Traits; 
	if ($.inArray("Grazing", traits) != -1) {
        card.events.onInputUp.add(function (card) {
			executeActionGrazing(card.getActiveProperty().Id);
		}, card);
    } 
	if ($.inArray("Hibernation", traits) != -1) {
        card.events.onInputUp.add(function (card) {
			executeActionHibernation(card.parent.Id);
		}, card);
    } else if ($.inArray("Piracy", traits) != -1) {
        card.events.onInputDown.add(function (card) {
			startSelection(card.parent, card.parent, onSelectPiracyTarget);
		}, card);
    } else if ($.inArray("Carnivorous", traits) != -1) {
		card.events.onInputDown.add(function (card) {
			startSelection(card.parent, card.parent, onSelectAttackTarget);
		}, card);
    }
}

function propertyOver(card, pointer) {
	if (card.selection == null) {
		card.selection = game.add.graphics();
		card.parent.parent.add(card.selection);
		var creature = card.parent;
		card.selection.lineStyle(1, 0x000000, 1);
		card.selection.drawRoundedRect(creature.position.x + 4 - cardWidth/4, creature.position.y + card.position.y - cardHeight/4  + 4 , cardWidth/2-8, cardEdgeWidth/2-5, 3);
		var property = card.getActiveProperty();
		if (property.pair) {
			var pairCard = getPairProperty(card);
			if (pairCard != null) {
				var creature = pairCard.parent;
				card.selection.drawRoundedRect(creature.position.x + 4 - cardWidth/4, creature.position.y + pairCard.position.y - cardHeight/4 + 4, cardWidth/2-8, cardEdgeWidth/2-5, 3);
			}
		}
	}
}

function backOver(card, pointer) {
	if (card.selection == null) {
		card.selection = game.add.graphics();
		card.parent.parent.add(card.selection);
		var creature = card.parent;
		card.selection.lineStyle(1, 0x000000, 1);
		card.selection.drawRoundedRect(creature.position.x + 8 - cardWidth/4, creature.position.y + card.position.y - cardHeight/4  + 8 , cardWidth/2-16, cardHeight/2-16, 3);
	}
}

function propertyOut(card, pointer) {
	card.selection.destroy();
	card.selection = null;
}

function cardOver(card, pointer) {
	card.bringToTop();
    card.scale.y = 1;
    card.scale.x = 1;
}

function cardUp(card, pointer) {
	if (card.input.pointerTimeUp()-card.input.pointerTimeDown() < 70) {
		card.flipped = !card.flipped;
		card.rotation = Math.PI - card.rotation;
	}
}

function cardOut(card, pointer) {
    card.scale.setTo(0.5, 0.5);
}

function cardDragStart(card) {
	card.scale.setTo(0.5, 0.5);
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
	for (var i = creature.children.length-2; i>=0; i++) {
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
		this.rotation = Math.PI - this.rotation;
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
    	alert(JSON.stringify(this.properties[0].Traits));
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

function onSelectAttackTarget(source, pointer) {
	var target = getCreatureAtPoint(pointer.position);
	if (source == null || target == null) {
		return
	}
	executeActionAttack(currentPlayerId, source.Id, target.Id);
}

function onSelectPiracyTarget(source, pointer) {
	var target = getCreatureAtPoint(pointer.position);
	if (source == null || target == null) {
		return
	}
	if ($.inArray("Food", target.Traits) != -1) {
		executeActionPiracy(source.Id, target.Id, "Food");
	} else if ($.inArray("Additional food", target.Traits) != -1) {
		executeActionPiracy(source.Id, target.Id, "Additional food");
	} else {
		return
	}
}