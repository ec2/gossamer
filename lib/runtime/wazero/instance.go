// Copyright 2023 ChainSafe Systems (ON)
// SPDX-License-Identifier: LGPL-3.0-only

package wazero_runtime

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/internal/log"
	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/lib/crypto"
	"github.com/ChainSafe/gossamer/lib/crypto/ed25519"
	"github.com/ChainSafe/gossamer/lib/keystore"
	"github.com/ChainSafe/gossamer/lib/runtime"
	"github.com/ChainSafe/gossamer/lib/runtime/allocator"
	"github.com/ChainSafe/gossamer/lib/runtime/offchain"
	"github.com/ChainSafe/gossamer/lib/transaction"
	"github.com/ChainSafe/gossamer/pkg/scale"
	"github.com/ChainSafe/gossamer/pkg/trie"
	"github.com/klauspost/compress/zstd"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

// Name represents the name of the interpreter
const Name = "wazero"

type contextKey string

const runtimeContextKey = contextKey("runtime.Context")

var _ runtime.Instance = &Instance{}

// Instance backed by wazero.Runtime
type Instance struct {
	Runtime  wazero.Runtime
	Module   api.Module
	Context  *runtime.Context
	codeHash common.Hash
	heapBase uint32
	sync.Mutex
}

// Config is the configuration used to create a Wasmer runtime instance.
type Config struct {
	Storage        runtime.Storage
	Keystore       *keystore.GlobalKeystore
	LogLvl         log.Level
	Role           common.NetworkRole
	NodeStorage    runtime.NodeStorage
	Network        runtime.BasicNetwork
	Transaction    runtime.TransactionState
	CodeHash       common.Hash
	DefaultVersion *runtime.Version
}

func decompressWasm(code []byte) ([]byte, error) {
	compressionFlag := []byte{82, 188, 83, 118, 70, 219, 142, 5}
	if !bytes.HasPrefix(code, compressionFlag) {
		return code, nil
	}

	decoder, err := zstd.NewReader(nil)
	if err != nil {
		return nil, fmt.Errorf("creating zstd reader: %s", err)
	}
	bytes, err := decoder.DecodeAll(code[len(compressionFlag):], nil)
	if err != nil {
		return nil, err
	}
	return bytes, err
}

// NewRuntimeFromGenesis creates a runtime instance from the genesis data
func NewRuntimeFromGenesis(cfg Config) (instance *Instance, err error) {
	if cfg.Storage == nil {
		return nil, errors.New("storage is nil")
	}

	code := cfg.Storage.LoadCode()
	if len(code) == 0 {
		return nil, fmt.Errorf("cannot find :code in state")
	}

	return NewInstance(code, cfg)
}

// NewInstanceFromTrie returns a new runtime instance with the code provided in the given trie
func NewInstanceFromTrie(t trie.Trie, cfg Config) (*Instance, error) {
	code := t.Get(common.CodeKey)
	if len(code) == 0 {
		return nil, fmt.Errorf("cannot find :code in trie")
	}

	return NewInstance(code, cfg)
}

