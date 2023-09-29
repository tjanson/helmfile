package state

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/helmfile/helmfile/pkg/helmexec"
	"github.com/helmfile/helmfile/pkg/testutil"
)

func TestAppendCascadeFlags(t *testing.T) {
	type args struct {
		flags    []string
		release  *ReleaseSpec
		cascade  string
		helm     helmexec.Interface
		helmSpec HelmSpec
		expected []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "no cascade when helm less than 3.11.0",
			args: args{
				flags:    []string{},
				release:  &ReleaseSpec{},
				cascade:  "background",
				helm:     testutil.NewVersionHelmExec("3.11.0"),
				expected: []string{},
			},
		},
		{
			name: "cascade from release",
			args: args{
				flags:    []string{},
				release:  &ReleaseSpec{Cascade: &[]string{"background", "background"}[0]},
				cascade:  "",
				helm:     testutil.NewVersionHelmExec("3.12.1"),
				expected: []string{"--cascade", "background"},
			},
		},
		{
			name: "cascade from cmd flag",
			args: args{
				flags:    []string{},
				release:  &ReleaseSpec{},
				cascade:  "background",
				helm:     testutil.NewVersionHelmExec("3.12.1"),
				expected: []string{"--cascade", "background"},
			},
		},
		{
			name: "cascade from helm defaults",
			args: args{
				flags:    []string{},
				release:  &ReleaseSpec{},
				helmSpec: HelmSpec{Cascade: &[]string{"background", "background"}[0]},
				cascade:  "",
				helm:     testutil.NewVersionHelmExec("3.12.1"),
				expected: []string{"--cascade", "background"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &HelmState{}
			st.HelmDefaults = tt.args.helmSpec
			got := st.appendCascadeFlags(tt.args.flags, tt.args.helm, tt.args.release, tt.args.cascade)
			require.Equalf(t, tt.args.expected, got, "appendCascadeFlags() = %v, want %v", got, tt.args.expected)
		})
	}
}

func TestFormatLabels(t *testing.T) {
	tests := []struct {
		name   string
		labels map[string]string
		want   string
	}{
		{
			name:   "empty labels",
			labels: map[string]string{},
			want:   "",
		},
		{
			name:   "single label",
			labels: map[string]string{"foo": "bar"},
			want:   "foo=bar",
		},
		{
			name:   "multiple labels",
			labels: map[string]string{"foo": "bar", "baz": "qux"},
			want:   "baz=qux,foo=bar",
		},
		{
			name:   "multiple labels with empty value",
			labels: map[string]string{"foo": "bar", "baz": "qux", "quux": ""},
			want:   "baz=qux,foo=bar,quux=null",
		},
		{
			name:   "multiple labels with empty key",
			labels: map[string]string{"foo": "bar", "baz": "qux", "": "quux"},
			want:   "baz=qux,foo=bar",
		},
		{
			name:   "empty label value",
			labels: map[string]string{"foo": ""},
			want:   "foo=null",
		},
		{
			name:   "empty label key",
			labels: map[string]string{"": "bar"},
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatLabels(tt.labels)
			require.Equal(t, tt.want, got, "formatLabels() = %v, want %v", got, tt.want)
		})
	}
}
