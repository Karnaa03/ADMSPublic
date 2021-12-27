package model

import "testing"

func TestShelf_GetFsmGraph(t *testing.T) {
	shelf := Shelf{}
	t.Log(shelf.GetFsmGraph())
}

func Test_equalDiffOrder(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "same", args: struct {
			a []string
			b []string
		}{a: []string{"a", "a"}, b: []string{"a", "a"}}, want: true},
		{name: "same order", args: struct {
			a []string
			b []string
		}{a: []string{"a", "b"}, b: []string{"a", "b"}}, want: true},
		{name: "different order", args: struct {
			a []string
			b []string
		}{a: []string{"a", "b"}, b: []string{"b", "a"}}, want: true},
		{name: "with unknown", args: struct {
			a []string
			b []string
		}{a: []string{"a", "b", "c"}, b: []string{"b", "a"}}, want: false},
		{name: "with less a", args: struct {
			a []string
			b []string
		}{a: []string{"a"}, b: []string{"b", "a"}}, want: true},
		{name: "with less b", args: struct {
			a []string
			b []string
		}{a: []string{"b"}, b: []string{"b", "a"}}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := in(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("in() = %v, want %v", got, tt.want)
			}
		})
	}
}
