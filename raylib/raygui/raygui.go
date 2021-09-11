package raygui

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Style property
type GuiStyleProp struct {
	ControlId     uint
	PropertyId    uint
	PropertyValue int
}

// Gui control state
type ControlState int

const (
	StateNormal ControlState = iota
	StateFocused
	StatePressed
	StateDisabled
)

// Gui control text alignment
type TextAlignment int

const (
	TextAlignLeft TextAlignment = iota
	TextAlignCenter
	TextAlightRight
)

// Gui controls
type Control int

const (
	DefaultControl Control = iota
	LabelControl
	ButtonControl
	ToggleControl
	SliderControl
	ProgressBarControl
	CheckBoxControl
	ComboBoxControl
	DropdownBoxControl
	TextBoxControl
	ValueBoxControl
	SpinnerControl
	ListViewControl
	ColorPickerControl
	ScrollBarControl
	StatusBarControl
)

// Gui base properties for every control
type ControlProperty int

const (
	BorderColorNormalProp ControlProperty = iota
	BaseColorNormalProp
	TextColorNormalProp
	BorderColorFocusedProp
	BaseColorFocusedProp
	TextColorFocusedProp
	BorderColorPressedProp
	BaseColorPressedProp
	TextColorPressedProp
	BorderColorDisabledProp
	BaseColorDisabledProp
	TextColorDisabledProp
	BorderWidthProp
	TextPaddingProp
	TextAlignmentProp
	ReservedProp
)

// Gui extended properties depend on control
// NOTE: We reserve a fixed size of additional properties per control

// DEFAULT properties\
const (
	TextSizeProp ControlProperty = iota + 16
	TextSpacingProp
	LineColorProp
	BackgroundColorProp
)

// Label
//typedef enum { } GuiLabelProperty;

// Button
//typedef enum { } GuiButtonProperty;

// Toggle / ToggleGroup
const (
	GroupPadding ControlProperty = iota + 16
)

// Slider / SliderBar
const (
	SliderWidth ControlProperty = iota + 16
	SliderPadding
)

// ProgressBar
const (
	ProgressPadding ControlProperty = iota + 16
)

// CheckBox
const (
	CheckPadding ControlProperty = iota + 16
)

// ComboBox
const (
	ComboButtonWidth ControlProperty = iota + 16
	ComboButtonPadding
)

// DropdownBox
const (
	ArrowPadding ControlProperty = iota + 16
	DropdownItemsPadding
)

// TextBox / TextBoxMulti / ValueBox / Spinner
const (
	TextInnerPadding ControlProperty = iota + 16
	TextLinesPadding
	ColorSelectedFG
	ColorSelectedBG
)

// Spinner
const (
	SpinButtonWidth ControlProperty = iota + 16
	SpinButtonPadding
)

// ScrollBar
const (
	ArrowsSize ControlProperty = iota + 16
	ArrowsVisible
	ScrollSliderPadding
	ScrollSliderSize
	ScrollPadding
	ScrollSpeed
)

// ScrollBar side
type ScrollBarSide int

const (
	ScrollBarLeftSide ScrollBarSide = iota
	ScrollBarRightSide
)

// ListView
const (
	ListItemsHeight ControlProperty = iota + 16
	ListItemsPadding
	ScrollBarWidth
	ScrollBarSideProp // Renamed from ScrollBarSide due to naming conflict
)

// ColorPicker
const (
	ColorSelectorSize      ControlProperty = iota + 16
	HueBarWidth                            // Right hue bar width
	HueBarPadding                          // Right hue bar separation from panel
	HueBarSelectorHeight                   // Right hue bar selector height
	HueBarSelectorOverflow                 // Right hue bar selector overflow
)

const MaxControls = 16     // Maximum number of standard controls
const MaxPropsDefault = 16 // Maximum number of standard properties
const MaxPropsExtended = 8 // Maximum number of extended properties

//----------------------------------------------------------------------------------
// Types and Structures Definition
//----------------------------------------------------------------------------------

