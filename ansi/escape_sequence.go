package ansi

type ANSIEscapeSequenceByte byte

const (
	ESCAPEByte ANSIEscapeSequenceByte = 0x1b
)

type ANSIEscapeSequence string

const (
	ESCAPE                        ANSIEscapeSequence = "\x1b"
	SS3                           ANSIEscapeSequence = ESCAPE + "O"
	CSI                           ANSIEscapeSequence = ESCAPE + "["
	RESET                         ANSIEscapeSequence = CSI + "0m"
	CLEAR_SCREEN                  ANSIEscapeSequence = CSI + "2J"
	CLEAR_SCROLLBACK              ANSIEscapeSequence = CSI + "3J"
	CURSOR_HOME                   ANSIEscapeSequence = CSI + "H"
	HIDE_CURSOR                   ANSIEscapeSequence = CSI + "?25l"
	SHOW_CURSOR                   ANSIEscapeSequence = CSI + "?25h"
	ENTER_ALTERNATE_SCREEN        ANSIEscapeSequence = CSI + "?1049h"
	EXIT_ALTERNATE_SCREEN         ANSIEscapeSequence = CSI + "?1049l"
	FG_BLACK                      ANSIEscapeSequence = CSI + "0;30m"
	FG_RED                        ANSIEscapeSequence = CSI + "0;31m"
	FG_GREEN                      ANSIEscapeSequence = CSI + "0;32m"
	FG_YELLOW                     ANSIEscapeSequence = CSI + "0;33m"
	FG_BLUE                       ANSIEscapeSequence = CSI + "0;34m"
	FG_PURPLE                     ANSIEscapeSequence = CSI + "0;35m"
	FG_CYAN                       ANSIEscapeSequence = CSI + "0;36m"
	FG_WHITE                      ANSIEscapeSequence = CSI + "0;37m"
	FG_BOLD_BLACK                 ANSIEscapeSequence = CSI + "1;30m"
	FG_BOLD_RED                   ANSIEscapeSequence = CSI + "1;31m"
	FG_BOLD_GREEN                 ANSIEscapeSequence = CSI + "1;32m"
	FG_BOLD_YELLOW                ANSIEscapeSequence = CSI + "1;33m"
	FG_BOLD_BLUE                  ANSIEscapeSequence = CSI + "1;34m"
	FG_BOLD_PURPLE                ANSIEscapeSequence = CSI + "1;35m"
	FG_BOLD_CYAN                  ANSIEscapeSequence = CSI + "1;36m"
	FG_BOLD_WHITE                 ANSIEscapeSequence = CSI + "1;37m"
	FG_UNDERLINE_BLACK            ANSIEscapeSequence = CSI + "4;30m"
	FG_UNDERLINE_RED              ANSIEscapeSequence = CSI + "4;31m"
	FG_UNDERLINE_GREEN            ANSIEscapeSequence = CSI + "4;32m"
	FG_UNDERLINE_YELLOW           ANSIEscapeSequence = CSI + "4;33m"
	FG_UNDERLINE_BLUE             ANSIEscapeSequence = CSI + "4;34m"
	FG_UNDERLINE_PURPLE           ANSIEscapeSequence = CSI + "4;35m"
	FG_UNDERLINE_CYAN             ANSIEscapeSequence = CSI + "4;36m"
	FG_UNDERLINE_WHITE            ANSIEscapeSequence = CSI + "4;37m"
	BG_BLACK                      ANSIEscapeSequence = CSI + "40m"
	BG_RED                        ANSIEscapeSequence = CSI + "41m"
	BG_GREEN                      ANSIEscapeSequence = CSI + "42m"
	BG_YELLOW                     ANSIEscapeSequence = CSI + "43m"
	BG_BLUE                       ANSIEscapeSequence = CSI + "44m"
	BG_PURPLE                     ANSIEscapeSequence = CSI + "45m"
	BG_CYAN                       ANSIEscapeSequence = CSI + "46m"
	BG_WHITE                      ANSIEscapeSequence = CSI + "47m"
	FG_HIGH_INTENSITY_BLACK       ANSIEscapeSequence = CSI + "0;90m"
	FG_HIGH_INTENSITY_RED         ANSIEscapeSequence = CSI + "0;91m"
	FG_HIGH_INTENSITY_GREEN       ANSIEscapeSequence = CSI + "0;92m"
	FG_HIGH_INTENSITY_YELLOW      ANSIEscapeSequence = CSI + "0;93m"
	FG_HIGH_INTENSITY_BLUE        ANSIEscapeSequence = CSI + "0;94m"
	FG_HIGH_INTENSITY_PURPLE      ANSIEscapeSequence = CSI + "0;95m"
	FG_HIGH_INTENSITY_CYAN        ANSIEscapeSequence = CSI + "0;96m"
	FG_HIGH_INTENSITY_WHITE       ANSIEscapeSequence = CSI + "0;97m"
	FG_BOLD_HIGH_INTENSITY_BLACK  ANSIEscapeSequence = CSI + "1;90m"
	FG_BOLD_HIGH_INTENSITY_RED    ANSIEscapeSequence = CSI + "1;91m"
	FG_BOLD_HIGH_INTENSITY_GREEN  ANSIEscapeSequence = CSI + "1;92m"
	FG_BOLD_HIGH_INTENSITY_YELLOW ANSIEscapeSequence = CSI + "1;93m"
	FG_BOLD_HIGH_INTENSITY_BLUE   ANSIEscapeSequence = CSI + "1;94m"
	FG_BOLD_HIGH_INTENSITY_PURPLE ANSIEscapeSequence = CSI + "1;95m"
	FG_BOLD_HIGH_INTENSITY_CYAN   ANSIEscapeSequence = CSI + "1;96m"
	FG_BOLD_HIGH_INTENSITY_WHITE  ANSIEscapeSequence = CSI + "1;97m"
	BG_HIGH_INTENSITY_BLACK       ANSIEscapeSequence = CSI + "0;100m"
	BG_HIGH_INTENSITY_RED         ANSIEscapeSequence = CSI + "0;101m"
	BG_HIGH_INTENSITY_GREEN       ANSIEscapeSequence = CSI + "0;102m"
	BG_HIGH_INTENSITY_YELLOW      ANSIEscapeSequence = CSI + "0;103m"
	BG_HIGH_INTENSITY_BLUE        ANSIEscapeSequence = CSI + "0;104m"
	BG_HIGH_INTENSITY_PURPLE      ANSIEscapeSequence = CSI + "0;105m"
	BG_HIGH_INTENSITY_CYAN        ANSIEscapeSequence = CSI + "0;106m"
	BG_HIGH_INTENSITY_WHITE       ANSIEscapeSequence = CSI + "0;107m"
)

