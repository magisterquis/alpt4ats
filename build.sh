#!/bin/sh
#
# build.sh
# Build ALL the libraries
# By J. Stuart McMurray
# Created 20240910
# Last Modified 20240910

set -e

BINDIR=bin
SRCDIR=src

# Build ALL the Go things
DIRS=$(cd $SRCDIR && ls */*.go* | cut -f 1 -d / | sort -u)
mkdir -p $BINDIR
GOLDFLAGS="-w -s"
if [ -n "$LINKFLAGS" ]; then
        GOLDFLAGS="$GOLDFLAGS $LINKFLAGS"
fi
for D in $DIRS; do
        BIN="$BINDIR/$D"
        (
                set -x
                go build \
                        -trimpath \
                        -ldflags "${GOLDFLAGS}" \
                        -o "$BINDIR/$D" \
                        "./$SRCDIR/$D"
        )
        ls -l "$BIN"
done

# ...and fake_edr.
BIN="$BINDIR/fake_edr" 
(
        set -x
        cc -O2 -Wall --pedantic -Wextra -o "$BIN" $SRCDIR/fake_edr/*.c
)
ls -l $BIN
