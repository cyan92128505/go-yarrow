package main

import (
	"fmt"

	"github.com/go-drift/drift/pkg/core"
	"github.com/go-drift/drift/pkg/graphics"
	"github.com/go-drift/drift/pkg/layout"
	"github.com/go-drift/drift/pkg/navigation"
	"github.com/go-drift/drift/pkg/platform"
	"github.com/go-drift/drift/pkg/theme"
	"github.com/go-drift/drift/pkg/widgets"
)

// resultPage displays the full hexagram result.
type resultPage struct {
	core.StatefulBase
	RecordID string
}

func (r resultPage) CreateState() core.State {
	return &resultState{recordID: r.RecordID}
}

type resultState struct {
	core.StateBase
	recordID string
}

func (s *resultState) Build(ctx core.BuildContext) core.Widget {
	colors := theme.ColorsOf(ctx)

	record := store.FindByID(s.recordID)
	if record == nil {
		return widgets.Container{
			Color:     colors.Background,
			Alignment: layout.AlignmentCenter,
			Child: widgets.Text{
				Content: "記錄不存在",
				Style:   graphics.TextStyle{Color: colors.OnSurfaceVariant, FontSize: 16},
			},
		}
	}

	hex := record.ToHexagram()

	return widgets.Container{
		Color: colors.Background,
		Child: widgets.Column{
			CrossAxisAlignment: widgets.CrossAxisAlignmentStretch,
			Children: []core.Widget{
				// App bar
				s.buildAppBar(ctx, colors, record),
				// Scrollable content
				widgets.Expanded{
					Child: widgets.ListView{
						Padding: layout.EdgeInsets{Left: 20, Right: 20, Top: 8, Bottom: 40},
						Physics: widgets.BouncingScrollPhysics{},
						Children: []core.Widget{
							// Question
							s.buildQuestionSection(colors, record),
							widgets.VSpace(24),
							// 本卦
							s.buildHexagramSection(ctx, colors, "本卦", hex.Original, hex.MovingLines),
							widgets.VSpace(24),
							// 之卦 (only show if there are moving lines)
							s.buildChangedSection(ctx, colors, hex),
							// 解卦
							s.buildInterpretSection(colors, record),
							widgets.VSpace(24),
							// 數值記錄
							s.buildRawDataSection(colors, hex),
							widgets.VSpace(24),
							// Action buttons
							s.buildActions(ctx, colors, record, hex),
						},
					},
				},
			},
		},
	}
}

func (s *resultState) buildAppBar(ctx core.BuildContext, colors theme.ColorScheme, record *DivinationRecord) core.Widget {
	return widgets.Container{
		Color:   colors.Surface,
		Padding: layout.EdgeInsets{Top: 48, Bottom: 16, Left: 20, Right: 20},
		Child: widgets.Row{
			MainAxisAlignment:  widgets.MainAxisAlignmentSpaceBetween,
			CrossAxisAlignment: widgets.CrossAxisAlignmentCenter,
			Children: []core.Widget{
				widgets.GestureDetector{
					OnTap: func() {

					},
					Child: widgets.Container{
						Padding: layout.EdgeInsets{Right: 16, Top: 4, Bottom: 4},
						Child: widgets.Text{
							Content: "← 返回",
							Style: graphics.TextStyle{
								Color:    colors.Primary,
								FontSize: 16,
							},
						},
					},
				},
				widgets.Text{
					Content: record.CreatedAt.Format("2006-01-02 15:04"),
					Style: graphics.TextStyle{
						Color:    colors.OnSurfaceVariant,
						FontSize: 13,
					},
				},
			},
		},
	}
}

func (s *resultState) buildQuestionSection(colors theme.ColorScheme, record *DivinationRecord) core.Widget {
	question := record.Question
	if question == "" {
		question = "（未記錄事由）"
	}

	return widgets.Container{
		Color:        colors.Surface,
		BorderRadius: 12,
		Padding:      layout.EdgeInsetsAll(16),
		Child: widgets.Column{
			CrossAxisAlignment: widgets.CrossAxisAlignmentStart,
			MainAxisSize:       widgets.MainAxisSizeMin,
			Children: []core.Widget{
				widgets.Text{
					Content: "占卜事由",
					Style: graphics.TextStyle{
						Color:      colors.OnSurfaceVariant,
						FontSize:   12,
						FontWeight: graphics.FontWeightMedium,
					},
				},
				widgets.VSpace(6),
				widgets.Text{
					Content: question,
					Style: graphics.TextStyle{
						Color:      colors.OnSurface,
						FontSize:   18,
						FontWeight: graphics.FontWeightMedium,
					},
				},
			},
		},
	}
}

