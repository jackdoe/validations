gorm validation callback that has Validate() interface and https://github.com/asaskevich/govalidator support

Example:

```

type User struct {
	gorm.Model
	Name           string `valid:"required"`
	Password       string `valid:"length(6|20)"`
	SecurePassword string `valid:"numeric"`
	Email          string `valid:"email"`
	CompanyID      uint
	Company        Company
}

func (u *User) Validate() error {
	if u.Name == "invalid" {
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


func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	validations.RegisterCallbacks(db)
}
```