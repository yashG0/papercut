package service

// RGB defines color components
type RGB struct {
	R, G, B int
}

// ThemeConfig defines visual rules for document rendering
type ThemeConfig struct {
	FontFamily  string
	MarginLeft  float64
	MarginTop   float64
	MarginRight float64
	HeaderColor RGB
	BodyColor   RGB
	AccentColor RGB
	RuleColor   RGB
}

// GetTheme returns the palette and typography matching the requested theme key
func GetTheme(name string) ThemeConfig {
	switch name {
	case "invoice":
		return ThemeConfig{
			FontFamily:  "Courier",
			MarginLeft:  18,
			MarginTop:   18,
			MarginRight: 18,
			HeaderColor: RGB{R: 24, G: 43, B: 73},  // Deep Navy
			BodyColor:   RGB{R: 51, G: 51, B: 51},  // Charcoal
			AccentColor: RGB{R: 37, G: 99, B: 235}, // Royal Blue
			RuleColor:   RGB{R: 37, G: 99, B: 235}, // Royal Blue
		}
	case "modern":
		return ThemeConfig{
			FontFamily:  "Helvetica",
			MarginLeft:  25,
			MarginTop:   25,
			MarginRight: 25,
			HeaderColor: RGB{R: 15, G: 23, B: 42},    // Slate 900
			BodyColor:   RGB{R: 71, G: 85, B: 105},   // Slate 600
			AccentColor: RGB{R: 13, G: 148, B: 136},  // Emerald / Teal
			RuleColor:   RGB{R: 226, G: 232, B: 240}, // Light Slate
		}
	default: // "default"
		return ThemeConfig{
			FontFamily:  "Helvetica",
			MarginLeft:  20,
			MarginTop:   20,
			MarginRight: 20,
			HeaderColor: RGB{R: 30, G: 30, B: 30},    // Dark Gray
			BodyColor:   RGB{R: 60, G: 60, B: 60},    // Medium Gray
			AccentColor: RGB{R: 80, G: 80, B: 80},    // Neutral Gray
			RuleColor:   RGB{R: 210, G: 210, B: 210}, // Light Gray
		}
	}
}
