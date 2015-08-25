// Evolution
package EvolutionEngine

import (
	"container/ring"
	"fmt"
	"math/rand"
	"log"
	"container/list"
	"time"
)

type Game struct {
	Players       *ring.Ring
	ObserversChs  []chan struct{Action *Action; Game *Game}
	PlayersCount  int
	Deck          []*Card
	Filters       []Filter
	CurrentPhase  PhaseType
	CurrentPlayer *Player
	Food          int  
}

func (g *Game) AddObserver() chan struct{Action *Action; Game *Game} {
	ch := make(chan struct{Action *Action; Game *Game})
	g.ObserversChs = append(g.ObserversChs, ch)
	return ch
}

func (g *Game) RemoveObserver(ch chan struct{Action *Action; Game *Game}) {
	for i, tmp := range g.ObserversChs {
		if tmp == ch {
			g.ObserversChs = append(g.ObserversChs[:i], g.ObserversChs[i+1:]...)
			return
		}
	}
}

func (g *Game) NotifyAll(action *Action) {
	log.Printf("%#v\n", action)
	msg := struct{Action *Action; Game *Game}{action, g}
	g.Players.Do(func (val interface{}) {
		val.(*Player).NotifyCh <- msg
	})
	for _,ch := range g.ObserversChs {
		ch <- msg
	}
}

type WithTraits interface {
	GetTraits() []TraitType
	AddTrait(trait TraitType)
	RemoveTrait(trait TraitType)
	ContainsTrait(trait TraitType) bool
}

type Source interface{}

type Card struct {
	ActiveProperty *Property
	Properties     []*Property
	Owners         []Source
}