// Gui control property style color element
const (
	Border ControlProperty = iota
	Base
	Text
	Other
)

//----------------------------------------------------------------------------------
// Global Variables Definition
//----------------------------------------------------------------------------------
var guiState = StateNormal

var guiFont rl.Font      // Gui current font (WARNING: highly coupled to raylib)
var guiLocked = false    // Gui lock state (no inputs processed)
var guiAlpha float32 = 1 // Gui element transpacency on drawing

// Global gui style array (allocated on data segment by default)
// NOTE: In raygui we manage a single int array with all the possible style properties.
// When a new style is loaded, it loads over the global style... but default gui style
// could always be recovered with GuiLoadStyleDefault()
var guiStyle [MaxControls * (MaxPropsDefault + MaxPropsExtended)]uint
var guiStyleLoaded = false // Style loaded flag for lazy style initialization

//----------------------------------------------------------------------------------
// Gui Setup Functions Definition
//----------------------------------------------------------------------------------

// Enable gui global state
func Enable() {
	guiState = StateNormal
}

// Disable gui global state
func Disable() {
	guiState = StateDisabled
}

// Lock gui global state
func Lock() {
	guiLocked = true
}

// Unlock gui global state
func Unlock() {
	guiLocked = false
}

// Set gui controls alpha global state
func Fade(alpha float32) {
	if alpha < 0 {
		alpha = 0
	} else if alpha > 1 {
		alpha = 1
	}

	alpha = alpha
}

// Set gui state (global state)
func SetState(state ControlState) {
	guiState = state
}

// Get gui state (global state)
func GetState() ControlState {
	return guiState
}

// Set custom gui font
// NOTE: Font loading/unloading is external to raygui
func SetFont(font rl.Font) {
	if font.Texture.ID > 0 {
		// NOTE: If we try to setup a font but default style has not been
		// lazily loaded before, it will be overwritten, so we need to force
		// default style loading first
		if !guiStyleLoaded {
			LoadStyleDefault()
		}

		guiFont = font
		SetStyle(DefaultControl, TextSizeProp, uint(font.BaseSize))
	}
}

// Get custom gui font
func GetFont() rl.Font {
	return guiFont
}

// Set control style property value
func SetStyle(control Control, property ControlProperty, value uint) {
	if !guiStyleLoaded {
		LoadStyleDefault()
	}
	guiStyle[int(control)*(MaxPropsDefault+MaxPropsExtended)+int(property)] = value

	// Default properties are propagated to all controls
	if (control == 0) && (property < MaxPropsDefault) {
		for i := 1; i < MaxControls; i++ {
			guiStyle[i*(MaxPropsDefault+MaxPropsExtended)+int(property)] = value
		}
	}
}

// Get control style property value
func GetStyle(control Control, property ControlProperty) uint {
	if !guiStyleLoaded {
		LoadStyleDefault()
	}
	return guiStyle[int(control)*(MaxPropsDefault+MaxPropsExtended)+int(property)]
}

//----------------------------------------------------------------------------------
// Gui Controls Functions Definition
//----------------------------------------------------------------------------------

// NOTE: This define is also used by GuiMessageBox() and GuiTextInputBox()
const WindowStatusBarHeight = 22

