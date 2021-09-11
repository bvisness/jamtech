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
	Default Control = iota
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
		SetStyle(Default, TextSizeProp, uint(font.BaseSize))
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
	borderColor := rl.GetColor(int32(GetStyle(Default, borderColorProp)))

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
	colorStyle := rl.GetColor(int32(GetStyle(Default, colorProp)))

	color := rl.Fade(colorStyle, guiAlpha)

	// Draw control
	//--------------------------------------------------------------------
	textBounds := rl.Rectangle{
		Width:  GetTextWidth(text), // TODO: Consider text icon
		Height: float32(GetStyle(Default, TextSizeProp)),
		X:      bounds.X + LineTextPadding,
		Y:      bounds.Y - float32(GetStyle(Default, TextSizeProp))/2,
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
	borderColor := rl.Fade(rl.GetColor(int32(GetStyle(Default, borderColorProp))), guiAlpha)

	colorProp := BackgroundColorProp
	if state == StateDisabled {
		colorProp = BaseColorDisabledProp
	}
	color := rl.Fade(rl.GetColor(int32(GetStyle(Default, colorProp))), guiAlpha)

	// Draw control
	//--------------------------------------------------------------------
	DrawRectangle(bounds, PanelBorderWidth, borderColor, color)
	//--------------------------------------------------------------------
}

// Scroll Panel control
func ScrollPanel(bounds, content rl.Rectangle, scroll *rl.Vector2) rl.Rectangle {
	state := guiState

	bw := float32(GetStyle(Default, BorderWidthProp))
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
	DrawRectangle(bounds, 0, rl.Blank, rl.GetColor(int32(GetStyle(Default, BackgroundColorProp)))) // Draw background

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
	DrawRectangle(bounds, int(GetStyle(Default, BorderWidthProp)), rl.Fade(rl.GetColor(int32(GetStyle(ListViewControl, Border+(ControlProperty(state)*3)))), guiAlpha), rl.Blank)

	// Set scrollbar slider size back to the way it was before
	SetStyle(ScrollBarControl, ScrollSliderSize, slider)
	//--------------------------------------------------------------------

	if scroll != nil {
		*scroll = scrollPos
	}

	return view
}

// Label control
func Label(bounds rl.Rectangle, text string) {
	state := guiState

	// Update control
	//--------------------------------------------------------------------
	// ...
	//--------------------------------------------------------------------

	// Draw control
	//--------------------------------------------------------------------
	colorProp := TextColorNormalProp
	if state == StateDisabled {
		colorProp = TextColorDisabledProp
	}
	DrawText(text, GetTextBounds(LabelControl, bounds), TextAlignment(GetStyle(LabelControl, TextAlignmentProp)), rl.Fade(rl.GetColor(int32(GetStyle(LabelControl, colorProp))), guiAlpha))
	//--------------------------------------------------------------------
}

// Button control, returns true when clicked
func Button(bounds rl.Rectangle, text string) bool {
	state := guiState
	pressed := false

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

			if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
				pressed = true
			}
		}
	}
	//--------------------------------------------------------------------

	// Draw control
	//--------------------------------------------------------------------
	DrawRectangle(bounds, int(GetStyle(ButtonControl, BorderWidthProp)), rl.Fade(rl.GetColor(int32(GetStyle(ButtonControl, Border+(ControlProperty(state)*3)))), guiAlpha), rl.Fade(rl.GetColor(int32(GetStyle(ButtonControl, Base+(ControlProperty(state)*3)))), guiAlpha))
	DrawText(text, GetTextBounds(ButtonControl, bounds), TextAlignment(GetStyle(ButtonControl, TextAlignmentProp)), rl.Fade(rl.GetColor(int32(GetStyle(ButtonControl, Text+(ControlProperty(state)*3)))), guiAlpha))
	//------------------------------------------------------------------

	return pressed
}

