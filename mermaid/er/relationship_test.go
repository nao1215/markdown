package er

import "testing"

func TestRelationship_string(t *testing.T) {
	t.Parallel()

	type args struct {
		lr bool
	}
	tests := []struct {
		name string
		r    Relationship
		args args
		want string
	}{
		{
			name: "ZeroToOneRelationship, left",
			r:    ZeroToOneRelationship,
			args: args{lr: left},
			want: "|o",
		},
		{
			name: "ZeroToOneRelationship, right",
			r:    ZeroToOneRelationship,
			args: args{lr: right},
			want: "o|",
		},
		{
			name: "ExactlyOneRelationship, left",
			r:    ExactlyOneRelationship,
			args: args{lr: left},
			want: "||",
		},
		{
			name: "ExactlyOneRelationship, right",
			r:    ExactlyOneRelationship,
			args: args{lr: right},
			want: "||",
		},
		{
			name: "ZeroToMoreRelationship, left",
			r:    ZeroToMoreRelationship,
			args: args{lr: left},
			want: "}o",
		},
		{
			name: "ZeroToMoreRelationship, right",
			r:    ZeroToMoreRelationship,
			args: args{lr: right},
			want: "o{",
		},
		{
			name: "OneToMoreRelationship, left",
			r:    OneToMoreRelationship,
			args: args{lr: left},
			want: "}|",
		},
		{
			name: "OneToMoreRelationship, right",
			r:    OneToMoreRelationship,
			args: args{lr: right},
			want: "|{",
		},
		{
			name: "default",
			r:    "default",
			args: args{lr: left},
			want: "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.r.string(tt.args.lr); got != tt.want {
				t.Errorf("Relationship.string() = %v, want %v", got, tt.want)
			}
		})
	}
}
