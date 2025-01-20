package v2

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

const (
	UpgradeName = "v2"
)

// Enter the list of currentAddress -> newAddress to replace the pending vesgin to new address.
var addressReplacements = map[string]string{
	//"self12ugrzmzmk7zj56cjrt7dwjrwgatajyqvnepwzx": "self1scmpmsrv74r47fhj2fzcgeuque6pudam59prw8",
	//"self1krxfd67wmrjksq20xww53rm0wqmyxcew22whah": "self1kr30hqm2ezdjapspemdjgrt5lkxhsmwwr6ujtr",
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

func CreateUpgradeHandler(mm *module.Manager, configurator module.Configurator, accountKeeper authkeeper.AccountKeeper, bankkeeper bankkeeper.Keeper, ) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Starting upgrade v2")

		// 1. Run all module migrations first
		newVM, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			ctx.Logger().Error("Failed to run module migrations for v2", "error", err)
			return nil, err
		}

		// 2. After all modules are migrated, run your custom logic
		if err := updateVestingSchedules(ctx, accountKeeper, bankkeeper); err != nil {
			ctx.Logger().Error("Failed to execute v2 upgrade (vestings)", "error", err)
			return nil, err
		}

		ctx.Logger().Info("Completed upgrade v2 successfully")
		return newVM, nil
	}
}

func updateVestingSchedules(ctx sdk.Context, k authkeeper.AccountKeeper, bankkeeper bankkeeper.Keeper) error {
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
		if err := replaceAccountAddress(ctx, k, oldAddr, newAddr, bankkeeper); err != nil {
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

func replaceAccountAddress(ctx sdk.Context, k authkeeper.AccountKeeper, oldAddress, newAddress string,  bankKeeper bankkeeper.Keeper) error {
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
