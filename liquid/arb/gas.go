package arb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"starbase.ag/liquidity/pkg/logger"
)

const (
	owlracleAPIKEY = "d1119011285240eab8d98d9580a07338"
	// owlracleAPISECRET = "94ef5dd57a1a442a96653d715da4ccf6"
	owlracleEndpoint = "https://api.owlracle.info/v4/%s/gas?apikey=%s"
)

type GasPrice struct {
	Acceptance           float64 `json:"acceptance"`
	MaxFeePerGas         float64 `json:"maxFeePerGas"`
	MaxPriorityFeePerGas float64 `json:"maxPriorityFeePerGas"`
	BaseFee              float64 `json:"baseFee"`
	EstimatedFee         float64 `json:"estimatedFee"`
}

// {"timestamp":"2024-10-31T06:37:49.270Z","lastBlock":21783659,
// "avgTime":2.0100502512562812,"avgTx":157.745,"avgGas":232767.86301263285,"avgL1Fee":0.0009805774589231432,
// "speeds":[
// {"acceptance":0.35,"maxFeePerGas":0.00547022,"maxPriorityFeePerGas":0.0010733779999999998,"baseFee":0.004396842,"estimatedFee":0.0033699058040512816},
// {"acceptance":0.6,"maxFeePerGas":0.005843164,"maxPriorityFeePerGas":0.0014192199999999997,"baseFee":0.004423944,"estimatedFee":0.0035996563717041554},
// {"acceptance":0.9,"maxFeePerGas":0.006303545,"maxPriorityFeePerGas":0.0018515829999999995,"baseFee":0.004451962,"estimatedFee":0.0038832721319432195},
// {"acceptance":1,"maxFeePerGas":0.009530283,"maxPriorityFeePerGas":0.005070314,"baseFee":0.004459969,"estimatedFee":0.005871090375880909}]}
func FetchGasPrice(network string) {
	resp, err := http.Get(fmt.Sprintf(owlracleEndpoint, network, owlracleAPIKEY))
	if err != nil {
		logger.Warn().Msgf("get gas price failed: %v", err)
		return
	}
	defer resp.Body.Close()
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Warn().Msgf("read gas response failed: %v", err)
		return
	}
	var v map[string]interface{}
	if err := json.Unmarshal(buf, &v); err != nil {
		logger.Warn().Msgf("unmarshal gas response failed: %v %v", err, string(buf))
		return
	}
	logger.Info().Msgf("gas: %v", v)
}
