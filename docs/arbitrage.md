## from v2 to v2

x: ETH
y: USD

```
dy1 = dx1 * f1 * y1 / (x1 + dx1*f1)

dx2 = dy1 * f2 * x2 / (y2 + dy1*f2) 
    = dx1*f1*f2*y1*x2 / [(y2 + dy1*f2) (x1 + dx1*f1)]
    = dx1*f1*f2*y1*x2 / [(y2 + dx1 * f1 * y1 *f2 / (x1 + dx1*f1)) * (x1 + dx1*f1) ]
    = dx1*f1*f2*y1*x2 / [y2*x1 + dx1*(y2*f1 + y1*f1*f2)]
    
    f = dx1*f1*f2*y1*x2
    g = [y2*x1 + dx1*(y2*f1 + y1*f1*f2)]
    dx2= f/g

k1 = f1*f2*y1*x2
k2 = y2*x1
k3 = y2*f1 + y1*f1*f2

profit(dx1) = dx2 - dx1 

profit' = dx2' - 1 = 0

dx2' = (f'g-fg')/ g*g

(f'g-fg') = 
= f1*f2*y1*x2 * [y2*x1 + dx1*(y2*f1 + y1*f1*f2)] - dx1*f1*f2*y1*x2 * (y2*f1 + y1*f1*f2)
= f1*f2*y1*x2*y2*x1

k = y2*f1 + y1*f1*f2
dx2' = f1*f2*y1*x2*y2*x1 / sqr(y2*x1 + dx1*k) = 1

sqr(y2*x1 + dx1*k) - f1*f2*y1*x2*y2*x1 = 0

k*k * dx1**2 + 2*k*y2*x1 * dx1 + y2*x1*(y2*x1 - f1*f2*y1*x2) = 0

k = y2*f1 + y1*f1*f2
a = k*k
b = 2*k*y2*x1
c = y2*x1*(y2*x1 - f1*f2*y1*x2)
r = (-b + sqrt(b**2 - 4*a*c)) / (2*a)

```

## v3 formula
### 0->1

卖出token0, 买入token1, 价格降低

```
SPNext = L*Q96*SP / (L*Q96 + dx*f*SP)
dy = L*(Sp-SPNext) / Q96
dy = dx*f*SP*SP*L / (Q192*L + dx*Q96*f*SP)
```

### 1->0
```
SPNext = dy*f*Q96/L + SP
dx = L*Q96*(SPNext - SP)/SP/SPNext
dx = dy*f*Q192*L / (Sp*Sp*L + dy*Q96*f*Sp)
```


## from v2 to v3

x: ETH
y: USD

