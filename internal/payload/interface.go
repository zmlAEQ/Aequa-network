package payload

type Payload interface { Kind() string }

// Manager placeholder decouples consensus from payload specifics.
