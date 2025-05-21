package v2

import (
	"context"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"fmt"
	cmttypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	"sort"
)

const (
	UpgradeName = "v2"

	// BlockMaxBytes is the max bytes for a block, 3mb
	BlockMaxBytes = int64(22020096)

	// BlockMaxGas is the max gas allowed in a block
	BlockMaxGas = int64(-1)
)

// AddressReplacement defines a struct for address replacement to ensure deterministic processing
type AddressReplacement struct {
	OldAddress string
	NewAddress string
}

// Enter the list of currentAddress -> newAddress to replace the pending vesgin to new address.
var addressReplacements = []AddressReplacement{
	{OldAddress: "self1jezc4atme56v75x5njqe4zuaccc4secug25wd3", NewAddress: "self1yt9pefssr0gzggmhlx30fmuqze0j6sh900xx3x"},
	{OldAddress: "self1fun8q0xuncfef6nkwh9njvvp4xqf4276x5sxgf", NewAddress: "self1veztmkrcrwf0ff49fu4y6mjd0wqpf4pcv8ruja"},
	{OldAddress: "self1krxfd67wmrjksq20xww53rm0wqmyxcew22whah", NewAddress: "self1gw586ensunqrk7x8yajqs3w4lcgmgtgugngqax"},
	{OldAddress: "self1qxjrq22m0gkcz7h73q4jvhmysmgja54s70amcp", NewAddress: "self1j3l8rersmt2p2fcv6zy2g6qmy2th7jkau4w7le"},
	{OldAddress: "self1e20929j3gng6cy72qapar630977vffqqzwxj75", NewAddress: "self1e5ux63egmatg42sn7ujr5ar0qg83pnukgl9q8y"},
	{OldAddress: "self12xes3fhuhfdech9gkyjhl526l6gdh3n3kwe3ml", NewAddress: "self17qf0ssjuvemeknrf9tspd0uatrpqhfhwvus7ml"},
	{OldAddress: "self1p9zmq9f5ftxwke6urd3vr98rypjhettfrsnna3", NewAddress: "self1xh72xjsy3c79s0u9mrhzehwm065c632ljrgtjc"},
	{OldAddress: "self1c0h75n6pfnl9pk80dktqnjwvqgz0tu2trfwg40", NewAddress: "self1havmjneetz96xdftg89nv5537g9tddnsn382fj"},
	{OldAddress: "self14ga5vmrskscuj3yktvjksm93sdt2f8r9k35pm0", NewAddress: "self1mwesu486zeu27xtrdl74nka8vhusk0tn34tslw"},
	{OldAddress: "self1sah0w5e2a2nxrru4t6e6n3v47xulklwvru7hmh", NewAddress: "self1rle4cakzj849xhg7zj86rscwrmm83cpganlf4z"},
	{OldAddress: "self1ychdx0fl0gt9c74afeeqr6ykv5j5rcqawxx2me", NewAddress: "self17xz6v4vtxcwfv793hj0cx2myav4f2lnycqyv2s"},
	{OldAddress: "self1vwvjfg8ezhuspk5lamkakahc32yudf4wkgrsh6", NewAddress: "self1vgl693sr0m8w76ycd9k8knhxydh4y9h5eg5sdy"},
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

func CreateUpgradeHandler(mm *module.Manager, configurator module.Configurator, paramsKeeper paramskeeper.Keeper, accountKeeper authkeeper.AccountKeeper, bankkeeper bankkeeper.Keeper, stakingkeeper *stakingkeeper.Keeper, distrkeeper distrkeeper.Keeper, consensuskeeper consensusparamkeeper.Keeper, ibckeeper ibckeeper.Keeper, distkeeper distrkeeper.Keeper) upgradetypes.UpgradeHandler {
	return func(context context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(context)
		ctx.Logger().Info("Starting upgrade v2")

		// 02-client ---------------------------------------------------------------
		clientSS, found := paramsKeeper.GetSubspace(ibcclienttypes.SubModuleName)
		if !found {
			clientSS = paramsKeeper.Subspace(ibcclienttypes.SubModuleName) // only if not present
		}
		if !clientSS.HasKeyTable() {
			clientSS = clientSS.WithKeyTable(ibcclienttypes.ParamKeyTable())
		}
		if !clientSS.Has(ctx, ibcclienttypes.KeyAllowedClients) {
			def := ibcclienttypes.DefaultParams()
			clientSS.SetParamSet(ctx, &def)
		}

		// 03-connection -----------------------------------------------------------
		connSS, found := paramsKeeper.GetSubspace(ibcconnectiontypes.SubModuleName)
		if !found {
			connSS = paramsKeeper.Subspace(ibcconnectiontypes.SubModuleName)
		}
		if !connSS.HasKeyTable() {
			connSS = connSS.WithKeyTable(ibcconnectiontypes.ParamKeyTable())
		}
		if !connSS.Has(ctx, ibcconnectiontypes.KeyMaxExpectedTimePerBlock) {
			def := ibcconnectiontypes.DefaultParams()
			connSS.SetParamSet(ctx, &def)
		}

		//1. Run all module migrations first
		newVM, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			ctx.Logger().Error("Failed to run module migrations for v2", "error", err)
			return nil, err
		}

		// 2. After all modules are migrated, run your custom logic
		if err := updateVestingSchedules(ctx, accountKeeper, bankkeeper, *stakingkeeper, distrkeeper, plan.Height); err != nil {
			ctx.Logger().Error("Failed to execute v2 upgrade (vestings)", "error", err)
			return nil, err
		}

		// Set the next block limits
		defaultConsensusParams := cmttypes.DefaultConsensusParams().ToProto()
		defaultConsensusParams.Block.MaxBytes = BlockMaxBytes
		defaultConsensusParams.Block.MaxGas = BlockMaxGas
		err = consensuskeeper.ParamsStore.Set(ctx, defaultConsensusParams)
		if err != nil {
			return nil, err
		}

		ibcClientSubspace, _ := paramsKeeper.GetSubspace(ibcclienttypes.SubModuleName)

		if !ibcClientSubspace.HasKeyTable() {
			ibcClientSubspace = ibcClientSubspace.WithKeyTable(ibcclienttypes.ParamKeyTable())
		}

		clientParams := ibcclienttypes.DefaultParams()
		ibckeeper.ClientKeeper.SetParams(ctx, clientParams)

		getParams := ibckeeper.ClientKeeper.GetParams(ctx)
		getParams.AllowedClients = append(getParams.AllowedClients, exported.Localhost)
		ibckeeper.ClientKeeper.SetParams(ctx, getParams)

		feePool, err := distkeeper.FeePool.Get(ctx)
		if err != nil {
			ctx.Logger().Info("app.DistrKeeper.FeePool.Get(ctx)", err)
			panic(err)
		}

		err = distkeeper.FeePool.Set(ctx, feePool)
		if err != nil {
			return nil, err
		}

		ctx.Logger().Info("Completed upgrade v2 successfully")
		return newVM, nil
	}
}

