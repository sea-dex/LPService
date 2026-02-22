
# v3 formula
## v3 0->1:
```
dy = dx*f*L*SP*SP / (Q192*L + dx*Q96*f*SP)
```
## v3 1->0:
```
dx = dy*f*L*Q192 / (Sp*Sp*L + dy*Q96*f*Sp)
```

# CalcAmountInV2ToV3

## v2 0->1, v3 1->0
```
dy1 = dx1*f1*y1 / (x1 + dx1*f1)
dx2 = dy1*f2*L2*Q192 / (Sp2*Sp2*L + dy1*Q96*f2*Sp2)
dx2 = (dy1-N)*f2*L2*Q192 / (Sp2*Sp2*L2 + (dy1-N)*Q96*f2*Sp2)
dx2 = (dx1*f1*y1/(x1+dx1*f1)-N)*f2*L2*Q192 / (Sp2*Sp2*L2 + (dx1*f1*y1/(x1+dx1*f1)-N)*Q96*f2*Sp2)
dx2 = (dx1*f1*y1-N*(x1+dx1*f1))*f2*L2*Q192 / (Sp2*Sp2*L2*(x1+dx1*f1) + (dx1*f1*y1-N*(x1+dx1*f1))*Q96*f2*Sp2)
dx2 = (dx1*f1*y1-N*x1-dx1*f1*N)*f2*L2*Q192 / (Sp2*Sp2*L2*x1+dx1*f1*Sp2*Sp2*L2 + (dx1*f1*y1-N*x1-dx1*f1*N)*Q96*f2*Sp2)
dx2 = (dx1*f1*f2*L2*Q192*(y1-N) - N*x1*f2*L2*Q192) / (Sp2*Sp2*L2*x1 + dx1*f1*Sp2*Sp2*L2 + (dx1*f1*(y1-N)-N*x1)*Q96*f2*Sp2)
dx2 = (dx1*f1*f2*L2*Q192*(y1-N) - N*x1*f2*L2*Q192) / (Sp2*Sp2*L2*x1 + dx1*f1*Sp2*Sp2*L2 + dx1*f1*(y1-N)*Q96*f2*Sp2 - N*x1*Q96*f2*Sp2)
dx2 = (dx1*f1*f2*L2*Q192*(y1-N) - N*x1*f2*L2*Q192) / (Sp2*Sp2*L2*x1 - N*x1*Q96*f2*Sp2 + dx1*f1*Sp2*(Sp2*L2 + (y1-N)*Q96*f2))
k1 = f1*f2*L2*Q192*(y1-N)
k2 = Sp2*Sp2*L2*x1 - N*x1*Q96*f2*Sp2
k3 = f1*Sp2*(Sp2*L2 + (y1-N)*Q96*f2)
k4 = N*x1*f2*L2*Q192

dx2 = (dx1*k1-k4)/(k2+dx1*k3)
f = dx1*k1-k4
g = k2 + dx1*k3
profit = dx2 - dx1
dx2' - 1 = 0
f'g - fg' = k1*g - f*k3 = k1*(k2 + dx1*k3) - (dx1*k1 -k4) * k3 = k1*k2 + k3*k4 = g*g
a = k3*k3
b = 2*k2*k3
c = k2*k2 - k1*k2 - k3*k4
r = (sqrt(b**2 - 4*a*c) - b) / (2*a)
```

