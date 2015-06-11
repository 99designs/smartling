package smartling

type Locale struct {
	Name       string `json:"name"`
	Locale     string `json:"locale"`
	Translated string `json:"translated"`
}

type LocalesResponse struct {
	Locales []Locale `json:"locales"`
}

func (c *Client) Locales() ([]Locale, error) {
	r := LocalesResponse{}
	err := c.doRequestAndUnmarshalData("/project/locale/list", nil, &r)

	return r.Locales, err
}
