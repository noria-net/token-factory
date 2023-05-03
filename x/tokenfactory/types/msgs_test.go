package types_test

import (
	fmt "fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/noria-net/token-factory/x/tokenfactory/types"

	"github.com/cometbft/cometbft/crypto/ed25519"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// TestMsgCreateDenom tests if valid/invalid create denom messages are properly validated/invalidated
func TestMsgCreateDenom(t *testing.T) {
	// generate a private/public key pair and get the respective address
	pk1 := ed25519.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pk1.Address())

	// make a proper createDenom message
	createMsg := func(after func(msg types.MsgTokenFactoryCreateDenom) types.MsgTokenFactoryCreateDenom) types.MsgTokenFactoryCreateDenom {
		properMsg := *types.NewMsgCreateDenom(
			addr1.String(),
			"bitcoin",
		)

		return after(properMsg)
	}

	// validate createDenom message was created as intended
	msg := createMsg(func(msg types.MsgTokenFactoryCreateDenom) types.MsgTokenFactoryCreateDenom {
		return msg
	})
	require.Equal(t, msg.Route(), types.RouterKey)
	require.Equal(t, msg.Type(), "create_denom")
	signers := msg.GetSigners()
	require.Equal(t, len(signers), 1)
	require.Equal(t, signers[0].String(), addr1.String())

	tests := []struct {
		name       string
		msg        types.MsgTokenFactoryCreateDenom
		expectPass bool
	}{
		{
			name: "proper msg",
			msg: createMsg(func(msg types.MsgTokenFactoryCreateDenom) types.MsgTokenFactoryCreateDenom {
				return msg
			}),
			expectPass: true,
		},
		{
			name: "empty sender",
			msg: createMsg(func(msg types.MsgTokenFactoryCreateDenom) types.MsgTokenFactoryCreateDenom {
				msg.Sender = ""
				return msg
			}),
			expectPass: false,
		},
		{
			name: "invalid subdenom",
			msg: createMsg(func(msg types.MsgTokenFactoryCreateDenom) types.MsgTokenFactoryCreateDenom {
				msg.Subdenom = "thissubdenomismuchtoolongasdkfjaasdfdsafsdlkfnmlksadmflksmdlfmlsakmfdsafasdfasdf"
				return msg
			}),
			expectPass: false,
		},
	}

	for _, test := range tests {
		if test.expectPass {
			require.NoError(t, test.msg.ValidateBasic(), "test: %v", test.name)
		} else {
			require.Error(t, test.msg.ValidateBasic(), "test: %v", test.name)
		}
	}
}

// TestMsgMint tests if valid/invalid create denom messages are properly validated/invalidated
func TestMsgMint(t *testing.T) {
	// generate a private/public key pair and get the respective address
	pk1 := ed25519.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pk1.Address())

	// make a proper mint message
	createMsg := func(after func(msg types.MsgTokenFactoryMint) types.MsgTokenFactoryMint) types.MsgTokenFactoryMint {
		properMsg := *types.NewMsgMint(
			addr1.String(),
			sdk.NewCoin("bitcoin", sdk.NewInt(500000000)),
		)

		return after(properMsg)
	}

	// validate mint message was created as intended
	msg := createMsg(func(msg types.MsgTokenFactoryMint) types.MsgTokenFactoryMint {
		return msg
	})
	require.Equal(t, msg.Route(), types.RouterKey)
	require.Equal(t, msg.Type(), "tf_mint")
	signers := msg.GetSigners()
	require.Equal(t, len(signers), 1)
	require.Equal(t, signers[0].String(), addr1.String())

	tests := []struct {
		name       string
		msg        types.MsgTokenFactoryMint
		expectPass bool
	}{
		{
			name: "proper msg",
			msg: createMsg(func(msg types.MsgTokenFactoryMint) types.MsgTokenFactoryMint {
				return msg
			}),
			expectPass: true,
		},
		{
			name: "empty sender",
			msg: createMsg(func(msg types.MsgTokenFactoryMint) types.MsgTokenFactoryMint {
				msg.Sender = ""
				return msg
			}),
			expectPass: false,
		},
		{
			name: "zero amount",
			msg: createMsg(func(msg types.MsgTokenFactoryMint) types.MsgTokenFactoryMint {
				msg.Amount = sdk.NewCoin("bitcoin", sdk.ZeroInt())
				return msg
			}),
			expectPass: false,
		},
		{
			name: "negative amount",
			msg: createMsg(func(msg types.MsgTokenFactoryMint) types.MsgTokenFactoryMint {
				msg.Amount.Amount = sdk.NewInt(-10000000)
				return msg
			}),
			expectPass: false,
		},
	}

	for _, test := range tests {
		if test.expectPass {
			require.NoError(t, test.msg.ValidateBasic(), "test: %v", test.name)
		} else {
			require.Error(t, test.msg.ValidateBasic(), "test: %v", test.name)
		}
	}
}