## v2 1->0, v3 0->1
```
dx1 = dy1*f1*x1/(y1+dy1*f1)
dy2 = dx1*f2*L2*SP2*SP2 / (Q192*L2 + dx1*Q96*f2*SP2)
dy2 = (dx1-N)*f2*L2*SP2*SP2 / (Q192*L2 + (dx1-N)*Q96*f2*Sp2)
dy2 = (dy1*f1*x1/(y1+dy1*f1)-N)*f2*L2*SP2*SP2 / (Q192*L2 + (dy1*f1*x1/(y1+dy1*f1)-N)*Q96*f2*Sp2)
dy2 = (dy1*f1*x1-N*(y1+dy1*f1))*f2*L2*SP2*SP2 / (Q192*L2*(y1+dy1*f1) + (dy1*f1*x1-N*(y1+dy1*f1))*Q96*f2*Sp2)
dy2 = (dy1*f1*(x1-N)-N*y1)*f2*L2*SP2*SP2 / (Q192*L2*y1+dy1*f1*Q192*L2 + (dy1*f1*(x1-N)-N*y1)*Q96*f2*Sp2)
dy2 = (dy1*f1*(x1-N)-N*y1)*f2*L2*SP2*SP2 / (Q192*L2*y1 + dy1*f1*Q192*L2 + dy1*f1*(x1-N)*Q96*f2*Sp2 - N*y1*Q96*f2*Sp2)
dy2 = (dy1*f1*f2*L2*SP2*SP2*(x1-N)-N*y1*f2*L2*SP2*SP2) / (Q192*L2*y1-N*y1*Q96*f2*Sp2 + dy1*f1*Q96(Q96*L2 + (x1-N)*f2*Sp2))
k1 = f1*f2*L2*SP2*SP2*(x1-N)
k2 = Q192*L2*y1-N*y1*Q96*f2*Sp2
k3 = f1*Q96(Q96*L2 + (x1-N)*f2*Sp2)
k4 = N*y1*f2*L2*SP2*SP2
a = k3*k3
b = 2*k2*k3
c = k2*k2 - k1*k2 - k3*k4
r = (-b + sqrt(b**2 - 4*a*c)) / (2*a)
```

# CalcAmountInV3ToV2

## v3 0->1, v2 1->0
```
dy1 = dx1*f1*L1*SP1*SP1 / (Q192*L1 + dx1*Q96*f1*SP1)
dx2 = dy1*f2*x2/(y2+dy1*f2)
+N

dx2 = (dy1+N)*f2*x2/(y2+(dy1+N)*f2)
    = (dx1*f1*L1*SP1*SP1 / (Q192*L1 + dx1*Q96*f1*SP1) + N)*f2*x2 / (y2+(dx1*f1*L1*SP1*SP1 / (Q192*L1 + dx1*Q96*f1*SP1) + N)*f2)
    = (dx1*f1*L1*SP1*SP1 + Q192*L1*N + dx1*Q96*f1*SP1*N)*f2*x2 / (y2 * (Q192*L1 + dx1*Q96*f1*SP1) + (dx1*f1*L1*SP1*SP1 + N*(Q192*L1 + dx1*Q96*f1*SP1))*f2)
    = (dx1*f1*L1*SP1*SP1 + Q192*L1*N + dx1*Q96*f1*SP1*N)*f2*x2 / (y2*Q192*L1 + dx1*Q96*f1*SP1*y2 + dx1*f1*f2*L1*SP1*SP1 + N*f2*(Q192*L1 + dx1*Q96*f1*SP1))
    = [dx1*(L1*SP1+Q96*N)*f1*SP1*f2*x2 + Q192*L1*N*f2*x2] / (y2*Q192*L1 + dx1*Q96*f1*SP1*y2 + dx1*f1*f2*L1*SP1*SP1 + N*f2*(Q192*L1 + dx1*Q96*f1*SP1))
    = [dx1*(L1*SP1+Q96*N)*f1*SP1*f2*x2 + Q192*L1*N*f2*x2] / (y2*Q192*L1 + dx1*Q96*f1*SP1*y2 + dx1*f1*f2*L1*SP1*SP1 + N*f2*Q192*L1 + dx1*Q96*f1*SP1*N*f2)
    = [dx1*(L1*SP1+Q96*N)*f1*SP1*f2*x2 + Q192*L1*N*f2*x2] / (y2*Q192*L1+N*f2*Q192*L1 + dx1*f1*SP1(Q96*y2 + f2*L1*SP1 + Q96*N*f2))
k1 = (L1*SP1+Q96*N)*f1*SP1*f2*x2
k2 = y2*Q192*L1+N*f2*Q192*L1
k3 = f1*SP1*(Q96*y2 + f2*L1*SP1 + Q96*N*f2)
k4 = -Q192*L1*N*f2*x2
a = k3*k3
b = 2*k2*k3
c = k2*k2 - k1*k2 - k3*k4
r = (-b + sqrt(b**2 - 4*a*c)) / (2*a)
```

## v3 1->0, v2 0->1

