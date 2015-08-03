// Evolution
package main

import (
	"container/list"
	"container/ring"
	"fmt"
	"math/rand"
	"time"
)

type Game struct {
	Players       *ring.Ring
	PlayersCount	int
	Deck          []*Card
	Filters       []Filter
	CurrentPhase  PhaseType
	CurrentPlayer *Player
	Food          int  
}

type WithTraits interface {
	GetTraits() []TraitType
	AddTrait(trait TraitType)
	RemoveTrait(trait TraitType)
}

type Source interface{}

type Card struct {
	ActiveProperty *Property
	Properties     []*Property
	Owners         []Source
}

func (c *Card) GetTraits() []TraitType {
	return c.ActiveProperty.Traits
}

func (c *Card) GoString() string {
	propertiesCount := len(c.Properties)
	if propertiesCount == 0 {
		return "()"
	}
	result := "(" + c.Properties[0].GoString()
	for i := 1; i < propertiesCount; i++ {
		result += "/" + c.Properties[i].GoString()
	}
	result += ")"
	return result
}

type Property struct {
	ContainingCard *Card
	Traits         []TraitType
}

func (p *Property) equals(property *Property) bool {
	if len(p.Traits) != len(property.Traits) {
		return false
	}
	equals := false
	for _,firstTrait := range p.Traits {
		equals = false
		for _,secondTrait := range property.Traits {
			if firstTrait == secondTrait {
				equals = true
				break
			}
		}
		if !equals {
			return false
		}
	}
	return true
}

func (c *Property) AddTrait(trait TraitType) {
	c.Traits = append(c.Traits, trait)
}

func (c *Property) RemoveTrait(trait TraitType) {
	for i, t := range c.Traits {
		if t == trait {
			c.Traits = append(c.Traits[:i], c.Traits[i+1:]...)
			return
		}
	}
}

func (c *Property) GetTraits() []TraitType {
	return c.Traits
}

func (c *Property) GoString() string {
	len := len(c.Traits)
	if len == 0 {
		return "()"
	}
	result := "(" + c.Traits[0].GoString()
	for i := 1; i<len;i++ {
		result += "/" + c.Traits[i].GoString()
	}
	result += ")"
	return result
}

type Creature struct {
	Head   *Card
	Tail   []*Card
	Owner  *Player
	Traits []TraitType
}

func (c *Creature) GoString() string {
	traits := c.GetTraits()
	len := len(traits)
	if len == 0 {
		return "(Creature)"
	}
	result := "(Creature : " + traits[0].GoString()
	for i := 1; i<len;i++ {
		result += "/" + traits[i].GoString()
	}
	result += "))"
	return result
}

func (c *Creature) GetTraits() []TraitType {
	result := make([]TraitType, 0, len(c.Traits))
	for _,trait := range c.Traits {
		result = append(result, trait)
	}
	for _,card := range c.Tail {
		result = append(result, card.GetTraits()...)
	}
	return result
}

func (c *Creature) AddTrait(trait TraitType) {
	c.Traits = append(c.Traits, trait)
}

func (c *Creature) RemoveTrait(trait TraitType) {
	for i := range c.Traits {
		if c.Traits[i] == trait {
			c.Traits = append(c.Traits[:i], c.Traits[i+1:]...)
			return
		}
	}
}

func (c *Creature) RemoveCard(card *Card) {
	for i := range c.Tail {
		if c.Tail[i] == card {
			c.Tail = append(c.Tail[:i], c.Tail[i+1:]...)
			return
		}
	}
}

type Player struct {
	ChoiceMaker
	Name      string
	Creatures []*Creature
	Cards     []*Card
	Traits    []TraitType
}

func (p *Player) RemoveCard(card *Card) {
	for i, c := range p.Cards {
		if c == card {
			p.Cards = append(p.Cards[:i], p.Cards[i+1:]...)
			return
		}
	}
}

func (p *Player) GetTraits() []TraitType {
	return p.Traits
}

func (p *Player) AddTrait(trait TraitType) {
	p.Traits = append(p.Traits, trait)
}

func (p *Player) RemoveTrait(trait TraitType) {
	for i := range p.Traits {
		if p.Traits[i] == trait {
			p.Traits = append(p.Traits[:i], p.Traits[i+1:]...)
			return	
		}
	}
}

func (p *Player) RemoveCreature(creature *Creature) {
	for i := range p.Creatures {
		if p.Creatures[i] == creature {
			p.Creatures = append(p.Creatures[:i], p.Creatures[i+1:]...)
			return	
		}
	}
}

