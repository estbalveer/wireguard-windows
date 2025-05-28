/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2019-2022 WireGuard LLC. All Rights Reserved.
 */

package ui

import (
	"github.com/lxn/walk"
)

type VPNKeyDialog struct {
	*walk.Dialog
	vpnKeyEdit *walk.LineEdit
	connectBtn *walk.PushButton
	cancelBtn  *walk.PushButton
	result     string
}

func runVPNKeyDialog(owner walk.Form) (string, bool) {
	dlg, err := newVPNKeyDialog(owner)
	if err != nil {
		return "", false
	}

	if dlg.Run() == walk.DlgCmdOK {
		return dlg.result, true
	}

	return "", false
}

func newVPNKeyDialog(owner walk.Form) (*VPNKeyDialog, error) {
	var err error
	var disposables walk.Disposables
	defer disposables.Treat()

	dlg := new(VPNKeyDialog)

	if dlg.Dialog, err = walk.NewDialog(owner); err != nil {
		return nil, err
	}
	disposables.Add(dlg)

	dlg.SetTitle("Secure VPN Connection")
	dlg.SetLayout(walk.NewVBoxLayout())
	dlg.SetMinMaxSize(walk.Size{400, 200}, walk.Size{200, 200})
	dlg.SetSize(walk.Size{400, 200})

	// Subtitle
	subtitleLabel, err := walk.NewTextLabel(dlg)
	if err != nil {
		return nil, err
	}
	subtitleLabel.SetText("Enter your VPN server key to establish a secure and private connection.")
	subtitleLabel.SetTextColor(walk.RGB(100, 100, 100))

	// VPN Key Input
	if dlg.vpnKeyEdit, err = walk.NewLineEdit(dlg); err != nil {
		return nil, err
	}
	dlg.vpnKeyEdit.SetMinMaxSize(walk.Size{200, 0}, walk.Size{200, 0})

	// Buttons container
	buttonsContainer, err := walk.NewComposite(dlg)
	if err != nil {
		return nil, err
	}
	buttonsContainer.SetLayout(walk.NewHBoxLayout())
	walk.NewHSpacer(buttonsContainer)

	// Connect button
	if dlg.connectBtn, err = walk.NewPushButton(buttonsContainer); err != nil {
		return nil, err
	}
	dlg.connectBtn.SetText("Connect")
	dlg.connectBtn.Clicked().Attach(func() {
		dlg.result = dlg.vpnKeyEdit.Text()
		dlg.Accept()
	})

	// Cancel button
	if dlg.cancelBtn, err = walk.NewPushButton(buttonsContainer); err != nil {
		return nil, err
	}
	dlg.cancelBtn.SetText("Cancel")
	dlg.cancelBtn.Clicked().Attach(dlg.Cancel)

	dlg.SetDefaultButton(dlg.connectBtn)
	dlg.SetCancelButton(dlg.cancelBtn)

	disposables.Spare()

	return dlg, nil
}
