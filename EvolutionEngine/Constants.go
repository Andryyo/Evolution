// Constants
package EvolutionEngine

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
	FILTER_SOURCE_PARAMETER_ONE_OF_CREATURES
	FILTER_SOURCE_PARAMETER_BANK_CARDS_COUNT
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
		case FILTER_SOURCE_PARAMETER_ONE_OF_CREATURES:
			return "One of creatures"
		default:
			return string(t)
	}
}

type TraitType string

const (
	TRAIT_PASS TraitType = "Pass"
	TRAIT_USED TraitType = "Used"
	TRAIT_TOOK_FOOD TraitType = "Food"
	TRAIT_SHARP_VISION TraitType = "Sharp vision"
	TRAIT_CAMOUFLAGE TraitType = "Camouflage"
	TRAIT_BURROWING TraitType = "Burowing"
	TRAIT_FED TraitType = "Fed"
	TRAIT_PAIR TraitType = "Pair"
	TRAIT_FAT_TISSUE TraitType = "Fat tissue"
	TRAIT_FOOD TraitType = "Food"
	TRAIT_ADDITIONAL_FOOD TraitType = "Additional food"
	TRAIT_REQUIRE_FOOD TraitType = "Require food"
	TRAIT_FAT TraitType = "Fat"
	TRAIT_SYMBIOSIS TraitType = "Symbiosys"
	TRAIT_PIRACY TraitType = "Piracy"
	TRAIT_HIGH_BODY_WEIGHT TraitType = "High body weight"
	TRAIT_GRAZING TraitType = "Grazing"
	TRAIT_TAIL_LOSS TraitType = "Tail loss"
	TRAIT_PARASITE TraitType = "Parasite"
	TRAIT_CARNIVOROUS TraitType = "Carnivorous"
	TRAIT_BURNED_FAT TraitType = "Burned fat"
	TRAIT_ADDITIONAL_GET_FOOD_FROM_BANK TraitType = "Additional get food from bank"
	TRAIT_HIBERNATION TraitType = "Hibernation"
	TRAIT_POISONOUS TraitType = "Poisonous"
	TRAIT_COMMUNICATION TraitType = "Communication"
	TRAIT_SCAVENGER TraitType = "Scavenger"
	TRAIT_SWIMMING TraitType = "Swimming"
	TRAIT_COOPERATION TraitType = "Cooperation"
	TRAIT_MIMICRY TraitType = "Mimicry"
	TRAIT_RUNNING TraitType = "Running"
)

type PhaseType string

const (
	PHASE_DEVELOPMENT PhaseType = "Development"
	PHASE_FOOD_BANK_DETERMINATION PhaseType = "Food bank determination"
	PHASE_FEEDING PhaseType = "Feeding"
	PHASE_EXTINCTION PhaseType = "Extinction"
	PHASE_FINAL PhaseType = "Final"
)

type ActionType string

const (
	ACTION_SEQUENCE ActionType = "Sequence"
	ACTION_SELECT ActionType = "Select"
	ACTION_START_TURN ActionType = "Start turn"
	ACTION_NEXT_PLAYER ActionType = "Next player"
	ACTION_ADD_CREATURE ActionType = "Add creature"
	ACTION_ADD_SINGLE_PROPERTY ActionType = "Add single property"
	ACTION_ADD_PAIR_PROPERTY ActionType = "Add pair property"
	ACTION_PASS ActionType = "Pass"
	ACTION_END_TURN ActionType = "End turn"
	ACTION_NEW_PHASE ActionType = "New phase"
	ACTION_ADD_TRAIT ActionType = "Add trait"
	ACTION_REMOVE_TRAIT ActionType = "Remove trait"
	ACTION_ADD_FILTERS ActionType = "Add filters"
	ACTION_ATTACK ActionType = "Attack"
	ACTION_DETERMINE_FOOD_BANK ActionType = "Determine food bank"
	ACTION_REMOVE_CREATURE ActionType = "Remove creature"
	ACTION_REMOVE_CARD ActionType = "Remove card"
	ACTION_REMOVE_PROPERTY ActionType = "Remove property"
	ACTION_GET_FOOD_FROM_BANK ActionType = "Get food from bank"
	ACTION_BURN_FAT ActionType = "Burn fat"
	ACTION_PIRACY ActionType = "Piracy"
	ACTION_DESTROY_BANK_FOOD ActionType = "Destroy bank food"
	ACTION_SELECT_FROM_AVAILABLE_ACTIONS ActionType = "Select from allowed actions"
	ACTION_EXTINCT ActionType = "Extinct"
	ACTION_TAKE_CARDS ActionType = "Take cards"
	ACTION_GAIN_FOOD ActionType = "Gain food"
	ACTION_GAIN_ADDITIONAL_FOOD ActionType = "Gain additional food"
	ACTION_EAT ActionType = "Eat"
	ACTION_HIBERNATE ActionType = "Hibernate"
	ACTION_RANDOM_ATTACK ActionType = "Random attack"
	ACTION_SCAVENGE ActionType = "Scavenge"
)

type ArgumentName string

const (
	PARAMETER_PROPERTY ArgumentName = "Property"
	PARAMETER_PHASE ArgumentName = "Phase"
	PARAMETER_PLAYER ArgumentName = "Player"
	PARAMETER_PAIR ArgumentName = "Pair"
	PARAMETER_CARD ArgumentName = "Card"
	PARAMETER_ACTIONS ArgumentName = "Actions"
	PARAMETER_CREATURE ArgumentName = "Creature"
	PARAMETER_TRAIT ArgumentName = "Trait"
	PARAMETER_SOURCE ArgumentName = "Source"
	PARAMETER_FILTERS ArgumentName = "Filters"
	PARAMETER_SOURCE_CREATURE ArgumentName = "Source creature"
	PARAMETER_TARGET_CREATURE ArgumentName = "Taget creature"
	PARAMETER_COUNT ArgumentName = "Count"
)

type AccessorMode int

const (
	ACCESSOR_MODE_ONE_OF_CREATURE_PROPERTIES AccessorMode = iota
	ACCESSOR_MODE_CREATURE_OWNER
	ACCESSOR_MODE_PROPERTY_OWNER
	ACCESSOR_MODE_CREATURES
)

type Type int

const (
	TYPE_CREATURE Type = iota
	TYPE_PROPERTY
)