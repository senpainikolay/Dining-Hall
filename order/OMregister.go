package order

type RegisterPayload struct {
	RestaurantId int     `json:"restaurant_id"`
	Address      string  `json:"address"`
	Name         string  `json:"name"`
	MenuItems    int     `json:"menu_items"`
	Menu         []Food  `json:"menu"`
	Rating       float64 `json:"rating"`
}
