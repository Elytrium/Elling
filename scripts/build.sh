DIR="build"
if [ -d "$DIR" ]; then
  rm build/elling
  rm build/plugins/*.so
else
  mkdir build
fi

echo Building..
go build -buildmode=exe -o $DIR/elling main.go
echo Building modules..
go build -buildmode=plugin -o $DIR/plugins/oauth.so basic/oauth/oauth.go
go build -buildmode=plugin -o $DIR/plugins/topup.so basic/topup/topup.go