// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package encryptview

import (
	"microsoft.com/pkggen/imagegen/attendedinstaller/primitives/navigationbar"
	"microsoft.com/pkggen/imagegen/attendedinstaller/uitext"
	"microsoft.com/pkggen/imagegen/attendedinstaller/uiutils"
	"microsoft.com/pkggen/imagegen/configuration"

	"github.com/gdamore/tcell"
	"github.com/muesli/crunchy"
	"github.com/rivo/tview"
)

// UI constants.
const (
	navButtonNext = 1
	noSelection   = -1

	formProportion = 0

	passwordFieldWidth = 64
)

// EncryptView contains the encrypt UI
type EncryptView struct {
	form                 *tview.Form
	passwordField        *tview.InputField
	confirmPasswordField *tview.InputField
	navBar               *navigationbar.NavigationBar
	flex                 *tview.Flex
	centeredFlex         *tview.Flex
	passwordValidator    *crunchy.Validator

	encryption *configuration.RootEncryption
}

// New creates and returns a new EncryptView.
func New() *EncryptView {
	return &EncryptView{
		passwordValidator: crunchy.NewValidator(),
	}
}

// Initialize initializes the view.
func (ev *EncryptView) Initialize(backButtonText string, sysConfig *configuration.SystemConfig, cfg *configuration.Config, app *tview.Application, nextPage, previousPage, quit, refreshTitle func()) (err error) {

	ev.encryption = &sysConfig.Encryption

	ev.passwordField = tview.NewInputField().
		SetFieldWidth(passwordFieldWidth).
		SetLabel(uitext.EncryptPasswordLabel).
		SetMaskCharacter('*')

	ev.confirmPasswordField = tview.NewInputField().
		SetFieldWidth(passwordFieldWidth).
		SetLabel(uitext.ConfirmEncryptPasswordLabel).
		SetMaskCharacter('*')

	ev.navBar = navigationbar.NewNavigationBar().
		AddButton(backButtonText, previousPage).
		AddButton(uitext.ButtonNext, func() {
			ev.onNextButton(nextPage, cfg)
		}).
		AddButton(uitext.SkipEncryption, func() {
			ev.encryption.Enable = false
			ev.encryption.Password = ""
			nextPage()
		}).
		SetAlign(tview.AlignCenter).
		SetOnFocusFunc(func() {
			ev.navBar.SetSelectedButton(navButtonNext)
		}).
		SetOnBlurFunc(func() {
			ev.navBar.SetSelectedButton(noSelection)
		})

	ev.form = tview.NewForm().
		SetButtonsAlign(tview.AlignCenter).
		AddFormItem(ev.passwordField).
		AddFormItem(ev.confirmPasswordField).
		AddFormItem(ev.navBar)

	ev.flex = tview.NewFlex().
		SetDirection(tview.FlexRow)

	formWidth, formHeight := uiutils.MinFormSize(ev.form)
	centeredForm := uiutils.CenterHorizontally(formWidth, ev.form)

	ev.flex.AddItem(centeredForm, formHeight+ev.navBar.GetHeight(), formProportion, true)
	ev.centeredFlex = uiutils.CenterVerticallyDynamically(ev.flex)

	// Box styling
	ev.centeredFlex.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)

	return
}

// HandleInput handles custom input.
func (ev *EncryptView) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	ev.navBar.ClearUserFeedback()

	// Allow Up-Down to navigate the form
	switch event.Key() {
	case tcell.KeyUp:
		return tcell.NewEventKey(tcell.KeyBacktab, 0, tcell.ModNone)
	case tcell.KeyDown:
		return tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)
	}

	return event
}

// Reset resets the page, undoing any user input.
func (ev *EncryptView) Reset() (err error) {
	ev.navBar.ClearUserFeedback()
	ev.navBar.SetSelectedButton(noSelection)

	ev.encryption.Password = ""
	ev.encryption.Enable = false

	return
}

// Name returns the friendly name of the view.
func (ev *EncryptView) Name() string {
	return "ENCRYPTDISK"
}

// Title returns the title of the view.
func (ev *EncryptView) Title() string {
	return uitext.EncryptTitle
}

// Primitive returns the primary primitive to be rendered for the view.
func (ev *EncryptView) Primitive() tview.Primitive {
	return ev.centeredFlex
}

// OnShow gets called when the view is shown to the user
func (ev *EncryptView) OnShow() {
}

func (ev *EncryptView) onNextButton(nextPage func(), cfg *configuration.Config) {
	enteredPassword := ev.passwordField.GetText()

	if enteredPassword != ev.confirmPasswordField.GetText() {
		ev.navBar.SetUserFeedback(uitext.PasswordMismatchFeedback, tview.Styles.TertiaryTextColor)
		return
	}

	err := ev.passwordValidator.Check(enteredPassword)
	if err != nil {
		ev.navBar.SetUserFeedback(uiutils.ErrorToUserFeedback(err), tview.Styles.TertiaryTextColor)
		return
	}

	ev.encryption.Enable = true
	ev.encryption.Password = enteredPassword
	nextPage()
}
