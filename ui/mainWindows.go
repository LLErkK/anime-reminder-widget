package ui

import (
	"anime-reminder/controllers"
	"anime-reminder/models"
	"fmt"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type MainWindow struct {
	window             fyne.Window
	animeController    *controllers.AnimeController
	ringToneController *controllers.RingToneController
}

// NewMainWindow creates a new main window (receives app from main.go)
func NewMainWindow(app fyne.App) *MainWindow {
	w := app.NewWindow("Anime Reminder")
	w.Resize(fyne.NewSize(800, 600))

	return &MainWindow{
		window:             w,
		animeController:    &controllers.AnimeController{},
		ringToneController: &controllers.RingToneController{},
	}
}

// Show displays the main window
func (mw *MainWindow) Show() {
	tabs := container.NewAppTabs(
		container.NewTabItem("Anime List", mw.createAnimeListTab()),
		container.NewTabItem("Add Anime", mw.createAddAnimeTab()),
		container.NewTabItem("Ringtones", mw.createRingToneTab()),
	)

	mw.window.SetContent(tabs)
	mw.window.Show()
}

// ShowAndRun displays and runs the main window (blocking)
func (mw *MainWindow) ShowAndRun() {
	tabs := container.NewAppTabs(
		container.NewTabItem("Anime List", mw.createAnimeListTab()),
		container.NewTabItem("Add Anime", mw.createAddAnimeTab()),
		container.NewTabItem("Ringtones", mw.createRingToneTab()),
	)

	mw.window.SetContent(tabs)
	mw.window.ShowAndRun()
}

func (mw *MainWindow) createAnimeListTab() fyne.CanvasObject {
	var animeList *widget.List

	animeList = widget.NewList(
		func() int {
			animes, _ := mw.animeController.GetAllAnimes()
			return len(animes)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Template"),
				widget.NewButton("Edit", func() {}),
				widget.NewButton("Delete", func() {}),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			animes, _ := mw.animeController.GetAllAnimes()
			if id < len(animes) {
				anime := animes[id]

				cont := item.(*fyne.Container)
				label := cont.Objects[0].(*widget.Label)
				label.SetText(fmt.Sprintf("%s - %s at %s", anime.Title, anime.Day, anime.Time.Format("15:04")))

				editBtn := cont.Objects[1].(*widget.Button)
				editBtn.OnTapped = func() {
					mw.showEditAnimeDialog(anime)
				}

				deleteBtn := cont.Objects[2].(*widget.Button)
				deleteBtn.OnTapped = func() {
					mw.deleteAnime(anime.Id)
					animeList.Refresh()
				}
			}
		},
	)

	refreshBtn := widget.NewButton("Refresh", func() {
		animeList.Refresh()
	})

	return container.NewBorder(nil, refreshBtn, nil, nil, animeList)
}

func (mw *MainWindow) createAddAnimeTab() fyne.CanvasObject {
	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("Anime Title")

	daySelect := widget.NewSelect(models.Days, func(value string) {})
	daySelect.PlaceHolder = "Select Day"

	hourEntry := widget.NewEntry()
	hourEntry.SetPlaceHolder("Hour (00-23)")

	minuteEntry := widget.NewEntry()
	minuteEntry.SetPlaceHolder("Minute (00-59)")

	imagePathLabel := widget.NewLabel("No image selected")
	var selectedImagePath string

	selectImageBtn := widget.NewButton("Select Image", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()

			sourcePath := reader.URI().Path()
			uploadedPath, err := UploadImage(sourcePath)
			if err != nil {
				dialog.ShowError(fmt.Errorf("failed to upload image: %v", err), mw.window)
				return
			}

			selectedImagePath = uploadedPath
			imagePathLabel.SetText(filepath.Base(uploadedPath))
		}, mw.window)
	})

	ringTones, _ := mw.ringToneController.GetAllRingTone()
	ringToneNames := make([]string, len(ringTones))
	ringToneMap := make(map[string]uint)

	for i, rt := range ringTones {
		ringToneNames[i] = rt.Name
		ringToneMap[rt.Name] = rt.Id
	}

	var selectedRingToneId uint
	ringToneSelect := widget.NewSelect(ringToneNames, func(value string) {
		selectedRingToneId = ringToneMap[value]
	})
	ringToneSelect.PlaceHolder = "Select Ringtone"

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Title", Widget: titleEntry},
			{Text: "Day", Widget: daySelect},
			{Text: "Hour", Widget: hourEntry},
			{Text: "Minute", Widget: minuteEntry},
			{Text: "Image", Widget: container.NewVBox(selectImageBtn, imagePathLabel)},
			{Text: "Ringtone", Widget: ringToneSelect},
		},
		OnSubmit: func() {
			if titleEntry.Text == "" || daySelect.Selected == "" {
				dialog.ShowError(fmt.Errorf("please fill all required fields"), mw.window)
				return
			}

			var hour, minute int
			fmt.Sscanf(hourEntry.Text, "%d", &hour)
			fmt.Sscanf(minuteEntry.Text, "%d", &minute)

			animeTime := time.Date(0, 1, 1, hour, minute, 0, 0, time.Local)

			_, err := mw.animeController.Create(
				titleEntry.Text,
				daySelect.Selected,
				selectedImagePath,
				animeTime,
				selectedRingToneId,
			)

			if err != nil {
				dialog.ShowError(err, mw.window)
				return
			}

			dialog.ShowInformation("Success", "Anime added successfully!", mw.window)
			titleEntry.SetText("")
			daySelect.SetSelected("")
			hourEntry.SetText("")
			minuteEntry.SetText("")
			imagePathLabel.SetText("No image selected")
			selectedImagePath = ""
			ringToneSelect.ClearSelected()
		},
	}

	return container.NewVScroll(form)
}

