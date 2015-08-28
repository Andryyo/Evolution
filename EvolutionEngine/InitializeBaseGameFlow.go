// InitializeBaseGameFlow
package EvolutionEngine

func (g *Game) InitializeBaseGameFlow() {
	
	//Remove all pass trait on phase start
	g.AddFilter(&FilterAction{
			FILTER_ACTION_EXECUTE_BEFORE, 
			NewANDCondition(
				&ConditionActionType{ACTION_NEW_PHASE}, 
				NewORCondition(&ConditionPhase{PHASE_DEVELOPMENT}, &ConditionPhase{PHASE_FEEDING})),
			nil,
			NewActionRemoveTrait(FILTER_SOURCE_PARAMETER_ALL_PLAYERS, TRAIT_PASS)})
	
	
	//Start player turn on phase start
	g.AddFilter(&FilterAction{
			FILTER_ACTION_EXECUTE_AFTER, 
			NewANDCondition(
				NewORCondition(&ConditionPhase{PHASE_DEVELOPMENT}, &ConditionPhase{PHASE_FEEDING}),
				&ConditionActionType{ACTION_NEW_PHASE}), 
			nil,
			NewActionStartTurn(FILTER_SOURCE_PARAMETER_CURRENT_PLAYER)})
	
	//Start selecting after turn start
	g.AddFilter(&FilterAction{
			FILTER_ACTION_EXECUTE_AFTER, 
			NewANDCondition(
				NewORCondition(&ConditionPhase{PHASE_DEVELOPMENT}, &ConditionPhase{PHASE_FEEDING}),
				&ConditionActionType{ACTION_START_TURN}), 
			nil,
			NewActionSelectFromAvailableActions()})
	
	//In feeding phase player make turns, until pass
	g.AddFilter(&FilterAction{
			FILTER_ACTION_EXECUTE_AFTER,
			NewANDCondition(&ConditionPhase{PHASE_FEEDING},&ConditionActionType{ACTION_SELECT_FROM_AVAILABLE_ACTIONS}),
			nil,
			NewActionSelectFromAvailableActions()})
	
	// In feeding phase, allow only one food gain action per turn
	g.AddFilter(&FilterAction{
		FILTER_ACTION_EXECUTE_AFTER,
		NewANDCondition(
			&ConditionPhase{PHASE_FEEDING},
			&ConditionActionType{ACTION_GET_FOOD_FROM_BANK}),
		nil,
		NewActionAddFilters(
			&FilterDeny{
				NewORCondition(
					NewANDCondition(
						&ConditionActionType{ACTION_GET_FOOD_FROM_BANK},
						NewConditionEqual(InstantiationOff{TraitsCount{FILTER_SOURCE_PARAMETER_CREATURE, TRAIT_ADDITIONAL_GET_FOOD_FROM_BANK}}, 0)),
					&ConditionActionType{ACTION_ATTACK},
					&ConditionActionType{ACTION_BURN_FAT}),
				&ConditionActionType{ACTION_START_TURN}})})
			
	//In development phase player pass turn to next player
	g.AddFilter(&FilterAction{
		FILTER_ACTION_EXECUTE_AFTER,
		NewANDCondition(&ConditionPhase{PHASE_DEVELOPMENT}, &ConditionActionType{ACTION_SELECT_FROM_AVAILABLE_ACTIONS}),
		nil,
		NewActionNextPlayer(g)})
					
	//Allow adding creatures in develompent phase
	g.AddFilter(NewFilterAllow(
		&ConditionPhase{PHASE_DEVELOPMENT}, 
		nil, 
		NewActionAddCreature(
			FILTER_SOURCE_PARAMETER_CURRENT_PLAYER, 
			FILTER_SOURCE_PARAMETER_ONE_OF_CURRENT_PLAYER_CARDS)))
	//Allow adding pair properties in development phase
	g.AddFilter(NewFilterAllow(
		&ConditionPhase{PHASE_DEVELOPMENT}, 
		nil, 
		NewActionAddPairProperty(
			FILTER_SOURCE_PARAMETER_ONE_OF_CURRENT_PLAYER_CREATURES_PAIR, 
			FILTER_SOURCE_PARAMETER_ONE_OF_CURRENT_PLAYER_CARDS_PROPERTIES)))
			
	g.AddFilter(&FilterAction{
		FILTER_ACTION_EXECUTE_AFTER,
		&ConditionActionType{ACTION_ADD_PAIR_PROPERTY},
		nil,
		NewActionAddFilters(
			&FilterDeny{
				NewANDCondition(
					&ConditionActionType{ACTION_ADD_PAIR_PROPERTY},
					NewConditionDeepEqual(
						Accessor{FILTER_SOURCE_PARAMETER_PROPERTY, ACCESSOR_MODE_TRAITS}, 
						InstantiationOff{Accessor{FILTER_SOURCE_PARAMETER_PROPERTY, ACCESSOR_MODE_TRAITS}}),
					NewConditionDeepEqual(FILTER_SOURCE_PARAMETER_PAIR, InstantiationOff{FILTER_SOURCE_PARAMETER_PAIR})),
				NewANDCondition(
						&ConditionActionType{ACTION_REMOVE_PROPERTY},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY))})})
						
	g.AddFilter(&FilterAction{
		FILTER_ACTION_EXECUTE_AFTER,
		NewANDCondition(
			&ConditionActionType{ACTION_ADD_SINGLE_PROPERTY},
			&NOTCondition{NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_FAT_TISSUE}, 1)}),
		nil,
		NewActionAddFilters(
			&FilterDeny{
				NewANDCondition(
					&ConditionActionType{ACTION_ADD_SINGLE_PROPERTY},
					NewConditionDeepEqual(
						Accessor{FILTER_SOURCE_PARAMETER_PROPERTY, ACCESSOR_MODE_TRAITS}, 
						InstantiationOff{Accessor{FILTER_SOURCE_PARAMETER_PROPERTY, ACCESSOR_MODE_TRAITS}}),
					NewConditionDeepEqual(FILTER_SOURCE_PARAMETER_CREATURE, InstantiationOff{FILTER_SOURCE_PARAMETER_CREATURE})),
				NewANDCondition(
						&ConditionActionType{ACTION_REMOVE_PROPERTY},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY))})})
		
	//Allow adding single properties in development phase
	g.AddFilter(NewFilterAllow(
		&ConditionPhase{PHASE_DEVELOPMENT}, 
		nil, 
		NewActionAddSingleProperty(
			FILTER_SOURCE_PARAMETER_ONE_OF_CREATURES, 
			FILTER_SOURCE_PARAMETER_ONE_OF_CURRENT_PLAYER_CARDS_PROPERTIES)))			
	
	//Deny adding single properties if
	g.AddFilter(&FilterDeny{
			NewANDCondition(
				&ConditionActionType{ACTION_ADD_SINGLE_PROPERTY},
				NewORCondition(
					NewANDCondition(
						NewConditionEqual(Accessor{FILTER_SOURCE_PARAMETER_CREATURE, ACCESSOR_MODE_CREATURE_OWNER}, FILTER_SOURCE_PARAMETER_CURRENT_PLAYER),
						NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_PARASITE}, 1)),
					NewANDCondition(
						&NOTCondition{NewConditionEqual(Accessor{FILTER_SOURCE_PARAMETER_CREATURE, ACCESSOR_MODE_CREATURE_OWNER}, FILTER_SOURCE_PARAMETER_CURRENT_PLAYER)},
						NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_PARASITE}, 0)),
					NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_PAIR}, 1),
					NewANDCondition(
						NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_SCAVENGER}, 1),
						NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_CREATURE, TRAIT_CARNIVOROUS}, 1)),
					NewANDCondition(
						NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_CREATURE, TRAIT_SCAVENGER}, 1),
						NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_CARNIVOROUS}, 1)))),
			nil})	
	
	//Deny adding pair properties is
	g.AddFilter(&FilterDeny{
			NewANDCondition(
				&ConditionActionType{ACTION_ADD_PAIR_PROPERTY},
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_PAIR}, 0)),
			nil,
		})
		
	//Allow transform fat to food
	g.AddFilter(NewFilterAllow(
		&ConditionPhase{PHASE_FEEDING},
		nil,
		NewActionBurnFat(FILTER_SOURCE_PARAMETER_ONE_OF_CURRENT_PLAYER_CREATURES)))
	
	//Add ACTION_BURN_FAT related filters
	g.AddFilter(&FilterAction{
		FILTER_ACTION_EXECUTE_AFTER,
		NewANDCondition(
			&ConditionPhase{PHASE_FEEDING},
			&ConditionActionType{ACTION_BURN_FAT},
			NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_CREATURE, TRAIT_BURNED_FAT}, 0)),
		nil,
		NewActionSequence(
			NewActionAddTrait(FILTER_SOURCE_PARAMETER_CREATURE, TRAIT_BURNED_FAT),
			NewActionAddFilters(
				&FilterDeny{
					NewORCondition(
						NewANDCondition(
							&ConditionActionType{ACTION_BURN_FAT},
							&NOTCondition{NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_CREATURE}, FILTER_SOURCE_PARAMETER_CREATURE)}),
						&ConditionActionType{ACTION_ATTACK},
						&ConditionActionType{ACTION_GET_FOOD_FROM_BANK}),
					&ConditionActionType{ACTION_START_TURN}},
				&FilterAction{
					FILTER_ACTION_EXECUTE_BEFORE,
					&ConditionActionType{ACTION_START_TURN},
					&ConditionActionType{ACTION_START_TURN},
					NewActionRemoveTrait(FILTER_SOURCE_PARAMETER_CREATURE, TRAIT_BURNED_FAT)}))})
	
	//Deny transform fat to food
	g.AddFilter(&FilterDeny{
		NewANDCondition(
			&ConditionActionType{ACTION_BURN_FAT},
			NewORCondition(
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_CREATURE, TRAIT_FED}, 1),
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_CREATURE, TRAIT_FAT}, 0),
				&ConditionActionDenied{NewActionAddTrait(FILTER_SOURCE_PARAMETER_CREATURE, TRAIT_FOOD)})),
		nil})
	
	//Alow pass turn to next player in feeding mode
	g.AddFilter(NewFilterAllow(
			&ConditionPhase{PHASE_FEEDING},
			nil,
			NewActionEndTurn()))	
	
	//Allow pass in development and feeding phase
	g.AddFilter(NewFilterAllow(
		NewORCondition(&ConditionPhase{PHASE_DEVELOPMENT}, 
		&ConditionPhase{PHASE_FEEDING}), 
		nil, 
		&Action{ACTION_PASS, map[ArgumentName]Source {}}))
	
	//If all players pass in development phase, start food bank determination
	g.AddFilter(&FilterAction{
			FILTER_ACTION_REPLACE, 
			NewANDCondition(
				&ConditionPhase{PHASE_DEVELOPMENT},
				&ConditionActionType{ACTION_NEXT_PLAYER},
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_ALL_PLAYERS, TRAIT_PASS}, 1)), 
			nil,
			NewActionNewPhase(PHASE_FOOD_BANK_DETERMINATION)})
	
	//If all players pass in feeding phase, start extinction
	g.AddFilter( 
		&FilterAction{
			FILTER_ACTION_REPLACE, 
			NewANDCondition(
				&ConditionPhase{PHASE_FEEDING},
				&ConditionActionType{ACTION_NEXT_PLAYER},
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_ALL_PLAYERS, TRAIT_PASS}, 1)), 
			nil,
			NewActionNewPhase(PHASE_EXTINCTION)})
	
	g.AddFilter(&FilterAction{
		FILTER_ACTION_REPLACE,
		NewANDCondition(
			&ConditionActionType{ACTION_TAKE_CARDS},
			NewConditionEqual(FILTER_SOURCE_PARAMETER_BANK_CARDS_COUNT, 0)),
		nil,
		NewActionNewPhase(PHASE_FINAL)})
	
	//Execute extinction action
	g.AddFilter(&FilterAction{
			FILTER_ACTION_EXECUTE_AFTER,
			NewANDCondition(
				&ConditionActionType{ACTION_NEW_PHASE},
				NewConditionEqual(FILTER_SOURCE_PARAMETER_PHASE, PHASE_EXTINCTION)),
			nil,
			&Action{ACTION_EXTINCT, map[ArgumentName]Source{}}})

	//After extinct, start take card
	g.AddFilter(&FilterAction{
			FILTER_ACTION_EXECUTE_AFTER,
			&ConditionActionType{ACTION_EXTINCT},
			nil,
			NewActionTakeCards()})


	//After take cards, start development again
	g.AddFilter(&FilterAction{
			FILTER_ACTION_EXECUTE_AFTER,
			&ConditionActionType{ACTION_TAKE_CARDS},
			nil,
			NewActionNewPhase(PHASE_DEVELOPMENT)})

	//If player pass - replace his turn with NextTurn
	g.AddFilter(&FilterAction{
			FILTER_ACTION_REPLACE, 
			NewANDCondition(
				&ConditionActionType{ACTION_SELECT_FROM_AVAILABLE_ACTIONS}, 
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_CURRENT_PLAYER, TRAIT_PASS}, 1)), 
			nil,
			NewActionNextPlayer(g)})
			
	//If player pass - replace his turn with NextTurn
	g.AddFilter(&FilterAction{
			FILTER_ACTION_REPLACE, 
			NewANDCondition(
				&ConditionActionType{ACTION_START_TURN}, 
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_PLAYER, TRAIT_PASS}, 1)), 
			nil,
			NewActionNextPlayer(g)})
			
	//Determine food bank
	g.AddFilter(&FilterAction{
			FILTER_ACTION_EXECUTE_AFTER,
			NewANDCondition(
				&ConditionActionType{ACTION_NEW_PHASE},
				NewConditionEqual(FILTER_SOURCE_PARAMETER_PHASE, PHASE_FOOD_BANK_DETERMINATION)),
			nil,
			&Action{ACTION_DETERMINE_FOOD_BANK, map[ArgumentName]Source {}}})
		
	//After food bank determination, start feeding phase
	g.AddFilter(&FilterAction{
				FILTER_ACTION_EXECUTE_AFTER,
				&ConditionActionType{ACTION_DETERMINE_FOOD_BANK},
				nil,
				NewActionNewPhase(PHASE_FEEDING)})
	
	//Allow get food from bank for creatures
	g.AddFilter(NewFilterAllow(
			&ConditionPhase{PHASE_FEEDING},
			nil,
			NewActionGetFoodFromBank(FILTER_SOURCE_PARAMETER_ONE_OF_CURRENT_PLAYER_CREATURES)))
		
	//Deny get food from bank
	g.AddFilter(&FilterDeny{
			NewANDCondition(
				&ConditionActionType{ACTION_GET_FOOD_FROM_BANK}, 
				NewORCondition(
					&ConditionActionDenied{NewActionGainFood(FILTER_SOURCE_PARAMETER_CREATURE)},
					NewConditionEqual(FILTER_SOURCE_PARAMETER_FOOD_BANK_COUNT, 0))),
			nil})
	
	//Deny food get if creature already full
	g.AddFilter(&FilterDeny{
			NewANDCondition(
				&ConditionActionType{ACTION_ADD_TRAIT},
				NewConditionEqual(FILTER_SOURCE_PARAMETER_TRAIT, FILTER_SOURCE_PARAMETER_ANY_FOOD),
				NewConditionEqual(
					TraitsCount{FILTER_SOURCE_PARAMETER_SOURCE, FILTER_SOURCE_PARAMETER_ALL_FOOD_AND_FAT},
					TraitsCount{FILTER_SOURCE_PARAMETER_SOURCE, FILTER_SOURCE_PARAMETER_FOOD_AND_FAT_LIMIT})),
			nil})
	
	//Replace food get with fat get
	g.AddFilter(&FilterAction{
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
	g.AddFilter(&FilterAction{
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
	g.AddFilter(&FilterDeny{
			NewANDCondition(
				&ConditionPhase{PHASE_FEEDING},
				&ConditionActionType{ACTION_REMOVE_TRAIT},
				NewConditionEqual(FILTER_SOURCE_PARAMETER_TRAIT, TRAIT_FOOD),
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_SOURCE , TRAIT_FOOD}, 0)),
			nil})
	
	//Deny remove additional food from creature, if it have none
	g.AddFilter(
		&FilterDeny{
			NewANDCondition(
				&ConditionPhase{PHASE_FEEDING},
				&ConditionActionType{ACTION_REMOVE_TRAIT},
				NewConditionEqual(FILTER_SOURCE_PARAMETER_TRAIT, TRAIT_ADDITIONAL_FOOD),
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_SOURCE , TRAIT_ADDITIONAL_FOOD}, 0)),
			nil})	
			
	g.AddFilter(
		&FilterDeny{
			NewANDCondition(
				&ConditionActionType{ACTION_GAIN_FOOD},
				&ConditionActionDenied{NewActionAddTrait(FILTER_SOURCE_PARAMETER_CREATURE, TRAIT_FOOD)}),
			nil})
			
	g.AddFilter(
		&FilterDeny{
			NewANDCondition(
				&ConditionActionType{ACTION_GAIN_ADDITIONAL_FOOD},
				&ConditionActionDenied{NewActionAddTrait(FILTER_SOURCE_PARAMETER_CREATURE, TRAIT_ADDITIONAL_FOOD)}),
			nil})
}