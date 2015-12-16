var game = new Phaser.Game(1000, 800, Phaser.AUTO, 'game_holder', { preload: preload, create: create, update: update, render: render});
var gameOverlay;
var debug;
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
var voteStart = false;
var selectionArrow;
var selectionRect = null;
var currentGameState = null;
var messagesQueue = [];
var messageInProcessing = null;
var currentPhaseText = null;
var turnIndicatorText = null;
var click;

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
	game.load.spritesheet('pass','assets/pass.png', 150, 31, 2);
	game.load.spritesheet('endturn', 'assets/endturn.png', 150, 31, 2);
	game.load.spritesheet('vote', 'assets/vote.png', 150, 31, 2);
	game.load.image('chain', 'assets/copper-chain-btf-0292-sm.png');
	game.load.image('text', 'assets/text.png');
	game.load.audio('click', 'assets/click.wav');
}

function create() {
	game.input.addMoveCallback(mouseMoveCallback,this);
	game.input.onUp.add(mouseUp, this);
	game.add.tileSprite(0, 0, game.width, game.height, 'table');
	mainArea = new Phaser.Rectangle(10, 10, game.width-20, game.height-cardHeight-10);
	handArea = new Phaser.Rectangle(10, game.height-cardHeight+10, game.width-controlAreaWidth-30, cardHeight-20);
	controlArea = new Phaser.Rectangle(game.width-controlAreaWidth-10, game.height-cardHeight+10, controlAreaWidth, cardHeight-20);
	var currentPhaseTextBackground = game.add.sprite(controlArea.x + 10, controlArea.y + 10, 'text');
	var turnIndicatorTextBackground = game.add.sprite(controlArea.x + 10, controlArea.y + 50, 'text');
	var button = game.add.button(controlArea.x + 10, controlArea.y + 90, 'pass', pass, this, 0, 0, 1, 0);
	button = game.add.button(controlArea.x + 10, controlArea.y + 130, 'endturn', endTurn, this, 0, 0, 1, 0);
	button = game.add.button(controlArea.x + 10, controlArea.y + 170, 'vote', vote, this, 0, 0, 1, 0);
	var style = {
		font : "bold 20px Calibri",
		wordWrap: true,
		wordWrapWidth: 150,
		align: "center"
	};
	currentPhaseText = game.add.text(0, 0, "-", style);
	currentPhaseText.anchor.set(0.5);
	currentPhaseText.x = currentPhaseTextBackground.x + 75;
	currentPhaseText.y = currentPhaseTextBackground.y + 15;
	turnIndicatorText = game.add.text(0, 0, "-", style);
	turnIndicatorText.anchor.set(0.5);
	turnIndicatorText.x = turnIndicatorTextBackground.x + 75;
	turnIndicatorText.y = turnIndicatorTextBackground.y+ 15;
	game.physics.startSystem(Phaser.Physics.ARCADE);
	debug = game.add.graphics(0, 0);
	debug.lineStyle(2, 0xFFFFFF, 1);
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
	click = game.add.audio('click');
}

function addMessage(message) {
	messagesQueue.push(message);
}

function processMessage(message) {
	if (message.Type == MESSAGE_EXECUTED_ACTION) {
		showAction(message.Value, message.Value.State);
		updateGameState(message.Value.State)
		return
	}
	if (message.Type == MESSAGE_CHOICES_LIST) {
		updateGameState(message.Value.State)
		availableActions = message.Value.Actions;
		turnIndicatorText.setText("You choose");
		click.play();
		messageInProcessing = null;
		return
	}
	if (message.Type == MESSAGE_LOBBIES_LIST) {
		updateLobbiesList(message.Value);
		messageInProcessing = null;
		return
	}
	updateGameState(message.Value.State)
}