### v2(1->0), v3(0->1)
```
dx1 = dy1*r1*x1 / (y1+dy1*r1)
Spnext = L*Q96*SP/(L*Q96+dx1*r2*SP)
dy2 = L*(Sp-Spnext)/Q96
    = L*(Sp - L*Q96*SP/(L*Q96+dx1*r2*SP))/Q96
    = L*(SP*(L*Q96+dx1*r2*SP) - L*Q96*SP)/(Q96*(L*Q96+dx1*r2*SP))
    = dx1*L*SP*r2*SP / (Q96*(L*Q96+dx1*r2*SP))
    = dx1*L*SP*r2*SP / (Q192*L+dx1*Q96*r2*SP)

    = dy1*r1*x1 * L*SP*r2*SP / [ (y1+dy1*r1) * (Q96*(L*Q96+dx1*r2*SP)) ]
    = dy1*r1*r2*x1*L*SP*SP / [ (y1+dy1*r1) * Q96 * (L*Q96+dx1*r2*SP) ]
    = dy1*r1*r2*x1*L*SP*SP / [ (y1+dy1*r1) * Q96 * (L*Q96 + r2*SP * dy1*r1*x1 / (y1+dy1*r1)) ]
    = dy1*r1*r2*x1*L*SP*SP / [ Q96 * (y1+dy1*r1) * (L*Q96 + r2*SP * dy1*r1*x1 / (y1+dy1*r1)) ]
    = dy1*r1*r2*x1*L*SP*SP / [ Q96 * (y1+dy1*r1) * L*Q96 + Q96 * r2*SP * dy1*r1*x1 ]
    = dy1*r1*r2*x1*L*SP*SP / [ Q192 * (y1+dy1*r1) * L + dy1* Q96 * r1*r2*SP*x1 ]
    = dy1*r1*r2*x1*L*SP*SP / [ Q192*y1*L + dy1*r1*L*Q192 + dy1* Q96*r1*r2*SP*x1 ]
    = dy1*r1*r2*x1*L*SP*SP / [ Q192*y1*L + dy1* (r1*L*Q192 + Q96*r1*r2*SP*x1) ]

如果超过当前的 liquidity:

dy2 = (dy1-N)*r1*r2*x1*L*SP*SP / [ Q192*y1*L + (dy1-N)* (r1*L*Q192 + Q96*r1*r2*SP*x1) ]
k1 = r1*r2*x1*L*SP*SP
k2 = Q192*y1*L
k3 = (r1*L*Q192 + Q96*r1*r2*SP*x1)
f = (dy1-N) * k1
g = k2 + (dy1-N)*k3

f'g - fg' = k1*g - f*k3 = k1*(k2 + (dy1-N)*k3) - (dy1-N) * k1 * k3 = k1*k2

k2' = k2 - N*k3


dy2 = (dx1-N)*L*SP*r2*SP / (Q192*L+(dx1-N)*Q96*r2*SP)
    = (dx1*L*SP*r2*SP - N*L*SP*r2*SP) / (Q192*L+(dx1-N)*Q96*r2*SP)
    = (dy1*r1*x1*L*SP*r2*SP - N*L*SP*r2*SP*(y1+dy1*r1)) / (y1+dy1*r1)*(Q192*L-N*Q96*r2*SP+dx1*Q96*r2*SP)
    = (dy1*r1*x1*L*SP*r2*SP - N*L*SP*r2*SP*(y1+dy1*r1)) / (y1+dy1*r1)*(Q192*L-N*Q96*r2*SP+dx1*Q96*r2*SP)
    = (dy1*r1*x1*L*SP*r2*SP - N*L*SP*r2*SP*(y1+dy1*r1)) / [(y1+dy1*r1)*(Q192*L-N*Q96*r2*SP)+dy1*r1*x1*Q96*r2*SP]
    = (dy1*r1*x1*L*SP*r2*SP - N*L*SP*r2*SP*y1 - dy1*r1*N*L*SP*r2*SP) / [(y1+dy1*r1)*(Q192*L-N*Q96*r2*SP)+dy1*r1*x1*Q96*r2*SP]
    = (dy1*r1*x1*L*SP*r2*SP - N*L*SP*r2*SP*y1 - dy1*r1*N*L*SP*r2*SP) / [y1*(Q192*L-N*Q96*r2*SP)+dy1*(r1*(Q192*L-N*Q96*r2*SP)+r1*x1*Q96*r2*SP)]
    = (dy1*(r1*x1*L*SP*r2*SP - r1*N*L*SP*r2*SP) - N*L*SP*r2*SP*y1) / [y1*(Q192*L-N*Q96*r2*SP)+dy1*(r1*(Q192*L-N*Q96*r2*SP)+r1*x1*Q96*r2*SP)]

k1=r1*SP*r2*SP*L(x1-N)
k2=y1*(Q192*L-N*Q96*r2*SP)
k3=r1*(Q192*L-N*Q96*r2*SP+x1*Q96*r2*SP)
k4=N*L*SP*r2*SP*y1

dy2 = (dy1*k1-k4)/(k2+dy1*k3)
f = dy1*k1-k4
g = k2 + dy1*k3
f'g - fg' = k1*g - f*k3 = k1*(k2 + dy1*k3) - (dy1*k1 -k4) * k3 = k1*k2 + k3*k4

(f'g - fg')/g*g = 1
g*g =  k3*k3*dy1*dy1 + 2k2*k3*dy1 + k2*k2
k3*k3*dy1*dy1 + 2k2*k3*dy1 + k2*k2 - k1*k2 - k3*k4 = 0

a = k3*k3
b = 2*k2*k3
c = k2*k2 - k1*k2 - k3*k4
r = (-b + sqrt(b**2 - 4*a*c)) / (2*a)
---------

dy1 = dx1 * f1 * y1 / (x1 + dx1*f1)

spNext = L * Q96 * SP / (L * Q96 + dy1 * f2 * SP)
dx2 = L * (Sp-spNext) / Q96 = (L/Q96) * (Sp - [L * Q96 * SP / (L * Q96 + dy1 * f2 * SP)])
    = (L/Q96) * Sp*dy1*f2*Sp/(L * Q96 + dy1 * f2 * SP)
    = (L*Sp*dy1*f2*Sp)/ (Q96* L * Q96 + dy1 * f2 * SP * Q96)
    = (dx1*L*Sp*f1*y1*f2*Sp) / [(x1 + dx1*f1) * (Q96 * L * Q96 + dy1 * f2 * SP * Q96)]
    = (dx1*L*Sp*f1*y1*f2*Sp) / [(x1 + dx1*f1) * Q96 * L * Q96 + (x1 + dx1*f1) * dy1 * f2 * SP * Q96]
    = (dx1*L*Sp*f1*y1*f2*Sp) / [(x1 + dx1*f1) * Q96 * L * Q96 + dx1 * f1 * y1 * f2 * SP * Q96]
    = (dx1*L*Sp*f1*y1*f2*Sp) / [x1 * Q96 * L * Q96 + dx1*(f1 * Q96 * L * Q96 +  f1 * Q96 * y1 * f2 * SP)]

k1 = L*Sp*f1*y1*f2*Sp
k2 = x1 * Q96 * L * Q96
k3 = (f1 * Q96 * L * Q96 +  f1 * Q96 * y1 * f2 * SP)
dx2 = dx1*k1 / (k2 + dx1*k3) = f/g
f = dx1*k1
g = k2 + dx1*k3

dx2' = f'g - fg'/g**2 

f'g - fg' = k1 * (k2+dx1*k3) - dx1*k1* k3 = k1*k2

dx2' = k1*k2 / (k2 + dx1*k3)**2 = 1

(k2 + dx1*k3)**2 - k1*k2 = 0
k3**2 * dx1**2 + 2*k2*k3*dx1 + k2*(k2-k1) = 0

a = k3**2
b = 2*k2*k3
c = k2*(k2-k1)
r = (-b + sqrt(b**2 - 4*a*c)) / (2*a)
```

