package tui

import "github.com/rivo/tview"

const (
	bearerSchema   = "bearer"
	authentication = "authorization"
)

const (
	startMenuPageName    = "StartMenuPage"
	authPageName         = "AuthPage"
	secretsPanelPageName = "SecretsPanelPage"
	errorWindowName      = "ErrorWindow"
	createPageName       = "CreatePageName"
	editPageName         = "EditPageName"
	deleteWindowName     = "DeleteWindow"
)

const (
	loginLabel  = "Login"
	signUpLabel = "Sign Up"
	okLabel     = "OK"
	submitLabel = "Submit"
	backLabel   = "Back"
	quitLabel   = "Quit"
	updateLabel = "Update"
	saveLabel   = "Save"
	createLabel = "Create"
	syncLabel   = "Sync"
	deleteLabel = "Delete"
	editLabel   = "Edit"
)

func newButton(label string, selectedFunc func()) *tview.Button {
	return tview.NewButton(label).SetSelectedFunc(func() {
		selectedFunc()
	})
}
