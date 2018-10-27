package checkpoint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgCheckpoint:
			// redirect to handle msg checkpoint
			return handleMsgCheckpoint(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("Invalid message in checkpoint module").Result()
		}
	}
}

func handleMsgCheckpoint(ctx sdk.Context, msg MsgCheckpoint, k Keeper) sdk.Result {
	// TODO add validation
	// if err := msg.ValidateBasic(); err != nil { // return failed  }

	// check if the roothash provided is valid for start and end
	valid := validateCheckpoint(int(msg.StartBlock), int(msg.EndBlock), msg.RootHash.String())

	// check msg.proposer with tm proposer
	var key int64
	if valid {

		// add checkpoint to state if rootHash matches
		key = k.AddCheckpoint(ctx, msg.StartBlock, msg.EndBlock, msg.RootHash, msg.Proposer)
		CheckpointLogger.Debug("RootHash matched!", "key", key)
	} else {
		CheckpointLogger.Debug("Root hash doesn't match ;(")
		// return Bad Block Error
		return ErrBadBlockDetails(k.codespace).Result()
	}

	// send tags
	return sdk.Result{}
}
