include $(GOROOT)/src/Make.inc

TARG=mime/message
GOFILES=\
	message.go\
	utils.go\
	multipart.go\
	encodings.go

include $(GOROOT)/src/Make.pkg