func NewGame(players ...string) *Game {
	fmt.Println("Here is library start!")
	game := new(Game)
	game.InitializeDeck()
	game.InitializePlayers(players...)
	game.InitializeFilters()
	game.ExecuteAction(NewActionNewPhase(PHASE_DEVELOPMENT))
	return game
}

func (g *Game) TakeCards(player *Player, count int) {
	for i := 0; i < count; i++ {
		g.TakeCard(player)
	}
}

func (g *Game) TakeCard(player *Player) {
	deckLen := len(g.Deck)
	player.Cards = append(player.Cards, g.Deck[deckLen-1])
	player.Cards[len(player.Cards)-1].Owners = []Source {player}
	g.Deck = g.Deck[:deckLen-1]
}

func (g *Game) InitializeDeck() {
	camouflage := &Property{Traits : []TraitType {TRAIT_CAMOUFLAGE}}
	burrowing := &Property{Traits : []TraitType {TRAIT_BURROWING}}
	sharpVision := &Property{Traits : []TraitType {TRAIT_SHART_VISION}}
	symbiosys := &Property{Traits : []TraitType {TRAIT_PAIR, TRAIT_SIMBIOSYS}}
	piracy := &Property{Traits : []TraitType {TRAIT_PIRACY}}
	grazing := &Property{Traits : []TraitType {TRAIT_GRAZING}}
	tailLoss := &Property{Traits : []TraitType {TRAIT_TAIL_LOSS}}
	/*hibernation := &Property{Name: "hibernation"}
	poisonous := &Property{Name: "poisonous"}
	communication := &Property{Name: "communication"}
	scavenger := &Property{Name: "scavenger"}
	running := &Property{Name: "running"}
	mimicry := &Property{Name: "mimicry"}
	swimming := &Property{Name: "swimming"}*/
	parasite := &Property{Traits : []TraitType {TRAIT_PARASITE, TRAIT_REQUIRE_FOOD, TRAIT_REQUIRE_FOOD}}
	carnivorous := &Property{Traits : []TraitType {TRAIT_CARNIVOROUS, TRAIT_REQUIRE_FOOD}}
	fatTissue := &Property{Traits : []TraitType {TRAIT_FAT_TISSUE}}
	//cooperation := &Property{Name: "cooperation"}
	highBodyWeight := &Property{Traits : []TraitType {TRAIT_HIGH_BODY_WEIGHT, TRAIT_REQUIRE_FOOD}}
	g.Deck = make([]*Card, 0, 84)
	g.AddCard(4, camouflage)
	g.AddCard(4, burrowing)
	g.AddCard(4, sharpVision)
	g.AddCard(4, symbiosys)
	g.AddCard(4, piracy)
	g.AddCard(4, grazing)
	g.AddCard(4, tailLoss)
	/*g.AddCard(4, hibernation)
	g.AddCard(4, poisonous)
	g.AddCard(4, communication)
	g.AddCard(4, scavenger)
	g.AddCard(4, running)
	g.AddCard(4, mimicry)
	g.AddCard(8, swimming)
	*/
	g.AddCard(4, parasite, carnivorous)
	/*g.AddCard(4, parasite, fatTissue)
	g.AddCard(4, cooperation, carnivorous)
	g.AddCard(4, cooperation, fatTissue)
	g.AddCard(4, highBodyWeight, carnivorous)*/
	g.AddCard(4, highBodyWeight, fatTissue)
	g.ShuffleDeck()
}

func (g *Game) InitializePlayers(names ...string) {
	g.Players = ring.New(len(names))
	g.PlayersCount = len(names)
	for _, name := range names {
		player := &Player{Name: name, ChoiceMaker: ConsoleChoiceMaker{}}
		g.Players.Value = player
		g.TakeCards(player, 12)
		g.Players = g.Players.Next()
	}
	g.CurrentPlayer = g.Players.Value.(*Player)
}

