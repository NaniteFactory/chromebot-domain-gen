package main

//go:generate qtc -dir gen/gotpl -ext qtpl
//go:generate gofmt -w -s gen/gotpl/

import (
	"github.com/nanitefactory/chromebot-domain-gen/gen/gotpl"
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	glob "github.com/ryanuber/go-glob"
	"golang.org/x/sync/errgroup"
	"golang.org/x/tools/imports"

	"github.com/chromedp/cdproto-gen/diff"
	"github.com/chromedp/cdproto-gen/fixup"
	"github.com/nanitefactory/chromebot-domain-gen/gen"
	"github.com/chromedp/cdproto-gen/gen/genutil"
	"github.com/chromedp/cdproto-gen/pdl"
	"github.com/chromedp/cdproto-gen/util"
)

const (
	easyjsonGo = "easyjson.go"
)

var (
	flagDebug = flag.Bool("debug", false, "toggle debug (writes generated files to disk without post-processing)")

	flagTTL = flag.Duration("ttl", 24*time.Hour, "file retrieval caching ttl")

	flagChromium = flag.String("chromium", "", "chromium protocol version")
	flagV8       = flag.String("v8", "", "v8 protocol version")
	flagLatest   = flag.Bool("latest", false, "use latest protocol")

	flagPdl = flag.String("pdl", "", "path to pdl file to use")

	flagCache = flag.String("cache", "", "protocol cache directory")
	flagOut   = flag.String("out", "", "package out directory")

	flagNoClean = flag.Bool("no-clean", false, "toggle not cleaning (removing) existing directories")
	flagNoDump  = flag.Bool("no-dump", false, "toggle not dumping generated protocol file to out directory")

	flagGoPkg = flag.String("go-pkg", "github.com/nanitefactory/chromebot/domain", "go base package name")
	flagGoWl  = flag.String("go-wl", "LICENSE,README.md,*.pdl,go.mod,go.sum,"+easyjsonGo, "comma-separated list of files to whitelist (ignore)")

	// flagWorkers = flag.Int("workers", runtime.NumCPU(), "number of workers")
)

