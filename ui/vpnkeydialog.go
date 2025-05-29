/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2019-2022 WireGuard LLC. All Rights Reserved.
 */

package ui

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/lxn/walk"
)

type VPNKeyDialog struct {
	*walk.Dialog
	vpnKeyEdit  *walk.LineEdit
	connectBtn  *walk.PushButton
	cancelBtn   *walk.PushButton
	progressBar *walk.ProgressBar
	result      string
	apiResponse *APIResponse
}

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		UUID      string `json:"uuid"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		UserType  string `json:"user_type"`
		IsActive  int    `json:"is_active"`
		UserCode  string `json:"user_code"`
		Server    struct {
			UUID                string `json:"uuid"`
			Address             string `json:"address"`
			DNS                 string `json:"dns"`
			PrivateKey          string `json:"private_key"`
			PublicKey           string `json:"public_key"`
			PresharedKey        string `json:"preshared_key"`
			AllowedIPs          string `json:"allowed_ips"`
			Endpoint            string `json:"endpoint"`
			PersistentKeepAlive string `json:"persistent_keep_alive"`
		} `json:"server"`
	} `json:"data"`
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
	dlg.SetMinMaxSize(walk.Size{400, 250}, walk.Size{200, 250})
	dlg.SetSize(walk.Size{400, 250})

	// Subtitle
	subtitleLabel, err := walk.NewTextLabel(dlg)
	if err != nil {
		return nil, err
	}
	subtitleLabel.SetText("Enter your VPN server key to establish a secure\n and private connection.")
	subtitleLabel.SetTextColor(walk.RGB(100, 100, 100))

	// VPN Key Input
	if dlg.vpnKeyEdit, err = walk.NewLineEdit(dlg); err != nil {
		return nil, err
	}
	dlg.vpnKeyEdit.SetMinMaxSize(walk.Size{200, 0}, walk.Size{200, 0})

	// Add text changed handler to convert to uppercase
	dlg.vpnKeyEdit.TextChanged().Attach(func() {
		text := dlg.vpnKeyEdit.Text()
		// Convert to uppercase
		upperText := strings.ToUpper(text)
		// Remove any non-alphanumeric characters
		cleanText := ""
		for _, char := range upperText {
			if (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
				cleanText += string(char)
			}
		}
		// Only update if the text has changed
		if text != cleanText {
			// Get current cursor position
			start, _ := dlg.vpnKeyEdit.TextSelection()
			// Set new text
			dlg.vpnKeyEdit.SetText(cleanText)
			// Restore cursor position
			dlg.vpnKeyEdit.SetTextSelection(start, start)
		}
	})

	// Progress Bar
	if dlg.progressBar, err = walk.NewProgressBar(dlg); err != nil {
		return nil, err
	}
	dlg.progressBar.SetVisible(false)
	dlg.progressBar.SetMarqueeMode(true)

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
		key := dlg.vpnKeyEdit.Text()
		if key == "" {
			walk.MsgBox(dlg, "Error", "Please enter a valid key", walk.MsgBoxIconError)
			return
		}

		dlg.result = key
		dlg.progressBar.SetVisible(true)
		dlg.connectBtn.SetEnabled(false)
		dlg.cancelBtn.SetEnabled(false)

		go func() {
			client := &http.Client{
				Timeout: 10 * time.Second,
			}

			apiURL := fmt.Sprintf("http://185.237.100.130/api/%s", key)
			resp, err := client.Get(apiURL)
			if err != nil {
				walk.MsgBox(dlg, "Error", "Failed to connect to server: "+err.Error(), walk.MsgBoxIconError)
				dlg.progressBar.SetVisible(false)
				dlg.connectBtn.SetEnabled(true)
				dlg.cancelBtn.SetEnabled(true)
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				walk.MsgBox(dlg, "Error", "Failed to read server response: "+err.Error(), walk.MsgBoxIconError)
				dlg.progressBar.SetVisible(false)
				dlg.connectBtn.SetEnabled(true)
				dlg.cancelBtn.SetEnabled(true)
				return
			}

			var apiResp APIResponse
			if err := json.Unmarshal(body, &apiResp); err != nil {
				walk.MsgBox(dlg, "Error", "Failed to parse server response: "+err.Error(), walk.MsgBoxIconError)
				dlg.progressBar.SetVisible(false)
				dlg.connectBtn.SetEnabled(true)
				dlg.cancelBtn.SetEnabled(true)
				return
			}

			if !apiResp.Success {
				walk.MsgBox(dlg, "Error", apiResp.Message, walk.MsgBoxIconError)
				dlg.progressBar.SetVisible(false)
				dlg.connectBtn.SetEnabled(true)
				dlg.cancelBtn.SetEnabled(true)
				return
			}

			dlg.apiResponse = &apiResp
			dlg.Accept()
		}()
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
