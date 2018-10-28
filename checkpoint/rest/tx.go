package rest

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gorilla/mux"

	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/helper"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase) {
	r.HandleFunc(
		"/checkpoint/new",
		newCheckpointHandler(),
	).Methods("POST")
}

type EpochCheckpoint struct {
	RootHash        string `json:"root_hash"`
	StartBlock      uint64 `json:"start_block"`
	EndBlock        uint64 `json:"end_block"`
	ProposerAddress string `json:"proposer_address"`
}

func newCheckpointHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m EpochCheckpoint

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		err = json.Unmarshal(body, &m)
		if err != nil {
			RestLogger.Error("Error unmarshalling json epoch checkpoint", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		txBytes, err := CreateTxBytes(m)
		if err != nil {
			RestLogger.Error("Unable to create txBytes", "EndBlock", m.EndBlock, "StartBlock", m.StartBlock, "RootHash", m.RootHash)
		}

		RestLogger.Info("Sending request to Tendermint", "txBytes", hex.EncodeToString(txBytes), "url", helper.GetConfig().TendermintEndpoint)

		resp, err := SendTendermintRequest(txBytes)
		if err != nil {
			RestLogger.Error("Error while sending request to Tendermint", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		var bodyString string
		if resp.StatusCode == http.StatusOK {
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			bodyString = string(bodyBytes)
		}
		w.Write([]byte(bodyString))
	}
}

func SendTendermintRequest(txBytes []byte) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", helper.GetConfig().TendermintEndpoint, "broadcast_tx_commit"), nil)
	if err != nil {
		RestLogger.Error("Error while drafting request for tendermint", "Error", err)
		return nil, err
	}

	queryParams := req.URL.Query()
	queryParams.Add("tx", fmt.Sprintf("0x%s", hex.EncodeToString(txBytes)))
	req.URL.RawQuery = queryParams.Encode()

	return client.Do(req)
}

func CreateTxBytes(m EpochCheckpoint) ([]byte, error) {
	msg := checkpoint.NewMsgCheckpointBlock(
		m.StartBlock,
		m.EndBlock,
		common.HexToHash(m.RootHash),
		m.ProposerAddress,
	)

	tx := checkpoint.NewBaseTx(msg)

	txBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		RestLogger.Error("Error generating TX Bytes", "error", err)

		return []byte(""), err
	}
	return txBytes, nil
}
