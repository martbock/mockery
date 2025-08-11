package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindPkgPath(t *testing.T) {
	pkgPath, err := findPkgPath("./fixtures")
	require.NoError(t, err)
	assert.NotEqual(t, "", pkgPath)
}

func Test_goimports(t *testing.T) {
	before := `package foo

import (
	"fmt"
	"golang.org/x/text"
	"github.com/example/bar"
)
`

	type args struct {
		src     []byte
		options map[string]any
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "no formatter options",
			args: args{
				src:     []byte(before),
				options: nil,
			},
			want: `package foo

import (
	"fmt"

	"github.com/example/bar"
	"golang.org/x/text"
)
`,
			wantErr: assert.NoError,
		},
		{
			name: "empty formatter options",
			args: args{
				src:     []byte(before),
				options: map[string]any{},
			},
			want: `package foo

import (
	"fmt"

	"github.com/example/bar"
	"golang.org/x/text"
)
`,
			wantErr: assert.NoError,
		},
		{
			name: "local prefix configured",
			args: args{
				src:     []byte(before),
				options: map[string]any{"local-prefix": "github.com/example"},
			},
			want: `package foo

import (
	"fmt"

	"golang.org/x/text"

	"github.com/example/bar"
)
`,
			wantErr: assert.NoError,
		},
		{
			name: "invalid local prefix option type",
			args: args{
				src:     []byte(before),
				options: map[string]any{"local-prefix": 42},
			},
			want: "",
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...any) bool {
				return assert.ErrorContains(t, err, "42", msgAndArgs...)
			},
		},
		{
			name: "incomplete func declaration",
			args: args{
				src:     []byte("package foo\n func ("),
				options: map[string]any{"local-prefix": "github.com/example"},
			},
			want: "",
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...any) bool {
				return assert.ErrorContains(t, err, "goimports", msgAndArgs...)
			},
		},
		{
			name: "unknown option",
			args: args{
				src:     []byte(before),
				options: map[string]any{"some-option": "foo"},
			},
			want: `package foo

import (
	"fmt"

	"github.com/example/bar"
	"golang.org/x/text"
)
`,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := goimports(tt.args.src, tt.args.options)
			tt.wantErr(t, err)
			assert.Equalf(t, tt.want, string(got), "formatted output mismatch for %q", tt.name)
		})
	}
}
