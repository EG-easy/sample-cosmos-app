# Initialize configuration files and genesis file
nsd init --chain-id testchain

# Copy the `Address` output here and save it for later use
nscli keys add jack

# Copy the `Address` output here and save it for later use
nscli keys add alice

# Add both accounts, with coins to the genesis file
nsd add-genesis-account $(nscli keys show jack -a) 1000mycoin,1000jackcoin
nsd add-genesis-account $(nscli keys show alice -a) 1000mycoin,1000alicecoin

# Configure your CLI to eliminate need for chain-id flag
nscli config chain-id testchain
nscli config output json
nscli config indent true
nscli config trust-node true
