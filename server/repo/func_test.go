package repo

import (
	"testing"

	D "github.com/NeoJRotary/describe-go"
)

func TestResolveSrc(t *testing.T) {
	src := "https://github.com/NeoJRotary/FormulaCI.git"
	hub, user, name, err := ResolveSrc(src)
	if D.IsErr(err) {
		t.Error(err)
		return
	}
	if hub != "github.com" || user != "NeoJRotary" || name != "FormulaCI" {
		t.Error("Wrong result : ", hub, user, name, " . with src : ", src)
	}

	src = "git@github.com:NeoJRotary/FormulaCI.git"
	hub, user, name, err = ResolveSrc(src)
	if D.IsErr(err) {
		t.Error(err)
		return
	}
	if hub != "github.com" || user != "NeoJRotary" || name != "FormulaCI" {
		t.Error("Wrong result : ", hub, user, name, " . with src : ", src)
	}

	src = "git@gitlab.com:google.com/golang-beta.git"
	hub, user, name, err = ResolveSrc(src)
	if D.IsErr(err) {
		t.Error(err)
		return
	}
	if hub != "gitlab.com" || user != "google.com" || name != "golang-beta" {
		t.Error("Wrong result : ", hub, user, name, " . with src : ", src)
	}
}