// NewInstance instantiates a runtime from raw wasm bytecode
func NewInstance(code []byte, cfg Config) (instance *Instance, err error) {
	logger.Info("instantiating a runtime!")
	logger.Patch(log.SetLevel(cfg.LogLvl), log.SetCallerFunc(true))

	ctx := context.Background()
	rt := wazero.NewRuntime(ctx)

	_, err = rt.NewHostModuleBuilder("env").
		// values from newer kusama/polkadot runtimes
		ExportMemory("memory", 23).
		NewFunctionBuilder().
		WithFunc(ext_logging_log_version_1).
		Export("ext_logging_log_version_1").
		NewFunctionBuilder().
		WithFunc(func() int32 {
			return 4
		}).
		Export("ext_logging_max_level_version_1").
		NewFunctionBuilder().
		WithFunc(func(a int32, b int32, c int32) {
			panic("unimplemented")
		}).
		Export("ext_transaction_index_index_version_1").
		NewFunctionBuilder().
		WithFunc(func(a int32, b int32) {
			panic("unimplemented")
		}).
		Export("ext_transaction_index_renew_version_1").
		NewFunctionBuilder().
		WithFunc(func(a int32) {
			panic("unimplemented")
		}).
		Export("ext_sandbox_instance_teardown_version_1").
		NewFunctionBuilder().
		WithFunc(func(a int32, b int64, c int64, d int32) int32 {
			panic("unimplemented")
		}).
		Export("ext_sandbox_instantiate_version_1").
		NewFunctionBuilder().
		WithFunc(func(a int32, b int64, c int64, d int32, e int32, f int32) int32 {
			panic("unimplemented")
		}).
		Export("ext_sandbox_invoke_version_1").
		NewFunctionBuilder().
		WithFunc(func(a int32, b int32, c int32, d int32) int32 {
			panic("unimplemented")
		}).
		Export("ext_sandbox_memory_get_version_1").
		NewFunctionBuilder().
		WithFunc(func(a int32, b int32, c int32, d int32) int32 {
			panic("unimplemented")
		}).
		Export("ext_sandbox_memory_set_version_1").
		NewFunctionBuilder().
		WithFunc(func(a int32) {
			panic("unimplemented")
		}).
		Export("ext_sandbox_memory_teardown_version_1").
		NewFunctionBuilder().
		WithFunc(ext_crypto_ed25519_generate_version_1).
		Export("ext_crypto_ed25519_generate_version_1").
		NewFunctionBuilder().
		WithFunc(ext_crypto_ed25519_public_keys_version_1).
		Export("ext_crypto_ed25519_public_keys_version_1").
		NewFunctionBuilder().
		WithFunc(ext_crypto_ed25519_sign_version_1).
		Export("ext_crypto_ed25519_sign_version_1").
		NewFunctionBuilder().
		WithFunc(ext_crypto_ed25519_verify_version_1).
		Export("ext_crypto_ed25519_verify_version_1").
		NewFunctionBuilder().
		WithFunc(ext_crypto_secp256k1_ecdsa_recover_version_1).
		Export("ext_crypto_secp256k1_ecdsa_recover_version_1").
		NewFunctionBuilder().
		WithFunc(ext_crypto_secp256k1_ecdsa_recover_version_2).
		Export("ext_crypto_secp256k1_ecdsa_recover_version_2").
		NewFunctionBuilder().
		WithFunc(ext_crypto_ecdsa_verify_version_2).
		Export("ext_crypto_ecdsa_verify_version_2").
		NewFunctionBuilder().
		WithFunc(ext_crypto_secp256k1_ecdsa_recover_compressed_version_1).
		Export("ext_crypto_secp256k1_ecdsa_recover_compressed_version_1").
		NewFunctionBuilder().
		WithFunc(ext_crypto_secp256k1_ecdsa_recover_compressed_version_2).
		Export("ext_crypto_secp256k1_ecdsa_recover_compressed_version_2").
		NewFunctionBuilder().
		WithFunc(ext_crypto_sr25519_generate_version_1).
		Export("ext_crypto_sr25519_generate_version_1").
		NewFunctionBuilder().
		WithFunc(ext_crypto_sr25519_public_keys_version_1).
		Export("ext_crypto_sr25519_public_keys_version_1").
		NewFunctionBuilder().
		WithFunc(ext_crypto_sr25519_sign_version_1).
		Export("ext_crypto_sr25519_sign_version_1").
		NewFunctionBuilder().
		WithFunc(ext_crypto_sr25519_verify_version_1).
		Export("ext_crypto_sr25519_verify_version_1").
		NewFunctionBuilder().
		WithFunc(ext_crypto_sr25519_verify_version_2).
		Export("ext_crypto_sr25519_verify_version_2").
		NewFunctionBuilder().
		WithFunc(ext_crypto_start_batch_verify_version_1).
		Export("ext_crypto_start_batch_verify_version_1").
		NewFunctionBuilder().
		WithFunc(ext_crypto_finish_batch_verify_version_1).
		Export("ext_crypto_finish_batch_verify_version_1").
		NewFunctionBuilder().
		WithFunc(ext_trie_blake2_256_root_version_1).
		Export("ext_trie_blake2_256_root_version_1").
		NewFunctionBuilder().
		WithFunc(ext_trie_blake2_256_root_version_2).
		Export("ext_trie_blake2_256_root_version_2").
		NewFunctionBuilder().
		WithFunc(ext_trie_blake2_256_ordered_root_version_1).
		Export("ext_trie_blake2_256_ordered_root_version_1").
		NewFunctionBuilder().
		WithFunc(ext_trie_blake2_256_ordered_root_version_2).
		Export("ext_trie_blake2_256_ordered_root_version_2").
		NewFunctionBuilder().
		WithFunc(ext_trie_blake2_256_verify_proof_version_1).
		Export("ext_trie_blake2_256_verify_proof_version_1").
		NewFunctionBuilder().
		WithFunc(ext_trie_blake2_256_verify_proof_version_2).
		Export("ext_trie_blake2_256_verify_proof_version_2").
		NewFunctionBuilder().
		WithFunc(ext_misc_print_hex_version_1).
		Export("ext_misc_print_hex_version_1").
		NewFunctionBuilder().
		WithFunc(ext_misc_print_num_version_1).
		Export("ext_misc_print_num_version_1").
		NewFunctionBuilder().
		WithFunc(ext_misc_print_utf8_version_1).
		Export("ext_misc_print_utf8_version_1").
		NewFunctionBuilder().
		WithFunc(ext_misc_runtime_version_version_1).
		Export("ext_misc_runtime_version_version_1").
		NewFunctionBuilder().
		WithFunc(ext_default_child_storage_set_version_1).
		Export("ext_default_child_storage_set_version_1").
		NewFunctionBuilder().
		WithFunc(ext_default_child_storage_read_version_1).
		Export("ext_default_child_storage_read_version_1").
		NewFunctionBuilder().
		WithFunc(ext_default_child_storage_clear_version_1).
		Export("ext_default_child_storage_clear_version_1").
		NewFunctionBuilder().
		WithFunc(ext_default_child_storage_clear_prefix_version_1).
		Export("ext_default_child_storage_clear_prefix_version_1").
		NewFunctionBuilder().
		WithFunc(ext_default_child_storage_clear_prefix_version_2).
		Export("ext_default_child_storage_clear_prefix_version_2").
		NewFunctionBuilder().
		WithFunc(ext_default_child_storage_exists_version_1).
		Export("ext_default_child_storage_exists_version_1").
		NewFunctionBuilder().
		WithFunc(ext_default_child_storage_get_version_1).
		Export("ext_default_child_storage_get_version_1").
		NewFunctionBuilder().
		WithFunc(ext_default_child_storage_next_key_version_1).
		Export("ext_default_child_storage_next_key_version_1").
		NewFunctionBuilder().
		WithFunc(ext_default_child_storage_root_version_1).
		Export("ext_default_child_storage_root_version_1").
		NewFunctionBuilder().
		WithFunc(ext_default_child_storage_root_version_2).
		Export("ext_default_child_storage_root_version_2").
		NewFunctionBuilder().
		WithFunc(ext_default_child_storage_storage_kill_version_1).
		Export("ext_default_child_storage_storage_kill_version_1").
		NewFunctionBuilder().
		WithFunc(ext_default_child_storage_storage_kill_version_2).
		Export("ext_default_child_storage_storage_kill_version_2").
		NewFunctionBuilder().
		WithFunc(ext_default_child_storage_storage_kill_version_3).
		Export("ext_default_child_storage_storage_kill_version_3").
		NewFunctionBuilder().
		WithFunc(ext_allocator_free_version_1).
		Export("ext_allocator_free_version_1").
		NewFunctionBuilder().
		WithFunc(ext_allocator_malloc_version_1).
		Export("ext_allocator_malloc_version_1").
		NewFunctionBuilder().
		WithFunc(ext_hashing_blake2_128_version_1).
		Export("ext_hashing_blake2_128_version_1").
		NewFunctionBuilder().
		WithFunc(ext_hashing_blake2_256_version_1).
		Export("ext_hashing_blake2_256_version_1").
		NewFunctionBuilder().
		WithFunc(ext_hashing_keccak_256_version_1).
		Export("ext_hashing_keccak_256_version_1").
		NewFunctionBuilder().
		WithFunc(ext_hashing_sha2_256_version_1).
		Export("ext_hashing_sha2_256_version_1").
		NewFunctionBuilder().
		WithFunc(ext_hashing_twox_256_version_1).
		Export("ext_hashing_twox_256_version_1").
		NewFunctionBuilder().
		WithFunc(ext_hashing_twox_128_version_1).
		Export("ext_hashing_twox_128_version_1").
		NewFunctionBuilder().
		WithFunc(ext_hashing_twox_64_version_1).
		Export("ext_hashing_twox_64_version_1").
		NewFunctionBuilder().
		WithFunc(ext_offchain_index_set_version_1).
		Export("ext_offchain_index_set_version_1").
		NewFunctionBuilder().
		WithFunc(ext_offchain_index_clear_version_1).
		Export("ext_offchain_index_clear_version_1").
		NewFunctionBuilder().
		WithFunc(ext_offchain_local_storage_clear_version_1).
		Export("ext_offchain_local_storage_clear_version_1").
		NewFunctionBuilder().
		WithFunc(ext_offchain_is_validator_version_1).
		Export("ext_offchain_is_validator_version_1").
		NewFunctionBuilder().
		WithFunc(ext_offchain_local_storage_compare_and_set_version_1).
		Export("ext_offchain_local_storage_compare_and_set_version_1").
		NewFunctionBuilder().
		WithFunc(ext_offchain_local_storage_get_version_1).
		Export("ext_offchain_local_storage_get_version_1").
		NewFunctionBuilder().
		WithFunc(ext_offchain_local_storage_set_version_1).
		Export("ext_offchain_local_storage_set_version_1").
		NewFunctionBuilder().
		WithFunc(ext_offchain_network_state_version_1).
		Export("ext_offchain_network_state_version_1").
		NewFunctionBuilder().
		WithFunc(ext_offchain_random_seed_version_1).
		Export("ext_offchain_random_seed_version_1").
		NewFunctionBuilder().
		WithFunc(ext_offchain_submit_transaction_version_1).
		Export("ext_offchain_submit_transaction_version_1").
		NewFunctionBuilder().
		WithFunc(ext_offchain_timestamp_version_1).
		Export("ext_offchain_timestamp_version_1").
		NewFunctionBuilder().
		WithFunc(ext_offchain_sleep_until_version_1).
		Export("ext_offchain_sleep_until_version_1").
		NewFunctionBuilder().
		WithFunc(ext_offchain_http_request_start_version_1).
		Export("ext_offchain_http_request_start_version_1").
		NewFunctionBuilder().
		WithFunc(ext_offchain_http_request_add_header_version_1).
		Export("ext_offchain_http_request_add_header_version_1").
		NewFunctionBuilder().
		WithFunc(ext_storage_append_version_1).
		Export("ext_storage_append_version_1").
		NewFunctionBuilder().
		WithFunc(ext_storage_changes_root_version_1).
		Export("ext_storage_changes_root_version_1").
		NewFunctionBuilder().
		WithFunc(ext_storage_clear_version_1).
		Export("ext_storage_clear_version_1").
		NewFunctionBuilder().
		WithFunc(ext_storage_clear_prefix_version_1).
		Export("ext_storage_clear_prefix_version_1").
		NewFunctionBuilder().
		WithFunc(ext_storage_clear_prefix_version_2).
		Export("ext_storage_clear_prefix_version_2").
		NewFunctionBuilder().
		WithFunc(ext_storage_exists_version_1).
		Export("ext_storage_exists_version_1").
		NewFunctionBuilder().
		WithFunc(ext_storage_get_version_1).
		Export("ext_storage_get_version_1").
		NewFunctionBuilder().
		WithFunc(ext_storage_next_key_version_1).
		Export("ext_storage_next_key_version_1").
		NewFunctionBuilder().
		WithFunc(ext_storage_read_version_1).
		Export("ext_storage_read_version_1").
		NewFunctionBuilder().
		WithFunc(ext_storage_root_version_1).
		Export("ext_storage_root_version_1").
		NewFunctionBuilder().
		WithFunc(ext_storage_root_version_2).
		Export("ext_storage_root_version_2").
		NewFunctionBuilder().
		WithFunc(ext_storage_set_version_1).
		Export("ext_storage_set_version_1").
		NewFunctionBuilder().
		WithFunc(ext_storage_start_transaction_version_1).
		Export("ext_storage_start_transaction_version_1").
		NewFunctionBuilder().
		WithFunc(ext_storage_rollback_transaction_version_1).
		Export("ext_storage_rollback_transaction_version_1").
		NewFunctionBuilder().
		WithFunc(ext_storage_commit_transaction_version_1).
		Export("ext_storage_commit_transaction_version_1").
		NewFunctionBuilder().
		WithFunc(ext_crypto_ecdsa_generate_version_1).
		Export("ext_crypto_ecdsa_generate_version_1").
		Instantiate(ctx)

	if err != nil {
		return nil, err
	}

	code, err = decompressWasm(code)
	if err != nil {
		return nil, err
	}

	mod, err := rt.Instantiate(ctx, code)
	if err != nil {
		return nil, err
	}

	encodedHeapBase := mod.ExportedGlobal("__heap_base")
	if encodedHeapBase == nil {
		return nil, fmt.Errorf("wazero error: nil global for __heap_base")
	}

	heapBase := api.DecodeU32(encodedHeapBase.Get())
	// hb = runtime.DefaultHeapBase

	mem := mod.Memory()
	if mem == nil {
		return nil, fmt.Errorf("wazero error: nil memory for module")
	}

	instance = &Instance{
		heapBase: heapBase,
		Runtime:  rt,
		Context: &runtime.Context{
			Keystore:        cfg.Keystore,
			Validator:       cfg.Role == common.AuthorityRole,
			NodeStorage:     cfg.NodeStorage,
			Network:         cfg.Network,
			Transaction:     cfg.Transaction,
			SigVerifier:     crypto.NewSignatureVerifier(logger),
			OffchainHTTPSet: offchain.NewHTTPSet(),
		},
		Module:   mod,
		codeHash: cfg.CodeHash,
	}

	if cfg.DefaultVersion == nil {
		err = instance.version()
		if err != nil {
			return nil, fmt.Errorf("while getting runtime version: %w", err)
		}
	} else {
		instance.Context.Version = cfg.DefaultVersion
	}

	if cfg.Storage != nil {
		instance.SetContextStorage(cfg.Storage)
	}

	return instance, nil
}

