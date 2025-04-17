package v2

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

const (
	UpgradeName = "v3"
)

// Enter the list of currentAddress -> newAddress to replace the pending vesgin to new address.
var addressReplacements = map[string]string{
	"self1jezc4atme56v75x5njqe4zuaccc4secug25wd3": "self1yt9pefssr0gzggmhlx30fmuqze0j6sh900xx3x",
	"self1fun8q0xuncfef6nkwh9njvvp4xqf4276x5sxgf": "self1veztmkrcrwf0ff49fu4y6mjd0wqpf4pcv8ruja",
	"self1krxfd67wmrjksq20xww53rm0wqmyxcew22whah": "self1gw586ensunqrk7x8yajqs3w4lcgmgtgugngqax",
	"self1qxjrq22m0gkcz7h73q4jvhmysmgja54s70amcp": "self1j3l8rersmt2p2fcv6zy2g6qmy2th7jkau4w7le",
	"self1e20929j3gng6cy72qapar630977vffqqzwxj75": "self1e5ux63egmatg42sn7ujr5ar0qg83pnukgl9q8y",
	"self12xes3fhuhfdech9gkyjhl526l6gdh3n3kwe3ml": "self17qf0ssjuvemeknrf9tspd0uatrpqhfhwvus7ml",
	"self1p9zmq9f5ftxwke6urd3vr98rypjhettfrsnna3": "self1xh72xjsy3c79s0u9mrhzehwm065c632ljrgtjc",
	"self1c0h75n6pfnl9pk80dktqnjwvqgz0tu2trfwg40": "self1havmjneetz96xdftg89nv5537g9tddnsn382fj",
	"self14ga5vmrskscuj3yktvjksm93sdt2f8r9k35pm0": "self1mwesu486zeu27xtrdl74nka8vhusk0tn34tslw",
	"self1sah0w5e2a2nxrru4t6e6n3v47xulklwvru7hmh": "self1rle4cakzj849xhg7zj86rscwrmm83cpganlf4z",
	"self1ychdx0fl0gt9c74afeeqr6ykv5j5rcqawxx2me": "self17xz6v4vtxcwfv793hj0cx2myav4f2lnycqyv2s",
	"self1vwvjfg8ezhuspk5lamkakahc32yudf4wkgrsh6": "self1vgl693sr0m8w76ycd9k8knhxydh4y9h5eg5sdy",
	// Add more address mappings as needed
}

// Enter the list of address for which the vesting schedule needs to be postponed by 3 months.
var vestingAddresses = []string{
	"self1fcahhgtw2llk06am4rala6khxjtj24zhhxn449",
	"self1uuemdshpnzsvh9fmnn3gnltvgqullc23kvr0ha",
	"self15mqsmcgnfgzsscxm5unv896cdxaq49sxqpmkfh",
	"self10cp5243eghvhpngm8yh5j5l7w2jeras670hazk",
	"self1kdwx0egjyhukrh57wycsj695c9my3q52w0j3vx",
	"self1kr30hqm2ezdjapspemdjgrt5lkxhsmwwr6ujtr",
	"self1v5hrqlv8dqgzvy0pwzqzg0gxy899rm4kwugwdu",
	"self1s7m8pytctmpejfpe4t06d05s29dmwgfaujrgnf",
	"self1gf2r58wmnt0wxa0lels39c2q2eeddcaqal403v",
	"self1wv2h4syek04xs5ya70nte64hzyg87r682he7rf",
	"self12ugrzmzmk7zj56cjrt7dwjrwgatajyqvnepwzx",
	"self1zh3dye5m5u78utp3kyvephg63h3zkgw333l9sv",
	"self1t50hj98rmnusr2yywvp7aaq4jwr98tfnercplu",
	"self102fgcqwkhcrwf6yv8jgen7v2gd0k4e0szpfh3d",
	"self1st6dracnvq7g203rpvzqtw0cyy86nj03w3uwza",
	"self132vnr88qpa4gkmtdv3ly0kpehr7e8zxmanflx2",
	"self1q3a9nggfp94wv6yntjt3xq9wc6gfnnkrhlr6uf",
	"self17cfx0jp7dur420aavhf23xxgxydk2c34nc3n33",
	"self1scmpmsrv74r47fhj2fzcgeuque6pudam59prw8",
	"self1cyf2hyjtpwkeh7y466t4ravyfcv0ze4yyk3a3r",
	"self1qvsus5qg8yhre7k2c78xkkw4nvqqgev7qu2n9f",
	"self1fxezqx5w9aw7rfteswm6uzdej56sp3sleup43l",
	"self1jezc4atme56v75x5njqe4zuaccc4secug25wd3",
	"self1fun8q0xuncfef6nkwh9njvvp4xqf4276x5sxgf",
	"self1krxfd67wmrjksq20xww53rm0wqmyxcew22whah",
	"self1qxjrq22m0gkcz7h73q4jvhmysmgja54s70amcp",
}

