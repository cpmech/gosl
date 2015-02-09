// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tsr

// Right Cauchy-Green deformation tensor: C := Ft * F
// Symmetric and positive definite: det(C) = det(F)^2
func RightCauchyGreenDef(C, F [][]float64) {
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            C[i][j] = 0.0
            for k := 0; k < 3; k++ {
                C[i][j] += F[k][i] * F[k][j]
            }
        }
    }
}

// Left Cauchy-Green deformation tensor: b := F * Ft
// Symmetric and positive definite: det(b) = det(F)^2
func LeftCauchyGreenDef(b, F [][]float64) {
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            b[i][j] = 0.0
            for k := 0; k < 3; k++ {
                b[i][j] += F[i][k] * F[j][k]
            }
        }
    }
}

// Lagrangian or Green strain tensor: E := 0.5 * (Ft * F - I)
// Symmetric
func GreenStrain(E, F [][]float64) {
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            E[i][j] = 0.0
            if i == j {
                E[i][j] = -0.5
            }
            for k := 0; k < 3; k++ {
                E[i][j] += 0.5 * F[k][i] * F[k][j]
            }
        }
    }
}

// Eulerian or Almansi strain tensor: e := 0.5 * (I - inv(F)^t * inv(F))
// Symmetric
func AlmansiStrain(e, Fi [][]float64) {
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            e[i][j] = 0.0
            if i == j {
                e[i][j] = 0.5
            }
            for k := 0; k < 3; k++ {
                e[i][j] -= 0.5 * Fi[k][i] * Fi[k][j]
            }
        }
    }
}

// Linear strain tensor: ε := 0.5 * (H + Ht) = 0.5 * (F + Ft) - I
func LinStrain(ε, F [][]float64) {
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            ε[i][j] = 0.5 * (F[i][j] + F[j][i])
            if i == j {
                ε[i][j] -= 1.0
            }
        }
    }
}

// Cauchy stress => first Piola-Kirchhoff: P := σ * inv(F)^t * J
func CauchyToPK1(P, σ, F, Fi [][]float64, J float64) {
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            P[i][j] = 0.0
            for k := 0; k < 3; k++ {
                P[i][j] += σ[i][k] * Fi[j][k] * J
            }
        }
    }
}

// First Piola-Kirchhoff => Cauchy stress: σ := P * Ft / J
func PK1ToCauchy(σ, P, F, Fi [][]float64, J float64) {
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            σ[i][j] = 0.0
            for k := 0; k < 3; k++ {
                σ[i][j] += P[i][k] * F[j][k] / J
            }
        }
    }
}

// Cauchy stress => second Piola-Kirchhoff: S := inv(F) * σ * inv(F)^t * J
func CauchyToPK2(S, σ, F, Fi [][]float64, J float64) {
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            S[i][j] = 0.0
            for k := 0; k < 3; k++ {
                for l := 0; l < 3; l++ {
                    S[i][j] += Fi[i][k] * σ[k][l] * Fi[j][l] * J
                }
            }
        }
    }
}

// Second Piola-Kirchhoff => Cauchy stress: σ := F * S * Ft / J
func PK2ToCauchy(σ, S, F, Fi [][]float64, J float64) {
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            σ[i][j] = 0.0
            for k := 0; k < 3; k++ {
                for l := 0; l < 3; l++ {
                    σ[i][j] += F[i][k] * S[k][l] * F[j][l] / J
                }
            }
        }
    }
}

// Push-forward (type A/cov): res := push-forward(a) = inv(F)^t * a * inv(F)
func PushForward(res, a, F, Fi [][]float64) {
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            res[i][j] = 0.0
            for k := 0; k < 3; k++ {
                for l := 0; l < 3; l++ {
                    res[i][j] += Fi[k][i] * a[k][l] * Fi[l][j]
                }
            }
        }
    }
}

// Pull-back (type A/cov): res := push-back(a) = Ft * a * F
func PullBack(res, a, F, Fi [][]float64) {
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            res[i][j] = 0.0
            for k := 0; k < 3; k++ {
                for l := 0; l < 3; l++ {
                    res[i][j] += F[k][i] * a[k][l] * F[l][j]
                }
            }
        }
    }
}

// Push-forward (type B/contra): res := push-forward(a) = F * a * Ft
func PushForwardB(res, a, F, Fi [][]float64) {
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            res[i][j] = 0.0
            for k := 0; k < 3; k++ {
                for l := 0; l < 3; l++ {
                    res[i][j] += F[i][k] * a[k][l] * F[j][l]
                }
            }
        }
    }
}

// Pull-back (type B/contra): res := pull-back(a) = inv(F) * a * inv(F)^t
func PullBackB(res, a, F, Fi [][]float64) {
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            res[i][j] = 0.0
            for k := 0; k < 3; k++ {
                for l := 0; l < 3; l++ {
                    res[i][j] += Fi[i][k] * a[k][l] * Fi[j][l]
                }
            }
        }
    }
}