// Label button control
func LabelButton(bounds rl.Rectangle, text string) bool {
	state := guiState
	pressed := false

	// NOTE: We force bounds.width to be all text
	textWidth := rl.MeasureTextEx(guiFont, text, float32(GetStyle(Default, TextSizeProp)), float32(GetStyle(Default, TextSpacingProp))).X
	if bounds.Width < textWidth {
		bounds.Width = textWidth
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

			if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
				pressed = true
			}
		}
	}
	//--------------------------------------------------------------------

	// Draw control
	//--------------------------------------------------------------------
	DrawText(text, GetTextBounds(LabelControl, bounds), TextAlignment(GetStyle(LabelControl, TextAlignmentProp)), rl.Fade(rl.GetColor(int32(GetStyle(LabelControl, Text+(ControlProperty(state)*3)))), guiAlpha))
	//--------------------------------------------------------------------

	return pressed
}

// Image button control, returns true when clicked
func ImageButton(bounds rl.Rectangle, text string, texture rl.Texture2D) bool {
	return ImageButtonEx(bounds, text, texture, rl.Rectangle{0, 0, float32(texture.Width), float32(texture.Height)})
}

// Image button control, returns true when clicked
func ImageButtonEx(bounds rl.Rectangle, text string, texture rl.Texture2D, texSource rl.Rectangle) bool {
	state := guiState
	clicked := false

	// Update control
	//--------------------------------------------------------------------
	if state != StateDisabled && !guiLocked {
		mousePoint := rl.GetMousePosition()

		// Check button state
		if rl.CheckCollisionPointRec(mousePoint, bounds) {
			if rl.IsMouseButtonDown(rl.MouseLeftButton) {
				state = StatePressed
			} else if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
				clicked = true
			} else {
				state = StateFocused
			}
		}
	}
	//--------------------------------------------------------------------

	// Draw control
	//--------------------------------------------------------------------
	DrawRectangle(bounds, int(GetStyle(ButtonControl, BorderWidthProp)), rl.Fade(rl.GetColor(int32(GetStyle(ButtonControl, Border+(ControlProperty(state)*3)))), guiAlpha), rl.Fade(rl.GetColor(int32(GetStyle(ButtonControl, Base+(ControlProperty(state)*3)))), guiAlpha))

	DrawText(text, GetTextBounds(ButtonControl, bounds), TextAlignment(GetStyle(ButtonControl, TextAlignmentProp)), rl.Fade(rl.GetColor(int32(GetStyle(ButtonControl, Text+(ControlProperty(state)*3)))), guiAlpha))
	if texture.ID > 0 {
		rl.DrawTextureRec(texture, texSource, rl.Vector2{bounds.X + bounds.Width/2 - texSource.Width/2, bounds.Y + bounds.Height/2 - texSource.Height/2}, rl.Fade(rl.GetColor(int32(GetStyle(ButtonControl, Text+(ControlProperty(state)*3)))), guiAlpha))
	}
	//------------------------------------------------------------------

	return clicked
}