func CreateUpgradeHandler(mm *module.Manager, configurator module.Configurator, accountKeeper authkeeper.AccountKeeper, bankkeeper bankkeeper.Keeper, stakingkeeper *stakingkeeper.Keeper, distrkeeper distrkeeper.Keeper) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Starting upgrade v2")

		// 1. Run all module migrations first
		newVM, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			ctx.Logger().Error("Failed to run module migrations for v2", "error", err)
			return nil, err
		}

		// 2. After all modules are migrated, run your custom logic
		if err := updateVestingSchedules(ctx, accountKeeper, bankkeeper, *stakingkeeper, distrkeeper); err != nil {
			ctx.Logger().Error("Failed to execute v2 upgrade (vestings)", "error", err)
			return nil, err
		}

		ctx.Logger().Info("Completed upgrade v2 successfully")
		return newVM, nil
	}
}

func updateVestingSchedules(ctx sdk.Context, k authkeeper.AccountKeeper, bankkeeper bankkeeper.Keeper, stakingkeeper stakingkeeper.Keeper, distrkeeper distrkeeper.Keeper) error {
	monthsToAdd := int64(3)
	for _, addr := range vestingAddresses {
		if err := updateVestingAccount(ctx, k, addr, monthsToAdd); err != nil {
			if err.Error() == fmt.Sprintf("account not found: %s", addr) {
				ctx.Logger().Info("Skipping non-existent account for vesting update",
					"address", addr)
				continue
			}

			return fmt.Errorf("failed to update vesting for %s: %w", addr, err)
		}
	}

	// Then handle address replacements separately
	ctx.Logger().Info("Processing address replacements", "count", len(addressReplacements))
	for oldAddr, newAddr := range addressReplacements {
		if err := replaceAccountAddress2(ctx, k, oldAddr, newAddr, bankkeeper, stakingkeeper, distrkeeper); err != nil {
			if err.Error() == fmt.Sprintf("account not found: %s", oldAddr) {
				ctx.Logger().Info("Skipping non-existent account for address replacement",
					"old_address", oldAddr,
					"new_address", newAddr)
				continue
			}
			return fmt.Errorf("failed to replace address for %s: %w", oldAddr, err)
		}
	}

	return nil
}