var ErrExportFunctionNotFound = errors.New("export function not found")

func (i *Instance) Exec(function string, data []byte) (result []byte, err error) {
	i.Lock()
	i.Context.Allocator = allocator.NewFreeingBumpHeapAllocator(i.heapBase)

	defer func() {
		i.Context.Allocator = nil
		i.Unlock()
	}()
	// instantiate a new allocator on every execution func

	dataLength := uint32(len(data))
	inputPtr, err := i.Context.Allocator.Allocate(i.Module.Memory(), dataLength)
	if err != nil {
		return nil, fmt.Errorf("allocating input memory: %w", err)
	}

	// Store the data into memory
	mem := i.Module.Memory()
	if mem == nil {
		panic("nil memory")
	}

	ok := mem.Write(inputPtr, data)
	if !ok {
		panic("write overflow")
	}

	runtimeFunc := i.Module.ExportedFunction(function)
	if runtimeFunc == nil {
		return nil, fmt.Errorf("%w: %s", ErrExportFunctionNotFound, function)
	}

	ctx := context.WithValue(context.Background(), runtimeContextKey, i.Context)
	values, err := runtimeFunc.Call(ctx, api.EncodeU32(inputPtr), api.EncodeU32(dataLength))
	if err != nil {
		return nil, fmt.Errorf("running runtime function: %w", err)
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("no returned values from runtime function: %s", function)
	}
	wasmValue := values[0]

	outputPtr, outputLength := splitPointerSize(wasmValue)
	result, ok = mem.Read(outputPtr, outputLength)
	if !ok {
		panic("write overflow")
	}

	return result, nil
}

