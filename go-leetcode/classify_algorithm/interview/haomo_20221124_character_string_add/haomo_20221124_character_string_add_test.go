package haomo_20221124_character_string_add

import "testing"

func TestAdd(t *testing.T) {
	type args struct {
		s1 string
		s2 string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "用例1",
			args: args{
				s1: "123",
				s2: "345",
			},
			want: "468",
		},
		{
			name: "用例2",
			args: args{
				s1: "876",
				s2: "128",
			},
			want: "1004",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Add(tt.args.s1, tt.args.s2); got != tt.want {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}
