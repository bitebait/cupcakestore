package models

import (
	"strings"
	"testing"
)

func TestPasswordChangeRequiresCurrentPasswordAndMinimumLength(t *testing.T) {
	user := User{Password: "original-password"}
	if err := user.HashPassword(); err != nil {
		t.Fatal(err)
	}
	original := user.Password
	for _, tc := range []struct{ old, next string }{
		{"wrong-password", "valid-new-password"},
		{"original-password", "short"},
		{"original-password", strings.Repeat("x", 73)},
	} {
		if err := user.UpdatePassword(tc.old, tc.next); err == nil {
			t.Fatal("invalid password change accepted")
		}
		if user.Password != original {
			t.Fatal("rejected change modified password")
		}
	}
	if err := user.UpdatePassword("original-password", "new-password"); err != nil {
		t.Fatal(err)
	}
	if err := user.CheckPassword("new-password"); err != nil {
		t.Fatal(err)
	}
}