func (s *resultState) buildHexagramSection(ctx core.BuildContext, colors theme.ColorScheme, title string, lines []Line, movingLines []int) core.Widget {
	movingSet := make(map[int]bool)
	for _, pos := range movingLines {
		movingSet[pos] = true
	}

	lineWidgets := make([]core.Widget, 0, 12)
	// Display top to bottom (line 6 -> line 1)
	for i := 5; i >= 0; i-- {
		lineNum := i + 1
		line := lines[i]
		isMoving := movingSet[lineNum]

		lineWidgets = append(lineWidgets, s.buildLineRow(colors, lineNum, line, isMoving))
		if i > 0 {
			lineWidgets = append(lineWidgets, widgets.VSpace(4))
		}
	}

	return widgets.Container{
		Color:        colors.Surface,
		BorderRadius: 12,
		Padding:      layout.EdgeInsetsAll(16),
		Child: widgets.Column{
			CrossAxisAlignment: widgets.CrossAxisAlignmentStart,
			MainAxisSize:       widgets.MainAxisSizeMin,
			Children: append([]core.Widget{
				widgets.Text{
					Content: title,
					Style: graphics.TextStyle{
						Color:      colors.OnSurfaceVariant,
						FontSize:   12,
						FontWeight: graphics.FontWeightMedium,
					},
				},
				widgets.VSpace(12),
			}, lineWidgets...),
		},
	}
}

func (s *resultState) buildLineRow(colors theme.ColorScheme, lineNum int, line Line, isMoving bool) core.Widget {
	symbolColor := colorYangLine
	if !line.IsYang() {
		symbolColor = colorYinLine
	}

	// Moving line marker
	markerText := "  "
	markerColor := colors.Background
	if isMoving {
		markerText = "◯"
		markerColor = colorMovingMarker
	}

	return widgets.Row{
		CrossAxisAlignment: widgets.CrossAxisAlignmentCenter,
		Children: []core.Widget{
			// Line number
			widgets.SizedBox{
				Width: 36,
				Child: widgets.Text{
					Content: fmt.Sprintf("第%d爻", lineNum),
					Align:   graphics.TextAlignRight,
					Style: graphics.TextStyle{
						Color:    colors.OnSurfaceVariant,
						FontSize: 11,
					},
				},
			},
			widgets.HSpace(10),
			// Line symbol (centered, monospaced-like)
			widgets.SizedBox{
				Width: 80,
				Child: widgets.Text{
					Content: line.Symbol(),
					Align:   graphics.TextAlignCenter,
					Style: graphics.TextStyle{
						Color:      symbolColor,
						FontSize:   18,
						FontWeight: graphics.FontWeightBold,
					},
				},
			},
			widgets.HSpace(10),
			// Label
			widgets.Text{
				Content: line.Label(),
				Style: graphics.TextStyle{
					Color:    colors.OnSurfaceVariant,
					FontSize: 12,
				},
			},
			widgets.HSpace(8),
			// Moving marker
			widgets.Text{
				Content: markerText,
				Style: graphics.TextStyle{
					Color:      markerColor,
					FontSize:   14,
					FontWeight: graphics.FontWeightBold,
				},
			},
		},
	}
}

func (s *resultState) buildChangedSection(ctx core.BuildContext, colors theme.ColorScheme, hex *Hexagram) core.Widget {
	if len(hex.MovingLines) == 0 {
		return widgets.VSpace(0)
	}
	return widgets.Column{
		MainAxisSize: widgets.MainAxisSizeMin,
		Children: []core.Widget{
			s.buildHexagramSection(ctx, colors, "之卦（變卦）", hex.Changed, nil),
			widgets.VSpace(24),
		},
	}
}

func (s *resultState) buildInterpretSection(colors theme.ColorScheme, record *DivinationRecord) core.Widget {
	return widgets.Container{
		Color:        colors.Primary,
		BorderRadius: 12,
		Padding:      layout.EdgeInsetsAll(16),
		Child: widgets.Column{
			CrossAxisAlignment: widgets.CrossAxisAlignmentStart,
			MainAxisSize:       widgets.MainAxisSizeMin,
			Children: []core.Widget{
				widgets.Text{
					Content: "朱熹法解卦",
					Style: graphics.TextStyle{
						Color:      colors.OnPrimary,
						FontSize:   12,
						FontWeight: graphics.FontWeightMedium,
					},
				},
				widgets.VSpace(8),
				widgets.Text{
					Content: record.Interpret,
					Style: graphics.TextStyle{
						Color:      colors.OnPrimary,
						FontSize:   16,
						FontWeight: graphics.FontWeightBold,
					},
				},
			},
		},
	}
}

