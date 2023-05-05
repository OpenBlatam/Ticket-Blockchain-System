# Phase 0 -- The Beacon Chain

## Table of contents
<!-- TOC -->
<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Introduction](#introduction)
- [Notation](#notation)
- [Custom types](#custom-types)
- [Constants](#constants)
    - [Misc](#misc)
    - [Withdrawal prefixes](#withdrawal-prefixes)
    - [Domain types](#domain-types)
- [Preset](#preset)
    - [Misc](#misc-1)
    - [Gwei values](#gwei-values)
    - [Time parameters](#time-parameters)
    - [State list lengths](#state-list-lengths)
    - [Rewards and penalties](#rewards-and-penalties)
    - [Max operations per block](#max-operations-per-block)
- [Configuration](#configuration)
    - [Genesis settings](#genesis-settings)
    - [Time parameters](#time-parameters-1)
    - [Validator cycle](#validator-cycle)
- [Containers](#containers)
    - [Misc dependencies](#misc-dependencies)
        - [`Fork`](#fork)
        - [`ForkData`](#forkdata)
        - [`Checkpoint`](#checkpoint)
        - [`Validator`](#validator)
        - [`AttestationData`](#attestationdata)
        - [`IndexedAttestation`](#indexedattestation)
        - [`PendingAttestation`](#pendingattestation)
        - [`Eth1Data`](#eth1data)
        - [`HistoricalBatch`](#historicalbatch)
        - [`DepositMessage`](#depositmessage)
        - [`DepositData`](#depositdata)
        - [`BeaconBlockHeader`](#beaconblockheader)
        - [`SigningData`](#signingdata)
    - [Beacon operations](#beacon-operations)
        - [`ProposerSlashing`](#proposerslashing)
        - [`AttesterSlashing`](#attesterslashing)
        - [`Attestation`](#attestation)
        - [`Deposit`](#deposit)
        - [`VoluntaryExit`](#voluntaryexit)
    - [Beacon blocks](#beacon-blocks)
        - [`BeaconBlockBody`](#beaconblockbody)
        - [`BeaconBlock`](#beaconblock)
    - [Beacon state](#beacon-state)
        - [`BeaconState`](#beaconstate)
    - [Signed envelopes](#signed-envelopes)
        - [`SignedVoluntaryExit`](#signedvoluntaryexit)
        - [`SignedBeaconBlock`](#signedbeaconblock)
        - [`SignedBeaconBlockHeader`](#signedbeaconblockheader)
- [Helper functions](#helper-functions)
    - [Math](#math)
        - [`integer_squareroot`](#integer_squareroot)
        - [`xor`](#xor)
        - [`uint_to_bytes`](#uint_to_bytes)
        - [`bytes_to_uint64`](#bytes_to_uint64)
    - [Crypto](#crypto)
        - [`hash`](#hash)
        - [`hash_tree_root`](#hash_tree_root)
        - [BLS signatures](#bls-signatures)
    - [Predicates](#predicates)
        - [`is_active_validator`](#is_active_validator)
        - [`is_eligible_for_activation_queue`](#is_eligible_for_activation_queue)
        - [`is_eligible_for_activation`](#is_eligible_for_activation)
        - [`is_slashable_validator`](#is_slashable_validator)
        - [`is_slashable_attestation_data`](#is_slashable_attestation_data)
        - [`is_valid_indexed_attestation`](#is_valid_indexed_attestation)
        - [`is_valid_merkle_branch`](#is_valid_merkle_branch)
    - [Misc](#misc-2)
        - [`compute_shuffled_index`](#compute_shuffled_index)
        - [`compute_proposer_index`](#compute_proposer_index)
        - [`compute_committee`](#compute_committee)
        - [`compute_epoch_at_slot`](#compute_epoch_at_slot)
        - [`compute_start_slot_at_epoch`](#compute_start_slot_at_epoch)
        - [`compute_activation_exit_epoch`](#compute_activation_exit_epoch)
        - [`compute_fork_data_root`](#compute_fork_data_root)
        - [`compute_fork_digest`](#compute_fork_digest)
        - [`compute_domain`](#compute_domain)
        - [`compute_signing_root`](#compute_signing_root)
    - [Beacon state accessors](#beacon-state-accessors)
        - [`get_current_epoch`](#get_current_epoch)
        - [`get_previous_epoch`](#get_previous_epoch)
        - [`get_block_root`](#get_block_root)
        - [`get_block_root_at_slot`](#get_block_root_at_slot)
        - [`get_randao_mix`](#get_randao_mix)
        - [`get_active_validator_indices`](#get_active_validator_indices)
        - [`get_validator_churn_limit`](#get_validator_churn_limit)
        - [`get_seed`](#get_seed)
        - [`get_committee_count_per_slot`](#get_committee_count_per_slot)
        - [`get_beacon_committee`](#get_beacon_committee)
        - [`get_beacon_proposer_index`](#get_beacon_proposer_index)
        - [`get_total_balance`](#get_total_balance)
        - [`get_total_active_balance`](#get_total_active_balance)
        - [`get_domain`](#get_domain)
        - [`get_indexed_attestation`](#get_indexed_attestation)
        - [`get_attesting_indices`](#get_attesting_indices)
    - [Beacon state mutators](#beacon-state-mutators)
        - [`increase_balance`](#increase_balance)
        - [`decrease_balance`](#decrease_balance)
        - [`initiate_validator_exit`](#initiate_validator_exit)
        - [`slash_validator`](#slash_validator)
- [Genesis](#genesis)
    - [Genesis state](#genesis-state)
    - [Genesis block](#genesis-block)
- [Beacon chain state transition function](#beacon-chain-state-transition-function)
    - [Epoch processing](#epoch-processing)
        - [Helper functions](#helper-functions-1)
        - [Justification and finalization](#justification-and-finalization)
        - [Rewards and penalties](#rewards-and-penalties-1)
            - [Helpers](#helpers)
            - [Components of attestation deltas](#components-of-attestation-deltas)
            - [`get_attestation_deltas`](#get_attestation_deltas)
            - [`process_rewards_and_penalties`](#process_rewards_and_penalties)
        - [Registry updates](#registry-updates)
        - [Slashings](#slashings)
        - [Eth1 data votes updates](#eth1-data-votes-updates)
        - [Effective balances updates](#effective-balances-updates)
        - [Slashings balances updates](#slashings-balances-updates)
        - [Randao mixes updates](#randao-mixes-updates)
        - [Historical roots updates](#historical-roots-updates)
        - [Participation records rotation](#participation-records-rotation)
    - [Block processing](#block-processing)
        - [Block header](#block-header)
        - [RANDAO](#randao)
        - [Eth1 data](#eth1-data)
        - [Operations](#operations)
            - [Proposer slashings](#proposer-slashings)
            - [Attester slashings](#attester-slashings)
            - [Attestations](#attestations)
            - [Deposits](#deposits)
            - [Voluntary exits](#voluntary-exits)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->
<!-- /TOC -->

## Introduction

This document represents the specification for Phase 0 -- The Beacon Chain.

At the core of Ethereum proof-of-stake is a system chain called the "beacon chain". The beacon chain stores and manages the registry of validators. In the initial deployment phases of proof-of-stake, the only mechanism to become a validator is to make a one-way ETH transaction to a deposit contract on the Ethereum proof-of-work chain. Activation as a validator happens when deposit receipts are processed by the beacon chain, the activation balance is reached, and a queuing process is completed. Exit is either voluntary or done forcibly as a penalty for misbehavior.
The primary source of load on the beacon chain is "attestations". Attestations are simultaneously availability votes for a shard block (in a later upgrade) and proof-of-stake votes for a beacon block (Phase 0).

## Notation

Code snippets appearing in `this style` are to be interpreted as Python 3 code.

## Custom types

We define the following Python custom types for type hinting and readability:

| Name       | SSZ equivalent | Description     |
|------------|----------------|-----------------|
| TicketID   | uint64         | a ticket number |
| Validation |                |                 |
|            |                |                 |
|            |                |                 |
|            |                |                 |
|            |                |                 |

| Name | SSZ equivalent | Description |
| - | - | - |
| `Slot` | `uint64` | a slot number |
| `Epoch` | `uint64` | an epoch number |
| `CommitteeIndex` | `uint64` | a committee index at a slot |
| `ValidatorIndex` | `uint64` | a validator registry index |
| `Gwei` | `uint64` | an amount in Gwei |
| `Root` | `Bytes32` | a Merkle root |
| `Hash32` | `Bytes32` | a 256-bit hash |
| `Version` | `Bytes4` | a fork version number |
| `DomainType` | `Bytes4` | a domain type |
| `ForkDigest` | `Bytes4` | a digest of the current fork data |
| `Domain` | `Bytes32` | a signature domain |
| `BLSPubkey` | `Bytes48` | a BLS12-381 public key |
| `BLSSignature` | `Bytes96` | a BLS12-381 signature |

## Constants

The following values are (non-configurable) constants used throughout the specification.

### Misc

| Name | Value |
| - | - |
| `GENESIS_SLOT` | `Slot(0)` |
| `GENESIS_EPOCH` | `Epoch(0)` |
| `FAR_FUTURE_EPOCH` | `Epoch(2**64 - 1)` |
| `BASE_REWARDS_PER_EPOCH` | `uint64(4)` |
| `DEPOSIT_CONTRACT_TREE_DEPTH` | `uint64(2**5)` (= 32) |
| `JUSTIFICATION_BITS_LENGTH` | `uint64(4)` |
| `ENDIANNESS` | `'little'` |

### Withdrawal prefixes

| Name | Value |
| - | - |
| `BLS_WITHDRAWAL_PREFIX` | `Bytes1('0x00')` |
| `ETH1_ADDRESS_WITHDRAWAL_PREFIX` | `Bytes1('0x01')` |

### Domain types

| Name | Value |
| - | - |
| `DOMAIN_BEACON_PROPOSER`     | `DomainType('0x00000000')` |
| `DOMAIN_BEACON_ATTESTER`     | `DomainType('0x01000000')` |
| `DOMAIN_RANDAO`              | `DomainType('0x02000000')` |
| `DOMAIN_DEPOSIT`             | `DomainType('0x03000000')` |
| `DOMAIN_VOLUNTARY_EXIT`      | `DomainType('0x04000000')` |
| `DOMAIN_SELECTION_PROOF`     | `DomainType('0x05000000')` |
| `DOMAIN_AGGREGATE_AND_PROOF` | `DomainType('0x06000000')` |
| `DOMAIN_APPLICATION_MASK`    | `DomainType('0x00000001')` |

*Note*: `DOMAIN_APPLICATION_MASK` reserves the rest of the bitspace in `DomainType` for application usage. This means for some `DomainType` `DOMAIN_SOME_APPLICATION`, `DOMAIN_SOME_APPLICATION & DOMAIN_APPLICATION_MASK` **MUST** be non-zero. This expression for any other `DomainType` in the consensus specs **MUST** be zero.

## Preset

*Note*: The below configuration is bundled as a preset: a bundle of configuration variables which are expected to differ
between different modes of operation, e.g. testing, but not generally between different networks.
Additional preset configurations can be found in the [`configs`](../../configs) directory.

### Misc

| Name | Value |
| - | - |
| `MAX_COMMITTEES_PER_SLOT` | `uint64(2**6)` (= 64) |
| `TARGET_COMMITTEE_SIZE` | `uint64(2**7)` (= 128) |
| `MAX_VALIDATORS_PER_COMMITTEE` | `uint64(2**11)` (= 2,048) |
| `SHUFFLE_ROUND_COUNT` | `uint64(90)` |
| `HYSTERESIS_QUOTIENT` | `uint64(4)` |
| `HYSTERESIS_DOWNWARD_MULTIPLIER` | `uint64(1)` |
| `HYSTERESIS_UPWARD_MULTIPLIER` | `uint64(5)` |

- For the safety of committees, `TARGET_COMMITTEE_SIZE` exceeds [the recommended minimum committee size of 111](http://web.archive.org/web/20190504131341/https://vitalik.ca/files/Ithaca201807_Sharding.pdf); with sufficient active validators (at least `SLOTS_PER_EPOCH * TARGET_COMMITTEE_SIZE`), the shuffling algorithm ensures committee sizes of at least `TARGET_COMMITTEE_SIZE`. (Unbiasable randomness with a Verifiable Delay Function (VDF) will improve committee robustness and lower the safe minimum committee size.)

### Gwei values

| Name | Value |
| - | - |
| `MIN_DEPOSIT_AMOUNT` | `Gwei(2**0 * 10**9)` (= 1,000,000,000) |
| `MAX_EFFECTIVE_BALANCE` | `Gwei(2**5 * 10**9)` (= 32,000,000,000) |
| `EFFECTIVE_BALANCE_INCREMENT` | `Gwei(2**0 * 10**9)` (= 1,000,000,000) |

### Time parameters

| Name | Value | Unit | Duration |
| - | - | :-: | :-: |
| `MIN_ATTESTATION_INCLUSION_DELAY` | `uint64(2**0)` (= 1) | slots | 12 seconds |
| `SLOTS_PER_EPOCH` | `uint64(2**5)` (= 32) | slots | 6.4 minutes |
| `MIN_SEED_LOOKAHEAD` | `uint64(2**0)` (= 1) | epochs | 6.4 minutes |
| `MAX_SEED_LOOKAHEAD` | `uint64(2**2)` (= 4) | epochs | 25.6 minutes |
| `MIN_EPOCHS_TO_INACTIVITY_PENALTY` | `uint64(2**2)` (= 4) | epochs | 25.6 minutes |
| `EPOCHS_PER_ETH1_VOTING_PERIOD` | `uint64(2**6)` (= 64) | epochs | ~6.8 hours |
| `SLOTS_PER_HISTORICAL_ROOT` | `uint64(2**13)` (= 8,192) | slots | ~27 hours |

### State list lengths

| Name | Value | Unit | Duration |
| - | - | :-: | :-: |
| `EPOCHS_PER_HISTORICAL_VECTOR` | `uint64(2**16)` (= 65,536) | epochs | ~0.8 years |
| `EPOCHS_PER_SLASHINGS_VECTOR` | `uint64(2**13)` (= 8,192) | epochs | ~36 days |
| `HISTORICAL_ROOTS_LIMIT` | `uint64(2**24)` (= 16,777,216) | historical roots | ~52,262 years |
| `VALIDATOR_REGISTRY_LIMIT` | `uint64(2**40)` (= 1,099,511,627,776) | validators |

### Rewards and penalties

| Name | Value |
| - | - |
| `BASE_REWARD_FACTOR` | `uint64(2**6)` (= 64) |
| `WHISTLEBLOWER_REWARD_QUOTIENT` | `uint64(2**9)` (= 512) |
| `PROPOSER_REWARD_QUOTIENT` | `uint64(2**3)` (= 8) |
| `INACTIVITY_PENALTY_QUOTIENT` | `uint64(2**26)` (= 67,108,864) |
| `MIN_SLASHING_PENALTY_QUOTIENT` | `uint64(2**7)` (= 128) |
| `PROPORTIONAL_SLASHING_MULTIPLIER` | `uint64(1)` |

- The `INACTIVITY_PENALTY_QUOTIENT` equals `INVERSE_SQRT_E_DROP_TIME**2` where `INVERSE_SQRT_E_DROP_TIME := 2**13` epochs (about 36 days) is the time it takes the inactivity penalty to reduce the balance of non-participating validators to about `1/sqrt(e) ~= 60.6%`. Indeed, the balance retained by offline validators after `n` epochs is about `(1 - 1/INACTIVITY_PENALTY_QUOTIENT)**(n**2/2)`; so after `INVERSE_SQRT_E_DROP_TIME` epochs, it is roughly `(1 - 1/INACTIVITY_PENALTY_QUOTIENT)**(INACTIVITY_PENALTY_QUOTIENT/2) ~= 1/sqrt(e)`. Note this value will be upgraded to `2**24` after Phase 0 mainnet stabilizes to provide a faster recovery in the event of an inactivity leak.

- The `PROPORTIONAL_SLASHING_MULTIPLIER` is set to `1` at initial mainnet launch, resulting in one-third of the minimum accountable safety margin in the event of a finality attack. After Phase 0 mainnet stabilizes, this value will be upgraded to `3` to provide the maximal minimum accountable safety margin.

### Max operations per block

| Name | Value |
| - | - |
| `MAX_PROPOSER_SLASHINGS` | `2**4` (= 16) |
| `MAX_ATTESTER_SLASHINGS` | `2**1` (= 2) |
| `MAX_ATTESTATIONS` | `2**7` (= 128) |
| `MAX_DEPOSITS` | `2**4` (= 16) |
| `MAX_VOLUNTARY_EXITS` | `2**4` (= 16) |

## Configuration

*Note*: The default mainnet configuration values are included here for illustrative purposes.
Defaults for this more dynamic type of configuration are available with the presets in the [`configs`](../../configs) directory.
Testnets and other types of chain instances may use a different configuration.

### Genesis settings

| Name | Value |
| - | - |
| `MIN_GENESIS_ACTIVE_VALIDATOR_COUNT` | `uint64(2**14)` (= 16,384) |
| `MIN_GENESIS_TIME` | `uint64(1606824000)` (Dec 1, 2020, 12pm UTC) |
| `GENESIS_FORK_VERSION` | `Version('0x00000000')` |
| `GENESIS_DELAY` | `uint64(604800)` (7 days) |

### Time parameters

| Name | Value | Unit | Duration |
| - | - | :-: | :-: |
| `SECONDS_PER_SLOT` | `uint64(12)` | seconds | 12 seconds |
| `SECONDS_PER_ETH1_BLOCK` | `uint64(14)` | seconds | 14 seconds |
| `MIN_VALIDATOR_WITHDRAWABILITY_DELAY` | `uint64(2**8)` (= 256) | epochs | ~27 hours |
| `SHARD_COMMITTEE_PERIOD` | `uint64(2**8)` (= 256) | epochs | ~27 hours |
| `ETH1_FOLLOW_DISTANCE` | `uint64(2**11)` (= 2,048) | Eth1 blocks | ~8 hours |

### Validator cycle

| Name | Value |
| - | - |
| `EJECTION_BALANCE` | `Gwei(2**4 * 10**9)` (= 16,000,000,000) |
| `MIN_PER_EPOCH_CHURN_LIMIT` | `uint64(2**2)` (= 4) |
| `CHURN_LIMIT_QUOTIENT` | `uint64(2**16)` (= 65,536) |

## Containers

The following types are [SimpleSerialize (SSZ)](../../ssz/simple-serialize.md) containers.

*Note*: The definitions are ordered topologically to facilitate execution of the spec.

*Note*: Fields missing in container instantiations default to their zero value.

### Misc dependencies

#### `Fork`

```python
class Fork(Container):
    previous_version: Version
    current_version: Version
    epoch: Epoch  # Epoch of latest fork
```

#### `ForkData`

```python
class ForkData(Container):
    current_version: Version
    genesis_validators_root: Root
```

#### `Checkpoint`

```python
class Checkpoint(Container):
    epoch: Epoch
    root: Root
```

#### `Validator`

```python
class Validator(Container):
    pubkey: BLSPubkey
    withdrawal_credentials: Bytes32  # Commitment to pubkey for withdrawals
    effective_balance: Gwei  # Balance at stake
    slashed: boolean
    # Status epochs
    activation_eligibility_epoch: Epoch  # When criteria for activation were met
    activation_epoch: Epoch
    exit_epoch: Epoch
    withdrawable_epoch: Epoch  # When validator can withdraw funds
```

#### `AttestationData`

```python
class AttestationData(Container):
    slot: Slot
    index: CommitteeIndex
    # LMD GHOST vote
    beacon_block_root: Root
    # FFG vote
    source: Checkpoint
    target: Checkpoint
```

#### `IndexedAttestation`

```python
class IndexedAttestation(Container):
    attesting_indices: List[ValidatorIndex, MAX_VALIDATORS_PER_COMMITTEE]
    data: AttestationData
    signature: BLSSignature
```

#### `PendingAttestation`

```python
class PendingAttestation(Container):
    aggregation_bits: Bitlist[MAX_VALIDATORS_PER_COMMITTEE]
    data: AttestationData
    inclusion_delay: Slot
    proposer_index: ValidatorIndex
```

#### `Eth1Data`

```python
class Eth1Data(Container):
    deposit_root: Root
    deposit_count: uint64
    block_hash: Hash32
```

#### `HistoricalBatch`

```python
class HistoricalBatch(Container):
    block_roots: Vector[Root, SLOTS_PER_HISTORICAL_ROOT]
    state_roots: Vector[Root, SLOTS_PER_HISTORICAL_ROOT]
```

#### `DepositMessage`

```python
class DepositMessage(Container):
    pubkey: BLSPubkey
    withdrawal_credentials: Bytes32
    amount: Gwei
```

#### `DepositData`

```python
class DepositData(Container):
    pubkey: BLSPubkey
    withdrawal_credentials: Bytes32
    amount: Gwei
    signature: BLSSignature  # Signing over DepositMessage
```

#### `BeaconBlockHeader`

```python
class BeaconBlockHeader(Container):
    slot: Slot
    proposer_index: ValidatorIndex
    parent_root: Root
    state_root: Root
    body_root: Root
```

#### `SigningData`

```python
class SigningData(Container):
    object_root: Root
    domain: Domain
```

### Beacon operations

#### `ProposerSlashing`

```python
class ProposerSlashing(Container):
    signed_header_1: SignedBeaconBlockHeader
    signed_header_2: SignedBeaconBlockHeader
```

#### `AttesterSlashing`

```python
class AttesterSlashing(Container):
    attestation_1: IndexedAttestation
    attestation_2: IndexedAttestation
```

#### `Attestation`

```python
class Attestation(Container):
    aggregation_bits: Bitlist[MAX_VALIDATORS_PER_COMMITTEE]
    data: AttestationData
    signature: BLSSignature
```

#### `Deposit`

```python
class Deposit(Container):
    proof: Vector[Bytes32, DEPOSIT_CONTRACT_TREE_DEPTH + 1]  # Merkle path to deposit root
    data: DepositData
```

#### `VoluntaryExit`

```python
class VoluntaryExit(Container):
    epoch: Epoch  # Earliest epoch when voluntary exit can be processed
    validator_index: ValidatorIndex
```

### Beacon blocks

#### `BeaconBlockBody`

```python
class BeaconBlockBody(Container):
    randao_reveal: BLSSignature
    eth1_data: Eth1Data  # Eth1 data vote
    graffiti: Bytes32  # Arbitrary data
    # Operations
    proposer_slashings: List[ProposerSlashing, MAX_PROPOSER_SLASHINGS]
    attester_slashings: List[AttesterSlashing, MAX_ATTESTER_SLASHINGS]
    attestations: List[Attestation, MAX_ATTESTATIONS]
    deposits: List[Deposit, MAX_DEPOSITS]
    voluntary_exits: List[SignedVoluntaryExit, MAX_VOLUNTARY_EXITS]
```

#### `BeaconBlock`

```python
class BeaconBlock(Container):
    slot: Slot
    proposer_index: ValidatorIndex
    parent_root: Root
    state_root: Root
    body: BeaconBlockBody
```

### Beacon state

#### `BeaconState`

```python
class BeaconState(Container):
    # Versioning
    genesis_time: uint64
    genesis_validators_root: Root
    slot: Slot
    fork: Fork
    # History
    latest_block_header: BeaconBlockHeader
    block_roots: Vector[Root, SLOTS_PER_HISTORICAL_ROOT]
    state_roots: Vector[Root, SLOTS_PER_HISTORICAL_ROOT]
    historical_roots: List[Root, HISTORICAL_ROOTS_LIMIT]
    # Eth1
    eth1_data: Eth1Data
    eth1_data_votes: List[Eth1Data, EPOCHS_PER_ETH1_VOTING_PERIOD * SLOTS_PER_EPOCH]
    eth1_deposit_index: uint64
    # Registry
    validators: List[Validator, VALIDATOR_REGISTRY_LIMIT]
    balances: List[Gwei, VALIDATOR_REGISTRY_LIMIT]
    # Randomness
    randao_mixes: Vector[Bytes32, EPOCHS_PER_HISTORICAL_VECTOR]
    # Slashings
    slashings: Vector[Gwei, EPOCHS_PER_SLASHINGS_VECTOR]  # Per-epoch sums of slashed effective balances
    # Attestations
    previous_epoch_attestations: List[PendingAttestation, MAX_ATTESTATIONS * SLOTS_PER_EPOCH]
    current_epoch_attestations: List[PendingAttestation, MAX_ATTESTATIONS * SLOTS_PER_EPOCH]
    # Finality
    justification_bits: Bitvector[JUSTIFICATION_BITS_LENGTH]  # Bit set for every recent justified epoch
    previous_justified_checkpoint: Checkpoint  # Previous epoch snapshot
    current_justified_checkpoint: Checkpoint
    finalized_checkpoint: Checkpoint
```