func (c Card) GoString() string {
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

func (c *Property) ContainsTrait(trait TraitType) bool {
	for _,t := range c.Traits {
		if trait == t {
			return true
		}
	}
	return false
}

func (c Property) GoString() string {
	len := len(c.Traits)
	if len == 0 {
		return "()"
	}
	result := "(" + string(c.Traits[0])
	for i := 1; i<len;i++ {
		result += "/" + string(c.Traits[i])
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

func (c Creature) GoString() string {
	traits := c.GetTraits()
	len := len(traits)
	if len == 0 {
		return "(Creature)"
	}
	result := "(Creature : " + string(traits[0])
	for i := 1; i<len;i++ {
		result += "/" + string(traits[i])
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
		result = append(result, card.ActiveProperty.GetTraits()...)
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

func (c *Creature) ContainsTrait(trait TraitType) bool {
	for _,t := range c.Traits {
		if trait == t {
			return true
		}
	}
	for _,card := range c.Tail {
		if card.ActiveProperty.ContainsTrait(trait) {
			return true
		}
	}
	return false
}

type Choice struct {
	Actions []*Action
	Game *Game
}

type Player struct {
	Id string
	PendingChoice *Choice
	AvailableChoiceCh chan *Choice
	ChoiceCh chan int
	NotifyCh chan struct{Action *Action; Game *Game}
	Creatures []*Creature
	Cards     []*Card
	Traits    []TraitType
	Occupied  bool
}

func (p *Player) MakeChoice(game *Game, actions []*Action) *Action {
	if len(actions) == 0 {
		return nil
	}
	if len(actions) == 1 {
		return actions[0]
	}
	choice := &Choice{actions, game}
	p.PendingChoice = choice
	choiceNum := -1
	for choiceNum < 0 || choiceNum >= len(actions) {
		p.AvailableChoiceCh <- choice
		choiceNum = <- p.ChoiceCh
	}
	p.PendingChoice = nil
	return actions[choiceNum]
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

func (c *Player) ContainsTrait(trait TraitType) bool {
	for _,t := range c.Traits {
		if trait == t {
			return true
		}
	}
	return false
}

func (p *Player) RemoveCreature(creature *Creature) {
	for i := range p.Creatures {
		if p.Creatures[i] == creature {
			p.Creatures = append(p.Creatures[:i], p.Creatures[i+1:]...)
			return	
		}
	}
}

func NewGame(playersCount int) *Game {
	game := new(Game)
	log.Println("Initializing cards filters")
	game.InitializeCardsFilters()
	log.Println("Initializing base game flow")
	game.InitializeBaseGameFlow()
	log.Println("Initializing deck")
	game.InitializeDeck()
	log.Println("Initializing players")
	game.InitializePlayers(playersCount)
	return game
}

func (g *Game) Start() {
	g.ExecuteAction(NewActionNewPhase(PHASE_DEVELOPMENT))
}

func (g *Game) TakeCards(player *Player, count int) {
	for i := 0; i < count; i++ {
		g.TakeCard(player)
	}
}

func (g *Game) TakeCard(player *Player) {
	deckLen := len(g.Deck)
	if deckLen == 0 {
		return
	}
	card := g.Deck[deckLen-1]
	card.Owners = []Source{player}
	player.Cards = append(player.Cards, card)
	player.Cards[len(player.Cards)-1].Owners = []Source {player}
	g.Deck = g.Deck[:deckLen-1]
}

func (g *Game) InitializeDeck() {
	camouflage := &Property{Traits : []TraitType {TRAIT_CAMOUFLAGE}}
	burrowing := &Property{Traits : []TraitType {TRAIT_BURROWING}}
	sharpVision := &Property{Traits : []TraitType {TRAIT_SHARP_VISION}}
	symbiosis := &Property{Traits : []TraitType {TRAIT_PAIR, TRAIT_SYMBIOSIS}}
	piracy := &Property{Traits : []TraitType {TRAIT_PIRACY}}
	grazing := &Property{Traits : []TraitType {TRAIT_GRAZING}}
	tailLoss := &Property{Traits : []TraitType {TRAIT_TAIL_LOSS}}
	hibernation := &Property{Traits : []TraitType{TRAIT_HIBERNATION}}
	poisonous := &Property{Traits: []TraitType{TRAIT_POISONOUS}}
	communication := &Property{Traits: []TraitType{TRAIT_COMMUNICATION, TRAIT_PAIR}}
	scavenger := &Property{Traits: []TraitType{TRAIT_SCAVENGER}}
	running := &Property{Traits: []TraitType{TRAIT_RUNNING}}
	mimicry := &Property{Traits: []TraitType{TRAIT_MIMICRY}}
	swimming := &Property{Traits: []TraitType{TRAIT_SWIMMING}}
	parasite := &Property{Traits : []TraitType {TRAIT_PARASITE, TRAIT_REQUIRE_FOOD, TRAIT_REQUIRE_FOOD}}
	carnivorous := &Property{Traits : []TraitType {TRAIT_CARNIVOROUS, TRAIT_REQUIRE_FOOD}}
	fatTissue := &Property{Traits : []TraitType {TRAIT_FAT_TISSUE}}
	cooperation := &Property{Traits : []TraitType {TRAIT_COOPERATION, TRAIT_PAIR}}
	highBodyWeight := &Property{Traits : []TraitType {TRAIT_HIGH_BODY_WEIGHT, TRAIT_REQUIRE_FOOD}}
	
	g.Deck = make([]*Card, 0, 84)
	g.AddCard(4, camouflage, fatTissue)
	g.AddCard(4, burrowing, fatTissue)
	g.AddCard(4, sharpVision, fatTissue)
	g.AddCard(4, symbiosis)
	g.AddCard(4, piracy)
	g.AddCard(4, grazing, fatTissue)
	g.AddCard(4, tailLoss)
	g.AddCard(4, hibernation, carnivorous)
	g.AddCard(4, poisonous, carnivorous)
	g.AddCard(4, communication, carnivorous)
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

func (g *Game) InitializePlayers(playersCount int) {
	if playersCount == 0 {
		return
	}
	g.PlayersCount = playersCount
	g.Players = ring.New(g.PlayersCount)
	for i := 0 ; i<g.PlayersCount; i++ {
		player := &Player{}
		player.Id = fmt.Sprintf("%p", player)
		player.Occupied = false
		player.AvailableChoiceCh = make(chan *Choice)
		player.ChoiceCh = make(chan int)
		player.NotifyCh = make(chan struct{Action *Action; Game *Game})
		g.Players.Value = player
		g.TakeCards(player, 6)
		g.Players = g.Players.Next()
	}
	g.CurrentPlayer = g.Players.Value.(*Player)
}

func (g *Game) AddFilter(filter Filter) {
	g.Filters = append(g.Filters, filter)
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

func (g *Game) GetUnoccupiedPlayer() *Player {
	if !g.Players.Value.(*Player).Occupied {
		return g.Players.Value.(*Player)
	}
	for p:= g.Players.Next(); p!=g.Players; p = p.Next() {
		if !p.Value.(*Player).Occupied {
			return p.Value.(*Player)
		}
	}
	return nil
}

func (g *Game) GetPlayerById(id string) *Player {
	log.Printf("%s %s\n", g.Players.Value.(*Player).Id, id)
	if g.Players.Value.(*Player).Id == id {
		return g.Players.Value.(*Player)
	}
	for p:= g.Players.Next(); p!=g.Players; p = p.Next() {
		if p.Value.(*Player).Id == id {
			return p.Value.(*Player)
		}
	}
	return nil
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

func (g *Game) GetAllowedActions() []*Action {
	var result []*Action
	for _, filter := range g.Filters {
		if filter.GetType() == FILTER_ALLOW && filter.CheckCondition(g, nil) {
			actions := filter.(*FilterAllow).GetActions(g)
			for _, action := range actions {
				if action.Type == ACTION_SELECT {
					result = append(result, g.ExpandActionSelect(action)...)
				} else {
					if !g.ActionDenied(action) {
						result = append(result, action)
					}
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
		stack.Remove(stackFront)
		action := stackFront.Value.(*Action)
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
		if replaced {
			continue
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
		action.Execute(g)
		g.NotifyAll(action)
		for _, filter := range g.Filters {
			if filter.GetType() == FILTER_ACTION_EXECUTE_AFTER && filter.CheckCondition(g, action) {
				stack.PushBack(filter.InstantiateFilterPrototype(g, action, true).(*FilterAction).GetAction())
			}
		}
		time.Sleep(250 * time.Millisecond)
	}
}