// Toggle Button control, returns true when active
func Toggle(bounds rl.Rectangle, text string, active bool) bool {
	state := guiState

	// Update control
	//--------------------------------------------------------------------
	if state != StateDisabled && !guiLocked {
		mousePoint := rl.GetMousePosition()

		// Check toggle button state
		if rl.CheckCollisionPointRec(mousePoint, bounds) {
			if rl.IsMouseButtonDown(rl.MouseLeftButton) {
				state = StatePressed
			} else if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
				state = StateNormal
				active = !active
			} else {
				state = StateFocused
			}
		}
	}
	//--------------------------------------------------------------------

	// Draw control
	//--------------------------------------------------------------------
	if state == StateNormal {
		var borderColorProp, baseColorProp, textColorProp ControlProperty
		if active {
			borderColorProp = BorderColorPressedProp
			baseColorProp = BaseColorPressedProp
			textColorProp = TextColorPressedProp
		} else {
			borderColorProp = Border + ControlProperty(state)*3
			baseColorProp = Base + ControlProperty(state)*3
			textColorProp = Base + ControlProperty(state)*3
		}
		DrawRectangle(bounds, int(GetStyle(ToggleControl, BorderWidthProp)), rl.Fade(rl.GetColor(int32(GetStyle(ToggleControl, borderColorProp))), guiAlpha), rl.Fade(rl.GetColor(int32(GetStyle(ToggleControl, baseColorProp))), guiAlpha))
		DrawText(text, GetTextBounds(ToggleControl, bounds), TextAlignment(GetStyle(ToggleControl, TextAlignmentProp)), rl.Fade(rl.GetColor(int32(GetStyle(ToggleControl, textColorProp))), guiAlpha))
	} else {
		DrawRectangle(bounds, int(GetStyle(ToggleControl, BorderWidthProp)), rl.Fade(rl.GetColor(int32(GetStyle(ToggleControl, Border+ControlProperty(state)*3))), guiAlpha), rl.Fade(rl.GetColor(int32(GetStyle(ToggleControl, Base+ControlProperty(state)*3))), guiAlpha))
		DrawText(text, GetTextBounds(ToggleControl, bounds), TextAlignment(GetStyle(ToggleControl, TextAlignmentProp)), rl.Fade(rl.GetColor(int32(GetStyle(ToggleControl, Text+ControlProperty(state)*3))), guiAlpha))
	}
	//--------------------------------------------------------------------

	return active
}

const ToggleGroupMaxElements = 32

// Toggle Group control, returns toggled button index
func ToggleGroup(bounds rl.Rectangle, text string, active int) int {
	initBoundsX := bounds.X

	// Get substrings items from text (items pointers)
	var rows [ToggleGroupMaxElements]int
	itemCount := 0
	items := GuiTextSplit(text, &itemCount, rows)

	prevRow := rows[0]

	for i := 0; i < itemCount; i++ {
		if prevRow != rows[i] {
			bounds.X = initBoundsX
			bounds.Y += bounds.Height + float32(GetStyle(ToggleControl, GroupPadding))
			prevRow = rows[i]
		}

		if i == active {
			Toggle(bounds, items[i], true)
		} else if Toggle(bounds, items[i], false) {
			active = i
		}

		bounds.X += bounds.Width + float32(GetStyle(ToggleControl, GroupPadding))
	}

	return active
}

//----------------------------------------------------------------------------------
// Module specific Functions Definition
//----------------------------------------------------------------------------------

// Gui get text width using default font
func GetTextWidth(text string) int {
	var size rl.Vector2

	if text != "" {
		size = rl.MeasureTextEx(guiFont, text, float32(GetStyle(Default, TextSizeProp)), float32(GetStyle(Default, TextSpacingProp)))
	}

	// TODO: Consider text icon width here???

	return int(size.X)
}

// Get text bounds considering control bounds
func GetTextBounds(control Control, bounds rl.Rectangle) rl.Rectangle {
	// TODO(port)
}

/*
Rectangle GetTextBounds(int control, Rectangle bounds)
{
    Rectangle textBounds = bounds;

    textBounds.x = bounds.x + GuiGetStyle(control, BORDER_WIDTH);
    textBounds.y = bounds.y + GuiGetStyle(control, BORDER_WIDTH);
    textBounds.width = bounds.width - 2*GuiGetStyle(control, BORDER_WIDTH);
    textBounds.height = bounds.height - 2*GuiGetStyle(control, BORDER_WIDTH);

    // Consider TEXT_PADDING properly, depends on control type and TEXT_ALIGNMENT
    switch (control)
    {
        case COMBOBOX: bounds.width -= (GuiGetStyle(control, COMBO_BUTTON_WIDTH) + GuiGetStyle(control, COMBO_BUTTON_PADDING)); break;
        case VALUEBOX: break;   // NOTE: ValueBox text value always centered, text padding applies to label
        default:
        {
            if (GuiGetStyle(control, TEXT_ALIGNMENT) == GUI_TEXT_ALIGN_RIGHT) textBounds.x -= GuiGetStyle(control, TEXT_PADDING);
            else textBounds.x += GuiGetStyle(control, TEXT_PADDING);
        } break;
    }

    // TODO: Special cases (no label): COMBOBOX, DROPDOWNBOX, LISTVIEW (scrollbar?)
    // More special cases (label side): CHECKBOX, SLIDER, VALUEBOX, SPINNER

    return textBounds;
}
*/

