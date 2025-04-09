with import <nixpkgs> {};

stdenv.mkDerivation {
    name = "go";
    buildInputs = [
        go
    ];
    shellHook = ''
        #export "Go version: $(go version)"
        #export GOPATH=$HOME/go
        #export PATH=$PATH:$GOPATH/bin
    '';
}