```
dx1 = dy1*f1*L1*Q192 / (Sp1*Sp1*L1 + dy1*Q96*f1*Sp1)
dy2 = dx1*f2*y2 / (x2 + dx1*f2)

dx1 -> dx1 + N
dy2 = (dx1+N)*f2*y2 / (x2 + (dx1+N)*f2)
    = (dy1*f1*L1*Q192 / (Sp1*Sp1*L1 + dy1*Q96*f1*Sp1)+N)*f2*y2 / (x2 + (dy1*f1*L1*Q192 / (Sp1*Sp1*L1 + dy1*Q96*f1*Sp1) + N)*f2)
    = (dy1*f1*L1*Q192 + N*(Sp1*Sp1*L1 + dy1*Q96*f1*Sp1))*f2*y2 / (x2*(Sp1*Sp1*L1 + dy1*Q96*f1*Sp1) + (dy1*f1*L1*Q192 + N*(Sp1*Sp1*L1 + dy1*Q96*f1*Sp1))*f2)
    = (dy1*f1*L1*Q192 + N*Sp1*Sp1*L1 + dy1*Q96*f1*Sp1*N)*f2*y2 / (x2*Sp1*Sp1*L1 + dy1*Q96*f1*Sp1*x2 + (dy1*f1*L1*Q192 + N*Sp1*Sp1*L1 + dy1*Q96*f1*Sp1*N)*f2)
    = (dy1*f1*Q96*L1*Q96 + dy1*Q96*f1*Sp1*N + N*Sp1*Sp1*L1)*f2*y2 / (x2*Sp1*Sp1*L1+N*f2*Sp1*Sp1*L1 + dy1*Q96*f1*Sp1*x2 + dy1*(L1*Q192+ Q96*Sp1*N)*f1*f2)
    = [dy1*f1*Q96*f2*y2*(L1*Q96+Sp1*N) + N*Sp1*Sp1*L1*f2*y2] / (x2*Sp1*Sp1*L1+N*f2*Sp1*Sp1*L1 + dy1*Q96*f1*Sp1*x2 + dy1*(L1*Q96*f2+ N*f2*Sp1)*Q96*f1)
    = [dy1*f1*Q96*f2*y2*(L1*Q96+Sp1*N) + N*Sp1*Sp1*L1*f2*y2] / (x2*Sp1*Sp1*L1+N*f2*Sp1*Sp1*L1 + dy1*Q96*f1*(Sp1*x2+L1*Q96*f2+N*f2*Sp1))
k1 = f1*Q96*f2*y2*(L1*Q96+Sp1*N)
k2 = x2*Sp1*Sp1*L1+N*f2*Sp1*Sp1*L1
k3 = Q96*f1*(Sp1*x2+L1*Q96*f2+N*f2*Sp1)
k4 = -N*Sp1*Sp1*L1*f2*y2
a = k3*k3
b = 2*k2*k3
c = k2*k2 - k1*k2 - k3*k4
r = (-b + sqrt(b**2 - 4*a*c)) / (2*a)
```

# CalcAmountInV3ToV3
## 0->1, 1->0

