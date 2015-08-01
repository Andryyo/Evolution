// Constants
package main

type FilterType int

const (
	FILTER_DENY FilterType = iota
	FILTER_ALLOW
	FILTER_MODIFY
	FILTER_ACTION_REPLACE
	FILTER_ACTION_EXECUTE_BEFORE
	FILTER_ACTION_EXECUTE_AFTER
)

type FilterSourcePrototype int

const (
	FILTER_SOURCE_PARAMETER_PLAYER FilterSourcePrototype = iota
	FILTER_SOURCE_PARAMETER_PROPERTY 
	FILTER_SOURCE_PARAMETER_SOURCE_CREATURE
	FILTER_SOURCE_PARAMETER_TARGET_CREATURE
	FILTER_SOURCE_PARAMETER_CREATURE
	FILTER_SOURCE_PARAMETER_TRAIT
	FILTER_SOURCE_PARAMETER_ALL_PLAYERS
	FILTER_SOURCE_PARAMETER_LEFT_CREATURE
	FILTER_SOURCE_PARAMETER_RIGHT_CREATURE
	FILTER_SOURCE_PARAMETER_PAIR
	FILTER_SOURCE_PARAMETER_ANY_FOOD
	FILTER_SOURCE_PARAMETER_ALL_FOOD
	FILTER_SOURCE_PARAMETER_FOOD_AND_FAT_LIMIT
	FILTER_SOURCE_PARAMETER_FOOD_AND_FAT
	FILTER_SOURCE_PARAMETER_ALL_FOOD_AND_FAT
	FILTER_SOURCE_PARAMETER_PHASE
	FILTER_SOURCE_PARAMETER_SOURCE
	FILTER_SOURCE_PARAMETER_CREATURE_PROPERTIES
	FILTER_SOURCE_PARAMETER_FOOD_BANK_COUNT
	FILTER_SOURCE_PARAMETER_CURRENT_PLAYER
	FILTER_SOURCE_PARAMETER_ONE_OF_CURRENT_PLAYER_CARDS
	FILTER_SOURCE_PARAMETER_ONE_OF_CURRENT_PLAYER_CARDS_PROPERTIES
	FILTER_SOURCE_PARAMETER_ONE_OF_CURRENT_PLAYER_CREATURES_PAIR
	FILTER_SOURCE_PARAMETER_ONE_OF_CURRENT_PLAYER_CREATURES
)

func (t FilterSourcePrototype) GoString() string {
	switch t {
		case FILTER_SOURCE_PARAMETER_PHASE:
			return "Phase"
		case FILTER_SOURCE_PARAMETER_PLAYER:
			return "Player"
		case FILTER_SOURCE_PARAMETER_PROPERTY:
			return "Property"
		case FILTER_SOURCE_PARAMETER_SOURCE_CREATURE:
			return "Source creature"
		case FILTER_SOURCE_PARAMETER_TARGET_CREATURE:
			return "Target creature"
		case FILTER_SOURCE_PARAMETER_CREATURE:
			return "Creature"
		case FILTER_SOURCE_PARAMETER_TRAIT:
			return "Trait"
		case FILTER_SOURCE_PARAMETER_ALL_PLAYERS:
			return "All players"
		case FILTER_SOURCE_PARAMETER_RIGHT_CREATURE:
			return "Second creature in pair"
		case FILTER_SOURCE_PARAMETER_LEFT_CREATURE:
			return "First creature in pait"
		case FILTER_SOURCE_PARAMETER_ANY_FOOD:
			return "Any food"
		case FILTER_SOURCE_PARAMETER_ALL_FOOD:
			return "All food"
		case FILTER_SOURCE_PARAMETER_FOOD_AND_FAT_LIMIT:
			return "Food and fat"
		case FILTER_SOURCE_PARAMETER_SOURCE:
			return "Source"
		case FILTER_SOURCE_PARAMETER_PAIR:
			return "Pair"
		case FILTER_SOURCE_PARAMETER_CREATURE_PROPERTIES:
			return "Creature properties"
		case FILTER_SOURCE_PARAMETER_FOOD_BANK_COUNT:
			return "Food in bank"
		case FILTER_SOURCE_PARAMETER_CURRENT_PLAYER:
			return "Current player"
		case FILTER_SOURCE_PARAMETER_ONE_OF_CURRENT_PLAYER_CARDS:
			return "One of current player cards"
		case FILTER_SOURCE_PARAMETER_ONE_OF_CURRENT_PLAYER_CARDS_PROPERTIES:
			return "One of players properties"
		case FILTER_SOURCE_PARAMETER_ONE_OF_CURRENT_PLAYER_CREATURES_PAIR:
			return "Creatures pair"
		case FILTER_SOURCE_PARAMETER_ONE_OF_CURRENT_PLAYER_CREATURES:
			return "Player creature"
		default:
			return string(t)
	}
}

