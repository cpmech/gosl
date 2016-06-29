# Copyright 2016 The Gosl Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

X <- seq(0, 3, 0.25)
N <- c(0, 0.5, 1)    # meanlog
Z <- c(0.25, 0.5, 1) # sdlog
Y <- matrix(ncol=4)
first <- TRUE
for (n in N) {
    for (z in Z) {
        ypdf <- dlnorm(X, n, z)
        ycdf <- plnorm(X, n, z)
        for (i in 1:length(X)) {
            if (first) {
                Y <- rbind(c(X[i], n, z, ypdf[i], ycdf[i]))
                first <- FALSE
            } else {
                Y <- rbind(Y, c(X[i], n, z, ypdf[i], ycdf[i]))
            }
        }
    }
}
write.table(Y, "/tmp/gosl-rnd-lognormal.dat", row.names=FALSE, col.names=c("x","n","z","ypdf","ycdf"), quote=FALSE)
print("file </tmp/gosl-rnd-lognormal.dat> written")