```
dy1 = dx1*f1*L1*SP1*SP1 / (Q192*L1 + dx1*Q96*f1*SP1)
dx2 = dy1*f2*L2*Q192 / (Sp2*Sp2*L2 + dy1*Q96*f2*Sp2)

dx1 = dx1 - n1
dy2 = dy1 - n2
dy1 = (dx1-n1)*f1*L1*SP1*SP1 / (Q192*L1 + (dx1-n1)*Q96*f1*SP1)
dx2 = (dy1-n2)*f2*L2*Q192 / (Sp2*Sp2*L2 + (dy1-n2)*Q96*f2*Sp2)
    = ((dx1-n1)*f1*L1*SP1*SP1 / (Q192*L1 + (dx1-n1)*Q96*f1*SP1) - n2)*f2*L2*Q192 /
     (Sp2*Sp2*L2 + ((dx1-n1)*f1*L1*SP1*SP1 / (Q192*L1 + (dx1-n1)*Q96*f1*SP1) - n2)*Q96*f2*Sp2)
    
    = ((dx1-n1)*f1*L1*SP1*SP1 - n2*(Q192*L1 + (dx1-n1)*Q96*f1*SP1))*f2*L2*Q192 / 
    (Sp2*Sp2*L2*(Q192*L1 + (dx1-n1)*Q96*f1*SP1) + ((dx1-n1)*f1*L1*SP1*SP1 - n2*(Q192*L1 + (dx1-n1)*Q96*f1*SP1))*Q96*f2*Sp2)

    = (dx1*f1*L1*SP1*SP1-n1*f1*L1*SP1*SP1 - n2*Q192*L1 - dx1*n2*Q96*f1*SP1+n1*n2*Q96*f1*SP1)*f2*L2*Q192 / 
        (Sp2*Sp2*L2*(Q192*L1-n1*Q96*f1*SP1 + dx1*Q96*f1*SP1) + ((dx1-n1)*f1*L1*SP1*SP1 - n2*(Q192*L1 + (dx1-n1)*Q96*f1*SP1)) *Q96*f2*Sp2 )

    = (dx1*f1*SP1*f2*L2*Q192*(L1*SP1-n2*Q96) - f2*L2*Q192*(n1*f1*L1*SP1*SP1+n2*Q192*L1-n1*n2*Q96*f1*SP1)) / 
        (Sp2*Sp2*L2*(Q192*L1-n1*Q96*f1*SP1 + dx1*Q96*f1*SP1) + ((dx1-n1)*f1*L1*SP1*SP1 - n2*(Q192*L1 + (dx1-n1)*Q96*f1*SP1)) *Q96*f2*Sp2 )

    = (dx1*f1*SP1*f2*L2*Q192*(L1*SP1-n2*Q96) - f2*L2*Q192*(n1*f1*L1*SP1*SP1+n2*Q192*L1-n1*n2*Q96*f1*SP1)) / 
        (Sp2*Sp2*L2*(Q192*L1-n1*Q96*f1*SP1) + dx1*Sp2*Sp2*L2*Q96*f1*SP1 + (dx1*f1*L1*SP1*SP1 - n1*f1*L1*SP1*SP1 - n2*Q192*L1 - (dx1-n1)*n2*Q96*f1*SP1) *Q96*f2*Sp2 )

    = (dx1*f1*SP1*f2*L2*Q192*(L1*SP1-n2*Q96) - f2*L2*Q192*(n1*f1*L1*SP1*SP1+n2*Q192*L1-n1*n2*Q96*f1*SP1)) / 
        (Sp2*Sp2*L2*(Q192*L1-n1*Q96*f1*SP1) + dx1*Sp2*Sp2*L2*Q96*f1*SP1 + (dx1*f1*SP1*(L1*SP1-n2*Q96) - n1*f1*L1*SP1*SP1 - n2*Q192*L1+n1*n2*Q96*f1*SP1) *Q96*f2*Sp2 )

    = (dx1*f1*SP1*f2*L2*Q192*(L1*SP1-n2*Q96) - f2*L2*Q192*(n1*f1*L1*SP1*SP1+n2*Q192*L1-n1*n2*Q96*f1*SP1)) / 
        (Sp2*Sp2*L2*Q96*(Q96*L1-n1*f1*SP1) -(n1*f1*L1*SP1*SP1+n2*Q192*L1-n1*n2*Q96*f1*SP1)*Q96*f2*Sp2 + dx1*f1*SP1*Sp2*Q96*(Sp2*L2+f2*L1*SP1-f2*n2*Q96) )


dx1
dy2 = dy1+cumY1-cumY2 = dy1 + cumY
dx2 = dy2*f2*L2*Q192 / (Sp2*Sp2*L2 + dy2*Q96*f2*Sp2)
    = (dy1+cumY)*f2*L2*Q192 / (Sp2*Sp2*L2 + (dy1+cumY)*Q96*f2*Sp2)
    = (dx1*f1*L1*SP1*SP1 / (Q192*L1 + dx1*Q96*f1*SP1)+cumY)*f2*L2*Q192 / (Sp2*Sp2*L2 + (dx1*f1*L1*SP1*SP1 / (Q192*L1 + dx1*Q96*f1*SP1)+cumY)*Q96*f2*Sp2)
    = (dx1*f1*L1*SP1*SP1 + cumY*(Q192*L1 + dx1*Q96*f1*SP1))*f2*L2*Q192 / 
        (Sp2*Sp2*L2*(Q192*L1 + dx1*Q96*f1*SP1) + (dx1*f1*L1*SP1*SP1 + cumY*(Q192*L1 + dx1*Q96*f1*SP1))*Q96*f2*Sp2)
    = (dx1*f1*L1*SP1*SP1 + cumY*Q192*L1 + dx1*cumY*Q96*f1*SP1)*f2*L2*Q192 /
        (Sp2*Sp2*L2*Q192*L1 + dx1*Sp2*Sp2*L2*Q96*f1*SP1 + (dx1*f1*L1*SP1*SP1 + cumY*Q192*L1 + dx1*cumY*Q96*f1*SP1)*Q96*f2*Sp2)
    = (dx1*f1*SP1*f2*L2*Q192*(L1*SP1+cumY*Q96) + cumY*Q192*L1*f2*L2*Q192) /
        ( (Sp2*L2+cumY*Q96*f2)*Q192*L1*Sp2 + dx1*f1*SP1*Q96*Sp2*(Sp2*L2 + f2*(L1*SP1+cumY*Q96)) )

k1 = f1*SP1*f2*L2*Q192*(L1*SP1+cumY*Q96)
k2 = (Sp2*L2+cumY*Q96*f2)*Q192*L1*Sp2
k3 = f1*SP1*Q96*Sp2*(Sp2*L2 + f2*(L1*SP1+cumY*Q96))
k4 = -cumY*Q192*L1*f2*L2*Q192
```

