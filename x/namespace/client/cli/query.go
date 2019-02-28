package cli

import (
	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"fmt"
	"strings"
)

func GetCmdResolveName(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: "resolve [name]",
		Short: "resolve name",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := args[0]

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/resolve/%s", queryRoute, name), nil)
			if err != nil {
				fmt.Printf("could not resolve name - %s \n", string(name))
				return nil
			}

			return cliCtx.PrintOutput(resolveRes{string(res)})
		},
	}
}

type resolveRes struct {
	Value string `json:"value"`
}

func (r resolveRes) String() string {
	return r.Value
}

func GetCmdWhois(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: "whois [name]",
		Short: "Query whois info of name",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := args[0]

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/whois/%s", queryRoute, name), nil)
			if err != nil {
				fmt.Printf("could not resolve whois - %s \n", string(name))
				return nil
			}

			var out whoIsRes
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
type whoIsRes struct {
	Value string `json:"value"`
	Owner sdk.AccAddress `json:"owner"`
	Price sdk.Coins `json:"price"`
}

func (w whoIsRes) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
Value: %s
Price: %s`, w.Owner, w.Value, w.Price))
}