function updateGameState(state) {
	if (currentGameState == null) {
		currentPlayerId=state.CurrentPlayerId;
		playerId = state.PlayerId;
		localStorage.setItem("PlayerId", playerId);
		updateFoodBank(state.FoodBank);
		updatePlayers(state.Players);
		updateHand(state.PlayerCards);
		currentGameState = state;
	}
	currentGameState = state;
}

function showAction(msg, state) {
	switch(msg.Action.Type) {
		case "Add creature":
			showActionAddCreature(msg, state);
			break;	
		case "Add single property":
			showActionAddSingleProperty(msg, state);
			break;
		case "Add pair property":
			showActionAddPairProperty(msg, state);
			break;
		case "Remove creature":
			var creature = findCreature(msg.Action.Arguments.Creature);
			var player = creature.parent.parent;
			for (var i in creature.Properties.children) {
				if (creature.Properties.getChildAt(i).selection != null) {
					creature.Properties.getChildAt(i).selection.destroy();
				}
			}
			if (creature.back.selection != null) {
					creature.back.selection.destroy();
			}
			creature.destroy();
			arrangePlayerCreatures(player);
			messageInProcessing = null;
			break;
		case "New phase":
			currentPhaseText.setText(msg.Action.Arguments.Phase);	
			if (msg.Action.Arguments.Phase = "Development") {
				for (var i in players.children) {
					for (var j in players.getChildAt(i).Creatures.children) {
						var creature = players.getChildAt(i).Creatures.getChildAt(j);
						while (creature.Food.children.length > 0)
							creature.Food.getChildAt(0).destroy();
						while (creature.AdditionalFood.children.length > 0)
							creature.AdditionalFood.getChildAt(0).destroy();
					}
				}
			}
			messageInProcessing = null;
			break;
		case "Determine food bank":
			updateFoodBank(state.FoodBank)
			messageInProcessing = null;
			break;
		case "Remove card":
			var card = null;
			for (var i in players.children) {
				for (var j in players.getChildAt(i).Creatures.children) {
					for (var k in players.getChildAt(i).Creatures.getChildAt(j).Properties.children) {
						if (players.getChildAt(i).Creatures.getChildAt(j).Properties.getChildAt(k).Id == msg.Action.Arguments.Card) {
							card = players.getChildAt(i).Creatures.getChildAt(j).Properties.getChildAt(k);
							break;
						}
					}
					if (card != null) break;
				}
				if (card != null) break;
			}
			if (card != null) {
				card.destroy();
			}
			messageInProcessing = null;
			break;
		case "Destroy bank food":
			if (foodBank.children.length > 0)
				foodBank.getChildAt(0).destroy();
			messageInProcessing = null;
			break;
		case "Take cards":
			updateHand(state.PlayerCards)
		case "Get food from bank":
			if (foodBank.children.length > 0)
				foodBank.getChildAt(0).destroy();
			messageInProcessing = null;
			break;
		case "Add trait":
			var creature = findCreature(msg.Action.Arguments.Source);
			if (creature != null) {
				switch (msg.Action.Arguments.Trait) {
					case "Additional food":
						var backBounds = new Phaser.Rectangle(-cardWidth/8, -cardHeight/8, cardWidth/4, cardHeight/4);
						var circle = game.add.graphics();
						creature.AdditionalFood.add(circle);
						circle.x = backBounds.randomX;
						circle.y = backBounds.randomY;
						circle.beginFill(0x0000FF, 1);
						circle.drawCircle(0, 0, 10);
						circle.endFill();
						break;
					case "Food":
						var backBounds = new Phaser.Rectangle(-cardWidth/8, -cardHeight/8, cardWidth/4, cardHeight/4);
						var circle = game.add.graphics();
						creature.Food.add(circle);
						circle.x = backBounds.randomX;
						circle.y = backBounds.randomY;
						circle.beginFill(0xFF0000, 1);
						circle.drawCircle(0, 0, 10);
						circle.endFill();
						break;
					case "Fat":
						var backBounds = new Phaser.Rectangle(-cardWidth/8, -cardHeight/8, cardWidth/4, cardHeight/4);
						var backBounds = new Phaser.Rectangle(-cardWidth/8, -cardHeight/8, cardWidth/4, cardHeight/4);
						var circle = game.add.graphics();
						creature.Fat.add(circle);
						circle.x = backBounds.randomX;
						circle.y = backBounds.randomY;
						circle.beginFill(0xFFFF00, 1);
						circle.drawCircle(0, 0, 10);
						circle.endFill();
						break;
				}
			}
			var card = null;
			for (var i in players.children) {
				for (var j in players.getChildAt(i).Creatures.children) {
					for (var k in players.getChildAt(i).Creatures.getChildAt(j).Properties.children) {
						if (players.getChildAt(i).Creatures.getChildAt(j).Properties.getChildAt(k).getActiveProperty().Id == msg.Action.Arguments.Source) {
							card = players.getChildAt(i).Creatures.getChildAt(j).Properties.getChildAt(k);
							break;
						}
					}
					if (card != null) break;
				}
				if (card != null) break;
			}
			if (card != null && msg.Action.Arguments.Trait == "Used") {
				card.rotation = Math.PI/2;
			}
			messageInProcessing = null;
			break;
		case "Remove trait":
			var creature = findCreature(msg.Action.Arguments.Source);
			if (creature != null) {
				switch (msg.Action.Arguments.Trait) {
					case "Additional food":
						if (creature.AdditionFood.children.length != 0)
							creature.AdditionFood.getChildAt(0).destroy();
						break;
					case "Food":
						if (creature.Food.children.length != 0)
							creature.Food.getChildAt(0).destroy();
						break;
					case "Fat":
						if (creature.Fat.children.length != 0)
							creature.Fat.getChildAt(0).destroy();
						break;
				}
			}
			var card = null;
			for (var i in players.children) {
				for (var j in players.getChildAt(i).Creatures.children) {
					for (var k in players.getChildAt(i).Creatures.getChildAt(j).Properties.children) {
						if (players.getChildAt(i).Creatures.getChildAt(j).Properties.getChildAt(k).getActiveProperty().Id == msg.Action.Arguments.Source) {
							card = players.getChildAt(i).Creatures.getChildAt(j).Properties.getChildAt(k);
							break;
						}
					}
					if (card != null) break;
				}
				if (card != null) break;
			}
			if (card != null && msg.Action.Arguments.Trait == "Used") {
				card.rotation = 0;
			}
			messageInProcessing = null;
			break;
		default:
			messageInProcessing = null;
	}
}

