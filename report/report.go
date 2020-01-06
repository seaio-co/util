package report

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

type Options struct {
	ShortRange      bool
	FilterGenerated bool
	Fixes           []analysis.SuggestedFix
	Related         []analysis.RelatedInformation
}

type Option func(*Options)

func ShortRange() Option {
	return func(opts *Options) {
		opts.ShortRange = true
	}
}

func FilterGenerated() Option {
	return func(opts *Options) {
		opts.FilterGenerated = true
	}
}

func Fixes(fixes ...analysis.SuggestedFix) Option {
	return func(opts *Options) {
		opts.Fixes = append(opts.Fixes, fixes...)
	}
}

func Related(node Positioner, message string) Option {
	return func(opts *Options) {
		pos, end := getRange(node, opts.ShortRange)
		r := analysis.RelatedInformation{
			Pos:     pos,
			End:     end,
			Message: message,
		}
		opts.Related = append(opts.Related, r)
	}
}

type Positioner interface {
	Pos() token.Pos
}

type fullPositioner interface {
	Pos() token.Pos
	End() token.Pos
}

type sourcer interface {
	Source() ast.Node
}
