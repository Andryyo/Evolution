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


function executeAddCreatureAction(cardId) {
	var action = {
		Type: "Add creature",
		Arguments: {
			Card: cardId,
			Player: playerId
		}
	};
	return executeAction(action);
}

function executeAddPropertyAction(creatureId, propertyId) {
	var action = {
		Type: "Add single property",
		Arguments: {
			Creature: creatureId,
			Property: propertyId
		}
	};
	return executeAction(action);
}

function executeAddPairPropertyAction(firstCreatureId, secondCreatureId, propertyId) {
	var action = {
		Type: "Add pair property",
		Arguments: {
			Pair: [
				firstCreatureId,
				secondCreatureId
			],
			Property: propertyId,
		}
	};
	return executeAction(action);
}

function executeActionGrazing(propertyId) {
	var action = {
		Type: "Destroy bank food",
		Arguments: {
			Property: propertyId
		}
	};
	return executeAction(action);
}

function executeActionHibernation(creatureId) {
	var action = {
		Type: "Hibernate",
		Arguments: {
			Creature: creatureId
		}
	};
	return executeAction(action);
}

function executeActionAttack(playerId, sourceCreatureId, targetCreatureId) {
	var action = {
		Type: "Attack",
		Arguments: {
			Player: playerId,
			SourceCreature: sourceCreatureId,
			TargetCreature: targetCreatureId
		}
	};
	return executeAction(action);
}

function executeActionPiracy(sourceCreatureId, targetCreatureId, trait) {
	var action = {
		Type: "Piracy",
		Arguments: {
			SourceCreature: sourceCreatureId,
			TargetCreature: targetCreatureId,
			Trait: trait
		}
	};
	return executeAction(action);
}

function executeActionGrabFood(creatureId) {
	var action = {
		Type: "Get food from bank",
		Arguments: {
			Creature: creatureId
		}};
	return executeAction(action);
}