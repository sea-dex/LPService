// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@uniswap/v2-core/contracts/interfaces/IUniswapV2Pair.sol";
import "@uniswap/v3-core/contracts/interfaces/IUniswapV3Pool.sol";
import "@uniswap/v3-core/contracts/libraries/TickMath.sol";
import "@uniswap/v3-core/contracts/libraries/FullMath.sol";

contract ArbitrageContract {
    address public owner;

    constructor() {
        owner = msg.sender;
    }

    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }

    function executeArbitrage(
        address[] memory path,
        uint256 amountIn,
        uint256 minAmountOut
    ) external onlyOwner {
        require(path.length >= 3, "Path must have at least 3 tokens");

        address firstToken = path[0];
        address lastToken = path[path.length - 1];
        
        // Start the flash swap from the first pool
        address firstPoolAddress = getPoolAddress(path[0], path[1]);
        bool isFirstPoolV3 = isUniswapV3Pool(firstPoolAddress);

        if (isFirstPoolV3) {
            flashSwapV3(firstPoolAddress, path[0], path[1], amountIn, path, minAmountOut);
        } else {
            flashSwapV2(firstPoolAddress, path[0], path[1], amountIn, path, minAmountOut);
        }

        // Transfer the profit to the caller
        uint256 profit = IERC20(lastToken).balanceOf(address(this));
        require(profit >= minAmountOut, "Slippage too high");
        IERC20(lastToken).transfer(msg.sender, profit);
    }

    function flashSwapV2(
        address poolAddress,
        address tokenIn,
        address tokenOut,
        uint256 amountIn,
        address[] memory path,
        uint256 minAmountOut
    ) internal {
        IUniswapV2Pair pair = IUniswapV2Pair(poolAddress);
        
        (uint amount0Out, uint amount1Out) = tokenIn < tokenOut 
            ? (uint(0), amountIn) 
            : (amountIn, uint(0));

        bytes memory data = abi.encode(path, minAmountOut);
        pair.swap(amount0Out, amount1Out, address(this), data);
    }

    function flashSwapV3(
        address poolAddress,
        address tokenIn,
        address tokenOut,
        uint256 amountIn,
        address[] memory path,
        uint256 minAmountOut
    ) internal {
        IUniswapV3Pool pool = IUniswapV3Pool(poolAddress);
        
        bool zeroForOne = tokenIn < tokenOut;

        bytes memory data = abi.encode(path, minAmountOut);

        pool.flash(
            address(this),
            zeroForOne ? amountIn : 0,
            zeroForOne ? 0 : amountIn,
            data
        );
    }

    function uniswapV2Call(address sender, uint amount0, uint amount1, bytes calldata data) external {
        require(msg.sender == getPoolAddress(IUniswapV2Pair(msg.sender).token0(), IUniswapV2Pair(msg.sender).token1()), "Unauthorized");
        
        (address[] memory path, uint256 minAmountOut) = abi.decode(data, (address[], uint256));
        
        uint256 amountIn = amount0 > 0 ? amount0 : amount1;
        address tokenIn = amount0 > 0 ? IUniswapV2Pair(msg.sender).token0() : IUniswapV2Pair(msg.sender).token1();

        // Execute the arbitrage
        uint256 currentAmount = executeArbitrageSteps(path, amountIn, 1);

        // Repay the flash swap
        uint256 amountToRepay = amountIn + ((amountIn * 3) / 997) + 1; // 0.3% fee
        IERC20(tokenIn).transfer(msg.sender, amountToRepay);
    }

    function uniswapV3FlashCallback(uint256 fee0, uint256 fee1, bytes calldata data) external {
        require(msg.sender == getPoolAddress(IUniswapV3Pool(msg.sender).token0(), IUniswapV3Pool(msg.sender).token1()), "Unauthorized");
        
        (address[] memory path, uint256 minAmountOut) = abi.decode(data, (address[], uint256));
        
        uint256 amountIn = fee0 > 0 ? fee0 : fee1;
        address tokenIn = fee0 > 0 ? IUniswapV3Pool(msg.sender).token0() : IUniswapV3Pool(msg.sender).token1();

        // Execute the arbitrage
        uint256 currentAmount = executeArbitrageSteps(path, amountIn, 1);

        // Repay the flash swap
        uint256 amountToRepay = amountIn + fee0 + fee1;
        IERC20(tokenIn).transfer(msg.sender, amountToRepay);
    }

    function executeArbitrageSteps(address[] memory path, uint256 amountIn, uint256 startIndex) internal returns (uint256) {
        uint256 currentAmount = amountIn;

        for (uint i = startIndex; i < path.length - 1; i++) {
            address tokenIn = path[i];
            address tokenOut = path[i + 1];
            
            address poolAddress = getPoolAddress(tokenIn, tokenOut);
            bool isV3 = isUniswapV3Pool(poolAddress);

            if (isV3) {
                currentAmount = swapV3(poolAddress, tokenIn, tokenOut, currentAmount);
            } else {
                currentAmount = swapV2(poolAddress, tokenIn, tokenOut, currentAmount);
            }
        }

        return currentAmount;
    }

    function swapV2(
        address poolAddress,
        address tokenIn,
        address tokenOut,
        uint256 amountIn
    ) internal returns (uint256 amountOut) {
        IUniswapV2Pair pair = IUniswapV2Pair(poolAddress);
        (uint reserve0, uint reserve1,) = pair.getReserves();
        
        (uint reserveIn, uint reserveOut) = tokenIn < tokenOut 
            ? (reserve0, reserve1) 
            : (reserve1, reserve0);

        amountOut = getAmountOut(amountIn, reserveIn, reserveOut);

        (uint amount0Out, uint amount1Out) = tokenIn < tokenOut 
            ? (uint(0), amountOut) 
            : (amountOut, uint(0));

        IERC20(tokenIn).transfer(poolAddress, amountIn);
        pair.swap(amount0Out, amount1Out, address(this), new bytes(0));

        return amountOut;
    }

    function swapV3(
        address poolAddress,
        address tokenIn,
        address tokenOut,
        uint256 amountIn
    ) internal returns (uint256 amountOut) {
        IUniswapV3Pool pool = IUniswapV3Pool(poolAddress);
        
        bool zeroForOne = tokenIn < tokenOut;

        (uint160 sqrtPriceX96, int24 tick,,,,,) = pool.slot0();

        uint160 sqrtPriceLimitX96 = zeroForOne 
            ? TickMath.MIN_SQRT_RATIO + 1 
            : TickMath.MAX_SQRT_RATIO - 1;

        IERC20(tokenIn).transfer(poolAddress, amountIn);

        (int256 amount0, int256 amount1) = pool.swap(
            address(this),
            zeroForOne,
            int256(amountIn),
            sqrtPriceLimitX96,
            abi.encode(tokenOut)
        );

        amountOut = uint256(-(zeroForOne ? amount1 : amount0));

        return amountOut;
    }

    function getAmountOut(
        uint amountIn, 
        uint reserveIn, 
        uint reserveOut
    ) internal pure returns (uint amountOut) {
        uint amountInWithFee = amountIn * 997;
        uint numerator = amountInWithFee * reserveOut;
        uint denominator = reserveIn * 1000 + amountInWithFee;
        amountOut = numerator / denominator;
    }

    function getPoolAddress(address token0, address token1) internal view returns (address) {
        // Implement logic to get pool address
        // This could involve looking up in a mapping or calling a factory contract
    }

    function isUniswapV3Pool(address poolAddress) internal view returns (bool) {
        // Implement logic to determine if the pool is V3
        // This could involve trying to call a V3-specific function and catching the revert
    }

    function uniswapV3SwapCallback(
        int256 amount0Delta,
        int256 amount1Delta,
        bytes calldata data
    ) external {
        address tokenOut = abi.decode(data, (address));
        uint256 amountToPay = uint256(amount0Delta > 0 ? amount0Delta : amount1Delta);
        IERC20(tokenOut).transfer(msg.sender, amountToPay);
    }

    function withdrawToken(address token, uint256 amount) external onlyOwner {
        IERC20(token).transfer(msg.sender, amount);
    }
}
