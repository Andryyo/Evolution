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
	Actions       list.List
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
	Properties     []Property
	Owners          []*Player
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

func (p Property) equals(property Property) bool {
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

func (c Property) AddTrait(trait TraitType) {
	c.Traits = append(c.Traits, trait)
}

func (c Property) RemoveTrait(trait TraitType) {
	for i, t := range c.Traits {
		if t == trait {
			c.Traits = append(c.Traits[:i], c.Traits[i+1:]...)
			return
		}
	}
}

func (c Property) GetTraits() []TraitType {
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

func (g *Game) GetInstantiationVariants(arguments map[ArgumentName]Source, argumentsNames []ArgumentName, argumentNumber int) []map[ArgumentName]Source {
	argumentsNamesLen := len(argumentsNames)
	if argumentsNamesLen == 0 {
		return []map[ArgumentName]Source{}
	}
	argumentName := argumentsNames[argumentNumber]
	argument := arguments[argumentName]
	instantiatedArguments := g.InstantiateArgument(argument)
	if argumentNumber == argumentsNamesLen-1 {
		result := make([]map[ArgumentName]Source, 0, argumentsNamesLen)
		for _, instantiatedArgument := range instantiatedArguments {
			tmp := make(map[ArgumentName]Source)
			tmp[argumentName] = instantiatedArgument
			result = append(result, tmp)
		}
		return result
	}
	completedVariants := g.GetInstantiationVariants(arguments, argumentsNames, argumentNumber+1)
	result := make([]map[ArgumentName]Source, 0, len(completedVariants)*len(instantiatedArguments))
	for _, instantiatedArgument := range instantiatedArguments {
		for _, completedVariant := range completedVariants {
			tmp := make(map[ArgumentName]Source)
			for key := range completedVariant {
				tmp[key] = completedVariant[key]
			}
			tmp[argumentName] = instantiatedArgument
			result = append(result, tmp)
		}
	}
	return result

}

func (g *Game) InstantiateArgument(argument Source) []Source {
	result := make([]Source, 0, 8)
	if _, ok := argument.(SourcePrototype); !ok {
		return []Source{argument}
	}
	switch argument {
	case SOURCE_PROTOTYPE_PLAYER:
		result = append(result, g.CurrentPlayer)
	case SOURCE_PROTOTYPE_PLAYER_CARD:
		for _, card := range g.CurrentPlayer.Cards {
			result = append(result, card)
		}
	case SOURCE_PROTOTYPE_PLAYER_CREATURE:
		for _, creature := range g.CurrentPlayer.Creatures {
			result = append(result, creature)
		}
	case SOURCE_PROTOTYPE_PLAYER_CARD_PROPERTY:
		for _, card := range g.CurrentPlayer.Cards {
			for _, property := range card.Properties {
				result = append(result, property)
			}
		}
	case SOURCE_PROTOTYPE_CREATURES_PAIR:
		for _, firstCreature := range g.CurrentPlayer.Creatures {
			for _,secondCreature := range g.CurrentPlayer.Creatures {
				if firstCreature != secondCreature {
					result = append(result, []*Creature {firstCreature, secondCreature})
				}
			}
		}
	}
	return result
}

func (g *Game) InstantiateActionPrototype(prototype *Action) []*Action {
	var result []*Action
	var definedArgumentsNames []ArgumentName
	var undefinedArgumentsNames []ArgumentName
	for key, argument := range prototype.Arguments {
		if _, ok := argument.(SourcePrototype); ok {
			undefinedArgumentsNames = append(undefinedArgumentsNames, key)
		} else {
			definedArgumentsNames = append(definedArgumentsNames, key)
		}
	}
	if len(undefinedArgumentsNames) == 0 {
		if !g.ActionDenied(prototype) {
			return []*Action{prototype}
		} else {
			return []*Action{}
		}
	}
	var variants []map[ArgumentName]Source = g.GetInstantiationVariants(prototype.Arguments, undefinedArgumentsNames, 0)
	for _, variant := range variants {
		for _, definedArgumentName := range definedArgumentsNames {
			variant[definedArgumentName] = prototype.Arguments[definedArgumentName]
		}
		action := &Action{prototype.Type, variant}
		if !g.ActionDenied(action) {
			result = append(result, action)
		}
	}
	return result
}

func NewGame(players ...string) *Game {
	fmt.Println("Here is library start!")
	game := new(Game)
	game.InitializeDeck()
	game.InitializePlayers(players...)
	game.InitializeFilters()
	game.Actions.Init()
	game.Actions.PushBack(NewActionNewPhase(PHASE_DEVELOPMENT))

	for action := game.Actions.Front(); action != nil; action = game.Actions.Front() {
		fmt.Println("Stack trace:")
		i := 0
		for a := game.Actions.Front(); a != nil; a = a.Next() {
			fmt.Printf("%v) %#v\n", i, a.Value)
			i++
		}
		game.Actions.Remove(action)
		game.ExecuteAction(action.Value.(*Action))
		time.Sleep(250 * time.Millisecond)
	}
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
	player.Cards[len(player.Cards)-1].Owners = []*Player{player}
	g.Deck = g.Deck[:deckLen-1]
}

func (g *Game) InitializeDeck() {
	camouflage := Property{Traits : []TraitType {TRAIT_CAMOUFLAGE}}
	burrowing := Property{Traits : []TraitType {TRAIT_BURROWING}}
	sharpVision := Property{Traits : []TraitType {TRAIT_SHART_VISION}}
	symbiosys := Property{Traits : []TraitType {TRAIT_PAIR, TRAIT_SIMBIOSYS}}
	/*piracy := &Property{Name: "piracy"}
	grazing := &Property{Name: "grazing"}
	tailLoss := &Property{Name: "tailLoss"}
	hibernation := &Property{Name: "hibernation"}
	poisonous := &Property{Name: "poisonous"}
	communication := &Property{Name: "communication"}
	scavenger := &Property{Name: "scavenger"}
	running := &Property{Name: "running"}
	mimicry := &Property{Name: "mimicry"}
	swimming := &Property{Name: "swimming"}
	parasite := &Property{Name: "parasite"}
	carnivorous := &Property{Name: "carnivorous"}
	fatTissue := &Property{Name: "fatTissue"}
	cooperation := &Property{Name: "cooperation"}
	highBodyWeight := &Property{Name: "highBodyWeight"}*/
	g.Deck = make([]*Card, 0, 84)
	g.AddCard(4, camouflage)
	g.AddCard(4, burrowing)
	g.AddCard(4, sharpVision)
	g.AddCard(4, symbiosys)
	/*g.AddCard(4, piracy)
	g.AddCard(4, grazing)
	g.AddCard(4, tailLoss)
	g.AddCard(4, hibernation)
	g.AddCard(4, poisonous)
	g.AddCard(4, communication)
	g.AddCard(4, scavenger)
	g.AddCard(4, running)
	g.AddCard(4, mimicry)
	g.AddCard(8, swimming)
	g.AddCard(4, parasite, carnivorous)
	g.AddCard(4, parasite, fatTissue)
	g.AddCard(4, cooperation, carnivorous)
	g.AddCard(4, cooperation, fatTissue)
	g.AddCard(4, highBodyWeight, carnivorous)
	g.AddCard(4, highBodyWeight, fatTissue)*/
	g.ShuffleDeck()
}

func (g *Game) InitializePlayers(names ...string) {
	g.Players = ring.New(len(names))
	g.PlayersCount = len(names)
	for _, name := range names {
		player := &Player{Name: name, ChoiceMaker: ConsoleChoiceMaker{}}
		g.Players.Value = player
		g.TakeCards(player, 6)
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
			NewActionStartTurn(SOURCE_PROTOTYPE_PLAYER)})
			
	//Alow pass turn to next player in feeding mode
	g.Filters = append(g.Filters,
		NewFilterAllow(
			&ConditionPhase{PHASE_FEEDING},
			nil,
			NewActionAddFilters(&FilterAction{
					FILTER_ACTION_REPLACE,
					&ConditionActionType{ACTION_START_TURN},
					&ConditionActionType{ACTION_NEXT_PLAYER},
					NewActionNextPlayer(g)})))
			
	//In feeding phase player make turns, until pass
	g.Filters = append(g.Filters,
		&FilterAction{
			FILTER_ACTION_EXECUTE_AFTER,
			NewANDCondition(&ConditionPhase{PHASE_FEEDING},&ConditionActionType{ACTION_START_TURN}),
			nil,
			NewActionStartTurn(FILTER_SOURCE_PARAMETER_PLAYER)})
			
	//In development phase player pass turn to next player
	g.Filters = append(g.Filters,
	&FilterAction{
		FILTER_ACTION_EXECUTE_AFTER,
		NewANDCondition(&ConditionPhase{PHASE_DEVELOPMENT}, &ConditionActionType{ACTION_START_TURN}),
		nil,
		NewActionNextPlayer(g)})
			
	//Allow adding creatures in develompent phase
	g.Filters = append(g.Filters, NewFilterAllow(&ConditionPhase{PHASE_DEVELOPMENT}, nil, NewActionAddCreature(SOURCE_PROTOTYPE_PLAYER, SOURCE_PROTOTYPE_PLAYER_CARD)))
	//Allow adding pair properties in development phase
	g.Filters = append(g.Filters, NewFilterAllow(&ConditionPhase{PHASE_DEVELOPMENT}, nil, NewActionAddPairProperty(SOURCE_PROTOTYPE_CREATURES_PAIR, SOURCE_PROTOTYPE_PLAYER_CARD_PROPERTY)))
	//Allow adding single properties in development phase
	g.Filters = append(g.Filters, NewFilterAllow(&ConditionPhase{PHASE_DEVELOPMENT}, nil, NewActionAddSingleProperty(SOURCE_PROTOTYPE_PLAYER_CREATURE, SOURCE_PROTOTYPE_PLAYER_CARD_PROPERTY)))
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
			NewActionGetFoodFromBank(SOURCE_PROTOTYPE_PLAYER_CREATURE)))
		
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
						NewConditionEqual(SourceWrapper{FILTER_SOURCE_PARAMETER_TARGET_CREATURE}, FILTER_SOURCE_PARAMETER_RIGHT_CREATURE)),
					NewANDCondition(
						&ConditionActionType{ACTION_REMOVE_PROPERTY},
						NewConditionEqual(SourceWrapper{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY))},
				&FilterDeny{
					NewANDCondition(
						&ConditionActionType{ACTION_ADD_TRAIT},
						NewConditionEqual(SourceWrapper{FILTER_SOURCE_PARAMETER_TRAIT}, FILTER_SOURCE_PARAMETER_ANY_FOOD),
						NewConditionEqual(SourceWrapper{FILTER_SOURCE_PARAMETER_SOURCE}, FILTER_SOURCE_PARAMETER_RIGHT_CREATURE),
						NewConditionEqual(SourceWrapper{TraitsCount{FILTER_SOURCE_PARAMETER_SOURCE, TRAIT_FED}}, 0)),
					NewANDCondition(
						&ConditionActionType{ACTION_REMOVE_PROPERTY},
						NewConditionEqual(SourceWrapper{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY))},
				&FilterDeny{
					NewANDCondition(
						&ConditionActionType{ACTION_ADD_PAIR_PROPERTY},
						NewConditionEqual(SourceWrapper{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY),
						NewConditionEqual(SourceWrapper{FILTER_SOURCE_PARAMETER_PAIR}, FILTER_SOURCE_PARAMETER_PAIR)),
					NewANDCondition(
						&ConditionActionType{ACTION_REMOVE_PROPERTY},
						NewConditionEqual(SourceWrapper{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY))})})
	
}

