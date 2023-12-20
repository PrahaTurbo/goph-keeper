// Package tui provides the terminal user interface for the application
package tui

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/PrahaTurbo/goph-keeper/api/proto"
	"github.com/PrahaTurbo/goph-keeper/internal/client/config"
)

// Application holds all the components necessary for the terminal interface of the application.
type Application struct {
	appContext     context.Context
	secretsClient  pb.SecretClient
	authClient     pb.AuthClient
	secretText     *tview.TextView
	createForm     *tview.Form
	errorWindow    *tview.Modal
	secretsPanel   *tview.Flex
	secretsList    *tview.List
	startMenu      *tview.Modal
	secretsDetails *tview.Flex
	authForm       *tview.Form
	editForm       *tview.Form
	deleteWindow   *tview.Modal
	selectedSecret *pb.SecretData
	Pages          *tview.Pages
	App            *tview.Application
	authStatus     string
	secrets        []*pb.SecretData
}

// NewApplication is a constructor function for Application.
// It initializes Application with necessary tview and ProtoBuf clients. It also sets
// up the pages and the menu in this function.
func NewApplication(authClient pb.AuthClient, secretsClient pb.SecretClient) Application {
	c := Application{
		App:            tview.NewApplication(),
		Pages:          tview.NewPages(),
		startMenu:      tview.NewModal(),
		authForm:       tview.NewForm(),
		errorWindow:    tview.NewModal(),
		secretsPanel:   tview.NewFlex(),
		secretsDetails: tview.NewFlex(),
		secretsList:    tview.NewList(),
		secretText:     tview.NewTextView(),
		createForm:     tview.NewForm(),
		editForm:       tview.NewForm(),
		deleteWindow:   tview.NewModal(),
		authClient:     authClient,
		secretsClient:  secretsClient,
	}

	c.setupPages()
	c.setupStartMenu()
	c.setupSecretsPanel()

	c.appContext = context.Background()

	return c
}

func (a *Application) setupPages() {
	a.Pages.AddPage(startMenuPageName, a.startMenu, true, true)
	a.Pages.AddPage(authPageName, a.authForm, true, false)
	a.Pages.AddPage(errorWindowName, a.errorWindow, true, false)
	a.Pages.AddPage(secretsPanelPageName, a.secretsPanel, true, false)
	a.Pages.AddPage(createPageName, a.createForm, true, false)
	a.Pages.AddPage(editPageName, a.editForm, true, false)
	a.Pages.AddPage(deleteWindowName, a.deleteWindow, true, false)
}

func (a *Application) setupStartMenu() {
	a.startMenu.SetText("Authorization").
		AddButtons([]string{loginLabel, signUpLabel, quitLabel}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == quitLabel {
				a.App.Stop()
			}

			a.authStatus = buttonLabel

			a.authForm.Clear(true)
			a.setupAuthForm()
			a.Pages.SwitchToPage(authPageName)
		})
}

func (a *Application) setupSecretsPanel() {
	footer := tview.NewTextView()
	footer.SetText(fmt.Sprintf("Build Version: %s, Build Date: %s", config.BuildVersion, config.BuildDate))
	footer.SetTextAlign(tview.AlignRight)

	secretsListFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	a.secretsDetails.SetDirection(tview.FlexRow)

	quitButton := newButton(quitLabel, a.App.Stop)
	createButton := newButton(createLabel, a.addCreateForm)
	syncButton := newButton(syncLabel, a.addSecretsList)
	editButton := newButton(editLabel, a.addEditForm)
	deleteButton := newButton(deleteLabel, a.addDeleteWindow)
	deleteButton.SetStyle(tcell.StyleDefault.Background(tcell.ColorRed))

	a.secretsPanel.SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(secretsListFlex, 0, 2, true).
			AddItem(a.secretsDetails, 0, 4, false), 0, 6, true).
		AddItem(footer, 1, 0, true)

	secretsListFlex.Box = tview.NewBox().SetBorder(true).SetTitle("Secrets")
	secretsListFlex.AddItem(tview.NewFlex().
		AddItem(quitButton, 0, 1, false).
		AddItem(tview.NewBox(), 1, 0, false).
		AddItem(createButton, 0, 1, false).
		AddItem(tview.NewBox(), 1, 0, false).
		AddItem(syncButton, 0, 1, false), 1, 0, false).
		AddItem(a.secretsList, 0, 10, true)

	a.secretsDetails.Box = tview.NewBox().SetBorder(true).SetTitle("Details")

	a.secretText.SetDynamicColors(true)

	a.secretsList.SetBorderPadding(1, 0, 0, 0)
	a.secretsList.SetSelectedFunc(func(index int, name string, secondName string, shortcut rune) {
		a.secretsDetails.Clear()
		a.secretsDetails.AddItem(tview.NewFlex().
			AddItem(editButton, 0, 1, false).
			AddItem(tview.NewBox(), 1, 0, false).
			AddItem(deleteButton, 0, 1, false).
			AddItem(tview.NewBox(), 0, 4, false), 1, 0, false).
			AddItem(tview.NewBox(), 1, 0, true).
			AddItem(a.secretText, 0, 10, true)

		a.selectedSecret = a.secrets[index]
		a.setSecretText(a.secrets[index])
	})
}

