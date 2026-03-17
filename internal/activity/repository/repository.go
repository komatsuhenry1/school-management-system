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
	GetActiveActivities() ([]model.Activity, error)
	GetSubmissionsByUserID(userID string) ([]model.ActivitySubmission, error)
	HasUserSubmittedActivity(userID string, activityID string) (bool, error)
	UpdateAlternative(alternativeID string, updates map[string]interface{}) error
	UpdateActivityFull(activity *model.Activity) error
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

func (r *activityRepository) GetActiveActivities() ([]model.Activity, error) {
	var activities []model.Activity
	// Only fetch exercises, we don't need alternatives here based on requirement
	if err := r.db.Preload("Exercises").Where("status = ?", "ACTIVE").Find(&activities).Error; err != nil {
		return nil, err
	}
	return activities, nil
}

func (r *activityRepository) GetSubmissionsByUserID(userID string) ([]model.ActivitySubmission, error) {
	var submissions []model.ActivitySubmission
	if err := r.db.Preload("Answers").Where("user_id = ?", userID).Find(&submissions).Error; err != nil {
		return nil, err
	}
	return submissions, nil
}

func (r *activityRepository) HasUserSubmittedActivity(userID string, activityID string) (bool, error) {
	var count int64
	err := r.db.Model(&model.ActivitySubmission{}).
		Where("user_id = ? AND activity_id = ?", userID, activityID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *activityRepository) UpdateAlternative(alternativeID string, updates map[string]interface{}) error {
	return r.db.Model(&model.Alternative{}).Where("id = ?", alternativeID).Updates(updates).Error
}

func (r *activityRepository) UpdateActivityFull(activity *model.Activity) error {
	// Use GORM's Full Save Associations to replace everything
	// Alternatively, we delete existing exercises and recreate them
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Delete existing Exercises (and cascade Alternatives if relying on DB constraints, 
		// otherwise GORM Unscoped to forcefully clean up before insert)
		if err := tx.Where("activity_id = ?", activity.ID).Delete(&model.Exercise{}).Error; err != nil {
			return err
		}

		// 2. Update the parent Activity record
		if err := tx.Model(activity).Select("Title", "Description", "ActivityValue").Updates(activity).Error; err != nil {
			return err
		}

		// 3. Insert the new Exercises correctly associated
		for i := range activity.Exercises {
			activity.Exercises[i].ActivityID = activity.ID
			if err := tx.Create(&activity.Exercises[i]).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
