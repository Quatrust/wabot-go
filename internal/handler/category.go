package handler

type Category struct {
	Name        string
}

var (
	UtilitiesCategory = &Category{
		Name:        "Utilities",
	}
	MiscCategory = &Category{
		Name:        "Misc",
	}
	Uncategorized = &Category{
		Name:        "Uncategorized",
	}
)