// Version returns the instance version.
// This is cheap to call since the instance version is cached.
// Note the instance version is set at creation and on code update.
func (in *Instance) Version() (runtime.Version, error) {
	if in.Context.Version != nil {
		return *in.Context.Version, nil
	}

	err := in.version()
	if err != nil {
		return runtime.Version{}, err
	}

	return *in.Context.Version, nil
}

// version calls runtime function Core_Version and returns the
// decoded version structure.
func (in *Instance) version() error { //skipcq: RVV-B0001
	res, err := in.Exec(runtime.CoreVersion, []byte{})
	if err != nil {
		return err
	}

	version, err := runtime.DecodeVersion(res)
	if err != nil {
		return fmt.Errorf("decoding version: %w", err)
	}

	in.Context.Version = &version
	return nil
}

// ValidateTransaction runs the extrinsic through the runtime function
// TaggedTransactionQueue_validate_transaction and returns *transaction.Validity. The error can
// be a VDT of either transaction.InvalidTransaction or transaction.UnknownTransaction, or can represent
// a normal error i.e. unmarshalling error
func (in *Instance) ValidateTransaction(e types.Extrinsic) (*transaction.Validity, error) {
	ret, err := in.Exec(runtime.TaggedTransactionQueueValidateTransaction, e)
	if err != nil {
		return nil, err
	}

	return runtime.UnmarshalTransactionValidity(ret)
}