// TestMsgBurn tests if valid/invalid create denom messages are properly validated/invalidated
func TestMsgBurn(t *testing.T) {
	// generate a private/public key pair and get the respective address
	pk1 := ed25519.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pk1.Address())

	// make a proper burn message
	baseMsg := types.NewMsgBurn(
		addr1.String(),
		sdk.NewCoin("bitcoin", sdk.NewInt(500000000)),
	)

	// validate burn message was created as intended
	require.Equal(t, baseMsg.Route(), types.RouterKey)
	require.Equal(t, baseMsg.Type(), "tf_burn")
	signers := baseMsg.GetSigners()
	require.Equal(t, len(signers), 1)
	require.Equal(t, signers[0].String(), addr1.String())

	tests := []struct {
		name       string
		msg        func() *types.MsgTokenFactoryBurn
		expectPass bool
	}{
		{
			name: "proper msg",
			msg: func() *types.MsgTokenFactoryBurn {
				msg := baseMsg
				return msg
			},
			expectPass: true,
		},
		{
			name: "empty sender",
			msg: func() *types.MsgTokenFactoryBurn {
				msg := baseMsg
				msg.Sender = ""
				return msg
			},
			expectPass: false,
		},
		{
			name: "zero amount",
			msg: func() *types.MsgTokenFactoryBurn {
				msg := baseMsg
				msg.Amount.Amount = sdk.ZeroInt()
				return msg
			},
			expectPass: false,
		},
		{
			name: "negative amount",
			msg: func() *types.MsgTokenFactoryBurn {
				msg := baseMsg
				msg.Amount.Amount = sdk.NewInt(-10000000)
				return msg
			},
			expectPass: false,
		},
	}

	for _, test := range tests {
		if test.expectPass {
			require.NoError(t, test.msg().ValidateBasic(), "test: %v", test.name)
		} else {
			require.Error(t, test.msg().ValidateBasic(), "test: %v", test.name)
		}
	}
}

