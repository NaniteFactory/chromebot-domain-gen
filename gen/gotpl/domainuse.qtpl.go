// Code generated by qtc from "domainuse.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line domainuse.qtpl:1
package gotpl

//line domainuse.qtpl:1
import (
	"github.com/chromedp/cdproto-gen/gen/genutil"
	"github.com/chromedp/cdproto-gen/pdl"
	"strings"
)

// DomainTemplate is the template for a single domain file.

//line domainuse.qtpl:8
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line domainuse.qtpl:8
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line domainuse.qtpl:8
func StreamDomainTemplate(qw422016 *qt422016.Writer, d *pdl.Domain, domains []*pdl.Domain) {
//line domainuse.qtpl:8
	qw422016.N().S(`
// `)
//line domainuse.qtpl:9
	qw422016.N().S(d.Domain.String())
//line domainuse.qtpl:9
	qw422016.N().S(` executes a cdproto command under `)
//line domainuse.qtpl:9
	qw422016.N().S(d.Domain.String())
//line domainuse.qtpl:9
	qw422016.N().S(` domain.
type `)
//line domainuse.qtpl:10
	qw422016.N().S(d.Domain.String())
//line domainuse.qtpl:10
	qw422016.N().S(` struct {
	ctxWithExecutor context.Context
}
`)
//line domainuse.qtpl:13
	for _, c := range d.Commands {
//line domainuse.qtpl:13
		qw422016.N().S(`
`)
//line domainuse.qtpl:14
		qw422016.N().S(CommandTemplate(c, d, domains))
//line domainuse.qtpl:14
		qw422016.N().S(`
`)
//line domainuse.qtpl:15
	}
//line domainuse.qtpl:15
	qw422016.N().S(`
`)
//line domainuse.qtpl:16
}

//line domainuse.qtpl:16
func WriteDomainTemplate(qq422016 qtio422016.Writer, d *pdl.Domain, domains []*pdl.Domain) {
//line domainuse.qtpl:16
	qw422016 := qt422016.AcquireWriter(qq422016)
//line domainuse.qtpl:16
	StreamDomainTemplate(qw422016, d, domains)
//line domainuse.qtpl:16
	qt422016.ReleaseWriter(qw422016)
//line domainuse.qtpl:16
}

//line domainuse.qtpl:16
func DomainTemplate(d *pdl.Domain, domains []*pdl.Domain) string {
//line domainuse.qtpl:16
	qb422016 := qt422016.AcquireByteBuffer()
//line domainuse.qtpl:16
	WriteDomainTemplate(qb422016, d, domains)
//line domainuse.qtpl:16
	qs422016 := string(qb422016.B)
//line domainuse.qtpl:16
	qt422016.ReleaseByteBuffer(qb422016)
//line domainuse.qtpl:16
	return qs422016
//line domainuse.qtpl:16
}

// CommandTemplate defines a single command.

