# Makefile
# Build alpt4ats
# By J. Stuart McMurray
# Created 20240901
# Last Modified 20240901

BINDIR      = bin
BINS        = ${DIRS:S,^,${BINDIR}/,}
GOLDFLAGS   = -w -s
.ifdef LINKFLAGS
GOLDFLAGS  += ${LINKFLAGS}
.endif
BUILDFLAGS  = -trimpath -ldflags '${GOLDFLAGS:S/'/'\''/g}'
SRCDIR      = src
DIRS       != cd ${SRCDIR} && ls */*.go* | cut -f 1 -d / | sort -u
TESTFLAGS  += -timeout 3s
VETFLAGS    = -printf.funcs 'debugf,errorf,erorrlogf,logf,printf'
FROMDIFFS   =

.PHONY: all test clean

# all builds ALL the things
all: ${DIRS} fromdiffs fake_edr


# Rules to build binaries for each directory.
.for D in ${DIRS}
# Source directory
SD=${SRCDIR}/$D

# Allow for building individual binaries without specifying bin/
.PHONY: $D
$D: ${BINDIR}/$D

# Update the sources which come from diffs.  There's gotta be a better way
# to get paths nicely.
DIFFS != find ${SD} -type f -name '*.diff'
.for DIFF in ${DIFFS}
DSRC = Getting source from ${DIFF} ...      # For debugging
DSRC != awk '/^--- /   {print $$2}' ${DIFF} # Source file for the diff
DDST != awk '/^\+\+\+/ {print $$2}' ${DIFF} # Where the diff goes
.poison empty DSRC
.poison empty DDST
${DIFF}_SRC   := ${DSRC} 
${DIFF}_DST   := ${DDST}
${DIFF}_FSRC !:= cd ${SD}/${DSRC:H} && pwd
${DIFF}_FSRC  := ${${DIFF}_FSRC:S,${.CURDIR}/,,}/${DSRC:T}
${DIFF}_FDST !:= cd ${SD}/${DDST:H} && pwd
${DIFF}_FDST  := ${${DIFF}_FDST:S,${.CURDIR}/,,}/${DDST:T}

.ifmake diffs || ${DIFF} # We .ifmake guard this to avoid circular dependencies.
${DIFF}: ${${DIFF}_FSRC} ${${DIFF}_FDST}
	cd ${@:H} && diff -u ${$@_SRC} ${$@_DST} > ${@:T} || [ $$? -eq 1 ]
diffs: ${DIFF}

.else # .ifmake diffs
${DIFF:R}: ${DIFF} ${${DIFF}_FSRC}
	cp ${>:N*.diff} $@
	patch --quiet --directory ${@:H} --input ${>:M*.diff:T} &&\
		rm $@.orig
FROMDIFFS += ${DIFF:R}
.endif

.endfor # .for DIFF IN ${DIFFS}

# fromdiffs builds all of the things that come from diffs.
fromdiffs: ${FROMDIFFS}

# Build the binary from the directory
FILES  != ls ${SD}/*
$D_DIR := ${SD}
${BINDIR}/$D: ${FILES:M*.go} ${FILES:M*.diff:R:M*.go}
	go build ${BUILDFLAGS} -o ./$@ ./${$D_DIR}

.endfor # .for DIR in ${DIRS}

# fake_edr is special, as it's not written in Go.
.PHONY: fake_edr

fake_edr: ${BINDIR}/fake_edr

${BINDIR}/fake_edr: ${SRCDIR}/fake_edr/*.c
	cc -O2 -Wall --pedantic -Wextra -o $@ $>

# test runs a series of tests on the code
test:
	go test ${BUILDFLAGS} ${TESTFLAGS} ./...
	go vet  ${BUILDFLAGS} ${VETFLAGS} ./...
	staticcheck ./...
	deadcode ./...

# clean removes built binaries
clean: cleanbins
	rm -f ${FROMDIFFS:S/^ //}

cleanbins:
	rm -rf ${BINDIR}
