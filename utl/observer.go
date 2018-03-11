// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

// Observer is an interface to objects that need to observe something
type Observer interface {
	Name() string // returns the unique name of this observer (used for initialization only)
	Update()      // the data observed by this observer is being update
}

// Observable indicates that an object is observable; i.e. it has a list of interested observers
type Observable struct {
	observers []Observer     // list of interested parties
	nameToObs map[string]int // maps name into observers list; used only during initialization for efficiency
}

// AddObserver adds an object to the list of interested observers
func (o *Observable) AddObserver(obs Observer) {
	if o.nameToObs == nil {
		o.nameToObs = make(map[string]int)
	}
	name := obs.Name()
	if _, ok := o.nameToObs[name]; ok {
		return // skip
	}
	o.nameToObs[name] = len(o.observers)
	o.observers = append(o.observers, obs)
}

// NotifyUpdate notifies observers of updates
func (o *Observable) NotifyUpdate() {
	for _, obs := range o.observers {
		obs.Update()
	}
}
