package test
//package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type Period struct {
	Coins         string `json:"coins"`
	LengthSeconds int64  `json:"length_seconds"`
}

type VestingSchedule struct {
	StartTime int64    `json:"start_time"`
	Periods   []Period `json:"periods"`
}

type AccountInfo struct {
	Address         string
	VestingSchedule VestingSchedule
	IsReplacement   bool
	NewAddress      string // Only for replacement addresses
	PostponePeriods bool   // True if this account needs 3-month postponement
}

type VestingAccountResponse struct {
	Type string `json:"@type"`
	BaseVestingAccount struct {
		BaseAccount struct {
			AccountNumber string      `json:"account_number"`
			Address      string      `json:"address"`
			PubKey       interface{} `json:"pub_key"`
			Sequence     string      `json:"sequence"`
		} `json:"base_account"`
		DelegatedFree    []interface{} `json:"delegated_free"`
		DelegatedVesting []interface{} `json:"delegated_vesting"`
		EndTime          string `json:"end_time"`
		OriginalVesting  []Coin `json:"original_vesting"`
	} `json:"base_vesting_account"`
	StartTime      string `json:"start_time"`
	VestingPeriods []struct {
		Amount []Coin `json:"amount"`
		Length string `json:"length"`
	} `json:"vesting_periods"`
}

type BalanceResponse struct {
	Balances []Coin `json:"balances"`
	Pagination struct {
		NextKey interface{} `json:"next_key"`
		Total   string      `json:"total"`
	} `json:"pagination"`
}

type Coin struct {
	Amount string `json:"amount"`
	Denom  string `json:"denom"`
}

type AccountState struct {
	VestingAccount   *VestingAccountResponse
	SpendableBalance *BalanceResponse
	TotalBalance     *BalanceResponse
}

