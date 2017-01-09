package main

import "testing"
import "github.com/stretchr/testify/assert"

func TestCodegenComment(t *testing.T) {
	tt := []struct {
		in  []string
		out string
	}{
		{[]string{"file_1.avsc", "file_2.avsc"}, `/*
 * CODE GENERATED AUTOMATICALLY WITH github.com/alanctgardner/gogen-avro
 * THIS FILE SHOULD NOT BE EDITED BY HAND
 *
 * SOURCES:
 *     file_1.avsc
 *     file_2.avsc
 */`},
		{[]string{"file_1.avsc"}, `/*
 * CODE GENERATED AUTOMATICALLY WITH github.com/alanctgardner/gogen-avro
 * THIS FILE SHOULD NOT BE EDITED BY HAND
 *
 * SOURCE:
 *     file_1.avsc
 */`},
	}

	for _, tc := range tt {
		assert.Equal(t, tc.out, codegenComment(tc.in))
	}
}