function showActionAddCreature(msg, state) {
	var player = findPlayer(msg.Action.Arguments.Player)
			var creature = null;
			var oldCreatures = currentGameState.Players[msg.Action.Arguments.Player].Creatures;
			var newCreatures = state.Players[msg.Action.Arguments.Player].Creatures;
			for (var i in newCreatures) {
				var found = false
				for (var j in oldCreatures)
				{
					if (newCreatures[i].Id == oldCreatures[j].Id) {
						found = true;
						break;
					}
				}
				if (!found) {
					creature = new Creature(newCreatures[i], 0, 0);
					game.add.existing(creature);
					player.Creatures.add(creature);
					break;
				}
			}
			if (player.Id == playerId) {
				var card = findCardInHand(msg.Action.Arguments.Card);
				hand.remove(card);
				arrangeCardsInHand();
				arrangePlayerCreatures(player);
			} else {
				creature.y = mainArea.height;
				arrangePlayerCreatures(player);
			}
			messageInProcessing = null;
}

function showActionAddSingleProperty(msg, state) {
	var creature = null;
			var player = null;
			var card = null;
			for (var i in players.children) {
				for (var j in players.getChildAt(i).Creatures.children) {
					if (players.getChildAt(i).Creatures.getChildAt(j).Id == msg.Action.Arguments.Creature) {
						creature = players.getChildAt(i).Creatures.getChildAt(j);
						break;
					}
				}
				if (creature != null) {
					player = players.getChildAt(i);
					break;
				}
			}
			var oldCards = null;
			if (currentGameState.Players[player.Id].Creatures[creature.Id]) {
				oldCards = currentGameState.Players[player.Id].Creatures[creature.Id].Cards;
			}
			var newCards = null;
			if (state.Players[player.Id].Creatures[creature.Id]) {
				newCards = state.Players[player.Id].Creatures[creature.Id].Cards;
			}
			for (var i in newCards) {
				var found = false
				for (var j in oldCards)
				{
					if (newCards[i].Id == oldCards[j].Id) {
						found = true;
						break;
					}
				}
				if (!found) {
					card = new Card(newCards[i], 0, 0);
					game.add.existing(card);
					addCard(creature, card)
					break;
				}
			}
			if (player.Id == playerId) {
				var card = findCardInHand(card.Id);
				hand.remove(card);
				arrangeCardsInHand();
			} else {
			}
			messageInProcessing = null;
}

