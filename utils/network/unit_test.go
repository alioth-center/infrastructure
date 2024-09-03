package network

import "testing"

func TestIsPrivateIP(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "private ip 1",
			args: args{ip: "10.0.0.1"},
			want: true,
		},
		{
			name: "private ip 2",
			args: args{ip: "172.16.0.1"},
			want: true,
		},
		{
			name: "public ip",
			args: args{ip: "8.8.8.8"},
			want: false,
		},
		{
			name: "invalid ip",
			args: args{ip: "999.999.999.999"},
			want: false,
		},
		{
			name: "empty ip",
			args: args{ip: ""},
			want: false,
		},
		{
			name: "ipv6 private ip",
			args: args{ip: "fd00::"},
			want: true,
		},
		{
			name: "ipv6 public ip",
			args: args{ip: "2001:4860:4860::8888"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPrivateIP(tt.args.ip); got != tt.want {
				t.Errorf("IsPrivateIP() = %v, want %v", got, tt.want)
			}
		})
	}
}
