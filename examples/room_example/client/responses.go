package client

type AddTODOResponse struct {
	Id        int    `json:"id"`
	Todo      string `json:"todo"`
	Completed bool   `json:"completed"`
	UserId    int    `json:"userId"`
}
