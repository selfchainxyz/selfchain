package types

type AddBeneficiaryRequest struct {
	Beneficiary string
	Cliff       uint64
	Duration    uint64
	Amount      string
}