// Gui draw text using default font
func DrawText(text string, bounds rl.Rectangle, alignment TextAlignment, tint rl.Color) {
	// TODO(port)
}

/*
void GuiDrawText(const char *text, Rectangle bounds, int alignment, Color tint)
{
    #define TEXT_VALIGN_PIXEL_OFFSET(h)  ((int)h%2)     // Vertical alignment for pixel perfect

    if ((text != NULL) && (text[0] != '\0'))
    {
        int iconId = 0;
        text = GetTextIcon(text, &iconId);  // Check text for icon and move cursor

        // Get text position depending on alignment and iconId
        //---------------------------------------------------------------------------------
        #define RICON_TEXT_PADDING   4

        Vector2 position = { bounds.x, bounds.y };

        // NOTE: We get text size after icon been processed
        int textWidth = GetTextWidth(text);
        int textHeight = GuiGetStyle(DEFAULT, TEXT_SIZE);

        // If text requires an icon, add size to measure
        if (iconId >= 0)
        {
            textWidth += RICON_SIZE;

            // WARNING: If only icon provided, text could be pointing to eof character!
            if ((text != NULL) && (text[0] != '\0')) textWidth += RICON_TEXT_PADDING;
        }

        // Check guiTextAlign global variables
        switch (alignment)
        {
            case GUI_TEXT_ALIGN_LEFT:
            {
                position.x = bounds.x;
                position.y = bounds.y + bounds.height/2 - textHeight/2 + TEXT_VALIGN_PIXEL_OFFSET(bounds.height);
            } break;
            case GUI_TEXT_ALIGN_CENTER:
            {
                position.x = bounds.x + bounds.width/2 - textWidth/2;
                position.y = bounds.y + bounds.height/2 - textHeight/2 + TEXT_VALIGN_PIXEL_OFFSET(bounds.height);
            } break;
            case GUI_TEXT_ALIGN_RIGHT:
            {
                position.x = bounds.x + bounds.width - textWidth;
                position.y = bounds.y + bounds.height/2 - textHeight/2 + TEXT_VALIGN_PIXEL_OFFSET(bounds.height);
            } break;
            default: break;
        }

        // NOTE: Make sure we get pixel-perfect coordinates,
        // In case of decimals we got weird text positioning
        position.x = (float)((int)position.x);
        position.y = (float)((int)position.y);
        //---------------------------------------------------------------------------------

        // Draw text (with icon if available)
        //---------------------------------------------------------------------------------
#if defined(RAYGUI_SUPPORT_RICONS)
        if (iconId >= 0)
        {
            // NOTE: We consider icon height, probably different than text size
            GuiDrawIcon(iconId, RAYGUI_CLITERAL(Vector2){ position.x, bounds.y + bounds.height/2 - RICON_SIZE/2 + TEXT_VALIGN_PIXEL_OFFSET(bounds.height) }, 1, tint);
            position.x += (RICON_SIZE + RICON_TEXT_PADDING);
        }
#endif
        DrawTextEx(guiFont, text, position, (float)GuiGetStyle(DEFAULT, TEXT_SIZE), (float)GuiGetStyle(DEFAULT, TEXT_SPACING), tint);
        //---------------------------------------------------------------------------------
    }
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
