#!/bin/sh

# 在 https://github.com/iqlusioninc/relayer/tree/master/testnets/relayer-alpha-2 目录中执行
rm -rf paths
mkdir -p paths
echo '{"src":{"chain-id":"irishub","client-id":"ibconeclient","connection-id":"ibconeconnection","channel-id":"ibconexfer","port-id":"transfer","order":"ordered"},"dst":{"chain-id":"ibc1","client-id":"ibczeroclient","connection-id":"ibczeroconnection","channel-id":"ibczeroxfer","port-id":"transfer","order":"ordered"},"strategy":{"type":"naive"}}' > paths/temp.json
for f in `ls *.json`; do
	chainid=$(jq --raw-output '."chain-id"' $f)
	subid=${chainid//-/''}
	key=$(jq --raw-output '."key"' $f)

	if [ $chainid ]; then
		echo "Connecting" $chainid

		if [ ${#subid} -gt 12 ]; then
			subid=${subid:0:12}
		fi

		rly chains add -f $f
		rly keys add $chainid $key
		rly chains edit $chainid key $key
		timeout 10 rly testnets request $chainid $key
		if [ $? -ne 0 ]; then
			continue
		fi
		rly lite init irishub -f
		rly lite init $chainid -f

		icliid=$(cat /dev/urandom | tr -dc 'a-z' | fold -w 10 | head -n 1)
		iconnid=$(cat /dev/urandom | tr -dc 'a-z' | fold -w 10 | head -n 1)
		ichanid=$(cat /dev/urandom | tr -dc 'a-z' | fold -w 10 | head -n 1)

		cliid=$(cat /dev/urandom | tr -dc 'a-z' | fold -w 10 | head -n 1)
		connid=$(cat /dev/urandom | tr -dc 'a-z' | fold -w 10 | head -n 1)
		chanid=$(cat /dev/urandom | tr -dc 'a-z' | fold -w 10 | head -n 1)

		jq '.src."client-id"="iris'$icliid'" | .src."connection-id"="iris'$iconnid'" | .src."channel-id"="iris'$ichanid'" | .dst."chain-id"="'$chainid'" | .dst."client-id"="iris'$cliid'" | .dst."connection-id"="iris'$connid'" | .dst."channel-id"="iris'$chanid'"' paths/temp.json > paths/$chainid.json
		rly paths add irishub $chainid i$subid -f paths/$chainid.json
		rly tx link i$subid
		if [ $? -ne 0 ]; then
			continue
		fi
		rly tx transfer irishub $chainid 100iris true $(rly ch addr $chainid)
		if [ $? -ne 0 ]; then
			continue
		fi

	fi
done;