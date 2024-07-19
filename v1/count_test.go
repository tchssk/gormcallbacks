package v1

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestCount(t *testing.T) {
	db, err := gorm.Open("sqlite3", "file::memory:?cache=shared")
	require.NoError(t, err)
	defer db.Close()
	db.DB().SetMaxOpenConns(1)

	type (
		Product struct{ gorm.Model }
		User    struct{ gorm.Model }
		Order   struct{ gorm.Model }
	)
	require.NoError(t, db.AutoMigrate(&Product{}, &User{}, &Order{}).Error)

	var (
		createCount = make(map[string]int)
		queryCount  = make(map[string]int)
		deleteCount = make(map[string]int)
		updateCount = make(map[string]int)
	)
	RegisterCountCallbacks(db, &CountConfig{createCount, queryCount, deleteCount, updateCount})

	var g = new(errgroup.Group)
	for i := 0; i < 12; i++ {
		g.Go(func() error {
			if i < 4 {
				var p Product
				if err := db.Create(&p).Error; err != nil {
					return err
				}
				if i > 0 {
					if err := db.First(&p, p.ID).Error; err != nil {
						return err
					}
				}
				if i > 1 {
					if err := db.Save(&p).Error; err != nil {
						return err
					}
				}
				if i > 2 {
					if err := db.Delete(&p).Error; err != nil {
						return err
					}
				}
			}
			if i < 8 {
				var u User
				if err := db.Create(&u).Error; err != nil {
					return err
				}
				if i > 1 {
					if err := db.First(&u, u.ID).Error; err != nil {
						return err
					}
				}
				if i > 3 {
					if err := db.Save(&u).Error; err != nil {
						return err
					}
				}
				if i > 5 {
					if err := db.Delete(&u).Error; err != nil {
						return err
					}
				}
			}
			{
				var o Order
				if err := db.Create(&o).Error; err != nil {
					return err
				}
				if i > 2 {
					if err := db.First(&o, o.ID).Error; err != nil {
						return err
					}
				}
				if i > 5 {
					if err := db.Save(&o).Error; err != nil {
						return err
					}
				}
				if i > 8 {
					if err := db.Delete(&o).Error; err != nil {
						return err
					}
				}
			}
			return nil
		})
	}
	require.NoError(t, g.Wait())

	assert.Equal(t, map[string]int{
		"products": 4,
		"users":    8,
		"orders":   12,
	}, createCount)
	assert.Equal(t, map[string]int{
		"products": 3,
		"users":    6,
		"orders":   9,
	}, queryCount)
	assert.Equal(t, map[string]int{
		"products": 2,
		"users":    4,
		"orders":   6,
	}, updateCount)
	assert.Equal(t, map[string]int{
		"products": 1,
		"users":    2,
		"orders":   3,
	}, deleteCount)
}
