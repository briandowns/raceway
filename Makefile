# vim:ft=make:

GOCMD = go
GOBUILD = $(GOCMD) build
GOGET = $(GOCMD) get
GOCLEAN = $(GOCMD) clean
GOINSTALL = $(GOCMD) install
GOTEST = $(GOCMD) test

define GIT_ERROR

FATAL: Git (git) is required to download scanii-go dependencies.
endef

define HG_ERROR

FATAL: Mercurial (hg) is required to download scanii-go dependencies.
endef

all: install

test:
	$(GOTEST) -v -cover github.com/briandowns/raceway

dep:
	$(GOGET) -v .

install: dep
	$(GOINSTALL) -v

clean:
	$(GOCLEAN) -n -i -x

build: dep 
	$(GOBUILD) -v 

# check for git
git:
	$(if $(shell git), , $(error $(GIT_ERROR)))

# check for mercurial
hg:
	$(if $(shell hg), , $(error $(HG_ERROR)))
