include $(GOROOT)/src/Make.inc

TARG=message
GOFILES=\
	message.go\
	utils.go\
	multipart.go\
	encodings.go

include $(GOROOT)/src/Make.pkg
