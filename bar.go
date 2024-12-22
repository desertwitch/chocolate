package chocolate

// chocolateBar is the main workhorse that will
// provide most of the functionality and is doing
// all the calculations and handling of the layout
// it further holds the tea.Models and wrap the
// calls around, so that it acts at the end just
// a view container
// type chocolateBar struct {
// 	id string
//
// 	// backref to the Chocolate the bar
// 	// belongs to
// 	// this is used to centralize some
// 	// parts for the unified theming
// 	// and controls like selector
// 	choc *Chocolate
//
// 	// bars in order for the layout
// 	// backref to the parent bar
// 	// this is used to propagate the
// 	// dynamic sizing back to the parent
// 	// as well to also adjust depending
// 	// on the parent layout
// 	parent *chocolateBar
//
// 	// layout parameters
// 	layoutType LayoutType
//
// 	// possible maximum content size
// 	// This is used to have a maximum for the content
// 	// the real size will be calculated during the
// 	// view rendering as it is the only possible
// 	// place to handle the scaling which depends on
// 	// possible dynmic content
// 	maxWidth  int
// 	maxHeight int
//
// 	// content size after calculation of the whole
// 	// layout
// 	width  int
// 	height int
//
// 	// model of the entry this can only be set, if
// 	// there are no sub bars and is the final leaf
// 	// of the whole tree and provides the real content
// 	// model tea.Model
//
// 	// models map so that the bar can have multiple
// 	// models to select from
// 	// models map[string]*BarModel
// 	// running actModel
// 	// actModel *BarModel
//
// 	// pre rendered view with maximum content sizes
// 	// this is used to get the correct sizes of the
// 	// models view to be used for dynamic scaling
// 	// and the following calculations
// 	preRendered   bool
// 	preView       string
// 	contentWidth  int
// 	contentHeight int
//
// 	// rendered
// 	view     string
// 	rendered bool
//
// 	// flavourPrefs generation function
// 	// this can be used to override the default
// 	// flavour preferences
// 	FlavourCustomzieHandler BarFlavourCustomizeHandlerFct
//
// 	// custom update function
// 	// this can be used to override the default
// 	// behavior which will only let the bar
// 	// take input focus when a model is attached
// 	// and just pass the tea messages through
// 	// UpdateHandlerFct func(*ChocolateBar) func(tea.Msg) tea.Cmd
//
// 	// if the bar is hidden
// 	// hidden bars are removed from the layout
// 	// rendering and the space is used for the
// 	// other bars
// 	hidden bool
//
// 	// if this bar can be selected
// 	selectable bool
// 	// if this bar should receive input when
// 	// selected
// 	inputOnSelect bool
// }
//
// func (b *chocolateBar) GetStyle() lipgloss.Style {
// 	ret := flavour.GetPresetNoErr(flavour.PRESET_PRIMARY_NOBORDER)
//
// 	if b.HasModel() || b.IsRoot() {
// 		ret = flavour.GetPresetNoErr(flavour.PRESET_PRIMARY)
// 		// ret = ret.BorderType(b.GetChoc().GetFlavour().GetBorderType())
// 	}
// 	if b.GetChocolate().IsSelected(b) && !b.IsRoot() {
// 		ret = ret.BorderForeground(flavour.GetColorNoErr(flavour.COLOR_SECONDARY))
// 	}
// 	if b.GetChocolate().IsFocused(b) && !b.IsRoot() {
// 		ret = flavour.GetPresetNoErr(flavour.PRESET_SECONDARY).
// 			BorderBackground(flavour.GetColorNoErr(flavour.COLOR_PRIMARY_BG))
// 	}
//
// 	if b.HasModel() && b.actModel.FlavourCustomizeHandler != nil {
// 		ret = b.actModel.FlavourCustomizeHandler(b, b.actModel.Model, ret)()
// 	} else if b.FlavourCustomzieHandler != nil {
// 		ret = b.FlavourCustomzieHandler(b, ret)()
// 	}
//
// 	return ret
// }
//
// func (b chocolateBar) IsRoot() bool {
// 	return b.parent == nil
// }
//
// func (b chocolateBar) GetParent() CChocolateBar {
// 	return b.parent
// }
//
// func (b chocolateBar) GetScaling() Scaling {
// 	return Scaling{b.X, b.Y}
// }
//
// func (b *chocolateBar) SetScaling(v Scaling) {
// 	b.X = v.X
// 	b.Y = v.Y
// }
//
// func (b chocolateBar) GetLayout() LayoutType {
// 	return b.layoutType
// }
//
// func (b *chocolateBar) setChocolate(v *Chocolate) {
// 	b.choc = v
// 	for _, c := range b.bars {
// 		c.(*chocolateBar).setChocolate(v)
// 	}
// }
//
// func (b *chocolateBar) SetChocolate(v *Chocolate) {
// 	if b.IsRoot() {
// 		b.setChocolate(v)
// 	} else {
// 		b.parent.SetChocolate(v)
// 	}
// }
//
// func (b chocolateBar) CanFocus() bool {
// 	return b.actModel != nil
// }
//
// func (b chocolateBar) CanSelect() bool {
// 	return b.selectable
// }
//
// func (b *chocolateBar) GetSelectables() []string {
// 	var ret []string
// 	for _, c := range b.bars {
// 		ret = append(ret, c.GetSelectables()...)
// 	}
// 	if b.CanSelect() {
// 		ret = append(ret, b.GetID())
// 	}
// 	if len(ret) == 0 && b.IsRoot() {
// 		ret = append(ret, b.GetID())
// 	}
// 	return ret
// }
//
// func (b chocolateBar) GetBars() []CChocolateBar {
// 	return b.bars
// }
//
// func (b chocolateBar) InputOnSelect() bool {
// 	return b.inputOnSelect
// }
//
// func (b chocolateBar) GetChocolate() *Chocolate {
// 	return b.choc
// }
//
// func (b chocolateBar) HasModel() bool {
// 	if b.actModel != nil {
// 		return b.actModel.Model != nil
// 	}
// 	return false
// }
//
// func (b chocolateBar) GetModel() tea.Model {
// 	if b.HasModel() {
// 		return b.actModel.Model
// 	}
// 	return nil
// }
//
// func (b chocolateBar) GetID() string {
// 	return b.id
// }
//
// func (b *chocolateBar) SelectModel(v string) {
// 	if b.models == nil {
// 		return
// 	}
// 	if m, ok := b.models[v]; ok {
// 		b.actModel = m
// 	}
// }
//
// func (b *chocolateBar) Resize(w, h int) {
// 	// if there is a frame set for the bar
// 	// this has to be removed from the available
// 	// content size
// 	width := w - b.GetStyle().GetHorizontalFrameSize()
// 	height := h - b.GetStyle().GetVerticalFrameSize()
//
// 	// if this is a fixed scaling than we don't have
// 	// to calculate anything
// 	if b.X.IsFixed() {
// 		width = b.X.GetValue()
// 		// b.width = width
// 	}
// 	if b.Y.IsFixed() {
// 		height = b.Y.GetValue()
// 		// b.height = height
// 	}
//
// 	b.maxWidth = width
// 	b.maxHeight = height
//
// 	// the root bar doesn't have to rescale itself
// 	if b.IsRoot() {
// 		b.width = width
// 		b.height = height
// 	}
//
// 	if b.models != nil {
// 		for _, m := range b.models {
// 			m.Model, _ = m.Model.Update(tea.WindowSizeMsg{Width: width, Height: height})
// 		}
// 	}
// 	if b.HasModel() {
// 		b.actModel.Model, _ = b.actModel.Model.Update(tea.WindowSizeMsg{Width: width, Height: height})
// 	} else {
// 		for _, c := range b.bars {
// 			c.Resize(width, height)
// 		}
// 	}
// }
//
// // pre render all models with their actual sizes
// // this is a temporary task that has to be done
// // so that it is possible to calculate the dynamic
// // sizes
// // TODO: Is there a better way to avoid calling models view?
// func (b *chocolateBar) preRender() {
// 	// skip hidden bars
// 	if b.hidden {
// 		return
// 	}
//
// 	if b.HasModel() {
// 		if !b.preRendered {
// 			b.preView = b.actModel.Model.View()
// 			b.contentWidth, b.contentHeight = lipgloss.Size(b.preView)
//
// 			b.preRendered = true
//
// 			if !b.IsRoot() {
// 				t, v := b.X.Get()
// 				switch t {
// 				case DYNAMIC:
// 					b.parent.contentWidth += b.contentWidth + b.GetStyle().GetHorizontalFrameSize()
// 					b.width = b.contentWidth
// 				case FIXED:
// 					b.parent.contentWidth += v + b.GetStyle().GetHorizontalFrameSize()
// 					b.width = v
// 				}
// 				t, v = b.Y.Get()
// 				switch t {
// 				case DYNAMIC:
// 					b.parent.contentHeight += b.contentHeight + b.GetStyle().GetVerticalFrameSize()
// 					b.height = b.contentHeight
// 				case FIXED:
// 					b.parent.contentHeight += v + b.GetStyle().GetVerticalFrameSize()
// 					b.height = v
// 				}
// 			}
// 		}
// 		return
// 	}
//
// 	// must be a bar without model
// 	// so go recursive to generate
// 	// all preViews of models
// 	for _, c := range b.bars {
// 		c.(*chocolateBar).preRender()
// 	}
//
// 	// all sub bars of this model are now pre rendered
// 	// we can build up the used sizes of the fixed
// 	// and dynamic sub bars
// 	if !b.IsRoot() {
// 		t, v := b.X.Get()
// 		switch t {
// 		case DYNAMIC:
// 			b.parent.contentWidth += b.contentWidth + b.GetStyle().GetHorizontalFrameSize()
// 		case FIXED:
// 			b.parent.contentHeight += v + b.GetStyle().GetHorizontalFrameSize()
// 		}
// 		t, v = b.Y.Get()
// 		switch t {
// 		case DYNAMIC:
// 			b.parent.contentHeight += b.contentHeight + b.GetStyle().GetVerticalFrameSize()
// 		case FIXED:
// 			b.parent.contentHeight += v + b.GetStyle().GetVerticalFrameSize()
// 		}
// 	}
// }
//
// func (b *chocolateBar) recalcSizes() {
// 	// skip hidden bars
// 	if b.hidden {
// 		return
// 	}
//
// 	// already done so just return
// 	if b.preRendered {
// 		return
// 	}
//
// 	switch b.layoutType {
// 	case LIST:
// 		b.recalcVerticalSizes()
// 	case LINEAR:
// 		b.recalcHorizontalSizes()
// 	}
// }
//
// func (b *chocolateBar) recalcVerticalSizes() {
// 	// after pre render all leafs with models
// 	// this must be a bar holding subs
// 	// so go recursive till we reach the last
// 	// layers
// 	for _, c := range b.bars {
// 		c.(*chocolateBar).recalcSizes()
// 	}
//
// 	// go over again and start calculation
// 	totalParts := 0
// 	totalParents := 0
// 	for _, c := range b.bars {
// 		if c.(*chocolateBar).Y.IsParent() && !c.(*chocolateBar).hidden {
// 			totalParts += c.(*chocolateBar).Y.GetValue()
// 			totalParents++
// 		}
// 	}
//
// 	if totalParts > 0 {
// 		partSize := (b.maxHeight - b.contentHeight) / totalParts
// 		partLast := (b.maxHeight - b.contentHeight) % totalParts
//
// 		for _, c := range b.bars {
// 			if c.(*chocolateBar).Y.IsParent() && !c.(*chocolateBar).hidden {
// 				totalParents--
// 				height := c.(*chocolateBar).Y.GetValue() * partSize
// 				if totalParents == 0 {
// 					height += partLast
// 				}
// 				c.(*chocolateBar).height = height - c.GetStyle().GetVerticalFrameSize()
// 				b.contentHeight += c.(*chocolateBar).height
// 			}
// 		}
// 	}
//
// 	if !b.IsRoot() {
// 		b.height = b.contentHeight
// 		b.parent.contentHeight += b.height
// 	}
//
// 	b.preRendered = true
// }
//
// func (b *chocolateBar) recalcHorizontalSizes() {
// 	// after pre render all leafs with models
// 	// this must be a bar holding subs
// 	// so go recursive till we reach the last
// 	// layers
// 	for _, c := range b.bars {
// 		c.(*chocolateBar).recalcSizes()
// 	}
//
// 	// go over again and start calculation
// 	totalParts := 0
// 	totalParents := 0
// 	for _, c := range b.bars {
// 		if c.(*chocolateBar).X.IsParent() && !c.(*chocolateBar).hidden {
// 			totalParts += c.(*chocolateBar).X.GetValue()
// 			totalParents++
// 		}
// 	}
//
// 	if totalParts > 0 {
// 		partSize := (b.maxWidth - b.contentWidth) / totalParts
// 		partLast := (b.maxWidth - b.contentWidth) % totalParts
//
// 		for _, c := range b.bars {
// 			if c.(*chocolateBar).X.IsParent() && !c.(*chocolateBar).hidden {
// 				totalParents--
// 				width := c.(*chocolateBar).X.GetValue() * partSize
// 				if totalParents == 0 {
// 					width += partLast
// 				}
// 				c.(*chocolateBar).width = width - c.GetStyle().GetHorizontalFrameSize()
// 				b.contentWidth += c.(*chocolateBar).width
// 			}
// 		}
// 	}
//
// 	if !b.IsRoot() {
// 		b.width = b.contentWidth
// 		b.parent.contentWidth += b.width
// 	}
//
// 	b.preRendered = true
// }
//
// func (b *chocolateBar) finalizeSizing() {
// 	// skip hidden bars
// 	if b.hidden {
// 		return
// 	}
//
// 	for _, c := range b.bars {
// 		c.(*chocolateBar).finalizeSizing()
// 	}
//
// 	if !b.IsRoot() {
// 		width := b.parent.width
// 		height := b.parent.height
// 		if width <= 0 {
// 			width = b.parent.maxWidth
// 		}
// 		if height <= 0 {
// 			height = b.parent.maxHeight
// 		}
// 		if b.width <= 0 {
// 			b.width = width - b.GetStyle().GetHorizontalFrameSize()
// 		}
// 		if b.height <= 0 {
// 			b.height = height - b.GetStyle().GetVerticalFrameSize()
// 		}
//
// 	}
// 	if b.HasModel() {
// 		b.actModel.Model, _ = b.actModel.Model.Update(tea.WindowSizeMsg{Width: b.width, Height: b.height})
// 	}
// }
//
// func (b *chocolateBar) render() {
// 	b.preRender()
// 	b.recalcSizes()
// 	b.finalizeSizing()
//
// 	// skip hidden bars
// 	if b.hidden {
// 		return
// 	}
//
// 	if b.HasModel() {
// 		b.view = b.GetStyle().
// 			Width(b.width).
// 			Height(b.height).
// 			Render(b.actModel.Model.View())
// 		b.rendered = true
// 		return
// 	}
//
// 	for _, c := range b.bars {
// 		c.(*chocolateBar).render()
// 	}
// }
//
// func (b *chocolateBar) joinBars() {
// 	// skip hidden bars
// 	if b.hidden {
// 		return
// 	}
//
// 	if b.rendered {
// 		return
// 	}
//
// 	switch b.layoutType {
// 	case LIST:
// 		b.joinVerticalBars()
// 	case LINEAR:
// 		b.joinHorizontalBars()
// 	}
// }
//
// func (b *chocolateBar) joinVerticalBars() {
// 	var bars []string
// 	if !b.rendered {
// 		for _, c := range b.bars {
// 			c.(*chocolateBar).joinBars()
// 			if c.(*chocolateBar).hidden {
// 				continue
// 			}
// 			s := b.GetStyle().
// 				BorderTop(false).
// 				BorderBottom(false).
// 				BorderLeft(false).
// 				BorderRight(false).
// 				Width(b.width)
// 			bars = append(bars, s.Render(c.(*chocolateBar).view))
// 		}
// 		s := b.GetStyle()
// 		if b.IsRoot() {
// 			s = s.Height(b.height)
// 		}
// 		b.view = s.
// 			Render(lipgloss.JoinVertical(0, bars...))
// 	}
// 	b.rendered = true
// }
//
// func (b *chocolateBar) joinHorizontalBars() {
// 	var bars []string
// 	if !b.rendered {
// 		for _, c := range b.bars {
// 			c.(*chocolateBar).joinBars()
// 			if c.(*chocolateBar).hidden {
// 				continue
// 			}
// 			s := b.GetStyle().
// 				BorderTop(false).
// 				BorderBottom(false).
// 				BorderLeft(false).
// 				BorderRight(false).
// 				Height(b.height)
// 			bars = append(bars, s.Render(c.(*chocolateBar).view))
// 		}
// 		s := b.GetStyle()
// 		if b.IsRoot() {
// 			s = s.Width(b.width)
// 		}
// 		b.view = s.
// 			Render(lipgloss.JoinHorizontal(0, bars...))
// 	}
// 	b.rendered = true
// }
//
// func (b *chocolateBar) resetRender() {
// 	for _, c := range b.bars {
// 		c.(*chocolateBar).resetRender()
// 	}
//
// 	// the root bar must not reset it's size
// 	if !b.IsRoot() {
// 		b.width = 0
// 		b.height = 0
// 	}
// 	//  if !b.IsRoot() || !b.X.IsFixed() {
// 	// 	b.width = 0
// 	// }
// 	// if !b.IsRoot() || !b.Y.IsFixed() {
// 	// 	b.height = 0
// 	// }
// 	b.preRendered = false
// 	b.contentHeight = 0
// 	b.contentWidth = 0
// 	b.preView = ""
// 	b.rendered = false
// 	b.view = ""
// }
//
// func (b *chocolateBar) Render() string {
// 	defer b.resetRender()
//
// 	b.resetRender()
// 	b.render()
// 	b.joinBars()
// 	w, h := lipgloss.Size(b.view)
// 	w -= b.GetStyle().GetHorizontalFrameSize()
// 	h -= b.GetStyle().GetVerticalFrameSize()
// 	if w > b.width || h > b.height {
// 		return "Window too small"
// 	}
// 	return b.view
// }
//
// func (b *chocolateBar) Hide(v bool) {
// 	b.hidden = v
// }
//
// func (b *chocolateBar) defaultUpdateHandler(msg tea.Msg) tea.Cmd {
// 	var cmds []tea.Cmd
// 	var cmd tea.Cmd
//
// 	switch msg := msg.(type) {
// 	case ModelChangeMsg:
// 		b.SelectModel(msg.Model)
// 	}
//
// 	if b.HasModel() {
// 		b.actModel.Model, cmd = b.actModel.Model.Update(msg)
// 		cmds = append(cmds, cmd)
// 	}
// 	return tea.Batch(cmds...)
// }
//
// func (b *chocolateBar) HandleUpdate(msg tea.Msg) tea.Cmd {
// 	var cmds []tea.Cmd
//
// 	cmds = append(cmds, b.defaultUpdateHandler(msg))
// 	if b.actModel.UpdateHandlerFct != nil {
// 		cmds = append(cmds, b.actModel.UpdateHandlerFct(b, b.actModel.Model)(msg))
// 	}
// 	return tea.Batch(cmds...)
// }
//
// type ChocolateBarOptions func(*chocolateBar)
//
// func WithLayout(v LayoutType) func(*chocolateBar) {
// 	return func(b *chocolateBar) {
// 		b.layoutType = v
// 	}
// }
//
// func WithID(v string) func(*chocolateBar) {
// 	return func(b *chocolateBar) {
// 		b.id = v
// 	}
// }
//
// func WithModels(v map[string]*BarModel, a string) func(*chocolateBar) {
// 	return func(b *chocolateBar) {
// 		b.models = v
// 		b.actModel = v[a]
// 		b.bars = nil
// 	}
// }
//
// func WithModel(v *BarModel) func(*chocolateBar) {
// 	return func(b *chocolateBar) {
// 		b.actModel = v
// 		b.bars = nil
// 	}
// }
//
// func WithSelectable() func(*chocolateBar) {
// 	return func(b *chocolateBar) {
// 		b.selectable = true
// 	}
// }
//
// func Hidden() func(*chocolateBar) {
// 	return func(b *chocolateBar) {
// 		b.hidden = true
// 	}
// }
//
// func WithFlavourCustomizeHandler(v BarFlavourCustomizeHandlerFct) func(*chocolateBar) {
// 	return func(b *chocolateBar) {
// 		b.FlavourCustomzieHandler = v
// 	}
// }
//
// func WithInputOnSelect() func(*chocolateBar) {
// 	return func(b *chocolateBar) {
// 		b.inputOnSelect = true
// 	}
// }
//
// func NewChocolateBar(bars []CChocolateBar, opts ...ChocolateBarOptions) *chocolateBar {
// 	ret := &chocolateBar{
// 		id:            uuid.NewString(),
// 		layoutType:    LIST,
// 		preRendered:   false,
// 		preView:       "",
// 		rendered:      false,
// 		view:          "",
// 		width:         0,
// 		height:        0,
// 		contentWidth:  0,
// 		contentHeight: 0,
// 		hidden:        false,
// 		selectable:    false,
// 		inputOnSelect: false,
// 	}
//
// 	for _, c := range bars {
// 		c.(*chocolateBar).parent = ret
// 	}
//
// 	for _, opt := range opts {
// 		opt(ret)
// 	}
//
// 	return ret
// }
