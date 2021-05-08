package generator

import "testing"

func TestToPublicName(t *testing.T) {
	testCases := []struct {
		name     string
		expected string
	}{
		{
			name:     "IntField",
			expected: "IntField",
		},
		{
			name:     "_intField_",
			expected: "IntField",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			name := ToPublicName(tc.name)
			if name != tc.expected {
				t.Errorf(`ToPublicName("%s") Expected %v, got %v`, tc.name, tc.expected, name)
			}
		})
	}
}

func TestToSnake(t *testing.T) {
	testCases := []struct {
		in  string
		out string
	}{
		{"a", "a"},
		{"snake", "snake"},
		{"A", "a"},
		{"ID", "id"},
		{"MOTD", "motd"},
		{"Snake", "snake"},
		{"SnakeTest", "snake_test"},
		{"SnakeID", "snake_id"},
		{"SnakeIDGoogle", "snake_id_google"},
		{"LinuxMOTD", "linux_motd"},
		{"OMGWTFBBQ", "omgwtfbbq"},
		{"omg_wtf_bbq", "omg_wtf_bbq"},
	}

	for _, tc := range testCases {
		t.Run(tc.in, func(t *testing.T) {
			output := ToSnake(tc.in)
			if output != tc.out {
				t.Errorf(`ToSnake("%s"), wanted "%s", got \%s"`, tc.in, tc.out, output)
			}
		})
	}
}