### v2 0-1, v3(1->0)
```
dy1 = dx1 * f1 * y1 / (x1 + dx1*f1)

# 1->0
SpNext = SP + dy1*f2*Q96/L
dx2 = L * Q96 * (dy1*f2*Q96/L)/Sp/SpNext = dy1*f2*Q192/Sp/SpNext
    = dy1*f2*Q192/ (Sp*Sp + dy1*f2*Q96*Sp/L)

    = dx1 * f1 * y1 * f2 * Q192 / (x1 + dx1*f1) * (Sp*Sp + dy1*f2*Q96*Sp/L)
    = dx1 * f1 * y1 * f2 * Q192 / (x1 + dx1*f1) * (Sp*Sp + dx1 * f1 * y1 *f2*Q96*Sp / (x1 + dx1*f1) * L)
    = dx1 * f1 * y1 * f2 * Q192 / [(x1 + dx1*f1) * Sp*Sp + dx1 * f1 * y1 *f2*Q96*Sp/L]
    = dx1 * f1 * y1 * f2 * Q192 / [x1  * Sp*Sp + dx1*f1*Sp * (Sp+y1*f2*Q96/L)]
    = dx1 * f1 * y1 * f2 * Q192 * L / [x1  * Sp*Sp * L + dx1*f1*Sp * (Sp*L+y1*f2*Q96)]
    = dx1*f1*y1*f2*L*Q192 / [x1*Sp*Sp*L + dx1 * (f1*Sp*Sp*L + f1*Sp*y1*f2*Q96)]

dx2 = (dx1*f1*f2*L2*Q192*(y1-N) - N*x1*f2*L2*Q192) / (Sp2*Sp2*L2*x1 - N*x1*Q96*f2*Sp2 + dx1*f1*Sp2*(Sp2*L2 + (y1-N)*Q96*f2))

k1 = f1*y1*f2*L*Q192
k2 = x1*Sp*Sp*L
k3 = f1*Sp*Sp*L + f1*Sp*y1*f2*Q96
dx2 = dx1*k1 / (k2 + dx1*k3) = f/g
f = dx1*k1
g = k2 + dx1*k3

dx2' = f'g - fg'/g**2 

f'g - fg' = k1 * (k2+dx1*k3) - dx1*k1* k3 = k1*k2

dx2' = k1*k2 / (k2 + dx1*k3)**2 = 1

(k2 + dx1*k3)**2 - k1*k2 = 0
k3**2 * dx1**2 + 2*k2*k3*dx1 + k2*(k2-k1) = 0

a = k3**2
b = 2*k2*k3
c = k2*(k2-k1)
r = (-b + sqrt(b**2 - 4*a*c)) / (2*a)
```