func main() {
	/*
	 replace in handler
	var addressReplacements = map[string]string{
		"self102fgcqwkhcrwf6yv8jgen7v2gd0k4e0szpfh3d": "self1scmpmsrv74r47fhj2fzcgeuque6pudam59prw8",
		"self1fcahhgtw2llk06am4rala6khxjtj24zhhxn449": "self1kr30hqm2ezdjapspemdjgrt5lkxhsmwwr6ujtr",
		// Add more address mappings as needed
	}

	// Enter the list of address for which the vesting schedule needs to be postponed by 3 months.
	var vestingAddresses = []string{
		"self1krxfd67wmrjksq20xww53rm0wqmyxcew22whah",
		"self1fcahhgtw2llk06am4rala6khxjtj24zhhxn449",
	}

	 selfchaind tendermint unsafe-reset-all

	 cp oldselfchaind /Users/shivakumarhg/go/bin/selfchaind

	 selfchaind start

	 selfchaind tx gov submit-proposal proposal2.json \
	 --node http://0.0.0.0:26657 \
	    --gas="auto" \
	    --gas-prices="0.5uslf" \
	    --gas-adjustment="1.5" \
	    --chain-id selfchain \
	    --broadcast-mode sync \
	    --from alice

	 selfchaind tx gov vote 1 yes --from alice --chain-id selfchain --fees 1000000uslf -y


	cat proposal2.json
	{
	 "messages": [
	  {
	   "@type": "/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
	   "authority": "self10d07y265gmmuvt4z0w9aw880jnsr700jlfwec6",
	   "plan": {
	    "name": "v2",
	    "height": "80",
	    "info": "",
	    "upgraded_client_state": null
	   }
	  }
	 ],
	 "metadata": "ipfs://CID",
	 "deposit": "10000000uslf",
	 "title": "Add CosmWasm Support",
	 "summary": "Adding CosmWasm module support",
	 "voting_period": "60s"
	}


	 */
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	currentTime := time.Now().Unix()

	// Test accounts with different vesting schedules
	testAccounts := []AccountInfo{
		{
			// Address replacement case 1
			Address: "self102fgcqwkhcrwf6yv8jgen7v2gd0k4e0szpfh3d",
			VestingSchedule: VestingSchedule{
				StartTime: currentTime - 120, // Started 2 minutes ago
				Periods: generatePeriods(16, map[string]interface{}{
					"coins":          "1000000000uslf",
					"length_seconds": int64(60),
				}),
			},
			IsReplacement: true,
			NewAddress:    "self1scmpmsrv74r47fhj2fzcgeuque6pudam59prw8",
		},
		{
			// Address replacement + postponement case
			Address: "self1fcahhgtw2llk06am4rala6khxjtj24zhhxn449",
			VestingSchedule: VestingSchedule{
				StartTime: currentTime - 4000,
				Periods: []Period{
					{Coins: "1000000000uslf", LengthSeconds: 3600},    // 1 hour
					{Coins: "2000000000uslf", LengthSeconds: 3600},    // 1 hour
					{Coins: "3000000000uslf", LengthSeconds: 3600},    // 1 hour
					{Coins: "4000000000uslf", LengthSeconds: 2628000}, // ~1 month
				},
			},
			IsReplacement:   true,
			NewAddress:      "self1kr30hqm2ezdjapspemdjgrt5lkxhsmwwr6ujtr",
			PostponePeriods: true,
		},
		{
			// Postponement only case
			Address: "self1krxfd67wmrjksq20xww53rm0wqmyxcew22whah",
			VestingSchedule: VestingSchedule{
				StartTime: currentTime - 7200, // Started 2 hours ago
				Periods: []Period{
					{Coins: "2000000000uslf", LengthSeconds: 3600},    // 1 hour
					{Coins: "3000000000uslf", LengthSeconds: 3600},    // 1 hour
					{Coins: "4000000000uslf", LengthSeconds: 2628000}, // ~1 month
				},
			},
			PostponePeriods: true,
		},
	}

	fmt.Println("=== Phase 1: Initial Setup and Vesting Creation ===")

	// Step 1: Check if chain is running
	if err := checkChainStatus(); err != nil {
		log.Fatalf("Chain status check failed: %v", err)
	}

	// Step 2: Create vesting accounts
	fmt.Println("\n=== Creating Vesting Accounts ===")
	for _, account := range testAccounts {
		if err := createVestingAccount(account); err != nil {
			log.Printf("Failed to create vesting account for %s: %v", account.Address, err)
			continue
		}
		fmt.Printf("Created vesting account for %s\n", account.Address)
		if account.IsReplacement {
			fmt.Printf("  - Will be replaced with: %s\n", account.NewAddress)
		}
		if account.PostponePeriods {
			fmt.Printf("  - Periods will be postponed by 3 months\n")
		}
	}

	fmt.Println("\n=== Phase 2: Pre-Upgrade State Recording ===")
	fmt.Println("Please create and submit the upgrade proposal now:")
	fmt.Println("1. Create and pass upgrade proposal:")
	fmt.Println("   selfchaind tx gov submit-proposal software-upgrade v2 --upgrade-height=<N> --from=alice --chain-id=selfchain")
	fmt.Println("2. Deposit and vote:")
	fmt.Println("   selfchaind tx gov deposit <proposal_id> 10000000uslf --from=alice --chain-id=selfchain")
	fmt.Println("   selfchaind tx gov vote <proposal_id> yes --from=alice --chain-id=selfchain")
	fmt.Println("\nWhen the chain reaches the upgrade height and stops, press Enter to record pre-upgrade states...")
	fmt.Scanln()

	// Record pre-upgrade states
	fmt.Println("\n=== Recording Pre-upgrade States ===")
	preUpgradeStates := make(map[string]AccountState)
	for _, account := range testAccounts {
		state, err := getAccountState(account.Address)
		if err != nil {
			log.Printf("Failed to get pre-upgrade state for %s: %v", account.Address, err)
			continue
		}
		preUpgradeStates[account.Address] = state
		fmt.Printf("\nPre-upgrade state for %s:\n", account.Address)
		printAccountState(account.Address, state)
	}

	fmt.Println("\n=== Phase 3: Post-Upgrade Verification ===")
	fmt.Println("Now please:")
	fmt.Println("1. Replace the binary with the new version")
	fmt.Println("2. Start the chain again")
	fmt.Println("\nAfter the chain has restarted with the new binary, press Enter to verify the changes...")
	fmt.Scanln()

	// Step 5: Verify post-upgrade states
	fmt.Println("\n=== Verifying Post-upgrade States ===")

	// Track all addresses to verify
	type AddressRole struct {
		isOldAddr bool
		isNewAddr bool
		oldAddr   string
		newAddr   string
	}
	addressRoles := make(map[string]AddressRole)

	// Build the address roles map
	for _, account := range testAccounts {
		// Add the main account
		role := AddressRole{}
		if account.IsReplacement {
			role.isOldAddr = true
			role.oldAddr = account.Address
			role.newAddr = account.NewAddress
		}
		addressRoles[account.Address] = role

		// Add the replacement address if it exists
		if account.IsReplacement {
			newRole := AddressRole{
				isNewAddr: true,
				oldAddr:   account.Address,
				newAddr:   account.NewAddress,
			}
			addressRoles[account.NewAddress] = newRole
		}
	}

	// Verify each address
	for addr, role := range addressRoles {
		fmt.Printf("\n=== Verifying address: %s ===\n", addr)

		var preState AccountState
		if role.isNewAddr {
			// For new addresses, we expect them not to exist pre-upgrade
			fmt.Println("New address - no pre-upgrade state expected")
		} else {
			// For existing addresses, get their pre-upgrade state
			preState = preUpgradeStates[addr]
		}

		postState, err := getAccountState(addr)
		if err != nil {
			log.Printf("Failed to get post-upgrade state for %s: %v", addr, err)
			continue
		}

		// Print states
		if !role.isNewAddr {
			fmt.Println("\nPre-upgrade state:")
			printAccountState(addr, preState)
		}
		fmt.Println("\nPost-upgrade state:")
		printAccountState(addr, postState)

		// Verify based on role
		if role.isOldAddr {
			oldAddrState, err := getAccountState(addr)
			if err != nil {
				log.Printf("Failed to get state for old address %s: %v", addr, err)
				continue
			}
			newAddrState, err := getAccountState(role.newAddr)
			if err != nil {
				log.Printf("Failed to get state for new address %s: %v", role.newAddr, err)
				continue
			}
			verifyReplacementTransfer(addr, role.newAddr, preState, oldAddrState, newAddrState)
		}

		// Check for postponement if this is not a new address
		if !role.isNewAddr {
			for _, account := range testAccounts {
				if account.PostponePeriods && account.Address == addr {
					verifyPostponement(account, preState, postState)
					break
				}
			}
		}
	}
}