// Metadata calls runtime function Metadata_metadata
func (in *Instance) Metadata() ([]byte, error) {
	return in.Exec(runtime.Metadata, []byte{})
}

// BabeConfiguration gets the configuration data for BABE from the runtime
func (in *Instance) BabeConfiguration() (*types.BabeConfiguration, error) {
	data, err := in.Exec(runtime.BabeAPIConfiguration, []byte{})
	if err != nil {
		return nil, err
	}

	bc := new(types.BabeConfiguration)
	err = scale.Unmarshal(data, bc)
	if err != nil {
		return nil, err
	}

	return bc, nil
}

// GrandpaAuthorities returns the genesis authorities from the runtime
func (in *Instance) GrandpaAuthorities() ([]types.Authority, error) {
	ret, err := in.Exec(runtime.GrandpaAuthorities, []byte{})
	if err != nil {
		return nil, err
	}

	var gar []types.GrandpaAuthoritiesRaw
	err = scale.Unmarshal(ret, &gar)
	if err != nil {
		return nil, err
	}

	return types.GrandpaAuthoritiesRawToAuthorities(gar)
}

// BabeGenerateKeyOwnershipProof returns the babe key ownership proof from the runtime.
func (in *Instance) BabeGenerateKeyOwnershipProof(slot uint64, authorityID [32]byte) (
	types.OpaqueKeyOwnershipProof, error) {

	// scale encoded slot uint64 + scale encoded array of 32 bytes
	const maxBufferLength = 8 + 33
	buffer := bytes.NewBuffer(make([]byte, 0, maxBufferLength))
	encoder := scale.NewEncoder(buffer)
	err := encoder.Encode(slot)
	if err != nil {
		return nil, fmt.Errorf("encoding slot: %w", err)
	}
	err = encoder.Encode(authorityID)
	if err != nil {
		return nil, fmt.Errorf("encoding authority id: %w", err)
	}

	encodedKeyOwnershipProof, err := in.Exec(runtime.BabeAPIGenerateKeyOwnershipProof, buffer.Bytes())
	if err != nil {
		return nil, fmt.Errorf("executing %s: %w", runtime.BabeAPIGenerateKeyOwnershipProof, err)
	}

	var keyOwnershipProof *types.OpaqueKeyOwnershipProof
	err = scale.Unmarshal(encodedKeyOwnershipProof, &keyOwnershipProof)
	if err != nil {
		return nil, fmt.Errorf("scale decoding key ownership proof: %w", err)
	}

	if keyOwnershipProof == nil {
		return nil, nil
	}

	return *keyOwnershipProof, nil
}

