DIRCMD := cmd/mydir/
DIRSTORAGE := StorageDir
TESTSTORAGE := StorageDirTest
TESTRESULT := testResult
FILEUSER := Usuarios.json
FILEUSERTEST := UsuariosTest.json
BINNAME := main
GO := go
TEST := test -coverprofile=
OUTPUTTEST := coverage.out
EXTENSIONS := build -o

all: build

dirs:
	mkdir $(DIRSTORAGE)

dirs_test:
	mkdir $(TESTRESULT)

json:
	touch $(FILEUSER)

build: dirs json
	cd $(DIRCMD) && $(GO) $(EXTENSIONS) $(BINNAME)
	mv $(DIRCMD)$(BINNAME) .

clean:
	rm -r $(DIRSTORAGE)
	rm $(FILEUSER)
	rm $(BINNAME)

test: dirs_test
	cd $(DIRCMD) && $(GO) $(TEST)$(OUTPUTTEST) && rm $(FILEUSERTEST) && rm -r $(TESTSTORAGE)
	mv $(DIRCMD)$(OUTPUTTEST) $(TESTRESULT)/

showTest:
	$(GO) tool cover -html=$(TESTRESULT)/$(OUTPUTTEST)

cleanTest:
	rm -r $(TESTRESULT)

run: 
	./$(BINNAME)