func main() {
	// add generator parameters
	var genTypes []string
	generators := gen.Generators()
	for n, g := range generators {
		genTypes = append(genTypes, n)
		g = g
	}

	flag.Parse()

	// run
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// run runs the generator.
func run() error {
	var err error

	// set cache path
	if *flagCache == "" {
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return err
		}
		*flagCache = filepath.Join(cacheDir, "cdproto-gen")
	}

	// get latest versions
	if *flagChromium == "" {
		if *flagChromium, err = util.GetLatestVersion(util.Cache{
			URL:  util.ChromiumBase,
			Path: filepath.Join(*flagCache, "html", "chromium.html"),
			TTL:  *flagTTL,
		}); err != nil {
			return err
		}
	}
	if *flagV8 == "" {
		if *flagLatest {
			if *flagV8, err = util.GetLatestVersion(util.Cache{
				URL:  util.V8Base,
				Path: filepath.Join(*flagCache, "html", "v8.html"),
				TTL:  *flagTTL,
			}); err != nil {
				return err
			}
		} else {
			if *flagV8, err = util.GetDepVersion("v8", *flagChromium, util.Cache{
				URL:    fmt.Sprintf(util.ChromiumDeps+"?format=TEXT", *flagChromium),
				Path:   filepath.Join(*flagCache, "deps", "chromium", *flagChromium),
				TTL:    *flagTTL,
				Decode: true,
			}, util.Cache{
				URL:  util.V8Base + "/+refs?format=JSON",
				Path: filepath.Join(*flagCache, "refs", "v8.json"),
				TTL:  *flagTTL,
			}); err != nil {
				return err
			}
		}
	}

	// load protocol definitions
	protoDefs, err := loadProtoDefs()
	if err != nil {
		return err
	}
	sort.Slice(protoDefs.Domains, func(i, j int) bool {
		return strings.Compare(protoDefs.Domains[i].Domain.String(), protoDefs.Domains[j].Domain.String()) <= 0
	})

	if *flagOut == "" {
		*flagOut = filepath.Join(os.Getenv("GOPATH"), "src", *flagGoPkg)
	} else {
		*flagOut, err = filepath.Abs(*flagOut)
		if err != nil {
			return err
		}
	}

	// create out directory
	if err = os.MkdirAll(*flagOut, 0755); err != nil {
		return err
	}

	combinedDir := filepath.Join(*flagCache, "pdl", "combined")
	if err = os.MkdirAll(combinedDir, 0755); err != nil {
		return err
	}
	protoFile := filepath.Join(combinedDir, fmt.Sprintf("%s_%s.pdl", *flagChromium, *flagV8))

	// write protocol definitions
	if *flagPdl == "" {
		util.Logf("WRITING: %s", protoFile)
		if err = ioutil.WriteFile(protoFile, protoDefs.Bytes(), 0644); err != nil {
			return err
		}

		// display differences between generated definitions and previous version on disk
		if runtime.GOOS != "windows" {
			diffBuf, err := diff.WalkAndCompare(combinedDir, `^([0-9_.]+)\.pdl$`, protoFile, func(a, b *diff.FileInfo) bool {
				n := strings.Split(strings.TrimSuffix(filepath.Base(a.Name), ".pdl"), "_")
				m := strings.Split(strings.TrimSuffix(filepath.Base(b.Name), ".pdl"), "_")
				if n[0] == m[0] {
					return util.CompareSemver(n[1], m[1])
				}
				return util.CompareSemver(n[0], m[0])
			})
			if err != nil {
				return err
			}
			if diffBuf != nil {
				os.Stdout.Write(diffBuf)
			}
		}
	}

	// determine what to process
	pkgs := []string{"", "cdp"}
	var processed []*pdl.Domain
	for _, d := range protoDefs.Domains {
		// skip if not processing
		if d.Deprecated {
			var extra []string
			if d.Deprecated {
				extra = append(extra, "deprecated")
			}
			util.Logf("SKIPPING(%s): %s %v", pad("domain", 7), d.Domain.String(), extra)
			continue
		}

		// will process
		pkgs = append(pkgs, genutil.PackageName(d))
		processed = append(processed, d)

		// cleanup types, events, commands
		d.Types = cleanupTypes("type", d.Domain.String(), d.Types)
		d.Events = cleanupTypes("event", d.Domain.String(), d.Events)
		d.Commands = cleanupTypes("command", d.Domain.String(), d.Commands)
	}

	// fixup
	fixup.FixDomains(processed)

	// get generator
	generator := gen.Generators()["go"]
	if generator == nil {
		return errors.New("no generator")
	}

	// emit
	emitter, err := generator(processed, *flagGoPkg)
	if err != nil {
		return err
	}
	files := emitter.Emit()

	// clean up files
	if !*flagNoClean {
		util.Logf("CLEANING: %s", *flagOut)
		outpath := *flagOut + string(filepath.Separator)
		err = filepath.Walk(outpath, func(n string, fi os.FileInfo, err error) error {
			switch {
			case os.IsNotExist(err) || n == outpath:
				return nil
			case err != nil:
				return err
			}

			// skip if file or path starts with ., is whitelisted, or is one of
			// the files whose output will be overwritten
			pn, fn := n[len(outpath):], fi.Name()
			if pn == "" || strings.HasPrefix(pn, ".") || strings.HasPrefix(fn, ".") || whitelisted(fn) || contains(files, pn) {
				return nil
			}

			util.Logf("REMOVING: %s", n)
			return os.RemoveAll(n)
		})
		if err != nil {
			return err
		}
	}

	util.Logf("WRITING: %d files", len(files))

	// dump files and exit
	if *flagDebug {
		return write(files)
	}

	// goimports (also writes to disk)
	if err = goimports(files); err != nil {
		return err
	}

	// gofmt
	if err = gofmt(fmtFiles(files, pkgs)); err != nil {
		return err
	}
	
	// domain manager
	if err := func() error {
		strFilepath := filepath.Join(*flagOut, "domain.go")
		util.Logf("WRITING: %s", strFilepath)
		f, err := os.Create(strFilepath)
		if err != nil {
			return err
		}
		defer f.Close()
		gotpl.WriteDomainManagerTemplate(f, processed)	
		return nil
	}(); err != nil {
		return err
	}

	util.Logf("done.")
	return nil
}