func (a *Application) addDeleteWindow() {
	a.deleteWindow.ClearButtons()
	a.Pages.SwitchToPage(deleteWindowName)

	a.deleteWindow.SetText("Are you sure?").
		AddButtons([]string{deleteLabel, backLabel}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonIndex {
			case 0:
				req := &pb.DeleteRequest{SecretId: a.selectedSecret.Id}
				_, err := a.secretsClient.Delete(a.appContext, req)
				if err != nil {
					s := status.Convert(err)
					a.addErrorWindow(fmt.Sprintf("%s: %s", s.Err(), s.Message()), secretsPanelPageName)
					return
				}

				a.addSecretsList()
				a.Pages.SwitchToPage(secretsPanelPageName)
			case 1:
				a.Pages.SwitchToPage(secretsPanelPageName)
			}
		})
}

func (a *Application) addEditForm() {
	a.editForm.Clear(true)
	a.Pages.SwitchToPage(editPageName)

	req := &pb.UpdateRequest{
		SecretId: a.selectedSecret.Id,
		Type:     a.selectedSecret.Type,
		Content:  a.selectedSecret.Content,
		MetaData: a.selectedSecret.MetaData,
	}

	typeList := []string{
		pb.SecretType_CREDENTIALS.String(),
		pb.SecretType_TEXT.String(),
		pb.SecretType_CARD.String(),
		pb.SecretType_BINARY.String(),
	}

	var initialOption int
	for i := range typeList {
		if typeList[i] == a.selectedSecret.Type.String() {
			initialOption = i
		}
	}

	a.editForm.AddDropDown("Type", typeList, initialOption, func(option string, optionIndex int) {
		secretType := pb.SecretType_UNSPECIFIED
		if v, ok := pb.SecretType_value[option]; ok {
			secretType = pb.SecretType(v)
		}

		req.Type = secretType
	})

	a.editForm.AddTextArea("Content", a.selectedSecret.Content, 40, 0, 0, func(text string) {
		req.Content = text
	})

	a.editForm.AddTextArea("Additional info", a.selectedSecret.MetaData, 40, 0, 0, func(text string) {
		req.MetaData = text
	})

	a.editForm.AddButton(updateLabel, func() {
		_, err := a.secretsClient.Update(a.appContext, req)
		if err != nil {
			s := status.Convert(err)
			a.addErrorWindow(fmt.Sprintf("%s: %s", s.Err(), s.Message()), editPageName)
			return
		}

		a.addSecretsList()
		a.Pages.SwitchToPage(secretsPanelPageName)
	})

	a.editForm.AddButton(backLabel, func() {
		a.Pages.SwitchToPage(secretsPanelPageName)
	})
}