func (g *Game) InitializeFilters() {
	//Remove all pass trait on phase start
	g.Filters = append(g.Filters, 
		&FilterAction{
			FILTER_ACTION_EXECUTE_BEFORE, 
			&ConditionActionType{ACTION_NEW_PHASE}, 
			nil,
			NewActionRemoveTrait(FILTER_SOURCE_PARAMETER_ALL_PLAYERS, TRAIT_PASS)})
	
	//Start player turn on phase start
	g.Filters = append(g.Filters, 
		&FilterAction{
			FILTER_ACTION_EXECUTE_AFTER, 
			NewANDCondition(
				NewORCondition(&ConditionPhase{PHASE_DEVELOPMENT}, &ConditionPhase{PHASE_FEEDING}),
				&ConditionActionType{ACTION_NEW_PHASE}), 
			nil,
			NewActionStartTurn(FILTER_SOURCE_PARAMETER_CURRENT_PLAYER)})
			
	//Start selecting after turn start
	g.Filters = append(g.Filters, 
		&FilterAction{
			FILTER_ACTION_EXECUTE_AFTER, 
			NewANDCondition(
				NewORCondition(&ConditionPhase{PHASE_DEVELOPMENT}, &ConditionPhase{PHASE_FEEDING}),
				&ConditionActionType{ACTION_START_TURN}), 
			nil,
			NewActionSelectFromAvailableActions()})
			
	//Alow pass turn to next player in feeding mode
	g.Filters = append(g.Filters,
		NewFilterAllow(
			&ConditionPhase{PHASE_FEEDING},
			nil,
			NewActionAddFilters(&FilterAction{
					FILTER_ACTION_REPLACE,
					&ConditionActionType{ACTION_SELECT_FROM_AVAILABLE_ACTIONS},
					&ConditionActionType{ACTION_NEXT_PLAYER},
					NewActionNextPlayer(g)})))
			
	//In feeding phase player make turns, until pass
	g.Filters = append(g.Filters,
		&FilterAction{
			FILTER_ACTION_EXECUTE_AFTER,
			NewANDCondition(&ConditionPhase{PHASE_FEEDING},&ConditionActionType{ACTION_SELECT_FROM_AVAILABLE_ACTIONS}),
			nil,
			NewActionSelectFromAvailableActions()})
			
	//In development phase player pass turn to next player
	g.Filters = append(g.Filters,
	&FilterAction{
		FILTER_ACTION_EXECUTE_AFTER,
		NewANDCondition(&ConditionPhase{PHASE_DEVELOPMENT}, &ConditionActionType{ACTION_START_TURN}),
		nil,
		NewActionNextPlayer(g)})
			
	//Allow adding creatures in develompent phase
	g.Filters = append(g.Filters, NewFilterAllow(&ConditionPhase{PHASE_DEVELOPMENT}, nil, NewActionAddCreature(FILTER_SOURCE_PARAMETER_CURRENT_PLAYER, FILTER_SOURCE_PARAMETER_ONE_OF_CURRENT_PLAYER_CARDS)))
	//Allow adding pair properties in development phase
	g.Filters = append(g.Filters, NewFilterAllow(&ConditionPhase{PHASE_DEVELOPMENT}, nil, NewActionAddPairProperty(FILTER_SOURCE_PARAMETER_ONE_OF_CURRENT_PLAYER_CREATURES_PAIR, FILTER_SOURCE_PARAMETER_ONE_OF_CURRENT_PLAYER_CARDS_PROPERTIES)))
	//Allow adding single properties in development phase
	g.Filters = append(g.Filters, NewFilterAllow(&ConditionPhase{PHASE_DEVELOPMENT}, nil, NewActionAddSingleProperty(FILTER_SOURCE_PARAMETER_ONE_OF_CURRENT_PLAYER_CREATURES, FILTER_SOURCE_PARAMETER_ONE_OF_CURRENT_PLAYER_CARDS_PROPERTIES)))
	//Deny adding single properties if
	g.Filters = append(g.Filters, 
		&FilterDeny{
			NewANDCondition(
				&ConditionActionType{ACTION_ADD_SINGLE_PROPERTY},
				NewORCondition(
					NewANDCondition(
						NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_FAT_TISSUE}, 0),
						&ConditionContains{FILTER_SOURCE_PARAMETER_CREATURE_PROPERTIES, FILTER_SOURCE_PARAMETER_PROPERTY}),
					NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_PAIR}, 1))),
			nil,
			})
	//Deny adding pair properties is
	g.Filters = append(g.Filters,
		&FilterDeny{
			NewANDCondition(
				&ConditionActionType{ACTION_ADD_PAIR_PROPERTY},
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_PAIR}, 0)),
			nil,
		})
		
	//Allow pass in development and feeding phase
	g.Filters = append(g.Filters, NewFilterAllow(NewORCondition(&ConditionPhase{PHASE_DEVELOPMENT}, &ConditionPhase{PHASE_FEEDING}), nil, &Action{ACTION_PASS, map[ArgumentName]Source {}}))
	
	//If all players pass in development phase, start food bank determination
	g.Filters = append(g.Filters, 
		&FilterAction{
			FILTER_ACTION_REPLACE, 
			NewANDCondition(
				&ConditionPhase{PHASE_DEVELOPMENT},
				&ConditionActionType{ACTION_NEXT_PLAYER},
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_ALL_PLAYERS, TRAIT_PASS}, 1)), 
			nil,
			NewActionNewPhase(PHASE_FOOD_BANK_DETERMINATION)})
	
	
	//If player pass - replace his turn with NextTurn
	g.Filters = append(g.Filters, 
		&FilterAction{
			FILTER_ACTION_REPLACE, 
			NewANDCondition(
				&ConditionActionType{ACTION_START_TURN}, 
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_PLAYER, TRAIT_PASS}, 1)), 
			nil,
			NewActionNextPlayer(g)})
	
	//Determine food bank
	g.Filters = append(g.Filters,
		&FilterAction{
			FILTER_ACTION_EXECUTE_AFTER,
			NewANDCondition(
				&ConditionActionType{ACTION_NEW_PHASE},
				NewConditionEqual(FILTER_SOURCE_PARAMETER_PHASE, PHASE_FOOD_BANK_DETERMINATION)),
			nil,
			&Action{ACTION_DETERMINE_FOOD_BANK, map[ArgumentName]Source {}}})
		
	//After food bank determination, start feeding phase
	g.Filters = append(g.Filters,
		&FilterAction{
				FILTER_ACTION_EXECUTE_AFTER,
				&ConditionActionType{ACTION_DETERMINE_FOOD_BANK},
				nil,
				NewActionNewPhase(PHASE_FEEDING)})
		
	//Allow get food from bank for creatures
	g.Filters = append(g.Filters,
		NewFilterAllow(
			&ConditionPhase{PHASE_FEEDING},
			nil,
			NewActionGetFoodFromBank(FILTER_SOURCE_PARAMETER_ONE_OF_CURRENT_PLAYER_CREATURES)))
		
	//Deny get food from bank
	g.Filters = append(g.Filters,
		&FilterDeny{
			NewANDCondition(
				&ConditionActionType{ACTION_GET_FOOD_FROM_BANK}, 
				NewORCondition(
					&ConditionActionDenied{NewActionAddTrait(FILTER_SOURCE_PARAMETER_CREATURE, TRAIT_FOOD)},
					NewConditionEqual(FILTER_SOURCE_PARAMETER_FOOD_BANK_COUNT, 0))),
			nil})
	
	//Deny food get if creature already full
	g.Filters = append(g.Filters,
		&FilterDeny{
			NewANDCondition(
				&ConditionActionType{ACTION_ADD_TRAIT},
				NewConditionEqual(FILTER_SOURCE_PARAMETER_TRAIT, FILTER_SOURCE_PARAMETER_ANY_FOOD),
				NewConditionEqual(
					TraitsCount{FILTER_SOURCE_PARAMETER_SOURCE, FILTER_SOURCE_PARAMETER_ALL_FOOD_AND_FAT},
					TraitsCount{FILTER_SOURCE_PARAMETER_SOURCE, FILTER_SOURCE_PARAMETER_FOOD_AND_FAT_LIMIT})),
			nil})
	
	//Replace food get with fat get
	g.Filters = append(g.Filters,
		&FilterAction{
			FILTER_ACTION_REPLACE,
			NewANDCondition(
				&ConditionActionType{ACTION_ADD_TRAIT},
				NewConditionEqual(FILTER_SOURCE_PARAMETER_TRAIT, FILTER_SOURCE_PARAMETER_ANY_FOOD),
				NewConditionEqual(
					TraitsCount{FILTER_SOURCE_PARAMETER_SOURCE, FILTER_SOURCE_PARAMETER_ALL_FOOD},
					TraitsCount{FILTER_SOURCE_PARAMETER_SOURCE, TRAIT_REQUIRE_FOOD})),
			nil,
			NewActionAddTrait(FILTER_SOURCE_PARAMETER_SOURCE, TRAIT_FAT)})
	
	//Set fed trait
	g.Filters = append(g.Filters,
		&FilterAction{
			FILTER_ACTION_EXECUTE_AFTER,
			NewANDCondition(
				&ConditionActionType{ACTION_ADD_TRAIT},
				NewConditionEqual(FILTER_SOURCE_PARAMETER_TRAIT, FILTER_SOURCE_PARAMETER_ANY_FOOD),
				NewConditionEqual(
					TraitsCount{FILTER_SOURCE_PARAMETER_SOURCE,FILTER_SOURCE_PARAMETER_ALL_FOOD},
					TraitsCount{FILTER_SOURCE_PARAMETER_SOURCE,TRAIT_REQUIRE_FOOD})),
			nil,
			NewActionAddTrait(FILTER_SOURCE_PARAMETER_SOURCE, TRAIT_FED)})

	//Deny remove food from creature, if it have none
	g.Filters = append(g.Filters,
		&FilterDeny{
			NewANDCondition(
				&ConditionPhase{PHASE_FEEDING},
				&ConditionActionType{ACTION_REMOVE_TRAIT},
				NewConditionEqual(FILTER_SOURCE_PARAMETER_TRAIT, TRAIT_FOOD),
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_SOURCE , TRAIT_FOOD}, 0)),
			nil})
	
	g.Filters = append(g.Filters,
		&FilterDeny{
			NewANDCondition(
				&ConditionPhase{PHASE_FEEDING},
				&ConditionActionType{ACTION_REMOVE_TRAIT},
				NewConditionEqual(FILTER_SOURCE_PARAMETER_TRAIT, TRAIT_ADDITIONAL_FOOD),
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_SOURCE , TRAIT_ADDITIONAL_FOOD}, 0)),
			nil})	
		
	//camouflage
	g.Filters = append(g.Filters,
		&FilterDeny{
			NewANDCondition(
				&ConditionActionType{ACTION_ATTACK},
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_SOURCE_CREATURE, TRAIT_CAMOUFLAGE}, 1),
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_TARGET_CREATURE, TRAIT_SHART_VISION}, 0)),
			nil})
	//burrowing
	g.Filters = append(g.Filters,
		&FilterDeny{
			NewANDCondition(
				&ConditionActionType{ACTION_ATTACK},
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_SOURCE_CREATURE, TRAIT_BURROWING}, 1),
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_SOURCE_CREATURE, TRAIT_FED}, 1)),
				nil})
	//cymbiosys
	g.Filters = append(g.Filters,
		&FilterAction{
			FILTER_ACTION_EXECUTE_AFTER,
			NewANDCondition(
				&ConditionActionType{ACTION_ADD_PAIR_PROPERTY},
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_SIMBIOSYS}, 1)),
			nil,
			NewActionAddFilters(
				&FilterDeny{
					NewANDCondition(
						&ConditionActionType{ACTION_ATTACK},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_TARGET_CREATURE}, FILTER_SOURCE_PARAMETER_RIGHT_CREATURE)),
					NewANDCondition(
						&ConditionActionType{ACTION_REMOVE_PROPERTY},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY))},
				&FilterDeny{
					NewANDCondition(
						&ConditionActionType{ACTION_ADD_TRAIT},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_TRAIT}, FILTER_SOURCE_PARAMETER_ANY_FOOD),
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_SOURCE}, FILTER_SOURCE_PARAMETER_RIGHT_CREATURE),
						NewConditionEqual(InstantiationOff{TraitsCount{InstantiationOn{FILTER_SOURCE_PARAMETER_LEFT_CREATURE}, TRAIT_FED}}, 0)),
					NewANDCondition(
						&ConditionActionType{ACTION_REMOVE_PROPERTY},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY))},
				&FilterDeny{
					NewANDCondition(
						&ConditionActionType{ACTION_ADD_PAIR_PROPERTY},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY),
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PAIR}, FILTER_SOURCE_PARAMETER_PAIR)),
					NewANDCondition(
						&ConditionActionType{ACTION_REMOVE_PROPERTY},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY))})})
	
	//piracy
	g.Filters = append(g.Filters,
		&FilterAction {
			FILTER_ACTION_EXECUTE_AFTER,
			NewANDCondition(
				&ConditionActionType{ACTION_ADD_SINGLE_PROPERTY},
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_PIRACY}, 1)),
			nil,
			NewActionAddFilters(
				NewFilterAllow(
					NewANDCondition(
						&ConditionPhase{PHASE_FEEDING},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_CURRENT_PLAYER}, FILTER_SOURCE_PARAMETER_CURRENT_PLAYER)),
					NewANDCondition(
						&ConditionActionType{ACTION_REMOVE_PROPERTY},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY)),
					NewActionPiracy(FILTER_SOURCE_PARAMETER_CREATURE, InstantiationOff{FILTER_SOURCE_PARAMETER_ONE_OF_CREATURES}, FILTER_SOURCE_PARAMETER_ANY_FOOD)),
				&FilterAction{
					FILTER_ACTION_EXECUTE_AFTER,
					&ConditionActionType{ACTION_PIRACY},
					NewANDCondition(
						&ConditionActionType{ACTION_REMOVE_PROPERTY},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY)),
					NewActionAddTrait(FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_USED)},
				&FilterAction{
					FILTER_ACTION_EXECUTE_BEFORE,
					NewANDCondition(
						&ConditionActionType{ACTION_START_TURN},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PLAYER}, FILTER_SOURCE_PARAMETER_PLAYER)),
					NewANDCondition(
						&ConditionActionType{ACTION_REMOVE_PROPERTY},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY)),
					NewActionRemoveTrait(FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_USED)},
				&FilterDeny{
					NewANDCondition(
						&ConditionActionType{ACTION_PIRACY},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_SOURCE_CREATURE}, FILTER_SOURCE_PARAMETER_SOURCE_CREATURE),
						NewConditionEqual(InstantiationOff{TraitsCount{FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_USED}}, 1)),
					NewANDCondition(
						&ConditionActionType{ACTION_REMOVE_PROPERTY},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY))})})
	g.Filters = append(g.Filters,
		&FilterDeny{
			NewANDCondition(
				&ConditionActionType{ACTION_PIRACY},
				NewORCondition(
					NewConditionEqual(FILTER_SOURCE_PARAMETER_SOURCE_CREATURE, FILTER_SOURCE_PARAMETER_TARGET_CREATURE),
					NewANDCondition(
						NewConditionEqual(FILTER_SOURCE_PARAMETER_TRAIT, TRAIT_FOOD),
						&ConditionActionDenied{NewActionRemoveTrait(FILTER_SOURCE_PARAMETER_TARGET_CREATURE, TRAIT_FOOD)}),
					NewANDCondition(
						NewConditionEqual(FILTER_SOURCE_PARAMETER_TRAIT, TRAIT_ADDITIONAL_FOOD),
						&ConditionActionDenied{NewActionRemoveTrait(FILTER_SOURCE_PARAMETER_TARGET_CREATURE, TRAIT_ADDITIONAL_FOOD)}),
					&ConditionActionDenied{NewActionAddTrait(FILTER_SOURCE_PARAMETER_SOURCE_CREATURE, TRAIT_ADDITIONAL_FOOD)})),	
			nil})
	//grazing
	g.Filters = append(g.Filters,
		&FilterAction {
			FILTER_ACTION_EXECUTE_AFTER,
			NewANDCondition(
				&ConditionActionType{ACTION_ADD_SINGLE_PROPERTY},
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_GRAZING}, 1)),
			nil,
			NewActionAddFilters(
				NewFilterAllow(
					NewANDCondition(
						&ConditionPhase{PHASE_FEEDING},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_CURRENT_PLAYER}, FILTER_SOURCE_PARAMETER_CURRENT_PLAYER)),
					NewANDCondition(
						&ConditionActionType{ACTION_REMOVE_PROPERTY},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY)),
					NewActionGrazing(FILTER_SOURCE_PARAMETER_PROPERTY)),
				&FilterAction{
					FILTER_ACTION_EXECUTE_AFTER,
					&ConditionActionType{ACTION_DESTROY_BANK_FOOD},
					NewANDCondition(
						&ConditionActionType{ACTION_REMOVE_PROPERTY},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY)),
					NewActionAddTrait(FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_USED)},
				&FilterAction{
					FILTER_ACTION_EXECUTE_BEFORE,
					NewANDCondition(
						&ConditionActionType{ACTION_START_TURN},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_CURRENT_PLAYER}, FILTER_SOURCE_PARAMETER_CURRENT_PLAYER)),
					NewANDCondition(
						&ConditionActionType{ACTION_REMOVE_PROPERTY},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY)),
					NewActionRemoveTrait(FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_USED)},
				&FilterDeny{
					NewANDCondition(
						&ConditionActionType{ACTION_DESTROY_BANK_FOOD},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_SOURCE_CREATURE}, FILTER_SOURCE_PARAMETER_SOURCE_CREATURE),
						NewConditionEqual(InstantiationOff{TraitsCount{FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_USED}}, 1)),
					NewANDCondition(
						&ConditionActionType{ACTION_REMOVE_PROPERTY},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY))})})
	g.Filters = append(g.Filters,
		&FilterDeny{
			NewANDCondition(
				&ConditionActionType{ACTION_DESTROY_BANK_FOOD},
				NewConditionEqual(FILTER_SOURCE_PARAMETER_FOOD_BANK_COUNT, 0)),	
			nil})
			
	//High body wieght
	g.Filters = append(g.Filters,
		&FilterDeny{
			NewANDCondition(
				&ConditionActionType{ACTION_ATTACK},
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_SOURCE_CREATURE, TRAIT_HIGH_BODY_WEIGHT}, 1),
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_TARGET_CREATURE, TRAIT_HIGH_BODY_WEIGHT}, 0)),
			nil})
			
	//Tail loss
	g.Filters = append(g.Filters,
		&FilterAction{
			FILTER_ACTION_REPLACE,
			NewANDCondition(
				&ConditionActionType{ACTION_ATTACK},
				NewConditionEqual(TypeOf{FILTER_SOURCE_PARAMETER_TARGET_CREATURE},TYPE_CREATURE),
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_TARGET_CREATURE, TRAIT_TAIL_LOSS}, 1)),
			nil,
			NewActionAttack(
				Accessor{FILTER_SOURCE_PARAMETER_TARGET_CREATURE, ACCESSOR_MODE_CREATURE_OWNER}, 
				FILTER_SOURCE_PARAMETER_SOURCE_CREATURE, 
				Accessor{FILTER_SOURCE_PARAMETER_TARGET_CREATURE, ACCESSOR_MODE_ONE_OF_CREATURE_PROPERTIES})})
	
	//Carnivorous
	g.Filters = append(g.Filters,
		&FilterAction {
			FILTER_ACTION_EXECUTE_AFTER,
			NewANDCondition(
				&ConditionActionType{ACTION_ADD_SINGLE_PROPERTY},
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_CARNIVOROUS}, 1)),
			nil,
			NewActionAddFilters(
				NewFilterAllow(
					NewANDCondition(
						&ConditionPhase{PHASE_FEEDING},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_CURRENT_PLAYER}, FILTER_SOURCE_PARAMETER_CURRENT_PLAYER)),
					NewANDCondition(
						&ConditionActionType{ACTION_REMOVE_PROPERTY},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY)),
					NewActionAttack(
						FILTER_SOURCE_PARAMETER_CURRENT_PLAYER, 
						FILTER_SOURCE_PARAMETER_CREATURE, 
						InstantiationOff{FILTER_SOURCE_PARAMETER_ONE_OF_CREATURES})),
				&FilterAction{
					FILTER_ACTION_EXECUTE_AFTER,
					&ConditionActionType{ACTION_PIRACY},
					NewANDCondition(
						&ConditionActionType{ACTION_REMOVE_PROPERTY},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY)),
					NewActionAddTrait(FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_USED)},
				&FilterAction{
					FILTER_ACTION_EXECUTE_BEFORE,
					NewANDCondition(
						&ConditionActionType{ACTION_START_TURN},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PLAYER}, FILTER_SOURCE_PARAMETER_PLAYER)),
					NewANDCondition(
						&ConditionActionType{ACTION_REMOVE_PROPERTY},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY)),
					NewActionRemoveTrait(FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_USED)},
				&FilterDeny{
					NewANDCondition(
						&ConditionActionType{ACTION_ATTACK},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_SOURCE_CREATURE}, FILTER_SOURCE_PARAMETER_SOURCE_CREATURE),
						NewConditionEqual(InstantiationOff{TraitsCount{FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_USED}}, 1)),
					NewANDCondition(
						&ConditionActionType{ACTION_REMOVE_PROPERTY},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY))})})
	g.Filters = append(g.Filters,
		&FilterDeny{
			NewANDCondition(
				&ConditionActionType{ACTION_ATTACK},
				NewORCondition(
					NewConditionEqual(FILTER_SOURCE_PARAMETER_SOURCE_CREATURE, FILTER_SOURCE_PARAMETER_TARGET_CREATURE),
					&ConditionActionDenied{NewActionAddTrait(FILTER_SOURCE_PARAMETER_SOURCE_CREATURE, TRAIT_ADDITIONAL_FOOD)})),	
			nil})
}

