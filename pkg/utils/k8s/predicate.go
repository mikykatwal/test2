/*
Copyright 2022 Mondoo, Inc.

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package k8s

import (
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var _ predicate.Predicate = IgnoreGenericEventsPredicate{}

// CreateOrDeletePredicate will completely ignore generic events.
type IgnoreGenericEventsPredicate struct{}

func (p IgnoreGenericEventsPredicate) Create(e event.CreateEvent) bool {
	return true
}

func (p IgnoreGenericEventsPredicate) Update(e event.UpdateEvent) bool {
	return true
}

func (p IgnoreGenericEventsPredicate) Delete(e event.DeleteEvent) bool {
	return true
}

func (p IgnoreGenericEventsPredicate) Generic(e event.GenericEvent) bool {
	return false
}

var _ predicate.Predicate = CreateUpdateEventsPredicate{}

// CreateUpdateEventsPredicate will allow only create and update events.
// Update events caused by object deletion are also ignored.
type CreateUpdateEventsPredicate struct{}

func (p CreateUpdateEventsPredicate) Create(e event.CreateEvent) bool {
	return true
}

func (p CreateUpdateEventsPredicate) Update(e event.UpdateEvent) bool {
	// If the deletion timestamp is set, the object is being deleted so we
	// can ignore the event.
	return e.ObjectNew.GetDeletionTimestamp() == nil
}

func (p CreateUpdateEventsPredicate) Delete(e event.DeleteEvent) bool {
	return false
}

func (p CreateUpdateEventsPredicate) Generic(e event.GenericEvent) bool {
	return false
}
