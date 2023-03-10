export MSYS_NO_PATHCONV=1
starttime=$(date +%s)

# launch network; create channel and join peer to channel
pushd bc-network
./network.sh down
./network.sh up createChannel -c mychannel -ca
./network.sh deployCC -ccn basic -ccp ../chaincode-go/ -ccl go
popd

# run gateway 
pushd application-gateway-go
go run evoting.go
popd

cat <<EOF
Total setup execution time : $(($(date +%s) - starttime)) secs ...
EOF