func checkChainStatus() error {
	cmd := exec.Command("selfchaind", "status", "--node", "http://0.0.0.0:26657")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("chain is not running: %v\nOutput: %s", err, string(out))
	}
	return nil
}

func createVestingAccount(account AccountInfo) error {
	// Create temporary file for vesting schedule
	schedule, err := json.Marshal(account.VestingSchedule)
	if err != nil {
		return fmt.Errorf("failed to marshal vesting schedule: %v", err)
	}

	tmpfile, err := os.CreateTemp("", "vesting-*.json")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(schedule); err != nil {
		return fmt.Errorf("failed to write schedule to file: %v", err)
	}

	if err := tmpfile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %v", err)
	}

	// Execute create-vesting-account command with updated flags
	cmd := exec.Command("selfchaind", "tx", "vesting", "create-periodic-vesting-account",
		account.Address,
		tmpfile.Name(),
		"--node", "http://0.0.0.0:26657",
		"--gas", "auto",
		"--gas-prices", "0.5uslf",
		"--gas-adjustment", "1.5",
		"--chain-id", "selfchain",
		"--broadcast-mode", "sync",
		"--from", "alice",
		"--keyring-backend", "test",
		"--yes")

	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create vesting account: %v\nOutput: %s", err, string(out))
	}

	// Wait for transaction to be processed
	time.Sleep(6 * time.Second)
	return nil
}

func getAccountState(address string) (AccountState, error) {
	var state AccountState

	// Get account info
	cmd := exec.Command("selfchaind", "query", "account", address, "--output", "json", "--node", "http://0.0.0.0:26657")
	out, err := cmd.CombinedOutput()
	if err == nil {
		var vestingAcc VestingAccountResponse
		if err := json.Unmarshal(out, &vestingAcc); err == nil {
			state.VestingAccount = &vestingAcc
		}
	}

	// Get spendable balance
	cmd = exec.Command("selfchaind", "query", "bank", "spendable-balances", address, "--output", "json", "--node", "http://0.0.0.0:26657")
	out, err = cmd.CombinedOutput()
	if err != nil {
		return state, fmt.Errorf("failed to query spendable balance: %v", err)
	}
	var spendable BalanceResponse
	if err := json.Unmarshal(out, &spendable); err != nil {
		return state, fmt.Errorf("failed to parse spendable balance: %v", err)
	}
	state.SpendableBalance = &spendable

	// Get total balance
	cmd = exec.Command("selfchaind", "query", "bank", "balances", address, "--output", "json", "--node", "http://0.0.0.0:26657")
	out, err = cmd.CombinedOutput()
	if err != nil {
		return state, fmt.Errorf("failed to query total balance: %v", err)
	}
	var total BalanceResponse
	if err := json.Unmarshal(out, &total); err != nil {
		return state, fmt.Errorf("failed to parse total balance: %v", err)
	}
	state.TotalBalance = &total

	return state, nil
}

