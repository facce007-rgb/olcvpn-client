package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// OLCTheme — тёмная тема в стиле v2RayTun/Material Design
type OLCTheme struct{}

var _ fyne.Theme = (*OLCTheme)(nil)

// Color возвращает цвет для указанного имени
func (t *OLCTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	// Используем тёмную тему как базу
	if variant == theme.VariantLight {
		variant = theme.VariantDark
	}

	switch name {
	// Основные цвета
	case theme.ColorNameBackground:
		return color.NRGBA{R: 18, G: 18, B: 18, A: 255} // #121212 - Material Dark
	case theme.ColorNameForeground:
		return color.NRGBA{R: 255, G: 255, B: 255, A: 255}

	// Акцентный цвет - голубой как в v2RayTun
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 0, G: 229, B: 255, A: 255} // #00E5FF - Cyan A400

	// Кнопки
	case theme.ColorNameButton:
		return color.NRGBA{R: 33, G: 33, B: 33, A: 255} // #212121
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 66, G: 66, B: 66, A: 255}

	// Поверхности (карточки)
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 30, G: 30, B: 30, A: 255} // #1E1E1E

	// Успех/Ошибка
	case theme.ColorNameSuccess:
		return color.NRGBA{R: 76, G: 175, B: 80, A: 255} // #4CAF50 - Green
	case theme.ColorNameError:
		return color.NRGBA{R: 244, G: 67, B: 54, A: 255} // #F44336 - Red
	case theme.ColorNameWarning:
		return color.NRGBA{R: 255, G: 193, B: 7, A: 255} // #FFC107 - Amber

	// Разделители
	case theme.ColorNameSeparator:
		return color.NRGBA{R: 66, G: 66, B: 66, A: 255}

	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

// Font возвращает шрифт для указанного стиля
func (t *OLCTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Icon возвращает иконку для указанного имени
func (t *OLCTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size возвращает размер для указанного имени
func (t *OLCTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameInlineIcon:
		return 24
	case theme.SizeNameScrollBar:
		return 12
	case theme.SizeNameScrollBarSmall:
		return 6
	case theme.SizeNameSeparatorThickness:
		return 1
	case theme.SizeNameText:
		return 14
	case theme.SizeNameHeadingText:
		return 20
	case theme.SizeNameSubHeadingText:
		return 16
	case theme.SizeNameCaptionText:
		return 12
	case theme.SizeNameInputBorder:
		return 2
	default:
		return theme.DefaultTheme().Size(name)
	}
}
