package global

import (
	"github.com/SSripilaipong/muto/core/mutation/rule/mutator"
)

func newForeignNormalMutators() []mutator.NamedUnit {
	return []mutator.NamedUnit{
		stringMutator,
		addMutator,
		subtractMutator,
		stringToRunesMutator,
		parseRunesToStringMutator,
		concatMutator,
		isStringMutator,
		isNumberMutator,
		isBooleanMutator,
		multiplyMutator,
		divideMutator,
		divIntegerMutator,
		modIntegerMutator,
		equalMutator,
		notEqualMutator,
		greaterThanMutator,
		greaterThanOrEqualMutator,
		lessThanMutator,
		lessThanOrEqualMutator,
		//cliInputMutator(cliReader),
		andMutator,
		orMutator,
		notMutator,
		newTryMutator(),
		newPortalMutator(),
		newSelectMutator(),
	}
}