func (mw *MainWindow) createRingToneTab() fyne.CanvasObject {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Ringtone Name")

	songPathLabel := widget.NewLabel("No file selected")
	var selectedSongPath string

	selectSongBtn := widget.NewButton("Select Audio File", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()

			sourcePath := reader.URI().Path()
			uploadedPath, err := UploadAudio(sourcePath)
			if err != nil {
				dialog.ShowError(fmt.Errorf("failed to upload audio: %v", err), mw.window)
				return
			}

			selectedSongPath = uploadedPath
			songPathLabel.SetText(filepath.Base(uploadedPath))
		}, mw.window)
	})

	var ringToneList *widget.List

	ringToneList = widget.NewList(
		func() int {
			ringTones, _ := mw.ringToneController.GetAllRingTone()
			return len(ringTones)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Template"),
				widget.NewButton("Delete", func() {}),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			ringTones, _ := mw.ringToneController.GetAllRingTone()
			if id < len(ringTones) {
				ringTone := ringTones[id]
				row := item.(*fyne.Container)

				label := row.Objects[0].(*widget.Label)
				label.SetText(fmt.Sprintf("%s - %s",
					ringTone.Name,
					filepath.Base(ringTone.SongPath),
				))

				deleteBtn := row.Objects[1].(*widget.Button)
				deleteBtn.OnTapped = func() {
					dialog.ShowConfirm("Delete", "Are you sure you want to delete this ringtone?",
						func(ok bool) {
							if ok {
								mw.ringToneController.DeleteRingTone(ringTone.Id)
								ringToneList.Refresh()
							}
						}, mw.window)
				}
			}
		},
	)

	addBtn := widget.NewButton("Add Ringtone", func() {
		if nameEntry.Text == "" || selectedSongPath == "" {
			dialog.ShowError(fmt.Errorf("please fill all fields"), mw.window)
			return
		}

		_, err := mw.ringToneController.Create(nameEntry.Text, selectedSongPath)
		if err != nil {
			dialog.ShowError(err, mw.window)
			return
		}

		dialog.ShowInformation("Success", "Ringtone added successfully!", mw.window)

		nameEntry.SetText("")
		songPathLabel.SetText("No file selected")
		selectedSongPath = ""

		ringToneList.Refresh()
	})

	form := container.NewVBox(
		widget.NewLabel("Add New Ringtone"),
		nameEntry,
		selectSongBtn,
		songPathLabel,
		addBtn,
		widget.NewSeparator(),
		widget.NewLabel("Ringtone List"),
		ringToneList,
	)

	return container.NewVScroll(form)
}

func (mw *MainWindow) showEditAnimeDialog(anime models.Anime) {
	titleEntry := widget.NewEntry()
	titleEntry.SetText(anime.Title)

	daySelect := widget.NewSelect(models.Days, func(value string) {})
	daySelect.SetSelected(anime.Day)

	hourEntry := widget.NewEntry()
	hourEntry.SetText(fmt.Sprintf("%02d", anime.Time.Hour()))

	minuteEntry := widget.NewEntry()
	minuteEntry.SetText(fmt.Sprintf("%02d", anime.Time.Minute()))

	imagePathLabel := widget.NewLabel(anime.ImagePath)
	selectedImagePath := anime.ImagePath

	selectImageBtn := widget.NewButton("Change Image", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()

			sourcePath := reader.URI().Path()
			uploadedPath, err := UploadImage(sourcePath)
			if err != nil {
				dialog.ShowError(fmt.Errorf("failed to upload image: %v", err), mw.window)
				return
			}

			selectedImagePath = uploadedPath
			imagePathLabel.SetText(filepath.Base(uploadedPath))
		}, mw.window)
	})

	ringTones, _ := mw.ringToneController.GetAllRingTone()
	ringToneNames := make([]string, len(ringTones))
	ringToneMap := make(map[string]uint)

	for i, rt := range ringTones {
		ringToneNames[i] = rt.Name
		ringToneMap[rt.Name] = rt.Id
	}

	selectedRingToneId := anime.RingToneId
	ringToneSelect := widget.NewSelect(ringToneNames, func(value string) {
		selectedRingToneId = ringToneMap[value]
	})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Title", Widget: titleEntry},
			{Text: "Day", Widget: daySelect},
			{Text: "Hour", Widget: hourEntry},
			{Text: "Minute", Widget: minuteEntry},
			{Text: "Image", Widget: container.NewVBox(selectImageBtn, imagePathLabel)},
			{Text: "Ringtone", Widget: ringToneSelect},
		},
		OnSubmit: func() {
			var hour, minute int
			fmt.Sscanf(hourEntry.Text, "%d", &hour)
			fmt.Sscanf(minuteEntry.Text, "%d", &minute)

			animeTime := time.Date(0, 1, 1, hour, minute, 0, 0, time.Local)

			_, err := mw.animeController.UpdateAnime(
				titleEntry.Text,
				daySelect.Selected,
				selectedImagePath,
				animeTime,
				anime.Id,
				selectedRingToneId,
			)

			if err != nil {
				dialog.ShowError(err, mw.window)
				return
			}

			dialog.ShowInformation("Success", "Anime updated successfully!", mw.window)
		},
	}

	d := dialog.NewCustom("Edit Anime", "Close", form, mw.window)
	d.Resize(fyne.NewSize(400, 500))
	d.Show()
}

func (mw *MainWindow) deleteAnime(id uint) {
	confirm := dialog.NewConfirm("Delete Anime",
		"Are you sure you want to delete this anime?",
		func(ok bool) {
			if ok {
				err := mw.animeController.DeleteAnime(id)
				if err != nil {
					dialog.ShowError(err, mw.window)
				}
			}
		}, mw.window)
	confirm.Show()
}
