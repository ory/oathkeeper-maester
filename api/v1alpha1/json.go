// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"bytes"
	"encoding/json"
)

func unescapedMarshalIndent(in interface{}, prefix, indent string) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	enc.SetIndent(prefix, indent)
	if err := enc.Encode(in); err != nil {
		return nil, err
	}

	result := b.Bytes()
	return result[:len(result)-1], nil
}

func unescapedMarshal(in interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(in); err != nil {
		return nil, err
	}

	result := b.Bytes()
	return result[:len(result)-1], nil
}
