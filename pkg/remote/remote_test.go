package remote

import (
	"fmt"
	"io"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"

	"github.com/helmfile/helmfile/pkg/filesystem"
	"github.com/helmfile/helmfile/pkg/helmexec"
	"github.com/helmfile/helmfile/pkg/testhelper"
)

func TestRemote_HttpsGitHub(t *testing.T) {
	cleanfs := map[string]string{
		CacheDir(): "",
	}
	cachefs := map[string]string{
		filepath.Join(CacheDir(), "https_github_com_cloudposse_helmfiles_git.ref=0.40.0/origin/releases/kiam.yaml"): "foo: bar",
	}

	type testcase struct {
		files          map[string]string
		expectCacheHit bool
	}

	testcases := []testcase{
		{files: cleanfs, expectCacheHit: false},
		{files: cachefs, expectCacheHit: true},
	}

	for i := range testcases {
		testcase := testcases[i]

		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			testfs := testhelper.NewTestFs(testcase.files)

			hit := true

			get := func(wd, src, dst string) error {
				if wd != CacheDir() {
					return fmt.Errorf("unexpected wd: %s", wd)
				}
				if src != "git::https://github.com/cloudposse/helmfiles.git?ref=0.40.0" {
					return fmt.Errorf("unexpected src: %s", src)
				}

				hit = false

				testfs.AddFiles(cachefs)

				return nil
			}

			getter := &testGetter{
				get: get,
			}
			remote := &Remote{
				Logger: helmexec.NewLogger(io.Discard, "debug"),
				Home:   CacheDir(),
				Getter: getter,
				fs:     testfs.ToFileSystem(),
			}

			// FYI, go-getter in the `dir` mode accepts URL like the below. So helmfile expects URLs similar to it:
			//   go-getter -mode dir git::https://github.com/cloudposse/helmfiles.git?ref=0.40.0 gettertest1/b

			// We use `@` to separate dir and the file path. This is a good idea borrowed from helm-git:
			//   https://github.com/aslafy-z/helm-git

			url := "git::https://github.com/cloudposse/helmfiles.git@releases/kiam.yaml?ref=0.40.0"
			file, err := remote.Locate(url)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			expectedFile := filepath.Join(CacheDir(), "https_github_com_cloudposse_helmfiles_git.ref=0.40.0/origin/releases/kiam.yaml")
			if file != expectedFile {
				t.Errorf("unexpected file located: %s vs expected: %s", file, expectedFile)
			}

			if testcase.expectCacheHit && !hit {
				t.Errorf("unexpected result: unexpected cache miss")
			}
			if !testcase.expectCacheHit && hit {
				t.Errorf("unexpected result: unexpected cache hit")
			}
		})
	}
}

func TestRemote_SShGitHub(t *testing.T) {
	cleanfs := map[string]string{
		CacheDir(): "",
	}
	cachefs := map[string]string{
		filepath.Join(CacheDir(), "ssh_github_com_cloudposse_helmfiles_git.ref=0.40.0/origin/releases/kiam.yaml"): "foo: bar",
	}

	type testcase struct {
		files          map[string]string
		expectCacheHit bool
	}

	testcases := []testcase{
		{files: cleanfs, expectCacheHit: false},
		{files: cachefs, expectCacheHit: true},
	}

	for i := range testcases {
		testcase := testcases[i]

		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			testfs := testhelper.NewTestFs(testcase.files)

			hit := true

			get := func(wd, src, dst string) error {
				if wd != CacheDir() {
					return fmt.Errorf("unexpected wd: %s", wd)
				}
				if src != "git::ssh://git@github.com/cloudposse/helmfiles.git?ref=0.40.0" {
					return fmt.Errorf("unexpected src: %s", src)
				}

				hit = false

				testfs.AddFiles(cachefs)

				return nil
			}

			getter := &testGetter{
				get: get,
			}
			remote := &Remote{
				Logger: helmexec.NewLogger(io.Discard, "debug"),
				Home:   CacheDir(),
				Getter: getter,
				fs:     testfs.ToFileSystem(),
			}

			url := "git::ssh://git@github.com/cloudposse/helmfiles.git@releases/kiam.yaml?ref=0.40.0"
			file, err := remote.Locate(url)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			expectedFile := filepath.Join(CacheDir(), "ssh_github_com_cloudposse_helmfiles_git.ref=0.40.0/origin/releases/kiam.yaml")
			if file != expectedFile {
				t.Errorf("unexpected file located: %s vs expected: %s", file, expectedFile)
			}

			if testcase.expectCacheHit && !hit {
				t.Errorf("unexpected result: unexpected cache miss")
			}
			if !testcase.expectCacheHit && hit {
				t.Errorf("unexpected result: unexpected cache hit")
			}
		})
	}
}