// Window Box control
func WindowBox(bounds rl.Rectangle, title string) bool {
	//GuiControlState state = guiState;
	clicked := false

	statusBarHeight := WindowStatusBarHeight + 2*GetStyle(StatusBarControl, BorderWidthProp)
	statusBarHeight += statusBarHeight % 2

	statusBar := rl.Rectangle{bounds.X, bounds.Y, bounds.Width, float32(statusBarHeight)}
	if bounds.Height < float32(statusBarHeight)*2 {
		bounds.Height = float32(statusBarHeight) * 2
	}

	windowPanel := rl.Rectangle{bounds.X, bounds.Y + float32(statusBarHeight) - 1, bounds.Width, bounds.Height - float32(statusBarHeight)}
	closeButtonRec := rl.Rectangle{
		statusBar.X + statusBar.Width - float32(GetStyle(StatusBarControl, BorderWidthProp)) - 20,
		statusBar.Y + float32(statusBarHeight)/2 - 18/2,
		18, 18,
	}

	// Update control
	//--------------------------------------------------------------------
	// NOTE: Logic is directly managed by button
	//--------------------------------------------------------------------

	// Draw control
	//--------------------------------------------------------------------
	StatusBar(statusBar, title) // Draw window header as status bar
	Panel(windowPanel)          // Draw window base

	// Draw window close button
	tempBorderWidth := GetStyle(ButtonControl, BorderWidthProp)
	tempTextAlignment := GetStyle(ButtonControl, TextAlignmentProp)
	SetStyle(ButtonControl, BorderWidthProp, 1)
	SetStyle(ButtonControl, TextAlignmentProp, uint(TextAlignCenter))
	clicked = Button(closeButtonRec, "x")
	/*
		// TODO(icons)
		#if defined(RAYGUI_SUPPORT_RICONS)
			clicked = GuiButton(closeButtonRec, GuiIconText(RICON_CROSS_SMALL, NULL));
		#else
			clicked = GuiButton(closeButtonRec, "x");
		#endif
	*/
	SetStyle(ButtonControl, BorderWidthProp, tempBorderWidth)
	SetStyle(ButtonControl, TextAlignmentProp, tempTextAlignment)

	return clicked
}

const GroupBoxLineThick = 1
const GroupBoxTextPadding = 10

// Group Box control with text name
func GroupBox(bounds rl.Rectangle, text string) {
	state := guiState

	borderColorProp := LineColorProp
	if state == StateDisabled {
		borderColorProp = BorderColorDisabledProp
	}
	borderColor := rl.GetColor(int32(GetStyle(DefaultControl, borderColorProp)))

	// Draw control
	//--------------------------------------------------------------------
	DrawRectangle(rl.Rectangle{bounds.X, bounds.Y, GroupBoxLineThick, bounds.Height}, 0, rl.Blank, rl.Fade(borderColor, guiAlpha))
	DrawRectangle(rl.Rectangle{bounds.X, bounds.Y + bounds.Height - 1, bounds.Width, GroupBoxLineThick}, 0, rl.Blank, rl.Fade(borderColor, guiAlpha))
	DrawRectangle(rl.Rectangle{bounds.X + bounds.Width - 1, bounds.Y, GroupBoxLineThick, bounds.Height}, 0, rl.Blank, rl.Fade(borderColor, guiAlpha))

	Line(rl.Rectangle{bounds.X, bounds.Y, bounds.Width, 1}, text)
	//--------------------------------------------------------------------
}

const LineTextPadding = 10

// Line control
func Line(bounds rl.Rectangle, text string) {
	state := guiState

	colorProp := LineColorProp
	if state == StateDisabled {
		colorProp = BorderColorDisabledProp
	}
	colorStyle := rl.GetColor(int32(GetStyle(DefaultControl, colorProp)))

	color := rl.Fade(colorStyle, guiAlpha)

	// Draw control
	//--------------------------------------------------------------------
	textBounds := rl.Rectangle{
		Width:  GetTextWidth(text), // TODO: Consider text icon
		Height: float32(GetStyle(DefaultControl, TextSizeProp)),
		X:      bounds.X + LineTextPadding,
		Y:      bounds.Y - float32(GetStyle(DefaultControl, TextSizeProp))/2,
	}

	// Draw line with embedded text label: "--- text --------------"
	DrawRectangle(rl.Rectangle{bounds.X, bounds.Y, LineTextPadding - 2, 1}, 0, rl.Blank, color)
	Label(textBounds, text)
	DrawRectangle(rl.Rectangle{bounds.X + LineTextPadding + textBounds.Width + 4, bounds.Y, bounds.Width - textBounds.Width - LineTextPadding - 4, 1}, 0, rl.Blank, color)
	//--------------------------------------------------------------------
}

