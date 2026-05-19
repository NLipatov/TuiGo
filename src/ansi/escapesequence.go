package ansi

type ANSIEscapeSequence string

const (
	RESET                         ANSIEscapeSequence = "\x1b[0m"
	FG_BLACK                      ANSIEscapeSequence = "\x1b[0;30m"
	FG_RED                        ANSIEscapeSequence = "\x1b[0;31m"
	FG_GREEN                      ANSIEscapeSequence = "\x1b[0;32m"
	FG_YELLOW                     ANSIEscapeSequence = "\x1b[0;33m"
	FG_BLUE                       ANSIEscapeSequence = "\x1b[0;34m"
	FG_PURPLE                     ANSIEscapeSequence = "\x1b[0;35m"
	FG_CYAN                       ANSIEscapeSequence = "\x1b[0;36m"
	FG_WHITE                      ANSIEscapeSequence = "\x1b[0;37m"
	FG_BOLD_BLACK                 ANSIEscapeSequence = "\x1b[1;30m"
	FG_BOLD_RED                   ANSIEscapeSequence = "\x1b[1;31m"
	FG_BOLD_GREEN                 ANSIEscapeSequence = "\x1b[1;32m"
	FG_BOLD_YELLOW                ANSIEscapeSequence = "\x1b[1;33m"
	FG_BOLD_BLUE                  ANSIEscapeSequence = "\x1b[1;34m"
	FG_BOLD_PURPLE                ANSIEscapeSequence = "\x1b[1;35m"
	FG_BOLD_CYAN                  ANSIEscapeSequence = "\x1b[1;36m"
	FG_BOLD_WHITE                 ANSIEscapeSequence = "\x1b[1;37m"
	FG_UNDERLINE_BLACK            ANSIEscapeSequence = "\x1b[4;30m"
	FG_UNDERLINE_RED              ANSIEscapeSequence = "\x1b[4;31m"
	FG_UNDERLINE_GREEN            ANSIEscapeSequence = "\x1b[4;32m"
	FG_UNDERLINE_YELLOW           ANSIEscapeSequence = "\x1b[4;33m"
	FG_UNDERLINE_BLUE             ANSIEscapeSequence = "\x1b[4;34m"
	FG_UNDERLINE_PURPLE           ANSIEscapeSequence = "\x1b[4;35m"
	FG_UNDERLINE_CYAN             ANSIEscapeSequence = "\x1b[4;36m"
	FG_UNDERLINE_WHITE            ANSIEscapeSequence = "\x1b[4;37m"
	BG_BLACK                      ANSIEscapeSequence = "\x1b[40m"
	BG_RED                        ANSIEscapeSequence = "\x1b[41m"
	BG_GREEN                      ANSIEscapeSequence = "\x1b[42m"
	BG_YELLOW                     ANSIEscapeSequence = "\x1b[43m"
	BG_BLUE                       ANSIEscapeSequence = "\x1b[44m"
	BG_PURPLE                     ANSIEscapeSequence = "\x1b[45m"
	BG_CYAN                       ANSIEscapeSequence = "\x1b[46m"
	BG_WHITE                      ANSIEscapeSequence = "\x1b[47m"
	FG_HIGH_INTENSITY_BLACK       ANSIEscapeSequence = "\x1b[0;90m"
	FG_HIGH_INTENSITY_RED         ANSIEscapeSequence = "\x1b[0;91m"
	FG_HIGH_INTENSITY_GREEN       ANSIEscapeSequence = "\x1b[0;92m"
	FG_HIGH_INTENSITY_YELLOW      ANSIEscapeSequence = "\x1b[0;93m"
	FG_HIGH_INTENSITY_BLUE        ANSIEscapeSequence = "\x1b[0;94m"
	FG_HIGH_INTENSITY_PURPLE      ANSIEscapeSequence = "\x1b[0;95m"
	FG_HIGH_INTENSITY_CYAN        ANSIEscapeSequence = "\x1b[0;96m"
	FG_HIGH_INTENSITY_WHITE       ANSIEscapeSequence = "\x1b[0;97m"
	FG_BOLD_HIGH_INTENSITY_BLACK  ANSIEscapeSequence = "\x1b[1;90m"
	FG_BOLD_HIGH_INTENSITY_RED    ANSIEscapeSequence = "\x1b[1;91m"
	FG_BOLD_HIGH_INTENSITY_GREEN  ANSIEscapeSequence = "\x1b[1;92m"
	FG_BOLD_HIGH_INTENSITY_YELLOW ANSIEscapeSequence = "\x1b[1;93m"
	FG_BOLD_HIGH_INTENSITY_BLUE   ANSIEscapeSequence = "\x1b[1;94m"
	FG_BOLD_HIGH_INTENSITY_PURPLE ANSIEscapeSequence = "\x1b[1;95m"
	FG_BOLD_HIGH_INTENSITY_CYAN   ANSIEscapeSequence = "\x1b[1;96m"
	FG_BOLD_HIGH_INTENSITY_WHITE  ANSIEscapeSequence = "\x1b[1;97m"
	BG_HIGH_INTENSITY_BLACK       ANSIEscapeSequence = "\x1b[0;100m"
	BG_HIGH_INTENSITY_RED         ANSIEscapeSequence = "\x1b[0;101m"
	BG_HIGH_INTENSITY_GREEN       ANSIEscapeSequence = "\x1b[0;102m"
	BG_HIGH_INTENSITY_YELLOW      ANSIEscapeSequence = "\x1b[0;103m"
	BG_HIGH_INTENSITY_BLUE        ANSIEscapeSequence = "\x1b[0;104m"
	BG_HIGH_INTENSITY_PURPLE      ANSIEscapeSequence = "\x1b[0;105m"
	BG_HIGH_INTENSITY_CYAN        ANSIEscapeSequence = "\x1b[0;106m"
	BG_HIGH_INTENSITY_WHITE       ANSIEscapeSequence = "\x1b[0;107m"
)