func TestRemote_SShGitHub_WithSshKey(t *testing.T) {
	cleanfs := map[string]string{
		CacheDir(): "",
	}
	cachefs := map[string]string{
		filepath.Join(CacheDir(), "ssh_github_com_cloudposse_helmfiles_git.ref=0.40.0_sshkey=redacted/origin/releases/kiam.yaml"): "foo: bar",
	}

	type testcase struct {
		files          map[string]string
		expectCacheHit bool
	}

	testcases := []testcase{
		{files: cleanfs, expectCacheHit: false},
		{files: cachefs, expectCacheHit: true},
	}

	for i := range testcases {
		testcase := testcases[i]

		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			testfs := testhelper.NewTestFs(testcase.files)

			hit := true

			get := func(wd, src, dst string) error {
				if wd != CacheDir() {
					return fmt.Errorf("unexpected wd: %s", wd)
				}
				if src != "git::ssh://git@github.com/cloudposse/helmfiles.git?ref=0.40.0&sshkey=ZWNkc2Etc2hhMi1uaXN0cDI1NiBBQUFBRTJWalpITmhMWE5vWVRJdGJtbHpkSEF5TlRZQUFBQUlibWx6ZEhBeU5UWUFBQUJCQkJTU3dOY2xoVzQ2Vm9VR3dMQ3JscVRHYUdOVWdRVUVEUEptc1ZzdUViL2RBNUcrQk9YMWxGaUVMYU9HQ2F6bS9KQkR2V3Y2Y0ZDQUtVRjVocVJOUjdJPSA=" {
					return fmt.Errorf("unexpected src: %s", src)
				}

				hit = false

				testfs.AddFiles(cachefs)

				return nil
			}

			getter := &testGetter{
				get: get,
			}
			remote := &Remote{
				Logger: helmexec.NewLogger(io.Discard, "debug"),
				Home:   CacheDir(),
				Getter: getter,
				fs:     testfs.ToFileSystem(),
			}

			url := "git::ssh://git@github.com/cloudposse/helmfiles.git@releases/kiam.yaml?ref=0.40.0&sshkey=ZWNkc2Etc2hhMi1uaXN0cDI1NiBBQUFBRTJWalpITmhMWE5vWVRJdGJtbHpkSEF5TlRZQUFBQUlibWx6ZEhBeU5UWUFBQUJCQkJTU3dOY2xoVzQ2Vm9VR3dMQ3JscVRHYUdOVWdRVUVEUEptc1ZzdUViL2RBNUcrQk9YMWxGaUVMYU9HQ2F6bS9KQkR2V3Y2Y0ZDQUtVRjVocVJOUjdJPSA="
			file, err := remote.Locate(url)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			expectedFile := filepath.Join(CacheDir(), "ssh_github_com_cloudposse_helmfiles_git.ref=0.40.0_sshkey=redacted/origin/releases/kiam.yaml")
			if file != expectedFile {
				t.Errorf("unexpected file located: %s vs expected: %s", file, expectedFile)
			}

			if testcase.expectCacheHit && !hit {
				t.Errorf("unexpected result: unexpected cache miss")
			}
			if !testcase.expectCacheHit && hit {
				t.Errorf("unexpected result: unexpected cache hit")
			}
		})
	}
}

func TestParse(t *testing.T) {
	type testcase struct {
		name                                   string
		input                                  string
		getter, scheme, dir, file, query, user string
		err                                    string
	}

	testcases := []testcase{
		{
			name:  "local path",
			input: "raw/incubator",
			err:   "parse url: missing scheme - probably this is a local file path? raw/incubator",
		},
		{
			name:   "remote path with full args",
			input:  "git::https://user:password@github.com/stakater/Forecastle.git@deployments/kubernetes/chart/forecastle?ref=v1.0.54",
			getter: "git",
			scheme: "https",
			dir:    "/stakater/Forecastle.git",
			file:   "deployments/kubernetes/chart/forecastle",
			query:  "ref=v1.0.54",
			user:   "user:password",
		},
		{
			name:   "remote path with no file",
			input:  "git::https://github.com/stakater/Forecastle.git",
			getter: "git",
			scheme: "https",
			dir:    "/stakater/Forecastle.git",
		},
	}

	for i := range testcases {
		tc := testcases[i]

		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			src, err := Parse(tc.input)

			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}

			if diff := cmp.Diff(tc.err, errMsg); diff != "" {
				t.Fatalf("Unexpected error:\n%s", diff)
			}

			var getter, scheme, dir, file, query, user string
			if src != nil {
				getter = src.Getter
				scheme = src.Scheme
				dir = src.Dir
				file = src.File
				query = src.RawQuery
				user = src.User
			}

			if diff := cmp.Diff(tc.getter, getter); diff != "" {
				t.Fatalf("Unexpected getter:\n%s", diff)
			}

			if diff := cmp.Diff(tc.scheme, scheme); diff != "" {
				t.Fatalf("Unexpected scheme:\n%s", diff)
			}

			if diff := cmp.Diff(tc.file, file); diff != "" {
				t.Fatalf("Unexpected file:\n%s", diff)
			}

			if diff := cmp.Diff(tc.dir, dir); diff != "" {
				t.Fatalf("Unexpected dir:\n%s", diff)
			}

			if diff := cmp.Diff(tc.query, query); diff != "" {
				t.Fatalf("Unexpected query:\n%s", diff)
			}
			if diff := cmp.Diff(tc.user, user); diff != "" {
				t.Fatalf("Unexpected user:\n%s", diff)
			}
		})
	}
}

