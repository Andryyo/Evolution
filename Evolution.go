// Evolution
package main

import (
	"container/ring"
	"fmt"
	"math/rand"
)

type Game struct {
	Players       *ring.Ring
	Deck          []*Card
	Filters       []Filter
	CurrentPhase  Phase
	CurrentPlayer *Player
}

type Source interface{}

type SourcePrototype int

const (
	SOURCE_PROTOTYPE_NONE SourcePrototype = iota
	SOURCE_PROTOTYPE_CREATURES_PAIR
	SOURCE_PROTOTYPE_PROPERTY
	SOURCE_PROTOTYPE_PLAYER_CREATURE_PROPERTY
	SOURCE_PROTOTYPE_PLAYER_CARD
	SOURCE_PROTOTYPE_PLAYER_CARD_PROPERTY_NOT_SELECTED
	SOURCE_PROTOTYPE_PLAYER_CARD_PROPERTY
	SOURCE_PROTOTYPE_OPPONENT_CARD
	SOURCE_PROTOTYPE_PLAYER
	SOURCE_PROTOTYPE_OPPONENT
	SOURCE_PROTOTYPE_ANY_PLAYER
	SOURCE_PROTOTYPE_PLAYER_CREATURE
	SOURCE_PROTOTYPE_PLAYER_CREATURE_NOT_SENDER
	SOURCE_PROTOTYPE_DECK
	SOURCE_PROTOTYPE_OWN_CREATURE
	SOURCE_PROTOTYPE_OWN_CREATURE_NOT_SENDER
	SOURCE_PROTOTYPE_DESK
)

type PhaseType int

const (
	PHASE_DEVELOPMENT PhaseType = iota
	PHASE_FOOD_BANK_DETERMINATION
	PHASE_FEEDING
	PHASE_EXTINCTION
)

type Card struct {
	ActiveProperty *Property
	Properties     []*Property
	Owner          *Player
}

func (c *Card) GoString() string {
	propertiesCount := len(c.Properties)
	if propertiesCount == 0 {
		return ""
	}
	result := c.Properties[0].Name
	for i := 1; i < propertiesCount; i++ {
		result += "/" + c.Properties[i].Name
	}
	return result
}

type Property struct {
	Name           string
	Filters        []Filter
	Actions        []*Action
	ContainingCard *Card
}

func (p *Property) String() string {
	return p.Name
}

type Creature struct {
	Head  *Card
	Tail  []*Card
	Owner *Player
}

type Player struct {
	Name      string
	Creatures []*Creature
	Cards     []*Card
}

type Phase interface{}

type DevelopmentPhase struct {
	game *Game
}

func NewDevelopmentPhase(game *Game) *DevelopmentPhase {
	return &DevelopmentPhase{game}
}

func (p *DevelopmentPhase) ChooseAction(player *Player, action *Action) {

}

func (g *Game) GetInstantiationVariants(arguments map[ArgumentName]Source, argumentsNames []ArgumentName, argumentNumber int) []map[ArgumentName]Source {
	argumentsLen := len(arguments)
	if argumentsLen == 0 {
		return []map[ArgumentName]Source{}
	}
	argumentName := argumentsNames[argumentNumber]
	argument := arguments[argumentName]
	instantiatedArguments := g.InstantiateArgument(argument)
	if argumentNumber == argumentsLen-1 {
		result := make([]map[ArgumentName]Source, 0, argumentsLen)
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
	var variants []map[ArgumentName]Source = g.GetInstantiationVariants(prototype.Arguments, undefinedArgumentsNames, 0)
	for _, variant := range variants {
		for _, definedArgumentName := range definedArgumentsNames {
			variant[definedArgumentName] = prototype.Arguments[definedArgumentName]
		}
		result = append(result, &Action{prototype.Type, variant})
	}
	return result
}

func NewGame(players ...string) *Game {
	fmt.Println("Here is library start!")
	game := new(Game)
	game.InitializeDeck()
	game.InitializePlayers(players...)
	game.InitializeFilters()
	game.ExecuteAction(NewActionSequence(NewActionNewPhase(PHASE_DEVELOPMENT), NewActionStartTurn(game.CurrentPlayer)))
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
	camouflage := &Property{Name: "camouflage"}
	burrowing := &Property{Name: "burrowing"}
	sharpVision := &Property{Name: "sharpVision"}
	symbiosys := &Property{Name: "symbiosys"}
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
	highBodyWeight := &Property{Name: "highBodyWeight"}
	g.Deck = make([]*Card, 0, 84)
	g.AddCard(4, camouflage)
	g.AddCard(4, burrowing)
	g.AddCard(4, sharpVision)
	g.AddCard(4, symbiosys)
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
	g.AddCard(4, highBodyWeight, fatTissue)
	g.ShuffleDeck()
}

func (g *Game) InitializePlayers(names ...string) {
	g.Players = ring.New(len(names))
	for _, name := range names {
		player := &Player{Name: name}
		g.Players.Value = player
		g.TakeCards(player, 6)
		g.Players = g.Players.Next()
	}
	g.CurrentPlayer = g.Players.Value.(*Player)
}

func (g *Game) InitializeFilters() {
	g.Filters = append(g.Filters, &FilterAllow{NewActionAddCreature(SOURCE_PROTOTYPE_PLAYER, SOURCE_PROTOTYPE_PLAYER_CARD)})
	g.Filters = append(g.Filters, &FilterAllow{NewActionAddProperty(SOURCE_PROTOTYPE_PLAYER_CREATURE, SOURCE_PROTOTYPE_PLAYER_CARD_PROPERTY)})
}

func (g *Game) AddCard(count int, properties ...*Property) {
	for i := 0; i < count; i++ {
		g.Deck = append(g.Deck, g.NewCard(properties...))
	}
}

func (g *Game) NewCard(properties ...*Property) *Card {
	if len(properties) == 0 {
		return nil
	}
	card := new(Card)
	card.Properties = properties
	card.ActiveProperty = card.Properties[0]
	for i := range card.Properties {
		card.Properties[i].ContainingCard = card
	}
	card.Owner = nil
	return &Card{properties[0], properties, nil}
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
		if filter.GetType() == FILTER_DENY && (filter.CheckCondition(g, action)) {
			return true
		}
	}
	return false
}

func (g *Game) GetAlowedActions() []*Action {
	var result []*Action
	for _, filter := range g.Filters {
		if filter.GetType() == FILTER_ALLOW {
			action := filter.(*FilterAllow).GetAction()
			if !g.ActionDenied(action) {
				result = append(result, g.InstantiateActionPrototype(action)...)
			}
		}
	}
	return result
}

func (g *Game) ExecuteAction(action *Action) {
	if g.ActionDenied(action) {
		return
	}
	for _, filter := range g.Filters {
		if filter.CheckCondition(g, action) && filter.GetType() == FILTER_MODIFY {
			filter.(FilterModify).ModifyAction(action)
		}
	}
	fmt.Printf("Executing action: %#v\n", action)
	action.Execute(g)
	for _, filter := range g.Filters {
		if filter.CheckCondition(g, action) && filter.GetType() == FILTER_ACTION {
			g.ExecuteAction(filter.(*FilterAction).GetAction())
		}
	}
}
