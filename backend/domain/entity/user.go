package entity

import (
	"time"

	"caltrack/domain/vo"
)

type User struct {
	id             vo.UserID
	email          vo.Email
	hashedPassword vo.HashedPassword
	nickname       vo.Nickname
	weight         vo.Weight
	height         vo.Height
	birthDate      vo.BirthDate
	gender         vo.Gender
	activityLevel  vo.ActivityLevel
	createdAt      time.Time
	updatedAt      time.Time
}

func NewUser(
	emailStr string,
	passwordStr string,
	nicknameStr string,
	weightVal float64,
	heightVal float64,
	birthDateVal time.Time,
	genderStr string,
	activityLevelStr string,
) (*User, []error) {
	var errs []error

	email, err := vo.NewEmail(emailStr)
	errs = appendIfErr(errs, err)

	password, err := vo.NewPassword(passwordStr)
	errs = appendIfErr(errs, err)

	var hashedPassword vo.HashedPassword
	if err == nil {
		hashedPassword, err = password.Hash()
		if err != nil {
			errs = append(errs, err)
		}
	}

	nickname, err := vo.NewNickname(nicknameStr)
	errs = appendIfErr(errs, err)

	weight, err := vo.NewWeight(weightVal)
	errs = appendIfErr(errs, err)

	height, err := vo.NewHeight(heightVal)
	errs = appendIfErr(errs, err)

	birthDate, err := vo.NewBirthDate(birthDateVal)
	errs = appendIfErr(errs, err)

	gender, err := vo.NewGender(genderStr)
	errs = appendIfErr(errs, err)

	activityLevel, err := vo.NewActivityLevel(activityLevelStr)
	errs = appendIfErr(errs, err)

	if len(errs) > 0 {
		return nil, errs
	}

	now := time.Now()
	return &User{
		id:             vo.NewUserID(),
		email:          email,
		hashedPassword: hashedPassword,
		nickname:       nickname,
		weight:         weight,
		height:         height,
		birthDate:      birthDate,
		gender:         gender,
		activityLevel:  activityLevel,
		createdAt:      now,
		updatedAt:      now,
	}, nil
}

// ReconstructUser はDBからの復元用。
// NOTE: DBデータは保存時にバリデーション済みなのでVO変換は本来不要だが、
// データ破損検知のため一旦バリデーションありで実装。
// パフォーマンス問題が出たらReconstruct系関数（バリデーションなし）を検討。
func ReconstructUser(
	idStr string,
	emailStr string,
	hashedPasswordStr string,
	nicknameStr string,
	weightVal float64,
	heightVal float64,
	birthDateVal time.Time,
	genderStr string,
	activityLevelStr string,
	createdAt time.Time,
	updatedAt time.Time,
) (*User, error) {
	id := vo.ReconstructUserID(idStr)

	email, err := vo.NewEmail(emailStr)
	if err != nil {
		return nil, err
	}

	hashedPassword := vo.NewHashedPassword(hashedPasswordStr)

	nickname, err := vo.NewNickname(nicknameStr)
	if err != nil {
		return nil, err
	}

	weight, err := vo.NewWeight(weightVal)
	if err != nil {
		return nil, err
	}

	height, err := vo.NewHeight(heightVal)
	if err != nil {
		return nil, err
	}

	birthDate, err := vo.NewBirthDate(birthDateVal)
	if err != nil {
		return nil, err
	}

	gender, err := vo.NewGender(genderStr)
	if err != nil {
		return nil, err
	}

	activityLevel, err := vo.NewActivityLevel(activityLevelStr)
	if err != nil {
		return nil, err
	}

	return &User{
		id:             id,
		email:          email,
		hashedPassword: hashedPassword,
		nickname:       nickname,
		weight:         weight,
		height:         height,
		birthDate:      birthDate,
		gender:         gender,
		activityLevel:  activityLevel,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}, nil
}

func (u *User) ID() vo.UserID {
	return u.id
}

func (u *User) Email() vo.Email {
	return u.email
}

func (u *User) HashedPassword() vo.HashedPassword {
	return u.hashedPassword
}

func (u *User) Nickname() vo.Nickname {
	return u.nickname
}

func (u *User) Weight() vo.Weight {
	return u.weight
}

func (u *User) Height() vo.Height {
	return u.height
}

func (u *User) BirthDate() vo.BirthDate {
	return u.birthDate
}

func (u *User) Gender() vo.Gender {
	return u.gender
}

func (u *User) ActivityLevel() vo.ActivityLevel {
	return u.activityLevel
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}