func replaceAccountAddress(ctx sdk.Context, k authkeeper.AccountKeeper, oldAddress, newAddress string, bankKeeper bankkeeper.Keeper) error {
	// Validate new address doesn't exist
	newAddr, err := sdk.AccAddressFromBech32(newAddress)
	if err != nil {
		return fmt.Errorf("invalid new address %s: %w", newAddress, err)
	}
	if k.GetAccount(ctx, newAddr) != nil {
		return fmt.Errorf("new address %s already exists", newAddress)
	}

	oldAddr, err := sdk.AccAddressFromBech32(oldAddress)
	if err != nil {
		return fmt.Errorf("invalid old address %s: %w", oldAddress, err)
	}

	// Get the old account
	oldAcc, err := getPeriodicVestingAccount(ctx, k, oldAddress)
	if err != nil {
		return err
	}

	currentTime := ctx.BlockTime().Unix()

	// Calculate unvested periods and coins
	var unvestedPeriods []vestingtypes.Period
	var vestedPeriods []vestingtypes.Period
	var unvestedCoins sdk.Coins
	cumulativeTime := oldAcc.StartTime
	partialElapsed := int64(0)
	firstUnVested := true

	// Find unvested periods
	for i, period := range oldAcc.VestingPeriods {
		if cumulativeTime+period.Length > currentTime {
			// This and all subsequent periods are unvested
			unvestedPeriods = append(unvestedPeriods, oldAcc.VestingPeriods[i:]...)
			for _, p := range oldAcc.VestingPeriods[i:] {
				unvestedCoins = unvestedCoins.Add(p.Amount...)
			}

			if firstUnVested {
				usedInThisPeriod := currentTime - cumulativeTime
				partialElapsed = usedInThisPeriod
				firstUnVested = false
			}
			break
		}
		vestedPeriods = append(vestedPeriods, period)
		cumulativeTime += period.Length
	}

	// If no unvested periods, nothing to migrate
	if len(unvestedPeriods) == 0 {
		ctx.Logger().Info("No unvested periods to migrate", "address", oldAddress)
		return nil
	}

	if partialElapsed > 0 {
		if partialElapsed > unvestedPeriods[0].Length {
			partialElapsed = unvestedPeriods[0].Length
		}

		unvestedPeriods[0].Length -= partialElapsed
	}

	// Create new base account with proper account number
	baseAcc := authtypes.NewBaseAccountWithAddress(newAddr)

	convertVestingToBaseAccount(ctx, k, oldAcc)
	if err := bankKeeper.SendCoins(ctx, oldAddr, newAddr, unvestedCoins); err != nil {
		return fmt.Errorf("failed to send coins: %w", err)
	}

	// Create new account with only unvested amounts
	newAcc := &vestingtypes.PeriodicVestingAccount{
		BaseVestingAccount: &vestingtypes.BaseVestingAccount{
			BaseAccount:      baseAcc,
			OriginalVesting:  unvestedCoins,
			DelegatedFree:    sdk.NewCoins(),
			DelegatedVesting: sdk.NewCoins(),
			EndTime:          oldAcc.EndTime,
		},
		StartTime:      currentTime,
		VestingPeriods: unvestedPeriods,
	}

	// Update old account to only contain vested periods
	oldAcc.VestingPeriods = vestedPeriods
	oldAcc.OriginalVesting = oldAcc.OriginalVesting.Sub(unvestedCoins...)
	oldAcc.EndTime = currentTime

	// Save both accounts
	k.SetAccount(ctx, oldAcc)
	k.SetAccount(ctx, newAcc)

	ctx.Logger().Info("Successfully split vesting account",
		"old_address", oldAddress,
		"new_address", newAddress,
		"vested_periods", len(vestedPeriods),
		"unvested_periods", len(unvestedPeriods),
		"unvested_coins", unvestedCoins)

	return nil
}

func convertVestingToBaseAccount(
	ctx sdk.Context,
	accountKeeper authkeeper.AccountKeeper,
	vestAcc *vestingtypes.PeriodicVestingAccount,
) authtypes.AccountI {
	oldBaseAcc := vestAcc.BaseVestingAccount.BaseAccount
	// Construct a fresh BaseAccount with same address, account number, sequence
	newBase := authtypes.NewBaseAccount(
		oldBaseAcc.GetAddress(),
		oldBaseAcc.GetPubKey(),
		oldBaseAcc.GetAccountNumber(),
		oldBaseAcc.GetSequence(),
	)

	// Overwrite the store so that address is now a BaseAccount
	accountKeeper.SetAccount(ctx, newBase)
	return newBase
}

func updateVestingAccount(ctx sdk.Context, k authkeeper.AccountKeeper, address string, monthsToAdd int64) error {
	acc, err := getPeriodicVestingAccount(ctx, k, address)
	if err != nil {
		return err
	}

	currentTime := ctx.BlockTime().Unix()
	secondsToAdd := monthsToAdd * 2628000 // ~30.44 days per month

	ctx.Logger().Info("Account details",
		"address", address,
		"current_start_time", acc.StartTime,
		"current_end_time", acc.EndTime,
		"periods", len(acc.VestingPeriods),
		"current_block_time", currentTime)

	// Find the first unvested period
	cumulativeTime := acc.StartTime
	firstUnvestedIdx := 0
	for i, period := range acc.VestingPeriods {
		cumulativeTime += period.Length
		if cumulativeTime > currentTime {
			firstUnvestedIdx = i
			break
		}
	}

	ctx.Logger().Info("Vesting analysis",
		"address", address,
		"first_unvested_period", firstUnvestedIdx,
		"total_periods", len(acc.VestingPeriods),
		"cumulative_time", cumulativeTime)

	if firstUnvestedIdx < len(acc.VestingPeriods) {
		// Only modify unvested periods
		newVestingPeriods := make([]vestingtypes.Period, len(acc.VestingPeriods))
		newEndTime := acc.StartTime

		// Copy all periods
		for i := 0; i < len(acc.VestingPeriods); i++ {
			newVestingPeriods[i] = acc.VestingPeriods[i]

			if i == firstUnvestedIdx {
				// Add the delay to the first unvested period
				newVestingPeriods[i].Length += secondsToAdd
			}
			newEndTime += newVestingPeriods[i].Length
		}

		acc.VestingPeriods = newVestingPeriods
		acc.BaseVestingAccount.EndTime = newEndTime

		ctx.Logger().Info("Updating vesting schedule",
			"address", address,
			"old_start_time", acc.StartTime,
			"new_end_time", acc.EndTime,
			"seconds_added", secondsToAdd)

		k.SetAccount(ctx, acc)
		ctx.Logger().Info("Successfully updated account", "address", address)
	} else {
		ctx.Logger().Info("No unvested periods found - skipping", "address", address)
	}

	return nil
}

