var game = new Phaser.Game(1366, 768, Phaser.AUTO, 'game_holder', { preload: preload, create: create, update: update, render: render});
var cardHeight = 254;
var cardWidth = 182;
var cardEdgeWidth = 35;
var handArea = new Phaser.Rectangle(0, 768-cardHeight, 1366, cardHeight);
var mainArea = new Phaser.Rectangle(0, 0, 1366, 768-cardHeight);
var hand = null;
var players = null;
var chosenSprite = null
var availableActions = null;
var currentPlayerId;
var playerId;
var selectionRect;

function preload() {
	game.load.spritesheet('cards','assets/spritesheet.png',cardWidth,cardHeight,20);
	game.load.image('back','assets/back.png');
}

function create() {
	game.physics.startSystem(Phaser.Physics.ARCADE);
	hand = game.add.group();
	hand.x = handArea.x;
	hand.y = handArea.y;
	players = game.add.group();
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
	if (chosenSprite != null) {
		game.debug.spriteInputInfo(chosenSprite, 32, 32);
	}
}

var socket = new WebSocket("ws://127.0.0.1:8081/client");

socket.onopen = function() {
	var textArea = document.getElementById("log");
    textArea.value = "";
};

socket.onmessage = function(event) {
	var textArea = document.getElementById("log");
	textArea.value = textArea.value + '\n' + event.data;
	textArea.scrollTop = textArea.scrollHeight;
	var obj = JSON.parse(event.data);
	if (obj.Type == 0) {
		showAction(obj.Value);
	}
	if (obj.Type == 1) {
		availableActions = obj.Value;
	}
};

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
	for (var i in availableActions) {
		if (JSON.stringify(availableActions[i]) === JSON.stringify(action)) {
			availableActions = null;
			socket.send(i);
			return true;
		}
	}
	return false;
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
	for (var i in availableActions) {
		if (JSON.stringify(availableActions[i]) === JSON.stringify(action)) {
			availableActions = null;
			socket.send(i);
			return true;
		}
	}
	return false;
}

function showAction(action) {
	updateGameState(action.State)
}

function updateGameState(state) {
	currentPlayerId=state.CurrentPlayerId;
	playerId = state.PlayerId;
	updatePlayers(state.Players);
	updateHand(state.PlayerCards);
}

function updateHand(handDTO) {
	hand.removeAll();
	var y = handArea.halfHeight;
	var startX = (handArea.width-(cardWidth*handDTO.length/2*3/2))/2;
	if (startX < 0) {
		startX = 0;
	}
	var offset = (handArea.width-startX*2)/(handDTO.length);
	
	for (var i in handDTO) {
		hand.add(new Card(handDTO[i], startX + i*offset, y));
	}
}

function updatePlayers(playersDTO) {
	players.removeAll();
	var startAngle = 180;
	var deltaAngle = 360/playersDTO.length;
	var radius = mainArea.halfHeight - cardHeight/4;
	for (var i in playersDTO) {
		if (playersDTO[i].Id == playerId) {
			var playerIndex = i;
			break;
		}
	}
	var angle = 0;
	for (var i = playerIndex; i<playersDTO.length; i++) {
		var playersCreatures = new PlayerCreatures(playersDTO[i], mainArea.halfWidth-Math.sin(angle*Math.PI/180)*radius, mainArea.halfHeight+Math.cos(angle*Math.PI/180)*radius, angle)
        game.add.existing(playersCreatures);
        players.add(playersCreatures);
        angle += deltaAngle;
	}
	for (var i = 0; i<playerIndex; i++) {
		var playersCreatures = new PlayerCreatures(playersDTO[i], mainArea.halfWidth-Math.sin(angle*Math.PI/180)*radius, mainArea.halfHeight+Math.cos(angle*Math.PI/180)*radius, angle)
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
		var creature = new Creature(playerDTO.Creatures[i], (+i + +0.5)*cardWidth/2-totalCreatureWidthHalf, 0);
		game.add.existing(creature);
		this.add(creature);
	}
};

PlayerCreatures.prototype = Object.create(Phaser.Group.prototype);
PlayerCreatures.prototype.constructor = PlayerCreatures;

Creature = function(creatureDTO, x, y) {
	Phaser.Group.call(this, game);
	this.x = x;
	this.y = y;
	this.id = creatureDTO.Id;
	for (var i in creatureDTO.Cards) {
		var card = new Card(creatureDTO.Cards[i], 0, cardEdgeWidth/2 * i);
		game.add.existing(card);
		this.add(card);
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
	chosenSprite = card;
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
	if (Phaser.Rectangle.intersects(card.getBounds(),mainArea)) {
		for (var i in players.children) {
			for (var j in players.getChildAt(i).children) {
				var creature = players.getChildAt(i).getChildAt(j);
				if (Phaser.Rectangle.intersects(card.getBounds(), creature.getBounds())) {
					if (card.properties.length == 1 || !card.flipped) {
						var property = card.properties[0];
					} else {
						var property = card.properties[1];
					}
					if (executeAddPropertyAction(creature.id, property.Id)) {
						return;
					}
				}
			}
		}
		if (executeAddCreatureAction(card.id)) {
			return;
		}
	}
	card.position = card.input.dragStartPoint.clone();
}

function cardDragUpdate(card) {
	if (Phaser.Rectangle.intersects(card.getBounds(),mainArea)) {
		for (var i in players.children) {
			for (var j in players.getChildAt(i).children) {
				var creature = players.getChildAt(i).getChildAt(j);
				if (Phaser.Rectangle.intersects(card.getBounds(), creature.getBounds())) {
					if (selectionRect == null) {
						selectionRect = game.add.graphics();
						//selectionRect.drawRect(-cardWidth/4-10, -cardHeight/4-10, cardWidth/2+20, cardHeight/2+20);
						selectionRect.lineColor = "#ffffff";
						selectionRect.lineStyle(2, 0x0000FF, 1);
						selectionRect.drawRect(-10, -10, cardWidth/2+20, cardHeight/2+20);
					}
					selectionRect.graphicsData[0].width = creature.getBounds().width+20;
					selectionRect.graphicsData[0].height = creature.getBounds().height+20;
					selectionRect.x = creature.getBounds().x
					selectionRect.y = creature.getBounds().y
					return;
				}
			}
		}
	}
	if (selectionRect != null) {
		selectionRect.destroy();
		selectionRect = null;
	}
}

Card = function(cardDTO, x, y) {
	Phaser.Sprite.call(this, game, x, y, 'cards');
	this.anchor.setTo(0.5, 0.5);
	this.scale.setTo(0.5, 0.5);
	game.physics.arcade.enable(this);
    this.inputEnabled = true;
    this.input.enableDrag();
    this.id = cardDTO.Id;
    this.properties = cardDTO.Properties;
    this.events.onInputOver.add(cardOver, this);
    this.events.onInputOut.add(cardOut, this);
    this.events.onInputUp.add(cardUp, this);
    this.events.onDragStart.add(cardDragStart, this);
    this.events.onDragStop.add(cardDragStop, this);
    this.events.onDragUpdate.add(cardDragUpdate, this);
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
