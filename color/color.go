package color

type Color uint8

const (
	Invalid Color = iota

	FgBlack
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgPurple
	FgCyan
	FgWhite

	FgBoldBlack
	FgBoldRed
	FgBoldGreen
	FgBoldYellow
	FgBoldBlue
	FgBoldPurple
	FgBoldCyan
	FgBoldWhite

	FgUnderlineBlack
	FgUnderlineRed
	FgUnderlineGreen
	FgUnderlineYellow
	FgUnderlineBlue
	FgUnderlinePurple
	FgUnderlineCyan
	FgUnderlineWhite

	FgHighIntensityBlack
	FgHighIntensityRed
	FgHighIntensityGreen
	FgHighIntensityYellow
	FgHighIntensityBlue
	FgHighIntensityPurple
	FgHighIntensityCyan
	FgHighIntensityWhite

	FgBoldHighIntensityBlack
	FgBoldHighIntensityRed
	FgBoldHighIntensityGreen
	FgBoldHighIntensityYellow
	FgBoldHighIntensityBlue
	FgBoldHighIntensityPurple
	FgBoldHighIntensityCyan
	FgBoldHighIntensityWhite

	BgBlack
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgPurple
	BgCyan
	BgWhite

	BgHighIntensityBlack
	BgHighIntensityRed
	BgHighIntensityGreen
	BgHighIntensityYellow
	BgHighIntensityBlue
	BgHighIntensityPurple
	BgHighIntensityCyan
	BgHighIntensityWhite
)
