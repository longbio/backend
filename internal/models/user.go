package models

import (
	"database/sql/driver"
	"errors"
	"time"

	"gorm.io/gorm"
)

// Gender enum
type Gender string

const (
	GenderMale      Gender = "male"
	GenderFemale    Gender = "female"
	GenderNonBinary Gender = "nonbinary"
)

func (g Gender) IsValid() bool {
	return g == GenderMale || g == GenderFemale || g == GenderNonBinary
}

func (g *Gender) Scan(value interface{}) error {
	if value == nil {
		*g = ""
		return nil
	}
	if val, ok := value.(string); ok {
		*g = Gender(val)
		return nil
	}
	return errors.New("cannot scan value into Gender")
}

func (g Gender) Value() (driver.Value, error) {
	return string(g), nil
}

// MaritalStatus enum
type MaritalStatus string

const (
	MaritalSingle   MaritalStatus = "single"
	MaritalMarried  MaritalStatus = "married"
	MaritalDivorced MaritalStatus = "divorced"
	MaritalWidowed  MaritalStatus = "widowed"
)

func (m MaritalStatus) IsValid() bool {
	return m == MaritalSingle || m == MaritalMarried || m == MaritalDivorced || m == MaritalWidowed
}

func (m *MaritalStatus) Scan(value interface{}) error {
	if value == nil {
		*m = ""
		return nil
	}
	if val, ok := value.(string); ok {
		*m = MaritalStatus(val)
		return nil
	}
	return errors.New("cannot scan value into MaritalStatus")
}

func (m MaritalStatus) Value() (driver.Value, error) {
	return string(m), nil
}

// EducationalStatus enum
type EducationalStatus string

const (
	EduHighSchool EducationalStatus = "high_school"
	EduBachelor   EducationalStatus = "bachelor"
	EduMaster     EducationalStatus = "master"
	EduPhD        EducationalStatus = "phd"
	EduDiploma    EducationalStatus = "diploma"
	EduNone       EducationalStatus = "none"
)

func (e EducationalStatus) IsValid() bool {
	return e == EduHighSchool || e == EduBachelor || e == EduMaster || e == EduPhD || e == EduDiploma || e == EduNone
}

func (e *EducationalStatus) Scan(value interface{}) error {
	if value == nil {
		*e = ""
		return nil
	}
	if val, ok := value.(string); ok {
		*e = EducationalStatus(val)
		return nil
	}
	return errors.New("cannot scan value into EducationalStatus")
}

func (e EducationalStatus) Value() (driver.Value, error) {
	return string(e), nil
}

// Skill enum
type Skill string

const (
	SkillProgramming Skill = "programming"
	SkillDesign      Skill = "design"
	SkillWriting     Skill = "writing"
	SkillMarketing   Skill = "marketing"
	SkillManagement  Skill = "management"
	SkillTeaching    Skill = "teaching"
	SkillCooking     Skill = "cooking"
	SkillPhotography Skill = "photography"
	SkillMusic       Skill = "music"
	SkillLanguages   Skill = "languages"
)

type Skills []Skill

func (s *Skills) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	// Implementation depends on how you store array in DB (JSON, comma-separated, etc.)
	return nil
}

func (s Skills) Value() (driver.Value, error) {
	// Implementation depends on how you want to store array in DB
	return nil, nil
}

// Interest enum
type Interest string

const (
	InterestSports     Interest = "sports"
	InterestTechnology Interest = "technology"
	InterestArts       Interest = "arts"
	InterestMusic      Interest = "music"
	InterestMovies     Interest = "movies"
	InterestBooks      Interest = "books"
	InterestTravel     Interest = "travel"
	InterestCooking    Interest = "cooking"
	InterestGaming     Interest = "gaming"
	InterestNature     Interest = "nature"
)

type Interests []Interest

func (i *Interests) Scan(value interface{}) error {
	if value == nil {
		*i = nil
		return nil
	}
	// Implementation depends on how you store array in DB
	return nil
}

func (i Interests) Value() (driver.Value, error) {
	// Implementation depends on how you want to store array in DB
	return nil, nil
}

type User struct {
	ID                uint              `gorm:"primaryKey" json:"id"`
	BirthDate         *time.Time        `gorm:"type:date" json:"birthDate"`
	Email             string            `gorm:"varchar(150);not null" json:"email"`
	FullName          string            `gorm:"varchar(150)" json:"fullName"`
	Gender            Gender            `gorm:"varchar(20);check:gender IN ('male','female','non-binary')" json:"gender"`
	MaritalStatus     MaritalStatus     `gorm:"varchar(20);check:marital_status IN ('single','married','divorced','widowed')" json:"maritalStatus"`
	EducationalStatus EducationalStatus `gorm:"varchar(20);check:educational_status IN ('high_school','bachelor','master','phd','diploma','none')" json:"educationalStatus"`
	ProfilePic        string            `gorm:"varchar(255)" json:"profilePic"`
	Height            float64           `gorm:"type:decimal(5,2)" json:"height"` // in cm
	Weight            float64           `gorm:"type:decimal(5,2)" json:"weight"` // in kg
	BornPlace         string            `gorm:"varchar(100)" json:"bornPlace"`
	LivePlace         string            `gorm:"varchar(100)" json:"livePlace"`
	HasPet            bool              `gorm:"default:false" json:"hasPet"`
	DoesExercise      bool              `gorm:"default:false" json:"doesExercise"`
	Skills            Skills            `gorm:"type:json" json:"skills"`
	Interests         Interests         `gorm:"type:json" json:"interests"`
	Details           string            `gorm:"varchar(250)" json:"details"`
	CreatedAt         time.Time         `json:"createdAt"`
	UpdatedAt         time.Time         `json:"updatedAt"`
}

// BeforeSave validates enum fields before saving
func (u *User) BeforeSave(tx *gorm.DB) error {
	if u.Gender != "" && !u.Gender.IsValid() {
		return errors.New("invalid gender value")
	}
	if u.MaritalStatus != "" && !u.MaritalStatus.IsValid() {
		return errors.New("invalid marital status value")
	}
	if u.EducationalStatus != "" && !u.EducationalStatus.IsValid() {
		return errors.New("invalid educational status value")
	}
	return nil
}
