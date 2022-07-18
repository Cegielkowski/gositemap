package link

import (
	"bytes"
	"golang.org/x/net/html"
	"gositemap/test/mocks"
	"gositemap/test/utils"
	"io"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestGetHrefs(t *testing.T) {
	type args struct {
		r       io.Reader
		baseUrl string
	}
	tests := []struct {
		name    string
		args    args
		want    HrefSlice
		wantErr bool
	}{
		{
			name: "Test with Aurox Href",
			args: args{
				r:       ioutil.NopCloser(bytes.NewReader([]byte(mocks.PageAurox))),
				baseUrl: "https://getaurox.com/",
			},
			want: HrefSlice{
				"https://getaurox.com/wallet",
				"https://getaurox.com/pdf/whitepaper.pdf",
				"https://getaurox.com/protocol",
				"https://getaurox.com/html/privacy.html",
				"https://getaurox.com/html/tos.html",
				"https://getaurox.com/terminal",
			},
			wantErr: false,
		},
		{
			name: "Test with no Href",
			args: args{
				r:       ioutil.NopCloser(bytes.NewReader([]byte(mocks.PageWithNoHref))),
				baseUrl: "https://isaprotege.com.br",
			},
			want:    HrefSlice{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHrefs(tt.args.r, tt.args.baseUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHrefs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !utils.SameStringSlice(got, tt.want) {
				t.Errorf("GetHrefs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHrefSlice_safeAppend(t *testing.T) {
	type args struct {
		value   string
		baseUrl string
	}
	tests := []struct {
		name string
		s    HrefSlice
		args args
		want HrefSlice
	}{
		{
			name: "Should append",
			s:    HrefSlice{},
			args: args{
				value:   "https://getaurox.com/hire-me",
				baseUrl: "https://getaurox.com/",
			},
			want: HrefSlice{"https://getaurox.com/hire-me"},
		},
		{
			name: "Should not append",
			s:    HrefSlice{},
			args: args{
				value:   "https://google.com/",
				baseUrl: "https://getaurox.com/",
			},
			want: HrefSlice{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.safeAppend(tt.args.value, tt.args.baseUrl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("safeAppend() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Should parse correctly",
			args: args{
				r: ioutil.NopCloser(bytes.NewReader([]byte(mocks.PageWithNoHref))),
			},
			want: []string{
				"https://www.instagram.com/bravemindsolutions/",
				"mailto:comercial@bravemind.com.br",
				"#body",
				"#functionalities",
				"https://api.whatsapp.com/send?phone=5514998583391&text=Ol%C3%A1%2C%20fiquei%20interessado%20sobre%20a%20Brave%20Mind%2C%20voc%C3%AA%20pode%20me%20tirar%20algumas%20d%C3%BAvidas%3F",
				"https://www.linkedin.com/company/brave-mind/",
				"https://www.facebook.com/bravemindsoftware/",
				"#conteudo",
				"#benefits",
				"#contact",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !utils.SameStringSlice(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlice_ToSlice(t *testing.T) {
	tests := []struct {
		name string
		ms   Slice
		want HrefSlice
	}{
		{
			name: "Should convert properly",
			ms:   Slice{"test": struct{}{}, "test2": struct{}{}},
			want: HrefSlice{"test", "test2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ms.ToSlice(); !utils.SameStringSlice(got, tt.want) {
				t.Errorf("ToSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlice_ToSliceSafe(t *testing.T) {
	type args struct {
		baseUrl string
	}
	tests := []struct {
		name string
		ms   Slice
		args args
		want HrefSlice
	}{
		{
			name: "Should bring nothing",
			ms:   Slice{"https://google.com": struct{}{}, "https://facebook.com": struct{}{}},
			args: args{
				baseUrl: "https://getaurox.com/",
			},
			want: HrefSlice{},
		},
		{
			name: "Should bring just one",
			ms:   Slice{"https://google.com": struct{}{}, "https://getaurox.com/golang": struct{}{}},
			args: args{
				baseUrl: "https://getaurox.com/",
			},
			want: HrefSlice{"https://getaurox.com/golang"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ms.ToSliceSafe(tt.args.baseUrl); !utils.SameStringSlice(got, tt.want) {
				t.Errorf("ToSliceSafe() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildLink(t *testing.T) {
	doc, _ := html.Parse(ioutil.NopCloser(bytes.NewReader([]byte(mocks.PageWithNoHref))))

	nodes := linkNodes(doc)
	node := nodes[9]

	type args struct {
		n *html.Node
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should build the link",
			args: args{
				n: node,
			},
			want: "https://www.instagram.com/bravemindsolutions/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildLink(tt.args.n); got != tt.want {
				t.Errorf("buildLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hasBaseUrl(t *testing.T) {
	type args struct {
		link    string
		baseurl string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Should have the baseurl",
			args: args{
				link:    "https://getaurox.com/golang",
				baseurl: "https://getaurox.com/",
			},
			want: true,
		},
		{
			name: "Should noot have the baseurl",
			args: args{
				link:    "https://google.com/golang",
				baseurl: "https://getaurox.com/",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasBaseUrl(tt.args.link, tt.args.baseurl); got != tt.want {
				t.Errorf("hasBaseUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hasMoreThanBaseUrl(t *testing.T) {
	type args struct {
		link    string
		baseUrl string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Should have more then baseurl",
			args: args{
				link:    "https://getaurox.com/golang",
				baseUrl: "https://getaurox.com/",
			},
			want: true,
		},
		{
			name: "Should not have more then baseurl",
			args: args{
				link:    "https://getaurox.com",
				baseUrl: "https://getaurox.com",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasMoreThanBaseUrl(tt.args.link, tt.args.baseUrl); got != tt.want {
				t.Errorf("hasMoreThenBaseUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isNotRelatedToBaseUrl(t *testing.T) {
	type args struct {
		link    string
		baseUrl string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Should be related",
			args: args{
				link:    "https://google.com/golang",
				baseUrl: "https://getaurox.com",
			},
			want: true,
		},
		{
			name: "Should be related",
			args: args{
				link:    "https://getaurox.com/contact-us",
				baseUrl: "https://getaurox.com/",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNotRelatedToBaseUrl(tt.args.link, tt.args.baseUrl); got != tt.want {
				t.Errorf("isNotRelatedToBaseUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_linkNodes(t *testing.T) {
	doc, _ := html.Parse(ioutil.NopCloser(bytes.NewReader([]byte(mocks.PageWithNoHref))))

	type args struct {
		n *html.Node
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "",
			args: args{
				n: doc,
			},
			want: 15,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := linkNodes(tt.args.n); len(got) != tt.want {
				t.Errorf("linkNodes len  = %v, want %v", len(got), tt.want)
			}
		})
	}
}
