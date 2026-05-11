package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// OLCTheme — тёмная тема для OLC VPN
type OLCTheme struct{}

var _ fyne.Theme = (*OLCTheme)(nil)

func (t *OLCTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{R: 13, G: 13, B: 13, A: 255} // #0D0D0D
	case theme.ColorNameButton:
		return color.RGBA{R: 26, G: 26, B: 26, A: 255} // #1A1A1A
	case theme.ColorNameDisabledButton:
		return color.RGBA{R: 40, G: 40, B: 40, A: 255}
	case theme.ColorNamePrimary:
		return color.RGBA{R: 0, G: 229, B: 255, A: 255} // #00E5FF
	case theme.ColorNameHover:
		return color.RGBA{R: 0, G: 200, B: 220, A: 255}
	case theme.ColorNameFocus:
		return color.RGBA{R: 0, G: 229, B: 255, A: 255}
	case theme.ColorNameForeground:
		return color.RGBA{R: 255, G: 255, B: 255, A: 255} // #FFFFFF
	case theme.ColorNameDisabled:
		return color.RGBA{R: 136, G: 136, B: 136, A: 255} // #888888
	case theme.ColorNameError:
		return color.RGBA{R: 255, G: 23, B: 68, A: 255} // #FF1744
	case theme.ColorNameSuccess:
		return color.RGBA{R: 0, G: 200, B: 83, A: 255} // #00C853
	case theme.ColorNameInputBackground:
		return color.RGBA{R: 26, G: 26, B: 26, A: 255}
	case theme.ColorNameOverlayBackground:
		return color.RGBA{R: 0, G: 0, B: 0, A: 200}
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (t *OLCTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *OLCTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *OLCTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameInlineIcon:
		return 20
	case theme.SizeNameScrollBar:
		return 12
	case theme.SizeNameScrollBarSmall:
		return 6
	case theme.SizeNameText:
		return 14
	case theme.SizeNameHeadingText:
		return 20
	case theme.SizeNameSubHeadingText:
		return 16
	case theme.SizeNameCaptionText:
		return 12
	case theme.SizeNameInputBorder:
		return 1
	default:
		return theme.DefaultTheme().Size(name)
	}
}
