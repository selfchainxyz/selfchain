package app

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// UpgradeName defines the on-chain upgrade name for this chain upgrade from v001 to v002.
const UpgradeName = "v001-to-v002"

// Naturally, this Handler is application specific and not defined on a per-module basis.
func (app App) RegisterUpgradeHandlers() {
  app.UpgradeKeeper.SetUpgradeHandler(
    UpgradeName,
    /// This UpgradeHandler functios retrievse the VersionMap from x/upgrade's state and return the new VersionMap
    /// to be stored in x/upgrade after the upgrade. The diff between the two VersionMaps determines
    /// which modules need upgrading.
    ///
    /// Inside this function, you must perform any upgrade logic to include in the provided plan.
    /// All upgrade handler functions must end with the following line of code  `return app.mm.RunMigrations(ctx, cfg, fromVM)`
    ///
    /// New modules are recognized because they have not yet been registered in x/upgrade's VersionMap store. In this case, RunMigrations
    /// calls the InitGenesis function from the corresponding module to set up its initial state.
    func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
      // the new duration is 2 minutes
      voting_period, _ := time.ParseDuration("2m")

      // Here we can change the params of the gov module
      app.GovKeeper.SetVotingParams(ctx, v1.VotingParams{
        VotingPeriod: &voting_period,
      })

      // This will cause the module manager to run the store migration of the modules that are being upgraded.
      return app.mm.RunMigrations(ctx, app.configurator, fromVM)
    },
  )

  /// // All chains preparing to run in-place store migrations will need to manually add store upgrades for new modules
  /// // and then configure the store loader to apply those upgrades. This ensures that the new module's stores are added
  /// // to the multistore before the migrations begin.
  /// upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk(); if err != nil {
  /// 	panic(err)
  /// }

  /// if upgradeInfo.Name == UpgradeName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
  /// 	storeUpgrades := storetypes.StoreUpgrades{
  /// 		// add store upgrades for new modules
  /// 		// Example:
  /// 		//    Added: []string{"foo", "bar"},
  /// 		// ...
  /// 	}

  /// 	// configure store loader that checks if version == upgradeHeight and applies store upgrades
  /// 	app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
  /// }
}
