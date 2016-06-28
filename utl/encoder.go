// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"encoding/gob"
	"encoding/json"
	goio "io"
)

// Encoder defines encoders; e.g. gob or json
type Encoder interface {
	Encode(e interface{}) error
}

// Decoder defines decoders; e.g. gob or json
type Decoder interface {
	Decode(e interface{}) error
}

// GetEncoder returns a new encoder
func GetEncoder(w goio.Writer, enctype string) Encoder {
	if enctype == "json" {
		return json.NewEncoder(w)
	}
	return gob.NewEncoder(w)
}

// GetDecoder returns a new decoder
func GetDecoder(r goio.Reader, enctype string) Decoder {
	if enctype == "json" {
		return json.NewDecoder(r)
	}
	return gob.NewDecoder(r)
}