## 1->0, 0->1
```
dx1 = dy1*f1*L1*Q192 / (Sp1*Sp1*L1 + dy1*Q96*f1*Sp1)
dy2 = dx2*f2*L2*SP2*SP2 / (Q192*L2 + dx2*Q96*f2*SP2)


dx2 = dx1+cumX1-cumX2 = dx1+cumX
dy2 = dx2*f2*L2*SP2*SP2 / (Q192*L2 + dx2*Q96*f2*SP2)
    = (dx1+cumX)*f2*L2*SP2*SP2 / (Q192*L2 + (dx1+cumX)*Q96*f2*SP2)
    = (dy1*f1*L1*Q192 / (Sp1*Sp1*L1 + dy1*Q96*f1*Sp1)+cumX)*f2*L2*SP2*SP2 / (Q192*L2 + (dy1*f1*L1*Q192 / (Sp1*Sp1*L1 + dy1*Q96*f1*Sp1)+cumX)*Q96*f2*SP2)
    = (dy1*f1*L1*Q192 +cumX*(Sp1*Sp1*L1 + dy1*Q96*f1*Sp1))*f2*L2*SP2*SP2 /
        (Q192*L2*(Sp1*Sp1*L1 + dy1*Q96*f1*Sp1) + (dy1*f1*L1*Q192 + cumX*(Sp1*Sp1*L1 + dy1*Q96*f1*Sp1))*Q96*f2*SP2)
    = (dy1*f1*L1*Q192 +cumX*Sp1*Sp1*L1 + dy1*cumX*Q96*f1*Sp1)*f2*L2*SP2*SP2 /
        (Q192*L2*Sp1*Sp1*L1 + dy1*Q192*L2*Q96*f1*Sp1 + (dy1*f1*L1*Q192 + cumX*Sp1*Sp1*L1 + dy1*cumX*Q96*f1*Sp1)*Q96*f2*SP2)
    = (dy1*f1*Q96*(L1*Q96+cumX*Sp1) +cumX*Sp1*Sp1*L1)*f2*L2*SP2*SP2 /
        (Q192*L2*Sp1*Sp1*L1 + dy1*Q192*L2*Q96*f1*Sp1 + (dy1*f1*Q96*(L1*Q96+cumX*Sp1) + cumX*Sp1*Sp1*L1)*Q96*f2*SP2)
    = (dy1*f1*Q96*f2*L2*SP2*SP2*(L1*Q96+cumX*Sp1) + cumX*Sp1*Sp1*L1*f2*L2*SP2*SP2) /
        (Q192*L2*Sp1*Sp1*L1 + dy1*Q192*L2*Q96*f1*Sp1 + dy1*f1*Q96*Q96*f2*SP2*(L1*Q96+cumX*Sp1) + cumX*Sp1*Sp1*L1*Q96*f2*SP2)
    = (dy1*f1*Q96*f2*L2*SP2*SP2*(L1*Q96+cumX*Sp1) + cumX*Sp1*Sp1*L1*f2*L2*SP2*SP2) /
        (Sp1*Sp1*L1*Q96*(Q96*L2+cumX*f2*SP2) + dy1*f1*Q192*(L2*Q96*Sp1 + f2*SP2*(L1*Q96+cumX*Sp1)))


------

dy1 = dy1 - n1
dx2 = dx1 - n2

dx1 = (dy1-n1)*f1*L1*Q192 / (Sp1*Sp1*L1 + (dy1-n1)*Q96*f1*Sp1)
dy2 = (dx2-n2)*f2*L2*SP2*SP2 / (Q192*L2 + (dx2-n2)*Q96*f2*SP2)
    = (dx1-n2)*f2*L2*SP2*SP2 / (Q192*L2 + (dx1-n2)*Q96*f2*SP2)

    = ((dy1-n1)*f1*L1*Q192 / (Sp1*Sp1*L1 + (dy1-n1)*Q96*f1*Sp1)-n2)*f2*L2*SP2*SP2 / 
        (Q192*L2 + ((dy1-n1)*f1*L1*Q192 / (Sp1*Sp1*L1 + (dy1-n1)*Q96*f1*Sp1)-n2)*Q96*f2*SP2)

    = ((dy1-n1)*f1*L1*Q192 - n2*(Sp1*Sp1*L1 + (dy1-n1)*Q96*f1*Sp1))*f2*L2*SP2*SP2 / 
        (Q192*L2*(Sp1*Sp1*L1 + (dy1-n1)*Q96*f1*Sp1) + ((dy1-n1)*f1*L1*Q192 - n2*(Sp1*Sp1*L1 + (dy1-n1)*Q96*f1*Sp1))*Q96*f2*SP2)

    = (dy1*f1*L1*Q192-n1*f1*L1*Q192 - (n2*Sp1*Sp1*L1 + (dy1-n1)*n2*Q96*f1*Sp1))*f2*L2*SP2*SP2 / 
        (Q192*L2*Sp1*Sp1*L1 - n1*Q192*L2*Q96*f1*Sp1 + dy1*Q192*L2*Q96*f1*Sp1 + ((dy1-n1)*f1*Q96*(L1*Q96-n2*SP1)- n2*Sp1*Sp1*L1)*Q96*f2*SP2)
    
    = (dy1*f1*L1*Q192 - dy1*n2*Q96*f1*Sp1 - n1*f1*L1*Q192 - n2*Sp1*Sp1*L1 + n1*n2*Q96*f1*Sp1)*f2*L2*SP2*SP2 /
        (Q192*L2*Sp1*Sp1*L1 - Q192*L2*Sp1*n1*Q96*f1 + dy1*Q192*L2*Q96*f1*Sp1 + (dy1-n1)*f1*Q96*(L1*Q96-n2*SP1)*Q96*f2*SP2 - n2*Sp1*Sp1*L1*Q96*f2*SP2)

    = (dy1*f1*Q96*f2*L2*SP2*SP2*(L1*Q96-n2*Sp1) - f2*L2*SP2*SP2*(n1*f1*L1*Q192+n2*Sp1*Sp1*L1-n1*n2*Q96*f1*Sp1)) / 
        (Q192*L2*Sp1*Sp1*L1 - Q192*L2*Sp1*n1*Q96*f1 - n2*Sp1*Sp1*L1*Q96*f2*SP2 - n1*f1*Q96*(L1*Q96-n2*SP1)*Q96*f2*SP2 + dy1*f1*Q192*(L2*Q96*Sp1 + (L1*Q96-n2*SP1)*f2*SP2) )

    = (dy1*f1*Q96*f2*L2*SP2*SP2*(L1*Q96-n2*Sp1) - f2*L2*SP2*SP2*(n1*f1*L1*Q192+n2*Sp1*Sp1*L1-n1*n2*Q96*f1*Sp1)) / 
        (Q192*L2*Sp1*Sp1*L1 + (dy1-n1)*Q192*L2*Q96*f1*Sp1 + ((dy1-n1)*f1*Q96*(L1*Q96-n2*SP1)- n2*Sp1*Sp1*L1)*Q96*f2*SP2)
        (Q192*L2*Sp1*(Sp1*L1 - n1*Q96*f1) - (n2*Sp1*Sp1*L1 - n1*f1*Q96*(L1*Q96-n2*SP1))*Q96*f2*SP2 + dy1*f1*Q192*(L2*Q96*Sp1 + (L1*Q96-n2*SP1)*f2*SP2) )

    = (dy1*f1*Q96*f2*L2*SP2*SP2*(L1*Q96-n2*Sp1) - f2*L2*SP2*SP2*(n1*f1*L1*Q192+n2*Sp1*Sp1*L1-n1*n2*Q96*f1*Sp1)) / 
        (Q96*Sp1*Sp1*L1*(Q96*L2-n2*f2*SP2) - n1*Q192*f1*(Q96*L2*Sp1+f2*SP2*(L1*Q96-n2*SP1)) + dy1*Q192*f1*(Q96*L2*Sp1 + f2*SP2*(L1*Q96-n2*SP1)))
```