func (s *resultState) buildRawDataSection(colors theme.ColorScheme, hex *Hexagram) core.Widget {
	origStr := ""
	changedStr := ""
	for i := 0; i < 6; i++ {
		if i > 0 {
			origStr += " "
			changedStr += " "
		}
		origStr += fmt.Sprintf("%d", hex.Original[i])
		changedStr += fmt.Sprintf("%d", hex.Changed[i])
	}

	return widgets.Container{
		Color:        colors.SurfaceVariant,
		BorderRadius: 12,
		Padding:      layout.EdgeInsetsAll(16),
		Child: widgets.Column{
			CrossAxisAlignment: widgets.CrossAxisAlignmentStart,
			MainAxisSize:       widgets.MainAxisSizeMin,
			Children: []core.Widget{
				widgets.Text{
					Content: "數值（由下到上）",
					Style: graphics.TextStyle{
						Color:      colors.OnSurfaceVariant,
						FontSize:   12,
						FontWeight: graphics.FontWeightMedium,
					},
				},
				widgets.VSpace(6),
				widgets.Text{
					Content: "本卦：" + origStr,
					Style: graphics.TextStyle{
						Color:    colors.OnSurface,
						FontSize: 14,
					},
				},
				widgets.VSpace(2),
				widgets.Text{
					Content: "之卦：" + changedStr,
					Style: graphics.TextStyle{
						Color:    colors.OnSurface,
						FontSize: 14,
					},
				},
			},
		},
	}
}

func (s *resultState) buildActions(ctx core.BuildContext, colors theme.ColorScheme, record *DivinationRecord, hex *Hexagram) core.Widget {
	return widgets.Row{
		MainAxisAlignment: widgets.MainAxisAlignmentCenter,
		Children: []core.Widget{
			// Share button
			widgets.GestureDetector{
				OnTap: func() {
					shareText := hex.ToShareText(record.Question, record.CreatedAt)
					platform.Share.ShareText(shareText)
				},
				Child: widgets.Container{
					Color:        colors.Surface,
					BorderRadius: 12,
					Padding:      layout.EdgeInsetsSymmetric(24, 14),
					Child: widgets.Text{
						Content: "分享",
						Style: graphics.TextStyle{
							Color:      colors.Primary,
							FontSize:   15,
							FontWeight: graphics.FontWeightMedium,
						},
					},
				},
			},
			widgets.HSpace(12),
			// Delete button
			widgets.GestureDetector{
				OnTap: func() {
					s.confirmDelete(ctx, colors, record)
				},
				Child: widgets.Container{
					Color:        colors.Surface,
					BorderRadius: 12,
					Padding:      layout.EdgeInsetsSymmetric(24, 14),
					Child: widgets.Text{
						Content: "刪除",
						Style: graphics.TextStyle{
							Color:      colors.Error,
							FontSize:   15,
							FontWeight: graphics.FontWeightMedium,
						},
					},
				},
			},
		},
	}
}

func (s *resultState) confirmDelete(ctx core.BuildContext, colors theme.ColorScheme, record *DivinationRecord) {
	navigation.ShowModalBottomSheet(ctx, func(ctx core.BuildContext) core.Widget {
		return widgets.Padding{
			Padding: layout.EdgeInsetsAll(24),
			Child: widgets.Column{
				MainAxisSize:       widgets.MainAxisSizeMin,
				CrossAxisAlignment: widgets.CrossAxisAlignmentStretch,
				Children: []core.Widget{
					widgets.Text{
						Content: "確認刪除",
						Style: graphics.TextStyle{
							Color:      colors.Error,
							FontSize:   20,
							FontWeight: graphics.FontWeightBold,
						},
					},
					widgets.VSpace(12),
					widgets.Text{
						Content: "確定要刪除這筆占卜記錄嗎？此操作無法復原。",
						Style: graphics.TextStyle{
							Color:    colors.OnSurfaceVariant,
							FontSize: 14,
						},
					},
					widgets.VSpace(20),
					widgets.Row{
						MainAxisAlignment: widgets.MainAxisAlignmentCenter,
						Children: []core.Widget{
							widgets.GestureDetector{
								OnTap: func() {
									widgets.BottomSheetScope{}.Of(ctx).Close(nil)
								},
								Child: widgets.Container{
									Color:        colors.SurfaceVariant,
									BorderRadius: 8,
									Padding:      layout.EdgeInsetsSymmetric(24, 12),
									Child: widgets.Text{
										Content: "取消",
										Style: graphics.TextStyle{
											Color:    colors.OnSurfaceVariant,
											FontSize: 15,
										},
									},
								},
							},
							widgets.HSpace(12),
							widgets.GestureDetector{
								OnTap: func() {
									store.Delete(record.ID)
									widgets.BottomSheetScope{}.Of(ctx).Close(nil)
									platform.Haptics.HeavyImpact()
									// Go back to home

									if nav := navigation.NavigatorOf(ctx); nav != nil {
										nav.Pop(nil)
									}
								},
								Child: widgets.Container{
									Color:        colors.Error,
									BorderRadius: 8,
									Padding:      layout.EdgeInsetsSymmetric(24, 12),
									Child: widgets.Text{
										Content: "刪除",
										Style: graphics.TextStyle{
											Color:      colors.OnError,
											FontSize:   15,
											FontWeight: graphics.FontWeightBold,
										},
									},
								},
							},
						},
					},
				},
			},
		}
	})
}