func updateVestingSchedules(ctx sdk.Context, k authkeeper.AccountKeeper, bankkeeper bankkeeper.Keeper, stakingkeeper stakingkeeper.Keeper, distrkeeper distrkeeper.Keeper, upgradeHeight int64) error {
	monthsToAdd := int64(3)
	sort.Strings(vestingAddresses)
	for _, addr := range vestingAddresses {
		if err := updateVestingAccount(ctx, k, addr, monthsToAdd, upgradeHeight); err != nil {
			return fmt.Errorf("failed to update vesting for %s: %w", addr, err)
		}
	}

	sort.SliceStable(addressReplacements, func(i, j int) bool {
		return addressReplacements[i].OldAddress < addressReplacements[j].OldAddress
	})
	ctx.Logger().Info("Processing address replacements", "count", len(addressReplacements))
	for _, replacement := range addressReplacements {
		newAddr, err := sdk.AccAddressFromBech32(replacement.NewAddress)
		if err != nil {
			// This indicates a problem with the hardcoded new address string itself.
			return fmt.Errorf("critical configuration error: invalid new address format %s: %w", replacement.NewAddress, err)
		}
		if k.GetAccount(ctx, newAddr) != nil {
			return fmt.Errorf("pre-check failed: new address %s (meant for old %s) already exists. Halting upgrade to prevent conflicts.", replacement.NewAddress, replacement.OldAddress)
		}

		if err := replaceAccountAddress(ctx, k, replacement.OldAddress, replacement.NewAddress, bankkeeper, stakingkeeper, distrkeeper, upgradeHeight); err != nil {
			return fmt.Errorf("failed to replace address for %s: %w", replacement.OldAddress, err)
		}
	}

	return nil
}