func (g *Game) AddCard(count int, properties ...*Property) {
	for i := 0; i < count; i++ {
		card := g.NewCard(properties...)
		g.Deck = append(g.Deck, card)
	}
}

func (g *Game) NewCard(properties ...*Property) *Card {
	if len(properties) == 0 {
		return &Card{}
	}
	card := new(Card)
	for _,property := range properties {
		card.Properties = append(card.Properties, &Property {Traits : property.Traits})
	}
	card.ActiveProperty = card.Properties[0]
	for i := range card.Properties {
		card.Properties[i].ContainingCard = card
	}
	return card
}

func (g *Game) ShuffleDeck() {
	rand.Seed(3)
	for i := range g.Deck {
		j := rand.Intn(i + 1)
		g.Deck[i], g.Deck[j] = g.Deck[j], g.Deck[i]
	}
}

func (g *Game) ActionDenied(action *Action) (result bool) {
	for _, filter := range g.Filters {
		if filter.GetType() == FILTER_DENY {
			if filter.CheckCondition(g, action) {
				fmt.Printf("%#v denied because %#v\n", action, filter.GetCondition())
				return true
			}
		}
	}
	return false
}

func (g *Game) GetAlowedActions() []*Action {
	var result []*Action
	for _, filter := range g.Filters {
		if filter.GetType() == FILTER_ALLOW && filter.CheckCondition(g, nil) {
			actions := filter.(*FilterAllow).GetActions(g)
			for _, action := range actions {
				if action.Type == ACTION_SELECT {
					result = append(result, g.ExpandActionSelect(action)...)
					//fmt.Printf("%#v\n",result)
				} else {
					if !g.ActionDenied(action) {
						result = append(result, action)
					}
					//fmt.Printf("%#v\n",result)
				}
			}
		}
	}
	return result
}

