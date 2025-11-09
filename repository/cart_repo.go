package repository

import (
	"gorm.io/gorm"
	
	"ecommerce-api/domain"
)


func (r *PostgresRepository) FindByUserID(userID uint) (*domain.Cart, error) {
	var Cart domain.Cart
	err := r.DB.Where("user_id = ?", userID).Preload("Items").First(&Cart).Error
	if err == gorm.ErrRecordNotFound {
		Cart = domain.Cart{UserID: userID}
		r.DB.Create(&Cart)
		return &Cart, nil
	}
	return &Cart, err
}

func (r *PostgresRepository) Save(Cart *domain.Cart) error {
	return r.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(Cart).Error
}

func (r *PostgresRepository) Clear(userID uint) error {
	if err := r.DB.Where("user_id = ?", userID).Delete(&domain.CartItem{}).Error; err != nil {
		return err
	}
	return r.DB.Where("user_id = ?", userID).Delete(&domain.Cart{}).Error
}

func (r *PostgresRepository) UpdateInventory(productID uint, quantityChange int) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&domain.Product{}).Where("id = ? AND inventory >= ?", productID, -quantityChange).
			UpdateColumn("inventory", gorm.Expr("inventory + ?", quantityChange))

		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return domain.ErrInsufficientInv // Fails if inventory would go below zero
		}
		return nil
	})
}