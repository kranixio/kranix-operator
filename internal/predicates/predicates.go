package predicates

import (
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// ResourceGenerationChangedPredicate filters events based on generation changes
func ResourceGenerationChangedPredicate() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			if e.ObjectOld == nil || e.ObjectNew == nil {
				return false
			}
			return e.ObjectNew.GetGeneration() != e.ObjectOld.GetGeneration()
		},
	}
}

// ResourceChangedOrDeletedPredicate filters for create, update, or delete events
func ResourceChangedOrDeletedPredicate() predicate.Predicate {
	return predicate.Funcs{
		CreateFunc:  func(e event.CreateEvent) bool { return true },
		UpdateFunc:  func(e event.UpdateEvent) bool { return true },
		DeleteFunc:  func(e event.DeleteEvent) bool { return true },
		GenericFunc: func(e event.GenericEvent) bool { return false },
	}
}

// AnnotationChangedPredicate filters events based on annotation changes
func AnnotationChangedPredicate() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			if e.ObjectOld == nil || e.ObjectNew == nil {
				return false
			}
			oldAnnotations := e.ObjectOld.GetAnnotations()
			newAnnotations := e.ObjectNew.GetAnnotations()
			
			if len(oldAnnotations) != len(newAnnotations) {
				return true
			}
			
			for k, v := range newAnnotations {
				if oldAnnotations[k] != v {
					return true
				}
			}
			return false
		},
	}
}
