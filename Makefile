BUILDDIR=build
ROMSDIR=roms

GDRIVE_URL="https://drive.google.com/uc?export=download&id="
TETRIS_ROM_ID=1__s6AgHw6Nh6ovToG_q5DkV4AhE6yRuD
PKM_GREEN_ID=1oyxxU1ZSmWWS2KbSxeS05xCJR1hXqYGR
PKM_RED_ID=1htBeWHPB3fnjlsLJJ2DB_Zl7yjr0AU6F

directories:
	mkdir -p $(BUILDDIR)
	mkdir -p $(ROMSDIR)

download-roms: directories
	wget -c $(GDRIVE_URL)$(TETRIS_ROM_ID) -O $(ROMSDIR)/tetris.gb
	wget -c $(GDRIVE_URL)$(PKM_GREEN_ID) -O $(ROMSDIR)/pkmn_green.gb
	wget -c $(GDRIVE_URL)$(PKM_RED_ID) -O $(ROMSDIR)/pkmn_red.gb

.PHONY:test
test:
	go test -v ./test

.PHONY:clean
clean:
	rm -rf $(ROMSDIR)
	rm -rf $(build)
