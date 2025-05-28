/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2019-2022 WireGuard LLC. All Rights Reserved.
 */

package ui

import (
	"github.com/lxn/walk"
	"golang.zx2c4.com/wireguard/windows/l18n"
)

type HomePage struct {
	*walk.TabPage
	imageView *walk.ImageView
}

const (
	DarkBlue = 0xFF020C1B
	Black    = 0xFF000000
)

func NewHomePage() (*HomePage, error) {
	var err error
	var disposables walk.Disposables
	defer disposables.Treat()

	hp := new(HomePage)

	if hp.TabPage, err = walk.NewTabPage(); err != nil {
		return nil, err
	}
	disposables.Add(hp)

	hp.SetTitle(l18n.Sprintf("Home"))

	// Main vertical layout
	mainLayout := walk.NewVBoxLayout()
	mainLayout.SetMargins(walk.Margins{0, 0, 0, 0})
	mainLayout.SetSpacing(0)
	hp.SetLayout(mainLayout)

	// Set dark background color
	if brush, err := walk.NewSolidColorBrush(walk.RGB(2, 12, 27)); err == nil {
		hp.SetBackground(brush)
		disposables.Add(brush)
	}

	// Top bar with settings button
	topBar, _ := walk.NewComposite(hp)
	topBarLayout := walk.NewHBoxLayout()
	topBarLayout.SetMargins(walk.Margins{10, 10, 10, 0})
	topBar.SetLayout(topBarLayout)
	settingsBtn, _ := walk.NewPushButton(topBar)
	settingsBtn.SetText("⚙") // Placeholder for settings icon
	settingsBtn.SetMinMaxSize(walk.Size{32, 32}, walk.Size{32, 32})
	walk.NewHSpacer(topBar)

	// Welcome section
	welcomeComposite, _ := walk.NewComposite(hp)
	welcomeLayout := walk.NewVBoxLayout()
	welcomeLayout.SetAlignment(walk.AlignHCenterVNear)
	welcomeComposite.SetLayout(welcomeLayout)
	welcomeComposite.Layout().SetSpacing(2)

	welcomeLabel, _ := walk.NewTextLabel(welcomeComposite)
	welcomeLabel.SetText("Welcome to")
	if font, err := walk.NewFont("Segoe UI", 10, 0); err == nil {
		welcomeLabel.SetFont(font)
	}
	welcomeLabel.SetTextColor(walk.RGB(180, 180, 180))
	welcomeLabel.SetTextAlignment(walk.AlignHCenterVCenter)

	appNameLabel, _ := walk.NewTextLabel(welcomeComposite)
	appNameLabel.SetText("CloakStream")
	if font, err := walk.NewFont("Segoe UI", 22, walk.FontBold); err == nil {
		appNameLabel.SetFont(font)
	}
	appNameLabel.SetTextColor(walk.RGB(255, 255, 255))
	appNameLabel.SetTextAlignment(walk.AlignHCenterVCenter)

	statusComposite, _ := walk.NewComposite(welcomeComposite)
	statusLayout := walk.NewHBoxLayout()
	statusLayout.SetAlignment(walk.AlignHCenterVNear)
	statusComposite.SetLayout(statusLayout)

	// Add spacers to center the status text block
	walk.NewHSpacer(statusComposite)

	statusLabel, _ := walk.NewTextLabel(statusComposite)
	statusLabel.SetText("Status: ")
	statusLabel.SetTextColor(walk.RGB(180, 180, 180))
	statusValue, _ := walk.NewTextLabel(statusComposite)
	statusValue.SetText("Disconnected")
	statusValue.SetTextColor(walk.RGB(0, 122, 255)) // Blue
	if font, err := walk.NewFont("Segoe UI", 10, walk.FontBold); err == nil {
		statusValue.SetFont(font)
	}

	// Add another spacer to center the status text block
	walk.NewHSpacer(statusComposite)

	// Centered animation/icon in a circle
	centerComposite, _ := walk.NewComposite(hp)
	centerLayout := walk.NewHBoxLayout()
	centerLayout.SetAlignment(walk.AlignHCenterVCenter)
	centerComposite.SetLayout(centerLayout)
	circleComposite, _ := walk.NewComposite(centerComposite)
	circleComposite.SetLayout(walk.NewVBoxLayout())
	circleComposite.SetMinMaxSize(walk.Size{140, 140}, walk.Size{140, 140})
	if brush, err := walk.NewSolidColorBrush(walk.RGB(40, 80, 180)); err == nil {
		circleComposite.SetBackground(brush)
		disposables.Add(brush)
	}
	// Placeholder for link icon (could be a GIF or static image)
	iconLabel, _ := walk.NewTextLabel(circleComposite)
	iconLabel.SetText("🔗") // Placeholder for broken link icon
	if font, err := walk.NewFont("Segoe UI", 48, 0); err == nil {
		iconLabel.SetFont(font)
	}
	iconLabel.SetTextAlignment(walk.AlignHCenterVCenter)

	// Status message
	statusMsg, _ := walk.NewTextLabel(hp)
	statusMsg.SetText("You are disconnected")
	if font, err := walk.NewFont("Segoe UI", 16, walk.FontBold); err == nil {
		statusMsg.SetFont(font)
	}
	statusMsg.SetTextColor(walk.RGB(255, 255, 255))
	statusMsg.SetTextAlignment(walk.AlignHCenterVCenter)

	// Spacer
	spacer, _ := walk.NewComposite(hp)
	spacer.SetMinMaxSize(walk.Size{0, 40}, walk.Size{0, 40})

	// Bottom buttons
	bottomComposite, _ := walk.NewComposite(hp)
	bottomLayout := walk.NewVBoxLayout()
	bottomLayout.SetAlignment(walk.AlignHCenterVNear)
	bottomLayout.SetSpacing(10)
	bottomComposite.SetLayout(bottomLayout)

	connectBtn, _ := walk.NewPushButton(bottomComposite)
	connectBtn.SetText("Connect")
	connectBtn.SetMinMaxSize(walk.Size{260, 44}, walk.Size{260, 44})
	if font, err := walk.NewFont("Segoe UI", 14, walk.FontBold); err == nil {
		connectBtn.SetFont(font)
	}

	returnBtn, _ := walk.NewPushButton(bottomComposite)
	returnBtn.SetText("Return to Streaming app   >")
	returnBtn.SetMinMaxSize(walk.Size{260, 44}, walk.Size{260, 44})
	if font, err := walk.NewFont("Segoe UI", 12, 0); err == nil {
		returnBtn.SetFont(font)
	}

	disposables.Spare()

	return hp, nil
}
