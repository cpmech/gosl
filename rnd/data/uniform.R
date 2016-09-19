# Copyright 2016 The Gosl Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

a <- 1.5 # min
b <- 2.5 # max

X <- seq(0.5, 3.0, 0.5)
Y <- matrix(ncol=5)
first <- TRUE
ypdf <- dunif(X, a, b)
ycdf <- punif(X, a, b)
for (i in 1:length(X)) {
    if (first) {
        Y <- rbind(c(X[i], a, b, ypdf[i], ycdf[i]))
        first <- FALSE
    } else {
        Y <- rbind(Y, c(X[i], a, b, ypdf[i], ycdf[i]))
    }
}
write.table(Y, "/tmp/gosl-rnd-uniform.dat", row.names=FALSE, col.names=c("x","a","b","ypdf","ycdf"), quote=FALSE)
print("file </tmp/gosl-rnd-uniform.dat> written")
