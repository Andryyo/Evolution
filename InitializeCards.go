// InitializeCards
package main


func (g *Game) InitializeCardsFilters() {
	
	//camouflage
	g.Filters = append(g.Filters,
		&FilterDeny{
			NewANDCondition(
				&ConditionActionType{ACTION_ATTACK},
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_SOURCE_CREATURE, TRAIT_CAMOUFLAGE}, 0),
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_TARGET_CREATURE, TRAIT_SHART_VISION}, 1)),
			nil})
	//burrowing
	g.Filters = append(g.Filters,
		&FilterDeny{
			NewANDCondition(
				&ConditionActionType{ACTION_ATTACK},
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_TARGET_CREATURE, TRAIT_BURROWING}, 1),
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_TARGET_CREATURE, TRAIT_FED}, 1)),
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
						NewORCondition(&ConditionActionType{ACTION_GAIN_FOOD},&ConditionActionType{ACTION_GAIN_ADDITIONAL_FOOD}),
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
					&ConditionActionDenied{NewActionGainAdditionalFood(FILTER_SOURCE_PARAMETER_SOURCE_CREATURE, 1)})),	
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
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_SOURCE_CREATURE, TRAIT_HIGH_BODY_WEIGHT}, 0),
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_TARGET_CREATURE, TRAIT_HIGH_BODY_WEIGHT}, 1)),
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
					NewANDCondition(
						&ConditionActionType{ACTION_ATTACK},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_SOURCE_CREATURE}, Accessor{FILTER_SOURCE_PARAMETER_PROPERTY, ACCESSOR_MODE_PROPERTY_OWNER})),
					NewANDCondition(
						&ConditionActionType{ACTION_REMOVE_PROPERTY},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY)),
					NewActionAddTrait(FILTER_SOURCE_PARAMETER_PROPERTY,TRAIT_USED)},
				&FilterAction{
					FILTER_ACTION_EXECUTE_BEFORE,
					NewANDCondition(
						&ConditionActionType{ACTION_NEW_PHASE},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PHASE}, PHASE_EXTINCTION)),
					NewANDCondition(
						&ConditionActionType{ACTION_REMOVE_PROPERTY},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY)),
					NewActionRemoveTrait(FILTER_SOURCE_PARAMETER_PROPERTY,TRAIT_USED)},
				&FilterDeny{
					NewANDCondition(
						&ConditionActionType{ACTION_ATTACK},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_SOURCE_CREATURE}, Accessor{FILTER_SOURCE_PARAMETER_PROPERTY, ACCESSOR_MODE_PROPERTY_OWNER}),
						NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_USED}, 1)),
					NewANDCondition(
						&ConditionActionType{ACTION_REMOVE_PROPERTY},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY))})})
	g.Filters = append(g.Filters,
		&FilterDeny{
			NewANDCondition(
				&ConditionActionType{ACTION_ATTACK},
				NewORCondition(
					NewConditionEqual(FILTER_SOURCE_PARAMETER_SOURCE_CREATURE, FILTER_SOURCE_PARAMETER_TARGET_CREATURE),
					&ConditionActionDenied{NewActionGainAdditionalFood(FILTER_SOURCE_PARAMETER_SOURCE_CREATURE, 2)})),	
			nil})
	g.Filters = append(g.Filters, 
		&FilterAction{	
			FILTER_ACTION_EXECUTE_AFTER,
			&ConditionActionType{ACTION_ATTACK},
			nil,
			NewActionAddFilters(
				&FilterDeny{
					NewORCondition(	
						&ConditionActionType{ACTION_BURN_FAT},
						&ConditionActionType{ACTION_ATTACK},
						&ConditionActionType{ACTION_GET_FOOD_FROM_BANK}),
					&ConditionActionType{ACTION_START_TURN}},
				&FilterDeny{
					NewANDCondition(
						&ConditionActionType{ACTION_ATTACK},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_SOURCE_CREATURE}, FILTER_SOURCE_PARAMETER_SOURCE_CREATURE)),
					&ConditionActionType{ACTION_NEW_PHASE}})})
				
	//Hibernation
	g.AddFilter(&FilterAction{
		FILTER_ACTION_EXECUTE_AFTER,
		NewANDCondition(
				&ConditionActionType{ACTION_ADD_SINGLE_PROPERTY},
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_HIBERNATION}, 1)),
		nil,
		NewActionAddFilters(
			NewFilterAllow(	
				&ConditionPhase{PHASE_FEEDING},
				nil,
				NewActionHibernate(Accessor{FILTER_SOURCE_PARAMETER_PROPERTY, ACCESSOR_MODE_PROPERTY_OWNER})),
			&FilterDeny{
				NewANDCondition(
					&ConditionActionType{ACTION_HIBERNATE},
					NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_CREATURE}, Accessor{FILTER_SOURCE_PARAMETER_PROPERTY, ACCESSOR_MODE_PROPERTY_OWNER}),
					&NOTCondition{NewConditionEqual(InstantiationOff{TraitsCount{InstantiationOn{FILTER_SOURCE_PARAMETER_PROPERTY}, TRAIT_USED}}, 0)}),
				NewANDCondition(
						&ConditionActionType{ACTION_REMOVE_PROPERTY},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY))},
			&FilterAction{
				FILTER_ACTION_EXECUTE_BEFORE,
				NewANDCondition(
					&ConditionActionType{ACTION_NEW_PHASE},
					NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PHASE}, PHASE_FEEDING),
					&NOTCondition{NewConditionEqual(InstantiationOff{TraitsCount{InstantiationOn{FILTER_SOURCE_PARAMETER_PROPERTY},  TRAIT_USED}}, 0)}),
				NewANDCondition(
						&ConditionActionType{ACTION_REMOVE_PROPERTY},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY)),
				NewActionRemoveTrait(FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_USED)},
			&FilterAction{
				FILTER_ACTION_EXECUTE_AFTER,
				NewANDCondition(
					&ConditionActionType{ACTION_HIBERNATE},
					NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_CREATURE}, Accessor{FILTER_SOURCE_PARAMETER_PROPERTY, ACCESSOR_MODE_PROPERTY_OWNER})),
				NewANDCondition(
						&ConditionActionType{ACTION_REMOVE_PROPERTY},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY)),
				NewActionSequence(
					NewActionAddTrait(InstantiationOff{FILTER_SOURCE_PARAMETER_CREATURE}, TRAIT_FED),
					NewActionAddTrait(FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_USED),
					NewActionAddTrait(FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_USED),
					NewActionAddFilters(
						&FilterDeny{
							NewANDCondition(
								NewORCondition(&ConditionActionType{ACTION_GAIN_FOOD}, &ConditionActionType{ACTION_GAIN_ADDITIONAL_FOOD}),
								NewConditionEqual(InstantiationOff{InstantiationOff{InstantiationOff{FILTER_SOURCE_PARAMETER_CREATURE}}}, InstantiationOff{FILTER_SOURCE_PARAMETER_CREATURE})),
							NewANDCondition(
								&ConditionActionType{ACTION_NEW_PHASE},
								NewConditionEqual(InstantiationOff{InstantiationOff{InstantiationOff{FILTER_SOURCE_PARAMETER_PHASE}}}, PHASE_FEEDING))}))})})
	g.AddFilter(&FilterDeny{
		NewANDCondition(
			&ConditionActionType{ACTION_HIBERNATE},
			NewConditionEqual(FILTER_SOURCE_PARAMETER_BANK_CARDS_COUNT, 0)),
		nil})
	
	//Poisonous
	g.AddFilter(&FilterAction{
		FILTER_ACTION_EXECUTE_AFTER,
		NewANDCondition(
				&ConditionActionType{ACTION_ADD_SINGLE_PROPERTY},
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_POISONOUS}, 1)),
		nil,
		NewActionAddFilters(
			&FilterAction{
				FILTER_ACTION_EXECUTE_BEFORE,
				NewANDCondition(
					&ConditionActionType{ACTION_ATTACK},
					NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_TARGET_CREATURE}, Accessor{FILTER_SOURCE_PARAMETER_PROPERTY, ACCESSOR_MODE_PROPERTY_OWNER})),
				NewANDCondition(
						&ConditionActionType{ACTION_REMOVE_PROPERTY},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY)),
				NewActionAddFilters(
					&FilterAction{
						FILTER_ACTION_EXECUTE_BEFORE,
						NewANDCondition(
							&ConditionActionType{ACTION_NEW_PHASE},
							NewConditionEqual(InstantiationOff{InstantiationOff{FILTER_SOURCE_PARAMETER_PHASE}}, PHASE_EXTINCTION)),
						&ConditionActionType{ACTION_NEW_PHASE},
						NewActionRemoveCreature(InstantiationOff{FILTER_SOURCE_PARAMETER_SOURCE_CREATURE})})})})
	
	//Communication
	g.AddFilter(&FilterAction{
		FILTER_ACTION_EXECUTE_AFTER,
		NewANDCondition(
				&ConditionActionType{ACTION_ADD_PAIR_PROPERTY},
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_COMMUNICATION}, 1)),
		nil,
		NewActionAddFilters(
			&FilterAction{
				FILTER_ACTION_EXECUTE_AFTER,
				NewANDCondition(
					&ConditionActionType{ACTION_GET_FOOD_FROM_BANK},
					NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_CREATURE}, FILTER_SOURCE_PARAMETER_LEFT_CREATURE),
					NewConditionEqual(InstantiationOff{TraitsCount{InstantiationOn{FILTER_SOURCE_PARAMETER_PROPERTY}, TRAIT_USED}}, 0)),
				NewANDCondition(
					&ConditionActionType{ACTION_REMOVE_PROPERTY},
					NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY)),
				NewActionAddTrait(FILTER_SOURCE_PARAMETER_RIGHT_CREATURE, TRAIT_ADDITIONAL_GET_FOOD_FROM_BANK)},
			&FilterAction{
				FILTER_ACTION_EXECUTE_AFTER,
				NewANDCondition(
					&ConditionActionType{ACTION_GET_FOOD_FROM_BANK},
					NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_CREATURE}, FILTER_SOURCE_PARAMETER_RIGHT_CREATURE),
					NewConditionEqual(InstantiationOff{TraitsCount{InstantiationOn{FILTER_SOURCE_PARAMETER_PROPERTY}, TRAIT_USED}}, 0)),
				NewANDCondition(
					&ConditionActionType{ACTION_REMOVE_PROPERTY},
					NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY)),
				NewActionAddTrait(FILTER_SOURCE_PARAMETER_LEFT_CREATURE, TRAIT_ADDITIONAL_GET_FOOD_FROM_BANK)},
			&FilterAction{
				FILTER_ACTION_EXECUTE_AFTER,
				NewANDCondition(
					NewORCondition(
						NewANDCondition(
							&ConditionActionType{ACTION_GET_FOOD_FROM_BANK},
							NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_CREATURE}, FILTER_SOURCE_PARAMETER_LEFT_CREATURE)),
						&ConditionActionType{ACTION_START_TURN}),
					NewConditionEqual(InstantiationOff{TraitsCount{InstantiationOn{FILTER_SOURCE_PARAMETER_LEFT_CREATURE}, TRAIT_ADDITIONAL_GET_FOOD_FROM_BANK}}, 1)),
				NewANDCondition(
					&ConditionActionType{ACTION_REMOVE_PROPERTY},
					NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY)),
				NewActionRemoveTrait(FILTER_SOURCE_PARAMETER_LEFT_CREATURE, TRAIT_ADDITIONAL_GET_FOOD_FROM_BANK)},
			&FilterAction{
				FILTER_ACTION_EXECUTE_AFTER,
				NewANDCondition(
					NewORCondition(
						NewANDCondition(
							&ConditionActionType{ACTION_GET_FOOD_FROM_BANK},
							NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_CREATURE}, FILTER_SOURCE_PARAMETER_RIGHT_CREATURE)),
						&ConditionActionType{ACTION_START_TURN}),
					NewConditionEqual(InstantiationOff{TraitsCount{InstantiationOn{FILTER_SOURCE_PARAMETER_RIGHT_CREATURE}, TRAIT_ADDITIONAL_GET_FOOD_FROM_BANK}}, 1)),
				NewANDCondition(
					&ConditionActionType{ACTION_REMOVE_PROPERTY},
					NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY)),
				NewActionRemoveTrait(FILTER_SOURCE_PARAMETER_RIGHT_CREATURE, TRAIT_ADDITIONAL_GET_FOOD_FROM_BANK)},
			&FilterAction{
				FILTER_ACTION_EXECUTE_AFTER,
				NewANDCondition(
					&ConditionActionType{ACTION_GET_FOOD_FROM_BANK},
					NewConditionEqual(InstantiationOff{TraitsCount{InstantiationOn{FILTER_SOURCE_PARAMETER_PROPERTY}, TRAIT_USED}}, 0),
					NewORCondition(
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_CREATURE}, FILTER_SOURCE_PARAMETER_LEFT_CREATURE),
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_CREATURE}, FILTER_SOURCE_PARAMETER_RIGHT_CREATURE))),
				NewANDCondition(
					&ConditionActionType{ACTION_REMOVE_PROPERTY},
					NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY)),
				NewActionSequence(
					NewActionAddTrait(FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_USED),
					NewActionAddFilters(
						&FilterAction{
							FILTER_ACTION_EXECUTE_BEFORE,
							&ConditionActionType{ACTION_START_TURN},
							&ConditionActionType{ACTION_START_TURN},
							NewActionRemoveTrait(FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_USED)}))})})
							
	//Swimming
	
	g.AddFilter(&FilterDeny{
		NewANDCondition(
			&ConditionActionType{ACTION_ATTACK},
			NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_SOURCE_CREATURE, TRAIT_SWIMMING}, 1),
			NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_TARGET_CREATURE, TRAIT_SWIMMING}, 0)),
		nil})
		
	g.AddFilter(&FilterDeny{
		NewANDCondition(
			&ConditionActionType{ACTION_ATTACK},
			NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_SOURCE_CREATURE, TRAIT_SWIMMING}, 0),
			NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_TARGET_CREATURE, TRAIT_SWIMMING}, 1)),
		nil})
		
	//Cooperation
	
	g.AddFilter(&FilterAction{
		FILTER_ACTION_EXECUTE_AFTER,
		NewANDCondition(
				&ConditionActionType{ACTION_ADD_PAIR_PROPERTY},
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_COOPERATION}, 1)),
		nil,
		NewActionAddFilters(
			&FilterAction{
				FILTER_ACTION_EXECUTE_AFTER,
				NewANDCondition(
					NewORCondition(
						&ConditionActionType{ACTION_GAIN_FOOD},
						&ConditionActionType{ACTION_GAIN_ADDITIONAL_FOOD}),
					NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_CREATURE}, FILTER_SOURCE_PARAMETER_LEFT_CREATURE),
					NewConditionEqual(InstantiationOff{TraitsCount{InstantiationOn{FILTER_SOURCE_PARAMETER_PROPERTY}, TRAIT_USED}}, 0),
					&NOTCondition{&ConditionActionDenied{NewActionGainAdditionalFood(FILTER_SOURCE_PARAMETER_RIGHT_CREATURE, 1)}}),
				NewANDCondition(
					&ConditionActionType{ACTION_REMOVE_PROPERTY},
					NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY)),
				NewActionSequence(
					NewActionAddTrait(FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_USED),
					NewActionGainAdditionalFood(FILTER_SOURCE_PARAMETER_RIGHT_CREATURE, 1))},
			&FilterAction{
				FILTER_ACTION_EXECUTE_AFTER,
				NewANDCondition(
					NewORCondition(
						&ConditionActionType{ACTION_GAIN_FOOD},
						&ConditionActionType{ACTION_GAIN_ADDITIONAL_FOOD}),
					NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_CREATURE}, FILTER_SOURCE_PARAMETER_RIGHT_CREATURE),
					NewConditionEqual(InstantiationOff{TraitsCount{InstantiationOn{FILTER_SOURCE_PARAMETER_PROPERTY}, TRAIT_USED}}, 0),
					&NOTCondition{&ConditionActionDenied{NewActionGainAdditionalFood(FILTER_SOURCE_PARAMETER_LEFT_CREATURE, 1)}}),
				NewANDCondition(
					&ConditionActionType{ACTION_REMOVE_PROPERTY},
					NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY)),
				NewActionSequence(
					NewActionAddTrait(FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_USED),
					NewActionGainAdditionalFood(FILTER_SOURCE_PARAMETER_LEFT_CREATURE, 1))},
			&FilterAction{
				FILTER_ACTION_EXECUTE_AFTER,
				NewANDCondition(
					NewORCondition(
						&ConditionActionType{ACTION_GAIN_FOOD},
						&ConditionActionType{ACTION_GAIN_ADDITIONAL_FOOD}),
					NewConditionEqual(InstantiationOff{TraitsCount{InstantiationOn{FILTER_SOURCE_PARAMETER_PROPERTY}, TRAIT_USED}}, 0),
					NewORCondition(
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_CREATURE}, FILTER_SOURCE_PARAMETER_LEFT_CREATURE),
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_CREATURE}, FILTER_SOURCE_PARAMETER_RIGHT_CREATURE))),
				NewANDCondition(
					&ConditionActionType{ACTION_REMOVE_PROPERTY},
					NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY)),
				NewActionAddFilters(
					&FilterAction{
						FILTER_ACTION_EXECUTE_BEFORE,
						&ConditionActionType{ACTION_START_TURN},
						&ConditionActionType{ACTION_START_TURN},
						NewActionRemoveTrait(FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_USED)})})})
							
	// Mimicry
	g.Filters = append(g.Filters,
		&FilterAction {
			FILTER_ACTION_EXECUTE_AFTER,
			NewANDCondition(
				&ConditionActionType{ACTION_ADD_SINGLE_PROPERTY},
				NewConditionEqual(TraitsCount{FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_MIMICRY}, 1)),
			nil,
			NewActionAddFilters(
				&FilterAction{
					FILTER_ACTION_REPLACE,
					NewANDCondition(
						&ConditionActionType{ACTION_ATTACK},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_TARGET_CREATURE}, Accessor{FILTER_SOURCE_PARAMETER_PROPERTY, ACCESSOR_MODE_PROPERTY_OWNER}),
						NewConditionEqual(InstantiationOff{TraitsCount{InstantiationOn{FILTER_SOURCE_PARAMETER_PROPERTY}, TRAIT_USED}}, 0)),
					NewANDCondition(
						&ConditionActionType{ACTION_REMOVE_PROPERTY},
						NewConditionEqual(InstantiationOff{FILTER_SOURCE_PARAMETER_PROPERTY}, FILTER_SOURCE_PARAMETER_PROPERTY)),
					NewActionSequence(
						NewActionAddTrait(FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_USED),
						NewActionAddFilters(&FilterDeny{
							NewANDCondition(
								&ConditionActionType{ACTION_ATTACK},
								NewConditionEqual(
									InstantiationOff{InstantiationOff{InstantiationOff{FILTER_SOURCE_PARAMETER_TARGET_CREATURE}}}, 
									InstantiationOff{FILTER_SOURCE_PARAMETER_TARGET_CREATURE})),
							NewORCondition(
								&ConditionActionDenied{NewActionAttack(
									InstantiationOff{Accessor{FILTER_SOURCE_PARAMETER_TARGET_CREATURE, ACCESSOR_MODE_CREATURE_OWNER}}, 
										InstantiationOff{FILTER_SOURCE_PARAMETER_SOURCE_CREATURE},
											InstantiationOff{
												InstantiationOff{
														Accessor{
															InstantiationOn{InstantiationOn{InstantiationOff{Accessor{
																FILTER_SOURCE_PARAMETER_TARGET_CREATURE, 
																ACCESSOR_MODE_CREATURE_OWNER}}}},
														ACCESSOR_MODE_CREATURES}}})},
								&ConditionActionType{ACTION_START_TURN})}),
						NewActionAddFilters(
							&FilterAction{
								FILTER_ACTION_EXECUTE_BEFORE,
								&ConditionActionType{ACTION_NEW_PHASE},
								&ConditionActionType{ACTION_NEW_PHASE},
								NewActionRemoveTrait(FILTER_SOURCE_PARAMETER_PROPERTY, TRAIT_USED)}),
						NewActionAttack(
							InstantiationOff{Accessor{FILTER_SOURCE_PARAMETER_TARGET_CREATURE, ACCESSOR_MODE_CREATURE_OWNER}}, 
							InstantiationOff{FILTER_SOURCE_PARAMETER_SOURCE_CREATURE},
							InstantiationOff{
								InstantiationOff{
									Accessor{
										InstantiationOn{InstantiationOn{InstantiationOff{Accessor{
											FILTER_SOURCE_PARAMETER_TARGET_CREATURE, 
											ACCESSOR_MODE_CREATURE_OWNER}}}},
										ACCESSOR_MODE_CREATURES}}}))})})
}