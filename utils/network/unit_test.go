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

func TestIsValidIP(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid ip 1",
			args: args{ip: "1.1.1.1"},
			want: true,
		},
		{
			name: "valid ip 2",
			args: args{ip: "192.168.1.1"},
			want: true,
		},
		{
			name: "valid ipv6",
			args: args{ip: "2001:4860:4860::8888"},
			want: true,
		},
		{
			name: "invalid ip",
			args: args{ip: "999.999.999.999"},
			want: false,
		},
		{
			name: "invalid ipv6",
			args: args{ip: "2001:4860:4860::8888:"},
			want: false,
		},
		{
			name: "empty ip",
			args: args{ip: ""},
			want: false,
		},
		{
			name: "ip with mask",
			args: args{ip: "192.168.1.10/24"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidIP(tt.args.ip); got != tt.want {
				t.Errorf("IsValidIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidIPOrCIDR(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		want bool
	}{
		{
			name: "valid ip 1",
			ip:   "1.1.1.1",
			want: true,
		},
		{
			name: "valid ip 2",
			ip:   "192.168.1.1",
			want: true,
		},
		{
			name: "valid ipv6",
			ip:   "2001:4860:4860::8888",
			want: true,
		},
		{
			name: "invalid ip",
			ip:   "999.999.999.999",
			want: false,
		},
		{
			name: "invalid ipv6",
			ip:   "2001:4860:4860::8888:",
			want: false,
		},
		{
			name: "empty ip",
			ip:   "",
			want: false,
		},
		{
			name: "cidr valid",
			ip:   "192.168.1.10/24",
			want: true,
		},
		{
			name: "cidr invalid",
			ip:   "192.168.1.10/33",
			want: false,
		},
		{
			name: "ipv6 cidr valid",
			ip:   "2001:4860:4860::8888/64",
			want: true,
		},
		{
			name: "ipv6 cidr invalid",
			ip:   "2001:4860:4860::8888/129",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidIPOrCIDR(tt.ip); got != tt.want {
				t.Errorf("IsValidIPOrCIDR() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPInCIDR(t *testing.T) {
	type args struct {
		ip   string
		cidr string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ip in cidr",
			args: args{cidr: "192.168.1.0/24", ip: "192.168.1.25"},
			want: true,
		},
		{
			name: "ip not in cidr",
			args: args{cidr: "192.168.1.0/24", ip: "192.168.2.1"},
			want: false,
		},
		{
			name: "ipv6 ip in cidr",
			args: args{cidr: "2001:db8::/32", ip: "2001:db8:0:1::1"},
			want: true,
		},
		{
			name: "ipv6 ip not in cidr",
			args: args{cidr: "2001:db8::/32", ip: "2001:db9::1"},
			want: false,
		},
		{
			name: "another ipv6 ip in cidr",
			args: args{cidr: "fd00::/8", ip: "fd12:3456:789a:1::1"},
			want: true,
		},
		{
			name: "another ipv6 ip not in cidr",
			args: args{cidr: "fd00::/8", ip: "fe80::1"},
			want: false,
		},
		{
			name: "invalid ip",
			args: args{cidr: "192.168.1.0/24", ip: "192.168.1.256"},
			want: false,
		},
		{
			name: "invalid cidr",
			args: args{cidr: "192.168.1.0/64", ip: "192.168.1.1"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IPInCIDR(tt.args.ip, tt.args.cidr); got != tt.want {
				t.Errorf("isIPInCIDR() = %v, want %v", got, tt.want)
			}
		})
	}
}