type TraitType int

const (
	TRAIT_PASS TraitType = iota
	TRAIT_USED
	TRAIT_TOOK_FOOD
	TRAIT_SHART_VISION
	TRAIT_CAMOUFLAGE
	TRAIT_BURROWING
	TRAIT_FED
	TRAIT_PAIR
	TRAIT_FAT_TISSUE
	TRAIT_FOOD
	TRAIT_ADDITIONAL_FOOD
	TRAIT_REQUIRE_FOOD
	TRAIT_FAT
	TRAIT_SIMBIOSYS
	TRAIT_PIRACY
)

func (t TraitType) GoString() string {
	switch t {
		case TRAIT_PASS:
			return "Pass"
		case TRAIT_TOOK_FOOD:
			return "Already took food"
		case TRAIT_SHART_VISION:
			return "Sharp vision"
		case TRAIT_CAMOUFLAGE:
			return "Camouflage"
		case TRAIT_BURROWING:
			return "Burrowing"
		case TRAIT_FED:
			return "Fed"
		case TRAIT_PAIR:
			return "Pair"
		case TRAIT_SIMBIOSYS:
			return "Simbiosys"
		case TRAIT_FAT_TISSUE:
			return "Fat tissue"
		case TRAIT_FOOD:
			return "Food"
		case TRAIT_ADDITIONAL_FOOD:
			return "Additional food"
		case TRAIT_REQUIRE_FOOD:
			return "Require food"
		case TRAIT_FAT:
			return "Fat"
		case TRAIT_PIRACY:
			return "Piracy"
		case TRAIT_USED:
			return "Used"
		default:
			return string(t)
	}
}

type PhaseType int

const (
	PHASE_DEVELOPMENT PhaseType = iota
	PHASE_FOOD_BANK_DETERMINATION
	PHASE_FEEDING
	PHASE_EXTINCTION
)

type ActionType int

const (
	ACTION_SEQUENCE ActionType = iota
	ACTION_SELECT
	ACTION_START_TURN
	ACTION_NEXT_PLAYER
	ACTION_ADD_CREATURE
	ACTION_ADD_SINGLE_PROPERTY
	ACTION_ADD_PAIR_PROPERTY
	ACTION_PASS
	ACTION_NEW_PHASE
	ACTION_ADD_TRAIT
	ACTION_REMOVE_TRAIT
	ACTION_ADD_FILTERS
	ACTION_ATTACK
	ACTION_DETERMINE_FOOD_BANK
	ACTION_REMOVE_CREATURE
	ACTION_REMOVE_CARD
	ACTION_REMOVE_PROPERTY
	ACTION_GET_FOOD_FROM_BANK
)

type ArgumentName int

const (
	PARAMETER_PROPERTY ArgumentName = iota
	PARAMETER_PHASE
	PARAMETER_PLAYER
	PARAMETER_PAIR
	PARAMETER_CARD
	PARAMETER_ACTIONS
	PARAMETER_CREATURE
	PARAMETER_TRAIT
	PARAMETER_SOURCE
	PARAMETER_FILTERS
	PARAMETER_SOURCE_CREATURE
	PARAMETER_TARGET_CREATURE
	
)