// BabeSubmitReportEquivocationUnsignedExtrinsic reports equivocation report to the runtime.
func (in *Instance) BabeSubmitReportEquivocationUnsignedExtrinsic(
	equivocationProof types.BabeEquivocationProof, keyOwnershipProof types.OpaqueKeyOwnershipProof,
) error {
	buffer := bytes.NewBuffer(nil)
	encoder := scale.NewEncoder(buffer)
	err := encoder.Encode(equivocationProof)
	if err != nil {
		return fmt.Errorf("encoding equivocation proof: %w", err)
	}
	err = encoder.Encode(keyOwnershipProof)
	if err != nil {
		return fmt.Errorf("encoding key ownership proof: %w", err)
	}
	_, err = in.Exec(runtime.BabeAPISubmitReportEquivocationUnsignedExtrinsic, buffer.Bytes())
	return err
}

// InitializeBlock calls runtime API function Core_initialise_block
func (in *Instance) InitializeBlock(header *types.Header) error {
	encodedHeader, err := scale.Marshal(*header)
	if err != nil {
		return fmt.Errorf("cannot encode header: %w", err)
	}

	_, err = in.Exec(runtime.CoreInitializeBlock, encodedHeader)
	return err
}

// InherentExtrinsics calls runtime API function BlockBuilder_inherent_extrinsics
func (in *Instance) InherentExtrinsics(data []byte) ([]byte, error) {
	return in.Exec(runtime.BlockBuilderInherentExtrinsics, data)
}

