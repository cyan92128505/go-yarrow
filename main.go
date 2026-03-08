package main

import (
	"context"
	"log"
	"strings"

	"github.com/go-drift/drift/pkg/core"
	"github.com/go-drift/drift/pkg/drift"
	"github.com/go-drift/drift/pkg/engine"
	"github.com/go-drift/drift/pkg/graphics"
	"github.com/go-drift/drift/pkg/layout"
	"github.com/go-drift/drift/pkg/navigation"
	"github.com/go-drift/drift/pkg/platform"
	"github.com/go-drift/drift/pkg/theme"
	"github.com/go-drift/drift/pkg/widgets"
)

// Global record store, initialized in OnInit.
var store *RecordStore

// yarrowApp is the root widget.
type yarrowApp struct{ core.StatefulBase }

func (yarrowApp) CreateState() core.State { return &yarrowAppState{} }

type yarrowAppState struct {
	core.StateBase
	isDark          bool
	isCupertino     bool
	cachedThemeData *theme.AppThemeData
}

func (s *yarrowAppState) InitState() {
	s.isDark = true
	s.updateSystemUI()
}

func (s *yarrowAppState) Build(ctx core.BuildContext) core.Widget {
	appThemeData := s.getThemeData()
	routes := s.buildRoutes()

	router := navigation.Router{
		InitialPath:  "/",
		Routes:       routes,
		ErrorBuilder: buildErrorPage,
	}

	return theme.AppTheme{
		Data:  appThemeData,
		Child: router,
	}
}

func (s *yarrowAppState) buildRoutes() []navigation.ScreenRoute {
	return []navigation.ScreenRoute{
		{
			Path: "/",
			Screen: func(ctx core.BuildContext, settings navigation.RouteSettings) core.Widget {
				return homePage{}
			},
		},
		{
			Path: "/divine",
			Screen: func(ctx core.BuildContext, settings navigation.RouteSettings) core.Widget {
				return divinePage{}
			},
		},
		{
			Path: "/result/:id",
			Screen: func(ctx core.BuildContext, settings navigation.RouteSettings) core.Widget {
				// Extract ID from path
				id := extractIDFromPath(settings.Name)
				return resultPage{RecordID: id}
			},
		},
		{
			Path: "/test",
			Screen: func(ctx core.BuildContext, settings navigation.RouteSettings) core.Widget {
				return testFormPage{}
			},
		},
	}
}

// extractIDFromPath pulls the record ID from "/result/123456".
func extractIDFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 2 {
		return parts[len(parts)-1]
	}
	return ""
}

func (s *yarrowAppState) getThemeData() *theme.AppThemeData {
	targetPlatform := theme.TargetPlatformMaterial
	if s.isCupertino {
		targetPlatform = theme.TargetPlatformCupertino
	}

	brightness := theme.BrightnessLight
	if s.isDark {
		brightness = theme.BrightnessDark
	}

	// Only recreate if values changed
	if s.cachedThemeData == nil ||
		s.cachedThemeData.Platform != targetPlatform ||
		s.cachedThemeData.Brightness() != brightness {

		// Create new theme data
		var material *theme.ThemeData
		var cupertino *theme.CupertinoThemeData

		if s.isDark {
			// Use dark theme
			material = DarkTheme()
			cupertino = theme.DefaultCupertinoDarkTheme()
		} else {
			// Use light theme for light mode
			material = LightTheme()
			cupertino = theme.DefaultCupertinoLightTheme()
		}

		s.cachedThemeData = &theme.AppThemeData{
			Platform:  targetPlatform,
			Material:  material,
			Cupertino: cupertino,
		}
	}
	return s.cachedThemeData
}

func (s *yarrowAppState) updateSystemUI() {
	themeData := s.getThemeData()
	statusStyle := platform.StatusBarStyleDark
	if s.isDark {
		statusStyle = platform.StatusBarStyleLight
	}
	bgColor := themeData.Material.ColorScheme.Surface
	engine.SetBackgroundColor(graphics.Color(themeData.Material.ColorScheme.Background))
	_ = platform.SetSystemUI(platform.SystemUIStyle{
		StatusBarStyle:  statusStyle,
		BackgroundColor: &bgColor,
		Transparent:     true,
	})
}

func buildErrorPage(ctx core.BuildContext, settings navigation.RouteSettings) core.Widget {
	colors := theme.ColorsOf(ctx)
	return widgets.Container{
		Color:     colors.Background,
		Alignment: layout.AlignmentCenter,
		Child: widgets.Text{
			Content: "頁面不存在: " + settings.Name,
			Style: graphics.TextStyle{
				Color:    colors.OnSurfaceVariant,
				FontSize: 16,
			},
		},
	}
}

func main() {
	app := drift.NewApp(yarrowApp{})

	app.OnInit = func(ctx context.Context) error {
		store = NewRecordStore()
		log.Println("蓍草占卜 App 初始化完成")
		return nil
	}

	app.OnDispose = func() {
		log.Println("蓍草占卜 App 關閉")
	}

	app.Run()
}
