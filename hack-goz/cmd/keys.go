package cmd

/*
Copyright © 2020 Yelong Zhang yelong@bianjie.ai

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

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"

	"github.com/iqlusioninc/relayer/relayer"
)

const (
	mnemonicEntropySize = 256
)

type keyOutput struct {
	Mnemonic string `json:"mnemonic" yaml:"mnemonic"`
	Address  string `json:"address" yaml:"address"`
}

// KeysCmd represents the keys command
func KeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "keys",
		Aliases: []string{"k"},
		Short:   "helps users manage keys for multiple chains",
	}

	cmd.AddCommand(
		keysChainAddCmd(),
		keysChainRestoreCmd(),
		keysChainShowCmd(),
		keysAddCmd(),
		keysDeleteCmd(),
		keysListCmd(),
	)

	return cmd
}

func keysChainAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "chain-add [chain-id] [[name]]",
		Aliases: []string{"a"},
		Short:   "adds a key to the keychain associated with a particular chain",
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			chain, err := config.Chains.Get(args[0])
			if err != nil {
				return err
			}

			var keyName string
			if len(args) == 2 {
				keyName = args[1]
			} else {
				keyName = chain.Key
			}

			if chain.KeyExists(keyName) {
				return errKeyExists(keyName)
			}

			mnemonic, err := relayer.CreateMnemonic()
			if err != nil {
				return err
			}

			info, err := chain.Keybase.NewAccount(keyName, mnemonic, "", hd.CreateHDPath(118, 0, 0).String(), hd.Secp256k1)
			if err != nil {
				return err
			}

			ko := keyOutput{Mnemonic: mnemonic, Address: info.GetAddress().String()}

			return chain.Print(ko, false, false)
		},
	}

	return cmd
}

// keysChainRestoreCmd respresents the `keys add` command
func keysChainRestoreCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "chain-restore [chain-id] [name] [mnemonic]",
		Aliases: []string{"r"},
		Short:   "restores a mnemonic to the keychain associated with a particular chain",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			keyName := args[1]
			chain, err := config.Chains.Get(args[0])
			if err != nil {
				return err
			}

			done := chain.UseSDKContext()
			defer done()

			if chain.KeyExists(keyName) {
				return errKeyExists(keyName)
			}

			info, err := chain.Keybase.NewAccount(keyName, args[2], "", hd.CreateHDPath(118, 0, 0).String(), hd.Secp256k1)
			if err != nil {
				return err
			}

			fmt.Println(info.GetAddress().String())
			return nil
		},
	}

	return cmd
}

func keysChainShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "chain-show [chain-id] [[name]]",
		Aliases: []string{"s"},
		Short:   "shows a key from the keychain associated with a particular chain",
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			chain, err := config.Chains.Get(args[0])
			if err != nil {
				return err
			}

			var keyName string
			if len(args) == 2 {
				keyName = args[1]
			} else {
				keyName = chain.Key
			}

			if !chain.KeyExists(keyName) {
				return errKeyDoesntExist(keyName)
			}

			info, err := chain.Keybase.Key(keyName)
			if err != nil {
				return err
			}

			fmt.Println(info.GetAddress().String())
			return nil
		},
	}

	return cmd
}

// keysAddCmd respresents the `keys add` command
func keysAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add [prefix] [number]",
		Aliases: []string{"a"},
		Short:   "adds keys to the keychain",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			prefix := args[0]
			num, err := strconv.ParseInt(args[1], 10, 32)
			if err != nil {
				return err
			}

			// read entropy seed straight from crypto.Rand and convert to mnemonic
			entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
			if err != nil {
				return err
			}

			mnemonic, err := bip39.NewMnemonic(entropySeed)
			if err != nil {
				return err
			}

			kb, err := keyring.New(prefix, keyring.BackendTest, homePath, nil)
			if err != nil {
				return err
			}

			for i := 0; i < int(num); i++ {
				_, err := kb.NewAccount(fmt.Sprintf("%s%d", prefix, i), mnemonic, "", hd.CreateHDPath(sdk.CoinType, 0, uint32(i)).String(), hd.Secp256k1)
				if err != nil {
					return err
				}
			}

			fmt.Println(mnemonic)
			return nil
		},
	}

	return cmd
}

// keysDeleteCmd respresents the `keys delete` command
func keysDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete [prefix]",
		Aliases: []string{"d"},
		Short:   "deletes a key from the keychain associated with a particular chain",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			prefix := args[0]
			kb, err := keyring.New(prefix, keyring.BackendTest, homePath, nil)
			if err != nil {
				return err
			}

			infos, err := kb.List()
			if err != nil {
				return err
			}
			for d, i := range infos {
				_ = kb.Delete(i.GetName())
				fmt.Printf("key(%d): %s -> %s deleted!\n", d, i.GetName(), i.GetAddress().String())
			}

			return nil
		},
	}

	return cmd
}

// keysListCmd respresents the `keys list` command
func keysListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list [prefix]",
		Aliases: []string{"l"},
		Short:   "lists keys of a specified prefix",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			prefix := args[0]
			kb, err := keyring.New(prefix, keyring.BackendTest, homePath, nil)
			if err != nil {
				return err
			}

			infos, err := kb.List()
			if err != nil {
				return err
			}
			for d, i := range infos {
				fmt.Printf("key(%d): %s -> %s\n", d, i.GetName(), i.GetAddress().String())
			}

			return nil
		},
	}

	return cmd
}