## from v3 to v2
### 0->1
```
spNext = L*Q96*SP / (L*Q96 + dx1*f1*SP)
dy1 = L * (Sp-spNext) / Q96
    = L * (Sp - L*Q96*SP / (L*Q96 + dx1*f1*SP)) / Q96
    = L*Sp*dx1*f1*SP / [Q96 * (L * Q96 + dx1 * f1 * SP)]
    = dx1 * L*Sp*f1*SP / (Q192*L + dx1*Q96*f1*SP)

dx2 = dy1 * f2 * x2 / (y2 + dy1*f2)
    = dx1 * L*Sp*f1*SP * f2 * x2 / [(Q96*L*Q96 + dx1 * Q96*f1*SP) * (y2 + dy1*f2)]
    = dx1 * L*Sp*f1*SP * f2 * x2 / [(Q192*L + dx1 * Q96*f1*SP) * (y2 + dx1 * L*Sp*f1*SP*f2 / (Q96*L*Q96 + dx1 * Q96*f1*SP))]
    = dx1 * L*Sp*f1*SP * f2 * x2 / [(Q192*L + dx1 * Q96*f1*SP) * y2 + dx1 * L*Sp*f1*SP*f2]
    = dx1 * L*Sp*f1*SP * f2 * x2 / (Q192*L * y2 + dx1 * Q96*f1*SP * y2 + dx1 * L*Sp*f1*SP*f2)
    = dx1 * L*Sp*f1*SP * f2 * x2 / [Q192*L*y2 + dx1 * (Q96*f1*SP*y2 + L*Sp*f1*SP*f2)]

k1 = L*Sp*f1*SP*f2*x2
k2 = Q192*L*y2
k3 = (Q96*f1*SP*y2 + L*Sp*f1*SP*f2)
```

### v3(1->0), v2(0->1)

```
v2 (0->1), v2 (1->0)
dx2 = dx1*f1*f2*y1*x2 / [y2*x1 + dx1*(y2*f1 + y1*f1*f2)]
v2 (1->0), v2 (0->1)
dy2 = dy1*f1*f2*x1*y2 / [x2*y1 + dy1 * (x2*f1 + x1*f1*f2)]
```

```
spNext = SP + dy1 * f1 * Q96/L
dx1 = L * Q96 * (dy1*f1*Q96/L)/Sp/SpNext
    = dy1*f1*Q192/ (Sp*Sp + dy1*f1*Q96*Sp/L)

dy2 = dx1*f2*y2 / (x2 + dx1*f2)
    = dy1*f1*Q192* f2 * y2 / [(x2 + dx1*f2) * (Sp*Sp + dy1*f1*Q96*Sp/L)]
    = dy1*f1*Q192*f2*y2 / [x2*(Sp*Sp + dy1*f1*Q96*Sp/L) + dx1*f2*(Sp*Sp + dy1*f1*Q96*Sp/L)]
    = dy1*f1*Q192*f2*y2 / [x2*(Sp*Sp + dy1*f1*Q96*Sp/L) + dy1*f1*Q192*f2]
    = dy1*f1*Q192*f2*y2 / [x2*Sp*Sp + dy1*f1*Q96*Sp*x2/L + dy1*f1*Q192*f2]
    = dy1*f1*Q192*f2*y2*L / [x2*Sp*Sp*L + dy1*f1*Q96*Sp*x2 + dy1*f1*Q192*f2*L]
    = dy1*f1*Q192*f2*y2*L / [x2*Sp*Sp*L + dy1*(f1*Q96*Sp*x2 + f1*Q192*f2*L)]

k1 = f1*Q192*f2*y2*L
k2 = x2*Sp*Sp*L
k3 = f1*Q96*Sp*x2 + f1*Q192*f2*L
```

