# liquidity

Fetch and subscribe blocks events, parse events, build pool/pair liquidity in realtime.

# Subscribe events

1. 全量获取事件, 不做任何条件过滤

全量获取事件，有很多公共 RPC 节点不支持，因此，当我们需要获取全量事件时，必须使用付费节点，或者选择支持全量获取数据的节点。

    大多数节点提供商的公共节点(例如publicnode)都不支持全量获取事件，但是 L2 的官方节点肯定支持全量获取。 此外， tenderly 的公共节点支持全量获取事件，且速度不错。

2. 最快的获取事件

最快的获取事件的方式是通过 subscribe 方法订阅事件。

3. 确保不丢事件

* 启动时，从上次记录的块开始继续，因此，可能会重复处理部分事件
* 订阅流程
    a. 获取上次结束时的块号a，当前最新区块高度b; 
    b. 使用 get_logs RPC 调用获取 [a, b] 之间的所有事件
    c. 完成后，开始订阅事件
    d. 当收到第一个订阅事件(块号c)时，通过 get_logs RPC 调用获取 [b, c] 之间的所有事件;

4. Reorg

暂不考虑


**其他注意事项：**

1. RPC 节点 rate limit 限制
2. RPC get_logs 的限制
3. RPC 节点失败时的处理

# New Pools/Pairs

Watch on every factory contract, event `CreatePair` or `CreatePool`

# Liquidity

## Uniswap v2 or similar

Uniswap v2 or similar is pretty simple. We just care about the `reserves`.

Every swap or mint or burn will emit event `Sync`, so we just watch all pool's `Sync` event is Ok.

## Uniswap v3 or similar

Uniswap v3 or similar use CAMM, it's much complex than AMM.

CAMM split whole liquidity into ticks, so we should reconstruct every tick liquidity base on events.

我们在链下对每个 pool 维持一个流动性， 该流动性结构体类似如下:

```
{
    tick: int32,           // current active tick
    liquidity: uint256,    // current active liquidity
    sqrtPriceX96: uint256, // current sqrt price
    ticks: {
        [tick]: {
            tick: int,
            liquidityNet: uint256,
            liquidityGross: uint256,
        }
    },                     // all ticks liquidity map
    tickList: []int,       // just like uniswapv3 tickBitmap
                           // but more simple to iterate
}
```

当 pool 被创建时，初始化 pool 的数据。然后根据 `Mint`, `Burn`, `Swap` 事件，链下程序使用链上合约相同的逻辑，对 pool 的 liquidity, ticks liquidity 进行修改，如此，理论上即可保证链下重建的流动性与链上合约一致。


### Initialize event

初始化 pool 的 tick 和 sqrtPriceX96

### Mint event

Mint 事件如下：

```
    /// @notice Emitted when liquidity is minted for a given position
    /// @param sender The address that minted the liquidity
    /// @param owner The owner of the position and recipient of any minted liquidity
    /// @param tickLower The lower tick of the position
    /// @param tickUpper The upper tick of the position
    /// @param amount The amount of liquidity minted to the position range
    /// @param amount0 How much token0 was required for the minted liquidity
    /// @param amount1 How much token1 was required for the minted liquidity
    event Mint(
        address sender,
        address indexed owner,
        int24 indexed tickLower,
        int24 indexed tickUpper,
        uint128 amount,
        uint256 amount0,
        uint256 amount1
    );
```

根据 Mint 事件：
1. 对 pool 的 ticks map 进行更新
2. 更新当前 active liquidity

更新流程参考 `UniswapV3Pool` 合约的 `_modifyPosition`,  `_updatePosition` 逻辑

### Burn event

Burn 事件如下：
```
    /// @notice Emitted when a position's liquidity is removed
    /// @dev Does not withdraw any fees earned by the liquidity position, which must be withdrawn via #collect
    /// @param owner The owner of the position for which liquidity is removed
    /// @param tickLower The lower tick of the position
    /// @param tickUpper The upper tick of the position
    /// @param amount The amount of liquidity to remove
    /// @param amount0 The amount of token0 withdrawn
    /// @param amount1 The amount of token1 withdrawn
    event Burn(
        address indexed owner,
        int24 indexed tickLower,
        int24 indexed tickUpper,
        uint128 amount,
        uint256 amount0,
        uint256 amount1
    );
```

根据 Burn 事件：
1. 对 pool 的 ticks map 进行更新
2. 更新当前 active liquidity

更新流程参考 `UniswapV3Pool` 合约的 `_modifyPosition`, `_updatePosition` 逻辑


### Swap event

swap 不改变任何 tick 的liquidity, 只改变 pool 的 tick, liquidity, sqrtPriceX96

## Verify liquidity

如何验证通过事件重构的流动性的正确性？这是至关重要的问题。

**方法1**

当一个块的所有事件处理完后，列出该块中流动性发生变化的 pool，通过读取合约的方式构建该 pool 的流动性，并与通过事件构建的流动性对比，来验证正确性。

uniswap v3-periphery `TickLens` 合约的 `getPopulatedTicksInWord` 方法可以获取指定 tick 所在的uint256 word的流动性列表, 通过对比 tick 所在的 word 及其相邻的 word 的流动性，可以快速验证。


**方法2**

模拟一笔交易，看看得到的结果是否相同。

链上可以使用链上 `QuoterV2` 合约的 `quoteExactInput` 和 `quoteExactOutput` 来得到结果，然后与链下计算的数据做比较。 参数 amount 可以设置为 pool TVL 的 [1%, 5%, 10%, 20%, 50%, 90%] 等值来对比。
