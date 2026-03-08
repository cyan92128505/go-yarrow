package main

import (
	"github.com/go-drift/drift/pkg/core"
	"github.com/go-drift/drift/pkg/theme"
	"github.com/go-drift/drift/pkg/widgets"
)

type testFormPage struct{ core.StatefulBase }

func (testFormPage) CreateState() core.State { return &testFormState{} }

type testFormState struct {
	core.StateBase
}

func (s *testFormState) Build(ctx core.BuildContext) core.Widget {
	return widgets.Form{
		Autovalidate: true,
		Child:        testFormContent{},
	}
}

type testFormContent struct{ core.StatelessBase }

func (f testFormContent) Build(ctx core.BuildContext) core.Widget {
	return widgets.Column{
		MainAxisSize: widgets.MainAxisSizeMin,
		Children: []core.Widget{
			theme.TextFormFieldOf(ctx).
				WithLabel("Test").
				WithPlaceholder("Type here"),
		},
	}
}