func updateVestingAccount(ctx sdk.Context, k authkeeper.AccountKeeper, address string, monthsToAdd int64, upgradeHeight int64) error {
	acc, err := getPeriodicVestingAccount(ctx, k, address)
	if err != nil {
		return err
	}

	// Use a deterministic reference time based on upgrade height instead of current block time
	// This ensures all validators use exactly the same time reference
	referenceTime := ctx.BlockHeader().Time.Unix()
	secondsToAdd := monthsToAdd * 2628000 // ~30.44 days per month

	ctx.Logger().Info("Account details",
		"address", address,
		"current_start_time", acc.StartTime,
		"current_end_time", acc.EndTime,
		"periods", len(acc.VestingPeriods),
		"reference_time", referenceTime,
		"upgrade_height", upgradeHeight)

	// Find the first unvested period
	cumulativeTime := acc.StartTime
	firstUnvestedIdx := 0
	for i, period := range acc.VestingPeriods {
		cumulativeTime += period.Length
		if cumulativeTime > referenceTime {
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
			"seconds_added", secondsToAdd,
			"upgrade_height", upgradeHeight)

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

func replaceAccountAddress(
	ctx sdk.Context,
	ak authkeeper.AccountKeeper,
	oldAddrStr, newAddrStr string,
	bk bankkeeper.Keeper,
	sk stakingkeeper.Keeper,
	dk distrkeeper.Keeper,
	upgradeHeight int64,
) error {
	// ---------- Validate addresses -----------------------------------------
	newAddr, err := sdk.AccAddressFromBech32(newAddrStr)
	if err != nil {
		return fmt.Errorf("invalid new address %s: %w", newAddrStr, err)
	}
	if ak.GetAccount(ctx, newAddr) != nil {
		return fmt.Errorf("new address %s already exists", newAddrStr)
	}

	oldAddr, err := sdk.AccAddressFromBech32(oldAddrStr)
	if err != nil {
		return fmt.Errorf("invalid old address %s: %w", oldAddrStr, err)
	}

	// Get the old account
	oldAcc := ak.GetAccount(ctx, oldAddr)
	if oldAcc == nil {
		return fmt.Errorf("account not found: %s", oldAddrStr)
	}

	// Log start of migration
	ctx.Logger().Info("Starting account migration",
		"old_address", oldAddrStr,
		"new_address", newAddrStr,
		"upgrade_height", upgradeHeight,
		"gas_remaining", ctx.GasMeter().Limit()-ctx.GasMeter().GasConsumed())

	var oldBaseAcc *authtypes.BaseAccount
	var accType string
	var accountHasVesting bool = false

	// Handle each vesting account type
	switch acc := oldAcc.(type) {
	case *vestingtypes.PeriodicVestingAccount:
		accType = "periodic"
		accountHasVesting = true
		oldBaseAcc = acc.BaseVestingAccount.BaseAccount

	case *vestingtypes.PermanentLockedAccount:
		accType = "permanent"
		accountHasVesting = true
		oldBaseAcc = acc.BaseVestingAccount.BaseAccount

	default:
		// For base accounts or any other type
		accType = "base"
		if baseAcc, ok := oldAcc.(*authtypes.BaseAccount); ok {
			oldBaseAcc = baseAcc
		} else {
			return fmt.Errorf("unsupported account type: %T", oldAcc)
		}
	}

	// ---------- Get withdraw address before any changes --------------------
	withdrawAddr, _ := dk.GetDelegatorWithdrawAddr(ctx, oldAddr)
	hasCustomWithdrawAddr := !withdrawAddr.Equals(oldAddr)

	/*
		There is no way to breach maxretrieve(65535) in our production,
		also the migration address are know and tracked which dont even have 100 delegations
		so no pagination retrival is needed here
	*/
	const maxRetrieve = uint16(1<<16 - 1) // 65535
	delegations, err := sk.GetDelegatorDelegations(ctx, oldAddr, maxRetrieve)
	if err != nil {
		return err
	}
	if len(delegations) == int(maxRetrieve) {
		return fmt.Errorf(
			"upgrade v2: delegations for %s may be truncated (got %d entries, limit %d)",
			oldAddrStr, len(delegations), maxRetrieve,
		)
	}

	// Store delegation info for later use
	type DelInfo struct {
		delegation stakingtypes.Delegation
		validator  stakingtypes.Validator
		tokens     sdk.DecCoin
		valAddr    sdk.ValAddress
	}

	var delsToMove []DelInfo
	var validatorAddressList []string

	// First pass - collect all delegations and validators
	for _, del := range delegations {
		valAddr, _ := sdk.ValAddressFromBech32(del.ValidatorAddress)
		val, err1 := sk.GetValidator(ctx, valAddr)
		found := true
		if err1 != nil {
			found = false
		}
		if !found {
			ctx.Logger().Info("Skipping delegation - validator not found",
				"validator", del.ValidatorAddress)
			continue
		}

		// Use deterministic token calculation with consistent precision
		tokens := sdk.NewDecCoinFromDec("uslf", val.TokensFromShares(del.Shares).TruncateDec())

		// Create a deterministic sort key

		delsToMove = append(delsToMove, DelInfo{
			delegation: del,
			validator:  val,
			tokens:     tokens,
			valAddr:    valAddr,
		})

		// Keep track of validator addresses in a deterministic way
		found = false
		for _, existingAddr := range validatorAddressList {
			if existingAddr == del.ValidatorAddress {
				found = true
				break
			}
		}
		if !found {
			validatorAddressList = append(validatorAddressList, del.ValidatorAddress)
		}
	}

	// Sort validator addresses for deterministic processing
	sort.Strings(validatorAddressList)

	// Sort delegations by the deterministic key for consistent processing
	sort.SliceStable(delsToMove, func(i, j int) bool {
		if delsToMove[i].delegation.ValidatorAddress != delsToMove[j].delegation.ValidatorAddress {
			return delsToMove[i].delegation.ValidatorAddress < delsToMove[j].delegation.ValidatorAddress
		}
		return delsToMove[i].delegation.Shares.LT(delsToMove[j].delegation.Shares)
	})

	// ---------- Create new base account ----------------------------------
	newBaseAcc := authtypes.NewBaseAccount(
		newAddr,
		nil, // This will be set later if needed
		ak.NextAccountNumber(ctx),
		0, // Start with sequence 0
	)

	// ---------- Create new vesting account (if applicable) ---------------
	var newAcc authtypes.AccountI

	if accountHasVesting {
		// Create the exact same type of vesting account with same parameters
		switch accType {
		case "periodic":
			oldPeriodicAcc := oldAcc.(*vestingtypes.PeriodicVestingAccount)

			// Create new periodic vesting account with SAME start time
			newPeriodicAcc, _ := vestingtypes.NewPeriodicVestingAccount(
				newBaseAcc,
				oldPeriodicAcc.OriginalVesting,
				oldPeriodicAcc.StartTime,
				oldPeriodicAcc.VestingPeriods,
			)

			// Keep the same end time
			newPeriodicAcc.EndTime = oldPeriodicAcc.EndTime

			// Set delegation amounts
			newPeriodicAcc.DelegatedFree = oldPeriodicAcc.DelegatedFree
			newPeriodicAcc.DelegatedVesting = oldPeriodicAcc.DelegatedVesting

			newAcc = newPeriodicAcc

		case "permanent":
			oldPermAcc := oldAcc.(*vestingtypes.PermanentLockedAccount)

			// Create new permanent locked account
			newPermAcc, _ := vestingtypes.NewPermanentLockedAccount(
				newBaseAcc,
				oldPermAcc.OriginalVesting,
			)

			// Set delegation amounts
			newPermAcc.DelegatedFree = oldPermAcc.DelegatedFree
			newPermAcc.DelegatedVesting = oldPermAcc.DelegatedVesting

			newAcc = newPermAcc
		}
	} else {
		// For non-vesting accounts, just use a base account
		newAcc = newBaseAcc
	}

	// ---------- Save the new account first --------------------------------------
	// Create the new account before modifying the old one to reduce risk of data loss
	ak.SetAccount(ctx, newAcc)
	ctx.Logger().Info("Created new account",
		"address", newAddrStr,
		"type", accType,
		"upgrade_height", upgradeHeight,
		"gas_remaining", ctx.GasMeter().Limit()-ctx.GasMeter().GasConsumed())

	// ---------- Convert old account to base for transfers -----------------
	baseAcc := authtypes.NewBaseAccount(
		oldAddr,
		oldBaseAcc.GetPubKey(),
		oldBaseAcc.GetAccountNumber(),
		oldBaseAcc.GetSequence(),
	)
	ak.SetAccount(ctx, baseAcc)

	// ---------- Set withdraw address for new account if needed ----------
	if hasCustomWithdrawAddr {
		dk.SetDelegatorWithdrawAddr(ctx, newAddr, withdrawAddr)
		ctx.Logger().Info("Set withdraw address",
			"delegator", newAddrStr,
			"withdraw_addr", withdrawAddr.String(),
			"upgrade_height", upgradeHeight)
	}

	// ---------- Store validator reward information before any changes -----
	// Store current validator periods and historicals for proper migration
	type ValidatorPeriodInfo struct {
		validatorAddr string
		period        uint64
	}

	var validatorPeriods []ValidatorPeriodInfo

	// Collect all validator periods first before making any changes
	// Use the sorted validator address list for deterministic processing
	for _, valAddrStr := range validatorAddressList {
		valAddr, _ := sdk.ValAddressFromBech32(valAddrStr)

		// Store current period for each validator
		valCurrentRewards, _ := dk.GetValidatorCurrentRewards(ctx, valAddr)
		validatorPeriods = append(validatorPeriods, ValidatorPeriodInfo{
			validatorAddr: valAddrStr,
			period:        valCurrentRewards.Period,
		})

		ctx.Logger().Info("Captured validator reward state",
			"validator", valAddrStr,
			"current_period", valCurrentRewards.Period,
			"upgrade_height", upgradeHeight)
	}

	// Sort validator periods for deterministic lookups
	sort.SliceStable(validatorPeriods, func(i, j int) bool {
		return validatorPeriods[i].validatorAddr < validatorPeriods[j].validatorAddr
	})

	// ---------- First, handle existing rewards ---------------------------
	totalRewardsWithdrawn := sdk.NewCoins()

	// Force withdrawal of any existing rewards - this is important to "reset"
	// the reward state and prevent double-counting
	for _, delInfo := range delsToMove {
		valAddr, _ := sdk.ValAddressFromBech32(delInfo.delegation.ValidatorAddress)

		// Withdraw rewards
		rewards, err := dk.WithdrawDelegationRewards(ctx, oldAddr, valAddr)
		if err != nil {
			// Log error but continue - we don't want to fail the entire migration for one reward issue
			ctx.Logger().Error("Failed to withdraw rewards",
				"validator", delInfo.delegation.ValidatorAddress,
				"error", err.Error(),
				"upgrade_height", upgradeHeight)
		} else if !rewards.IsZero() {
			totalRewardsWithdrawn = totalRewardsWithdrawn.Add(rewards...)
			ctx.Logger().Info("Withdrawn rewards",
				"validator", delInfo.delegation.ValidatorAddress,
				"amount", rewards.String(),
				"upgrade_height", upgradeHeight)
		}
	}

	// ---------- Transfer all balances ------------------------------------
	// Transfer the full balance, including any withdrawn rewards
	allCoins := bk.GetAllBalances(ctx, oldAddr)
	if !allCoins.IsZero() {
		if err := bk.SendCoins(ctx, oldAddr, newAddr, allCoins); err != nil {
			return fmt.Errorf("failed to transfer balances: %w", err)
		}
		ctx.Logger().Info("Transferred balance",
			"amount", allCoins.String(),
			"from", oldAddrStr,
			"to", newAddrStr,
			"upgrade_height", upgradeHeight,
			"gas_remaining", ctx.GasMeter().Limit()-ctx.GasMeter().GasConsumed())
	}

	// ---------- Handle unbonding delegations -------------------------------
	// Get and sort unbonding delegations for deterministic processing
	unbondingDels, err := sk.GetUnbondingDelegations(ctx, oldAddr, maxRetrieve)
	if err != nil {
		return err
	}

	// Sort unbonding delegations by validator address for deterministic processing
	sort.SliceStable(unbondingDels, func(i, j int) bool {
		return unbondingDels[i].ValidatorAddress < unbondingDels[j].ValidatorAddress
	})

	// Create new unbonding delegations before removing old ones
	for _, ubd := range unbondingDels {
		// Create new unbonding delegation first
		valAddr, _ := sdk.ValAddressFromBech32(ubd.ValidatorAddress)
		newUBD := stakingtypes.UnbondingDelegation{
			DelegatorAddress: newAddr.String(),
			ValidatorAddress: valAddr.String(),
			Entries:          ubd.Entries,
		}
		sk.SetUnbondingDelegation(ctx, newUBD)

		// Then remove old unbonding delegation
		sk.RemoveUnbondingDelegation(ctx, ubd)

		ctx.Logger().Info("Moved unbonding delegation",
			"validator", ubd.ValidatorAddress,
			"entries", len(ubd.Entries),
			"upgrade_height", upgradeHeight)
	}

	// ---------- Handle redelegations --------------------------------------
	// Get and sort redelegations for deterministic processing
	redelegations, err := sk.GetRedelegations(ctx, oldAddr, maxRetrieve)
	if err != nil {
		return err
	}

	// Sort redelegations by source and destination validator addresses
	sort.SliceStable(redelegations, func(i, j int) bool {
		if redelegations[i].ValidatorSrcAddress == redelegations[j].ValidatorSrcAddress {
			return redelegations[i].ValidatorDstAddress < redelegations[j].ValidatorDstAddress
		}
		return redelegations[i].ValidatorSrcAddress < redelegations[j].ValidatorSrcAddress
	})

	// Create new redelegations before removing old ones
	for _, red := range redelegations {
		// Create new redelegation first
		srcVal, _ := sdk.ValAddressFromBech32(red.ValidatorSrcAddress)
		dstVal, _ := sdk.ValAddressFromBech32(red.ValidatorDstAddress)
		newRed := stakingtypes.Redelegation{
			DelegatorAddress:    newAddr.String(),
			ValidatorSrcAddress: srcVal.String(),
			ValidatorDstAddress: dstVal.String(),
			Entries:             red.Entries,
		}
		sk.SetRedelegation(ctx, newRed)

		// Then remove old redelegation
		sk.RemoveRedelegation(ctx, red)

		ctx.Logger().Info("Moved redelegation",
			"src_validator", red.ValidatorSrcAddress,
			"dst_validator", red.ValidatorDstAddress,
			"entries", len(red.Entries),
			"upgrade_height", upgradeHeight,
			"gas_remaining", ctx.GasMeter().Limit()-ctx.GasMeter().GasConsumed())
	}

	// ---------- Move delegations and set up reward state properly ---------
	// Need to update validator for each delegation
	for _, delInfo := range delsToMove {
		del := delInfo.delegation
		valAddr, _ := sdk.ValAddressFromBech32(del.ValidatorAddress)

		// Create new delegation with the same shares first
		newDel := stakingtypes.NewDelegation(
			newAddr.String(),
			valAddr.String(),
			del.Shares,
		)
		sk.SetDelegation(ctx, newDel)

		// CRITICAL FIX: Set up proper reward state for new delegation
		// We need the current period from before the migration
		// Find the period using binary search on the sorted validator periods
		var currentPeriod uint64

		// Binary search for the validator period
		left, right := 0, len(validatorPeriods)-1
		found := false
		for left <= right {
			mid := (left + right) / 2
			if validatorPeriods[mid].validatorAddr == del.ValidatorAddress {
				currentPeriod = validatorPeriods[mid].period
				found = true
				break
			} else if validatorPeriods[mid].validatorAddr < del.ValidatorAddress {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}

		// Fallback to linear search if binary search fails
		if !found {
			for _, vp := range validatorPeriods {
				if vp.validatorAddr == del.ValidatorAddress {
					currentPeriod = vp.period
					break
				}
			}
		}

		// Create starting info with current validator period
		// This is critical to ensure rewards accrue correctly
		startInfo := distrtypes.NewDelegatorStartingInfo(
			currentPeriod,         // Use current period for proper tracking
			delInfo.tokens.Amount, // Current token value of shares
			uint64(upgradeHeight), // Use upgrade height instead of current block height
		)

		// Set the starting info for the new delegator
		dk.SetDelegatorStartingInfo(ctx, valAddr, newAddr, startInfo)

		// Then remove old delegation (this cleans up distribution state too)
		sk.RemoveDelegation(ctx, del)

		ctx.Logger().Info("Set up rewards for new delegation",
			"validator", del.ValidatorAddress,
			"shares", del.Shares.String(),
			"tokens", delInfo.tokens.String(),
			"current_period", currentPeriod,
			"upgrade_height", upgradeHeight)
	}

	// ---------- Final step: Force a rewards claim to initialize properly -----
	// This ensures the delegator starts with a clean slate for future rewards
	for _, delInfo := range delsToMove {
		valAddr, _ := sdk.ValAddressFromBech32(delInfo.delegation.ValidatorAddress)

		// Withdraw any rewards that might have accrued during migration
		// This ensures a completely clean starting point
		rewards, err := dk.WithdrawDelegationRewards(ctx, newAddr, valAddr)
		if err != nil {
			ctx.Logger().Error("Failed to initialize rewards state - non-critical",
				"validator", delInfo.delegation.ValidatorAddress,
				"error", err.Error(),
				"upgrade_height", upgradeHeight)
		} else if !rewards.IsZero() {
			ctx.Logger().Info("Initialized rewards state with withdrawal",
				"validator", delInfo.delegation.ValidatorAddress,
				"amount", rewards.String(),
				"upgrade_height", upgradeHeight)
		}
	}

	// Log completion
	ctx.Logger().Info("Account migration complete",
		"old_address", oldAddrStr,
		"new_address", newAddrStr,
		"account_type", accType,
		"delegations", len(delsToMove),
		"unbonding_delegations", len(unbondingDels),
		"redelegations", len(redelegations),
		"rewards_withdrawn", totalRewardsWithdrawn,
		"upgrade_height", upgradeHeight,
		"final_gas_remaining", ctx.GasMeter().Limit()-ctx.GasMeter().GasConsumed())

	return nil
}
