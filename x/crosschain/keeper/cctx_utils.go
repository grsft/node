package keeper

import (
	"fmt"

	cosmoserrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// UpdateNonce sets the CCTX outbound nonce to the next nonce, and updates the nonce of blockchain state.
// It also updates the PendingNonces that is used to track the unfulfilled outbound txs.
func (k Keeper) UpdateNonce(ctx sdk.Context, receiveChainID int64, cctx *types.CrossChainTx) error {
	chain := k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, receiveChainID)
	if chain == nil {
		return zetaObserverTypes.ErrSupportedChains
	}

	nonce, found := k.GetObserverKeeper().GetChainNonces(ctx, chain.ChainName.String())
	if !found {
		return cosmoserrors.Wrap(types.ErrCannotFindReceiverNonce, fmt.Sprintf("Chain(%s) | Identifiers : %s ", chain.ChainName.String(), cctx.LogIdentifierForCCTX()))
	}

	// SET nonce
	cctx.GetCurrentOutTxParam().OutboundTxTssNonce = nonce.Nonce
	tss, found := k.zetaObserverKeeper.GetTSS(ctx)
	if !found {
		return cosmoserrors.Wrap(types.ErrCannotFindTSSKeys, fmt.Sprintf("Chain(%s) | Identifiers : %s ", chain.ChainName.String(), cctx.LogIdentifierForCCTX()))
	}

	p, found := k.GetObserverKeeper().GetPendingNonces(ctx, tss.TssPubkey, receiveChainID)
	if !found {
		return cosmoserrors.Wrap(types.ErrCannotFindPendingNonces, fmt.Sprintf("chain_id %d, nonce %d", receiveChainID, nonce.Nonce))
	}

	// #nosec G701 always in range
	if p.NonceHigh != int64(nonce.Nonce) {
		return cosmoserrors.Wrap(types.ErrNonceMismatch, fmt.Sprintf("chain_id %d, high nonce %d, current nonce %d", receiveChainID, p.NonceHigh, nonce.Nonce))
	}

	nonce.Nonce++
	p.NonceHigh++
	k.GetObserverKeeper().SetChainNonces(ctx, nonce)
	k.GetObserverKeeper().SetPendingNonces(ctx, p)
	return nil
}

// GetRevertGasLimit returns the gas limit for the revert transaction in a CCTX
// It returns 0 if there is no error but the gas limit can't be determined from the CCTX data
func (k Keeper) GetRevertGasLimit(ctx sdk.Context, cctx types.CrossChainTx) (uint64, error) {
	if cctx.InboundTxParams == nil {
		return 0, nil
	}

	if cctx.InboundTxParams.CoinType == common.CoinType_Gas {
		// get the gas limit of the gas token
		fc, found := k.fungibleKeeper.GetGasCoinForForeignCoin(ctx, cctx.InboundTxParams.SenderChainId)
		if !found {
			return 0, types.ErrForeignCoinNotFound
		}
		gasLimit, err := k.fungibleKeeper.QueryGasLimit(ctx, ethcommon.HexToAddress(fc.Zrc20ContractAddress))
		if err != nil {
			return 0, errors.Wrap(fungibletypes.ErrContractCall, err.Error())
		}
		return gasLimit.Uint64(), nil
	} else if cctx.InboundTxParams.CoinType == common.CoinType_ERC20 {
		// get the gas limit of the associated asset
		fc, found := k.fungibleKeeper.GetForeignCoinFromAsset(ctx, cctx.InboundTxParams.Asset, cctx.InboundTxParams.SenderChainId)
		if !found {
			return 0, types.ErrForeignCoinNotFound
		}
		gasLimit, err := k.fungibleKeeper.QueryGasLimit(ctx, ethcommon.HexToAddress(fc.Zrc20ContractAddress))
		if err != nil {
			return 0, errors.Wrap(fungibletypes.ErrContractCall, err.Error())
		}
		return gasLimit.Uint64(), nil
	}

	return 0, nil
}

func IsPending(cctx types.CrossChainTx) bool {
	// pending inbound is not considered a "pending" state because it has not reached consensus yet
	return cctx.CctxStatus.Status == types.CctxStatus_PendingOutbound || cctx.CctxStatus.Status == types.CctxStatus_PendingRevert
}

// GetAbortedAmount returns the amount to refund for a given CCTX .
// If the CCTX has an outbound transaction, it returns the amount of the outbound transaction.
// If OutTxParams is nil or the amount is zero, it returns the amount of the inbound transaction.
// This is because there might be a case where the transaction is set to be aborted before paying gas or creating an outbound transaction.In such a situation we can refund the entire amount that has been locked in connector or TSS
func GetAbortedAmount(cctx types.CrossChainTx) sdkmath.Uint {
	if cctx.OutboundTxParams != nil && !cctx.GetCurrentOutTxParam().Amount.IsZero() {
		return cctx.GetCurrentOutTxParam().Amount
	}
	if cctx.InboundTxParams != nil {
		return cctx.InboundTxParams.Amount
	}

	return sdkmath.ZeroUint()
}
