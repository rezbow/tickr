package users

import (
	"context"

	"github.com/google/uuid"
	"github.com/rezbow/tickr/internal/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (service *UsersService) createUser(ctx context.Context, user *entities.User) error {
	user.ID = uuid.New()
	err := gorm.G[entities.User](service.db).Create(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

// deleteUser deletes a user from the database.
func (service *UsersService) deleteUser(ctx context.Context, userID uuid.UUID) error {
	rowsAffected, err := gorm.G[entities.User](service.db).Where("id = ?", userID).Delete(ctx)
	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	} else if err != nil {
		return err
	}
	return nil
}

// getUser retrieves a user from the database.
func (service *UsersService) getUser(ctx context.Context, userID uuid.UUID) (*entities.User, error) {
	user, err := gorm.G[entities.User](service.db).Where("id = ?", userID).First(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// getUsers: retrieves a list of users from the database.
func (service *UsersService) getUsers(ctx context.Context, page, limit int) ([]entities.User, int64, error) {
	var total int64
	if res := service.db.Model(&entities.User{}).Count(&total); res.Error != nil {
		return nil, 0, res.Error
	}
	users, err := gorm.G[entities.User](service.db).Offset((page - 1) * limit).Limit(limit).Find(ctx)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// getUsersCursorPagination: retrieves a list of users with cursor pagination
func (service *UsersService) getUsersCursorPagination(ctx context.Context, cursor string, limit int) ([]entities.User, string, error) {
	var users []entities.User
	var nextCursor string

	if cursor != "" {
		cursorID, err := uuid.Parse(cursor)
		if err != nil {
			return nil, "", err
		}
		users, err = gorm.G[entities.User](service.db).Where("id > ?", cursorID).Limit(limit).Order("id ASC").Find(ctx)
		if err != nil {
			return nil, "", err
		}
		if len(users) > 0 {
			nextCursor = users[len(users)-1].ID.String()
		}
	} else {
		users, err := gorm.G[entities.User](service.db).Limit(limit).Find(ctx)
		if err != nil {
			return nil, "", err
		}
		if len(users) > 0 {
			nextCursor = users[len(users)-1].ID.String()
		}
	}

	return users, nextCursor, nil
}

// updateUser updates a user in the database.
func (service *UsersService) updateUser(ctx context.Context, userID uuid.UUID, user entities.User) error {
	rowsAffected, err := gorm.G[entities.User](service.db).Where("id = ?", userID).Updates(ctx, user)
	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	} else if err != nil {
		return err
	}
	return nil
}

// updateUserAtomic updates a user in the database atomically.
func (service *UsersService) updateUserAtomic(userID uuid.UUID, updatedUser map[string]any) (*entities.User, error) {
	tx := service.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // re-throw panic
		}
	}()
	var user entities.User
	if res := service.db.Model(&entities.User{}).Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, userID); res.Error != nil {
		tx.Rollback()
		return nil, res.Error
	}
	if res := service.db.Model(user).Updates(updatedUser); res.Error != nil {
		tx.Rollback()
		return nil, res.Error
	}

	if err := service.db.Model(&entities.User{}).First(&user, userID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return &user, nil
}
