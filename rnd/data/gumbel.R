# Copyright 2016 The Gosl Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# needs r-cran-evd

library(evd)

X <- seq(-3, 3, 0.5)
U <- c(0, 0.5, 1) # location
B <- c(0.5, 1, 2) # scale
Y <- matrix(ncol=4)
first <- TRUE
for (u in U) {
    for (b in B) {
        ypdf <- dgumbel(X, u, b)
        ycdf <- pgumbel(X, u, b)
        for (i in 1:length(X)) {
            if (first) {
                Y <- rbind(c(X[i], u, b, ypdf[i], ycdf[i]))
                first <- FALSE
            } else {
                Y <- rbind(Y, c(X[i], u, b, ypdf[i], ycdf[i]))
            }
        }
    }
}
write.table(Y, "/tmp/gosl-rnd-gumbel.dat", row.names=FALSE, col.names=c("x","u","b","ypdf","ycdf"), quote=FALSE)
print("file </tmp/gosl-rnd-gumbel.dat> written")
