package main

import (
	"github.com/go-drift/drift/pkg/graphics"
	"github.com/go-drift/drift/pkg/theme"
)

// Color palette inspired by traditional Chinese ink painting
var (
	colorInk          = graphics.Color(0xff1a1a2e) // 墨色 deep background
	colorPaper        = graphics.Color(0xfff5f0e8) // 宣紙 light background
	colorCinnabar     = graphics.Color(0xffc73e1d) // 硃砂 primary accent
	colorGold         = graphics.Color(0xffb8860b) // 金色 secondary accent
	colorJade         = graphics.Color(0xff2d6a4f) // 玉色 success/static
	colorCloudWhite   = graphics.Color(0xfffafafa) // 雲白
	colorLightInk     = graphics.Color(0xff2d2d44) // 淡墨 surface dark
	colorMediumInk    = graphics.Color(0xff3d3d5c) // 中墨 surface variant dark
	colorWarmGray     = graphics.Color(0xffe8e0d4) // 暖灰 surface variant light
	colorDarkText     = graphics.Color(0xff1a1a1a) // 濃墨文字
	colorLightText    = graphics.Color(0xffe8e0d4) // 淡色文字
	colorSubtleLight  = graphics.Color(0xff8a8278) // 淡灰文字 light mode
	colorSubtleDark   = graphics.Color(0xff9a96a6) // 淡灰文字 dark mode
	colorErrorRed     = graphics.Color(0xffd32f2f) // 錯誤紅
	colorYangLine     = graphics.Color(0xfff0e6d2) // 陽爻顏色 (light on dark)
	colorYinLine      = graphics.Color(0xfff0e6d2) // 陰爻顏色
	colorMovingMarker = graphics.Color(0xffc73e1d) // 變爻標記
)

// DarkTheme returns a complete dark ThemeData with Chinese ink painting colors.
func DarkTheme() *theme.ThemeData {
	colors := theme.ColorScheme{
		// Primary - 硃砂
		Primary:            graphics.Color(0xffc73e1d),
		OnPrimary:          graphics.Color(0xfffafafa),
		PrimaryContainer:   graphics.Color(0xffb8860b), // 金色
		OnPrimaryContainer: graphics.Color(0xff1a1a2e),

		// Secondary - 金色
		Secondary:            graphics.Color(0xffb8860b),
		OnSecondary:          graphics.Color(0xff1a1a2e),
		SecondaryContainer:   graphics.Color(0xff5c4306),
		OnSecondaryContainer: graphics.Color(0xffe8d5a0),

		// Tertiary - 玉色
		Tertiary:            graphics.Color(0xff2d6a4f),
		OnTertiary:          graphics.Color(0xfffafafa),
		TertiaryContainer:   graphics.Color(0xff1b4030),
		OnTertiaryContainer: graphics.Color(0xffa0d4b8),

		// Surface - 墨色系
		Surface:                 graphics.Color(0xff2d2d44), // 淡墨
		OnSurface:               graphics.Color(0xffe8e0d4), // 淡色文字
		SurfaceVariant:          graphics.Color(0xff3d3d5c), // 中墨
		OnSurfaceVariant:        graphics.Color(0xff9a96a6), // 淡灰文字
		SurfaceDim:              graphics.Color(0xff1a1a2e),
		SurfaceBright:           graphics.Color(0xff4a4a6a),
		SurfaceContainerLowest:  graphics.Color(0xff141428),
		SurfaceContainerLow:     graphics.Color(0xff1e1e36),
		SurfaceContainer:        graphics.Color(0xff24243e),
		SurfaceContainerHigh:    graphics.Color(0xff2e2e4a),
		SurfaceContainerHighest: graphics.Color(0xff383856),

		// Background - 墨色
		Background:   graphics.Color(0xff1a1a2e),
		OnBackground: graphics.Color(0xffe8e0d4),

		// Error
		Error:            graphics.Color(0xffffb3af),
		OnError:          graphics.Color(0xff68000e),
		ErrorContainer:   graphics.Color(0xffd32f2f),
		OnErrorContainer: graphics.Color(0xfffafafa),

		// Outline
		Outline:        graphics.Color(0xff6a6880),
		OutlineVariant: graphics.Color(0xff3d3d5c),

		// Shadow and Scrim
		Shadow: graphics.Color(0xff000000),
		Scrim:  graphics.Color(0xff000000),

		// Inverse
		InverseSurface:   graphics.Color(0xffe8e0d4),
		OnInverseSurface: graphics.Color(0xff2d2d44),
		InversePrimary:   graphics.Color(0xff8b2010),

		SurfaceTint: graphics.Color(0xffc73e1d),
		Brightness:  theme.BrightnessDark,
	}

	return &theme.ThemeData{
		ColorScheme: colors,
		TextTheme:   theme.DefaultTextTheme(colors.OnBackground),
		Brightness:  theme.BrightnessDark,
	}
}

