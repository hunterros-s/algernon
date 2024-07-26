package text

import "fmt"

// Content types
const (
	TextType      = "text"
	TranslateType = "translate"
	KeybindType   = "keybind"
	ScoreType     = "score"
	SelectorType  = "selector"
	NBTType       = "nbt"
)

// Font types
const (
	DefaultFont    = "minecraft:default"
	UniformFont    = "minecraft:uniform"
	AltFont        = "minecraft:alt"
	IllageraltFont = "minecraft:illageralt"
)

// Click event actions
const (
	OpenURLAction         = "open_url"
	RunCommandAction      = "run_command"
	SuggestCommandAction  = "suggest_command"
	ChangePageAction      = "change_page"
	CopyToClipboardAction = "copy_to_clipboard"
)

// Hover event actions
const (
	ShowTextAction   = "show_text"
	ShowItemAction   = "show_item"
	ShowEntityAction = "show_entity"
)

type ClickEventData struct {
	Action string `json:"action" nbt:"action"`
	Value  string `json:"value" nbt:"value"`
}

type HoverEventData struct {
	Action   string        `json:"action" nbt:"action"`
	Contents HoverContents `json:"contents" nbt:"contents"`
}

type HoverContents struct {
	ID    string `json:"id,omitempty" nbt:"id,omitempty"`
	Count int    `json:"count,omitempty" nbt:"count,omitempty"`
	Tag   string `json:"tag,omitempty" nbt:"tag,omitempty"`
	Type  string `json:"type,omitempty" nbt:"type,omitempty"`
	Name  string `json:"name,omitempty" nbt:"name,omitempty"`
}

type TextComponent struct {
	Type     string          `json:"type,omitempty" nbt:"type,omitempty"`
	Children []TextComponent `json:"extra,omitempty" nbt:"extra,omitempty"`
	Text     string          `json:"text" nbt:"text"`

	// Styling
	Color         string         `json:"color,omitempty" nbt:"color,omitempty"`
	Bold          bool           `json:"bold,omitempty" nbt:"bold,omitempty"`
	Italic        bool           `json:"italic,omitempty" nbt:"italic,omitempty"`
	Underlined    bool           `json:"underlined,omitempty" nbt:"underlined,omitempty"`
	Strikethrough bool           `json:"strikethrough,omitempty" nbt:"strikethrough,omitempty"`
	Obfuscated    bool           `json:"obfuscated,omitempty" nbt:"obfuscated,omitempty"`
	Font          string         `json:"font,omitempty" nbt:"font,omitempty"`
	Insertion     string         `json:"insertion,omitempty" nbt:"insertion,omitempty"`
	ClickEvent    ClickEventData `json:"click_event,omitempty" nbt:"click_event,omitempty"`
	HoverEvent    HoverEventData `json:"hover_event,omitempty" nbt:"hover_event,omitempty"`
}

func ParseFormatted(format string, args ...interface{}) TextComponent {
	return Parse(fmt.Sprintf(format, args...), '&')
}

func Parse(text string, formatChar rune) TextComponent {
	root := TextComponent{}
	current := TextComponent{}
	runes := []rune(text)

	for i := 0; i < len(runes); i++ {
		if runes[i] == formatChar && i+1 < len(runes) && isFormatCode(runes[i+1]) {
			if current.Text != "" {
				root.Children = append(root.Children, current)
				current = TextComponent{}
			}
			applyFormatCode(runes[i+1], &current)
			i++
		} else {
			current.Text += string(runes[i])
		}
	}

	root.Children = append(root.Children, current)
	return root
}

func Serialize(comp TextComponent, formatChar rune) string {
	components := append([]TextComponent{comp}, comp.Children...)
	result := ""

	for _, c := range components {
		result += serializeComponent(c, formatChar)
		result += c.Text
	}

	return result
}

func isFormatCode(r rune) bool {
	return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'k' && r <= 'o') || r == 'r'
}

func applyFormatCode(code rune, comp *TextComponent) {
	if code >= '0' && code <= 'f' {
		comp.Color = codeToColor(code)
	}
	switch code {
	case 'k':
		comp.Obfuscated = true
	case 'l':
		comp.Bold = true
	case 'm':
		comp.Strikethrough = true
	case 'n':
		comp.Underlined = true
	case 'o':
		comp.Italic = true
	}
}

func serializeComponent(comp TextComponent, formatChar rune) string {
	result := colorToCode(comp.Color, formatChar)
	if comp.Obfuscated {
		result += string(formatChar) + "k"
	}
	if comp.Bold {
		result += string(formatChar) + "l"
	}
	if comp.Strikethrough {
		result += string(formatChar) + "m"
	}
	if comp.Underlined {
		result += string(formatChar) + "n"
	}
	if comp.Italic {
		result += string(formatChar) + "o"
	}
	return result
}

func codeToColor(code rune) string {
	colors := map[rune]string{
		'0': "black", '1': "dark_blue", '2': "dark_green", '3': "dark_aqua",
		'4': "dark_red", '5': "dark_purple", '6': "gold", '7': "gray",
		'8': "dark_gray", '9': "blue", 'a': "green", 'b': "aqua",
		'c': "red", 'd': "light_purple", 'e': "yellow", 'f': "white",
	}
	return colors[code]
}

func colorToCode(color string, formatChar rune) string {
	codes := map[string]string{
		"black": "0", "dark_blue": "1", "dark_green": "2", "dark_aqua": "3",
		"dark_red": "4", "dark_purple": "5", "gold": "6", "gray": "7",
		"dark_gray": "8", "blue": "9", "green": "a", "aqua": "b",
		"red": "c", "light_purple": "d", "yellow": "e", "white": "f",
	}
	if code, ok := codes[color]; ok {
		return string(formatChar) + code
	}
	return ""
}
