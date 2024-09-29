package lianying_20221209_graph_cycle

import "testing"

func Test_hasCycle(t *testing.T) {
	type args struct {
		graph [][]int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "用例-1",
			args: args{
				// (0,1)-->(1,2)-->(2,0)
				graph: [][]int{
					{0, 1, 0},
					{0, 0, 1},
					{1, 0, 0},
				},
			},
			want: true,
		},
		{
			name: "用例-2",
			args: args{
				// (0,1)-->(1,2)
				graph: [][]int{
					{0, 1, 0},
					{0, 0, 1},
					{0, 0, 0},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasCycle(tt.args.graph); got != tt.want {
				t.Errorf("hasCycle() = %v, want %v", got, tt.want)
			}
		})
	}
}
