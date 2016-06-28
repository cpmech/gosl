// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msh

// Find finds an edge
//  Returns nil if not found
func (o EdgesMap) Find(key EdgeKey) *Edge {
	return nil
}

// Functions to sort VertSet
func (o VertSet) Len() int           { return len(o) }
func (o VertSet) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o VertSet) Less(i, j int) bool { return o[i].Id < o[j].Id }

// Functions to sort CellSet
func (o CellSet) Len() int           { return len(o) }
func (o CellSet) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o CellSet) Less(i, j int) bool { return o[i].Id < o[j].Id }
