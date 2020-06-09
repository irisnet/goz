/*
Copyright © 2020 Jack Zampolin <jack.zampolin@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	neturl "net/url"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	lite "github.com/tendermint/tendermint/lite2"

	"github.com/cosmos/cosmos-sdk/client/flags"
)

var (
	flagHash  = "hash"
	flagURL   = "url"
	flagForce = "force"
)

// chainCmd represents the keys command
func liteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "lite",
		Aliases: []string{"l"},
		Short:   "manage lite clients held by the relayer for each chain",
	}

	cmd.AddCommand(
		initLiteCmd(),
	)

	return cmd
}

func initLiteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init [chain-id]",
		Aliases: []string{"i"},
		Short:   "Initiate the light client",
		Long: `Initiate the light client by:
	1. passing it a root of trust as a --hash/-x and --height
	2. via --url/-u where trust options can be found
	3. Use --force/-f to initalize from the configured node`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			chain, err := config.Chains.Get(args[0])
			if err != nil {
				return err
			}

			db, df, err := chain.NewLiteDB()
			if err != nil {
				return err
			}
			defer df()

			url, err := cmd.Flags().GetString(flagURL)
			if err != nil {
				return err
			}
			force, err := cmd.Flags().GetBool(flagForce)
			if err != nil {
				return err
			}
			height, err := cmd.Flags().GetInt64(flags.FlagHeight)
			if err != nil {
				return err
			}
			hash, err := cmd.Flags().GetBytesHex(flagHash)
			if err != nil {
				return err
			}

			switch {
			case force: // force initialization from trusted node
				if _, err = chain.TrustNodeInitClient(db); err != nil {
					return err
				}
			case height > 0 && len(hash) > 0: // height and hash are given
				if _, err = chain.InitLiteClient(db, chain.TrustOptions(height, hash)); err != nil {
					return wrapInitFailed(err)
				}
			case len(url) > 0: // URL is given, query trust options
				if _, err := neturl.Parse(url); err != nil {
					return wrapIncorrectURL(err)
				}

				to, err := queryTrustOptions(url)
				if err != nil {
					return err
				}

				if _, err = chain.InitLiteClient(db, to); err != nil {
					return wrapInitFailed(err)
				}
			default: // return error
				return errInitWrongFlags
			}

			return nil
		},
	}

	return forceFlag(liteFlags(cmd))
}

func queryTrustOptions(url string) (out lite.TrustOptions, err error) {
	// fetch from URL
	res, err := http.Get(url)
	if err != nil {
		return
	}

	// read in the res body
	bz, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	// close the response body
	if err = res.Body.Close(); err != nil {
		return
	}

	// unmarshal the data into the trust options hash
	if err = json.Unmarshal(bz, &out); err != nil {
		return
	}

	return
}

func forceFlag(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().BoolP(flagForce, "f", false, "option to force initialization of lite client from configured chain")
	if err := viper.BindPFlag(flagForce, cmd.Flags().Lookup(flagForce)); err != nil {
		panic(err)
	}
	return cmd
}

func liteFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().Int64(flags.FlagHeight, -1, "Trusted header's height")
	cmd.Flags().BytesHexP(flagHash, "x", []byte{}, "Trusted header's hash")
	cmd.Flags().StringP(flagURL, "u", "", "Optional URL to fetch trusted-hash and trusted-height")
	if err := viper.BindPFlag(flags.FlagHeight, cmd.Flags().Lookup(flags.FlagHeight)); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag(flagHash, cmd.Flags().Lookup(flagHash)); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag(flagURL, cmd.Flags().Lookup(flagURL)); err != nil {
		panic(err)
	}
	return cmd
}
