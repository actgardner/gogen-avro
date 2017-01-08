package main

import "testing"

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

	for i, tc := range tt {
		got := codegenComment(tc.in)
		expected := tc.out

		if got != expected {
			t.Errorf("#%d: got\n%sexpected\n%s", i, got, expected)
		}
	}

}
