all: give_ownership take_ownership

give_ownership: give_ownership.c
	$(CC) -o $@ $^

take_ownership: take_ownership.c
	$(CC) -o $@ $^
	sudo setcap cap_chown=ep $@

give_ownership take_ownership: ownership_attrs.h
