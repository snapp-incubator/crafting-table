ROOT=$1

if [ ! -d "$ROOT/.build" ]; then
  mkdir "$ROOT/.build"
fi


./crafting-table generate -s $ROOT/example/src/example.go -d $ROOT/.build/example.go --get "[ var1, (var1, var3) ]" --update "[[(var3),(var2, var1)], [(var2, var3), (var1)]]" --create true

expectedFile="$ROOT/example/dst/builtin.go"
buildFile="$ROOT"/.build/example.go

expected=$(cat $expectedFile)
build=$(cat $buildFile)

if [ "$expected" != "$build" ]; then
  cmp "$expectedFile" "$buildFile"
  echo "ERROR: integration test failed. Expected file and created file from example.go are different"
  rm -r "$ROOT/.build"
  exit 1
fi

rm -r "$ROOT/.build"
echo "PASS: integration test passed"