// LightTheme returns a complete light ThemeData with Chinese ink painting colors.
func LightTheme() *theme.ThemeData {
	colors := theme.ColorScheme{
		// Primary - 硃砂
		Primary:            graphics.Color(0xffc73e1d),
		OnPrimary:          graphics.Color(0xfffafafa),
		PrimaryContainer:   graphics.Color(0xffb8860b),
		OnPrimaryContainer: graphics.Color(0xff1a1a1a),

		// Secondary - 玉色
		Secondary:            graphics.Color(0xff2d6a4f),
		OnSecondary:          graphics.Color(0xfffafafa),
		SecondaryContainer:   graphics.Color(0xffc8e6d6),
		OnSecondaryContainer: graphics.Color(0xff1a3a28),

		// Tertiary - 金色
		Tertiary:            graphics.Color(0xffb8860b),
		OnTertiary:          graphics.Color(0xfffafafa),
		TertiaryContainer:   graphics.Color(0xfff0e0b8),
		OnTertiaryContainer: graphics.Color(0xff4a3500),

		// Surface - 宣紙系
		Surface:                 graphics.Color(0xfffafafa), // 雲白
		OnSurface:               graphics.Color(0xff1a1a1a), // 濃墨
		SurfaceVariant:          graphics.Color(0xffe8e0d4), // 暖灰
		OnSurfaceVariant:        graphics.Color(0xff8a8278), // 淡灰文字
		SurfaceDim:              graphics.Color(0xffddd5c8),
		SurfaceBright:           graphics.Color(0xfff5f0e8),
		SurfaceContainerLowest:  graphics.Color(0xffffffff),
		SurfaceContainerLow:     graphics.Color(0xfff8f3eb),
		SurfaceContainer:        graphics.Color(0xfff2ede5),
		SurfaceContainerHigh:    graphics.Color(0xffece7df),
		SurfaceContainerHighest: graphics.Color(0xffe6e1d9),

		// Background - 宣紙
		Background:   graphics.Color(0xfff5f0e8),
		OnBackground: graphics.Color(0xff1a1a1a),

		// Error
		Error:            graphics.Color(0xffd32f2f),
		OnError:          graphics.Color(0xfffafafa),
		ErrorContainer:   graphics.Color(0xffffdad6),
		OnErrorContainer: graphics.Color(0xff410002),

		// Outline
		Outline:        graphics.Color(0xff8a8278),
		OutlineVariant: graphics.Color(0xffd4cec4),

		// Shadow and Scrim
		Shadow: graphics.Color(0xff000000),
		Scrim:  graphics.Color(0xff000000),

		// Inverse
		InverseSurface:   graphics.Color(0xff2d2d44),
		OnInverseSurface: graphics.Color(0xfff5f0e8),
		InversePrimary:   graphics.Color(0xffff8a70),

		SurfaceTint: graphics.Color(0xffc73e1d),
		Brightness:  theme.BrightnessLight,
	}

	return &theme.ThemeData{
		ColorScheme: colors,
		TextTheme:   theme.DefaultTextTheme(colors.OnBackground),
		Brightness:  theme.BrightnessLight,
	}
}
