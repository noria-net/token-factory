# Token Factory

The token factory allows creating tokens using the bank module from the Cosmos SDK.

This Noria fork has upgraded dependencies to Cosmos SDK v0.47.2 and wasmvm 1.2.3.

## Messages

### CreateDenom

Creates a denom of `factory/{creator address}/{subdenom}` given the denom creator
address and the subdenom. Subdenoms can contain `[a-zA-Z0-9./]`.

```go
message MsgCreateDenom {
  string sender = 1 [ (gogoproto.moretags) = "yaml:\"sender\"" ];
  string subdenom = 2 [ (gogoproto.moretags) = "yaml:\"subdenom\"" ];
}
```

**State Modifications:**

- Fund community pool with the denom creation fee from the creator address, set
  in `Params`.
- Set `DenomMetaData` via bank keeper.
- Set `AuthorityMetadata` for the given denom to store the admin for the created
  denom `factory/{creator address}/{subdenom}`. Admin is automatically set as the
  Msg sender.
- Add denom to the `CreatorPrefixStore`, where a state of denoms created per
  creator is kept.

### Mint

Minting of a specific denom is only allowed for the current admin.
Note, the current admin is defaulted to the creator of the denom.

```go
message MsgMint {
  string sender = 1 [ (gogoproto.moretags) = "yaml:\"sender\"" ];
  cosmos.base.v1beta1.Coin amount = 2 [
    (gogoproto.moretags) = "yaml:\"amount\"",
    (gogoproto.nullable) = false
  ];
}
```

**State Modifications:**

- Safety check the following
  - Check that the denom minting is created via `tokenfactory` module
  - Check that the sender of the message is the admin of the denom
- Mint designated amount of tokens for the denom via `bank` module

### Burn

Burning of a specific denom is only allowed for the current admin.
Note, the current admin is defaulted to the creator of the denom.

```go
message MsgBurn {
  string sender = 1 [ (gogoproto.moretags) = "yaml:\"sender\"" ];
  cosmos.base.v1beta1.Coin amount = 2 [
    (gogoproto.moretags) = "yaml:\"amount\"",
    (gogoproto.nullable) = false
  ];
}
```

**State Modifications:**

- Saftey check the following
  - Check that the denom minting is created via `tokenfactory` module
  - Check that the sender of the message is the admin of the denom
- Burn designated amount of tokens for the denom via `bank` module

### ChangeAdmin

Change the admin of a denom. Note, this is only allowed to be called by the current admin of the denom.

```go
message MsgChangeAdmin {
  string sender = 1 [ (gogoproto.moretags) = "yaml:\"sender\"" ];
  string denom = 2 [ (gogoproto.moretags) = "yaml:\"denom\"" ];
  string newAdmin = 3 [ (gogoproto.moretags) = "yaml:\"new_admin\"" ];
}
```

### SetDenomMetadata

Setting of metadata for a specific denom is only allowed for the admin of the denom.
It allows the overwriting of the denom metadata in the bank module.

```go
message MsgChangeAdmin {
  string sender = 1 [ (gogoproto.moretags) = "yaml:\"sender\"" ];
  cosmos.bank.v1beta1.Metadata metadata = 2 [ (gogoproto.moretags) = "yaml:\"metadata\"", (gogoproto.nullable)   = false ];
}
```

**State Modifications:**

- Check that sender of the message is the admin of denom
- Modify `AuthorityMetadata` state entry to change the admin of the denom

## Expectations from the chain

The chain's bech32 prefix for addresses can be at most 16 characters long.

This comes from denoms having a 128 byte maximum length, enforced from the SDK,
and us setting longest_subdenom to be 44 bytes.

A token factory token's denom is: `factory/{creator address}/{subdenom}`

Splitting up into sub-components, this has:

- `len(factory) = 7`
- `2 * len("/") = 2`
- `len(longest_subdenom)`
- `len(creator_address) = len(bech32(longest_addr_length, chain_addr_prefix))`.

Longest addr length at the moment is `32 bytes`. Due to SDK error correction
settings, this means `len(bech32(32, chain_addr_prefix)) = len(chain_addr_prefix) + 1 + 58`.
Adding this all, we have a total length constraint of `128 = 7 + 2 + len(longest_subdenom) + len(longest_chain_addr_prefix) + 1 + 58`.
Therefore `len(longest_subdenom) + len(longest_chain_addr_prefix) = 128 - (7 + 2 + 1 + 58) = 60`.

The choice between how we standardized the split these 60 bytes between maxes
from longest_subdenom and longest_chain_addr_prefix is somewhat arbitrary.
Considerations going into this:

- Per [BIP-0173](https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki#bech32)
  the technically longest HRP for a 32 byte address ('data field') is 31 bytes.
  (Comes from encode(data) = 59 bytes, and max length = 90 bytes)
- subdenom should be at least 32 bytes so hashes can go into it
- longer subdenoms are very helpful for creating human readable denoms
- chain addresses should prefer being smaller. The longest HRP in cosmos to date is 11 bytes. (`persistence`)

For explicitness, its currently set to `len(longest_subdenom) = 44` and `len(longest_chain_addr_prefix) = 16`.

Please note, if the SDK increases the maximum length of a denom from 128 bytes,
these caps should increase.

So please don't make code rely on these max lengths for parsing.

## WASM Bindings for smart contracts

### ExecuteMsgs

```rust

#[cw_serde]
/// A number of Custom messages that can call into the chain custom modules bindings
pub enum ChainMsg { // can be renamed to your chain
    Token(TokenFactoryMsg),
}


#[cw_serde]
/// A number of Custom messages that can call into the TokenFactory bindings
pub enum TokenFactoryMsg {
    /// CreateDenom creates a new factory denom, of denomination:
    /// factory/{creating contract address}/{Subdenom}
    /// Subdenom can be of length at most 44 characters, in [0-9a-zA-Z./]
    /// The (creating contract address, subdenom) pair must be unique.
    /// The created denom's admin is the creating contract address,
    /// but this admin can be changed using the ChangeAdmin binding.
    CreateDenom {
        subdenom: String,
        metadata: Option<DenomMetadata>,
    },
    /// ChangeAdmin changes the admin for a factory denom.
    /// If the NewAdminAddress is empty, the denom has no admin.
    ChangeAdmin { denom: String, new_admin: Addr },
    /// Contracts can mint native tokens for an existing factory denom
    /// that they are the admin of.
    Mint {
        denom: String,
        amount: String,
        mint_to_address: Addr,
    },
    /// Contracts can burn native tokens for an existing factory denom
    /// that they are the admin of.
    /// Currently, the burn from address must be the admin contract.
    Burn {
        denom: String,
        amount: String,
        burn_from_address: Addr,
    },
    /// Sets the metadata on a denom which the contract controls
    SetDenomMetadata { metadata: DenomMetadata },
    /// Forces a transfer of tokens from one address to another.
    ForceTransfer {
        denom: String,
        from_address: Addr,
        to_address: Addr,
        amount: String,
    },
}

/// DenomUnit is used to describe a token for the Bank module; part of the SetDenomMetadata message
#[cw_serde]
pub struct DenomUnit {
    /// Denom represents the string name of the given denom unit (e.g uatom). pub denom: String,
    pub denom: String,
    /// Exponent represents power of 10 exponent that one must
    /// raise the base_denom to in order to equal the given DenomUnit's denom
    /// 1 denom = 1^exponent base_denom
    /// (e.g. with a base_denom of uatom, one can create a DenomUnit of 'atom' with
    /// exponent = 6, thus: 1 atom = 10^6 uatom).
    pub exponent: u32,
    /// Aliases is a list of string aliases for the given denom
    pub aliases: Vec<String>,
}

/// DenomMetadata is used to describe a token for the Bank module; part of the SetDenomMetadata message
#[cw_serde]
pub struct DenomMetadata {
    pub description: String,
    /// DenomUnits represents the list of DenomUnit's for a given coin
    pub denom_units: Vec<DenomUnit>,
    /// Base represents the base denom (should be the DenomUnit with exponent = 0).
    pub base: String,
    /// Display indicates the suggested denom that should be displayed in clients.
    pub display: String,
    /// Name defines the name of the token (eg: Cosmos Atom)
    pub name: String,
    /// Symbol is the token symbol usually shown on exchanges (eg: ATOM).
    /// This can be the same as the display.
    pub symbol: String,
}
```

### Queries

```rust
#[cw_serde]
/// A number of Custom queries that can call into the chain custom modules bindings
pub enum ChainQuery { // can be renamed to your chain
    Token(TokenFactoryQuery),
}

#[cw_serde]
/// A number of Custom queries that can call into the TokenFactory bindings
pub enum TokenFactoryQuery {
    FullDenom { subdenom: String, creator_addr: Addr },
    Admin { denom: String },
    Metadata { denom: String },
    DenomsByCreator { creator: Addr },
    Params {}
}
```

### Query Responses

```rust
#[cw_serde]
pub FullDenomResponse struct {
	pub denom: String;
}

#[cw_serde]
pub AdminResponse struct {
	pub admin: String;
}

#[cw_serde]
pub MetadataResponse struct {
	pub metadata Option<Metadata>;
}

#[cw_serde]
pub DenomsByCreatorResponse struct {
	pub denoms Vec<String>;
}

#[cw_serde]
pub ParamsResponse struct {
	pub params TokenParams;
}

#[cw_serde]
pub struct TokenParams {
    pub denom_creation_fee: Vec<DenomCreationFee>,
}

#[cw_serde]
pub struct DenomCreationFee {
    pub amount: String,
    pub denom: String,
}
```