type testGetter struct {
	get func(wd, src, dst string) error
}

func (t *testGetter) Get(wd, src, dst string) error {
	return t.get(wd, src, dst)
}

func TestRemote_Fetch(t *testing.T) {
	cleanfs := map[string]string{
		CacheDir(): "",
	}
	cachefs := map[string]string{
		filepath.Join(CacheDir(), "https_github_com_helmfile_helmfile_git.ref=v0.151.0/origin/README.md"): "foo: bar",
	}

	type testcase struct {
		files          map[string]string
		expectCacheHit bool
		cacheDirOpt    string
	}

	testcases := []testcase{
		{files: cleanfs, expectCacheHit: false, cacheDirOpt: ""},
		{files: cachefs, expectCacheHit: true, cacheDirOpt: ""},
		{files: cleanfs, expectCacheHit: false, cacheDirOpt: "states"},
		{files: cachefs, expectCacheHit: true, cacheDirOpt: "states"},
	}

	for i := range testcases {
		testcase := testcases[i]

		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			testfs := testhelper.NewTestFs(testcase.files)

			hit := true

			get := func(wd, src, dst string) error {
				if wd != CacheDir() {
					return fmt.Errorf("unexpected wd: %s", wd)
				}
				if src != "git::https://github.com/helmfile/helmfile.git?ref=v0.151.0" {
					return fmt.Errorf("unexpected src: %s", src)
				}

				hit = false

				testfs.AddFiles(cachefs)

				return nil
			}

			getter := &testGetter{
				get: get,
			}
			remote := &Remote{
				Logger: helmexec.NewLogger(io.Discard, "debug"),
				Home:   CacheDir(),
				Getter: getter,
				fs:     testfs.ToFileSystem(),
			}

			url := "git::https://github.com/helmfile/helmfile.git@README.md?ref=v0.151.0"
			file, err := remote.Fetch(url)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			expectedFile := filepath.Join(CacheDir(), "https_github_com_helmfile_helmfile_git.ref=v0.151.0/origin/README.md")
			if file != expectedFile {
				t.Errorf("unexpected file located: %s vs expected: %s", file, expectedFile)
			}

			if testcase.expectCacheHit && !hit {
				t.Errorf("unexpected result: unexpected cache miss")
			}
			if !testcase.expectCacheHit && hit {
				t.Errorf("unexpected result: unexpected cache hit")
			}
		})
	}
}

func TestRemote_CommonHttpUrl(t *testing.T) {
	testcases := []struct {
		name   string
		input  string
		rpath  string
		errStr string
	}{
		{
			name:  "common git url",
			input: "git::https://github.com/helmfile/helmfile.git?ref=v0.153.1",
			rpath: filepath.Join(CacheDir(), "https_github_com_helmfile_helmfile_git.ref=v0.153.1/origin"),
		},
		{
			name:  "common git url with exist subpath",
			input: "git::https://github.com/dragonflyoss/helm-charts.git@charts?ref=dragonfly-1.0.2",
			rpath: filepath.Join(CacheDir(), "https_github_com_dragonflyoss_helm-charts_git.ref=dragonfly-1.0.2/origin/charts"),
		},
		{
			name:  "common git url with no-exist subpath",
			input: "git::https://github.com/dragonflyoss/helm-charts.git@no-existcharts?ref=dragonfly-1.0.2",
			rpath: filepath.Join(CacheDir(), "https_github_com_dragonflyoss_helm-charts_git.ref=dragonfly-1.0.2/origin/no-existcharts"),
		},
		{
			name:  "common http url",
			input: "https://raw.githubusercontent.com/helmfile/testdata/main/remote-values/value.yaml",
			rpath: filepath.Join(CacheDir(), "https_raw_githubusercontent_com_helmfile_testdata_main_remote-values_value_yaml/origin"),
		},
		{
			name:   "common http no-exist url",
			input:  "https://raw.githubusercontent.com/helmfile/testdata/main/remote-values/no-exist-value.yaml",
			errStr: "get: bad response code: 404",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			remote := NewRemote(helmexec.NewLogger(io.Discard, "debug"), CacheDir(), filesystem.DefaultFileSystem())

			rPath, err := remote.Fetch(tt.input)
			errStr := ""
			if err != nil {
				errStr = err.Error()
			}

			require.Equalf(t, tt.errStr, errStr, "unexpected error: %s", err)
			require.Equalf(t, tt.rpath, rPath, "unexpected rpath: %s", rPath)
		})
	}
}