// ApplyExtrinsic calls runtime API function BlockBuilder_apply_extrinsic
func (in *Instance) ApplyExtrinsic(data types.Extrinsic) ([]byte, error) {
	return in.Exec(runtime.BlockBuilderApplyExtrinsic, data)
}

// FinalizeBlock calls runtime API function BlockBuilder_finalize_block
func (in *Instance) FinalizeBlock() (*types.Header, error) {
	data, err := in.Exec(runtime.BlockBuilderFinalizeBlock, []byte{})
	if err != nil {
		return nil, err
	}

	bh := types.NewEmptyHeader()
	err = scale.Unmarshal(data, bh)
	if err != nil {
		return nil, err
	}

	return bh, nil
}

// ExecuteBlock calls runtime function Core_execute_block
func (in *Instance) ExecuteBlock(block *types.Block) ([]byte, error) {
	// copy block since we're going to modify it
	b, err := block.DeepCopy()
	if err != nil {
		return nil, err
	}

	b.Header.Digest = types.NewDigest()

	// remove seal digest only
	for _, d := range block.Header.Digest {
		digestValue, err := d.Value()
		if err != nil {
			return nil, fmt.Errorf("getting digest type value: %w", err)
		}
		switch digestValue.(type) {
		case types.SealDigest:
			continue
		default:
			err = b.Header.Digest.Add(digestValue)
			if err != nil {
				return nil, err
			}
		}
	}

	bdEnc, err := b.Encode()
	if err != nil {
		return nil, err
	}

	return in.Exec(runtime.CoreExecuteBlock, bdEnc)
}

// DecodeSessionKeys decodes the given public session keys. Returns a list of raw public keys including their key type.
func (in *Instance) DecodeSessionKeys(enc []byte) ([]byte, error) {
	return in.Exec(runtime.DecodeSessionKeys, enc)
}

// PaymentQueryInfo returns information of a given extrinsic
func (in *Instance) PaymentQueryInfo(ext []byte) (*types.RuntimeDispatchInfo, error) {
	encLen, err := scale.Marshal(uint32(len(ext)))
	if err != nil {
		return nil, err
	}

	resBytes, err := in.Exec(runtime.TransactionPaymentAPIQueryInfo, append(ext, encLen...))
	if err != nil {
		return nil, err
	}

	dispatchInfo := new(types.RuntimeDispatchInfo)
	if err = scale.Unmarshal(resBytes, dispatchInfo); err != nil {
		return nil, err
	}

	return dispatchInfo, nil
}

// QueryCallInfo returns information of a given extrinsic
func (in *Instance) QueryCallInfo(ext []byte) (*types.RuntimeDispatchInfo, error) {
	encLen, err := scale.Marshal(uint32(len(ext)))
	if err != nil {
		return nil, err
	}

	resBytes, err := in.Exec(runtime.TransactionPaymentCallAPIQueryCallInfo, append(ext, encLen...))
	if err != nil {
		return nil, err
	}

	dispatchInfo := new(types.RuntimeDispatchInfo)
	if err = scale.Unmarshal(resBytes, dispatchInfo); err != nil {
		return nil, err
	}

	return dispatchInfo, nil
}

// QueryCallFeeDetails returns call fee details for given call
func (in *Instance) QueryCallFeeDetails(ext []byte) (*types.FeeDetails, error) {
	encLen, err := scale.Marshal(uint32(len(ext)))
	if err != nil {
		return nil, err
	}

	resBytes, err := in.Exec(runtime.TransactionPaymentCallAPIQueryCallFeeDetails, append(ext, encLen...))
	if err != nil {
		return nil, err
	}

	dispatchInfo := new(types.FeeDetails)
	if err = scale.Unmarshal(resBytes, dispatchInfo); err != nil {
		return nil, err
	}

	return dispatchInfo, nil
}

