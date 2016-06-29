# Copyright 2016 The Gosl Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# needs r-cran-evd

library(evd)

X <- seq(0, 4, 0.5)
L <- c(0, 0.5)  # location
C <- c(1, 2)    # scale
A <- c(1, 2, 3) # shape
Y <- matrix(ncol=5)
first <- TRUE
for (l in L) {
    for (c in C) {
        for (a in A) {
            ypdf <- dfrechet(X, l, c, a)
            ycdf <- pfrechet(X, l, c, a)
            for (i in 1:length(X)) {
                if (first) {
                    Y <- rbind(c(X[i], l, c, a, ypdf[i], ycdf[i]))
                    first <- FALSE
                } else {
                    Y <- rbind(Y, c(X[i], l, c, a, ypdf[i], ycdf[i]))
                }
            }
        }
    }
}
write.table(Y, "/tmp/gosl-rnd-frechet.dat", row.names=FALSE, col.names=c("x","l","c","a","ypdf","ycdf"), quote=FALSE)
print("file </tmp/gosl-rnd-frechet.dat> written")
