package jobspec2

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

func TestDecodeMapInterfaceType(t *testing.T) {
	type args struct {
		V map[string]interface{}
	}
	tests := map[string]struct {
		args  args
		want  interface{}
		diags hcl.Diagnostics
	}{
		"empty-map": {
			args: args{
				V: map[string]interface{}{},
			},
			want: map[string]interface{}{},
		},
		"bool": {
			args: args{
				V: map[string]interface{}{
					"attr": &hcl.Attribute{
						Expr: hcl.StaticExpr(cty.BoolVal(true), hcl.Range{}),
					},
				},
			},
			want: map[string]interface{}{
				"attr": true,
			},
		},
		"number": {
			args: args{
				V: map[string]interface{}{
					"attr": &hcl.Attribute{
						Expr: hcl.StaticExpr(cty.NumberIntVal(1), hcl.Range{}),
					},
				},
			},
			want: map[string]interface{}{
				"attr": 1,
			},
		},
		"string": {
			args: args{
				V: map[string]interface{}{
					"attr": &hcl.Attribute{
						Expr: hcl.StaticExpr(cty.StringVal("hello"), hcl.Range{}),
					},
				},
			},
			want: map[string]interface{}{
				"attr": "hello",
			},
		},
		"null": {
			args: args{
				V: map[string]interface{}{
					"attr": &hcl.Attribute{
						Expr: hcl.StaticExpr(cty.NullVal(cty.DynamicPseudoType), hcl.Range{}),
					},
				},
			},
			want: map[string]interface{}{
				"attr": nil,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			diags := DecodeMapInterfaceType(&tt.args, nil)
			require.Equal(t, diags, tt.diags)
			require.Equal(t, tt.want, tt.args.V)
		})
	}
}
