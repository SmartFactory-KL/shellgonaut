package traverse

import (
	"strings"

	"github.com/aas-core-works/aas-core3.1-golang/types"
)

// HasSemanticID checks if semanticID is used as the value for a key within
// SemanticID() or SupplementalSemanticIDs() of target - without checking the key type.
func HasSemanticID(target types.IHasSemantics, semanticID string) bool {
	if target == nil {
		return false
	}

	targetSemantic := target.SemanticID()

	if targetSemantic == nil {
		return false
	}

	keys := targetSemantic.Keys()
	if DoKeysContainSemanticValue(keys, semanticID) {
		return true
	}

	supplementalSemantics := target.SupplementalSemanticIDs()
	for _, supplementalEntry := range supplementalSemantics {
		if supplementalEntry == nil {
			continue
		}

		if DoKeysContainSemanticValue(supplementalEntry.Keys(), semanticID) {
			return true
		}
	}

	return false
}

// DoKeysContainSemanticValue checks if semanticID is used as the value of any key from keys (keyTypes are ignored)
func DoKeysContainSemanticValue(keys []types.IKey, semanticID string) bool {
	// not sure if this is even possible
	if len(strings.TrimSpace(semanticID)) == 0 {
		return false
	}

	// no keys -> no match
	if len(keys) == 0 {
		return false
	}

	for _, keyEntry := range keys {
		if keyEntry == nil {
			continue
		}

		if strings.EqualFold(keyEntry.Value(), semanticID) {
			return true
		}
	}

	return false
}