func (a *Application) addCreateForm() {
	a.createForm.Clear(true)
	a.Pages.SwitchToPage(createPageName)

	req := &pb.CreateRequest{}

	a.createForm.AddDropDown("Type *", []string{
		pb.SecretType_CREDENTIALS.String(),
		pb.SecretType_TEXT.String(),
		pb.SecretType_CARD.String(),
		pb.SecretType_BINARY.String(),
	}, -1, func(option string, optionIndex int) {
		secretType := pb.SecretType_UNSPECIFIED
		if v, ok := pb.SecretType_value[option]; ok {
			secretType = pb.SecretType(v)
		}

		req.Type = secretType
	})

	a.createForm.AddTextArea("Content *", "", 40, 0, 0, func(text string) {
		req.Content = text
	})

	a.createForm.AddTextArea("Additional info", "", 40, 0, 0, func(text string) {
		req.MetaData = text
	})

	a.createForm.AddButton(saveLabel, func() {
		if req.Type == pb.SecretType_UNSPECIFIED || req.Content == "" {
			a.addErrorWindow("You have to fill all required fields", createPageName)
			return
		}

		_, err := a.secretsClient.Create(a.appContext, req)
		if err != nil {
			s := status.Convert(err)
			a.addErrorWindow(fmt.Sprintf("%s: %s", s.Err(), s.Message()), createPageName)
			return
		}

		a.addSecretsList()
		a.Pages.SwitchToPage(secretsPanelPageName)
	})

	a.createForm.AddButton(backLabel, func() {
		a.Pages.SwitchToPage(secretsPanelPageName)
	})
}

func (a *Application) setSecretText(secret *pb.SecretData) {
	a.secretText.Clear()

	text := fmt.Sprintf("[green]TYPE[white]\n%s\n\n", secret.Type) +
		fmt.Sprintf("[green]CONTENT[white]\n%s\n\n", secret.Content)

	if secret.MetaData != "" {
		text += fmt.Sprintf("[green]META DATA[white]\n%s\n\n", secret.MetaData)
	}

	a.secretText.SetText(text)
}

func (a *Application) addSecretsList() {
	a.secretsList.Clear()
	a.secretText.Clear()
	a.secretsDetails.Clear()

	resp, err := a.secretsClient.GetSecrets(a.appContext, &pb.GetSecretsRequest{})
	if err != nil {
		s := status.Convert(err)
		a.addErrorWindow(fmt.Sprintf("%s: %s", s.Err(), s.Message()), secretsPanelPageName)
		return
	}

	a.secrets = resp.Secrets
	for i, s := range resp.Secrets {
		a.secretsList.AddItem(s.Type.String(), s.Content, rune(49+i), nil)
	}
}

func (a *Application) setupAuthForm() {
	req := &pb.AuthRequest{}

	a.authForm.AddInputField("Login", "", 20, nil, func(login string) {
		req.Login = login
	})

	a.authForm.AddPasswordField("Password", "", 20, '*', func(password string) {
		req.Password = password
	})

	a.authForm.AddButton(submitLabel, func() {
		var resp *pb.AuthResponse
		var err error

		switch a.authStatus {
		case loginLabel:
			resp, err = a.authClient.Login(context.Background(), req)
			if err != nil {
				s := status.Convert(err)
				a.addErrorWindow(fmt.Sprintf("%s: %s", s.Err(), s.Message()), authPageName)
				return
			}
		default:
			resp, err = a.authClient.Register(context.Background(), req)
			if err != nil {
				s := status.Convert(err)
				a.addErrorWindow(fmt.Sprintf("%s: %s", s.Err(), s.Message()), authPageName)
				return
			}
		}

		md := metadata.Pairs(authentication, fmt.Sprintf("%s %s", bearerSchema, resp.Token))
		a.appContext = metadata.NewOutgoingContext(a.appContext, md)

		a.addSecretsList()
		a.Pages.SwitchToPage(secretsPanelPageName)
	})

	a.authForm.AddButton(backLabel, func() {
		a.Pages.SwitchToPage(startMenuPageName)
	})
}

func (a *Application) addErrorWindow(err string, parentPage string) {
	a.errorWindow.ClearButtons()
	a.errorWindow.SetBackgroundColor(tcell.ColorRed)

	a.Pages.SwitchToPage(errorWindowName)

	a.errorWindow.SetText(err).
		AddButtons([]string{okLabel}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.Pages.SwitchToPage(parentPage)
		})
}