const PanelBorderWidth = 1

// Panel control
func Panel(bounds rl.Rectangle) {
	state := guiState

	borderColorProp := LineColorProp
	if state == StateDisabled {
		borderColorProp = BorderColorDisabledProp
	}
	borderColor := rl.Fade(rl.GetColor(int32(GetStyle(DefaultControl, borderColorProp))), guiAlpha)

	colorProp := BackgroundColorProp
	if state == StateDisabled {
		colorProp = BaseColorDisabledProp
	}
	color := rl.Fade(rl.GetColor(int32(GetStyle(DefaultControl, colorProp))), guiAlpha)

	// Draw control
	//--------------------------------------------------------------------
	DrawRectangle(bounds, PanelBorderWidth, borderColor, color)
	//--------------------------------------------------------------------
}

// Scroll Panel control
func ScrollPanel(bounds, content rl.Rectangle, scroll *rl.Vector2) rl.Rectangle {
	state := guiState

	bw := float32(GetStyle(DefaultControl, BorderWidthProp))
	side := ScrollBarSide(GetStyle(ListViewControl, ScrollBarSideProp))

	scrollPos := rl.Vector2{0, 0}
	if scroll != nil {
		scrollPos = *scroll
	}

	hasHorizontalScrollBar := content.Width > bounds.Width-2*bw
	hasVerticalScrollBar := content.Height > bounds.Height-2*bw

	// Recheck to account for the other scrollbar being visible
	if !hasHorizontalScrollBar {
		hasHorizontalScrollBar = hasVerticalScrollBar && (content.Width > bounds.Width-2*bw-float32(GetStyle(ListViewControl, ScrollBarWidth)))
	}
	if !hasVerticalScrollBar {
		hasVerticalScrollBar = hasHorizontalScrollBar && (content.Height > bounds.Height-2*bw-float32(GetStyle(ListViewControl, ScrollBarWidth)))
	}

	var horizontalScrollBarWidth int = 0
	if hasHorizontalScrollBar {
		horizontalScrollBarWidth = int(GetStyle(ListViewControl, ScrollBarWidth))
	}
	var verticalScrollBarWidth int = 0
	if hasVerticalScrollBar {
		verticalScrollBarWidth = int(GetStyle(ListViewControl, ScrollBarWidth))
	}

	hx := bounds.X
	if side == ScrollBarLeftSide {
		hx = bounds.X + float32(verticalScrollBarWidth)
	}
	horizontalScrollBar := rl.Rectangle{
		X:      hx + bw,
		Y:      bounds.Y + bounds.Height - float32(horizontalScrollBarWidth) - bw,
		Width:  bounds.Width - float32(verticalScrollBarWidth) - 2*bw,
		Height: float32(horizontalScrollBarWidth),
	}

	vx := bounds.X + bounds.Width - float32(verticalScrollBarWidth) - bw
	if side == ScrollBarLeftSide {
		vx = bounds.X + bw
	}
	verticalScrollBar := rl.Rectangle{
		X:      vx,
		Y:      bounds.Y + bw,
		Width:  float32(verticalScrollBarWidth),
		Height: bounds.Height - float32(horizontalScrollBarWidth) - 2*bw,
	}

	// Calculate view area (area without the scrollbars)
	view := rl.Rectangle{bounds.X + bw, bounds.Y + bw, bounds.Width - 2*bw - float32(verticalScrollBarWidth), bounds.Height - 2*bw - float32(horizontalScrollBarWidth)}
	if side == ScrollBarLeftSide {
		view = rl.Rectangle{bounds.X + float32(verticalScrollBarWidth) + bw, bounds.Y + bw, bounds.Width - 2*bw - float32(verticalScrollBarWidth), bounds.Height - 2*bw - float32(horizontalScrollBarWidth)}
	}

	// Clip view area to the actual content size
	if view.Width > content.Width {
		view.Width = content.Width
	}
	if view.Height > content.Height {
		view.Height = content.Height
	}

	// TODO: Review!
	var horizontalMin float32
	var horizontalMax float32
	if hasHorizontalScrollBar {
		var minOffset float32
		var maxOffset float32
		if side == ScrollBarLeftSide {
			minOffset = float32(-verticalScrollBarWidth)
			maxOffset = float32(verticalScrollBarWidth)
		}
		horizontalMin = minOffset - bw
		horizontalMax = content.Width - bounds.Width + float32(verticalScrollBarWidth) + bw - maxOffset
	} else {
		var minOffset float32
		if side == ScrollBarLeftSide {
			minOffset = float32(-verticalScrollBarWidth)
		}
		horizontalMin = minOffset - bw
		horizontalMax = -bw
	}

	var verticalMin float32 = -bw
	var verticalMax float32
	if hasVerticalScrollBar {
		verticalMin = -bw
		verticalMax = content.Height - bounds.Height + float32(horizontalScrollBarWidth) + bw
	} else {
		verticalMin = -bw
		verticalMax = -bw
	}

	// Update control
	//--------------------------------------------------------------------
	if state != StateDisabled && !guiLocked {
		mousePoint := rl.GetMousePosition()

		// Check button state
		if rl.CheckCollisionPointRec(mousePoint, bounds) {
			if rl.IsMouseButtonDown(rl.MouseLeftButton) {
				state = StatePressed
			} else {
				state = StateFocused
			}

			if hasHorizontalScrollBar {
				if rl.IsKeyDown(rl.KeyRight) {
					scrollPos.X -= float32(GetStyle(ScrollBarControl, ScrollSpeed))
				}
				if rl.IsKeyDown(rl.KeyLeft) {
					scrollPos.X += float32(GetStyle(ScrollBarControl, ScrollSpeed))
				}
			}

			if hasVerticalScrollBar {
				if rl.IsKeyDown(rl.KeyDown) {
					scrollPos.Y -= float32(GetStyle(ScrollBarControl, ScrollSpeed))
				}
				if rl.IsKeyDown(rl.KeyUp) {
					scrollPos.Y += float32(GetStyle(ScrollBarControl, ScrollSpeed))
				}
			}

			wheelMove := rl.GetMouseWheelMove()

			// Horizontal scroll (Shift + Mouse wheel)
			if hasHorizontalScrollBar && (rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift)) {
				scrollPos.X += float32(wheelMove) * 20
			} else {
				// Vertical scroll
				scrollPos.Y += float32(wheelMove) * 20
			}
		}
	}

	// Normalize scroll values
	if scrollPos.X > -horizontalMin {
		scrollPos.X = -horizontalMin
	}
	if scrollPos.X < -horizontalMax {
		scrollPos.X = -horizontalMax
	}
	if scrollPos.Y > -verticalMin {
		scrollPos.Y = -verticalMin
	}
	if scrollPos.Y < -verticalMax {
		scrollPos.Y = -verticalMax
	}
	//--------------------------------------------------------------------

	// Draw control
	//--------------------------------------------------------------------
	DrawRectangle(bounds, 0, rl.Blank, rl.GetColor(int32(GetStyle(DefaultControl, BackgroundColorProp)))) // Draw background

	// Save size of the scrollbar slider
	slider := GetStyle(ScrollBarControl, ScrollSliderSize)

	// Draw horizontal scrollbar if visible
	if hasHorizontalScrollBar {
		// Change scrollbar slider size to show the diff in size between the content width and the widget width
		SetStyle(ScrollBarControl, ScrollSliderSize, uint(((bounds.Width-2*bw-float32(verticalScrollBarWidth))/floor32(content.Width))*(floor32(bounds.Width)-2*bw-float32(verticalScrollBarWidth))))
		scrollPos.X = float32(-ScrollBar(horizontalScrollBar, int(-scrollPos.X), int(horizontalMin), int(horizontalMax)))
	}

	// Draw vertical scrollbar if visible
	if hasVerticalScrollBar {
		// Change scrollbar slider size to show the diff in size between the content height and the widget height
		SetStyle(ScrollBarControl, ScrollSliderSize, uint(((bounds.Height-2*bw-float32(horizontalScrollBarWidth))/floor32(content.Height))*(floor32(bounds.Height)-2*bw-float32(horizontalScrollBarWidth))))
		scrollPos.Y = float32(-ScrollBar(verticalScrollBar, int(-scrollPos.Y), int(verticalMin), int(verticalMax)))
	}

	// Draw detail corner rectangle if both scroll bars are visible
	if hasHorizontalScrollBar && hasVerticalScrollBar {
		var x float32
		if side == ScrollBarLeftSide {
			x = bounds.X + bw + 2
		} else {
			x = horizontalScrollBar.X + horizontalScrollBar.Width + 2
		}
		corner := rl.Rectangle{x, verticalScrollBar.Y + verticalScrollBar.Height + 2, float32(horizontalScrollBarWidth) - 4, float32(verticalScrollBarWidth) - 4}
		DrawRectangle(corner, 0, rl.Blank, rl.Fade(rl.GetColor(int32(GetStyle(ListViewControl, Text+(ControlProperty(state)*3)))), guiAlpha))
	}

	// Draw scrollbar lines depending on current state
	DrawRectangle(bounds, int(GetStyle(DefaultControl, BorderWidthProp)), rl.Fade(rl.GetColor(int32(GetStyle(ListViewControl, Border+(ControlProperty(state)*3)))), guiAlpha), rl.Blank)

	// Set scrollbar slider size back to the way it was before
	SetStyle(ScrollBarControl, ScrollSliderSize, slider)
	//--------------------------------------------------------------------

	if scroll != nil {
		*scroll = scrollPos
	}

	return view
}

