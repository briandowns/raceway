# vim:ft=make:

GOCMD = go
GOBUILD = $(GOCMD) build
GOGET = $(GOCMD) get -v 
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
	$(GOGET) github.com/codegangsta/negroni
	$(GOGET) github.com/goincremental/negroni-sessions
	$(GOGET) github.com/goincremental/negroni-sessions/cookiestore
	$(GOGET) github.com/gorilla/mux
	$(GOGET) github.com/pborman/uuid
	$(GOGET) github.com/thoas/stats
	$(GOGET) github.com/unrolled/render
	$(GOGET) github.com/boltdb/bolt
	$(GOGET) github.com/pborman/uuid
	$(GOGET) github.com/robfig/cron
	$(GOGET) github.com/go-sql-driver/mysql
	$(GOGET) github.com/jinzhu/gorm

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