func getPeriodicVestingAccount(ctx sdk.Context, k authkeeper.AccountKeeper, address string) (*vestingtypes.PeriodicVestingAccount, error) {
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		ctx.Logger().Error("Invalid address format", "address", address, "error", err)
		return nil, fmt.Errorf("invalid address %s: %w", address, err)
	}

	acc := k.GetAccount(ctx, addr)
	if acc == nil {
		ctx.Logger().Error("Account not found", "address", address)
		return nil, fmt.Errorf("account not found: %s", address)
	}

	periodicAcc, ok := acc.(*vestingtypes.PeriodicVestingAccount)
	if !ok {
		ctx.Logger().Error("Account is not a periodic vesting account",
			"address", address,
			"actual_type", fmt.Sprintf("%T", acc))
		return nil, fmt.Errorf("account %s is not a periodic vesting account", address)
	}

	return periodicAcc, nil
}

// -----------------------------------------------------------------------------
// helper constant – staking getters take uint16 for `limit`
// -----------------------------------------------------------------------------
const maxRetrieve = uint16(1<<16 - 1) // 65535

// First, we need a more flexible account retrieval function
func getVestingAccount(ctx sdk.Context, k authkeeper.AccountKeeper, address string) (authtypes.AccountI, string, error) {
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		ctx.Logger().Error("Invalid address format", "address", address, "error", err)
		return nil, "", fmt.Errorf("invalid address %s: %w", address, err)
	}

	acc := k.GetAccount(ctx, addr)
	if acc == nil {
		ctx.Logger().Error("Account not found", "address", address)
		return nil, "", fmt.Errorf("account not found: %s", address)
	}

	// Check for PeriodicVestingAccount
	if periodicAcc, ok := acc.(*vestingtypes.PeriodicVestingAccount); ok {
		return periodicAcc, "periodic", nil
	}

	// Check for PermanentLockedAccount
	if permanentAcc, ok := acc.(*vestingtypes.PermanentLockedAccount); ok {
		return permanentAcc, "permanent", nil
	}

	// Neither periodic nor permanent
	ctx.Logger().Error("Account is not a supported vesting account type",
		"address", address,
		"actual_type", fmt.Sprintf("%T", acc))
	return nil, "", fmt.Errorf("account %s is not a supported vesting account type", address)
}

