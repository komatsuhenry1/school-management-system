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
	SubmitActivity(submission *model.ActivitySubmission) error
	GetSubmissionsByActivityID(activityID string) ([]model.ActivitySubmission, error)
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
	if err := r.db.Preload("Exercises.Alternatives").Find(&activities).Error; err != nil {
		return nil, err
	}
	return activities, nil
}

func (r *activityRepository) GetActivityByID(id string) (*model.Activity, error) {
	var activity model.Activity
	if err := r.db.Preload("Exercises.Alternatives").First(&activity, "id = ?", id).Error; err != nil {
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

func (r *activityRepository) SubmitActivity(submission *model.ActivitySubmission) error {
	return r.db.Create(submission).Error
}

func (r *activityRepository) GetSubmissionsByActivityID(activityID string) ([]model.ActivitySubmission, error) {
	var submissions []model.ActivitySubmission
	if err := r.db.Preload("Answers").Where("activity_id = ?", activityID).Find(&submissions).Error; err != nil {
		return nil, err
	}
	return submissions, nil
}