// CheckInherents checks inherents in the block verification process.
// TODO: use this in block verification process (#1873)
func (*Instance) CheckInherents() {}

// GrandpaGenerateKeyOwnershipProof returns grandpa key ownership proof from the runtime.
func (in *Instance) GrandpaGenerateKeyOwnershipProof(authSetID uint64, authorityID ed25519.PublicKeyBytes) (
	types.GrandpaOpaqueKeyOwnershipProof, error) {
	const bufferSize = 8 + 32 // authSetID uint64 + ed25519.PublicKeyBytes
	buffer := bytes.NewBuffer(make([]byte, 0, bufferSize))
	encoder := scale.NewEncoder(buffer)
	err := encoder.Encode(authSetID)
	if err != nil {
		return nil, fmt.Errorf("encoding auth set id: %w", err)
	}
	err = encoder.Encode(authorityID)
	if err != nil {
		return nil, fmt.Errorf("encoding authority id: %w", err)
	}
	encodedOpaqueKeyOwnershipProof, err := in.Exec(runtime.GrandpaGenerateKeyOwnershipProof, buffer.Bytes())
	if err != nil {
		return nil, err
	}

	var keyOwnershipProof *types.GrandpaOpaqueKeyOwnershipProof
	err = scale.Unmarshal(encodedOpaqueKeyOwnershipProof, &keyOwnershipProof)
	if err != nil {
		return nil, fmt.Errorf("scale decoding: %w", err)
	}

	if keyOwnershipProof == nil {
		return nil, nil
	}

	return *keyOwnershipProof, nil
}

// GrandpaSubmitReportEquivocationUnsignedExtrinsic reports an equivocation report to the runtime.
func (in *Instance) GrandpaSubmitReportEquivocationUnsignedExtrinsic(
	equivocationProof types.GrandpaEquivocationProof, keyOwnershipProof types.GrandpaOpaqueKeyOwnershipProof,
) error {
	buffer := bytes.NewBuffer(nil)
	encoder := scale.NewEncoder(buffer)
	err := encoder.Encode(equivocationProof)
	if err != nil {
		return fmt.Errorf("encoding equivocation proof: %w", err)
	}
	err = encoder.Encode(keyOwnershipProof)
	if err != nil {
		return fmt.Errorf("encoding key ownership proof: %w", err)
	}
	_, err = in.Exec(runtime.GrandpaSubmitReportEquivocation, buffer.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (*Instance) RandomSeed() {
	panic("unimplemented")
}
func (*Instance) OffchainWorker() {
	panic("unimplemented")
}
func (*Instance) GenerateSessionKeys() {
	panic("unimplemented")
}

// GetCodeHash returns the code of the instance
func (in *Instance) GetCodeHash() common.Hash {
	return in.codeHash
}

// NodeStorage to get reference to runtime node service
func (in *Instance) NodeStorage() runtime.NodeStorage {
	return in.Context.NodeStorage
}

// NetworkService to get referernce to runtime network service
func (in *Instance) NetworkService() runtime.BasicNetwork {
	return in.Context.Network
}

// Keystore to get reference to runtime keystore
func (in *Instance) Keystore() *keystore.GlobalKeystore {
	return in.Context.Keystore
}

// Validator returns the context's Validator
func (in *Instance) Validator() bool {
	return in.Context.Validator
}

// SetContextStorage sets the runtime's storage.
func (in *Instance) SetContextStorage(s runtime.Storage) {
	in.Lock()
	defer in.Unlock()

	if in.Context.Version == nil {
		panic("expected runtime version got nil")
	}

	runtimeStateVersion, err := trie.ParseVersion(in.Context.Version.StateVersion)
	if err != nil {
		panic(err)
	}

	s.SetVersion(runtimeStateVersion)
	in.Context.Storage = s
}

// Stop closes the WASM instance, its imports and clears
// the context allocator in a thread-safe way.
func (in *Instance) Stop() {
	in.Lock()
	defer in.Unlock()
	err := in.Runtime.Close(context.Background())
	if err != nil {
		log.Errorf("runtime failed to close: %v", err)
	}
}