// Updated replaceAccountAddress2 function to handle both account types
func replaceAccountAddress2(
	ctx sdk.Context,
	k authkeeper.AccountKeeper,
	oldAddress string,
	newAddress string,
	bankKeeper bankkeeper.Keeper,
	stakingKeeper stakingkeeper.Keeper,
	distrKeeper distrkeeper.Keeper,
) error {
	// ---------- sanity -------------------------------------------------------
	newAddr, err := sdk.AccAddressFromBech32(newAddress)
	if err != nil {
		return fmt.Errorf("invalid new address %s: %w", newAddress, err)
	}
	if k.GetAccount(ctx, newAddr) != nil {
		return fmt.Errorf("new address %s already exists", newAddress)
	}

	oldAddr, err := sdk.AccAddressFromBech32(oldAddress)
	if err != nil {
		return fmt.Errorf("invalid old address %s: %w", oldAddress, err)
	}

	oldAcc, accType, err := getVestingAccount(ctx, k, oldAddress)
	if err != nil {
		return err
	}

	bondDenom := stakingKeeper.BondDenom(ctx)
	now := ctx.BlockTime().Unix()

	// ---------- unlock the old account up-front ------------------------------
	// Convert the account to a BaseAccount temporarily to allow transfers
	var oldBaseAcc *authtypes.BaseAccount
	var originalVesting sdk.Coins
	var delegatedFree, delegatedVesting sdk.Coins
	var vestingSchedule interface{} // Will hold either vesting periods or lock

	switch accType {
	case "periodic":
		periodicAcc := oldAcc.(*vestingtypes.PeriodicVestingAccount)
		oldBaseAcc = periodicAcc.BaseVestingAccount.BaseAccount
		originalVesting = periodicAcc.OriginalVesting
		delegatedFree = periodicAcc.DelegatedFree
		delegatedVesting = periodicAcc.DelegatedVesting
		_ = periodicAcc.EndTime
		vestingSchedule = periodicAcc.VestingPeriods
		// Convert to base account
		baseAcc := authtypes.NewBaseAccount(
			oldBaseAcc.GetAddress(),
			oldBaseAcc.GetPubKey(),
			oldBaseAcc.GetAccountNumber(),
			oldBaseAcc.GetSequence(),
		)
		k.SetAccount(ctx, baseAcc)

	case "permanent":
		permanentAcc := oldAcc.(*vestingtypes.PermanentLockedAccount)
		oldBaseAcc = permanentAcc.BaseVestingAccount.BaseAccount
		originalVesting = permanentAcc.OriginalVesting
		delegatedFree = permanentAcc.DelegatedFree
		delegatedVesting = permanentAcc.DelegatedVesting
		_ = permanentAcc.EndTime
		vestingSchedule = "permanent" // Special marker for permanent locked account
		// Convert to base account
		baseAcc := authtypes.NewBaseAccount(
			oldBaseAcc.GetAddress(),
			oldBaseAcc.GetPubKey(),
			oldBaseAcc.GetAccountNumber(),
			oldBaseAcc.GetSequence(),
		)
		k.SetAccount(ctx, baseAcc)

	default:
		return fmt.Errorf("unsupported account type: %s", accType)
	}

	// Get critical account state values
	totalBalance := bankKeeper.GetBalance(ctx, oldAddr, bondDenom).Amount
	originalVestingAmount := originalVesting.AmountOf(bondDenom)

	// For permanent accounts, all coins are unvested
	// For periodic accounts, we need to calculate unvested amount
	var unvestedCoins sdk.Coins
	var unvestedAmount sdk.Int
	var movedDelVesting sdk.Coins

	if accType == "permanent" {
		// For permanent locked accounts, everything is unvested
		unvestedCoins = originalVesting
		unvestedAmount = originalVestingAmount
		movedDelVesting = delegatedVesting
	} else {
		// For periodic vesting, calculate unvested coins
		periodicAcc := oldAcc.(*vestingtypes.PeriodicVestingAccount)
		var unvestedPeriods []vestingtypes.Period
		var vestedPeriods []vestingtypes.Period
		var elapsedFirst int64

		cumTime := periodicAcc.StartTime
		for i, p := range periodicAcc.VestingPeriods {
			if cumTime+p.Length > now {
				unvestedPeriods = append(unvestedPeriods, periodicAcc.VestingPeriods[i:]...)
				for _, up := range unvestedPeriods {
					unvestedCoins = unvestedCoins.Add(up.Amount...)
				}
				elapsedFirst = now - cumTime
				break
			}
			vestedPeriods = append(vestedPeriods, p)
			cumTime += p.Length
		}
		if elapsedFirst > 0 && len(unvestedPeriods) > 0 {
			if elapsedFirst > unvestedPeriods[0].Length {
				elapsedFirst = unvestedPeriods[0].Length
			}
			unvestedPeriods[0].Length -= elapsedFirst
		}

		// Store updated vesting schedule
		vestingSchedule = unvestedPeriods
		unvestedAmount = unvestedCoins.AmountOf(bondDenom)

		// Move ALL delegated_vesting (these are all unvested tokens)
		movedDelVesting = delegatedVesting
	}

	vestedAmount := originalVestingAmount.Sub(unvestedAmount)

	// Log what we're starting with
	ctx.Logger().Info("account state pre-migration",
		"address", oldAddress,
		"account_type", accType,
		"total_balance", totalBalance,
		"original_vesting", originalVestingAmount,
		"unvested_amount", unvestedAmount,
		"vested_amount", vestedAmount,
		"delegated_vesting", delegatedVesting,
		"delegated_free", delegatedFree,
		"now", now,
		" ctx.BlockTime()", ctx.BlockHeight())

	// ---------- Handle delegated tokens --------------------------------------
	// Do NOT move any delegated_free
	delegatedFreeToMove := sdk.NewCoins()

	ctx.Logger().Info("delegation amounts to move",
		"delegated_vesting_to_move", movedDelVesting,
		"delegated_free_to_move", delegatedFreeToMove)

	// ---------- Update delegations to reflect the split --------------------
	var newDelegations []stakingtypes.Delegation
	delList := stakingKeeper.GetDelegatorDelegations(ctx, oldAddr, maxRetrieve)

	// Calculate total delegated vesting amount
	totalDelegatedVesting := movedDelVesting.AmountOf(bondDenom)

	// Skip if nothing to move
	if totalDelegatedVesting.IsZero() {
		ctx.Logger().Info("No delegated vesting to move")
	} else {
		// Track processed validators to avoid double processing
		processedValidators := make(map[string]bool)

		// Find tokens per validator by processing all delegations
		validatorTokens := make(map[string]sdk.Int)
		validatorShares := make(map[string]sdk.Dec)
		totalTokens := sdk.ZeroInt()

		// First pass: calculate total delegated tokens and tokens per validator
		for _, del := range delList {
			valAddr := del.ValidatorAddress

			// Skip if already processed
			if processedValidators[valAddr] {
				continue
			}
			processedValidators[valAddr] = true

			val, found := stakingKeeper.GetValidator(ctx, del.GetValidatorAddr())
			if !found {
				continue
			}

			// Calculate tokens from shares
			tokens := val.TokensFromShares(del.Shares).TruncateInt()
			validatorTokens[valAddr] = tokens
			validatorShares[valAddr] = del.Shares
			totalTokens = totalTokens.Add(tokens)
		}

		// Only proceed if we have delegations
		if !totalTokens.IsZero() {
			// Reset processed validators map
			processedValidators = make(map[string]bool)

			// Second pass: move vesting shares proportionally
			for _, del := range delList {
				valAddr := del.ValidatorAddress

				// Skip if already processed
				if processedValidators[valAddr] {
					continue
				}
				processedValidators[valAddr] = true

				val, found := stakingKeeper.GetValidator(ctx, del.GetValidatorAddr())
				if !found {
					continue
				}

				// Calculate tokens to move (proportional to vesting amount)
				var vestingToMove sdk.Int
				if totalTokens.IsZero() {
					vestingToMove = sdk.ZeroInt()
				} else {
					// Calculate proportional amount of delegated_vesting to move
					ratio := sdk.NewDecFromInt(validatorTokens[valAddr]).QuoInt(totalTokens)
					vestingToMove = sdk.NewDecFromInt(totalDelegatedVesting).Mul(ratio).TruncateInt()
				}

				ctx.Logger().Info("validator vesting calculation",
					"validator", valAddr,
					"vesting_tokens_to_move", vestingToMove,
					"total_tokens", validatorTokens[valAddr])

				// Skip if nothing to move
				if vestingToMove.IsZero() {
					continue
				}

				// Calculate shares to move
				sharesToMove, _ := val.SharesFromTokensTruncated(vestingToMove)

				// Make sure we don't try to move more shares than exist
				if sharesToMove.GT(del.Shares) {
					sharesToMove = del.Shares
				}

				// Create new delegation for the target address
				newDelegations = append(newDelegations,
					stakingtypes.NewDelegation(newAddr, del.GetValidatorAddr(), sharesToMove))

				// Reduce shares from the old delegation
				del.Shares = del.Shares.Sub(sharesToMove)

				// Update or remove the old delegation
				if del.Shares.IsZero() {
					stakingKeeper.RemoveDelegation(ctx, del)
				} else {
					stakingKeeper.SetDelegation(ctx, del)
				}

				ctx.Logger().Info("delegation update",
					"validator", valAddr,
					"old_shares", del.Shares.Add(sharesToMove),
					"new_shares", del.Shares,
					"moved_shares", sharesToMove)
			}
		}
	}

	// ---------- Now handle liquid balance ----------------------------------
	// Calculate how much unvested liquid balance to transfer
	// Total unvested minus what's already in delegations
	totalUnvestedAmount := unvestedAmount
	delegatedVestingAmount := movedDelVesting.AmountOf(bondDenom)

	// Calculate liquid unvested amount to transfer
	liquidUnvestedToTransfer := totalUnvestedAmount.Sub(delegatedVestingAmount)

	// Ensure we don't go negative due to rounding
	if liquidUnvestedToTransfer.IsNegative() {
		liquidUnvestedToTransfer = sdk.ZeroInt()
	}

	ctx.Logger().Info("liquid calculation",
		"total_unvested", totalUnvestedAmount,
		"delegated_vesting_amount", delegatedVestingAmount,
		"liquid_unvested_to_transfer", liquidUnvestedToTransfer)

	// Only transfer if there's something to transfer
	var transferredLiquid sdk.Int
	if liquidUnvestedToTransfer.IsPositive() {
		transfer := sdk.NewCoins(sdk.NewCoin(bondDenom, liquidUnvestedToTransfer))
		ctx.Logger().Info("Attempting to send coins",
			"from", oldAddr.String(),
			"to", newAddr.String(),
			"amount", transfer.String())

		if err := bankKeeper.SendCoins(ctx, oldAddr, newAddr, transfer); err != nil {
			return fmt.Errorf("sending unvested liquid coins: %w", err)
		}
		transferredLiquid = liquidUnvestedToTransfer

		// Verify the old account has the expected balance
		expectedOldBalance := totalBalance.Sub(liquidUnvestedToTransfer)
		gotOldBalance := bankKeeper.GetBalance(ctx, oldAddr, bondDenom).Amount

		ctx.Logger().Info("Post-transfer balance check",
			"expected_old_balance", expectedOldBalance,
			"actual_old_balance", gotOldBalance)

		// Add more tolerance for rounding errors
		diff := expectedOldBalance.Sub(gotOldBalance)
		if diff.Abs().GT(sdk.NewInt(100)) { // Allow small difference due to rounding
			ctx.Logger().Info("Balance discrepancy detected",
				"difference", diff,
				"expected", expectedOldBalance,
				"actual", gotOldBalance)
		}
	} else {
		transferredLiquid = sdk.ZeroInt()
		ctx.Logger().Info("No liquid unvested to transfer")
	}

	// ---------- unbonding delegations ---------------------------------------
	ubdList := stakingKeeper.GetUnbondingDelegations(ctx, oldAddr, maxRetrieve)
	for _, ubd := range ubdList {
		// remove the old object
		stakingKeeper.RemoveUnbondingDelegation(ctx, ubd)

		// build a mirror object under the new delegator
		valAddr, _ := sdk.ValAddressFromBech32(ubd.ValidatorAddress)

		newUBD := stakingtypes.UnbondingDelegation{
			DelegatorAddress: newAddr.String(),
			ValidatorAddress: valAddr.String(),
			Entries:          ubd.Entries, // copy-by-value is fine
		}
		stakingKeeper.SetUnbondingDelegation(ctx, newUBD)
	}

	// ---------- redelegations -----------------------------------------------
	redList := stakingKeeper.GetRedelegations(ctx, oldAddr, maxRetrieve)
	for _, red := range redList {
		stakingKeeper.RemoveRedelegation(ctx, red)

		srcVal, _ := sdk.ValAddressFromBech32(red.ValidatorSrcAddress)
		dstVal, _ := sdk.ValAddressFromBech32(red.ValidatorDstAddress)

		newRed := stakingtypes.Redelegation{
			DelegatorAddress:    newAddr.String(),
			ValidatorSrcAddress: srcVal.String(),
			ValidatorDstAddress: dstVal.String(),
			Entries:             red.Entries, // keep every entry intact
		}
		stakingKeeper.SetRedelegation(ctx, newRed)
	}

	// ---------- create new vesting account ----------------------------------
	baseAcc := authtypes.NewBaseAccount(newAddr, nil, k.NextAccountNumber(ctx), 0)

	// Log the final calculated amounts for debugging
	ctx.Logger().Info("final migration amounts",
		"old_address", oldAddress,
		"new_address", newAddress,
		"account_type", accType,
		"unvested_coins", unvestedCoins.String(),
		"moved_del_vesting", movedDelVesting.String(),
		"moved_del_free", delegatedFreeToMove.String(),
		"liquid_transferred", transferredLiquid.String())

	// Create the appropriate new account based on the original type
	var newVestingAcc authtypes.AccountI

	if accType == "permanent" {
		// Create a new permanent locked account
		permanentAcc := vestingtypes.NewPermanentLockedAccount(
			baseAcc,
			unvestedCoins,
		)

		// Manually set delegated balances since the constructor doesn't accept them
		permanentAcc.DelegatedFree = delegatedFreeToMove
		permanentAcc.DelegatedVesting = movedDelVesting

		newVestingAcc = permanentAcc
	} else {
		// Create a new periodic vesting account
		periodicVestingSchedule := vestingSchedule.([]vestingtypes.Period)
		newVestingAcc = vestingtypes.NewPeriodicVestingAccount(
			baseAcc,
			unvestedCoins,
			now,
			periodicVestingSchedule,
		)

		// Manually set delegated balances since the constructor doesn't accept them
		if periodicAcc, ok := newVestingAcc.(*vestingtypes.PeriodicVestingAccount); ok {
			periodicAcc.DelegatedFree = delegatedFreeToMove
			periodicAcc.DelegatedVesting = movedDelVesting
		}
	}

	// Update old account
	if accType == "permanent" {
		// For permanent locked accounts, just convert to a base account
		// since all tokens are either transferred or in delegations
		oldBaseAccUpdated := authtypes.NewBaseAccount(
			oldAddr,
			oldBaseAcc.GetPubKey(),
			oldBaseAcc.GetAccountNumber(),
			oldBaseAcc.GetSequence(),
		)
		k.SetAccount(ctx, oldBaseAccUpdated)
	} else {
		// For periodic vesting accounts, update with remaining vested periods
		periodicAcc := oldAcc.(*vestingtypes.PeriodicVestingAccount)
		vestedPeriods := []vestingtypes.Period{}

		// Extract vested periods from the original account
		cumTime := periodicAcc.StartTime
		for _, p := range periodicAcc.VestingPeriods {
			if cumTime+p.Length <= now {
				vestedPeriods = append(vestedPeriods, p)
			}
			cumTime += p.Length
		}

		// Create updated old account
		updatedOldAcc := vestingtypes.NewPeriodicVestingAccount(
			oldBaseAcc,
			originalVesting.Sub(unvestedCoins...),
			periodicAcc.StartTime,
			vestedPeriods,
		)

		// Manually set delegated balances
		updatedOldAcc.DelegatedFree = delegatedFree
		updatedOldAcc.DelegatedVesting = delegatedVesting.Sub(movedDelVesting...)
		k.SetAccount(ctx, updatedOldAcc)
	}

	// Save the new account
	k.SetAccount(ctx, newVestingAcc)

	// ---------- persist new delegations last ---------------------------------
	// In replaceAccountAddress2 function, update the delegations section:

	// ---------- persist new delegations last ---------------------------------
	for _, nd := range newDelegations {
		stakingKeeper.SetDelegation(ctx, nd)

		// Properly set delegator starting info for new delegations
		valAddr := nd.GetValidatorAddr()

		// Need to lookup the validator for each new delegation
		val, found := stakingKeeper.GetValidator(ctx, valAddr)
		if !found {
			return fmt.Errorf("validator %s not found", valAddr)
		}

		// Calculate stake in bond tokens for the new shares
		stakeDec := val.TokensFromSharesTruncated(nd.Shares)

		// Get the current period from validator rewards
		period := distrKeeper.GetValidatorCurrentRewards(ctx, valAddr).Period
		height := uint64(ctx.BlockHeight())

		// Create new starting info with correct values for new delegator
		startInfo := distrtypes.NewDelegatorStartingInfo(period, stakeDec, height)

		// Set the correct starting info for the new delegator
		distrKeeper.SetDelegatorStartingInfo(ctx, valAddr, newAddr, startInfo)

		// UPDATE: Find the remaining delegation for the old address
		oldDel, found := stakingKeeper.GetDelegation(ctx, oldAddr, valAddr)
		if found && !oldDel.Shares.IsZero() {
			// Calculate the reduced stake for the old delegator
			oldStakeDec := val.TokensFromSharesTruncated(oldDel.Shares)

			// Update the old delegator's starting info to match reduced shares
			updatedOldInfo := distrtypes.NewDelegatorStartingInfo(period, oldStakeDec, height)
			distrKeeper.SetDelegatorStartingInfo(ctx, valAddr, oldAddr, updatedOldInfo)

			ctx.Logger().Info("Updated old delegator starting info",
				"old_address", oldAddr.String(),
				"validator", valAddr.String(),
				"updated_stake", oldStakeDec.String())
		}
	}

	// Verify final balances
	newBalance := bankKeeper.GetBalance(ctx, newAddr, bondDenom).Amount
	oldBalance := bankKeeper.GetBalance(ctx, oldAddr, bondDenom).Amount

	ctx.Logger().Info("final balances",
		"old_address_balance", oldBalance,
		"new_address_balance", newBalance)

	ctx.Logger().Info("account migration complete",
		"old", oldAddress,
		"new", newAddress,
		"account_type", accType,
		"liquid_moved", transferredLiquid,
		"unvested_moved", unvestedCoins.String(),
		"del_vesting_moved", movedDelVesting.String(),
		"del_free_moved", delegatedFreeToMove.String(),
		"delegations", len(newDelegations),
		"ubd", len(ubdList),
		"red", len(redList),
	)

	return nil
}
