package main

import (
	"fmt"

	"github.com/go-drift/drift/pkg/core"
	"github.com/go-drift/drift/pkg/graphics"
	"github.com/go-drift/drift/pkg/layout"
	"github.com/go-drift/drift/pkg/navigation"
	"github.com/go-drift/drift/pkg/theme"
	"github.com/go-drift/drift/pkg/widgets"
)

// homePage is the main page showing divination history.
type homePage struct{ core.StatefulBase }

func (homePage) CreateState() core.State { return &homeState{} }

type homeState struct {
	core.StateBase
}

func (s *homeState) InitState() {
	core.UseObservable(s, recordsObservable)
}

func (s *homeState) Build(ctx core.BuildContext) core.Widget {
	colors := theme.ColorsOf(ctx)
	allRecords := recordsObservable.Value()

	// App bar
	appBar := widgets.Container{
		Color:   colors.Surface,
		Padding: layout.EdgeInsets{Top: 48, Bottom: 16, Left: 20, Right: 20},
		Child: widgets.Row{
			MainAxisAlignment:  widgets.MainAxisAlignmentSpaceBetween,
			CrossAxisAlignment: widgets.CrossAxisAlignmentCenter,
			Children: []core.Widget{
				widgets.Column{
					CrossAxisAlignment: widgets.CrossAxisAlignmentStart,
					MainAxisSize:       widgets.MainAxisSizeMin,
					Children: []core.Widget{
						widgets.Text{
							Content: "蓍草占卜",
							Style: graphics.TextStyle{
								Color:      colors.OnSurface,
								FontSize:   28,
								FontWeight: graphics.FontWeightBold,
							},
						},
						widgets.VSpace(2),
						widgets.Text{
							Content: fmt.Sprintf("共 %d 筆記錄", len(allRecords)),
							Style: graphics.TextStyle{
								Color:    colors.OnSurfaceVariant,
								FontSize: 13,
							},
						},
					},
				},
				// 起卦 button
				widgets.GestureDetector{
					OnTap: func() {
						if nav := navigation.NavigatorOf(ctx); nav != nil {
							nav.PushNamed("/divine", nil)
						}
					},
					Child: widgets.Container{
						Color:        colors.Primary,
						BorderRadius: 16,
						Padding:      layout.EdgeInsetsSymmetric(20, 12),
						Child: widgets.Text{
							Content: "起卦",
							Style: graphics.TextStyle{
								Color:      colors.OnPrimary,
								FontSize:   16,
								FontWeight: graphics.FontWeightBold,
							},
						},
					},
				},
			},
		},
	}

	// Content
	var content core.Widget
	if len(allRecords) == 0 {
		content = s.buildEmptyState(ctx, colors)
	} else {
		content = s.buildRecordList(ctx, colors, allRecords)
	}

	return widgets.Column{
		CrossAxisAlignment: widgets.CrossAxisAlignmentStretch,
		Children: []core.Widget{
			appBar,
			widgets.Expanded{Child: content},
		},
	}
}

func (s *homeState) buildEmptyState(ctx core.BuildContext, colors theme.ColorScheme) core.Widget {
	return widgets.Container{
		Color: colors.Background,
		Child: widgets.Column{
			MainAxisAlignment:  widgets.MainAxisAlignmentCenter,
			CrossAxisAlignment: widgets.CrossAxisAlignmentCenter,
			Children: []core.Widget{
				widgets.Text{
					Content: "大衍之數五十",
					Style: graphics.TextStyle{
						Color:    colors.OnSurfaceVariant,
						FontSize: 24,
					},
				},
				widgets.VSpace(8),
				widgets.Text{
					Content: "其用四十有九",
					Style: graphics.TextStyle{
						Color:    colors.OnSurfaceVariant,
						FontSize: 18,
					},
				},
				widgets.VSpace(24),
				widgets.Text{
					Content: "點擊「起卦」開始占卜",
					Style: graphics.TextStyle{
						Color:    colors.OnSurfaceVariant,
						FontSize: 14,
					},
				},
			},
		},
	}
}

func (s *homeState) buildRecordList(ctx core.BuildContext, colors theme.ColorScheme, records []*DivinationRecord) core.Widget {
	return widgets.Container{
		Color: colors.Background,
		Child: widgets.ListViewBuilder{
			Padding:    layout.EdgeInsets{Left: 16, Right: 16, Top: 8, Bottom: 24},
			ItemCount:  len(records),
			ItemExtent: 96,
			ItemBuilder: func(ctx core.BuildContext, index int) core.Widget {
				record := records[index]
				return widgets.Padding{
					Padding: layout.EdgeInsets{Bottom: 8},
					Child:   s.buildRecordCard(ctx, colors, record),
				}
			},
		},
	}
}

func (s *homeState) buildRecordCard(ctx core.BuildContext, colors theme.ColorScheme, record *DivinationRecord) core.Widget {
	hex := record.ToHexagram()

	// Build mini hexagram display (6 lines, top to bottom)
	miniLines := make([]core.Widget, 0, 6)
	for i := 5; i >= 0; i-- {
		lineColor := colors.OnSurfaceVariant
		if hex.Original[i].IsMoving() {
			lineColor = colorMovingMarker
		}
		miniLines = append(miniLines, widgets.Text{
			Content: hex.Original[i].Symbol(),
			Style: graphics.TextStyle{
				Color:    lineColor,
				FontSize: 8,
			},
		})
	}

	// Question text (truncated)
	questionText := record.Question
	if questionText == "" {
		questionText = "無事由"
	}
	if len([]rune(questionText)) > 20 {
		questionText = string([]rune(questionText)[:20]) + "…"
	}

	return widgets.GestureDetector{
		OnTap: func() {
			if nav := navigation.NavigatorOf(ctx); nav != nil {
				nav.PushNamed("/result/"+record.ID, nil)
			}
		},
		Child: widgets.Container{
			Color:        colors.Surface,
			BorderRadius: 12,
			Padding:      layout.EdgeInsetsAll(14),
			Child: widgets.Row{
				CrossAxisAlignment: widgets.CrossAxisAlignmentCenter,
				Children: []core.Widget{
					// Mini hexagram
					widgets.Container{
						Width:        48,
						Height:       64,
						Color:        colors.SurfaceVariant,
						BorderRadius: 8,
						Padding:      layout.EdgeInsetsAll(6),
						Alignment:    layout.AlignmentCenter,
						Child: widgets.Column{
							MainAxisAlignment: widgets.MainAxisAlignmentCenter,
							MainAxisSize:      widgets.MainAxisSizeMin,
							Children:          miniLines,
						},
					},
					widgets.HSpace(14),
					// Text info
					widgets.Expanded{
						Child: widgets.Column{
							CrossAxisAlignment: widgets.CrossAxisAlignmentStart,
							MainAxisSize:       widgets.MainAxisSizeMin,
							Children: []core.Widget{
								widgets.Text{
									Content: questionText,
									Style: graphics.TextStyle{
										Color:      colors.OnSurface,
										FontSize:   16,
										FontWeight: graphics.FontWeightMedium,
									},
								},
								widgets.VSpace(4),
								widgets.Text{
									Content: record.CreatedAt.Format("2006-01-02 15:04"),
									Style: graphics.TextStyle{
										Color:    colors.OnSurfaceVariant,
										FontSize: 12,
									},
								},
								widgets.VSpace(2),
								widgets.Text{
									Content: record.Interpret,
									Style: graphics.TextStyle{
										Color:    colors.OnSurfaceVariant,
										FontSize: 11,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
