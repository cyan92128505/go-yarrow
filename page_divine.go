package main

import (
	"time"

	"github.com/go-drift/drift/pkg/core"
	"github.com/go-drift/drift/pkg/graphics"
	"github.com/go-drift/drift/pkg/layout"
	"github.com/go-drift/drift/pkg/navigation"
	"github.com/go-drift/drift/pkg/platform"
	"github.com/go-drift/drift/pkg/theme"
	"github.com/go-drift/drift/pkg/widgets"
)

// divinePage handles question input and triggers divination.
type divinePage struct{ core.StatefulBase }

func (divinePage) CreateState() core.State { return &divineState{} }

func buildDivinePage(_ core.BuildContext) core.Widget { return divinePage{} }

// formData holds the collected form values after validation.
type formData struct {
	Question string
}

type divineState struct {
	core.StateBase
	data formData
}

func (s *divineState) InitState() {
	s.data = formData{Question: ""}
}

func (s *divineState) Build(ctx core.BuildContext) core.Widget {
	colors := theme.ColorsOf(ctx)
	textTheme := theme.TextThemeOf(ctx)

	// Back button — same pattern as showcase pageScaffold
	backButton := widgets.Button{
		Label: "← 返回",
		OnTap: func() {
			if nav := navigation.NavigatorOf(ctx); nav != nil {
				nav.Pop(nil)
			}
		},
		Color:        colors.SurfaceContainerHigh,
		TextColor:    colors.OnSurface,
		Padding:      layout.EdgeInsetsSymmetric(16, 10),
		BorderRadius: 8,
		FontSize:     14,
		Haptic:       true,
	}

	// Header — same as showcase pageScaffold
	headerPadding := widgets.SafeAreaPadding(ctx).OnlyTop().Add(16)
	header := widgets.Container{
		Color: colors.Surface,
		Child: widgets.Padding{
			Padding: headerPadding,
			Child: widgets.Row{
				CrossAxisAlignment: widgets.CrossAxisAlignmentCenter,
				Children: []core.Widget{
					backButton,
					widgets.HSpace(16),
					widgets.Text{Content: "起卦", Style: textTheme.HeadlineMedium},
				},
			},
		},
	}

	// Content in ScrollView — same structure as showcase demoPage
	content := widgets.ScrollView{
		ScrollDirection: widgets.AxisVertical,
		Physics:         widgets.BouncingScrollPhysics{},
		Padding:         layout.EdgeInsetsAll(20),
		Child: widgets.Column{
			MainAxisSize: widgets.MainAxisSizeMin,
			Children: []core.Widget{
				widgets.VSpace(20),
				widgets.Text{
					Content: "分二以象兩，掛一以象三",
					Align:   graphics.TextAlignCenter,
					Style: graphics.TextStyle{
						Color:    colors.OnSurfaceVariant,
						FontSize: 16,
					},
				},
				widgets.VSpace(4),
				widgets.Text{
					Content: "揲之以四以象四時，歸奇於扐以象閏",
					Align:   graphics.TextAlignCenter,
					Style: graphics.TextStyle{
						Color:    colors.OnSurfaceVariant,
						FontSize: 14,
					},
				},
				widgets.VSpace(40),
				widgets.Text{
					Content: "占卜事由",
					Style: graphics.TextStyle{
						Color:      colors.OnSurface,
						FontSize:   14,
						FontWeight: graphics.FontWeightMedium,
					},
				},
				widgets.VSpace(8),
				// Form + formContent — exact showcase pattern
				widgets.Form{
					Autovalidate: false,
					Child:        divineFormContent{state: s},
				},
				widgets.VSpace(40),
			},
		},
	}

	// Assemble with pageScaffold structure
	return widgets.Expanded{
		Child: widgets.Container{
			Color: colors.Background,
			Child: widgets.Column{
				Children: []core.Widget{
					header,
					widgets.Expanded{Child: content},
				},
			},
		},
	}
}

// divineFormContent — same pattern as showcase formContent
type divineFormContent struct {
	core.StatelessBase
	state *divineState
}

func (f divineFormContent) Build(ctx core.BuildContext) core.Widget {
	form := widgets.FormOf(ctx)

	return widgets.Column{
		CrossAxisAlignment: widgets.CrossAxisAlignmentStretch,
		MainAxisSize:       widgets.MainAxisSizeMin,
		Children: []core.Widget{
			theme.TextFormFieldOf(ctx).
				WithPlaceholder("心中所想之事（可留空）").
				WithOnSaved(func(value string) {
					f.state.data.Question = value
				}),
			widgets.VSpace(32),
			// Submit button
			theme.ButtonOf(ctx, "起卦", func() {
				if form != nil {
					form.Save()
				}
				f.state.performDivination(ctx)
			}),
		},
	}
}

func (s *divineState) performDivination(ctx core.BuildContext) {
	platform.Haptics.MediumImpact()

	seed := time.Now().UnixNano()
	hex := generateHexagram(seed)

	record := NewRecordFromHexagram(s.data.Question, seed, hex)
	if err := store.Add(record); err != nil {
		_ = err
	}

	if nav := navigation.NavigatorOf(ctx); nav != nil {
		nav.PushNamed("/result/"+record.ID, nil)
	}
}
