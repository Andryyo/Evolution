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
	FoodBank	  FoodBank
}

type WithTraits interface {
	GetTraits() []TraitType
	AddTrait(trait TraitType)
	RemoveTrait(trait TraitType)
}

type FoodBank struct {
	Count int
}

type Source interface{}

type Card struct {
	ActiveProperty *Property
	Properties     []Property
	Owner          *Player
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
	result := c.Traits
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
		game.Actions.Remove(action)
		game.ExecuteAction(action.Value.(*Action))
		time.Sleep(time.Second)
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
	player.Cards[len(player.Cards)-1].Owner = player
	g.Deck = g.Deck[:deckLen-1]
}

func (g *Game) InitializeDeck() {
	camouflage := Property{Traits : []TraitType {TRAIT_CAMOUFLAGE}}
	burrowing := Property{Traits : []TraitType {TRAIT_BURROWING}}
	sharpVision := Property{Traits : []TraitType {TRAIT_SHART_VISION}}
	/*symbiosys := &Property{Name: "symbiosys"}
	piracy := &Property{Name: "piracy"}
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
	/*g.AddCard(4, symbiosys)
	g.AddCard(4, piracy)
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
	for _, name := range names {
		player := &Player{Name: name, ChoiceMaker: ConsoleChoiceMaker{}}
		g.Players.Value = player
		g.TakeCards(player, 6)
		g.Players = g.Players.Next()
	}
	g.CurrentPlayer = g.Players.Value.(*Player)
}

func (g *Game) InitializeFilters() {
	g.Filters = append(g.Filters, &FilterAction{FILTER_ACTION_EXECUTE_AFTER, &ConditionActionType{ACTION_NEW_PHASE}, NewActionStartTurn(SOURCE_PROTOTYPE_PLAYER)})
	g.Filters = append(g.Filters, &FilterAction{FILTER_ACTION_REPLACE, NewANDCondition(&ConditionActionType{ACTION_START_TURN}, &ConditionTraitCountEqual{PARAMETER_PLAYER, TRAIT_PASS, 1}), NewActionNextPlayer(g)})
	g.Filters = append(g.Filters, &FilterAction{FILTER_ACTION_EXECUTE_BEFORE, &ConditionActionType{ACTION_NEW_PHASE}, NewActionRemoveTrait(FILTER_SOURCE_PARAMETER_ALL_PLAYERS, TRAIT_PASS)})
	g.Filters = append(g.Filters, &FilterAction{FILTER_ACTION_REPLACE, NewANDCondition(&ConditionActionType{ACTION_NEXT_PLAYER},&ConditionTraitCountEqual{FILTER_SOURCE_PARAMETER_ALL_PLAYERS, TRAIT_PASS, 1}), NewActionNewPhase(PHASE_FOOD_BANK_DETERMINATION)})
	g.Filters = append(g.Filters, &FilterAllow{&ConditionPhase{PHASE_DEVELOPMENT}, NewActionAddCreature(SOURCE_PROTOTYPE_PLAYER, SOURCE_PROTOTYPE_PLAYER_CARD)})
	g.Filters = append(g.Filters, &FilterAllow{&ConditionPhase{PHASE_DEVELOPMENT}, NewActionAddCreature(SOURCE_PROTOTYPE_PLAYER, SOURCE_PROTOTYPE_PLAYER_CARD)})
	g.Filters = append(g.Filters, &FilterAllow{&ConditionPhase{PHASE_DEVELOPMENT}, NewActionAddProperty(SOURCE_PROTOTYPE_PLAYER_CREATURE, SOURCE_PROTOTYPE_PLAYER_CARD_PROPERTY)})
	g.Filters = append(g.Filters, &FilterAllow{nil, &Action{ACTION_PASS, map[ArgumentName]Source {}}})
	
	g.Filters = append(g.Filters, 
		&FilterDeny{
			NewANDCondition(
				&ConditionActionType{ACTION_ADD_SINGLE_PROPERTY},
				&NOTCondition{
					&ConditionTraitCountEqual{FILTER_SOURCE_PARAMETER_CREATURE, TRAIT_FAT_TISSUE, 1}},
				&ConditionTraitCountEqual{FILTER_SOURCE_PARAMETER_CREATURE, FILTER_SOURCE_PARAMETER_TRAIT, 1},
			)})
	
	//camouflage
	g.Filters = append(g.Filters,
		&FilterDeny{
			NewANDCondition(
				&ConditionActionType{ACTION_ATTACK},
				&ConditionTraitCountEqual{FILTER_SOURCE_PARAMETER_SOURCE_CREATURE, TRAIT_CAMOUFLAGE, 1},
				&ConditionTraitCountEqual{FILTER_SOURCE_PARAMETER_TARGET_CREATURE, TRAIT_SHART_VISION, 0},	
			)})
	//burrowing
	g.Filters = append(g.Filters,
		&FilterDeny{
			NewANDCondition(
				&ConditionActionType{ACTION_ATTACK},
				&ConditionTraitCountEqual{FILTER_SOURCE_PARAMETER_SOURCE_CREATURE, TRAIT_BURROWING, 1},
				&ConditionTraitCountEqual{FILTER_SOURCE_PARAMETER_SOURCE_CREATURE, TRAIT_FED, 1},
			)})
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

func (g *Game) ActionDenied(action *Action) bool {
	for _, filter := range g.Filters {
		if filter.GetType() == FILTER_DENY {
			instantiatedFilter := filter.InstantiateFilterPrototype(g, action)
			if instantiatedFilter.CheckCondition(g, action) {
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
			action := filter.(*FilterAllow).GetAction()
			result = append(result, g.InstantiateActionPrototype(action)...)
		}
	}
	return result
}

func (g *Game) ExecuteAction(rawAction *Action) {
	variants := g.InstantiateActionPrototype(rawAction)
	var action *Action
	if len(variants) > 1 {
		if player, ok := variants[0].Arguments[PARAMETER_PLAYER]; ok {
			action = player.(*Player).MakeChoice(variants)
		} else {
			fmt.Println("Something went wrong")
			return
		}
	} else {
		action = variants[0]
	}
	for _, filter := range g.Filters {
		if filter.GetType() == FILTER_ACTION_EXECUTE_AFTER && filter.CheckCondition(g, action) {
			g.Actions.PushFront(filter.(*FilterAction).GetAction().InstantiateFilterPrototypeAction(g, action))
		}
		if filter.GetType() == FILTER_ACTION_REPLACE && filter.CheckCondition(g, action) {
			g.Actions.PushFront(filter.(*FilterAction).GetAction().InstantiateFilterPrototypeAction(g, action))
			fmt.Printf("Replaced %#v with %#v because %#v\n", action, filter.(*FilterAction).GetAction().InstantiateFilterPrototypeAction(g, action), filter.GetCondition())
			return
		} 
	}
	fmt.Printf("Executing action: %#v\n", action)
	action.Execute(g)
	for _, filter := range g.Filters {
		if filter.GetType() == FILTER_ACTION_EXECUTE_BEFORE && filter.CheckCondition(g, action) {
			g.Actions.PushFront(filter.(*FilterAction).GetAction().InstantiateFilterPrototypeAction(g, action))
		}
	}
}