func (c ANSIEscapeSequence) IsColor() bool {
	return c.IsValid()
}

func (c ANSIEscapeSequence) IsValid() bool {
	switch c {
	case FG_BLACK, FG_RED, FG_GREEN, FG_YELLOW, FG_BLUE, FG_PURPLE, FG_CYAN, FG_WHITE,
		FG_BOLD_BLACK, FG_BOLD_RED, FG_BOLD_GREEN, FG_BOLD_YELLOW, FG_BOLD_BLUE, FG_BOLD_PURPLE, FG_BOLD_CYAN, FG_BOLD_WHITE,
		FG_UNDERLINE_BLACK, FG_UNDERLINE_RED, FG_UNDERLINE_GREEN, FG_UNDERLINE_YELLOW, FG_UNDERLINE_BLUE, FG_UNDERLINE_PURPLE, FG_UNDERLINE_CYAN, FG_UNDERLINE_WHITE,
		BG_BLACK, BG_RED, BG_GREEN, BG_YELLOW, BG_BLUE, BG_PURPLE, BG_CYAN, BG_WHITE,
		FG_HIGH_INTENSITY_BLACK, FG_HIGH_INTENSITY_RED, FG_HIGH_INTENSITY_GREEN, FG_HIGH_INTENSITY_YELLOW, FG_HIGH_INTENSITY_BLUE, FG_HIGH_INTENSITY_PURPLE, FG_HIGH_INTENSITY_CYAN, FG_HIGH_INTENSITY_WHITE,
		FG_BOLD_HIGH_INTENSITY_BLACK, FG_BOLD_HIGH_INTENSITY_RED, FG_BOLD_HIGH_INTENSITY_GREEN, FG_BOLD_HIGH_INTENSITY_YELLOW, FG_BOLD_HIGH_INTENSITY_BLUE, FG_BOLD_HIGH_INTENSITY_PURPLE, FG_BOLD_HIGH_INTENSITY_CYAN, FG_BOLD_HIGH_INTENSITY_WHITE,
		BG_HIGH_INTENSITY_BLACK, BG_HIGH_INTENSITY_RED, BG_HIGH_INTENSITY_GREEN, BG_HIGH_INTENSITY_YELLOW, BG_HIGH_INTENSITY_BLUE, BG_HIGH_INTENSITY_PURPLE, BG_HIGH_INTENSITY_CYAN, BG_HIGH_INTENSITY_WHITE:
		return true
	default:
		return false
	}
}
