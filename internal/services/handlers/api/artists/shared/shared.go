package shared

type Albumn struct {
	ID   string
	Name string
	Year int
}
type Artist struct {
	Id     string
	Name   string
	Albums []Albumn
}

var Artists []Artist

func init() {
	Artists = []Artist{
		{
			Id:   "1",
			Name: "Metallica",
			Albums: []Albumn{
				{ID: "2", Name: "Ride the Lightning", Year: 1984},
				{ID: "1", Name: "Master of Puppets", Year: 1986},
				{ID: "3", Name: "Metallica", Year: 1988},
			},
		},
		{
			Id:   "2",
			Name: "Iron Maiden",
			Albums: []Albumn{
				{ID: "4", Name: "Ace of Spades", Year: 1984},
				{ID: "5", Name: "Fear of the Dark", Year: 1985},
				{ID: "6", Name: "The Number of the Beast", Year: 1986},
			},
		},
		{
			Id:   "3",
			Name: "AC/DC",
			Albums: []Albumn{
				{ID: "7", Name: "Back in Black", Year: 1980},
				{ID: "8", Name: "Highway to Hell", Year: 1983},
				{ID: "9", Name: "High Voltage", Year: 1984},
			},
		},
	}
}