function showActionAddPairProperty(msg, state) {
	var firstCreature = null;
	var secondCreature = null;
			var player = null;
			var firstCard = null;
			var secondCard = null;
			for (var i in players.children) {
				for (var j in players.getChildAt(i).Creatures.children) {
					if (players.getChildAt(i).Creatures.getChildAt(j).Id == msg.Action.Arguments.Pair[0]) {
						firstCreature = players.getChildAt(i).Creatures.getChildAt(j);
					}
					if (players.getChildAt(i).Creatures.getChildAt(j).Id == msg.Action.Arguments.Pair[1]) {
						secondCreature = players.getChildAt(i).Creatures.getChildAt(j);
					}
				}
				if (firstCreature != null) {
					player = players.getChildAt(i);
					break;
				}
			}
			var oldCards = null;
			if (currentGameState.Players[player.Id].Creatures[firstCreature.Id]) {
				oldCards = currentGameState.Players[player.Id].Creatures[firstCreature.Id].Cards;
			}
			var newCards = null;
			if (state.Players[player.Id].Creatures[firstCreature.Id]) {
				newCards = state.Players[player.Id].Creatures[firstCreature.Id].Cards;
			}
			for (var i in newCards) {
				var found = false
				for (var j in oldCards)
				{
					if (newCards[i].Id == oldCards[j].Id) {
						found = true;
						break;
					}
				}
				if (!found) {
					firstCard = new Card(newCards[i], 0, 0);
					secondCard = new Card(newCards[i], 0, 0);
					game.add.existing(firstCard);
					game.add.existing(secondCard);
					addCard(firstCreature, firstCard);
					addCard(secondCreature, secondCard);
					break;
				}
			}
			if (player.Id == playerId) {
				var card = findCardInHand(firstCard.Id);
				hand.remove(card);
				arrangeCardsInHand();
			} else {
			}
			messageInProcessing = null;
}

