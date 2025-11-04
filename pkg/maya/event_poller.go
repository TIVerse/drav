package maya

import (
	"context"
	"fmt"

	"github.com/TIVerse/drav/pkg/agni"
	"github.com/gdamore/tcell/v2"
)

// EventPoller polls tcell events and emits them to the event hub.
type EventPoller struct {
	driver     *TcellDriver
	dispatcher *agni.Dispatcher
}

// NewEventPoller creates a new event poller.
func NewEventPoller(driver *TcellDriver, dispatcher *agni.Dispatcher) *EventPoller {
	return &EventPoller{
		driver:     driver,
		dispatcher: dispatcher,
	}
}

// Poll polls events from tcell and emits them to the dispatcher.
// This should be called in a goroutine.
func (p *EventPoller) Poll(ctx context.Context) error {
	screen := p.driver.Screen()
	if screen == nil {
		return fmt.Errorf("screen not initialized")
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Poll event with timeout to check context
			ev := screen.PollEvent()
			if ev == nil {
				continue
			}

			// Convert tcell event to agni event
			agniEvent := p.convertEvent(ev)
			if agniEvent != nil {
				// Emit to dispatcher
				if err := p.dispatcher.Emit(ctx, agniEvent); err != nil {
					// Log error but continue
					continue
				}
			}
		}
	}
}

// convertEvent converts a tcell event to an agni event.
func (p *EventPoller) convertEvent(ev tcell.Event) agni.Event {
	switch e := ev.(type) {
	case *tcell.EventKey:
		return p.convertKeyEvent(e)
	case *tcell.EventMouse:
		return p.convertMouseEvent(e)
	case *tcell.EventResize:
		return p.convertResizeEvent(e)
	case *tcell.EventInterrupt:
		return agni.NewQuitEvent("interrupt")
	default:
		return nil
	}
}

// convertKeyEvent converts a tcell key event.
func (p *EventPoller) convertKeyEvent(e *tcell.EventKey) agni.Event {
	key, r := p.tcellKeyToAgniKey(e.Key(), e.Rune())
	mods := p.tcellModsToAgniMods(e.Modifiers())
	
	return agni.NewKeyEvent(key, r, mods, false)
}

// tcellKeyToAgniKey converts tcell key to agni key.
func (p *EventPoller) tcellKeyToAgniKey(k tcell.Key, r rune) (agni.Key, rune) {
	switch k {
	case tcell.KeyRune:
		return agni.KeyRune, r
	case tcell.KeyEnter:
		return agni.KeyEnter, 0
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		return agni.KeyBackspace, 0
	case tcell.KeyTab:
		return agni.KeyTab, 0
	case tcell.KeyEscape:
		return agni.KeyEscape, 0
	case tcell.KeyUp:
		return agni.KeyUp, 0
	case tcell.KeyDown:
		return agni.KeyDown, 0
	case tcell.KeyLeft:
		return agni.KeyLeft, 0
	case tcell.KeyRight:
		return agni.KeyRight, 0
	case tcell.KeyHome:
		return agni.KeyHome, 0
	case tcell.KeyEnd:
		return agni.KeyEnd, 0
	case tcell.KeyPgUp:
		return agni.KeyPageUp, 0
	case tcell.KeyPgDn:
		return agni.KeyPageDown, 0
	case tcell.KeyDelete:
		return agni.KeyDelete, 0
	case tcell.KeyInsert:
		return agni.KeyInsert, 0
	case tcell.KeyF1:
		return agni.KeyF1, 0
	case tcell.KeyF2:
		return agni.KeyF2, 0
	case tcell.KeyF3:
		return agni.KeyF3, 0
	case tcell.KeyF4:
		return agni.KeyF4, 0
	case tcell.KeyF5:
		return agni.KeyF5, 0
	case tcell.KeyF6:
		return agni.KeyF6, 0
	case tcell.KeyF7:
		return agni.KeyF7, 0
	case tcell.KeyF8:
		return agni.KeyF8, 0
	case tcell.KeyF9:
		return agni.KeyF9, 0
	case tcell.KeyF10:
		return agni.KeyF10, 0
	case tcell.KeyF11:
		return agni.KeyF11, 0
	case tcell.KeyF12:
		return agni.KeyF12, 0
	case tcell.KeyCtrlC:
		return agni.KeyCtrlC, 0
	case tcell.KeyCtrlD:
		return agni.KeyCtrlD, 0
	case tcell.KeyCtrlZ:
		return agni.KeyCtrlZ, 0
	default:
		return agni.KeyRune, r
	}
}

// tcellModsToAgniMods converts tcell modifiers to agni modifiers.
func (p *EventPoller) tcellModsToAgniMods(m tcell.ModMask) agni.ModMask {
	var mods agni.ModMask
	if m&tcell.ModShift != 0 {
		mods |= agni.ModShift
	}
	if m&tcell.ModCtrl != 0 {
		mods |= agni.ModCtrl
	}
	if m&tcell.ModAlt != 0 {
		mods |= agni.ModAlt
	}
	if m&tcell.ModMeta != 0 {
		mods |= agni.ModMeta
	}
	return mods
}

// convertMouseEvent converts a tcell mouse event.
func (p *EventPoller) convertMouseEvent(e *tcell.EventMouse) agni.Event {
	x, y := e.Position()
	button := p.tcellButtonToAgniButton(e.Buttons())
	action := p.tcellActionToAgniAction(e.Buttons())
	mods := p.tcellModsToAgniMods(e.Modifiers())
	
	// Handle wheel events
	wheel := 0
	if e.Buttons()&tcell.WheelUp != 0 {
		wheel = 1
	} else if e.Buttons()&tcell.WheelDown != 0 {
		wheel = -1
	}
	
	return agni.NewMouseEvent(x, y, button, action, wheel, mods)
}

// tcellButtonToAgniButton converts tcell button to agni button.
func (p *EventPoller) tcellButtonToAgniButton(b tcell.ButtonMask) agni.MouseButton {
	if b&tcell.Button1 != 0 {
		return agni.MouseLeft
	}
	if b&tcell.Button2 != 0 {
		return agni.MouseMiddle
	}
	if b&tcell.Button3 != 0 {
		return agni.MouseRight
	}
	return agni.MouseNone
}

// tcellActionToAgniAction converts tcell button state to agni action.
func (p *EventPoller) tcellActionToAgniAction(b tcell.ButtonMask) agni.MouseAction {
	if b&tcell.WheelUp != 0 || b&tcell.WheelDown != 0 {
		return agni.MouseScroll
	}
	if b&tcell.Button1 != 0 || b&tcell.Button2 != 0 || b&tcell.Button3 != 0 {
		return agni.MousePress
	}
	return agni.MouseMove
}

// convertResizeEvent converts a tcell resize event.
func (p *EventPoller) convertResizeEvent(e *tcell.EventResize) agni.Event {
	width, height := e.Size()
	return agni.NewResizeEvent(width, height)
}
