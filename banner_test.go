package banner

import (
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	nowtime := time.Now()
	type args struct {
		srs []DisplayBanner
	}
	tests := []struct {
		name    string
		args    args
		want    *banner
		wantErr bool
	}{
		{
			name: "Pass zero banners",
			args: args{
				srs: make([]DisplayBanner, 0),
			},
			wantErr: true,
		},
		{
			name: "Single banner",
			args: args{
				srs: []DisplayBanner{{
					Promotion: Promotion{
						Name:  "snowflakes",
						Start: nowtime.Add(-1 * time.Hour),
						End:   nowtime.Add(1 * time.Hour),
					},
					img: "new year.png",
					src: "http://google.com",
				}},
			},
			want: &banner{
				srs: []DisplayBanner{{
					Promotion: Promotion{
						Name:  "snowflakes",
						Start: nowtime.Add(-1 * time.Hour),
						End:   nowtime.Add(1 * time.Hour),
					},
					img: "new year.png",
					src: "http://google.com",
				}},
			},
			wantErr: false,
		},
		{
			name: "Multiple banners - sorted with earlier expiration first",
			args: args{
				srs: []DisplayBanner{{
					Promotion: Promotion{
						Name:  "new year",
						Start: nowtime.Add(-1 * time.Hour),
						End:   nowtime.Add(1 * time.Hour),
					},
					img: "new year.png",
					src: "http://google.com",
				},
					{
						Promotion: Promotion{
							Name:  "sale 90%",
							Start: nowtime.Add(-1 * time.Hour),
							End:   nowtime.Add(30 * time.Minute),
						},
						img: "sale.png",
						src: "http://amazon.com",
					}},
			},
			want: &banner{
				srs: []DisplayBanner{
					{
						Promotion: Promotion{
							Name:  "sale 90%",
							Start: nowtime.Add(-1 * time.Hour),
							End:   nowtime.Add(30 * time.Minute),
						},
						img: "sale.png",
						src: "http://amazon.com",
					},
					{
						Promotion: Promotion{
							Name:  "snowflakes",
							Start: nowtime.Add(-1 * time.Hour),
							End:   nowtime.Add(1 * time.Hour),
						},
						img: "new year.png",
						src: "http://google.com",
					}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.srs)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				return
			}
			// Match first element to check order
			if !reflect.DeepEqual(got.srs[0], tt.want.srs[0]) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_banner_DisplayFor(t *testing.T) {
	nowtime := time.Now()
	type args struct {
		addr string
	}
	tests := []struct {
		name string
		b    *banner
		args args
		want *DisplayBanner
	}{
		{
			name: "Internal IP",
			b: &banner{
				srs: []DisplayBanner{
					{
						Promotion: Promotion{
							Name:  "sale 90%",
							Start: nowtime.Add(5 * time.Minute),
							End:   nowtime.Add(30 * time.Minute),
						},
						img: "sale.png",
						src: "http://amazon.com",
					},
					{
						Promotion: Promotion{
							Name:  "snowflakes",
							Start: nowtime.Add(-1 * time.Hour),
							End:   nowtime.Add(1 * time.Hour),
						},
						img: "new year.png",
						src: "http://google.com",
					}},
			},
			args: args{
				addr: "10.0.0.1",
			},
			want: &DisplayBanner{
				Promotion: Promotion{
					Name:  "sale 90%",
					Start: nowtime.Add(5 * time.Minute),
					End:   nowtime.Add(30 * time.Minute),
				},
				img: "sale.png",
				src: "http://amazon.com",
			},
		},
		{
			name: "External IP",
			b: &banner{
				srs: []DisplayBanner{
					{
						Promotion: Promotion{
							Name:  "sale 90%",
							Start: nowtime.Add(5 * time.Minute),
							End:   nowtime.Add(30 * time.Minute),
						},
						img: "sale.png",
						src: "http://amazon.com",
					},
					{
						Promotion: Promotion{
							Name:  "snowflakes",
							Start: nowtime.Add(-1 * time.Hour),
							End:   nowtime.Add(1 * time.Hour),
						},
						img: "new year.png",
						src: "http://google.com",
					}},
			},
			args: args{
				addr: "210.130.169.196",
			},
			want: &DisplayBanner{
				Promotion: Promotion{
					Name:  "snowflakes",
					Start: nowtime.Add(-1 * time.Hour),
					End:   nowtime.Add(1 * time.Hour),
				},
				img: "new year.png",
				src: "http://google.com",
			},
		}, {
			name: "External IP - no available banner",
			b: &banner{
				srs: []DisplayBanner{
					{
						Promotion: Promotion{
							Name:  "sale 90%",
							Start: nowtime.Add(5 * time.Minute),
							End:   nowtime.Add(30 * time.Minute),
						},
						img: "sale.png",
						src: "http://amazon.com",
					}},
			},
			args: args{
				addr: "210.130.169.196",
			},
			want: nil,
		},
		{
			name: "Internal IP - no available banner",
			b: &banner{
				srs: []DisplayBanner{
					{
						Promotion: Promotion{
							Name:  "sale 90%",
							Start: nowtime.Add(-5 * time.Hour),
							End:   nowtime.Add(-30 * time.Minute),
						},
						img: "sale.png",
						src: "http://amazon.com",
					},
					{
						Promotion: Promotion{
							Name:  "snowflakes",
							Start: nowtime.Add(-1 * time.Hour),
							End:   nowtime.Add(-1 * time.Minute),
						},
						img: "new year.png",
						src: "http://google.com",
					}},
			},
			args: args{
				addr: "10.0.0.1",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.DisplayFor(tt.args.addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("banner.DisplayFor() = %v, want %v", got, tt.want)
			}
		})
	}
}
