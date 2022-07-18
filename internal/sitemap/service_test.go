package sitemap

import (
	"gositemap/internal/link"
	"reflect"
	"testing"
)

func Test_siteMap_AddPages(t *testing.T) {
	type fields struct {
		Xmlns string
		Urls  []Loc
	}
	type args struct {
		pages link.HrefSlice
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Loc
	}{
		{
			name: "Should add page",
			fields: fields{
				Xmlns: "test",
				Urls:  []Loc{},
			},
			args: args{
				pages: link.HrefSlice{"www.facebook.com", "www.instagram.com"},
			},
			want: []Loc{{Value: "www.facebook.com"}, {Value: "www.instagram.com"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := &siteMap{
				Xmlns: tt.fields.Xmlns,
				Urls:  tt.fields.Urls,
			}
			sm.AddPages(tt.args.pages)
			if !reflect.DeepEqual(sm.Urls, tt.want) {
				t.Errorf("AddPages() got = %v, want %v", sm.Urls, tt.want)

			}
		})
	}
}
