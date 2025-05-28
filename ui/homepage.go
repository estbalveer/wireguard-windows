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

	// Create a vertical layout for the home page
	vlayout := walk.NewVBoxLayout()
	vlayout.SetMargins(walk.Margins{10, 10, 10, 10})
	vlayout.SetSpacing(10)
	hp.SetLayout(vlayout)

	// Set dark background color
	if brush, err := walk.NewSolidColorBrush(walk.RGB(2, 12, 27)); err == nil {
		hp.SetBackground(brush)
		disposables.Add(brush)
	}

	// Add a welcome label
	welcomeLabel, err := walk.NewTextLabel(hp)
	if err != nil {
		return nil, err
	}
	welcomeLabel.SetText("Welcome to CloakStream")
	if font, err := walk.NewFont("Segoe UI", 16, walk.FontBold); err == nil {
		welcomeLabel.SetFont(font)
	}
	welcomeLabel.SetTextColor(walk.RGB(255, 255, 255)) // White text

	// Add a description label
	descLabel, err := walk.NewTextLabel(hp)
	if err != nil {
		return nil, err
	}
	descLabel.SetText("Your secure connection manager")
	if font, err := walk.NewFont("Segoe UI", 10, 0); err == nil {
		descLabel.SetFont(font)
	}
	descLabel.SetTextColor(walk.RGB(200, 200, 200)) // Light gray text

	// Create a composite for centering the image
	composite, err := walk.NewComposite(hp)
	if err != nil {
		return nil, err
	}
	composite.SetLayout(walk.NewHBoxLayout())
	composite.SetMinMaxSize(walk.Size{Width: 300, Height: 300}, walk.Size{Width: 300, Height: 300})

	// Add the image view
	hp.imageView, err = walk.NewImageView(composite)
	if err != nil {
		return nil, err
	}
	hp.imageView.SetMinMaxSize(walk.Size{Width: 200, Height: 200}, walk.Size{Width: 200, Height: 200})
	hp.imageView.SetMode(walk.ImageViewModeShrink)

	// Load and set the GIF
	if img, err := walk.NewImageFromFile("resources/map.gif"); err == nil {
		hp.imageView.SetImage(img)
		disposables.Add(img)
	}

	disposables.Spare()

	return hp, nil
}