func (g *Game) ExpandActionSelect(action *Action) []*Action {
	result := make([]*Action, 0, 4)
	for _,a := range action.Arguments[PARAMETER_ACTIONS].([]*Action) {
		if a.Type == ACTION_SELECT {
			result = append(result, g.ExpandActionSelect(a)...)
		} else {
			if !g.ActionDenied(a) {
				result = append(result, a)
			}
		}
	}
	return result
}

func (g *Game) ExecuteAction(rawAction *Action) {
	stack := list.New()
	stack.PushFront(rawAction)
	for stackFront := stack.Front(); stackFront != nil ; stackFront = stack.Front() {
		/*fmt.Println("Stack trace:")
		i := 0
		for a := stack.Front(); a != nil; a = a.Next() {
			fmt.Printf("%v) %#v\n", i, a.Value)
			i++
		}*/
		stack.Remove(stackFront)
		action := stackFront.Value.(*Action)
		if action.Type == ACTION_SELECT {
			action = g.CurrentPlayer.ChoiceMaker.MakeChoice(g, action.Arguments[PARAMETER_ACTIONS].([]*Action))
		}
		replaced := false
		for _, filter := range g.Filters {
			if filter.GetType() == FILTER_ACTION_REPLACE && filter.CheckCondition(g, action) {
				stack.PushFront(filter.(*FilterAction).GetAction().InstantiateFilterPrototypeAction(g, action, true))
				fmt.Printf("Replaced %#v with %#v because %#v\n", action, filter.(*FilterAction).GetAction().InstantiateFilterPrototypeAction(g, action, true), filter.GetCondition())
				replaced = true
				break
			}
			if filter.GetType() == FILTER_ACTION_EXECUTE_BEFORE && filter.CheckCondition(g, action) {
				g.ExecuteAction(filter.InstantiateFilterPrototype(g, action, true).(*FilterAction).GetAction())
			} 
		}
		removed := true
		for removed {
			removed = false
			for i, filter := range g.Filters {
				if filter.CheckRemoveCondition(g, action) {
					fmt.Printf("Removing filter %#v because %#v\n", filter, filter.GetCondition())
					g.Filters = append(g.Filters[:i], g.Filters[i+1:]...)
					removed = true
					break
				}
			}
		}
		if replaced {
			continue
		}
		fmt.Printf("Executing action: %#v\n", action)
		action.Execute(g)
		for _, filter := range g.Filters {
			if filter.GetType() == FILTER_ACTION_EXECUTE_AFTER && filter.CheckCondition(g, action) {
				stack.PushBack(filter.InstantiateFilterPrototype(g, action, true).(*FilterAction).GetAction())
			}
		}
		time.Sleep(250 * time.Millisecond)
	}
}