func (g *Game) AddCard(count int, properties ...Property) {
	for i := 0; i < count; i++ {
		card := g.NewCard(properties...)
		g.Deck = append(g.Deck, card)
	}
}

func (g *Game) NewCard(properties ...Property) *Card {
	if len(properties) == 0 {
		return &Card{}
	}
	card := new(Card)
	for _,property := range properties {
		card.Properties = append(card.Properties, Property {Traits : property.Traits})
	}
	card.ActiveProperty = &card.Properties[0]
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
			//instantiatedFilter := filter.InstantiateFilterPrototype(g, action)
			if filter.CheckCondition(g, action) {
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
			actions := filter.InstantiateFilterPrototype(g, nil).(*FilterAllow).GetActions()
			for _, action := range actions {
				result = append(result, g.InstantiateActionPrototype(action)...)
			}
		}
	}
	return result
}

func (g *Game) ExecuteAction(rawAction *Action) {
	variants := g.InstantiateActionPrototype(rawAction)
	var action *Action
	if len(variants) == 0 {
		return
	}
	if len(variants) > 1 {
		if player, ok := variants[0].Arguments[PARAMETER_PLAYER]; ok {
			action = player.(*Player).MakeChoice(g, variants)
		} else {
			fmt.Println("Something went wrong")
			return
		}
	} else {
		action = variants[0]
	}
	for i, filter := range g.Filters {
		if filter.GetType() == FILTER_ACTION_REPLACE && filter.CheckCondition(g, action) {
			g.Actions.PushFront(filter.InstantiateFilterPrototype(g, action).(*FilterAction).GetAction())
			fmt.Printf("Replaced %#v with %#v because %#v\n", action, filter.InstantiateFilterPrototype(g, action).(*FilterAction).GetAction(), filter.GetCondition())
			return
		}
		if filter.CheckRemoveCondition(g, action) {
			fmt.Printf("Removing filter %#v because &#v", filter, filter.GetCondition())
			g.Filters = append(g.Filters[:i], g.Filters[i+1:]...)
		}
		if filter.GetType() == FILTER_ACTION_EXECUTE_BEFORE && filter.CheckCondition(g, action) {
			g.ExecuteAction(filter.InstantiateFilterPrototype(g, action).(*FilterAction).GetAction())
		} 
	}
	fmt.Printf("Executing action: %#v\n", action)
	action.Execute(g)
	for _, filter := range g.Filters {
		if filter.GetType() == FILTER_ACTION_EXECUTE_AFTER && filter.CheckCondition(g, action) {
			g.Actions.PushFront(filter.InstantiateFilterPrototype(g, action).(*FilterAction).GetAction())
		}
	}
}
