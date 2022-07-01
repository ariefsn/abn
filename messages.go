package abn

type messages struct {
	GuidRequired     string
	AbnRequired      string
	AcnRequired      string
	NameRequired     string
	AbnInvalidLength string
	AbnInvalidType   string
	AbnInvalid       string
}

func NewMessages() *messages {
	m := messages{
		GuidRequired:     "guid is required",
		AbnRequired:      "abn is required",
		AcnRequired:      "acn is required",
		NameRequired:     "name is required",
		AbnInvalidLength: "abn must be 11 digits",
		AbnInvalidType:   "abn must be a number",
		AbnInvalid:       "invalid abn",
	}

	return &m
}