func printAccountState(address string, state AccountState) {
	fmt.Printf("\nState for %s:\n", address)
	fmt.Printf("=== Account Details ===\n")

	if state.VestingAccount != nil {
		fmt.Printf("Account Type: %s\n", state.VestingAccount.Type)
		fmt.Printf("Start Time: %s\n", state.VestingAccount.StartTime)
		fmt.Printf("End Time: %s\n", state.VestingAccount.BaseVestingAccount.EndTime)
		fmt.Printf("Original Vesting:\n")
		for _, coin := range state.VestingAccount.BaseVestingAccount.OriginalVesting {
			fmt.Printf("  - %s %s\n", coin.Amount, coin.Denom)
		}
		fmt.Printf("Vesting Periods:\n")
		for i, period := range state.VestingAccount.VestingPeriods {
			fmt.Printf("  Period %d:\n", i+1)
			for _, coin := range period.Amount {
				fmt.Printf("    - Amount: %s %s\n", coin.Amount, coin.Denom)
			}
			fmt.Printf("    - Length: %s seconds\n", period.Length)
		}
	} else {
		fmt.Printf("No vesting account details found\n")
	}

	fmt.Printf("\n=== Balances ===\n")
	fmt.Printf("Total Balance:\n")
	if state.TotalBalance != nil && len(state.TotalBalance.Balances) > 0 {
		for _, coin := range state.TotalBalance.Balances {
			fmt.Printf("  - %s %s\n", coin.Amount, coin.Denom)
		}
	} else {
		fmt.Printf("  No balance found\n")
	}

	fmt.Printf("Spendable Balance:\n")
	if state.SpendableBalance != nil && len(state.SpendableBalance.Balances) > 0 {
		for _, coin := range state.SpendableBalance.Balances {
			fmt.Printf("  - %s %s\n", coin.Amount, coin.Denom)
		}
	} else {
		fmt.Printf("  No spendable balance found\n")
	}
}

func generatePeriods(count int, template map[string]interface{}) []Period {
	periods := make([]Period, count)
	for i := 0; i < count; i++ {
		periods[i] = Period{
			Coins:         template["coins"].(string),
			LengthSeconds: template["length_seconds"].(int64),
		}
	}
	return periods
}

func verifyPostponement(account AccountInfo, pre, post AccountState) {
	fmt.Printf("\nVerifying postponement for %s:\n", account.Address)
	currentTime := time.Now().Unix()

	if account.IsReplacement {
		fmt.Println("Skipping postponement check because this address also has replacement logic.")
		return
	}

	// Check total balance preservation
	preTotal := sumCoins(pre.TotalBalance.Balances)
	postTotal := sumCoins(post.TotalBalance.Balances)
	if preTotal != postTotal {
		fmt.Printf("Warn- Total balance changed: %d -> %d\n", preTotal, postTotal)
	} else {
		fmt.Printf("✅ Total balance unchanged: %d\n", preTotal)
	}

	if post.VestingAccount == nil {
		fmt.Printf("❌ Vesting account structure lost\n")
		return
	}

	// Check each period
	preStart, _ := strconv.ParseInt(pre.VestingAccount.StartTime, 10, 64)
	cumulativeTime := preStart
	threeMonths := int64(7884000)

	for i, period := range pre.VestingAccount.VestingPeriods {
		preLength, _ := strconv.ParseInt(period.Length, 10, 64)
		periodEndTime := cumulativeTime + preLength

		if i >= len(post.VestingAccount.VestingPeriods) {
			if periodEndTime > currentTime {
				fmt.Printf("❌ Missing unvested period %d\n", i)
			}
			continue
		}

		postLength, _ := strconv.ParseInt(post.VestingAccount.VestingPeriods[i].Length, 10, 64)

		if periodEndTime <= currentTime {
			// Vested period should be unchanged
			if preLength != postLength {
				fmt.Printf("❌ Vested period %d modified: %d -> %d\n",
					i, preLength, postLength)
			} else {
				fmt.Printf("✅ Vested period %d unchanged\n", i)
			}
		} else {
			// Unvested period should be extended
			expectedLength := preLength + threeMonths
			if postLength != expectedLength {
				fmt.Printf("❌ Period %d not properly extended: expected %d, got %d\n",
					i, expectedLength, postLength)
			} else {
				fmt.Printf("✅ Period %d properly extended\n", i)
			}
		}

		cumulativeTime += preLength
	}
}


