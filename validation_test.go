package validations

import (
	"errors"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

type User struct {
	gorm.Model
	Name           string `valid:"required"`
	Password       string `valid:"length(6|20)"`
	SecurePassword string `valid:"numeric"`
	Email          string `valid:"email"`
	CompanyID      uint
	Company        Company
}

func (c *User) Validate() error {
	if c.Name == "invalid" {
		return errors.New("bad user name")
	}

	return nil
}

type Company struct {
	gorm.Model
	Name string
}

func (c *Company) Validate() error {
	if c.Name == "invalid" {
		return errors.New("bad company name")
	}

	return nil
}

func init() {
	db = MakeTestDB()
	RegisterCallbacks(db)
	db.AutoMigrate(&User{}, &Company{})
}

func TestGoValidation(t *testing.T) {
	type testCase struct {
		User          User
		ExpectedError string
	}

	for _, c := range []testCase{
		testCase{
			User:          User{Name: "", Password: "123123", Email: "a@gmail.com"},
			ExpectedError: "Name: non zero value required",
		},
		testCase{
			User:          User{Name: "zz", Password: "123123", Email: "zail.com"},
			ExpectedError: "Email: zail.com does not validate as email",
		},
		testCase{
			User:          User{Name: "invalid", Password: "123123", Email: "a@zail.com"},
			ExpectedError: "bad user name",
		},
		testCase{
			User:          User{Name: "zzz", Password: "123", Email: "a@zail.com"},
			ExpectedError: "Password: 123 does not validate as length(6|20)",
		},
		testCase{
			User:          User{Name: "valid", Password: "123123", Email: "a@zail.com", Company: Company{Name: "invalid"}},
			ExpectedError: "bad company name",
		},
		testCase{
			User:          User{Name: "valid", Password: "123123", Email: "a@zail.com", Company: Company{Name: "valid"}},
			ExpectedError: "",
		},
	} {
		result := db.Save(&c.User)
		if c.ExpectedError == "" {
			if result.Error != nil {
				t.Fatalf("expected nil error but got <%v>", result.Error.Error())
			}
		} else {
			if result.Error.Error() != c.ExpectedError {
				t.Fatalf("expected <%v> got <%v>", c.ExpectedError, result.Error.Error())
			}

		}

	}
}

func TestSaveInvalidUser(t *testing.T) {
	user := User{Name: "invalid"}

	if db.Save(&user).Error.Error() != "bad user name" {
		t.Errorf("missing bad user name")
	}

	user = User{Name: "valid", Company: Company{Name: "invalid"}}
	if db.Save(&user).Error.Error() != "bad company name" {
		t.Errorf("missing bad company name")
	}

	user = User{Name: "valid", Company: Company{Name: "valid"}}
	if result := db.Save(&user); result.Error != nil {
		t.Errorf("shouldnt get error, but got: %v", result.Error)
	}
}

func MakeTestDB() *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      false,
		},
	)
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{Logger: newLogger})

	if err != nil {
		panic(err)
	}

	return db
}