// TestMsgChangeAdmin tests if valid/invalid create denom messages are properly validated/invalidated
func TestMsgChangeAdmin(t *testing.T) {
	// generate a private/public key pair and get the respective address
	pk1 := ed25519.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pk1.Address())
	pk2 := ed25519.GenPrivKey().PubKey()
	addr2 := sdk.AccAddress(pk2.Address())
	tokenFactoryDenom := fmt.Sprintf("factory/%s/bitcoin", addr1.String())

	// make a proper changeAdmin message
	baseMsg := types.NewMsgChangeAdmin(
		addr1.String(),
		tokenFactoryDenom,
		addr2.String(),
	)

	// validate changeAdmin message was created as intended
	require.Equal(t, baseMsg.Route(), types.RouterKey)
	require.Equal(t, baseMsg.Type(), "change_admin")
	signers := baseMsg.GetSigners()
	require.Equal(t, len(signers), 1)
	require.Equal(t, signers[0].String(), addr1.String())

	tests := []struct {
		name       string
		msg        func() *types.MsgTokenFactoryChangeAdmin
		expectPass bool
	}{
		{
			name: "proper msg",
			msg: func() *types.MsgTokenFactoryChangeAdmin {
				msg := baseMsg
				return msg
			},
			expectPass: true,
		},
		{
			name: "empty sender",
			msg: func() *types.MsgTokenFactoryChangeAdmin {
				msg := baseMsg
				msg.Sender = ""
				return msg
			},
			expectPass: false,
		},
		{
			name: "empty newAdmin",
			msg: func() *types.MsgTokenFactoryChangeAdmin {
				msg := baseMsg
				msg.NewAdmin = ""
				return msg
			},
			expectPass: false,
		},
		{
			name: "invalid denom",
			msg: func() *types.MsgTokenFactoryChangeAdmin {
				msg := baseMsg
				msg.Denom = "bitcoin"
				return msg
			},
			expectPass: false,
		},
	}

	for _, test := range tests {
		if test.expectPass {
			require.NoError(t, test.msg().ValidateBasic(), "test: %v", test.name)
		} else {
			require.Error(t, test.msg().ValidateBasic(), "test: %v", test.name)
		}
	}
}

// TestMsgSetDenomMetadata tests if valid/invalid create denom messages are properly validated/invalidated
func TestMsgSetDenomMetadata(t *testing.T) {
	// generate a private/public key pair and get the respective address
	pk1 := ed25519.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pk1.Address())
	tokenFactoryDenom := fmt.Sprintf("factory/%s/bitcoin", addr1.String())
	denomMetadata := banktypes.Metadata{
		Description: "nakamoto",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    tokenFactoryDenom,
				Exponent: 0,
			},
			{
				Denom:    "sats",
				Exponent: 6,
			},
		},
		Display: "sats",
		Base:    tokenFactoryDenom,
		Name:    "bitcoin",
		Symbol:  "BTC",
	}
	invalidDenomMetadata := banktypes.Metadata{
		Description: "nakamoto",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "bitcoin",
				Exponent: 0,
			},
			{
				Denom:    "sats",
				Exponent: 6,
			},
		},
		Display: "sats",
		Base:    "bitcoin",
		Name:    "bitcoin",
		Symbol:  "BTC",
	}

	// make a proper setDenomMetadata message
	baseMsg := types.NewMsgSetDenomMetadata(
		addr1.String(),
		denomMetadata,
	)

	// validate setDenomMetadata message was created as intended
	require.Equal(t, baseMsg.Route(), types.RouterKey)
	require.Equal(t, baseMsg.Type(), "set_denom_metadata")
	signers := baseMsg.GetSigners()
	require.Equal(t, len(signers), 1)
	require.Equal(t, signers[0].String(), addr1.String())

	tests := []struct {
		name       string
		msg        func() *types.MsgTokenFactorySetDenomMetadata
		expectPass bool
	}{
		{
			name: "proper msg",
			msg: func() *types.MsgTokenFactorySetDenomMetadata {
				msg := baseMsg
				return msg
			},
			expectPass: true,
		},
		{
			name: "empty sender",
			msg: func() *types.MsgTokenFactorySetDenomMetadata {
				msg := baseMsg
				msg.Sender = ""
				return msg
			},
			expectPass: false,
		},
		{
			name: "invalid metadata",
			msg: func() *types.MsgTokenFactorySetDenomMetadata {
				msg := baseMsg
				msg.Metadata.Name = ""
				return msg
			},

			expectPass: false,
		},
		{
			name: "invalid base",
			msg: func() *types.MsgTokenFactorySetDenomMetadata {
				msg := baseMsg
				msg.Metadata = invalidDenomMetadata
				return msg
			},
			expectPass: false,
		},
	}

	for _, test := range tests {
		if test.expectPass {
			require.NoError(t, test.msg().ValidateBasic(), "test: %v", test.name)
		} else {
			require.Error(t, test.msg().ValidateBasic(), "test: %v", test.name)
		}
	}
}
