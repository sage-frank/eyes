package domain

type (
	Article struct {
		ID      int64
		Author  int64
		Title   string
		Content string
	}

	DA struct {
		ID int64 `gorm:"primaryKey,autoIncrement"`
		// 存的是作者的用户 ID
		Author  int64  `gorm:"not null"`
		Title   string `form:"title"`
		Content string `form:"content"`
		Ctime   int64
		Utime   int64
	}
)
