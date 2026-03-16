package repository

import (
	"schoolmanagement/internal/activity/model"

	"gorm.io/gorm"
)

type ActivityRepository interface {
	CreateActivity(activity *model.Activity) error
	GetAllActivities() ([]model.Activity, error)
	GetActivityByID(id string) (*model.Activity, error)
	UpdateActivity(id string, updates map[string]interface{}) (*model.Activity, error)
	DeleteActivity(id string) error
}

type activityRepository struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return &activityRepository{db: db}
}

func (r *activityRepository) CreateActivity(activity *model.Activity) error {
	return r.db.Create(activity).Error
}

func (r *activityRepository) GetAllActivities() ([]model.Activity, error) {
	var activities []model.Activity
	if err := r.db.Find(&activities).Error; err != nil {
		return nil, err
	}
	return activities, nil
}

func (r *activityRepository) GetActivityByID(id string) (*model.Activity, error) {
	var activity model.Activity
	if err := r.db.First(&activity, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &activity, nil
}

func (r *activityRepository) UpdateActivity(id string, updates map[string]interface{}) (*model.Activity, error) {
	var activity model.Activity
	if err := r.db.First(&activity, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&activity).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &activity, nil
}

func (r *activityRepository) DeleteActivity(id string) error {
	return r.db.Delete(&model.Activity{}, "id = ?", id).Error
}
