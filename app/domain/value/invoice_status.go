package value

import "fmt"

type InvoiceStatus string

const (
	InvoiceStatusUnprocessed InvoiceStatus = "未処理"
	InvoiceStatusProcessing  InvoiceStatus = "処理中"
	InvoiceStatusError       InvoiceStatus = "エラー"
	InvoiceStatusProcessed   InvoiceStatus = "処理済"
)

func (s *InvoiceStatus) String() string {
	return string(*s)
}

func (s *InvoiceStatus) UnmarshalJSON(b []byte) error {
	*s = InvoiceStatus(b)
	return nil
}

func (s *InvoiceStatus) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, s.String())), nil
}