## from v3 to v3

### 0->1, then 1->0

```
0->1:
spNext1 = L1 * Q96 * SP1 / (L1 * Q96 + dx1 * f1 * SP1)
dy1 = L1 * (Sp1-spNext1) / Q96
    = (L1*Sp1*dx1*f1*Sp1)/ (Q96* L1 * Q96 + dx1 * f1 * SP1 * Q96)
    = (dx1*L1*Sp1*f1*Sp1)/ (Q192* L1 + dx1 * f1 * SP1 * Q96)

1->0:
spNext2 = SP2 + dy1 * f2 * Q96/L2
dx2 = L2 * Q96 * (dy1*f2*Q96/L2)/Sp2/SpNext2
    = dy1*f2*Q192/ (Sp2*Sp2 + dy1*f2*Q96*Sp2/L2)
    = dy1*f2*Q192*L2/ (Sp2*Sp2*L2 + dy1*f2*Q96*Sp2)
    = dy1*f2*Q192*L2/ (Sp2*Sp2*L2 + dy1*f2*Q96*Sp2)
    = (dx1*L1*Sp1*f1*Sp1) * f2*Q192*L2 / [(Sp2*Sp2*L2 + dy1*f2*Q96*Sp2) * (Q192* L1 + dx1 * f1 * SP1 * Q96)]
    = (dx1*L1*Sp1*f1*Sp1) * f2*Q192*L2 / [Sp2*Sp2*L2* (Q192* L1 + dx1 * f1 * SP1 * Q96) + dy1*f2*Q96*Sp2 * (Q192* L1 + dx1 * f1 * SP1 * Q96)]
    = (dx1*L1*Sp1*f1*Sp1) * f2*Q192*L2 / [Sp2*Sp2*L2* (Q192* L1 + dx1 * f1 * SP1 * Q96) + (dx1*L1*Sp1*f1*Sp1) *f2*Q96*Sp2 ]
    = (dx1*L1*Sp1*f1*Sp1) * f2*Q192*L2 / [Sp2*Sp2*L2* Q192* L1 + dx1 * Sp2*Sp2*L2* f1 * SP1 * Q96 + dx1 * L1*Sp1*f1*Sp1 *f2*Q96*Sp2 ]
    = (dx1*L1*Sp1*f1*Sp1) * f2*Q96*L2 / [Sp2*Sp2*L2* Q96* L1 + dx1 * Sp2*Sp2*L2* f1 * SP1  + dx1 * L1*Sp1*f1*Sp1 *f2*Sp2]
    = dx1 * (L1*Sp1*f1*Sp1*f2*Q96*L2) / [Sp2*Sp2*L2* Q96* L1 + dx1 * (Sp2*Sp2*L2*f1*SP1  + L1*Sp1*f1*Sp1*f2*Sp2)]

dy1 = dx1*L1*f1*SP1*SP1 / (Q192*L1 + dx1*Q96*f1*SP1)
dx2 = dy1*f2*Q192*L2 / (Sp2*Sp2*L2 + dy1*Q96*f2*Sp2)
    = dx1*L1*f1*SP1*SP1*f2*Q192*L2 / [(Q192*L1 + dx1*Q96*f1*SP1) * (Sp2*Sp2*L2 + dy1*Q96*f2*Sp2)]
    = dx1*L1*f1*SP1*SP1*f2*Q192*L2 / [(Q192*L1 + dx1*Q96*f1*SP1)*Sp2*Sp2*L2 + (Q192*L1 + dx1*Q96*f1*SP1) * dy1*Q96*f2*Sp2 ]
    = dx1*L1*f1*SP1*SP1*f2*Q192*L2 / [(Q192*L1 + dx1*Q96*f1*SP1)*Sp2*Sp2*L2 + dx1*L1*f1*SP1*SP1*Q96*f2*Sp2]
    = dx1*L1*f1*SP1*SP1*f2*Q96*L2 / [(Q96*L1 + dx1*f1*SP1)*Sp2*Sp2*L2 + dx1*L1*f1*SP1*SP1*f2*Sp2]
    = dx1*L1*f1*SP1*SP1*f2*Q96*L2 / [Q96*L1*Sp2*Sp2*L2 + dx1*f1*SP1*Sp2*Sp2*L2 + dx1*L1*f1*SP1*SP1*f2*Sp2]
    = dx1*L1*f1*SP1*SP1*f2*Q96*L2 / [Q96*L1*Sp2*Sp2*L2 + dx1 * (f1*SP1*Sp2*Sp2*L2 + L1*f1*SP1*SP1*f2*Sp2)]


k1 = (L1*Sp1*f1*Sp1*f2*Q96*L2)
k2 = Sp2*Sp2*L2* Q96* L1
k3 = (Sp2*Sp2*L2*f1*SP1  + L1*Sp1*f1*Sp1*f2*Sp2)
```