func (c ANSIEscapeSequence) IsColor() bool {
	return c.IsValid()
}

func (c ANSIEscapeSequence) IsValid() bool {
	switch c {
	case FG_BLACK, FG_RED, FG_GREEN, FG_YELLOW, FG_BLUE, FG_PURPLE, FG_CYAN,
		FG_WHITE, FG_BOLD_BLACK, FG_BOLD_RED, FG_BOLD_GREEN, FG_BOLD_YELLOW,
		FG_BOLD_BLUE, FG_BOLD_PURPLE, FG_BOLD_CYAN, FG_BOLD_WHITE, FG_UNDERLINE_BLACK,
		FG_UNDERLINE_RED, FG_UNDERLINE_GREEN, FG_UNDERLINE_YELLOW, FG_UNDERLINE_BLUE,
		FG_UNDERLINE_PURPLE, FG_UNDERLINE_CYAN, FG_UNDERLINE_WHITE, BG_BLACK,
		BG_RED, BG_GREEN, BG_YELLOW, BG_BLUE, BG_PURPLE, BG_CYAN, BG_WHITE,
		FG_HIGH_INTENSITY_BLACK, FG_HIGH_INTENSITY_RED, FG_HIGH_INTENSITY_GREEN,
		FG_HIGH_INTENSITY_YELLOW, FG_HIGH_INTENSITY_BLUE, FG_HIGH_INTENSITY_PURPLE,
		FG_HIGH_INTENSITY_CYAN, FG_HIGH_INTENSITY_WHITE, FG_BOLD_HIGH_INTENSITY_BLACK,
		FG_BOLD_HIGH_INTENSITY_RED, FG_BOLD_HIGH_INTENSITY_GREEN, FG_BOLD_HIGH_INTENSITY_YELLOW,
		FG_BOLD_HIGH_INTENSITY_BLUE, FG_BOLD_HIGH_INTENSITY_PURPLE, FG_BOLD_HIGH_INTENSITY_CYAN,
		FG_BOLD_HIGH_INTENSITY_WHITE, BG_HIGH_INTENSITY_BLACK, BG_HIGH_INTENSITY_RED,
		BG_HIGH_INTENSITY_GREEN, BG_HIGH_INTENSITY_YELLOW, BG_HIGH_INTENSITY_BLUE,
		BG_HIGH_INTENSITY_PURPLE, BG_HIGH_INTENSITY_CYAN, BG_HIGH_INTENSITY_WHITE:
		return true
	default:
		return false
	}
}

func (c ANSIEscapeSequence) IsCommand() bool {
	switch c {
	case RESET, CLEAR_SCREEN, CLEAR_SCROLLBACK, CURSOR_HOME, HIDE_CURSOR,
		SHOW_CURSOR, ENTER_ALTERNATE_SCREEN, EXIT_ALTERNATE_SCREEN:
		return true
	default:
		return false
	}
}
