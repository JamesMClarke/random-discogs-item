package models

type Record struct {
	ID               int              `json:"id"`
	InstanceID       int              `json:"instance_id"`
	FolderID         int              `json:"folder_id"`
	Rating           int              `json:"rating"`
	BasicInformation BasicInformation `json:"basic_information"`
	Notes            []Note           `json:"notes"`
	FolderName       string           // Not from JSON, set manually
}

type BasicInformation struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Year        int      `json:"year"`
	ResourceURL string   `json:"resource_url"`
	Thumb       string   `json:"thumb"`
	CoverImage  string   `json:"cover_image"`
	Formats     []Format `json:"formats"`
	Labels      []Label  `json:"labels"`
	Artists     []Artist `json:"artists"`
	Genres      []string `json:"genres"`
	Styles      []string `json:"styles"`
}

type Format struct {
	Qty          string   `json:"qty"`
	Descriptions []string `json:"descriptions"`
	Name         string   `json:"name"`
}

type Label struct {
	ResourceURL string `json:"resource_url"`
	EntityType  string `json:"entity_type"`
	CatNo       string `json:"catno"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
}

type Artist struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Join        string `json:"join"`
	ResourceURL string `json:"resource_url"`
	Anv         string `json:"anv"`
	Tracks      string `json:"tracks"`
	Role        string `json:"role"`
}

type Note struct {
	FieldID int    `json:"field_id"`
	Value   string `json:"value"`
}
