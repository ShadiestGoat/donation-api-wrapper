package donations

import "net/http"

type Donor struct {
	// The donor's internal donor id
	ID string `json:"id"`
	// The donor's discord id
	DiscordID string `json:"discordID"`
	// The donor's paypal idi
	PayPal string `json:"PayPal"`
	// The day of the month their pay cycle expires on.
	CycleDay int `json:"payCycle"`
}

// A profile refers to the overall profile of a certain donor.
// You can think of this as the aggregate data based on a query.
type DonorProfile struct {
	// These are all the donor profiles that matched with a query.
	// If the query is based on a donor ID, this would respond with an array of length 1 (if there is a donor under such ID)
	Donors []*Donor `json:"donors"`
	// The aggregate count by all the donors matches
	Total *AggregateDonations `json:"total"`
	// The donations made by the donors.
	// Please note that this can be nil!
	Donations *[]*Donation `json:"donations,omitempty"`
}

type AggregateDonations struct {
	// The amount that was donated in all time
	Total float64 `json:"total"`
	// The amount that was donated this pay cycle
	Month float64 `json:"month"`
}

func (c *Client) fetchDonor(id string, idType string, resolve bool) (*DonorProfile, error) {
	resp := &DonorProfile{}
	q := ""
	if resolve {
		q = `?resolve=true`
	}

	err := c.fetch(http.MethodGet, `/donors/`+idType+`/`+id+q, nil, &resp)

	return withError(resp, err)
}

// Fetch a donor by their discord id.
// If resolve = true, it will return information about specific donations
func (c *Client) DonorByDiscord(discordID string, resolve bool) (*DonorProfile, error) {
	return c.fetchDonor(discordID, `discord`, resolve)
}

// Fetch a donor by their donor id.
// If resolve = true, it will return information about specific donations
func (c *Client) DonorByID(id string, resolve bool) (*DonorProfile, error) {
	return c.fetchDonor(id, `donor`, resolve)
}

// Fetch a donor by their paypal id.
// If resolve = true, it will return information about specific donations
func (c *Client) DonorByPayPalID(paypalID string, resolve bool) (*DonorProfile, error) {
	return c.fetchDonor(paypalID, `paypal`, resolve)
}
