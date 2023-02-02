package mid

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWithDomains(t *testing.T) {
	type args struct {
		starter CORSOptions
		domains []string
		headers []string
	}

	tests := []struct {
		name string
		args args
		want CORSOptions
	}{
		// Nil starters.
		{
			name: "nil options lists, nil starter",
			args: args{
				starter: CORSOptions{
					domains:           nil,
					additionalHeaders: nil,
				},
				domains: nil,
				headers: nil,
			},
			want: CORSOptions{
				domains:           nil,
				additionalHeaders: nil,
			},
		},
		{
			name: "empty options lists, nil starter",
			args: args{
				starter: CORSOptions{
					domains:           nil,
					additionalHeaders: nil,
				},
				domains: []string{},
				headers: []string{},
			},
			want: CORSOptions{
				domains:           nil,
				additionalHeaders: nil,
			},
		},
		{
			name: "not empty lists, nil starter",
			args: args{
				starter: CORSOptions{
					domains:           nil,
					additionalHeaders: nil,
				},
				domains: []string{
					"newitem.xyz",
				},
				headers: []string{
					"New-Header",
				},
			},
			want: CORSOptions{
				domains: []string{
					"newitem.xyz",
				},
				additionalHeaders: []string{
					"New-Header",
				},
			},
		},

		// Empty starters.
		{
			name: "nil options lists, empty starter",
			args: args{
				starter: CORSOptions{
					domains:           []string{},
					additionalHeaders: []string{},
				},
				domains: nil,
				headers: nil,
			},
			want: CORSOptions{
				domains:           []string{},
				additionalHeaders: []string{},
			},
		},
		{
			name: "empty options lists, empty starter",
			args: args{
				starter: CORSOptions{
					domains:           []string{},
					additionalHeaders: []string{},
				},
				domains: []string{},
				headers: []string{},
			},
			want: CORSOptions{
				domains:           []string{},
				additionalHeaders: []string{},
			},
		},
		{
			name: "not empty options list, empty starter",
			args: args{
				starter: CORSOptions{
					domains:           []string{},
					additionalHeaders: []string{},
				},
				domains: []string{
					"newitem.xyz",
				},
				headers: []string{
					"New-Header",
				},
			},
			want: CORSOptions{
				domains: []string{
					"newitem.xyz",
				},
				additionalHeaders: []string{
					"New-Header",
				},
			},
		},

		// Not empty starters.
		{
			name: "nil options lists, not empty starter",
			args: args{
				starter: CORSOptions{
					domains: []string{
						"*",
						"example.com",
						"localhost",
					},
					additionalHeaders: []string{
						"Some-Header-X",
						"Other-Header",
					},
				},
				domains: nil,
				headers: nil,
			},
			want: CORSOptions{
				domains: []string{
					"*",
					"example.com",
					"localhost",
				},
				additionalHeaders: []string{
					"Some-Header-X",
					"Other-Header",
				},
			},
		},
		{
			name: "empty options lists, not empty starter",
			args: args{
				starter: CORSOptions{
					domains: []string{
						"*",
						"example.com",
						"localhost",
					},
					additionalHeaders: []string{
						"Some-Header-X",
						"Other-Header",
					},
				},
				domains: []string{},
				headers: []string{},
			},
			want: CORSOptions{
				domains: []string{
					"*",
					"example.com",
					"localhost",
				},
				additionalHeaders: []string{
					"Some-Header-X",
					"Other-Header",
				},
			},
		},

		{
			name: "not empty options list, not empty starter",
			args: args{
				starter: CORSOptions{
					domains: []string{
						"*",
						"example.com",
						"localhost",
					},
					additionalHeaders: []string{
						"Some-Header-X",
						"Other-Header",
					},
				},
				domains: []string{
					"newitem.xyz",
				},
				headers: []string{
					"New-Header",
				},
			},
			want: CORSOptions{
				domains: []string{
					"*",
					"example.com",
					"localhost",
					"newitem.xyz",
				},
				additionalHeaders: []string{
					"Some-Header-X",
					"Other-Header",
					"New-Header",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			od := WithDomains(tt.args.domains...)
			oh := WithHeaders(tt.args.headers...)

			od(&tt.args.starter)
			oh(&tt.args.starter)

			if !cmp.Equal(tt.want.domains, tt.args.starter.domains) {
				t.Errorf("WithDomains().domains = %v, want %v", tt.args.starter.domains, tt.want.domains)
			}

			if !cmp.Equal(tt.want.additionalHeaders, tt.args.starter.additionalHeaders) {
				t.Errorf("WithDomains().additionalHeaders = %v, want %v", tt.args.starter.additionalHeaders, tt.want.additionalHeaders)
			}
		})
	}
}