function update() {
	if (messageInProcessing == null && messagesQueue.length != 0) {
		messageInProcessing = messagesQueue.shift();
		processMessage(messageInProcessing);
	}
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

function updateFoodBank(count) {
	while (foodBank.children.length > 0)
		foodBank.getChildAt(0).destroy();
	var rectangle = new Phaser.Rectangle(-50, -50, 100, 100);
	for (var i = 0; i<count; i++) {
		var circle = game.add.graphics();
		foodBank.add(circle);
		circle.lineStyle(0);
		circle.beginFill(0xFF0000, 1);
		circle.drawCircle(rectangle.randomX, rectangle.randomY, 10);
		circle.endFill();
	}
}

function updateHand(handDTO) {
	hand.removeAll(true);
	for (var i in handDTO) {
		var card = new Card(handDTO[i], 0, 0);
		card.events.onInputOver.add(cardOver, card);
    	card.events.onInputOut.add(cardOut, card);
	    card.events.onInputUp.add(cardUp, card);
	    card.events.onDragStart.add(cardDragStart, card);
	    card.events.onDragStop.add(cardDragStop, card);
	    card.events.onDragUpdate.add(cardDragUpdate, card);
	    card.input.enableDrag();
		hand.add(card);
	}
	arrangeCardsInHand();
}

function arrangeCardsInHand() {
	var y = handArea.halfHeight;
	var startX = (handArea.width-(cardWidth*hand.children.length/2*3/2))/2;
	if (startX < 0) {
		startX = cardWidth/4;
	}
	var offset = (handArea.width-startX*2)/(hand.children.length);
	
	for (var i in hand.children) {
		hand.getChildAt(i).x = startX + (+i + +0.5)*offset;
		hand.getChildAt(i).y = y;
	}
}

function updatePlayers(playersDTO) {
	if (selectionArrow != null) {
		selectionArrow.arrow.destroy();
		selectionArrow = null;
	}
	players.removeAll(true);
	var startAngle = 180;
	var deltaAngle = 360/Object.keys(playersDTO).length;
	var radiusX = mainArea.halfWidth - cardHeight/4;
	var radiusY = mainArea.halfHeight - cardHeight/4;
	var playerIndex = 0;
	for (var i in playersDTO) {
		if (playersDTO[i].Id == playerId) {
			var playersCreatures = new PlayerCreatures(playersDTO[i], mainArea.halfWidth, mainArea.halfHeight+radiusY, 0)
    		game.add.existing(playersCreatures);
    		players.add(playersCreatures);
			break;
		}
	}
	var angle = 0;
	for (var i in playersDTO) {
		if (playersDTO[i].Id == playerId) {
			continue;	
		}
		angle += deltaAngle;
		var playersCreatures = new PlayerCreatures(playersDTO[i], mainArea.halfWidth-Math.sin(angle*Math.PI/180)*radiusX, mainArea.halfHeight+Math.cos(angle*Math.PI/180)*radiusY, angle)
        game.add.existing(playersCreatures);
        players.add(playersCreatures);
	}
}

function findPlayer(id) {
	for (var i in players.children) {
		if (players.getChildAt(i).Id == id) {
			return players.getChildAt(i);
		}
	}
}

function findCreature(id) {
	for (var i in players.children) {
		for (var j in players.getChildAt(i).Creatures.children) {
			if (players.getChildAt(i).Creatures.getChildAt(j).Id == id) {
				return players.getChildAt(i).Creatures.getChildAt(j);
			}
		}
	}
}

function findCardInHand(id) {
	for (var i in hand.children) {
		if (hand.getChildAt(i).Id == id) {
			return hand.getChildAt(i);
		}
	}
}

PlayerCreatures = function(playerDTO, x, y, angle) {
	Phaser.Group.call(this, game);
	this.x = x;
	this.y = y;
	this.angle = angle;
	this.Id = playerDTO.Id;
	this.Creatures = game.add.group();
	this.add(this.Creatures)
	for (var i in playerDTO.Creatures) {
		var creature = new Creature(playerDTO.Creatures[i], 0, 0);
		game.add.existing(creature);
		this.Creatures.add(creature);
	}
	arrangePlayerCreatures(this);
};

function arrangePlayerCreatures(player) {
	var totalCreatureWidthHalf = cardWidth/2 * player.Creatures.children.length/2;
	for (var j in player.Creatures.children) {
		game.add.tween(player.Creatures.getChildAt(j)).to({x: (+j + +1)*cardWidth/2-totalCreatureWidthHalf, y: 0}, 1000, Phaser.Easing.Quadratic.InOut, true);
	}
}

PlayerCreatures.prototype = Object.create(Phaser.Group.prototype);
PlayerCreatures.prototype.constructor = PlayerCreatures;

Creature = function(creatureDTO, x, y) {
	Phaser.Group.call(this, game);
	this.x = x-cardWidth/4;
	this.y = y-cardHeight/4;
	this.Id = creatureDTO.Id;
	this.Traits = creatureDTO.Traits;
	this.Properties = game.add.group();
	this.add(this.Properties);
	//var back = new Phaser.Sprite(game, 0, creatureDTO.Cards.length*cardEdgeWidth/2, 'back');
	this.back = new Phaser.Sprite(game, 0, 0, 'back');
	this.back.inputEnabled = true;
	this.back.anchor.setTo(0.5, 0.5);
    this.back.scale.setTo(0.5, 0.5);
    this.back.events.onInputUp.add(function (card) {
    	executeActionGrabFood(card.parent.Id);
    }, card);
    this.back.events.onInputOut.add(propertyOut, this.back);
    this.back.events.onInputOver.add(backOver,this.back);
	game.add.existing(this.back);
	this.add(this.back);
	for (var i in creatureDTO.Cards) {
		var card = new Card(creatureDTO.Cards[i], 0, cardEdgeWidth/2 * i);
		addCard(this, card);
	}
	var backBounds = new Phaser.Rectangle(-cardWidth/8, -cardHeight/8, cardWidth/4, cardHeight/4);
	this.Food = game.add.group();
	this.Food.x = this.back.x;
	this.Food.y = this.back.y;
	this.AdditionalFood = game.add.group();
	this.AdditionalFood.x = this.back.x;
	this.AdditionalFood.y = this.back.y;
	this.Fat = game.add.group();
	this.Fat.x = this.back.x;
	this.Fat.y = this.back.y;
	this.add(this.Food);
	this.add(this.AdditionalFood);
	this.add(this.Fat);
    for (var i in creatureDTO.Traits) {
        if (creatureDTO.Traits[i] == "Food") {
			var circle = game.add.graphics();
			this.Food.add(circle);
			circle.x = backBounds.randomX;
			circle.y = backBounds.randomY;
			circle.beginFill(0xFF0000, 1);
			circle.drawCircle(0, 0, 10);
			circle.endFill();
       	}
    }
    for (var i in creatureDTO.Traits) {
    	if (creatureDTO.Traits[i] == "Additional food") {
        	var circle = game.add.graphics();
			this.Food.add(circle);
			circle.x = backBounds.randomX;
			circle.y = backBounds.randomY;
			circle.beginFill(0x0000FF, 1);
			circle.drawCircle(0, 0, 10);
			circle.endFill();;
        }
    }
    for (var i in creatureDTO.Traits) {
    	if (creatureDTO.Traits[i] == "Fat") {
            var circle = game.add.graphics();
			this.Food.add(circle);
			circle.x = backBounds.randomX;
			circle.y = backBounds.randomY;
			circle.beginFill(0xFFFF00, 1);
			circle.drawCircle(0, 0, 10);
			circle.endFill();
        }
    }
};

function addCard(creature, card) {
	card.y = -cardEdgeWidth/2 * (creature.Properties.children.length + 1);
	card.inputEnabled = true;
	creature.Properties.add(card);
	card.bringToTop();
	card.selection = null;
	if ($.inArray("Used", card.getActiveProperty().Traits) != -1) {
		card.rotation = Math.PI/2;
	} else {
		card.events.onInputOver.add(propertyOver, card);
		card.events.onInputOut.add(propertyOut, card);
		addPropertyEvents(card);
	}
	card.sendToBack();
}

Creature.prototype = Object.create(Phaser.Group.prototype);
Creature.prototype.constructor = Creature;
Creature.prototype.back = null;

function addPropertyEvents(card) {
	var traits = card.getActiveProperty().Traits; 
	if ($.inArray("Grazing", traits) != -1) {
        card.events.onInputUp.add(function (card) {
			executeActionGrazing(card.getActiveProperty().Id);
		}, card);
    } 
	if ($.inArray("Hibernation", traits) != -1) {
        card.events.onInputUp.add(function (card) {
			executeActionHibernation(card.parent.parent.Id);
		}, card);
    } else if ($.inArray("Piracy", traits) != -1) {
        card.events.onInputDown.add(function (card) {
			startSelection(card.parent.parent, card.parent.parent, onSelectPiracyTarget);
		}, card);
    } else if ($.inArray("Carnivorous", traits) != -1) {
		card.events.onInputDown.add(function (card) {
			startSelection(card.parent.parent, card.parent.parent, onSelectAttackTarget);
		}, card);
    } else if ($.inArray("Fat tissue", traits) != -1) {
        card.events.onInputUp.add(function (card) {
			executeActionBurnFat(card.parent.parent.Id);
		}, card);
	}
}

function propertyOver(card, pointer) {
	if (card.selection == null) {
		card.selection = game.add.graphics();
		card.parent.parent.parent.parent.add(card.selection);
		var creature = card.parent.parent;
		card.selection.lineStyle(1, 0x000000, 1);
		card.selection.drawRoundedRect(creature.position.x + 4 - cardWidth/4, creature.position.y + card.position.y - cardHeight/4  + 4 , cardWidth/2-8, cardEdgeWidth/2-8, 3);
		var property = card.getActiveProperty();
		if (property.pair) {
			var pairProperty = getPairProperty(card);
			if (pairProperty != null) {
				var creature = pairProperty.parent.parent;
				card.selection.drawRoundedRect(creature.position.x + 4 - cardWidth/4, creature.position.y + pairProperty.position.y - cardHeight/4 + 4, cardWidth/2-8, cardEdgeWidth/2-8, 3);
			}
		}
	}
}

function backOver(card, pointer) {
	if (card.selection == null) {
		card.selection = game.add.graphics();
		card.parent.parent.parent.add(card.selection);
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
				startSelection(creature.Properties, arguments, onSelectSecondPairCreature);
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
		intersectedCreature.parent.parent.add(selectionRect);
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
	var maxIntersectObject = null;
	var maxIntersectArea = 0;
	if (Phaser.Rectangle.intersects(rectangle,mainArea)) {
		for (var i in players.children) {
			for (var j in players.getChildAt(i).Creatures.children) {
				var creature = players.getChildAt(i).Creatures.getChildAt(j);
				var bounds = creature.back.getBounds();
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
			for (var j in players.getChildAt(i).Creatures.children) {
				var creature = players.getChildAt(i).Creatures.getChildAt(j);
				var bounds = creature.back.getBounds();
				if (Phaser.Rectangle.containsPoint(bounds, point)) {
					return creature;
				}
			}
		}
	}
	return null;
}

function getPairProperty(firstProperty) {
	var player = firstProperty.parent.parent.parent.parent;
	for (var j in player.Creatures.children) {
		var creature = player.Creatures.getChildAt(j);
		if (firstProperty.parent.parent.Id == creature.Id) {
			continue;
		}
		for (var k = 0; k<creature.Properties.children.length; k++) {
			if (creature.Properties.getChildAt(k).Id == firstProperty.Id) {
				return creature.Properties.getChildAt(k);
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
	for (var i = creature.Properties.children.length-2; i>=0; i++) {
		if (Phaser.Rectangle.containsPoint(creature.Properties.getChildAt(i).getBounds(), point)) {
			return creature.getChildAt(i);
		}
	}
	if (Phaser.Rectangle.containsPoint(creature.back.getBounds(), point)) {
			return creature.back;
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
	//arrow.x = startObject.worldPosition.x + startObject.getBounds().width/2;
	//arrow.y = startObject.worldPosition.y + startObject.getBounds().height/2;
	arrow.x = startObject.worldPosition.x;
	arrow.y = startObject.worldPosition.y;
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