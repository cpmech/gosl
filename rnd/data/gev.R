# Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# needs apt-get install r-cran-fextremes

suppressPackageStartupMessages(require(fExtremes))
X <- seq(-4, 4, 0.5)
K <- c(-0.5, 0, 0.5)
M <- c(-0.5, 0, 0.5)
S <- c(0.5, 1)
Y <- matrix(ncol=5)
first <- TRUE
for (ksi in K) {
    for (mu in M) {
        for (sig in S) {
            ypdf <- dgev(X, ksi, mu, sig)
            ycdf <- pgev(X, ksi, mu, sig)
            for (i in 1:length(X)) {
                if (first) {
                    Y <- rbind(c(X[i], ksi, mu, sig, ypdf[i], ycdf[i]))
                    first <- FALSE
                } else {
                    Y <- rbind(Y, c(X[i], ksi, mu, sig, ypdf[i], ycdf[i]))
                }
            }
        }
    }
}
write.table(Y, "/tmp/gosl-rnd-gev.dat", row.names=FALSE, col.names=c("x","ksi","mu","sig","ypdf","ycdf"), quote=FALSE)
print("file </tmp/gosl-rnd-gev.dat> written")
