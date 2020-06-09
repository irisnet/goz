package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	codecstd "github.com/cosmos/cosmos-sdk/std"
	gaia "github.com/cosmos/gaia/app"

	rly "github.com/iqlusioninc/relayer/cmd"

	"github.com/irisnet/hack-goz/types"
)

var (
	cfgPath    string
	homePath   string
	debug      bool
	config     *rly.Config
	cdc        *codec.Codec
	appCodec   *codecstd.Codec
	flagConfig = "config"

	// Default identifiers for dummy usage
	// dcli = "defaultclientid"
	// dcon = "defaultconnectionid"
	// dcha = "defaultchannelid"
	// dpor = "defaultportid"
)

var rootCmd = &cobra.Command{
	Use:   "hackcli",
	Short: "Command line interface for hacking game of zones",
}

func init() {
	// Register top level flags --home and --config
	rootCmd.PersistentFlags().StringVar(&homePath, flags.FlagHome, types.DefaultCLIHome, "set home directory")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug output")
	rootCmd.PersistentFlags().StringVar(&cfgPath, flagConfig, "config.yaml", "set config file")
	if err := viper.BindPFlag(flags.FlagHome, rootCmd.Flags().Lookup(flags.FlagHome)); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag(flagConfig, rootCmd.Flags().Lookup(flagConfig)); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("debug", rootCmd.Flags().Lookup("debug")); err != nil {
		panic(err)
	}

	// Construct Root Command
	rootCmd.AddCommand(
		KeysCmd(),
		ConfigCmd(),
		liteCmd(),
		FaucetCmd(),
		TxCmd(),
		AutoTxCmd(),
		queryCmd(),
	)

	appCodec, cdc = gaia.MakeCodecs()
}

func Execute() {

	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		// reads `homeDir/config/config.yaml` into `var config *Config` before each command
		return initConfig(rootCmd)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