func calculateVestedUnvested(preState AccountState, currentTime int64) (vested, unvested int64) {
	if preState.VestingAccount == nil {
		return 0, 0
	}

	startTime, _ := strconv.ParseInt(preState.VestingAccount.StartTime, 10, 64)
	cumulativeTime := startTime

	for _, period := range preState.VestingAccount.VestingPeriods {
		amount, _ := strconv.ParseInt(period.Amount[0].Amount, 10, 64)
		length, _ := strconv.ParseInt(period.Length, 10, 64)

		periodEndTime := cumulativeTime + length
		if periodEndTime <= currentTime {
			vested += amount
		} else {
			unvested += amount
		}
		cumulativeTime = periodEndTime
	}
	return vested, unvested
}

func verifyReplacementTransfer(oldAddress, newAddress string, pre, postOld, postNew AccountState) {
	fmt.Printf("\nVerifying replacement from %s to %s:\n", oldAddress, newAddress)
	currentTime := time.Now().Unix()

	// Calculate spendable/unvested based on spendable balance
	vested := sumCoins(pre.SpendableBalance.Balances)
	total := sumCoins(pre.TotalBalance.Balances)
	unvested := total - vested

	// Get actual post-upgrade splits
	postOldTotal := sumCoins(postOld.TotalBalance.Balances)
	postNewTotal := sumCoins(postNew.TotalBalance.Balances)

	fmt.Printf("Expected split based on spendable balance:\n")
	fmt.Printf("  - Vested (should stay): %d\n", vested)
	fmt.Printf("  - Unvested (should move): %d\n", unvested)
	fmt.Printf("Actual split after upgrade:\n")
	fmt.Printf("  - Old address balance: %d\n", postOldTotal)
	fmt.Printf("  - New address balance: %d\n", postNewTotal)

	if vested == postOldTotal {
		fmt.Printf("✅ Old address retained correct vested amount\n")
	} else {
		fmt.Printf("❌ Old address balance mismatch: expected=%d, got=%d\n", vested, postOldTotal)
	}

	if unvested == postNewTotal {
		fmt.Printf("✅ New address received correct unvested amount\n")
	} else {
		fmt.Printf("❌ New address balance mismatch: expected=%d, got=%d\n", unvested, postNewTotal)
	}

	if postNew.VestingAccount != nil {
		// Verify start time
		newStart, _ := strconv.ParseInt(postNew.VestingAccount.StartTime, 10, 64)
		if newStart < currentTime-300 || newStart > currentTime+300 {
			fmt.Printf("❌ New vesting schedule start time incorrect: expected ~%d, got %d\n",
				currentTime, newStart)
		} else {
			fmt.Printf("✅ New vesting schedule starts at correct time\n")
		}

		// Count remaining unvested periods from original account
		preStart, _ := strconv.ParseInt(pre.VestingAccount.StartTime, 10, 64)
		expectedUnvestedPeriods := 0
		cumulativeTime := preStart
		firstUnvestedFound := false
		//var firstUnvestedPartialLength int64 = 0

		for _, period := range pre.VestingAccount.VestingPeriods {
			length, _ := strconv.ParseInt(period.Length, 10, 64)
			if cumulativeTime+length > currentTime {
				if !firstUnvestedFound {
					// Calculate remaining time in first unvested period
					//elapsed := currentTime - cumulativeTime
					//firstUnvestedPartialLength = length - elapsed
					firstUnvestedFound = true
				}
				expectedUnvestedPeriods++
			}
			cumulativeTime += length
		}

		// Verify period count and structure
		actualPeriods := len(postNew.VestingAccount.VestingPeriods)
		if expectedUnvestedPeriods != actualPeriods || actualPeriods == (expectedUnvestedPeriods+1) {
			fmt.Printf("❌ Number of unvested periods incorrect: expected=%d, got=%d\n",
				expectedUnvestedPeriods, actualPeriods)
		} else {
			fmt.Printf("✅ Correct number of unvested periods\n")
		}

		// Check if first period is partial
		//if firstUnvestedFound && firstUnvestedPartialLength > 0 {
		//	firstPeriodLength, _ := strconv.ParseInt(postNew.VestingAccount.VestingPeriods[0].Length, 10, 64)
		//	if firstPeriodLength != firstUnvestedPartialLength {
		//		fmt.Printf("❌ First period length incorrect: expected=%d, got=%d\n",
		//			firstUnvestedPartialLength, firstPeriodLength)
		//	} else {
		//		fmt.Printf("✅ First period has correct partial length\n")
		//	}
		//}
	}
}

func sumCoins(coins []Coin) int64 {
	var total int64
	for _, coin := range coins {
		amount, _ := strconv.ParseInt(coin.Amount, 10, 64)
		total += amount
	}
	return total
}