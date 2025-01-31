package cli

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func main() {
	// Configure the bech32 prefix
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("cosmos", "cosmospub")

	// Generate some test addresses
	addr1 := sdk.AccAddress([]byte("test1_______________"))
	addr2 := sdk.AccAddress([]byte("test2_______________"))
	addr3 := sdk.AccAddress([]byte("test3_______________"))

	fmt.Printf("Test Address 1: %s\n", addr1.String())
	fmt.Printf("Test Address 2: %s\n", addr2.String())
	fmt.Printf("Test Address 3: %s\n", addr3.String())
}