### 1->0, then 0->1

```
1->0:
spNext1 = SP1 + dy1 * f1 * Q96/L1
dx1 = L * Q96 * (dy1*f1*Q96/L1)/Sp1/SpNext1
    = dy1*f1*Q192/ (Sp1*Sp1 + dy1*f1*Q96*Sp1/L1)

0->1:
spNext2 = L2 * Q96 * SP2 / (L2 * Q96 + dx1 * f2 * SP2)
dy2 = L2 * (Sp2-spNext2) / Q96
    = (L2*Sp2*dx1*f2*Sp2)/ (Q96* L2 * Q96 + dx1 * f2 * SP2 * Q96)
    = (dx1*L2*Sp2*f2*Sp2)/ (Q96* L2 * Q96 + dx1 * f2 * SP2 * Q96)
    = (dy1*f1*Q192 *L2*Sp2*f2*Sp2)/ [(Sp1*Sp1 + dy1*f1*Q96*Sp1/L1) * (Q96* L2 * Q96 + dx1 * f2 * SP2 * Q96)]
    = (dy1*f1*Q192 *L2*Sp2*f2*Sp2) / [(Sp1*Sp1 + dy1*f1*Q96*Sp1/L1) * (Q192*L2 + dx1 * f2 * SP2 * Q96)]
    = (dy1*f1*Q192 *L2*Sp2*f2*Sp2) / [(Sp1*Sp1 + dy1*f1*Q96*Sp1/L1) * Q192*L2 + (Sp1*Sp1 + dy1*f1*Q96*Sp1/L1) * dx1 * f2 * SP2 * Q96]
    = (dy1*f1*Q192 *L2*Sp2*f2*Sp2) / [(Sp1*Sp1 + dy1*f1*Q96*Sp1/L1) * Q192*L2 + dy1*f1*Q192 * f2 * SP2 * Q96]
    = (dy1*f1 *L2*Sp2*f2*Sp2) / [Sp1*Sp1* L2  + dy1*f1*Q96*Sp1/L1 *L2 + dy1*f1 * f2 * SP2 * Q96]
    = (dy1 *f1*L2*Sp2*f2*Sp2) / [Sp1*Sp1*L2  + dy1* (f1*Q96*Sp1 *L2/L1 + f1*f2*SP2*Q96)]
    = (dy1 *f1*L2*Sp2*f2*Sp2*L1) / [Sp1*Sp1*L2*L1  + dy1* (f1*Q96*Sp1*L2 + L1*f1*f2*SP2*Q96)]
k1 = f1*L2*Sp2*f2*Sp2*L1
k2 = Sp1*Sp1*L2*L1
k3 = (f1*Q96*Sp1*L2 + L1*f1*f2*SP2*Q96)

```