// loadProtoDefs loads the protocol definitions either from the path specified
// in -proto or by retrieving the versions specified in the -browser and -js
// files.
func loadProtoDefs() (*pdl.PDL, error) {
	var err error

	if *flagPdl != "" {
		util.Logf("PROTOCOL: %s", *flagPdl)
		buf, err := ioutil.ReadFile(*flagPdl)
		if err != nil {
			return nil, err
		}
		return pdl.Parse(buf)
	}

	var protoDefs []*pdl.PDL
	load := func(urlstr, typ, ver string) error {
		buf, err := util.Get(util.Cache{
			URL:    fmt.Sprintf(urlstr+"?format=TEXT", ver),
			Path:   filepath.Join(*flagCache, "pdl", typ, ver+".pdl"),
			TTL:    *flagTTL,
			Decode: true,
		})
		if err != nil {
			return err
		}

		// parse
		protoDef, err := pdl.Parse(buf)
		if err != nil {
			return err
		}
		protoDefs = append(protoDefs, protoDef)
		return nil
	}

	// grab browser + js definition
	if err = load(util.ChromiumURL, "chromium", *flagChromium); err != nil {
		return nil, err
	}
	if err = load(util.V8URL, "v8", *flagV8); err != nil {
		return nil, err
	}

	// grab har definition
	har, err := pdl.Parse([]byte(pdl.HAR))
	if err != nil {
		return nil, err
	}

	return pdl.Combine(append(protoDefs, har)...), nil
}

// cleanupTypes removes deprecated and redirected types.
func cleanupTypes(n string, dtyp string, typs []*pdl.Type) []*pdl.Type {
	var ret []*pdl.Type

	for _, t := range typs {
		typ := dtyp + "." + t.Name
		if t.Deprecated {
			util.Logf("SKIPPING(%s): %s [deprecated]", pad(n, 7), typ)
			continue
		}

		if t.Redirect != nil {
			util.Logf("SKIPPING(%s): %s [redirect:%s]", pad(n, 7), typ, t.Redirect)
			continue
		}

		if t.Properties != nil {
			t.Properties = cleanupTypes(n[0:1]+" property", typ, t.Properties)
		}

		if t.Parameters != nil {
			t.Parameters = cleanupTypes(n[0:1]+" param", typ, t.Parameters)
		}

		if t.Returns != nil {
			t.Returns = cleanupTypes(n[0:1]+" return param", typ, t.Returns)
		}

		ret = append(ret, t)
	}

	return ret
}

// write writes all file buffer to disk.
func write(fileBuffers map[string]*bytes.Buffer) error {
	var keys []string
	for k := range fileBuffers {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		// add out path
		_, filename := filepath.Split(k)
		n := filepath.Join(*flagOut, filename)

		// write file
		if err := ioutil.WriteFile(n, fileBuffers[k].Bytes(), 0644); err != nil {
			return err
		}
	}
	return nil
}

// goimports formats all the output file buffers on disk using goimports.
func goimports(fileBuffers map[string]*bytes.Buffer) error {
	util.Logf("RUNNING: goimports")

	var keys []string
	for k := range fileBuffers {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	eg, _ := errgroup.WithContext(context.Background())
	for _, k := range keys {
		eg.Go(func(n string) func() error {
			return func() error {
				_, filename := filepath.Split(n)
				fn := filepath.Join(*flagOut, filename)
				buf, err := imports.Process(fn, fileBuffers[n].Bytes(), nil)
				if err != nil {
					return err
				}
				return ioutil.WriteFile(fn, buf, 0644)
			}
		}(k))
	}
	return eg.Wait()
}

// gofmt go formats all files on disk.
func gofmt(files []string) error {
	util.Logf("RUNNING: gofmt")
	eg, _ := errgroup.WithContext(context.Background())
	for _, k := range files {
		eg.Go(func(n string) func() error {
			return func() error {
				_, filename := filepath.Split(n)
				n = filepath.Join(*flagOut, filename)
				in, err := ioutil.ReadFile(n)
				if err != nil {
					return err
				}
				out, err := format.Source(in)
				if err != nil {
					return err
				}
				return ioutil.WriteFile(n, out, 0644)
			}
		}(k))
	}
	return eg.Wait()
}

// fmtFiles returns the list of all files to format from the specified file
// buffers and packages.
func fmtFiles(files map[string]*bytes.Buffer, pkgs []string) []string {
	filelen := len(files)
	f := make([]string, filelen)

	var i int
	for n := range files {
		f[i] = n
		i++
	}

	sort.Strings(f)
	return f
}

// contains determines if any key in m is equal to n or starts with the path
// prefix equal to n.
func contains(m map[string]*bytes.Buffer, n string) bool {
	d := n + string(filepath.Separator)
	for k := range m {
		if n == k || strings.HasPrefix(k, d) {
			return true
		}
	}
	return false
}

// pad pads a string.
func pad(s string, n int) string {
	return s + strings.Repeat(" ", n-len(s))
}

// whitelisted checks if n is a whitelisted file.
func whitelisted(n string) bool {
	for _, z := range strings.Split(*flagGoWl, ",") {
		if z == n || glob.Glob(z, n) {
			return true
		}
	}
	return false
}
