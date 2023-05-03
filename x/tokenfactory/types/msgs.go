package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// constants
const (
	TypeMsgCreateDenom      = "create_denom"
	TypeMsgMint             = "tf_mint"
	TypeMsgBurn             = "tf_burn"
	TypeMsgForceTransfer    = "force_transfer"
	TypeMsgChangeAdmin      = "change_admin"
	TypeMsgSetDenomMetadata = "set_denom_metadata"
)

// NewMsgCreateDenom creates a msg to create a new denom
func NewMsgCreateDenom(sender, subdenom string) *MsgTokenFactoryCreateDenom {
	return &MsgTokenFactoryCreateDenom{
		Sender:   sender,
		Subdenom: subdenom,
	}
}

func (m MsgTokenFactoryCreateDenom) Route() string { return RouterKey }
func (m MsgTokenFactoryCreateDenom) Type() string  { return TypeMsgCreateDenom }
func (m MsgTokenFactoryCreateDenom) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	_, err = GetTokenDenom(m.Sender, m.Subdenom)
	if err != nil {
		return sdkerrors.Wrap(ErrInvalidDenom, err.Error())
	}

	return nil
}

func (m MsgTokenFactoryCreateDenom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgTokenFactoryCreateDenom) GetSigners() []sdk.AccAddress {
	sender, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{sender}
}

// NewMsgMint creates a message to mint tokens
func NewMsgMint(sender string, amount sdk.Coin) *MsgTokenFactoryMint {
	return &MsgTokenFactoryMint{
		Sender: sender,
		Amount: amount,
	}
}

func NewMsgMintTo(sender string, amount sdk.Coin, mintToAddress string) *MsgTokenFactoryMint {
	return &MsgTokenFactoryMint{
		Sender:        sender,
		Amount:        amount,
		MintToAddress: mintToAddress,
	}
}

func (m MsgTokenFactoryMint) Route() string { return RouterKey }
func (m MsgTokenFactoryMint) Type() string  { return TypeMsgMint }
func (m MsgTokenFactoryMint) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	if m.MintToAddress != "" {
		_, err = sdk.AccAddressFromBech32(m.MintToAddress)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid mint to address (%s)", err)
		}
	}

	if !m.Amount.IsValid() || m.Amount.Amount.Equal(sdk.ZeroInt()) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}

	return nil
}

func (m MsgTokenFactoryMint) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgTokenFactoryMint) GetSigners() []sdk.AccAddress {
	sender, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{sender}
}

// NewMsgBurn creates a message to burn tokens
func NewMsgBurn(sender string, amount sdk.Coin) *MsgTokenFactoryBurn {
	return &MsgTokenFactoryBurn{
		Sender: sender,
		Amount: amount,
	}
}

// NewMsgBurn creates a message to burn tokens
func NewMsgBurnFrom(sender string, amount sdk.Coin, burnFromAddress string) *MsgTokenFactoryBurn {
	return &MsgTokenFactoryBurn{
		Sender:          sender,
		Amount:          amount,
		BurnFromAddress: burnFromAddress,
	}
}

func (m MsgTokenFactoryBurn) Route() string { return RouterKey }
func (m MsgTokenFactoryBurn) Type() string  { return TypeMsgBurn }
func (m MsgTokenFactoryBurn) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	if !m.Amount.IsValid() || m.Amount.Amount.Equal(sdk.ZeroInt()) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}

	if m.BurnFromAddress != "" {
		_, err = sdk.AccAddressFromBech32(m.BurnFromAddress)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid burn from address (%s)", err)
		}
	}

	return nil
}

func (m MsgTokenFactoryBurn) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgTokenFactoryBurn) GetSigners() []sdk.AccAddress {
	sender, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{sender}
}

// NewMsgForceTransfer creates a transfer funds from one account to another
func NewMsgForceTransfer(sender string, amount sdk.Coin, fromAddr, toAddr string) *MsgTokenFactoryForceTransfer {
	return &MsgTokenFactoryForceTransfer{
		Sender:              sender,
		Amount:              amount,
		TransferFromAddress: fromAddr,
		TransferToAddress:   toAddr,
	}
}

func (m MsgTokenFactoryForceTransfer) Route() string { return RouterKey }
func (m MsgTokenFactoryForceTransfer) Type() string  { return TypeMsgForceTransfer }
func (m MsgTokenFactoryForceTransfer) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(m.TransferFromAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid from address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(m.TransferToAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid to address (%s)", err)
	}

	if !m.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}

	return nil
}

func (m MsgTokenFactoryForceTransfer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgTokenFactoryForceTransfer) GetSigners() []sdk.AccAddress {
	sender, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{sender}
}

// NewMsgChangeAdmin creates a message to burn tokens
func NewMsgChangeAdmin(sender, denom, newAdmin string) *MsgTokenFactoryChangeAdmin {
	return &MsgTokenFactoryChangeAdmin{
		Sender:   sender,
		Denom:    denom,
		NewAdmin: newAdmin,
	}
}

func (m MsgTokenFactoryChangeAdmin) Route() string { return RouterKey }
func (m MsgTokenFactoryChangeAdmin) Type() string  { return TypeMsgChangeAdmin }
func (m MsgTokenFactoryChangeAdmin) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(m.NewAdmin)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid address (%s)", err)
	}

	_, _, err = DeconstructDenom(m.Denom)
	if err != nil {
		return err
	}

	return nil
}

func (m MsgTokenFactoryChangeAdmin) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgTokenFactoryChangeAdmin) GetSigners() []sdk.AccAddress {
	sender, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{sender}
}

// NewMsgChangeAdmin creates a message to burn tokens
func NewMsgSetDenomMetadata(sender string, metadata banktypes.Metadata) *MsgTokenFactorySetDenomMetadata {
	return &MsgTokenFactorySetDenomMetadata{
		Sender:   sender,
		Metadata: metadata,
	}
}

func (m MsgTokenFactorySetDenomMetadata) Route() string { return RouterKey }
func (m MsgTokenFactorySetDenomMetadata) Type() string  { return TypeMsgSetDenomMetadata }
func (m MsgTokenFactorySetDenomMetadata) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	err = m.Metadata.Validate()
	if err != nil {
		return err
	}

	_, _, err = DeconstructDenom(m.Metadata.Base)
	if err != nil {
		return err
	}

	return nil
}

func (m MsgTokenFactorySetDenomMetadata) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgTokenFactorySetDenomMetadata) GetSigners() []sdk.AccAddress {
	sender, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{sender}
}
