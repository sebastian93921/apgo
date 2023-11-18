package system

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type MyTheme struct{}

// Color implements fyne.Theme.
func (*MyTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameDisabled {
		if variant == theme.VariantLight {
			return color.Black
		}
		return color.White
	}

	return theme.DefaultTheme().Color(name, variant)
}

// Font implements fyne.Theme.
func (*MyTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Icon implements fyne.Theme.
func (*MyTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size implements fyne.Theme.
func (*MyTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

var _ fyne.Theme = (*MyTheme)(nil)
