PKG := core

all: clean
	c-for-go -debug $(PKG).yml
	sed -i. 's/diveCallback[A-Z0-9^\)]*/diveCallback/g' core/cgo_helpers.c

clean:
	rm -f $(PKG)/cgo_helpers.go $(PKG)/cgo_helpers.h $(PKG)/cgo_helpers.c
	rm -f $(PKG)/const.go $(PKG)/doc.go $(PKG)/types.go
	rm -f $(PKG)/$(PKG).go

test:
	cd $(PKG) && go build -v
