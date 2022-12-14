package handler

type Category struct {
	Name        string
}

var (
	UtilitiesCategory = &Category{
		Name:        "Utilities",
	}
	GroupCategory = &Category{
		Name:        "Group",
	}
	MiscCategory = &Category{
		Name:        "Misc",
	}
	Uncategorized = &Category{
		Name:        "Uncategorized",
	}
)
