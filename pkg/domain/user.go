package domain

type User struct {
	Id       int    `json:"-" db:"id"`
	Name     string `json:"name"     binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CourseFilter struct {
	IsFree    bool
	PriceFrom int64 // от 1000_00 коп.
	PriceTo   int64 //
}

type Course struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	ImageURL string `json:"imageUrl"`
}
