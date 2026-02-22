package pool

import (
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

const (
	// uniswap v2 SYNC event.
	TopicSYNC     = "0x1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1"
	TopicSYNCAero = "0xcf2aa50876cdfbb541206f89af0ee78d44a2abf8d328e37fa4917f982149848a"
	// uniswap v2 CreatePair event.
	TopicPairCreated = "0x0d3648bd0f6ba80134a33ba9275ac585d9d315f0ad8355cddefde31afa28d0e9"
	// aero v2 create pool.
	TopicAeroV2PoolCreated = "0x2128d88d14c80cb081c1252a5acff7a264671bf199ce226b53788fb26065005e"
	// infusion/baso.finance pair created.
	TopicInfusionPairCreated = "0xc4805696c66d7cf352fc1d6bb633ad5ee82f6cb577c453024b6e0eb8306c6fc9"

	TopicAeroSwapV2 = "0xb3e2773606abfd36b5bd91394b3a54d1398336c65005baf7bf7a05efeffaf75b"
	TopicSwapV2     = "0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822"

	// uniswap v3 CreatePool event.
	TopicPoolCreated = "0x783cca1c0412dd0d695e784568c96da2e9c22ff989357a2e8b1d9b2b4e6b7118"
	// uniswap v3 Initialize event.
	TopicInitialize = "0x98636036cb66a9c19a37435efc1e90142190214e8abeb821bdba3f2990dd4c95"
	// uniswap v3 Mint event.
	TopicMint = "0x7a53080ba414158be7ec69b987b5fb7d07dee101fe85488f0853ae16239d0bde"
	// uniswap v3 Burn event.
	TopicBurn = "0x0c396cd989a39f4459b5fa1aed6a9a8dcdbc45908acfd67e028cd568da98982c"
	// uniswap v3 Swap event.
	TopicSwap = "0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67"
	//
	TopicCollect = "0x70935338e69775456a85ddef226c395fb668b63fa0115f5f20610b388e6ca9c0"
	// pancake v3 swap.
	TopicPancakeSwap = "0x19b47279256b2a23a1665c810c8d55a1758940ee09377d4f8d26497a3577dc83"
)

var (
	// 0xae468ce586f9a87660fdffc1448cee942042c16ae2f02046b134b5224f31936b.
	TopicAeroV2SetFee = strings.ToLower(crypto.Keccak256Hash([]byte(`SetCustomFee(address,uint256)`)).String())
	// SetCustomFee 0xd444e1b10a2a0c61e10ee9f0167820955df343074f16b69614952caef34de21d.
	TopicAeroV3SetFee = strings.ToLower(crypto.Keccak256Hash([]byte(`SetCustomFee(address,uint24)`)).String())
)
