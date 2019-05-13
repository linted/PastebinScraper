/*
    Attempt to locate things that look like crypto wallets
*/

rule crypto_wallet_private
{
    meta:
        author= "@linted"
    strings:
        $btc_private_key = /[\b\s:]([0-9a-zA-Z]{64})[\b\s:]/
        $WIF = /[\b\s:]([59][0-9a-zA-Z]{50})[\b\s:]/
    condition:
        any of them
}

rule crypto_wallet_public
{
    meta:
        author = "@linted"
    strings:
        $btc_address = /[\b\s:]([0-9a-zA-Z]{34})[\b\s:]/
    condition:
        any of them
}