package ip

import "testing"

func TestIsInternal(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Localhost",
			args: args{
				ip: "127.0.0.1",
			},
			want: true,
		},
		{
			name: "Some private",
			args: args{
				ip: "10.231.0.11",
			},
			want: true,
		},
		{
			name: "bad ip",
			args: args{
				ip: "someip",
			},
			want: false,
		},
		{
			name: "external ip",
			args: args{
				ip: "	210.130.169.196",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInternal(tt.args.ip); got != tt.want {
				t.Errorf("IsInternal() = %v, want %v", got, tt.want)
			}
		})
	}
}
