package common

// PoolType pool type.
type PoolType uint

// [200 - 250): AMM and AMM variety
// [250, 260): Curve and variety
// [300, 350): CAMM and CAMM variety.
const (
	PoolTypeUnknown     = PoolType(0)
	PoolTypeAMM         = PoolType(200)
	PoolTypeAeroAMM     = PoolType(201)
	PoolTypeInfusionAMM = PoolType(202)
	PoolTypeCAMM        = PoolType(300)
	PoolTypeAeroCAMM    = PoolType(301)
	PoolTypePancakeCAMM = PoolType(302)
	PoolTypeReloadAll   = PoolType(9999)
)

var poolTypeNames = map[PoolType]string{
	PoolTypeUnknown:     "Unknown",
	PoolTypeAMM:         "Standard AMM",
	PoolTypeAeroAMM:     "Aero AMM",
	PoolTypeInfusionAMM: "Infusion AMM",
	PoolTypeCAMM:        "Standard CAMM",
	PoolTypeAeroCAMM:    "Aero CAMM",
	PoolTypePancakeCAMM: "Pancake CAMM",
}

var validPoolTypes = map[PoolType]bool{
	PoolTypeUnknown:     true,
	PoolTypeAMM:         true,
	PoolTypeAeroAMM:     true,
	PoolTypeInfusionAMM: true,
	PoolTypeCAMM:        true,
	PoolTypeAeroCAMM:    true,
	PoolTypePancakeCAMM: true,
}

// SwapFactory swap factory.
type SwapFactory struct {
	Address         string   `toml:"address"`
	Name            string   `toml:"name"`
	Typ             PoolType `toml:"pool_type"`
	Fee             int      `toml:"fee"`
	StableFee       int      `toml:"stable_fee"`       // for areodrome v2
	PositionManager string   `toml:"position_manager"` // used for load all pools
	Known           bool
	Pools           uint32
}

// Token ERC20 token.
type Token struct {
	Address     string `json:"address"`
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Decimals    uint   `json:"decimals"`
	PoolCount   uint
	MaxTVLPools []string
}

// IsValidPoolType pool type is valid.
func IsValidPoolType(pt PoolType) bool {
	return validPoolTypes[pt]
}

// String pool type readable name.
func (pt PoolType) String() string {
	return poolTypeNames[pt]
}

// IsAMMVariety pool is AMM variety.
func (pt PoolType) IsAMMVariety() bool {
	return pt >= 200 && pt < 250
}

// IsAMMVariety pool is CAMM or CAMM variety.
func (pt PoolType) IsCAMMVariety() bool {
	return pt >= 300 && pt < 350
}
