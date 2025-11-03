package decimal

func (d Decimal) MarshalJSON() ([]byte, error) {
	return d.bytes(), nil
}

func (d *Decimal) UnmarshalJSON(data []byte) error {
	if len(data) >= 2 && data[0] == '"' && data[len(data)-1] == '"' {
		data = data[1 : len(data)-1]
	}
	return d.FromString(string(data))
}
