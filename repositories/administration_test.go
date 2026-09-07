package repositories

import (
	"errors"
	"sync"
	"testing"

	"github.com/bitebait/cupcakestore/models"
	"gorm.io/gorm"
)

func TestLastActiveAdministratorCannotBeRemoved(t *testing.T) {
	for _, action := range []string{"deactivate", "demote", "delete"} {
		t.Run(action, func(t *testing.T) {
			db, profile, _ := commerceDB(t)
			if err := db.Model(&models.User{}).Where("id = ?", profile.UserID).Update("is_staff", true).Error; err != nil {
				t.Fatal(err)
			}
			repo := NewUserRepository(db)
			admin, err := repo.FindById(profile.UserID)
			if err != nil {
				t.Fatal(err)
			}
			switch action {
			case "deactivate":
				admin.IsActive = false
				err = repo.Update(&admin)
			case "demote":
				admin.IsStaff = false
				err = repo.Update(&admin)
			case "delete":
				err = repo.Delete(&admin)
			}
			if !errors.Is(err, ErrLastAdministrator) {
				t.Fatalf("lost last administrator: %v", err)
			}
			var active int64
			if err := db.Model(&models.User{}).Where("is_active = ? AND is_staff = ?", true, true).Count(&active).Error; err != nil || active != 1 {
				t.Fatalf("active administrators = %d, %v", active, err)
			}
		})
	}
}

func TestConcurrentAdministratorDemotionsKeepAccess(t *testing.T) {
	db, profile, _ := commerceDB(t)
	if err := db.Model(&models.User{}).Where("id = ?", profile.UserID).Update("is_staff", true).Error; err != nil {
		t.Fatal(err)
	}
	second := models.User{Email: "second@example.com", IsStaff: true, IsActive: true}
	if err := db.Session(&gorm.Session{SkipHooks: true}).Create(&second).Error; err != nil {
		t.Fatal(err)
	}
	repo := NewUserRepository(db)
	first, err := repo.FindById(profile.UserID)
	if err != nil {
		t.Fatal(err)
	}
	start := make(chan struct{})
	var wait sync.WaitGroup
	for _, admin := range []models.User{first, second} {
		wait.Add(1)
		go func(admin models.User) {
			defer wait.Done()
			<-start
			admin.IsStaff = false
			_ = repo.Update(&admin) // A concurrent SQLite writer may safely fail with SQLITE_BUSY.
		}(admin)
	}
	close(start)
	wait.Wait()
	var active int64
	if err := db.Model(&models.User{}).Where("is_active = ? AND is_staff = ?", true, true).Count(&active).Error; err != nil || active < 1 {
		t.Fatalf("concurrent updates removed administrative access: %d, %v", active, err)
	}
}

func TestSettingsRejectStaleUpdatesAndNeverInsert(t *testing.T) {
	db, _, _ := commerceDB(t)
	repo := NewStoreConfigRepository(db)
	current, err := repo.GetStoreConfig()
	if err != nil {
		t.Fatal(err)
	}
	stale := current
	current.DeliveryPrice = 4
	if err := repo.Update(&current); err != nil {
		t.Fatal(err)
	}
	stale.PhysicalStoreCity = "Outra cidade"
	if err := repo.Update(&stale); !errors.Is(err, ErrStoreConfigChanged) {
		t.Fatalf("stale settings accepted: %v", err)
	}
	if err := repo.Update(&models.StoreConfig{}); err == nil {
		t.Fatal("update created settings")
	}
	saved, err := repo.GetStoreConfig()
	if err != nil || saved.DeliveryPrice != 4 || saved.PhysicalStoreCity != current.PhysicalStoreCity {
		t.Fatalf("lost settings: %+v, %v", saved, err)
	}
}