/*

Rectangle GuiScrollPanel(Rectangle bounds, Rectangle content, Vector2 *scroll)
{


    // Draw scrollbar lines depending on current state
    GuiDrawRectangle(bounds, GuiGetStyle(DEFAULT, BORDER_WIDTH), Fade(GetColor(GuiGetStyle(LISTVIEW, BORDER + (state*3))), guiAlpha), BLANK);

    // Set scrollbar slider size back to the way it was before
    GuiSetStyle(SCROLLBAR, SCROLL_SLIDER_SIZE, slider);
    //--------------------------------------------------------------------

    if (scroll != NULL) *scroll = scrollPos;

    return view;
}
*/

// Gui draw rectangle using default raygui plain style with borders
func DrawRectangle(rec rl.Rectangle, borderWidth int, borderColor, color rl.Color) {
	if color.A > 0 {
		// Draw rectangle filled with color
		rl.DrawRectangle(int32(rec.X), int32(rec.Y), int32(rec.Width), int32(rec.Height), color)
	}

	if borderWidth > 0 {
		// Draw rectangle border lines with color
		rl.DrawRectangle(int32(rec.X), int32(rec.Y), int32(rec.Width), int32(borderWidth), borderColor)
		rl.DrawRectangle(int32(rec.X), int32(rec.Y)+int32(borderWidth), int32(borderWidth), int32(rec.Height)-2*int32(borderWidth), borderColor)
		rl.DrawRectangle(int32(rec.X)+int32(rec.Width)-int32(borderWidth), int32(rec.Y)+int32(borderWidth), int32(borderWidth), int32(rec.Height)-2*int32(borderWidth), borderColor)
		rl.DrawRectangle(int32(rec.X), int32(rec.Y)+int32(rec.Height)-int32(borderWidth), int32(rec.Width), int32(borderWidth), borderColor)
	}

	// TODO: For n-patch-based style we would need: [state] and maybe [control]
	// In this case all controls drawing logic should be moved to this function... I don't like it...
}

func floor32(f float32) float32 {
	return float32(math.Floor(float64(f)))
}