//line domainuse.qtpl:19
func StreamCommandTemplate(qw422016 *qt422016.Writer, c *pdl.Type, d *pdl.Domain, domains []*pdl.Domain) {
//line domainuse.qtpl:20
	domainName := d.Domain.String()
	packageName := strings.ToLower(domainName)
	cmdName := CamelName(c)
	hasEmptyRet := len(c.Returns) == 0
	retTypeList := RetTypeList(c, d, domains)
	if retTypeList != "" {
		retTypeList += ", "
	}

//line domainuse.qtpl:28
	qw422016.N().S(`
`)
//line domainuse.qtpl:29
	qw422016.N().S(genutil.FormatComment(c.Description, "", cmdName+" "))
//line domainuse.qtpl:29
	qw422016.N().S(`
//
// See: `)
//line domainuse.qtpl:31
	qw422016.N().S(DocRefLink(c))
//line domainuse.qtpl:31
	if len(c.Parameters) > 0 {
//line domainuse.qtpl:31
		qw422016.N().S(`
//
// parameters:`)
//line domainuse.qtpl:33
		for _, p := range c.Parameters {
//line domainuse.qtpl:33
			qw422016.N().S(`
//  - `)
//line domainuse.qtpl:34
			qw422016.N().S(ParamDesc(p))
//line domainuse.qtpl:34
		}
//line domainuse.qtpl:34
	}
//line domainuse.qtpl:34
	if !hasEmptyRet {
//line domainuse.qtpl:34
		qw422016.N().S(`
//
// returns:`)
//line domainuse.qtpl:36
		for _, p := range c.Returns {
//line domainuse.qtpl:36
			if p.Name == Base64EncodedParamName {
//line domainuse.qtpl:36
				continue
//line domainuse.qtpl:36
			}
//line domainuse.qtpl:36
			qw422016.N().S(`
//  - `)
//line domainuse.qtpl:37
			qw422016.N().S(RetParamDesc(p))
//line domainuse.qtpl:37
		}
//line domainuse.qtpl:37
	}
//line domainuse.qtpl:37
	qw422016.N().S(`
func (do`)
//line domainuse.qtpl:38
	qw422016.N().S(domainName)
//line domainuse.qtpl:38
	qw422016.N().S(` `)
//line domainuse.qtpl:38
	qw422016.N().S(domainName)
//line domainuse.qtpl:38
	qw422016.N().S(`) `)
//line domainuse.qtpl:38
	qw422016.N().S(cmdName)
//line domainuse.qtpl:38
	qw422016.N().S(`(`)
//line domainuse.qtpl:38
	qw422016.N().S(ParamList(c, d, domains, true))
//line domainuse.qtpl:38
	qw422016.N().S(`) (`)
//line domainuse.qtpl:38
	qw422016.N().S(retTypeList)
//line domainuse.qtpl:38
	qw422016.N().S(`err error) {
	b := `)
//line domainuse.qtpl:39
	qw422016.N().S(packageName)
//line domainuse.qtpl:39
	qw422016.N().S(`.`)
//line domainuse.qtpl:39
	qw422016.N().S(cmdName)
//line domainuse.qtpl:39
	qw422016.N().S(`(`)
//line domainuse.qtpl:39
	qw422016.N().S(ArgList(c, d, domains, false))
//line domainuse.qtpl:39
	qw422016.N().S(`)`)
//line domainuse.qtpl:39
	if len(c.Parameters) != 0 {
//line domainuse.qtpl:40
		for _, p := range c.Parameters {
//line domainuse.qtpl:40
			if !p.Optional {
//line domainuse.qtpl:40
				continue
//line domainuse.qtpl:40
			}
//line domainuse.qtpl:41
			optName := OptionFuncPrefix + GoName(p, false) + OptionFuncSuffix
			v := strings.TrimSpace(GoName(p, true))

//line domainuse.qtpl:43
			qw422016.N().S(`
	if `)
//line domainuse.qtpl:44
			qw422016.N().S(v)
//line domainuse.qtpl:44
			qw422016.N().S(` != nil {
		b = b.`)
//line domainuse.qtpl:45
			qw422016.N().S(optName)
//line domainuse.qtpl:45
			if IsTypeOriginallyNilable(p, d, domains) {
//line domainuse.qtpl:45
				qw422016.N().S(`(`)
//line domainuse.qtpl:45
				qw422016.N().S(v)
//line domainuse.qtpl:45
				qw422016.N().S(`)`)
//line domainuse.qtpl:45
			} else {
//line domainuse.qtpl:45
				qw422016.N().S(`(*`)
//line domainuse.qtpl:45
				qw422016.N().S(v)
//line domainuse.qtpl:45
				qw422016.N().S(`)`)
//line domainuse.qtpl:45
			}
//line domainuse.qtpl:45
			qw422016.N().S(`
	}`)
//line domainuse.qtpl:46
		}
//line domainuse.qtpl:46
	}
//line domainuse.qtpl:46
	qw422016.N().S(`
	return b.Do(do`)
//line domainuse.qtpl:47
	qw422016.N().S(domainName)
//line domainuse.qtpl:47
	qw422016.N().S(`.ctxWithExecutor)
}
`)
//line domainuse.qtpl:49
}

//line domainuse.qtpl:49
func WriteCommandTemplate(qq422016 qtio422016.Writer, c *pdl.Type, d *pdl.Domain, domains []*pdl.Domain) {
//line domainuse.qtpl:49
	qw422016 := qt422016.AcquireWriter(qq422016)
//line domainuse.qtpl:49
	StreamCommandTemplate(qw422016, c, d, domains)
//line domainuse.qtpl:49
	qt422016.ReleaseWriter(qw422016)
//line domainuse.qtpl:49
}

//line domainuse.qtpl:49
func CommandTemplate(c *pdl.Type, d *pdl.Domain, domains []*pdl.Domain) string {
//line domainuse.qtpl:49
	qb422016 := qt422016.AcquireByteBuffer()
//line domainuse.qtpl:49
	WriteCommandTemplate(qb422016, c, d, domains)
//line domainuse.qtpl:49
	qs422016 := string(qb422016.B)
//line domainuse.qtpl:49
	qt422016.ReleaseByteBuffer(qb422016)
//line domainuse.qtpl:49
	return qs422016
//line domainuse.qtpl:49
}
